# EigenDA Throughput Characteristics

## What is the total throughput of EigenDA? 

Let $S_i$ denote the total stake delegated to an operator $i \in O$, and $n = |O|$ denote the total number of operators. Let $C_i$ be the throughput/bandwidth capacity by the individual operator $i$. 

Recall that the system specifies a security parameter $\alpha$ corresponding to the maximum amount of adversarial stake tolerated. For a given data store with a percentage of a stake signed $\beta > \alpha$, this means that we assume that a percentage $\gamma = \beta-\alpha$ of the total stake is both 1) honest and 2) in possession of their chunk.

This writeup will assume the use of the (unoptimized) basic assignment scheme specified in the [Assignment Module](./protocol-modules/storage/assignment.md#standard-assignment-security-logic), with the quantization parameter $\rho$. 

This means that for every blob, an operator $i$ is assigned a number of chunks given by 

$$m_i = \text{ceil}\left(\frac{n\rho S_i}{\sum_j S_j}\right) = \text{ceil}\left(n\rho\alpha_i\right),$$

where we define $\alpha_i$ to be the fraction of total stake held by $i$. 

Moreover, blobs will be encoded with the reconstruction threshold $m$ (number of chunks needed to reconstruct blob, i.e. number of systematic chunks) satisfying

$$m \ge n\rho\gamma$$

For a given unit of throughput $T$ accepted by the system, a fraction $\eta_i = \frac{T_i}{T} = \frac{m_i}{m}$ must be accepted by operator $i$. Since 

$$T_i = T\frac{m_i}{m} \le C_i$$

must hold for all $i \in O$, rearranging we see that 

$$T = \min_{i\in O} C_i \frac{m}{m_i} = \min_{i\in O} \frac{n\rho C_i}{\text{ceil}(n\rho \alpha_i)}\gamma \approx \min_{i\in O} \frac{C_i}{\max(1/n\rho,\alpha_i)}\gamma.\tag{1}$$


## Discussion

Notice that Eq. (1) implies two dynamics by which the total throughput of EigenDA can increase. With $j \in O$ as operator for which $C_i/\max(1/n\rho,\alpha_i)$ is currently the smallest:

1. Operator $j$ can increase its capacity. 
2. If $\alpha_j > 1/n\rho$, then $\alpha_j$ can decrease either by stakers redelegating from $j$ to some operator with more capacity per stake or by new stakers entering the ecosystem and delegating to such an operator. 

In the next section, we will consider whether we can use these observations to construct a mechanism which incentivizes operators to marginally increase their capacity (when the demand exists) in order to gain new stake/rewards or in order to avoid losing their delegated stake. 

Note that the quantization parameter $\rho$ plays a critical role in capping the system throughput. In particular, it is likely that the throughput will be given by 

$$
T = n\rho\gamma C_j
$$

where $j$ is the smallest operator, in which case reallocations of stake as in 2) will no longer suffice to scale system throughput. Thus, it makes sense to explore approaches for ensuring that $\rho$ can be scaled to prevent smaller operators from bottlenecking the system throughput. For the purposes of this writeup, we will assume that quantization is not a limiting factor. 
