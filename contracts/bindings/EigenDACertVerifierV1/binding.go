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

// BatchHeader is an auto generated low-level Go binding around an user-defined struct.
type BatchHeader struct {
	BlobHeadersRoot       [32]byte
	QuorumNumbers         []byte
	SignedStakeForQuorums []byte
	ReferenceBlockNumber  uint32
}

// BatchMetadata is an auto generated low-level Go binding around an user-defined struct.
type BatchMetadata struct {
	BatchHeader             BatchHeader
	SignatoryRecordHash     [32]byte
	ConfirmationBlockNumber uint32
}

// BlobHeader is an auto generated low-level Go binding around an user-defined struct.
type BlobHeader struct {
	Commitment       BN254G1Point
	DataLength       uint32
	QuorumBlobParams []QuorumBlobParam
}

// BlobVerificationProof is an auto generated low-level Go binding around an user-defined struct.
type BlobVerificationProof struct {
	BatchId        uint32
	BlobIndex      uint32
	BatchMetadata  BatchMetadata
	InclusionProof []byte
	QuorumIndices  []byte
}

// QuorumBlobParam is an auto generated low-level Go binding around an user-defined struct.
type QuorumBlobParam struct {
	QuorumNumber                    uint8
	AdversaryThresholdPercentage    uint8
	ConfirmationThresholdPercentage uint8
	ChunkLength                     uint32
}

// ContractEigenDACertVerifierV1MetaData contains all meta data concerning the ContractEigenDACertVerifierV1 contract.
var ContractEigenDACertVerifierV1MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_eigenDAThresholdRegistryV1\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"},{\"name\":\"_eigenDABatchMetadataStorageV1\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"checkDACertV1\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"success\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDABatchMetadataStorageV1\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistryV1\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV1\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"error\",\"name\":\"BatchMetadataMismatch\",\"inputs\":[{\"name\":\"actualHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expectedHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"BlobQuorumsNotSubset\",\"inputs\":[{\"name\":\"blobQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmedQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ConfirmationThresholdNotMet\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"requiredThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"actualThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"InvalidInclusionProof\",\"inputs\":[{\"name\":\"blobIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blobHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"rootHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InvalidThresholdPercentages\",\"inputs\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"QuorumNumberMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"actual\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RelayKeyNotSet\",\"inputs\":[{\"name\":\"relayKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RequiredQuorumsNotSubset\",\"inputs\":[{\"name\":\"requiredQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmedQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SecurityAssumptionsNotMet\",\"inputs\":[{\"name\":\"errParams\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"StakeThresholdNotMet\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"requiredThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"actualThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}]",
	Bin: "0x60c06040523480156200001157600080fd5b50604051620019bc380380620019bc833981016040819052620000349162000065565b6001600160a01b039182166080521660a052620000a4565b6001600160a01b03811681146200006257600080fd5b50565b600080604083850312156200007957600080fd5b825162000086816200004c565b602084015190925062000099816200004c565b809150509250929050565b60805160a0516118df620000dd6000396000818160af015261020001526000818160560152818161017301526102a401526118df6000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80634cff90c4146100515780637d644cad14610095578063a9c823e1146100aa578063c084bcbf146100d1575b600080fd5b6100787f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b0390911681526020015b60405180910390f35b6100a86100a3366004611017565b6100f4565b005b6100787f000000000000000000000000000000000000000000000000000000000000000081565b6100e46100df366004611017565b61012e565b604051901515815260200161008c565b60008061011a61010261016f565b61010b856101fc565b86866101156102a0565b610300565b9150915061012882826103ab565b50505050565b60008061013c61010261016f565b509050600081600a8111156101535761015361108a565b1415610163576001915050610169565b60009150505b92915050565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663bafa91076040518163ffffffff1660e01b8152600401600060405180830381865afa1580156101cf573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526101f791908101906111aa565b905090565b60007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663eccbbfc961023a602085018561124b565b6040516001600160e01b031960e084901b16815263ffffffff919091166004820152602401602060405180830381865afa15801561027c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610169919061126f565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e15234ff6040518163ffffffff1660e01b8152600401600060405180830381865afa1580156101cf573d6000803e3d6000fd5b6000606061030e868561078e565b9092509050600082600a8111156103275761032761108a565b14610331576103a1565b61033b8585610808565b9092509050600082600a8111156103545761035461108a565b1461035e576103a1565b600061036b888787610947565b91945092509050600083600a8111156103865761038661108a565b1461039157506103a1565b61039b84826109ed565b92509250505b9550959350505050565b600082600a8111156103bf576103bf61108a565b14156103c9575050565b600182600a8111156103dd576103dd61108a565b141561042757600080828060200190518101906103fa9190611288565b6040516315d615c960e21b8152600481018390526024810182905291935091506044015b60405180910390fd5b600282600a81111561043b5761043b61108a565b141561048c5760008060008380602001905181019061045a91906112ac565b60405163d54d727760e01b8152600481018490526024810183905260448101829052929550909350915060640161041e565b600382600a8111156104a0576104a061108a565b14156104e857600080828060200190518101906104bd91906112e9565b6040516314fa310760e31b815260ff808416600483015282166024820152919350915060440161041e565b600482600a8111156104fc576104fc61108a565b1415610544576000808280602001905181019061051991906112e9565b604051631b00235d60e01b815260ff808416600483015282166024820152919350915060440161041e565b600582600a8111156105585761055861108a565b14156105ad576000806000838060200190518101906105779190611318565b604051638aa11c4360e01b815260ff8085166004830152808416602483015282166044820152929550909350915060640161041e565b600682600a8111156105c1576105c161108a565b1415610616576000806000838060200190518101906105e09190611318565b60405163a4ad875560e01b815260ff8085166004830152808416602483015282166044820152929550909350915060640161041e565b600782600a81111561062a5761062a61108a565b141561066f57600080828060200190518101906106479190611288565b60405163114b085b60e21b81526004810183905260248101829052919350915060440161041e565b600a82600a8111156106835761068361108a565b14156106c35760008180602001905181019061069f9190611365565b6040516309efaa0b60e41b815263ffffffff8216600482015290915060240161041e565b600882600a8111156106d7576106d761108a565b14156106f85780604051638c59c92f60e01b815260040161041e91906113ae565b600982600a81111561070c5761070c61108a565b141561075157600080828060200190518101906107299190611288565b604051634a47030360e11b81526004810183905260248101829052919350915060440161041e565b60405162461bcd60e51b8152602060048201526012602482015271556e6b6e6f776e206572726f7220636f646560701b604482015260640161041e565b60006060816107b16107a360408601866113c1565b6107ac90611432565b610a3d565b9050848114156107d4575050604080516020810190915260008082529150610801565b60408051602081018390529081018690526001906060015b60405160208183030381529060405292509250505b9250929050565b600060608161081e6108198661150a565b610aae565b905060008160405160200161083591815260200190565b604051602081830303815290604052805190602001209050600085806040019061085f91906113c1565b610869908061165b565b35905060006108d261087e6060890189611671565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508692508791506108c7905060408c0160208d0161124b565b63ffffffff16610ade565b905080156108f9576000604051806020016040528060008152509550955050505050610801565b600261090b6040890160208a0161124b565b6040805163ffffffff90921660208301528101859052606081018490526080016040516020818303038152906040529550955050505050610801565b600060608180610959868401876116b8565b9050905060005b818110156109cd57600080600061097a8b8b8b8788610af6565b91945092509050600083600a8111156109955761099561108a565b146109ac57509095509350600092506109e4915050565b600160ff82161b8617955050505080806109c590611718565b915050610960565b505060408051602081019091526000808252935091505b93509350939050565b6000606060006109fc85610d87565b9050838116811415610a21575050604080516020810190915260008082529150610801565b60408051602081018390529081018590526007906060016107ec565b60006101698260000151604051602001610a579190611733565b60408051808303601f1901815282825280516020918201208682015187840151838601929092528484015260e01b6001600160e01b0319166060840152815160448185030181526064909301909152815191012090565b600081604051602001610ac19190611793565b604051602081830303815290604052805190602001209050919050565b600083610aec868585610f14565b1495945050505050565b600060608136610b08888401896116b8565b87818110610b1857610b1861183c565b608002919091019150610b3090506020820182611852565b91506000610b416080890189611671565b87818110610b5157610b5161183c565b919091013560f81c915060009050610b6c60408a018a6113c1565b610b76908061165b565b610b84906020810190611671565b8360ff16818110610b9757610b9761183c565b919091013560f81c91505060ff84168114610be4576040805160ff958616602082015291909416818501528351808203850181526060909101909352506003935090915060009050610d7c565b6000610bf66060850160408601611852565b90506000610c0a6040860160208701611852565b90508060ff168260ff1611610c53576040805160ff93841660208201529190921681830152815180820383018152606090910190915260049650945060009350610d7c92505050565b60008d8760ff1681518110610c6a57610c6a61183c565b016020015160f81c905060ff8316811115610cc7576040805160ff808a1660208301528084169282019290925290841660608201526005906080016040516020818303038152906040526000985098509850505050505050610d7c565b6000610cd660408e018e6113c1565b610ce0908061165b565b610cee906040810190611671565b8760ff16818110610d0157610d0161183c565b919091013560f81c91505060ff8416811015610d60576040805160ff808b166020830152808716928201929092529082166060820152600690608001604051602081830303815290604052600099509950995050505050505050610d7c565b5050604080516020810190915260008082529850965050505050505b955095509592505050565b600061010082511115610e105760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a40161041e565b8151610e1e57506000919050565b60008083600081518110610e3457610e3461183c565b0160200151600160f89190911c81901b92505b8451811015610f0b57848181518110610e6257610e6261183c565b0160200151600160f89190911c1b9150828211610ef75760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a40161041e565b91811791610f0481611718565b9050610e47565b50909392505050565b600060208451610f24919061186f565b15610fab5760405162461bcd60e51b815260206004820152604b60248201527f4d65726b6c652e70726f63657373496e636c7573696f6e50726f6f664b65636360448201527f616b3a2070726f6f66206c656e6774682073686f756c642062652061206d756c60648201526a3a34b836329037b310199960a91b608482015260a40161041e565b8260205b8551811161100e57610fc260028561186f565b610fe357816000528086015160205260406000209150600284049350610ffc565b8086015160005281602052604060002091506002840493505b611007602082611891565b9050610faf565b50949350505050565b6000806040838503121561102a57600080fd5b823567ffffffffffffffff8082111561104257600080fd5b908401906080828703121561105657600080fd5b9092506020840135908082111561106c57600080fd5b50830160a0818603121561107f57600080fd5b809150509250929050565b634e487b7160e01b600052602160045260246000fd5b634e487b7160e01b600052604160045260246000fd5b6040516060810167ffffffffffffffff811182821017156110d9576110d96110a0565b60405290565b6040516080810167ffffffffffffffff811182821017156110d9576110d96110a0565b6040805190810167ffffffffffffffff811182821017156110d9576110d96110a0565b604051601f8201601f1916810167ffffffffffffffff8111828210171561114e5761114e6110a0565b604052919050565b600067ffffffffffffffff821115611170576111706110a0565b50601f01601f191660200190565b60005b83811015611199578181015183820152602001611181565b838111156101285750506000910152565b6000602082840312156111bc57600080fd5b815167ffffffffffffffff8111156111d357600080fd5b8201601f810184136111e457600080fd5b80516111f76111f282611156565b611125565b81815285602083850101111561120c57600080fd5b61121d82602083016020860161117e565b95945050505050565b63ffffffff8116811461123857600080fd5b50565b803561124681611226565b919050565b60006020828403121561125d57600080fd5b813561126881611226565b9392505050565b60006020828403121561128157600080fd5b5051919050565b6000806040838503121561129b57600080fd5b505080516020909101519092909150565b6000806000606084860312156112c157600080fd5b8351925060208401519150604084015190509250925092565b60ff8116811461123857600080fd5b600080604083850312156112fc57600080fd5b8251611307816112da565b602084015190925061107f816112da565b60008060006060848603121561132d57600080fd5b8351611338816112da565b6020850151909350611349816112da565b604085015190925061135a816112da565b809150509250925092565b60006020828403121561137757600080fd5b815161126881611226565b6000815180845261139a81602086016020860161117e565b601f01601f19169290920160200192915050565b6020815260006112686020830184611382565b60008235605e198336030181126113d757600080fd5b9190910192915050565b600082601f8301126113f257600080fd5b81356114006111f282611156565b81815284602083860101111561141557600080fd5b816020850160208301376000918101602001919091529392505050565b60006060823603121561144457600080fd5b61144c6110b6565b823567ffffffffffffffff8082111561146457600080fd5b81850191506080823603121561147957600080fd5b6114816110df565b8235815260208301358281111561149757600080fd5b6114a3368286016113e1565b6020830152506040830135828111156114bb57600080fd5b6114c7368286016113e1565b604083015250606083013592506114dd83611226565b826060820152808452505050602083013560208201526114ff6040840161123b565b604082015292915050565b600081360360808082121561151e57600080fd5b6115266110b6565b60408084121561153557600080fd5b61153d611102565b863581526020808801358183015290835281870135945061155d85611226565b8481840152606094508487013567ffffffffffffffff8082111561158057600080fd5b9088019036601f83011261159357600080fd5b8135818111156115a5576115a56110a0565b6115b3848260051b01611125565b818152848101925060079190911b8301840190368211156115d357600080fd5b928401925b81841015611647578784360312156115f05760008081fd5b6115f86110df565b8435611603816112da565b815284860135611612816112da565b8187015284870135611623816112da565b81880152848a013561163481611226565b818b0152835292870192918401916115d8565b948601949094525092979650505050505050565b60008235607e198336030181126113d757600080fd5b6000808335601e1984360301811261168857600080fd5b83018035915067ffffffffffffffff8211156116a357600080fd5b60200191503681900382131561080157600080fd5b6000808335601e198436030181126116cf57600080fd5b83018035915067ffffffffffffffff8211156116ea57600080fd5b6020019150600781901b360382131561080157600080fd5b634e487b7160e01b600052601160045260246000fd5b600060001982141561172c5761172c611702565b5060010190565b6020815281516020820152600060208301516080604084015261175960a0840182611382565b90506040840151601f198483030160608501526117768282611382565b91505063ffffffff60608501511660808401528091505092915050565b6000602080835260a08301845180518386015282810151905060408181870152838701519150606063ffffffff80841682890152828901519350608080818a015285855180885260c08b0191508887019750600096505b8087101561182d578751805160ff90811684528a82015181168b8501528782015116878401528501518416858301529688019660019690960195908201906117ea565b509a9950505050505050505050565b634e487b7160e01b600052603260045260246000fd5b60006020828403121561186457600080fd5b8135611268816112da565b60008261188c57634e487b7160e01b600052601260045260246000fd5b500690565b600082198211156118a4576118a4611702565b50019056fea2646970667358221220cc39e190efa31888c3602dfa88f9df0f64881fc7db835235b29dee426a436fde64736f6c634300080c0033",
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

// CheckDACertV1 is a free data retrieval call binding the contract method 0xc084bcbf.
//
// Solidity: function checkDACertV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns(bool success)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Caller) CheckDACertV1(opts *bind.CallOpts, blobHeader BlobHeader, blobVerificationProof BlobVerificationProof) (bool, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1.contract.Call(opts, &out, "checkDACertV1", blobHeader, blobVerificationProof)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckDACertV1 is a free data retrieval call binding the contract method 0xc084bcbf.
//
// Solidity: function checkDACertV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns(bool success)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Session) CheckDACertV1(blobHeader BlobHeader, blobVerificationProof BlobVerificationProof) (bool, error) {
	return _ContractEigenDACertVerifierV1.Contract.CheckDACertV1(&_ContractEigenDACertVerifierV1.CallOpts, blobHeader, blobVerificationProof)
}

// CheckDACertV1 is a free data retrieval call binding the contract method 0xc084bcbf.
//
// Solidity: function checkDACertV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns(bool success)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1CallerSession) CheckDACertV1(blobHeader BlobHeader, blobVerificationProof BlobVerificationProof) (bool, error) {
	return _ContractEigenDACertVerifierV1.Contract.CheckDACertV1(&_ContractEigenDACertVerifierV1.CallOpts, blobHeader, blobVerificationProof)
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

// VerifyDACertV1 is a free data retrieval call binding the contract method 0x7d644cad.
//
// Solidity: function verifyDACertV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns()
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Caller) VerifyDACertV1(opts *bind.CallOpts, blobHeader BlobHeader, blobVerificationProof BlobVerificationProof) error {
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
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Session) VerifyDACertV1(blobHeader BlobHeader, blobVerificationProof BlobVerificationProof) error {
	return _ContractEigenDACertVerifierV1.Contract.VerifyDACertV1(&_ContractEigenDACertVerifierV1.CallOpts, blobHeader, blobVerificationProof)
}

// VerifyDACertV1 is a free data retrieval call binding the contract method 0x7d644cad.
//
// Solidity: function verifyDACertV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns()
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1CallerSession) VerifyDACertV1(blobHeader BlobHeader, blobVerificationProof BlobVerificationProof) error {
	return _ContractEigenDACertVerifierV1.Contract.VerifyDACertV1(&_ContractEigenDACertVerifierV1.CallOpts, blobHeader, blobVerificationProof)
}
