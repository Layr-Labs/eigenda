# Encoding Benchmark Suite

This testing package holds various benchmarks related to operations performed for encoding
that are important to make the entire EigenDA network fast. The benchmarks are separated
into high-level and low-level operations.

## High-Level Operations

`benchmark_eigenda_test.go` contains benchmarks for the high-level math/crypto operations that are
performed by different actors of the EigenDA network:
- Clients: PayloadToBlob conversion, Commitment generation
- Dispersers: Frame generation (RS encoding into chunks + KZG multiproof generation)
- Validators: Verification of commitments and proofs (TODO: write benchmark for this)

## Low-Level Operations

`benchmark_primitives_test.go` contains benchmarks for the typical crypto primitives: FFTFr, FFTG1, MSMG1/G2.
Speeding up any of the primitives leads to speedups in the higher level operations.

### GPU

`benchmark_icicle_test.go` contains benchmarks to test GPU implementations of the primitives using the icicle library.