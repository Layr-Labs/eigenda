
# DA Node

The DA node has responsibilities in both the storage protocol and the retrieval protocol.

```
type Service interface{
    StoreChunks(ctx context.Context, in *StoreChunksRequest) (*StoreChunksReply, error)
}
```

## RPC Interface

The `StoreChunks` method is used for sending encoded chunks to the node.

- method: `StoreChunks`
- params:
    - paymentsRoot: bytes
    - payments: repeated bytes
    - headerRoot: bytes
    - header: [`DataStoreHeader`](./types/node-types.md#datastoreheader)
    - chunks: repeated bytes
- returns:
    - signature: bytes


When the `StoreChunks` method is called, the node performs the following checks:
1. Check that all payments are correct (See [Payment Constraints](./node-payments.md)).
2. Check that its own chunks are correct (See [Blob Encoding Constraints](./node-encoding.md))

Provided that both checks are successful, the node will sign the concatenation of the paymentRoot and blobRoot using the BLS key registered with the `BLSRegistry` and then return the signature. 



## Adapters 

### Indexer Adapter

The DA Node utilizes an adapter on top of the [Indexer](./indexer.md) interface which provides a view of the registration status and delegated stake for each operator at a given block number. The relevant structs are [Registrant](./types/node-types.md#registrant), [TotalStake](./types/node-types.md#totalstake), [TotalOperator](./types/node-types.md#totaloperator), [RegistrantView](./types/node-types.md#registrantview), and [StateView](./types/node-types.md#stateview).

```go

type IndexerAdapter interface{
    GetStateView(ctx context.Context, blockNumber uint32) (*StateView, error)
}
```
