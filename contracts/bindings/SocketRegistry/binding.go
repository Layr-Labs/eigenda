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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getOperatorSocket\",\"inputs\":[{\"name\":\"_operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"migrateOperatorSockets\",\"inputs\":[{\"name\":\"_operatorIds\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"_sockets\",\"type\":\"string[]\",\"internalType\":\"string[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"operatorIdToSocket\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registryCoordinator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setOperatorSocket\",\"inputs\":[{\"name\":\"_operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_socket\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"}]",
	Bin: "0x60a060405234801561001057600080fd5b5060405161099738038061099783398101604081905261002f91610040565b6001600160a01b0316608052610070565b60006020828403121561005257600080fd5b81516001600160a01b038116811461006957600080fd5b9392505050565b6080516108ff61009860003960008181609f015281816101a301526103f501526108ff6000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c806310bea0d71461005c578063639b9957146100855780636d14a9871461009a578063af65fdfc146100d9578063f043367e146100ec575b600080fd5b61006f61006a366004610550565b6100ff565b60405161007c9190610569565b60405180910390f35b610098610093366004610729565b6101a1565b005b6100c17f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b03909116815260200161007c565b61006f6100e7366004610550565b610350565b6100986100fa3660046107e2565b6103ea565b600081815260208190526040902080546060919061011c9061081f565b80601f01602080910402602001604051908101604052809291908181526020018280546101489061081f565b80156101955780601f1061016a57610100808354040283529160200191610195565b820191906000526020600020905b81548152906001019060200180831161017857829003601f168201915b50505050509050919050565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316638da5cb5b6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156101ff573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610223919061085a565b6001600160a01b0316336001600160a01b0316146102d45760405162461bcd60e51b815260206004820152605760248201527f536f636b657452656769737472792e6f6e6c79436f6f7264696e61746f724f7760448201527f6e65723a2063616c6c6572206973206e6f7420746865206f776e6572206f662060648201527f746865207265676973747279436f6f7264696e61746f72000000000000000000608482015260a4015b60405180910390fd5b60005b825181101561034b578181815181106102f2576102f261088a565b602002602001015160008085848151811061030f5761030f61088a565b6020026020010151815260200190815260200160002090805190602001906103389291906104b7565b5080610343816108a0565b9150506102d7565b505050565b600060208190529081526040902080546103699061081f565b80601f01602080910402602001604051908101604052809291908181526020018280546103959061081f565b80156103e25780601f106103b7576101008083540402835291602001916103e2565b820191906000526020600020905b8154815290600101906020018083116103c557829003601f168201915b505050505081565b336001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161461049e5760405162461bcd60e51b815260206004820152604d60248201527f536f636b657452656769737472792e6f6e6c795265676973747279436f6f726460448201527f696e61746f723a2063616c6c6572206973206e6f74207468652052656769737460648201526c393ca1b7b7b93234b730ba37b960991b608482015260a4016102cb565b600082815260208181526040909120825161034b928401905b8280546104c39061081f565b90600052602060002090601f0160209004810192826104e5576000855561052b565b82601f106104fe57805160ff191683800117855561052b565b8280016001018555821561052b579182015b8281111561052b578251825591602001919060010190610510565b5061053792915061053b565b5090565b5b80821115610537576000815560010161053c565b60006020828403121561056257600080fd5b5035919050565b600060208083528351808285015260005b818110156105965785810183015185820160400152820161057a565b818111156105a8576000604083870101525b50601f01601f1916929092016040019392505050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff811182821017156105fd576105fd6105be565b604052919050565b600067ffffffffffffffff82111561061f5761061f6105be565b5060051b60200190565b600082601f83011261063a57600080fd5b813567ffffffffffffffff811115610654576106546105be565b610667601f8201601f19166020016105d4565b81815284602083860101111561067c57600080fd5b816020850160208301376000918101602001919091529392505050565b600082601f8301126106aa57600080fd5b813560206106bf6106ba83610605565b6105d4565b82815260059290921b840181019181810190868411156106de57600080fd5b8286015b8481101561071e57803567ffffffffffffffff8111156107025760008081fd5b6107108986838b0101610629565b8452509183019183016106e2565b509695505050505050565b6000806040838503121561073c57600080fd5b823567ffffffffffffffff8082111561075457600080fd5b818501915085601f83011261076857600080fd5b813560206107786106ba83610605565b82815260059290921b8401810191818101908984111561079757600080fd5b948201945b838610156107b55785358252948201949082019061079c565b965050860135925050808211156107cb57600080fd5b506107d885828601610699565b9150509250929050565b600080604083850312156107f557600080fd5b82359150602083013567ffffffffffffffff81111561081357600080fd5b6107d885828601610629565b600181811c9082168061083357607f821691505b6020821081141561085457634e487b7160e01b600052602260045260246000fd5b50919050565b60006020828403121561086c57600080fd5b81516001600160a01b038116811461088357600080fd5b9392505050565b634e487b7160e01b600052603260045260246000fd5b60006000198214156108c257634e487b7160e01b600052601160045260246000fd5b506001019056fea2646970667358221220e893cad91c4a5486104415b7ee09c65fdf8b97a68219a30bf01d9d8c9ca5d9c964736f6c634300080c0033",
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

// MigrateOperatorSockets is a paid mutator transaction binding the contract method 0x639b9957.
//
// Solidity: function migrateOperatorSockets(bytes32[] _operatorIds, string[] _sockets) returns()
func (_ContractSocketRegistry *ContractSocketRegistryTransactor) MigrateOperatorSockets(opts *bind.TransactOpts, _operatorIds [][32]byte, _sockets []string) (*types.Transaction, error) {
	return _ContractSocketRegistry.contract.Transact(opts, "migrateOperatorSockets", _operatorIds, _sockets)
}

// MigrateOperatorSockets is a paid mutator transaction binding the contract method 0x639b9957.
//
// Solidity: function migrateOperatorSockets(bytes32[] _operatorIds, string[] _sockets) returns()
func (_ContractSocketRegistry *ContractSocketRegistrySession) MigrateOperatorSockets(_operatorIds [][32]byte, _sockets []string) (*types.Transaction, error) {
	return _ContractSocketRegistry.Contract.MigrateOperatorSockets(&_ContractSocketRegistry.TransactOpts, _operatorIds, _sockets)
}

// MigrateOperatorSockets is a paid mutator transaction binding the contract method 0x639b9957.
//
// Solidity: function migrateOperatorSockets(bytes32[] _operatorIds, string[] _sockets) returns()
func (_ContractSocketRegistry *ContractSocketRegistryTransactorSession) MigrateOperatorSockets(_operatorIds [][32]byte, _sockets []string) (*types.Transaction, error) {
	return _ContractSocketRegistry.Contract.MigrateOperatorSockets(&_ContractSocketRegistry.TransactOpts, _operatorIds, _sockets)
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
