package ejector

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	dataapiv2 "github.com/Layr-Labs/eigenda/disperser/dataapi/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

var _ SigningRateLookup = (*dynamoSigningRateLookup)(nil)

// Uses batch information in dynamoDB to determine signing rates.
type dynamoSigningRateLookup struct {
	logger     logging.Logger
	url        string
	httpClient *http.Client
}

func NewDynamoSigningRateLookup(
	logger logging.Logger,
	url string,
	httpTimeout time.Duration,
) *dynamoSigningRateLookup {

	httpClient := &http.Client{
		Timeout: httpTimeout,
	}

	return &dynamoSigningRateLookup{
		logger:     logger,
		url:        url,
		httpClient: httpClient,
	}
}

func (srl *dynamoSigningRateLookup) GetSigningRates(
	timeSpan time.Duration,
	quorums []core.QuorumID,
	version ProtocolVersion,
	omitPerfectSigners bool,
) ([]*validator.ValidatorSigningRate, error) {
	switch version {
	case ProtocolVersionV1:
		if !omitPerfectSigners {
			srl.logger.Warn(
				"omitPerfectSigners flag is ignored for ProtocolVersionV1, will never return perfect signers")
		}
		return srl.getV1SigningRates(timeSpan, quorums)
	case ProtocolVersionV2:
		return srl.getV2SigningRates(timeSpan, quorums, omitPerfectSigners)
	default:
		return nil, fmt.Errorf("unsupported protocol version: %d", version)
	}
}

// Look up signing rates for v1.
func (srl *dynamoSigningRateLookup) getV1SigningRates(
	timeSpan time.Duration,
	quorums []core.QuorumID,
) ([]*validator.ValidatorSigningRate, error) {

	quorumSet := make(map[core.QuorumID]struct{})
	for _, q := range quorums {
		quorumSet[q] = struct{}{}
	}

	now := time.Now()

	path := "api/v1/metrics/operator-nonsigning-percentage"
	urlStr, err := url.JoinPath(srl.url, path)
	if err != nil {
		return nil, fmt.Errorf("error joining URL path with %s and %s: %w", srl.url, path, err)
	}
	url, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %w", err)
	}
	// add query parameters
	q := url.Query()
	// end: datetime formatted in "2006-01-02T15:04:05Z"
	q.Set("end", now.UTC().Format("2006-01-02T15:04:05Z"))
	// interval: lookback window in seconds
	q.Set("interval", strconv.Itoa(int(timeSpan.Seconds())))
	url.RawQuery = q.Encode()
	srl.logger.Debug("making request to DataAPI", "url", url.String())

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	resp, err := srl.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	srl.logger.Info("Received response", "responseBody", string(respBody))

	if resp.StatusCode != http.StatusOK {
		var errResp dataapi.ErrorResponse
		err = json.Unmarshal(respBody, &errResp)
		if err != nil {
			return nil, fmt.Errorf("error parsing error response: %w", err)
		}
		return nil, fmt.Errorf(
			"error response (%d) from dataapi: %s",
			resp.StatusCode,
			errResp.Error,
		)
	}

	var response dataapi.OperatorsNonsigningPercentage
	err = json.NewDecoder(strings.NewReader(string(respBody))).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error parsing response body: %w", err)
	}

	// Use a map to combine results from multiple quorums.
	signingRateMap := make(map[core.OperatorID]*validator.ValidatorSigningRate)

	for _, data := range response.Data {

		if len(quorumSet) > 0 {
			if _, ok := quorumSet[core.QuorumID(data.QuorumId)]; !ok {
				// This quorum is not in the requested set, skip it.
				continue
			}
		}

		signingRate, err := translateV1ToProto(data)
		if err != nil {
			return nil, fmt.Errorf("error translating dataapi rate to proto: %w", err)
		}

		signingRateMap[core.OperatorID(signingRate.ValidatorId)] =
			combineSigningRates(
				signingRateMap[core.OperatorID(signingRate.ValidatorId)],
				signingRate)
	}

	signingRates := make([]*validator.ValidatorSigningRate, 0, len(signingRateMap))
	for _, rate := range signingRateMap {
		signingRates = append(signingRates, rate)
	}

	return signingRates, nil
}

// Look up signing rates for v2.
func (srl *dynamoSigningRateLookup) getV2SigningRates(
	timeSpan time.Duration,
	quorums []core.QuorumID,
	omitPerfectSigners bool,
) ([]*validator.ValidatorSigningRate, error) {

	quorumSet := make(map[core.QuorumID]struct{})
	for _, q := range quorums {
		quorumSet[q] = struct{}{}
	}

	now := time.Now()

	path := "api/v2/operators/signing-info"
	urlStr, err := url.JoinPath(srl.url, path)
	if err != nil {
		return nil, fmt.Errorf("error joining URL path with %s and %s: %w", srl.url, path, err)
	}
	url, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %w", err)
	}
	// add query parameters
	q := url.Query()
	// end: datetime formatted in "2006-01-02T15:04:05Z"
	q.Set("end", now.UTC().Format("2006-01-02T15:04:05Z"))
	// interval: lookback window in seconds
	q.Set("interval", strconv.Itoa(int(timeSpan.Seconds())))
	if omitPerfectSigners {
		q.Set("nonsigner_only", "true")
	}
	url.RawQuery = q.Encode()
	srl.logger.Debug("making request to DataAPI", "url", url.String())

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	resp, err := srl.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	srl.logger.Info("Received response", "responseBody", string(respBody))

	if resp.StatusCode != http.StatusOK {
		var errResp dataapi.ErrorResponse
		err = json.Unmarshal(respBody, &errResp)
		if err != nil {
			return nil, fmt.Errorf("error parsing error response: %w", err)
		}
		return nil, fmt.Errorf(
			"error response (%d) from dataapi: %s",
			resp.StatusCode,
			errResp.Error,
		)
	}

	var response dataapiv2.OperatorsSigningInfoResponse
	err = json.NewDecoder(strings.NewReader(string(respBody))).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error parsing response body: %w", err)
	}

	// Use a map to combine results from multiple quorums.
	signingRateMap := make(map[core.OperatorID]*validator.ValidatorSigningRate)

	for _, data := range response.OperatorSigningInfo {
		if len(quorumSet) > 0 {
			if _, ok := quorumSet[core.QuorumID(data.QuorumId)]; !ok {
				// This quorum is not in the requested set, skip it.
				continue
			}
		}

		signingRate, err := translateV2ToProto(data)
		if err != nil {
			return nil, fmt.Errorf("error translating dataapi rate to proto: %w", err)
		}

		signingRateMap[core.OperatorID(signingRate.ValidatorId)] =
			combineSigningRates(
				signingRateMap[core.OperatorID(signingRate.ValidatorId)],
				signingRate)
	}

	signingRates := make([]*validator.ValidatorSigningRate, 0, len(signingRateMap))
	for _, rate := range signingRateMap {
		signingRates = append(signingRates, rate)
	}

	return signingRates, nil
}

// Translates a single DataAPI OperatorNonsigningPercentageMetrics to a ValidatorSigningRate protobuf.
func translateV1ToProto(data *dataapi.OperatorNonsigningPercentageMetrics) (*validator.ValidatorSigningRate, error) {
	validatorID, err := core.OperatorIDFromHex(data.OperatorId)
	if err != nil {
		return nil, fmt.Errorf("error parsing operator ID %s: %w", data.OperatorId, err)
	}

	signedBatches := data.TotalBatches - data.TotalUnsignedBatches
	unsignedBatches := data.TotalUnsignedBatches

	signingRate := &validator.ValidatorSigningRate{
		ValidatorId:     validatorID[:],
		SignedBatches:   uint64(signedBatches),
		UnsignedBatches: uint64(unsignedBatches),
		SignedBytes:     uint64(signedBatches),   // Not accurate, but we don't have byte info from DataAPI.
		UnsignedBytes:   uint64(unsignedBatches), // Not accurate, but we don't have byte info from DataAPI.
		SigningLatency:  0,                       // Not available from DataAPI.
	}

	return signingRate, nil
}

// Translates a single DataAPI v2 OperatorSigningInfo to a ValidatorSigningRate protobuf.
func translateV2ToProto(data *dataapiv2.OperatorSigningInfo) (*validator.ValidatorSigningRate, error) {
	validatorID, err := core.OperatorIDFromHex(data.OperatorId)
	if err != nil {
		return nil, fmt.Errorf("error parsing operator ID %s: %w", data.OperatorId, err)
	}

	signedBatches := data.TotalBatches - data.TotalUnsignedBatches
	unsignedBatches := data.TotalUnsignedBatches

	signingRate := &validator.ValidatorSigningRate{
		ValidatorId:     validatorID[:],
		SignedBatches:   uint64(signedBatches),
		UnsignedBatches: uint64(unsignedBatches),
		SignedBytes:     uint64(signedBatches),   // Not accurate, but we don't have byte info from DataAPI.
		UnsignedBytes:   uint64(unsignedBatches), // Not accurate, but we don't have byte info from DataAPI.
		SigningLatency:  0,                       // Not available from DataAPI.
	}

	return signingRate, nil
}
