// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractBLSPublicKeyCompendium

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

// ContractBLSPublicKeyCompendiumMetaData contains all meta data concerning the ContractBLSPublicKeyCompendium contract.
var ContractBLSPublicKeyCompendiumMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"getMessageHash\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operatorToPubkeyHash\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pubkeyHashToOperator\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerBLSPublicKey\",\"inputs\":[{\"name\":\"signedMessageHash\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"pubkeyG1\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"pubkeyG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"NewPubkeyRegistration\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"pubkeyG1\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"pubkeyG2\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]}],\"anonymous\":false}]",
	Bin: "0x608060405234801561001057600080fd5b50610f81806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c8063161a334d146100515780631f5ac1b214610066578063de29fac01461008f578063e8bb9ae6146100bd575b600080fd5b61006461005f366004610d03565b6100fe565b005b610079610074366004610d7c565b61045b565b6040516100869190610dac565b60405180910390f35b6100af61009d366004610d7c565b60006020819052908152604090205481565b604051908152602001610086565b6100e66100cb366004610dc3565b6001602052600090815260409020546001600160a01b031681565b6040516001600160a01b039091168152602001610086565b6000610109836104f6565b33600090815260208190526040902054909150156101ac5760405162461bcd60e51b815260206004820152604f60248201527f424c535075626c69634b6579436f6d70656e6469756d2e72656769737465724260448201527f4c535075626c69634b65793a206f70657261746f7220616c726561647920726560648201526e6769737465726564207075626b657960881b608482015260a4015b60405180910390fd5b6000818152600160205260409020546001600160a01b03161561024a5760405162461bcd60e51b815260206004820152604a60248201527f424c535075626c69634b6579436f6d70656e6469756d2e72656769737465724260448201527f4c535075626c69634b65793a207075626c6963206b657920616c7265616479206064820152691c9959da5cdd195c995960b21b608482015260a4016101a3565b60006102553361045b565b8551602080880151875188830151885189850151875186890151604051999a506000997f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001996102ad9990989796959493929101610e05565b6040516020818303038152906040528051906020012060001c6102d09190610e51565b90506103376102e96102e28784610539565b88906105d0565b6102f1610664565b61033161032a85610324604080518082018252600080825260209182015281518083019092526001825260029082015290565b90610539565b86906105d0565b87610724565b6103ec5760405162461bcd60e51b815260206004820152607460248201527f424c535075626c69634b6579436f6d70656e6469756d2e72656769737465724260448201527f4c535075626c69634b65793a2065697468657220746865204731207369676e6160648201527f747572652069732077726f6e672c206f7220473120616e6420473220707269766084820152730c2e8ca40d6caf240c8de40dcdee840dac2e8c6d60631b60a482015260c4016101a3565b3360008181526020818152604080832087905586835260019091529081902080546001600160a01b03191683179055517fe3fb6613af2e8930cf85d47fcf6db10192224a64c6cbe8023e0eee1ba38280419061044b9088908890610e73565b60405180910390a2505050505050565b60408051808201909152600080825260208201526040516bffffffffffffffffffffffff19606084811b8216602084015230901b1660348201524660488201527f456967656e4c617965725f424e3235345f5075626b65795f52656769737472616068820152633a34b7b760e11b60888201526104f090608c0160405160208183030381529060405280519060200120610991565b92915050565b60008160000151826020015160405160200161051c929190918252602082015260400190565b604051602081830303815290604052805190602001209050919050565b6040805180820190915260008082526020820152610555610b4b565b835181526020808501519082015260408082018490526000908360608460076107d05a03fa90508080156105885761058a565bfe5b50806105c85760405162461bcd60e51b815260206004820152600d60248201526c1958cb5b5d5b0b59985a5b1959609a1b60448201526064016101a3565b505092915050565b60408051808201909152600080825260208201526105ec610b69565b835181526020808501518183015283516040808401919091529084015160608301526000908360808460066107d05a03fa90508080156105885750806105c85760405162461bcd60e51b815260206004820152600d60248201526c1958cb5859190b59985a5b1959609a1b60448201526064016101a3565b61066c610b87565b50604080516080810182527f198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c28183019081527f1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed6060830152815281518083019092527f275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec82527f1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d60208381019190915281019190915290565b604080518082018252858152602080820185905282518084019093528583528201839052600091610753610bac565b60005b600281101561091857600061076c826006610ed9565b905084826002811061078057610780610ead565b60200201515183610792836000610ef8565b600c81106107a2576107a2610ead565b60200201528482600281106107b9576107b9610ead565b602002015160200151838260016107d09190610ef8565b600c81106107e0576107e0610ead565b60200201528382600281106107f7576107f7610ead565b602002015151518361080a836002610ef8565b600c811061081a5761081a610ead565b602002015283826002811061083157610831610ead565b602002015151600160200201518361084a836003610ef8565b600c811061085a5761085a610ead565b602002015283826002811061087157610871610ead565b60200201516020015160006002811061088c5761088c610ead565b60200201518361089d836004610ef8565b600c81106108ad576108ad610ead565b60200201528382600281106108c4576108c4610ead565b6020020151602001516001600281106108df576108df610ead565b6020020151836108f0836005610ef8565b600c811061090057610900610ead565b6020020152508061091081610f10565b915050610756565b50610921610bcb565b60006020826101808560086107d05a03fa90508080156105885750806109815760405162461bcd60e51b81526020600482015260156024820152741c185a5c9a5b99cb5bdc18dbd9194b59985a5b1959605a1b60448201526064016101a3565b5051151598975050505050505050565b6040805180820190915260008082526020820152600080806109c1600080516020610f2c83398151915286610e51565b90505b6109cd81610a21565b9093509150600080516020610f2c833981519152828309831415610a07576040805180820190915290815260208101919091529392505050565b600080516020610f2c8339815191526001820890506109c4565b60008080600080516020610f2c8339815191526003600080516020610f2c83398151915286600080516020610f2c833981519152888909090890506000610a97827f0c19139cb84c680a6e14116da060561765e05aa45a1c72a34f082305b61f3f52600080516020610f2c833981519152610aa3565b91959194509092505050565b600080610aae610bcb565b610ab6610be9565b602080825281810181905260408201819052606082018890526080820187905260a082018690528260c08360056107d05a03fa9250828015610588575082610b405760405162461bcd60e51b815260206004820152601a60248201527f424e3235342e6578704d6f643a2063616c6c206661696c75726500000000000060448201526064016101a3565b505195945050505050565b60405180606001604052806003906020820280368337509192915050565b60405180608001604052806004906020820280368337509192915050565b6040518060400160405280610b9a610c07565b8152602001610ba7610c07565b905290565b604051806101800160405280600c906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180604001604052806002906020820280368337509192915050565b634e487b7160e01b600052604160045260246000fd5b6040805190810167ffffffffffffffff81118282101715610c5e57610c5e610c25565b60405290565b600060408284031215610c7657600080fd5b6040516040810181811067ffffffffffffffff82111715610c9957610c99610c25565b604052823581526020928301359281019290925250919050565b600082601f830112610cc457600080fd5b610ccc610c3b565b806040840185811115610cde57600080fd5b845b81811015610cf8578035845260209384019301610ce0565b509095945050505050565b6000806000838503610100811215610d1a57600080fd5b610d248686610c64565b9350610d338660408701610c64565b92506080607f1982011215610d4757600080fd5b50610d50610c3b565b610d5d8660808701610cb3565b8152610d6c8660c08701610cb3565b6020820152809150509250925092565b600060208284031215610d8e57600080fd5b81356001600160a01b0381168114610da557600080fd5b9392505050565b8151815260208083015190820152604081016104f0565b600060208284031215610dd557600080fd5b5035919050565b8060005b6002811015610dff578151845260209384019390910190600101610de0565b50505050565b888152876020820152866040820152856060820152610e276080820186610ddc565b610e3460c0820185610ddc565b610100810192909252610120820152610140019695505050505050565b600082610e6e57634e487b7160e01b600052601260045260246000fd5b500690565b825181526020808401519082015260c08101610e93604083018451610ddc565b6020830151610ea56080840182610ddc565b509392505050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b6000816000190483118215151615610ef357610ef3610ec3565b500290565b60008219821115610f0b57610f0b610ec3565b500190565b6000600019821415610f2457610f24610ec3565b506001019056fe30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47a2646970667358221220977090042295bd82a2aa3914aa3ffcc077d62695f06783c3ff2d0cc5f80f6b2864736f6c634300080c0033",
}

// ContractBLSPublicKeyCompendiumABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractBLSPublicKeyCompendiumMetaData.ABI instead.
var ContractBLSPublicKeyCompendiumABI = ContractBLSPublicKeyCompendiumMetaData.ABI

// ContractBLSPublicKeyCompendiumBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractBLSPublicKeyCompendiumMetaData.Bin instead.
var ContractBLSPublicKeyCompendiumBin = ContractBLSPublicKeyCompendiumMetaData.Bin

// DeployContractBLSPublicKeyCompendium deploys a new Ethereum contract, binding an instance of ContractBLSPublicKeyCompendium to it.
func DeployContractBLSPublicKeyCompendium(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ContractBLSPublicKeyCompendium, error) {
	parsed, err := ContractBLSPublicKeyCompendiumMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractBLSPublicKeyCompendiumBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractBLSPublicKeyCompendium{ContractBLSPublicKeyCompendiumCaller: ContractBLSPublicKeyCompendiumCaller{contract: contract}, ContractBLSPublicKeyCompendiumTransactor: ContractBLSPublicKeyCompendiumTransactor{contract: contract}, ContractBLSPublicKeyCompendiumFilterer: ContractBLSPublicKeyCompendiumFilterer{contract: contract}}, nil
}

// ContractBLSPublicKeyCompendium is an auto generated Go binding around an Ethereum contract.
type ContractBLSPublicKeyCompendium struct {
	ContractBLSPublicKeyCompendiumCaller     // Read-only binding to the contract
	ContractBLSPublicKeyCompendiumTransactor // Write-only binding to the contract
	ContractBLSPublicKeyCompendiumFilterer   // Log filterer for contract events
}

// ContractBLSPublicKeyCompendiumCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractBLSPublicKeyCompendiumCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractBLSPublicKeyCompendiumTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractBLSPublicKeyCompendiumTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractBLSPublicKeyCompendiumFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractBLSPublicKeyCompendiumFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractBLSPublicKeyCompendiumSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractBLSPublicKeyCompendiumSession struct {
	Contract     *ContractBLSPublicKeyCompendium // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                   // Call options to use throughout this session
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// ContractBLSPublicKeyCompendiumCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractBLSPublicKeyCompendiumCallerSession struct {
	Contract *ContractBLSPublicKeyCompendiumCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                         // Call options to use throughout this session
}

// ContractBLSPublicKeyCompendiumTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractBLSPublicKeyCompendiumTransactorSession struct {
	Contract     *ContractBLSPublicKeyCompendiumTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                         // Transaction auth options to use throughout this session
}

// ContractBLSPublicKeyCompendiumRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractBLSPublicKeyCompendiumRaw struct {
	Contract *ContractBLSPublicKeyCompendium // Generic contract binding to access the raw methods on
}

// ContractBLSPublicKeyCompendiumCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractBLSPublicKeyCompendiumCallerRaw struct {
	Contract *ContractBLSPublicKeyCompendiumCaller // Generic read-only contract binding to access the raw methods on
}

// ContractBLSPublicKeyCompendiumTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractBLSPublicKeyCompendiumTransactorRaw struct {
	Contract *ContractBLSPublicKeyCompendiumTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractBLSPublicKeyCompendium creates a new instance of ContractBLSPublicKeyCompendium, bound to a specific deployed contract.
func NewContractBLSPublicKeyCompendium(address common.Address, backend bind.ContractBackend) (*ContractBLSPublicKeyCompendium, error) {
	contract, err := bindContractBLSPublicKeyCompendium(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractBLSPublicKeyCompendium{ContractBLSPublicKeyCompendiumCaller: ContractBLSPublicKeyCompendiumCaller{contract: contract}, ContractBLSPublicKeyCompendiumTransactor: ContractBLSPublicKeyCompendiumTransactor{contract: contract}, ContractBLSPublicKeyCompendiumFilterer: ContractBLSPublicKeyCompendiumFilterer{contract: contract}}, nil
}

// NewContractBLSPublicKeyCompendiumCaller creates a new read-only instance of ContractBLSPublicKeyCompendium, bound to a specific deployed contract.
func NewContractBLSPublicKeyCompendiumCaller(address common.Address, caller bind.ContractCaller) (*ContractBLSPublicKeyCompendiumCaller, error) {
	contract, err := bindContractBLSPublicKeyCompendium(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractBLSPublicKeyCompendiumCaller{contract: contract}, nil
}

// NewContractBLSPublicKeyCompendiumTransactor creates a new write-only instance of ContractBLSPublicKeyCompendium, bound to a specific deployed contract.
func NewContractBLSPublicKeyCompendiumTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractBLSPublicKeyCompendiumTransactor, error) {
	contract, err := bindContractBLSPublicKeyCompendium(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractBLSPublicKeyCompendiumTransactor{contract: contract}, nil
}

// NewContractBLSPublicKeyCompendiumFilterer creates a new log filterer instance of ContractBLSPublicKeyCompendium, bound to a specific deployed contract.
func NewContractBLSPublicKeyCompendiumFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractBLSPublicKeyCompendiumFilterer, error) {
	contract, err := bindContractBLSPublicKeyCompendium(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractBLSPublicKeyCompendiumFilterer{contract: contract}, nil
}

// bindContractBLSPublicKeyCompendium binds a generic wrapper to an already deployed contract.
func bindContractBLSPublicKeyCompendium(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractBLSPublicKeyCompendiumMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractBLSPublicKeyCompendium *ContractBLSPublicKeyCompendiumRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractBLSPublicKeyCompendium.Contract.ContractBLSPublicKeyCompendiumCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractBLSPublicKeyCompendium *ContractBLSPublicKeyCompendiumRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractBLSPublicKeyCompendium.Contract.ContractBLSPublicKeyCompendiumTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractBLSPublicKeyCompendium *ContractBLSPublicKeyCompendiumRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractBLSPublicKeyCompendium.Contract.ContractBLSPublicKeyCompendiumTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractBLSPublicKeyCompendium *ContractBLSPublicKeyCompendiumCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractBLSPublicKeyCompendium.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractBLSPublicKeyCompendium *ContractBLSPublicKeyCompendiumTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractBLSPublicKeyCompendium.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractBLSPublicKeyCompendium *ContractBLSPublicKeyCompendiumTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractBLSPublicKeyCompendium.Contract.contract.Transact(opts, method, params...)
}

// GetMessageHash is a free data retrieval call binding the contract method 0x1f5ac1b2.
//
// Solidity: function getMessageHash(address operator) view returns((uint256,uint256))
func (_ContractBLSPublicKeyCompendium *ContractBLSPublicKeyCompendiumCaller) GetMessageHash(opts *bind.CallOpts, operator common.Address) (BN254G1Point, error) {
	var out []interface{}
	err := _ContractBLSPublicKeyCompendium.contract.Call(opts, &out, "getMessageHash", operator)

	if err != nil {
		return *new(BN254G1Point), err
	}

	out0 := *abi.ConvertType(out[0], new(BN254G1Point)).(*BN254G1Point)

	return out0, err

}

// GetMessageHash is a free data retrieval call binding the contract method 0x1f5ac1b2.
//
// Solidity: function getMessageHash(address operator) view returns((uint256,uint256))
func (_ContractBLSPublicKeyCompendium *ContractBLSPublicKeyCompendiumSession) GetMessageHash(operator common.Address) (BN254G1Point, error) {
	return _ContractBLSPublicKeyCompendium.Contract.GetMessageHash(&_ContractBLSPublicKeyCompendium.CallOpts, operator)
}

// GetMessageHash is a free data retrieval call binding the contract method 0x1f5ac1b2.
//
// Solidity: function getMessageHash(address operator) view returns((uint256,uint256))
func (_ContractBLSPublicKeyCompendium *ContractBLSPublicKeyCompendiumCallerSession) GetMessageHash(operator common.Address) (BN254G1Point, error) {
	return _ContractBLSPublicKeyCompendium.Contract.GetMessageHash(&_ContractBLSPublicKeyCompendium.CallOpts, operator)
}

// OperatorToPubkeyHash is a free data retrieval call binding the contract method 0xde29fac0.
//
// Solidity: function operatorToPubkeyHash(address ) view returns(bytes32)
func (_ContractBLSPublicKeyCompendium *ContractBLSPublicKeyCompendiumCaller) OperatorToPubkeyHash(opts *bind.CallOpts, arg0 common.Address) ([32]byte, error) {
	var out []interface{}
	err := _ContractBLSPublicKeyCompendium.contract.Call(opts, &out, "operatorToPubkeyHash", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// OperatorToPubkeyHash is a free data retrieval call binding the contract method 0xde29fac0.
//
// Solidity: function operatorToPubkeyHash(address ) view returns(bytes32)
func (_ContractBLSPublicKeyCompendium *ContractBLSPublicKeyCompendiumSession) OperatorToPubkeyHash(arg0 common.Address) ([32]byte, error) {
	return _ContractBLSPublicKeyCompendium.Contract.OperatorToPubkeyHash(&_ContractBLSPublicKeyCompendium.CallOpts, arg0)
}

// OperatorToPubkeyHash is a free data retrieval call binding the contract method 0xde29fac0.
//
// Solidity: function operatorToPubkeyHash(address ) view returns(bytes32)
func (_ContractBLSPublicKeyCompendium *ContractBLSPublicKeyCompendiumCallerSession) OperatorToPubkeyHash(arg0 common.Address) ([32]byte, error) {
	return _ContractBLSPublicKeyCompendium.Contract.OperatorToPubkeyHash(&_ContractBLSPublicKeyCompendium.CallOpts, arg0)
}

// PubkeyHashToOperator is a free data retrieval call binding the contract method 0xe8bb9ae6.
//
// Solidity: function pubkeyHashToOperator(bytes32 ) view returns(address)
func (_ContractBLSPublicKeyCompendium *ContractBLSPublicKeyCompendiumCaller) PubkeyHashToOperator(opts *bind.CallOpts, arg0 [32]byte) (common.Address, error) {
	var out []interface{}
	err := _ContractBLSPublicKeyCompendium.contract.Call(opts, &out, "pubkeyHashToOperator", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PubkeyHashToOperator is a free data retrieval call binding the contract method 0xe8bb9ae6.
//
// Solidity: function pubkeyHashToOperator(bytes32 ) view returns(address)
func (_ContractBLSPublicKeyCompendium *ContractBLSPublicKeyCompendiumSession) PubkeyHashToOperator(arg0 [32]byte) (common.Address, error) {
	return _ContractBLSPublicKeyCompendium.Contract.PubkeyHashToOperator(&_ContractBLSPublicKeyCompendium.CallOpts, arg0)
}

// PubkeyHashToOperator is a free data retrieval call binding the contract method 0xe8bb9ae6.
//
// Solidity: function pubkeyHashToOperator(bytes32 ) view returns(address)
func (_ContractBLSPublicKeyCompendium *ContractBLSPublicKeyCompendiumCallerSession) PubkeyHashToOperator(arg0 [32]byte) (common.Address, error) {
	return _ContractBLSPublicKeyCompendium.Contract.PubkeyHashToOperator(&_ContractBLSPublicKeyCompendium.CallOpts, arg0)
}

// RegisterBLSPublicKey is a paid mutator transaction binding the contract method 0x161a334d.
//
// Solidity: function registerBLSPublicKey((uint256,uint256) signedMessageHash, (uint256,uint256) pubkeyG1, (uint256[2],uint256[2]) pubkeyG2) returns()
func (_ContractBLSPublicKeyCompendium *ContractBLSPublicKeyCompendiumTransactor) RegisterBLSPublicKey(opts *bind.TransactOpts, signedMessageHash BN254G1Point, pubkeyG1 BN254G1Point, pubkeyG2 BN254G2Point) (*types.Transaction, error) {
	return _ContractBLSPublicKeyCompendium.contract.Transact(opts, "registerBLSPublicKey", signedMessageHash, pubkeyG1, pubkeyG2)
}

// RegisterBLSPublicKey is a paid mutator transaction binding the contract method 0x161a334d.
//
// Solidity: function registerBLSPublicKey((uint256,uint256) signedMessageHash, (uint256,uint256) pubkeyG1, (uint256[2],uint256[2]) pubkeyG2) returns()
func (_ContractBLSPublicKeyCompendium *ContractBLSPublicKeyCompendiumSession) RegisterBLSPublicKey(signedMessageHash BN254G1Point, pubkeyG1 BN254G1Point, pubkeyG2 BN254G2Point) (*types.Transaction, error) {
	return _ContractBLSPublicKeyCompendium.Contract.RegisterBLSPublicKey(&_ContractBLSPublicKeyCompendium.TransactOpts, signedMessageHash, pubkeyG1, pubkeyG2)
}

// RegisterBLSPublicKey is a paid mutator transaction binding the contract method 0x161a334d.
//
// Solidity: function registerBLSPublicKey((uint256,uint256) signedMessageHash, (uint256,uint256) pubkeyG1, (uint256[2],uint256[2]) pubkeyG2) returns()
func (_ContractBLSPublicKeyCompendium *ContractBLSPublicKeyCompendiumTransactorSession) RegisterBLSPublicKey(signedMessageHash BN254G1Point, pubkeyG1 BN254G1Point, pubkeyG2 BN254G2Point) (*types.Transaction, error) {
	return _ContractBLSPublicKeyCompendium.Contract.RegisterBLSPublicKey(&_ContractBLSPublicKeyCompendium.TransactOpts, signedMessageHash, pubkeyG1, pubkeyG2)
}

// ContractBLSPublicKeyCompendiumNewPubkeyRegistrationIterator is returned from FilterNewPubkeyRegistration and is used to iterate over the raw logs and unpacked data for NewPubkeyRegistration events raised by the ContractBLSPublicKeyCompendium contract.
type ContractBLSPublicKeyCompendiumNewPubkeyRegistrationIterator struct {
	Event *ContractBLSPublicKeyCompendiumNewPubkeyRegistration // Event containing the contract specifics and raw log

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
func (it *ContractBLSPublicKeyCompendiumNewPubkeyRegistrationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractBLSPublicKeyCompendiumNewPubkeyRegistration)
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
		it.Event = new(ContractBLSPublicKeyCompendiumNewPubkeyRegistration)
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
func (it *ContractBLSPublicKeyCompendiumNewPubkeyRegistrationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractBLSPublicKeyCompendiumNewPubkeyRegistrationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractBLSPublicKeyCompendiumNewPubkeyRegistration represents a NewPubkeyRegistration event raised by the ContractBLSPublicKeyCompendium contract.
type ContractBLSPublicKeyCompendiumNewPubkeyRegistration struct {
	Operator common.Address
	PubkeyG1 BN254G1Point
	PubkeyG2 BN254G2Point
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterNewPubkeyRegistration is a free log retrieval operation binding the contract event 0xe3fb6613af2e8930cf85d47fcf6db10192224a64c6cbe8023e0eee1ba3828041.
//
// Solidity: event NewPubkeyRegistration(address indexed operator, (uint256,uint256) pubkeyG1, (uint256[2],uint256[2]) pubkeyG2)
func (_ContractBLSPublicKeyCompendium *ContractBLSPublicKeyCompendiumFilterer) FilterNewPubkeyRegistration(opts *bind.FilterOpts, operator []common.Address) (*ContractBLSPublicKeyCompendiumNewPubkeyRegistrationIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ContractBLSPublicKeyCompendium.contract.FilterLogs(opts, "NewPubkeyRegistration", operatorRule)
	if err != nil {
		return nil, err
	}
	return &ContractBLSPublicKeyCompendiumNewPubkeyRegistrationIterator{contract: _ContractBLSPublicKeyCompendium.contract, event: "NewPubkeyRegistration", logs: logs, sub: sub}, nil
}

// WatchNewPubkeyRegistration is a free log subscription operation binding the contract event 0xe3fb6613af2e8930cf85d47fcf6db10192224a64c6cbe8023e0eee1ba3828041.
//
// Solidity: event NewPubkeyRegistration(address indexed operator, (uint256,uint256) pubkeyG1, (uint256[2],uint256[2]) pubkeyG2)
func (_ContractBLSPublicKeyCompendium *ContractBLSPublicKeyCompendiumFilterer) WatchNewPubkeyRegistration(opts *bind.WatchOpts, sink chan<- *ContractBLSPublicKeyCompendiumNewPubkeyRegistration, operator []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ContractBLSPublicKeyCompendium.contract.WatchLogs(opts, "NewPubkeyRegistration", operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractBLSPublicKeyCompendiumNewPubkeyRegistration)
				if err := _ContractBLSPublicKeyCompendium.contract.UnpackLog(event, "NewPubkeyRegistration", log); err != nil {
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

// ParseNewPubkeyRegistration is a log parse operation binding the contract event 0xe3fb6613af2e8930cf85d47fcf6db10192224a64c6cbe8023e0eee1ba3828041.
//
// Solidity: event NewPubkeyRegistration(address indexed operator, (uint256,uint256) pubkeyG1, (uint256[2],uint256[2]) pubkeyG2)
func (_ContractBLSPublicKeyCompendium *ContractBLSPublicKeyCompendiumFilterer) ParseNewPubkeyRegistration(log types.Log) (*ContractBLSPublicKeyCompendiumNewPubkeyRegistration, error) {
	event := new(ContractBLSPublicKeyCompendiumNewPubkeyRegistration)
	if err := _ContractBLSPublicKeyCompendium.contract.UnpackLog(event, "NewPubkeyRegistration", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
