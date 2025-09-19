// TODO(samlaf): hexify the G2 points and move G1/G2PowerOf2SRS to encoding/constants.go
package srs

import (
	_ "embed"
	"encoding/hex"
	"fmt"

	"github.com/consensys/gnark-crypto/ecc/bn254"
)

//go:embed g2.point.powerOf2
var serializedG2PowerOf2Data []byte

// G2PowerOf2SRS contains 28 G2 points: [1] [tau^1] [tau^2]..[tau^{2^27}]
//
// the powerOf2 file, only [tau^exp] are stored.
// exponent    0,    1,       2,    , ..
// actual pow [tau],[tau^2],[tau^4],.. (stored in the file)
// In our convention SRSOrder contains the total number of series of g1, g2 starting with generator
// i.e. [1] [tau] [tau^2]..
// So the actual power of tau is SRSOrder - 1
// The mainnet SRS, the max power is 2^28-1, so the last power in powerOf2 file is [tau^(2^27)]
var G2PowerOf2SRS []bn254.G2Affine

// Contains 28 points represented by G1SRSFile[2^28 - 2^i] for i in 0..28
var G1PowerOf2SRS []bn254.G1Affine

func init() {
	// Note that we can't use bn254.NewDecoder(bytes.NewReader(serializedG2PowerOf2Data)).Decode(&G2PowerOf2Data)
	// because the file was not encoded using gnark-crypto's encoder.
	// It only contains the 28 raw serialized points, each taking 64 bytes.
	// gnark-crypto's encoder/decoder adds a 4 bytes header.
	for pointIndex := 0; pointIndex < len(serializedG2PowerOf2Data); pointIndex += 64 {
		serializedPoint := serializedG2PowerOf2Data[pointIndex : pointIndex+64]
		var p bn254.G2Affine
		if _, err := p.SetBytes(serializedPoint); err != nil {
			panic(fmt.Sprintf("error deserializing G2 point at index %d: %v", pointIndex, err))
		}
		G2PowerOf2SRS = append(G2PowerOf2SRS, p)
	}
	G1PowerOf2HexStrings := []string{
		"d902537c5ac68b39468f8cfcc46b00da353024b618b0454e6847d2aee530e850",
		"c672cec11e3d0c0096d550635d28c4b51dd3c2deb407a5985f458f5a8610fe94",
		"c2ee739ae261af377f3c65362049fef9402013bd898351eb60e0d7429a56f880",
		"d9d7adde8d4e55e87312fcb7713fdd1e8713358347d7c8bc1ac21fa1ef1db34a",
		"8f89ac9e77846d29af1c58af31cfc3fe61c45ca14cf0a10f7aa71605b868c0e7",
		"e3adc6dbd3cb5c694b17bdb974860d07222f9ffd75655d8800702f455ddffc97",
		"cb36c6a5486f20baf22e4ca283ee475b70447c35e16b749d710caa4bc7bf2ac9",
		"86aa256a555695131099bc873d2d23dd0f31d6fa9e0117f2dc913565380f8536",
		"ac25b7bfc07913ee0aa062624b6c06f0218dcf1f03f232397761463d079be7fe",
		"a219ecf6ee97fe6b32ab434e3daff9ed6cbe8f467979ad1c6f39ee9dc660b212",
		"97840fcc4dd766fa0748bf2d50dd85242d4bb031fac39f8bb8f12d1146b0d443",
		"95b23343161483eed8834b7d595f5ac18b8c5dfe85d40538c83027a12b7e0ca6",
		"82e393213e64aef7726afb4ee823b58634087a47330ef6150c6e9e496a18cc5a",
		"e11a52ebf3628f51663a2fb41de1b755d28b764dde082521072531cf53bbb895",
		"c6318eb9dfa5f5627ceddde2f026af4e3ef79b7f1702f497b353437f57f188e4",
		"cfdc1c150ef291fa5eb1bdd2743815eebea02b4a89b3b0b1cc801269fb2502d6",
		"88723a42d3025fbb3beb27a75cf1266e37c59959a434ccabd04c332c888afc7c",
		"ebc857866d0cbb6ce20c2abd612cb99d1d2f4446e5330255c77f69bcfc56c8be",
		"ecc72a85bca27e6e6d9dcc73e15bd528b9bb1b4dc5b158e87d821571859820eb",
		"899e0d8eda7fcd5d0fcb7488574361663a7d05e920c11643101c1c996aad21b7",
		"87f814468e6e5b08526830fe3ce8fde7b5385f53e7654d3c061f5f602c5452b6",
		"a26d44f770db3207696477d61e7feedc3f0a83ea58f37b9ef914834fb32895b8",
		"90ec8f5ba15034bf2faa5b650606b7786e3c5c16201488c4411de3a40476874c",
		"ad10e224d82572833b2854c327a5db10a0b6c617c367e3aff58f5862aab90a41",
		"8ccb85c07ad9092316ea6f95161e0a64ed7cca863f23bc22300225bf456d094a",
		"d20165a1b364337df11a35fb687aa62382236938f8f740cb7059b656e1f4dd1c",
		"a9d669092e951729fcc2eaf05ff706cf372e04cbde166f48833337fa37b69537",
		"eb34e5696bcd208899dbd9d1e7604ec39cc594eeedae3eaf40ff8695ab25ca72",
		"8000000000000000000000000000000000000000000000000000000000000001",
	}
	for i := 0; i < 28; i++ {
		G1PowerOf2SRS = append(G1PowerOf2SRS, toG1Affine(G1PowerOf2HexStrings[i]))
	}
}

func toG1Affine(pointHex string) bn254.G1Affine {
	var p bn254.G1Affine
	pointBytes, err := hex.DecodeString(pointHex)
	if err != nil {
		panic(fmt.Sprintf("error decoding hex string %s: %v", pointHex, err))
	}
	if _, err := p.SetBytes(pointBytes); err != nil {
		panic(fmt.Sprintf("error deserializing G1 point %s: %v", pointHex, err))
	}
	return p
}
