package config

import "time"

// Config configures the traffic generator workers.
type WorkerConfig struct {
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

	// The amount of time between attempts by the verifier to confirm the status of blobs.
	VerifierInterval time.Duration
	// The amount of time to wait for a blob status to be fetched.
	GetBlobStatusTimeout time.Duration
	// The size of the channel used to communicate between the writer and verifier.
	VerificationChannelCapacity uint

	// The number of worker threads that generate read traffic.
	NumReadInstances uint
	// The period of the submission rate of read requests for each read worker thread.
	ReadRequestInterval time.Duration
	// For each blob, how many times should it be downloaded? If between 0.0 and 1.0, blob will be downloaded
	// 0 or 1 times with the specified probability (e.g. 0.2 means each blob has a 20% chance of being downloaded).
	// If greater than 1.0, then each blob will be downloaded the specified number of times.
	RequiredDownloads float64
	// The size of a table of blobs to optionally read when we run out of blobs that we are required to read. Blobs
	// that are no longer required are added to this table, and when the table is at capacity they are randomly retired.
	// Set this to 0 to disable this feature.
	ReadOverflowTableSize uint
	// The amount of time to wait for a batch header to be fetched.
	FetchBatchHeaderTimeout time.Duration
	// The amount of time to wait for a blob to be retrieved.
	RetrieveBlobChunksTimeout time.Duration

	// The address of the EigenDA service manager smart contract, in hex.
	EigenDAServiceManager string
	// The private key to use for signing requests.
	SignerPrivateKey string
	// Custom quorum numbers to use for the traffic generator.
	CustomQuorums []uint8
}
