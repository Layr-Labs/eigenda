#![allow(missing_docs)]

use criterion::{Criterion, black_box, criterion_group, criterion_main};
use eigenda_verification::verification::cert::{self};

fn cert_verification_bench(c: &mut Criterion) {
    let inputs = cert::test_utils::success_inputs();

    let mut group = c.benchmark_group("cert_verification");
    group.sample_size(1000);
    group.measurement_time(std::time::Duration::from_secs(10));

    group.bench_function("verify", |b| {
        b.iter(|| {
            let inputs_clone = black_box(inputs.clone());
            cert::verify(inputs_clone)
                .expect("Certificate verification should succeed with test data")
        })
    });

    group.finish();
}

criterion_group!(benches, cert_verification_bench);
criterion_main!(benches);
