package disperser

import "time"

const (
	Localhost = "0.0.0.0"
)

type ServerConfig struct {
	GrpcPort    string
	GrpcTimeout time.Duration

	PprofHttpPort string
	EnablePprof   bool
	
	// Version checking config
	OperatorVersionCheck         bool          // If true, enforce node version rollout check
	RequiredNodeVersion          string        // Required minimum node version (e.g. ">=0.9.0-rc.1")
	VersionStakeThreshold        float64       // Percentage of stake that needs to be running the required version (e.g. 0.8 for 80%)
	VersionCheckInterval         time.Duration // How often to check version info
}
