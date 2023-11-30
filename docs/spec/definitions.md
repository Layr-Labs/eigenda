# Definitions


## Data Packaging

**Blob**. Blobs are the fundamental unit of data posted to EigenDA by users.

**Batch**. 

## System Components

**DA Node**. The DA Node is an off-chain component which is run by an EigenLayer operator. EigenDA operators are responsible for accepting data from the disperser, certifying its availability, and following a protocol for distributing the data to registered retrievers. It is assumed that honest DA nodes are delegated a threshold proportion of the stake from EigenLayer restakers, where this threshold may be defined per DA end-user.

**Disperser**. The Disperser is an off-chain component which is responsible for packaging data blobs in a specific way, distributing their data among the DA nodes, aggregating certifications from the nodes, and then pushing the aggregated certificate to the chain. The disperser is an untrusted system component. 

**Retriever**. The Retriever is an off-chain component which implements a protocol for receiving data blobs from the set of DA nodes. 

**Smart contracts**. The smart contracts track the amount of EigenLayer stake held by each operator, and allow for DA certifications to be verified on-chain and consumed by downstream applications such as rollup smart contracts. 

## Staking Concepts

**Quorum**
