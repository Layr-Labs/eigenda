This folder contains the API documentation for the gRPC services included in the EigenDA platform. Each markdown file contains the protobuf definitions for each respective service including:
- Churner: a hosted service responsible for maintaining the active set of Operators in the EigenDA network based on their delegated TVL.
- Disperser: the hosted service and primary point of interaction for Rollup users.
- Node: individual EigenDA nodes run on the network by EigenLayer Operators.
- Retriever: a service that users can run on their own infrastructure, which exposes a gRPC endpoint for retrieval of blobs from EigenDA nodes.

