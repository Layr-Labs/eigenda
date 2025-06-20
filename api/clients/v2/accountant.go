package clients

import (
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

type Accountant struct {
	// on-chain states
	accountID    gethcommon.Address
	reservations map[core.QuorumID]*core.ReservedPayment
	// OnDemand is initially only enabled on quorum 0. Accountant must be updated to be quorum specific
	// after the protocol decides to support onDemand on custom quorums and decentralized ratelimiting.
	onDemand *core.OnDemandPayment

	paymentVaultParams *meterer.PaymentVaultParams

	// local accounting
	// contains a fixed meterer.MinNumBins bins per quorum
	periodRecords     meterer.QuorumPeriodRecords
	cumulativePayment *big.Int

	// locks for concurrent access to period records
	periodRecordsLock sync.Mutex
	// lock for concurrent access to on-demand payment
	onDemandLock sync.Mutex
}

// NewAccountant initializes an accountant with the given account ID. The accountant must call SetPaymentState to populate the state.
func NewAccountant(accountID gethcommon.Address) *Accountant {
	reservations := make(map[core.QuorumID]*core.ReservedPayment)
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(0),
	}
	periodRecords := make(meterer.QuorumPeriodRecords)
	return &Accountant{
		accountID:    accountID,
		reservations: reservations,
		onDemand:     onDemand,
		paymentVaultParams: &meterer.PaymentVaultParams{
			QuorumPaymentConfigs:  make(map[uint8]*core.PaymentQuorumConfig),
			QuorumProtocolConfigs: make(map[uint8]*core.PaymentQuorumProtocolConfig),
			OnDemandQuorumNumbers: make([]uint8, 0),
		},
		periodRecords:     periodRecords,
		cumulativePayment: big.NewInt(0),
	}
}

// ReservationUsage attempts to use the reservation for the requested quorums; if any quorum fails to use the reservation, the entire operation is rolled back.
func (a *Accountant) reservationUsage(numSymbols uint64, quorumNumbers []core.QuorumID, timestampNs int64) error {
	if err := meterer.ValidateReservations(a.reservations, a.paymentVaultParams.QuorumProtocolConfigs, quorumNumbers, timestampNs, time.Now()); err != nil {
		return err
	}

	a.periodRecordsLock.Lock()
	defer a.periodRecordsLock.Unlock()

	periodRecordsCopy := a.periodRecords.DeepCopy()

	for _, quorumNumber := range quorumNumbers {
		reservation, exists := a.reservations[quorumNumber]
		if !exists {
			// this case should never happen because ValidateReservations should have already checked this; handle it just in case
			return fmt.Errorf("reservation not found for quorum %d", quorumNumber)
		}
		_, protocolConfig, err := a.paymentVaultParams.GetQuorumConfigs(quorumNumber)
		if err != nil {
			return err
		}
		if err := periodRecordsCopy.UpdateUsage(quorumNumber, timestampNs, numSymbols, reservation, protocolConfig); err != nil {
			return err
		}
	}

	a.periodRecords = periodRecordsCopy
	return nil
}

// onDemandUsage attempts to use on-demand payment for the given request.
// Returns the cumulative payment if successful, or an error if on-demand cannot be used.
func (a *Accountant) onDemandUsage(numSymbols uint64, quorumNumbers []core.QuorumID) (*big.Int, error) {
	if err := meterer.ValidateQuorum(quorumNumbers, a.paymentVaultParams.OnDemandQuorumNumbers); err != nil {
		return nil, err
	}

	paymentQuorumConfig, protocolConfig, err := a.paymentVaultParams.GetQuorumConfigs(meterer.OnDemandQuorumID)
	if err != nil {
		return nil, err
	}
	symbolsCharged := meterer.SymbolsCharged(numSymbols, protocolConfig.MinNumSymbols)
	paymentCharged := meterer.PaymentCharged(symbolsCharged, paymentQuorumConfig.OnDemandPricePerSymbol)

	a.onDemandLock.Lock()
	defer a.onDemandLock.Unlock()
	// calculate the increment required to add to the cumulative payment
	resultingPayment := new(big.Int).Add(a.cumulativePayment, paymentCharged)
	if resultingPayment.Cmp(a.onDemand.CumulativePayment) <= 0 {
		a.cumulativePayment.Add(a.cumulativePayment, paymentCharged)
		return a.cumulativePayment, nil
	}

	return nil, errors.New("insufficient ondemand payment")
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

	if onchainCumulativePayment == nil {
		a.onDemand = &core.OnDemandPayment{
			CumulativePayment: big.NewInt(0),
		}
	} else {
		a.onDemand = &core.OnDemandPayment{
			CumulativePayment: new(big.Int).Set(onchainCumulativePayment),
		}
	}

	if cumulativePayment == nil {
		a.cumulativePayment = big.NewInt(0)
	} else {
		a.cumulativePayment = new(big.Int).Set(cumulativePayment)
	}

	if reservations == nil {
		a.reservations = make(map[core.QuorumID]*core.ReservedPayment)
		a.periodRecords = make(meterer.QuorumPeriodRecords)
	} else {
		a.reservations = reservations
		if periodRecords == nil {
			a.periodRecords = make(meterer.QuorumPeriodRecords)
		} else {
			a.periodRecords = periodRecords
		}
	}

	return nil
}
