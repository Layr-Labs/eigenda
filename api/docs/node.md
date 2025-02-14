# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [node/node.proto](#node_node-proto)
    - [AttestBatchReply](#node-AttestBatchReply)
    - [AttestBatchRequest](#node-AttestBatchRequest)
    - [BatchHeader](#node-BatchHeader)
    - [Blob](#node-Blob)
    - [BlobHeader](#node-BlobHeader)
    - [BlobQuorumInfo](#node-BlobQuorumInfo)
    - [Bundle](#node-Bundle)
    - [G2Commitment](#node-G2Commitment)
    - [GetBlobHeaderReply](#node-GetBlobHeaderReply)
    - [GetBlobHeaderRequest](#node-GetBlobHeaderRequest)
    - [MerkleProof](#node-MerkleProof)
    - [NodeInfoReply](#node-NodeInfoReply)
    - [NodeInfoRequest](#node-NodeInfoRequest)
    - [RetrieveChunksReply](#node-RetrieveChunksReply)
    - [RetrieveChunksRequest](#node-RetrieveChunksRequest)
    - [StoreBlobsReply](#node-StoreBlobsReply)
    - [StoreBlobsRequest](#node-StoreBlobsRequest)
    - [StoreChunksReply](#node-StoreChunksReply)
    - [StoreChunksRequest](#node-StoreChunksRequest)
  
    - [ChunkEncodingFormat](#node-ChunkEncodingFormat)
  
    - [Dispersal](#node-Dispersal)
    - [Retrieval](#node-Retrieval)
  
- [Scalar Value Types](#scalar-value-types)



<a name="node_node-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## node/node.proto



<a name="node-AttestBatchReply"></a>

### AttestBatchReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| signature | [bytes](#bytes) |  |  |






<a name="node-AttestBatchRequest"></a>

### AttestBatchRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| batch_header | [BatchHeader](#node-BatchHeader) |  | header of the batch |
| blob_header_hashes | [bytes](#bytes) | repeated | the header hashes of all blobs in the batch |






<a name="node-BatchHeader"></a>

### BatchHeader
BatchHeader (see core/data.go#BatchHeader)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| batch_root | [bytes](#bytes) |  | The root of the merkle tree with hashes of blob headers as leaves. |
| reference_block_number | [uint32](#uint32) |  | The Ethereum block number at which the batch is dispersed. |






<a name="node-Blob"></a>

### Blob
In EigenDA, the original blob to disperse is encoded as a polynomial via taking
taking different point evaluations (i.e. erasure coding). These points are split
into disjoint subsets which are assigned to different operator nodes in the EigenDA
network.
The data in this message is a subset of these points that are assigned to a
single operator node.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| header | [BlobHeader](#node-BlobHeader) |  | Which (original) blob this is for. |
| bundles | [Bundle](#node-Bundle) | repeated | Each bundle contains all chunks for a single quorum of the blob. The number of bundles must be equal to the total number of quorums associated with the blob, and the ordering must be the same as BlobHeader.quorum_headers. Note: an operator may be in some but not all of the quorums; in that case the bundle corresponding to that quorum will be empty. |






<a name="node-BlobHeader"></a>

### BlobHeader



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| commitment | [common.G1Commitment](#common-G1Commitment) |  | The KZG commitment to the polynomial representing the blob. |
| length_commitment | [G2Commitment](#node-G2Commitment) |  | The KZG commitment to the polynomial representing the blob on G2, it is used for proving the degree of the polynomial |
| length_proof | [G2Commitment](#node-G2Commitment) |  | The low degree proof. It&#39;s the KZG commitment to the polynomial shifted to the largest SRS degree. |
| length | [uint32](#uint32) |  | The length of the original blob in number of symbols (in the field where the polynomial is defined). |
| quorum_headers | [BlobQuorumInfo](#node-BlobQuorumInfo) | repeated | The params of the quorums that this blob participates in. |
| account_id | [string](#string) |  | The ID of the user who is dispersing this blob to EigenDA. |
| reference_block_number | [uint32](#uint32) |  | The reference block number whose state is used to encode the blob |






<a name="node-BlobQuorumInfo"></a>

### BlobQuorumInfo
See BlobQuorumParam as defined in
api/proto/disperser/disperser.proto


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| quorum_id | [uint32](#uint32) |  |  |
| adversary_threshold | [uint32](#uint32) |  |  |
| confirmation_threshold | [uint32](#uint32) |  |  |
| chunk_length | [uint32](#uint32) |  |  |
| ratelimit | [uint32](#uint32) |  |  |






<a name="node-Bundle"></a>

### Bundle
A Bundle is the collection of chunks associated with a single blob, for a single
operator and a single quorum.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| chunks | [bytes](#bytes) | repeated | Each chunk corresponds to a collection of points on the polynomial. Each chunk has same number of points. |
| bundle | [bytes](#bytes) |  | All chunks of the bundle encoded in a byte array. |






<a name="node-G2Commitment"></a>

### G2Commitment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| x_a0 | [bytes](#bytes) |  | The A0 element of the X coordinate of G2 point. |
| x_a1 | [bytes](#bytes) |  | The A1 element of the X coordinate of G2 point. |
| y_a0 | [bytes](#bytes) |  | The A0 element of the Y coordinate of G2 point. |
| y_a1 | [bytes](#bytes) |  | The A1 element of the Y coordinate of G2 point. |






<a name="node-GetBlobHeaderReply"></a>

### GetBlobHeaderReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob_header | [BlobHeader](#node-BlobHeader) |  | The header of the blob requested per GetBlobHeaderRequest. |
| proof | [MerkleProof](#node-MerkleProof) |  | Merkle proof that returned blob header belongs to the batch and is the batch&#39;s MerkleProof.index-th blob. This can be checked against the batch root on chain. |






<a name="node-GetBlobHeaderRequest"></a>

### GetBlobHeaderRequest
See RetrieveChunksRequest for documentation of each parameter of GetBlobHeaderRequest.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| batch_header_hash | [bytes](#bytes) |  |  |
| blob_index | [uint32](#uint32) |  |  |
| quorum_id | [uint32](#uint32) |  |  |






<a name="node-MerkleProof"></a>

### MerkleProof



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| hashes | [bytes](#bytes) | repeated | The proof itself. |
| index | [uint32](#uint32) |  | Which index (the leaf of the Merkle tree) this proof is for. |






<a name="node-NodeInfoReply"></a>

### NodeInfoReply
Node info reply


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| semver | [string](#string) |  |  |
| arch | [string](#string) |  |  |
| os | [string](#string) |  |  |
| num_cpu | [uint32](#uint32) |  |  |
| mem_bytes | [uint64](#uint64) |  |  |






<a name="node-NodeInfoRequest"></a>

### NodeInfoRequest
Node info request






<a name="node-RetrieveChunksReply"></a>

### RetrieveChunksReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| chunks | [bytes](#bytes) | repeated | All chunks the Node is storing for the requested blob per RetrieveChunksRequest. |
| chunk_encoding_format | [ChunkEncodingFormat](#node-ChunkEncodingFormat) |  | How the above chunks are encoded. |






<a name="node-RetrieveChunksRequest"></a>

### RetrieveChunksRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| batch_header_hash | [bytes](#bytes) |  | The hash of the ReducedBatchHeader defined onchain, see: https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/interfaces/IEigenDAServiceManager.sol#L43 This identifies which batch to retrieve for. |
| blob_index | [uint32](#uint32) |  | Which blob in the batch to retrieve for (note: a batch is logically an ordered list of blobs). |
| quorum_id | [uint32](#uint32) |  | Which quorum of the blob to retrieve for (note: a blob can have multiple quorums and the chunks for different quorums at a Node can be different). The ID must be in range [0, 254]. |






<a name="node-StoreBlobsReply"></a>

### StoreBlobsReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| signatures | [google.protobuf.BytesValue](#google-protobuf-BytesValue) | repeated | The operator&#39;s BLS sgnature signed on the blob header hashes. The ordering of the signatures must match the ordering of the blobs sent in the request, with empty signatures in the places for discarded blobs. |






<a name="node-StoreBlobsRequest"></a>

### StoreBlobsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blobs | [Blob](#node-Blob) | repeated | Blobs to store |
| reference_block_number | [uint32](#uint32) |  | The reference block number whose state is used to encode the blobs |






<a name="node-StoreChunksReply"></a>

### StoreChunksReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| signature | [bytes](#bytes) |  | The operator&#39;s BLS signature signed on the batch header hash. |






<a name="node-StoreChunksRequest"></a>

### StoreChunksRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| batch_header | [BatchHeader](#node-BatchHeader) |  | Which batch this request is for. |
| blobs | [Blob](#node-Blob) | repeated | The chunks for each blob in the batch to be stored in an EigenDA Node. |





 


<a name="node-ChunkEncodingFormat"></a>

### ChunkEncodingFormat
This describes how the chunks returned in RetrieveChunksReply are encoded.
Used to facilitate the decoding of chunks.

| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN | 0 |  |
| GNARK | 1 |  |
| GOB | 2 |  |


 

 


<a name="node-Dispersal"></a>

### Dispersal


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| StoreChunks | [StoreChunksRequest](#node-StoreChunksRequest) | [StoreChunksReply](#node-StoreChunksReply) | StoreChunks validates that the chunks match what the Node is supposed to receive ( different Nodes are responsible for different chunks, as EigenDA is horizontally sharded) and is correctly coded (e.g. each chunk must be a valid KZG multiproof) according to the EigenDA protocol. It also stores the chunks along with metadata for the protocol-defined length of custody. It will return a signature at the end to attest to the data in this request it has processed. |
| StoreBlobs | [StoreBlobsRequest](#node-StoreBlobsRequest) | [StoreBlobsReply](#node-StoreBlobsReply) | StoreBlobs is simiar to StoreChunks, but it stores the blobs using a different storage schema so that the stored blobs can later be aggregated by AttestBatch method to a bigger batch. StoreBlobs &#43; AttestBatch will eventually replace and deprecate StoreChunks method. DEPRECATED: StoreBlobs method is not used |
| AttestBatch | [AttestBatchRequest](#node-AttestBatchRequest) | [AttestBatchReply](#node-AttestBatchReply) | AttestBatch is used to aggregate the batches stored by StoreBlobs method to a bigger batch. It will return a signature at the end to attest to the aggregated batch. DEPRECATED: AttestBatch method is not used |
| NodeInfo | [NodeInfoRequest](#node-NodeInfoRequest) | [NodeInfoReply](#node-NodeInfoReply) | Retrieve node info metadata |


<a name="node-Retrieval"></a>

### Retrieval


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| RetrieveChunks | [RetrieveChunksRequest](#node-RetrieveChunksRequest) | [RetrieveChunksReply](#node-RetrieveChunksReply) | RetrieveChunks retrieves the chunks for a blob custodied at the Node. |
| GetBlobHeader | [GetBlobHeaderRequest](#node-GetBlobHeaderRequest) | [GetBlobHeaderReply](#node-GetBlobHeaderReply) | GetBlobHeader is similar to RetrieveChunks, this just returns the header of the blob. |
| NodeInfo | [NodeInfoRequest](#node-NodeInfoRequest) | [NodeInfoReply](#node-NodeInfoReply) | Retrieve node info metadata |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

