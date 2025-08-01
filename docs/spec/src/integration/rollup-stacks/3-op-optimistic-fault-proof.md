# OP Optimistic Fault Proof

This document explains how to integrate **EigenDA** blob derivation (via **Hokulea**) into the OP derivation pipeline and secure it with the default OP Fault‑Proof VM (FPVM).

Upgrade 16’s [Interop Contracts proposal](https://gov.optimism.io/t/upgrade-16-proposal-interop-contracts-stage-1-and-go-1-23-support-in-cannon/10037) adds **Kona** fault‑proof programs to **Cannon**, enabling MIPS‑ELF binaries compiled from Kona. We therefore extend Kona with Hokulea so EigenDA‑based rollups can rely on the official OP fault‑proof system. *Spec is still work‑in‑progress.*

---

## OP Fault‑Proof Recap

1. Any party may dispute an L2 output by running **op‑challenger**.  
2. Players alternate moves within fixed clock deadlines; If the clock expires without a move, the last mover wins.
3. A bounded game depth reduces the dispute to one VM step, which **Cannon** re‑executes—so every step must be fault‑provable.

---

## L2 Consensus with EigenDA

| Component | Purpose | Executed in |
|-----------|---------|-------------|
| **Kona**  | OP derivation pipeline | Cannon |
| **Hokulea** | EigenDA blob derivation | Cannon |

Both parts compile into a single MIPS‑ELF. Cannon runs it whenever a challenge is raised.

---

## Proving one Step in EigenDA Blob Derivation on L1

| Type of VM Step | Verification type | Handling |
|--------------------|--------------------|----------|
| Execution step     | Logic              | The MIPS instructions are implemented in the smart contract to re-execute any incorrect logic (e.g., any incorrect execution when converting an encoded payload to a rollup payload)|
| Pre‑image lookup   | Data               | Requires correct key–value pair on L1 |

If the disputed step is a pre‑image lookup, to resolve the final setp at the maximum game depth, the player must first upload the valid key-value pair. The Preimage oracle ignores the uploaded value(preimage) if the key-value is not hold. If a party cannot provide a valid preimage on time, the party loses after its clock expire.  

---

## Onchain Pre‑Image Infrastructure

Cannon relies on [`PreimageOracle.sol`](https://github.com/ethereum-optimism/optimism/blob/develop/packages/contracts-bedrock/src/cannon/PreimageOracle.sol):

```solidity
mapping(bytes32 => uint256) public preimageLengths;
mapping(bytes32 => mapping(uint256 => bytes32)) public preimageParts;
mapping(bytes32 => mapping(uint256 => bool)) public preimagePartOk;
```

EigenDA derivation requires three pre‑images:

1. **Recency window**  
2. **Certificate validity**  
3. **Point opening on blob**

Keys are `keccak256(address)` of the reserved addresses (prefixed as *type 3* per the [OP spec](https://specs.optimism.io/fault-proof/index.html#type-3-global-generic-key)).

The pre‑image and relation between (key-value) pair can be specified by an upgradeable contract that:

- stores the recency‑window parameter as the preimage;  
- uses **certVerifier Router** to establish the validity of the DA certificate, the preimage is a boolean;  
- verifies KZG point openings on a blob (using EigenDA’s `BN254` library), the preimage is bytes from the EigenDA blob;

---

## OP‑Challenger Duties

- **Logic steps:** automatically prepares proof data for re-executing the MIPS instruction
- **Pre‑image steps:** downloads the blob from EigenDA, constructs a point‑opening proof, and submits the pre‑image to L1.

This integration ensures every logic or data step is fault‑provable, allowing EigenDA rollups to benefit from the official OP security model.
