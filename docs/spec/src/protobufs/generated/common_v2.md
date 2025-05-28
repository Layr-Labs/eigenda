# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [common/v2/common_v2.proto](#common_v2_common_v2-proto)
    - [Batch](#common-v2-Batch)
    - [BatchHeader](#common-v2-BatchHeader)
    - [BlobCertificate](#common-v2-BlobCertificate)
    - [BlobHeader](#common-v2-BlobHeader)
    - [PaymentHeader](#common-v2-PaymentHeader)
    - [PaymentQuorumConfig](#common-v2-PaymentQuorumConfig)
    - [PaymentQuorumProtocolConfig](#common-v2-PaymentQuorumProtocolConfig)
    - [PaymentVaultParams](#common-v2-PaymentVaultParams)
    - [PaymentVaultParams.QuorumPaymentConfigsEntry](#common-v2-PaymentVaultParams-QuorumPaymentConfigsEntry)
    - [PaymentVaultParams.QuorumProtocolConfigsEntry](#common-v2-PaymentVaultParams-QuorumProtocolConfigsEntry)
  
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
| signature | [bytes](#bytes) |  | signature is an ECDSA signature signed by the blob request signer&#39;s account ID over the BlobHeader&#39;s blobKey, which is a keccak hash of the serialized BlobHeader, and used to verify against blob dispersal request&#39;s account ID |
| relay_keys | [uint32](#uint32) | repeated | relay_keys is the list of relay keys that are in custody of the blob. The relays custodying the data are chosen by the Disperser to which the DisperseBlob request was submitted. It needs to contain at least 1 relay number. To retrieve a blob from the relay, one can find that relay&#39;s URL in the EigenDARelayRegistry contract: https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/core/EigenDARelayRegistry.sol |






<a name="common-v2-BlobHeader"></a>

### BlobHeader
BlobHeader contains the information describing a blob and the way it is to be dispersed.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| version | [uint32](#uint32) |  | The blob version. Blob versions are pushed onchain by EigenDA governance in an append only fashion and store the maximum number of operators, number of chunks, and coding rate for a blob. On blob verification, these values are checked against supplied or default security thresholds to validate the security assumptions of the blob&#39;s availability. |
| quorum_numbers | [uint32](#uint32) | repeated | quorum_numbers is the list of quorum numbers that the blob is part of. Each quorum will store the data, hence adding quorum numbers adds redundancy, making the blob more likely to be retrievable. Each quorum requires separate payment.

On-demand dispersal is currently limited to using a subset of the following quorums: - 0: ETH - 1: EIGEN

Reserved-bandwidth dispersal is free to use multiple quorums, however those must be reserved ahead of time. The quorum_numbers specified here must be a subset of the ones allowed by the on-chain reservation. Check the allowed quorum numbers by looking up reservation struct: https://github.com/Layr-Labs/eigenda/blob/1430d56258b4e814b388e497320fd76354bfb478/contracts/src/interfaces/IPaymentVault.sol#L10 |
| commitment | [common.BlobCommitment](#common-BlobCommitment) |  | commitment is the KZG commitment to the blob |
| payment_header | [PaymentHeader](#common-v2-PaymentHeader) |  | payment_header contains payment information for the blob |






<a name="common-v2-PaymentHeader"></a>

### PaymentHeader
PaymentHeader contains payment information for a blob, which is crucial for validating and processing dispersal requests.
The PaymentHeader is designed to support two distinct payment methods within the EigenDA protocol:

1. Reservation-based payment system:
   This system allows users to reserve bandwidth in advance for a specified time period. It&#39;s designed for
   users who need predictable throughput with a fixed ratelimit bin in required or custom quorums.
   Under this method, the user pre-arranges a reservation with specific parameters on the desired quorums:
   - symbolsPerSecond: The rate at which they can disperse data
   - startTimestamp and endTimestamp: The timeframe during which the reservation is active

2. On-demand payment system:
   This is a pay-as-you-go model where users deposit funds into the PaymentVault contract and
   payments are deducted as they make dispersal requests. This system is more flexible but has
   more restrictions on which quorums can be used (currently limited to quorums 0 and 1).

The disperser client always attempts to use a reservation-based payment first if one exists for the account.
If no valid reservation exists or if the reservation doesn&#39;t have enough remaining bandwidth,
the client will fall back to on-demand payment, provided the user has deposited sufficient funds
in the PaymentVault contract.

The distinction between these two payment methods is made by examining:
- For reservation-based: The timestamp must be within an active reservation period, and cumulative_payment is zero or empty
- For on-demand: The cumulative_payment field contains a non-zero value representing the total payment for all dispersals

Every dispersal request is metered based on the size of the data being dispersed, rounded up to the
nearest multiple of the minNumSymbols parameter defined in the PaymentVault contract. The size is calculated as:
symbols_charged = ceiling(blob_size / minNumSymbols) * minNumSymbols
On-demand payments take a step further by calculating the specific cost
cost = symbols_charged * price_per_symbol

Security and Authentication:
The payment header is protected by a cryptographic signature that covers the entire BlobHeader.
This signature is verified during request processing to ensure that:
1. The request is genuinely from the holder of the private key corresponding to account_id
2. The payment information hasn&#39;t been tampered with
3. The same request isn&#39;t being resubmitted (replay protection)

This signature verification happens in core/auth/v2/authenticator.go where:
- The BlobKey (a hash of the serialized BlobHeader) is computed
- The signature is verified against this key
- The recovered public key is checked against the account_id in the payment header

Once a payment has been processed and the signature verified, the disperser server will not
roll back the payment or usage records, even if subsequent processing fails. This design choice
prevents double-spending and ensures payment integrity.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| account_id | [string](#string) |  | The account ID of the disperser client, represented as an Ethereum wallet address in hex format (e.g., &#34;0x1234...abcd&#34;). This field is critical for both payment methods as it:

1. Identifies whose reservation to check for reservation-based payments 2. Identifies whose on-chain deposit balance to check for on-demand payments 3. Provides the address against which the BlobHeader signature is verified

The account_id has special significance in the authentication flow: - When a client signs a BlobHeader, they use their private key - The disperser server recovers the public key from this signature - The recovered public key is converted to an Ethereum address - This derived address must exactly match the account_id in this field

This verification process (implemented in core/auth/v2/authenticator.go&#39;s AuthenticateBlobRequest method) ensures that only the legitimate owner of the account can submit dispersal requests charged to that account. It prevents unauthorized payments or impersonation attacks where someone might try to use another user&#39;s reservation or on-chain balance.

The account_id is typically set by the client&#39;s Accountant when constructing the PaymentMetadata (see api/clients/v2/accountant.go - AccountBlob method). |
| timestamp | [int64](#int64) |  | The timestamp represents the UNIX timestamp in nanoseconds at the time the dispersal request is created. This high-precision timestamp serves multiple critical functions in the protocol:

For reservation-based payments: 1. Reservation Period Determination: The timestamp is used to calculate which reservation period the request belongs to using the formula: reservation_period = floor(timestamp_ns / (reservationPeriodInterval_s * 1e9)) * reservationPeriodInterval_s where reservationPeriodInterval_s is in seconds, and the result is in seconds.

2. Reservation Validity Check: The timestamp must fall within an active reservation window: - It must be &gt;= the reservation&#39;s startTimestamp (in seconds) - It must be &lt; the reservation&#39;s endTimestamp (in seconds)

3. Period Window Check: The server validates that the request&#39;s reservation period is either: - The current period (based on server time) - The immediately previous period This prevents requests with future timestamps or very old timestamps.

4. Rate Limiting: The server uses the timestamp to allocate the request to the appropriate rate-limiting bucket. Each reservation period has a fixed bandwidth limit (symbolsPerSecond * reservationPeriodInterval).

For on-demand payments: 1. Replay Protection: The timestamp helps ensure each request is unique and prevent replay attacks.

2. Global Ratelimiting (TO BE IMPLEMENTED): Treating all on-demand requests as an user-agnostic more frequent reservation, timestamp is checked against the OnDemandSymbolsPerSecond and OnDemandPeriodInterval.

The timestamp is typically acquired by calling time.Now().UnixNano() in Go and accounted for NTP offsets by periodically syncing with a configuratble NTP server endpoint. The client&#39;s Accountant component (api/clients/v2/accountant.go) expects the caller to provide this timestamp, which it then uses to determine the correct reservation period and check bandwidth availability. |
| cumulative_payment | [bytes](#bytes) |  | The cumulative_payment field is a serialized uint256 big integer representing the total amount of tokens paid by the requesting account across all their dispersal requests, including the current one. The unit is in wei. This field is exclusively used for on-demand payments and should be zero or empty for reservation-based payments. If this field is zero or empty, disperser server&#39;s meterer will treat this request as reservation-based. For the current implementation, the choice of quorum doesn&#39;t affect the payment calculations. A client may choose to use any or all of the required quorums.

Detailed Payment Mechanics: 1. Cumulative Design: Rather than sending incremental payment amounts, the protocol uses a cumulative approach where each request states the total amount paid by the account so far. This design: - Prevents double-spending even with concurrent requests - Simplifies verification logic - Requests are enforced by a strictly increasing order

2. Calculation Formula: For a new dispersal request, the cumulative_payment is calculated as: new_cumulative = previous_cumulative &#43; (symbols_charged * price_per_symbol)

 Where: - previous_cumulative: The highest cumulative payment value from previous dispersals - symbols_charged: The blob size rounded up to the nearest multiple of minNumSymbols - price_per_symbol: The cost per symbol set in the PaymentVault contract

3. Validation Process: When the disperser receives a request with a cumulative_payment, it performs multiple validations: - Checks that the on-chain deposit balance in the PaymentVault is sufficient to cover this payment - Verifies the cumulative_payment is greater than the highest previous payment from this account - Verifies the increase from the previous cumulative payment is appropriate for the blob size - If other requests from the same account are currently processing, ensures this new cumulative value is consistent with those (preventing double-spending)

4. On-chain Implementation: The PaymentVault contract maintains: - A deposit balance for each account - Global parameters including minNumSymbols, GlobalSymbolsPerSecond and pricePerSymbol Due to the use of cumulative payments, if a client loses track of their current cumulative payment value, they can query the disperser server for their current payment state using the GetPaymentState RPC. |






<a name="common-v2-PaymentQuorumConfig"></a>

### PaymentQuorumConfig
PaymentQuorumConfig contains the configuration for a quorum&#39;s payment configurations


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| reservation_symbols_per_second | [uint64](#uint64) |  | reservation_symbols_per_second is the total symbols per second that can be reserved for this quorum |
| on_demand_symbols_per_second | [uint64](#uint64) |  | on_demand_symbols_per_second is the symbols per second allowed for on-demand payments for this quorum |
| on_demand_price_per_symbol | [uint64](#uint64) |  | on_demand_price_per_symbol is the price per symbol for on-demand payments in wei |






<a name="common-v2-PaymentQuorumProtocolConfig"></a>

### PaymentQuorumProtocolConfig
PaymentQuorumProtocolConfig contains the configuration for a quorum&#39;s protocol-level configurations


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| min_num_symbols | [uint64](#uint64) |  | min_num_symbols is the minimum number of symbols that must be charged for any request |
| reservation_advance_window | [uint64](#uint64) |  | reservation_advance_window is the window in seconds before a reservation starts that it can be activated |
| reservation_rate_limit_window | [uint64](#uint64) |  | reservation_rate_limit_window is the time window in seconds for reservation rate limiting |
| on_demand_rate_limit_window | [uint64](#uint64) |  | on_demand_rate_limit_window is the time window in seconds for on-demand rate limiting |
| on_demand_enabled | [bool](#bool) |  | on_demand_enabled indicates whether on-demand payments are enabled for this quorum |






<a name="common-v2-PaymentVaultParams"></a>

### PaymentVaultParams
PaymentVaultParams contains the global payment configuration parameters from the payment vault


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| quorum_payment_configs | [PaymentVaultParams.QuorumPaymentConfigsEntry](#common-v2-PaymentVaultParams-QuorumPaymentConfigsEntry) | repeated | quorum_payment_configs maps quorum IDs to their payment configurations |
| quorum_protocol_configs | [PaymentVaultParams.QuorumProtocolConfigsEntry](#common-v2-PaymentVaultParams-QuorumProtocolConfigsEntry) | repeated | quorum_protocol_configs maps quorum IDs to their protocol configurations |
| on_demand_quorum_numbers | [uint32](#uint32) | repeated | on_demand_quorum_numbers lists the quorum numbers that support on-demand payments |






<a name="common-v2-PaymentVaultParams-QuorumPaymentConfigsEntry"></a>

### PaymentVaultParams.QuorumPaymentConfigsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [uint32](#uint32) |  |  |
| value | [PaymentQuorumConfig](#common-v2-PaymentQuorumConfig) |  |  |






<a name="common-v2-PaymentVaultParams-QuorumProtocolConfigsEntry"></a>

### PaymentVaultParams.QuorumProtocolConfigsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [uint32](#uint32) |  |  |
| value | [PaymentQuorumProtocolConfig](#common-v2-PaymentQuorumProtocolConfig) |  |  |





 

 

 

 



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

