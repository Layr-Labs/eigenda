// This file contains the structs which are converted to JSON and stored in the blacklist store

package node

import (
	"encoding/json"
	"time"
)

// Blacklist contains entries of blacklisted dispersers
type Blacklist struct {
	Entries []BlacklistEntry `json:"entries"`
}

// BlacklistEntry represents a single blacklist record
type BlacklistEntry struct {
	DisperserAddress []byte            `json:"disperser_address"`
	Metadata         BlacklistMetadata `json:"metadata"`
	Timestamp        uint64            `json:"timestamp"`
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
func (b *Blacklist) AddEntry(disperserAddr []byte, contextID, reason string) {
	b.Entries = append(b.Entries, BlacklistEntry{
		DisperserAddress: disperserAddr,
		Metadata: BlacklistMetadata{
			ContextId: contextID,
			Reason:    reason,
		},
		Timestamp: uint64(time.Now().Unix()),
	})
}
