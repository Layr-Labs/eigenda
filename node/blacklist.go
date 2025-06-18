// This file contains the structs which are converted to JSON and stored in the blacklist store

package node

import (
	"encoding/json"
	"time"
)

// Blacklist contains entries of blacklisted dispersers
type Blacklist struct {
	Entries     []BlacklistEntry `json:"entries"`
	LastUpdated uint64           `json:"last_updated"`
}

// BlacklistEntry represents a single blacklist record
type BlacklistEntry struct {
	DisperserID uint32            `json:"disperser_id"`
	Metadata    BlacklistMetadata `json:"metadata"`
	Timestamp   uint64            `json:"timestamp"`
}

// BlacklistMetadata contains additional information about the blacklisting
type BlacklistMetadata struct {
	ContextId string `json:"context_id"`
	Reason    string `json:"reason"`
}

// ToBytes serializes the Blacklist to JSON bytes
func (b *Blacklist) ToBytes() ([]byte, error) {
	return json.Marshal(b)
}

// FromBytes deserializes JSON bytes into the current Blacklist
func (b *Blacklist) FromBytes(data []byte) error {
	return json.Unmarshal(data, b)
}

// AddEntry adds a new blacklist entry
func (b *Blacklist) AddEntry(disperserId uint32, contextId, reason string) {
	b.LastUpdated = uint64(time.Now().Unix())
	b.Entries = append(b.Entries, BlacklistEntry{
		DisperserID: disperserId,
		Metadata: BlacklistMetadata{
			ContextId: contextId,
			Reason:    reason,
		},
		Timestamp: b.LastUpdated,
	})
}
