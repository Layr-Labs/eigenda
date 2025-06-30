package clients

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

type Accountant struct {
	// account identification
	accountID gethcommon.Address

	// payment system configuration
	paymentVaultParams *meterer.PaymentVaultParams

	// unified ledger managing all account state (on-chain settings + local tracking)
	accountLedger meterer.AccountLedger

	// lock for concurrent access to account state
	accountLock sync.Mutex
}

// NewAccountant initializes an accountant with the given account ID. The accountant must call SetPaymentState to populate the state.
func NewAccountant(accountID gethcommon.Address) *Accountant {
	return &Accountant{
		accountID: accountID,
		paymentVaultParams: &meterer.PaymentVaultParams{
			QuorumPaymentConfigs:  make(map[uint8]*core.PaymentQuorumConfig),
			QuorumProtocolConfigs: make(map[uint8]*core.PaymentQuorumProtocolConfig),
			OnDemandQuorumNumbers: make([]uint8, 0),
		},
		accountLedger: meterer.NewLocalAccountLedger(),
	}
}

// ReservationUsage attempts to use the reservation for the requested quorums; if any quorum fails to use the reservation, the entire operation is rolled back.
func (a *Accountant) reservationUsage(numSymbols uint64, quorumNumbers []core.QuorumID, paymentHeaderTimestampNs int64) error {
	a.accountLock.Lock()
	defer a.accountLock.Unlock()

	return a.accountLedger.RecordReservationUsage(
		context.Background(),
		a.accountID,
		paymentHeaderTimestampNs,
		numSymbols,
		quorumNumbers,
		a.paymentVaultParams,
	)
}

// onDemandUsage attempts to use on-demand payment for the given request.
// Returns the cumulative payment if successful, or an error if on-demand cannot be used.
func (a *Accountant) onDemandUsage(numSymbols uint64, quorumNumbers []core.QuorumID) (*big.Int, error) {
	a.accountLock.Lock()
	defer a.accountLock.Unlock()

	return a.accountLedger.RecordOnDemandUsage(
		context.Background(),
		a.accountID,
		numSymbols,
		quorumNumbers,
		a.paymentVaultParams,
	)
}

// AccountBlob accountant generates payment information for a request. The accountant
// takes in a timestamp at the current UNIX time in nanoseconds, number of symbols of the request,
// and the quorums to disperse the request to. It will attempt to use the active reservation first
// and then on-demand if the reservation is not available or insufficient for the request.
// It returns a payment metadata object that will be used to create the payment header and signature,
// as specified in api/proto/common/v2/common_v2.proto
func (a *Accountant) AccountBlob(
	timestamp int64,
	numSymbols uint64,
	quorums []uint8) (*core.PaymentMetadata, error) {
	if len(quorums) == 0 {
		return nil, fmt.Errorf("no quorums provided")
	}
	if numSymbols == 0 {
		return nil, fmt.Errorf("zero symbols requested")
	}

	// Always try to use reservation first
	rezErr := a.reservationUsage(numSymbols, quorums, timestamp)
	if rezErr == nil {
		return &core.PaymentMetadata{
			AccountID:         a.accountID,
			Timestamp:         timestamp,
			CumulativePayment: big.NewInt(0),
		}, nil
	}

	// Fall back to on-demand payment if reservation fails
	cumulativePayment, onDemandErr := a.onDemandUsage(numSymbols, quorums)
	if onDemandErr != nil {
		return nil, fmt.Errorf("cannot create payment information for reservation or on-demand. Consider depositing more eth to the PaymentVault contract for your account. For more details, see https://docs.eigenda.xyz/core-concepts/payments#disperser-client-requirements. Account: %s, Reservation Error: %w, On-demand Error: %w", a.accountID.Hex(), rezErr, onDemandErr)
	}

	pm := &core.PaymentMetadata{
		AccountID:         a.accountID,
		Timestamp:         timestamp,
		CumulativePayment: cumulativePayment,
	}

	return pm, nil
}

// SetPaymentState sets the accountant's state, requiring valid payment vault parameters, but
// optional account level on/off-chain state. If on-chain fields are not present, we use dummy
// values that disable accountant from using the corresponding payment method. If off-chain
// fields are not present, we assume the account has no payment history and set accountant state
// to use initial values.
func (a *Accountant) SetPaymentState(
	paymentVaultParams *meterer.PaymentVaultParams,
	reservations map[core.QuorumID]*core.ReservedPayment,
	cumulativePayment *big.Int,
	onchainCumulativePayment *big.Int,
	periodRecords meterer.QuorumPeriodRecords,
) error {
	if paymentVaultParams == nil {
		return fmt.Errorf("payment vault params cannot be nil")
	}

	a.paymentVaultParams = paymentVaultParams

	// Create on-demand payment state
	var onDemand *core.OnDemandPayment
	if onchainCumulativePayment == nil {
		onDemand = &core.OnDemandPayment{
			CumulativePayment: big.NewInt(0),
		}
	} else {
		onDemand = &core.OnDemandPayment{
			CumulativePayment: new(big.Int).Set(onchainCumulativePayment),
		}
	}

	// Set the complete account ledger state including on-chain settings and local tracking
	accountState := meterer.AccountState{
		Reservations:      reservations,
		OnDemand:          onDemand,
		PeriodRecords:     periodRecords,
		CumulativePayment: cumulativePayment,
	}
	a.accountLedger.SetAccountState(accountState)

	return nil
}

// GetPeriodRecords returns a copy of the current period records for testing purposes
func (a *Accountant) GetPeriodRecords() meterer.QuorumPeriodRecords {
	return a.accountLedger.GetAccountState().PeriodRecords
}
