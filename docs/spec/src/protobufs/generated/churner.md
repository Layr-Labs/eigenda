# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [churner/churner.proto](#churner_churner-proto)
    - [ChurnReply](#churner-ChurnReply)
    - [ChurnRequest](#churner-ChurnRequest)
    - [OperatorToChurn](#churner-OperatorToChurn)
    - [SignatureWithSaltAndExpiry](#churner-SignatureWithSaltAndExpiry)
  
    - [Churner](#churner-Churner)
  
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

