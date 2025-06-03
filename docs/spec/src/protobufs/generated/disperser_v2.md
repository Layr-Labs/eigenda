# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [disperser/v2/disperser_v2.proto](#disperser_v2_disperser_v2-proto)
    - [Attestation](#disperser-v2-Attestation)
    - [BlobCommitmentReply](#disperser-v2-BlobCommitmentReply)
    - [BlobCommitmentRequest](#disperser-v2-BlobCommitmentRequest)
    - [BlobInclusionInfo](#disperser-v2-BlobInclusionInfo)
    - [BlobStatusReply](#disperser-v2-BlobStatusReply)
    - [BlobStatusRequest](#disperser-v2-BlobStatusRequest)
    - [DisperseBlobReply](#disperser-v2-DisperseBlobReply)
    - [DisperseBlobRequest](#disperser-v2-DisperseBlobRequest)
    - [GetPaymentStateForAllQuorumsReply](#disperser-v2-GetPaymentStateForAllQuorumsReply)
    - [GetPaymentStateForAllQuorumsReply.PeriodRecordsEntry](#disperser-v2-GetPaymentStateForAllQuorumsReply-PeriodRecordsEntry)
    - [GetPaymentStateForAllQuorumsReply.ReservationsEntry](#disperser-v2-GetPaymentStateForAllQuorumsReply-ReservationsEntry)
    - [GetPaymentStateForAllQuorumsRequest](#disperser-v2-GetPaymentStateForAllQuorumsRequest)
    - [GetPaymentStateReply](#disperser-v2-GetPaymentStateReply)
    - [GetPaymentStateRequest](#disperser-v2-GetPaymentStateRequest)
    - [PaymentGlobalParams](#disperser-v2-PaymentGlobalParams)
    - [PaymentQuorumConfig](#disperser-v2-PaymentQuorumConfig)
    - [PaymentQuorumProtocolConfig](#disperser-v2-PaymentQuorumProtocolConfig)
    - [PaymentVaultParams](#disperser-v2-PaymentVaultParams)
    - [PaymentVaultParams.QuorumPaymentConfigsEntry](#disperser-v2-PaymentVaultParams-QuorumPaymentConfigsEntry)
    - [PaymentVaultParams.QuorumProtocolConfigsEntry](#disperser-v2-PaymentVaultParams-QuorumProtocolConfigsEntry)
    - [PeriodRecord](#disperser-v2-PeriodRecord)
    - [PeriodRecords](#disperser-v2-PeriodRecords)
    - [QuorumReservation](#disperser-v2-QuorumReservation)
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
| quorum_signed_percentages | [bytes](#bytes) |  | The attestation rate for each quorum. Each quorum&#39;s signing percentage is represented by an 8 bit unsigned integer. The integer is the fraction of the quorum that has signed, with 100 representing 100% of the quorum signing, and 0 representing 0% of the quorum signing. The first byte in the byte array corresponds to the first quorum in the quorum_numbers array, the second byte corresponds to the second quorum, and so on. |






<a name="disperser-v2-BlobCommitmentReply"></a>

### BlobCommitmentReply
The result of a BlobCommitmentRequest().


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob_commitment | [common.BlobCommitment](#common-BlobCommitment) |  | The commitment of the blob. |






<a name="disperser-v2-BlobCommitmentRequest"></a>

### BlobCommitmentRequest
The input for a BlobCommitmentRequest().
This can be used to construct a BlobHeader.commitment.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob | [bytes](#bytes) |  | The blob data to compute the commitment for. |






<a name="disperser-v2-BlobInclusionInfo"></a>

### BlobInclusionInfo
BlobInclusionInfo is the information needed to verify the inclusion of a blob in a batch.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob_certificate | [common.v2.BlobCertificate](#common-v2-BlobCertificate) |  |  |
| blob_index | [uint32](#uint32) |  | blob_index is the index of the blob in the batch |
| inclusion_proof | [bytes](#bytes) |  | inclusion_proof is the inclusion proof of the blob in the batch |






<a name="disperser-v2-BlobStatusReply"></a>

### BlobStatusReply
BlobStatusReply is the reply to a BlobStatusRequest.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [BlobStatus](#disperser-v2-BlobStatus) |  | The status of the blob. |
| signed_batch | [SignedBatch](#disperser-v2-SignedBatch) |  | The signed batch. Only set if the blob status is GATHERING_SIGNATURES or COMPLETE. signed_batch and blob_inclusion_info are only set if the blob status is GATHERING_SIGNATURES or COMPLETE. When blob is in GATHERING_SIGNATURES status, the attestation object in signed_batch contains attestation information at the point in time. As it gathers more signatures, attestation object will be updated according to the latest attestation status. The client can use this intermediate attestation to verify a blob if it has gathered enough signatures. Otherwise, it should should poll the GetBlobStatus API until the desired level of attestation has been gathered or status is COMPLETE. When blob is in COMPLETE status, the attestation object in signed_batch contains the final attestation information. If the final attestation does not meet the client&#39;s requirement, the client should try a new dispersal. |
| blob_inclusion_info | [BlobInclusionInfo](#disperser-v2-BlobInclusionInfo) |  | BlobInclusionInfo is the information needed to verify the inclusion of a blob in a batch. Only set if the blob status is GATHERING_SIGNATURES or COMPLETE. |






<a name="disperser-v2-BlobStatusRequest"></a>

### BlobStatusRequest
BlobStatusRequest is used to query the status of a blob.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob_key | [bytes](#bytes) |  | The unique identifier for the blob. |






<a name="disperser-v2-DisperseBlobReply"></a>

### DisperseBlobReply
A reply to a DisperseBlob request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| result | [BlobStatus](#disperser-v2-BlobStatus) |  | The status of the blob associated with the blob key. |
| blob_key | [bytes](#bytes) |  | The unique 32 byte identifier for the blob.

The blob_key is the keccak hash of the rlp serialization of the BlobHeader, as computed here: https://github.com/Layr-Labs/eigenda/blob/0f14d1c90b86d29c30ff7e92cbadf2762c47f402/core/v2/serialization.go#L30 The blob_key must thus be unique for every request, even if the same blob is being dispersed. Meaning the blob_header must be different for each request.

Note that attempting to disperse a blob with the same blob key as a previously dispersed blob may cause the disperser to reject the blob (DisperseBlob() RPC will return an error). |






<a name="disperser-v2-DisperseBlobRequest"></a>

### DisperseBlobRequest
A request to disperse a blob.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob | [bytes](#bytes) |  | The blob to be dispersed.

The size of this byte array may be any size as long as it does not exceed the maximum length of 16MiB. While the data being dispersed is only required to be greater than 0 bytes, the blob size charged against the payment method will be rounded up to the nearest multiple of `minNumSymbols` defined by the payment vault contract (https://github.com/Layr-Labs/eigenda/blob/1430d56258b4e814b388e497320fd76354bfb478/contracts/src/payments/PaymentVaultStorage.sol#L9).

Every 32 bytes of data is interpreted as an integer in big endian format where the lower address has more significant bits. The integer must stay in the valid range to be interpreted as a field element on the bn254 curve. The valid range is 0 &lt;= x &lt; 21888242871839275222246405745257275088548364400416034343698204186575808495617. If any one of the 32 bytes elements is outside the range, the whole request is deemed as invalid, and rejected. |
| blob_header | [common.v2.BlobHeader](#common-v2-BlobHeader) |  | The header contains metadata about the blob.

This header can be thought of as an &#34;eigenDA tx&#34;, in that it plays a purpose similar to an eth_tx to disperse a 4844 blob. Note that a call to DisperseBlob requires the blob and the blobHeader, which is similar to how dispersing a blob to ethereum requires sending a tx whose data contains the hash of the kzg commit of the blob, which is dispersed separately. |
| signature | [bytes](#bytes) |  | signature over keccak hash of the blob_header that can be verified by blob_header.payment_header.account_id |






<a name="disperser-v2-GetPaymentStateForAllQuorumsReply"></a>

### GetPaymentStateForAllQuorumsReply
GetPaymentStateForAllQuorumsReply contains the payment state of an account. EigenLabs disperser is the only disperser that allows
for ondemand usages, and it will provide the latest on-demand offchain payment records for the request account.
Other dispersers will refuse to serve ondemand requests and serve 0 for off-chain on-demand payment usage (`cumulative_payment`). A client using
non-EigenDA dispersers should only request with reserved usages and disregard the cumulative_payment shared by the non EigenLabs dispersers.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| payment_vault_params | [PaymentVaultParams](#disperser-v2-PaymentVaultParams) |  | payment vault parameters with per-quorum configurations |
| period_records | [GetPaymentStateForAllQuorumsReply.PeriodRecordsEntry](#disperser-v2-GetPaymentStateForAllQuorumsReply-PeriodRecordsEntry) | repeated | period_records maps quorum IDs to the off-chain account reservation usage records for the current and next two periods |
| reservations | [GetPaymentStateForAllQuorumsReply.ReservationsEntry](#disperser-v2-GetPaymentStateForAllQuorumsReply-ReservationsEntry) | repeated | reservations maps quorum IDs to the on-chain account reservation record |
| cumulative_payment | [bytes](#bytes) |  | off-chain on-demand payment usage |
| onchain_cumulative_payment | [bytes](#bytes) |  | on-chain on-demand payment deposited |






<a name="disperser-v2-GetPaymentStateForAllQuorumsReply-PeriodRecordsEntry"></a>

### GetPaymentStateForAllQuorumsReply.PeriodRecordsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [uint32](#uint32) |  |  |
| value | [PeriodRecords](#disperser-v2-PeriodRecords) |  |  |






<a name="disperser-v2-GetPaymentStateForAllQuorumsReply-ReservationsEntry"></a>

### GetPaymentStateForAllQuorumsReply.ReservationsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [uint32](#uint32) |  |  |
| value | [QuorumReservation](#disperser-v2-QuorumReservation) |  |  |






<a name="disperser-v2-GetPaymentStateForAllQuorumsRequest"></a>

### GetPaymentStateForAllQuorumsRequest
GetPaymentStateForAllQuorumsRequest contains parameters to query the payment state of an account.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| account_id | [string](#string) |  | The ID of the account being queried. This account ID is an eth wallet address of the user. |
| timestamp | [uint64](#uint64) |  | Timestamp of the request in nanoseconds since the Unix epoch. If too far out of sync with the server&#39;s clock, request may be rejected. |
| signature | [bytes](#bytes) |  | Signature over the payment account ID and timestamp. |






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
| account_id | [string](#string) |  | The ID of the account being queried. This account ID is an eth wallet address of the user. |
| signature | [bytes](#bytes) |  | Signature over the account ID and timestamp. |
| timestamp | [uint64](#uint64) |  | Timestamp of the request in nanoseconds since the Unix epoch. If too far out of sync with the server&#39;s clock, request may be rejected. |






<a name="disperser-v2-PaymentGlobalParams"></a>

### PaymentGlobalParams
Global constant parameters defined by the payment vault.
This message type will soon be deprecated in replacement of PaymentVaultParams. During endpoint migration, this will be filled
with the parameters on quorum 0, quorum configurations will be the same across quorums for the foreseeable future.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| global_symbols_per_second | [uint64](#uint64) |  | Global ratelimit for on-demand dispersals |
| min_num_symbols | [uint64](#uint64) |  | Minimum number of symbols accounted for all dispersals |
| price_per_symbol | [uint64](#uint64) |  | Price charged per symbol for on-demand dispersals |
| reservation_window | [uint64](#uint64) |  | Reservation window for all reservations |
| on_demand_quorum_numbers | [uint32](#uint32) | repeated | quorums allowed to make on-demand dispersals |






<a name="disperser-v2-PaymentQuorumConfig"></a>

### PaymentQuorumConfig
PaymentQuorumConfig contains the configuration for a quorum&#39;s payment configurations


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| reservation_symbols_per_second | [uint64](#uint64) |  | reservation_symbols_per_second is the total symbols per second that can be reserved for this quorum |
| on_demand_symbols_per_second | [uint64](#uint64) |  | on_demand_symbols_per_second is the symbols per second allowed for on-demand payments for this quorum |
| on_demand_price_per_symbol | [uint64](#uint64) |  | on_demand_price_per_symbol is the price per symbol for on-demand payments in wei |






<a name="disperser-v2-PaymentQuorumProtocolConfig"></a>

### PaymentQuorumProtocolConfig
PaymentQuorumProtocolConfig contains the configuration for a quorum&#39;s protocol-level configurations


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| min_num_symbols | [uint64](#uint64) |  | min_num_symbols is the minimum number of symbols that must be charged for any request |
| reservation_advance_window | [uint64](#uint64) |  | reservation_advance_window is the window in seconds before a reservation starts that it can be activated |
| reservation_rate_limit_window | [uint64](#uint64) |  | reservation_rate_limit_window is the time window in seconds for reservation rate limiting |
| on_demand_rate_limit_window | [uint64](#uint64) |  | on_demand_rate_limit_window is the time window in seconds for on-demand rate limiting |
| on_demand_enabled | [bool](#bool) |  | on_demand_enabled indicates whether on-demand payments are enabled for this quorum |






<a name="disperser-v2-PaymentVaultParams"></a>

### PaymentVaultParams
PaymentVaultParams contains the global payment configuration parameters from the payment vault
This is the new version of


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| quorum_payment_configs | [PaymentVaultParams.QuorumPaymentConfigsEntry](#disperser-v2-PaymentVaultParams-QuorumPaymentConfigsEntry) | repeated | quorum_payment_configs maps quorum IDs to their payment configurations |
| quorum_protocol_configs | [PaymentVaultParams.QuorumProtocolConfigsEntry](#disperser-v2-PaymentVaultParams-QuorumProtocolConfigsEntry) | repeated | quorum_protocol_configs maps quorum IDs to their protocol configurations |
| on_demand_quorum_numbers | [uint32](#uint32) | repeated | on_demand_quorum_numbers lists the quorum numbers that support on-demand payments |






<a name="disperser-v2-PaymentVaultParams-QuorumPaymentConfigsEntry"></a>

### PaymentVaultParams.QuorumPaymentConfigsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [uint32](#uint32) |  |  |
| value | [PaymentQuorumConfig](#disperser-v2-PaymentQuorumConfig) |  |  |






<a name="disperser-v2-PaymentVaultParams-QuorumProtocolConfigsEntry"></a>

### PaymentVaultParams.QuorumProtocolConfigsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [uint32](#uint32) |  |  |
| value | [PaymentQuorumProtocolConfig](#disperser-v2-PaymentQuorumProtocolConfig) |  |  |






<a name="disperser-v2-PeriodRecord"></a>

### PeriodRecord
PeriodRecord is the usage record of an account in a bin. The API should return the active bin
record and the subsequent two records that contains potential overflows.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| index | [uint32](#uint32) |  | Period index of the reservation |
| usage | [uint64](#uint64) |  | symbol usage recorded |






<a name="disperser-v2-PeriodRecords"></a>

### PeriodRecords
An array of period records. Typically this is used to include 3 records, from the current period to the next two periods.
The next two period records are included because they may include spillage usages from the previous period or the current period.
The client should be aware of the spillage so they account for them as they disperse during those periods.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| records | [PeriodRecord](#disperser-v2-PeriodRecord) | repeated |  |






<a name="disperser-v2-QuorumReservation"></a>

### QuorumReservation
Reservation parameters of an account, used to determine the rate limit for the account.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| symbols_per_second | [uint64](#uint64) |  | rate limit for the account |
| start_timestamp | [uint32](#uint32) |  | start timestamp of the reservation |
| end_timestamp | [uint32](#uint32) |  | end timestamp of the reservation |






<a name="disperser-v2-Reservation"></a>

### Reservation
Reservation parameters of an account, used to determine the rate limit for the account.
This message type will soon be deprecated. During the migration time, we will maintain the usage by returning the
most restrictive reservation parameters across the quroums: symbols_per_second will be the lowest rate of across all quroums
with latest start and earliest end timestamp, all the quorum numbers with a reservation, and a dummy quorum_splits which was never used.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| symbols_per_second | [uint64](#uint64) |  | rate limit for the account |
| start_timestamp | [uint32](#uint32) |  | start timestamp of the reservation |
| end_timestamp | [uint32](#uint32) |  | end timestamp of the reservation |
| quorum_numbers | [uint32](#uint32) | repeated | quorums allowed to make reserved dispersals |
| quorum_splits | [uint32](#uint32) | repeated | quorum splits describes how the payment is split among the quorums |






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
Intermediate states are states that the blob can be in while being processed, and it can be updated to a different state:
- QUEUED
- ENCODED
- GATHERING_SIGNATURES
Terminal states are states that will not be updated to a different state:
- UNKNOWN
- COMPLETE
- FAILED

| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN | 0 | UNKNOWN means that the status of the blob is unknown. This is a catch all and should not be encountered absent a bug.

This status is functionally equivalent to FAILED, but is used to indicate that the failure is due to an unanticipated bug. |
| QUEUED | 1 | QUEUED means that the blob has been queued by the disperser for processing. The DisperseBlob API is asynchronous, meaning that after request validation, but before any processing, the blob is stored in a queue of some sort, and a response immediately returned to the client. |
| ENCODED | 2 | ENCODED means that the blob has been Reed-Solomon encoded into chunks and is ready to be dispersed to DA Nodes. |
| GATHERING_SIGNATURES | 3 | GATHERING_SIGNATURES means that the blob chunks are currently actively being transmitted to validators, and in doing so requesting that the validators sign to acknowledge receipt of the blob. Requests that timeout or receive errors are resubmitted to DA nodes for some period of time set by the disperser, after which the BlobStatus becomes COMPLETE. |
| COMPLETE | 4 | COMPLETE means the blob has been dispersed to DA nodes, and the GATHERING_SIGNATURES period of time has completed. This status does not guarantee any signer percentage, so a client should check that the signature has met its required threshold, and resubmit a new blob dispersal request if not. |
| FAILED | 5 | FAILED means that the blob has failed permanently. Note that this is a terminal state, and in order to retry the blob, the client must submit the blob again (blob key is required to be unique). |


 

 


<a name="disperser-v2-Disperser"></a>

### Disperser
Disperser defines the public APIs for dispersing blobs.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| DisperseBlob | [DisperseBlobRequest](#disperser-v2-DisperseBlobRequest) | [DisperseBlobReply](#disperser-v2-DisperseBlobReply) | DisperseBlob accepts blob to disperse from clients. This executes the dispersal asynchronously, i.e. it returns once the request is accepted. The client could use GetBlobStatus() API to poll the the processing status of the blob. |
| GetBlobStatus | [BlobStatusRequest](#disperser-v2-BlobStatusRequest) | [BlobStatusReply](#disperser-v2-BlobStatusReply) | GetBlobStatus is meant to be polled for the blob status. |
| GetBlobCommitment | [BlobCommitmentRequest](#disperser-v2-BlobCommitmentRequest) | [BlobCommitmentReply](#disperser-v2-BlobCommitmentReply) | GetBlobCommitment is a utility method that calculates commitment for a blob payload. It is provided to help clients who are trying to construct a DisperseBlobRequest.blob_header and don&#39;t have the ability to calculate the commitment themselves (expensive operation which requires SRS points).

For an example usage, see how our disperser_client makes a call to this endpoint when it doesn&#39;t have a local prover: https://github.com/Layr-Labs/eigenda/blob/6059c6a068298d11c41e50f5bcd208d0da44906a/api/clients/v2/disperser_client.go#L166 |
| GetPaymentState | [GetPaymentStateRequest](#disperser-v2-GetPaymentStateRequest) | [GetPaymentStateReply](#disperser-v2-GetPaymentStateReply) | GetPaymentState is a utility method to get the payment state of a given account, at a given disperser. EigenDA&#39;s payment system for v2 is currently centralized, meaning that each disperser does its own accounting. As reservation moves to be quorum specific and served by permissionless dispersers, GetPaymentState will soon be deprecated in replacement of GetPaymentStateForAllQuorums to include more specifications. During the endpoint migration time, the response uses quorum 0 for the global parameters, and the most retrictive reservation parameters of a user across quorums. For OnDemand, EigenDA disperser is the only allowed disperser, so it will provide real values tracked for on-demand offchain payment records. For other dispersers, they will refuse to serve ondemand requests and serve 0 as the on-demand offchain records. A client using non-EigenDA dispersers should only request with reserved usages.

A client wanting to disperse a blob would thus need to synchronize its local accounting state with that of the disperser. That typically only needs to be done once, and the state can be updated locally as the client disperses blobs. The accounting rules are simple and can be updated locally, but periodic checks with the disperser can&#39;t hurt.

For an example usage, see how our disperser_client makes a call to this endpoint to populate its local accountant struct: https://github.com/Layr-Labs/eigenda/blob/6059c6a068298d11c41e50f5bcd208d0da44906a/api/clients/v2/disperser_client.go#L298 |
| GetPaymentStateForAllQuorums | [GetPaymentStateForAllQuorumsRequest](#disperser-v2-GetPaymentStateForAllQuorumsRequest) | [GetPaymentStateForAllQuorumsReply](#disperser-v2-GetPaymentStateForAllQuorumsReply) | GetPaymentStateForAllQuorums is a utility method to get the payment state of a given account, at a given disperser. EigenDA&#39;s dispersers and validators each does its own accounting for reservation usages, indexed by the account and quorum id. A client wanting to disperse a blob would thus need to synchronize its local accounting state with the disperser it plans to disperse to. That typically only needs to be done once, and the state can be updated locally as the client disperses blobs. The accounting rules are simple and can be updated locally, but periodic checks with the disperser can&#39;t hurt.

For an example usage, see how our disperser_client makes a call to this endpoint to populate its local accountant struct: https://github.com/Layr-Labs/eigenda/blob/6059c6a068298d11c41e50f5bcd208d0da44906a/api/clients/v2/disperser_client.go#L298 |

 



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

