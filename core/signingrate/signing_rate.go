package signingrate

import (
	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/core"
)

// SigningRate records information about validator signing rate during a time period.
//
// This struct mirrors the validator.ValidatorSigningRate protobuf message. Although it's generally not good practice
// to mirror protobufs when not necessary, in this case it is needed in order to attach methods to the object.
type SigningRate struct {
	// The unique identifier of the validator (i.e. the operator ID).
	id core.OperatorID

	// The number of signed by the validator during the period.
	signedBatches uint64

	// The number of unsigned by the validator during the period.
	unsignedBatches uint64

	// The total number of bytes signed during the period.
	signedBytes uint64

	// The total number of bytes unsigned during the period.
	unsignedBytes uint64

	// Contains the sum of the time spent by the validator waiting for signing requests to be processed, in nanoseconds.
	// Only batches that are signed are considered (i.e. if the validator does not succeed in signing a batch,
	// the time spend in the attempt is not counted).
	signingLatency uint64

	// The number of nanoseconds that the validator is considered "up" by the disperser. When a validator successfully
	// signs for a batch, it is considered "up" until it misses a batch. This number can be used to differentiate
	// between a validator that occasionally misses batches, and one that is fully offline or malfunctioning.
	uptime uint64

	// A cached protobuf representation of this signing rate. Set to nil whenever the signing rate is modified.
	cachedProtobuf *validator.ValidatorSigningRate
}

// NewSigningRate creates a new SigningRate instance.
func NewSigningRate(id core.OperatorID) *SigningRate {
	return &SigningRate{
		id: id,
	}
}

// NewSigningRateFromProtobuf creates a new SigningRate instance from a protobuf representation.
func NewSigningRateFromProtobuf(proto *validator.ValidatorSigningRate) *SigningRate {
	return &SigningRate{
		id:              (core.OperatorID)(proto.Id),
		signedBatches:   proto.SignedBatches,
		unsignedBatches: proto.UnsignedBatches,
		signedBytes:     proto.SignedBytes,
		unsignedBytes:   proto.UnsignedBytes,
		signingLatency:  proto.SigningLatency,
		uptime:          proto.Uptime,
		cachedProtobuf:  proto,
	}
}

// ID returns the unique identifier of the validator.
func (s *SigningRate) ID() core.OperatorID {
	return s.id
}

// SignedBatches returns the number of batches signed by the validator during the period.
func (s *SigningRate) SignedBatches() uint64 {
	return s.signedBatches
}

// Add one to the number of batches signed by the validator during the period.
func (s *SigningRate) IncrementSignedBatches() {
	s.signedBatches++
	s.cachedProtobuf = nil
}

// UnsignedBatches returns the number of batches unsigned by the validator during the period.
func (s *SigningRate) UnsignedBatches() uint64 {
	return s.unsignedBatches
}

// Add one to the number of batches unsigned by the validator during the period.
func (s *SigningRate) IncrementUnsignedBatches() {
	s.unsignedBatches++
	s.cachedProtobuf = nil
}

// SignedBytes returns the total number of bytes signed during the period.
func (s *SigningRate) SignedBytes() uint64 {
	return s.signedBytes
}

// AddSignedBytes adds the given number of bytes to the total number of bytes signed during the period.
func (s *SigningRate) AddSignedBytes(bytes uint64) {
	s.signedBytes += bytes
	s.cachedProtobuf = nil
}

// UnsignedBytes returns the total number of bytes unsigned during the period.
func (s *SigningRate) UnsignedBytes() uint64 {
	return s.unsignedBytes
}

func (s *SigningRate) AddUnsignedBytes(bytes uint64) {
	s.unsignedBytes += bytes
	s.cachedProtobuf = nil
}

// SigningLatency returns the sum of the time spent by the validator waiting for signing requests to be processed,
// in nanoseconds.
func (s *SigningRate) SigningLatency() uint64 {
	return s.signingLatency
}

// AddSigningLatency adds the given number of nanoseconds to the sum of the time spent by the validator waiting
// for signing requests to be processed.
func (s *SigningRate) AddSigningLatency(latency uint64) {
	s.signingLatency += latency
	s.cachedProtobuf = nil
}

// Uptime returns the number of nanoseconds that the validator is considered "up" by the disperser.
func (s *SigningRate) Uptime() uint64 {
	return s.uptime
}

// AddUptime adds the given number of nanoseconds to the uptime of the validator.
func (s *SigningRate) AddUptime(uptime uint64) {
	s.uptime += uptime
	s.cachedProtobuf = nil
}

// ToProtobuf returns a protobuf representation of this signing rate. The protobuf is read safe even if
// this SigningRate object is modified concurrently.
func (s *SigningRate) ToProtobuf() *validator.ValidatorSigningRate {
	if s.cachedProtobuf != nil {
		return s.cachedProtobuf
	}

	s.cachedProtobuf = &validator.ValidatorSigningRate{
		Id:              s.id[:],
		SignedBatches:   s.signedBatches,
		UnsignedBatches: s.unsignedBatches,
		SignedBytes:     s.signedBytes,
		UnsignedBytes:   s.unsignedBytes,
		SigningLatency:  s.signingLatency,
		Uptime:          s.uptime,
	}
	return s.cachedProtobuf
}
