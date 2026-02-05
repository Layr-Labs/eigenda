package version

// The semantic defaultVersion string of the code in this branch. Sometimes a more specific version may be provided
// by the build toolchain.
var defaultVersion = NewSemver(2, 5, 2, "")

// Get the default version of the code in this branch.
func DefaultVersion() *Semver {
	return defaultVersion
}
