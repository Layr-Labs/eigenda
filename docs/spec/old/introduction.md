# Introduction

The EigenDA system is a scalable Data Availability (DA) service that is secured by Ethereum stakers via EigenLayer. In informal terms, DA is a guarantee that a given piece of data is available to anyone who wishes to retrieve it. EigenDA is focused on providing DA with both high security and throughput. 

At a high level, a DA system is one which accepts blobs of data via some interface and then makes them available to retrievers through another interface. 

Two important aspects of a DA system are 
1. Security: The security of a DA system constitutes the set of conditions which are sufficient to ensure that all data blobs certified by the system as available are indeed available for honest retrievers to download. 
2. Throughput: The throughput of a DA system is the rate at which the system is able to accept blobs of data, typically measured in bytes/second. 

## EigenLayer Quorums

Most baseline EigenDA security guarantees are derived under a Byzantine model which stipulates that a maximum percentage of validators will behave adversarially at any given moment in time. As an EigenLayer AVS, EigenDA makes use of the validator set represented by  validators who have restaked Ether or other staking assets via the EigenLayer platform. Consequently, all constraints on adversarial behavior per the Byzantine modeling approach take the form of a maximum amount of stake which can be held by adversarial agents. 

An important aspect of restaking on EigenLayer is the notion of a quorum. EigenLayer supports restaking of various types of assets, from natively staked Ether and Liquid Staking Tokens (LSTs) to the wrapped assets of other protocols such as Ethereum rollups. Since these different categories of assets can have arbitrary and variable exchange rates, EigenLayer supports the quorum as a means for the users of protocols such as EigenDA to specify the nominal level of security that each staking token is taken to provide. In practice, a quorum is a vector specifying the relative weight of each staking strategy supported by EigenLayer. 

[TODO: Reference EigenLayer SDK quorum documentation]

## Essential Properties of EigenDA

### EigenDA Security

When an end user posts a blob of data to EigenDA, they can specify a list of [security parameters](./data-model.md#quorum-information), each of which consists of a `QuorumID` identifying a particular quorum registered with EigenDA and an `AdversaryThreshold` which specifies the Byzantine adversarial tolerance that the user expects the blobs availability to respect. 

For such a blob accepted by the system (See [Dispersal Flow](./flows/dispersal.md)), EigenDA delivers the following security guarantee: Unless more than an `AdversaryThreshold` of stakers acts adversarially in every quorum associated with the blob, the blob will be available to any honest retriever. How this guarantee is supported is discussed further in [The Modules of Data Availability](./protocol-modules/overview.md)

### EigenDA Throughput

[TODO: Complete this section]

