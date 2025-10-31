package enablement_test

import (
	"testing"

	"github.com/Layr-Labs/eigenda/api/proxy/config/enablement"
	"github.com/stretchr/testify/assert"
)

func TestToAPIStrings(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		config   enablement.EnabledServersConfig
		expected []string
	}{
		{
			name: "All APIs enabled",
			config: enablement.EnabledServersConfig{
				Metric:      true,
				ArbCustomDA: true,
				RestAPIConfig: enablement.RestApisEnabled{
					Admin:               true,
					OpGenericCommitment: true,
					OpKeccakCommitment:  true,
					StandardCommitment:  true,
				},
			},
			expected: []string{"metrics", "arb", "admin", "op-generic", "op-keccak", "standard"},
		},
		{
			name: "No APIs enabled",
			config: enablement.EnabledServersConfig{
				Metric:      false,
				ArbCustomDA: false,
				RestAPIConfig: enablement.RestApisEnabled{
					Admin:               false,
					OpGenericCommitment: false,
					OpKeccakCommitment:  false,
					StandardCommitment:  false,
				},
			},
			expected: []string{},
		},
		{
			name: "Only Metric enabled",
			config: enablement.EnabledServersConfig{
				Metric:      true,
				ArbCustomDA: false,
				RestAPIConfig: enablement.RestApisEnabled{
					Admin:               false,
					OpGenericCommitment: false,
					OpKeccakCommitment:  false,
					StandardCommitment:  false,
				},
			},
			expected: []string{"metrics"},
		},
		{
			name: "Only ArbCustomDA enabled",
			config: enablement.EnabledServersConfig{
				Metric:      false,
				ArbCustomDA: true,
				RestAPIConfig: enablement.RestApisEnabled{
					Admin:               false,
					OpGenericCommitment: false,
					OpKeccakCommitment:  false,
					StandardCommitment:  false,
				},
			},
			expected: []string{"arb"},
		},
		{
			name: "Only REST APIs enabled",
			config: enablement.EnabledServersConfig{
				Metric:      false,
				ArbCustomDA: false,
				RestAPIConfig: enablement.RestApisEnabled{
					Admin:               true,
					OpGenericCommitment: true,
					OpKeccakCommitment:  true,
					StandardCommitment:  true,
				},
			},
			expected: []string{"admin", "op-generic", "op-keccak", "standard"},
		},
		{
			name: "Mixed configuration - Metric and some REST APIs",
			config: enablement.EnabledServersConfig{
				Metric:      true,
				ArbCustomDA: false,
				RestAPIConfig: enablement.RestApisEnabled{
					Admin:               true,
					OpGenericCommitment: false,
					OpKeccakCommitment:  true,
					StandardCommitment:  false,
				},
			},
			expected: []string{"metrics", "admin", "op-keccak"},
		},
		{
			name: "Mixed configuration - ArbCustomDA and some REST APIs",
			config: enablement.EnabledServersConfig{
				Metric:      false,
				ArbCustomDA: true,
				RestAPIConfig: enablement.RestApisEnabled{
					Admin:               false,
					OpGenericCommitment: true,
					OpKeccakCommitment:  false,
					StandardCommitment:  true,
				},
			},
			expected: []string{"arb", "op-generic", "standard"},
		},
		{
			name: "Only Admin enabled",
			config: enablement.EnabledServersConfig{
				Metric:      false,
				ArbCustomDA: false,
				RestAPIConfig: enablement.RestApisEnabled{
					Admin:               true,
					OpGenericCommitment: false,
					OpKeccakCommitment:  false,
					StandardCommitment:  false,
				},
			},
			expected: []string{"admin"},
		},
		{
			name: "Only OpGenericCommitment enabled",
			config: enablement.EnabledServersConfig{
				Metric:      false,
				ArbCustomDA: false,
				RestAPIConfig: enablement.RestApisEnabled{
					Admin:               false,
					OpGenericCommitment: true,
					OpKeccakCommitment:  false,
					StandardCommitment:  false,
				},
			},
			expected: []string{"op-generic"},
		},
		{
			name: "Only OpKeccakCommitment enabled",
			config: enablement.EnabledServersConfig{
				Metric:      false,
				ArbCustomDA: false,
				RestAPIConfig: enablement.RestApisEnabled{
					Admin:               false,
					OpGenericCommitment: false,
					OpKeccakCommitment:  true,
					StandardCommitment:  false,
				},
			},
			expected: []string{"op-keccak"},
		},
		{
			name: "Only StandardCommitment enabled",
			config: enablement.EnabledServersConfig{
				Metric:      false,
				ArbCustomDA: false,
				RestAPIConfig: enablement.RestApisEnabled{
					Admin:               false,
					OpGenericCommitment: false,
					OpKeccakCommitment:  false,
					StandardCommitment:  true,
				},
			},
			expected: []string{"standard"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.config.ToAPIStrings()
			assert.Equal(t, tc.expected, got)
		})
	}
}
