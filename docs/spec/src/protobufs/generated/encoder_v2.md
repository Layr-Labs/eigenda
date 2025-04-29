# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [encoder/v2/encoder_v2.proto](#encoder_v2_encoder_v2-proto)
    - [EncodeBlobReply](#encoder-v2-EncodeBlobReply)
    - [EncodeBlobRequest](#encoder-v2-EncodeBlobRequest)
    - [EncodingParams](#encoder-v2-EncodingParams)
    - [FragmentInfo](#encoder-v2-FragmentInfo)
  
    - [Encoder](#encoder-v2-Encoder)
  
- [Scalar Value Types](#scalar-value-types)



<a name="encoder_v2_encoder_v2-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## encoder/v2/encoder_v2.proto



<a name="encoder-v2-EncodeBlobReply"></a>

### EncodeBlobReply
EncodeBlobReply contains metadata about the encoded chunks


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| fragment_info | [FragmentInfo](#encoder-v2-FragmentInfo) |  |  |






<a name="encoder-v2-EncodeBlobRequest"></a>

### EncodeBlobRequest
EncodeBlobRequest contains the reference to the blob to be encoded and the encoding parameters
determined by the control plane.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob_key | [bytes](#bytes) |  |  |
| encoding_params | [EncodingParams](#encoder-v2-EncodingParams) |  |  |
| blob_size | [uint64](#uint64) |  |  |






<a name="encoder-v2-EncodingParams"></a>

### EncodingParams
EncodingParams specifies how the blob should be encoded into chunks


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| chunk_length | [uint64](#uint64) |  |  |
| num_chunks | [uint64](#uint64) |  |  |






<a name="encoder-v2-FragmentInfo"></a>

### FragmentInfo
FragmentInfo contains metadata about the encoded fragments


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| total_chunk_size_bytes | [uint32](#uint32) |  |  |
| fragment_size_bytes | [uint32](#uint32) |  |  |





 

 

 


<a name="encoder-v2-Encoder"></a>

### Encoder


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| EncodeBlob | [EncodeBlobRequest](#encoder-v2-EncodeBlobRequest) | [EncodeBlobReply](#encoder-v2-EncodeBlobReply) | EncodeBlob encodes a blob into chunks using specified encoding parameters. The blob is retrieved using the provided blob key and the encoded chunks are persisted for later retrieval. |

 



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

