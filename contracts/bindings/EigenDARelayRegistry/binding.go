// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractEigenDARelayRegistry

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

// RelayInfo is an auto generated low-level Go binding around an user-defined struct.
type RelayInfo struct {
	RelayAddress common.Address
	RelayURL     string
}

// ContractEigenDARelayRegistryMetaData contains all meta data concerning the ContractEigenDARelayRegistry contract.
var ContractEigenDARelayRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addRelayInfo\",\"inputs\":[{\"name\":\"relayInfo\",\"type\":\"tuple\",\"internalType\":\"structRelayInfo\",\"components\":[{\"name\":\"relayAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"relayURL\",\"type\":\"string\",\"internalType\":\"string\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_initialOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"nextRelayKey\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"relayKeyToAddress\",\"inputs\":[{\"name\":\"key\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"relayKeyToInfo\",\"inputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"relayAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"relayURL\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"relayKeyToUrl\",\"inputs\":[{\"name\":\"key\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RelayAdded\",\"inputs\":[{\"name\":\"relay\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"key\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"},{\"name\":\"relayURL\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false}]",
	Bin: "0x60806040523480156200001157600080fd5b50620000226200002860201b60201c565b620001d3565b600060019054906101000a900460ff16156200007b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401620000729062000176565b60405180910390fd5b60ff801660008054906101000a900460ff1660ff161015620000ed5760ff6000806101000a81548160ff021916908360ff1602179055507f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb384740249860ff604051620000e49190620001b6565b60405180910390a15b565b600082825260208201905092915050565b7f496e697469616c697a61626c653a20636f6e747261637420697320696e69746960008201527f616c697a696e6700000000000000000000000000000000000000000000000000602082015250565b60006200015e602783620000ef565b91506200016b8262000100565b604082019050919050565b6000602082019050818103600083015262000191816200014f565b9050919050565b600060ff82169050919050565b620001b08162000198565b82525050565b6000602082019050620001cd6000830184620001a5565b92915050565b61105180620001e36000396000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c8063841f6a2e11610066578063841f6a2e146101205780638da5cb5b14610151578063b5a872da1461016f578063c4d66de81461019f578063f2fde38b146101bb57610093565b806315ddaa5d146100985780632fc35013146100b6578063631eabb8146100e6578063715018a614610116575b600080fd5b6100a06101d7565b6040516100ad9190610945565b60405180910390f35b6100d060048036038101906100cb9190610b8e565b6101ed565b6040516100dd9190610945565b60405180910390f35b61010060048036038101906100fb9190610c03565b610346565b60405161010d9190610cb8565b60405180910390f35b61011e6103fa565b005b61013a60048036038101906101359190610c03565b61040e565b604051610148929190610ce9565b60405180910390f35b6101596104da565b6040516101669190610d19565b60405180910390f35b61018960048036038101906101849190610c03565b610504565b6040516101969190610d19565b60405180910390f35b6101b960048036038101906101b49190610d34565b610550565b005b6101d560048036038101906101d09190610d34565b610690565b005b606660009054906101000a900463ffffffff1681565b60006101f7610714565b8160656000606660009054906101000a900463ffffffff1663ffffffff1663ffffffff16815260200190815260200160002060008201518160000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550602082015181600101908051906020019061028c929190610883565b50905050606660009054906101000a900463ffffffff1663ffffffff16826000015173ffffffffffffffffffffffffffffffffffffffff167f01c289e409d41a712a615bf286126433da55c193bbe64fc8e77af5f1ff13db9984602001516040516102f79190610cb8565b60405180910390a36066600081819054906101000a900463ffffffff168092919061032190610d90565b91906101000a81548163ffffffff021916908363ffffffff1602179055509050919050565b6060606560008363ffffffff1663ffffffff168152602001908152602001600020600101805461037590610dec565b80601f01602080910402602001604051908101604052809291908181526020018280546103a190610dec565b80156103ee5780601f106103c3576101008083540402835291602001916103ee565b820191906000526020600020905b8154815290600101906020018083116103d157829003601f168201915b50505050509050919050565b610402610714565b61040c6000610792565b565b60656020528060005260406000206000915090508060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169080600101805461045790610dec565b80601f016020809104026020016040519081016040528092919081815260200182805461048390610dec565b80156104d05780601f106104a5576101008083540402835291602001916104d0565b820191906000526020600020905b8154815290600101906020018083116104b357829003601f168201915b5050505050905082565b6000603360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6000606560008363ffffffff1663ffffffff16815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050919050565b60008060019054906101000a900460ff161590508080156105815750600160008054906101000a900460ff1660ff16105b806105ae575061059030610858565b1580156105ad5750600160008054906101000a900460ff1660ff16145b5b6105ed576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105e490610e90565b60405180910390fd5b60016000806101000a81548160ff021916908360ff160217905550801561062a576001600060016101000a81548160ff0219169083151502179055505b61063382610792565b801561068c5760008060016101000a81548160ff0219169083151502179055507f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb384740249860016040516106839190610f02565b60405180910390a15b5050565b610698610714565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415610708576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106ff90610f8f565b60405180910390fd5b61071181610792565b50565b61071c61087b565b73ffffffffffffffffffffffffffffffffffffffff1661073a6104da565b73ffffffffffffffffffffffffffffffffffffffff1614610790576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161078790610ffb565b60405180910390fd5b565b6000603360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905081603360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b6000808273ffffffffffffffffffffffffffffffffffffffff163b119050919050565b600033905090565b82805461088f90610dec565b90600052602060002090601f0160209004810192826108b157600085556108f8565b82601f106108ca57805160ff19168380011785556108f8565b828001600101855582156108f8579182015b828111156108f75782518255916020019190600101906108dc565b5b5090506109059190610909565b5090565b5b8082111561092257600081600090555060010161090a565b5090565b600063ffffffff82169050919050565b61093f81610926565b82525050565b600060208201905061095a6000830184610936565b92915050565b6000604051905090565b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6109c282610979565b810181811067ffffffffffffffff821117156109e1576109e061098a565b5b80604052505050565b60006109f4610960565b9050610a0082826109b9565b919050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610a3582610a0a565b9050919050565b610a4581610a2a565b8114610a5057600080fd5b50565b600081359050610a6281610a3c565b92915050565b600080fd5b600080fd5b600067ffffffffffffffff821115610a8d57610a8c61098a565b5b610a9682610979565b9050602081019050919050565b82818337600083830152505050565b6000610ac5610ac084610a72565b6109ea565b905082815260208101848484011115610ae157610ae0610a6d565b5b610aec848285610aa3565b509392505050565b600082601f830112610b0957610b08610a68565b5b8135610b19848260208601610ab2565b91505092915050565b600060408284031215610b3857610b37610974565b5b610b4260406109ea565b90506000610b5284828501610a53565b600083015250602082013567ffffffffffffffff811115610b7657610b75610a05565b5b610b8284828501610af4565b60208301525092915050565b600060208284031215610ba457610ba361096a565b5b600082013567ffffffffffffffff811115610bc257610bc161096f565b5b610bce84828501610b22565b91505092915050565b610be081610926565b8114610beb57600080fd5b50565b600081359050610bfd81610bd7565b92915050565b600060208284031215610c1957610c1861096a565b5b6000610c2784828501610bee565b91505092915050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610c6a578082015181840152602081019050610c4f565b83811115610c79576000848401525b50505050565b6000610c8a82610c30565b610c948185610c3b565b9350610ca4818560208601610c4c565b610cad81610979565b840191505092915050565b60006020820190508181036000830152610cd28184610c7f565b905092915050565b610ce381610a2a565b82525050565b6000604082019050610cfe6000830185610cda565b8181036020830152610d108184610c7f565b90509392505050565b6000602082019050610d2e6000830184610cda565b92915050565b600060208284031215610d4a57610d4961096a565b5b6000610d5884828501610a53565b91505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610d9b82610926565b915063ffffffff821415610db257610db1610d61565b5b600182019050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680610e0457607f821691505b60208210811415610e1857610e17610dbd565b5b50919050565b7f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160008201527f647920696e697469616c697a6564000000000000000000000000000000000000602082015250565b6000610e7a602e83610c3b565b9150610e8582610e1e565b604082019050919050565b60006020820190508181036000830152610ea981610e6d565b9050919050565b6000819050919050565b600060ff82169050919050565b6000819050919050565b6000610eec610ee7610ee284610eb0565b610ec7565b610eba565b9050919050565b610efc81610ed1565b82525050565b6000602082019050610f176000830184610ef3565b92915050565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b6000610f79602683610c3b565b9150610f8482610f1d565b604082019050919050565b60006020820190508181036000830152610fa881610f6c565b9050919050565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b6000610fe5602083610c3b565b9150610ff082610faf565b602082019050919050565b6000602082019050818103600083015261101481610fd8565b905091905056fea2646970667358221220c58fb8ecca7f4d0825d24f2bd02c0478d79873b67af502dc301076db33c3a47864736f6c634300080c0033",
}

// ContractEigenDARelayRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEigenDARelayRegistryMetaData.ABI instead.
var ContractEigenDARelayRegistryABI = ContractEigenDARelayRegistryMetaData.ABI

// ContractEigenDARelayRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractEigenDARelayRegistryMetaData.Bin instead.
var ContractEigenDARelayRegistryBin = ContractEigenDARelayRegistryMetaData.Bin

// DeployContractEigenDARelayRegistry deploys a new Ethereum contract, binding an instance of ContractEigenDARelayRegistry to it.
func DeployContractEigenDARelayRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ContractEigenDARelayRegistry, error) {
	parsed, err := ContractEigenDARelayRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractEigenDARelayRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractEigenDARelayRegistry{ContractEigenDARelayRegistryCaller: ContractEigenDARelayRegistryCaller{contract: contract}, ContractEigenDARelayRegistryTransactor: ContractEigenDARelayRegistryTransactor{contract: contract}, ContractEigenDARelayRegistryFilterer: ContractEigenDARelayRegistryFilterer{contract: contract}}, nil
}

// ContractEigenDARelayRegistry is an auto generated Go binding around an Ethereum contract.
type ContractEigenDARelayRegistry struct {
	ContractEigenDARelayRegistryCaller     // Read-only binding to the contract
	ContractEigenDARelayRegistryTransactor // Write-only binding to the contract
	ContractEigenDARelayRegistryFilterer   // Log filterer for contract events
}

// ContractEigenDARelayRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractEigenDARelayRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDARelayRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractEigenDARelayRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDARelayRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractEigenDARelayRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDARelayRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractEigenDARelayRegistrySession struct {
	Contract     *ContractEigenDARelayRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                 // Call options to use throughout this session
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// ContractEigenDARelayRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractEigenDARelayRegistryCallerSession struct {
	Contract *ContractEigenDARelayRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                       // Call options to use throughout this session
}

// ContractEigenDARelayRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractEigenDARelayRegistryTransactorSession struct {
	Contract     *ContractEigenDARelayRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                       // Transaction auth options to use throughout this session
}

// ContractEigenDARelayRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractEigenDARelayRegistryRaw struct {
	Contract *ContractEigenDARelayRegistry // Generic contract binding to access the raw methods on
}

// ContractEigenDARelayRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractEigenDARelayRegistryCallerRaw struct {
	Contract *ContractEigenDARelayRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// ContractEigenDARelayRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractEigenDARelayRegistryTransactorRaw struct {
	Contract *ContractEigenDARelayRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractEigenDARelayRegistry creates a new instance of ContractEigenDARelayRegistry, bound to a specific deployed contract.
func NewContractEigenDARelayRegistry(address common.Address, backend bind.ContractBackend) (*ContractEigenDARelayRegistry, error) {
	contract, err := bindContractEigenDARelayRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDARelayRegistry{ContractEigenDARelayRegistryCaller: ContractEigenDARelayRegistryCaller{contract: contract}, ContractEigenDARelayRegistryTransactor: ContractEigenDARelayRegistryTransactor{contract: contract}, ContractEigenDARelayRegistryFilterer: ContractEigenDARelayRegistryFilterer{contract: contract}}, nil
}

// NewContractEigenDARelayRegistryCaller creates a new read-only instance of ContractEigenDARelayRegistry, bound to a specific deployed contract.
func NewContractEigenDARelayRegistryCaller(address common.Address, caller bind.ContractCaller) (*ContractEigenDARelayRegistryCaller, error) {
	contract, err := bindContractEigenDARelayRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDARelayRegistryCaller{contract: contract}, nil
}

// NewContractEigenDARelayRegistryTransactor creates a new write-only instance of ContractEigenDARelayRegistry, bound to a specific deployed contract.
func NewContractEigenDARelayRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractEigenDARelayRegistryTransactor, error) {
	contract, err := bindContractEigenDARelayRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDARelayRegistryTransactor{contract: contract}, nil
}

// NewContractEigenDARelayRegistryFilterer creates a new log filterer instance of ContractEigenDARelayRegistry, bound to a specific deployed contract.
func NewContractEigenDARelayRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractEigenDARelayRegistryFilterer, error) {
	contract, err := bindContractEigenDARelayRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDARelayRegistryFilterer{contract: contract}, nil
}

// bindContractEigenDARelayRegistry binds a generic wrapper to an already deployed contract.
func bindContractEigenDARelayRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractEigenDARelayRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDARelayRegistry.Contract.ContractEigenDARelayRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.Contract.ContractEigenDARelayRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.Contract.ContractEigenDARelayRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDARelayRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.Contract.contract.Transact(opts, method, params...)
}

// NextRelayKey is a free data retrieval call binding the contract method 0x15ddaa5d.
//
// Solidity: function nextRelayKey() view returns(uint32)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryCaller) NextRelayKey(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _ContractEigenDARelayRegistry.contract.Call(opts, &out, "nextRelayKey")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// NextRelayKey is a free data retrieval call binding the contract method 0x15ddaa5d.
//
// Solidity: function nextRelayKey() view returns(uint32)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistrySession) NextRelayKey() (uint32, error) {
	return _ContractEigenDARelayRegistry.Contract.NextRelayKey(&_ContractEigenDARelayRegistry.CallOpts)
}

// NextRelayKey is a free data retrieval call binding the contract method 0x15ddaa5d.
//
// Solidity: function nextRelayKey() view returns(uint32)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryCallerSession) NextRelayKey() (uint32, error) {
	return _ContractEigenDARelayRegistry.Contract.NextRelayKey(&_ContractEigenDARelayRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDARelayRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistrySession) Owner() (common.Address, error) {
	return _ContractEigenDARelayRegistry.Contract.Owner(&_ContractEigenDARelayRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryCallerSession) Owner() (common.Address, error) {
	return _ContractEigenDARelayRegistry.Contract.Owner(&_ContractEigenDARelayRegistry.CallOpts)
}

// RelayKeyToAddress is a free data retrieval call binding the contract method 0xb5a872da.
//
// Solidity: function relayKeyToAddress(uint32 key) view returns(address)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryCaller) RelayKeyToAddress(opts *bind.CallOpts, key uint32) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDARelayRegistry.contract.Call(opts, &out, "relayKeyToAddress", key)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RelayKeyToAddress is a free data retrieval call binding the contract method 0xb5a872da.
//
// Solidity: function relayKeyToAddress(uint32 key) view returns(address)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistrySession) RelayKeyToAddress(key uint32) (common.Address, error) {
	return _ContractEigenDARelayRegistry.Contract.RelayKeyToAddress(&_ContractEigenDARelayRegistry.CallOpts, key)
}

// RelayKeyToAddress is a free data retrieval call binding the contract method 0xb5a872da.
//
// Solidity: function relayKeyToAddress(uint32 key) view returns(address)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryCallerSession) RelayKeyToAddress(key uint32) (common.Address, error) {
	return _ContractEigenDARelayRegistry.Contract.RelayKeyToAddress(&_ContractEigenDARelayRegistry.CallOpts, key)
}

// RelayKeyToInfo is a free data retrieval call binding the contract method 0x841f6a2e.
//
// Solidity: function relayKeyToInfo(uint32 ) view returns(address relayAddress, string relayURL)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryCaller) RelayKeyToInfo(opts *bind.CallOpts, arg0 uint32) (struct {
	RelayAddress common.Address
	RelayURL     string
}, error) {
	var out []interface{}
	err := _ContractEigenDARelayRegistry.contract.Call(opts, &out, "relayKeyToInfo", arg0)

	outstruct := new(struct {
		RelayAddress common.Address
		RelayURL     string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RelayAddress = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.RelayURL = *abi.ConvertType(out[1], new(string)).(*string)

	return *outstruct, err

}

// RelayKeyToInfo is a free data retrieval call binding the contract method 0x841f6a2e.
//
// Solidity: function relayKeyToInfo(uint32 ) view returns(address relayAddress, string relayURL)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistrySession) RelayKeyToInfo(arg0 uint32) (struct {
	RelayAddress common.Address
	RelayURL     string
}, error) {
	return _ContractEigenDARelayRegistry.Contract.RelayKeyToInfo(&_ContractEigenDARelayRegistry.CallOpts, arg0)
}

// RelayKeyToInfo is a free data retrieval call binding the contract method 0x841f6a2e.
//
// Solidity: function relayKeyToInfo(uint32 ) view returns(address relayAddress, string relayURL)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryCallerSession) RelayKeyToInfo(arg0 uint32) (struct {
	RelayAddress common.Address
	RelayURL     string
}, error) {
	return _ContractEigenDARelayRegistry.Contract.RelayKeyToInfo(&_ContractEigenDARelayRegistry.CallOpts, arg0)
}

// RelayKeyToUrl is a free data retrieval call binding the contract method 0x631eabb8.
//
// Solidity: function relayKeyToUrl(uint32 key) view returns(string)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryCaller) RelayKeyToUrl(opts *bind.CallOpts, key uint32) (string, error) {
	var out []interface{}
	err := _ContractEigenDARelayRegistry.contract.Call(opts, &out, "relayKeyToUrl", key)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// RelayKeyToUrl is a free data retrieval call binding the contract method 0x631eabb8.
//
// Solidity: function relayKeyToUrl(uint32 key) view returns(string)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistrySession) RelayKeyToUrl(key uint32) (string, error) {
	return _ContractEigenDARelayRegistry.Contract.RelayKeyToUrl(&_ContractEigenDARelayRegistry.CallOpts, key)
}

// RelayKeyToUrl is a free data retrieval call binding the contract method 0x631eabb8.
//
// Solidity: function relayKeyToUrl(uint32 key) view returns(string)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryCallerSession) RelayKeyToUrl(key uint32) (string, error) {
	return _ContractEigenDARelayRegistry.Contract.RelayKeyToUrl(&_ContractEigenDARelayRegistry.CallOpts, key)
}

// AddRelayInfo is a paid mutator transaction binding the contract method 0x2fc35013.
//
// Solidity: function addRelayInfo((address,string) relayInfo) returns(uint32)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryTransactor) AddRelayInfo(opts *bind.TransactOpts, relayInfo RelayInfo) (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.contract.Transact(opts, "addRelayInfo", relayInfo)
}

// AddRelayInfo is a paid mutator transaction binding the contract method 0x2fc35013.
//
// Solidity: function addRelayInfo((address,string) relayInfo) returns(uint32)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistrySession) AddRelayInfo(relayInfo RelayInfo) (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.Contract.AddRelayInfo(&_ContractEigenDARelayRegistry.TransactOpts, relayInfo)
}

// AddRelayInfo is a paid mutator transaction binding the contract method 0x2fc35013.
//
// Solidity: function addRelayInfo((address,string) relayInfo) returns(uint32)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryTransactorSession) AddRelayInfo(relayInfo RelayInfo) (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.Contract.AddRelayInfo(&_ContractEigenDARelayRegistry.TransactOpts, relayInfo)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _initialOwner) returns()
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryTransactor) Initialize(opts *bind.TransactOpts, _initialOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.contract.Transact(opts, "initialize", _initialOwner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _initialOwner) returns()
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistrySession) Initialize(_initialOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.Contract.Initialize(&_ContractEigenDARelayRegistry.TransactOpts, _initialOwner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _initialOwner) returns()
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryTransactorSession) Initialize(_initialOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.Contract.Initialize(&_ContractEigenDARelayRegistry.TransactOpts, _initialOwner)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistrySession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.Contract.RenounceOwnership(&_ContractEigenDARelayRegistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.Contract.RenounceOwnership(&_ContractEigenDARelayRegistry.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.Contract.TransferOwnership(&_ContractEigenDARelayRegistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.Contract.TransferOwnership(&_ContractEigenDARelayRegistry.TransactOpts, newOwner)
}

// ContractEigenDARelayRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ContractEigenDARelayRegistry contract.
type ContractEigenDARelayRegistryInitializedIterator struct {
	Event *ContractEigenDARelayRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *ContractEigenDARelayRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDARelayRegistryInitialized)
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
		it.Event = new(ContractEigenDARelayRegistryInitialized)
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
func (it *ContractEigenDARelayRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDARelayRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDARelayRegistryInitialized represents a Initialized event raised by the ContractEigenDARelayRegistry contract.
type ContractEigenDARelayRegistryInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*ContractEigenDARelayRegistryInitializedIterator, error) {

	logs, sub, err := _ContractEigenDARelayRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDARelayRegistryInitializedIterator{contract: _ContractEigenDARelayRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ContractEigenDARelayRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDARelayRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDARelayRegistryInitialized)
				if err := _ContractEigenDARelayRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryFilterer) ParseInitialized(log types.Log) (*ContractEigenDARelayRegistryInitialized, error) {
	event := new(ContractEigenDARelayRegistryInitialized)
	if err := _ContractEigenDARelayRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDARelayRegistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ContractEigenDARelayRegistry contract.
type ContractEigenDARelayRegistryOwnershipTransferredIterator struct {
	Event *ContractEigenDARelayRegistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ContractEigenDARelayRegistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDARelayRegistryOwnershipTransferred)
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
		it.Event = new(ContractEigenDARelayRegistryOwnershipTransferred)
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
func (it *ContractEigenDARelayRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDARelayRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDARelayRegistryOwnershipTransferred represents a OwnershipTransferred event raised by the ContractEigenDARelayRegistry contract.
type ContractEigenDARelayRegistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ContractEigenDARelayRegistryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractEigenDARelayRegistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDARelayRegistryOwnershipTransferredIterator{contract: _ContractEigenDARelayRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ContractEigenDARelayRegistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractEigenDARelayRegistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDARelayRegistryOwnershipTransferred)
				if err := _ContractEigenDARelayRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*ContractEigenDARelayRegistryOwnershipTransferred, error) {
	event := new(ContractEigenDARelayRegistryOwnershipTransferred)
	if err := _ContractEigenDARelayRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDARelayRegistryRelayAddedIterator is returned from FilterRelayAdded and is used to iterate over the raw logs and unpacked data for RelayAdded events raised by the ContractEigenDARelayRegistry contract.
type ContractEigenDARelayRegistryRelayAddedIterator struct {
	Event *ContractEigenDARelayRegistryRelayAdded // Event containing the contract specifics and raw log

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
func (it *ContractEigenDARelayRegistryRelayAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDARelayRegistryRelayAdded)
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
		it.Event = new(ContractEigenDARelayRegistryRelayAdded)
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
func (it *ContractEigenDARelayRegistryRelayAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDARelayRegistryRelayAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDARelayRegistryRelayAdded represents a RelayAdded event raised by the ContractEigenDARelayRegistry contract.
type ContractEigenDARelayRegistryRelayAdded struct {
	Relay    common.Address
	Key      uint32
	RelayURL string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterRelayAdded is a free log retrieval operation binding the contract event 0x01c289e409d41a712a615bf286126433da55c193bbe64fc8e77af5f1ff13db99.
//
// Solidity: event RelayAdded(address indexed relay, uint32 indexed key, string relayURL)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryFilterer) FilterRelayAdded(opts *bind.FilterOpts, relay []common.Address, key []uint32) (*ContractEigenDARelayRegistryRelayAddedIterator, error) {

	var relayRule []interface{}
	for _, relayItem := range relay {
		relayRule = append(relayRule, relayItem)
	}
	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _ContractEigenDARelayRegistry.contract.FilterLogs(opts, "RelayAdded", relayRule, keyRule)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDARelayRegistryRelayAddedIterator{contract: _ContractEigenDARelayRegistry.contract, event: "RelayAdded", logs: logs, sub: sub}, nil
}

// WatchRelayAdded is a free log subscription operation binding the contract event 0x01c289e409d41a712a615bf286126433da55c193bbe64fc8e77af5f1ff13db99.
//
// Solidity: event RelayAdded(address indexed relay, uint32 indexed key, string relayURL)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryFilterer) WatchRelayAdded(opts *bind.WatchOpts, sink chan<- *ContractEigenDARelayRegistryRelayAdded, relay []common.Address, key []uint32) (event.Subscription, error) {

	var relayRule []interface{}
	for _, relayItem := range relay {
		relayRule = append(relayRule, relayItem)
	}
	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _ContractEigenDARelayRegistry.contract.WatchLogs(opts, "RelayAdded", relayRule, keyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDARelayRegistryRelayAdded)
				if err := _ContractEigenDARelayRegistry.contract.UnpackLog(event, "RelayAdded", log); err != nil {
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

// ParseRelayAdded is a log parse operation binding the contract event 0x01c289e409d41a712a615bf286126433da55c193bbe64fc8e77af5f1ff13db99.
//
// Solidity: event RelayAdded(address indexed relay, uint32 indexed key, string relayURL)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryFilterer) ParseRelayAdded(log types.Log) (*ContractEigenDARelayRegistryRelayAdded, error) {
	event := new(ContractEigenDARelayRegistryRelayAdded)
	if err := _ContractEigenDARelayRegistry.contract.UnpackLog(event, "RelayAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
