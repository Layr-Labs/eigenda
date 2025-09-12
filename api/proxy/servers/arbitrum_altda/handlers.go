package arbitrum_altda

import (
	"context"
	"fmt"

	proxy_common "github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/certs"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/commitments"
	"github.com/Layr-Labs/eigenda/api/proxy/store"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

/*
	This is a (hopefully) comprehensive handlers blue print for introducing a new ALT DA server type
	that's compatible with Arbitrum's upcoming Custom DA spec.

	TODO: Understand the best way to incorporate the ALT DA server into the
	 	  existing eigenda-proxy E2E testing framework. Vendoring is rather challenging given
		  the amount of arbitrum specific dependencies in their DA Client type. The DA Client is a wrapper
		  on-top of the go-ethereum rpc client. It could be sufficient for E2E testing intra-monorepo to
		  to utilize the generic rpc client while maintaining a more comprehensive framework in our layr-labs/nitro
		  fork.

	TODO: Understand what fork management for our Arbitrum forks will look like; at a high level we need to:
			1. test E2E correctness of the nitro stack with EigenDA
			2. introduce missing key security checks that could impact the integration's L2 Beat assessment

	TODO: Method implementations:
		[X] IsValidHeaderByte // trusted integration
		[-] Store // trusted integration
		[-] RecoverPayloadFromBatch // trusted integration
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
	//       this and the REST ALT DA Server. op-geth has added interception to provide arbitrary
	//       preprocessing callbacks on the incoming/outgoing RPC message:
	//       https://github.com/ethereum-optimism/optimism/blob/
	//       8749b77f4d6b4767e40d11371ac3d37cb7f2f2d8/op-service/metrics/rpc_metrics.go
	//
	//      This is something we could leverage but would further solidify our reliance on op-geth which
	//      would be a major footgun for long-term monorepo mgmt. Therefore manually adding metric expressions
	//      to each method function is the only viable solution - although having general modularity through
	//      callback injection would be nice :/
	//
	// TODO: Logging - the underlying go-ethereum (geth) RPC server framework uses geth logging for capturing
	//       invalid namespace/method and deserialization errors when targeting through meta-level reflection.
	///      This can result in consistency issues since this is a geth native logger where we use a custom logger
	//       maintained in https://github.com/Layr-Labs/eigensdk-go/tree/dev/logging.
	//
	//       We should dig into this underlying logging and see if there's a way to intuitively override, disable,
	//       or enforce consistency between log outputs.

	eigenDAManager *store.EigenDAManager
}

func NewHandlers(m *store.EigenDAManager) *Handlers {
	return &Handlers{
		eigenDAManager: m,
	}
}

// IsValidHeaderByte determines whether or not the sequencer message header byte is an EigenDAV2 cert type.
// Arbitrum Nitro does this check via a bitwise AND which can cause overlapping and requires careful future
// management. while we could determine a byte value with bits that don't overlap - it's more maintainable
// to do a literal comparison and assume OCL NOR our competitors would never introduce a conflicting byte value
func (h *Handlers) IsValidHeaderByte(ctx context.Context, headerByte byte) (IsValidHeaderByteResult, error) {
	return IsValidHeaderByteResult{
		IsValid: headerByte == EigenDAV2MessageHeaderByte,
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
// TODO: Map 418 Im A Teapot or "Invalid Cert" status code into a "drop cert" signal
// that's processed by the Nitro Inbox Reader IF validateSeqMsg=true.
// If validateSeqMsg=false then it's assumed that the chain will halt in the presence of an invalid cert
//
// TODO: Populate preimage mapping with EigenDA V2 batch using KECCAK256 hash of DA Cert bytes as key.
// Preimages mapping is used for powering the validation pipeline when doing defensive validations or
// challenge block executions. the raw Payload returned is used for standard STF.
// It might make sense to populate the preimage mapping with the payload poly representation of the blob
// if we need to prove the decoding within the replay script.
func (h *Handlers) RecoverPayloadFromBatch(
	ctx context.Context,
	batchNum hexutil.Uint64,
	batchBlockHash common.Hash,
	sequencerMsg hexutil.Bytes,
	preimages PreimagesMap,
	validateSeqMsg bool,
) (RecoverPayloadFromBatchResult, error) {
	if len(sequencerMsg) <= 1 {
		return RecoverPayloadFromBatchResult{},
			fmt.Errorf("sequencer message expected to be >1 byte, got: %d", len(sequencerMsg))
	}

	// strip version byte from sequencer message
	//
	// TODO: There will be additional bytes encoded here (i.e, CustomDA byte, EigenDAV2 Message Header byte).
	//       Given the lack of response from OCL, it's still unknown what the exact expected schema should.
	//       Once their interface is more hardended we can make a safer deduction for how to introduce this.
	//       Will likely require introducing a new ArbCustomDA commitment type
	//
	versionByte := sequencerMsg[0]
	versionedCert := certs.NewVersionedCert([]byte(sequencerMsg[1:]), certs.VersionByte(versionByte))

	payload, err := h.eigenDAManager.Get(ctx, versionedCert, proxy_common.GETOpts{})
	if err != nil {
		return RecoverPayloadFromBatchResult{}, fmt.Errorf("get rollup payload from DA Cert: %w", err)
	}

	return RecoverPayloadFromBatchResult{
		Payload: payload,
	}, nil
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
// TODO: Map 503 Service Unavailable status code error returned from EigenDA manager into an Arbitrum
// failover error message if disableFallbackStoreDataOnChain=true
//
// TODO: Determine the encoding standard to use for the returned DA Commitment. It's assumed that an EigenDAV2 message
// header byte will be prefixed. We can likely reuse the Standard Commitment mode but will require some analysis.
//
// TODO: Add processing for client provided timeout value
func (h *Handlers) Store(
	ctx context.Context,
	message hexutil.Bytes,
	timeout hexutil.Uint64,
	disableFallbackStoreDataOnChain bool,
) (*StoreResult, error) {
	dispersalBackend := h.eigenDAManager.GetDispersalBackend()
	if dispersalBackend != proxy_common.V2EigenDABackend {
		return nil, fmt.Errorf("expected EigenDAV2 backend, got: %v", dispersalBackend)
	}

	if len(message) == 0 {
		return nil, fmt.Errorf("received empty rollup payload")
	}

	certBytes, err := h.eigenDAManager.Put(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("put rollup payload: %w", err)
	}

	// TODO: This should eventually be propagated by the Put method given the actual
	//       version byte assumed is dictated by the EigenDACertVerifier used
	versionedCert := certs.NewVersionedCert(certBytes, certs.V2VersionByte)
	daCommitment := commitments.NewStandardCommitment(versionedCert)

	result := &StoreResult{
		SerializedDACert: daCommitment.Encode(),
	}

	return result, nil
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
func (h *Handlers) GenerateProof(
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
// and is verified against memory pre-state agreed upon by all challenging parties
//
// There’s no need for appending additional proof metadata for a one step proof tx
// contesting DA Cert validity
//
// TODO: Assuming we have to manage a custom fork of nitro, should we remove the proof enhancement step for
// ValidateCert opcode given the client<>server latency introduced given its noop? Then again,
// this is only ever called in the worst case one step proof WHEN the determined canonnical prestate between
// challengers is the step before calling a ValidateCert type opcode so performance considerations are rather
// irrelevant
func (h *Handlers) GenerateCertificateValidityProof(
	ctx context.Context,
	preimageType hexutil.Uint,
	certificate hexutil.Bytes,
) (GenerateCertificateValidityProofResult, error) {
	return GenerateCertificateValidityProofResult{
		Proof: []byte{},
	}, nil
}
