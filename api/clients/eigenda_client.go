package clients

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"

	"github.com/Layr-Labs/eigenda/api"
	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	grpcdisperser "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	edasm "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/auth"
)

// IEigenDAClient is a wrapper around the DisperserClient interface which
// encodes blobs before dispersing them, and decodes them after retrieving them.
type IEigenDAClient interface {
	GetBlob(ctx context.Context, batchHeaderHash []byte, blobIndex uint32) ([]byte, error)
	PutBlob(ctx context.Context, txData []byte) (*grpcdisperser.BlobInfo, error)
	GetCodec() codecs.BlobCodec
	Close() error
}

// See the NewEigenDAClient constructor's documentation for details and usage examples.
// TODO: Refactor this struct and interface above to use same naming convention as disperser client.
//
//	Also need to make the fields private and use the constructor in the tests.
type EigenDAClient struct {
	// TODO: all of these should be private, to prevent users from using them directly,
	// which breaks encapsulation and makes it hard for us to do refactors or changes
	Config      EigenDAClientConfig
	Log         log.Logger
	Client      DisperserClient
	ethClient   *ethclient.Client
	edasmCaller *edasm.ContractEigenDAServiceManagerCaller
	Codec       codecs.BlobCodec
}

var _ IEigenDAClient = &EigenDAClient{}

// EigenDAClient is a wrapper around the DisperserClient which
// encodes blobs before dispersing them, and decodes them after retrieving them.
// It also turns the disperser's async polling-based API (disperseBlob + poll GetBlobStatus)
// into a sync API where PutBlob will poll for the blob to be confirmed or finalized.
//
// DisperserClient is safe to be used concurrently by multiple goroutines.
// Don't forget to call Close() on the client when you're done with it, to close the
// underlying grpc connection maintained by the DiserserClient.
//
// Example usage:
//
//	client, err := NewEigenDAClient(log, EigenDAClientConfig{...})
//	if err != nil {
//	  return err
//	}
//	defer client.Close()
//
//	blobData := []byte("hello world")
//	blobInfo, err := client.PutBlob(ctx, blobData)
//	if err != nil {
//	  return err
//	}
//
//	retrievedData, err := client.GetBlob(ctx, blobInfo.BatchMetadata.BatchHeaderHash, blobInfo.BlobIndex)
//	if err != nil {
//	  return err
//	}
func NewEigenDAClient(log log.Logger, config EigenDAClientConfig) (*EigenDAClient, error) {
	err := config.CheckAndSetDefaults()
	if err != nil {
		return nil, err
	}

	var ethClient *ethclient.Client
	var edasmCaller *edasm.ContractEigenDAServiceManagerCaller
	ethClient, err = ethclient.Dial(config.EthRpcUrl)
	if err != nil {
		return nil, fmt.Errorf("dial ETH RPC node: %w", err)
	}
	edasmCaller, err = edasm.NewContractEigenDAServiceManagerCaller(common.HexToAddress(config.SvcManagerAddr), ethClient)
	if err != nil {
		return nil, fmt.Errorf("new EigenDAServiceManagerCaller: %w", err)
	}

	host, port, err := net.SplitHostPort(config.RPC)
	if err != nil {
		return nil, fmt.Errorf("parse EigenDA RPC: %w", err)
	}

	var signer core.BlobRequestSigner
	if len(config.SignerPrivateKeyHex) == 64 {
		signer = auth.NewLocalBlobRequestSigner(config.SignerPrivateKeyHex)
	} else if len(config.SignerPrivateKeyHex) == 0 {
		// noop signer is used when we need a read-only eigenda client
		signer = auth.NewLocalNoopSigner()
	} else {
		return nil, fmt.Errorf("invalid length for signer private key")
	}

	disperserConfig := NewConfig(host, port, config.ResponseTimeout, !config.DisableTLS)

	disperserClient, err := NewDisperserClient(disperserConfig, signer)
	if err != nil {
		return nil, fmt.Errorf("new disperser-client: %w", err)
	}

	lowLevelCodec, err := codecs.BlobEncodingVersionToCodec(config.PutBlobEncodingVersion)
	if err != nil {
		return nil, fmt.Errorf("create low level codec: %w", err)
	}

	var codec codecs.BlobCodec
	if config.DisablePointVerificationMode {
		codec = codecs.NewNoIFFTCodec(lowLevelCodec)
	} else {
		codec = codecs.NewIFFTCodec(lowLevelCodec)
	}

	return &EigenDAClient{
		Log:         log,
		Config:      config,
		Client:      disperserClient,
		ethClient:   ethClient,
		edasmCaller: edasmCaller,
		Codec:       codec,
	}, nil
}

func (m *EigenDAClient) GetCodec() codecs.BlobCodec {
	return m.Codec
}

// GetBlob retrieves a blob from the EigenDA service using the provided context,
// batch header hash, and blob index.  If decode is set to true, the function
// decodes the retrieved blob data. If set to false it returns the encoded blob
// data, which is necessary for generating KZG proofs for data's correctness.
// The function handles potential errors during blob retrieval, data length
// checks, and decoding processes.
func (m *EigenDAClient) GetBlob(ctx context.Context, batchHeaderHash []byte, blobIndex uint32) ([]byte, error) {
	data, err := m.Client.RetrieveBlob(ctx, batchHeaderHash, blobIndex)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve blob: %w", err)
	}

	if len(data) == 0 {
		// This should never happen, because empty blobs are rejected from even entering the system:
		// https://github.com/Layr-Labs/eigenda/blob/master/disperser/apiserver/server.go#L930
		return nil, fmt.Errorf("blob has length zero - this should not be possible")
	}

	decodedData, err := m.Codec.DecodeBlob(data)
	if err != nil {
		return nil, fmt.Errorf("error decoding blob: %w", err)
	}

	return decodedData, nil
}

// PutBlob encodes and writes a blob to EigenDA, waiting for a desired blob status
// to be reached (guarded by WaitForFinalization config param) before returning.
//
// TODO: describe retry/timeout behavior
//
// Upon return the blob is guaranteed to be:
//   - finalized onchain (if Config.WaitForFinalization is true), or
//   - confirmed at a certain depth (if Config.WaitForFinalization is false,
//     in which case Config.WaitForConfirmationDepth specifies the depth).
//
// Errors returned all either grpc errors, or api.ErrorFailover, for eg:
//
//	 blobInfo, err := client.PutBlob(ctx, blobData)
//	 if err != nil {
//	   if errors.Is(err, api.ErrorFailover) {
//	     // failover to ethda
//		  }
//	   st, isGRPCError := status.FromError(err)
//	   if isGRPCError {
//	     // use st.Code() and st.Message()
//	   } else {
//	     // assert this shouldn't happen
//	   }
//	 }
//
// An api.ErrorFailover error returned is used to signify that eigenda is temporarily unavailable,
// and suggest to the caller (most likely some rollup batcher via the eigenda-proxy)
// to fallback to ethda for some amount of time. Three reasons for returning api.ErrorFailover:
//  1. Failed to put the blob in the disperser's queue (disperser is down)
//  2. Timed out before getting confirmed onchain (batcher is down)
//  3. Insufficient signatures (eigenda network is down)
//
// See https://github.com/ethereum-optimism/specs/issues/434 for more details.
func (m *EigenDAClient) PutBlob(ctx context.Context, data []byte) (*grpcdisperser.BlobInfo, error) {
	resultChan, errorChan := m.PutBlobAsync(ctx, data)
	select { // no timeout here because we depend on the configured timeout in PutBlobAsync
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}

func (m *EigenDAClient) PutBlobAsync(ctx context.Context, data []byte) (resultChan chan *grpcdisperser.BlobInfo, errChan chan error) {
	resultChan = make(chan *grpcdisperser.BlobInfo, 1)
	errChan = make(chan error, 1)
	go m.putBlob(ctx, data, resultChan, errChan)
	return
}

func (m *EigenDAClient) putBlob(ctxFinality context.Context, rawData []byte, resultChan chan *grpcdisperser.BlobInfo, errChan chan error) {
	m.Log.Info("Attempting to disperse blob to EigenDA")

	// encode blob
	if m.Codec == nil {
		errChan <- api.NewErrorInternal("codec not initialized")
		return
	}

	data, err := m.Codec.EncodeBlob(rawData)
	if err != nil {
		// Encode can only fail if there is something wrong with the data, so we return a 400 error
		errChan <- api.NewErrorInvalidArg(fmt.Sprintf("error encoding blob: %v", err))
		return
	}

	customQuorumNumbers := make([]uint8, len(m.Config.CustomQuorumIDs))
	for i, e := range m.Config.CustomQuorumIDs {
		customQuorumNumbers[i] = uint8(e)
	}
	// disperse blob
	// TODO: would be nice to add a trace-id key to the context, to be able to follow requests from batcher->proxy->eigenda
	// clients with a payment signer setting can disperse paid blobs
	_, requestID, err := m.Client.DisperseBlobAuthenticated(ctxFinality, data, customQuorumNumbers)
	if err != nil {
		// DisperserClient returned error is already a grpc error which can be a 400 (eg rate limited) or 500,
		// so we wrap the error such that clients can still use grpc's status.FromError() function to get the status code.
		errChan <- fmt.Errorf("error submitting authenticated blob to disperser: %w", err)
		return
	}

	base64RequestID := base64.StdEncoding.EncodeToString(requestID)
	m.Log.Info("Blob accepted by EigenDA disperser, now polling for status updates", "requestID", base64RequestID)

	ticker := time.NewTicker(m.Config.StatusQueryRetryInterval)
	defer ticker.Stop()

	confirmationCh := time.NewTimer(m.Config.ConfirmationTimeout).C
	var cancel context.CancelFunc
	// finality here can either mean reaching some confirmationDepth or reaching actual finality
	// depending on the WaitForFinalization config param.
	ctxFinality, cancel = context.WithTimeout(ctxFinality, m.Config.StatusQueryTimeout)
	defer cancel()

	alreadyWaitingForDispersal := false
	alreadyWaitingForConfirmationOrFinality := false
	var latestBlobStatus grpcdisperser.BlobStatus
	for {
		select {
		// The two first timeout cases can only happen while blob is still in
		// 1. processing or dispersing status: waiting to land onchain
		// 2. or confirmed status: landed onchain, waiting for finalization
		// because all other statuses return immediately once reached (see below).
		case <-confirmationCh:
			if latestBlobStatus == grpcdisperser.BlobStatus_PROCESSING || latestBlobStatus == grpcdisperser.BlobStatus_DISPERSING {
				errChan <- api.NewErrorFailover(fmt.Errorf("eigenda might be down. timed out waiting for blob to land onchain (request id=%s): %w", base64RequestID, ctxFinality.Err()))
			}
			// set to nil so this case doesn't get triggered again
			confirmationCh = nil
		case <-ctxFinality.Done():
			// this should have been triggered above because confirmationTimeout < ctxFinality timeout,
			// but we leave this assert here as a safety net.
			if latestBlobStatus == grpcdisperser.BlobStatus_PROCESSING || latestBlobStatus == grpcdisperser.BlobStatus_DISPERSING {
				errChan <- api.NewErrorFailover(fmt.Errorf("eigenda might be down. timed out waiting for blob to land onchain (request id=%s): %w", base64RequestID, ctxFinality.Err()))
			} else if latestBlobStatus == grpcdisperser.BlobStatus_CONFIRMED {
				// Assuming that the ctxFinality timeout is correctly set (long enough for batch to land onchain + finalize),
				// still being in confirmed state here means that there is a problem with Ethereum, so we return DeadlineExceeded (504).
				// batcher would most likely resubmit another blob, which is not ideal but there isn't much to be done...
				// eigenDA v2 will have idempotency so one can just resubmit the same blob safely.
				// TODO: (if timeout was not long enough to finalize in normal conditions): eigenda-client is badly configured, should be a 400 (INVALID_ARGUMENT)
				errChan <- api.NewErrorDeadlineExceeded(
					fmt.Sprintf("timed out waiting for blob that landed onchain to finalize (request id=%s). "+
						"Either timeout not long enough, or ethereum might be experiencing difficulties: %v. ", base64RequestID, ctxFinality.Err()))
			} else {
				// this should not be reachable... indicates something wrong with either this client or eigenda, so we failover to ethda
				errChan <- api.NewErrorFailover(fmt.Errorf("timed out in a state that shouldn't be possible (request id=%s): %w", base64RequestID, ctxFinality.Err()))
			}
			return
		case <-ticker.C:
			statusRes, err := m.Client.GetBlobStatus(ctxFinality, requestID)
			if err != nil {
				m.Log.Warn("Unable to retrieve blob dispersal status, will retry", "requestID", base64RequestID, "err", err)
				continue
			}
			latestBlobStatus = statusRes.Status
			switch statusRes.Status {
			case grpcdisperser.BlobStatus_PROCESSING, grpcdisperser.BlobStatus_DISPERSING:
				// to prevent log clutter, we only log at info level once
				if alreadyWaitingForDispersal {
					m.Log.Debug("Blob is being processed by the EigenDA network", "requestID", base64RequestID)
				} else {
					m.Log.Info("Blob is being processed by the EigenDA network", "requestID", base64RequestID)
					alreadyWaitingForDispersal = true
				}
			case grpcdisperser.BlobStatus_FAILED:
				// This can happen for a few reasons:
				// 1. blob has expired, a client retrieve after 14 days. Sounds like 400 errors, but not sure this can happen during dispersal...
				// 2. internal logic error while requesting encoding (shouldn't happen), but should probably return api.ErrorFailover
				// 3. wait for blob finalization from confirmation and blob retry has exceeded its limit.
				//    Probably from a chain re-org. See https://github.com/Layr-Labs/eigenda/blob/master/disperser/batcher/finalizer.go#L179-L189.
				//    So we should be returning 500 to force a blob resubmission (not eigenda's fault but until
				//    we have idempotency this is unfortunately the only solution)
				// TODO: we should create new BlobStatus categories to separate these cases out. For now returning 500 is fine.
				errChan <- api.NewErrorInternal(fmt.Sprintf("blob dispersal (requestID=%s) reached failed status. please resubmit the blob.", base64RequestID))
				return
			case grpcdisperser.BlobStatus_INSUFFICIENT_SIGNATURES:
				// Some quorum failed to sign the blob, indicating that the whole network is having issues.
				// We hence return api.ErrorFailover to let the batcher failover to ethda. This could however be a very unlucky
				// temporary issue, so the caller should retry at least one more time before failing over.
				errChan <- api.NewErrorFailover(fmt.Errorf("blob dispersal (requestID=%s) failed with insufficient signatures. eigenda nodes are probably down.", base64RequestID))
				return
			case grpcdisperser.BlobStatus_CONFIRMED:
				if m.Config.WaitForFinalization {
					// to prevent log clutter, we only log at info level once
					if alreadyWaitingForConfirmationOrFinality {
						m.Log.Debug("EigenDA blob included onchain, waiting for finalization", "requestID", base64RequestID)
					} else {
						m.Log.Info("EigenDA blob included onchain, waiting for finalization", "requestID", base64RequestID)
						alreadyWaitingForConfirmationOrFinality = true
					}
				} else {
					batchId := statusRes.Info.BlobVerificationProof.GetBatchId()
					batchConfirmed, err := m.batchIdConfirmedAtDepth(ctxFinality, batchId, m.Config.WaitForConfirmationDepth)
					if err != nil {
						m.Log.Warn("Error checking if batch ID is confirmed at depth. Will retry...", "requestID", base64RequestID, "err", err)
					}
					if batchConfirmed {
						m.Log.Info("EigenDA blob confirmed", "requestID", base64RequestID, "confirmationDepth", m.Config.WaitForConfirmationDepth)
						resultChan <- statusRes.Info
						return
					}
					// to prevent log clutter, we only log at info level once
					if alreadyWaitingForConfirmationOrFinality {
						m.Log.Debug("EigenDA blob included onchain, waiting for confirmation", "requestID", base64RequestID, "confirmationDepth", m.Config.WaitForConfirmationDepth)
					} else {
						m.Log.Info("EigenDA blob included onchain, waiting for confirmation", "requestID", base64RequestID, "confirmationDepth", m.Config.WaitForConfirmationDepth)
						alreadyWaitingForConfirmationOrFinality = true
					}
				}
			case grpcdisperser.BlobStatus_FINALIZED:
				batchHeaderHashHex := fmt.Sprintf("0x%s", hex.EncodeToString(statusRes.Info.BlobVerificationProof.BatchMetadata.BatchHeaderHash))
				m.Log.Info("EigenDA blob finalized", "requestID", base64RequestID, "batchHeaderHash", batchHeaderHashHex)
				resultChan <- statusRes.Info
				return
			default:
				// This should never happen. If it does, the blob is in a heisenberg state... it could either eventually get confirmed or fail.
				// However, this doesn't mean there's a major outage with EigenDA, so we return a 500 error to let the caller redisperse the blob,
				// rather than an api.ErrorFailover to failover to EthDA.
				errChan <- api.NewErrorInternal(fmt.Sprintf("unknown reply status %d. ask for assistance from EigenDA team, using requestID %s", statusRes.Status, base64RequestID))
				return
			}
		}
	}
}

// Close simply calls Close() on the wrapped disperserClient, to close the grpc connection to the disperser server.
// It is thread safe and can be called multiple times.
func (c *EigenDAClient) Close() error {
	return c.Client.Close()
}

// getConfDeepBlockNumber returns the block number that is `depth` blocks behind the current block number.
func (m EigenDAClient) getConfDeepBlockNumber(ctx context.Context, depth uint64) (*big.Int, error) {
	curBlockNumber, err := m.ethClient.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block number: %w", err)
	}
	// If curBlock < depth, this will return the genesis block number (0),
	// which would cause to accept as confirmed a block that isn't depth deep.
	// TODO: there's prob a better way to deal with this, like returning a special error
	if curBlockNumber < depth {
		return big.NewInt(0), nil
	}
	return new(big.Int).SetUint64(curBlockNumber - depth), nil
}

// batchIdConfirmedAtDepth checks if a batch ID has been confirmed at a certain depth.
// It returns true if the batch ID has been confirmed at the given depth, and false otherwise,
// or returns an error if any of the network calls fail.
func (m EigenDAClient) batchIdConfirmedAtDepth(ctx context.Context, batchId uint32, depth uint64) (bool, error) {
	confDeepBlockNumber, err := m.getConfDeepBlockNumber(ctx, depth)
	if err != nil {
		return false, fmt.Errorf("failed to get confirmation deep block number: %w", err)
	}
	onchainBatchMetadataHash, err := m.edasmCaller.BatchIdToBatchMetadataHash(&bind.CallOpts{BlockNumber: confDeepBlockNumber}, batchId)
	if err != nil {
		return false, fmt.Errorf("failed to get batch metadata hash: %w", err)
	}
	if bytes.Equal(onchainBatchMetadataHash[:], make([]byte, 32)) {
		return false, nil
	}
	return true, nil
}
