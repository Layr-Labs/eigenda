# Security Parameters

In this page, we prove the relationship between the blob parameters and security thresholds. 
We also point the reader to the code where the constraint of the security thresholds are implemented.

## Encoding Rate and Reconstruction Threshold
Recall the blob parameters in [Encoding](./encoding.md):
- $c$: The total number of encoded chunks.  
- $\gamma$: The ratio of original data to total encoded chunks, providing redundancy.
- $r$ (`ReconstructionThreshold`): The minimum fraction of total stake required to reconstruct the blob.  
- $n$: Maximum number of validators.

In this section, we prove that, with our [assignment algorithm](./assignment.md), the encoding rate and the reconstruction threshold satisfy the following euqation:
$$
r = \frac{c}{c-n} \gamma 
$$

In other words, we want to prove that any subset of validators with $\frac{c}{c-n} \gamma$ of toal stake own enough chunks to reconstruct the original blob. 
Formally, we need to show that for any set of validators $H$ with total stake $\sum_{i \in H} \eta_i \geq \frac{c}{c-n} \gamma$, the chunks assigned to $H$ satisfy $\sum_{i \in H} c_i \geq c\gamma$. 

**Proof:**

By the chunk assignment scheme, we have
$$c_i \geq c'_i = \lceil \eta_i(c - n) \rceil $$
$$\geq \eta_i(c - n)$$

Therefore, since $\sum_{i \in H} \eta_i \geq \frac{c}{c-n} \gamma$, we have
$$ \sum_{i \in H} c_i \geq \sum_{i \in H} \eta_i (c-n) \geq \frac{c}{c-n} \gamma \cdot (c - n) = \gamma c$$


## BFT Security 
### Definition of Security Thresholds
In this section, we define and prove the safety and liveness properties of EigenDA, building on the reconstruction property established above.

The Byzantine liveness and safety properties of a blob are specified by a collection of `SecurityThresholds`.

- `ConfirmationThreshold` - The confirmation threshold defines the minimum percentage of stake which needs to sign to make the DA certificate valid.
- `SafetyThreshold` - The safety threshold refers to the minimum percentage of total stake an attacker must control to make a blob with a valid DA certificate unavailable.
- `LivenessThreshold` - The liveness threshold refers to the minimum percentage of total stake an attacker must control to cause a liveness failure.

### How to Set the Confirmation Threshold

In the BFT security model, the `SafetyThreshold` and `LivenessThreshold` are estimated by the client, where `SafetyThreshold` is the maximum stake controlled by an adversary that sign the certificate and refuse to serve the data and `SafetyThreshold` is the maximum stake controlled by an adversary that does not sign a certificate.

The `ConfirmationThreshold` is set based on the following two criterias.

**1. Confirmation Threshold and Safety Threshold**

To ensure that each blob with a valid DA certificate is available, the following ineqality must be satisfied when setting the `ConfirmationThreshold`: 

`ConfirmationThreshold` - `SafetyThreshold` >= `ReconstructionThreshold` (1)

Intuitively, since the adversary control less than `SafetyThreshold` of stake, at least  `ConfirmationThreshold` - `SafetyThreshold` honest validators need to sign to form a valid DA certificate. 
Therefore, as long as `ConfirmationThreshold` - `SafetyThreshold` >= `ReconstructionThreshold`, the honest validators should own a large enough set of chunks to reconstruct the blob.

**2. Confirmation Threshold and Liveness Threshold**

The `ConfirmationThreshold` and `LivenessThreshold` satisfies the following inequality:

`ConfirmationThreshold` <= 1 - `LivenessThreshold` (2)

This is because a valid certificate require signature from at least `ConfirmationThreshold` of stake. If `ConfirmationThreshold` is greater than 1 - `LivenessThreshold`, the adversary can cause a liveness failure by simply not signing the certificate.

### Implementation

**1. Safety Threshold Check**

The check for the inequality (1) above is implemented [here](https://github.com/Layr-Labs/eigenda/blob/6cd192ecbe5f0abfe73fc08df306cf00e32ef010/contracts/src/integrations/cert/libraries/EigenDACertVerificationLib.sol#L188).
Specifically, the code implements the check for the following inequality:

`ConfirmationThreshold` - `SafetyThreshold` >=  `ReconstructionThreshold`$ = \frac{c}{c-n} \gamma$,

with the following mapping of notation in the doc and variables in the code:

- `ConfirmationThreshold` : `securityThresholds.confirmationThreshold`
- `SafetyThreshold` : `securityThresholds.adversaryThreshold`
- $c$ : `blobParams.numChunks`
- $n$ : `blobParams.maxNumOperators`
- $\gamma$: 1 / `blobParams.codingRate`

**2. Liveness Threshold**

The `LivenessThreshold` does not show in the code, but it's something that the user should keep in mind if they want to change the default setting. We strongly recommend the user to set a `SafetyThreshold` >= 33%. 

By default, the `ConfirmationThreshold` is 55%. With the default `ReconstructionThreshold` = 13%, it gives a `SafetyThreshold` of 42% and a `LivenessThreshold` of 45%. 