package ejector

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"sort"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// Ejector is responsible for periodically evaluating validators and deciding which ones to eject.
type Ejector struct {
	ctx    context.Context
	logger logging.Logger
	config *EjectorConfig

	// Responsible for executing ejections and managing the ejection lifecycle.
	ejectionManager *ThreadedEjectionManager

	// Used for looking up signing rates for V1.
	// TODO(cody.littley): remove this after V1 sunset
	signingRateLookupV1 SigningRateLookup

	// Used for looking up signing rates for V2.
	signingRateLookupV2 SigningRateLookup

	// Used to convert validator IDs to validator addresses.
	validatorIDToAddressCache eth.ValidatorIDToAddressConverter

	// Used to look up the latest reference number.
	referenceBlockProvider eth.ReferenceBlockProvider

	// Used to look up which quorums a validator is a member of.
	validatorQuorumLookup eth.ValidatorQuorumLookup

	// Used to look up validator stake fractions.
	validatorStakeLookup eth.ValidatorStakeLookup

	// The last time signing rates were logged verbosely.
	lastSigningRateLogTime time.Time
}

// NewEjector creates a new Ejector.
func NewEjector(
	ctx context.Context,
	logger logging.Logger,
	config *EjectorConfig,
	ejectionManager *ThreadedEjectionManager,
	signingRateLookupV1 SigningRateLookup,
	signingRateLookupV2 SigningRateLookup,
	validatorIDToAddressCache eth.ValidatorIDToAddressConverter,
	referenceBlockProvider eth.ReferenceBlockProvider,
	validatorQuorumLookup eth.ValidatorQuorumLookup,
	validatorStakeLookup eth.ValidatorStakeLookup,
) *Ejector {
	e := &Ejector{
		ctx:                       ctx,
		logger:                    logger,
		config:                    config,
		ejectionManager:           ejectionManager,
		signingRateLookupV1:       signingRateLookupV1,
		signingRateLookupV2:       signingRateLookupV2,
		validatorIDToAddressCache: validatorIDToAddressCache,
		referenceBlockProvider:    referenceBlockProvider,
		validatorQuorumLookup:     validatorQuorumLookup,
		validatorStakeLookup:      validatorStakeLookup,
	}

	go e.mainLoop()

	return e
}

// The main loop periodically evaluates validators for ejection.
func (e *Ejector) mainLoop() {
	e.logger.Debugf("ejector started, evaluating validators for ejection every %s", e.config.EjectionPeriod.String())

	ticker := time.NewTicker(e.config.EjectionPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-e.ctx.Done():
			e.logger.Info("ejector shutting down")
			return
		case <-ticker.C:
			err := e.evaluateValidators()
			if err != nil {
				e.logger.Error("error evaluating validators", "error", err)
			}
		}
	}
}

// evaluateValidators looks up signing rates and evaluates which validators should be ejected.
func (e *Ejector) evaluateValidators() error {

	e.logger.Debug("evaluating validators for ejection")

	omitPerfectSigners := true
	if e.config.SigningRateLogPeriod > 0 {
		omitPerfectSigners = false
	}

	v1SigningRates, err := e.signingRateLookupV1.GetSigningRates(
		e.config.EjectionCriteriaTimeWindow,
		nil, // all quorums
		ProtocolVersionV1,
		omitPerfectSigners,
	)
	if err != nil {
		return fmt.Errorf("error looking up v1 signing rates: %w", err)
	}

	v2SigningRates, err := e.signingRateLookupV2.GetSigningRates(
		e.config.EjectionCriteriaTimeWindow,
		nil, // all quorums
		ProtocolVersionV2,
		omitPerfectSigners,
	)
	if err != nil {
		return fmt.Errorf("error looking up v2 signing rates: %w", err)
	}

	e.logSigningRates(v1SigningRates, v2SigningRates)

	// Combine data from v1 and v2 lookups, since the validator is likely to cancel ejection if it is active in either.
	signingRates, err := combineSigningRateSlices(v1SigningRates, v2SigningRates)
	if err != nil {
		return fmt.Errorf("error combining signing rates: %w", err)
	}
	sortByUnsignedBytesDescending(signingRates)

	for _, signingRate := range signingRates {
		err := e.evaluateValidator(signingRate)
		if err != nil {
			e.logger.Error("error evaluating validator", "validatorID", signingRate.GetValidatorId(), "error", err)
		}
	}

	return nil
}

// evaluateValidator evaluates a single validator's signing rate and decides whether to eject it.
func (e *Ejector) evaluateValidator(signingRate *validator.ValidatorSigningRate) error {
	isEjectable := signingRate.GetSignedBatches() == 0 && signingRate.GetUnsignedBatches() > 0
	if !isEjectable {
		return nil
	}

	if len(signingRate.GetValidatorId()) != 32 {
		return fmt.Errorf("invalid validator ID %s, length is not 32", hex.EncodeToString(signingRate.GetValidatorId()))
	}

	validatorID := core.OperatorID(signingRate.GetValidatorId()[:])
	validatorAddress, err := e.validatorIDToAddressCache.ValidatorIDToAddress(e.ctx, validatorID)
	if err != nil {
		return fmt.Errorf("error converting validator ID to address: %w", err)
	}

	stakeFractions, err := e.getStakeFractionMap(validatorID)
	if err != nil {
		return fmt.Errorf("error calculating stake fractions: %w", err)
	}

	e.logger.Debug("Validator is eligible for ejection",
		"validatorID", core.OperatorID(signingRate.GetValidatorId()).Hex(),
		"validatorAddress", validatorAddress.Hex(),
		"signedBatches", signingRate.GetSignedBatches(),
		"unsignedBatches", signingRate.GetUnsignedBatches(),
		"signedBytes", signingRate.GetSignedBytes(),
		"unsignedBytes", signingRate.GetUnsignedBytes(),
		"stakeFractions", stakeFractions,
	)

	// The ejection manager is responsible for deduplicating ejection requests, and deciding if
	// there are other factors that may prevent ejection (e.g. too many ejection attempts, etc.).
	err = e.ejectionManager.EjectValidator(validatorAddress, stakeFractions)
	if err != nil {
		return fmt.Errorf("error requesting ejection: %w", err)
	}

	return nil
}

// Get the stake fraction map for a given validator.
func (e *Ejector) getStakeFractionMap(validatorID core.OperatorID) (map[core.QuorumID]float64, error) {

	rbn, err := e.referenceBlockProvider.GetReferenceBlockNumber(e.ctx)
	if err != nil {
		return nil, fmt.Errorf("error looking up latest reference block number: %w", err)
	}

	quorums, err := e.validatorQuorumLookup.GetQuorumsForValidator(
		e.ctx,
		validatorID,
		rbn,
	)
	if err != nil {
		return nil, fmt.Errorf("error looking up quorums for validator: %w", err)
	}

	stakeFractions := make(map[core.QuorumID]float64, len(quorums))

	for _, quorumID := range quorums {
		stakeFraction, err := e.validatorStakeLookup.GetValidatorStakeFraction(
			e.ctx,
			quorumID,
			validatorID,
			rbn,
		)
		if err != nil {
			return nil, fmt.Errorf("error looking up stake fraction for validator %x in quorum %d: %w",
				validatorID, quorumID, err)
		}
		stakeFractions[quorumID] = stakeFraction
	}

	return stakeFractions, nil
}

// Verbose logging of signing rates. Useful for debugging problems with signing rate tracking/lookup.
func (e *Ejector) logSigningRates(
	v1SigningRates []*validator.ValidatorSigningRate,
	v2SigningRates []*validator.ValidatorSigningRate,
) {

	if e.config.SigningRateLogPeriod == 0 {
		// Signing rate logging disabled.
		return
	}

	if time.Since(e.lastSigningRateLogTime) < e.config.SigningRateLogPeriod {
		// Not time to log yet.
		return
	}

	e.lastSigningRateLogTime = time.Now()

	// Create an index of all validators seen in either v1 or v2 data.
	validatorSet := make(map[core.OperatorID]struct{})
	v1Data := make(map[core.OperatorID]*validator.ValidatorSigningRate)
	v2Data := make(map[core.OperatorID]*validator.ValidatorSigningRate)

	for _, rate := range v1SigningRates {
		id := core.OperatorID(rate.GetValidatorId()[:])
		validatorSet[id] = struct{}{}
		v1Data[id] = rate
	}
	for _, rate := range v2SigningRates {
		id := core.OperatorID(rate.GetValidatorId()[:])
		validatorSet[id] = struct{}{}
		v2Data[id] = rate
	}

	validatorList := make([]core.OperatorID, 0, len(validatorSet))
	for id := range validatorSet {
		validatorList = append(validatorList, id)
	}

	// Sort by validator ID for consistent logging order.
	sort.Slice(validatorList, func(i, j int) bool {
		return bytes.Compare(validatorList[i][:], validatorList[j][:]) < 0
	})

	// Build JSON string manually for compact output.
	var jsonBuilder bytes.Buffer
	jsonBuilder.WriteString("[")

	for i, id := range validatorList {
		if i > 0 {
			jsonBuilder.WriteString(",")
		}

		v1Rate, v1Exists := v1Data[id]
		v2Rate, v2Exists := v2Data[id]

		jsonBuilder.WriteString(fmt.Sprintf(`{"id":"%s"`, id.Hex()))

		if v1Exists {
			jsonBuilder.WriteString(fmt.Sprintf(
				`,"v1":{"sb":%d,"ub":%d,"sy":%d,"uy":%d}`,
				v1Rate.GetSignedBatches(),
				v1Rate.GetUnsignedBatches(),
				v1Rate.GetSignedBytes(),
				v1Rate.GetUnsignedBytes(),
			))
		}

		if v2Exists {
			jsonBuilder.WriteString(fmt.Sprintf(
				`,"v2":{"sb":%d,"ub":%d,"sy":%d,"uy":%d}`,
				v2Rate.GetSignedBatches(),
				v2Rate.GetUnsignedBatches(),
				v2Rate.GetSignedBytes(),
				v2Rate.GetUnsignedBytes(),
			))
		}

		jsonBuilder.WriteString("}")
	}

	jsonBuilder.WriteString("]")

	e.logger.Info("Signing rates",
		"key", "id=validator_id, sb=signed_batches, ub=unsigned_batches, sy=signed_bytes, uy=unsigned_bytes",
		"data", jsonBuilder.String())
}
