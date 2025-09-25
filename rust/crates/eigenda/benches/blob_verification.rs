#![allow(missing_docs)]

use std::sync::LazyLock;

use criterion::{Criterion, black_box, criterion_group, criterion_main};
use eigenda::verification::blob::{self, SRS};

fn blob_verification_bench(c: &mut Criterion) {
    LazyLock::force(&SRS);

    // testing with very large (but realistic) input
    let (blob_commitment, encoded_payload) = blob::success_inputs(&[123; 1_048_600]);
    blob::verify(&blob_commitment, &encoded_payload)
        .expect("Blob verification should succeed with test data (warm-up)");

    let mut group = c.benchmark_group("blob_verification");
    group.sample_size(10);
    group.measurement_time(std::time::Duration::from_secs(20));

    group.bench_function("verify", |b| {
        b.iter(|| {
            let blob_commitment = black_box(&blob_commitment);
            let encoded_payload = black_box(&encoded_payload);
            blob::verify(blob_commitment, encoded_payload)
                .expect("Blob verification should succeed with test data")
        })
    });

    group.finish();
}

criterion_group!(benches, blob_verification_bench);
criterion_main!(benches);
