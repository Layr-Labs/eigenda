package examples

import (
	"crypto/rand"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
)

// createRandomPayload creates a payload with random data of the specified size in bytes
func createRandomPayload(byteCount int) (*coretypes.Payload, error) {
	payloadBytes := make([]byte, byteCount)
	_, err := rand.Read(payloadBytes)
	if err != nil {
		return nil, err
	}
	return coretypes.NewPayload(payloadBytes), nil
}
