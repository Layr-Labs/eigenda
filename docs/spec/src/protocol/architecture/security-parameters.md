# Security Parameters

This page proves the relationship between blob parameters and security thresholds. 
We also point readers to the code where security threshold constraints are implemented.

## Blob Parameters and Reconstruction Threshold

In this part, we present the blob parameters and use these parameters to derive the reconstricution threshold.

### Blob Parameters

We define the **Blob parameters** as a tuple **$(n, c, \gamma)$** where:


- $n$ (`MaxNumOperators`): Maximum number of validators allowed in EigenDA.  
- $c$ (`NumChunks`): The total number of encoded chunks after erasure coding (must be a power of 2).  
- $\gamma$ (`1/CodingRate`): The ratio of original data to total encoded chunks, providing redundancy (must be an inverse power of 2). Note that for representational purposes, the `CodingRate` in our code is the inverse of  $\gamma$, the standard coding rate used in coding theory.

Among the blob parameters, `CodingRate` and `NumChunks` are used in the [encoding](./encoding.md) process, while `NumChunks` and `MaxNumOperators` are used in the chunk [assignment](./assignment.md) process.

This tumple is stored in the struct shown below ([see in the code](https://github.com/Layr-Labs/eigenda/blob/d8090af76ed69920983bb3781399a91d84d20d10/contracts/src/core/libraries/v1/EigenDATypesV1.sol#L7)):

```solidity
struct VersionedBlobParams {
    uint32 maxNumOperators;
    uint32 numChunks;
    uint8 codingRate;
}
```
The blob parameters for each version is stored in deploying the `EigenDAThresholdRegistry` contract.
It's configured [here](https://github.com/Layr-Labs/eigenda/blob/556dc34fcd4774b683cbc78590bccee66a096b42/contracts/script/deploy/eigenda/mainnet.beta.config.toml#L69) and the default parameters are shown below.
```
versionedBlobParams = [
    { 0_maxNumOperators = 3537, 1_numChunks = 8192, 2_codingRate = 8 }
]
```

### Reconstruction Threshold

We define `ReconstructionThreshold`, also denoted as $r$, the minimum fraction of total stake required to reconstruct the blob. 
In this section, we prove that, with our [chunk assignment algorithm](./assignment.md), the reconstruction threshold is:
$$
r = \frac{c}{c-n} \gamma 
$$

In other words, we want to prove that any subset of validators with $\frac{c}{c-n} \gamma$ of total stake collectively own enough chunks to reconstruct the original blob. 
Formally, we need to show that for any set of validators $H$ with total stake $\sum_{i \in H} \eta_i \geq \frac{c}{c-n} \gamma$, the chunks assigned to $H$ satisfy $\sum_{i \in H} c_i \geq \gamma c$. 

**Proof:**

By the chunk assignment scheme, we have:
$$c_i \geq c'_i = \lceil \eta_i(c - n) \rceil $$
$$\geq \eta_i(c - n)$$

Therefore, since $\sum_{i \in H} \eta_i \geq \frac{c}{c-n} \gamma$, we have:
$$ \sum_{i \in H} c_i \geq \sum_{i \in H} \eta_i (c-n) \geq \frac{c}{c-n} \gamma \cdot (c - n) = \gamma c$$

Now, we prove that any subset of validators with $r$ of the total stake own at least $\gamma c$ chunks, which is guaranteed to reconstruct the origianl blob due to the property of Reed-Solomon encoding.

As we show in the previous subsection, by default, $n = 3537$, $c = 8192$ and $\gamma = 1/8$, which gives us the reconstruction threshold $r = 22\%$.

## BFT Security

Having established the relationship between the blob parameters and the reconstruction threshold, we now turn to the Byzantine Fault Tolerant (BFT) security model and how it relates to the blob parameters. 
### Definition of Security Thresholds
In this section, we define and prove the safety and liveness properties of EigenDA, building on the reconstruction property established above.

The Byzantine liveness and safety properties of a blob are specified by a collection of `SecurityThresholds`:

- `ConfirmationThreshold` - The confirmation threshold defines the minimum percentage of stake that must sign to make the DA certificate valid.
- `SafetyThreshold` - The safety threshold refers to the minimum percentage of total stake an attacker must control to make a blob with a valid DA certificate unavailable.
- `LivenessThreshold` - The liveness threshold refers to the minimum percentage of total stake an attacker must control to cause a liveness failure.

### How to Set the Confirmation Threshold

In the BFT security model, the `SafetyThreshold` and `LivenessThreshold` are estimated by the client. The `SafetyThreshold` is the maximum stake controlled by an adversary that signs the certificate but fails to serve the data, while the `LivenessThreshold` is the maximum stake controlled by an adversary that does not sign the certificates.

The `ConfirmationThreshold` is set based on the following two criteria:

**1. Confirmation Threshold and Safety Threshold**

To ensure that each blob with a valid DA certificate is available, the following inequality must be satisfied when setting the `ConfirmationThreshold`: 

`ConfirmationThreshold` - `SafetyThreshold` >= `ReconstructionThreshold` (1)

Intuitively, since the adversary controls less than `SafetyThreshold` of stake, at least `ConfirmationThreshold` - `SafetyThreshold` honest validators need to sign to form a valid DA certificate. 
Therefore, as long as `ConfirmationThreshold` - `SafetyThreshold` >= `ReconstructionThreshold`, the honest validators should own a large enough set of chunks to reconstruct the blob.

**2. Confirmation Threshold and Liveness Threshold**

The `ConfirmationThreshold` and `LivenessThreshold` satisfy the following inequality:

`ConfirmationThreshold` <= 1 - `LivenessThreshold` (2)

This is because a valid certificate requires signatures from at least `ConfirmationThreshold` of stake. If `ConfirmationThreshold` is greater than 1 - `LivenessThreshold`, the adversary can cause a liveness failure by simply not signing the certificate.

In summary, the `SafetyThreshold` and `LivenessThreshold` depends on the choice of `ConfirmationThreshold`. The picture below shows the relationship between these security thresholds.

![image](../../assets/security_thresholds.png)

### Implementation Details

**1. Safety Threshold**

The check for the inequality (1) above is implemented [here](https://github.com/Layr-Labs/eigenda/blob/6cd192ecbe5f0abfe73fc08df306cf00e32ef010/contracts/src/integrations/cert/libraries/EigenDACertVerificationLib.sol#L188).
Specifically, the code implements the check for the following inequality:

`ConfirmationThreshold` - `SafetyThreshold` >=  `ReconstructionThreshold`$ = \frac{c}{c-n} \gamma$,

with the following mapping of notation in the doc and variables in the code:

- `ConfirmationThreshold` : `securityThresholds.confirmationThreshold`
- `SafetyThreshold` : `securityThresholds.adversaryThreshold`
- $c$ : `blobParams.numChunks`
- $n$ : `blobParams.maxNumOperators`
- $\gamma$: 1 / `blobParams.codingRate`

We strongly recommend that users set a `SafetyThreshold` >= 33% if they ever want to change the default settings.

**2. Liveness Threshold**

The `LivenessThreshold` does not appear in the code, but users should keep it in mind when changing the default settings. 

**System Default**

By default, the `ConfirmationThreshold` is 55%. With the default `ReconstructionThreshold` = 22%, the default  `ConfirmationThreshold` gives a `SafetyThreshold` of 33% and a `LivenessThreshold` of 45%. 