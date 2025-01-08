# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [node/v2/node_v2.proto](#node_v2_node_v2-proto)
    - [GetChunksReply](#node-v2-GetChunksReply)
    - [GetChunksRequest](#node-v2-GetChunksRequest)
    - [NodeInfoReply](#node-v2-NodeInfoReply)
    - [NodeInfoRequest](#node-v2-NodeInfoRequest)
    - [StoreChunksReply](#node-v2-StoreChunksReply)
    - [StoreChunksRequest](#node-v2-StoreChunksRequest)
  
    - [Dispersal](#node-v2-Dispersal)
    - [Retrieval](#node-v2-Retrieval)
  
- [Scalar Value Types](#scalar-value-types)



<a name="node_v2_node_v2-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## node/v2/node_v2.proto



<a name="node-v2-GetChunksReply"></a>

### GetChunksReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| chunks | [bytes](#bytes) | repeated | All chunks the Node is storing for the requested blob per RetrieveChunksRequest. |






<a name="node-v2-GetChunksRequest"></a>

### GetChunksRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob_key | [bytes](#bytes) |  |  |
| quorum_id | [uint32](#uint32) |  | Which quorum of the blob to retrieve for (note: a blob can have multiple quorums and the chunks for different quorums at a Node can be different). The ID must be in range [0, 254]. |






<a name="node-v2-NodeInfoReply"></a>

### NodeInfoReply
Node info reply


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| semver | [string](#string) |  |  |
| arch | [string](#string) |  |  |
| os | [string](#string) |  |  |
| num_cpu | [uint32](#uint32) |  |  |
| mem_bytes | [uint64](#uint64) |  |  |






<a name="node-v2-NodeInfoRequest"></a>

### NodeInfoRequest
Node info request






<a name="node-v2-StoreChunksReply"></a>

### StoreChunksReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| signature | [bytes](#bytes) |  |  |






<a name="node-v2-StoreChunksRequest"></a>

### StoreChunksRequest
Request that the Node store a batch of chunks.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| batch | [common.v2.Batch](#common-v2-Batch) |  | batch of blobs to store |
| disperserID | [uint32](#uint32) |  | ID of the disperser that is requesting the storage of the batch. |
| signature | [bytes](#bytes) |  | Signature using the disperser&#39;s ECDSA key over keccak hash of the batch. The purpose of this signature is to prevent hooligans from tricking DA nodes into storing data that they shouldn&#39;t be storing.

Algorithm for computing the hash is as follows. All integer values are serialized in big-endian order (unsigned). A reference implementation (golang) can be found at https://github.com/Layr-Labs/eigenda/blob/master/disperser/auth/request_signing.go

1. digest batch.BatchHeader.BatchRoot 2. digest batch.BatchHeader.ReferenceBlockNumber (8 bytes, unsigned big endian) 3. for each certificate in batch.BlobCertificates: a. digest certificate.BlobHeader.Version (4 bytes, unsigned big endian) b. for each quorum_number in certificate.BlobHeader.QuorumNumbers: i. digest quorum_number (4 bytes, unsigned big endian) c. digest certificate.BlobHeader.Commitment.Commitment d. digest certificate.BlobHeader.Commitment.LengthCommitment e. digest certificate.BlobHeader.Commitment.LengthProof f. digest certificate.BlobHeader.Commitment.Length (4 bytes, unsigned big endian) g. digest certificate.BlobHeader.PaymentHeader.AccountId h. digest certificate.BlobHeader.PaymentHeader.ReservationPeriod (4 bytes, unsigned big endian) i. digest certificate.BlobHeader.PaymentHeader.CumulativePayment j. digest certificate.BlobHeader.PaymentHeader.Salt (4 bytes, unsigned big endian) k. digest certificate.BlobHeader.Signature l. for each relay in certificate.Relays: i. digest relay (4 bytes, unsigned big endian) 4. digest disperserID (4 bytes, unsigned big endian)

Note that this signature is not included in the hash for obvious reasons. |





 

 

 


<a name="node-v2-Dispersal"></a>

### Dispersal
WARNING: the following RPCs are experimental and subject to change.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| StoreChunks | [StoreChunksRequest](#node-v2-StoreChunksRequest) | [StoreChunksReply](#node-v2-StoreChunksReply) |  |
| NodeInfo | [NodeInfoRequest](#node-v2-NodeInfoRequest) | [NodeInfoReply](#node-v2-NodeInfoReply) |  |


<a name="node-v2-Retrieval"></a>

### Retrieval


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetChunks | [GetChunksRequest](#node-v2-GetChunksRequest) | [GetChunksReply](#node-v2-GetChunksReply) | GetChunks retrieves the chunks for a blob custodied at the Node. |
| NodeInfo | [NodeInfoRequest](#node-v2-NodeInfoRequest) | [NodeInfoReply](#node-v2-NodeInfoReply) | Retrieve node info metadata |

 



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

