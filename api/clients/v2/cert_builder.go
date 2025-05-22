package clients

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	coreEth "github.com/Layr-Labs/eigenda/core/eth"

	disperser "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/common"
	certTypesBinding "github.com/Layr-Labs/eigenda/contracts/bindings/IEigenDACertTypeBindings"
	opsrbinding "github.com/Layr-Labs/eigenda/contracts/bindings/OperatorStateRetriever"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

type CertBuilder struct {
	logger                  logging.Logger
	opsrCaller              *opsrbinding.ContractOperatorStateRetrieverCaller
	registryCoordinatorAddr gethcommon.Address
}

// NewCertBuilder constructs a new CertBuilder instance used to build EigenDA certificates
// across different versions.
func NewCertBuilder(
	logger logging.Logger,
	opsrAddr gethcommon.Address,
	registryCoordinatorAddr gethcommon.Address,
	ethClient common.EthClient,
) (*CertBuilder, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}
	// Create the Operator State Retriever caller
	opsrCaller, err := opsrbinding.NewContractOperatorStateRetrieverCaller(opsrAddr, ethClient)
	if err != nil {
		return nil, fmt.Errorf("create operator state retriever caller: %w", err)
	}

	return &CertBuilder{
		logger:                  logger,
		opsrCaller:              opsrCaller,
		registryCoordinatorAddr: registryCoordinatorAddr,
	}, nil
}

// BuildCert builds an EigenDA certificate of the specified version using the provided blob key and blob status reply.
func (cb *CertBuilder) BuildCert(
	ctx context.Context,
	certVersion coretypes.CertificateVersion,
	blobStatusReply *disperser.BlobStatusReply,
) (coretypes.EigenDACert, error) {
	switch certVersion {
	case coretypes.VersionThreeCert:
		return cb.buildEigenDAV3Cert(ctx, blobStatusReply)
	default:
		return nil, fmt.Errorf("unsupported EigenDA cert version: %d", certVersion)
	}
}

// buildEigenDAV3Cert builds an EigenDA certificate of version 3 using the provided blob key and blob status reply.
func (cb *CertBuilder) buildEigenDAV3Cert(
	ctx context.Context,
	blobStatusReply *disperser.BlobStatusReply,
) (*coretypes.EigenDACertV3, error) {
	nonSignerStakesAndSignature, err := cb.getNonSignerStakesAndSignature(
		ctx, blobStatusReply.GetSignedBatch())
	if err != nil {
		return nil, fmt.Errorf("get non signer stake and signature: %w", err)
	}

	eigenDACert, err := coretypes.BuildEigenDACertV3(blobStatusReply, nonSignerStakesAndSignature)
	if err != nil {
		return nil, fmt.Errorf("build eigenda v3 cert: %w", err)
	}

	return eigenDACert, nil
}

// GetNonSignerStakesAndSignature constructs a NonSignerStakesAndSignature object by calling an
// onchain OperatorStateRetriever retriever to fetch necessary non-signer metadata
func (cb *CertBuilder) getNonSignerStakesAndSignature(
	ctx context.Context,
	signedBatch *disperser.SignedBatch,
) (*certTypesBinding.EigenDATypesV1NonSignerStakesAndSignature, error) {
	// 1 - Pre-process inputs for operator state retriever call
	signedBatchBinding, err := coretypes.SignedBatchProtoToV2CertBinding(signedBatch)
	if err != nil {
		return nil, fmt.Errorf("convert signed batch: %w", err)
	}

	nonSignerPubKeys := signedBatchBinding.Attestation.NonSignerPubkeys

	// 2a - create operator IDs by hashing non-signer public keys
	nonSignerOperatorIDs := make([][32]byte, len(nonSignerPubKeys))
	for i, pubKeySet := range nonSignerPubKeys {
		g1Point := core.NewG1Point(pubKeySet.X, pubKeySet.Y)
		nonSignerOperatorIDs[i] = coreEth.HashPubKeyG1(g1Point)
	}

	// 2b - cast []uint32 to []byte for quorum numbers
	quorumNumbers, err := coretypes.QuorumNumbersUint32ToUint8(signedBatchBinding.Attestation.QuorumNumbers)
	if err != nil {
		return nil, fmt.Errorf("convert quorum numbers: %w", err)
	}

	// use the reference block # from the disperser generated signed batch header
	// for referencing operator states at a specific block checkpoint
	rbn := signedBatch.GetHeader().GetReferenceBlockNumber()

	// 3 - call operator state retriever to fetch signature indices
	checkSigIndices, err := cb.opsrCaller.GetCheckSignaturesIndices(&bind.CallOpts{Context: ctx, BlockNumber: big.NewInt(int64(rbn))},
		cb.registryCoordinatorAddr, uint32(rbn), quorumNumbers, nonSignerOperatorIDs)

	if err != nil {
		return nil, fmt.Errorf("check sig indices call: %w", err)
	}

	// 4 - translate from CertVerifier binding types to cert type
	// TODO: Should probably put SignedBatch into the types directly to avoid this downstream conversion
	nonSignerPubKeysBN254 := make([]certTypesBinding.BN254G1Point, len(signedBatchBinding.Attestation.NonSignerPubkeys))
	for i, pubKeySet := range signedBatchBinding.Attestation.NonSignerPubkeys {
		nonSignerPubKeysBN254[i] = certTypesBinding.BN254G1Point{
			X: pubKeySet.X,
			Y: pubKeySet.Y,
		}
	}

	quorumApksBN254 := make([]certTypesBinding.BN254G1Point, len(signedBatchBinding.Attestation.QuorumApks))
	for i, apkSet := range signedBatchBinding.Attestation.QuorumApks {
		quorumApksBN254[i] = certTypesBinding.BN254G1Point{
			X: apkSet.X,
			Y: apkSet.Y,
		}
	}

	apkG2BN254 := certTypesBinding.BN254G2Point{
		X: signedBatchBinding.Attestation.ApkG2.X,
		Y: signedBatchBinding.Attestation.ApkG2.Y,
	}

	sigmaBN254 := certTypesBinding.BN254G1Point{
		X: signedBatchBinding.Attestation.Sigma.X,
		Y: signedBatchBinding.Attestation.Sigma.Y,
	}

	// 5 - construct non signer stakes and signature
	return &certTypesBinding.EigenDATypesV1NonSignerStakesAndSignature{
		NonSignerQuorumBitmapIndices: checkSigIndices.NonSignerQuorumBitmapIndices,
		NonSignerPubkeys:             nonSignerPubKeysBN254,
		QuorumApks:                   quorumApksBN254,
		ApkG2:                        apkG2BN254,
		Sigma:                        sigmaBN254,
		QuorumApkIndices:             checkSigIndices.QuorumApkIndices,
		TotalStakeIndices:            checkSigIndices.TotalStakeIndices,
		NonSignerStakeIndices:        checkSigIndices.NonSignerStakeIndices,
	}, nil
}
