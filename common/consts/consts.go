package consts

// EthHappyPathFinalizationDepth is the number of blocks that must be included on top of a block for it to be considered "final",
// under happy-path aka normal network conditions.
//
// See https://www.alchemy.com/overviews/ethereum-commitment-levels for a quick TLDR explanation,
// or https://eth2book.info/capella/part3/transition/epoch/#finalisation for full details.
var EthHappyPathFinalizationDepthBlocks = uint8(64)
