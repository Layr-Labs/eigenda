package node

import (
	"fmt"
	"log"

	"github.com/Layr-Labs/eigenda/common/version"
)

var (
	// Possibly set by go build -ldflags="-X 'github.com/Layr-Labs/eigenda/node.SemVer=${SEMVER}'
	// If not set, then the version defined in common/version will be used.
	// If not empty, then the default version defined in common/version will be overridden.
	SemVer = ""
	// Similar to SemVer, possibly set by go build -ldflags.
	GitCommit = ""
	// Similar to SemVer, possibly set by go build -ldflags.
	GitDate = ""
)

// Determine the software version, possibly using build-time variables.
func GetSoftwareVersion() *version.Semver {
	softwareVersion := version.DefaultVersion()

	if SemVer != "" {
		semver := SemVer
		if GitCommit != "" {
			semver = fmt.Sprintf("%s-%s", semver, GitCommit)
		}
		if GitDate != "" {
			semver = fmt.Sprintf("%s-%s", semver, GitDate)
		}

		var err error
		softwareVersion, err = version.SemverFromString(semver)
		if err != nil {
			log.Printf("Version string \"%s\" is invalid, falling back to hard coded version", SemVer)
			softwareVersion = version.DefaultVersion()
		}
	}

	return softwareVersion
}
