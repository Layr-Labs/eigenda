package verifier_test

import (
	"runtime"

	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/committer"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover/v2"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier/v2"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
)

type testHarness struct {
	verifierV2KzgConfig          *verifier.Config
	committerConfig              *committer.Config
	proverV2KzgConfig            *prover.KzgConfig
	numNode                      uint64
	numSys                       uint64
	numPar                       uint64
	paddedGettysburgAddressBytes []byte
}

func getTestHarness() *testHarness {
	kzgConfig := &kzg.KzgConfig{
		G1Path:          "../../../../resources/srs/g1.point",
		G2Path:          "../../../../resources/srs/g2.point",
		G2TrailingPath:  "../../../../resources/srs/g2.trailing.point",
		CacheDir:        "../../../../resources/srs/SRSTables",
		SRSOrder:        4096,
		SRSNumberToLoad: 4096,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
		LoadG2Points:    true,
	}
	committerConfig := &committer.Config{
		SRSNumberToLoad:   4096,
		G1SRSPath:         "../../../../resources/srs/g1.point",
		G2SRSPath:         "../../../../resources/srs/g2.point",
		G2TrailingSRSPath: "../../../../resources/srs/g2.trailing.point",
	}
	numNode := uint64(4)
	numSys := uint64(3)
	numPar := numNode - numSys
	paddedGettysburgAddressBytes := codec.ConvertByPaddingEmptyByte([]byte("Fourscore and seven years ago our fathers brought forth, on this continent, a new nation, conceived in liberty, and dedicated to the proposition that all men are created equal. Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived, and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting-place for those who here gave their lives, that that nation might live. It is altogether fitting and proper that we should do this. But, in a larger sense, we cannot dedicate, we cannot consecrate—we cannot hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they here gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom, and that government of the people, by the people, for the people, shall not perish from the earth."))
	return &testHarness{
		verifierV2KzgConfig:          verifier.ConfigFromV1KzgConfig(kzgConfig),
		proverV2KzgConfig:            prover.KzgConfigFromV1Config(kzgConfig),
		committerConfig:              committerConfig,
		numNode:                      numNode,
		numSys:                       numSys,
		numPar:                       numPar,
		paddedGettysburgAddressBytes: paddedGettysburgAddressBytes,
	}
}
