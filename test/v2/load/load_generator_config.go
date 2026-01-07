package load

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/Layr-Labs/eigenda/test/v2/client"
)

var _ config.DocumentedConfig = (*TrafficGeneratorConfig)(nil)

// Configuration for the traffic generator.
//
// TODO(cody.littley): This parent struct is not currently used for deploying a traffic generator,
// but that will soon change. When the change is made, I will also do some renaming to make things cleaner.
type TrafficGeneratorConfig struct {
	// Configures the environment towards which the traffic generator will run.
	Environment client.TestClientConfig
	// Configures the load the traffic generator will produce.
	Load LoadGeneratorConfig
}

// DefaultTrafficGeneratorConfig returns a default configuration for the traffic generator.
func DefaultTrafficGeneratorConfig() *TrafficGeneratorConfig {
	return &TrafficGeneratorConfig{
		Environment: *client.DefaultTestClientConfig(),
		Load:        *DefaultLoadGeneratorConfig(),
	}
}

var _ config.VerifiableConfig = (*LoadGeneratorConfig)(nil)

// LoadGeneratorConfig is the configuration for the load generator.
type LoadGeneratorConfig struct {
	// The desired number of megabytes bytes per second to write.
	MbPerSecond float64
	// The size of the blobs to write, in megabytes.
	BlobSizeMb float64
	// By default, this utility reads each blob back from each relay once. The number of
	// reads per relay is multiplied by this factor. For example, If this is set to 3,
	// then each blob is read back from each relay 3 times. If less than 1, then this value
	// is treated as a probability. For example, if this is set to 0.5, then each blob is read back
	// from each relay with a 50% chance. If running with the proxy, this value is used to determine
	// how many times to read each blob back from the proxy (since in the normal case, proxy reads translate
	// to relay reads).
	RelayReadAmplification float64
	// By default, this utility reads chunks once. The number of chunk reads is multiplied
	// by this factor. If this is set to 3, then chunks are read back 3 times. If less than 1,
	// then this value is treated as a probability. For example, if this is set to 0.5, then
	// each chunk is read back from validators with a 50% chance. Ignored if the load generator is configured
	// to use the proxy.
	ValidatorReadAmplification float64
	// A number between 0 and 1.0 that specifies the fraction of blobs that are verified by the validator.
	// If 1.0, all blobs are verified. If 0.0, no blobs are verified. If 0.5, half of the blobs are verified.
	ValidatorVerificationFraction float64
	// The maximum number of parallel blobs submissions in flight.
	SubmissionParallelism uint64
	// The maximum number of parallel blob relay read operations in flight.
	RelayReadParallelism uint64
	// The maximum number of parallel blob validator read operations in flight.
	ValidatorReadParallelism uint64
	// The maximum number of parallel gas estimation operations in flight.
	GasEstimationParallelism uint64
	// The timeout for each blob dispersal, in seconds.
	DispersalTimeout uint32
	// The timeout for reading a blob from a relay, in seconds. This is the timeout per individual read.
	RelayReadTimeout uint32
	// The timeout for reading a blob from the validators, in seconds. This is the timeout per individual read.
	ValidatorReadTimeout uint32
	// The timeout for gas estimation operations, in seconds.
	GasEstimationTimeout uint32
	// EnablePprof enables the pprof HTTP server for profiling
	EnablePprof bool
	// PprofHttpPort is the port that the pprof HTTP server listens on
	PprofHttpPort int
	// FrequencyAcceleration determines the speed at which the frequency of blob submissions accelerates at startup
	// time, in HZ/s. Frequency will start at 0 and accelerate to the target frequency at this rate. If 0, then
	// the frequency will immediately be set to the target frequency.
	FrequencyAcceleration float64
	// If true, then route traffic through the proxy instead of directly using the GRPC clients.
	UseProxy bool
}

// DefaultLoadGeneratorConfig returns a default configuration for the load generator.
func DefaultLoadGeneratorConfig() *LoadGeneratorConfig {
	return &LoadGeneratorConfig{
		MbPerSecond:                   0.5,
		BlobSizeMb:                    2.0,
		RelayReadAmplification:        1.0,
		ValidatorReadAmplification:    1.0,
		ValidatorVerificationFraction: 0.01,
		SubmissionParallelism:         300,
		RelayReadParallelism:          300,
		ValidatorReadParallelism:      300,
		GasEstimationParallelism:      300,
		DispersalTimeout:              600,
		RelayReadTimeout:              600,
		ValidatorReadTimeout:          600,
		GasEstimationTimeout:          15,
		EnablePprof:                   false,
		PprofHttpPort:                 6060,
		FrequencyAcceleration:         0.0025,
		UseProxy:                      false,
	}
}

func (c *TrafficGeneratorConfig) GetEnvVarPrefix() string {
	return "TRAFFIC_GENERATOR"
}

func (c *TrafficGeneratorConfig) GetName() string {
	return "TrafficGenerator"
}

func (c *TrafficGeneratorConfig) GetPackagePaths() []string {
	return []string{
		"github.com/Layr-Labs/eigenda/test/v2/client",
		"github.com/Layr-Labs/eigenda/test/v2/load",
	}
}

func (l *LoadGeneratorConfig) Verify() error {
	if l.MbPerSecond <= 0 {
		return fmt.Errorf("MbPerSecond must be greater than 0")
	}
	if l.BlobSizeMb <= 0 {
		return fmt.Errorf("BlobSizeMb must be greater than 0")
	}
	if l.RelayReadAmplification < 0 {
		return fmt.Errorf("RelayReadAmplification must be non-negative")
	}
	if l.ValidatorReadAmplification < 0 {
		return fmt.Errorf("ValidatorReadAmplification must be non-negative")
	}
	if l.ValidatorVerificationFraction < 0 || l.ValidatorVerificationFraction > 1.0 {
		return fmt.Errorf("ValidatorVerificationFraction must be between 0 and 1.0")
	}
	if l.SubmissionParallelism == 0 {
		return fmt.Errorf("SubmissionParallelism must be greater than 0")
	}
	if l.RelayReadParallelism == 0 {
		return fmt.Errorf("RelayReadParallelism must be greater than 0")
	}
	if l.ValidatorReadParallelism == 0 {
		return fmt.Errorf("ValidatorReadParallelism must be greater than 0")
	}
	if l.GasEstimationParallelism == 0 {
		return fmt.Errorf("GasEstimationParallelism must be greater than 0")
	}
	if l.DispersalTimeout == 0 {
		return fmt.Errorf("DispersalTimeout must be greater than 0")
	}
	if l.RelayReadTimeout == 0 {
		return fmt.Errorf("RelayReadTimeout must be greater than 0")
	}
	if l.ValidatorReadTimeout == 0 {
		return fmt.Errorf("ValidatorReadTimeout must be greater than 0")
	}
	if l.GasEstimationTimeout == 0 {
		return fmt.Errorf("GasEstimationTimeout must be greater than 0")
	}
	if l.EnablePprof && (l.PprofHttpPort <= 0 || l.PprofHttpPort > 65535) {
		return fmt.Errorf("PprofHttpPort must be a valid port number when EnablePprof is true")
	}
	if l.FrequencyAcceleration < 0 {
		return fmt.Errorf("FrequencyAcceleration must be non-negative")
	}
	return nil
}

func (c *TrafficGeneratorConfig) Verify() error {
	err := c.Load.Verify()
	if err != nil {
		return fmt.Errorf("load generator config verification failed: %w", err)
	}
	err = c.Environment.Verify()
	if err != nil {
		return fmt.Errorf("environment config verification failed: %w", err)
	}
	return nil
}
