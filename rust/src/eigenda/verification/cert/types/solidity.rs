use alloy_sol_types::sol;

sol! {
    /// Parameters defining blob size and encoding constraints for a specific version.
    ///
    /// These parameters control the operational limits for data blobs at different
    /// protocol versions, ensuring proper encoding and operator capacity constraints.
    #[derive(Default, Debug)]
    struct VersionedBlobParams {
        /// Maximum number of operators that can participate in this blob version
        uint32 maxNumOperators;
        /// Number of data chunks the blob is divided into for encoding
        uint32 numChunks;
        /// Coding rate used for erasure coding (affects redundancy level)
        uint8 codingRate;
    }

    /// Security thresholds defining minimum requirements for certificate validity.
    ///
    /// These thresholds determine the minimum stake percentages required for
    /// valid certificate signatures in the EigenDA protocol.
    #[derive(Default, Debug)]
    struct SecurityThresholds {
        /// Minimum percentage of stake required to confirm a certificate
        uint8 confirmationThreshold;
        /// Maximum percentage of adversarial stake that can be tolerated
        uint8 adversaryThreshold;
    }

    /// Historical update entry for quorum membership bitmaps.
    ///
    /// Tracks changes in an operator's quorum membership over time,
    /// allowing verification of which quorums an operator belonged to
    /// at any given block number.
    #[derive(Default, Debug)]
    struct QuorumBitmapUpdate {
        /// Block number when this membership update became active
        uint32 updateBlockNumber;
        /// Block number when this update was superseded (0 if current)
        uint32 nextUpdateBlockNumber;
        /// Bitmap indicating which quorums the operator belongs to
        uint192 quorumBitmap;
    }

    /// Historical update entry for aggregate public key hashes.
    ///
    /// Tracks changes to quorum aggregate public keys over time,
    /// enabling verification of the correct APK at any historical block.
    #[derive(Default, Debug)]
    struct ApkUpdate {
        /// Truncated hash of the aggregate public key (24 bytes)
        bytes24 apkHash;
        /// Block number when this APK update became active
        uint32 updateBlockNumber;
        /// Block number when this update was superseded (0 if current)
        uint32 nextUpdateBlockNumber;
    }

    /// Historical update entry for operator stake amounts.
    ///
    /// Tracks changes in an operator's stake over time within a specific quorum,
    /// allowing verification of operator voting power at any historical point.
    #[derive(Default, Debug)]
    struct StakeUpdate {
        /// Block number when this stake update became active
        uint32 updateBlockNumber;
        /// Block number when this update was superseded (0 if current)
        uint32 nextUpdateBlockNumber;
        /// Stake amount in the quorum's denomination (96-bit precision)
        uint96 stake;
    }
}
