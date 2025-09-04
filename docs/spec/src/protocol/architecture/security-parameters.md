# Security Parameters

In this page, we prove the relationship between the blob parameters and security thresholds. 
We also point the reader to the code where the constraint of the security thresholds are implemented.

## Encoding Rate and Reconstruction Threshold
Recall the blob parameters in [Encoding](./encoding.md):
- $c$: The total number of encoded chunks.  
- $\gamma$: The ratio of original data to total encoded chunks, providing redundancy.
- $r$: The minimum fraction of total stake required to reconstruct the blob.  
- $n$: Maximum number of validators.

In this section, we prove that, with our [assignment algorithm](./assignment.md), the encoding rate and the reconstruction threshold satisfy the following inequality:
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
**Safety Threshold** 

**Liveness Threshold**

### Relationship Between Security Thresholds, Confirmation Threshold and Reconstruction Threshold

