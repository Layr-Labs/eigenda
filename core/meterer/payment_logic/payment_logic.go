package payment_logic

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"slices"
	"time"

	"github.com/Layr-Labs/eigenda/core"
)

// GetBinLimit returns the bin limit given the bin interval and the symbols per second.
// The BinLimit serves to check if a reservation or global on-demand usage is within the rate limit.
//
// Parameters:
// SymbolsPerSecond is the number of symbols that can be charged per second, and can be understood as the rate of a reservation or global on-demand usage.
// BinInterval is the duration of a single bin in seconds.
//
// Returns:
// BinLimit is the maximum number of symbols that can be charged in a single bin.
func GetBinLimit(symbolsPerSecond uint64, binInterval uint64) uint64 {
	// If the rate for this bin is 0 or the bin interval is 0, then the bin limit is 0.
	// These two parameters should never be zero, but added for safety as to prevent division by zero in the overflow check.
	if symbolsPerSecond == 0 || binInterval == 0 {
		return 0
	}

	// Check for overflow before multiplication by comparing against the maximum safe value.
	if symbolsPerSecond > math.MaxUint64/binInterval {
		return math.MaxUint64
	}

	return symbolsPerSecond * binInterval
}

// GetReservationPeriodByNanosecond returns the current reservation period by finding the nearest lower multiple of the bin interval;
// bin interval used by the disperser is publicly recorded on-chain at the payment vault contract
func GetReservationPeriodByNanosecond(nanosecondTimestamp int64, binInterval uint64) uint64 {
	if nanosecondTimestamp < 0 {
		return 0
	}
	return GetReservationPeriod(int64((time.Duration(nanosecondTimestamp) * time.Nanosecond).Seconds()), binInterval)
}

// GetReservationPeriod returns the current reservation period by finding the nearest lower multiple of the bin interval;
// bin interval used by the disperser is publicly recorded on-chain at the payment vault contract
func GetReservationPeriod(timestamp int64, binInterval uint64) uint64 {
	if binInterval == 0 {
		return 0
	}
	return uint64(timestamp) / binInterval * binInterval
}

// GetOverflowPeriod returns the overflow period by adding the overflow offset to the current reservation period
// the offset is 2*reservationWindow, skipping the immediate next period for the period that will be used for overflow from the current period
func GetOverflowPeriod(reservationPeriod uint64, reservationWindow uint64) uint64 {
	return reservationPeriod + reservationWindow*2
}

// PaymentCharged returns the chargeable price for a given number of symbols
func PaymentCharged(numSymbols, pricePerSymbol uint64) *big.Int {
	// directly convert to uint64 to avoid overflow
	numSymbolsInt := new(big.Int).SetUint64(numSymbols)
	pricePerSymbolInt := new(big.Int).SetUint64(pricePerSymbol)
	return new(big.Int).Mul(numSymbolsInt, pricePerSymbolInt)
}

// SymbolsCharged returns the number of symbols charged for a given data length
// being at least MinNumSymbols or the nearest rounded-up multiple of MinNumSymbols.
func SymbolsCharged(numSymbols uint64, minSymbols uint64) uint64 {
	if numSymbols <= minSymbols {
		return minSymbols
	}
	if minSymbols == 0 {
		return numSymbols
	}
	// Round up to the nearest multiple of MinNumSymbols
	roundedUp := core.RoundUpDivide(numSymbols, minSymbols) * minSymbols
	// Check for overflow; this case should never happen
	if roundedUp < numSymbols {
		return math.MaxUint64
	}
	return roundedUp
}

// ValidateQuorum ensures that the quorums listed in the blobHeader are present within allowedQuorums
// Note: A reservation that does not utilize all of the allowed quorums will be accepted. However, it
// will still charge against all of the allowed quorums. A on-demand requests require and only allow
// the ETH and EIGEN quorums.
func ValidateQuorum(headerQuorums []uint8, allowedQuorums []uint8) error {
	if len(headerQuorums) == 0 {
		return fmt.Errorf("no quorum numbers provided in the request")
	}

	// check that all the quorum ids are in ReservedPayment's
	for _, q := range headerQuorums {
		if !slices.Contains(allowedQuorums, q) {
			// fail the entire request if there's a quorum number mismatch
			return fmt.Errorf("quorum number mismatch: %d", q)
		}
	}
	return nil
}

// ValidateReservations ensures that the quorums listed in the blobHeader are present within allowedQuorums.
//
// Notes:
//   - Reservations that don't use all allowed quorums are still accepted
//   - Charges apply to ALL allowed quorums, even if not all are used
//   - On-demand requests have special requirements: they must use ETH and EIGEN quorums only
func ValidateReservations(
	reservations map[core.QuorumID]*core.ReservedPayment,
	quorumConfigs map[core.QuorumID]*core.PaymentQuorumProtocolConfig,
	quorumNumbers []uint8,
	paymentHeaderTimestampNs int64,
	receivedTimestampNs int64,
) error {
	reservationQuorums := make([]uint8, 0, len(reservations))
	reservationWindows := make(map[core.QuorumID]uint64, len(reservations))
	requestReservationPeriods := make(map[core.QuorumID]uint64, len(reservations))

	// Gather quorums the user had an reservations on and relevant quorum configurations
	for quorumID, reservation := range reservations {
		quorumConfig, ok := quorumConfigs[quorumID]
		if !ok {
			return fmt.Errorf("quorum config not found for quorum %d", quorumID)
		}
		if quorumConfig == nil {
			return fmt.Errorf("quorum config is nil for quorum %d", quorumID)
		}
		if reservation == nil {
			return fmt.Errorf("reservation not found for quorum %d", quorumID)
		}
		reservationQuorums = append(reservationQuorums, uint8(quorumID))
		reservationWindows[quorumID] = quorumConfigs[quorumID].ReservationRateLimitWindow
		requestReservationPeriods[quorumID] = GetReservationPeriodByNanosecond(paymentHeaderTimestampNs, quorumConfigs[quorumID].ReservationRateLimitWindow)
	}
	if err := ValidateQuorum(quorumNumbers, reservationQuorums); err != nil {
		return err
	}
	// Validate the used reservations are active and is of valid periods
	for _, quorumID := range quorumNumbers {
		reservation := reservations[core.QuorumID(quorumID)]
		if !reservation.IsActiveByNanosecond(paymentHeaderTimestampNs) {
			return errors.New("reservation not active")
		}
		if !ValidateReservationPeriod(reservation, requestReservationPeriods[quorumID], reservationWindows[quorumID], receivedTimestampNs) {
			return fmt.Errorf("invalid reservation period for reservation on quorum %d", quorumID)
		}
	}

	return nil
}

// ValidateReservationPeriod checks if the provided reservation period is valid
// Note: This is called per-quorum since reservation is for a single quorum.
func ValidateReservationPeriod(reservation *core.ReservedPayment, requestReservationPeriod uint64, reservationWindow uint64, receivedTimestampNs int64) bool {
	if reservation == nil {
		return false
	}
	currentReservationPeriod := GetReservationPeriodByNanosecond(receivedTimestampNs, reservationWindow)
	// Valid reservation periods are either the current bin or the previous bin
	isCurrentOrPreviousPeriod := requestReservationPeriod == currentReservationPeriod || requestReservationPeriod == (currentReservationPeriod-reservationWindow)
	startPeriod := GetReservationPeriod(int64(reservation.StartTimestamp), reservationWindow)
	endPeriod := GetReservationPeriod(int64(reservation.EndTimestamp), reservationWindow)
	isWithinReservationWindow := startPeriod <= requestReservationPeriod && requestReservationPeriod < endPeriod
	if !isCurrentOrPreviousPeriod || !isWithinReservationWindow {
		return false
	}
	return true
}

// IsOnDemandPayment explicitly determines if the payment is an on-demand payment by checking if the cumulative payment is greater than 0
// If the cumulative payment is 0, it is not an on-demand payment, but a reservation payment.
func IsOnDemandPayment(paymentMetadata *core.PaymentMetadata) bool {
	return paymentMetadata.CumulativePayment.Cmp(big.NewInt(0)) > 0
}
