package node

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
