package lightnode

import (
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
	"testing"
	"time"
)

func TestGetters(t *testing.T) {
	tu.InitializeRandom()

	id := rand.Uint64()
	seed := rand.Uint64()
	now := time.Unix(int64(rand.Int31()), 0)

	registration := NewRegistration(id, seed, now)

	assert.Equal(t, id, registration.ID())
	assert.Equal(t, seed, registration.Seed())
	assert.Equal(t, now, registration.RegistrationTime())
}
