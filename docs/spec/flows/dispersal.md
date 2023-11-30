
# Dispersal

Data is made available on EigenDA through the following flow: 
1. The [Disperser](./disperer.md) encodes the data in accordance with the [storage module](./protocol-modules/storage/overview.md) requirements, constructs the appropriate header, and sends the chunks to the DA nodes.
2. Upon receiving signatures from the DA nodes, the disperser aggregates these signatures. 
3. Next, the disperser sends the aggregated signatures and header to the `confirmBatch` method of the `ServiceManager`
4. Once retrievers see the confirmed Batch on chain, they can request to download the associated chunks from a set of DA nodes, in accordance with the [retrieval module](./protocol-modules/retrieval/retrieval.md) of the protocol. 
