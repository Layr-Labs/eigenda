# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [lightnode/lightnode.proto](#lightnode_lightnode-proto)
    - [StreamChunkAvailabilityReply](#lightnode-StreamChunkAvailabilityReply)
    - [StreamChunkAvailabilityRequest](#lightnode-StreamChunkAvailabilityRequest)
  
    - [LightNode](#lightnode-LightNode)
  
- [common/common.proto](#common_common-proto)
    - [BlobCommitment](#common-BlobCommitment)
    - [G1Commitment](#common-G1Commitment)
    - [PaymentHeader](#common-PaymentHeader)
  
- [Scalar Value Types](#scalar-value-types)



<a name="lightnode_lightnode-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lightnode/lightnode.proto



<a name="lightnode-StreamChunkAvailabilityReply"></a>

### StreamChunkAvailabilityReply
A reply to a StreamAvailabilityStatus request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| header_hash | [bytes](#bytes) |  | The hash of a blob header corresponding to a chunk the agent received and verified. From the light node&#39;s perspective, the blob is available if all chunks the light node wants to sample are available. |






<a name="lightnode-StreamChunkAvailabilityRequest"></a>

### StreamChunkAvailabilityRequest
A request from a DA node to an agent light node to stream the availability status of all chunks
assigned to the light node.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| authentication_token | [bytes](#bytes) |  |  |





 

 

 


<a name="lightnode-LightNode"></a>

### LightNode


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| StreamBlobAvailability | [StreamChunkAvailabilityRequest](#lightnode-StreamChunkAvailabilityRequest) | [StreamChunkAvailabilityReply](#lightnode-StreamChunkAvailabilityReply) stream | StreamBlobAvailability streams the availability status blobs from the light node&#39;s perspective. A light node considers a blob to be available if all chunks it wants to sample are available. This API is for use by a DA node for monitoring the availability of chunks through its constellation of agent light nodes. |

 



<a name="common_common-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## common/common.proto



<a name="common-BlobCommitment"></a>

### BlobCommitment
BlobCommitment represents commitment of a specific blob, containing its
KZG commitment, degree proof, the actual degree, and data length in number of symbols.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| commitment | [bytes](#bytes) |  |  |
| length_commitment | [bytes](#bytes) |  |  |
| length_proof | [bytes](#bytes) |  |  |
| length | [uint32](#uint32) |  |  |






<a name="common-G1Commitment"></a>

### G1Commitment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| x | [bytes](#bytes) |  | The X coordinate of the KZG commitment. This is the raw byte representation of the field element. |
| y | [bytes](#bytes) |  | The Y coordinate of the KZG commitment. This is the raw byte representation of the field element. |






<a name="common-PaymentHeader"></a>

### PaymentHeader



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| account_id | [string](#string) |  |  |
| bin_index | [uint32](#uint32) |  |  |
| cumulative_payment | [bytes](#bytes) |  |  |





 

 

 

 



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

