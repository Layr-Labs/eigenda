package arbitrum_altda

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

/*
	This is a (hopefully) comprehensive handlers blue print for introducing a new ALT DA server type
	that's compatible with Arbitrum's upcoming Custom DA spec.

	TODO: Understand the best way to incorporate the ALT DA server into the
	 	  existing eigenda-proxy E2E testing framework

	TODO: Understand what fork management for our Arbitrum forks will look like; at a high level we need:
			1. test E2E correctness of the nitro stack with EigenDA.
			2. introduce missing key security checks that could impact the integration's L2 Beat assessment

	TODO: Consider that we'll need access to the Arbitrum's ALT DA client for testing, we'll need to either
	      import (likely infeasible) or vendor. vendoring might be the only option for this - we've talked about
		  doing this on other code sections across the monorepo. There could be some way to leverage a programming
		  agent to do detection for when vendor'd code should update to help smooth the maintenance burden introduced.

	TODO: Method implementations:
		[X] IsValidHeaderByte // trusted integration
		[ ] Store // trusted integration
		[ ] RecoverPayloadFromBatch // trusted integration
		[ ] GenerateProof // secure integration
		[ ] GenerateCertificateValidityProof // secure integration
*/

// Handlers defines the Arbitrum ALT DA server spec's JSON RPC methods
// This method implementations should serve as a thin wrapper over the existing EigenDA manager construct
// with translation mapping 503 (failover) and 418 (invalid_cert) status codes into error messages that
// arbitrum nitro can understand to take actions preserving both rollup liveness and safety
//
// Some custom code / refactoring will likely be necessary for supporting the READPREIMAGE proof serialization logic
type Handlers struct {
	// TODO: Metrics support - makes sense to share metrics server between both rest and arbitrum alt da
	//       servers. There should exist some label used or tag that can be used to filter between
	//       this and the REST ALT DA Server.
	// TODO: Add EigenDA manager here
	// TODO: Add logging
}

// IsValidHeaderByte determines whether or not the sequencer message header byte is an EigenDAV2 cert type.
// Arbitrum Nitro does this check via a bitwise AND which can cause overlapping and requires careful future
// management. while we could determine a byte value with bits that don't overlap - it's more maintainable
// to do a literal comparison and assume OCL NOR our competitors would never introduce a conflicting byte value
func (s *Handlers) IsValidHeaderByte(ctx context.Context, headerByte byte) (IsValidHeaderByteResult, error) {
	isValid := headerByte == EigenDAV2MessageHeaderByte

	return IsValidHeaderByteResult{
		IsValid: isValid,
	}, nil
}

// RecoverPayloadFromBatch is used to fetch the rollup payload of
// of the dispersed batch provided the DA Cert bytes.
//
// @param batch_num: batch number position in global state sequence
// @param batch_block_hash: block hash of the certL1InclusionBlock
// @param sequencer_msg: The DA Certificate
// @param preimages: Preimage mapping
// @param validateSeqMsg: Whether or not to validate the DA Cert
//
// @return bytes: Rollup payload bytes
// @return error: A structured error message (if applicable)
//
// TODO: Implement method to use EigenDA manager's Get function for retrieving rollup payload from
// DA Cert.
//
// TODO: Map 418 Im A Teapot or "Invalid Cert" status code into a "drop cert" signal
// that's processed by the Nitro Inbox Reader IF validateSeqMsg=true.
// If validateSeqMsg=false then it's assumed that the chain will halt in the presence of an invalid cert
//
// TODO: Populate preimage mapping with EigenDA V2 batch using KECCAK256 hash of DA Cert bytes as key.
// Preimages mapping is used for powering the validation pipeline when doing defensive validations or
// challenge block executions. the raw Payload returned is used for standard STF.
// It might make sense to populate the preimage mapping with the payload poly representation of the blob
// if we need to prove the decoding within the replay script.
func (s *Server) RecoverPayloadFromBatch(
	ctx context.Context,
	batchNum hexutil.Uint64,
	batchBlockHash common.Hash,
	sequencerMsg hexutil.Bytes,
	preimages PreimagesMap,
	validateSeqMsg bool,
) (RecoverPayloadFromBatchResult, error) {
	panic("RecoverPayloadFromBatch is unimplemented")
}

// Store persists a rollup payload to EigenDA and returns an associated ABI encoded DA Cert.
//
// @param message: The rollup payload bytes
//
//	@param timeout: context timeout for how long the request can be processed up-to
//	@param disableFallbackStoreDataOnChain: whether or not to enable a failover
//	               signal in the event of a detected liveness outage
//
//	@return bytes: Arbitrum Custom DA commitment bytes
//	@return error: a structured error message (if applicable)
//
// TODO: Implement method to use EigenDA manager's Put function. Given the dispersal backend
// can be dynamically modified - we should add a check to ensure that we're always dispersing to EigenDA V2
//
// TODO: Map 503 Service Unavailable status code error returned from EigenDA manager into an Arbitrum
// failover error message if disableFallbackStoreDataOnChain=true
//
// TODO: Determine the encoding standard to use for the returned DA Commitment. It's assumed that an EigenDAV2 message
// header byte will be prefixed. We can likely reuse the Standard Commitment mode but will require some analysis.
//
// TODO: Add processing for client provided timeout value
func (s *Server) Store(
	ctx context.Context,
	message hexutil.Bytes,
	timeout hexutil.Uint64,
	disableFallbackStoreDataOnChain bool,
) (StoreResult, error) {
	panic("Store is unimplemented")
}

// Generate proof is used to prove a 32 byte CustomDA preimage type for READPREIMAGE
// The exact implementation here is still a bit TBD - but we'll prove availability of the 32 bytes
// by computing a kzg point opening proof using the data commitment provided in the DA Cert.
// This will be equivalent to what's already done in the arbitrator for serializing an EigenDA READPREIMAGE
// proof. The large difference is this is done on the Custom DA server in go code as an
// "extension" of the one step proof
// construction logic.
//
// READPREIMAGE only cares about the availability or corectness of an EigenDA blob wrt it's kzg data commitment that's
// persisted in the already agreed upon DA Cert.
// Let's assumes that the EigenDA disperser would never sign over a DA Cert with an invalid data commitment.
// Pulling that off would require majority corruption of the EigenDA operator quorums and collusion with disperser
// which is a highly improbable event.
// The data commitment is a tamper resistant field in the rollup domain since modification would result
// in an incorrect merkle leaf hash being constructed from the blob header and result in an invalid merkle inclusion
// proof which would be treated as an invalid DA Cert by the rollup.
//
// TODO: Generating the data witness "opening" proof requires access to the entire EigenDA blob
// which isn't provided by client here. We can do a storage retrieval operation through the EigenDA Manager
// to fetch the blob corresponding to the DA Cert. Redundantly performing DA Cert verification is a necessary
// invariant here to strictly enforce given that this function would only ever be called if checkDACert(DA Cert)=true.
// It's slow to do another storage lookup but performance considerations are irrelevant given this is only callable
// in the worst case one step proof.
//
// TODO: Determine encoding standard that's also understood for onchain verification
//

/*
current encoding proposal:

	Assumptions:
		- kzg commitment and preimage length are extractable
		  from the existing DA Cert

	Proposed schema:
		- [0:32]: root of unity @ field element offset
		- [32:64]: field element or preimageChunk being one step proven
		- [64:128]: point opening proof (g1 point)
		- [128:256]: g2TauMinusG2z
*/
func (s *Server) GenerateProof(
	ctx context.Context,
	preimageType hexutil.Uint,
	certHash common.Hash,
	offset hexutil.Uint64,
	certificate hexutil.Bytes,
) (*GenerateProofResult, error) {
	panic("GenerateProof method is unimplemented")
}

// Non operational implementation.
// The DA Cert is already tamper resistant given its already been pre-committed to a rollup inbox
// and is verified against memory pre-state in the
// There’s no need for appending additional proof metadata for a one step proof tx
// contesting DA Cert validity
//
// TODO: Assuming we have to manage a custom fork of nitro, should we remove the proof enhancement step for
// ValidateCert opcode given the client<>server latency introduced given its noop? Then again,
// this is only ever called in the worst case one step proof WHEN the determined canonnical prestate between
// challengers is the step before calling a ValidateCert type opcode
func (s *Server) GenerateCertificateValidityProof(
	ctx context.Context,
	preimageType hexutil.Uint,
	certificate hexutil.Bytes,
) (*GenerateCertificateValidityProofResult, error) {
	panic("GenerateCertificateValidityProof is unimplemented")
}
