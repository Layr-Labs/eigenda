# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [disperser.proto](#disperser-proto)
    - [BatchHeader](#disperser-BatchHeader)
    - [BatchMetadata](#disperser-BatchMetadata)
    - [BlobHeader](#disperser-BlobHeader)
    - [BlobInfo](#disperser-BlobInfo)
    - [BlobQuorumParam](#disperser-BlobQuorumParam)
    - [BlobStatusReply](#disperser-BlobStatusReply)
    - [BlobStatusRequest](#disperser-BlobStatusRequest)
    - [BlobVerificationProof](#disperser-BlobVerificationProof)
    - [DisperseBlobReply](#disperser-DisperseBlobReply)
    - [DisperseBlobRequest](#disperser-DisperseBlobRequest)
    - [RetrieveBlobReply](#disperser-RetrieveBlobReply)
    - [RetrieveBlobRequest](#disperser-RetrieveBlobRequest)
    - [SecurityParams](#disperser-SecurityParams)
  
    - [BlobStatus](#disperser-BlobStatus)
  
    - [Disperser](#disperser-Disperser)
  
- [Scalar Value Types](#scalar-value-types)



<a name="disperser-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## disperser.proto



<a name="disperser-BatchHeader"></a>

### BatchHeader



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| batch_root | [bytes](#bytes) |  | The root of the merkle tree with the hashes of blob headers as leaves. |
| quorum_numbers | [bytes](#bytes) |  | All quorums associated with blobs in this batch. |
| quorum_signed_percentages | [bytes](#bytes) |  | The percentage of stake that has signed for this batch. The quorum_signed_percentages[i] is percentage for the quorum_numbers[i]. |
| reference_block_number | [uint32](#uint32) |  | The Ethereum block number at which the batch was created. The Disperser will encode and disperse the blobs based on the onchain info (e.g. operator stakes) at this block number. |






<a name="disperser-BatchMetadata"></a>

### BatchMetadata



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| batch_header | [BatchHeader](#disperser-BatchHeader) |  |  |
| signatory_record_hash | [bytes](#bytes) |  | The hash of all public keys of the operators that did not sign the batch. |
| fee | [bytes](#bytes) |  | The gas fee of confirming this batch. It&#39;s the bytes representation of a big.Int value. |
| confirmation_block_number | [uint32](#uint32) |  | The Ethereum block number at which the batch is confirmed onchain. |
| batch_header_hash | [bytes](#bytes) |  | This is the hash of the ReducedBatchHeader defined onchain, see: https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/interfaces/IEigenDAServiceManager.sol#L43 The is the message that the operators will sign their signatures on. |






<a name="disperser-BlobHeader"></a>

### BlobHeader



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| commitment | [bytes](#bytes) |  | KZG commitment to the blob. |
| data_length | [uint32](#uint32) |  | The length of the blob in symbols (each symbol is 31 bytes). |
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
| adversary_threshold_percentage | [uint32](#uint32) |  | Same as SecurityParams.adversary_threshold. |
| quorum_threshold_percentage | [uint32](#uint32) |  | Same as SecurityParams.quorum_threshold. |
| quantization_param | [uint32](#uint32) |  | This determines the nominal number of chunks for the blob, which is nominal_num_chunks = quantization_param * num_operators. A chunk is the smallest unit that&#39;s distributed to DA Nodes, corresponding to a set of evaluations of the polynomial (representing the blob) and a KZG multiproof. See more details in data model of EigenDA: https://github.com/Layr-Labs/eigenda/blob/master/docs/spec/data-model.md |
| encoded_length | [uint64](#uint64) |  | The length of the blob after encoding (in number of symbols). |






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
| result | [BlobStatus](#disperser-BlobStatus) |  | The status of the blob associated with the request_id. |
| request_id | [bytes](#bytes) |  | The request ID generated by the disperser. Once a request is accepted (although not processed), a unique request ID will be generated. Two different DisperseBlobRequests (determined by the hash of the DisperseBlobRequest) will have different IDs, and the same DisperseBlobRequest sent repeatedly at different times will also have different IDs. The client should use this ID to query the processing status of the request (via the GetBlobStatus API). |






<a name="disperser-DisperseBlobRequest"></a>

### DisperseBlobRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [bytes](#bytes) |  | The data to be dispersed. The size of data must be &lt;= 512KiB. |
| security_params | [SecurityParams](#disperser-SecurityParams) | repeated | Security parameters allowing clients to customize the safety (via adversary threshold) and liveness (via quorum threshold). Clients can define one SecurityParams per quorum, and specify multiple quorums. The disperser will ensure that the encoded blobs for each quorum are all processed within the same batch. |






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






<a name="disperser-SecurityParams"></a>

### SecurityParams
SecurityParams contains the security parameters for a given quorum.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| quorum_id | [uint32](#uint32) |  | The ID of the quorum. The quorum must be already registered on EigenLayer. The ID must be in range [0, 255]. |
| adversary_threshold | [uint32](#uint32) |  | The max percentage of stake within the quorum that can be held by or delegated to adversarial operators.

Clients use this to customize the trust assumption (safety).

Requires: 1 &lt;= adversary_threshold &lt; 100 |
| quorum_threshold | [uint32](#uint32) |  | The min percentage of stake that must attest in order to consider the dispersal is successful.

Clients use this to customize liveness requirement. The higher this number, the more operators may need to be up for attesting the blob, so the chance the dispersal request to fail may be higher (liveness for dispersal).

Requires: 1 &lt;= quorum_threshld &lt;= 100 quorum_threshld &gt; adversary_threshold.

Note: The adversary_threshold and quorum_threshold will directly influence the cost of encoding for the blob to be dispersed, roughly by a factor of 100 / (quorum_threshold - adversary_threshold). See the spec for more details: https://github.com/Layr-Labs/eigenda/blob/master/docs/spec/protocol-modules/storage/overview.md |





 


<a name="disperser-BlobStatus"></a>

### BlobStatus


| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN | 0 |  |
| PROCESSING | 1 | PROCESSING means that the blob is currently being processed by the disperser |
| CONFIRMED | 2 | CONFIRMED means that the blob has been dispersed to DA Nodes and the dispersed batch containing the blob has been confirmed onchain |
| FAILED | 3 | FAILED means that the blob has failed permanently (for reasons other than insufficient signatures, which is a separate state) |
| FINALIZED | 4 | FINALIZED means that the block containing the blob&#39;s confirmation transaction has been finalized on Ethereum |
| INSUFFICIENT_SIGNATURES | 5 | INSUFFICIENT_SIGNATURES means that the quorum threshold for the blob was not met for at least one quorum. |


 

 


<a name="disperser-Disperser"></a>

### Disperser
Disperser defines the public APIs for dispersing blobs.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| DisperseBlob | [DisperseBlobRequest](#disperser-DisperseBlobRequest) | [DisperseBlobReply](#disperser-DisperseBlobReply) | This API accepts blob to disperse from clients. This executes the dispersal async, i.e. it returns once the request is accepted. The client could use GetBlobStatus() API to poll the the processing status of the blob. |
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

