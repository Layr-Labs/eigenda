package core

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/beevik/ntp"
)

var (
	ntpOffsetNano int64
)

// StartNtpSync launches a background goroutine to periodically sync with NTP using a ticker and context for cancellation.
func StartNtpSync(ctx context.Context, server string, interval time.Duration, logger logging.Logger) {
	// First, verify the server string is correct by doing a blocking NTP query
	offset, err := getNtpOffset(server)
	if err != nil {
		panic("NTP sync failed on startup: invalid server config or unreachable NTP server: " + err.Error())
	}
	atomic.StoreInt64(&ntpOffsetNano, int64(offset))
	logger.Info("Initial NTP sync success", "offset_ns", offset)

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				logger.Info("NTP sync stopped")
				return
			case <-ticker.C:
				func() {
					defer func() {
						if r := recover(); r != nil {
							logger.Error("panic in NTP sync", "recover", r)
						}
					}()
					offset, err := getNtpOffset(server)
					if err != nil {
						logger.Warn("NTP sync failed", "err", err)
					} else {
						atomic.StoreInt64(&ntpOffsetNano, int64(offset))
						logger.Info("NTP sync success", "offset_ns", offset)
					}
				}()
			}
		}
	}()
}

// getNtpOffset fetches the offset between NTP and local time (in nanoseconds)
func getNtpOffset(server string) (int64, error) {
	rsp, err := ntp.Query(server)
	if err != nil {
		return 0, err
	}
	return rsp.ClockOffset.Nanoseconds(), nil
}

// NowWithNtpOffset returns the current time compensated by the latest NTP offset
// If NTP has never synced, it falls back to system time by taking ntpOffsetNano as 0
func NowWithNtpOffset() time.Time {
	offset := atomic.LoadInt64(&ntpOffsetNano)
	return time.Now().Add(time.Duration(offset))
}
