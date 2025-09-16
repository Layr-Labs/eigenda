package enabled_apis

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Layr-Labs/eigenda/api/proxy/common"
)

// EnabledAPIs is a wrapper type of the APIs set enum and provides
// encapsulation for the entrypoint processing to make enablement determination
// through specific boolean accessors.
type EnabledAPIs struct {
	apis []API
}

func (e EnabledAPIs) Count() uint {
	return uint(len(e.apis))
}

// ArbCustomDA ... Is Arbitrum Custom DA enabled?
func (e EnabledAPIs) ArbCustomDA() bool {
	return e.has(ArbCustomDAServer)
}

// RestALTDA ... Is the REST ALTDA server enabled?
func (e EnabledAPIs) RestALTDA() bool {
	return e.has(Admin) || e.has(OpGenericCommitment) ||
		e.has(OpKeccakCommitment) || e.has(StandardCommitment)
}

// RestALTDAOPGeneric ... Is REST server enabled with op-generic commitment mode
func (e EnabledAPIs) RestALTDAOPGeneric() bool {
	return e.has(OpGenericCommitment)
}

// RestALTDAOPKeccak ... Is REST server enabled with op-keccak commitment mode
func (e EnabledAPIs) RestALTDAOPKeccak() bool {
	return e.has(OpKeccakCommitment)
}

// RestALTStandard ... Is REST server enabled with standard commitment mode
func (e EnabledAPIs) RestALTStandard() bool {
	return e.has(StandardCommitment)
}

// RestALTDAWithAdmin ... Is REST server enabled with admin mode
func (e EnabledAPIs) RestALTDAWithAdmin() bool {
	return e.has(Admin)
}

// Metrics ... Are metrics exposed on the metrics server?
func (e EnabledAPIs) Metrics() bool {
	return e.has(MetricsServer)
}

// Check ... Ensures that expression of the enabled API set is correct
func (e EnabledAPIs) Check() error {
	if e.Count() == 0 {
		return fmt.Errorf("expected at least one \"apis.enabled\" value to be provided")
	}

	if common.ContainsDuplicates(e.apis) {
		return fmt.Errorf("enabled apis contains duplicate: %+v", e.apis)
	}

	if !e.RestALTDA() && !e.ArbCustomDA() {
		return fmt.Errorf("an `arb` or REST ALT DA Server api type must be provided to start application")
	}

	if e.RestALTDAWithAdmin() &&
		(!e.RestALTDAOPGeneric() && !e.RestALTDAOPKeccak() && !e.RestALTStandard()) {
		return fmt.Errorf("admin mode for REST ALTDA server cannot be enabled without also " +
			"setting one of `op-generic`, `op-keccak`, `standard`")
	}

	return nil
}

func StringsToEnabledAPIs(strSlice []string) (*EnabledAPIs, error) {
	enabledAPIs := EnabledAPIs{
		make([]API, len(strSlice)),
	}

	for i, apiStr := range strSlice {
		enabledAPI, err := APIFromString(apiStr)
		if err != nil {
			return nil, fmt.Errorf("could not read string into API enum type: %w", err)
		}

		enabledAPIs.apis[i] = enabledAPI
	}

	return &enabledAPIs, nil
}

// New processes a string slice passed from user CLI
// into an API enum set
func New(apis []API) *EnabledAPIs {
	return &EnabledAPIs{
		apis,
	}
}

func (e EnabledAPIs) has(api API) bool {
	return slices.Contains(e.apis, api)
}

// API represents the different APIs that can be exposed on the proxy application
type API string

const (
	Admin               API = "admin"
	OpKeccakCommitment  API = "op-generic"
	OpGenericCommitment API = "op-keccak"
	StandardCommitment  API = "standard"
	ArbCustomDAServer   API = "arb"
	MetricsServer       API = "metrics"
)

func AllRestAPIs() []API {
	return []API{
		Admin, OpGenericCommitment, OpKeccakCommitment, StandardCommitment,
	}
}

func APIFromString(s string) (API, error) {
	// case insensitive
	s = strings.ToLower(s)

	switch s {
	case "admin":
		return Admin, nil
	case "op-generic":
		return OpGenericCommitment, nil
	case "op-keccak":
		return OpKeccakCommitment, nil
	case "standard":
		return StandardCommitment, nil
	case "arb":
		return ArbCustomDAServer, nil
	case "metrics":
		return MetricsServer, nil
	default:
		return "", fmt.Errorf("unknown API string: %s", s)
	}
}
