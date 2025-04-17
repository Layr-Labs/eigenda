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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_eigenDAThresholdRegistryV1\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"},{\"name\":\"_eigenDABatchMetadataStorageV1\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"eigenDABatchMetadataStorageV1\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistryV1\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV1\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV1ForZkProof\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"success\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertsV1\",\"inputs\":[{\"name\":\"blobHeaders\",\"type\":\"tuple[]\",\"internalType\":\"structBlobHeader[]\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProofs\",\"type\":\"tuple[]\",\"internalType\":\"structBlobVerificationProof[]\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"error\",\"name\":\"BatchMetadataMismatch\",\"inputs\":[{\"name\":\"actualHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expectedHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"BlobQuorumsNotSubset\",\"inputs\":[{\"name\":\"blobQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmedQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ConfirmationThresholdNotMet\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"requiredThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"actualThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"InvalidInclusionProof\",\"inputs\":[{\"name\":\"blobIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blobHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"rootHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InvalidThresholdPercentages\",\"inputs\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"LengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"QuorumNumberMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"actual\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RelayKeyNotSet\",\"inputs\":[{\"name\":\"relayKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RequiredQuorumsNotSubset\",\"inputs\":[{\"name\":\"requiredQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmedQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SecurityAssumptionsNotMet\",\"inputs\":[{\"name\":\"errParams\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"StakeThresholdNotMet\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"requiredThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"actualThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}]",
	Bin: "0x60c06040523480156200001157600080fd5b5060405162001b2438038062001b24833981016040819052620000349162000065565b6001600160a01b039182166080521660a052620000a4565b6001600160a01b03811681146200006257600080fd5b50565b600080604083850312156200007957600080fd5b825162000086816200004c565b602084015190925062000099816200004c565b809150509250929050565b60805160a051611a47620000dd6000396000818160cd01526102ae01526000818160760152818161022101526103520152611a476000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c806331a3479a1461005c5780634cff90c4146100715780637d644cad146100b5578063a9c823e1146100c8578063f88adbba146100ef575b600080fd5b61006f61006a366004611109565b610112565b005b6100987f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b0390911681526020015b60405180910390f35b61006f6100c3366004611174565b6101a2565b6100987f000000000000000000000000000000000000000000000000000000000000000081565b6101026100fd366004611174565b6101dc565b60405190151581526020016100ac565b828114610135576040516001621398b960e31b0319815260040160405180910390fd5b60005b8381101561019b5761018b858583818110610155576101556111e6565b905060200281019061016791906111fc565b848484818110610179576101796111e6565b90506020028101906100c3919061121c565b61019481611248565b9050610138565b5050505050565b6000806101c86101b061021d565b6101b9856102aa565b86866101c361034e565b6103ae565b915091506101d68282610459565b50505050565b6000806101ea6101b061021d565b509050600081600a81111561020157610201611263565b1415610211576001915050610217565b60009150505b92915050565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663bafa91076040518163ffffffff1660e01b8152600401600060405180830381865afa15801561027d573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526102a5919081019061137e565b905090565b60007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663eccbbfc96102e8602085018561141e565b6040516001600160e01b031960e084901b16815263ffffffff919091166004820152602401602060405180830381865afa15801561032a573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102179190611442565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e15234ff6040518163ffffffff1660e01b8152600401600060405180830381865afa15801561027d573d6000803e3d6000fd5b600060606103bc868561083c565b9092509050600082600a8111156103d5576103d5611263565b146103df5761044f565b6103e985856108b6565b9092509050600082600a81111561040257610402611263565b1461040c5761044f565b60006104198887876109f5565b91945092509050600083600a81111561043457610434611263565b1461043f575061044f565b6104498482610a9b565b92509250505b9550959350505050565b600082600a81111561046d5761046d611263565b1415610477575050565b600182600a81111561048b5761048b611263565b14156104d557600080828060200190518101906104a8919061145b565b6040516315d615c960e21b8152600481018390526024810182905291935091506044015b60405180910390fd5b600282600a8111156104e9576104e9611263565b141561053a57600080600083806020019051810190610508919061147f565b60405163d54d727760e01b815260048101849052602481018390526044810182905292955090935091506064016104cc565b600382600a81111561054e5761054e611263565b1415610596576000808280602001905181019061056b91906114bc565b6040516314fa310760e31b815260ff80841660048301528216602482015291935091506044016104cc565b600482600a8111156105aa576105aa611263565b14156105f257600080828060200190518101906105c791906114bc565b604051631b00235d60e01b815260ff80841660048301528216602482015291935091506044016104cc565b600582600a81111561060657610606611263565b141561065b5760008060008380602001905181019061062591906114eb565b604051638aa11c4360e01b815260ff808516600483015280841660248301528216604482015292955090935091506064016104cc565b600682600a81111561066f5761066f611263565b14156106c45760008060008380602001905181019061068e91906114eb565b60405163a4ad875560e01b815260ff808516600483015280841660248301528216604482015292955090935091506064016104cc565b600782600a8111156106d8576106d8611263565b141561071d57600080828060200190518101906106f5919061145b565b60405163114b085b60e21b8152600481018390526024810182905291935091506044016104cc565b600a82600a81111561073157610731611263565b14156107715760008180602001905181019061074d9190611538565b6040516309efaa0b60e41b815263ffffffff821660048201529091506024016104cc565b600882600a81111561078557610785611263565b14156107a65780604051638c59c92f60e01b81526004016104cc9190611581565b600982600a8111156107ba576107ba611263565b14156107ff57600080828060200190518101906107d7919061145b565b604051634a47030360e11b8152600481018390526024810182905291935091506044016104cc565b60405162461bcd60e51b8152602060048201526012602482015271556e6b6e6f776e206572726f7220636f646560701b60448201526064016104cc565b600060608161085f6108516040860186611594565b61085a906115fb565b610aeb565b9050848114156108825750506040805160208101909152600080825291506108af565b60408051602081018390529081018690526001906060015b60405160208183030381529060405292509250505b9250929050565b60006060816108cc6108c7866116d2565b610b5c565b90506000816040516020016108e391815260200190565b604051602081830303815290604052805190602001209050600085806040019061090d9190611594565b61091790806111fc565b359050600061098061092c6060890189611822565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250869250879150610975905060408c0160208d0161141e565b63ffffffff16610b8c565b905080156109a75760006040518060200160405280600081525095509550505050506108af565b60026109b96040890160208a0161141e565b6040805163ffffffff909216602083015281018590526060810184905260800160405160208183030381529060405295509550505050506108af565b600060608180610a0786840187611868565b9050905060005b81811015610a7b576000806000610a288b8b8b8788610ba4565b91945092509050600083600a811115610a4357610a43611263565b14610a5a5750909550935060009250610a92915050565b600160ff82161b861795505050508080610a7390611248565b915050610a0e565b505060408051602081019091526000808252935091505b93509350939050565b600060606000610aaa85610e35565b9050838116811415610acf5750506040805160208101909152600080825291506108af565b604080516020810183905290810185905260079060600161089a565b60006102178260000151604051602001610b0591906118b1565b60408051808303601f1901815282825280516020918201208682015187840151838601929092528484015260e01b6001600160e01b0319166060840152815160448185030181526064909301909152815191012090565b600081604051602001610b6f9190611911565b604051602081830303815290604052805190602001209050919050565b600083610b9a868585610fc2565b1495945050505050565b600060608136610bb688840189611868565b87818110610bc657610bc66111e6565b608002919091019150610bde905060208201826119ba565b91506000610bef6080890189611822565b87818110610bff57610bff6111e6565b919091013560f81c915060009050610c1a60408a018a611594565b610c2490806111fc565b610c32906020810190611822565b8360ff16818110610c4557610c456111e6565b919091013560f81c91505060ff84168114610c92576040805160ff958616602082015291909416818501528351808203850181526060909101909352506003935090915060009050610e2a565b6000610ca460608501604086016119ba565b90506000610cb860408601602087016119ba565b90508060ff168260ff1611610d01576040805160ff93841660208201529190921681830152815180820383018152606090910190915260049650945060009350610e2a92505050565b60008d8760ff1681518110610d1857610d186111e6565b016020015160f81c905060ff8316811115610d75576040805160ff808a1660208301528084169282019290925290841660608201526005906080016040516020818303038152906040526000985098509850505050505050610e2a565b6000610d8460408e018e611594565b610d8e90806111fc565b610d9c906040810190611822565b8760ff16818110610daf57610daf6111e6565b919091013560f81c91505060ff8416811015610e0e576040805160ff808b166020830152808716928201929092529082166060820152600690608001604051602081830303815290604052600099509950995050505050505050610e2a565b5050604080516020810190915260008082529850965050505050505b955095509592505050565b600061010082511115610ebe5760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a4016104cc565b8151610ecc57506000919050565b60008083600081518110610ee257610ee26111e6565b0160200151600160f89190911c81901b92505b8451811015610fb957848181518110610f1057610f106111e6565b0160200151600160f89190911c1b9150828211610fa55760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a4016104cc565b91811791610fb281611248565b9050610ef5565b50909392505050565b600060208451610fd291906119d7565b156110595760405162461bcd60e51b815260206004820152604b60248201527f4d65726b6c652e70726f63657373496e636c7573696f6e50726f6f664b65636360448201527f616b3a2070726f6f66206c656e6774682073686f756c642062652061206d756c60648201526a3a34b836329037b310199960a91b608482015260a4016104cc565b8260205b855181116110bc576110706002856119d7565b611091578160005280860151602052604060002091506002840493506110aa565b8086015160005281602052604060002091506002840493505b6110b56020826119f9565b905061105d565b50949350505050565b60008083601f8401126110d757600080fd5b5081356001600160401b038111156110ee57600080fd5b6020830191508360208260051b85010111156108af57600080fd5b6000806000806040858703121561111f57600080fd5b84356001600160401b038082111561113657600080fd5b611142888389016110c5565b9096509450602087013591508082111561115b57600080fd5b50611168878288016110c5565b95989497509550505050565b6000806040838503121561118757600080fd5b82356001600160401b038082111561119e57600080fd5b90840190608082870312156111b257600080fd5b909250602084013590808211156111c857600080fd5b50830160a081860312156111db57600080fd5b809150509250929050565b634e487b7160e01b600052603260045260246000fd5b60008235607e1983360301811261121257600080fd5b9190910192915050565b60008235609e1983360301811261121257600080fd5b634e487b7160e01b600052601160045260246000fd5b600060001982141561125c5761125c611232565b5060010190565b634e487b7160e01b600052602160045260246000fd5b634e487b7160e01b600052604160045260246000fd5b604051606081016001600160401b03811182821017156112b1576112b1611279565b60405290565b604051608081016001600160401b03811182821017156112b1576112b1611279565b604080519081016001600160401b03811182821017156112b1576112b1611279565b604051601f8201601f191681016001600160401b038111828210171561132357611323611279565b604052919050565b60006001600160401b0382111561134457611344611279565b50601f01601f191660200190565b60005b8381101561136d578181015183820152602001611355565b838111156101d65750506000910152565b60006020828403121561139057600080fd5b81516001600160401b038111156113a657600080fd5b8201601f810184136113b757600080fd5b80516113ca6113c58261132b565b6112fb565b8181528560208385010111156113df57600080fd5b6113f0826020830160208601611352565b95945050505050565b63ffffffff8116811461140b57600080fd5b50565b8035611419816113f9565b919050565b60006020828403121561143057600080fd5b813561143b816113f9565b9392505050565b60006020828403121561145457600080fd5b5051919050565b6000806040838503121561146e57600080fd5b505080516020909101519092909150565b60008060006060848603121561149457600080fd5b8351925060208401519150604084015190509250925092565b60ff8116811461140b57600080fd5b600080604083850312156114cf57600080fd5b82516114da816114ad565b60208401519092506111db816114ad565b60008060006060848603121561150057600080fd5b835161150b816114ad565b602085015190935061151c816114ad565b604085015190925061152d816114ad565b809150509250925092565b60006020828403121561154a57600080fd5b815161143b816113f9565b6000815180845261156d816020860160208601611352565b601f01601f19169290920160200192915050565b60208152600061143b6020830184611555565b60008235605e1983360301811261121257600080fd5b600082601f8301126115bb57600080fd5b81356115c96113c58261132b565b8181528460208386010111156115de57600080fd5b816020850160208301376000918101602001919091529392505050565b60006060823603121561160d57600080fd5b61161561128f565b82356001600160401b038082111561162c57600080fd5b81850191506080823603121561164157600080fd5b6116496112b7565b8235815260208301358281111561165f57600080fd5b61166b368286016115aa565b60208301525060408301358281111561168357600080fd5b61168f368286016115aa565b604083015250606083013592506116a5836113f9565b826060820152808452505050602083013560208201526116c76040840161140e565b604082015292915050565b60008136036080808212156116e657600080fd5b6116ee61128f565b6040808412156116fd57600080fd5b6117056112d9565b8635815260208088013581830152908352818701359450611725856113f9565b848184015260609450848701356001600160401b038082111561174757600080fd5b9088019036601f83011261175a57600080fd5b81358181111561176c5761176c611279565b61177a848260051b016112fb565b818152848101925060079190911b83018401903682111561179a57600080fd5b928401925b8184101561180e578784360312156117b75760008081fd5b6117bf6112b7565b84356117ca816114ad565b8152848601356117d9816114ad565b81870152848701356117ea816114ad565b81880152848a01356117fb816113f9565b818b01528352928701929184019161179f565b948601949094525092979650505050505050565b6000808335601e1984360301811261183957600080fd5b8301803591506001600160401b0382111561185357600080fd5b6020019150368190038213156108af57600080fd5b6000808335601e1984360301811261187f57600080fd5b8301803591506001600160401b0382111561189957600080fd5b6020019150600781901b36038213156108af57600080fd5b602081528151602082015260006020830151608060408401526118d760a0840182611555565b90506040840151601f198483030160608501526118f48282611555565b91505063ffffffff60608501511660808401528091505092915050565b6000602080835260a08301845180518386015282810151905060408181870152838701519150606063ffffffff80841682890152828901519350608080818a015285855180885260c08b0191508887019750600096505b808710156119ab578751805160ff90811684528a82015181168b850152878201511687840152850151841685830152968801966001969096019590820190611968565b509a9950505050505050505050565b6000602082840312156119cc57600080fd5b813561143b816114ad565b6000826119f457634e487b7160e01b600052601260045260246000fd5b500690565b60008219821115611a0c57611a0c611232565b50019056fea2646970667358221220532cef936e5c9941d1d7c22d30bc2c167145a20fb660354a45c2639e6147778464736f6c634300080c0033",
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

// VerifyDACertV1ForZkProof is a free data retrieval call binding the contract method 0xf88adbba.
//
// Solidity: function verifyDACertV1ForZkProof(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns(bool success)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Caller) VerifyDACertV1ForZkProof(opts *bind.CallOpts, blobHeader BlobHeader, blobVerificationProof BlobVerificationProof) (bool, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1.contract.Call(opts, &out, "verifyDACertV1ForZkProof", blobHeader, blobVerificationProof)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyDACertV1ForZkProof is a free data retrieval call binding the contract method 0xf88adbba.
//
// Solidity: function verifyDACertV1ForZkProof(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns(bool success)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Session) VerifyDACertV1ForZkProof(blobHeader BlobHeader, blobVerificationProof BlobVerificationProof) (bool, error) {
	return _ContractEigenDACertVerifierV1.Contract.VerifyDACertV1ForZkProof(&_ContractEigenDACertVerifierV1.CallOpts, blobHeader, blobVerificationProof)
}

// VerifyDACertV1ForZkProof is a free data retrieval call binding the contract method 0xf88adbba.
//
// Solidity: function verifyDACertV1ForZkProof(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns(bool success)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1CallerSession) VerifyDACertV1ForZkProof(blobHeader BlobHeader, blobVerificationProof BlobVerificationProof) (bool, error) {
	return _ContractEigenDACertVerifierV1.Contract.VerifyDACertV1ForZkProof(&_ContractEigenDACertVerifierV1.CallOpts, blobHeader, blobVerificationProof)
}

// VerifyDACertsV1 is a free data retrieval call binding the contract method 0x31a3479a.
//
// Solidity: function verifyDACertsV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[])[] blobHeaders, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes)[] blobVerificationProofs) view returns()
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Caller) VerifyDACertsV1(opts *bind.CallOpts, blobHeaders []BlobHeader, blobVerificationProofs []BlobVerificationProof) error {
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
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Session) VerifyDACertsV1(blobHeaders []BlobHeader, blobVerificationProofs []BlobVerificationProof) error {
	return _ContractEigenDACertVerifierV1.Contract.VerifyDACertsV1(&_ContractEigenDACertVerifierV1.CallOpts, blobHeaders, blobVerificationProofs)
}

// VerifyDACertsV1 is a free data retrieval call binding the contract method 0x31a3479a.
//
// Solidity: function verifyDACertsV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[])[] blobHeaders, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes)[] blobVerificationProofs) view returns()
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1CallerSession) VerifyDACertsV1(blobHeaders []BlobHeader, blobVerificationProofs []BlobVerificationProof) error {
	return _ContractEigenDACertVerifierV1.Contract.VerifyDACertsV1(&_ContractEigenDACertVerifierV1.CallOpts, blobHeaders, blobVerificationProofs)
}
