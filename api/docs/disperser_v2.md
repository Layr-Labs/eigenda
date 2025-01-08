# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [disperser/v2/disperser_v2.proto](#disperser_v2_disperser_v2-proto)
    - [Attestation](#disperser-v2-Attestation)
    - [BlobCommitmentReply](#disperser-v2-BlobCommitmentReply)
    - [BlobCommitmentRequest](#disperser-v2-BlobCommitmentRequest)
    - [BlobStatusReply](#disperser-v2-BlobStatusReply)
    - [BlobStatusRequest](#disperser-v2-BlobStatusRequest)
    - [BlobVerificationInfo](#disperser-v2-BlobVerificationInfo)
    - [DisperseBlobReply](#disperser-v2-DisperseBlobReply)
    - [DisperseBlobRequest](#disperser-v2-DisperseBlobRequest)
    - [GetPaymentStateReply](#disperser-v2-GetPaymentStateReply)
    - [GetPaymentStateRequest](#disperser-v2-GetPaymentStateRequest)
    - [PaymentGlobalParams](#disperser-v2-PaymentGlobalParams)
    - [PeriodRecord](#disperser-v2-PeriodRecord)
    - [Reservation](#disperser-v2-Reservation)
    - [SignedBatch](#disperser-v2-SignedBatch)
  
    - [BlobStatus](#disperser-v2-BlobStatus)
  
    - [Disperser](#disperser-v2-Disperser)
  
- [Scalar Value Types](#scalar-value-types)



<a name="disperser_v2_disperser_v2-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## disperser/v2/disperser_v2.proto



<a name="disperser-v2-Attestation"></a>

### Attestation



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| non_signer_pubkeys | [bytes](#bytes) | repeated | Serialized bytes of non signer public keys (G1 points) |
| apk_g2 | [bytes](#bytes) |  | Serialized bytes of G2 point that represents aggregate public key of all signers |
| quorum_apks | [bytes](#bytes) | repeated | Serialized bytes of aggregate public keys (G1 points) from all nodes for each quorum The order of the quorum_apks should match the order of the quorum_numbers |
| sigma | [bytes](#bytes) |  | Serialized bytes of aggregate signature |
| quorum_numbers | [uint32](#uint32) | repeated | Relevant quorum numbers for the attestation |
| quorum_signed_percentages | [bytes](#bytes) |  | The attestation rate for each quorum. The order of the quorum_signed_percentages should match the order of the quorum_numbers |






<a name="disperser-v2-BlobCommitmentReply"></a>

### BlobCommitmentReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob_commitment | [common.BlobCommitment](#common-BlobCommitment) |  |  |






<a name="disperser-v2-BlobCommitmentRequest"></a>

### BlobCommitmentRequest
Utility method used to generate the commitment of blob given its data.
This can be used to construct BlobHeader.commitment


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [bytes](#bytes) |  |  |






<a name="disperser-v2-BlobStatusReply"></a>

### BlobStatusReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [BlobStatus](#disperser-v2-BlobStatus) |  | The status of the blob. |
| signed_batch | [SignedBatch](#disperser-v2-SignedBatch) |  | The signed batch |
| blob_verification_info | [BlobVerificationInfo](#disperser-v2-BlobVerificationInfo) |  |  |






<a name="disperser-v2-BlobStatusRequest"></a>

### BlobStatusRequest
BlobStatusRequest is used to query the status of a blob.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob_key | [bytes](#bytes) |  |  |






<a name="disperser-v2-BlobVerificationInfo"></a>

### BlobVerificationInfo
BlobVerificationInfo is the information needed to verify the inclusion of a blob in a batch.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob_certificate | [common.v2.BlobCertificate](#common-v2-BlobCertificate) |  |  |
| blob_index | [uint32](#uint32) |  | blob_index is the index of the blob in the batch |
| inclusion_proof | [bytes](#bytes) |  | inclusion_proof is the inclusion proof of the blob in the batch |






<a name="disperser-v2-DisperseBlobReply"></a>

### DisperseBlobReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| result | [BlobStatus](#disperser-v2-BlobStatus) |  | The status of the blob associated with the blob key. |
| blob_key | [bytes](#bytes) |  |  |






<a name="disperser-v2-DisperseBlobRequest"></a>

### DisperseBlobRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [bytes](#bytes) |  | The data to be dispersed. The size of data must be &lt;= 16MiB. Every 32 bytes of data is interpreted as an integer in big endian format where the lower address has more significant bits. The integer must stay in the valid range to be interpreted as a field element on the bn254 curve. The valid range is 0 &lt;= x &lt; 21888242871839275222246405745257275088548364400416034343698204186575808495617 If any one of the 32 bytes elements is outside the range, the whole request is deemed as invalid, and rejected. |
| blob_header | [common.v2.BlobHeader](#common-v2-BlobHeader) |  |  |






<a name="disperser-v2-GetPaymentStateReply"></a>

### GetPaymentStateReply
GetPaymentStateReply contains the payment state of an account.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| payment_global_params | [PaymentGlobalParams](#disperser-v2-PaymentGlobalParams) |  | global payment vault parameters |
| period_records | [PeriodRecord](#disperser-v2-PeriodRecord) | repeated | off-chain account reservation usage records |
| reservation | [Reservation](#disperser-v2-Reservation) |  | on-chain account reservation setting |
| cumulative_payment | [bytes](#bytes) |  | off-chain on-demand payment usage |
| onchain_cumulative_payment | [bytes](#bytes) |  | on-chain on-demand payment deposited |






<a name="disperser-v2-GetPaymentStateRequest"></a>

### GetPaymentStateRequest
GetPaymentStateRequest contains parameters to query the payment state of an account.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| account_id | [string](#string) |  |  |
| signature | [bytes](#bytes) |  | Signature over the account ID TODO: sign over a reservation period or a nonce to mitigate signature replay attacks |






<a name="disperser-v2-PaymentGlobalParams"></a>

### PaymentGlobalParams



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| global_symbols_per_second | [uint64](#uint64) |  |  |
| min_num_symbols | [uint32](#uint32) |  |  |
| price_per_symbol | [uint32](#uint32) |  |  |
| reservation_window | [uint32](#uint32) |  |  |
| on_demand_quorum_numbers | [uint32](#uint32) | repeated |  |






<a name="disperser-v2-PeriodRecord"></a>

### PeriodRecord
PeriodRecord is the usage record of an account in a bin. The API should return the active bin 
record and the subsequent two records that contains potential overflows.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| index | [uint32](#uint32) |  |  |
| usage | [uint64](#uint64) |  |  |






<a name="disperser-v2-Reservation"></a>

### Reservation



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| symbols_per_second | [uint64](#uint64) |  |  |
| start_timestamp | [uint32](#uint32) |  |  |
| end_timestamp | [uint32](#uint32) |  |  |
| quorum_numbers | [uint32](#uint32) | repeated |  |
| quorum_splits | [uint32](#uint32) | repeated |  |






<a name="disperser-v2-SignedBatch"></a>

### SignedBatch
SignedBatch is a batch of blobs with a signature.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| header | [common.v2.BatchHeader](#common-v2-BatchHeader) |  | header contains metadata about the batch |
| attestation | [Attestation](#disperser-v2-Attestation) |  | attestation on the batch |





 


<a name="disperser-v2-BlobStatus"></a>

### BlobStatus
BlobStatus represents the status of a blob.
The status of a blob is updated as the blob is processed by the disperser.
The status of a blob can be queried by the client using the GetBlobStatus API.
Intermediate states are states that the blob can be in while being processed, and it can be updated to a differet state:
- QUEUED
- ENCODED
Terminal states are states that will not be updated to a different state:
- CERTIFIED
- FAILED
- INSUFFICIENT_SIGNATURES

| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN | 0 |  |
| QUEUED | 1 | QUEUED means that the blob has been queued by the disperser for processing |
| ENCODED | 2 | ENCODED means that the blob has been encoded and is ready to be dispersed to DA Nodes |
| CERTIFIED | 3 | CERTIFIED means the blob has been dispersed and attested by the DA nodes |
| FAILED | 4 | FAILED means that the blob has failed permanently |
| INSUFFICIENT_SIGNATURES | 5 | INSUFFICIENT_SIGNATURES means that the blob has failed to gather sufficient attestation |


 

 


<a name="disperser-v2-Disperser"></a>

### Disperser
Disperser defines the public APIs for dispersing blobs.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| DisperseBlob | [DisperseBlobRequest](#disperser-v2-DisperseBlobRequest) | [DisperseBlobReply](#disperser-v2-DisperseBlobReply) | DisperseBlob accepts blob to disperse from clients. This executes the dispersal asynchronously, i.e. it returns once the request is accepted. The client could use GetBlobStatus() API to poll the the processing status of the blob. |
| GetBlobStatus | [BlobStatusRequest](#disperser-v2-BlobStatusRequest) | [BlobStatusReply](#disperser-v2-BlobStatusReply) | GetBlobStatus is meant to be polled for the blob status. |
| GetBlobCommitment | [BlobCommitmentRequest](#disperser-v2-BlobCommitmentRequest) | [BlobCommitmentReply](#disperser-v2-BlobCommitmentReply) | GetBlobCommitment is a utility method that calculates commitment for a blob payload. |
| GetPaymentState | [GetPaymentStateRequest](#disperser-v2-GetPaymentStateRequest) | [GetPaymentStateReply](#disperser-v2-GetPaymentStateReply) | GetPaymentState is a utility method to get the payment state of a given account. |

 



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

