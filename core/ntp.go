package core

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/beevik/ntp"
)

// NTPSyncedClock provides synchronized time based on NTP offset.
type NTPSyncedClock struct {
	offset int64
	logger logging.Logger
}

// NewNTPSyncedClock creates a new NTP synchronized clock and starts background sync.
func NewNTPSyncedClock(ctx context.Context, server string, syncInterval time.Duration, logger logging.Logger) (*NTPSyncedClock, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger must not be nil")
	}

	clock := &NTPSyncedClock{
		logger: logger.With("component", "NTPSyncedClock"),
	}

	// Try to sync once, but don't fail if it doesn't work
	if err := clock.syncOnce(server); err != nil {
		clock.logger.Warn("Initial NTP sync failed, will retry in background and use system clock until sync succeeds", "err", err)
	}

	go clock.runSyncLoop(ctx, server, syncInterval)

	return clock, nil
}

// runSyncLoop periodically syncs the clock with NTP.
func (c *NTPSyncedClock) runSyncLoop(ctx context.Context, server string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("NTP sync stopped")
			return
		case <-ticker.C:
			_ = c.syncOnce(server)
		}
	}
}

// syncOnce performs a single NTP sync and updates the offset.
func (c *NTPSyncedClock) syncOnce(server string) error {
	offset, err := ntpOffset(server)
	if err != nil {
		c.logger.Warn("NTP sync failed", "err", err)
		return err
	}
	atomic.StoreInt64(&c.offset, offset)
	c.logger.Debug("NTP sync success", "offset_ns", offset)
	return nil
}

// Now returns the current time compensated by the latest NTP offset.
// If the NTP sync has not yet succeeded, it returns the current system time.
func (c *NTPSyncedClock) Now() time.Time {
	offset := atomic.LoadInt64(&c.offset)
	return time.Now().Add(time.Duration(offset))
}

// ntpOffset fetches the offset between NTP and local time (in nanoseconds).
func ntpOffset(server string) (int64, error) {
	rsp, err := ntp.Query(server)
	if err != nil {
		return 0, fmt.Errorf("ntp query to server %q failed: %w", server, err)
	}
	return rsp.ClockOffset.Nanoseconds(), nil
}
