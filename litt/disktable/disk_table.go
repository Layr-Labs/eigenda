package disktable

import (
	"github.com/Layr-Labs/eigenda/litt"
	"time"
)

var _ litt.ManagedTable = &diskTable{}

type diskTable struct {
}

func (d *diskTable) Name() string {
	//TODO implement me
	panic("implement me")
}

func (d *diskTable) Put(key []byte, value []byte) error {
	//TODO implement me
	panic("implement me")
}

func (d *diskTable) Get(key []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (d *diskTable) Flush() error {
	//TODO implement me
	panic("implement me")
}

func (d *diskTable) SetTTL(ttl time.Duration) {
	//TODO implement me
	panic("implement me")
}

func (d *diskTable) DoGarbageCollection() error {
	//TODO implement me
	panic("implement me")
}

func (d *diskTable) Destroy() error {
	//TODO implement me
	panic("implement me")
}
