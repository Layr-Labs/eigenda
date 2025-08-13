use alloy_sol_types::sol;

sol! {
    #[derive(Default, Debug)]
    struct VersionedBlobParams {
        uint32 maxNumOperators;
        uint32 numChunks;
        uint8 codingRate;
    }

    #[derive(Default, Debug)]
    struct SecurityThresholds {
        uint8 confirmationThreshold;
        uint8 adversaryThreshold;
    }

    #[derive(Default, Debug)]
    struct QuorumBitmapUpdate {
        uint32 updateBlockNumber;
        uint32 nextUpdateBlockNumber;
        uint192 quorumBitmap;
    }

    #[derive(Default, Debug)]
    struct ApkUpdate {
        bytes24 apkHash;
        uint32 updateBlockNumber;
        uint32 nextUpdateBlockNumber;
    }

    #[derive(Default, Debug)]
    struct StakeUpdate {
        uint32 updateBlockNumber;
        uint32 nextUpdateBlockNumber;
        uint96 stake;
    }
}
