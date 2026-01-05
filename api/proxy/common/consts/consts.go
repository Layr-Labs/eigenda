package consts

// EthHappyPathFinalizationDepth is the number of blocks that must be included on top of a block for it to be considered
// "final",
// under happy-path aka normal network conditions.
//
// See https://www.alchemy.com/overviews/ethereum-commitment-levels for a quick TLDR explanation,
// or https://eth2book.info/capella/part3/transition/epoch/#finalisation for full details.
var EthHappyPathFinalizationDepthBlocks = uint8(64)

// RBNRecencyWindowSizeV0 is the recency window size in L1 blocks for V4+ certs with a derivation version
// of 0. This value is used in the RBN recency check to determine if a certificate is too old
// compared to the L1 inclusion block number provided by the client. The value of 14400 represents 48 hours
// worth of blocks at an average block time of 12 seconds.
//
// See https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#1-rbn-recency-validation
var RBNRecencyWindowSizeV0 uint64 = 14400
