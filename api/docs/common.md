# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [common/common.proto](#common_common-proto)
    - [BlobCommitment](#common-BlobCommitment)
    - [G1Commitment](#common-G1Commitment)
    - [PaymentHeader](#common-PaymentHeader)
  
- [Scalar Value Types](#scalar-value-types)



<a name="common_common-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## common/common.proto



<a name="common-BlobCommitment"></a>

### BlobCommitment
BlobCommitment represents commitment of a specific blob, containing its
KZG commitment, degree proof, the actual degree, and data length in number of symbols (field elements).
It deserializes into https://github.com/Layr-Labs/eigenda/blob/ce89dab18d2f8f55004002e17dd3a18529277845/encoding/data.go#L27

See https://github.com/Layr-Labs/eigenda/blob/master/docs/spec/attestation/encoding.md#validation-via-kzg
to understand how this commitment is used to validate the blob.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| commitment | [bytes](#bytes) |  | Concatenation of the x and y coordinates of `common.G1Commitment`. |
| length_commitment | [bytes](#bytes) |  | A commitment to the blob data with G2 SRS, used to work with length_proof such that the claimed length below is verifiable. |
| length_proof | [bytes](#bytes) |  | A proof that the degree of the polynomial used to generate the blob commitment is valid. It is computed such that the coefficient of the polynomial is committing with the G2 SRS at the end of the highest order. |
| length | [uint32](#uint32) |  | The length of the blob in symbols (field elements), which must be a power of 2. This also specifies the degree of the polynomial used to generate the blob commitment, since length = degree &#43; 1. |






<a name="common-G1Commitment"></a>

### G1Commitment
G1Commitment represents the serialized coordinates of a G1 KZG commitment.
We use gnark-crypto so adopt its serialization, which is big-endian. See:
https://github.com/Consensys/gnark-crypto/blob/779e884dabb38b92e677f4891286637a3d2e5734/ecc/bn254/fp/element.go#L862


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| x | [bytes](#bytes) |  | The X coordinate of the KZG commitment. This is the raw byte representation of the field element. |
| y | [bytes](#bytes) |  | The Y coordinate of the KZG commitment. This is the raw byte representation of the field element. |






<a name="common-PaymentHeader"></a>

### PaymentHeader
PaymentHeader contains payment information for a blob.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| account_id | [string](#string) |  | The account ID of the disperser client. This account ID is an eth wallet address of the user, corresponding to the key used by the client to sign the BlobHeader. |
| reservation_period | [uint32](#uint32) |  | The reservation period of the dispersal request. |
| cumulative_payment | [bytes](#bytes) |  | The cumulative payment of the dispersal request. |
| salt | [uint32](#uint32) |  | The salt of the disperser request. This is used to ensure that the payment header is intentionally unique. |





 

 

 

 



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

