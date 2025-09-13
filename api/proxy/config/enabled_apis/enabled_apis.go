package enabled_apis

import (
	"fmt"
	"slices"
	"strings"
)

// EnabledAPIs is a wrapper type of the APIs set enum and provides
// encapsulation for the entrypoint processing to make enablement determination
// through specific boolean accessors.
type EnabledAPIs struct {
	apis []API
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
	if e.Metrics() && (!e.RestALTDA() && !e.ArbCustomDA()) {
		return fmt.Errorf("metrics cannot be enabled unless `arb` and/or `rest` also is")
	}

	if e.RestALTDAWithAdmin() &&
		(!e.RestALTDAOPGeneric() && !e.RestALTDAOPKeccak() && !e.RestALTStandard()) {
		return fmt.Errorf("admin mode for REST ALTDA server cannot be enabled without also " +
			"setting one of `op-generic`, `op-keccak`, `standard`")
	}

	return nil
}

// NewEnabledAPIs processes a string slice passed from user CLI
// into an API enum set
func NewEnabledAPIs(strSlice []string) (*EnabledAPIs, error) {
	enabledAPIs := EnabledAPIs{
		make([]API, len(strSlice)),
	}

	for i, apiStr := range strSlice {
		enabledAPI, err := APIFromString(apiStr)
		if err != nil {
			return nil, fmt.Errorf("could not read string into API enum type: %w", err)
		}

		// SET data structure enforcement
		if enabledAPIs.has(enabledAPI) {
			return nil, fmt.Errorf("cannot pass the same API type more than once: %s", enabledAPI.ToString())
		}

		enabledAPIs.apis[i] = enabledAPI
	}

	return &enabledAPIs, nil
}

func (e EnabledAPIs) has(api API) bool {
	return slices.Contains(e.apis, api)
}

// API represents the different APIs that can be exposed on the proxy application
type API uint8

const (
	Admin               API = 1
	OpKeccakCommitment  API = 2
	OpGenericCommitment API = 3
	StandardCommitment  API = 4
	ArbCustomDAServer   API = 5
	MetricsServer       API = 6
)

func (api API) ToString() string {
	switch api {
	case Admin:
		return "admin"

	case OpGenericCommitment:
		return "op-generic"

	case OpKeccakCommitment:
		return "op-keccak"

	case StandardCommitment:
		return "standard"

	case ArbCustomDAServer:
		return "arb"

	case MetricsServer:
		return "metrics"
	default:
		return "unknown"
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
		return 0, fmt.Errorf("unknown API string: %s", s)
	}
}
