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
    struct RelayInfo {
        address relayAddress;
        string relayURL;
    }
}
