package middleware

import (
	"sync"
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

// DisperserBlacklist is a temporary blacklist keyed by disperser ID.
//
// This is intended as an abuse mitigation: if a disperser sends invalid requests
// that an honest disperser should never forward, we temporarily refuse requests
// from that disperser to protect validator resources.
//
// The blacklist is best-effort and local to the process (not persisted).
type DisperserBlacklist struct {
	logger logging.Logger
	ttl    time.Duration

	mu    sync.Mutex
	until map[uint32]time.Time
}

func NewDisperserBlacklist(logger logging.Logger, ttl time.Duration) *DisperserBlacklist {
	return &DisperserBlacklist{
		logger: logger,
		ttl:    ttl,
		until:  make(map[uint32]time.Time),
	}
}

// IsBlacklisted returns true if the disperser is currently blacklisted.
// Expired entries are pruned on lookup.
func (b *DisperserBlacklist) IsBlacklisted(disperserID uint32, now time.Time) bool {
	if b == nil || b.ttl <= 0 {
		return false
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	t, ok := b.until[disperserID]
	if !ok {
		return false
	}
	if now.After(t) {
		delete(b.until, disperserID)
		return false
	}
	return true
}

// Blacklist temporarily blacklists the disperser for the configured TTL.
func (b *DisperserBlacklist) Blacklist(disperserID uint32, now time.Time, reason string) {
	if b == nil || b.ttl <= 0 {
		return
	}

	b.mu.Lock()
	b.until[disperserID] = now.Add(b.ttl)
	b.mu.Unlock()

	if b.logger != nil {
		b.logger.Warn("blacklisting disperser for invalid request",
			"disperserID", disperserID,
			"forgivenessWindow", b.ttl.String(),
			"reason", reason)
	}
}
