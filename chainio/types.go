package chainio

import (
	"encoding/hex"
	"errors"
	"strings"
)

const (
	// We use uint8 to count the number of quorums, so we can have at most 255 quorums,
	// which means the max ID can not be larger than 254 (from 0 to 254, there are 255
	// different IDs).
	MaxQuorumID = 254
)

type QuorumID = uint8

type OperatorID = [32]byte

func GetOperatorHex(id OperatorID) string {
	return hex.EncodeToString(id[:])
}

// The "s" is an operatorId in hex string format, which may or may not have the "0x" prefix.
func OperatorIDFromHex(s string) (OperatorID, error) {
	opID := [32]byte{}
	s = strings.TrimPrefix(s, "0x")
	if len(s) != 64 {
		return OperatorID(opID), errors.New("operatorID hex string must be 64 bytes, or 66 bytes if starting with 0x")
	}
	opIDslice, err := hex.DecodeString(s)
	if err != nil {
		return OperatorID(opID), err
	}
	copy(opID[:], opIDslice)
	return OperatorID(opID), nil
}
