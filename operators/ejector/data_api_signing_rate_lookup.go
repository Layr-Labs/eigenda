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
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
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
	version ProtocolVersion,
) ([]*validator.ValidatorSigningRate, error) {
	switch version {
	case ProtocolVersionV1:
		return srl.getV1SigningRates(timeSpan)
	case ProtocolVersionV2:
		return srl.getV2SigningRates(timeSpan)
	default:
		return nil, fmt.Errorf("unsupported protocol version: %d", version)
	}
}

// Look up signing rates for v1.
func (srl *dynamoSigningRateLookup) getV1SigningRates(
	timeSpan time.Duration,
) ([]*validator.ValidatorSigningRate, error) {

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

	// return &response, nil
	return nil, nil // TODO

}

// Look up signing rates for v2.
func (srl *dynamoSigningRateLookup) getV2SigningRates(
	timeSpan time.Duration,
) ([]*validator.ValidatorSigningRate, error) {
	//TODO implement me
	panic("implement me")
}
