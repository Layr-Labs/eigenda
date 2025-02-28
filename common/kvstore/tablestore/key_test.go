package tablestore

import (
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetName(t *testing.T) {
	tu.InitializeRandom()

	tableName := tu.RandomString(10)

	kb := newKeyBuilder(tableName, 0)
	assert.Equal(t, tableName, kb.TableName())
}

func TestBytesRoundTrip(t *testing.T) {
	tu.InitializeRandom()

	tableName := tu.RandomString(10)
	b := tu.RandomBytes(10)

	kb := newKeyBuilder(tableName, 0)
	k := kb.Key(b)
	assert.Equal(t, b, k.Bytes())
}
