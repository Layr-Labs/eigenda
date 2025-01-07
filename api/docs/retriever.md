# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [retriever/v2/retriever.proto](#retriever_v2_retriever-proto)
    - [BlobReply](#retriever-v2-BlobReply)
    - [BlobRequest](#retriever-v2-BlobRequest)
  
    - [Retriever](#retriever-v2-Retriever)
  
- [Scalar Value Types](#scalar-value-types)



<a name="retriever_v2_retriever-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## retriever/v2/retriever.proto



<a name="retriever-v2-BlobReply"></a>

### BlobReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [bytes](#bytes) |  | The blob retrieved and reconstructed from the EigenDA Nodes per BlobRequest. |






<a name="retriever-v2-BlobRequest"></a>

### BlobRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob_header | [common.v2.BlobHeader](#common-v2-BlobHeader) |  | header of the blob to be retrieved |
| reference_block_number | [uint32](#uint32) |  | The Ethereum block number at which the batch for this blob was constructed. |
| quorum_id | [uint32](#uint32) |  | Which quorum of the blob this is requesting for (note a blob can participate in multiple quorums). |





 

 

 


<a name="retriever-v2-Retriever"></a>

### Retriever
The Retriever is a service for retrieving chunks corresponding to a blob from
the EigenDA operator nodes and reconstructing the original blob from the chunks.
This is a client-side library that the users are supposed to operationalize.

Note: Users generally have two ways to retrieve a blob from EigenDA V2:
  1) Retrieve from the relay that the blob is assigned to: the API
     is Relay.GetBlob() as defined in api/proto/relay/relay.proto
  2) Retrieve directly from the EigenDA Nodes, which is supported by this Retriever.

The Relay.GetBlob() (the 1st approach) is generally faster and cheaper as the
relay manages the blobs that it has processed, whereas the Retriever.RetrieveBlob()
(the 2nd approach here) removes the need to trust the relay, with the downside of
worse cost and performance.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| RetrieveBlob | [BlobRequest](#retriever-v2-BlobRequest) | [BlobReply](#retriever-v2-BlobReply) | This fans out request to EigenDA Nodes to retrieve the chunks and returns the reconstructed original blob in response. |

 



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

