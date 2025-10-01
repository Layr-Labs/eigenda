package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type Foo struct {
	String    string
	Int       int
	Int64     int64
	Int32     int32
	Int16     int16
	Int8      int8
	Uint      uint
	Uint64    uint64
	Uint32    uint32
	Uint16    uint16
	Uint8     uint8
	Float64   float64
	Float32   float32
	Bool      bool
	Recursive *Foo
	Bar       Bar
	Baz       *Baz
}

func (f *Foo) Verify() error {
	return nil
}

type Bar struct {
	A   string
	B   int
	C   bool
	Baz *Baz
}

func (b *Bar) Verify() error {
	return nil
}

type Baz struct {
	X string
	Y int
	Z bool
}

func (b *Baz) Verify() error {
	return nil
}

func TestTOMLParsing(t *testing.T) {

	configFile := "testdata/config.toml"

	var foo Foo
	err := ParseConfig(&foo, "FOO", configFile)
	require.NoError(t, err)

	// TODO remove this debug print before merging
	fmt.Printf("%+v\n", foo)

}
