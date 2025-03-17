package keymap

import (
	"github.com/Layr-Labs/eigensdk-go/logging"
)

var _ KeymapBuilder = &MemKeymapBuilder{}

// MemKeymapBuilder is a KeymapBuilder that builds MemKeymap instances.
type MemKeymapBuilder struct {
}

// NewMemKeymapBuilder creates a new LevelDBKeymapBuilder.
func NewMemKeymapBuilder() *MemKeymapBuilder {
	return &MemKeymapBuilder{}
}

func (b *MemKeymapBuilder) Type() KeymapType {
	return MemKeymapType
}

func (b *MemKeymapBuilder) Build(logger logging.Logger, paths []string) (Keymap, bool, error) {
	return NewMemKeymap(logger), true, nil
}

func (b *MemKeymapBuilder) DeleteFiles(logger logging.Logger, paths []string) error {
	return nil
}
