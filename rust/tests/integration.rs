mod common;

use sov_eigenda_adapter::spec::BlobWithSender;
use sov_rollup_interface::{
    da::{BlobReaderTrait, DaVerifier, RelevantBlobs},
    node::da::DaService,
};

use crate::common::{proxy::start_proxy, setup_adapter};

#[tokio::test]
async fn submit_extract_verify_e2e() {
    let (proxy_url, _proxy_container) = start_proxy().await.unwrap();
    let (service, verifier) = setup_adapter(proxy_url).await.unwrap();

    // Block header before we publish the rollup related tx
    let mut latest_block_height = service
        .get_head_block_header()
        .await
        .unwrap()
        .as_ref()
        .number;

    // Post the rollup data to the network
    let mut rollup_batches = vec![vec![123; 1024], vec![50; 512], vec![1; 2048], vec![0; 256]];
    let mut rollup_proofs = vec![vec![200; 1536], vec![10; 1024], vec![3; 3072], vec![0; 128]];

    for blob in &rollup_batches {
        service.send_transaction(blob).await.await.unwrap().unwrap();
    }
    for proof in &rollup_proofs {
        service.send_proof(proof).await.await.unwrap().unwrap();
    }

    // Fail if blobs are not submitted in N blocks
    let limit_height = latest_block_height + 10;

    // Process blocks until we fetch all data
    loop {
        // Extract rollup data from the block
        let block = service.get_block_at(latest_block_height).await.unwrap();
        let mut blobs = service.extract_relevant_blobs(&block);
        let proofs = service.get_extraction_proof(&block, &blobs).await;

        // Simulate processing of data by the rollup
        rollup_process_relevant_blobs(&mut blobs);

        // Simulate we're sending this to zkvm
        let header = risc0_zkvm::serde::to_vec(&block.header).unwrap();
        let blobs = risc0_zkvm::serde::to_vec(&blobs).unwrap();
        let proofs = risc0_zkvm::serde::to_vec(&proofs).unwrap();

        // Receive on zkvm side
        let header = risc0_zkvm::serde::from_slice(&header).unwrap();
        let blobs = risc0_zkvm::serde::from_slice(&blobs).unwrap();
        let proofs = risc0_zkvm::serde::from_slice(&proofs).unwrap();

        // Verify
        verifier
            .verify_relevant_tx_list(&header, &blobs, proofs)
            .unwrap();

        // Remove verified blobs
        let RelevantBlobs {
            proof_blobs,
            batch_blobs,
        } = blobs;

        for mut proof_blob in proof_blobs {
            rollup_proofs.retain(|b| b != proof_blob.full_data());
        }
        for mut batch_blob in batch_blobs {
            rollup_batches.retain(|b| b != batch_blob.full_data());
        }

        // Success. All blobs were persisted and verified
        if rollup_batches.is_empty() && rollup_proofs.is_empty() {
            break;
        }

        if limit_height == latest_block_height {
            panic!("not all blobs were persisted to the chain");
        }

        // Process next block
        latest_block_height += 1;
    }
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
