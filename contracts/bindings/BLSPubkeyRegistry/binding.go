// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractBLSPubkeyRegistry

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

// IBLSPubkeyRegistryApkUpdate is an auto generated low-level Go binding around an user-defined struct.
type IBLSPubkeyRegistryApkUpdate struct {
	ApkHash               [24]byte
	UpdateBlockNumber     uint32
	NextUpdateBlockNumber uint32
}

// ContractBLSPubkeyRegistryMetaData contains all meta data concerning the ContractBLSPubkeyRegistry contract.
var ContractBLSPubkeyRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"},{\"name\":\"_pubkeyCompendium\",\"type\":\"address\",\"internalType\":\"contractIBLSPublicKeyCompendium\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deregisterOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubkey\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getApkForQuorum\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApkHashForQuorumAtBlockNumberFromIndex\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"blockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes24\",\"internalType\":\"bytes24\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApkIndicesForQuorumsAtBlockNumber\",\"inputs\":[{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApkUpdateForQuorumByIndex\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIBLSPubkeyRegistry.ApkUpdate\",\"components\":[{\"name\":\"apkHash\",\"type\":\"bytes24\",\"internalType\":\"bytes24\"},{\"name\":\"updateBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"nextUpdateBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorFromPubkeyHash\",\"inputs\":[{\"name\":\"pubkeyHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumApkHistoryLength\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pubkeyCompendium\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIBLSPublicKeyCompendium\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumApk\",\"inputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumApkUpdates\",\"inputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"apkHash\",\"type\":\"bytes24\",\"internalType\":\"bytes24\"},{\"name\":\"updateBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"nextUpdateBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubkey\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registryCoordinator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorAddedToQuorums\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorRemovedFromQuorums\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false}]",
	Bin: "0x60c06040523480156200001157600080fd5b506040516200158d3803806200158d833981016040819052620000349162000138565b6001600160a01b03808316608052811660a0528181620000536200005d565b5050505062000177565b600054610100900460ff1615620000ca5760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840160405180910390fd5b60005460ff90811610156200011d576000805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b6001600160a01b03811681146200013557600080fd5b50565b600080604083850312156200014c57600080fd5b825162000159816200011f565b60208401519092506200016c816200011f565b809150509250929050565b60805160a0516113dc620001b16000396000818160e401526105990152600081816101f10152818161033401526104ac01526113dc6000f3fe608060405234801561001057600080fd5b50600436106100b45760003560e01c80636d14a987116100715780636d14a987146101ec5780637225057e146102135780637f5eccbb14610261578063c1af6b24146102a2578063eda10763146102cf578063fb81a7be146102ef57600080fd5b806303ce4bad146100b9578063187548c8146100df57806324369b2a1461011e57806332de63081461013357806347b314e81461016f57806363a9451014610182575b600080fd5b6100cc6100c7366004610f4d565b610327565b6040519081526020015b60405180910390f35b6101067f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b0390911681526020016100d6565b61013161012c366004610f4d565b6104a1565b005b61015a610141366004611017565b6002602052600090815260409020805460019091015482565b604080519283526020830191909152016100d6565b61010661017d366004611039565b610580565b6101d1610190366004611017565b60408051808201909152600080825260208201525060ff16600090815260026020908152604091829020825180840190935280548352600101549082015290565b604080518251815260209283015192810192909252016100d6565b6101067f000000000000000000000000000000000000000000000000000000000000000081565b610226610221366004611052565b610612565b60408051825167ffffffffffffffff1916815260208084015163ffffffff9081169183019190915292820151909216908201526060016100d6565b61027461026f366004611052565b6106a4565b6040805167ffffffffffffffff19909416845263ffffffff92831660208501529116908201526060016100d6565b6102b56102b036600461107c565b6106ef565b60405167ffffffffffffffff1990911681526020016100d6565b6102e26102dd3660046110c4565b610778565b6040516100d6919061113c565b6103126102fd366004611017565b60ff1660009081526001602052604090205490565b60405163ffffffff90911681526020016100d6565b6000336001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161461037a5760405162461bcd60e51b815260040161037190611186565b60405180910390fd5b6000610385836109d0565b90507fad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb581141561041d5760405162461bcd60e51b815260206004820152603f60248201527f424c535075626b657952656769737472792e72656769737465724f706572617460448201527f6f723a2063616e6e6f74207265676973746572207a65726f207075626b6579006064820152608401610371565b846001600160a01b031661043082610580565b6001600160a01b0316146104565760405162461bcd60e51b8152600401610371906111fd565b6104608484610a13565b7f5358c5b42179178c8fc757734ac2a3198f9071c765ee0d8389211525f5005246858560405161049192919061125b565b60405180910390a1949350505050565b336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146104e95760405162461bcd60e51b815260040161037190611186565b60006104f4826109d0565b9050836001600160a01b031661050982610580565b6001600160a01b03161461052f5760405162461bcd60e51b8152600401610371906111fd565b6105418361053c84610bd2565b610a13565b7f14a5172b312e9d2c22b8468f9c70ec2caa9de934fe380734fbc6f3beff2b14ba848460405161057292919061125b565b60405180910390a150505050565b60405163745dcd7360e11b8152600481018290526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063e8bb9ae690602401602060405180830381865afa1580156105e8573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061060c91906112c0565b92915050565b604080516060810182526000808252602080830182905282840182905260ff86168252600190529190912080548390811061064f5761064f6112dd565b600091825260209182902060408051606081018252919092015467ffffffffffffffff1981841b16825263ffffffff600160c01b8204811694830194909452600160e01b900490921690820152905092915050565b600160205281600052604060002081815481106106c057600080fd5b600091825260209091200154604081901b925063ffffffff600160c01b820481169250600160e01b9091041683565b60ff83166000908152600160205260408120805482919084908110610716576107166112dd565b600091825260209182902060408051606081018252919092015467ffffffffffffffff1981841b16825263ffffffff600160c01b8204811694830194909452600160e01b900490921690820152905061076f8185610c91565b51949350505050565b606060008367ffffffffffffffff81111561079557610795610eb7565b6040519080825280602002602001820160405280156107be578160200160208202803683370190505b50905060005b848110156109c75760008686838181106107e0576107e06112dd565b919091013560f81c600081815260016020526040902054909250905063ffffffff81161580610849575060ff82166000908152600160205260408120805490919061082d5761082d6112dd565b600091825260209091200154600160c01b900463ffffffff1686105b156108e25760405162461bcd60e51b815260206004820152605e60248201527f424c535075626b657952656769737472792e67657441706b496e64696365734660448201527f6f7251756f72756d734174426c6f636b4e756d6265723a20626c6f636b4e756d60648201527f626572206973206265666f726520746865206669727374207570646174650000608482015260a401610371565b60005b8163ffffffff168163ffffffff1610156109b15760ff83166000908152600160208190526040909120889161091a8486611309565b6109249190611309565b63ffffffff168154811061093a5761093a6112dd565b600091825260209091200154600160c01b900463ffffffff161161099f5760016109648284611309565b61096e9190611309565b858581518110610980576109806112dd565b602002602001019063ffffffff16908163ffffffff16815250506109b1565b806109a98161132e565b9150506108e5565b50505080806109bf90611352565b9150506107c4565b50949350505050565b6000816000015182602001516040516020016109f6929190918252602082015260400190565b604051602081830303815290604052805190602001209050919050565b604080518082019091526000808252602082015260005b8351811015610bcc576000848281518110610a4757610a476112dd565b0160209081015160f81c600081815260019092526040909120549091508015610ac75760ff821660009081526001602081905260409091204391610a8b908461136d565b81548110610a9b57610a9b6112dd565b90600052602060002001600001601c6101000a81548163ffffffff021916908363ffffffff1602179055505b60ff82166000908152600260209081526040918290208251808401909352805483526001015490820152610afb9086610dde565b60ff8316600090815260026020908152604080832084518155828501516001909101558051606081018252838152918201839052810191909152909450610b41856109d0565b67ffffffffffffffff1916815263ffffffff438116602080840191825260ff90951660009081526001808752604080832080548084018255908452979092208551970180549351958301518516600160e01b026001600160e01b0396909516600160c01b026001600160e01b03199094169790921c9690961791909117929092161790555001610a2a565b50505050565b60408051808201909152600080825260208201528151158015610bf757506020820151155b15610c15575050604080518082019091526000808252602082015290565b6040518060400160405280836000015181526020017f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478460200151610c5a9190611384565b610c84907f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4761136d565b905292915050565b919050565b816020015163ffffffff168163ffffffff161015610d2a5760405162461bcd60e51b815260206004820152604a60248201527f424c535075626b657952656769737472792e5f76616c696461746541706b486160448201527f7368466f7251756f72756d4174426c6f636b4e756d6265723a20696e646578206064820152691d1bdbc81c9958d95b9d60b21b608482015260a401610371565b604082015163ffffffff161580610d505750816040015163ffffffff168163ffffffff16105b610dda5760405162461bcd60e51b815260206004820152604f60248201527f424c535075626b657952656769737472792e5f76616c696461746541706b486160448201527f7368466f7251756f72756d4174426c6f636b4e756d6265723a206e6f74206c6160648201526e746573742061706b2075706461746560881b608482015260a401610371565b5050565b6040805180820190915260008082526020820152610dfa610e81565b835181526020808501518183015283516040808401919091529084015160608301526000908360808460066107d05a03fa9050808015610e3957610e3b565bfe5b5080610e795760405162461bcd60e51b815260206004820152600d60248201526c1958cb5859190b59985a5b1959609a1b6044820152606401610371565b505092915050565b60405180608001604052806004906020820280368337509192915050565b6001600160a01b0381168114610eb457600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff81118282101715610ef657610ef6610eb7565b604052919050565b600060408284031215610f1057600080fd5b6040516040810181811067ffffffffffffffff82111715610f3357610f33610eb7565b604052823581526020928301359281019290925250919050565b600080600060808486031215610f6257600080fd5b8335610f6d81610e9f565b925060208481013567ffffffffffffffff80821115610f8b57600080fd5b818701915087601f830112610f9f57600080fd5b813581811115610fb157610fb1610eb7565b610fc3601f8201601f19168501610ecd565b91508082528884828501011115610fd957600080fd5b8084840185840137600084828401015250809450505050610ffd8560408601610efe565b90509250925092565b803560ff81168114610c8c57600080fd5b60006020828403121561102957600080fd5b61103282611006565b9392505050565b60006020828403121561104b57600080fd5b5035919050565b6000806040838503121561106557600080fd5b61106e83611006565b946020939093013593505050565b60008060006060848603121561109157600080fd5b61109a84611006565b9250602084013563ffffffff811681146110b357600080fd5b929592945050506040919091013590565b6000806000604084860312156110d957600080fd5b833567ffffffffffffffff808211156110f157600080fd5b818601915086601f83011261110557600080fd5b81358181111561111457600080fd5b87602082850101111561112657600080fd5b6020928301989097509590910135949350505050565b6020808252825182820181905260009190848201906040850190845b8181101561117a57835163ffffffff1683529284019291840191600101611158565b50909695505050505050565b60208082526051908201527f424c535075626b657952656769737472792e6f6e6c795265676973747279436f60408201527f6f7264696e61746f723a2063616c6c6572206973206e6f74207468652072656760608201527034b9ba393c9031b7b7b93234b730ba37b960791b608082015260a00190565b602080825260409082018190527f424c535075626b657952656769737472792e72656769737465724f7065726174908201527f6f723a206f70657261746f7220646f6573206e6f74206f776e207075626b6579606082015260800190565b60018060a01b038316815260006020604081840152835180604085015260005b818110156112975785810183015185820160600152820161127b565b818111156112a9576000606083870101525b50601f01601f191692909201606001949350505050565b6000602082840312156112d257600080fd5b815161103281610e9f565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b600063ffffffff83811690831681811015611326576113266112f3565b039392505050565b600063ffffffff80831681811415611348576113486112f3565b6001019392505050565b6000600019821415611366576113666112f3565b5060010190565b60008282101561137f5761137f6112f3565b500390565b6000826113a157634e487b7160e01b600052601260045260246000fd5b50069056fea26469706673582212202b869eb5c17656b23a1f20a21ed234d4b252b6a8cf629742fe7019133b5e390f64736f6c634300080c0033",
}

// ContractBLSPubkeyRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractBLSPubkeyRegistryMetaData.ABI instead.
var ContractBLSPubkeyRegistryABI = ContractBLSPubkeyRegistryMetaData.ABI

// ContractBLSPubkeyRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractBLSPubkeyRegistryMetaData.Bin instead.
var ContractBLSPubkeyRegistryBin = ContractBLSPubkeyRegistryMetaData.Bin

// DeployContractBLSPubkeyRegistry deploys a new Ethereum contract, binding an instance of ContractBLSPubkeyRegistry to it.
func DeployContractBLSPubkeyRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, _registryCoordinator common.Address, _pubkeyCompendium common.Address) (common.Address, *types.Transaction, *ContractBLSPubkeyRegistry, error) {
	parsed, err := ContractBLSPubkeyRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractBLSPubkeyRegistryBin), backend, _registryCoordinator, _pubkeyCompendium)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractBLSPubkeyRegistry{ContractBLSPubkeyRegistryCaller: ContractBLSPubkeyRegistryCaller{contract: contract}, ContractBLSPubkeyRegistryTransactor: ContractBLSPubkeyRegistryTransactor{contract: contract}, ContractBLSPubkeyRegistryFilterer: ContractBLSPubkeyRegistryFilterer{contract: contract}}, nil
}

// ContractBLSPubkeyRegistry is an auto generated Go binding around an Ethereum contract.
type ContractBLSPubkeyRegistry struct {
	ContractBLSPubkeyRegistryCaller     // Read-only binding to the contract
	ContractBLSPubkeyRegistryTransactor // Write-only binding to the contract
	ContractBLSPubkeyRegistryFilterer   // Log filterer for contract events
}

// ContractBLSPubkeyRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractBLSPubkeyRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractBLSPubkeyRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractBLSPubkeyRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractBLSPubkeyRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractBLSPubkeyRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractBLSPubkeyRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractBLSPubkeyRegistrySession struct {
	Contract     *ContractBLSPubkeyRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts              // Call options to use throughout this session
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// ContractBLSPubkeyRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractBLSPubkeyRegistryCallerSession struct {
	Contract *ContractBLSPubkeyRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                    // Call options to use throughout this session
}

// ContractBLSPubkeyRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractBLSPubkeyRegistryTransactorSession struct {
	Contract     *ContractBLSPubkeyRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                    // Transaction auth options to use throughout this session
}

// ContractBLSPubkeyRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractBLSPubkeyRegistryRaw struct {
	Contract *ContractBLSPubkeyRegistry // Generic contract binding to access the raw methods on
}

// ContractBLSPubkeyRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractBLSPubkeyRegistryCallerRaw struct {
	Contract *ContractBLSPubkeyRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// ContractBLSPubkeyRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractBLSPubkeyRegistryTransactorRaw struct {
	Contract *ContractBLSPubkeyRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractBLSPubkeyRegistry creates a new instance of ContractBLSPubkeyRegistry, bound to a specific deployed contract.
func NewContractBLSPubkeyRegistry(address common.Address, backend bind.ContractBackend) (*ContractBLSPubkeyRegistry, error) {
	contract, err := bindContractBLSPubkeyRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractBLSPubkeyRegistry{ContractBLSPubkeyRegistryCaller: ContractBLSPubkeyRegistryCaller{contract: contract}, ContractBLSPubkeyRegistryTransactor: ContractBLSPubkeyRegistryTransactor{contract: contract}, ContractBLSPubkeyRegistryFilterer: ContractBLSPubkeyRegistryFilterer{contract: contract}}, nil
}

// NewContractBLSPubkeyRegistryCaller creates a new read-only instance of ContractBLSPubkeyRegistry, bound to a specific deployed contract.
func NewContractBLSPubkeyRegistryCaller(address common.Address, caller bind.ContractCaller) (*ContractBLSPubkeyRegistryCaller, error) {
	contract, err := bindContractBLSPubkeyRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractBLSPubkeyRegistryCaller{contract: contract}, nil
}

// NewContractBLSPubkeyRegistryTransactor creates a new write-only instance of ContractBLSPubkeyRegistry, bound to a specific deployed contract.
func NewContractBLSPubkeyRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractBLSPubkeyRegistryTransactor, error) {
	contract, err := bindContractBLSPubkeyRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractBLSPubkeyRegistryTransactor{contract: contract}, nil
}

// NewContractBLSPubkeyRegistryFilterer creates a new log filterer instance of ContractBLSPubkeyRegistry, bound to a specific deployed contract.
func NewContractBLSPubkeyRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractBLSPubkeyRegistryFilterer, error) {
	contract, err := bindContractBLSPubkeyRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractBLSPubkeyRegistryFilterer{contract: contract}, nil
}

// bindContractBLSPubkeyRegistry binds a generic wrapper to an already deployed contract.
func bindContractBLSPubkeyRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractBLSPubkeyRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractBLSPubkeyRegistry.Contract.ContractBLSPubkeyRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractBLSPubkeyRegistry.Contract.ContractBLSPubkeyRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractBLSPubkeyRegistry.Contract.ContractBLSPubkeyRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractBLSPubkeyRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractBLSPubkeyRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractBLSPubkeyRegistry.Contract.contract.Transact(opts, method, params...)
}

// GetApkForQuorum is a free data retrieval call binding the contract method 0x63a94510.
//
// Solidity: function getApkForQuorum(uint8 quorumNumber) view returns((uint256,uint256))
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryCaller) GetApkForQuorum(opts *bind.CallOpts, quorumNumber uint8) (BN254G1Point, error) {
	var out []interface{}
	err := _ContractBLSPubkeyRegistry.contract.Call(opts, &out, "getApkForQuorum", quorumNumber)

	if err != nil {
		return *new(BN254G1Point), err
	}

	out0 := *abi.ConvertType(out[0], new(BN254G1Point)).(*BN254G1Point)

	return out0, err

}

// GetApkForQuorum is a free data retrieval call binding the contract method 0x63a94510.
//
// Solidity: function getApkForQuorum(uint8 quorumNumber) view returns((uint256,uint256))
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistrySession) GetApkForQuorum(quorumNumber uint8) (BN254G1Point, error) {
	return _ContractBLSPubkeyRegistry.Contract.GetApkForQuorum(&_ContractBLSPubkeyRegistry.CallOpts, quorumNumber)
}

// GetApkForQuorum is a free data retrieval call binding the contract method 0x63a94510.
//
// Solidity: function getApkForQuorum(uint8 quorumNumber) view returns((uint256,uint256))
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryCallerSession) GetApkForQuorum(quorumNumber uint8) (BN254G1Point, error) {
	return _ContractBLSPubkeyRegistry.Contract.GetApkForQuorum(&_ContractBLSPubkeyRegistry.CallOpts, quorumNumber)
}

// GetApkHashForQuorumAtBlockNumberFromIndex is a free data retrieval call binding the contract method 0xc1af6b24.
//
// Solidity: function getApkHashForQuorumAtBlockNumberFromIndex(uint8 quorumNumber, uint32 blockNumber, uint256 index) view returns(bytes24)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryCaller) GetApkHashForQuorumAtBlockNumberFromIndex(opts *bind.CallOpts, quorumNumber uint8, blockNumber uint32, index *big.Int) ([24]byte, error) {
	var out []interface{}
	err := _ContractBLSPubkeyRegistry.contract.Call(opts, &out, "getApkHashForQuorumAtBlockNumberFromIndex", quorumNumber, blockNumber, index)

	if err != nil {
		return *new([24]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([24]byte)).(*[24]byte)

	return out0, err

}

// GetApkHashForQuorumAtBlockNumberFromIndex is a free data retrieval call binding the contract method 0xc1af6b24.
//
// Solidity: function getApkHashForQuorumAtBlockNumberFromIndex(uint8 quorumNumber, uint32 blockNumber, uint256 index) view returns(bytes24)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistrySession) GetApkHashForQuorumAtBlockNumberFromIndex(quorumNumber uint8, blockNumber uint32, index *big.Int) ([24]byte, error) {
	return _ContractBLSPubkeyRegistry.Contract.GetApkHashForQuorumAtBlockNumberFromIndex(&_ContractBLSPubkeyRegistry.CallOpts, quorumNumber, blockNumber, index)
}

// GetApkHashForQuorumAtBlockNumberFromIndex is a free data retrieval call binding the contract method 0xc1af6b24.
//
// Solidity: function getApkHashForQuorumAtBlockNumberFromIndex(uint8 quorumNumber, uint32 blockNumber, uint256 index) view returns(bytes24)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryCallerSession) GetApkHashForQuorumAtBlockNumberFromIndex(quorumNumber uint8, blockNumber uint32, index *big.Int) ([24]byte, error) {
	return _ContractBLSPubkeyRegistry.Contract.GetApkHashForQuorumAtBlockNumberFromIndex(&_ContractBLSPubkeyRegistry.CallOpts, quorumNumber, blockNumber, index)
}

// GetApkIndicesForQuorumsAtBlockNumber is a free data retrieval call binding the contract method 0xeda10763.
//
// Solidity: function getApkIndicesForQuorumsAtBlockNumber(bytes quorumNumbers, uint256 blockNumber) view returns(uint32[])
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryCaller) GetApkIndicesForQuorumsAtBlockNumber(opts *bind.CallOpts, quorumNumbers []byte, blockNumber *big.Int) ([]uint32, error) {
	var out []interface{}
	err := _ContractBLSPubkeyRegistry.contract.Call(opts, &out, "getApkIndicesForQuorumsAtBlockNumber", quorumNumbers, blockNumber)

	if err != nil {
		return *new([]uint32), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint32)).(*[]uint32)

	return out0, err

}

// GetApkIndicesForQuorumsAtBlockNumber is a free data retrieval call binding the contract method 0xeda10763.
//
// Solidity: function getApkIndicesForQuorumsAtBlockNumber(bytes quorumNumbers, uint256 blockNumber) view returns(uint32[])
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistrySession) GetApkIndicesForQuorumsAtBlockNumber(quorumNumbers []byte, blockNumber *big.Int) ([]uint32, error) {
	return _ContractBLSPubkeyRegistry.Contract.GetApkIndicesForQuorumsAtBlockNumber(&_ContractBLSPubkeyRegistry.CallOpts, quorumNumbers, blockNumber)
}

// GetApkIndicesForQuorumsAtBlockNumber is a free data retrieval call binding the contract method 0xeda10763.
//
// Solidity: function getApkIndicesForQuorumsAtBlockNumber(bytes quorumNumbers, uint256 blockNumber) view returns(uint32[])
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryCallerSession) GetApkIndicesForQuorumsAtBlockNumber(quorumNumbers []byte, blockNumber *big.Int) ([]uint32, error) {
	return _ContractBLSPubkeyRegistry.Contract.GetApkIndicesForQuorumsAtBlockNumber(&_ContractBLSPubkeyRegistry.CallOpts, quorumNumbers, blockNumber)
}

// GetApkUpdateForQuorumByIndex is a free data retrieval call binding the contract method 0x7225057e.
//
// Solidity: function getApkUpdateForQuorumByIndex(uint8 quorumNumber, uint256 index) view returns((bytes24,uint32,uint32))
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryCaller) GetApkUpdateForQuorumByIndex(opts *bind.CallOpts, quorumNumber uint8, index *big.Int) (IBLSPubkeyRegistryApkUpdate, error) {
	var out []interface{}
	err := _ContractBLSPubkeyRegistry.contract.Call(opts, &out, "getApkUpdateForQuorumByIndex", quorumNumber, index)

	if err != nil {
		return *new(IBLSPubkeyRegistryApkUpdate), err
	}

	out0 := *abi.ConvertType(out[0], new(IBLSPubkeyRegistryApkUpdate)).(*IBLSPubkeyRegistryApkUpdate)

	return out0, err

}

// GetApkUpdateForQuorumByIndex is a free data retrieval call binding the contract method 0x7225057e.
//
// Solidity: function getApkUpdateForQuorumByIndex(uint8 quorumNumber, uint256 index) view returns((bytes24,uint32,uint32))
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistrySession) GetApkUpdateForQuorumByIndex(quorumNumber uint8, index *big.Int) (IBLSPubkeyRegistryApkUpdate, error) {
	return _ContractBLSPubkeyRegistry.Contract.GetApkUpdateForQuorumByIndex(&_ContractBLSPubkeyRegistry.CallOpts, quorumNumber, index)
}

// GetApkUpdateForQuorumByIndex is a free data retrieval call binding the contract method 0x7225057e.
//
// Solidity: function getApkUpdateForQuorumByIndex(uint8 quorumNumber, uint256 index) view returns((bytes24,uint32,uint32))
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryCallerSession) GetApkUpdateForQuorumByIndex(quorumNumber uint8, index *big.Int) (IBLSPubkeyRegistryApkUpdate, error) {
	return _ContractBLSPubkeyRegistry.Contract.GetApkUpdateForQuorumByIndex(&_ContractBLSPubkeyRegistry.CallOpts, quorumNumber, index)
}

// GetOperatorFromPubkeyHash is a free data retrieval call binding the contract method 0x47b314e8.
//
// Solidity: function getOperatorFromPubkeyHash(bytes32 pubkeyHash) view returns(address)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryCaller) GetOperatorFromPubkeyHash(opts *bind.CallOpts, pubkeyHash [32]byte) (common.Address, error) {
	var out []interface{}
	err := _ContractBLSPubkeyRegistry.contract.Call(opts, &out, "getOperatorFromPubkeyHash", pubkeyHash)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOperatorFromPubkeyHash is a free data retrieval call binding the contract method 0x47b314e8.
//
// Solidity: function getOperatorFromPubkeyHash(bytes32 pubkeyHash) view returns(address)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistrySession) GetOperatorFromPubkeyHash(pubkeyHash [32]byte) (common.Address, error) {
	return _ContractBLSPubkeyRegistry.Contract.GetOperatorFromPubkeyHash(&_ContractBLSPubkeyRegistry.CallOpts, pubkeyHash)
}

// GetOperatorFromPubkeyHash is a free data retrieval call binding the contract method 0x47b314e8.
//
// Solidity: function getOperatorFromPubkeyHash(bytes32 pubkeyHash) view returns(address)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryCallerSession) GetOperatorFromPubkeyHash(pubkeyHash [32]byte) (common.Address, error) {
	return _ContractBLSPubkeyRegistry.Contract.GetOperatorFromPubkeyHash(&_ContractBLSPubkeyRegistry.CallOpts, pubkeyHash)
}

// GetQuorumApkHistoryLength is a free data retrieval call binding the contract method 0xfb81a7be.
//
// Solidity: function getQuorumApkHistoryLength(uint8 quorumNumber) view returns(uint32)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryCaller) GetQuorumApkHistoryLength(opts *bind.CallOpts, quorumNumber uint8) (uint32, error) {
	var out []interface{}
	err := _ContractBLSPubkeyRegistry.contract.Call(opts, &out, "getQuorumApkHistoryLength", quorumNumber)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// GetQuorumApkHistoryLength is a free data retrieval call binding the contract method 0xfb81a7be.
//
// Solidity: function getQuorumApkHistoryLength(uint8 quorumNumber) view returns(uint32)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistrySession) GetQuorumApkHistoryLength(quorumNumber uint8) (uint32, error) {
	return _ContractBLSPubkeyRegistry.Contract.GetQuorumApkHistoryLength(&_ContractBLSPubkeyRegistry.CallOpts, quorumNumber)
}

// GetQuorumApkHistoryLength is a free data retrieval call binding the contract method 0xfb81a7be.
//
// Solidity: function getQuorumApkHistoryLength(uint8 quorumNumber) view returns(uint32)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryCallerSession) GetQuorumApkHistoryLength(quorumNumber uint8) (uint32, error) {
	return _ContractBLSPubkeyRegistry.Contract.GetQuorumApkHistoryLength(&_ContractBLSPubkeyRegistry.CallOpts, quorumNumber)
}

// PubkeyCompendium is a free data retrieval call binding the contract method 0x187548c8.
//
// Solidity: function pubkeyCompendium() view returns(address)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryCaller) PubkeyCompendium(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractBLSPubkeyRegistry.contract.Call(opts, &out, "pubkeyCompendium")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PubkeyCompendium is a free data retrieval call binding the contract method 0x187548c8.
//
// Solidity: function pubkeyCompendium() view returns(address)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistrySession) PubkeyCompendium() (common.Address, error) {
	return _ContractBLSPubkeyRegistry.Contract.PubkeyCompendium(&_ContractBLSPubkeyRegistry.CallOpts)
}

// PubkeyCompendium is a free data retrieval call binding the contract method 0x187548c8.
//
// Solidity: function pubkeyCompendium() view returns(address)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryCallerSession) PubkeyCompendium() (common.Address, error) {
	return _ContractBLSPubkeyRegistry.Contract.PubkeyCompendium(&_ContractBLSPubkeyRegistry.CallOpts)
}

// QuorumApk is a free data retrieval call binding the contract method 0x32de6308.
//
// Solidity: function quorumApk(uint8 ) view returns(uint256 X, uint256 Y)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryCaller) QuorumApk(opts *bind.CallOpts, arg0 uint8) (struct {
	X *big.Int
	Y *big.Int
}, error) {
	var out []interface{}
	err := _ContractBLSPubkeyRegistry.contract.Call(opts, &out, "quorumApk", arg0)

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

// QuorumApk is a free data retrieval call binding the contract method 0x32de6308.
//
// Solidity: function quorumApk(uint8 ) view returns(uint256 X, uint256 Y)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistrySession) QuorumApk(arg0 uint8) (struct {
	X *big.Int
	Y *big.Int
}, error) {
	return _ContractBLSPubkeyRegistry.Contract.QuorumApk(&_ContractBLSPubkeyRegistry.CallOpts, arg0)
}

// QuorumApk is a free data retrieval call binding the contract method 0x32de6308.
//
// Solidity: function quorumApk(uint8 ) view returns(uint256 X, uint256 Y)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryCallerSession) QuorumApk(arg0 uint8) (struct {
	X *big.Int
	Y *big.Int
}, error) {
	return _ContractBLSPubkeyRegistry.Contract.QuorumApk(&_ContractBLSPubkeyRegistry.CallOpts, arg0)
}

// QuorumApkUpdates is a free data retrieval call binding the contract method 0x7f5eccbb.
//
// Solidity: function quorumApkUpdates(uint8 , uint256 ) view returns(bytes24 apkHash, uint32 updateBlockNumber, uint32 nextUpdateBlockNumber)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryCaller) QuorumApkUpdates(opts *bind.CallOpts, arg0 uint8, arg1 *big.Int) (struct {
	ApkHash               [24]byte
	UpdateBlockNumber     uint32
	NextUpdateBlockNumber uint32
}, error) {
	var out []interface{}
	err := _ContractBLSPubkeyRegistry.contract.Call(opts, &out, "quorumApkUpdates", arg0, arg1)

	outstruct := new(struct {
		ApkHash               [24]byte
		UpdateBlockNumber     uint32
		NextUpdateBlockNumber uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ApkHash = *abi.ConvertType(out[0], new([24]byte)).(*[24]byte)
	outstruct.UpdateBlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.NextUpdateBlockNumber = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

// QuorumApkUpdates is a free data retrieval call binding the contract method 0x7f5eccbb.
//
// Solidity: function quorumApkUpdates(uint8 , uint256 ) view returns(bytes24 apkHash, uint32 updateBlockNumber, uint32 nextUpdateBlockNumber)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistrySession) QuorumApkUpdates(arg0 uint8, arg1 *big.Int) (struct {
	ApkHash               [24]byte
	UpdateBlockNumber     uint32
	NextUpdateBlockNumber uint32
}, error) {
	return _ContractBLSPubkeyRegistry.Contract.QuorumApkUpdates(&_ContractBLSPubkeyRegistry.CallOpts, arg0, arg1)
}

// QuorumApkUpdates is a free data retrieval call binding the contract method 0x7f5eccbb.
//
// Solidity: function quorumApkUpdates(uint8 , uint256 ) view returns(bytes24 apkHash, uint32 updateBlockNumber, uint32 nextUpdateBlockNumber)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryCallerSession) QuorumApkUpdates(arg0 uint8, arg1 *big.Int) (struct {
	ApkHash               [24]byte
	UpdateBlockNumber     uint32
	NextUpdateBlockNumber uint32
}, error) {
	return _ContractBLSPubkeyRegistry.Contract.QuorumApkUpdates(&_ContractBLSPubkeyRegistry.CallOpts, arg0, arg1)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryCaller) RegistryCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractBLSPubkeyRegistry.contract.Call(opts, &out, "registryCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistrySession) RegistryCoordinator() (common.Address, error) {
	return _ContractBLSPubkeyRegistry.Contract.RegistryCoordinator(&_ContractBLSPubkeyRegistry.CallOpts)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryCallerSession) RegistryCoordinator() (common.Address, error) {
	return _ContractBLSPubkeyRegistry.Contract.RegistryCoordinator(&_ContractBLSPubkeyRegistry.CallOpts)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0x24369b2a.
//
// Solidity: function deregisterOperator(address operator, bytes quorumNumbers, (uint256,uint256) pubkey) returns()
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryTransactor) DeregisterOperator(opts *bind.TransactOpts, operator common.Address, quorumNumbers []byte, pubkey BN254G1Point) (*types.Transaction, error) {
	return _ContractBLSPubkeyRegistry.contract.Transact(opts, "deregisterOperator", operator, quorumNumbers, pubkey)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0x24369b2a.
//
// Solidity: function deregisterOperator(address operator, bytes quorumNumbers, (uint256,uint256) pubkey) returns()
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistrySession) DeregisterOperator(operator common.Address, quorumNumbers []byte, pubkey BN254G1Point) (*types.Transaction, error) {
	return _ContractBLSPubkeyRegistry.Contract.DeregisterOperator(&_ContractBLSPubkeyRegistry.TransactOpts, operator, quorumNumbers, pubkey)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0x24369b2a.
//
// Solidity: function deregisterOperator(address operator, bytes quorumNumbers, (uint256,uint256) pubkey) returns()
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryTransactorSession) DeregisterOperator(operator common.Address, quorumNumbers []byte, pubkey BN254G1Point) (*types.Transaction, error) {
	return _ContractBLSPubkeyRegistry.Contract.DeregisterOperator(&_ContractBLSPubkeyRegistry.TransactOpts, operator, quorumNumbers, pubkey)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x03ce4bad.
//
// Solidity: function registerOperator(address operator, bytes quorumNumbers, (uint256,uint256) pubkey) returns(bytes32)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryTransactor) RegisterOperator(opts *bind.TransactOpts, operator common.Address, quorumNumbers []byte, pubkey BN254G1Point) (*types.Transaction, error) {
	return _ContractBLSPubkeyRegistry.contract.Transact(opts, "registerOperator", operator, quorumNumbers, pubkey)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x03ce4bad.
//
// Solidity: function registerOperator(address operator, bytes quorumNumbers, (uint256,uint256) pubkey) returns(bytes32)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistrySession) RegisterOperator(operator common.Address, quorumNumbers []byte, pubkey BN254G1Point) (*types.Transaction, error) {
	return _ContractBLSPubkeyRegistry.Contract.RegisterOperator(&_ContractBLSPubkeyRegistry.TransactOpts, operator, quorumNumbers, pubkey)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x03ce4bad.
//
// Solidity: function registerOperator(address operator, bytes quorumNumbers, (uint256,uint256) pubkey) returns(bytes32)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryTransactorSession) RegisterOperator(operator common.Address, quorumNumbers []byte, pubkey BN254G1Point) (*types.Transaction, error) {
	return _ContractBLSPubkeyRegistry.Contract.RegisterOperator(&_ContractBLSPubkeyRegistry.TransactOpts, operator, quorumNumbers, pubkey)
}

// ContractBLSPubkeyRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ContractBLSPubkeyRegistry contract.
type ContractBLSPubkeyRegistryInitializedIterator struct {
	Event *ContractBLSPubkeyRegistryInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractBLSPubkeyRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractBLSPubkeyRegistryInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractBLSPubkeyRegistryInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractBLSPubkeyRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractBLSPubkeyRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractBLSPubkeyRegistryInitialized represents a Initialized event raised by the ContractBLSPubkeyRegistry contract.
type ContractBLSPubkeyRegistryInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*ContractBLSPubkeyRegistryInitializedIterator, error) {

	logs, sub, err := _ContractBLSPubkeyRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ContractBLSPubkeyRegistryInitializedIterator{contract: _ContractBLSPubkeyRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ContractBLSPubkeyRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _ContractBLSPubkeyRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractBLSPubkeyRegistryInitialized)
				if err := _ContractBLSPubkeyRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryFilterer) ParseInitialized(log types.Log) (*ContractBLSPubkeyRegistryInitialized, error) {
	event := new(ContractBLSPubkeyRegistryInitialized)
	if err := _ContractBLSPubkeyRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractBLSPubkeyRegistryOperatorAddedToQuorumsIterator is returned from FilterOperatorAddedToQuorums and is used to iterate over the raw logs and unpacked data for OperatorAddedToQuorums events raised by the ContractBLSPubkeyRegistry contract.
type ContractBLSPubkeyRegistryOperatorAddedToQuorumsIterator struct {
	Event *ContractBLSPubkeyRegistryOperatorAddedToQuorums // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractBLSPubkeyRegistryOperatorAddedToQuorumsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractBLSPubkeyRegistryOperatorAddedToQuorums)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractBLSPubkeyRegistryOperatorAddedToQuorums)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractBLSPubkeyRegistryOperatorAddedToQuorumsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractBLSPubkeyRegistryOperatorAddedToQuorumsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractBLSPubkeyRegistryOperatorAddedToQuorums represents a OperatorAddedToQuorums event raised by the ContractBLSPubkeyRegistry contract.
type ContractBLSPubkeyRegistryOperatorAddedToQuorums struct {
	Operator      common.Address
	QuorumNumbers []byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOperatorAddedToQuorums is a free log retrieval operation binding the contract event 0x5358c5b42179178c8fc757734ac2a3198f9071c765ee0d8389211525f5005246.
//
// Solidity: event OperatorAddedToQuorums(address operator, bytes quorumNumbers)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryFilterer) FilterOperatorAddedToQuorums(opts *bind.FilterOpts) (*ContractBLSPubkeyRegistryOperatorAddedToQuorumsIterator, error) {

	logs, sub, err := _ContractBLSPubkeyRegistry.contract.FilterLogs(opts, "OperatorAddedToQuorums")
	if err != nil {
		return nil, err
	}
	return &ContractBLSPubkeyRegistryOperatorAddedToQuorumsIterator{contract: _ContractBLSPubkeyRegistry.contract, event: "OperatorAddedToQuorums", logs: logs, sub: sub}, nil
}

// WatchOperatorAddedToQuorums is a free log subscription operation binding the contract event 0x5358c5b42179178c8fc757734ac2a3198f9071c765ee0d8389211525f5005246.
//
// Solidity: event OperatorAddedToQuorums(address operator, bytes quorumNumbers)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryFilterer) WatchOperatorAddedToQuorums(opts *bind.WatchOpts, sink chan<- *ContractBLSPubkeyRegistryOperatorAddedToQuorums) (event.Subscription, error) {

	logs, sub, err := _ContractBLSPubkeyRegistry.contract.WatchLogs(opts, "OperatorAddedToQuorums")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractBLSPubkeyRegistryOperatorAddedToQuorums)
				if err := _ContractBLSPubkeyRegistry.contract.UnpackLog(event, "OperatorAddedToQuorums", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOperatorAddedToQuorums is a log parse operation binding the contract event 0x5358c5b42179178c8fc757734ac2a3198f9071c765ee0d8389211525f5005246.
//
// Solidity: event OperatorAddedToQuorums(address operator, bytes quorumNumbers)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryFilterer) ParseOperatorAddedToQuorums(log types.Log) (*ContractBLSPubkeyRegistryOperatorAddedToQuorums, error) {
	event := new(ContractBLSPubkeyRegistryOperatorAddedToQuorums)
	if err := _ContractBLSPubkeyRegistry.contract.UnpackLog(event, "OperatorAddedToQuorums", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractBLSPubkeyRegistryOperatorRemovedFromQuorumsIterator is returned from FilterOperatorRemovedFromQuorums and is used to iterate over the raw logs and unpacked data for OperatorRemovedFromQuorums events raised by the ContractBLSPubkeyRegistry contract.
type ContractBLSPubkeyRegistryOperatorRemovedFromQuorumsIterator struct {
	Event *ContractBLSPubkeyRegistryOperatorRemovedFromQuorums // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractBLSPubkeyRegistryOperatorRemovedFromQuorumsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractBLSPubkeyRegistryOperatorRemovedFromQuorums)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractBLSPubkeyRegistryOperatorRemovedFromQuorums)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractBLSPubkeyRegistryOperatorRemovedFromQuorumsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractBLSPubkeyRegistryOperatorRemovedFromQuorumsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractBLSPubkeyRegistryOperatorRemovedFromQuorums represents a OperatorRemovedFromQuorums event raised by the ContractBLSPubkeyRegistry contract.
type ContractBLSPubkeyRegistryOperatorRemovedFromQuorums struct {
	Operator      common.Address
	QuorumNumbers []byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOperatorRemovedFromQuorums is a free log retrieval operation binding the contract event 0x14a5172b312e9d2c22b8468f9c70ec2caa9de934fe380734fbc6f3beff2b14ba.
//
// Solidity: event OperatorRemovedFromQuorums(address operator, bytes quorumNumbers)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryFilterer) FilterOperatorRemovedFromQuorums(opts *bind.FilterOpts) (*ContractBLSPubkeyRegistryOperatorRemovedFromQuorumsIterator, error) {

	logs, sub, err := _ContractBLSPubkeyRegistry.contract.FilterLogs(opts, "OperatorRemovedFromQuorums")
	if err != nil {
		return nil, err
	}
	return &ContractBLSPubkeyRegistryOperatorRemovedFromQuorumsIterator{contract: _ContractBLSPubkeyRegistry.contract, event: "OperatorRemovedFromQuorums", logs: logs, sub: sub}, nil
}

// WatchOperatorRemovedFromQuorums is a free log subscription operation binding the contract event 0x14a5172b312e9d2c22b8468f9c70ec2caa9de934fe380734fbc6f3beff2b14ba.
//
// Solidity: event OperatorRemovedFromQuorums(address operator, bytes quorumNumbers)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryFilterer) WatchOperatorRemovedFromQuorums(opts *bind.WatchOpts, sink chan<- *ContractBLSPubkeyRegistryOperatorRemovedFromQuorums) (event.Subscription, error) {

	logs, sub, err := _ContractBLSPubkeyRegistry.contract.WatchLogs(opts, "OperatorRemovedFromQuorums")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractBLSPubkeyRegistryOperatorRemovedFromQuorums)
				if err := _ContractBLSPubkeyRegistry.contract.UnpackLog(event, "OperatorRemovedFromQuorums", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOperatorRemovedFromQuorums is a log parse operation binding the contract event 0x14a5172b312e9d2c22b8468f9c70ec2caa9de934fe380734fbc6f3beff2b14ba.
//
// Solidity: event OperatorRemovedFromQuorums(address operator, bytes quorumNumbers)
func (_ContractBLSPubkeyRegistry *ContractBLSPubkeyRegistryFilterer) ParseOperatorRemovedFromQuorums(log types.Log) (*ContractBLSPubkeyRegistryOperatorRemovedFromQuorums, error) {
	event := new(ContractBLSPubkeyRegistryOperatorRemovedFromQuorums)
	if err := _ContractBLSPubkeyRegistry.contract.UnpackLog(event, "OperatorRemovedFromQuorums", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
