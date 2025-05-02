# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [encoder/encoder.proto](#encoder_encoder-proto)
    - [BlobCommitment](#encoder-BlobCommitment)
    - [EncodeBlobReply](#encoder-EncodeBlobReply)
    - [EncodeBlobRequest](#encoder-EncodeBlobRequest)
    - [EncodingParams](#encoder-EncodingParams)
  
    - [ChunkEncodingFormat](#encoder-ChunkEncodingFormat)
  
    - [Encoder](#encoder-Encoder)
  
- [Scalar Value Types](#scalar-value-types)



<a name="encoder_encoder-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## encoder/encoder.proto



<a name="encoder-BlobCommitment"></a>

### BlobCommitment
BlobCommitments contains the blob&#39;s commitment, degree proof, and the actual degree
DEPRECATED: use common.BlobCommitment instead


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| commitment | [bytes](#bytes) |  |  |
| length_commitment | [bytes](#bytes) |  |  |
| length_proof | [bytes](#bytes) |  |  |
| length | [uint32](#uint32) |  |  |






<a name="encoder-EncodeBlobReply"></a>

### EncodeBlobReply
EncodeBlobReply returns all encoded chunks along with BlobCommitment for the same,
where Chunk is the smallest unit that is distributed to DA nodes


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| commitment | [BlobCommitment](#encoder-BlobCommitment) |  |  |
| chunks | [bytes](#bytes) | repeated |  |
| chunk_encoding_format | [ChunkEncodingFormat](#encoder-ChunkEncodingFormat) |  | How the above chunks are encoded. |






<a name="encoder-EncodeBlobRequest"></a>

### EncodeBlobRequest
EncodeBlobRequest contains data and pre-computed encoding params provided to Encoder


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [bytes](#bytes) |  |  |
| encoding_params | [EncodingParams](#encoder-EncodingParams) |  |  |






<a name="encoder-EncodingParams"></a>

### EncodingParams
Parameters needed by Encoder for encoding


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| chunk_length | [uint32](#uint32) |  |  |
| num_chunks | [uint32](#uint32) |  |  |





 


<a name="encoder-ChunkEncodingFormat"></a>

### ChunkEncodingFormat


| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN | 0 |  |
| GNARK | 1 |  |
| GOB | 2 |  |


 

 


<a name="encoder-Encoder"></a>

### Encoder


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| EncodeBlob | [EncodeBlobRequest](#encoder-EncodeBlobRequest) | [EncodeBlobReply](#encoder-EncodeBlobReply) |  |

 



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

