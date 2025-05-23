# EigenDA OP Secure integration

This document offers a concise, high-level overview of how we securely integrate EigenDA with the Optimism (OP) stack. We will touch topics on:
- Separtion of `write` and `read` path in a L2 consensus
- Requirement on the `read` path for a secure EigenDA integration for L2 consensus
  - Extending derivation pipeline with eigenda-proxy read path
  - Filtering invalid batches emitted by op-batcher
- Hokulea, our secure integration module that enables secure OP integration
- zkVMs such as op-succinct and op-kailua, and potential integration with Optimism’s interactive fault-proof system

## Write and Read path in L2 consensus

At the high level, a rollup can be splitted into two parts: write path to L1 and read path from L1

- The `write path` ensures the liveness of the L2 consensus. It consists of L2 blocks produced by op-batch and L1 induced transactions.
- The `read path` controls the safety of the L2 consensus. It ensures all L2 consensus node sees an identical list of L2 blocks, then a EVM can produce identical L2 state
- However, because L2 state is bridged into L1, to make any challenge in a fault proof game (zk or interactive), the L2 consensus must progress to the claimed L2 block to decide whether to challenge. So a secure integration requires liveness of L2 consensus to be able to reach to the block height. So in short, the read path must not stall regardless of op-batcher is maliicous.

### L2 Write path
An op-batcher is a critical part of the rollup sequencer, it is responsible for collecting L2 blocks containing user-transactions. In the diagram below,
the op-batcher sends rollup data into an eigenda proxy and uses its write-path to convert the data into an eigenda blob, and eventually sends it over to the EigenDA network. With EigenDA v2, a DA certificate is returned within seconds and given back to the op-batcher. The DA certificate is eventually sent into L1 inbox.

![](../../assets/integration/op-integration-normal.png)

### L2 Read path

The read path from L1 determines L2 consensus. An OP secure integration essentially means about a correct way to derive
rollup L2 blocks from L1. OP has defined a derivation pipeline in OP [spec](https://specs.optimism.io/protocol/derivation.html#l2-chain-derivation-pipeline), which is the definition of the correct way. 
Both [op-program](https://github.com/ethereum-optimism/optimism/tree/develop/op-program) in Golang and [Kona](https://github.com/op-rs/kona/tree/main)
in Rust implement the derivation pipeline. Like the diagram above, the derivation pipeline consists of stages that bring L1 transactions down to Payload Attributes which are L2 blocks.
To support secure integration, we have defined and inserted a Eigenda section in the OP derivation pipeline.
Both Eigenda proxy and Hokulea which support secure integration implement the eigenda blob derivation.

## L2 Read path with EigenDA

The op-node uses the `read-path` on the eigenda-proxy to fetch EigenDA blob from the EigenDA network. The proxy performs internal check to decide
if a given blob is elligible for retrieval, depending on if a certificate is valid or recent enough. The proxy decodes the blob and returns the data
to the next stage of OP derivation pipeline. The key properties which EigenDA derivation strives to keep are

- there is one deterministic output (decoded data) for every EigenDA certificate
- the EigenDA derivation pipeline must not stall

We define [spec](https://github.com/Layr-Labs/hokulea/tree/master/docs) to ensure those properties. We sometimes refer EigenDA derivation pipeline
as hokulea derivation pipeline as this pipeline is designed for OP stack.

### Discard Certs

We need to discard certs that can stall the derivation pipeline, as it endangers the safety of L2 state posted on L1.
A certificate can be incorrect for many reasons, such as wrong cert serialization, incorrect cert format and etc, more can be found in [spec](https://github.com/Layr-Labs/hokulea/tree/master/docs). Those errors has nature with stateless data processing, because they do not require L1
state to infer the correctness.
The other group of incorrectness requires L1 information such as
- whether the DA cert has been attested by sufficient stake
- whether the DA cert has been generated recent enough, such that it is not recent enough

### Deterministic input and outputs relation

In every step in the hokulea derivation pipeline, we ensure for every DA cert there is a unique eigenda blob. It is ensured
by KZG commitment, such that no adversary can generate an identical commitment from two distinct blobs.

## Proving L2 consensus to L1

So far, we have ensured all L2 nodes integrated with EigenDA can reach a common L2 state and the chain is always live even in the presence of
misbehaving op-batcher.
However, the security of rollup is determined by if there are provable ways to challenge incorrect L2 state posted on L1.
In this section, we discuss our OP secure integraton library **Hokulea**.

### Short intro to OP FPVM

The correctness of a L2 state is determined by the derivation rules, which are implemented in both Go [op-program](https://github.com/ethereum-optimism/optimism/tree/develop/op-program) and Rust [Kona](https://github.com/op-rs/kona/tree/main).

With interactive fault proof, the derivation logics are packaged into a binary ELF file, which can be run inside a FPVM. Cannon is a production
level implementation of FPVM, such that any incorrect execution are fault provable.

However, the FPVM requires both the ELF binary and data (including L2 batches and L1 deposits, etc) to be able to derive the final L2 blocks.
Conceptually, it is very similar to what op-node does to reach L2 consensus, except that the ELF only contains the core logics without actual data.

So a preimage oracle is also fed into the FPVM. And OP spec has ensured all types of data in the preimage oracle are verifiable and
self-consistently. 
The preimage oracle by itself is a key-value map, where the value has to hold a binding relation with key, and value has to pass certain logics
defined by the OP [Spec](https://specs.optimism.io/fault-proof/index.html#pre-image-oracle) in order to be considered valid for the key.
A FPVM creates a binding system, it creates a relation between (software logic, data, and VM opcode).

### Hokulea

Hokulea builds on top of OP Kona derivation pipeline to integrate EigenDA as a Data Availability Source. Hokulea provides traits, implementation
for EigenDA part of derivation pipeline, such that those logics can be compiled into ELF together with Kona.

Hokulea also extends preimage oracle for EigenDA. Recall previously we specify that a secure integration must be able to
- deterministically derive rollup payload from a EigenDA certificate
- discard DA certs that can stall the derivation pipeline

For the first point, the determinism is provided by the decoding algorithm and kzg commitment.

Previously, we categorized illegitimate DA certs into `stateless data processing` and `L1 state dependent` error

To obtain this information, the hokulea(EigenDA) derivation pipeline communicates to Preimage Oracle to obtain necessary information.
- if a DA cert is correct
- what is the current recency window (which we choose to be set as `sequencing_window_size`, more see [secure-integration](./6-secure-integration.md) page)

The communication spec can be found at [Hokulea](https://github.com/Layr-Labs/hokulea/tree/master/docs). Because EigenDA operators are native
on Ethereum L1, DA cert can be verified with a smart contract call.
We also developed a rust library called [**CANOE**](https://github.com/Layr-Labs/hokulea/tree/master/canoe#1protocol-overview) that uses zk validity proof to provide a way 
- to efficiently verify the Cert validity
- to verify Cert validity within zkVM by verifying the canoe zk validity proof

### Hokulea integration with zkVM

Unlike interactively challenge game with fault proof, a zk proof has a property that only the honest party can create a valid zk proof.
The incorrect party can raise a challenge but is unable to defend its position.
How Hokulea support zkVM is very similar to FPVM.
The Hokulea+Kona derivation and preimage oracle are fed into a zkVM (Risc0 or Sp1), the zkVM produces a zk validity proof, which can be
verified onchain.
We currently integrate with [OP-succinct](https://github.com/succinctlabs/op-succinct) and [OP-Kailua](https://github.com/risc0/kailua).
For an integration guide, please refer to the [preloader](https://github.com/Layr-Labs/hokulea/tree/master/example/preloader) example for zk integration.

ZkVM has additional contraint than FPVM that all the data in the preimage oracle must be verified within zkVM.

For Hokulea, it means DA cert must be verified by zkVM. We developed a rust library called [Canoe](https://github.com/Layr-Labs/hokulea/tree/master/canoe#1protocol-overview) that provides two implementations that take advantages of [Risc0 Steel](https://risczero.com/steel) and 
[Succinct Sp1 contract call](https://github.com/succinctlabs/sp1-contract-call).

The constraint also requires all eigenda blob with respect to the kzg commitment in the DA cert. We developed a similar library called
[rust-kzg-bn254](https://github.com/Layr-Labs/rust-kzg-bn254) that offers similar functionalities as [4844 spec](https://github.com/ethereum/consensus-specs/blob/86fb82b221474cc89387fa6436806507b3849d88/specs/deneb/polynomial-commitments.md).