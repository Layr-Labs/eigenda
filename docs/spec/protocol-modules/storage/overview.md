# Storage

## Overview

Within the modular structure of the EigenDA protocol, when we consider "storage," we are thinking about the following guarantee:

>> When the minimal adversarial threshold assumptions of a blob are met for any quorum, then on-chain acceptance of a blob implies a full blob is held by honest DA nodes of that quorum for the designated period.

We can further break this guarantee into two parts:
1. Acceptance guarantee: When a sufficient portion of stake is held by honest DA nodes of a given quorum, then off-chain acceptance implies that a full blob was transferred to the honest DA nodes of that quorum.
2. Storage guarantee: When honest nodes receive data that is accepted on-chain, they will store it for the designated period.

## Acceptance Guarantee

Let $\alpha$ denote the maximum proportion of adversarial stake that the system is able to tolerate (such that the actual amount will always be less than or equal to $\alpha$). Likewise, let $\beta$ represent the amount of stake held by the signing operators.

The acceptance guarantee requires that for any possible group of dishonest operators, we need to be able to reconstruct from its complement. That is for any $U_q \subseteq O$ such that

$$ \sum_{i \in U_q} S_i \ge \beta \sum_{i \in O}S_i$$

and any $U_a \subseteq U_q$ such

$$ \sum_{i \in U_a} S_i \le \alpha \sum_{i \in O}S_i$$

we need to be able to reconstruct from $U_q \setminus U_a$.

The guarantee is upheld by two smaller modules of encoding and assignment.

#### Encoding
Encoding is used to take a data blob and transform it into an extended representation consisting of a collection of chunks, such that the original blob can be reconstructed from any sufficiently large group of chunks. This must be done in a verifiable manner so that the agent performing the encoding does not need to be a trusted actor. See [Encoding](./encoding.md) for details and validation actions.

#### Assignment
The acceptance guarantee is only satisfied when chunks are properly assigned to DA nodes in proportion to the amount of stake held by the DA nodes within the required quorums. See [Assignment](./assignment.md) for details and validation actions.

## Storage Guarantee

The storage guarantee derives simply from the conditions under which a node will continue to store a blob which has been accepted. Once an EigenDA node has attested to a batch, it will store it until one of the following conditions is met:
- If there is no finalized confirmation transaction including the header of a batch within `BLOCK_STALE_MEASURE` blocks of the dataStore's `referenceBlockNumber`, the node knows that the dataStore cannot be confirmed on chain, so they prune it from their storage.
- If there is a finalized confirmation, TODO: how long will they store?
