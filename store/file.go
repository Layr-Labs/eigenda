package store

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path"

	plasma "github.com/Layr-Labs/op-plasma-eigenda"
)

type FileStore struct {
	directory string
}

func NewFileStore(directory string) *FileStore {
	return &FileStore{
		directory: directory,
	}
}

func (s *FileStore) Get(ctx context.Context, key []byte) ([]byte, error) {
	data, err := os.ReadFile(s.fileName(key))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, plasma.ErrNotFound
		}
		return nil, err
	}
	return data, nil
}

func (s *FileStore) Put(ctx context.Context, key []byte, value []byte) error {
	return os.WriteFile(s.fileName(key), value, 0600)
}

func (s *FileStore) fileName(key []byte) string {
	return path.Join(s.directory, hex.EncodeToString(key))
}

func (s *FileStore) PutWithComm(ctx context.Context, key []byte, value []byte) error {
	return s.Put(ctx, key, value)
}

func (s *FileStore) PutWithoutComm(ctx context.Context, value []byte) (key []byte, err error) {
	// make key fingerprint of value
	// this could result in collisions
	hasher := sha256.New()
	hasher.Write(value)
	bs := hasher.Sum(nil)

	if err := s.PutWithComm(ctx, bs, value); err != nil {
		return nil, err
	}

	return bs, err
}
