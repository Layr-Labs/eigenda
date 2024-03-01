// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractMockRollup

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

// EigenDARollupUtilsBlobVerificationProof is an auto generated low-level Go binding around an user-defined struct.
type EigenDARollupUtilsBlobVerificationProof struct {
	BatchId                uint32
	BlobIndex              uint8
	BatchMetadata          IEigenDAServiceManagerBatchMetadata
	InclusionProof         []byte
	QuorumThresholdIndexes []byte
}

// IEigenDAServiceManagerBatchHeader is an auto generated low-level Go binding around an user-defined struct.
type IEigenDAServiceManagerBatchHeader struct {
	BlobHeadersRoot            [32]byte
	QuorumNumbers              []byte
	QuorumThresholdPercentages []byte
	ReferenceBlockNumber       uint32
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
	QuorumNumber                 uint8
	AdversaryThresholdPercentage uint8
	QuorumThresholdPercentage    uint8
	ChunkLength                  uint32
}

// ContractMockRollupMetaData contains all meta data concerning the ContractMockRollup contract.
var ContractMockRollupMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_eigenDAServiceManager\",\"type\":\"address\",\"internalType\":\"contractIEigenDAServiceManager\"},{\"name\":\"_tau\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"_illegalValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_stakeRequired\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"blacklist\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"challengeCommitment\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"point\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"commitments\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"polynomialCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAServiceManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAServiceManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"illegalValue\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"postCommitment\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structIEigenDAServiceManager.QuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"quorumThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structEigenDARollupUtils.BlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumThresholdPercentages\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumThresholdIndexes\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerValidator\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"stakeRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tau\",\"inputs\":[],\"outputs\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validators\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161189c38038061189c83398101604081905261002f91610069565b600080546001600160a01b0319166001600160a01b0395909516949094179093558151600155602090910151600255600355600455610107565b60008060008084860360a081121561008057600080fd5b85516001600160a01b038116811461009757600080fd5b94506040601f19820112156100ab57600080fd5b50604080519081016001600160401b03811182821017156100dc57634e487b7160e01b600052604160045260246000fd5b6040908152602087810151835290870151908201526060860151608090960151949790965092505050565b611786806101166000396000f3fe6080604052600436106100915760003560e01c8063bcc6587f11610059578063bcc6587f1461018e578063cfc4af5514610196578063f9f92be4146101c6578063fa52c7d814610206578063fc30cad01461023657600080fd5b806305010105146100965780634440bc5c146100bf57806349ce8997146100d55780635a191f721461014c5780636281e63b1461016e575b600080fd5b3480156100a257600080fd5b506100ac60045481565b6040519081526020015b60405180910390f35b3480156100cb57600080fd5b506100ac60035481565b3480156100e157600080fd5b5061013d6100f0366004610ef6565b6007602090815260009182526040918290208054835180850190945260018201548452600290910154918301919091526001600160a01b03811691600160a01b90910463ffffffff169083565b6040516100b693929190610f0f565b34801561015857600080fd5b5061016c610167366004611238565b61026e565b005b34801561017a57600080fd5b5061016c610189366004611415565b610467565b61016c610704565b3480156101a257600080fd5b506001546002546101b1919082565b604080519283526020830191909152016100b6565b3480156101d257600080fd5b506101f66101e136600461147d565b60066020526000908152604090205460ff1681565b60405190151581526020016100b6565b34801561021257600080fd5b506101f661022136600461147d565b60056020526000908152604090205460ff1681565b34801561024257600080fd5b50600054610256906001600160a01b031681565b6040516001600160a01b0390911681526020016100b6565b3360009081526005602052604090205460ff166102ee5760405162461bcd60e51b815260206004820152603360248201527f4d6f636b526f6c6c75702e706f7374436f6d6d69746d656e743a2056616c6964604482015272185d1bdc881b9bdd081c9959da5cdd195c9959606a1b60648201526084015b60405180910390fd5b426000908152600760205260409020546001600160a01b0316156103715760405162461bcd60e51b815260206004820152603460248201527f4d6f636b526f6c6c75702e706f7374436f6d6d69746d656e743a20436f6d6d696044820152731d1b595b9d08185b1c9958591e481c1bdcdd195960621b60648201526084016102e5565b60005460405163219460e160e21b815273__$32f04d18c688c2c57b0347c20a77f3d6c9$__916386518384916103b99186916001600160a01b039091169086906004016115ce565b60006040518083038186803b1580156103d157600080fd5b505af41580156103e5573d6000803e3d6000fd5b50506040805160608101825233815260208681015163ffffffff90811682840190815297518385019081524260009081526007845294909420925183549851909116600160a01b026001600160c01b03199098166001600160a01b03919091161796909617815590518051600183015590940151600290940193909355505050565b600083815260076020908152604091829020825160608101845281546001600160a01b038082168352600160a01b90910463ffffffff16828501528451808601865260018401548152600290930154938301939093529283015281511661052e5760405162461bcd60e51b815260206004820152603560248201527f4d6f636b526f6c6c75702e6368616c6c656e6765436f6d6d69746d656e743a2060448201527410dbdb5b5a5d1b595b9d081b9bdd081c1bdcdd1959605a1b60648201526084016102e5565b806020015163ffffffff1683106105b95760405162461bcd60e51b815260206004820152604360248201527f4d6f636b526f6c6c75702e6368616c6c656e6765436f6d6d69746d656e743a2060448201527f506f696e74206d757374206265206c657373207468616e2064617461206c656e6064820152620cee8d60eb1b608482015260a4016102e5565b60035460408051808201825260015481526002546020820152908301516105e49286929091866108a4565b6106615760405162461bcd60e51b815260206004820152604260248201527f4d6f636b526f6c6c75702e6368616c6c656e6765436f6d6d69746d656e743a2060448201527f446f6573206e6f74206576616c7561746520746f20696c6c6567616c2076616c606482015261756560f01b608482015260a4016102e5565b80516001600160a01b039081166000908152600560209081526040808320805460ff19908116909155855190941683526006909152808220805490931660011790925590513390670de0b6b3a7640000908381818185875af1925050503d80600081146106ea576040519150601f19603f3d011682016040523d82523d6000602084013e6106ef565b606091505b50509050806106fd57600080fd5b5050505050565b60045434146107865760405162461bcd60e51b815260206004820152604260248201527f4d6f636b526f6c6c75702e726567697374657256616c696461746f723a204d7560448201527f73742073656e64207374616b6520726571756972656420746f2072656769737460648201526132b960f11b608482015260a4016102e5565b3360009081526005602052604090205460ff161561080c5760405162461bcd60e51b815260206004820152603a60248201527f4d6f636b526f6c6c75702e726567697374657256616c696461746f723a20566160448201527f6c696461746f7220616c7265616479207265676973746572656400000000000060648201526084016102e5565b3360009081526006602052604090205460ff16156108885760405162461bcd60e51b815260206004820152603360248201527f4d6f636b526f6c6c75702e726567697374657256616c696461746f723a2056616044820152721b1a59185d1bdc88189b1858dadb1a5cdd1959606a1b60648201526084016102e5565b336000908152600560205260409020805460ff19166001179055565b6000806108db6108d6604080518082018252600080825260209182015281518083019092526001825260029082015290565b610921565b90506109166108f46108ed838a6109e0565b8790610a77565b84610909610902858b6109e0565b8890610a77565b610911610b0b565b610bcb565b979650505050505050565b6040805180820190915260008082526020820152815115801561094657506020820151155b15610964575050604080518082019091526000808252602082015290565b6040518060400160405280836000015181526020017f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4784602001516109a99190611699565b6109d3907f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd476116d1565b905292915050565b919050565b60408051808201909152600080825260208201526109fc610e3a565b835181526020808501519082015260408082018490526000908360608460076107d05a03fa9050808015610a2f57610a31565bfe5b5080610a6f5760405162461bcd60e51b815260206004820152600d60248201526c1958cb5b5d5b0b59985a5b1959609a1b60448201526064016102e5565b505092915050565b6040805180820190915260008082526020820152610a93610e58565b835181526020808501518183015283516040808401919091529084015160608301526000908360808460066107d05a03fa9050808015610a2f575080610a6f5760405162461bcd60e51b815260206004820152600d60248201526c1958cb5859190b59985a5b1959609a1b60448201526064016102e5565b610b13610e76565b50604080516080810182527f198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c28183019081527f1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed6060830152815281518083019092527f275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec82527f1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d60208381019190915281019190915290565b604080518082018252858152602080820185905282518084019093528583528201839052600091610bfa610e9b565b60005b6002811015610dbf576000610c138260066116fe565b9050848260028110610c2757610c276116e8565b60200201515183610c3983600061171d565b600c8110610c4957610c496116e8565b6020020152848260028110610c6057610c606116e8565b60200201516020015183826001610c77919061171d565b600c8110610c8757610c876116e8565b6020020152838260028110610c9e57610c9e6116e8565b6020020151515183610cb183600261171d565b600c8110610cc157610cc16116e8565b6020020152838260028110610cd857610cd86116e8565b6020020151516001602002015183610cf183600361171d565b600c8110610d0157610d016116e8565b6020020152838260028110610d1857610d186116e8565b602002015160200151600060028110610d3357610d336116e8565b602002015183610d4483600461171d565b600c8110610d5457610d546116e8565b6020020152838260028110610d6b57610d6b6116e8565b602002015160200151600160028110610d8657610d866116e8565b602002015183610d9783600561171d565b600c8110610da757610da76116e8565b60200201525080610db781611735565b915050610bfd565b50610dc8610eba565b60006020826101808560086107d05a03fa9050808015610a2f575080610e285760405162461bcd60e51b81526020600482015260156024820152741c185a5c9a5b99cb5bdc18dbd9194b59985a5b1959605a1b60448201526064016102e5565b5051151593505050505b949350505050565b60405180606001604052806003906020820280368337509192915050565b60405180608001604052806004906020820280368337509192915050565b6040518060400160405280610e89610ed8565b8152602001610e96610ed8565b905290565b604051806101800160405280600c906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b60405180604001604052806002906020820280368337509192915050565b600060208284031215610f0857600080fd5b5035919050565b6001600160a01b038416815263ffffffff8316602082015260808101610e32604083018480518252602090810151910152565b634e487b7160e01b600052604160045260246000fd5b6040516060810167ffffffffffffffff81118282101715610f7b57610f7b610f42565b60405290565b6040516080810167ffffffffffffffff81118282101715610f7b57610f7b610f42565b60405160a0810167ffffffffffffffff81118282101715610f7b57610f7b610f42565b6040805190810167ffffffffffffffff81118282101715610f7b57610f7b610f42565b604051601f8201601f1916810167ffffffffffffffff8111828210171561101357611013610f42565b604052919050565b803563ffffffff811681146109db57600080fd5b803560ff811681146109db57600080fd5b600082601f83011261105157600080fd5b813567ffffffffffffffff81111561106b5761106b610f42565b61107e601f8201601f1916602001610fea565b81815284602083860101111561109357600080fd5b816020850160208301376000918101602001919091529392505050565b6000606082840312156110c257600080fd5b6110ca610f58565b9050813567ffffffffffffffff808211156110e457600080fd5b90830190608082860312156110f857600080fd5b611100610f81565b8235815260208301358281111561111657600080fd5b61112287828601611040565b60208301525060408301358281111561113a57600080fd5b61114687828601611040565b6040830152506111586060840161101b565b606082015283525050602082810135908201526111776040830161101b565b604082015292915050565b600060a0828403121561119457600080fd5b61119c610fa4565b90506111a78261101b565b81526111b56020830161102f565b6020820152604082013567ffffffffffffffff808211156111d557600080fd5b6111e1858386016110b0565b604084015260608401359150808211156111fa57600080fd5b61120685838601611040565b6060840152608084013591508082111561121f57600080fd5b5061122c84828501611040565b60808301525092915050565b600080604080848603121561124c57600080fd5b833567ffffffffffffffff8082111561126457600080fd5b9085019081870360808082121561127a57600080fd5b611282610f58565b8583121561128f57600080fd5b611297610fc7565b925084358352602080860135818501528382526112b587870161101b565b818301526060935083860135858111156112ce57600080fd5b8087019650508a601f8701126112e357600080fd5b8535858111156112f5576112f5610f42565b611303828260051b01610fea565b81815260079190911b8701820190828101908d83111561132257600080fd5b978301975b8289101561138e5785898f03121561133f5760008081fd5b611347610f81565b6113508a61102f565b815261135d858b0161102f565b8582015261136c8b8b0161102f565b8b82015261137b888b0161101b565b8189015282529785019790830190611327565b988401989098525090975088013594505050808311156113ad57600080fd5b50506113bb85828601611182565b9150509250929050565b600082601f8301126113d657600080fd5b6113de610fc7565b8060408401858111156113f057600080fd5b845b8181101561140a5780358452602093840193016113f2565b509095945050505050565b600080600083850360c081121561142b57600080fd5b84359350602085013592506080603f198201121561144857600080fd5b50611451610fc7565b61145e86604087016113c5565b815261146d86608087016113c5565b6020820152809150509250925092565b60006020828403121561148f57600080fd5b81356001600160a01b03811681146114a657600080fd5b9392505050565b6000815180845260005b818110156114d3576020818501810151868301820152016114b7565b818111156114e5576000602083870101525b50601f01601f19169290920160200192915050565b600063ffffffff80835116845260ff6020840151166020850152604083015160a060408601528051606060a087015280516101008701526020810151608061012088015261154c6101808801826114ad565b9050604082015160ff198883030161014089015261156a82826114ad565b91505083606083015116610160880152602083015160c08801528360408401511660e08801526060860151935086810360608801526115a981856114ad565b9350505050608083015184820360808601526115c582826114ad565b95945050505050565b6060808252845180518383015260200151608083015260009060e0830160208088015163ffffffff80821660a088015260409150818a015160808060c08a01528582518088526101008b0191508684019750600093505b80841015611668578751805160ff908116845288820151811689850152878201511687840152890151851689830152968601966001939093019290820190611625565b506001600160a01b038c168a870152898103858b0152611688818c6114fa565b9d9c50505050505050505050505050565b6000826116b657634e487b7160e01b600052601260045260246000fd5b500690565b634e487b7160e01b600052601160045260246000fd5b6000828210156116e3576116e36116bb565b500390565b634e487b7160e01b600052603260045260246000fd5b6000816000190483118215151615611718576117186116bb565b500290565b60008219821115611730576117306116bb565b500190565b6000600019821415611749576117496116bb565b506001019056fea26469706673582212205998f4655a03dda95969d9a1afc4b711c6cfac5605b54d1edff2b372a2a30fab64736f6c634300080c0033",
}

// ContractMockRollupABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractMockRollupMetaData.ABI instead.
var ContractMockRollupABI = ContractMockRollupMetaData.ABI

// ContractMockRollupBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractMockRollupMetaData.Bin instead.
var ContractMockRollupBin = ContractMockRollupMetaData.Bin

// DeployContractMockRollup deploys a new Ethereum contract, binding an instance of ContractMockRollup to it.
func DeployContractMockRollup(auth *bind.TransactOpts, backend bind.ContractBackend, _eigenDAServiceManager common.Address, _tau BN254G1Point, _illegalValue *big.Int, _stakeRequired *big.Int) (common.Address, *types.Transaction, *ContractMockRollup, error) {
	parsed, err := ContractMockRollupMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractMockRollupBin), backend, _eigenDAServiceManager, _tau, _illegalValue, _stakeRequired)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractMockRollup{ContractMockRollupCaller: ContractMockRollupCaller{contract: contract}, ContractMockRollupTransactor: ContractMockRollupTransactor{contract: contract}, ContractMockRollupFilterer: ContractMockRollupFilterer{contract: contract}}, nil
}

// ContractMockRollup is an auto generated Go binding around an Ethereum contract.
type ContractMockRollup struct {
	ContractMockRollupCaller     // Read-only binding to the contract
	ContractMockRollupTransactor // Write-only binding to the contract
	ContractMockRollupFilterer   // Log filterer for contract events
}

// ContractMockRollupCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractMockRollupCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractMockRollupTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractMockRollupTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractMockRollupFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractMockRollupFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractMockRollupSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractMockRollupSession struct {
	Contract     *ContractMockRollup // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ContractMockRollupCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractMockRollupCallerSession struct {
	Contract *ContractMockRollupCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// ContractMockRollupTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractMockRollupTransactorSession struct {
	Contract     *ContractMockRollupTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// ContractMockRollupRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractMockRollupRaw struct {
	Contract *ContractMockRollup // Generic contract binding to access the raw methods on
}

// ContractMockRollupCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractMockRollupCallerRaw struct {
	Contract *ContractMockRollupCaller // Generic read-only contract binding to access the raw methods on
}

// ContractMockRollupTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractMockRollupTransactorRaw struct {
	Contract *ContractMockRollupTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractMockRollup creates a new instance of ContractMockRollup, bound to a specific deployed contract.
func NewContractMockRollup(address common.Address, backend bind.ContractBackend) (*ContractMockRollup, error) {
	contract, err := bindContractMockRollup(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractMockRollup{ContractMockRollupCaller: ContractMockRollupCaller{contract: contract}, ContractMockRollupTransactor: ContractMockRollupTransactor{contract: contract}, ContractMockRollupFilterer: ContractMockRollupFilterer{contract: contract}}, nil
}

// NewContractMockRollupCaller creates a new read-only instance of ContractMockRollup, bound to a specific deployed contract.
func NewContractMockRollupCaller(address common.Address, caller bind.ContractCaller) (*ContractMockRollupCaller, error) {
	contract, err := bindContractMockRollup(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractMockRollupCaller{contract: contract}, nil
}

// NewContractMockRollupTransactor creates a new write-only instance of ContractMockRollup, bound to a specific deployed contract.
func NewContractMockRollupTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractMockRollupTransactor, error) {
	contract, err := bindContractMockRollup(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractMockRollupTransactor{contract: contract}, nil
}

// NewContractMockRollupFilterer creates a new log filterer instance of ContractMockRollup, bound to a specific deployed contract.
func NewContractMockRollupFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractMockRollupFilterer, error) {
	contract, err := bindContractMockRollup(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractMockRollupFilterer{contract: contract}, nil
}

// bindContractMockRollup binds a generic wrapper to an already deployed contract.
func bindContractMockRollup(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractMockRollupMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractMockRollup *ContractMockRollupRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractMockRollup.Contract.ContractMockRollupCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractMockRollup *ContractMockRollupRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.ContractMockRollupTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractMockRollup *ContractMockRollupRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.ContractMockRollupTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractMockRollup *ContractMockRollupCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractMockRollup.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractMockRollup *ContractMockRollupTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractMockRollup *ContractMockRollupTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.contract.Transact(opts, method, params...)
}

// Blacklist is a free data retrieval call binding the contract method 0xf9f92be4.
//
// Solidity: function blacklist(address ) view returns(bool)
func (_ContractMockRollup *ContractMockRollupCaller) Blacklist(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _ContractMockRollup.contract.Call(opts, &out, "blacklist", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Blacklist is a free data retrieval call binding the contract method 0xf9f92be4.
//
// Solidity: function blacklist(address ) view returns(bool)
func (_ContractMockRollup *ContractMockRollupSession) Blacklist(arg0 common.Address) (bool, error) {
	return _ContractMockRollup.Contract.Blacklist(&_ContractMockRollup.CallOpts, arg0)
}

// Blacklist is a free data retrieval call binding the contract method 0xf9f92be4.
//
// Solidity: function blacklist(address ) view returns(bool)
func (_ContractMockRollup *ContractMockRollupCallerSession) Blacklist(arg0 common.Address) (bool, error) {
	return _ContractMockRollup.Contract.Blacklist(&_ContractMockRollup.CallOpts, arg0)
}

// Commitments is a free data retrieval call binding the contract method 0x49ce8997.
//
// Solidity: function commitments(uint256 ) view returns(address validator, uint32 dataLength, (uint256,uint256) polynomialCommitment)
func (_ContractMockRollup *ContractMockRollupCaller) Commitments(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Validator            common.Address
	DataLength           uint32
	PolynomialCommitment BN254G1Point
}, error) {
	var out []interface{}
	err := _ContractMockRollup.contract.Call(opts, &out, "commitments", arg0)

	outstruct := new(struct {
		Validator            common.Address
		DataLength           uint32
		PolynomialCommitment BN254G1Point
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Validator = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.DataLength = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.PolynomialCommitment = *abi.ConvertType(out[2], new(BN254G1Point)).(*BN254G1Point)

	return *outstruct, err

}

// Commitments is a free data retrieval call binding the contract method 0x49ce8997.
//
// Solidity: function commitments(uint256 ) view returns(address validator, uint32 dataLength, (uint256,uint256) polynomialCommitment)
func (_ContractMockRollup *ContractMockRollupSession) Commitments(arg0 *big.Int) (struct {
	Validator            common.Address
	DataLength           uint32
	PolynomialCommitment BN254G1Point
}, error) {
	return _ContractMockRollup.Contract.Commitments(&_ContractMockRollup.CallOpts, arg0)
}

// Commitments is a free data retrieval call binding the contract method 0x49ce8997.
//
// Solidity: function commitments(uint256 ) view returns(address validator, uint32 dataLength, (uint256,uint256) polynomialCommitment)
func (_ContractMockRollup *ContractMockRollupCallerSession) Commitments(arg0 *big.Int) (struct {
	Validator            common.Address
	DataLength           uint32
	PolynomialCommitment BN254G1Point
}, error) {
	return _ContractMockRollup.Contract.Commitments(&_ContractMockRollup.CallOpts, arg0)
}

// EigenDAServiceManager is a free data retrieval call binding the contract method 0xfc30cad0.
//
// Solidity: function eigenDAServiceManager() view returns(address)
func (_ContractMockRollup *ContractMockRollupCaller) EigenDAServiceManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractMockRollup.contract.Call(opts, &out, "eigenDAServiceManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDAServiceManager is a free data retrieval call binding the contract method 0xfc30cad0.
//
// Solidity: function eigenDAServiceManager() view returns(address)
func (_ContractMockRollup *ContractMockRollupSession) EigenDAServiceManager() (common.Address, error) {
	return _ContractMockRollup.Contract.EigenDAServiceManager(&_ContractMockRollup.CallOpts)
}

// EigenDAServiceManager is a free data retrieval call binding the contract method 0xfc30cad0.
//
// Solidity: function eigenDAServiceManager() view returns(address)
func (_ContractMockRollup *ContractMockRollupCallerSession) EigenDAServiceManager() (common.Address, error) {
	return _ContractMockRollup.Contract.EigenDAServiceManager(&_ContractMockRollup.CallOpts)
}

// IllegalValue is a free data retrieval call binding the contract method 0x4440bc5c.
//
// Solidity: function illegalValue() view returns(uint256)
func (_ContractMockRollup *ContractMockRollupCaller) IllegalValue(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractMockRollup.contract.Call(opts, &out, "illegalValue")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// IllegalValue is a free data retrieval call binding the contract method 0x4440bc5c.
//
// Solidity: function illegalValue() view returns(uint256)
func (_ContractMockRollup *ContractMockRollupSession) IllegalValue() (*big.Int, error) {
	return _ContractMockRollup.Contract.IllegalValue(&_ContractMockRollup.CallOpts)
}

// IllegalValue is a free data retrieval call binding the contract method 0x4440bc5c.
//
// Solidity: function illegalValue() view returns(uint256)
func (_ContractMockRollup *ContractMockRollupCallerSession) IllegalValue() (*big.Int, error) {
	return _ContractMockRollup.Contract.IllegalValue(&_ContractMockRollup.CallOpts)
}

// StakeRequired is a free data retrieval call binding the contract method 0x05010105.
//
// Solidity: function stakeRequired() view returns(uint256)
func (_ContractMockRollup *ContractMockRollupCaller) StakeRequired(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractMockRollup.contract.Call(opts, &out, "stakeRequired")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StakeRequired is a free data retrieval call binding the contract method 0x05010105.
//
// Solidity: function stakeRequired() view returns(uint256)
func (_ContractMockRollup *ContractMockRollupSession) StakeRequired() (*big.Int, error) {
	return _ContractMockRollup.Contract.StakeRequired(&_ContractMockRollup.CallOpts)
}

// StakeRequired is a free data retrieval call binding the contract method 0x05010105.
//
// Solidity: function stakeRequired() view returns(uint256)
func (_ContractMockRollup *ContractMockRollupCallerSession) StakeRequired() (*big.Int, error) {
	return _ContractMockRollup.Contract.StakeRequired(&_ContractMockRollup.CallOpts)
}

// Tau is a free data retrieval call binding the contract method 0xcfc4af55.
//
// Solidity: function tau() view returns(uint256 X, uint256 Y)
func (_ContractMockRollup *ContractMockRollupCaller) Tau(opts *bind.CallOpts) (struct {
	X *big.Int
	Y *big.Int
}, error) {
	var out []interface{}
	err := _ContractMockRollup.contract.Call(opts, &out, "tau")

	outstruct := new(struct {
		X *big.Int
		Y *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.X = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Y = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Tau is a free data retrieval call binding the contract method 0xcfc4af55.
//
// Solidity: function tau() view returns(uint256 X, uint256 Y)
func (_ContractMockRollup *ContractMockRollupSession) Tau() (struct {
	X *big.Int
	Y *big.Int
}, error) {
	return _ContractMockRollup.Contract.Tau(&_ContractMockRollup.CallOpts)
}

// Tau is a free data retrieval call binding the contract method 0xcfc4af55.
//
// Solidity: function tau() view returns(uint256 X, uint256 Y)
func (_ContractMockRollup *ContractMockRollupCallerSession) Tau() (struct {
	X *big.Int
	Y *big.Int
}, error) {
	return _ContractMockRollup.Contract.Tau(&_ContractMockRollup.CallOpts)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) view returns(bool)
func (_ContractMockRollup *ContractMockRollupCaller) Validators(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _ContractMockRollup.contract.Call(opts, &out, "validators", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) view returns(bool)
func (_ContractMockRollup *ContractMockRollupSession) Validators(arg0 common.Address) (bool, error) {
	return _ContractMockRollup.Contract.Validators(&_ContractMockRollup.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) view returns(bool)
func (_ContractMockRollup *ContractMockRollupCallerSession) Validators(arg0 common.Address) (bool, error) {
	return _ContractMockRollup.Contract.Validators(&_ContractMockRollup.CallOpts, arg0)
}

// ChallengeCommitment is a paid mutator transaction binding the contract method 0x6281e63b.
//
// Solidity: function challengeCommitment(uint256 timestamp, uint256 point, (uint256[2],uint256[2]) proof) returns()
func (_ContractMockRollup *ContractMockRollupTransactor) ChallengeCommitment(opts *bind.TransactOpts, timestamp *big.Int, point *big.Int, proof BN254G2Point) (*types.Transaction, error) {
	return _ContractMockRollup.contract.Transact(opts, "challengeCommitment", timestamp, point, proof)
}

// ChallengeCommitment is a paid mutator transaction binding the contract method 0x6281e63b.
//
// Solidity: function challengeCommitment(uint256 timestamp, uint256 point, (uint256[2],uint256[2]) proof) returns()
func (_ContractMockRollup *ContractMockRollupSession) ChallengeCommitment(timestamp *big.Int, point *big.Int, proof BN254G2Point) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.ChallengeCommitment(&_ContractMockRollup.TransactOpts, timestamp, point, proof)
}

// ChallengeCommitment is a paid mutator transaction binding the contract method 0x6281e63b.
//
// Solidity: function challengeCommitment(uint256 timestamp, uint256 point, (uint256[2],uint256[2]) proof) returns()
func (_ContractMockRollup *ContractMockRollupTransactorSession) ChallengeCommitment(timestamp *big.Int, point *big.Int, proof BN254G2Point) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.ChallengeCommitment(&_ContractMockRollup.TransactOpts, timestamp, point, proof)
}

// PostCommitment is a paid mutator transaction binding the contract method 0x5a191f72.
//
// Solidity: function postCommitment(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint8,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) returns()
func (_ContractMockRollup *ContractMockRollupTransactor) PostCommitment(opts *bind.TransactOpts, blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof EigenDARollupUtilsBlobVerificationProof) (*types.Transaction, error) {
	return _ContractMockRollup.contract.Transact(opts, "postCommitment", blobHeader, blobVerificationProof)
}

// PostCommitment is a paid mutator transaction binding the contract method 0x5a191f72.
//
// Solidity: function postCommitment(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint8,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) returns()
func (_ContractMockRollup *ContractMockRollupSession) PostCommitment(blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof EigenDARollupUtilsBlobVerificationProof) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.PostCommitment(&_ContractMockRollup.TransactOpts, blobHeader, blobVerificationProof)
}

// PostCommitment is a paid mutator transaction binding the contract method 0x5a191f72.
//
// Solidity: function postCommitment(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint8,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) returns()
func (_ContractMockRollup *ContractMockRollupTransactorSession) PostCommitment(blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof EigenDARollupUtilsBlobVerificationProof) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.PostCommitment(&_ContractMockRollup.TransactOpts, blobHeader, blobVerificationProof)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0xbcc6587f.
//
// Solidity: function registerValidator() payable returns()
func (_ContractMockRollup *ContractMockRollupTransactor) RegisterValidator(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractMockRollup.contract.Transact(opts, "registerValidator")
}

// RegisterValidator is a paid mutator transaction binding the contract method 0xbcc6587f.
//
// Solidity: function registerValidator() payable returns()
func (_ContractMockRollup *ContractMockRollupSession) RegisterValidator() (*types.Transaction, error) {
	return _ContractMockRollup.Contract.RegisterValidator(&_ContractMockRollup.TransactOpts)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0xbcc6587f.
//
// Solidity: function registerValidator() payable returns()
func (_ContractMockRollup *ContractMockRollupTransactorSession) RegisterValidator() (*types.Transaction, error) {
	return _ContractMockRollup.Contract.RegisterValidator(&_ContractMockRollup.TransactOpts)
}
