package tablestore

import (
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestGetName(t *testing.T) {
	tu.InitializeRandom()

	tableName := tu.RandomString(10)

	kb := newKeyBuilder(tableName, 0)
	assert.Equal(t, tableName, kb.TableName())
}

func TestGetBuilder(t *testing.T) {
	tu.InitializeRandom()

	tableName := tu.RandomString(10)

	kb := newKeyBuilder(tableName, 0)
	k := kb.StringKey("asdf")
	assert.Same(t, kb, k.Builder())
}

func TestStringRoundTrip(t *testing.T) {
	tu.InitializeRandom()

	tableName := tu.RandomString(10)
	str := tu.RandomString(10)

	kb := newKeyBuilder(tableName, 0)
	k := kb.StringKey(str)
	assert.Equal(t, str, k.AsString())
}

func TestBytesRoundTrip(t *testing.T) {
	tu.InitializeRandom()

	tableName := tu.RandomString(10)
	b := tu.RandomBytes(10)

	kb := newKeyBuilder(tableName, 0)
	k := kb.Key(b)
	assert.Equal(t, b, k.AsBytes())
}

func TestUint64RoundTrip(t *testing.T) {
	tu.InitializeRandom()

	tableName := tu.RandomString(10)
	u := rand.Uint64()

	kb := newKeyBuilder(tableName, 0)
	k := kb.Uint64Key(u)
	u2, err := k.AsUint64()
	assert.NoError(t, err)
	assert.Equal(t, u, u2)
}

func TestInt64RoundTrip(t *testing.T) {
	tu.InitializeRandom()

	tableName := tu.RandomString(10)
	i := rand.Int63()

	kb := newKeyBuilder(tableName, 0)
	k := kb.Int64Key(i)
	u2, err := k.AsInt64()
	assert.NoError(t, err)
	assert.Equal(t, i, u2)
}

func TestUint32RoundTrip(t *testing.T) {
	tu.InitializeRandom()

	tableName := tu.RandomString(10)
	u := rand.Uint32()

	kb := newKeyBuilder(tableName, 0)
	k := kb.Uint32Key(u)
	u2, err := k.AsUint32()
	assert.NoError(t, err)
	assert.Equal(t, u, u2)
}

func TestInt32RoundTrip(t *testing.T) {
	tu.InitializeRandom()

	tableName := tu.RandomString(10)
	i := rand.Int31()

	kb := newKeyBuilder(tableName, 0)
	k := kb.Int32Key(i)
	u2, err := k.AsInt32()
	assert.NoError(t, err)
	assert.Equal(t, i, u2)
}
