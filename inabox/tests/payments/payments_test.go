package payments

import (
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth/directory"
	"github.com/Layr-Labs/eigenda/core/payments/clientledger"
	"github.com/Layr-Labs/eigenda/core/payments/reservation"
	"github.com/Layr-Labs/eigenda/core/payments/vault"
	integration "github.com/Layr-Labs/eigenda/inabox/tests"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/Layr-Labs/eigenda/test/random"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

// NOTE: Currently, it doesn't work to run these tests in sequence. Each test must be run as a separate command.
// The problem is that the cleanup logic sometimes randomly fails to free docker ports, so subsequent setups fail.
// Once we figure out why resources aren't being freed, then these tests will be runnable the "normal" way.

func TestLegacyController(t *testing.T) {
	t.Skip("Manual test for now")
	testWithControllerMode(t, false)
}

func TestNewController(t *testing.T) {
	t.Skip("Manual test for now")
	testWithControllerMode(t, true)
}

func testWithControllerMode(t *testing.T, controllerUseNewPayments bool) {
	// Save current working directory. The setup process in its current form changes working directory, which causes
	// subsequent executions to fail, since the process relies on relative paths. This is a workaround for now: we just
	// capture the original working directory, and switch back to it as a cleanup step.
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Logf("Failed to restore working directory: %v", err)
		}
	})

	infraConfig := &integration.InfrastructureConfig{
		TemplateName:             "testconfig-anvil.yaml",
		TestName:                 "",
		Logger:                   test.GetLogger(),
		RootPath:                 "../../../",
		RelayCount:               4,
		ControllerUseNewPayments: controllerUseNewPayments,
	}

	infra, err := integration.SetupInfrastructure(t.Context(), infraConfig)
	if infra != nil {
		t.Cleanup(func() {
			integration.TeardownInfrastructure(infra)
		})
	}
	require.NoError(t, err)

	testHarness, err := integration.NewTestHarnessWithSetup(infra)
	if testHarness != nil {
		t.Cleanup(func() {
			testHarness.Cleanup()
		})
	}
	require.NoError(t, err)

	// Subtests all use unique accountIDs, so they can run in parallel

	t.Run("Reservation only: Legacy client payments with reservation reduction", func(t *testing.T) {
		t.Parallel()
		testReservationReduction(t, infra.Logger, testHarness, clientledger.ClientLedgerModeLegacy)
	})

	t.Run("Reservation only: New client payments with reservation reduction", func(t *testing.T) {
		t.Parallel()
		testReservationReduction(t, infra.Logger, testHarness, clientledger.ClientLedgerModeReservationOnly)
	})

	t.Run("Reservation only: Legacy client payments with reservation increase", func(t *testing.T) {
		t.Parallel()
		testReservationIncrease(t, infra.Logger, testHarness, clientledger.ClientLedgerModeLegacy)
	})

	t.Run("Reservation only: New client payments with reservation increase", func(t *testing.T) {
		t.Parallel()
		testReservationIncrease(t, infra.Logger, testHarness, clientledger.ClientLedgerModeReservationOnly)
	})

	t.Run("On-demand only: Legacy client payments", func(t *testing.T) {
		t.Parallel()
		testOnDemandOnly(t, infra.Logger, testHarness, clientledger.ClientLedgerModeLegacy)
	})

	t.Run("On-demand only: New client payments", func(t *testing.T) {
		t.Parallel()
		testOnDemandOnly(t, infra.Logger, testHarness, clientledger.ClientLedgerModeOnDemandOnly)
	})
}

// - Submit blobs at a rate that is supported by the reservation, and assert that all dispersals succeed
// - Make the reservation smaller
// - Submit blobs at the same rate, and assert some dispersals fail
func testReservationReduction(
	t *testing.T,
	logger logging.Logger,
	testHarness *integration.TestHarness,
	clientLedgerMode clientledger.ClientLedgerMode,
) {
	// will be billed as a minimum size blob
	payloadBytes := 1000
	// long enough to approach expected averages
	submissionDuration := 30 * time.Second
	blobsPerSecond := float32(0.5)
	// how large a reservation (in symbols / second) is required to submit 1 minimum size blob / second
	// min billable blob size = 128KiB = 4096 symbols
	minSizeBlobPerSecondReservationSize := 4096
	// reservation required to exactly support blobsPerSecond
	reservationRequiredForRate := float32(minSizeBlobPerSecondReservationSize) * blobsPerSecond

	testRandom := random.NewTestRandom()
	publicKey, privateKey, err := testRandom.ECDSA()
	require.NoError(t, err)
	privateKeyHex := gethcommon.Bytes2Hex(crypto.FromECDSA(privateKey))
	accountID := crypto.PubkeyToAddress(*publicKey)

	payloadDisperserConfig := integration.GetDefaultTestPayloadDisperserConfig()
	payloadDisperserConfig.ClientLedgerMode = clientLedgerMode
	payloadDisperserConfig.PrivateKey = privateKeyHex

	clientReservation, err := reservation.NewReservation(
		// reservation larger than it needs to be
		uint64(reservationRequiredForRate*2.0),
		time.Now().Add(-1*time.Hour),
		time.Now().Add(24*time.Hour),
		[]core.QuorumID{0, 1},
	)
	require.NoError(t, err)
	registerReservation(t, testHarness, clientReservation, accountID)

	payloadDisperser, err := testHarness.CreatePayloadDisperser(t.Context(), logger, payloadDisperserConfig)
	require.NoError(t, err)

	// Since we're dispersing at half the supported rate, assert no failures
	resultChan := mustSubmitPayloads(
		t, testRandom, payloadDisperser, blobsPerSecond, payloadBytes, submissionDuration, 1.0, 0)
	// Drain the results channel. This test doesn't need the values.
	for range resultChan {
	}

	clientReservation, err = reservation.NewReservation(
		// reservation smaller than it needs to be
		uint64(reservationRequiredForRate/2.0),
		time.Now().Add(-1*time.Hour),
		time.Now().Add(24*time.Hour),
		[]core.QuorumID{0, 1},
	)
	require.NoError(t, err)
	registerReservation(t, testHarness, clientReservation, accountID)

	// The next part of the test asserts the same dispersal conditions now yield errors, considering the decreased
	// reservation. The legacy client ledger mode doesn't observe payment vault updates, so we need to build
	// a new payload disperser to pick up the change
	if clientLedgerMode == clientledger.ClientLedgerModeLegacy {
		payloadDisperser, err = testHarness.CreatePayloadDisperser(t.Context(), logger, payloadDisperserConfig)
		require.NoError(t, err)
	}

	// Since we're dispersing at double the supported rate, assert ~50% success rate
	resultChan = mustSubmitPayloads(
		t, testRandom, payloadDisperser, blobsPerSecond, payloadBytes, submissionDuration, 0.5, 0.25)
	for range resultChan {
	}
}

// - Submit blobs at a rate that is larger than the reservation, and assert some dispersals fail
// - Make the reservation larger
// - Submit blobs at the same rate, and assert that all dispersals succeed
func testReservationIncrease(
	t *testing.T,
	logger logging.Logger,
	testHarness *integration.TestHarness,
	clientLedgerMode clientledger.ClientLedgerMode,
) {
	// will be billed as a minimum size blob
	payloadBytes := 1000
	// long enough to approach expected averages
	submissionDuration := 30 * time.Second
	blobsPerSecond := float32(0.5)
	// how large a reservation (in symbols / second) is required to submit 1 minimum size blob / second
	// min billable blob size = 128KiB = 4096 symbols
	minSizeBlobPerSecondReservationSize := 4096
	// reservation required to exactly support blobsPerSecond
	reservationRequiredForRate := float32(minSizeBlobPerSecondReservationSize) * blobsPerSecond

	testRandom := random.NewTestRandom()
	publicKey, privateKey, err := testRandom.ECDSA()
	require.NoError(t, err)
	privateKeyHex := gethcommon.Bytes2Hex(crypto.FromECDSA(privateKey))
	accountID := crypto.PubkeyToAddress(*publicKey)

	payloadDisperserConfig := integration.GetDefaultTestPayloadDisperserConfig()
	payloadDisperserConfig.ClientLedgerMode = clientLedgerMode
	payloadDisperserConfig.PrivateKey = privateKeyHex

	clientReservation, err := reservation.NewReservation(
		// reservation smaller than it needs to be
		uint64(reservationRequiredForRate/2.0),
		time.Now().Add(-1*time.Hour),
		time.Now().Add(24*time.Hour),
		[]core.QuorumID{0, 1},
	)
	require.NoError(t, err)
	registerReservation(t, testHarness, clientReservation, accountID)

	payloadDisperser, err := testHarness.CreatePayloadDisperser(t.Context(), logger, payloadDisperserConfig)
	require.NoError(t, err)

	// Since we're dispersing at double the supported rate, assert ~50% success rate
	resultChan := mustSubmitPayloads(
		t, testRandom, payloadDisperser, blobsPerSecond, payloadBytes, submissionDuration, 0.5, 0.25)
	// Drain the results channel. This test doesn't need the values.
	for range resultChan {
	}

	clientReservation, err = reservation.NewReservation(
		// reservation larger than it needs to be
		uint64(reservationRequiredForRate*2.0),
		time.Now().Add(-1*time.Hour),
		time.Now().Add(24*time.Hour),
		[]core.QuorumID{0, 1},
	)
	require.NoError(t, err)
	registerReservation(t, testHarness, clientReservation, accountID)

	// The next part of the test asserts the same dispersal conditions no longer yield errors, considering the increased
	// reservation. The legacy client ledger mode doesn't observe payment vault updates, so we need to build
	// a new payload disperser to pick up the change
	if clientLedgerMode == clientledger.ClientLedgerModeLegacy {
		payloadDisperser, err = testHarness.CreatePayloadDisperser(t.Context(), logger, payloadDisperserConfig)
		require.NoError(t, err)
	}

	// Since we're dispersing at half the supported rate, assert no failures
	resultChan = mustSubmitPayloads(
		t, testRandom, payloadDisperser, blobsPerSecond, payloadBytes, submissionDuration, 1.0, 0)
	for range resultChan {
	}
}

func testOnDemandOnly(
	t *testing.T,
	logger logging.Logger,
	testHarness *integration.TestHarness,
	clientLedgerMode clientledger.ClientLedgerMode,
) {
	testRandom := random.NewTestRandom()
	publicKey, privateKey, err := testRandom.ECDSA()
	require.NoError(t, err)
	privateKeyHex := gethcommon.Bytes2Hex(crypto.FromECDSA(privateKey))
	accountID := crypto.PubkeyToAddress(*publicKey)

	paymentVaultAddress, err := testHarness.ContractDirectory.GetContractAddress(t.Context(), directory.PaymentVault)
	require.NoError(t, err)
	paymentVault, err := vault.NewPaymentVault(logger, testHarness.EthClient, paymentVaultAddress)
	require.NoError(t, err)
	pricePerSymbol, err := paymentVault.GetPricePerSymbol(t.Context())
	require.NoError(t, err)
	minNumSymbols, err := paymentVault.GetMinNumSymbols(t.Context())
	require.NoError(t, err)

	costPerMinSizeBlob := pricePerSymbol * uint64(minNumSymbols)
	blobsToDisperse := 5
	deposit := uint64(blobsToDisperse) * costPerMinSizeBlob

	depositOnDemand(t, testHarness, big.NewInt(int64(deposit)), accountID)

	payloadDisperserConfig := integration.GetDefaultTestPayloadDisperserConfig()
	payloadDisperserConfig.ClientLedgerMode = clientLedgerMode
	payloadDisperserConfig.PrivateKey = privateKeyHex

	payloadDisperser, err := testHarness.CreatePayloadDisperser(t.Context(), logger, payloadDisperserConfig)
	require.NoError(t, err)

	// will be billed as a minimum size blob
	payloadBytes := 1000

	// disperse the number of blobs that we expect to succeed
	for i := 0; i < blobsToDisperse; i++ {
		payload := coretypes.Payload(testRandom.Bytes(payloadBytes))
		_, err := payloadDisperser.SendPayload(t.Context(), payload)
		require.NoError(t, err)
	}

	// the very next dispersal should fail
	payload := coretypes.Payload(testRandom.Bytes(payloadBytes))
	_, err = payloadDisperser.SendPayload(t.Context(), payload)
	require.Error(t, err)

	depositOnDemand(t, testHarness, big.NewInt(int64(deposit)), accountID)

	// The next part of the test makes additional dispersals, considering the new deposit.
	// The legacy client ledger mode doesn't observe payment vault updates, so we need to rebuild it
	if clientLedgerMode == clientledger.ClientLedgerModeLegacy {
		payloadDisperser, err = testHarness.CreatePayloadDisperser(t.Context(), logger, payloadDisperserConfig)
		require.NoError(t, err)
	}

	// disperse the number of blobs that we expect to succeed
	for i := 0; i < blobsToDisperse; i++ {
		payload := coretypes.Payload(testRandom.Bytes(payloadBytes))
		_, err := payloadDisperser.SendPayload(t.Context(), payload)
		require.NoError(t, err)
	}

	// the very next dispersal should fail
	payload = coretypes.Payload(testRandom.Bytes(payloadBytes))
	_, err = payloadDisperser.SendPayload(t.Context(), payload)
	require.Error(t, err)
}

// Registers a reservation on-chain, then sleeps for a short time to wait for the updated value to be picked up by
// payment vault monitors
func registerReservation(
	t *testing.T,
	testHarness *integration.TestHarness,
	newReservation *reservation.Reservation,
	accountID gethcommon.Address,
) {
	err := testHarness.UpdateReservationOnChain(t, accountID, newReservation)
	require.NoError(t, err)
	// the vault monitor checks every 1 second, so this should be plenty of time
	time.Sleep(3 * time.Second)
}

// Makes an on-demand deposit for an account and waits for the vault monitor to pick it up
func depositOnDemand(
	t *testing.T,
	testHarness *integration.TestHarness,
	depositAmount *big.Int,
	accountID gethcommon.Address,
) {
	err := testHarness.DepositOnDemandOnChain(t, accountID, depositAmount)
	require.NoError(t, err)
	// the vault monitor checks every 1 second, so this should be plenty of time
	time.Sleep(3 * time.Second)
}
