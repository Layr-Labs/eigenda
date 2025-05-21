// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractEigenDACertVerifierV1

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

// EigenDATypesV1BatchHeader is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV1BatchHeader struct {
	BlobHeadersRoot       [32]byte
	QuorumNumbers         []byte
	SignedStakeForQuorums []byte
	ReferenceBlockNumber  uint32
}

// EigenDATypesV1BatchMetadata is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV1BatchMetadata struct {
	BatchHeader             EigenDATypesV1BatchHeader
	SignatoryRecordHash     [32]byte
	ConfirmationBlockNumber uint32
}

// EigenDATypesV1BlobHeader is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV1BlobHeader struct {
	Commitment       BN254G1Point
	DataLength       uint32
	QuorumBlobParams []EigenDATypesV1QuorumBlobParam
}

// EigenDATypesV1BlobVerificationProof is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV1BlobVerificationProof struct {
	BatchId        uint32
	BlobIndex      uint32
	BatchMetadata  EigenDATypesV1BatchMetadata
	InclusionProof []byte
	QuorumIndices  []byte
}

// EigenDATypesV1QuorumBlobParam is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV1QuorumBlobParam struct {
	QuorumNumber                    uint8
	AdversaryThresholdPercentage    uint8
	ConfirmationThresholdPercentage uint8
	ChunkLength                     uint32
}

// EigenDATypesV1VersionedBlobParams is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV1VersionedBlobParams struct {
	MaxNumOperators uint32
	NumChunks       uint32
	CodingRate      uint8
}

// ContractEigenDACertVerifierV1MetaData contains all meta data concerning the ContractEigenDACertVerifierV1 contract.
var ContractEigenDACertVerifierV1MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_eigenDAThresholdRegistryV1\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"},{\"name\":\"_eigenDABatchMetadataStorageV1\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"eigenDABatchMetadataStorageV1\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistryV1\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBlobParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.VersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIsQuorumRequired\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumAdversaryThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumConfirmationThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumAdversaryThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumConfirmationThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV1\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.BlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structEigenDATypesV1.QuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.BlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.BatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.BatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertsV1\",\"inputs\":[{\"name\":\"blobHeaders\",\"type\":\"tuple[]\",\"internalType\":\"structEigenDATypesV1.BlobHeader[]\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structEigenDATypesV1.QuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProofs\",\"type\":\"tuple[]\",\"internalType\":\"structEigenDATypesV1.BlobVerificationProof[]\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.BatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.BatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"}]",
	Bin: "0x60c06040523480156200001157600080fd5b50604051620026c4380380620026c4833981016040819052620000349162000065565b6001600160a01b039182166080521660a052620000a4565b6001600160a01b03811681146200006257600080fd5b50565b600080604083850312156200007957600080fd5b825162000086816200004c565b602084015190925062000099816200004c565b809150509250929050565b60805160a0516125ad62000117600039600081816101bf015281816103fb015261045701526000818161015801528181610208015281816102bd01528181610348015281816103da015281816104360152818161048d0152818161051a0152818161057a01526105f901526125ad6000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c80637d644cad116100715780637d644cad146101925780638687feae146101a5578063a9c823e1146101ba578063bafa9107146101e1578063e15234ff146101e9578063ee6c3bcf146101f157600080fd5b8063048886d2146100ae5780631429c7c2146100d65780632ecfe72b146100fb57806331a3479a1461013e5780634cff90c414610153575b600080fd5b6100c16100bc366004611b56565b610204565b60405190151581526020015b60405180910390f35b6100e96100e4366004611b56565b61029a565b60405160ff90911681526020016100cd565b61010e610109366004611b7a565b610329565b60408051825163ffffffff9081168252602080850151909116908201529181015160ff16908201526060016100cd565b61015161014c366004611be9565b6103d5565b005b61017a7f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b0390911681526020016100cd565b6101516101a0366004611c54565b610431565b6101ad610489565b6040516100cd9190611d1e565b61017a7f000000000000000000000000000000000000000000000000000000000000000081565b6101ad610516565b6101ad610576565b6100e96101ff366004611b56565b6105d6565b60007f0000000000000000000000000000000000000000000000000000000000000000604051630244436960e11b815260ff841660048201526001600160a01b03919091169063048886d290602401602060405180830381865afa158015610270573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102949190611d31565b92915050565b604051630a14e3e160e11b815260ff821660048201526000906001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690631429c7c2906024015b602060405180830381865afa158015610305573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102949190611d53565b60408051606081018252600080825260208201819052918101919091527f0000000000000000000000000000000000000000000000000000000000000000604051632ecfe72b60e01b815261ffff841660048201526001600160a01b039190911690632ecfe72b90602401606060405180830381865afa1580156103b1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102949190611e34565b61042b7f00000000000000000000000000000000000000000000000000000000000000007f000000000000000000000000000000000000000000000000000000000000000086868686610426610576565b610628565b50505050565b6104857f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000008484610480610576565b6110f6565b5050565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316638687feae6040518163ffffffff1660e01b8152600401600060405180830381865afa1580156104e9573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526105119190810190611ecc565b905090565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663bafa91076040518163ffffffff1660e01b8152600401600060405180830381865afa1580156104e9573d6000803e3d6000fd5b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e15234ff6040518163ffffffff1660e01b8152600401600060405180830381865afa1580156104e9573d6000803e3d6000fd5b60405163ee6c3bcf60e01b815260ff821660048201526000906001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063ee6c3bcf906024016102e8565b8382146106cc5760405162461bcd60e51b815260206004820152606d602482015260008051602061255883398151915260448201527f7269667944414365727473466f7251756f72756d733a20626c6f62486561646560648201527f727320616e6420626c6f62566572696669636174696f6e50726f6f6673206c6560848201526c0dccee8d040dad2e6dac2e8c6d609b1b60a482015260c4015b60405180910390fd5b6000876001600160a01b031663bafa91076040518163ffffffff1660e01b8152600401600060405180830381865afa15801561070c573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526107349190810190611ecc565b905060005b858110156110eb57876001600160a01b031663eccbbfc986868481811061076257610762611f47565b90506020028101906107749190611f5d565b610782906020810190611f8d565b6040516001600160e01b031960e084901b16815263ffffffff919091166004820152602401602060405180830381865afa1580156107c4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107e89190611faa565b61082b8686848181106107fd576107fd611f47565b905060200281019061080f9190611f5d565b61081d906040810190611fc3565b6108269061202a565b6117fb565b146108be5760405162461bcd60e51b8152602060048201526063602482015260008051602061255883398151915260448201527f7269667944414365727473466f7251756f72756d733a2062617463684d65746160648201527f6461746120646f6573206e6f74206d617463682073746f726564206d6574616460848201526261746160e81b60a482015260c4016106c3565b610a048585838181106108d3576108d3611f47565b90506020028101906108e59190611f5d565b6108f3906060810190612101565b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525089925088915085905081811061093c5761093c611f47565b905060200281019061094e9190611f5d565b61095c906040810190611fc3565b6109669080612147565b3561099c8a8a8681811061097c5761097c611f47565b905060200281019061098e9190612147565b6109979061215d565b61186c565b6040516020016109ae91815260200190565b604051602081830303815290604052805190602001208888868181106109d6576109d6611f47565b90506020028101906109e89190611f5d565b6109f9906040810190602001611f8d565b63ffffffff1661189c565b610a7e5760405162461bcd60e51b8152602060048201526051602482015260008051602061255883398151915260448201527f7269667944414365727473466f7251756f72756d733a20696e636c7573696f6e606482015270081c1c9bdbd9881a5cc81a5b9d985b1a59607a1b608482015260a4016106c3565b6000805b888884818110610a9457610a94611f47565b9050602002810190610aa69190612147565b610ab49060608101906122ad565b905081101561102657888884818110610acf57610acf611f47565b9050602002810190610ae19190612147565b610aef9060608101906122ad565b82818110610aff57610aff611f47565b610b159260206080909202019081019150611b56565b60ff16878785818110610b2a57610b2a611f47565b9050602002810190610b3c9190611f5d565b610b4a906040810190611fc3565b610b549080612147565b610b62906020810190612101565b898987818110610b7457610b74611f47565b9050602002810190610b869190611f5d565b610b94906080810190612101565b85818110610ba457610ba4611f47565b919091013560f81c9050818110610bbd57610bbd611f47565b9050013560f81c60f81b60f81c60ff1614610c495760405162461bcd60e51b8152602060048201526052602482015260008051602061255883398151915260448201527f7269667944414365727473466f7251756f72756d733a2071756f72756d4e756d6064820152710c4cae440c8decae640dcdee840dac2e8c6d60731b608482015260a4016106c3565b888884818110610c5b57610c5b611f47565b9050602002810190610c6d9190612147565b610c7b9060608101906122ad565b82818110610c8b57610c8b611f47565b9050608002016020016020810190610ca39190611b56565b60ff16898985818110610cb857610cb8611f47565b9050602002810190610cca9190612147565b610cd89060608101906122ad565b83818110610ce857610ce8611f47565b9050608002016040016020810190610d009190611b56565b60ff1611610d8a5760405162461bcd60e51b815260206004820152605a602482015260008051602061255883398151915260448201527f7269667944414365727473466f7251756f72756d733a207468726573686f6c6460648201527f2070657263656e746167657320617265206e6f742076616c6964000000000000608482015260a4016106c3565b83898985818110610d9d57610d9d611f47565b9050602002810190610daf9190612147565b610dbd9060608101906122ad565b83818110610dcd57610dcd611f47565b610de39260206080909202019081019150611b56565b60ff1681518110610df657610df6611f47565b016020015160f81c898985818110610e1057610e10611f47565b9050602002810190610e229190612147565b610e309060608101906122ad565b83818110610e4057610e40611f47565b9050608002016040016020810190610e589190611b56565b60ff161015610e795760405162461bcd60e51b81526004016106c3906122f6565b888884818110610e8b57610e8b611f47565b9050602002810190610e9d9190612147565b610eab9060608101906122ad565b82818110610ebb57610ebb611f47565b9050608002016040016020810190610ed39190611b56565b60ff16878785818110610ee857610ee8611f47565b9050602002810190610efa9190611f5d565b610f08906040810190611fc3565b610f129080612147565b610f20906040810190612101565b898987818110610f3257610f32611f47565b9050602002810190610f449190611f5d565b610f52906080810190612101565b85818110610f6257610f62611f47565b919091013560f81c9050818110610f7b57610f7b611f47565b9050013560f81c60f81b60f81c60ff161015610fa95760405162461bcd60e51b81526004016106c3906122f6565b611012828a8a86818110610fbf57610fbf611f47565b9050602002810190610fd19190612147565b610fdf9060608101906122ad565b84818110610fef57610fef611f47565b6110059260206080909202019081019150611b56565b600160ff919091161b1790565b91508061101e81612387565b915050610a82565b5061103a611033856118b4565b8281161490565b6110da5760405162461bcd60e51b8152602060048201526071602482015260008051602061255883398151915260448201527f7269667944414365727473466f7251756f72756d733a2072657175697265642060648201527f71756f72756d7320617265206e6f74206120737562736574206f662074686520608482015270636f6e6669726d65642071756f72756d7360781b60a482015260c4016106c3565b506110e481612387565b9050610739565b505050505050505050565b6001600160a01b03841663eccbbfc96111126020850185611f8d565b6040516001600160e01b031960e084901b16815263ffffffff919091166004820152602401602060405180830381865afa158015611154573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111789190611faa565b61118861081d6040850185611fc3565b1461121a5760405162461bcd60e51b8152602060048201526062602482015260008051602061255883398151915260448201527f72696679444143657274466f7251756f72756d733a2062617463684d6574616460648201527f61746120646f6573206e6f74206d617463682073746f726564206d6574616461608482015261746160f01b60a482015260c4016106c3565b6112be61122a6060840184612101565b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061126c925050506040850185611fc3565b6112769080612147565b356112836109978761215d565b60405160200161129591815260200190565b604051602081830303815290604052805190602001208560200160208101906109f99190611f8d565b6113375760405162461bcd60e51b8152602060048201526050602482015260008051602061255883398151915260448201527f72696679444143657274466f7251756f72756d733a20696e636c7573696f6e2060648201526f1c1c9bdbd9881a5cc81a5b9d985b1a5960821b608482015260a4016106c3565b6000805b61134860608601866122ad565b90508110156117475761135e60608601866122ad565b8281811061136e5761136e611f47565b6113849260206080909202019081019150611b56565b60ff166113946040860186611fc3565b61139e9080612147565b6113ac906020810190612101565b6113b96080880188612101565b858181106113c9576113c9611f47565b919091013560f81c90508181106113e2576113e2611f47565b9050013560f81c60f81b60f81c60ff161461146d5760405162461bcd60e51b8152602060048201526051602482015260008051602061255883398151915260448201527f72696679444143657274466f7251756f72756d733a2071756f72756d4e756d626064820152700cae440c8decae640dcdee840dac2e8c6d607b1b608482015260a4016106c3565b61147a60608601866122ad565b8281811061148a5761148a611f47565b90506080020160200160208101906114a29190611b56565b60ff166114b260608701876122ad565b838181106114c2576114c2611f47565b90506080020160400160208101906114da9190611b56565b60ff16116115645760405162461bcd60e51b8152602060048201526059602482015260008051602061255883398151915260448201527f72696679444143657274466f7251756f72756d733a207468726573686f6c642060648201527f70657263656e746167657320617265206e6f742076616c696400000000000000608482015260a4016106c3565b6001600160a01b038716631429c7c261158060608801886122ad565b8481811061159057611590611f47565b6115a69260206080909202019081019150611b56565b6040516001600160e01b031960e084901b16815260ff9091166004820152602401602060405180830381865afa1580156115e4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906116089190611d53565b60ff1661161860608701876122ad565b8381811061162857611628611f47565b90506080020160400160208101906116409190611b56565b60ff1610156116615760405162461bcd60e51b81526004016106c3906123a2565b61166e60608601866122ad565b8281811061167e5761167e611f47565b90506080020160400160208101906116969190611b56565b60ff166116a66040860186611fc3565b6116b09080612147565b6116be906040810190612101565b6116cb6080880188612101565b858181106116db576116db611f47565b919091013560f81c90508181106116f4576116f4611f47565b9050013560f81c60f81b60f81c60ff1610156117225760405162461bcd60e51b81526004016106c3906123a2565b61173382610fdf60608801886122ad565b91508061173f81612387565b91505061133b565b50611754611033836118b4565b6117f35760405162461bcd60e51b8152602060048201526070602482015260008051602061255883398151915260448201527f72696679444143657274466f7251756f72756d733a207265717569726564207160648201527f756f72756d7320617265206e6f74206120737562736574206f6620746865206360848201526f6f6e6669726d65642071756f72756d7360801b60a482015260c4016106c3565b505050505050565b600061029482600001516040516020016118159190612414565b60408051808303601f1901815282825280516020918201208682015187840151838601929092528484015260e01b6001600160e01b0319166060840152815160448185030181526064909301909152815191012090565b60008160405160200161187f9190612474565b604051602081830303815290604052805190602001209050919050565b6000836118aa868585611a41565b1495945050505050565b60006101008251111561193d5760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a4016106c3565b815161194b57506000919050565b6000808360008151811061196157611961611f47565b0160200151600160f89190911c81901b92505b8451811015611a385784818151811061198f5761198f611f47565b0160200151600160f89190911c1b9150828211611a245760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a4016106c3565b91811791611a3181612387565b9050611974565b50909392505050565b600060208451611a51919061251d565b15611ad85760405162461bcd60e51b815260206004820152604b60248201527f4d65726b6c652e70726f63657373496e636c7573696f6e50726f6f664b65636360448201527f616b3a2070726f6f66206c656e6774682073686f756c642062652061206d756c60648201526a3a34b836329037b310199960a91b608482015260a4016106c3565b8260205b85518111611b3b57611aef60028561251d565b611b1057816000528086015160205260406000209150600284049350611b29565b8086015160005281602052604060002091506002840493505b611b3460208261253f565b9050611adc565b50949350505050565b60ff81168114611b5357600080fd5b50565b600060208284031215611b6857600080fd5b8135611b7381611b44565b9392505050565b600060208284031215611b8c57600080fd5b813561ffff81168114611b7357600080fd5b60008083601f840112611bb057600080fd5b5081356001600160401b03811115611bc757600080fd5b6020830191508360208260051b8501011115611be257600080fd5b9250929050565b60008060008060408587031215611bff57600080fd5b84356001600160401b0380821115611c1657600080fd5b611c2288838901611b9e565b90965094506020870135915080821115611c3b57600080fd5b50611c4887828801611b9e565b95989497509550505050565b60008060408385031215611c6757600080fd5b82356001600160401b0380821115611c7e57600080fd5b9084019060808287031215611c9257600080fd5b90925060208401359080821115611ca857600080fd5b50830160a08186031215611cbb57600080fd5b809150509250929050565b60005b83811015611ce1578181015183820152602001611cc9565b8381111561042b5750506000910152565b60008151808452611d0a816020860160208601611cc6565b601f01601f19169290920160200192915050565b602081526000611b736020830184611cf2565b600060208284031215611d4357600080fd5b81518015158114611b7357600080fd5b600060208284031215611d6557600080fd5b8151611b7381611b44565b634e487b7160e01b600052604160045260246000fd5b604051606081016001600160401b0381118282101715611da857611da8611d70565b60405290565b604051608081016001600160401b0381118282101715611da857611da8611d70565b604080519081016001600160401b0381118282101715611da857611da8611d70565b604051601f8201601f191681016001600160401b0381118282101715611e1a57611e1a611d70565b604052919050565b63ffffffff81168114611b5357600080fd5b600060608284031215611e4657600080fd5b604051606081018181106001600160401b0382111715611e6857611e68611d70565b6040528251611e7681611e22565b81526020830151611e8681611e22565b60208201526040830151611e9981611b44565b60408201529392505050565b60006001600160401b03821115611ebe57611ebe611d70565b50601f01601f191660200190565b600060208284031215611ede57600080fd5b81516001600160401b03811115611ef457600080fd5b8201601f81018413611f0557600080fd5b8051611f18611f1382611ea5565b611df2565b818152856020838501011115611f2d57600080fd5b611f3e826020830160208601611cc6565b95945050505050565b634e487b7160e01b600052603260045260246000fd5b60008235609e19833603018112611f7357600080fd5b9190910192915050565b8035611f8881611e22565b919050565b600060208284031215611f9f57600080fd5b8135611b7381611e22565b600060208284031215611fbc57600080fd5b5051919050565b60008235605e19833603018112611f7357600080fd5b600082601f830112611fea57600080fd5b8135611ff8611f1382611ea5565b81815284602083860101111561200d57600080fd5b816020850160208301376000918101602001919091529392505050565b60006060823603121561203c57600080fd5b612044611d86565b82356001600160401b038082111561205b57600080fd5b81850191506080823603121561207057600080fd5b612078611dae565b8235815260208301358281111561208e57600080fd5b61209a36828601611fd9565b6020830152506040830135828111156120b257600080fd5b6120be36828601611fd9565b604083015250606083013592506120d483611e22565b826060820152808452505050602083013560208201526120f660408401611f7d565b604082015292915050565b6000808335601e1984360301811261211857600080fd5b8301803591506001600160401b0382111561213257600080fd5b602001915036819003821315611be257600080fd5b60008235607e19833603018112611f7357600080fd5b600081360360808082121561217157600080fd5b612179611d86565b60408084121561218857600080fd5b612190611dd0565b86358152602080880135818301529083528187013594506121b085611e22565b848184015260609450848701356001600160401b03808211156121d257600080fd5b9088019036601f8301126121e557600080fd5b8135818111156121f7576121f7611d70565b612205848260051b01611df2565b818152848101925060079190911b83018401903682111561222557600080fd5b928401925b81841015612299578784360312156122425760008081fd5b61224a611dae565b843561225581611b44565b81528486013561226481611b44565b818701528487013561227581611b44565b81880152848a013561228681611e22565b818b01528352928701929184019161222a565b948601949094525092979650505050505050565b6000808335601e198436030181126122c457600080fd5b8301803591506001600160401b038211156122de57600080fd5b6020019150600781901b3603821315611be257600080fd5b602080825260619082015260008051602061255883398151915260408201527f7269667944414365727473466f7251756f72756d733a20636f6e6669726d617460608201527f696f6e5468726573686f6c6450657263656e74616765206973206e6f74206d656080820152601d60fa1b60a082015260c00190565b634e487b7160e01b600052601160045260246000fd5b600060001982141561239b5761239b612371565b5060010190565b6020808252606090820181905260008051602061255883398151915260408301527f72696679444143657274466f7251756f72756d733a20636f6e6669726d617469908201527f6f6e5468726573686f6c6450657263656e74616765206973206e6f74206d6574608082015260a00190565b6020815281516020820152600060208301516080604084015261243a60a0840182611cf2565b90506040840151601f198483030160608501526124578282611cf2565b91505063ffffffff60608501511660808401528091505092915050565b6000602080835260a08301845180518386015282810151905060408181870152838701519150606063ffffffff80841682890152828901519350608080818a015285855180885260c08b0191508887019750600096505b8087101561250e578751805160ff90811684528a82015181168b8501528782015116878401528501518416858301529688019660019690960195908201906124cb565b509a9950505050505050505050565b60008261253a57634e487b7160e01b600052601260045260246000fd5b500690565b6000821982111561255257612552612371565b50019056fe456967656e444143657274566572696669636174696f6e56314c69622e5f7665a2646970667358221220712f40f1859529c6596409909e2cf36f820e9efad78cc1eb2fc9c077ebee943564736f6c634300080c0033",
}

// ContractEigenDACertVerifierV1ABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEigenDACertVerifierV1MetaData.ABI instead.
var ContractEigenDACertVerifierV1ABI = ContractEigenDACertVerifierV1MetaData.ABI

// ContractEigenDACertVerifierV1Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractEigenDACertVerifierV1MetaData.Bin instead.
var ContractEigenDACertVerifierV1Bin = ContractEigenDACertVerifierV1MetaData.Bin

// DeployContractEigenDACertVerifierV1 deploys a new Ethereum contract, binding an instance of ContractEigenDACertVerifierV1 to it.
func DeployContractEigenDACertVerifierV1(auth *bind.TransactOpts, backend bind.ContractBackend, _eigenDAThresholdRegistryV1 common.Address, _eigenDABatchMetadataStorageV1 common.Address) (common.Address, *types.Transaction, *ContractEigenDACertVerifierV1, error) {
	parsed, err := ContractEigenDACertVerifierV1MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractEigenDACertVerifierV1Bin), backend, _eigenDAThresholdRegistryV1, _eigenDABatchMetadataStorageV1)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractEigenDACertVerifierV1{ContractEigenDACertVerifierV1Caller: ContractEigenDACertVerifierV1Caller{contract: contract}, ContractEigenDACertVerifierV1Transactor: ContractEigenDACertVerifierV1Transactor{contract: contract}, ContractEigenDACertVerifierV1Filterer: ContractEigenDACertVerifierV1Filterer{contract: contract}}, nil
}

// ContractEigenDACertVerifierV1 is an auto generated Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV1 struct {
	ContractEigenDACertVerifierV1Caller     // Read-only binding to the contract
	ContractEigenDACertVerifierV1Transactor // Write-only binding to the contract
	ContractEigenDACertVerifierV1Filterer   // Log filterer for contract events
}

// ContractEigenDACertVerifierV1Caller is an auto generated read-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV1Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierV1Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV1Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierV1Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractEigenDACertVerifierV1Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierV1Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractEigenDACertVerifierV1Session struct {
	Contract     *ContractEigenDACertVerifierV1 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                  // Call options to use throughout this session
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// ContractEigenDACertVerifierV1CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractEigenDACertVerifierV1CallerSession struct {
	Contract *ContractEigenDACertVerifierV1Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                        // Call options to use throughout this session
}

// ContractEigenDACertVerifierV1TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractEigenDACertVerifierV1TransactorSession struct {
	Contract     *ContractEigenDACertVerifierV1Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                        // Transaction auth options to use throughout this session
}

// ContractEigenDACertVerifierV1Raw is an auto generated low-level Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV1Raw struct {
	Contract *ContractEigenDACertVerifierV1 // Generic contract binding to access the raw methods on
}

// ContractEigenDACertVerifierV1CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV1CallerRaw struct {
	Contract *ContractEigenDACertVerifierV1Caller // Generic read-only contract binding to access the raw methods on
}

// ContractEigenDACertVerifierV1TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV1TransactorRaw struct {
	Contract *ContractEigenDACertVerifierV1Transactor // Generic write-only contract binding to access the raw methods on
}

// NewContractEigenDACertVerifierV1 creates a new instance of ContractEigenDACertVerifierV1, bound to a specific deployed contract.
func NewContractEigenDACertVerifierV1(address common.Address, backend bind.ContractBackend) (*ContractEigenDACertVerifierV1, error) {
	contract, err := bindContractEigenDACertVerifierV1(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierV1{ContractEigenDACertVerifierV1Caller: ContractEigenDACertVerifierV1Caller{contract: contract}, ContractEigenDACertVerifierV1Transactor: ContractEigenDACertVerifierV1Transactor{contract: contract}, ContractEigenDACertVerifierV1Filterer: ContractEigenDACertVerifierV1Filterer{contract: contract}}, nil
}

// NewContractEigenDACertVerifierV1Caller creates a new read-only instance of ContractEigenDACertVerifierV1, bound to a specific deployed contract.
func NewContractEigenDACertVerifierV1Caller(address common.Address, caller bind.ContractCaller) (*ContractEigenDACertVerifierV1Caller, error) {
	contract, err := bindContractEigenDACertVerifierV1(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierV1Caller{contract: contract}, nil
}

// NewContractEigenDACertVerifierV1Transactor creates a new write-only instance of ContractEigenDACertVerifierV1, bound to a specific deployed contract.
func NewContractEigenDACertVerifierV1Transactor(address common.Address, transactor bind.ContractTransactor) (*ContractEigenDACertVerifierV1Transactor, error) {
	contract, err := bindContractEigenDACertVerifierV1(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierV1Transactor{contract: contract}, nil
}

// NewContractEigenDACertVerifierV1Filterer creates a new log filterer instance of ContractEigenDACertVerifierV1, bound to a specific deployed contract.
func NewContractEigenDACertVerifierV1Filterer(address common.Address, filterer bind.ContractFilterer) (*ContractEigenDACertVerifierV1Filterer, error) {
	contract, err := bindContractEigenDACertVerifierV1(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierV1Filterer{contract: contract}, nil
}

// bindContractEigenDACertVerifierV1 binds a generic wrapper to an already deployed contract.
func bindContractEigenDACertVerifierV1(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractEigenDACertVerifierV1MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDACertVerifierV1.Contract.ContractEigenDACertVerifierV1Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierV1.Contract.ContractEigenDACertVerifierV1Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierV1.Contract.ContractEigenDACertVerifierV1Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDACertVerifierV1.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierV1.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierV1.Contract.contract.Transact(opts, method, params...)
}

// EigenDABatchMetadataStorageV1 is a free data retrieval call binding the contract method 0xa9c823e1.
//
// Solidity: function eigenDABatchMetadataStorageV1() view returns(address)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Caller) EigenDABatchMetadataStorageV1(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1.contract.Call(opts, &out, "eigenDABatchMetadataStorageV1")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDABatchMetadataStorageV1 is a free data retrieval call binding the contract method 0xa9c823e1.
//
// Solidity: function eigenDABatchMetadataStorageV1() view returns(address)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Session) EigenDABatchMetadataStorageV1() (common.Address, error) {
	return _ContractEigenDACertVerifierV1.Contract.EigenDABatchMetadataStorageV1(&_ContractEigenDACertVerifierV1.CallOpts)
}

// EigenDABatchMetadataStorageV1 is a free data retrieval call binding the contract method 0xa9c823e1.
//
// Solidity: function eigenDABatchMetadataStorageV1() view returns(address)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1CallerSession) EigenDABatchMetadataStorageV1() (common.Address, error) {
	return _ContractEigenDACertVerifierV1.Contract.EigenDABatchMetadataStorageV1(&_ContractEigenDACertVerifierV1.CallOpts)
}

// EigenDAThresholdRegistryV1 is a free data retrieval call binding the contract method 0x4cff90c4.
//
// Solidity: function eigenDAThresholdRegistryV1() view returns(address)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Caller) EigenDAThresholdRegistryV1(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1.contract.Call(opts, &out, "eigenDAThresholdRegistryV1")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDAThresholdRegistryV1 is a free data retrieval call binding the contract method 0x4cff90c4.
//
// Solidity: function eigenDAThresholdRegistryV1() view returns(address)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Session) EigenDAThresholdRegistryV1() (common.Address, error) {
	return _ContractEigenDACertVerifierV1.Contract.EigenDAThresholdRegistryV1(&_ContractEigenDACertVerifierV1.CallOpts)
}

// EigenDAThresholdRegistryV1 is a free data retrieval call binding the contract method 0x4cff90c4.
//
// Solidity: function eigenDAThresholdRegistryV1() view returns(address)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1CallerSession) EigenDAThresholdRegistryV1() (common.Address, error) {
	return _ContractEigenDACertVerifierV1.Contract.EigenDAThresholdRegistryV1(&_ContractEigenDACertVerifierV1.CallOpts)
}

// GetBlobParams is a free data retrieval call binding the contract method 0x2ecfe72b.
//
// Solidity: function getBlobParams(uint16 version) view returns((uint32,uint32,uint8))
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Caller) GetBlobParams(opts *bind.CallOpts, version uint16) (EigenDATypesV1VersionedBlobParams, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1.contract.Call(opts, &out, "getBlobParams", version)

	if err != nil {
		return *new(EigenDATypesV1VersionedBlobParams), err
	}

	out0 := *abi.ConvertType(out[0], new(EigenDATypesV1VersionedBlobParams)).(*EigenDATypesV1VersionedBlobParams)

	return out0, err

}

// GetBlobParams is a free data retrieval call binding the contract method 0x2ecfe72b.
//
// Solidity: function getBlobParams(uint16 version) view returns((uint32,uint32,uint8))
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Session) GetBlobParams(version uint16) (EigenDATypesV1VersionedBlobParams, error) {
	return _ContractEigenDACertVerifierV1.Contract.GetBlobParams(&_ContractEigenDACertVerifierV1.CallOpts, version)
}

// GetBlobParams is a free data retrieval call binding the contract method 0x2ecfe72b.
//
// Solidity: function getBlobParams(uint16 version) view returns((uint32,uint32,uint8))
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1CallerSession) GetBlobParams(version uint16) (EigenDATypesV1VersionedBlobParams, error) {
	return _ContractEigenDACertVerifierV1.Contract.GetBlobParams(&_ContractEigenDACertVerifierV1.CallOpts, version)
}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Caller) GetIsQuorumRequired(opts *bind.CallOpts, quorumNumber uint8) (bool, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1.contract.Call(opts, &out, "getIsQuorumRequired", quorumNumber)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Session) GetIsQuorumRequired(quorumNumber uint8) (bool, error) {
	return _ContractEigenDACertVerifierV1.Contract.GetIsQuorumRequired(&_ContractEigenDACertVerifierV1.CallOpts, quorumNumber)
}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1CallerSession) GetIsQuorumRequired(quorumNumber uint8) (bool, error) {
	return _ContractEigenDACertVerifierV1.Contract.GetIsQuorumRequired(&_ContractEigenDACertVerifierV1.CallOpts, quorumNumber)
}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Caller) GetQuorumAdversaryThresholdPercentage(opts *bind.CallOpts, quorumNumber uint8) (uint8, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1.contract.Call(opts, &out, "getQuorumAdversaryThresholdPercentage", quorumNumber)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Session) GetQuorumAdversaryThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDACertVerifierV1.Contract.GetQuorumAdversaryThresholdPercentage(&_ContractEigenDACertVerifierV1.CallOpts, quorumNumber)
}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1CallerSession) GetQuorumAdversaryThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDACertVerifierV1.Contract.GetQuorumAdversaryThresholdPercentage(&_ContractEigenDACertVerifierV1.CallOpts, quorumNumber)
}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Caller) GetQuorumConfirmationThresholdPercentage(opts *bind.CallOpts, quorumNumber uint8) (uint8, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1.contract.Call(opts, &out, "getQuorumConfirmationThresholdPercentage", quorumNumber)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Session) GetQuorumConfirmationThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDACertVerifierV1.Contract.GetQuorumConfirmationThresholdPercentage(&_ContractEigenDACertVerifierV1.CallOpts, quorumNumber)
}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1CallerSession) GetQuorumConfirmationThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDACertVerifierV1.Contract.GetQuorumConfirmationThresholdPercentage(&_ContractEigenDACertVerifierV1.CallOpts, quorumNumber)
}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Caller) QuorumAdversaryThresholdPercentages(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1.contract.Call(opts, &out, "quorumAdversaryThresholdPercentages")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Session) QuorumAdversaryThresholdPercentages() ([]byte, error) {
	return _ContractEigenDACertVerifierV1.Contract.QuorumAdversaryThresholdPercentages(&_ContractEigenDACertVerifierV1.CallOpts)
}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1CallerSession) QuorumAdversaryThresholdPercentages() ([]byte, error) {
	return _ContractEigenDACertVerifierV1.Contract.QuorumAdversaryThresholdPercentages(&_ContractEigenDACertVerifierV1.CallOpts)
}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Caller) QuorumConfirmationThresholdPercentages(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1.contract.Call(opts, &out, "quorumConfirmationThresholdPercentages")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Session) QuorumConfirmationThresholdPercentages() ([]byte, error) {
	return _ContractEigenDACertVerifierV1.Contract.QuorumConfirmationThresholdPercentages(&_ContractEigenDACertVerifierV1.CallOpts)
}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1CallerSession) QuorumConfirmationThresholdPercentages() ([]byte, error) {
	return _ContractEigenDACertVerifierV1.Contract.QuorumConfirmationThresholdPercentages(&_ContractEigenDACertVerifierV1.CallOpts)
}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Caller) QuorumNumbersRequired(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1.contract.Call(opts, &out, "quorumNumbersRequired")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Session) QuorumNumbersRequired() ([]byte, error) {
	return _ContractEigenDACertVerifierV1.Contract.QuorumNumbersRequired(&_ContractEigenDACertVerifierV1.CallOpts)
}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1CallerSession) QuorumNumbersRequired() ([]byte, error) {
	return _ContractEigenDACertVerifierV1.Contract.QuorumNumbersRequired(&_ContractEigenDACertVerifierV1.CallOpts)
}

// VerifyDACertV1 is a free data retrieval call binding the contract method 0x7d644cad.
//
// Solidity: function verifyDACertV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns()
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Caller) VerifyDACertV1(opts *bind.CallOpts, blobHeader EigenDATypesV1BlobHeader, blobVerificationProof EigenDATypesV1BlobVerificationProof) error {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1.contract.Call(opts, &out, "verifyDACertV1", blobHeader, blobVerificationProof)

	if err != nil {
		return err
	}

	return err

}

// VerifyDACertV1 is a free data retrieval call binding the contract method 0x7d644cad.
//
// Solidity: function verifyDACertV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns()
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Session) VerifyDACertV1(blobHeader EigenDATypesV1BlobHeader, blobVerificationProof EigenDATypesV1BlobVerificationProof) error {
	return _ContractEigenDACertVerifierV1.Contract.VerifyDACertV1(&_ContractEigenDACertVerifierV1.CallOpts, blobHeader, blobVerificationProof)
}

// VerifyDACertV1 is a free data retrieval call binding the contract method 0x7d644cad.
//
// Solidity: function verifyDACertV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns()
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1CallerSession) VerifyDACertV1(blobHeader EigenDATypesV1BlobHeader, blobVerificationProof EigenDATypesV1BlobVerificationProof) error {
	return _ContractEigenDACertVerifierV1.Contract.VerifyDACertV1(&_ContractEigenDACertVerifierV1.CallOpts, blobHeader, blobVerificationProof)
}

// VerifyDACertsV1 is a free data retrieval call binding the contract method 0x31a3479a.
//
// Solidity: function verifyDACertsV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[])[] blobHeaders, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes)[] blobVerificationProofs) view returns()
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Caller) VerifyDACertsV1(opts *bind.CallOpts, blobHeaders []EigenDATypesV1BlobHeader, blobVerificationProofs []EigenDATypesV1BlobVerificationProof) error {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1.contract.Call(opts, &out, "verifyDACertsV1", blobHeaders, blobVerificationProofs)

	if err != nil {
		return err
	}

	return err

}

// VerifyDACertsV1 is a free data retrieval call binding the contract method 0x31a3479a.
//
// Solidity: function verifyDACertsV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[])[] blobHeaders, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes)[] blobVerificationProofs) view returns()
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Session) VerifyDACertsV1(blobHeaders []EigenDATypesV1BlobHeader, blobVerificationProofs []EigenDATypesV1BlobVerificationProof) error {
	return _ContractEigenDACertVerifierV1.Contract.VerifyDACertsV1(&_ContractEigenDACertVerifierV1.CallOpts, blobHeaders, blobVerificationProofs)
}

// VerifyDACertsV1 is a free data retrieval call binding the contract method 0x31a3479a.
//
// Solidity: function verifyDACertsV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[])[] blobHeaders, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes)[] blobVerificationProofs) view returns()
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1CallerSession) VerifyDACertsV1(blobHeaders []EigenDATypesV1BlobHeader, blobVerificationProofs []EigenDATypesV1BlobVerificationProof) error {
	return _ContractEigenDACertVerifierV1.Contract.VerifyDACertsV1(&_ContractEigenDACertVerifierV1.CallOpts, blobHeaders, blobVerificationProofs)
}
