#![cfg(feature = "native")]

mod common;

use std::{
    collections::HashMap,
    hash::{DefaultHasher, Hash, Hasher},
};

use anyhow::Ok;
use bytes::Bytes;
use rand::{Rng, RngCore};
use sov_eigenda_adapter::{
    service::EigenDaService,
    spec::{BlobWithSender, EthereumBlockHeader},
    verifier::{EigenDaCompletenessProof, EigenDaInclusionProof, EigenDaVerifier},
};
use sov_rollup_interface::{
    common::HexHash,
    da::{BlobReaderTrait, DaVerifier, RelevantBlobs, RelevantProofs},
    node::da::DaService,
};
use tracing::info;

use crate::common::{proxy::start_proxy, setup_adapter};

#[tokio::test]
async fn submit_extract_verify_e2e() {
    common::tracing::init_tracing();

    let (proxy_url, _proxy_container) = start_proxy().await.unwrap();
    let (service, verifier) = setup_adapter(proxy_url).await.unwrap();

    let mut rng = rand::thread_rng();
    let blobs_size_range = 1024..2048;

    let blobs = (0..5)
        .map(|_| {
            let size = rng.gen_range(blobs_size_range.clone());
            let mut blob = vec![0u8; size];
            rng.fill_bytes(&mut blob);
            Bytes::from(blob)
        })
        .collect::<Vec<_>>();

    let blobs_cumulative_size = blobs.iter().map(|b| b.len()).sum::<usize>();
    info!(
        "Going to submit {} blobs of total size {} bytes",
        blobs.len(),
        blobs_cumulative_size
    );

    check_blobs_roundtrip(&service, &verifier, &blobs)
        .await
        .unwrap();
}

async fn check_blobs_roundtrip(
    service: &EigenDaService,
    verifier: &EigenDaVerifier,
    blobs: &[Bytes],
) -> anyhow::Result<()> {
    let mut sent_blobs: HashMap<u64, HexHash> = HashMap::with_capacity(blobs.len());
    let sender = service.get_signer().await;

    // Height at which we will start looking for blobs
    let start_height = service.get_head_block_header().await?.as_ref().number;
    // The height at which we will stop looking for blobs
    let end_height = start_height + 10;

    // Submit blobs
    for blob in blobs {
        let receipt = service.send_transaction(blob).await.await??;
        let bytes_hash = hash_bytes(blob);

        info!(
            "Sent blob of size {} bytes hash {} blob hash {}",
            blob.len(),
            bytes_hash,
            receipt.blob_hash
        );

        if sent_blobs.insert(bytes_hash, receipt.blob_hash).is_some() {
            anyhow::bail!(
                "Non unique blob {} size={}, test cannot be performed",
                bytes_hash,
                blob.len()
            );
        }
    }

    // Look for blobs in the blocks and process them
    for height in start_height..=end_height {
        let block = service.get_block_at(height).await?;

        let (mut blobs, proofs) = service.extract_relevant_blobs_with_proof(&block).await;
        info!(
            "Inspecting block {:?}. batch blobs: {} proof blobs: {}",
            block.header,
            blobs.batch_blobs.len(),
            blobs.proof_blobs.len(),
        );

        // Simulate processing of the blobs
        rollup_process_relevant_blobs(&mut blobs);

        // Verify relevant blobs against proofs
        verify_relevant_blobs(&verifier, &block.header, &blobs, &proofs)?;

        for batch in blobs.batch_blobs {
            if batch.sender() != sender {
                continue;
            }

            info!("Received batch hash: {}, height={}", batch.hash(), height);
            let data = batch.verified_data();
            let bytes_hash = hash_bytes(data);

            match sent_blobs.remove(&bytes_hash) {
                None => {
                    anyhow::bail!(
                        "Received blob on height={} ethereum_hash={} bytes_hash={} len={} not found in sent blobs. Bug or there's another sender conflicting with the test.",
                        block.header.as_ref().number,
                        batch.hash(),
                        bytes_hash,
                        data.len(),
                    );
                }
                Some(sent_blob_hash) => {
                    if sent_blob_hash.0 != *batch.hash() {
                        anyhow::bail!("Blob hashes do not match for the same blob data");
                    }
                }
            }
        }

        if sent_blobs.is_empty() {
            return Ok(());
        }
    }

    if !sent_blobs.is_empty() {
        info!("Remaining blobs: {:?}", sent_blobs);
        anyhow::bail!("Error: {} blobs were not received", sent_blobs.len());
    }

    Ok(())
}

fn hash_bytes(bytes: &[u8]) -> u64 {
    let mut hasher = DefaultHasher::new();
    bytes.hash(&mut hasher);
    hasher.finish()
}

fn rollup_process_relevant_blobs(relevant_blobs: &mut RelevantBlobs<BlobWithSender>) {
    let read_full = |blob_with_sender: &mut BlobWithSender| {
        let total_len = blob_with_sender.blob.total_len();
        blob_with_sender.blob.advance(total_len);
        let data = blob_with_sender.blob.accumulator();
        assert_eq!(data.len(), total_len);
    };

    let blob_iters = relevant_blobs.as_iters();
    for batch in blob_iters.batch_blobs {
        read_full(batch);
    }
    for proof in blob_iters.proof_blobs {
        read_full(proof);
    }
}

fn verify_relevant_blobs(
    verifier: &EigenDaVerifier,
    header: &EthereumBlockHeader,
    blobs: &RelevantBlobs<BlobWithSender>,
    proofs: &RelevantProofs<EigenDaInclusionProof, EigenDaCompletenessProof>,
) -> anyhow::Result<()> {
    // Simulate we're sending this to zkvm
    let header = risc0_zkvm::serde::to_vec(header)?;
    let blobs = risc0_zkvm::serde::to_vec(blobs)?;
    let proofs = risc0_zkvm::serde::to_vec(proofs)?;

    // Receive on zkvm side
    let header = risc0_zkvm::serde::from_slice(&header)?;
    let blobs = risc0_zkvm::serde::from_slice(&blobs)?;
    let proofs = risc0_zkvm::serde::from_slice(&proofs)?;

    // Verify
    verifier.verify_relevant_tx_list(&header, &blobs, proofs)?;

    Ok(())
}
