# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [churner/churner.proto](#churner_churner-proto)
    - [ChurnReply](#churner-ChurnReply)
    - [ChurnRequest](#churner-ChurnRequest)
    - [OperatorToChurn](#churner-OperatorToChurn)
    - [SignatureWithSaltAndExpiry](#churner-SignatureWithSaltAndExpiry)
  
    - [Churner](#churner-Churner)
  
- [common/common.proto](#common_common-proto)
    - [BlobCommitment](#common-BlobCommitment)
    - [G1Commitment](#common-G1Commitment)
    - [PaymentHeader](#common-PaymentHeader)
  
- [common/v2/common.proto](#common_v2_common-proto)
    - [Batch](#common-v2-Batch)
    - [BatchHeader](#common-v2-BatchHeader)
    - [BlobCertificate](#common-v2-BlobCertificate)
    - [BlobHeader](#common-v2-BlobHeader)
  
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
  
- [node/node.proto](#node_node-proto)
    - [AttestBatchReply](#node-AttestBatchReply)
    - [AttestBatchRequest](#node-AttestBatchRequest)
    - [BatchHeader](#node-BatchHeader)
    - [Blob](#node-Blob)
    - [BlobHeader](#node-BlobHeader)
    - [BlobQuorumInfo](#node-BlobQuorumInfo)
    - [Bundle](#node-Bundle)
    - [G2Commitment](#node-G2Commitment)
    - [GetBlobHeaderReply](#node-GetBlobHeaderReply)
    - [GetBlobHeaderRequest](#node-GetBlobHeaderRequest)
    - [MerkleProof](#node-MerkleProof)
    - [NodeInfoReply](#node-NodeInfoReply)
    - [NodeInfoRequest](#node-NodeInfoRequest)
    - [RetrieveChunksReply](#node-RetrieveChunksReply)
    - [RetrieveChunksRequest](#node-RetrieveChunksRequest)
    - [StoreBlobsReply](#node-StoreBlobsReply)
    - [StoreBlobsRequest](#node-StoreBlobsRequest)
    - [StoreChunksReply](#node-StoreChunksReply)
    - [StoreChunksRequest](#node-StoreChunksRequest)
  
    - [ChunkEncodingFormat](#node-ChunkEncodingFormat)
  
    - [Dispersal](#node-Dispersal)
    - [Retrieval](#node-Retrieval)
  
- [node/v2/node_v2.proto](#node_v2_node_v2-proto)
    - [GetChunksReply](#node-v2-GetChunksReply)
    - [GetChunksRequest](#node-v2-GetChunksRequest)
    - [NodeInfoReply](#node-v2-NodeInfoReply)
    - [NodeInfoRequest](#node-v2-NodeInfoRequest)
    - [StoreChunksReply](#node-v2-StoreChunksReply)
    - [StoreChunksRequest](#node-v2-StoreChunksRequest)
  
    - [Dispersal](#node-v2-Dispersal)
    - [Retrieval](#node-v2-Retrieval)
  
- [relay/relay.proto](#relay_relay-proto)
    - [ChunkRequest](#relay-ChunkRequest)
    - [ChunkRequestByIndex](#relay-ChunkRequestByIndex)
    - [ChunkRequestByRange](#relay-ChunkRequestByRange)
    - [GetBlobReply](#relay-GetBlobReply)
    - [GetBlobRequest](#relay-GetBlobRequest)
    - [GetChunksReply](#relay-GetChunksReply)
    - [GetChunksRequest](#relay-GetChunksRequest)
  
    - [Relay](#relay-Relay)
  
- [retriever/retriever.proto](#retriever_retriever-proto)
    - [BlobReply](#retriever-BlobReply)
    - [BlobRequest](#retriever-BlobRequest)
  
    - [Retriever](#retriever-Retriever)
  
- [retriever/v2/retriever.proto](#retriever_v2_retriever-proto)
    - [BlobReply](#retriever-v2-BlobReply)
    - [BlobRequest](#retriever-v2-BlobRequest)
  
    - [Retriever](#retriever-v2-Retriever)
  
- [Scalar Value Types](#scalar-value-types)



<a name="churner_churner-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## churner/churner.proto



<a name="churner-ChurnReply"></a>

### ChurnReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| signature_with_salt_and_expiry | [SignatureWithSaltAndExpiry](#churner-SignatureWithSaltAndExpiry) |  | The signature signed by the Churner. |
| operators_to_churn | [OperatorToChurn](#churner-OperatorToChurn) | repeated | A list of existing operators that get churned out. This list will contain all quorums specified in the ChurnRequest even if some quorums may not have any churned out operators. If a quorum has available space, OperatorToChurn object will contain the quorum ID and empty operator and pubkey. The smart contract should only churn out the operators for quorums that are full.

For example, if the ChurnRequest specifies quorums 0 and 1 where quorum 0 is full and quorum 1 has available space, the ChurnReply will contain two OperatorToChurn objects with the respective quorums. OperatorToChurn for quorum 0 will contain the operator to churn out and OperatorToChurn for quorum 1 will contain empty operator (zero address) and pubkey. The smart contract should only churn out the operators for quorum 0 because quorum 1 has available space without having any operators churned. Note: it&#39;s possible an operator gets churned out just for one or more quorums (rather than entirely churned out for all quorums). |






<a name="churner-ChurnRequest"></a>

### ChurnRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operator_address | [string](#string) |  | The Ethereum address (in hex like &#34;0x123abcdef...&#34;) of the operator. |
| operator_to_register_pubkey_g1 | [bytes](#bytes) |  | The operator making the churn request. |
| operator_to_register_pubkey_g2 | [bytes](#bytes) |  |  |
| operator_request_signature | [bytes](#bytes) |  | The operator&#39;s BLS signature signed on the keccak256 hash of concat(&#34;ChurnRequest&#34;, operator address, g1, g2, salt). |
| salt | [bytes](#bytes) |  | The salt used as part of the message to sign on for operator_request_signature. |
| quorum_ids | [uint32](#uint32) | repeated | The quorums to register for. Note: - If any of the quorum here has already been registered, this entire request will fail to proceed. - If any of the quorum fails to register, this entire request will fail. - Regardless of whether the specified quorums are full or not, the Churner will return parameters for all quorums specified here. The smart contract will determine whether it needs to churn out existing operators based on whether the quorums have available space. The IDs must be in range [0, 254]. |






<a name="churner-OperatorToChurn"></a>

### OperatorToChurn
This describes an operator to churn out for a quorum.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| quorum_id | [uint32](#uint32) |  | The ID of the quorum of the operator to churn out. |
| operator | [bytes](#bytes) |  | The address of the operator. |
| pubkey | [bytes](#bytes) |  | BLS pubkey (G1 point) of the operator. |






<a name="churner-SignatureWithSaltAndExpiry"></a>

### SignatureWithSaltAndExpiry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| signature | [bytes](#bytes) |  | Churner&#39;s signature on the Operator&#39;s attributes. |
| salt | [bytes](#bytes) |  | Salt is the keccak256 hash of concat(&#34;churn&#34;, time.Now(), operatorToChurn&#39;s OperatorID, Churner&#39;s ECDSA private key) |
| expiry | [int64](#int64) |  | When this churn decision will expire. |





 

 

 


<a name="churner-Churner"></a>

### Churner
The Churner is a service that handles churn requests from new operators trying to
join the EigenDA network.
When the EigenDA network reaches the maximum number of operators, any new operator
trying to join will have to make a churn request to this Churner, which acts as the
sole decision maker to decide whether this new operator could join, and if so, which
existing operator will be churned out (so the max number of operators won&#39;t be
exceeded).
The max number of operators, as well as the rules to make churn decisions, are
defined onchain, see details in OperatorSetParam at:
https://github.com/Layr-Labs/eigenlayer-middleware/blob/master/src/interfaces/IBLSRegistryCoordinatorWithIndices.sol#L24.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Churn | [ChurnRequest](#churner-ChurnRequest) | [ChurnReply](#churner-ChurnReply) |  |

 



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
| account_id | [string](#string) |  | The account ID of the disperser client. This should be a hex-encoded string of the ECSDA public key corresponding to the key used by the client to sign the BlobHeader. |
| reservation_period | [uint32](#uint32) |  | The reservation period of the dispersal request. |
| cumulative_payment | [bytes](#bytes) |  | The cumulative payment of the dispersal request. |
| salt | [uint32](#uint32) |  | The salt of the disperser request. This is used to ensure that the payment header is intentionally unique. |





 

 

 

 



<a name="common_v2_common-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## common/v2/common.proto



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
BlobCertificate is what gets attested by the network


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob_header | [BlobHeader](#common-v2-BlobHeader) |  |  |
| relays | [uint32](#uint32) | repeated |  |






<a name="common-v2-BlobHeader"></a>

### BlobHeader



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| version | [uint32](#uint32) |  | Blob version |
| quorum_numbers | [uint32](#uint32) | repeated |  |
| commitment | [common.BlobCommitment](#common-BlobCommitment) |  |  |
| payment_header | [common.PaymentHeader](#common-PaymentHeader) |  |  |
| signature | [bytes](#bytes) |  | signature over keccak hash of the blob_header that can be verified by blob_header.account_id |





 

 

 

 



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

 



<a name="node_node-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## node/node.proto



<a name="node-AttestBatchReply"></a>

### AttestBatchReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| signature | [bytes](#bytes) |  |  |






<a name="node-AttestBatchRequest"></a>

### AttestBatchRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| batch_header | [BatchHeader](#node-BatchHeader) |  | header of the batch |
| blob_header_hashes | [bytes](#bytes) | repeated | the header hashes of all blobs in the batch |






<a name="node-BatchHeader"></a>

### BatchHeader
BatchHeader (see core/data.go#BatchHeader)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| batch_root | [bytes](#bytes) |  | The root of the merkle tree with hashes of blob headers as leaves. |
| reference_block_number | [uint32](#uint32) |  | The Ethereum block number at which the batch is dispersed. |






<a name="node-Blob"></a>

### Blob
In EigenDA, the original blob to disperse is encoded as a polynomial via taking
taking different point evaluations (i.e. erasure coding). These points are split
into disjoint subsets which are assigned to different operator nodes in the EigenDA
network.
The data in this message is a subset of these points that are assigned to a
single operator node.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| header | [BlobHeader](#node-BlobHeader) |  | Which (original) blob this is for. |
| bundles | [Bundle](#node-Bundle) | repeated | Each bundle contains all chunks for a single quorum of the blob. The number of bundles must be equal to the total number of quorums associated with the blob, and the ordering must be the same as BlobHeader.quorum_headers. Note: an operator may be in some but not all of the quorums; in that case the bundle corresponding to that quorum will be empty. |






<a name="node-BlobHeader"></a>

### BlobHeader



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| commitment | [common.G1Commitment](#common-G1Commitment) |  | The KZG commitment to the polynomial representing the blob. |
| length_commitment | [G2Commitment](#node-G2Commitment) |  | The KZG commitment to the polynomial representing the blob on G2, it is used for proving the degree of the polynomial |
| length_proof | [G2Commitment](#node-G2Commitment) |  | The low degree proof. It&#39;s the KZG commitment to the polynomial shifted to the largest SRS degree. |
| length | [uint32](#uint32) |  | The length of the original blob in number of symbols (in the field where the polynomial is defined). |
| quorum_headers | [BlobQuorumInfo](#node-BlobQuorumInfo) | repeated | The params of the quorums that this blob participates in. |
| account_id | [string](#string) |  | The ID of the user who is dispersing this blob to EigenDA. |
| reference_block_number | [uint32](#uint32) |  | The reference block number whose state is used to encode the blob |






<a name="node-BlobQuorumInfo"></a>

### BlobQuorumInfo
See BlobQuorumParam as defined in
api/proto/disperser/disperser.proto


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| quorum_id | [uint32](#uint32) |  |  |
| adversary_threshold | [uint32](#uint32) |  |  |
| confirmation_threshold | [uint32](#uint32) |  |  |
| chunk_length | [uint32](#uint32) |  |  |
| ratelimit | [uint32](#uint32) |  |  |






<a name="node-Bundle"></a>

### Bundle
A Bundle is the collection of chunks associated with a single blob, for a single
operator and a single quorum.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| chunks | [bytes](#bytes) | repeated | Each chunk corresponds to a collection of points on the polynomial. Each chunk has same number of points. |
| bundle | [bytes](#bytes) |  | All chunks of the bundle encoded in a byte array. |






<a name="node-G2Commitment"></a>

### G2Commitment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| x_a0 | [bytes](#bytes) |  | The A0 element of the X coordinate of G2 point. |
| x_a1 | [bytes](#bytes) |  | The A1 element of the X coordinate of G2 point. |
| y_a0 | [bytes](#bytes) |  | The A0 element of the Y coordinate of G2 point. |
| y_a1 | [bytes](#bytes) |  | The A1 element of the Y coordinate of G2 point. |






<a name="node-GetBlobHeaderReply"></a>

### GetBlobHeaderReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob_header | [BlobHeader](#node-BlobHeader) |  | The header of the blob requested per GetBlobHeaderRequest. |
| proof | [MerkleProof](#node-MerkleProof) |  | Merkle proof that returned blob header belongs to the batch and is the batch&#39;s MerkleProof.index-th blob. This can be checked against the batch root on chain. |






<a name="node-GetBlobHeaderRequest"></a>

### GetBlobHeaderRequest
See RetrieveChunksRequest for documentation of each parameter of GetBlobHeaderRequest.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| batch_header_hash | [bytes](#bytes) |  |  |
| blob_index | [uint32](#uint32) |  |  |
| quorum_id | [uint32](#uint32) |  |  |






<a name="node-MerkleProof"></a>

### MerkleProof



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| hashes | [bytes](#bytes) | repeated | The proof itself. |
| index | [uint32](#uint32) |  | Which index (the leaf of the Merkle tree) this proof is for. |






<a name="node-NodeInfoReply"></a>

### NodeInfoReply
Node info reply


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| semver | [string](#string) |  |  |
| arch | [string](#string) |  |  |
| os | [string](#string) |  |  |
| num_cpu | [uint32](#uint32) |  |  |
| mem_bytes | [uint64](#uint64) |  |  |






<a name="node-NodeInfoRequest"></a>

### NodeInfoRequest
Node info request






<a name="node-RetrieveChunksReply"></a>

### RetrieveChunksReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| chunks | [bytes](#bytes) | repeated | All chunks the Node is storing for the requested blob per RetrieveChunksRequest. |
| chunk_encoding_format | [ChunkEncodingFormat](#node-ChunkEncodingFormat) |  | How the above chunks are encoded. |






<a name="node-RetrieveChunksRequest"></a>

### RetrieveChunksRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| batch_header_hash | [bytes](#bytes) |  | The hash of the ReducedBatchHeader defined onchain, see: https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/interfaces/IEigenDAServiceManager.sol#L43 This identifies which batch to retrieve for. |
| blob_index | [uint32](#uint32) |  | Which blob in the batch to retrieve for (note: a batch is logically an ordered list of blobs). |
| quorum_id | [uint32](#uint32) |  | Which quorum of the blob to retrieve for (note: a blob can have multiple quorums and the chunks for different quorums at a Node can be different). The ID must be in range [0, 254]. |






<a name="node-StoreBlobsReply"></a>

### StoreBlobsReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| signatures | [google.protobuf.BytesValue](#google-protobuf-BytesValue) | repeated | The operator&#39;s BLS sgnature signed on the blob header hashes. The ordering of the signatures must match the ordering of the blobs sent in the request, with empty signatures in the places for discarded blobs. |






<a name="node-StoreBlobsRequest"></a>

### StoreBlobsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blobs | [Blob](#node-Blob) | repeated | Blobs to store |
| reference_block_number | [uint32](#uint32) |  | The reference block number whose state is used to encode the blobs |






<a name="node-StoreChunksReply"></a>

### StoreChunksReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| signature | [bytes](#bytes) |  | The operator&#39;s BLS signature signed on the batch header hash. |






<a name="node-StoreChunksRequest"></a>

### StoreChunksRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| batch_header | [BatchHeader](#node-BatchHeader) |  | Which batch this request is for. |
| blobs | [Blob](#node-Blob) | repeated | The chunks for each blob in the batch to be stored in an EigenDA Node. |





 


<a name="node-ChunkEncodingFormat"></a>

### ChunkEncodingFormat
This describes how the chunks returned in RetrieveChunksReply are encoded.
Used to facilitate the decoding of chunks.

| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN | 0 |  |
| GNARK | 1 |  |
| GOB | 2 |  |


 

 


<a name="node-Dispersal"></a>

### Dispersal


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| StoreChunks | [StoreChunksRequest](#node-StoreChunksRequest) | [StoreChunksReply](#node-StoreChunksReply) | StoreChunks validates that the chunks match what the Node is supposed to receive ( different Nodes are responsible for different chunks, as EigenDA is horizontally sharded) and is correctly coded (e.g. each chunk must be a valid KZG multiproof) according to the EigenDA protocol. It also stores the chunks along with metadata for the protocol-defined length of custody. It will return a signature at the end to attest to the data in this request it has processed. |
| StoreBlobs | [StoreBlobsRequest](#node-StoreBlobsRequest) | [StoreBlobsReply](#node-StoreBlobsReply) | StoreBlobs is simiar to StoreChunks, but it stores the blobs using a different storage schema so that the stored blobs can later be aggregated by AttestBatch method to a bigger batch. StoreBlobs &#43; AttestBatch will eventually replace and deprecate StoreChunks method. DEPRECATED: StoreBlobs method is not used |
| AttestBatch | [AttestBatchRequest](#node-AttestBatchRequest) | [AttestBatchReply](#node-AttestBatchReply) | AttestBatch is used to aggregate the batches stored by StoreBlobs method to a bigger batch. It will return a signature at the end to attest to the aggregated batch. DEPRECATED: AttestBatch method is not used |
| NodeInfo | [NodeInfoRequest](#node-NodeInfoRequest) | [NodeInfoReply](#node-NodeInfoReply) | Retrieve node info metadata |


<a name="node-Retrieval"></a>

### Retrieval


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| RetrieveChunks | [RetrieveChunksRequest](#node-RetrieveChunksRequest) | [RetrieveChunksReply](#node-RetrieveChunksReply) | RetrieveChunks retrieves the chunks for a blob custodied at the Node. |
| GetBlobHeader | [GetBlobHeaderRequest](#node-GetBlobHeaderRequest) | [GetBlobHeaderReply](#node-GetBlobHeaderReply) | GetBlobHeader is similar to RetrieveChunks, this just returns the header of the blob. |
| NodeInfo | [NodeInfoRequest](#node-NodeInfoRequest) | [NodeInfoReply](#node-NodeInfoReply) | Retrieve node info metadata |

 



<a name="node_v2_node_v2-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## node/v2/node_v2.proto



<a name="node-v2-GetChunksReply"></a>

### GetChunksReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| chunks | [bytes](#bytes) | repeated | All chunks the Node is storing for the requested blob per RetrieveChunksRequest. |






<a name="node-v2-GetChunksRequest"></a>

### GetChunksRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob_key | [bytes](#bytes) |  |  |
| quorum_id | [uint32](#uint32) |  | Which quorum of the blob to retrieve for (note: a blob can have multiple quorums and the chunks for different quorums at a Node can be different). The ID must be in range [0, 254]. |






<a name="node-v2-NodeInfoReply"></a>

### NodeInfoReply
Node info reply


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| semver | [string](#string) |  |  |
| arch | [string](#string) |  |  |
| os | [string](#string) |  |  |
| num_cpu | [uint32](#uint32) |  |  |
| mem_bytes | [uint64](#uint64) |  |  |






<a name="node-v2-NodeInfoRequest"></a>

### NodeInfoRequest
Node info request






<a name="node-v2-StoreChunksReply"></a>

### StoreChunksReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| signature | [bytes](#bytes) |  |  |






<a name="node-v2-StoreChunksRequest"></a>

### StoreChunksRequest
Request that the Node store a batch of chunks.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| batch | [common.v2.Batch](#common-v2-Batch) |  | batch of blobs to store |
| disperserID | [uint32](#uint32) |  | ID of the disperser that is requesting the storage of the batch. |
| signature | [bytes](#bytes) |  | Signature using the disperser&#39;s ECDSA key over keccak hash of the batch. The purpose of this signature is to prevent hooligans from tricking DA nodes into storing data that they shouldn&#39;t be storing.

Algorithm for computing the hash is as follows. All integer values are serialized in big-endian order (unsigned). A reference implementation (golang) can be found at https://github.com/Layr-Labs/eigenda/blob/master/disperser/auth/request_signing.go

1. digest batch.BatchHeader.BatchRoot 2. digest batch.BatchHeader.ReferenceBlockNumber (8 bytes, unsigned big endian) 3. for each certificate in batch.BlobCertificates: a. digest certificate.BlobHeader.Version (4 bytes, unsigned big endian) b. for each quorum_number in certificate.BlobHeader.QuorumNumbers: i. digest quorum_number (4 bytes, unsigned big endian) c. digest certificate.BlobHeader.Commitment.Commitment d. digest certificate.BlobHeader.Commitment.LengthCommitment e. digest certificate.BlobHeader.Commitment.LengthProof f. digest certificate.BlobHeader.Commitment.Length (4 bytes, unsigned big endian) g. digest certificate.BlobHeader.PaymentHeader.AccountId h. digest certificate.BlobHeader.PaymentHeader.ReservationPeriod (4 bytes, unsigned big endian) i. digest certificate.BlobHeader.PaymentHeader.CumulativePayment j. digest certificate.BlobHeader.PaymentHeader.Salt (4 bytes, unsigned big endian) k. digest certificate.BlobHeader.Signature l. for each relay in certificate.Relays: i. digest relay (4 bytes, unsigned big endian) 4. digest disperserID (4 bytes, unsigned big endian)

Note that this signature is not included in the hash for obvious reasons. |





 

 

 


<a name="node-v2-Dispersal"></a>

### Dispersal
WARNING: the following RPCs are experimental and subject to change.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| StoreChunks | [StoreChunksRequest](#node-v2-StoreChunksRequest) | [StoreChunksReply](#node-v2-StoreChunksReply) |  |
| NodeInfo | [NodeInfoRequest](#node-v2-NodeInfoRequest) | [NodeInfoReply](#node-v2-NodeInfoReply) |  |


<a name="node-v2-Retrieval"></a>

### Retrieval


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetChunks | [GetChunksRequest](#node-v2-GetChunksRequest) | [GetChunksReply](#node-v2-GetChunksReply) | GetChunks retrieves the chunks for a blob custodied at the Node. |
| NodeInfo | [NodeInfoRequest](#node-v2-NodeInfoRequest) | [NodeInfoReply](#node-v2-NodeInfoReply) | Retrieve node info metadata |

 



<a name="relay_relay-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## relay/relay.proto



<a name="relay-ChunkRequest"></a>

### ChunkRequest
A request for chunks within a specific blob. Requests are fulfilled in all-or-nothing fashion. If any of the
requested chunks are not found or are unable to be fetched, the entire request will fail.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| by_index | [ChunkRequestByIndex](#relay-ChunkRequestByIndex) |  | Request chunks by their individual indices. |
| by_range | [ChunkRequestByRange](#relay-ChunkRequestByRange) |  | Request chunks by a range of indices. |






<a name="relay-ChunkRequestByIndex"></a>

### ChunkRequestByIndex
A request for chunks within a specific blob. Each chunk is requested individually by its index.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob_key | [bytes](#bytes) |  | The blob key. |
| chunk_indices | [uint32](#uint32) | repeated | The index of the chunk within the blob. |






<a name="relay-ChunkRequestByRange"></a>

### ChunkRequestByRange
A request for chunks within a specific blob. Each chunk is requested a range of indices.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob_key | [bytes](#bytes) |  | The blob key. |
| start_index | [uint32](#uint32) |  | The first index to start fetching chunks from. |
| end_index | [uint32](#uint32) |  | One past the last index to fetch chunks from. Similar semantics to golang slices. |






<a name="relay-GetBlobReply"></a>

### GetBlobReply
The reply to a GetBlobs request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob | [bytes](#bytes) |  | The blob requested. |






<a name="relay-GetBlobRequest"></a>

### GetBlobRequest
A request to fetch one or more blobs.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob_key | [bytes](#bytes) |  | The key of the blob to fetch. |






<a name="relay-GetChunksReply"></a>

### GetChunksReply
The reply to a GetChunks request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [bytes](#bytes) | repeated | The chunks requested. The order of these chunks will be the same as the order of the requested chunks. data is the raw data of the bundle (i.e. serialized byte array of the frames) |






<a name="relay-GetChunksRequest"></a>

### GetChunksRequest
Request chunks from blobs stored by this relay.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| chunk_requests | [ChunkRequest](#relay-ChunkRequest) | repeated | The chunk requests. Chunks are returned in the same order as they are requested. |
| operator_id | [bytes](#bytes) |  | If this is an authenticated request, this should hold the ID of the operator. If this is an unauthenticated request, this field should be empty. Relays may choose to reject unauthenticated requests. |
| operator_signature | [bytes](#bytes) |  | If this is an authenticated request, this field will hold a BLS signature by the requester on the hash of this request. Relays may choose to reject unauthenticated requests.

The following describes the schema for computing the hash of this request This algorithm is implemented in golang using relay.auth.HashGetChunksRequest().

All integers are encoded as unsigned 4 byte big endian values.

Perform a keccak256 hash on the following data in the following order: 1. the operator id 2. for each chunk request: a. if the chunk request is a request by index: i. a one byte ASCII representation of the character &#34;i&#34; (aka Ox69) ii. the blob key iii. the start index iv. the end index b. if the chunk request is a request by range: i. a one byte ASCII representation of the character &#34;r&#34; (aka Ox72) ii. the blob key iii. each requested chunk index, in order |





 

 

 


<a name="relay-Relay"></a>

### Relay
Relay is a service that provides access to public relay functionality.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetBlob | [GetBlobRequest](#relay-GetBlobRequest) | [GetBlobReply](#relay-GetBlobReply) | GetBlob retrieves a blob stored by the relay. |
| GetChunks | [GetChunksRequest](#relay-GetChunksRequest) | [GetChunksReply](#relay-GetChunksReply) | GetChunks retrieves chunks from blobs stored by the relay. |

 



<a name="retriever_retriever-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## retriever/retriever.proto



<a name="retriever-BlobReply"></a>

### BlobReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [bytes](#bytes) |  | The blob retrieved and reconstructed from the EigenDA Nodes per BlobRequest. |






<a name="retriever-BlobRequest"></a>

### BlobRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| batch_header_hash | [bytes](#bytes) |  | The hash of the ReducedBatchHeader defined onchain, see: https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/interfaces/IEigenDAServiceManager.sol#L43 This identifies the batch that this blob belongs to. |
| blob_index | [uint32](#uint32) |  | Which blob in the batch this is requesting for (note: a batch is logically an ordered list of blobs). |
| reference_block_number | [uint32](#uint32) |  | The Ethereum block number at which the batch for this blob was constructed. |
| quorum_id | [uint32](#uint32) |  | Which quorum of the blob this is requesting for (note a blob can participate in multiple quorums). |





 

 

 


<a name="retriever-Retriever"></a>

### Retriever
The Retriever is a service for retrieving chunks corresponding to a blob from
the EigenDA operator nodes and reconstructing the original blob from the chunks.
This is a client-side library that the users are supposed to operationalize.

Note: Users generally have two ways to retrieve a blob from EigenDA:
  1) Retrieve from the Disperser that the user initially used for dispersal: the API
     is Disperser.RetrieveBlob() as defined in api/proto/disperser/disperser.proto
  2) Retrieve directly from the EigenDA Nodes, which is supported by this Retriever.

The Disperser.RetrieveBlob() (the 1st approach) is generally faster and cheaper as the
Disperser manages the blobs that it has processed, whereas the Retriever.RetrieveBlob()
(the 2nd approach here) removes the need to trust the Disperser, with the downside of
worse cost and performance.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| RetrieveBlob | [BlobRequest](#retriever-BlobRequest) | [BlobReply](#retriever-BlobReply) | This fans out request to EigenDA Nodes to retrieve the chunks and returns the reconstructed original blob in response. |

 



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

