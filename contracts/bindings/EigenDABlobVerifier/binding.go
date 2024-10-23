// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractEigenDABlobVerifier

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// BN254G1Point is an auto generated low-level Go binding around an user-defined struct.
type BN254G1Point struct {
	X *big.Int
	Y *big.Int
}

// BN254G2Point is an auto generated low-level Go binding around an user-defined struct.
type BN254G2Point struct {
	X [2]*big.Int
	Y [2]*big.Int
}

// EigenDABlobVerificationUtilsAttestation is an auto generated low-level Go binding around an user-defined struct.
type EigenDABlobVerificationUtilsAttestation struct {
	NonSignerQuorumBitmapIndices []uint32
	NonSignerPubkeys             []BN254G1Point
	QuorumApks                   []BN254G1Point
	ApkG2                        BN254G2Point
	Sigma                        BN254G1Point
	QuorumApkIndices             []uint32
	TotalStakeIndices            []uint32
	NonSignerStakeIndices        [][]uint32
}

// EigenDABlobVerificationUtilsBlobCertificate is an auto generated low-level Go binding around an user-defined struct.
type EigenDABlobVerificationUtilsBlobCertificate struct {
	BlobKey              []byte
	BlobHeader           IEigenDAServiceManagerBlobHeader
	ReferenceBlockNumber uint32
	RelayKeys            []string
}

// EigenDABlobVerificationUtilsBlobVerificationProof is an auto generated low-level Go binding around an user-defined struct.
type EigenDABlobVerificationUtilsBlobVerificationProof struct {
	BatchId        uint32
	BlobIndex      uint32
	BatchMetadata  IEigenDAServiceManagerBatchMetadata
	InclusionProof []byte
	QuorumIndices  []byte
}

// EigenDABlobVerificationUtilsSignedCertificate is an auto generated low-level Go binding around an user-defined struct.
type EigenDABlobVerificationUtilsSignedCertificate struct {
	BlobCertificate             EigenDABlobVerificationUtilsBlobCertificate
	NonSignerStakesAndSignature EigenDABlobVerificationUtilsAttestation
}

// IEigenDAServiceManagerBatchHeader is an auto generated low-level Go binding around an user-defined struct.
type IEigenDAServiceManagerBatchHeader struct {
	BlobHeadersRoot       [32]byte
	QuorumNumbers         []byte
	SignedStakeForQuorums []byte
	ReferenceBlockNumber  uint32
}

// IEigenDAServiceManagerBatchMetadata is an auto generated low-level Go binding around an user-defined struct.
type IEigenDAServiceManagerBatchMetadata struct {
	BatchHeader             IEigenDAServiceManagerBatchHeader
	SignatoryRecordHash     [32]byte
	ConfirmationBlockNumber uint32
}

// IEigenDAServiceManagerBlobHeader is an auto generated low-level Go binding around an user-defined struct.
type IEigenDAServiceManagerBlobHeader struct {
	Commitment       BN254G1Point
	DataLength       uint32
	QuorumBlobParams []IEigenDAServiceManagerQuorumBlobParam
}

// IEigenDAServiceManagerQuorumBlobParam is an auto generated low-level Go binding around an user-defined struct.
type IEigenDAServiceManagerQuorumBlobParam struct {
	QuorumNumber                    uint8
	AdversaryThresholdPercentage    uint8
	ConfirmationThresholdPercentage uint8
	ChunkLength                     uint32
}

// ContractEigenDABlobVerifierMetaData contains all meta data concerning the ContractEigenDABlobVerifier contract.
var ContractEigenDABlobVerifierMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_eigenDAThresholdRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"},{\"name\":\"_eigenDABatchMetadataStorage\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"},{\"name\":\"_eigenDASignatureVerifier\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"eigenDABatchMetadataStorage\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDASignatureVerifier\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIsQuorumRequired\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumAdversaryThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumConfirmationThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumAdversaryThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumConfirmationThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV1\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structIEigenDAServiceManager.QuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structEigenDABlobVerificationUtils.BlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV1\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structIEigenDAServiceManager.QuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structEigenDABlobVerificationUtils.BlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"additionalQuorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV2\",\"inputs\":[{\"name\":\"signedCertificate\",\"type\":\"tuple\",\"internalType\":\"structEigenDABlobVerificationUtils.SignedCertificate\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structEigenDABlobVerificationUtils.BlobCertificate\",\"components\":[{\"name\":\"blobKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structIEigenDAServiceManager.QuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"relayKeys\",\"type\":\"string[]\",\"internalType\":\"string[]\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structEigenDABlobVerificationUtils.Attestation\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV2\",\"inputs\":[{\"name\":\"signedCertificate\",\"type\":\"tuple\",\"internalType\":\"structEigenDABlobVerificationUtils.SignedCertificate\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structEigenDABlobVerificationUtils.BlobCertificate\",\"components\":[{\"name\":\"blobKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structIEigenDAServiceManager.QuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"relayKeys\",\"type\":\"string[]\",\"internalType\":\"string[]\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structEigenDABlobVerificationUtils.Attestation\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}]},{\"name\":\"additionalQuorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b5060405162001c6d38038062001c6d83398101604081905262000034916200006b565b6001600160a01b0392831660805290821660a0521660c052620000bf565b6001600160a01b03811681146200006857600080fd5b50565b6000806000606084860312156200008157600080fd5b83516200008e8162000052565b6020850151909350620000a18162000052565b6040850151909250620000b48162000052565b809150509250925092565b60805160a05160c051611b2a62000143600039600081816101d80152818161035c01526106050152600081816101260152818161043f01526104970152600081816102120152818161024f015281816102e3015281816103900152818161041e01528181610476015281816104f20152818161055201526105c90152611b2a6000f3fe608060405234801561001057600080fd5b50600436106100cf5760003560e01c80638f3a8f321161008c578063ee6c3bcf11610066578063ee6c3bcf146101c0578063efd4532b146101d3578063f78a6ea1146101fa578063f8c668141461020d57600080fd5b80638f3a8f321461019d578063bafa9107146101b0578063e15234ff146101b857600080fd5b8063048886d2146100d45780631429c7c2146100fc578063640f65d9146101215780637baca37c146101605780638687feae146101755780638d67b9091461018a575b600080fd5b6100e76100e23660046110fc565b610234565b60405190151581526020015b60405180910390f35b61010f61010a3660046110fc565b6102c8565b60405160ff90911681526020016100f3565b6101487f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b0390911681526020016100f3565b61017361016e366004611138565b610357565b005b61017d61038c565b6040516100f391906111cc565b610173610198366004611203565b610419565b6101736101ab3660046112ae565b610471565b61017d6104ee565b61017d61054e565b61010f6101ce3660046110fc565b6105ae565b6101487f000000000000000000000000000000000000000000000000000000000000000081565b610173610208366004611338565b610600565b6101487f000000000000000000000000000000000000000000000000000000000000000081565b604051630244436960e11b815260ff821660048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063048886d290602401602060405180830381865afa15801561029e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102c291906113a0565b92915050565b604051630a14e3e160e11b815260ff821660048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690631429c7c2906024015b602060405180830381865afa158015610333573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102c291906113c2565b6103897f00000000000000000000000000000000000000000000000000000000000000008261038461054e565b505050565b50565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316638687feae6040518163ffffffff1660e01b8152600401600060405180830381865afa1580156103ec573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261041491908101906114b8565b905090565b61046d7f00000000000000000000000000000000000000000000000000000000000000007f0000000000000000000000000000000000000000000000000000000000000000848461046861054e565b610653565b5050565b6104e87f00000000000000000000000000000000000000000000000000000000000000007f000000000000000000000000000000000000000000000000000000000000000086866104c061054e565b87876040516020016104d493929190611533565b604051602081830303815290604052610653565b50505050565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663bafa91076040518163ffffffff1660e01b8152600401600060405180830381865afa1580156103ec573d6000803e3d6000fd5b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e15234ff6040518163ffffffff1660e01b8152600401600060405180830381865afa1580156103ec573d6000803e3d6000fd5b60405163ee6c3bcf60e01b815260ff821660048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063ee6c3bcf90602401610316565b6103847f00000000000000000000000000000000000000000000000000000000000000008461062d61054e565b858560405160200161064193929190611533565b60408051601f19818403019052525050565b6001600160a01b03841663eccbbfc961066f6020850185611574565b6040516001600160e01b031960e084901b16815263ffffffff919091166004820152602401602060405180830381865afa1580156106b1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106d5919061158f565b6106f36106e560408501856115a8565b6106ee90611619565b610da4565b1461077f5760405162461bcd60e51b81526020600482015260606024820152600080516020611ad583398151915260448201527f72696679426c6f62466f7251756f72756d733a2062617463684d65746164617460648201527f6120646f6573206e6f74206d617463682073746f726564206d65746164617461608482015260a4015b60405180910390fd5b61083361078f60608401846116e9565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506107d19250505060408501856115a8565b6107db908061172f565b356107ed6107e887611745565b610e15565b6040516020016107ff91815260200190565b604051602081830303815290604052805190602001208560200160208101906108289190611574565b63ffffffff16610e45565b6108aa5760405162461bcd60e51b815260206004820152604e6024820152600080516020611ad583398151915260448201527f72696679426c6f62466f7251756f72756d733a20696e636c7573696f6e20707260648201526d1bdbd9881a5cc81a5b9d985b1a5960921b608482015260a401610776565b6000805b6108bb6060860186611890565b9050811015610ceb576108d16060860186611890565b828181106108e1576108e16118d9565b6108f792602060809092020190810191506110fc565b60ff1661090760408601866115a8565b610911908061172f565b61091f9060208101906116e9565b61092c60808801886116e9565b8581811061093c5761093c6118d9565b919091013560f81c9050818110610955576109556118d9565b9050013560f81c60f81b60f81c60ff16146109de5760405162461bcd60e51b815260206004820152604f6024820152600080516020611ad583398151915260448201527f72696679426c6f62466f7251756f72756d733a2071756f72756d4e756d62657260648201526e040c8decae640dcdee840dac2e8c6d608b1b608482015260a401610776565b6109eb6060860186611890565b828181106109fb576109fb6118d9565b9050608002016020016020810190610a1391906110fc565b60ff16610a236060870187611890565b83818110610a3357610a336118d9565b9050608002016040016020810190610a4b91906110fc565b60ff1611610ad55760405162461bcd60e51b81526020600482015260576024820152600080516020611ad583398151915260448201527f72696679426c6f62466f7251756f72756d733a207468726573686f6c6420706560648201527f7263656e746167657320617265206e6f742076616c6964000000000000000000608482015260a401610776565b6001600160a01b038716631429c7c2610af16060880188611890565b84818110610b0157610b016118d9565b610b1792602060809092020190810191506110fc565b6040516001600160e01b031960e084901b16815260ff9091166004820152602401602060405180830381865afa158015610b55573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b7991906113c2565b60ff16610b896060870187611890565b83818110610b9957610b996118d9565b9050608002016040016020810190610bb191906110fc565b60ff161015610bd25760405162461bcd60e51b8152600401610776906118ef565b610bdf6060860186611890565b82818110610bef57610bef6118d9565b9050608002016040016020810190610c0791906110fc565b60ff16610c1760408601866115a8565b610c21908061172f565b610c2f9060408101906116e9565b610c3c60808801886116e9565b85818110610c4c57610c4c6118d9565b919091013560f81c9050818110610c6557610c656118d9565b9050013560f81c60f81b60f81c60ff161015610c935760405162461bcd60e51b8152600401610776906118ef565b610cd782610ca46060880188611890565b84818110610cb457610cb46118d9565b610cca92602060809092020190810191506110fc565b600160ff919091161b1790565b915080610ce381611976565b9150506108ae565b50610cff610cf883610e5d565b8281161490565b610d9c5760405162461bcd60e51b815260206004820152606e6024820152600080516020611ad583398151915260448201527f72696679426c6f62466f7251756f72756d733a2072657175697265642071756f60648201527f72756d7320617265206e6f74206120737562736574206f662074686520636f6e60848201526d6669726d65642071756f72756d7360901b60a482015260c401610776565b505050505050565b60006102c28260000151604051602001610dbe9190611991565b60408051808303601f1901815282825280516020918201208682015187840151838601929092528484015260e01b6001600160e01b0319166060840152815160448185030181526064909301909152815191012090565b600081604051602001610e2891906119f1565b604051602081830303815290604052805190602001209050919050565b600083610e53868585610fea565b1495945050505050565b600061010082511115610ee65760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a401610776565b8151610ef457506000919050565b60008083600081518110610f0a57610f0a6118d9565b0160200151600160f89190911c81901b92505b8451811015610fe157848181518110610f3857610f386118d9565b0160200151600160f89190911c1b9150828211610fcd5760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a401610776565b91811791610fda81611976565b9050610f1d565b50909392505050565b600060208451610ffa9190611a9a565b156110815760405162461bcd60e51b815260206004820152604b60248201527f4d65726b6c652e70726f63657373496e636c7573696f6e50726f6f664b65636360448201527f616b3a2070726f6f66206c656e6774682073686f756c642062652061206d756c60648201526a3a34b836329037b310199960a91b608482015260a401610776565b8260205b855181116110e457611098600285611a9a565b6110b9578160005280860151602052604060002091506002840493506110d2565b8086015160005281602052604060002091506002840493505b6110dd602082611abc565b9050611085565b50949350505050565b60ff8116811461038957600080fd5b60006020828403121561110e57600080fd5b8135611119816110ed565b9392505050565b60006040828403121561113257600080fd5b50919050565b60006020828403121561114a57600080fd5b81356001600160401b0381111561116057600080fd5b61116c84828501611120565b949350505050565b60005b8381101561118f578181015183820152602001611177565b838111156104e85750506000910152565b600081518084526111b8816020860160208601611174565b601f01601f19169290920160200192915050565b60208152600061111960208301846111a0565b60006080828403121561113257600080fd5b600060a0828403121561113257600080fd5b6000806040838503121561121657600080fd5b82356001600160401b038082111561122d57600080fd5b611239868387016111df565b9350602085013591508082111561124f57600080fd5b5061125c858286016111f1565b9150509250929050565b60008083601f84011261127857600080fd5b5081356001600160401b0381111561128f57600080fd5b6020830191508360208285010111156112a757600080fd5b9250929050565b600080600080606085870312156112c457600080fd5b84356001600160401b03808211156112db57600080fd5b6112e7888389016111df565b955060208701359150808211156112fd57600080fd5b611309888389016111f1565b9450604087013591508082111561131f57600080fd5b5061132c87828801611266565b95989497509550505050565b60008060006040848603121561134d57600080fd5b83356001600160401b038082111561136457600080fd5b61137087838801611120565b9450602086013591508082111561138657600080fd5b5061139386828701611266565b9497909650939450505050565b6000602082840312156113b257600080fd5b8151801515811461111957600080fd5b6000602082840312156113d457600080fd5b8151611119816110ed565b634e487b7160e01b600052604160045260246000fd5b604051606081016001600160401b0381118282101715611417576114176113df565b60405290565b604051608081016001600160401b0381118282101715611417576114176113df565b604080519081016001600160401b0381118282101715611417576114176113df565b604051601f8201601f191681016001600160401b0381118282101715611489576114896113df565b604052919050565b60006001600160401b038211156114aa576114aa6113df565b50601f01601f191660200190565b6000602082840312156114ca57600080fd5b81516001600160401b038111156114e057600080fd5b8201601f810184136114f157600080fd5b80516115046114ff82611491565b611461565b81815285602083850101111561151957600080fd5b61152a826020830160208601611174565b95945050505050565b60008451611545818460208901611174565b8201838582376000930192835250909392505050565b803563ffffffff8116811461156f57600080fd5b919050565b60006020828403121561158657600080fd5b6111198261155b565b6000602082840312156115a157600080fd5b5051919050565b60008235605e198336030181126115be57600080fd5b9190910192915050565b600082601f8301126115d957600080fd5b81356115e76114ff82611491565b8181528460208386010111156115fc57600080fd5b816020850160208301376000918101602001919091529392505050565b60006060823603121561162b57600080fd5b6116336113f5565b82356001600160401b038082111561164a57600080fd5b81850191506080823603121561165f57600080fd5b61166761141d565b8235815260208301358281111561167d57600080fd5b611689368286016115c8565b6020830152506040830135828111156116a157600080fd5b6116ad368286016115c8565b6040830152506116bf6060840161155b565b606082015283525050602083810135908201526116de6040840161155b565b604082015292915050565b6000808335601e1984360301811261170057600080fd5b8301803591506001600160401b0382111561171a57600080fd5b6020019150368190038213156112a757600080fd5b60008235607e198336030181126115be57600080fd5b600081360360808082121561175957600080fd5b6117616113f5565b60408084121561177057600080fd5b61177861143f565b9350853584526020808701358186015284835261179682880161155b565b8184015260609450848701356001600160401b03808211156117b757600080fd5b9088019036601f8301126117ca57600080fd5b8135818111156117dc576117dc6113df565b6117ea848260051b01611461565b818152848101925060079190911b83018401903682111561180a57600080fd5b928401925b8184101561187c578784360312156118275760008081fd5b61182f61141d565b843561183a816110ed565b815284860135611849816110ed565b818701528487013561185a816110ed565b81880152611869858b0161155b565b818b01528352928701929184019161180f565b948601949094525092979650505050505050565b6000808335601e198436030181126118a757600080fd5b8301803591506001600160401b038211156118c157600080fd5b6020019150600781901b36038213156112a757600080fd5b634e487b7160e01b600052603260045260246000fd5b6020808252605e90820152600080516020611ad583398151915260408201527f72696679426c6f62466f7251756f72756d733a20636f6e6669726d6174696f6e60608201527f5468726573686f6c6450657263656e74616765206973206e6f74206d65740000608082015260a00190565b634e487b7160e01b600052601160045260246000fd5b600060001982141561198a5761198a611960565b5060010190565b602081528151602082015260006020830151608060408401526119b760a08401826111a0565b90506040840151601f198483030160608501526119d482826111a0565b91505063ffffffff60608501511660808401528091505092915050565b6000602080835260a08301845180518386015282810151905060408181870152838701519150606063ffffffff80841682890152828901519350608080818a015285855180885260c08b0191508887019750600096505b80871015611a8b578751805160ff90811684528a82015181168b850152878201511687840152850151841685830152968801966001969096019590820190611a48565b509a9950505050505050505050565b600082611ab757634e487b7160e01b600052601260045260246000fd5b500690565b60008219821115611acf57611acf611960565b50019056fe456967656e4441426c6f62566572696669636174696f6e5574696c732e5f7665a2646970667358221220555735b8d130e400b09a2120fef5d7e4d47a32f9f1661a00757d5c41561959b464736f6c634300080c0033",
}

// ContractEigenDABlobVerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEigenDABlobVerifierMetaData.ABI instead.
var ContractEigenDABlobVerifierABI = ContractEigenDABlobVerifierMetaData.ABI

// ContractEigenDABlobVerifierBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractEigenDABlobVerifierMetaData.Bin instead.
var ContractEigenDABlobVerifierBin = ContractEigenDABlobVerifierMetaData.Bin

// DeployContractEigenDABlobVerifier deploys a new Ethereum contract, binding an instance of ContractEigenDABlobVerifier to it.
func DeployContractEigenDABlobVerifier(auth *bind.TransactOpts, backend bind.ContractBackend, _eigenDAThresholdRegistry common.Address, _eigenDABatchMetadataStorage common.Address, _eigenDASignatureVerifier common.Address) (common.Address, *types.Transaction, *ContractEigenDABlobVerifier, error) {
	parsed, err := ContractEigenDABlobVerifierMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractEigenDABlobVerifierBin), backend, _eigenDAThresholdRegistry, _eigenDABatchMetadataStorage, _eigenDASignatureVerifier)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractEigenDABlobVerifier{ContractEigenDABlobVerifierCaller: ContractEigenDABlobVerifierCaller{contract: contract}, ContractEigenDABlobVerifierTransactor: ContractEigenDABlobVerifierTransactor{contract: contract}, ContractEigenDABlobVerifierFilterer: ContractEigenDABlobVerifierFilterer{contract: contract}}, nil
}

// ContractEigenDABlobVerifier is an auto generated Go binding around an Ethereum contract.
type ContractEigenDABlobVerifier struct {
	ContractEigenDABlobVerifierCaller     // Read-only binding to the contract
	ContractEigenDABlobVerifierTransactor // Write-only binding to the contract
	ContractEigenDABlobVerifierFilterer   // Log filterer for contract events
}

// ContractEigenDABlobVerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractEigenDABlobVerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDABlobVerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractEigenDABlobVerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDABlobVerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractEigenDABlobVerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDABlobVerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractEigenDABlobVerifierSession struct {
	Contract     *ContractEigenDABlobVerifier // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                // Call options to use throughout this session
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ContractEigenDABlobVerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractEigenDABlobVerifierCallerSession struct {
	Contract *ContractEigenDABlobVerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                      // Call options to use throughout this session
}

// ContractEigenDABlobVerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractEigenDABlobVerifierTransactorSession struct {
	Contract     *ContractEigenDABlobVerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                      // Transaction auth options to use throughout this session
}

// ContractEigenDABlobVerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractEigenDABlobVerifierRaw struct {
	Contract *ContractEigenDABlobVerifier // Generic contract binding to access the raw methods on
}

// ContractEigenDABlobVerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractEigenDABlobVerifierCallerRaw struct {
	Contract *ContractEigenDABlobVerifierCaller // Generic read-only contract binding to access the raw methods on
}

// ContractEigenDABlobVerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractEigenDABlobVerifierTransactorRaw struct {
	Contract *ContractEigenDABlobVerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractEigenDABlobVerifier creates a new instance of ContractEigenDABlobVerifier, bound to a specific deployed contract.
func NewContractEigenDABlobVerifier(address common.Address, backend bind.ContractBackend) (*ContractEigenDABlobVerifier, error) {
	contract, err := bindContractEigenDABlobVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDABlobVerifier{ContractEigenDABlobVerifierCaller: ContractEigenDABlobVerifierCaller{contract: contract}, ContractEigenDABlobVerifierTransactor: ContractEigenDABlobVerifierTransactor{contract: contract}, ContractEigenDABlobVerifierFilterer: ContractEigenDABlobVerifierFilterer{contract: contract}}, nil
}

// NewContractEigenDABlobVerifierCaller creates a new read-only instance of ContractEigenDABlobVerifier, bound to a specific deployed contract.
func NewContractEigenDABlobVerifierCaller(address common.Address, caller bind.ContractCaller) (*ContractEigenDABlobVerifierCaller, error) {
	contract, err := bindContractEigenDABlobVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDABlobVerifierCaller{contract: contract}, nil
}

// NewContractEigenDABlobVerifierTransactor creates a new write-only instance of ContractEigenDABlobVerifier, bound to a specific deployed contract.
func NewContractEigenDABlobVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractEigenDABlobVerifierTransactor, error) {
	contract, err := bindContractEigenDABlobVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDABlobVerifierTransactor{contract: contract}, nil
}

// NewContractEigenDABlobVerifierFilterer creates a new log filterer instance of ContractEigenDABlobVerifier, bound to a specific deployed contract.
func NewContractEigenDABlobVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractEigenDABlobVerifierFilterer, error) {
	contract, err := bindContractEigenDABlobVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDABlobVerifierFilterer{contract: contract}, nil
}

// bindContractEigenDABlobVerifier binds a generic wrapper to an already deployed contract.
func bindContractEigenDABlobVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractEigenDABlobVerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDABlobVerifier.Contract.ContractEigenDABlobVerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDABlobVerifier.Contract.ContractEigenDABlobVerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDABlobVerifier.Contract.ContractEigenDABlobVerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDABlobVerifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDABlobVerifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDABlobVerifier.Contract.contract.Transact(opts, method, params...)
}

// EigenDABatchMetadataStorage is a free data retrieval call binding the contract method 0x640f65d9.
//
// Solidity: function eigenDABatchMetadataStorage() view returns(address)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) EigenDABatchMetadataStorage(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "eigenDABatchMetadataStorage")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDABatchMetadataStorage is a free data retrieval call binding the contract method 0x640f65d9.
//
// Solidity: function eigenDABatchMetadataStorage() view returns(address)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) EigenDABatchMetadataStorage() (common.Address, error) {
	return _ContractEigenDABlobVerifier.Contract.EigenDABatchMetadataStorage(&_ContractEigenDABlobVerifier.CallOpts)
}

// EigenDABatchMetadataStorage is a free data retrieval call binding the contract method 0x640f65d9.
//
// Solidity: function eigenDABatchMetadataStorage() view returns(address)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) EigenDABatchMetadataStorage() (common.Address, error) {
	return _ContractEigenDABlobVerifier.Contract.EigenDABatchMetadataStorage(&_ContractEigenDABlobVerifier.CallOpts)
}

// EigenDASignatureVerifier is a free data retrieval call binding the contract method 0xefd4532b.
//
// Solidity: function eigenDASignatureVerifier() view returns(address)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) EigenDASignatureVerifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "eigenDASignatureVerifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDASignatureVerifier is a free data retrieval call binding the contract method 0xefd4532b.
//
// Solidity: function eigenDASignatureVerifier() view returns(address)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) EigenDASignatureVerifier() (common.Address, error) {
	return _ContractEigenDABlobVerifier.Contract.EigenDASignatureVerifier(&_ContractEigenDABlobVerifier.CallOpts)
}

// EigenDASignatureVerifier is a free data retrieval call binding the contract method 0xefd4532b.
//
// Solidity: function eigenDASignatureVerifier() view returns(address)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) EigenDASignatureVerifier() (common.Address, error) {
	return _ContractEigenDABlobVerifier.Contract.EigenDASignatureVerifier(&_ContractEigenDABlobVerifier.CallOpts)
}

// EigenDAThresholdRegistry is a free data retrieval call binding the contract method 0xf8c66814.
//
// Solidity: function eigenDAThresholdRegistry() view returns(address)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) EigenDAThresholdRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "eigenDAThresholdRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDAThresholdRegistry is a free data retrieval call binding the contract method 0xf8c66814.
//
// Solidity: function eigenDAThresholdRegistry() view returns(address)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) EigenDAThresholdRegistry() (common.Address, error) {
	return _ContractEigenDABlobVerifier.Contract.EigenDAThresholdRegistry(&_ContractEigenDABlobVerifier.CallOpts)
}

// EigenDAThresholdRegistry is a free data retrieval call binding the contract method 0xf8c66814.
//
// Solidity: function eigenDAThresholdRegistry() view returns(address)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) EigenDAThresholdRegistry() (common.Address, error) {
	return _ContractEigenDABlobVerifier.Contract.EigenDAThresholdRegistry(&_ContractEigenDABlobVerifier.CallOpts)
}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) GetIsQuorumRequired(opts *bind.CallOpts, quorumNumber uint8) (bool, error) {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "getIsQuorumRequired", quorumNumber)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) GetIsQuorumRequired(quorumNumber uint8) (bool, error) {
	return _ContractEigenDABlobVerifier.Contract.GetIsQuorumRequired(&_ContractEigenDABlobVerifier.CallOpts, quorumNumber)
}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) GetIsQuorumRequired(quorumNumber uint8) (bool, error) {
	return _ContractEigenDABlobVerifier.Contract.GetIsQuorumRequired(&_ContractEigenDABlobVerifier.CallOpts, quorumNumber)
}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) GetQuorumAdversaryThresholdPercentage(opts *bind.CallOpts, quorumNumber uint8) (uint8, error) {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "getQuorumAdversaryThresholdPercentage", quorumNumber)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) GetQuorumAdversaryThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDABlobVerifier.Contract.GetQuorumAdversaryThresholdPercentage(&_ContractEigenDABlobVerifier.CallOpts, quorumNumber)
}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) GetQuorumAdversaryThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDABlobVerifier.Contract.GetQuorumAdversaryThresholdPercentage(&_ContractEigenDABlobVerifier.CallOpts, quorumNumber)
}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) GetQuorumConfirmationThresholdPercentage(opts *bind.CallOpts, quorumNumber uint8) (uint8, error) {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "getQuorumConfirmationThresholdPercentage", quorumNumber)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) GetQuorumConfirmationThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDABlobVerifier.Contract.GetQuorumConfirmationThresholdPercentage(&_ContractEigenDABlobVerifier.CallOpts, quorumNumber)
}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) GetQuorumConfirmationThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDABlobVerifier.Contract.GetQuorumConfirmationThresholdPercentage(&_ContractEigenDABlobVerifier.CallOpts, quorumNumber)
}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) QuorumAdversaryThresholdPercentages(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "quorumAdversaryThresholdPercentages")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) QuorumAdversaryThresholdPercentages() ([]byte, error) {
	return _ContractEigenDABlobVerifier.Contract.QuorumAdversaryThresholdPercentages(&_ContractEigenDABlobVerifier.CallOpts)
}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) QuorumAdversaryThresholdPercentages() ([]byte, error) {
	return _ContractEigenDABlobVerifier.Contract.QuorumAdversaryThresholdPercentages(&_ContractEigenDABlobVerifier.CallOpts)
}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) QuorumConfirmationThresholdPercentages(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "quorumConfirmationThresholdPercentages")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) QuorumConfirmationThresholdPercentages() ([]byte, error) {
	return _ContractEigenDABlobVerifier.Contract.QuorumConfirmationThresholdPercentages(&_ContractEigenDABlobVerifier.CallOpts)
}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) QuorumConfirmationThresholdPercentages() ([]byte, error) {
	return _ContractEigenDABlobVerifier.Contract.QuorumConfirmationThresholdPercentages(&_ContractEigenDABlobVerifier.CallOpts)
}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) QuorumNumbersRequired(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "quorumNumbersRequired")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) QuorumNumbersRequired() ([]byte, error) {
	return _ContractEigenDABlobVerifier.Contract.QuorumNumbersRequired(&_ContractEigenDABlobVerifier.CallOpts)
}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) QuorumNumbersRequired() ([]byte, error) {
	return _ContractEigenDABlobVerifier.Contract.QuorumNumbersRequired(&_ContractEigenDABlobVerifier.CallOpts)
}

// VerifyBlobV1 is a free data retrieval call binding the contract method 0x8d67b909.
//
// Solidity: function verifyBlobV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) VerifyBlobV1(opts *bind.CallOpts, blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof EigenDABlobVerificationUtilsBlobVerificationProof) error {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "verifyBlobV1", blobHeader, blobVerificationProof)

	if err != nil {
		return err
	}

	return err

}

// VerifyBlobV1 is a free data retrieval call binding the contract method 0x8d67b909.
//
// Solidity: function verifyBlobV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) VerifyBlobV1(blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof EigenDABlobVerificationUtilsBlobVerificationProof) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV1(&_ContractEigenDABlobVerifier.CallOpts, blobHeader, blobVerificationProof)
}

// VerifyBlobV1 is a free data retrieval call binding the contract method 0x8d67b909.
//
// Solidity: function verifyBlobV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) VerifyBlobV1(blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof EigenDABlobVerificationUtilsBlobVerificationProof) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV1(&_ContractEigenDABlobVerifier.CallOpts, blobHeader, blobVerificationProof)
}

// VerifyBlobV10 is a free data retrieval call binding the contract method 0x8f3a8f32.
//
// Solidity: function verifyBlobV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof, bytes additionalQuorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) VerifyBlobV10(opts *bind.CallOpts, blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof EigenDABlobVerificationUtilsBlobVerificationProof, additionalQuorumNumbersRequired []byte) error {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "verifyBlobV10", blobHeader, blobVerificationProof, additionalQuorumNumbersRequired)

	if err != nil {
		return err
	}

	return err

}

// VerifyBlobV10 is a free data retrieval call binding the contract method 0x8f3a8f32.
//
// Solidity: function verifyBlobV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof, bytes additionalQuorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) VerifyBlobV10(blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof EigenDABlobVerificationUtilsBlobVerificationProof, additionalQuorumNumbersRequired []byte) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV10(&_ContractEigenDABlobVerifier.CallOpts, blobHeader, blobVerificationProof, additionalQuorumNumbersRequired)
}

// VerifyBlobV10 is a free data retrieval call binding the contract method 0x8f3a8f32.
//
// Solidity: function verifyBlobV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof, bytes additionalQuorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) VerifyBlobV10(blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof EigenDABlobVerificationUtilsBlobVerificationProof, additionalQuorumNumbersRequired []byte) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV10(&_ContractEigenDABlobVerifier.CallOpts, blobHeader, blobVerificationProof, additionalQuorumNumbersRequired)
}

// VerifyBlobV2 is a free data retrieval call binding the contract method 0x7baca37c.
//
// Solidity: function verifyBlobV2(((bytes,((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]),uint32,string[]),(uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][])) signedCertificate) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) VerifyBlobV2(opts *bind.CallOpts, signedCertificate EigenDABlobVerificationUtilsSignedCertificate) error {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "verifyBlobV2", signedCertificate)

	if err != nil {
		return err
	}

	return err

}

// VerifyBlobV2 is a free data retrieval call binding the contract method 0x7baca37c.
//
// Solidity: function verifyBlobV2(((bytes,((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]),uint32,string[]),(uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][])) signedCertificate) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) VerifyBlobV2(signedCertificate EigenDABlobVerificationUtilsSignedCertificate) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV2(&_ContractEigenDABlobVerifier.CallOpts, signedCertificate)
}

// VerifyBlobV2 is a free data retrieval call binding the contract method 0x7baca37c.
//
// Solidity: function verifyBlobV2(((bytes,((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]),uint32,string[]),(uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][])) signedCertificate) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) VerifyBlobV2(signedCertificate EigenDABlobVerificationUtilsSignedCertificate) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV2(&_ContractEigenDABlobVerifier.CallOpts, signedCertificate)
}

// VerifyBlobV20 is a free data retrieval call binding the contract method 0xf78a6ea1.
//
// Solidity: function verifyBlobV2(((bytes,((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]),uint32,string[]),(uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][])) signedCertificate, bytes additionalQuorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) VerifyBlobV20(opts *bind.CallOpts, signedCertificate EigenDABlobVerificationUtilsSignedCertificate, additionalQuorumNumbersRequired []byte) error {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "verifyBlobV20", signedCertificate, additionalQuorumNumbersRequired)

	if err != nil {
		return err
	}

	return err

}

// VerifyBlobV20 is a free data retrieval call binding the contract method 0xf78a6ea1.
//
// Solidity: function verifyBlobV2(((bytes,((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]),uint32,string[]),(uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][])) signedCertificate, bytes additionalQuorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) VerifyBlobV20(signedCertificate EigenDABlobVerificationUtilsSignedCertificate, additionalQuorumNumbersRequired []byte) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV20(&_ContractEigenDABlobVerifier.CallOpts, signedCertificate, additionalQuorumNumbersRequired)
}

// VerifyBlobV20 is a free data retrieval call binding the contract method 0xf78a6ea1.
//
// Solidity: function verifyBlobV2(((bytes,((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]),uint32,string[]),(uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][])) signedCertificate, bytes additionalQuorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) VerifyBlobV20(signedCertificate EigenDABlobVerificationUtilsSignedCertificate, additionalQuorumNumbersRequired []byte) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV20(&_ContractEigenDABlobVerifier.CallOpts, signedCertificate, additionalQuorumNumbersRequired)
}
