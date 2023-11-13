### `BlobStoreRequest`

```go
type BlobStoreRequest struct {
    BlobID [32]byte
    Data []byte
    HeaderQuorumParams []struct{
        QuorumID uint8
        AdversaryThresholdBPs uint16
        QuorumThresholdBPs    uint16
    }
}
```

### `BlobStoreResponse`

```go
type BlobStoreResponse struct {
    BlobID [32]byte
    StartDegree uint32
    EndDegree uint32
    DataStoreHeader DataStoreHeader
}