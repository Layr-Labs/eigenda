# KZG and FFT utils

This repo is *super experimental*.

This is an implementation in Go, initially aimed at chunkification and extension of data, 
and building/verifying KZG proofs for the output data.
The KZG proofs, or Kate proofs, are built on top of BLS12-381.

Part of a low-latency data-availability sampling network prototype for Eth2 Phase 1.
See https://github.com/protolambda/eth2-das

Code is based on:
- [KZG Data availability code by Dankrad](https://github.com/ethereum/research/tree/master/kzg_data_availability)
- [Verkle and FFT code by Dankrad and Vitalik](https://github.com/ethereum/research/tree/master/verkle)
- [Reed solomon erasure code recovery with FFTs by Vitalik](https://ethresear.ch/t/reed-solomon-erasure-code-recovery-in-n-log-2-n-time-with-ffts/3039)
- [FFT explainer by Vitalik](https://vitalik.ca/general/2019/05/12/fft.html)
- [Kate explainer by Dankrad](https://dankradfeist.de/ethereum/2020/06/16/kate-polynomial-commitments.html)
- [Kate amortized paper by Dankrad and Dmitry](https://github.com/khovratovich/Kate/blob/master/Kate_amortized.pdf)

Features:
- (I)FFT on `F_r`
- (I)FFT on `G1`
- Specialized FFT for extension of `F_r` data
- KZG
  - commitments
  - generate/verify proof for single point
  - generate/verify proofs for multiple points
  - generate/verify proofs for all points, using FK20
  - generate/verify proofs for ranges (cosets) of points, using FK20
- Data recovery: given an arbitrary subset of data (at least half), recover the rest
- Optimized for Data-availability usage
- Change Bignum / BLS with build tags.

## BLS

Currently supported BLS implementations: Herumi BLS and Kilic BLS (default).

## Field elements (Fr)

The BLS curve order is used for the modulo math, different libraries could be used to provide this functionality.
Note: some of these libraries do not have full BLS functionality, only Bignum / uint256. The KZG code will be excluded when compiling with a non-BLS build tag.

Build tag options:
- (no build tags, default): Use Kilic BLS library. Previously used by `bignum_kilic` build tag. [`kilic/bls12-381`](https://github.com/kilic/bls12-381)
- `-tags bignum_hbls`: use Herumi BLS library. [`herumi/bls-eth-go-binary`](https://github.com/herumi/bls-eth-go-binary/)
- `-tags bignum_hol256`: Use the uint256 code that Geth uses, [`holiman/uint256`](https://github.com/holiman/uint256)
- `-tags bignum_pure`: Use the native Go Bignum implementation.


## Benchmarks

See [`BENCH.md`](./BENCH.md) for benchmarks of FFT, FFT in G1, FFT-extension, zero polynomials, and sample recovery.

## License

MIT, see [`LICENSE`](./LICENSE) file.

