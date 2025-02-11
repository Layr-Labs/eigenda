package littdb

var _ LittDB = &littDB{}

// littDB is an implementation of LittDB.
type littDB struct {
}

// NewLittDB creates a new LittDB instance.
func NewLittDB(config *LittDBConfig) LittDB {
	return &littDB{}
}

func (l *littDB) Put(key *LittKey, value []byte) error {
	//TODO implement me
	panic("implement me")
}

func (l *littDB) Get(key *LittKey) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (l *littDB) Flush() error {
	//TODO implement me
	panic("implement me")
}

func (l *littDB) Start() error {
	//TODO implement me
	panic("implement me")
}

func (l *littDB) Stop() error {
	//TODO implement me
	panic("implement me")
}
