

# Organization

The EigenDA repo is organized as a monorepo, with each project adhering to the "Ben Johnson" project structure style. Within the core project directories (e.g., `core`, `disperser`, `node`, `retriever`, `indexer`), the main interfaces and data types are defined at the root of the project, while implementations are organized by dependency. For instance, the folder `indexer/inmem` contains implementations of the interfaces in `indexer` which use in-memory storage, while `indexer/leveldb` may contain implementations of the same interfaces that use `leveldb`. Mocks of all interfaces in the `indexer` project go in `indexer/mock`. 

The same pattern is used for intra-project and inter-project dependencies. For instance, the folder `indexer/indexer` contains implementations of the interfaces in `core` which depend on the `indexer` project. 

In general, the `core` project contains implementation of all the important business logic responsible for the security guarantees of the EigenDA protocol, while the other projects add the networking layers needed to run the distributed system. 


# Directory structure
<pre>
┌── <a href="./api">api</a> Protobuf definitions and contract bindings
├── <a href="./contracts">contracts</a>
|   ├── <a href="./contracts/eignlayer-contracts">eigenlayer-contracts</a>: Contracts for the EigenLayer restaking platform
┌── <a href="./core">core</a>: Core logic of the EigenDA protocol
├── <a href="./disperser">disperser</a>: Disperser service
├── <a href="./docs">docs</a>: Documentation and specification
── <a href="./indexer">indexer</a>: A simple indexer for efficiently tracking chain state and maintaining accumulators
├── <a href="./node">node</a>: DA node service
├── <a href="./pkg">pkg</a>
|   ├── <a href="./pkg/encoding">encoding</a>: Core encoding/decoding functionality and multiproof generation
|   └── <a href="./pkg/kzg">kzg</a>: kzg libraries
├── <a href="./retriever">retriever</a>: Retriever service
├── <a href="./test">test</a>: Tools for running integration tests
</pre>
