
# Chunk Assignment Design

## Definitions and assumptions

Let $O$ be the set of operators. We will assume that the set of operators has a fixed size $|O| = n$. For example, we can take $n=500$ due to current restrictions on signature verification.

For each operator $i \in O$, let $S_i$ signify the amount of stake held by that operator. 


## Primary Requirements

The system must observe the following primary design requirements

### 1. Upholds security guarantee

The system specifies a security parameter $\alpha$ corresponding to the maximum amount of adversarial stake. For a given data store with a percentage of a stake signed $\beta > \alpha$, this means that we assume that a percentage $\gamma = \beta-\alpha$ of the total stake is both 1) honest and 2) in possession of their chunk. 

For the system to be secure, we need to be sure that the chunks associated with this stake suffice to reconstruct the original blob. That is, for any set of operators $U \subseteq O$ such that the following condition holds

$$ \sum_{i \in U} S_i \ge \gamma \sum_{i \in O}S_i$$ 

we must be able to reconstruct the original blob from the chunks held by the operators in $U$. 


#### Alternative Statement
An alternative security requirement: For all possible groups of dishonest operators, we need to be able to reconstruct from the complement. That is for any $U_q \subseteq O$ such that 

$$ \sum_{i \in U_q} S_i \ge \beta \sum_{i \in O}S_i$$ 

and any $U_a \subseteq U_q$ such 

$$ \sum_{i \in U_a} S_i \le \alpha \sum_{i \in O}S_i$$


we need to be able to reconstruct from $U_q \setminus U_a$. But we can see that the total stake held by this group will satisfy

$$
\sum_{i \in U_q \setminus U_a} S_i = \sum_{i \in U_q}S_i - \sum_{i \in U_a}S_i \ge (\beta-\alpha)\sum_{i \in O}S_i  = \gamma \sum_{i \in O}S_i.
$$

This tells us that the family $\{U_q \setminus U_a \}$ which satisfies the above conditions is contained in the family $\{U\}$ satisfying the previous condition. Since the second family is strictly larger, the alternative statement can actually result in a **larger** m. 

### 2. Observes operator bandwidth usage limits

Let $B$ be the total size of a blob sent to EigenDA, e.g. in MB, and let $B_i$ denote the portion of the blob stored by operator $i$ having stake $S_i$. 

Note that if blob data were distributed exactly in accordance with stake, an operator's storage requirement would be given by

$$\gamma \tilde{B}_i = B\frac{S_i}{\sum_j S_j}.$$

We require that portion of the blob stored by an operator $i$ will exceed its proportional allocation by no more than $B/n\gamma$. That is 

$$\max_{\{S_j:j\in O\}} \gamma\frac{B_i - \tilde{B}_i}{B} \le 1/n.$$

### 3. Minimizes encoding complexity

The system should minimize coding and verification computational complexity for both the disperser and operators. The computational complexity roughly scales with the number of chunks (or more specifically, inversely with the chunk size) [clarification required]. Thus, the system should minimize the number of chunks, subject to requirements 1 and 2. 

## Proposed solution

We outline the essential structure of the proposed solution, leaving optimizations for later. 

We introduce a quantization parameter $\rho$. Given $\rho$, we allocate to each operator $i$ a number of chunks equal to 

$$m_i = \text{ceil}\left(\frac{\rho nS_i}{\sum_j S_j}\right).$$

We can think of $\rho n$ as the nominal number of chunks. 

For a given $\rho$, we will be able to find a set of encoding parameters that satisfy requirement 1. Then, the strategy will be to find the smallest $\rho$ (in keeping with 3) which satisfies requirement 2. Let us first consider the problem of determining the correct coding parameters for a given $\rho$. 

### Upholding security (Req. 1) for fixed quantization parameter

<!-- In the context of the chunk allocation scheme discussed above, we can write the security requirement as follows:  -->

We know that we need to be able to reconstruct the blob from the chunks held by all sets $\mathcal{U}$ of operators $U$ satisfying the condition 

$$ \sum_{i \in U} S_i \ge \gamma \sum_{i \in O}S_i.$$

For such a set of operators $U$, the number of chunks held is given by $m_U = \sum_{i \in U}m_i$. 

Let $m$ be the smallest such $m_U$, i.e., 

$$m = \min_{U \in \mathcal{U}} m_U .$$

Clearly, if the encoding is such that we can reconstruct from any $m$ chunks, then we can also reconstruct from $m_U$ chunks for any $U \in \mathcal{U}$. Thus, $m$ determines the required reconstruction property of the system. 

<!-- Fortunately, this minimization can be solved very simply in $O(nlogn)$ complexity, and fraud proved in $O(|U^*|)$ complexity. [Details to be provided after Gautham has a chance to try figuring it out] -->


### Finding the optimal quantization parameter

Let's find a lower bound for $m$. Using the fact that $m_i \ge \rho nS_i\sum_jS_j$ and substituting, we have 

$$m = \sum_{i \in U^\star} m_i \ge  n\rho\frac{\sum_{i \in U^\star} S_i}{\sum_jS_j} \ge n\rho\gamma$$

where $U^\star \in \mathcal{U}$ is the set solving the optimization in the previous section and the final inequality uses the first inequality of the previous section, which holds for all $U \in \mathcal{U}$.

Let $C$ be the size of an individual chunk and recall that $B$ is the blob size. Since 

$$C = B/m$$

we have that

$$C \le  \frac{B}{n\rho\gamma}.$$

Moreover, the total amount of data allocated to a given operator $i$ is bounded by 

$$B_i = m_i C \le \left(\frac{n\rho S_i}{\sum_j S_j}+1\right)C \le B\left(\frac{S_i}{\gamma \sum_jS_j} + \frac{1}{n\rho\gamma}\right) = \tilde{B}_i + \frac{B}{n\rho\gamma}$$

and further 

$$\gamma\frac{B_i - \tilde{B}_i}{B} \le \frac{1}{n\rho}.$$

We can therefore satisfy requirement 2 by letting $\rho=1$.

### Assessment of coding complexity

It turns out that to meet the desired requirements, we do not need to increase the encoding complexity (i.e. decrease chunk size) compared to the default case. An increase in the total number of chunks due to the `ceil()` function can be handled by increasing the number of parity symbols. 

Moreover, the optimization routine described for finding $m$ will serve only to improve beyond the baseline (lower bound), which already achieves desired performance. 

## FAQs

Q1. Can increasing the number of parity symbols increase the total degree of the polynomial, resulting in greater coding complexity. 

A1. This seems like a possibility. In general, interactions with constraints of the proving system are not covered here. However, if this is a concern it should be possible to adjust block size constraints accordingly to avoid pushing over some limit. 
