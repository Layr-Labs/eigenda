package metrics

const (
	namespace = "eigenda"
)

var (
	// Buckets for payload and blob size measurements
	// Starting from 0 up to 16MiB
	blobSizeBuckets = []float64{
		0,
		131072,   // 128KiB
		262144,   // 256KiB
		524288,   // 512KiB
		1048576,  // 1MiB
		2097152,  // 2MiB
		4194304,  // 4MiB
		8388608,  // 8MiB
		16777216, // 16MiB
	}
)
