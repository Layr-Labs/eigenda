package e2e

import (
	"testing"

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/common/types/commitments"
	"github.com/Layr-Labs/eigenda-proxy/testutils"
	altda "github.com/ethereum-optimism/optimism/op-alt-da"
	e2econfig "github.com/ethereum-optimism/optimism/op-e2e/config"
	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils"
	"github.com/ethereum-optimism/optimism/op-service/sources"
	"github.com/ethereum-optimism/optimism/op-service/testlog"
	gethcommon "github.com/ethereum/go-ethereum/common"
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
		AllocType:           e2econfig.AllocTypeAltDA,
	}

	log := testlog.Logger(t, log.LevelWarn)

	// config.DeployConfig.DACommitmentType = altda.GenericCommitmentString
	dp := e2eutils.MakeDeployParams(t, p)
	dp.DeployConfig.DAChallengeProxy = gethcommon.Address{0x42}
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

	l1F, err := sources.NewL1Client(
		miner.RPCClient(),
		log,
		nil,
		sources.L1ClientDefaultConfig(sd.RollupCfg, false, sources.RPCKindBasic))
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
	miner.ActL1SetFeeRecipient(gethcommon.Address{'A'})
	sequencer.ActL2PipelineFull(t)

	batcher := actions.NewL2Batcher(
		log,
		sd.RollupCfg,
		actions.AltDABatcherCfg(dp, storage),
		sequencer.RollupClient(),
		l1Client,
		engine.EthClient(),
		engCl)

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

func TestOptimismKeccak256CommitmentV1(t *testing.T) {
	testOptimismKeccak256Commitment(t, common.V1EigenDABackend)
}

func TestOptimismKeccak256CommitmentV2(t *testing.T) {
	testOptimismKeccak256Commitment(t, common.V2EigenDABackend)
}

func testOptimismKeccak256Commitment(t *testing.T, dispersalBackend common.EigenDABackend) {
	testCfg := testutils.NewTestConfig(testutils.GetBackend(), dispersalBackend, nil)
	testCfg.UseKeccak256ModeS3 = true

	tsConfig := testutils.BuildTestSuiteConfig(testCfg)
	proxyTS, shutDown := testutils.CreateTestSuite(tsConfig)
	defer shutDown()

	ot := actions.NewDefaultTesting(t)

	optimism := NewL2AltDA(ot, proxyTS.Address(), false)

	// build L1 block #1
	optimism.ActL1Blocks(ot, 1)
	optimism.miner.ActL1SafeNext(ot)

	// Fill with l2 blocks up to the L1 head
	optimism.sequencer.ActL1HeadSignal(ot)
	optimism.sequencer.ActBuildToL1Head(ot)

	optimism.sequencer.ActL2PipelineFull(ot)
	optimism.sequencer.ActL1SafeSignal(ot)
	require.Equal(ot, uint64(1), optimism.sequencer.SyncStatus().SafeL1.Number)

	// add L1 block #2
	optimism.ActL1Blocks(ot, 1)
	optimism.miner.ActL1SafeNext(ot)
	optimism.miner.ActL1FinalizeNext(ot)
	optimism.sequencer.ActL1HeadSignal(ot)
	optimism.sequencer.ActBuildToL1Head(ot)

	// Catch up derivation
	optimism.sequencer.ActL2PipelineFull(ot)
	optimism.sequencer.ActL1FinalizedSignal(ot)
	optimism.sequencer.ActL1SafeSignal(ot)

	// commit all the l2 blocks to L1
	optimism.batcher.ActSubmitAll(ot)
	optimism.miner.ActL1StartBlock(12)(ot)
	optimism.miner.ActL1IncludeTx(optimism.dp.Addresses.Batcher)(ot)
	optimism.miner.ActL1EndBlock(ot)

	// verify
	optimism.sequencer.ActL2PipelineFull(ot)
	optimism.ActL1Finalized(ot)

	requireDispersalRetrievalEigenDA(
		t,
		proxyTS.Metrics.HTTPServerRequestsTotal,
		commitments.OptimismKeccakCommitmentMode)
}

func TestOptimismGenericCommitmentV1(t *testing.T) {
	testOptimismGenericCommitment(t, common.V1EigenDABackend)
}

func TestOptimismGenericCommitmentV2(t *testing.T) {
	testOptimismGenericCommitment(t, common.V2EigenDABackend)
}

func testOptimismGenericCommitment(t *testing.T, dispersalBackend common.EigenDABackend) {
	testCfg := testutils.NewTestConfig(testutils.GetBackend(), dispersalBackend, nil)

	tsConfig := testutils.BuildTestSuiteConfig(testCfg)
	proxyTS, shutDown := testutils.CreateTestSuite(tsConfig)
	defer shutDown()

	ot := actions.NewDefaultTesting(t)

	optimism := NewL2AltDA(ot, proxyTS.Address(), true)
	exerciseGenericCommitments(t, ot, optimism)

	requireDispersalRetrievalEigenDA(
		t,
		proxyTS.Metrics.HTTPServerRequestsTotal,
		commitments.OptimismGenericCommitmentMode)
}

func exerciseGenericCommitments(
	t *testing.T,
	ot actions.StatefulTesting,
	optimism *L2AltDA,
) {
	expectedBlockNumber := optimism.sequencer.SyncStatus().SafeL1.Number

	// build L1 block #1
	optimism.ActL1Blocks(ot, 1)
	optimism.miner.ActL1SafeNext(ot)

	// Fill with l2 blocks up to the L1 head
	optimism.sequencer.ActL1HeadSignal(ot)
	optimism.sequencer.ActBuildToL1Head(ot)

	optimism.sequencer.ActL2PipelineFull(ot)
	optimism.sequencer.ActL1SafeSignal(ot)

	expectedBlockNumber++
	require.Equal(t, expectedBlockNumber, optimism.sequencer.SyncStatus().SafeL1.Number)

	// add L1 block #2
	optimism.ActL1Blocks(ot, 1)
	optimism.miner.ActL1SafeNext(ot)
	optimism.miner.ActL1FinalizeNext(ot)
	optimism.sequencer.ActL1HeadSignal(ot)
	optimism.sequencer.ActBuildToL1Head(ot)

	// Catch up derivation
	optimism.sequencer.ActL2PipelineFull(ot)
	optimism.sequencer.ActL1FinalizedSignal(ot)
	optimism.sequencer.ActL1SafeSignal(ot)

	expectedBlockNumber++
	require.Equal(t, expectedBlockNumber, optimism.sequencer.SyncStatus().SafeL1.Number)

	// commit all the l2 blocks to L1
	optimism.batcher.ActSubmitAll(ot)
	optimism.miner.ActL1StartBlock(12)(ot)
	optimism.miner.ActL1IncludeTx(optimism.dp.Addresses.Batcher)(ot)
	optimism.miner.ActL1EndBlock(ot)

	// verify
	optimism.sequencer.ActL2PipelineFull(ot)
	optimism.ActL1Finalized(ot)
}

func TestOptimismGenericCommitmentMigration(t *testing.T) {
	testCfg := testutils.NewTestConfig(
		testutils.GetBackend(),
		common.V1EigenDABackend,
		[]common.EigenDABackend{common.V1EigenDABackend, common.V2EigenDABackend})
	tsConfig := testutils.BuildTestSuiteConfig(testCfg)
	proxyTS, shutDown := testutils.CreateTestSuite(tsConfig)
	defer shutDown()

	expectedWriteCount := uint64(0)
	expectedReadCount := uint64(0)

	ot := actions.NewDefaultTesting(t)

	optimism := NewL2AltDA(ot, proxyTS.Address(), true)
	exerciseGenericCommitments(t, ot, optimism)
	expectedWriteCount++
	expectedReadCount++
	requireDispersalRetrievalEigenDACounts(
		t,
		proxyTS.Metrics.HTTPServerRequestsTotal,
		commitments.OptimismGenericCommitmentMode,
		expectedWriteCount,
		expectedReadCount)

	// turn on v2 dispersal
	proxyTS.Server.SetDispersalBackend(common.V2EigenDABackend)
	exerciseGenericCommitments(t, ot, optimism)
	expectedWriteCount++
	expectedReadCount++
	requireDispersalRetrievalEigenDACounts(
		t,
		proxyTS.Metrics.HTTPServerRequestsTotal,
		commitments.OptimismGenericCommitmentMode,
		expectedWriteCount,
		expectedReadCount)
}
