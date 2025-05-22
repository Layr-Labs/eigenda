# EigenDA OP secure integration

This doc serves as a high level guide to walk through our approach to OP secure integration. We will touch on
topics including 
- batch posting by op-batcher with eigenda-proxy write path
- derivation pipeline by op-node via eigenda-proxy read path
- filtering invalid batch posted from op-batcher from proxy
- hokulea - our secure integration module built on top of rust implmenetation of op derivation pipeline called [kona](https://github.com/op-rs/kona)
- secure integration with zkVM like op-succinct and op-kailua
- secure integration with OP interactive fault proof

We start with a system workflow diagram specifying the expected behavior of all parties in the system. Then we dive into a concept called derivation
pipeline, and present where and what EigenDA makes changes in the derivation pipeline.

Then we start to introduce malicious op-batcher and malicious L2 proposer working together to post wrong state roots on L1, we show how hokulea can
be used to provide security integration by allowing honest challengers to resolve the wrong state root.

## High Level System Workflow Diagram 

At the high level, a rollup can be splitted into two parts: write path to L1 and read path from L1

- The write path ensures the liveness of the L2 consensus, constant inflow of L2 blocks and L1 user deposit transactions.
- The read path controls the safety of the L2 consensus. As long as there is consensus on the data input, the deterministic EVM can produce identical outcome
- However, because L2 state is bridged into L1, to make any challenge in a fault proof game (zk or interactive), the L2 consensus must progress to the claimed L2 block to make the counter-claim, but it requires liveness of L2 consensus to be able to reach certain block height In short, the read path must not stall.

### Write path
An op-batcher is a critical part of the rollup sequencer, it is responsible for collecting L2 blocks containing user-transactions. In the diagram below,
the op-batcher sends rollup data into an eigenda proxy and uses its write-path to convert the data into an eigenda blob, and eventually sends it over to the EigenDA network. With EigenDA v2, a DA certificate is returned within seconds and given back to the op-batcher. The DA certificate is eventually sent into L1 inbox.

![](../../assets/integration/op-integration-normal.png)

### Read(consensus) path

The read path from L1 determines the data inputs. When we talk about OP secure integration, essentially we are talking about a correct way to derive
rollup L2 blocks deterministically from L1. OP has defined a derivation pipeline in OP [spec](https://specs.optimism.io/protocol/derivation.html#l2-chain-derivation-pipeline), which by definition of is the correct way to agree on a list of L2 transactions from L1 data source.

Both [op-program](https://github.com/ethereum-optimism/optimism/tree/develop/op-program) in Golang and [Kona](https://github.com/op-rs/kona/tree/main)
in rust implement the derivation pipeline

Like the diagram above, the derivation pipeline consists of stages that bring L1 transactions down to Payload Attributes which are L2 blocks without output properties like Merkle Patricia Tree Root. On the read path, the op-node uses the `read-path` of the eigenda-proxy to fetch EigenDA blob which
also decodes the blob in the same form as what op-batcher has sent it.

## L2 consensus Path with EigenDA

At the high level, a secure OP integration requires that on the L2 consensus read path, and for the EigenDA section of the derivation pipeline
- there is one deterministic output for every EigenDA certificate
- the EigenDA part of derivation pipeline must not stall, i.e the stage must handle all types of inputs, and none of them can stall the entire derivation pipeline.

### Discard invalid or expired Certs

We have defined the spec in the hokulea [repo](https://github.com/Layr-Labs/hokulea/tree/master/docs). But roughly speaking, a cert must 
- have correct serialization
- no be too old
- have been attested by sufficient stake
- corresponds to an eigenda blob whose rollup data can be correctly decoded

If any conditions are not met, the cert is discarded in the altda stage in the the op-node derivation.

### Deterministic input and outputs relation

In every step in the hokulea derivation pipeline, we ensure there is one-to-one correspondence between DA cert and eigenda blob. It is guaranteed
by KZG commitment, that no adversary can generate an identical commitment from two distinct blobs. All processing are deterministic.

So far, we have presented a rollup with EigenDA allowing all honest L2 nodes to reach consensus. But as mentioned, a secure integration requires
L1 be convinced that some L2 state is correctly derived based on the L2 deriviation rule.

## Proving L2 state is derived according to the L2 derivation rule

All fault proof or validity rollup requires a VM that is fault provable. The VM provides a constraint system, such that any incorrect execution
can be one step proven to be wrong.
So at the high level, a FPVM (fault proof VM) is given a L2 derivation software and a data source called preimage oracle.
The preimage oracle contains all the data necessary for the VM to execute the software logics.

Take hokulea derivation rule as an example, the derivation must correctly discard invalid cert that has wrong serialization or does not have
sufficient stake attesting it, assuming it has correct preimage inputs. 

In hokulea, there are three preimage inputs, more information can be found at [secure integration page](./6-secure-integration.md)
- a boolean indicating if a DA cert is valid
- a recency window stating how old a DA cert can be before discarding
- eigenDA blobs

The preimage communication spec in hokulea can be found at the hokulea doc [site](https://github.com/Layr-Labs/hokulea/tree/master/docs).

For zkVM, Hokulea has ensured all inputs are verified before giving the data for zkVM to consume.

