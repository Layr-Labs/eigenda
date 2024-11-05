package eth

import (
	"math/big"
	"slices"

	"github.com/Layr-Labs/eigenda/core"

	eigendasrvmg "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"

	"github.com/ethereum/go-ethereum/crypto"
)

var (
	maxNumberOfQuorums = 192
)

type BN254G1Point struct {
	X *big.Int
	Y *big.Int
}

type BN254G2Point struct {
	X [2]*big.Int
	Y [2]*big.Int
}

func signatureToBN254G1Point(s *core.Signature) eigendasrvmg.BN254G1Point {
	return eigendasrvmg.BN254G1Point{
		X: s.X.BigInt(new(big.Int)),
		Y: s.Y.BigInt(new(big.Int)),
	}
}

func pubKeyG1ToBN254G1Point(p *core.G1Point) eigendasrvmg.BN254G1Point {
	return eigendasrvmg.BN254G1Point{
		X: p.X.BigInt(new(big.Int)),
		Y: p.Y.BigInt(new(big.Int)),
	}
}

func pubKeyG2ToBN254G2Point(p *core.G2Point) eigendasrvmg.BN254G2Point {
	return eigendasrvmg.BN254G2Point{
		X: [2]*big.Int{p.X.A1.BigInt(new(big.Int)), p.X.A0.BigInt(new(big.Int))},
		Y: [2]*big.Int{p.Y.A1.BigInt(new(big.Int)), p.Y.A0.BigInt(new(big.Int))},
	}
}

func quorumIDsToQuorumNumbers(quorumIds []core.QuorumID) []byte {
	quorumNumbers := make([]byte, len(quorumIds))
	for i, quorumId := range quorumIds {
		quorumNumbers[i] = byte(quorumId)
	}
	return quorumNumbers
}

func quorumParamsToQuorumNumbers(quorumParams map[core.QuorumID]*core.QuorumResult) []byte {
	quorumNumbers := make([]byte, len(quorumParams))
	quorums := make([]uint8, len(quorumParams))
	i := 0
	for k := range quorumParams {
		quorums[i] = k
		i++
	}
	slices.Sort(quorums)
	i = 0
	for _, quorum := range quorums {
		qp := quorumParams[quorum]
		quorumNumbers[i] = byte(qp.QuorumID)
		i++
	}
	return quorumNumbers
}

func serializeSignedStakeForQuorums(quorumParams map[core.QuorumID]*core.QuorumResult) []byte {
	thresholdPercentages := make([]byte, len(quorumParams))
	quorums := make([]uint8, len(quorumParams))
	i := 0
	for k := range quorumParams {
		quorums[i] = k
		i++
	}
	slices.Sort(quorums)
	i = 0
	for _, quorum := range quorums {
		qp := quorumParams[quorum]
		thresholdPercentages[i] = byte(qp.PercentSigned)
		i++
	}
	return thresholdPercentages
}

func HashPubKeyG1(pk *core.G1Point) [32]byte {
	gp := pubKeyG1ToBN254G1Point(pk)
	xBytes := make([]byte, 32)
	yBytes := make([]byte, 32)
	gp.X.FillBytes(xBytes)
	gp.Y.FillBytes(yBytes)
	return crypto.Keccak256Hash(append(xBytes, yBytes...))
}

func BitmapToQuorumIds(bitmap *big.Int) []core.QuorumID {
	// loop through each index in the bitmap to construct the array

	quorumIds := make([]core.QuorumID, 0, maxNumberOfQuorums)
	for i := 0; i < maxNumberOfQuorums; i++ {
		if bitmap.Bit(i) == 1 {
			quorumIds = append(quorumIds, core.QuorumID(i))
		}
	}
	return quorumIds
}

func bitmapToBytesArray(bitmap *big.Int) []byte {
	// initialize an empty uint64 to be used as a bitmask inside the loop
	var (
		bytesArray []byte
	)
	// loop through each index in the bitmap to construct the array
	for i := 0; i < maxNumberOfQuorums; i++ {
		// check if the i-th bit is flipped in the bitmap
		if bitmap.Bit(i) == 1 {
			// if the i-th bit is flipped, then add a byte encoding the value 'i' to the `bytesArray`
			bytesArray = append(bytesArray, byte(uint8(i)))
		}
	}
	return bytesArray
}
