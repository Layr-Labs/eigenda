# encoding


- performs Reed Solomon Encoding using elliptic curve points. The library enables KZG multi-proof and reveal in O(n log n) time using FFT, based on FK20 algorithm.

- is built upon crypto primitive from https://pkg.go.dev/github.com/protolambda/go-kzg

- accepts arbitrary number of systematic nodes, parity nodes and data size, free of restriction on power of 2
