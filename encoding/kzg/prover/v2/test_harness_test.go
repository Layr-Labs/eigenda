package prover_test

import (
	"runtime"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/encoding/codec"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover/v2"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier/v2"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/committer"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/stretchr/testify/require"
)

type testHarness struct {
	logger                       logging.Logger
	verifierV2KzgConfig          *verifier.Config
	proverV2KzgConfig            *prover.KzgConfig
	committerConfig              *committer.Config
	numNode                      uint64
	numSys                       uint64
	numPar                       uint64
	paddedGettysburgAddressBytes []byte
}

func getTestHarness(t require.TestingT) *testHarness {
	proverConfig := &prover.KzgConfig{
		SRSNumberToLoad: 2900,
		G1Path:          "../../../../resources/srs/g1.point",
		PreloadEncoder:  true,
		CacheDir:        "../../../../resources/srs/SRSTables",
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
	}
	committerConfig := &committer.Config{
		SRSNumberToLoad:   proverConfig.SRSNumberToLoad,
		G1SRSPath:         proverConfig.G1Path,
		G2SRSPath:         "../../../../resources/srs/g2.point",
		G2TrailingSRSPath: "../../../../resources/srs/g2.trailing.point",
	}
	// Gettysburg address length is 1146 bytes.
	numNode := uint64(4)
	numSys := uint64(3)
	numPar := numNode - numSys
	paddedGettysburgAddressBytes := codec.ConvertByPaddingEmptyByte([]byte("Fourscore and seven years ago our fathers brought forth, on this continent, a new nation, conceived in liberty, and dedicated to the proposition that all men are created equal. Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived, and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting-place for those who here gave their lives, that that nation might live. It is altogether fitting and proper that we should do this. But, in a larger sense, we cannot dedicate, we cannot consecrate—we cannot hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they here gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom, and that government of the people, by the people, for the people, shall not perish from the earth."))

	return &testHarness{
		logger:                       common.TestLogger(t),
		verifierV2KzgConfig:          verifier.ConfigFromProverV2Config(proverConfig),
		proverV2KzgConfig:            proverConfig,
		committerConfig:              committerConfig,
		numNode:                      numNode,
		numSys:                       numSys,
		numPar:                       numPar,
		paddedGettysburgAddressBytes: paddedGettysburgAddressBytes,
	}
}
