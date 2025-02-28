# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [relay/relay.proto](#relay_relay-proto)
    - [ChunkRequest](#relay-ChunkRequest)
    - [ChunkRequestByIndex](#relay-ChunkRequestByIndex)
    - [ChunkRequestByRange](#relay-ChunkRequestByRange)
    - [GetBlobReply](#relay-GetBlobReply)
    - [GetBlobRequest](#relay-GetBlobRequest)
    - [GetChunksReply](#relay-GetChunksReply)
    - [GetChunksRequest](#relay-GetChunksRequest)
  
    - [Relay](#relay-Relay)
  
- [Scalar Value Types](#scalar-value-types)



<a name="relay_relay-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## relay/relay.proto



<a name="relay-ChunkRequest"></a>

### ChunkRequest
A request for chunks within a specific blob. Requests are fulfilled in all-or-nothing fashion. If any of the
requested chunks are not found or are unable to be fetched, the entire request will fail.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| by_index | [ChunkRequestByIndex](#relay-ChunkRequestByIndex) |  | Request chunks by their individual indices. |
| by_range | [ChunkRequestByRange](#relay-ChunkRequestByRange) |  | Request chunks by a range of indices. |






<a name="relay-ChunkRequestByIndex"></a>

### ChunkRequestByIndex
A request for chunks within a specific blob. Each chunk is requested individually by its index.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob_key | [bytes](#bytes) |  | The blob key. |
| chunk_indices | [uint32](#uint32) | repeated | The index of the chunk within the blob. |






<a name="relay-ChunkRequestByRange"></a>

### ChunkRequestByRange
A request for chunks within a specific blob. Each chunk is requested a range of indices.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob_key | [bytes](#bytes) |  | The blob key. |
| start_index | [uint32](#uint32) |  | The first index to start fetching chunks from. |
| end_index | [uint32](#uint32) |  | One past the last index to fetch chunks from. Similar semantics to golang slices. |






<a name="relay-GetBlobReply"></a>

### GetBlobReply
The reply to a GetBlobs request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob | [bytes](#bytes) |  | The blob requested. |






<a name="relay-GetBlobRequest"></a>

### GetBlobRequest
A request to fetch one or more blobs.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob_key | [bytes](#bytes) |  | The key of the blob to fetch. |






<a name="relay-GetChunksReply"></a>

### GetChunksReply
The reply to a GetChunks request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [bytes](#bytes) | repeated | The chunks requested. The order of these chunks will be the same as the order of the requested chunks. data is the raw data of the bundle (i.e. serialized byte array of the frames) |






<a name="relay-GetChunksRequest"></a>

### GetChunksRequest
Request chunks from blobs stored by this relay.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| chunk_requests | [ChunkRequest](#relay-ChunkRequest) | repeated | The chunk requests. Chunks are returned in the same order as they are requested. |
| operator_id | [bytes](#bytes) |  | If this is an authenticated request, this should hold the ID of the operator. If this is an unauthenticated request, this field should be empty. Relays may choose to reject unauthenticated requests. |
| operator_signature | [bytes](#bytes) |  | If this is an authenticated request, this field will hold a BLS signature by the requester on the hash of this request. Relays may choose to reject unauthenticated requests.

The following describes the schema for computing the hash of this request This algorithm is implemented in golang using relay.auth.HashGetChunksRequest().

All integers are encoded as unsigned 4 byte big endian values.

Perform a keccak256 hash on the following data in the following order: 1. the operator id 2. for each chunk request: a. if the chunk request is a request by index: i. a one byte ASCII representation of the character &#34;i&#34; (aka Ox69) ii. the blob key iii. the start index iv. the end index b. if the chunk request is a request by range: i. a one byte ASCII representation of the character &#34;r&#34; (aka Ox72) ii. the blob key iii. each requested chunk index, in order |





 

 

 


<a name="relay-Relay"></a>

### Relay
Relay is a service that provides access to public relay functionality.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetBlob | [GetBlobRequest](#relay-GetBlobRequest) | [GetBlobReply](#relay-GetBlobReply) | GetBlob retrieves a blob stored by the relay. |
| GetChunks | [GetChunksRequest](#relay-GetChunksRequest) | [GetChunksReply](#relay-GetChunksReply) | GetChunks retrieves chunks from blobs stored by the relay. |

 



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

