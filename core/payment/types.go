package payment

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"

	commonpbv2 "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

// PaymentQuorumConfig contains the configuration for a quorum's payment configurations
// This is pretty much the same as the PaymentVaultTypesQuorumConfig struct in the contracts/bindings/IPaymentVault/binding.go file
type PaymentQuorumConfig struct {
	ReservationSymbolsPerSecond uint64

	// OnDemand is initially only enabled on Quorum 0
	OnDemandSymbolsPerSecond uint64
	OnDemandPricePerSymbol   uint64
}

// PaymentQuorumProtocolConfig contains the configuration for a quorum's ratelimiting configurations
// This is pretty much the same as the PaymentVaultTypesQuorumProtocolConfig struct in the contracts/bindings/IPaymentVault/binding.go file
type PaymentQuorumProtocolConfig struct {
	MinNumSymbols              uint64
	ReservationAdvanceWindow   uint64
	ReservationRateLimitWindow uint64

	// OnDemand is initially only enabled on Quorum 0
	OnDemandRateLimitWindow uint64
	OnDemandEnabled         bool
}

// PaymentVaultParams contains all configuration parameters for the payment vault
type PaymentVaultParams struct {
	QuorumPaymentConfigs  map[uint8]*PaymentQuorumConfig
	QuorumProtocolConfigs map[uint8]*PaymentQuorumProtocolConfig
	OnDemandQuorumNumbers []uint8
}

// GetQuorumConfigs retrieves payment and protocol configurations for a specific quorum
func (pvp *PaymentVaultParams) GetQuorumConfigs(quorumID uint8) (*PaymentQuorumConfig, *PaymentQuorumProtocolConfig, error) {
	paymentConfig, ok := pvp.QuorumPaymentConfigs[quorumID]
	if !ok {
		return nil, nil, fmt.Errorf("payment config not found for quorum %d", quorumID)
	}
	protocolConfig, ok := pvp.QuorumProtocolConfigs[quorumID]
	if !ok {
		return nil, nil, fmt.Errorf("protocol config not found for quorum %d", quorumID)
	}
	return paymentConfig, protocolConfig, nil
}

// PaymentVaultParamsFromProtobuf converts a protobuf payment vault params to a core payment vault params
func PaymentVaultParamsFromProtobuf(vaultParams *disperser_rpc.PaymentVaultParams) (*PaymentVaultParams, error) {
	if vaultParams == nil {
		return nil, fmt.Errorf("payment vault params cannot be nil")
	}

	if vaultParams.GetQuorumPaymentConfigs() == nil {
		return nil, fmt.Errorf("payment quorum configs cannot be nil")
	}

	if vaultParams.GetQuorumProtocolConfigs() == nil {
		return nil, fmt.Errorf("payment quorum protocol configs cannot be nil")
	}

	// Convert protobuf configs to core types
	quorumPaymentConfigs := make(map[uint8]*PaymentQuorumConfig)
	quorumProtocolConfigs := make(map[uint8]*PaymentQuorumProtocolConfig)

	for quorumID, pbPaymentConfig := range vaultParams.GetQuorumPaymentConfigs() {
		quorumPaymentConfigs[uint8(quorumID)] = &PaymentQuorumConfig{
			ReservationSymbolsPerSecond: pbPaymentConfig.GetReservationSymbolsPerSecond(),
			OnDemandSymbolsPerSecond:    pbPaymentConfig.GetOnDemandSymbolsPerSecond(),
			OnDemandPricePerSymbol:      pbPaymentConfig.GetOnDemandPricePerSymbol(),
		}
	}

	for quorumID, pbProtocolConfig := range vaultParams.GetQuorumProtocolConfigs() {
		quorumProtocolConfigs[uint8(quorumID)] = &PaymentQuorumProtocolConfig{
			MinNumSymbols:              pbProtocolConfig.GetMinNumSymbols(),
			ReservationAdvanceWindow:   pbProtocolConfig.GetReservationAdvanceWindow(),
			ReservationRateLimitWindow: pbProtocolConfig.GetReservationRateLimitWindow(),
			OnDemandRateLimitWindow:    pbProtocolConfig.GetOnDemandRateLimitWindow(),
			OnDemandEnabled:            pbProtocolConfig.GetOnDemandEnabled(),
		}
	}
	// Convert uint32 slice to uint8 slice
	onDemandQuorumNumbers := make([]uint8, len(vaultParams.GetOnDemandQuorumNumbers()))
	for i, num := range vaultParams.GetOnDemandQuorumNumbers() {
		onDemandQuorumNumbers[i] = uint8(num)
	}
	return &PaymentVaultParams{
		QuorumPaymentConfigs:  quorumPaymentConfigs,
		QuorumProtocolConfigs: quorumProtocolConfigs,
		OnDemandQuorumNumbers: onDemandQuorumNumbers,
	}, nil
}

// PaymentVaultParamsToProtobuf converts core payment vault params to protobuf format
func (pvp *PaymentVaultParams) PaymentVaultParamsToProtobuf() (*disperser_rpc.PaymentVaultParams, error) {
	if pvp == nil {
		return nil, fmt.Errorf("payment vault params cannot be nil")
	}

	if pvp.QuorumPaymentConfigs == nil {
		return nil, fmt.Errorf("payment quorum configs cannot be nil")
	}

	if pvp.QuorumProtocolConfigs == nil {
		return nil, fmt.Errorf("payment quorum protocol configs cannot be nil")
	}

	quorumPaymentConfigs := make(map[uint32]*disperser_rpc.PaymentQuorumConfig)
	quorumProtocolConfigs := make(map[uint32]*disperser_rpc.PaymentQuorumProtocolConfig)

	for quorumID, paymentConfig := range pvp.QuorumPaymentConfigs {
		quorumPaymentConfigs[uint32(quorumID)] = &disperser_rpc.PaymentQuorumConfig{
			ReservationSymbolsPerSecond: paymentConfig.ReservationSymbolsPerSecond,
			OnDemandSymbolsPerSecond:    paymentConfig.OnDemandSymbolsPerSecond,
			OnDemandPricePerSymbol:      paymentConfig.OnDemandPricePerSymbol,
		}
	}

	for quorumID, protocolConfig := range pvp.QuorumProtocolConfigs {
		quorumProtocolConfigs[uint32(quorumID)] = &disperser_rpc.PaymentQuorumProtocolConfig{
			MinNumSymbols:              protocolConfig.MinNumSymbols,
			ReservationAdvanceWindow:   protocolConfig.ReservationAdvanceWindow,
			ReservationRateLimitWindow: protocolConfig.ReservationRateLimitWindow,
			OnDemandRateLimitWindow:    protocolConfig.OnDemandRateLimitWindow,
			OnDemandEnabled:            protocolConfig.OnDemandEnabled,
		}
	}

	onDemandQuorumNumbers := make([]uint32, len(pvp.OnDemandQuorumNumbers))
	for i, num := range pvp.OnDemandQuorumNumbers {
		onDemandQuorumNumbers[i] = uint32(num)
	}

	return &disperser_rpc.PaymentVaultParams{
		QuorumPaymentConfigs:  quorumPaymentConfigs,
		QuorumProtocolConfigs: quorumProtocolConfigs,
		OnDemandQuorumNumbers: onDemandQuorumNumbers,
	}, nil
}

// ReservationsFromProtobuf converts protobuf reservations to native types
func ReservationsFromProtobuf(pbReservations map[uint32]*disperser_rpc.QuorumReservation) map[uint8]*ReservedPayment {
	if pbReservations == nil {
		return nil
	}

	reservations := make(map[uint8]*ReservedPayment)
	for quorumNumber, reservation := range pbReservations {
		if reservation == nil {
			continue
		}
		quorumID := uint8(quorumNumber)
		reservations[quorumID] = &ReservedPayment{
			SymbolsPerSecond: reservation.GetSymbolsPerSecond(),
			StartTimestamp:   uint64(reservation.GetStartTimestamp()),
			EndTimestamp:     uint64(reservation.GetEndTimestamp()),
		}
	}
	return reservations
}

// CumulativePaymentFromProtobuf converts protobuf payment bytes to *big.Int
func CumulativePaymentFromProtobuf(paymentBytes []byte) *big.Int {
	if paymentBytes == nil {
		return nil
	}
	return new(big.Int).SetBytes(paymentBytes)
}

// ConvertPaymentStateFromProtobuf converts a protobuf GetPaymentStateForAllQuorumsReply to native types
func ConvertPaymentStateFromProtobuf(paymentStateProto *disperser_rpc.GetPaymentStateForAllQuorumsReply) (
	*PaymentVaultParams,
	map[uint8]*ReservedPayment,
	*big.Int,
	*big.Int,
	QuorumPeriodRecords,
	error,
) {
	if paymentStateProto == nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("payment state cannot be nil")
	}

	paymentVaultParams, err := PaymentVaultParamsFromProtobuf(paymentStateProto.GetPaymentVaultParams())
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error converting payment vault params: %w", err)
	}

	reservations := ReservationsFromProtobuf(paymentStateProto.GetReservations())

	cumulativePayment := CumulativePaymentFromProtobuf(paymentStateProto.GetCumulativePayment())
	onchainCumulativePayment := CumulativePaymentFromProtobuf(paymentStateProto.GetOnchainCumulativePayment())

	var periodRecords QuorumPeriodRecords
	if paymentStateProto.GetPeriodRecords() != nil {
		periodRecords = FromProtoRecords(paymentStateProto.GetPeriodRecords())
	}

	return paymentVaultParams, reservations, cumulativePayment, onchainCumulativePayment, periodRecords, nil
}

// PaymentMetadata represents the header information for a blob
type PaymentMetadata struct {
	// AccountID is the ETH account address for the payer
	AccountID gethcommon.Address `json:"account_id"`

	// Timestamp represents the nanosecond of the dispersal request creation
	Timestamp int64 `json:"timestamp"`
	// CumulativePayment represents the total amount of payment (in wei) made by the user up to this point
	CumulativePayment *big.Int `json:"cumulative_payment"`
}

// Hash returns the Keccak256 hash of the PaymentMetadata
func (pm *PaymentMetadata) Hash() ([32]byte, error) {
	if pm == nil {
		return [32]byte{}, errors.New("payment metadata is nil")
	}
	blobHeaderType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{
			Name: "accountID",
			Type: "string",
		},
		{
			Name: "timestamp",
			Type: "int64",
		},
		{
			Name: "cumulativePayment",
			Type: "uint256",
		},
	})
	if err != nil {
		return [32]byte{}, err
	}

	arguments := abi.Arguments{
		{
			Type: blobHeaderType,
		},
	}

	s := struct {
		AccountID         string
		Timestamp         int64
		CumulativePayment *big.Int
	}{
		AccountID:         pm.AccountID.Hex(),
		Timestamp:         pm.Timestamp,
		CumulativePayment: pm.CumulativePayment,
	}

	bytes, err := arguments.Pack(s)
	if err != nil {
		return [32]byte{}, err
	}

	var hash [32]byte
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(bytes)
	copy(hash[:], hasher.Sum(nil)[:32])

	return hash, nil
}

func (pm *PaymentMetadata) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	if pm == nil {
		return nil, errors.New("payment metadata is nil")
	}

	return &types.AttributeValueMemberM{
		Value: map[string]types.AttributeValue{
			"AccountID": &types.AttributeValueMemberS{Value: pm.AccountID.Hex()},
			"Timestamp": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", pm.Timestamp)},
			"CumulativePayment": &types.AttributeValueMemberN{
				Value: pm.CumulativePayment.String(),
			},
		},
	}, nil
}

func (pm *PaymentMetadata) UnmarshalDynamoDBAttributeValue(av types.AttributeValue) error {
	m, ok := av.(*types.AttributeValueMemberM)
	if !ok {
		return fmt.Errorf("expected *types.AttributeValueMemberM, got %T", av)
	}
	accountID, ok := m.Value["AccountID"].(*types.AttributeValueMemberS)
	if !ok {
		return fmt.Errorf("expected *types.AttributeValueMemberS for AccountID, got %T", m.Value["AccountID"])
	}
	pm.AccountID = gethcommon.HexToAddress(accountID.Value)
	rp, ok := m.Value["Timestamp"].(*types.AttributeValueMemberN)
	if !ok {
		return fmt.Errorf("expected *types.AttributeValueMemberN for Timestamp, got %T", m.Value["Timestamp"])
	}
	timestamp, err := strconv.ParseInt(rp.Value, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse Timestamp: %w", err)
	}
	pm.Timestamp = timestamp
	cp, ok := m.Value["CumulativePayment"].(*types.AttributeValueMemberN)
	if !ok {
		return fmt.Errorf("expected *types.AttributeValueMemberN for CumulativePayment, got %T", m.Value["CumulativePayment"])
	}
	pm.CumulativePayment, _ = new(big.Int).SetString(cp.Value, 10)
	return nil
}

// ReservedPayment contains information the onchain state about a reserved payment
type ReservedPayment struct {
	// reserve number of symbols per second
	SymbolsPerSecond uint64
	// reservation activation timestamp
	StartTimestamp uint64
	// reservation expiration timestamp
	EndTimestamp uint64

	// allowed quorums
	QuorumNumbers []uint8
	// ordered mapping of quorum number to payment split; on-chain validation should ensure split <= 100
	QuorumSplits []byte
}

type OnDemandPayment struct {
	// Total amount deposited by the user
	CumulativePayment *big.Int
}

// IsActive returns true if the reservation is active at the given timestamp
func (ar *ReservedPayment) IsActive(currentTimestamp uint64) bool {
	return WithinTime(time.Unix(int64(currentTimestamp), 0), time.Unix(int64(ar.StartTimestamp), 0), time.Unix(int64(ar.EndTimestamp), 0))
}

// IsActiveByNanosecond returns true if the reservation is active at the given timestamp
func (ar *ReservedPayment) IsActiveByNanosecond(currentTimestamp int64) bool {
	return WithinTime(time.Unix(0, currentTimestamp), time.Unix(int64(ar.StartTimestamp), 0), time.Unix(int64(ar.EndTimestamp), 0))
}

// WithinTime returns true if the timestamp is within the time range, inclusive of the start and end timestamps
func WithinTime(timestamp time.Time, startTimestamp time.Time, endTimestamp time.Time) bool {
	return !timestamp.Before(startTimestamp) && !timestamp.After(endTimestamp)
}

// ToProtobuf converts PaymentMetadata to protobuf format
func (pm *PaymentMetadata) ToProtobuf() *commonpbv2.PaymentHeader {
	if pm == nil {
		return nil
	}
	return &commonpbv2.PaymentHeader{
		AccountId:         pm.AccountID.Hex(),
		Timestamp:         pm.Timestamp,
		CumulativePayment: pm.CumulativePayment.Bytes(),
	}
}

// ConvertToPaymentMetadata converts a protobuf payment header to PaymentMetadata
func ConvertToPaymentMetadata(ph *commonpbv2.PaymentHeader) (*PaymentMetadata, error) {
	if ph == nil {
		return nil, nil
	}

	if !gethcommon.IsHexAddress(ph.AccountId) {
		return nil, fmt.Errorf("invalid account ID: %s", ph.AccountId)
	}

	return &PaymentMetadata{
		AccountID:         gethcommon.HexToAddress(ph.AccountId),
		Timestamp:         ph.Timestamp,
		CumulativePayment: new(big.Int).SetBytes(ph.CumulativePayment),
	}, nil
}
