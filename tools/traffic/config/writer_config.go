package config

import "time"

// BlobWriterConfig configures the blob writer.
type BlobWriterConfig struct {
	// The number of worker threads that generate write traffic.
	NumWriteInstances uint

	// The period of the submission rate of new blobs for each write worker thread.
	WriteRequestInterval time.Duration

	// The Size of each blob dispersed, in bytes.
	DataSize uint64

	// If true, then each blob will contain unique random data. If false, the same random data
	// will be dispersed for each blob by a particular worker thread.
	RandomizeBlobs bool

	// The amount of time to wait for a blob to be written.
	WriteTimeout time.Duration

	// Custom quorum numbers to use for the traffic generator.
	CustomQuorums []uint8
}
