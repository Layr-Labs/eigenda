
# The "Modules" of Data Availability

The overall security guarantee provided by EigenDA is actually a composite of many smaller guarantees, each with its own intricacies. For the purpose of exposition, in this section we break down the EigenDA protocol into a set of three modules, which divide roughly along the lines of the guarantees for which the modules are responsible. 

## Attestation
The main guarantee supported by the attestation module concerns the on-chain conditions under which a batch is able to be confirmed by the EigenDA smart contracts. In particular, the attestation module is responsible for upholding the following guarantee:
- Sufficient stake checking: A blob is only accepted on-chain when signatures from operators having sufficient stake on each quorum are presented. 

The Attestation module is largely implemented by the EigenDA smart contracts via bookkeeping of stake and associated checks performed at the batch confirmation phase of the [Disperal Flow](../flows/dispersal.md). For more details, see the [Attestation module documentation](./attestation/attestation.md)

## Storage
The main guarantee supported by the storage module concerns the off-chain conditions which mirror the on-chain conditions of the storage module. In particular, the storage module is responsible for upholding the following guarantee:
- Honest custody of complete blob: When the minimal adversarial threshold assumptions of a blob are met for any quorum, then on-chain acceptance of a blob implies a full blob is held by honest DA nodes of that quorum for the designated period.

The Storage module is largely implemented by the DA nodes, with an untrusted supporting role by the Disperser. For more details, see the 
[Storage module documentation](./storage/overview.md)

## Retrieval 
The main guarantee supported by the retrieval module concerns the retrievability of stored blob data by honest consumers of that data. In particular, the retrieval module is responsible for upholding the following guarantee:
- TODO: Articulate the retrieval guarantee that we support. 

For more details, see the [Retrieval module documentation](./retrieval/retrieval.md)
