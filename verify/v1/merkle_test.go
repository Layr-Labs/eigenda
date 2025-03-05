package verify

import (
	"testing"

	"github.com/wealdtech/go-merkletree/v2"
	"github.com/wealdtech/go-merkletree/v2/keccak256"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
)

func TestProcessInclusionProofPass(t *testing.T) {
	proof, err := hexutil.Decode("0xc455c1ea0e725d7ea3e5f29e9f48be8fc2787bb0a914d5a86710ba302c166ac4f626d76f67f1055bb960a514fb8923af2078fd84085d712655b58a19612e8cd15c3e4ac1cef57acde3438dbcf63f47c9fefe1221344c4d5c1a4943dd0d1803091ca81a270909dc0e146841441c9bd0e08e69ce6168181a3e4060ffacf3627480bec6abdd8d7bb92b49d33f180c42f49e041752aaded9c403db3a17b85e48a11e9ea9a08763f7f383dab6d25236f1b77c12b4c49c5cdbcbea32554a604e3f1d2f466851cb43fe73617b3d01e665e4c019bf930f92dea7394c25ed6a1e200d051fb0c30a2193c459f1cfef00bf1ba6656510d16725a4d1dc031cb759dbc90bab427b0f60ddc6764681924dda848824605a4f08b7f526fe6bd4572458c94e83fbf2150f2eeb28d3011ec921996dc3e69efa52d5fcf3182b20b56b5857a926aa66605808079b4d52c0c0cfe06923fa92e65eeca2c3e6126108e8c1babf5ac522f4d7")
	require.NoError(t, err)

	leaf := common.HexToHash("0xf6106e6ae4631e68abe0fa898cedbe97dbae6c7efb1b088c5aa2e8b91190ff96")
	index := uint64(580)

	expectedRoot, err := hexutil.Decode("0x7390b8023db8248123dcaeca57fa6c9340bef639e204f2278fc7ec3d46ad071b")
	require.NoError(t, err)

	actualRoot, err := ProcessInclusionProof(proof, leaf, index)
	require.NoError(t, err)

	require.Equal(t, expectedRoot, actualRoot.Bytes())
}

func TestProcessInclusionProofFail(t *testing.T) {
	proof, err := hexutil.Decode("0xc455c1ea0e725d7ea3e5f29e9f48be8fc2787bb0a914d5a86710ba302c166ac4f626d76f67f1055bb960a514fb8923af2078fd84085d712655b58a19612e8cd15c3e4ac1cef57acde3438dbcf63f47c9fefe1221344c4d5c1a4943dd0d1803091ca81a270909dc0e146841441c9bd0e08e69ce6168181a3e4060ffacf3627480bec6abdd8d7bb92b49d33f180c42f49e041752aaded9c403db3a17b85e48a11e9ea9a08763f7f383dab6d25236f1b77c12b4c49c5cdbcbea32554a604e3f1d2f466851cb43fe73617b3d01e665e4c019bf930f92dea7394c25ed6a1e200d051fb0c30a2193c459f1cfef00bf1ba6656510d16725a4d1dc031cb759dbc90bab427b0f60ddc6764681924dda848824605a4f08b7f526fe6bd4572458c94e83fbf2150f2eeb28d3011ec921996dc3e69efa52d5fcf3182b20b56b5857a926aa66605808079b4d52c0c0cfe06923fa92e65eeca2c3e6126108e8c1babf5ac522f4d7")
	require.NoError(t, err)

	leaf := common.HexToHash("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	index := uint64(580)

	expectedRoot, err := hexutil.Decode("0x7390b8023db8248123dcaeca57fa6c9340bef639e204f2278fc7ec3d46ad071b")
	require.NoError(t, err)

	actualRoot, err := ProcessInclusionProof(proof, leaf, index)
	require.NoError(t, err)

	require.NotEqual(t, expectedRoot, actualRoot.Bytes())
}

// TestProcessInclusionProofSingleNode confirms that a merkle tree containing a single node is successfully confirmed
func TestProcessInclusionProofSingleNode(t *testing.T) {
	leaf, err := hexutil.Decode("0x616C6C206861696C20746865206772656174207361746F736869")
	require.NotNil(t, leaf)
	require.NoError(t, err)

	tree, err := merkletree.NewTree(merkletree.WithData([][]byte{leaf}), merkletree.WithHashType(keccak256.New()))
	require.NotNil(t, tree)
	require.NoError(t, err)

	merkleProof, err := tree.GenerateProofWithIndex(0, 0)
	require.NotNil(t, merkleProof)
	require.NoError(t, err)

	// sanity check: there shouldn't be any sibling hashes for this tree
	require.Equal(t, 0, len(merkleProof.Hashes))

	emptyProof := make([]byte, 0)

	computedRoot, err := ProcessInclusionProof(
		emptyProof,
		common.BytesToHash(keccak256.New().Hash(leaf)),
		0)
	require.NotNil(t, computedRoot)
	require.NoError(t, err)
	require.Equal(t, computedRoot.Bytes(), tree.Root())

	// create an alternate leaf, and make sure that the inclusion proof fails the comparison check
	badLeaf, err := hexutil.Decode("0xab")
	require.NotNil(t, badLeaf)
	require.NoError(t, err)

	computedRoot, err = ProcessInclusionProof(
		emptyProof,
		common.BytesToHash(keccak256.New().Hash(badLeaf)),
		0)
	require.NotNil(t, computedRoot)
	require.NoError(t, err)
	require.NotEqual(t, computedRoot.Bytes(), tree.Root())
}
