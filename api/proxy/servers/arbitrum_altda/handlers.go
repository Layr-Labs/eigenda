package arbitrum_altda

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	proxy_common "github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/certs"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/commitments"
	"github.com/Layr-Labs/eigenda/api/proxy/store"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

/*
	This is a (hopefully) comprehensive handlers blue print for introducing a new ALT DA server type
	that's compatible with Arbitrum's upcoming Custom DA spec.

	TODO: Understand what fork management for our Arbitrum forks will look like; at a high level we need to:
			1. test E2E correctness of the nitro stack with EigenDA
			2. introduce missing key security checks that could impact the integration's L2 Beat assessment

	TODO: Method implementations:
		[X] GetSupportedHeaderBytes // trusted integration
		[-] Store // trusted integration
		[X] RecoveryPayload // trusted integration
		[-] CollectPreimages // trusted integration
		[ ] GenerateProof // trustless AND secure integration
		[ ] GenerateCertificateValidityProof // trustless AND secure integration
*/

// IHandlers defines the expected JSON RPC interface as defined per Arbitrum Nitro's Custom DA interface:
// https://github.com/OffchainLabs/nitro/blob/c1bdcd8c571c1b22fdcdd4cc030a8ff49cbc5184/daprovider/daclient/daclient.go
type IHandlers interface {
	GetSupportedHeaderBytes(ctx context.Context) (*SupportedHeaderBytesResult, error)

	RecoverPayload(
		ctx context.Context,
		batchNum hexutil.Uint64,
		batchBlockHash common.Hash,
		sequencerMsg hexutil.Bytes,
	) (*PayloadResult, error)

	CollectPreimages(
		ctx context.Context,
		batchNum hexutil.Uint64,
		batchBlockHash common.Hash,
		sequencerMsg hexutil.Bytes,
	) (*PreimagesResult, error)

	Store(
		ctx context.Context,
		message hexutil.Bytes,
		timeout hexutil.Uint64,
	) (*StoreResult, error)

	GenerateReadPreimageProof(
		ctx context.Context,
		certHash common.Hash,
		offset hexutil.Uint64,
		certificate hexutil.Bytes,
	) (*GenerateReadPreimageProofResult, error)

	GenerateCertificateValidityProof(
		ctx context.Context,
		certificate hexutil.Bytes,
	) (*GenerateCertificateValidityProofResult, error)

	CompatibilityConfig(ctx context.Context) (*CompatibilityConfigResult, error)
}

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
	///      This can result in std out consistency issues since this is a geth native logger where we use a
	//       custom logger maintained in https://github.com/Layr-Labs/eigensdk-go/tree/dev/logging.
	//
	//       We should dig into this underlying logging and see if there's a way to intuitively override, disable,
	//       or enforce consistency between log outputs.

	log              logging.Logger
	eigenDAManager   *store.EigenDAManager
	compatibilityCfg proxy_common.CompatibilityConfig
}

// NewHandlers is a constructor
func NewHandlers(m *store.EigenDAManager, l logging.Logger, compatCfg proxy_common.CompatibilityConfig) IHandlers {
	return &Handlers{
		log:              l,
		eigenDAManager:   m,
		compatibilityCfg: compatCfg,
	}
}

// GetSupportedHeaderBytes returns the supported DA Header bytes by the CustomDA server
// this method is designed to return a span of bytes for compatibility with
// Arbitrum AnyTrust where multiple message types are supported.
// For CustomDA the provider only returns the Arbitrum CustomDA header byte.
func (h *Handlers) GetSupportedHeaderBytes(ctx context.Context) (*SupportedHeaderBytesResult, error) {
	h.logMethodCall(MethodGetSupportedHeaderBytes)

	return &SupportedHeaderBytesResult{
		HeaderBytes: []byte{commitments.ArbCustomDAHeaderByte},
	}, nil
}

// deserializeCertFromSequencerMsg reads the VersionedCert from the raw sequencer message provided
// by the DA Client
func (h *Handlers) deserializeCertFromSequencerMsg(sequencerMsg hexutil.Bytes) (*certs.VersionedCert, error) {
	if len(sequencerMsg) <= DACertOffset {
		return nil,
			fmt.Errorf("sequencer message expected to be >%d bytes, got: %d",
				DACertOffset, len(sequencerMsg))
	}

	daCommit := sequencerMsg[MessageHeaderOffset:]

	daHeaderByte := daCommit[0]
	if daHeaderByte != commitments.ArbCustomDAHeaderByte {
		return nil,
			fmt.Errorf("expected CustomDAHeader byte (%x) for 0th index byte of message, instead got: %x ",
				commitments.ArbCustomDAHeaderByte, daHeaderByte)
	}

	daLayerByte := daCommit[1]
	if daLayerByte != commitments.EigenDALayerByte {
		return nil,
			fmt.Errorf("expected EigenDALayer byte (%x) for 1st index byte of message, instead got: %x ",
				commitments.EigenDALayerByte, daLayerByte)
	}

	certVersionByte := daCommit[2]
	versionedCert := certs.NewVersionedCert([]byte(daCommit[DACommitPrefixBytes+1:]), certs.VersionByte(certVersionByte))
	return versionedCert, nil
}

// logMethodCall logs the method call with timing information and allows caller to pass in
// method specific log context
func (h *Handlers) logMethodCall(method string, logValue ...string) func() {
	start := time.Now()

	return func() {
		h.log.Info(method, "ns", time.Since(start).Nanoseconds(), logValue)
	}
}

// RecoverPayload is used to fetch the rollup payload of
// of the dispersed batch provided the DA Cert bytes.
//
// @param batch_num: batch number position in global state sequence
// @param batch_block_hash: block hash of the certL1InclusionBlock
// @param sequencer_msg: The encoded rollup payload
//
// @return bytes: Rollup payload bytes
// @return error: A structured error message (if applicable)
func (h *Handlers) RecoverPayload(
	ctx context.Context,
	batchNum hexutil.Uint64,
	batchBlockHash common.Hash,
	sequencerMsg hexutil.Bytes,
) (*PayloadResult, error) {
	callBack := h.logMethodCall(MethodRecoverPayload, "sequencer_message", sequencerMsg.String())
	defer callBack()

	// if the DA Cert fails to be deserialized from the SequencerMessage
	// then it is treated as a DerivationError
	daCert, err := h.deserializeCertFromSequencerMsg(sequencerMsg)
	if err != nil {
		return nil, nil
	}

	payload, err := h.eigenDAManager.Get(ctx, daCert, coretypes.CertSerializationABI, proxy_common.GETOpts{})
	if err != nil {
		var dpError *coretypes.DerivationError
		if errors.As(err, &dpError) {
			// returning nil for the batch payload indicates to the
			// nitro derivation pipeline to "discard" this batch and move
			// onto the next DA Cert in the Sequencer Inbox
			return nil, nil
		}

		return nil, fmt.Errorf("get rollup payload from DA Cert: %w", err)
	}

	return &PayloadResult{
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
// TODO: Determine the encoding standard to use for the returned DA Commitment. It's assumed that an EigenDAV2 message
// header byte will be prefixed. We can likely reuse the Standard Commitment mode but will require some analysis.
//
// TODO: Add processing for client provided timeout value
func (h *Handlers) Store(
	ctx context.Context,
	message hexutil.Bytes,
	timeout hexutil.Uint64,
) (*StoreResult, error) {
	callBack := h.logMethodCall(MethodStore)
	defer callBack()

	dispersalBackend := h.eigenDAManager.GetDispersalBackend()
	if dispersalBackend != proxy_common.V2EigenDABackend {
		return nil, fmt.Errorf("expected EigenDAV2 backend, got: %v", dispersalBackend)
	}

	if len(message) == 0 {
		return nil, fmt.Errorf("received empty rollup payload")
	}

	versionedCert, err := h.eigenDAManager.Put(ctx, message, coretypes.CertSerializationABI)
	if err != nil {
		return nil, fmt.Errorf("put rollup payload: %w", err)
	}
	daCommitment := commitments.NewArbCommitment(*versionedCert)

	result := &StoreResult{
		SerializedDACert: daCommitment.Encode(),
	}

	return result, nil
}

// NOTE: The validation pipeline for CustomDA in Arbitrum is currently unimplemented
// meaning a consensus artifact cannot be generated which reads CustomDA rollup payloads
//
// CollectPreimages fetches the "polynomial evaluation form" (not yet) of the dispersed rollup payload
// and inserts it as a value into a PreimageMap using the hash of the DA Cert as the
// preimage key
//
// @param batch_num: batch number position in global state sequence
// @param batch_block_hash: block hash of the certL1InclusionBlock
// @param sequencer_msg: The DA Certificate
//
//	@return preimages_result: preimage mapping that contains EigenDA V2 entry
//	@return error: a structured error message (if applicable)
//
// TODO: Figure out whether there's value in determining "invalid cert" errors here.
//
//	In theory this is only ever be callable when a DA Cert is validated by the ValidateCert
//	opcode and is assumed to be correct and the associated blob is assumed to be available
//	making validation signaling not needed.
func (h *Handlers) CollectPreimages(
	ctx context.Context,
	batchNum hexutil.Uint64,
	batchBlockHash common.Hash,
	sequencerMsg hexutil.Bytes,
) (*PreimagesResult, error) {
	callBack := h.logMethodCall(MethodCollectPreimages, "sequencer_message", sequencerMsg.String())
	defer callBack()

	daCert, err := h.deserializeCertFromSequencerMsg(sequencerMsg)
	if err != nil {
		return nil, fmt.Errorf("deserialize cert: %w", err)
	}

	payload, err := h.eigenDAManager.Get(ctx, daCert,
		coretypes.CertSerializationABI, proxy_common.GETOpts{})
	if err != nil {
		var dpError *coretypes.DerivationError
		if errors.As(err, &dpError) {
			// returning nil for the batch payload indicates to the
			// nitro derivation pipeline to "discard" this batch and move
			// onto the next DA Cert in the Sequencer Inbox
			return nil, nil
		}

		return nil, fmt.Errorf("get rollup payload from DA Cert: %w", err)
	}

	preimages := make(PreimagesMap)
	preimageRecorder := RecordPreimagesTo(preimages)

	// Record the mapping from certificate hash to actual payload data
	// This is what the replay binary expects: keccak256(certificate) -> payload
	certHash := crypto.Keccak256Hash(sequencerMsg[MessageHeaderOffset:])
	preimageRecorder(certHash, payload, CustomDAPreimageType)

	return &PreimagesResult{
		Preimages: preimages,
	}, nil
}

// GenerateReadPreimageProof is used to prove a 32 byte CustomDA preimage type for READPREIMAGE
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
func (h *Handlers) GenerateReadPreimageProof(
	ctx context.Context,
	certHash common.Hash,
	offset hexutil.Uint64,
	certificate hexutil.Bytes,
) (*GenerateReadPreimageProofResult, error) {
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
	certificate hexutil.Bytes,
) (*GenerateCertificateValidityProofResult, error) {
	return &GenerateCertificateValidityProofResult{
		Proof: []byte{},
	}, nil
}

// CompatibilityConfig returns compatibility values an external service can use to verify compatibility between
// the proxy instance and itself. E.g version, recency window, apis enabled.
// Note: This is not part of the Custom DA spec.
func (h *Handlers) CompatibilityConfig(ctx context.Context) (*CompatibilityConfigResult, error) {
	return &CompatibilityConfigResult{
		CompatibilityConfig: h.compatibilityCfg,
	}, nil
}
