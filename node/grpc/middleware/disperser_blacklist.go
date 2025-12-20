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

	// strikeWindow is the duration in which invalid requests count toward blacklisting.
	// For example, strikeWindow=2m and maxInvalid=3 means 3 invalids within the last 2 minutes triggers a ban.
	strikeWindow time.Duration
	maxInvalid   int

	mu    sync.Mutex
	state map[uint32]*disperserBlacklistState
}

type disperserBlacklistState struct {
	blacklistedUntil time.Time
	invalidTimes     []time.Time
}

func NewDisperserBlacklist(
	logger logging.Logger,
	ttl time.Duration,
	strikeWindow time.Duration,
	maxInvalid int,
) *DisperserBlacklist {
	return &DisperserBlacklist{
		logger:       logger,
		ttl:          ttl,
		strikeWindow: strikeWindow,
		maxInvalid:   maxInvalid,
		state:        make(map[uint32]*disperserBlacklistState),
	}
}

// RecordInvalid records an invalid request from the disperser and blacklists them if the configured
// threshold is exceeded within the strikeWindow.
func (b *DisperserBlacklist) RecordInvalid(disperserID uint32, now time.Time, reason string) {
	if b == nil || b.ttl <= 0 || b.maxInvalid <= 0 {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	s := b.getOrCreateStateLocked(disperserID)

	// If currently blacklisted, keep it simple and do nothing (avoid extending bans).
	if !s.blacklistedUntil.IsZero() && now.Before(s.blacklistedUntil) {
		return
	}

	// Prune old invalid timestamps outside the strike window.
	if b.strikeWindow > 0 {
		cutoff := now.Add(-b.strikeWindow)
		s.invalidTimes = pruneTimesBefore(s.invalidTimes, cutoff)
	}

	s.invalidTimes = append(s.invalidTimes, now)

	if len(s.invalidTimes) < b.maxInvalid {
		return
	}

	// Threshold exceeded: blacklist and clear strikes so the disperser is "forgiven" after the TTL expires.
	s.invalidTimes = nil
	s.blacklistedUntil = now.Add(b.ttl)

	if b.logger != nil {
		b.logger.Warn("blacklisting disperser for invalid requests",
			"disperserID", disperserID,
			"strikeWindow", b.strikeWindow.String(),
			"maxInvalid", b.maxInvalid,
			"blacklistDuration", b.ttl.String(),
			"reason", reason)
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

	s, ok := b.state[disperserID]
	if !ok || s == nil {
		return false
	}

	if s.blacklistedUntil.IsZero() {
		return false
	}

	if !now.Before(s.blacklistedUntil) {
		// Expired: clear state (including strikes) to ensure forgiveness.
		delete(b.state, disperserID)
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
	s := b.getOrCreateStateLocked(disperserID)
	s.invalidTimes = nil
	s.blacklistedUntil = now.Add(b.ttl)
	b.mu.Unlock()

	if b.logger != nil {
		b.logger.Warn("blacklisting disperser for invalid request",
			"disperserID", disperserID,
			"blacklistDuration", b.ttl.String(),
			"reason", reason)
	}
}

func (b *DisperserBlacklist) getOrCreateStateLocked(disperserID uint32) *disperserBlacklistState {
	s := b.state[disperserID]
	if s == nil {
		s = &disperserBlacklistState{}
		b.state[disperserID] = s
	}
	return s
}

func pruneTimesBefore(times []time.Time, cutoff time.Time) []time.Time {
	// Find first index >= cutoff and reslice.
	i := 0
	for i < len(times) && times[i].Before(cutoff) {
		i++
	}
	if i == 0 {
		return times
	}
	if i >= len(times) {
		return nil
	}
	return times[i:]
}
