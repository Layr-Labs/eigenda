

# EigenDA

EigenDA is a Data Availability (DA) service, implemented as an actively validated service (AVS) on EigenLayer, that provides secure and scalable DA for L2s on Ethereum. 

## What is DA?

In informal terms, DA is a guarantee that a given piece of data will be available to anyone who wishes to retrieve it. 

A DA system accepts blobs of data (via some interface) and then makes them available to retrievers (through another interface). 

Two important aspects of a DA system are 
1. Security: The security of a DA system constitutes the set of conditions which are sufficient to ensure that all data blobs certified by the system as available are indeed available for honest retrievers to download. 
2. Throughput: The throughput of a DA system is the rate at which the system is able to accept blobs of data, typically measured in bytes/second. 

## An EigenLayer AVS for DA

EigenDA is implemented as an actively validated service on EigenLayer, which is a restaking protocol for Ethereum. 

Because of this, EigenDA makes use of the EigenLayer state, which is stored on Ethereum, for consensus about the state of operators and as a callback for consensus about the availability of data. This means that EigenDA can be simpler in implementation than many existing DA solutions: EigenDA doesn't need to build it's own chain or consensus protocol; it rides on the back of Ethereum. 

### A first of its kind, horizontally scalable DA solution

Among extant DA solutions, EigenDA takes an approach to scalability which is unique in that it yields true horizontal scalability: Every additional unit of capacity contributed by an operator can increase the total system capacity. 

This property is achieved by using a Reed Solomon erasure encoding scheme to shard the blob data across the DA nodes. While other systems such as Celestia and Danksharding (planned) also make use of Reed Solomon encoding, they do so only for the purpose of supporting certain observability properties of Data Availability Sampling (DAS) by light nodes. On the other hand, all incentivized/full nodes of the system download, store, and serve the full system bandwidth. 

Horizontal scalability provides the promise for the technological bottlenecks of DA capacity to continually track demand, which has enormous implications for Layer 2 ecosystems. 

## Security Model

EigenDA produces a DA attestation which asserts that a given blob or collection of blobs is available. Attestations are anchored to one or more "Quorums," each of which defines a set of EigenLayer stakers which underwrite the security of the attestation. Quorums should be considered as redundant: Each quorum linked to an attestation provides an independent guarantee of availability as if the other quorums did not exist. 


Each attestation is characterized by safety and liveness tolerances:
- Liveness tolerance: Conditions under which the system will produce an availability attestation. 
- Safety tolerance: Conditions under which an availability attestation implies that data is indeed available. 

EigenDA defines two properties of each blob attestation which relate to its liveness and safety tolerance: 
- Liveness threshold: The liveness threshold defines the minimum percentage of stake which an attacker must control in order to mount a liveness attack on the system. 
- Safety threshold: The safety threshold defines the total percentage of stake which an attacker must control in order to mount a first-order safety attack on the system. 

The term "first-order attack" alludes to the fact that exceeding the safety threshold may represent only a contingency rather than an actual safety failure due to the presence of recovery mechanisms that would apply during such a contingency. Discussion of such mechanisms is outside of the scope of the current documentation. 

Safety thresholds can translate directly into cryptoeconomic safety properties for quorums consisting of tokens which experience toxicity in the event of publicly observable attacks by a large coalition of token holders. This and other discussions of cryptoeconomic security are also beyond the scope of this technical documentation. We restrict the discussion to illustrating how the protocol preserves the given safety and liveness thresholds.