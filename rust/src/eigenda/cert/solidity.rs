use alloy_sol_types::sol;

sol! {
    #[derive(Debug)]
    struct BatchHeaderV2 {
        bytes32 batchRoot;
        uint32 referenceBlockNumber;
    }

    #[derive(Debug)]
    struct G1Point {
        uint256 X;
        uint256 Y;
    }

    // Encoding of field elements is: X[1] * i + X[0]
    #[derive(Debug)]
    struct G2Point {
        uint256[2] X;
        uint256[2] Y;
    }

    #[derive(Debug)]
    struct BlobCommitment {
        G1Point commitment;
        G2Point lengthCommitment;
        G2Point lengthProof;
        uint32 length;
    }
}
