package keymap

import (
	"github.com/Layr-Labs/eigensdk-go/logging"
)

var _ KeyMapBuilder = &MemKeyMapBuilder{}

// MemKeyMapBuilder is a KeyMapBuilder that builds MemKeyMap instances.
type MemKeyMapBuilder struct {
}

// NewMemKeyMapBuilder creates a new LevelDBKeyMapBuilder.
func NewMemKeyMapBuilder() *MemKeyMapBuilder {
	return &MemKeyMapBuilder{}
}

func (b *MemKeyMapBuilder) Type() KeyMapType {
	return MemKeyMapType
}

func (b *MemKeyMapBuilder) Build(logger logging.Logger, paths []string) (KeyMap, bool, error) {
	return NewMemKeyMap(logger), true, nil
}

func (b *MemKeyMapBuilder) DeleteFiles(logger logging.Logger, paths []string) error {
	return nil
}
