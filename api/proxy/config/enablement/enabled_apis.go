package enablement

import (
	"fmt"
	"strings"

	"github.com/Layr-Labs/eigenda/api/proxy/common"
)

// EnabledServersConfig is the highest level of code path dictation for
// a proxy application instance.
type EnabledServersConfig struct {
	Metric      bool
	ArbCustomDA bool

	RestAPIConfig RestApisEnabled
}

// RestApisEnabled stores boolean fields that dictate which
// commitment modes and routes to support.
// TODO: Add support for a `read-only` mode
type RestApisEnabled struct {
	Admin               bool
	OpGenericCommitment bool
	OpKeccakCommitment  bool
	StandardCommitment  bool
}

func (e *RestApisEnabled) Enabled() bool {
	return e.OpGenericCommitment ||
		e.OpKeccakCommitment || e.StandardCommitment
}

// Check ... Ensures that expression of the enabled API set is correct
func (e EnabledServersConfig) Check() error {
	if !e.RestAPIConfig.Enabled() && !e.ArbCustomDA {
		return fmt.Errorf("an `arb` or REST ALT DA Server api type must be provided to start application")
	}

	return nil
}

// APIStringsToEnabledServersConfig takes a dynamic array of strings provided from user CLI
// input and converts them into a high level enablement config
func APIStringsToEnabledServersConfig(strSlice []string) (*EnabledServersConfig, error) {
	if len(strSlice) == 0 {
		return nil, fmt.Errorf("cannot provide empty values for `apis.enabled`")
	}

	apis := make([]API, 0)

	for _, apiStr := range strSlice {
		enabledAPI, err := APIFromString(apiStr)
		if err != nil {
			return nil, fmt.Errorf("could not read string into API enum type: %w", err)
		}

		// no duplicate entries allowed
		if common.Contains(apis, enabledAPI) {
			return nil, fmt.Errorf("string api type already provided: %s", enabledAPI)
		}

		apis = append(apis, enabledAPI)
	}

	return &EnabledServersConfig{
		Metric:      common.Contains(apis, MetricsServer),
		ArbCustomDA: common.Contains(apis, ArbCustomDAServer),
		RestAPIConfig: RestApisEnabled{
			Admin:               common.Contains(apis, Admin),
			OpGenericCommitment: common.Contains(apis, OpGenericCommitment),
			OpKeccakCommitment:  common.Contains(apis, OpKeccakCommitment),
			StandardCommitment:  common.Contains(apis, StandardCommitment),
		},
	}, nil
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

func AllAPIsString() string {
	return fmt.Sprintf(
		"%s, %s, %s, %s, %s, %s", Admin, StandardCommitment,
		OpGenericCommitment, OpKeccakCommitment,
		ArbCustomDAServer, MetricsServer)
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
