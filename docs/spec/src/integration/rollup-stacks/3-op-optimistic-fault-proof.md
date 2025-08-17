# OP Optimistic Fault Proof with Cannon

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

## Proving one Instruction in EigenDA Blob Derivation on L1

| Type of VM Instruction | Verification type | Handling |
|--------------------|--------------------|----------|
| Execution step     | Logic              | The MIPS instructions are implemented in the smart contract to re-execute any processing logic (e.g., any incorrect execution when converting an encoded payload to a rollup payload)|
| Preimage lookup   | Data               | Requires correct key–value pair on L1 Preimage Oracle contract|

When the disputed instruction is a preimage lookup, the player must first submit the correct key-value pair to the preimage oracle contract, and then resolve the final instruction. The Preimage Oracle will disregard any submitted value if the required key-value pair relation does not hold. If a party fails to provide a valid preimage before its timer expires, that party forfeits the game.

---

## Onchain Pre‑Image Infrastructure

Cannon relies on [`PreimageOracle.sol`](https://github.com/ethereum-optimism/optimism/blob/develop/packages/contracts-bedrock/src/cannon/PreimageOracle.sol):

```solidity
mapping(bytes32 => uint256) public preimageLengths;
mapping(bytes32 => mapping(uint256 => bytes32)) public preimageParts;
mapping(bytes32 => mapping(uint256 => bool)) public preimagePartOk;
```

EigenDA blob derivation requires three pre‑images:

1. **Recency window**  
2. **Certificate validity**  
3. **Point opening on blob**

Keys of the preimages are `keccak256(address)` of the [reserved addresses](https://github.com/Layr-Labs/hokulea/tree/master/docs) (prefixed as *type 3* per the [OP spec](https://specs.optimism.io/fault-proof/index.html#type-3-global-generic-key)).

The preimage and relation between (key-value) pair can be specified by a contract that:

- stores the recency‑window parameter as the preimage;  
- uses **certVerifier Router** to establish the validity of the DA certificate, the preimage is a boolean;  
- verifies KZG point openings on a blob (using EigenDA’s `BN254` library), the preimage is 32 bytes from the EigenDA blob;

---

## OP‑Challenger Duties

- **Logic steps:** automatically prepares proof data for re-executing the MIPS instruction
- **Pre‑image steps:** downloads the blob from EigenDA, constructs a point‑opening proof, and submits the pre‑image to L1.

This integration ensures every logic or data step is fault‑provable, allowing EigenDA rollups to benefit from the official OP security model.
