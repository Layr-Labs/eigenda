package dataapi

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/dataapi/subgraph"
	"github.com/consensys/gnark-crypto/ecc/bn254"
)

func ConvertHexadecimalToBytes(byteHash []byte) ([32]byte, error) {
	hexString := strings.TrimPrefix(string(byteHash), "0x")

	// Now decode the hex string to bytes
	decodedBytes, err := hex.DecodeString(hexString)
	if err != nil {
		return [32]byte{}, err
	}

	// We expect the resulting byte slice to have a length of 32 bytes.
	if len(decodedBytes) != 32 {
		return [32]byte{}, errors.New("error decoding hash, invalid length")
	}

	// Convert the byte slice to a [32]byte array
	var byteArray [32]byte
	copy(byteArray[:], decodedBytes[:32])

	return byteArray, nil
}

func ConvertNanosecondToSecond(timestamp uint64) uint64 {
	return timestamp / uint64(time.Second)
}

func ConvertOperatorInfoGqlToIndexedOperatorInfo(operator *subgraph.IndexedOperatorInfo) (*core.IndexedOperatorInfo, error) {
	if operator == nil {
		return nil, errors.New("operator is nil")
	}

	if len(operator.SocketUpdates) == 0 {
		return nil, errors.New("no socket updates found for operator")
	}

	pubkeyG1 := new(bn254.G1Affine)
	_, err := pubkeyG1.X.SetString(string(operator.PubkeyG1_X))
	if err != nil {
		return nil, fmt.Errorf("failed to set PubkeyG1_X: %v", err)
	}
	_, err = pubkeyG1.Y.SetString(string(operator.PubkeyG1_Y))
	if err != nil {
		return nil, fmt.Errorf("failed to set PubkeyG1_Y: %v", err)
	}

	if len(operator.PubkeyG2_X) < 2 || len(operator.PubkeyG2_Y) < 2 {
		return nil, errors.New("incomplete PubkeyG2 coordinates")
	}

	pubkeyG2 := new(bn254.G2Affine)
	_, err = pubkeyG2.X.A1.SetString(string(operator.PubkeyG2_X[0]))
	if err != nil {
		return nil, fmt.Errorf("failed to set PubkeyG2_X[0]: %v", err)
	}
	_, err = pubkeyG2.X.A0.SetString(string(operator.PubkeyG2_X[1]))
	if err != nil {
		return nil, fmt.Errorf("failed to set PubkeyG2_X[1]: %v", err)
	}
	_, err = pubkeyG2.Y.A1.SetString(string(operator.PubkeyG2_Y[0]))
	if err != nil {
		return nil, fmt.Errorf("failed to set PubkeyG2_Y[0]: %v", err)
	}
	_, err = pubkeyG2.Y.A0.SetString(string(operator.PubkeyG2_Y[1]))
	if err != nil {
		return nil, fmt.Errorf("failed to set PubkeyG2_Y[1]: %v", err)
	}

	return &core.IndexedOperatorInfo{
		PubkeyG1: &core.G1Point{G1Affine: pubkeyG1},
		PubkeyG2: &core.G2Point{G2Affine: pubkeyG2},
		Socket:   string(operator.SocketUpdates[0].Socket),
	}, nil
}
