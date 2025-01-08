# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [disperser/disperser.proto](#disperser_disperser-proto)
    - [AuthenticatedReply](#disperser-AuthenticatedReply)
    - [AuthenticatedRequest](#disperser-AuthenticatedRequest)
    - [AuthenticationData](#disperser-AuthenticationData)
    - [BatchHeader](#disperser-BatchHeader)
    - [BatchMetadata](#disperser-BatchMetadata)
    - [BlobAuthHeader](#disperser-BlobAuthHeader)
    - [BlobHeader](#disperser-BlobHeader)
    - [BlobInfo](#disperser-BlobInfo)
    - [BlobQuorumParam](#disperser-BlobQuorumParam)
    - [BlobStatusReply](#disperser-BlobStatusReply)
    - [BlobStatusRequest](#disperser-BlobStatusRequest)
    - [BlobVerificationProof](#disperser-BlobVerificationProof)
    - [DisperseBlobReply](#disperser-DisperseBlobReply)
    - [DisperseBlobRequest](#disperser-DisperseBlobRequest)
    - [DispersePaidBlobRequest](#disperser-DispersePaidBlobRequest)
    - [RetrieveBlobReply](#disperser-RetrieveBlobReply)
    - [RetrieveBlobRequest](#disperser-RetrieveBlobRequest)
  
    - [BlobStatus](#disperser-BlobStatus)
  
    - [Disperser](#disperser-Disperser)
  
- [Scalar Value Types](#scalar-value-types)



<a name="disperser_disperser-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## disperser/disperser.proto



<a name="disperser-AuthenticatedReply"></a>

### AuthenticatedReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob_auth_header | [BlobAuthHeader](#disperser-BlobAuthHeader) |  |  |
| disperse_reply | [DisperseBlobReply](#disperser-DisperseBlobReply) |  |  |






<a name="disperser-AuthenticatedRequest"></a>

### AuthenticatedRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| disperse_request | [DisperseBlobRequest](#disperser-DisperseBlobRequest) |  |  |
| authentication_data | [AuthenticationData](#disperser-AuthenticationData) |  |  |






<a name="disperser-AuthenticationData"></a>

### AuthenticationData
AuthenticationData contains the signature of the BlobAuthHeader.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| authentication_data | [bytes](#bytes) |  |  |






<a name="disperser-BatchHeader"></a>

### BatchHeader



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| batch_root | [bytes](#bytes) |  | The root of the merkle tree with the hashes of blob headers as leaves. |
| quorum_numbers | [bytes](#bytes) |  | All quorums associated with blobs in this batch. Sorted in ascending order. Ex. [0, 2, 1] =&gt; 0x000102 |
| quorum_signed_percentages | [bytes](#bytes) |  | The percentage of stake that has signed for this batch. The quorum_signed_percentages[i] is percentage for the quorum_numbers[i]. |
| reference_block_number | [uint32](#uint32) |  | The Ethereum block number at which the batch was created. The Disperser will encode and disperse the blobs based on the onchain info (e.g. operator stakes) at this block number. |






<a name="disperser-BatchMetadata"></a>

### BatchMetadata



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| batch_header | [BatchHeader](#disperser-BatchHeader) |  |  |
| signatory_record_hash | [bytes](#bytes) |  | The hash of all public keys of the operators that did not sign the batch. |
| fee | [bytes](#bytes) |  | The fee payment paid by users for dispersing this batch. It&#39;s the bytes representation of a big.Int value. |
| confirmation_block_number | [uint32](#uint32) |  | The Ethereum block number at which the batch is confirmed onchain. |
| batch_header_hash | [bytes](#bytes) |  | This is the hash of the ReducedBatchHeader defined onchain, see: https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/interfaces/IEigenDAServiceManager.sol#L43 The is the message that the operators will sign their signatures on. |






<a name="disperser-BlobAuthHeader"></a>

### BlobAuthHeader
BlobAuthHeader contains information about the blob for the client to verify and sign.
- Once payments are enabled, the BlobAuthHeader will contain the KZG commitment to the blob, which the client
will verify and sign. Having the client verify the KZG commitment instead of calculating it avoids
the need for the client to have the KZG structured reference string (SRS), which can be large.
The signed KZG commitment prevents the disperser from sending a different blob to the DA Nodes
than the one the client sent.
- In the meantime, the BlobAuthHeader contains a simple challenge parameter is used to prevent
replay attacks in the event that a signature is leaked.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| challenge_parameter | [uint32](#uint32) |  |  |






<a name="disperser-BlobHeader"></a>

### BlobHeader



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| commitment | [common.G1Commitment](#common-G1Commitment) |  | KZG commitment of the blob. |
| data_length | [uint32](#uint32) |  | The length of the blob in symbols (each symbol is 32 bytes). |
| blob_quorum_params | [BlobQuorumParam](#disperser-BlobQuorumParam) | repeated | The params of the quorums that this blob participates in. |






<a name="disperser-BlobInfo"></a>

### BlobInfo
BlobInfo contains information needed to confirm the blob against the EigenDA contracts


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob_header | [BlobHeader](#disperser-BlobHeader) |  |  |
| blob_verification_proof | [BlobVerificationProof](#disperser-BlobVerificationProof) |  |  |






<a name="disperser-BlobQuorumParam"></a>

### BlobQuorumParam



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| quorum_number | [uint32](#uint32) |  | The ID of the quorum. |
| adversary_threshold_percentage | [uint32](#uint32) |  | The max percentage of stake within the quorum that can be held by or delegated to adversarial operators. Currently, this and the next parameter are standardized across the quorum using values read from the EigenDA contracts. |
| confirmation_threshold_percentage | [uint32](#uint32) |  | The min percentage of stake that must attest in order to consider the dispersal is successful. |
| chunk_length | [uint32](#uint32) |  | The length of each chunk. |






<a name="disperser-BlobStatusReply"></a>

### BlobStatusReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [BlobStatus](#disperser-BlobStatus) |  | The status of the blob. |
| info | [BlobInfo](#disperser-BlobInfo) |  | The blob info needed for clients to confirm the blob against the EigenDA contracts. |






<a name="disperser-BlobStatusRequest"></a>

### BlobStatusRequest
BlobStatusRequest is used to query the status of a blob.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| request_id | [bytes](#bytes) |  |  |






<a name="disperser-BlobVerificationProof"></a>

### BlobVerificationProof



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| batch_id | [uint32](#uint32) |  | batch_id is an incremental ID assigned to a batch by EigenDAServiceManager |
| blob_index | [uint32](#uint32) |  | The index of the blob in the batch (which is logically an ordered list of blobs). |
| batch_metadata | [BatchMetadata](#disperser-BatchMetadata) |  |  |
| inclusion_proof | [bytes](#bytes) |  | inclusion_proof is a merkle proof for a blob header&#39;s inclusion in a batch |
| quorum_indexes | [bytes](#bytes) |  | indexes of quorums in BatchHeader.quorum_numbers that match the quorums in BlobHeader.blob_quorum_params Ex. BlobHeader.blob_quorum_params = [ 	{ 		quorum_number = 0, 		... 	}, 	{ 		quorum_number = 3, 		... 	}, 	{ 		quorum_number = 5, 		... 	}, ] BatchHeader.quorum_numbers = [0, 5, 3] =&gt; 0x000503 Then, quorum_indexes = [0, 2, 1] =&gt; 0x000201 |






<a name="disperser-DisperseBlobReply"></a>

### DisperseBlobReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| result | [BlobStatus](#disperser-BlobStatus) |  | The status of the blob associated with the request_id. Will always be PROCESSING. |
| request_id | [bytes](#bytes) |  | The request ID generated by the disperser. Once a request is accepted (although not processed), a unique request ID will be generated. Two different DisperseBlobRequests (determined by the hash of the DisperseBlobRequest) will have different IDs, and the same DisperseBlobRequest sent repeatedly at different times will also have different IDs. The client should use this ID to query the processing status of the request (via the GetBlobStatus API). |






<a name="disperser-DisperseBlobRequest"></a>

### DisperseBlobRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [bytes](#bytes) |  | The data to be dispersed. The size of data must be &lt;= 16MiB. Every 32 bytes of data is interpreted as an integer in big endian format where the lower address has more significant bits. The integer must stay in the valid range to be interpreted as a field element on the bn254 curve. The valid range is 0 &lt;= x &lt; 21888242871839275222246405745257275088548364400416034343698204186575808495617 If any one of the 32 bytes elements is outside the range, the whole request is deemed as invalid, and rejected. |
| custom_quorum_numbers | [uint32](#uint32) | repeated | The quorums to which the blob will be sent, in addition to the required quorums which are configured on the EigenDA smart contract. If required quorums are included here, an error will be returned. The disperser will ensure that the encoded blobs for each quorum are all processed within the same batch. |
| account_id | [string](#string) |  | The account ID of the client. This should be a hex-encoded string of the ECSDA public key corresponding to the key used by the client to sign the BlobAuthHeader. |






<a name="disperser-DispersePaidBlobRequest"></a>

### DispersePaidBlobRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [bytes](#bytes) |  | The data to be dispersed. Same requirements as DisperseBlobRequest. |
| quorum_numbers | [uint32](#uint32) | repeated | The quorums to which the blob to be sent |
| payment_header | [common.PaymentHeader](#common-PaymentHeader) |  | Payment header contains account_id, reservation_period, cumulative_payment, and salt |
| payment_signature | [bytes](#bytes) |  | signature of payment_header |






<a name="disperser-RetrieveBlobReply"></a>

### RetrieveBlobReply
RetrieveBlobReply contains the retrieved blob data


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [bytes](#bytes) |  |  |






<a name="disperser-RetrieveBlobRequest"></a>

### RetrieveBlobRequest
RetrieveBlobRequest contains parameters to retrieve the blob.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| batch_header_hash | [bytes](#bytes) |  |  |
| blob_index | [uint32](#uint32) |  |  |





 


<a name="disperser-BlobStatus"></a>

### BlobStatus
BlobStatus represents the status of a blob.
The status of a blob is updated as the blob is processed by the disperser.
The status of a blob can be queried by the client using the GetBlobStatus API.
Intermediate states are states that the blob can be in while being processed, and it can be updated to a differet state:
- PROCESSING
- DISPERSING
- CONFIRMED
Terminal states are states that will not be updated to a different state:
- FAILED
- FINALIZED
- INSUFFICIENT_SIGNATURES

| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN | 0 |  |
| PROCESSING | 1 | PROCESSING means that the blob is currently being processed by the disperser |
| CONFIRMED | 2 | CONFIRMED means that the blob has been dispersed to DA Nodes and the dispersed batch containing the blob has been confirmed onchain |
| FAILED | 3 | FAILED means that the blob has failed permanently (for reasons other than insufficient signatures, which is a separate state). This status is somewhat of a catch-all category, containg (but not necessarily exclusively as errors can be added in the future): - blob has expired - internal logic error while requesting encoding - blob retry has exceeded its limit while waiting for blob finalization after confirmation. Most likely triggered by a chain reorg: see https://github.com/Layr-Labs/eigenda/blob/master/disperser/batcher/finalizer.go#L179-L189. |
| FINALIZED | 4 | FINALIZED means that the block containing the blob&#39;s confirmation transaction has been finalized on Ethereum |
| INSUFFICIENT_SIGNATURES | 5 | INSUFFICIENT_SIGNATURES means that the confirmation threshold for the blob was not met for at least one quorum. |
| DISPERSING | 6 | The DISPERSING state is comprised of two separate phases: - Dispersing to DA nodes and collecting signature - Submitting the transaction on chain and waiting for tx receipt |


 

 


<a name="disperser-Disperser"></a>

### Disperser
Disperser defines the public APIs for dispersing blobs.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| DisperseBlob | [DisperseBlobRequest](#disperser-DisperseBlobRequest) | [DisperseBlobReply](#disperser-DisperseBlobReply) | DisperseBlob accepts a single blob to be dispersed. This executes the dispersal async, i.e. it returns once the request is accepted. The client should use GetBlobStatus() API to poll the processing status of the blob.

If DisperseBlob returns the following error codes: INVALID_ARGUMENT (400): request is invalid for a reason specified in the error msg. RESOURCE_EXHAUSTED (429): request is rate limited for the quorum specified in the error msg. user should retry after the specified duration. INTERNAL (500): serious error, user should NOT retry. |
| DisperseBlobAuthenticated | [AuthenticatedRequest](#disperser-AuthenticatedRequest) stream | [AuthenticatedReply](#disperser-AuthenticatedReply) stream | DisperseBlobAuthenticated is similar to DisperseBlob, except that it requires the client to authenticate itself via the AuthenticationData message. The protocol is as follows: 1. The client sends a DisperseBlobAuthenticated request with the DisperseBlobRequest message 2. The Disperser sends back a BlobAuthHeader message containing information for the client to verify and sign. 3. The client verifies the BlobAuthHeader and sends back the signed BlobAuthHeader in an 	 AuthenticationData message. 4. The Disperser verifies the signature and returns a DisperseBlobReply message. |
| GetBlobStatus | [BlobStatusRequest](#disperser-BlobStatusRequest) | [BlobStatusReply](#disperser-BlobStatusReply) | This API is meant to be polled for the blob status. |
| RetrieveBlob | [RetrieveBlobRequest](#disperser-RetrieveBlobRequest) | [RetrieveBlobReply](#disperser-RetrieveBlobReply) | This retrieves the requested blob from the Disperser&#39;s backend. This is a more efficient way to retrieve blobs than directly retrieving from the DA Nodes (see detail about this approach in api/proto/retriever/retriever.proto). The blob should have been initially dispersed via this Disperser service for this API to work. |

 



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

