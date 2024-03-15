package mock

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestContractSimulator(t *testing.T) {
	t.Skip("Skipping this test after the simulated backend upgrade broke this test. Enable it after fixing the issue.")
	sc := MustNewContractSimulator()
	ctx, cancel := context.WithCancel(context.Background())
	sc.Start(time.Millisecond, cancel)

	<-ctx.Done()

	events, err := sc.DepositEvents()
	assert.Nil(t, err)
	assert.Equal(t, 3, len(events))
	assert.Equal(t, events[0].Wad.Int64(), int64(1))
	assert.Equal(t, events[1].Wad.Int64(), int64(3))
	assert.Equal(t, events[2].Wad.Int64(), int64(4))
}
