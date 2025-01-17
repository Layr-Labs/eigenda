# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [common/v2/common_v2.proto](#common_v2_common_v2-proto)
    - [Batch](#common-v2-Batch)
    - [BatchHeader](#common-v2-BatchHeader)
    - [BlobCertificate](#common-v2-BlobCertificate)
    - [BlobHeader](#common-v2-BlobHeader)
  
- [Scalar Value Types](#scalar-value-types)



<a name="common_v2_common_v2-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## common/v2/common_v2.proto



<a name="common-v2-Batch"></a>

### Batch
Batch is a batch of blob certificates


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| header | [BatchHeader](#common-v2-BatchHeader) |  | header contains metadata about the batch |
| blob_certificates | [BlobCertificate](#common-v2-BlobCertificate) | repeated | blob_certificates is the list of blob certificates in the batch |






<a name="common-v2-BatchHeader"></a>

### BatchHeader
BatchHeader is the header of a batch of blobs


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| batch_root | [bytes](#bytes) |  | batch_root is the root of the merkle tree of the hashes of blob certificates in the batch |
| reference_block_number | [uint64](#uint64) |  | reference_block_number is the block number that the state of the batch is based on for attestation |






<a name="common-v2-BlobCertificate"></a>

### BlobCertificate
BlobCertificate contains a full description of a blob and how it is dispersed. Part of the certificate
is provided by the blob submitter (i.e. the blob header), and part is provided by the disperser (i.e. the relays).
Validator nodes eventually sign the blob certificate once they are in custody of the required chunks
(note that the signature is indirect; validators sign the hash of a Batch, which contains the blob certificate).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob_header | [BlobHeader](#common-v2-BlobHeader) |  | blob_header contains data about the blob. |
| relays | [uint32](#uint32) | repeated | relays is the list of relays that are in custody of the blob. The relays custodying the data are chosen by the Disperser to which the DisperseBlob request was submitted. It needs to contain at least 1 relay number. To retrieve a blob from the relay, one can find that relay&#39;s URL in the EigenDARelayRegistry contract: https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/core/EigenDARelayRegistry.sol |






<a name="common-v2-BlobHeader"></a>

### BlobHeader
BlobHeader contains the information describing a blob and the way it is to be dispersed.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| version | [uint32](#uint32) |  | The blob version. Blob versions are pushed onchain by EigenDA governance in an append only fashion and store the maximum number of operators, number of chunks, and coding rate for a blob. On blob verification, these values are checked against supplied or default security thresholds to validate the security assumptions of the blob&#39;s availability. |
| quorum_numbers | [uint32](#uint32) | repeated | quorum_numbers is the list of quorum numbers that the blob is part of. All quorums must be specified (including required quorums).

The following quorums are currently required: - 0: ETH - 1: EIGEN |
| commitment | [common.BlobCommitment](#common-BlobCommitment) |  | commitment is the KZG commitment to the blob |
| payment_header | [common.PaymentHeader](#common-PaymentHeader) |  | payment_header contains payment information for the blob |
| signature | [bytes](#bytes) |  | signature over keccak hash of the blob_header that can be verified by blob_header.account_id |





 

 

 

 



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

