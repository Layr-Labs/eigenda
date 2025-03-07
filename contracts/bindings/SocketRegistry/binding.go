// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractSocketRegistry

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

// ContractSocketRegistryMetaData contains all meta data concerning the ContractSocketRegistry contract.
var ContractSocketRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getOperatorSocket\",\"inputs\":[{\"name\":\"_operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operatorIdToSocket\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registryCoordinator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setOperatorSocket\",\"inputs\":[{\"name\":\"_operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_socket\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"}]",
	Bin: "0x60a060405234801561001057600080fd5b5060405161094d38038061094d833981810160405281019061003291906100e1565b8073ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff16815250505061010e565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061009c82610071565b9050919050565b60006100ae82610091565b9050919050565b6100be816100a3565b81146100c957600080fd5b50565b6000815190506100db816100b5565b92915050565b6000602082840312156100f7576100f661006c565b5b6000610105848285016100cc565b91505092915050565b60805161081d610130600039600081816101910152610255015261081d6000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c806310bea0d7146100515780636d14a98714610081578063af65fdfc1461009f578063f043367e146100cf575b600080fd5b61006b600480360381019061006691906103f9565b6100eb565b60405161007891906104bf565b60405180910390f35b61008961018f565b6040516100969190610522565b60405180910390f35b6100b960048036038101906100b491906103f9565b6101b3565b6040516100c691906104bf565b60405180910390f35b6100e960048036038101906100e49190610672565b610253565b005b6060600080838152602001908152602001600020805461010a906106fd565b80601f0160208091040260200160405190810160405280929190818152602001828054610136906106fd565b80156101835780601f1061015857610100808354040283529160200191610183565b820191906000526020600020905b81548152906001019060200180831161016657829003601f168201915b50505050509050919050565b7f000000000000000000000000000000000000000000000000000000000000000081565b600060205280600052604060002060009150905080546101d2906106fd565b80601f01602080910402602001604051908101604052809291908181526020018280546101fe906106fd565b801561024b5780601f106102205761010080835404028352916020019161024b565b820191906000526020600020905b81548152906001019060200180831161022e57829003601f168201915b505050505081565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146102e1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102d8906107c7565b60405180910390fd5b80600080848152602001908152602001600020908051906020019061030792919061030c565b505050565b828054610318906106fd565b90600052602060002090601f01602090048101928261033a5760008555610381565b82601f1061035357805160ff1916838001178555610381565b82800160010185558215610381579182015b82811115610380578251825591602001919060010190610365565b5b50905061038e9190610392565b5090565b5b808211156103ab576000816000905550600101610393565b5090565b6000604051905090565b600080fd5b600080fd5b6000819050919050565b6103d6816103c3565b81146103e157600080fd5b50565b6000813590506103f3816103cd565b92915050565b60006020828403121561040f5761040e6103b9565b5b600061041d848285016103e4565b91505092915050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610460578082015181840152602081019050610445565b8381111561046f576000848401525b50505050565b6000601f19601f8301169050919050565b600061049182610426565b61049b8185610431565b93506104ab818560208601610442565b6104b481610475565b840191505092915050565b600060208201905081810360008301526104d98184610486565b905092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061050c826104e1565b9050919050565b61051c81610501565b82525050565b60006020820190506105376000830184610513565b92915050565b600080fd5b600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b61057f82610475565b810181811067ffffffffffffffff8211171561059e5761059d610547565b5b80604052505050565b60006105b16103af565b90506105bd8282610576565b919050565b600067ffffffffffffffff8211156105dd576105dc610547565b5b6105e682610475565b9050602081019050919050565b82818337600083830152505050565b6000610615610610846105c2565b6105a7565b90508281526020810184848401111561063157610630610542565b5b61063c8482856105f3565b509392505050565b600082601f8301126106595761065861053d565b5b8135610669848260208601610602565b91505092915050565b60008060408385031215610689576106886103b9565b5b6000610697858286016103e4565b925050602083013567ffffffffffffffff8111156106b8576106b76103be565b5b6106c485828601610644565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061071557607f821691505b60208210811415610729576107286106ce565b5b50919050565b7f536f636b657452656769737472792e6f6e6c795265676973747279436f6f726460008201527f696e61746f723a2063616c6c6572206973206e6f74207468652052656769737460208201527f7279436f6f7264696e61746f7200000000000000000000000000000000000000604082015250565b60006107b1604d83610431565b91506107bc8261072f565b606082019050919050565b600060208201905081810360008301526107e0816107a4565b905091905056fea2646970667358221220de365e0ede30082de57079b958a68b2e6f623fb93cb9324bb2992b0121b3ecfb64736f6c634300080c0033",
}

// ContractSocketRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractSocketRegistryMetaData.ABI instead.
var ContractSocketRegistryABI = ContractSocketRegistryMetaData.ABI

// ContractSocketRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractSocketRegistryMetaData.Bin instead.
var ContractSocketRegistryBin = ContractSocketRegistryMetaData.Bin

// DeployContractSocketRegistry deploys a new Ethereum contract, binding an instance of ContractSocketRegistry to it.
func DeployContractSocketRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, _registryCoordinator common.Address) (common.Address, *types.Transaction, *ContractSocketRegistry, error) {
	parsed, err := ContractSocketRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractSocketRegistryBin), backend, _registryCoordinator)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractSocketRegistry{ContractSocketRegistryCaller: ContractSocketRegistryCaller{contract: contract}, ContractSocketRegistryTransactor: ContractSocketRegistryTransactor{contract: contract}, ContractSocketRegistryFilterer: ContractSocketRegistryFilterer{contract: contract}}, nil
}

// ContractSocketRegistry is an auto generated Go binding around an Ethereum contract.
type ContractSocketRegistry struct {
	ContractSocketRegistryCaller     // Read-only binding to the contract
	ContractSocketRegistryTransactor // Write-only binding to the contract
	ContractSocketRegistryFilterer   // Log filterer for contract events
}

// ContractSocketRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractSocketRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSocketRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractSocketRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSocketRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractSocketRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSocketRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractSocketRegistrySession struct {
	Contract     *ContractSocketRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// ContractSocketRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractSocketRegistryCallerSession struct {
	Contract *ContractSocketRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// ContractSocketRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractSocketRegistryTransactorSession struct {
	Contract     *ContractSocketRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// ContractSocketRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractSocketRegistryRaw struct {
	Contract *ContractSocketRegistry // Generic contract binding to access the raw methods on
}

// ContractSocketRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractSocketRegistryCallerRaw struct {
	Contract *ContractSocketRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// ContractSocketRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractSocketRegistryTransactorRaw struct {
	Contract *ContractSocketRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractSocketRegistry creates a new instance of ContractSocketRegistry, bound to a specific deployed contract.
func NewContractSocketRegistry(address common.Address, backend bind.ContractBackend) (*ContractSocketRegistry, error) {
	contract, err := bindContractSocketRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractSocketRegistry{ContractSocketRegistryCaller: ContractSocketRegistryCaller{contract: contract}, ContractSocketRegistryTransactor: ContractSocketRegistryTransactor{contract: contract}, ContractSocketRegistryFilterer: ContractSocketRegistryFilterer{contract: contract}}, nil
}

// NewContractSocketRegistryCaller creates a new read-only instance of ContractSocketRegistry, bound to a specific deployed contract.
func NewContractSocketRegistryCaller(address common.Address, caller bind.ContractCaller) (*ContractSocketRegistryCaller, error) {
	contract, err := bindContractSocketRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractSocketRegistryCaller{contract: contract}, nil
}

// NewContractSocketRegistryTransactor creates a new write-only instance of ContractSocketRegistry, bound to a specific deployed contract.
func NewContractSocketRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractSocketRegistryTransactor, error) {
	contract, err := bindContractSocketRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractSocketRegistryTransactor{contract: contract}, nil
}

// NewContractSocketRegistryFilterer creates a new log filterer instance of ContractSocketRegistry, bound to a specific deployed contract.
func NewContractSocketRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractSocketRegistryFilterer, error) {
	contract, err := bindContractSocketRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractSocketRegistryFilterer{contract: contract}, nil
}

// bindContractSocketRegistry binds a generic wrapper to an already deployed contract.
func bindContractSocketRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractSocketRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractSocketRegistry *ContractSocketRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractSocketRegistry.Contract.ContractSocketRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractSocketRegistry *ContractSocketRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractSocketRegistry.Contract.ContractSocketRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractSocketRegistry *ContractSocketRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractSocketRegistry.Contract.ContractSocketRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractSocketRegistry *ContractSocketRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractSocketRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractSocketRegistry *ContractSocketRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractSocketRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractSocketRegistry *ContractSocketRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractSocketRegistry.Contract.contract.Transact(opts, method, params...)
}

// GetOperatorSocket is a free data retrieval call binding the contract method 0x10bea0d7.
//
// Solidity: function getOperatorSocket(bytes32 _operatorId) view returns(string)
func (_ContractSocketRegistry *ContractSocketRegistryCaller) GetOperatorSocket(opts *bind.CallOpts, _operatorId [32]byte) (string, error) {
	var out []interface{}
	err := _ContractSocketRegistry.contract.Call(opts, &out, "getOperatorSocket", _operatorId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetOperatorSocket is a free data retrieval call binding the contract method 0x10bea0d7.
//
// Solidity: function getOperatorSocket(bytes32 _operatorId) view returns(string)
func (_ContractSocketRegistry *ContractSocketRegistrySession) GetOperatorSocket(_operatorId [32]byte) (string, error) {
	return _ContractSocketRegistry.Contract.GetOperatorSocket(&_ContractSocketRegistry.CallOpts, _operatorId)
}

// GetOperatorSocket is a free data retrieval call binding the contract method 0x10bea0d7.
//
// Solidity: function getOperatorSocket(bytes32 _operatorId) view returns(string)
func (_ContractSocketRegistry *ContractSocketRegistryCallerSession) GetOperatorSocket(_operatorId [32]byte) (string, error) {
	return _ContractSocketRegistry.Contract.GetOperatorSocket(&_ContractSocketRegistry.CallOpts, _operatorId)
}

// OperatorIdToSocket is a free data retrieval call binding the contract method 0xaf65fdfc.
//
// Solidity: function operatorIdToSocket(bytes32 ) view returns(string)
func (_ContractSocketRegistry *ContractSocketRegistryCaller) OperatorIdToSocket(opts *bind.CallOpts, arg0 [32]byte) (string, error) {
	var out []interface{}
	err := _ContractSocketRegistry.contract.Call(opts, &out, "operatorIdToSocket", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// OperatorIdToSocket is a free data retrieval call binding the contract method 0xaf65fdfc.
//
// Solidity: function operatorIdToSocket(bytes32 ) view returns(string)
func (_ContractSocketRegistry *ContractSocketRegistrySession) OperatorIdToSocket(arg0 [32]byte) (string, error) {
	return _ContractSocketRegistry.Contract.OperatorIdToSocket(&_ContractSocketRegistry.CallOpts, arg0)
}

// OperatorIdToSocket is a free data retrieval call binding the contract method 0xaf65fdfc.
//
// Solidity: function operatorIdToSocket(bytes32 ) view returns(string)
func (_ContractSocketRegistry *ContractSocketRegistryCallerSession) OperatorIdToSocket(arg0 [32]byte) (string, error) {
	return _ContractSocketRegistry.Contract.OperatorIdToSocket(&_ContractSocketRegistry.CallOpts, arg0)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractSocketRegistry *ContractSocketRegistryCaller) RegistryCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractSocketRegistry.contract.Call(opts, &out, "registryCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractSocketRegistry *ContractSocketRegistrySession) RegistryCoordinator() (common.Address, error) {
	return _ContractSocketRegistry.Contract.RegistryCoordinator(&_ContractSocketRegistry.CallOpts)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractSocketRegistry *ContractSocketRegistryCallerSession) RegistryCoordinator() (common.Address, error) {
	return _ContractSocketRegistry.Contract.RegistryCoordinator(&_ContractSocketRegistry.CallOpts)
}

// SetOperatorSocket is a paid mutator transaction binding the contract method 0xf043367e.
//
// Solidity: function setOperatorSocket(bytes32 _operatorId, string _socket) returns()
func (_ContractSocketRegistry *ContractSocketRegistryTransactor) SetOperatorSocket(opts *bind.TransactOpts, _operatorId [32]byte, _socket string) (*types.Transaction, error) {
	return _ContractSocketRegistry.contract.Transact(opts, "setOperatorSocket", _operatorId, _socket)
}

// SetOperatorSocket is a paid mutator transaction binding the contract method 0xf043367e.
//
// Solidity: function setOperatorSocket(bytes32 _operatorId, string _socket) returns()
func (_ContractSocketRegistry *ContractSocketRegistrySession) SetOperatorSocket(_operatorId [32]byte, _socket string) (*types.Transaction, error) {
	return _ContractSocketRegistry.Contract.SetOperatorSocket(&_ContractSocketRegistry.TransactOpts, _operatorId, _socket)
}

// SetOperatorSocket is a paid mutator transaction binding the contract method 0xf043367e.
//
// Solidity: function setOperatorSocket(bytes32 _operatorId, string _socket) returns()
func (_ContractSocketRegistry *ContractSocketRegistryTransactorSession) SetOperatorSocket(_operatorId [32]byte, _socket string) (*types.Transaction, error) {
	return _ContractSocketRegistry.Contract.SetOperatorSocket(&_ContractSocketRegistry.TransactOpts, _operatorId, _socket)
}
