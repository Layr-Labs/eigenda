// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractBLSOperatorStateRetriever

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

// BLSOperatorStateRetrieverCheckSignaturesIndices is an auto generated low-level Go binding around an user-defined struct.
type BLSOperatorStateRetrieverCheckSignaturesIndices struct {
	NonSignerQuorumBitmapIndices []uint32
	QuorumApkIndices             []uint32
	TotalStakeIndices            []uint32
	NonSignerStakeIndices        [][]uint32
}

// BLSOperatorStateRetrieverOperator is an auto generated low-level Go binding around an user-defined struct.
type BLSOperatorStateRetrieverOperator struct {
	OperatorId [32]byte
	Stake      *big.Int
}

// ContractBLSOperatorStateRetrieverMetaData contains all meta data concerning the ContractBLSOperatorStateRetriever contract.
var ContractBLSOperatorStateRetrieverMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"getCheckSignaturesIndices\",\"inputs\":[{\"name\":\"registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIBLSRegistryCoordinatorWithIndices\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"nonSignerOperatorIds\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structBLSOperatorStateRetriever.CheckSignaturesIndices\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorState\",\"inputs\":[{\"name\":\"registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIBLSRegistryCoordinatorWithIndices\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"blockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[][]\",\"internalType\":\"structBLSOperatorStateRetriever.Operator[][]\",\"components\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"stake\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorState\",\"inputs\":[{\"name\":\"registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIBLSRegistryCoordinatorWithIndices\"},{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"blockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"tuple[][]\",\"internalType\":\"structBLSOperatorStateRetriever.Operator[][]\",\"components\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"stake\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"stateMutability\":\"view\"}]",
	Bin: "0x608060405234801561001057600080fd5b506114cc806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80633563b0d1146100465780634f739f741461006f578063cefdc1d41461008f575b600080fd5b610059610054366004610d22565b6100b0565b6040516100669190610e6a565b60405180910390f35b61008261007d366004610ed0565b61042b565b6040516100669190610fd4565b6100a261009d36600461108f565b610ab2565b6040516100669291906110d1565b60606000846001600160a01b031663683048356040518163ffffffff1660e01b8152600401602060405180830381865afa1580156100f2573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061011691906110f2565b90506000856001600160a01b0316639e9923c26040518163ffffffff1660e01b8152600401602060405180830381865afa158015610158573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061017c91906110f2565b90506000855167ffffffffffffffff81111561019a5761019a610cb9565b6040519080825280602002602001820160405280156101cd57816020015b60608152602001906001900390816101b85790505b50905060005b86518110156104205760008782815181106101f0576101f061110f565b016020015160405163889ae3e560e01b815260f89190911c6004820181905263ffffffff8916602483015291506000906001600160a01b0386169063889ae3e590604401600060405180830381865afa158015610251573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526102799190810190611149565b9050805167ffffffffffffffff81111561029557610295610cb9565b6040519080825280602002602001820160405280156102da57816020015b60408051808201909152600080825260208201528152602001906001900390816102b35790505b508484815181106102ed576102ed61110f565b602002602001018190525060005b815181101561040a5760008282815181106103185761031861110f565b6020908102919091018101516040805180820182528281529051631b32722560e01b81526004810183905260ff8816602482015263ffffffff8e166044820152919350918201906001600160a01b038b1690631b32722590606401602060405180830381865afa158015610390573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103b491906111df565b6001600160601b03168152508686815181106103d2576103d261110f565b602002602001015183815181106103eb576103eb61110f565b60200260200101819052505080806104029061121e565b9150506102fb565b50505080806104189061121e565b9150506101d3565b509695505050505050565b6104566040518060800160405280606081526020016060815260200160608152602001606081525090565b6000876001600160a01b031663683048356040518163ffffffff1660e01b8152600401602060405180830381865afa158015610496573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104ba91906110f2565b90506104e76040518060800160405280606081526020016060815260200160608152602001606081525090565b6040516385020d4960e01b81526001600160a01b038a16906385020d4990610517908b9089908990600401611239565b600060405180830381865afa158015610534573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261055c9190810190611283565b815260405163e192e9ad60e01b81526001600160a01b0383169063e192e9ad9061058e908b908b908b9060040161133b565b600060405180830381865afa1580156105ab573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526105d39190810190611283565b60408201528567ffffffffffffffff8111156105f1576105f1610cb9565b60405190808252806020026020018201604052801561062457816020015b606081526020019060019003908161060f5790505b50606082015260005b60ff81168711156109c35760008567ffffffffffffffff81111561065357610653610cb9565b60405190808252806020026020018201604052801561067c578160200160208202803683370190505b5083606001518360ff16815181106106965761069661110f565b602002602001018190525060005b868110156108c25760008c6001600160a01b0316633064620d8a8a858181106106cf576106cf61110f565b905060200201358e886000015186815181106106ed576106ed61110f565b60200260200101516040518463ffffffff1660e01b815260040161072a9392919092835263ffffffff918216602084015216604082015260600190565b602060405180830381865afa158015610747573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061076b9190611364565b90508a8a8560ff168181106107825761078261110f565b6001600160c01b03841692013560f81c9190911c6001908116141590506108af57856001600160a01b031663480858668a8a858181106107c4576107c461110f565b905060200201358d8d8860ff168181106107e0576107e061110f565b6040516001600160e01b031960e087901b1681526004810194909452919091013560f81c60248301525063ffffffff8f166044820152606401602060405180830381865afa158015610836573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061085a919061138d565b85606001518560ff16815181106108735761087361110f565b6020026020010151848151811061088c5761088c61110f565b63ffffffff90921660209283029190910190910152826108ab8161121e565b9350505b50806108ba8161121e565b9150506106a4565b5060008167ffffffffffffffff8111156108de576108de610cb9565b604051908082528060200260200182016040528015610907578160200160208202803683370190505b50905060005b828110156109885784606001518460ff168151811061092e5761092e61110f565b602002602001015181815181106109475761094761110f565b60200260200101518282815181106109615761096161110f565b63ffffffff90921660209283029190910190910152806109808161121e565b91505061090d565b508084606001518460ff16815181106109a3576109a361110f565b6020026020010181905250505080806109bb906113aa565b91505061062d565b506000896001600160a01b0316633561deb16040518163ffffffff1660e01b8152600401602060405180830381865afa158015610a04573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a2891906110f2565b60405163eda1076360e01b81529091506001600160a01b0382169063eda1076390610a5b908b908b908e906004016113ca565b600060405180830381865afa158015610a78573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610aa09190810190611283565b60208301525098975050505050505050565b6040805160018082528183019092526000916060918391602080830190803683370190505090508481600081518110610aed57610aed61110f565b60209081029190910101526040516385020d4960e01b81526000906001600160a01b038816906385020d4990610b2990889086906004016113f4565b600060405180830381865afa158015610b46573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610b6e9190810190611283565b600081518110610b8057610b8061110f565b6020908102919091010151604051633064620d60e01b81526004810188905263ffffffff87811660248301529091166044820181905291506000906001600160a01b03891690633064620d90606401602060405180830381865afa158015610bec573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c109190611364565b6001600160c01b031690506000610c2682610c44565b905081610c348a838a6100b0565b9550955050505050935093915050565b60606000805b610100811015610c9a576001811b915083821615610c8a57828160f81b604051602001610c78929190611448565b60405160208183030381529060405292505b610c938161121e565b9050610c4a565b5050919050565b6001600160a01b0381168114610cb657600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff81118282101715610cf857610cf8610cb9565b604052919050565b63ffffffff81168114610cb657600080fd5b8035610d1d81610d00565b919050565b600080600060608486031215610d3757600080fd5b8335610d4281610ca1565b925060208481013567ffffffffffffffff80821115610d6057600080fd5b818701915087601f830112610d7457600080fd5b813581811115610d8657610d86610cb9565b610d98601f8201601f19168501610ccf565b91508082528884828501011115610dae57600080fd5b8084840185840137600084828401015250809450505050610dd160408501610d12565b90509250925092565b600081518084526020808501808196508360051b810191508286016000805b86811015610e5c578385038a52825180518087529087019087870190845b81811015610e47578351805184528a01516001600160601b03168a84015292890192604090920191600101610e17565b50509a87019a95505091850191600101610df9565b509298975050505050505050565b602081526000610e7d6020830184610dda565b9392505050565b60008083601f840112610e9657600080fd5b50813567ffffffffffffffff811115610eae57600080fd5b6020830191508360208260051b8501011115610ec957600080fd5b9250929050565b60008060008060008060808789031215610ee957600080fd5b8635610ef481610ca1565b95506020870135610f0481610d00565b9450604087013567ffffffffffffffff80821115610f2157600080fd5b818901915089601f830112610f3557600080fd5b813581811115610f4457600080fd5b8a6020828501011115610f5657600080fd5b602083019650809550506060890135915080821115610f7457600080fd5b50610f8189828a01610e84565b979a9699509497509295939492505050565b600081518084526020808501945080840160005b83811015610fc957815163ffffffff1687529582019590820190600101610fa7565b509495945050505050565b600060208083528351608082850152610ff060a0850182610f93565b905081850151601f198086840301604087015261100d8383610f93565b9250604087015191508086840301606087015261102a8383610f93565b60608801518782038301608089015280518083529194508501925084840190600581901b8501860160005b82811015611081578487830301845261106f828751610f93565b95880195938801939150600101611055565b509998505050505050505050565b6000806000606084860312156110a457600080fd5b83356110af81610ca1565b92506020840135915060408401356110c681610d00565b809150509250925092565b8281526040602082015260006110ea6040830184610dda565b949350505050565b60006020828403121561110457600080fd5b8151610e7d81610ca1565b634e487b7160e01b600052603260045260246000fd5b600067ffffffffffffffff82111561113f5761113f610cb9565b5060051b60200190565b6000602080838503121561115c57600080fd5b825167ffffffffffffffff81111561117357600080fd5b8301601f8101851361118457600080fd5b805161119761119282611125565b610ccf565b81815260059190911b820183019083810190878311156111b657600080fd5b928401925b828410156111d4578351825292840192908401906111bb565b979650505050505050565b6000602082840312156111f157600080fd5b81516001600160601b0381168114610e7d57600080fd5b634e487b7160e01b600052601160045260246000fd5b600060001982141561123257611232611208565b5060010190565b63ffffffff84168152604060208201819052810182905260006001600160fb1b0383111561126657600080fd5b8260051b8085606085013760009201606001918252509392505050565b6000602080838503121561129657600080fd5b825167ffffffffffffffff8111156112ad57600080fd5b8301601f810185136112be57600080fd5b80516112cc61119282611125565b81815260059190911b820183019083810190878311156112eb57600080fd5b928401925b828410156111d457835161130381610d00565b825292840192908401906112f0565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b63ffffffff8416815260406020820152600061135b604083018486611312565b95945050505050565b60006020828403121561137657600080fd5b81516001600160c01b0381168114610e7d57600080fd5b60006020828403121561139f57600080fd5b8151610e7d81610d00565b600060ff821660ff8114156113c1576113c1611208565b60010192915050565b6040815260006113de604083018587611312565b905063ffffffff83166020830152949350505050565b60006040820163ffffffff851683526020604081850152818551808452606086019150828701935060005b8181101561143b5784518352938301939183019160010161141f565b5090979650505050505050565b6000835160005b81811015611469576020818701810151858301520161144f565b81811115611478576000828501525b506001600160f81b031993909316919092019081526001019291505056fea2646970667358221220bb7598eb8c072001f48565267ff944a12ddc3fb3ea780fbc65699e9c973fa08964736f6c634300080c0033",
}

// ContractBLSOperatorStateRetrieverABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractBLSOperatorStateRetrieverMetaData.ABI instead.
var ContractBLSOperatorStateRetrieverABI = ContractBLSOperatorStateRetrieverMetaData.ABI

// ContractBLSOperatorStateRetrieverBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractBLSOperatorStateRetrieverMetaData.Bin instead.
var ContractBLSOperatorStateRetrieverBin = ContractBLSOperatorStateRetrieverMetaData.Bin

// DeployContractBLSOperatorStateRetriever deploys a new Ethereum contract, binding an instance of ContractBLSOperatorStateRetriever to it.
func DeployContractBLSOperatorStateRetriever(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ContractBLSOperatorStateRetriever, error) {
	parsed, err := ContractBLSOperatorStateRetrieverMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractBLSOperatorStateRetrieverBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractBLSOperatorStateRetriever{ContractBLSOperatorStateRetrieverCaller: ContractBLSOperatorStateRetrieverCaller{contract: contract}, ContractBLSOperatorStateRetrieverTransactor: ContractBLSOperatorStateRetrieverTransactor{contract: contract}, ContractBLSOperatorStateRetrieverFilterer: ContractBLSOperatorStateRetrieverFilterer{contract: contract}}, nil
}

// ContractBLSOperatorStateRetriever is an auto generated Go binding around an Ethereum contract.
type ContractBLSOperatorStateRetriever struct {
	ContractBLSOperatorStateRetrieverCaller     // Read-only binding to the contract
	ContractBLSOperatorStateRetrieverTransactor // Write-only binding to the contract
	ContractBLSOperatorStateRetrieverFilterer   // Log filterer for contract events
}

// ContractBLSOperatorStateRetrieverCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractBLSOperatorStateRetrieverCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractBLSOperatorStateRetrieverTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractBLSOperatorStateRetrieverTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractBLSOperatorStateRetrieverFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractBLSOperatorStateRetrieverFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractBLSOperatorStateRetrieverSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractBLSOperatorStateRetrieverSession struct {
	Contract     *ContractBLSOperatorStateRetriever // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                      // Call options to use throughout this session
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// ContractBLSOperatorStateRetrieverCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractBLSOperatorStateRetrieverCallerSession struct {
	Contract *ContractBLSOperatorStateRetrieverCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                            // Call options to use throughout this session
}

// ContractBLSOperatorStateRetrieverTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractBLSOperatorStateRetrieverTransactorSession struct {
	Contract     *ContractBLSOperatorStateRetrieverTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                            // Transaction auth options to use throughout this session
}

// ContractBLSOperatorStateRetrieverRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractBLSOperatorStateRetrieverRaw struct {
	Contract *ContractBLSOperatorStateRetriever // Generic contract binding to access the raw methods on
}

// ContractBLSOperatorStateRetrieverCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractBLSOperatorStateRetrieverCallerRaw struct {
	Contract *ContractBLSOperatorStateRetrieverCaller // Generic read-only contract binding to access the raw methods on
}

// ContractBLSOperatorStateRetrieverTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractBLSOperatorStateRetrieverTransactorRaw struct {
	Contract *ContractBLSOperatorStateRetrieverTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractBLSOperatorStateRetriever creates a new instance of ContractBLSOperatorStateRetriever, bound to a specific deployed contract.
func NewContractBLSOperatorStateRetriever(address common.Address, backend bind.ContractBackend) (*ContractBLSOperatorStateRetriever, error) {
	contract, err := bindContractBLSOperatorStateRetriever(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractBLSOperatorStateRetriever{ContractBLSOperatorStateRetrieverCaller: ContractBLSOperatorStateRetrieverCaller{contract: contract}, ContractBLSOperatorStateRetrieverTransactor: ContractBLSOperatorStateRetrieverTransactor{contract: contract}, ContractBLSOperatorStateRetrieverFilterer: ContractBLSOperatorStateRetrieverFilterer{contract: contract}}, nil
}

// NewContractBLSOperatorStateRetrieverCaller creates a new read-only instance of ContractBLSOperatorStateRetriever, bound to a specific deployed contract.
func NewContractBLSOperatorStateRetrieverCaller(address common.Address, caller bind.ContractCaller) (*ContractBLSOperatorStateRetrieverCaller, error) {
	contract, err := bindContractBLSOperatorStateRetriever(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractBLSOperatorStateRetrieverCaller{contract: contract}, nil
}

// NewContractBLSOperatorStateRetrieverTransactor creates a new write-only instance of ContractBLSOperatorStateRetriever, bound to a specific deployed contract.
func NewContractBLSOperatorStateRetrieverTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractBLSOperatorStateRetrieverTransactor, error) {
	contract, err := bindContractBLSOperatorStateRetriever(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractBLSOperatorStateRetrieverTransactor{contract: contract}, nil
}

// NewContractBLSOperatorStateRetrieverFilterer creates a new log filterer instance of ContractBLSOperatorStateRetriever, bound to a specific deployed contract.
func NewContractBLSOperatorStateRetrieverFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractBLSOperatorStateRetrieverFilterer, error) {
	contract, err := bindContractBLSOperatorStateRetriever(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractBLSOperatorStateRetrieverFilterer{contract: contract}, nil
}

// bindContractBLSOperatorStateRetriever binds a generic wrapper to an already deployed contract.
func bindContractBLSOperatorStateRetriever(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractBLSOperatorStateRetrieverMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractBLSOperatorStateRetriever *ContractBLSOperatorStateRetrieverRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractBLSOperatorStateRetriever.Contract.ContractBLSOperatorStateRetrieverCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractBLSOperatorStateRetriever *ContractBLSOperatorStateRetrieverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractBLSOperatorStateRetriever.Contract.ContractBLSOperatorStateRetrieverTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractBLSOperatorStateRetriever *ContractBLSOperatorStateRetrieverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractBLSOperatorStateRetriever.Contract.ContractBLSOperatorStateRetrieverTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractBLSOperatorStateRetriever *ContractBLSOperatorStateRetrieverCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractBLSOperatorStateRetriever.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractBLSOperatorStateRetriever *ContractBLSOperatorStateRetrieverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractBLSOperatorStateRetriever.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractBLSOperatorStateRetriever *ContractBLSOperatorStateRetrieverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractBLSOperatorStateRetriever.Contract.contract.Transact(opts, method, params...)
}

// GetCheckSignaturesIndices is a free data retrieval call binding the contract method 0x4f739f74.
//
// Solidity: function getCheckSignaturesIndices(address registryCoordinator, uint32 referenceBlockNumber, bytes quorumNumbers, bytes32[] nonSignerOperatorIds) view returns((uint32[],uint32[],uint32[],uint32[][]))
func (_ContractBLSOperatorStateRetriever *ContractBLSOperatorStateRetrieverCaller) GetCheckSignaturesIndices(opts *bind.CallOpts, registryCoordinator common.Address, referenceBlockNumber uint32, quorumNumbers []byte, nonSignerOperatorIds [][32]byte) (BLSOperatorStateRetrieverCheckSignaturesIndices, error) {
	var out []interface{}
	err := _ContractBLSOperatorStateRetriever.contract.Call(opts, &out, "getCheckSignaturesIndices", registryCoordinator, referenceBlockNumber, quorumNumbers, nonSignerOperatorIds)

	if err != nil {
		return *new(BLSOperatorStateRetrieverCheckSignaturesIndices), err
	}

	out0 := *abi.ConvertType(out[0], new(BLSOperatorStateRetrieverCheckSignaturesIndices)).(*BLSOperatorStateRetrieverCheckSignaturesIndices)

	return out0, err

}

// GetCheckSignaturesIndices is a free data retrieval call binding the contract method 0x4f739f74.
//
// Solidity: function getCheckSignaturesIndices(address registryCoordinator, uint32 referenceBlockNumber, bytes quorumNumbers, bytes32[] nonSignerOperatorIds) view returns((uint32[],uint32[],uint32[],uint32[][]))
func (_ContractBLSOperatorStateRetriever *ContractBLSOperatorStateRetrieverSession) GetCheckSignaturesIndices(registryCoordinator common.Address, referenceBlockNumber uint32, quorumNumbers []byte, nonSignerOperatorIds [][32]byte) (BLSOperatorStateRetrieverCheckSignaturesIndices, error) {
	return _ContractBLSOperatorStateRetriever.Contract.GetCheckSignaturesIndices(&_ContractBLSOperatorStateRetriever.CallOpts, registryCoordinator, referenceBlockNumber, quorumNumbers, nonSignerOperatorIds)
}

// GetCheckSignaturesIndices is a free data retrieval call binding the contract method 0x4f739f74.
//
// Solidity: function getCheckSignaturesIndices(address registryCoordinator, uint32 referenceBlockNumber, bytes quorumNumbers, bytes32[] nonSignerOperatorIds) view returns((uint32[],uint32[],uint32[],uint32[][]))
func (_ContractBLSOperatorStateRetriever *ContractBLSOperatorStateRetrieverCallerSession) GetCheckSignaturesIndices(registryCoordinator common.Address, referenceBlockNumber uint32, quorumNumbers []byte, nonSignerOperatorIds [][32]byte) (BLSOperatorStateRetrieverCheckSignaturesIndices, error) {
	return _ContractBLSOperatorStateRetriever.Contract.GetCheckSignaturesIndices(&_ContractBLSOperatorStateRetriever.CallOpts, registryCoordinator, referenceBlockNumber, quorumNumbers, nonSignerOperatorIds)
}

// GetOperatorState is a free data retrieval call binding the contract method 0x3563b0d1.
//
// Solidity: function getOperatorState(address registryCoordinator, bytes quorumNumbers, uint32 blockNumber) view returns((bytes32,uint96)[][])
func (_ContractBLSOperatorStateRetriever *ContractBLSOperatorStateRetrieverCaller) GetOperatorState(opts *bind.CallOpts, registryCoordinator common.Address, quorumNumbers []byte, blockNumber uint32) ([][]BLSOperatorStateRetrieverOperator, error) {
	var out []interface{}
	err := _ContractBLSOperatorStateRetriever.contract.Call(opts, &out, "getOperatorState", registryCoordinator, quorumNumbers, blockNumber)

	if err != nil {
		return *new([][]BLSOperatorStateRetrieverOperator), err
	}

	out0 := *abi.ConvertType(out[0], new([][]BLSOperatorStateRetrieverOperator)).(*[][]BLSOperatorStateRetrieverOperator)

	return out0, err

}

// GetOperatorState is a free data retrieval call binding the contract method 0x3563b0d1.
//
// Solidity: function getOperatorState(address registryCoordinator, bytes quorumNumbers, uint32 blockNumber) view returns((bytes32,uint96)[][])
func (_ContractBLSOperatorStateRetriever *ContractBLSOperatorStateRetrieverSession) GetOperatorState(registryCoordinator common.Address, quorumNumbers []byte, blockNumber uint32) ([][]BLSOperatorStateRetrieverOperator, error) {
	return _ContractBLSOperatorStateRetriever.Contract.GetOperatorState(&_ContractBLSOperatorStateRetriever.CallOpts, registryCoordinator, quorumNumbers, blockNumber)
}

// GetOperatorState is a free data retrieval call binding the contract method 0x3563b0d1.
//
// Solidity: function getOperatorState(address registryCoordinator, bytes quorumNumbers, uint32 blockNumber) view returns((bytes32,uint96)[][])
func (_ContractBLSOperatorStateRetriever *ContractBLSOperatorStateRetrieverCallerSession) GetOperatorState(registryCoordinator common.Address, quorumNumbers []byte, blockNumber uint32) ([][]BLSOperatorStateRetrieverOperator, error) {
	return _ContractBLSOperatorStateRetriever.Contract.GetOperatorState(&_ContractBLSOperatorStateRetriever.CallOpts, registryCoordinator, quorumNumbers, blockNumber)
}

// GetOperatorState0 is a free data retrieval call binding the contract method 0xcefdc1d4.
//
// Solidity: function getOperatorState(address registryCoordinator, bytes32 operatorId, uint32 blockNumber) view returns(uint256, (bytes32,uint96)[][])
func (_ContractBLSOperatorStateRetriever *ContractBLSOperatorStateRetrieverCaller) GetOperatorState0(opts *bind.CallOpts, registryCoordinator common.Address, operatorId [32]byte, blockNumber uint32) (*big.Int, [][]BLSOperatorStateRetrieverOperator, error) {
	var out []interface{}
	err := _ContractBLSOperatorStateRetriever.contract.Call(opts, &out, "getOperatorState0", registryCoordinator, operatorId, blockNumber)

	if err != nil {
		return *new(*big.Int), *new([][]BLSOperatorStateRetrieverOperator), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new([][]BLSOperatorStateRetrieverOperator)).(*[][]BLSOperatorStateRetrieverOperator)

	return out0, out1, err

}

// GetOperatorState0 is a free data retrieval call binding the contract method 0xcefdc1d4.
//
// Solidity: function getOperatorState(address registryCoordinator, bytes32 operatorId, uint32 blockNumber) view returns(uint256, (bytes32,uint96)[][])
func (_ContractBLSOperatorStateRetriever *ContractBLSOperatorStateRetrieverSession) GetOperatorState0(registryCoordinator common.Address, operatorId [32]byte, blockNumber uint32) (*big.Int, [][]BLSOperatorStateRetrieverOperator, error) {
	return _ContractBLSOperatorStateRetriever.Contract.GetOperatorState0(&_ContractBLSOperatorStateRetriever.CallOpts, registryCoordinator, operatorId, blockNumber)
}

// GetOperatorState0 is a free data retrieval call binding the contract method 0xcefdc1d4.
//
// Solidity: function getOperatorState(address registryCoordinator, bytes32 operatorId, uint32 blockNumber) view returns(uint256, (bytes32,uint96)[][])
func (_ContractBLSOperatorStateRetriever *ContractBLSOperatorStateRetrieverCallerSession) GetOperatorState0(registryCoordinator common.Address, operatorId [32]byte, blockNumber uint32) (*big.Int, [][]BLSOperatorStateRetrieverOperator, error) {
	return _ContractBLSOperatorStateRetriever.Contract.GetOperatorState0(&_ContractBLSOperatorStateRetriever.CallOpts, registryCoordinator, operatorId, blockNumber)
}
