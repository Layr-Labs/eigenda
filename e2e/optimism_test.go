package e2e_test

import (
	"testing"

	"github.com/Layr-Labs/eigenda-proxy/commitments"
	"github.com/Layr-Labs/eigenda-proxy/e2e"
	altda "github.com/ethereum-optimism/optimism/op-alt-da"
	"github.com/ethereum-optimism/optimism/op-e2e/config"
	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils"
	"github.com/ethereum-optimism/optimism/op-service/sources"
	"github.com/ethereum-optimism/optimism/op-service/testlog"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"

	actions "github.com/ethereum-optimism/optimism/op-e2e/actions/helpers"
)

var defaultAlloc = &e2eutils.AllocParams{PrefundTestUsers: true}

// L2AltDA is a test harness for manipulating altda DA state.
type L2AltDA struct {
	log       log.Logger
	storage   *altda.DAClient
	daMgr     *altda.DA
	altdaCfg  altda.Config
	batcher   *actions.L2Batcher
	sequencer *actions.L2Sequencer
	engine    *actions.L2Engine
	engCl     *sources.EngineClient
	sd        *e2eutils.SetupData
	dp        *e2eutils.DeployParams
	miner     *actions.L1Miner
}

func (a *L2AltDA) ActL1Blocks(t actions.Testing, n uint64) {
	for i := uint64(0); i < n; i++ {
		a.miner.ActL1StartBlock(12)(t)
		a.miner.ActL1EndBlock(t)
	}
}

func NewL2AltDA(t actions.Testing, daHost string, altDA bool) *L2AltDA {
	p := &e2eutils.TestParams{
		MaxSequencerDrift:   40,
		SequencerWindowSize: 120,
		ChannelTimeout:      120,
		L1BlockTime:         12,
		UseAltDA:            true,
		AllocType:           config.AllocTypeAltDA,
	}

	log := testlog.Logger(t, log.LvlDebug)

	// config.DeployConfig.DACommitmentType = altda.GenericCommitmentString
	dp := e2eutils.MakeDeployParams(t, p)
	dp.DeployConfig.DAChallengeProxy = common.Address{0x42}
	sd := e2eutils.Setup(t, dp, defaultAlloc)

	require.True(t, sd.RollupCfg.AltDAEnabled())

	miner := actions.NewL1Miner(t, log, sd.L1Cfg)
	l1Client := miner.EthClient()

	jwtPath := e2eutils.WriteDefaultJWT(t)
	engine := actions.NewL2Engine(t, log, sd.L2Cfg, jwtPath)
	engCl := engine.EngineClient(t, sd.RollupCfg)

	var storage *altda.DAClient
	if !altDA {
		storage = altda.NewDAClient(daHost, true, true)
	} else {
		storage = altda.NewDAClient(daHost, false, false)
	}

	l1F, err := sources.NewL1Client(miner.RPCClient(), log, nil, sources.L1ClientDefaultConfig(sd.RollupCfg, false, sources.RPCKindBasic))
	require.NoError(t, err)

	altdaCfg, err := sd.RollupCfg.GetOPAltDAConfig()
	require.NoError(t, err)

	if altDA {
		altdaCfg.CommitmentType = altda.GenericCommitmentType
	} else {
		altdaCfg.CommitmentType = altda.Keccak256CommitmentType
	}

	daMgr := altda.NewAltDAWithStorage(log, altdaCfg, storage, &altda.NoopMetrics{})

	enabled := sd.RollupCfg.AltDAEnabled()
	require.True(t, enabled)

	sequencer := actions.NewL2Sequencer(t, log, l1F, miner.BlobStore(), daMgr, engCl, sd.RollupCfg, 0, nil)
	miner.ActL1SetFeeRecipient(common.Address{'A'})
	sequencer.ActL2PipelineFull(t)

	batcher := actions.NewL2Batcher(log, sd.RollupCfg, actions.AltDABatcherCfg(dp, storage), sequencer.RollupClient(), l1Client, engine.EthClient(), engCl)

	return &L2AltDA{
		log:       log,
		storage:   storage,
		daMgr:     daMgr,
		altdaCfg:  altdaCfg,
		batcher:   batcher,
		sequencer: sequencer,
		engine:    engine,
		engCl:     engCl,
		sd:        sd,
		dp:        dp,
		miner:     miner,
	}
}

func (a *L2AltDA) ActL1Finalized(t actions.Testing) {
	latest := uint64(2)
	a.miner.ActL1Safe(t, latest)
	a.miner.ActL1Finalize(t, latest)
	a.sequencer.ActL1FinalizedSignal(t)
}

func TestOptimismKeccak256Commitment(gt *testing.T) {
	if !runIntegrationTests && !runTestnetIntegrationTests {
		gt.Skip("Skipping test as INTEGRATION or TESTNET env var not set")
	}

	testCfg := e2e.TestConfig(useMemory())
	testCfg.UseKeccak256ModeS3 = true

	tsConfig := e2e.TestSuiteConfig(testCfg)
	proxyTS, shutDown := e2e.CreateTestSuite(tsConfig)
	defer shutDown()

	t := actions.NewDefaultTesting(gt)

	optimism := NewL2AltDA(t, proxyTS.Address(), false)

	// build L1 block #1
	optimism.ActL1Blocks(t, 1)
	optimism.miner.ActL1SafeNext(t)

	// Fill with l2 blocks up to the L1 head
	optimism.sequencer.ActL1HeadSignal(t)
	optimism.sequencer.ActBuildToL1Head(t)

	optimism.sequencer.ActL2PipelineFull(t)
	optimism.sequencer.ActL1SafeSignal(t)
	require.Equal(t, uint64(1), optimism.sequencer.SyncStatus().SafeL1.Number)

	// add L1 block #2
	optimism.ActL1Blocks(t, 1)
	optimism.miner.ActL1SafeNext(t)
	optimism.miner.ActL1FinalizeNext(t)
	optimism.sequencer.ActL1HeadSignal(t)
	optimism.sequencer.ActBuildToL1Head(t)

	// Catch up derivation
	optimism.sequencer.ActL2PipelineFull(t)
	optimism.sequencer.ActL1FinalizedSignal(t)
	optimism.sequencer.ActL1SafeSignal(t)

	// commit all the l2 blocks to L1
	optimism.batcher.ActSubmitAll(t)
	optimism.miner.ActL1StartBlock(12)(t)
	optimism.miner.ActL1IncludeTx(optimism.dp.Addresses.Batcher)(t)
	optimism.miner.ActL1EndBlock(t)

	// verify
	optimism.sequencer.ActL2PipelineFull(t)
	optimism.ActL1Finalized(t)

	requireDispersalRetrievalEigenDA(gt, proxyTS.Metrics.HTTPServerRequestsTotal, commitments.OptimismKeccak)
}

func TestOptimismGenericCommitment(gt *testing.T) {
	if !runIntegrationTests && !runTestnetIntegrationTests {
		gt.Skip("Skipping test as INTEGRATION or TESTNET env var not set")
	}

	tsConfig := e2e.TestSuiteConfig(e2e.TestConfig(useMemory()))
	proxyTS, shutDown := e2e.CreateTestSuite(tsConfig)
	defer shutDown()

	t := actions.NewDefaultTesting(gt)

	optimism := NewL2AltDA(t, proxyTS.Address(), true)

	// build L1 block #1
	optimism.ActL1Blocks(t, 1)
	optimism.miner.ActL1SafeNext(t)

	// Fill with l2 blocks up to the L1 head
	optimism.sequencer.ActL1HeadSignal(t)
	optimism.sequencer.ActBuildToL1Head(t)

	optimism.sequencer.ActL2PipelineFull(t)
	optimism.sequencer.ActL1SafeSignal(t)
	require.Equal(t, uint64(1), optimism.sequencer.SyncStatus().SafeL1.Number)

	// add L1 block #2
	optimism.ActL1Blocks(t, 1)
	optimism.miner.ActL1SafeNext(t)
	optimism.miner.ActL1FinalizeNext(t)
	optimism.sequencer.ActL1HeadSignal(t)
	optimism.sequencer.ActBuildToL1Head(t)

	// Catch up derivation
	optimism.sequencer.ActL2PipelineFull(t)
	optimism.sequencer.ActL1FinalizedSignal(t)
	optimism.sequencer.ActL1SafeSignal(t)

	// commit all the l2 blocks to L1
	optimism.batcher.ActSubmitAll(t)
	optimism.miner.ActL1StartBlock(12)(t)
	optimism.miner.ActL1IncludeTx(optimism.dp.Addresses.Batcher)(t)
	optimism.miner.ActL1EndBlock(t)

	// verify
	optimism.sequencer.ActL2PipelineFull(t)
	optimism.ActL1Finalized(t)

	requireDispersalRetrievalEigenDA(gt, proxyTS.Metrics.HTTPServerRequestsTotal, commitments.OptimismGeneric)
}
