package payments

import (
	"errors"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth/directory"
	"github.com/Layr-Labs/eigenda/core/payments"
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

func TestPayments(t *testing.T) {
	// manual test for now
	test.SkipInCI(t)

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
		TemplateName: "testconfig-anvil.yaml",
		TestName:     "",
		Logger:       test.GetLogger(),
		RootPath:     "../../../",
		RelayCount:   4,
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

	// - Submit blobs at a rate that is supported by the reservation, and assert that all dispersals succeed
	// - Make the reservation smaller
	// - Submit blobs at the same rate, and assert some dispersals fail
	t.Run("Reservation only with reservation reduction", func(t *testing.T) {
		t.Parallel()

		// will be billed as a minimum size blob
		payloadBytes := 1000
		// long enough to approach expected averages
		submissionDuration := 30 * time.Second
		blobsPerSecond := float32(0.5)

		paymentVault := getPaymentVault(t, testHarness, infra.Logger)
		minNumSymbols, err := paymentVault.GetMinNumSymbols(t.Context())
		require.NoError(t, err)

		// reservation required to exactly support blobsPerSecond
		reservationRequiredForRate := float32(minNumSymbols) * blobsPerSecond

		testRandom := random.NewTestRandom()
		accountID, privateKey, err := testRandom.EthAccount()
		require.NoError(t, err)
		privateKeyHex := gethcommon.Bytes2Hex(crypto.FromECDSA(privateKey))

		payloadDisperserConfig := integration.GetDefaultTestPayloadDisperserConfig()
		payloadDisperserConfig.ClientLedgerMode = clientledger.ClientLedgerModeReservationOnly
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

		payloadDisperser, err := testHarness.CreatePayloadDisperser(t.Context(), infra.Logger, payloadDisperserConfig)
		require.NoError(t, err)

		// Since we're dispersing at half the supported rate, assert no failures
		mustSubmitPayloads(t, testRandom, payloadDisperser, blobsPerSecond, payloadBytes, submissionDuration, 1.0, 0)

		clientReservation, err = reservation.NewReservation(
			// reservation smaller than it needs to be
			uint64(reservationRequiredForRate/2.0),
			time.Now().Add(-1*time.Hour),
			time.Now().Add(24*time.Hour),
			[]core.QuorumID{0, 1},
		)
		require.NoError(t, err)
		registerReservation(t, testHarness, clientReservation, accountID)

		// Since we're dispersing at double the supported rate, assert ~50% success rate
		mustSubmitPayloads(t, testRandom, payloadDisperser, blobsPerSecond, payloadBytes, submissionDuration, 0.5, 0.25)
	})

	// - Submit blobs at a rate that is larger than the reservation, and assert some dispersals fail
	// - Make the reservation larger
	// - Submit blobs at the same rate, and assert that all dispersals succeed
	t.Run("Reservation only with reservation increase", func(t *testing.T) {
		t.Parallel()

		// will be billed as a minimum size blob
		payloadBytes := 1000
		// long enough to approach expected averages
		submissionDuration := 30 * time.Second
		blobsPerSecond := float32(0.5)

		paymentVault := getPaymentVault(t, testHarness, infra.Logger)
		minNumSymbols, err := paymentVault.GetMinNumSymbols(t.Context())
		require.NoError(t, err)

		// reservation required to exactly support blobsPerSecond
		reservationRequiredForRate := float32(minNumSymbols) * blobsPerSecond

		testRandom := random.NewTestRandom()
		accountID, privateKey, err := testRandom.EthAccount()
		require.NoError(t, err)
		privateKeyHex := gethcommon.Bytes2Hex(crypto.FromECDSA(privateKey))

		payloadDisperserConfig := integration.GetDefaultTestPayloadDisperserConfig()
		payloadDisperserConfig.ClientLedgerMode = clientledger.ClientLedgerModeReservationOnly
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

		payloadDisperser, err := testHarness.CreatePayloadDisperser(t.Context(), infra.Logger, payloadDisperserConfig)
		require.NoError(t, err)

		// Since we're dispersing at double the supported rate, assert ~50% success rate
		mustSubmitPayloads(t, testRandom, payloadDisperser, blobsPerSecond, payloadBytes, submissionDuration, 0.5, 0.25)

		clientReservation, err = reservation.NewReservation(
			// reservation larger than it needs to be
			uint64(reservationRequiredForRate*2.0),
			time.Now().Add(-1*time.Hour),
			time.Now().Add(24*time.Hour),
			[]core.QuorumID{0, 1},
		)
		require.NoError(t, err)
		registerReservation(t, testHarness, clientReservation, accountID)

		// Since we're dispersing at half the supported rate, assert no failures
		mustSubmitPayloads(t, testRandom, payloadDisperser, blobsPerSecond, payloadBytes, submissionDuration, 1.0, 0)
	})

	t.Run("On-demand only", func(t *testing.T) {
		t.Parallel()

		testRandom := random.NewTestRandom()
		accountID, privateKey, err := testRandom.EthAccount()
		require.NoError(t, err)
		privateKeyHex := gethcommon.Bytes2Hex(crypto.FromECDSA(privateKey))

		paymentVault := getPaymentVault(t, testHarness, infra.Logger)
		pricePerSymbol, err := paymentVault.GetPricePerSymbol(t.Context())
		require.NoError(t, err)
		minNumSymbols, err := paymentVault.GetMinNumSymbols(t.Context())
		require.NoError(t, err)

		costPerMinSizeBlob := pricePerSymbol * uint64(minNumSymbols)
		blobsToDisperse := 5
		deposit := uint64(blobsToDisperse) * costPerMinSizeBlob

		depositOnDemand(t, testHarness, big.NewInt(int64(deposit)), accountID)

		payloadDisperserConfig := integration.GetDefaultTestPayloadDisperserConfig()
		payloadDisperserConfig.ClientLedgerMode = clientledger.ClientLedgerModeOnDemandOnly
		payloadDisperserConfig.PrivateKey = privateKeyHex

		payloadDisperser, err := testHarness.CreatePayloadDisperser(t.Context(), infra.Logger, payloadDisperserConfig)
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
	})

	t.Run("Reservation and on-demand", func(t *testing.T) {
		t.Parallel()

		testRandom := random.NewTestRandom()
		accountID, privateKey, err := testRandom.EthAccount()
		require.NoError(t, err)
		privateKeyHex := gethcommon.Bytes2Hex(crypto.FromECDSA(privateKey))

		paymentVault := getPaymentVault(t, testHarness, infra.Logger)
		pricePerSymbol, err := paymentVault.GetPricePerSymbol(t.Context())
		require.NoError(t, err)
		minNumSymbols, err := paymentVault.GetMinNumSymbols(t.Context())
		require.NoError(t, err)

		payloadBytes := 1000
		submissionDuration := 60 * time.Second
		blobsPerSecond := float32(0.5)

		// this is the total amount of billable symbols that are being dispersed
		billableSymbolsPerSecond := uint64(blobsPerSecond * float32(minNumSymbols))

		// Reservation covers 25% of the dispersal rate
		clientReservation, err := reservation.NewReservation(
			billableSymbolsPerSecond/4,
			time.Now().Add(-1*time.Hour),
			time.Now().Add(24*time.Hour),
			[]core.QuorumID{0, 1},
		)
		require.NoError(t, err)
		registerReservation(t, testHarness, clientReservation, accountID)

		// deposit enough on-demand funds to cover one entire dispersal duration
		onDemandDepositSymbols := billableSymbolsPerSecond * uint64(submissionDuration.Seconds())
		onDemandDeposit := big.NewInt(int64(onDemandDepositSymbols * pricePerSymbol))
		depositOnDemand(t, testHarness, onDemandDeposit, accountID)

		payloadDisperserConfig := integration.GetDefaultTestPayloadDisperserConfig()
		payloadDisperserConfig.ClientLedgerMode = clientledger.ClientLedgerModeReservationAndOnDemand
		payloadDisperserConfig.PrivateKey = privateKeyHex

		payloadDisperser, err := testHarness.CreatePayloadDisperser(t.Context(), infra.Logger, payloadDisperserConfig)
		require.NoError(t, err)

		// Phase 1: Since the reservation covers 25% of the dispersal rate, this is expected to use up 75% of the
		// deposited on-demand funds, but there shouldn't be any failures.
		mustSubmitPayloads(t, testRandom, payloadDisperser, blobsPerSecond, payloadBytes, submissionDuration, 1.0, 0)

		// Phase 2: 25% of the dispersals within this period are covered by the reservation. 25% are covered by
		// remaining on-demand funds. So expected failure rate is 50%
		mustSubmitPayloads(t, testRandom, payloadDisperser, blobsPerSecond, payloadBytes, submissionDuration, 0.5, 0.25)

		// Phase 3: This phase disperses at half the rate of the previous phases. Even with the decreased rate, only
		// half of dispersals are covered by the reservation. There are no on-demand funds remaining, so failure rate
		// should be 50%
		mustSubmitPayloads(t, testRandom, payloadDisperser, blobsPerSecond/2, payloadBytes, submissionDuration, 0.5, 0.25)
	})

	t.Run("Reservation only with reservation expiration", func(t *testing.T) {
		t.Parallel()
		testReservationExpiration(t, infra.Logger, testHarness, clientledger.ClientLedgerModeReservationOnly)
	})

	t.Run("Reservation and on-demand with reservation expiration", func(t *testing.T) {
		t.Parallel()
		testReservationExpiration(t, infra.Logger, testHarness, clientledger.ClientLedgerModeReservationAndOnDemand)
	})
}

// - Create a reservation that expires soon
// - Submit a blob and assert success
// - Sleep until reservation expires
// - Assert next blob submission fails appropriately based on client ledger mode
// - Register a new valid reservation
// - Create a new payload disperser
// - Submit a blob and assert success
func testReservationExpiration(
	t *testing.T,
	logger logging.Logger,
	testHarness *integration.TestHarness,
	clientLedgerMode clientledger.ClientLedgerMode,
) {
	payloadBytes := 1000
	// the reservation will be configured to expire shortly after the first dispersal
	reservationExpirationDelay := 20 * time.Second

	paymentVault := getPaymentVault(t, testHarness, logger)
	minNumSymbols, err := paymentVault.GetMinNumSymbols(t.Context())
	require.NoError(t, err)

	testRandom := random.NewTestRandom()
	accountID, privateKey, err := testRandom.EthAccount()
	require.NoError(t, err)
	privateKeyHex := gethcommon.Bytes2Hex(crypto.FromECDSA(privateKey))

	payloadDisperserConfig := integration.GetDefaultTestPayloadDisperserConfig()
	payloadDisperserConfig.ClientLedgerMode = clientLedgerMode
	payloadDisperserConfig.PrivateKey = privateKeyHex

	clientReservation, err := reservation.NewReservation(
		uint64(minNumSymbols)*100,
		time.Now().Add(-1*time.Hour),
		// expires soon
		time.Now().Add(reservationExpirationDelay),
		[]core.QuorumID{0, 1},
	)
	require.NoError(t, err)
	registerReservation(t, testHarness, clientReservation, accountID)

	payloadDisperser, err := testHarness.CreatePayloadDisperser(t.Context(), logger, payloadDisperserConfig)
	require.NoError(t, err)

	// Blob should succeed while reservation is active
	payload := coretypes.Payload(testRandom.Bytes(payloadBytes))
	_, err = payloadDisperser.SendPayload(t.Context(), payload)
	require.NoError(t, err)

	// Wait for reservation to expire
	time.Sleep(reservationExpirationDelay)

	payload = coretypes.Payload(testRandom.Bytes(payloadBytes))

	// Behavior differs based on client ledger mode:
	// - ReservationOnly: returns TimeOutOfRangeError
	// - ReservationAndOnDemand: panics to avoid inadvertently depleting on-demand funds
	switch clientLedgerMode {
	case clientledger.ClientLedgerModeReservationOnly:
		_, err = payloadDisperser.SendPayload(t.Context(), payload)
		require.Error(t, err, "dispersal should fail with expired reservation")
		var timeOutOfRangeError *reservation.TimeOutOfRangeError
		require.True(t, errors.As(err, &timeOutOfRangeError), "error should be TimeOutOfRangeError")
	case clientledger.ClientLedgerModeReservationAndOnDemand:
		require.Panics(t, func() {
			_, _ = payloadDisperser.SendPayload(t.Context(), payload)
		}, "dispersal should panic with expired reservation in ReservationAndOnDemand mode")
	case clientledger.ClientLedgerModeOnDemandOnly:
		panic("testReservationExpiration should not be called with OnDemandOnly")
	default:
		panic("testReservationExpiration called with unexpected client ledger mode")
	}

	// Register a new valid reservation
	clientReservation, err = reservation.NewReservation(
		uint64(minNumSymbols)*100,
		time.Now().Add(-reservationExpirationDelay),
		time.Now().Add(24*time.Hour),
		[]core.QuorumID{0, 1},
	)
	require.NoError(t, err)
	registerReservation(t, testHarness, clientReservation, accountID)

	payloadDisperser, err = testHarness.CreatePayloadDisperser(t.Context(), logger, payloadDisperserConfig)
	require.NoError(t, err)

	// Blob should succeed with the new valid reservation
	payload = coretypes.Payload(testRandom.Bytes(payloadBytes))
	_, err = payloadDisperser.SendPayload(t.Context(), payload)
	require.NoError(t, err)
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

func getPaymentVault(t *testing.T, testHarness *integration.TestHarness, logger logging.Logger) payments.PaymentVault {
	paymentVaultAddress, err := testHarness.ContractDirectory.GetContractAddress(t.Context(), directory.PaymentVault)
	require.NoError(t, err)
	paymentVault, err := vault.NewPaymentVault(logger, testHarness.EthClient, paymentVaultAddress)
	require.NoError(t, err)

	return paymentVault
}
