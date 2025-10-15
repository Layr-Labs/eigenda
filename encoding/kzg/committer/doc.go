// Package committer provides functions to create and verify EigenDA [encoding.BlobCommitments].
//
// Note that EigenDA blob commitments are not simply a single KZG commitment, but also
// include the blob's length, as well as a proof of this length (LengthCommitment + LengthProof).
// This complexity stems from the fact that EigenDA, unlike Ethereum which only allows 128KiB blobs,
// allows blobs of any power-of-2 size between 32B and 16MiB (currently).
//
// There are 2 facets to data availability:
// 1. Local (chunks) availability: validator attests to having received and being able to serve its chunks
// 2. Global (blob) availability: validator attests to the entire blob being available in the network.
//
// Because of the sharded nature of EigenDA, each validator only receives a subset of the blob's content.
// In order to attest to global availability, it thus needs to know how many chunks there are in total,
// and to make sure that the chunks it receives are actually proportional to its stake. This is why
// BlobCommitments contains a length field, as well as a proof of this length (LengthCommitment + LengthProof).
//
// Here's an example scenario which shows that EigenDA could go wrong without this Length.
// In the extreme case, a malicious disperser could just tell the validators that the blob size is 1,
// and ask all validators except for one to sign off on the commitment. For a slightly more involved but
// analogous scenario, assume a network of 8 DA nodes with uniform stake distribution, and coding ratio 1/2.
// For a blob containing 128 field elements (FEs), each node gets 128*2/8=32 FEs, meaning that any 4 nodes can
// join forces and reconstruct the data. Now assume a world without length proof; a malicious disperser colluding
// with a client disperses the same blob/commitment, but claims that the blob only has length of 4 FEs.
// He sends each node 4*2/8=1 FE. The chunks submitted to the nodes match the commitment, so the nodes accept
// and sign over the blob’s batch. But now there are only 8 FEs in the system, which is not enough to reconstruct
// the original blob (need at least 128 for that).
//
// ----- Length Commitment + Length Proof Explanation -----
//
// In theory, proving an upper bound on the actual blob length is very simple (assuming knowledge of pairings),
// and would require only a LengthCommitment (no LengthProof needed).
// - G1 and G2: generators of the bn254 curve groups
// - BL: blob length (power of 2)
// - BC_G1: blob commitment; [p(x)]_1 := p(s)G1 (this is the same as our [encoding.BlobCommitments].Commitment)
// - LC_G1: len commitment; q(x) = x^(2^28-BL)*p(x)
// Verification is simply e(BC_G1, s^(2^28-BL)*G2) = e(LC_G1, G2)
//
// Unfortunately, this simple strategy does not work, due to our (unfortunate) choice of SRS ceremony,
// which generated 2^29 G1 points but only 2^28 G2 points. Note that this is somehow not documented in
// https://github.com/privacy-ethereum/perpetualpowersoftau/tree/master itself for some unknown reason...
// but one can see that there are twice as many points in g1.point than in g2.point from the parsing code, e.g.
// https://github.com/iden3/snarkjs/blob/e0c7219bd69db078/src/powersoftau_challenge_contribute.js#L22
// Because of these extra available G1 points, a malicious client/disperser is able to claim that its blob
// is smaller than it really is, and it can generate a LC_G1 commitment for that smaller blob length,
// given the extra available G1 SRS points.
//
// Attack in practice:
// - BL: actual blob length, same as above
// - BC_G1: same as above
// - FBL: fake blob length = BL/2
// - FLC_G1: fake length commitment to q'(x) = x^(2^28-BL/2)*p(x)
// Note that if there were only 2^28 G1 points, then the malicious client/disperser would not be able to generate
// the commitment FLC_G1, because it has degree 2^28-BL/2+BL = 2^28+BL/2 > 2^28
// - Verification works: e(BC_G1, s^{2^28-BL/2}G2) = e(FLC_G1, G2)
//
// So our actual implementation is as follows:
// - C1 (G1): blob commitment
// - C2 (G2): len commitment to p(x)
// - P2 (G2): len proof to q(x) = x^{2^28-bloblen}p(x)
// - Verify e(s^{2^28-bloblen}, C2) = e(G1, P2)
// Note there is no C1 in above pairing, which is why we verify a second pairing e(C1,G2) = e(G1,C2)
// in [VerifyCommitEquivalenceBatch]!
//
// Note that we actually missed a simpler scheme when initially implementing this,
// whose proofs are two (smaller) G1 points instead of two G2 points:
// - shift = 2^28 - bloblen
// - proof1 = [s^(shift/2) * p(s)]_1
// - proof2 = [s^shift * p(s)]_1
// - verifier pairing1: e([p(s)]_1, [s^(shift/2)]_2) = e(proof1, [1]_2)
// - verifier pairing2: e(proof1, [s^(shift/2)]_2) = e(proof2, [1]_2)
// Note that we can even optimize to a single pairing by combining the two equations
// with gamma = random fiat shamir:
// e([p(s)]_1 + gamma + proof1, [s^(shift/2)]_2) = e(proof1 + gamma * proof2, [1]_2)
package committer

import (
	_ "github.com/Layr-Labs/eigenda/encoding"
)
