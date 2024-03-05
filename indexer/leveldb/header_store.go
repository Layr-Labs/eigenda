package leveldb

import (
	"errors"
	"os"

	"github.com/Layr-Labs/eigenda/indexer"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	ErrNotFound              = errors.New("not found")
	ErrPrevBlockHashNotFound = errors.New("previous block hash not found")
)

type headerEntryReader struct {
	db interface {
		Get(key []byte, value any) error
		Iter(key []byte) *iter
	}
}

func (r headerEntryReader) GetHeaderEntry(key []byte) (*headerEntry, error) {
	e := new(headerEntry)

	err := r.db.Get(key, e)
	if err == nil {
		return e, nil
	}
	if errors.Is(err, leveldb.ErrNotFound) {
		return nil, ErrNotFound
	}

	return nil, err
}

func (r headerEntryReader) GetLatestHeaderEntry() (*headerEntry, error) {
	it := r.db.Iter(headerKeyPrefix)
	defer it.Release()

	if !it.First() {
		return nil, ErrNotFound
	}

	entry := new(headerEntry)
	if err := it.Value(entry); err != nil {
		return nil, err
	}
	return entry, nil
}

type headerEntryWriter struct {
	tx     *transaction
	reader headerEntryReader
}

func (w headerEntryWriter) PutHeaderEntries(headers indexer.Headers) (indexer.Headers, error) {
	var err error

	w.putFinalizedHeaderEntry(headers)

	if !w.tx.Empty() {
		headers, err = w.filterNew(headers)
		if err != nil {
			return nil, err
		}
		if headers.Empty() {
			return nil, nil
		}
	}

	for _, header := range headers {
		if err := w.putHeaderEntry(header); err != nil {
			return nil, err
		}
	}

	return headers, nil
}

func (w headerEntryWriter) filterNew(headers indexer.Headers) (indexer.Headers, error) {
	entry, err := w.reader.GetHeaderEntry(newHeaderKey(headers.Last().Number))
	if err == nil && entry.Header.Equals(headers.Last()) {
		return nil, nil
	}
	if !errors.Is(err, ErrNotFound) {
		return nil, err
	}

	for i := headers.Len() - 1; i >= 0; i-- {
		entry, err = w.reader.GetHeaderEntry(newHeaderKey(headers[i].Number - 1))
		if errors.Is(err, ErrNotFound) {
			continue
		}
		if err != nil {
			return nil, err
		}

		if !headers[i].After(entry.Header) {
			continue
		}

		return headers[i:], nil
	}

	return nil, ErrPrevBlockHashNotFound
}

func (w headerEntryWriter) putFinalizedHeaderEntry(headers indexer.Headers) {
	var finalized *indexer.Header

	for i := headers.Len() - 1; i >= 0; i-- {
		if headers[i].Finalized {
			finalized = headers[i]
			break
		}
	}

	if finalized != nil {
		w.tx.Put(finalizedHeaderKey, newHeaderEntry(finalized))
	}
}

func (w headerEntryWriter) putHeaderEntry(header *indexer.Header) error {
	var oldEntry headerEntry

	err := w.tx.Get(newHeaderKey(header.Number), &oldEntry)
	if err != nil && !errors.Is(err, leveldb.ErrNotFound) {
		return err
	}

	for _, key := range oldEntry.AccumulatorKeys {
		w.tx.Delete(key)
	}

	w.tx.Put(newHeaderKey(header.Number), newHeaderEntry(header))
	return nil
}

type HeaderStore struct {
	db     *levelDB
	opener []opener
	reader headerEntryReader
}

var _ indexer.HeaderStore = (*HeaderStore)(nil)

func NewHeaderStore(path string, opener ...opener) (*HeaderStore, error) {
	db, err := newLevelDB(path, opener...)
	if err != nil {
		return nil, err
	}

	r := headerEntryReader{db: db}
	return &HeaderStore{
		db:     db,
		opener: opener,
		reader: r,
	}, nil
}

func (s *HeaderStore) Close() {
	s.db.Close()
}

func (s *HeaderStore) AddHeaders(headers indexer.Headers) (indexer.Headers, error) {
	if headers.Empty() {
		return headers, nil
	}
	if err := headers.OK(); err != nil {
		return nil, err
	}

	tx, err := s.db.Tx()
	if err != nil {
		return nil, err
	}
	defer tx.Discard()

	r := headerEntryReader{db: tx}
	w := headerEntryWriter{tx: tx, reader: r}

	headers, err = w.PutHeaderEntries(headers)
	if err != nil {
		return nil, err
	}
	if err := w.tx.Commit(); err != nil {
		return nil, err
	}

	return headers, nil
}

func (s *HeaderStore) GetLatestHeader(finalized bool) (*indexer.Header, error) {
	var (
		e   *headerEntry
		err error
	)

	if finalized {
		e, err = s.reader.GetHeaderEntry(finalizedHeaderKey)
	} else {
		e, err = s.reader.GetLatestHeaderEntry()
	}

	if errors.Is(err, ErrNotFound) {
		return nil, indexer.ErrNoHeaders
	}
	if err != nil {
		return nil, err
	}
	return e.Header, nil
}

func (s *HeaderStore) AttachObject(
	object indexer.AccumulatorObject,
	header *indexer.Header,
	acc indexer.Accumulator,
) error {
	accData, err := acc.SerializeObject(object, indexer.UpgradeFork(header.CurrentFork))
	if err != nil {
		return err
	}

	tx, err := s.db.Tx()
	if err != nil {
		return err
	}
	defer tx.Discard()

	accKey := newAccumulatorKey(acc, header)
	tx.Put(accKey, newAccumulatorEntry(header.Number, accData))

	hdrKey := newHeaderKey(header.Number)
	tx.Put(hdrKey, s.updateHeaderEntry(header, accKey, tx))

	return tx.Commit()
}

func (s *HeaderStore) GetLatestObject(
	acc indexer.Accumulator,
	finalized bool,
) (indexer.AccumulatorObject, *indexer.Header, error) {
	header, err := s.GetLatestHeader(finalized)
	if err != nil {
		return nil, nil, err
	}
	return s.GetObject(header, acc)
}

func (s *HeaderStore) GetObject(
	header *indexer.Header,
	acc indexer.Accumulator,
) (indexer.AccumulatorObject, *indexer.Header, error) {
	accEntry, err := s.getAccumulatorEntry(header, acc)
	if err != nil {
		return nil, nil, err
	}

	hdrEntry, err := s.reader.GetHeaderEntry(newHeaderKey(accEntry.HeaderNumber))
	if err != nil {
		return nil, nil, err
	}

	accObj, err := acc.DeserializeObject(accEntry.AccumulatorData, indexer.UpgradeFork(hdrEntry.Header.CurrentFork))
	if err != nil {
		return nil, nil, err
	}

	return accObj, hdrEntry.Header, nil
}

func (s *HeaderStore) FastForward() error {
	// TODO: make FastForward() return an error to avoid panics here
	finalized, err := s.GetLatestHeader(true)

	if err != nil {
		return err
	}

	path := s.db.Path
	s.Close()

	if err := os.RemoveAll(path); err != nil {
		return err
	}

	db, err := newLevelDB(path, s.opener...)
	if err != nil {
		return err
	}

	s.db = db
	s.reader = headerEntryReader{db: db}

	var headers indexer.Headers
	headers = append(headers, finalized)

	_, err = s.AddHeaders(headers)
	if err != nil {
		return err
	}
	return nil
}

func (s *HeaderStore) updateHeaderEntry(header *indexer.Header, accKey []byte, tx *transaction) *headerEntry {
	e, err := s.reader.GetHeaderEntry(newHeaderKey(header.Number))
	if err != nil {
		tx.SetErr(err)
		return nil
	}
	return e.UpdateAccumulatorKeys(accKey)
}

func (s *HeaderStore) getAccumulatorEntry(
	header *indexer.Header, acc indexer.Accumulator,
) (*accumulatorEntry, error) {
	var (
		entry = new(accumulatorEntry)
		it    = s.db.Iter(newAccumulatorKeyPrefix(acc))
	)
	defer it.Release()

	for ok := it.First(); ok; ok = it.Next() {
		if err := it.Value(entry); err != nil {
			return nil, err
		}

		if entry.HeaderNumber > header.Number {
			continue
		}

		return entry, nil
	}

	return nil, ErrNotFound
}
