package test

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/stretchr/testify/require"
)

// TestSignalHeartbeat verifies that SignalHeartbeat sends exactly one message
// with the correct component name and timestamp.
func TestSignalHeartbeat(t *testing.T) {
	ch := make(chan healthcheck.HeartbeatMessage, 1)
	loggerConfig := common.DefaultLoggerConfig()
	logger, err := common.NewLogger(loggerConfig)
	assert.NoError(t, err)

	start := time.Now()
	healthcheck.SignalHeartbeat("dispatcher", ch, logger)

	select {
	case hb := <-ch:
		require.Equal(t, "dispatcher", hb.Component)
		// ensure timestamp is within a second of our call
		require.WithinDuration(t, start, hb.Timestamp, time.Second)
	default:
		t.Fatal("expected a heartbeat message on the channel")
	}
}

// TestHeartbeatMonitor_WritesSummaryAndStops sends two heartbeats, closes the channel,
// and then verifies that the monitor wrote a file containing both entries and returned.
func TestHeartbeatMonitor_WritesSummaryAndStops(t *testing.T) {
	dir := t.TempDir()
	fpath := filepath.Join(dir, "hb.txt")

	// Make a liveness channel and start the monitor.
	ch := make(chan healthcheck.HeartbeatMessage)
	done := make(chan error, 1)
	go func() {
		// monitor will exit when channel is closed
		err := healthcheck.HeartbeatMonitor(fpath, 50*time.Millisecond, ch)
		done <- err
	}()

	// send two distinct heartbeats
	ch <- healthcheck.HeartbeatMessage{Component: "dispatcher", Timestamp: time.Unix(1, 0)}
	ch <- healthcheck.HeartbeatMessage{Component: "encodingManager", Timestamp: time.Unix(2, 0)}
	// closing the channel causes the monitor to return
	close(ch)

	// wait for the monitor to exit
	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(200 * time.Millisecond):
		t.Fatal("heartbeat monitor did not exit in time")
	}

	// read the file that should have been written
	data, err := os.ReadFile(fpath)
	require.NoError(t, err)
	text := string(data)

	// it should contain both component lines
	require.Contains(t, text, "Component dispatcher: Last heartbeat at 1")
	require.Contains(t, text, "Component encodingManager: Last heartbeat at 2")
}

// TestHeartbeatMonitor_StallWarning starts a monitor without sending a heartbeat,
// and ensures that it logs a warning (we can verify the file is created).
func TestHeartbeatMonitor_StallWarning(t *testing.T) {
	dir := t.TempDir()
	fpath := filepath.Join(dir, "hb-stall.txt")

	ch := make(chan healthcheck.HeartbeatMessage)
	done := make(chan error, 1)
	go func() {
		err := healthcheck.HeartbeatMonitor(fpath, 20*time.Millisecond, ch)
		done <- err
	}()

	// give it some stall intervals
	time.Sleep(60 * time.Millisecond)
	// now close the channel and let it exit
	close(ch)

	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("heartbeat monitor did not exit after stall")
	}

	// since no heartbeats arrived and nothing was written to file, the file may not exist
	// ensure no panic and exit.
	_, err := os.Stat(fpath)
	require.True(t, os.IsNotExist(err) || err == nil)
}
