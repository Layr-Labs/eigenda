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
	BatchId        uint32
	BlobIndex      uint32
	BatchMetadata  IEigenDAServiceManagerBatchMetadata
	InclusionProof []byte
	QuorumIndices  []byte
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

// ContractMockRollupMetaData contains all meta data concerning the ContractMockRollup contract.
var ContractMockRollupMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_eigenDAServiceManager\",\"type\":\"address\",\"internalType\":\"contractIEigenDAServiceManager\"},{\"name\":\"_tau\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"challengeCommitment\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"point\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"challengeValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"commitments\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"confirmer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"polynomialCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAServiceManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAServiceManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"postCommitment\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structIEigenDAServiceManager.QuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structEigenDARollupUtils.BlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"tau\",\"inputs\":[],\"outputs\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161140e38038061140e83398101604081905261002f91610061565b600080546001600160a01b0319166001600160a01b0393909316929092179091558051600155602001516002556100ef565b600080828403606081121561007557600080fd5b83516001600160a01b038116811461008c57600080fd5b92506040601f19820112156100a057600080fd5b50604080519081016001600160401b03811182821017156100d157634e487b7160e01b600052604160045260246000fd5b60409081526020858101518352940151938101939093525092909150565b611310806100fe6000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c806349ce89971461005c578063b5144c73146100cf578063cfc4af55146100e4578063d2d16eb214610107578063fc30cad01461012a575b600080fd5b6100b761006a366004610aab565b6003602090815260009182526040918290208054835180850190945260018201548452600290910154918301919091526001600160a01b03811691600160a01b90910463ffffffff169083565b6040516100c693929190610ac4565b60405180910390f35b6100e26100dd366004610ded565b610155565b005b6001546002546100f2919082565b604080519283526020830191909152016100c6565b61011a610115366004610fca565b6102d3565b60405190151581526020016100c6565b60005461013d906001600160a01b031681565b6040516001600160a01b0390911681526020016100c6565b426000908152600360205260409020546001600160a01b0316156101dd5760405162461bcd60e51b815260206004820152603460248201527f4d6f636b526f6c6c75702e706f7374436f6d6d69746d656e743a20436f6d6d696044820152731d1b595b9d08185b1c9958591e481c1bdcdd195960621b60648201526084015b60405180910390fd5b60005460405163219460e160e21b815273__$32f04d18c688c2c57b0347c20a77f3d6c9$__916386518384916102259186916001600160a01b03909116908690600401611158565b60006040518083038186803b15801561023d57600080fd5b505af4158015610251573d6000803e3d6000fd5b50506040805160608101825233815260208681015163ffffffff90811682840190815297518385019081524260009081526003845294909420925183549851909116600160a01b026001600160c01b03199098166001600160a01b03919091161796909617815590518051600183015590940151600290940193909355505050565b6000848152600360209081526040808320815160608101835281546001600160a01b038082168352600160a01b90910463ffffffff168286015283518085018552600184015481526002909301549483019490945291820152805190911661039b5760405162461bcd60e51b815260206004820152603560248201527f4d6f636b526f6c6c75702e6368616c6c656e6765436f6d6d69746d656e743a2060448201527410dbdb5b5a5d1b595b9d081b9bdd081c1bdcdd1959605a1b60648201526084016101d4565b806020015163ffffffff1685106104265760405162461bcd60e51b815260206004820152604360248201527f4d6f636b526f6c6c75702e6368616c6c656e6765436f6d6d69746d656e743a2060448201527f506f696e74206d757374206265206c657373207468616e2064617461206c656e6064820152620cee8d60eb1b608482015260a4016101d4565b604080518082018252600154815260025460208201529082015161044f9187918691908861045b565b9150505b949350505050565b60008061049261048d604080518082018252600080825260209182015281518083019092526001825260029082015290565b6104d8565b90506104cd6104ab6104a4838a610597565b879061062e565b846104c06104b9858b610597565b889061062e565b6104c86106c2565b610782565b979650505050505050565b604080518082019091526000808252602082015281511580156104fd57506020820151155b1561051b575050604080518082019091526000808252602082015290565b6040518060400160405280836000015181526020017f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4784602001516105609190611223565b61058a907f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4761125b565b905292915050565b919050565b60408051808201909152600080825260208201526105b36109ef565b835181526020808501519082015260408082018490526000908360608460076107d05a03fa90508080156105e6576105e8565bfe5b50806106265760405162461bcd60e51b815260206004820152600d60248201526c1958cb5b5d5b0b59985a5b1959609a1b60448201526064016101d4565b505092915050565b604080518082019091526000808252602082015261064a610a0d565b835181526020808501518183015283516040808401919091529084015160608301526000908360808460066107d05a03fa90508080156105e65750806106265760405162461bcd60e51b815260206004820152600d60248201526c1958cb5859190b59985a5b1959609a1b60448201526064016101d4565b6106ca610a2b565b50604080516080810182527f198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c28183019081527f1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed6060830152815281518083019092527f275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec82527f1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d60208381019190915281019190915290565b6040805180820182528581526020808201859052825180840190935285835282018390526000916107b1610a50565b60005b60028110156109765760006107ca826006611288565b90508482600281106107de576107de611272565b602002015151836107f08360006112a7565b600c811061080057610800611272565b602002015284826002811061081757610817611272565b6020020151602001518382600161082e91906112a7565b600c811061083e5761083e611272565b602002015283826002811061085557610855611272565b60200201515151836108688360026112a7565b600c811061087857610878611272565b602002015283826002811061088f5761088f611272565b60200201515160016020020151836108a88360036112a7565b600c81106108b8576108b8611272565b60200201528382600281106108cf576108cf611272565b6020020151602001516000600281106108ea576108ea611272565b6020020151836108fb8360046112a7565b600c811061090b5761090b611272565b602002015283826002811061092257610922611272565b60200201516020015160016002811061093d5761093d611272565b60200201518361094e8360056112a7565b600c811061095e5761095e611272565b6020020152508061096e816112bf565b9150506107b4565b5061097f610a6f565b60006020826101808560086107d05a03fa90508080156105e65750806109df5760405162461bcd60e51b81526020600482015260156024820152741c185a5c9a5b99cb5bdc18dbd9194b59985a5b1959605a1b60448201526064016101d4565b5051151598975050505050505050565b60405180606001604052806003906020820280368337509192915050565b60405180608001604052806004906020820280368337509192915050565b6040518060400160405280610a3e610a8d565b8152602001610a4b610a8d565b905290565b604051806101800160405280600c906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b60405180604001604052806002906020820280368337509192915050565b600060208284031215610abd57600080fd5b5035919050565b6001600160a01b038416815263ffffffff8316602082015260808101610453604083018480518252602090810151910152565b634e487b7160e01b600052604160045260246000fd5b6040516060810167ffffffffffffffff81118282101715610b3057610b30610af7565b60405290565b6040516080810167ffffffffffffffff81118282101715610b3057610b30610af7565b60405160a0810167ffffffffffffffff81118282101715610b3057610b30610af7565b6040805190810167ffffffffffffffff81118282101715610b3057610b30610af7565b604051601f8201601f1916810167ffffffffffffffff81118282101715610bc857610bc8610af7565b604052919050565b803563ffffffff8116811461059257600080fd5b803560ff8116811461059257600080fd5b600082601f830112610c0657600080fd5b813567ffffffffffffffff811115610c2057610c20610af7565b610c33601f8201601f1916602001610b9f565b818152846020838601011115610c4857600080fd5b816020850160208301376000918101602001919091529392505050565b600060608284031215610c7757600080fd5b610c7f610b0d565b9050813567ffffffffffffffff80821115610c9957600080fd5b9083019060808286031215610cad57600080fd5b610cb5610b36565b82358152602083013582811115610ccb57600080fd5b610cd787828601610bf5565b602083015250604083013582811115610cef57600080fd5b610cfb87828601610bf5565b604083015250610d0d60608401610bd0565b60608201528352505060208281013590820152610d2c60408301610bd0565b604082015292915050565b600060a08284031215610d4957600080fd5b610d51610b59565b9050610d5c82610bd0565b8152610d6a60208301610bd0565b6020820152604082013567ffffffffffffffff80821115610d8a57600080fd5b610d9685838601610c65565b60408401526060840135915080821115610daf57600080fd5b610dbb85838601610bf5565b60608401526080840135915080821115610dd457600080fd5b50610de184828501610bf5565b60808301525092915050565b6000806040808486031215610e0157600080fd5b833567ffffffffffffffff80821115610e1957600080fd5b90850190818703608080821215610e2f57600080fd5b610e37610b0d565b85831215610e4457600080fd5b610e4c610b7c565b92508435835260208086013581850152838252610e6a878701610bd0565b81830152606093508386013585811115610e8357600080fd5b8087019650508a601f870112610e9857600080fd5b853585811115610eaa57610eaa610af7565b610eb8828260051b01610b9f565b81815260079190911b8701820190828101908d831115610ed757600080fd5b978301975b82891015610f435785898f031215610ef45760008081fd5b610efc610b36565b610f058a610be4565b8152610f12858b01610be4565b85820152610f218b8b01610be4565b8b820152610f30888b01610bd0565b8189015282529785019790830190610edc565b98840198909852509097508801359450505080831115610f6257600080fd5b5050610f7085828601610d37565b9150509250929050565b600082601f830112610f8b57600080fd5b610f93610b7c565b806040840185811115610fa557600080fd5b845b81811015610fbf578035845260209384019301610fa7565b509095945050505050565b60008060008084860360e0811215610fe157600080fd5b85359450602086013593506080603f1982011215610ffe57600080fd5b50611007610b7c565b6110148760408801610f7a565b81526110238760808801610f7a565b60208201529396929550929360c00135925050565b6000815180845260005b8181101561105e57602081850181015186830182015201611042565b81811115611070576000602083870101525b50601f01601f19169290920160200192915050565b600063ffffffff808351168452806020840151166020850152604083015160a060408601528051606060a08701528051610100870152602081015160806101208801526110d6610180880182611038565b9050604082015160ff19888303016101408901526110f48282611038565b91505083606083015116610160880152602083015160c08801528360408401511660e08801526060860151935086810360608801526111338185611038565b93505050506080830151848203608086015261114f8282611038565b95945050505050565b6060808252845180518383015260200151608083015260009060e0830160208088015163ffffffff80821660a088015260409150818a015160808060c08a01528582518088526101008b0191508684019750600093505b808410156111f2578751805160ff9081168452888201518116898501528782015116878401528901518516898301529686019660019390930192908201906111af565b506001600160a01b038c168a870152898103858b0152611212818c611085565b9d9c50505050505050505050505050565b60008261124057634e487b7160e01b600052601260045260246000fd5b500690565b634e487b7160e01b600052601160045260246000fd5b60008282101561126d5761126d611245565b500390565b634e487b7160e01b600052603260045260246000fd5b60008160001904831182151516156112a2576112a2611245565b500290565b600082198211156112ba576112ba611245565b500190565b60006000198214156112d3576112d3611245565b506001019056fea2646970667358221220ed148c7173dd91988170fe356fbfcaa4f309e8b595bf034cf086ce0a847c2b2064736f6c634300080c0033",
}

// ContractMockRollupABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractMockRollupMetaData.ABI instead.
var ContractMockRollupABI = ContractMockRollupMetaData.ABI

// ContractMockRollupBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractMockRollupMetaData.Bin instead.
var ContractMockRollupBin = ContractMockRollupMetaData.Bin

// DeployContractMockRollup deploys a new Ethereum contract, binding an instance of ContractMockRollup to it.
func DeployContractMockRollup(auth *bind.TransactOpts, backend bind.ContractBackend, _eigenDAServiceManager common.Address, _tau BN254G1Point) (common.Address, *types.Transaction, *ContractMockRollup, error) {
	parsed, err := ContractMockRollupMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractMockRollupBin), backend, _eigenDAServiceManager, _tau)
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

// Commitments is a free data retrieval call binding the contract method 0x49ce8997.
//
// Solidity: function commitments(uint256 ) view returns(address confirmer, uint32 dataLength, (uint256,uint256) polynomialCommitment)
func (_ContractMockRollup *ContractMockRollupCaller) Commitments(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Confirmer            common.Address
	DataLength           uint32
	PolynomialCommitment BN254G1Point
}, error) {
	var out []interface{}
	err := _ContractMockRollup.contract.Call(opts, &out, "commitments", arg0)

	outstruct := new(struct {
		Confirmer            common.Address
		DataLength           uint32
		PolynomialCommitment BN254G1Point
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Confirmer = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.DataLength = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.PolynomialCommitment = *abi.ConvertType(out[2], new(BN254G1Point)).(*BN254G1Point)

	return *outstruct, err

}

// Commitments is a free data retrieval call binding the contract method 0x49ce8997.
//
// Solidity: function commitments(uint256 ) view returns(address confirmer, uint32 dataLength, (uint256,uint256) polynomialCommitment)
func (_ContractMockRollup *ContractMockRollupSession) Commitments(arg0 *big.Int) (struct {
	Confirmer            common.Address
	DataLength           uint32
	PolynomialCommitment BN254G1Point
}, error) {
	return _ContractMockRollup.Contract.Commitments(&_ContractMockRollup.CallOpts, arg0)
}

// Commitments is a free data retrieval call binding the contract method 0x49ce8997.
//
// Solidity: function commitments(uint256 ) view returns(address confirmer, uint32 dataLength, (uint256,uint256) polynomialCommitment)
func (_ContractMockRollup *ContractMockRollupCallerSession) Commitments(arg0 *big.Int) (struct {
	Confirmer            common.Address
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

// ChallengeCommitment is a paid mutator transaction binding the contract method 0xd2d16eb2.
//
// Solidity: function challengeCommitment(uint256 timestamp, uint256 point, (uint256[2],uint256[2]) proof, uint256 challengeValue) returns(bool)
func (_ContractMockRollup *ContractMockRollupTransactor) ChallengeCommitment(opts *bind.TransactOpts, timestamp *big.Int, point *big.Int, proof BN254G2Point, challengeValue *big.Int) (*types.Transaction, error) {
	return _ContractMockRollup.contract.Transact(opts, "challengeCommitment", timestamp, point, proof, challengeValue)
}

// ChallengeCommitment is a paid mutator transaction binding the contract method 0xd2d16eb2.
//
// Solidity: function challengeCommitment(uint256 timestamp, uint256 point, (uint256[2],uint256[2]) proof, uint256 challengeValue) returns(bool)
func (_ContractMockRollup *ContractMockRollupSession) ChallengeCommitment(timestamp *big.Int, point *big.Int, proof BN254G2Point, challengeValue *big.Int) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.ChallengeCommitment(&_ContractMockRollup.TransactOpts, timestamp, point, proof, challengeValue)
}

// ChallengeCommitment is a paid mutator transaction binding the contract method 0xd2d16eb2.
//
// Solidity: function challengeCommitment(uint256 timestamp, uint256 point, (uint256[2],uint256[2]) proof, uint256 challengeValue) returns(bool)
func (_ContractMockRollup *ContractMockRollupTransactorSession) ChallengeCommitment(timestamp *big.Int, point *big.Int, proof BN254G2Point, challengeValue *big.Int) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.ChallengeCommitment(&_ContractMockRollup.TransactOpts, timestamp, point, proof, challengeValue)
}

// PostCommitment is a paid mutator transaction binding the contract method 0xb5144c73.
//
// Solidity: function postCommitment(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) returns()
func (_ContractMockRollup *ContractMockRollupTransactor) PostCommitment(opts *bind.TransactOpts, blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof EigenDARollupUtilsBlobVerificationProof) (*types.Transaction, error) {
	return _ContractMockRollup.contract.Transact(opts, "postCommitment", blobHeader, blobVerificationProof)
}

// PostCommitment is a paid mutator transaction binding the contract method 0xb5144c73.
//
// Solidity: function postCommitment(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) returns()
func (_ContractMockRollup *ContractMockRollupSession) PostCommitment(blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof EigenDARollupUtilsBlobVerificationProof) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.PostCommitment(&_ContractMockRollup.TransactOpts, blobHeader, blobVerificationProof)
}

// PostCommitment is a paid mutator transaction binding the contract method 0xb5144c73.
//
// Solidity: function postCommitment(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) returns()
func (_ContractMockRollup *ContractMockRollupTransactorSession) PostCommitment(blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof EigenDARollupUtilsBlobVerificationProof) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.PostCommitment(&_ContractMockRollup.TransactOpts, blobHeader, blobVerificationProof)
}
