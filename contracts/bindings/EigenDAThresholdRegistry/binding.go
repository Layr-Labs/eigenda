// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractEigenDAThresholdRegistry

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

// SecurityThresholds is an auto generated low-level Go binding around an user-defined struct.
type SecurityThresholds struct {
	ConfirmationThreshold uint8
	AdversaryThreshold    uint8
}

// VersionedBlobParams is an auto generated low-level Go binding around an user-defined struct.
type VersionedBlobParams struct {
	MaxNumOperators uint32
	NumChunks       uint32
	CodingRate      uint8
}

// ContractEigenDAThresholdRegistryMetaData contains all meta data concerning the ContractEigenDAThresholdRegistry contract.
var ContractEigenDAThresholdRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addVersionedBlobParams\",\"inputs\":[{\"name\":\"_versionedBlobParams\",\"type\":\"tuple\",\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getBlobParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIsQuorumRequired\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumAdversaryThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumConfirmationThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_quorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_quorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_quorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_versionedBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structVersionedBlobParams[]\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"nextBlobVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumAdversaryThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumConfirmationThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"versionedBlobParams\",\"inputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"DefaultSecurityThresholdsV2Updated\",\"inputs\":[{\"name\":\"previousDefaultSecurityThresholdsV2\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"newDefaultSecurityThresholdsV2\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumAdversaryThresholdPercentagesUpdated\",\"inputs\":[{\"name\":\"previousQuorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumConfirmationThresholdPercentagesUpdated\",\"inputs\":[{\"name\":\"previousQuorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumNumbersRequiredUpdated\",\"inputs\":[{\"name\":\"previousQuorumNumbersRequired\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumNumbersRequired\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VersionedBlobParamsAdded\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"indexed\":true,\"internalType\":\"uint16\"},{\"name\":\"versionedBlobParams\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"anonymous\":false}]",
	Bin: "0x608060405234801561001057600080fd5b5061001961001e565b6100de565b603254610100900460ff161561008a5760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840160405180910390fd5b60325460ff90811610156100dc576032805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b610f1e806100ed6000396000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c80638a4769821161008c578063e15234ff11610066578063e15234ff14610265578063ee6c3bcf1461026d578063f2fde38b14610280578063f74e363c1461029357600080fd5b80638a4769821461022f5780638da5cb5b14610242578063bafa91071461025d57600080fd5b806332430f14116100c857806332430f14146101dc578063715018a6146101fd5780638491bad6146102075780638687feae1461021a57600080fd5b8063048886d2146100ef5780631429c7c2146101175780632ecfe72b1461013c575b600080fd5b6101026100fd366004610af5565b6102f8565b60405190151581526020015b60405180910390f35b61012a610125366004610af5565b6103a2565b60405160ff909116815260200161010e565b6101ac61014a366004610b17565b60408051606080820183526000808352602080840182905292840181905261ffff9490941684526004825292829020825193840183525463ffffffff808216855264010000000082041691840191909152600160401b900460ff169082015290565b60408051825163ffffffff9081168252602080850151909116908201529181015160ff169082015260600161010e565b6003546101ea9061ffff1681565b60405161ffff909116815260200161010e565b610205610410565b005b610205610215366004610c8c565b610424565b6102226105bc565b60405161010e9190610db8565b6101ea61023d366004610e0d565b61064a565b6065546040516001600160a01b03909116815260200161010e565b610222610663565b610222610670565b61012a61027b366004610af5565b61067d565b61020561028e366004610e29565b6106a9565b6102d26102a1366004610b17565b60046020526000908152604090205463ffffffff80821691640100000000810490911690600160401b900460ff1683565b6040805163ffffffff948516815293909216602084015260ff169082015260600161010e565b600080600160ff84161b9050806103986002805461031590610e44565b80601f016020809104026020016040519081016040528092919081815260200182805461034190610e44565b801561038e5780601f106103635761010080835404028352916020019161038e565b820191906000526020600020905b81548152906001019060200180831161037157829003601f168201915b5050505050610722565b9091161492915050565b60008160ff16600180546103b590610e44565b9050111561040b5760018260ff1681546103ce90610e44565b81106103dc576103dc610e7f565b8154600116156103fb5790600052602060002090602091828204019190065b9054901a600160f81b0260f81c90505b919050565b6104186108af565b6104226000610909565b565b603254610100900460ff16158080156104445750603254600160ff909116105b8061045e5750303b15801561045e575060325460ff166001145b6104c65760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b60648201526084015b60405180910390fd5b6032805460ff1916600117905580156104e9576032805461ff0019166101001790555b6104f286610909565b8451610505906000906020880190610a4b565b508351610519906001906020870190610a4b565b50825161052d906002906020860190610a4b565b5060005b825181101561056d5761055c83828151811061054f5761054f610e7f565b602002602001015161095b565b5061056681610eab565b9050610531565b5080156105b4576032805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b505050505050565b600080546105c990610e44565b80601f01602080910402602001604051908101604052809291908181526020018280546105f590610e44565b80156106425780601f1061061757610100808354040283529160200191610642565b820191906000526020600020905b81548152906001019060200180831161062557829003601f168201915b505050505081565b60006106546108af565b61065d8261095b565b92915050565b600180546105c990610e44565b600280546105c990610e44565b60008160ff166000805461069090610e44565b9050111561040b5760008260ff1681546103ce90610e44565b6106b16108af565b6001600160a01b0381166107165760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b60648201526084016104bd565b61071f81610909565b50565b6000610100825111156107ab5760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a4016104bd565b81516107b957506000919050565b600080836000815181106107cf576107cf610e7f565b0160200151600160f89190911c81901b92505b84518110156108a6578481815181106107fd576107fd610e7f565b0160200151600160f89190911c1b91508282116108925760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a4016104bd565b9181179161089f81610eab565b90506107e2565b50909392505050565b6065546001600160a01b031633146104225760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016104bd565b606580546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b6003805461ffff90811660009081526004602090815260408083208651815488850180518a8601805163ffffffff95861667ffffffffffffffff199095168517640100000000938716939093029290921768ff00000000000000001916600160401b60ff9384160217909555985485519283529051909216948101949094529051909516908201529092909116907fdbee9d337a6e5fde30966e157673aaeeb6a0134afaf774a4b6979b7c79d07da49060600160405180910390a26003805461ffff16906000610a2a83610ec6565b91906101000a81548161ffff021916908361ffff1602179055509050919050565b828054610a5790610e44565b90600052602060002090601f016020900481019282610a795760008555610abf565b82601f10610a9257805160ff1916838001178555610abf565b82800160010185558215610abf579182015b82811115610abf578251825591602001919060010190610aa4565b50610acb929150610acf565b5090565b5b80821115610acb5760008155600101610ad0565b803560ff8116811461040b57600080fd5b600060208284031215610b0757600080fd5b610b1082610ae4565b9392505050565b600060208284031215610b2957600080fd5b813561ffff81168114610b1057600080fd5b80356001600160a01b038116811461040b57600080fd5b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff81118282101715610b9157610b91610b52565b604052919050565b600082601f830112610baa57600080fd5b813567ffffffffffffffff811115610bc457610bc4610b52565b610bd7601f8201601f1916602001610b68565b818152846020838601011115610bec57600080fd5b816020850160208301376000918101602001919091529392505050565b803563ffffffff8116811461040b57600080fd5b600060608284031215610c2f57600080fd5b6040516060810181811067ffffffffffffffff82111715610c5257610c52610b52565b604052905080610c6183610c09565b8152610c6f60208401610c09565b6020820152610c8060408401610ae4565b60408201525092915050565b600080600080600060a08688031215610ca457600080fd5b610cad86610b3b565b945060208087013567ffffffffffffffff80821115610ccb57600080fd5b610cd78a838b01610b99565b96506040890135915080821115610ced57600080fd5b610cf98a838b01610b99565b9550606091508189013581811115610d1057600080fd5b610d1c8b828c01610b99565b955050608089013581811115610d3157600080fd5b8901601f81018b13610d4257600080fd5b803582811115610d5457610d54610b52565b610d62858260051b01610b68565b818152858101935090840282018501908c821115610d7f57600080fd5b918501915b81831015610da557610d968d84610c1d565b84529285019291840191610d84565b8096505050505050509295509295909350565b600060208083528351808285015260005b81811015610de557858101830151858201604001528201610dc9565b81811115610df7576000604083870101525b50601f01601f1916929092016040019392505050565b600060608284031215610e1f57600080fd5b610b108383610c1d565b600060208284031215610e3b57600080fd5b610b1082610b3b565b600181811c90821680610e5857607f821691505b60208210811415610e7957634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b6000600019821415610ebf57610ebf610e95565b5060010190565b600061ffff80831681811415610ede57610ede610e95565b600101939250505056fea26469706673582212207bc016ad7f136621f1a113f814f3932c0410c559a5a249ee557b6171aa88746d64736f6c634300080c0033",
}

// ContractEigenDAThresholdRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEigenDAThresholdRegistryMetaData.ABI instead.
var ContractEigenDAThresholdRegistryABI = ContractEigenDAThresholdRegistryMetaData.ABI

// ContractEigenDAThresholdRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractEigenDAThresholdRegistryMetaData.Bin instead.
var ContractEigenDAThresholdRegistryBin = ContractEigenDAThresholdRegistryMetaData.Bin

// DeployContractEigenDAThresholdRegistry deploys a new Ethereum contract, binding an instance of ContractEigenDAThresholdRegistry to it.
func DeployContractEigenDAThresholdRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ContractEigenDAThresholdRegistry, error) {
	parsed, err := ContractEigenDAThresholdRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractEigenDAThresholdRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractEigenDAThresholdRegistry{ContractEigenDAThresholdRegistryCaller: ContractEigenDAThresholdRegistryCaller{contract: contract}, ContractEigenDAThresholdRegistryTransactor: ContractEigenDAThresholdRegistryTransactor{contract: contract}, ContractEigenDAThresholdRegistryFilterer: ContractEigenDAThresholdRegistryFilterer{contract: contract}}, nil
}

// ContractEigenDAThresholdRegistry is an auto generated Go binding around an Ethereum contract.
type ContractEigenDAThresholdRegistry struct {
	ContractEigenDAThresholdRegistryCaller     // Read-only binding to the contract
	ContractEigenDAThresholdRegistryTransactor // Write-only binding to the contract
	ContractEigenDAThresholdRegistryFilterer   // Log filterer for contract events
}

// ContractEigenDAThresholdRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractEigenDAThresholdRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDAThresholdRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractEigenDAThresholdRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDAThresholdRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractEigenDAThresholdRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDAThresholdRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractEigenDAThresholdRegistrySession struct {
	Contract     *ContractEigenDAThresholdRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                     // Call options to use throughout this session
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// ContractEigenDAThresholdRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractEigenDAThresholdRegistryCallerSession struct {
	Contract *ContractEigenDAThresholdRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                           // Call options to use throughout this session
}

// ContractEigenDAThresholdRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractEigenDAThresholdRegistryTransactorSession struct {
	Contract     *ContractEigenDAThresholdRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                           // Transaction auth options to use throughout this session
}

// ContractEigenDAThresholdRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractEigenDAThresholdRegistryRaw struct {
	Contract *ContractEigenDAThresholdRegistry // Generic contract binding to access the raw methods on
}

// ContractEigenDAThresholdRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractEigenDAThresholdRegistryCallerRaw struct {
	Contract *ContractEigenDAThresholdRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// ContractEigenDAThresholdRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractEigenDAThresholdRegistryTransactorRaw struct {
	Contract *ContractEigenDAThresholdRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractEigenDAThresholdRegistry creates a new instance of ContractEigenDAThresholdRegistry, bound to a specific deployed contract.
func NewContractEigenDAThresholdRegistry(address common.Address, backend bind.ContractBackend) (*ContractEigenDAThresholdRegistry, error) {
	contract, err := bindContractEigenDAThresholdRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAThresholdRegistry{ContractEigenDAThresholdRegistryCaller: ContractEigenDAThresholdRegistryCaller{contract: contract}, ContractEigenDAThresholdRegistryTransactor: ContractEigenDAThresholdRegistryTransactor{contract: contract}, ContractEigenDAThresholdRegistryFilterer: ContractEigenDAThresholdRegistryFilterer{contract: contract}}, nil
}

// NewContractEigenDAThresholdRegistryCaller creates a new read-only instance of ContractEigenDAThresholdRegistry, bound to a specific deployed contract.
func NewContractEigenDAThresholdRegistryCaller(address common.Address, caller bind.ContractCaller) (*ContractEigenDAThresholdRegistryCaller, error) {
	contract, err := bindContractEigenDAThresholdRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAThresholdRegistryCaller{contract: contract}, nil
}

// NewContractEigenDAThresholdRegistryTransactor creates a new write-only instance of ContractEigenDAThresholdRegistry, bound to a specific deployed contract.
func NewContractEigenDAThresholdRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractEigenDAThresholdRegistryTransactor, error) {
	contract, err := bindContractEigenDAThresholdRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAThresholdRegistryTransactor{contract: contract}, nil
}

// NewContractEigenDAThresholdRegistryFilterer creates a new log filterer instance of ContractEigenDAThresholdRegistry, bound to a specific deployed contract.
func NewContractEigenDAThresholdRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractEigenDAThresholdRegistryFilterer, error) {
	contract, err := bindContractEigenDAThresholdRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAThresholdRegistryFilterer{contract: contract}, nil
}

// bindContractEigenDAThresholdRegistry binds a generic wrapper to an already deployed contract.
func bindContractEigenDAThresholdRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractEigenDAThresholdRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDAThresholdRegistry.Contract.ContractEigenDAThresholdRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.ContractEigenDAThresholdRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.ContractEigenDAThresholdRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDAThresholdRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.contract.Transact(opts, method, params...)
}

// GetBlobParams is a free data retrieval call binding the contract method 0x2ecfe72b.
//
// Solidity: function getBlobParams(uint16 version) view returns((uint32,uint32,uint8))
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCaller) GetBlobParams(opts *bind.CallOpts, version uint16) (VersionedBlobParams, error) {
	var out []interface{}
	err := _ContractEigenDAThresholdRegistry.contract.Call(opts, &out, "getBlobParams", version)

	if err != nil {
		return *new(VersionedBlobParams), err
	}

	out0 := *abi.ConvertType(out[0], new(VersionedBlobParams)).(*VersionedBlobParams)

	return out0, err

}

// GetBlobParams is a free data retrieval call binding the contract method 0x2ecfe72b.
//
// Solidity: function getBlobParams(uint16 version) view returns((uint32,uint32,uint8))
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) GetBlobParams(version uint16) (VersionedBlobParams, error) {
	return _ContractEigenDAThresholdRegistry.Contract.GetBlobParams(&_ContractEigenDAThresholdRegistry.CallOpts, version)
}

// GetBlobParams is a free data retrieval call binding the contract method 0x2ecfe72b.
//
// Solidity: function getBlobParams(uint16 version) view returns((uint32,uint32,uint8))
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCallerSession) GetBlobParams(version uint16) (VersionedBlobParams, error) {
	return _ContractEigenDAThresholdRegistry.Contract.GetBlobParams(&_ContractEigenDAThresholdRegistry.CallOpts, version)
}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCaller) GetIsQuorumRequired(opts *bind.CallOpts, quorumNumber uint8) (bool, error) {
	var out []interface{}
	err := _ContractEigenDAThresholdRegistry.contract.Call(opts, &out, "getIsQuorumRequired", quorumNumber)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) GetIsQuorumRequired(quorumNumber uint8) (bool, error) {
	return _ContractEigenDAThresholdRegistry.Contract.GetIsQuorumRequired(&_ContractEigenDAThresholdRegistry.CallOpts, quorumNumber)
}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCallerSession) GetIsQuorumRequired(quorumNumber uint8) (bool, error) {
	return _ContractEigenDAThresholdRegistry.Contract.GetIsQuorumRequired(&_ContractEigenDAThresholdRegistry.CallOpts, quorumNumber)
}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8 adversaryThresholdPercentage)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCaller) GetQuorumAdversaryThresholdPercentage(opts *bind.CallOpts, quorumNumber uint8) (uint8, error) {
	var out []interface{}
	err := _ContractEigenDAThresholdRegistry.contract.Call(opts, &out, "getQuorumAdversaryThresholdPercentage", quorumNumber)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8 adversaryThresholdPercentage)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) GetQuorumAdversaryThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDAThresholdRegistry.Contract.GetQuorumAdversaryThresholdPercentage(&_ContractEigenDAThresholdRegistry.CallOpts, quorumNumber)
}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8 adversaryThresholdPercentage)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCallerSession) GetQuorumAdversaryThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDAThresholdRegistry.Contract.GetQuorumAdversaryThresholdPercentage(&_ContractEigenDAThresholdRegistry.CallOpts, quorumNumber)
}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8 confirmationThresholdPercentage)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCaller) GetQuorumConfirmationThresholdPercentage(opts *bind.CallOpts, quorumNumber uint8) (uint8, error) {
	var out []interface{}
	err := _ContractEigenDAThresholdRegistry.contract.Call(opts, &out, "getQuorumConfirmationThresholdPercentage", quorumNumber)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8 confirmationThresholdPercentage)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) GetQuorumConfirmationThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDAThresholdRegistry.Contract.GetQuorumConfirmationThresholdPercentage(&_ContractEigenDAThresholdRegistry.CallOpts, quorumNumber)
}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8 confirmationThresholdPercentage)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCallerSession) GetQuorumConfirmationThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDAThresholdRegistry.Contract.GetQuorumConfirmationThresholdPercentage(&_ContractEigenDAThresholdRegistry.CallOpts, quorumNumber)
}

// NextBlobVersion is a free data retrieval call binding the contract method 0x32430f14.
//
// Solidity: function nextBlobVersion() view returns(uint16)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCaller) NextBlobVersion(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _ContractEigenDAThresholdRegistry.contract.Call(opts, &out, "nextBlobVersion")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// NextBlobVersion is a free data retrieval call binding the contract method 0x32430f14.
//
// Solidity: function nextBlobVersion() view returns(uint16)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) NextBlobVersion() (uint16, error) {
	return _ContractEigenDAThresholdRegistry.Contract.NextBlobVersion(&_ContractEigenDAThresholdRegistry.CallOpts)
}

// NextBlobVersion is a free data retrieval call binding the contract method 0x32430f14.
//
// Solidity: function nextBlobVersion() view returns(uint16)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCallerSession) NextBlobVersion() (uint16, error) {
	return _ContractEigenDAThresholdRegistry.Contract.NextBlobVersion(&_ContractEigenDAThresholdRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDAThresholdRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) Owner() (common.Address, error) {
	return _ContractEigenDAThresholdRegistry.Contract.Owner(&_ContractEigenDAThresholdRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCallerSession) Owner() (common.Address, error) {
	return _ContractEigenDAThresholdRegistry.Contract.Owner(&_ContractEigenDAThresholdRegistry.CallOpts)
}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCaller) QuorumAdversaryThresholdPercentages(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDAThresholdRegistry.contract.Call(opts, &out, "quorumAdversaryThresholdPercentages")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) QuorumAdversaryThresholdPercentages() ([]byte, error) {
	return _ContractEigenDAThresholdRegistry.Contract.QuorumAdversaryThresholdPercentages(&_ContractEigenDAThresholdRegistry.CallOpts)
}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCallerSession) QuorumAdversaryThresholdPercentages() ([]byte, error) {
	return _ContractEigenDAThresholdRegistry.Contract.QuorumAdversaryThresholdPercentages(&_ContractEigenDAThresholdRegistry.CallOpts)
}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCaller) QuorumConfirmationThresholdPercentages(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDAThresholdRegistry.contract.Call(opts, &out, "quorumConfirmationThresholdPercentages")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) QuorumConfirmationThresholdPercentages() ([]byte, error) {
	return _ContractEigenDAThresholdRegistry.Contract.QuorumConfirmationThresholdPercentages(&_ContractEigenDAThresholdRegistry.CallOpts)
}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCallerSession) QuorumConfirmationThresholdPercentages() ([]byte, error) {
	return _ContractEigenDAThresholdRegistry.Contract.QuorumConfirmationThresholdPercentages(&_ContractEigenDAThresholdRegistry.CallOpts)
}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCaller) QuorumNumbersRequired(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDAThresholdRegistry.contract.Call(opts, &out, "quorumNumbersRequired")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) QuorumNumbersRequired() ([]byte, error) {
	return _ContractEigenDAThresholdRegistry.Contract.QuorumNumbersRequired(&_ContractEigenDAThresholdRegistry.CallOpts)
}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCallerSession) QuorumNumbersRequired() ([]byte, error) {
	return _ContractEigenDAThresholdRegistry.Contract.QuorumNumbersRequired(&_ContractEigenDAThresholdRegistry.CallOpts)
}

// VersionedBlobParams is a free data retrieval call binding the contract method 0xf74e363c.
//
// Solidity: function versionedBlobParams(uint16 ) view returns(uint32 maxNumOperators, uint32 numChunks, uint8 codingRate)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCaller) VersionedBlobParams(opts *bind.CallOpts, arg0 uint16) (struct {
	MaxNumOperators uint32
	NumChunks       uint32
	CodingRate      uint8
}, error) {
	var out []interface{}
	err := _ContractEigenDAThresholdRegistry.contract.Call(opts, &out, "versionedBlobParams", arg0)

	outstruct := new(struct {
		MaxNumOperators uint32
		NumChunks       uint32
		CodingRate      uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.MaxNumOperators = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.NumChunks = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.CodingRate = *abi.ConvertType(out[2], new(uint8)).(*uint8)

	return *outstruct, err

}

// VersionedBlobParams is a free data retrieval call binding the contract method 0xf74e363c.
//
// Solidity: function versionedBlobParams(uint16 ) view returns(uint32 maxNumOperators, uint32 numChunks, uint8 codingRate)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) VersionedBlobParams(arg0 uint16) (struct {
	MaxNumOperators uint32
	NumChunks       uint32
	CodingRate      uint8
}, error) {
	return _ContractEigenDAThresholdRegistry.Contract.VersionedBlobParams(&_ContractEigenDAThresholdRegistry.CallOpts, arg0)
}

// VersionedBlobParams is a free data retrieval call binding the contract method 0xf74e363c.
//
// Solidity: function versionedBlobParams(uint16 ) view returns(uint32 maxNumOperators, uint32 numChunks, uint8 codingRate)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCallerSession) VersionedBlobParams(arg0 uint16) (struct {
	MaxNumOperators uint32
	NumChunks       uint32
	CodingRate      uint8
}, error) {
	return _ContractEigenDAThresholdRegistry.Contract.VersionedBlobParams(&_ContractEigenDAThresholdRegistry.CallOpts, arg0)
}

// AddVersionedBlobParams is a paid mutator transaction binding the contract method 0x8a476982.
//
// Solidity: function addVersionedBlobParams((uint32,uint32,uint8) _versionedBlobParams) returns(uint16)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactor) AddVersionedBlobParams(opts *bind.TransactOpts, _versionedBlobParams VersionedBlobParams) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.contract.Transact(opts, "addVersionedBlobParams", _versionedBlobParams)
}

// AddVersionedBlobParams is a paid mutator transaction binding the contract method 0x8a476982.
//
// Solidity: function addVersionedBlobParams((uint32,uint32,uint8) _versionedBlobParams) returns(uint16)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) AddVersionedBlobParams(_versionedBlobParams VersionedBlobParams) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.AddVersionedBlobParams(&_ContractEigenDAThresholdRegistry.TransactOpts, _versionedBlobParams)
}

// AddVersionedBlobParams is a paid mutator transaction binding the contract method 0x8a476982.
//
// Solidity: function addVersionedBlobParams((uint32,uint32,uint8) _versionedBlobParams) returns(uint16)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactorSession) AddVersionedBlobParams(_versionedBlobParams VersionedBlobParams) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.AddVersionedBlobParams(&_ContractEigenDAThresholdRegistry.TransactOpts, _versionedBlobParams)
}

// Initialize is a paid mutator transaction binding the contract method 0x8491bad6.
//
// Solidity: function initialize(address _initialOwner, bytes _quorumAdversaryThresholdPercentages, bytes _quorumConfirmationThresholdPercentages, bytes _quorumNumbersRequired, (uint32,uint32,uint8)[] _versionedBlobParams) returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactor) Initialize(opts *bind.TransactOpts, _initialOwner common.Address, _quorumAdversaryThresholdPercentages []byte, _quorumConfirmationThresholdPercentages []byte, _quorumNumbersRequired []byte, _versionedBlobParams []VersionedBlobParams) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.contract.Transact(opts, "initialize", _initialOwner, _quorumAdversaryThresholdPercentages, _quorumConfirmationThresholdPercentages, _quorumNumbersRequired, _versionedBlobParams)
}

// Initialize is a paid mutator transaction binding the contract method 0x8491bad6.
//
// Solidity: function initialize(address _initialOwner, bytes _quorumAdversaryThresholdPercentages, bytes _quorumConfirmationThresholdPercentages, bytes _quorumNumbersRequired, (uint32,uint32,uint8)[] _versionedBlobParams) returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) Initialize(_initialOwner common.Address, _quorumAdversaryThresholdPercentages []byte, _quorumConfirmationThresholdPercentages []byte, _quorumNumbersRequired []byte, _versionedBlobParams []VersionedBlobParams) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.Initialize(&_ContractEigenDAThresholdRegistry.TransactOpts, _initialOwner, _quorumAdversaryThresholdPercentages, _quorumConfirmationThresholdPercentages, _quorumNumbersRequired, _versionedBlobParams)
}

// Initialize is a paid mutator transaction binding the contract method 0x8491bad6.
//
// Solidity: function initialize(address _initialOwner, bytes _quorumAdversaryThresholdPercentages, bytes _quorumConfirmationThresholdPercentages, bytes _quorumNumbersRequired, (uint32,uint32,uint8)[] _versionedBlobParams) returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactorSession) Initialize(_initialOwner common.Address, _quorumAdversaryThresholdPercentages []byte, _quorumConfirmationThresholdPercentages []byte, _quorumNumbersRequired []byte, _versionedBlobParams []VersionedBlobParams) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.Initialize(&_ContractEigenDAThresholdRegistry.TransactOpts, _initialOwner, _quorumAdversaryThresholdPercentages, _quorumConfirmationThresholdPercentages, _quorumNumbersRequired, _versionedBlobParams)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.RenounceOwnership(&_ContractEigenDAThresholdRegistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.RenounceOwnership(&_ContractEigenDAThresholdRegistry.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.TransferOwnership(&_ContractEigenDAThresholdRegistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.TransferOwnership(&_ContractEigenDAThresholdRegistry.TransactOpts, newOwner)
}

// ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2UpdatedIterator is returned from FilterDefaultSecurityThresholdsV2Updated and is used to iterate over the raw logs and unpacked data for DefaultSecurityThresholdsV2Updated events raised by the ContractEigenDAThresholdRegistry contract.
type ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2UpdatedIterator struct {
	Event *ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2Updated // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2UpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2Updated)
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
		it.Event = new(ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2Updated)
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
func (it *ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2UpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2UpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2Updated represents a DefaultSecurityThresholdsV2Updated event raised by the ContractEigenDAThresholdRegistry contract.
type ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2Updated struct {
	PreviousDefaultSecurityThresholdsV2 SecurityThresholds
	NewDefaultSecurityThresholdsV2      SecurityThresholds
	Raw                                 types.Log // Blockchain specific contextual infos
}

// FilterDefaultSecurityThresholdsV2Updated is a free log retrieval operation binding the contract event 0xfe03afd62c76a6aed7376ae995cc55d073ba9d83d83ac8efc5446f8da4d50997.
//
// Solidity: event DefaultSecurityThresholdsV2Updated((uint8,uint8) previousDefaultSecurityThresholdsV2, (uint8,uint8) newDefaultSecurityThresholdsV2)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) FilterDefaultSecurityThresholdsV2Updated(opts *bind.FilterOpts) (*ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2UpdatedIterator, error) {

	logs, sub, err := _ContractEigenDAThresholdRegistry.contract.FilterLogs(opts, "DefaultSecurityThresholdsV2Updated")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2UpdatedIterator{contract: _ContractEigenDAThresholdRegistry.contract, event: "DefaultSecurityThresholdsV2Updated", logs: logs, sub: sub}, nil
}

// WatchDefaultSecurityThresholdsV2Updated is a free log subscription operation binding the contract event 0xfe03afd62c76a6aed7376ae995cc55d073ba9d83d83ac8efc5446f8da4d50997.
//
// Solidity: event DefaultSecurityThresholdsV2Updated((uint8,uint8) previousDefaultSecurityThresholdsV2, (uint8,uint8) newDefaultSecurityThresholdsV2)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) WatchDefaultSecurityThresholdsV2Updated(opts *bind.WatchOpts, sink chan<- *ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2Updated) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDAThresholdRegistry.contract.WatchLogs(opts, "DefaultSecurityThresholdsV2Updated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2Updated)
				if err := _ContractEigenDAThresholdRegistry.contract.UnpackLog(event, "DefaultSecurityThresholdsV2Updated", log); err != nil {
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

// ParseDefaultSecurityThresholdsV2Updated is a log parse operation binding the contract event 0xfe03afd62c76a6aed7376ae995cc55d073ba9d83d83ac8efc5446f8da4d50997.
//
// Solidity: event DefaultSecurityThresholdsV2Updated((uint8,uint8) previousDefaultSecurityThresholdsV2, (uint8,uint8) newDefaultSecurityThresholdsV2)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) ParseDefaultSecurityThresholdsV2Updated(log types.Log) (*ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2Updated, error) {
	event := new(ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2Updated)
	if err := _ContractEigenDAThresholdRegistry.contract.UnpackLog(event, "DefaultSecurityThresholdsV2Updated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAThresholdRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ContractEigenDAThresholdRegistry contract.
type ContractEigenDAThresholdRegistryInitializedIterator struct {
	Event *ContractEigenDAThresholdRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAThresholdRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAThresholdRegistryInitialized)
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
		it.Event = new(ContractEigenDAThresholdRegistryInitialized)
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
func (it *ContractEigenDAThresholdRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAThresholdRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAThresholdRegistryInitialized represents a Initialized event raised by the ContractEigenDAThresholdRegistry contract.
type ContractEigenDAThresholdRegistryInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*ContractEigenDAThresholdRegistryInitializedIterator, error) {

	logs, sub, err := _ContractEigenDAThresholdRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAThresholdRegistryInitializedIterator{contract: _ContractEigenDAThresholdRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ContractEigenDAThresholdRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDAThresholdRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAThresholdRegistryInitialized)
				if err := _ContractEigenDAThresholdRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) ParseInitialized(log types.Log) (*ContractEigenDAThresholdRegistryInitialized, error) {
	event := new(ContractEigenDAThresholdRegistryInitialized)
	if err := _ContractEigenDAThresholdRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAThresholdRegistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ContractEigenDAThresholdRegistry contract.
type ContractEigenDAThresholdRegistryOwnershipTransferredIterator struct {
	Event *ContractEigenDAThresholdRegistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAThresholdRegistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAThresholdRegistryOwnershipTransferred)
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
		it.Event = new(ContractEigenDAThresholdRegistryOwnershipTransferred)
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
func (it *ContractEigenDAThresholdRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAThresholdRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAThresholdRegistryOwnershipTransferred represents a OwnershipTransferred event raised by the ContractEigenDAThresholdRegistry contract.
type ContractEigenDAThresholdRegistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ContractEigenDAThresholdRegistryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractEigenDAThresholdRegistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAThresholdRegistryOwnershipTransferredIterator{contract: _ContractEigenDAThresholdRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ContractEigenDAThresholdRegistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractEigenDAThresholdRegistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAThresholdRegistryOwnershipTransferred)
				if err := _ContractEigenDAThresholdRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*ContractEigenDAThresholdRegistryOwnershipTransferred, error) {
	event := new(ContractEigenDAThresholdRegistryOwnershipTransferred)
	if err := _ContractEigenDAThresholdRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdatedIterator is returned from FilterQuorumAdversaryThresholdPercentagesUpdated and is used to iterate over the raw logs and unpacked data for QuorumAdversaryThresholdPercentagesUpdated events raised by the ContractEigenDAThresholdRegistry contract.
type ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdatedIterator struct {
	Event *ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdated // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdated)
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
		it.Event = new(ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdated)
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
func (it *ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdated represents a QuorumAdversaryThresholdPercentagesUpdated event raised by the ContractEigenDAThresholdRegistry contract.
type ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdated struct {
	PreviousQuorumAdversaryThresholdPercentages []byte
	NewQuorumAdversaryThresholdPercentages      []byte
	Raw                                         types.Log // Blockchain specific contextual infos
}

// FilterQuorumAdversaryThresholdPercentagesUpdated is a free log retrieval operation binding the contract event 0xf73542111561dc551cbbe9111c4dd3a040d53d7bc0339a53290f4d7f9a95c3cc.
//
// Solidity: event QuorumAdversaryThresholdPercentagesUpdated(bytes previousQuorumAdversaryThresholdPercentages, bytes newQuorumAdversaryThresholdPercentages)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) FilterQuorumAdversaryThresholdPercentagesUpdated(opts *bind.FilterOpts) (*ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdatedIterator, error) {

	logs, sub, err := _ContractEigenDAThresholdRegistry.contract.FilterLogs(opts, "QuorumAdversaryThresholdPercentagesUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdatedIterator{contract: _ContractEigenDAThresholdRegistry.contract, event: "QuorumAdversaryThresholdPercentagesUpdated", logs: logs, sub: sub}, nil
}

// WatchQuorumAdversaryThresholdPercentagesUpdated is a free log subscription operation binding the contract event 0xf73542111561dc551cbbe9111c4dd3a040d53d7bc0339a53290f4d7f9a95c3cc.
//
// Solidity: event QuorumAdversaryThresholdPercentagesUpdated(bytes previousQuorumAdversaryThresholdPercentages, bytes newQuorumAdversaryThresholdPercentages)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) WatchQuorumAdversaryThresholdPercentagesUpdated(opts *bind.WatchOpts, sink chan<- *ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDAThresholdRegistry.contract.WatchLogs(opts, "QuorumAdversaryThresholdPercentagesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdated)
				if err := _ContractEigenDAThresholdRegistry.contract.UnpackLog(event, "QuorumAdversaryThresholdPercentagesUpdated", log); err != nil {
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

// ParseQuorumAdversaryThresholdPercentagesUpdated is a log parse operation binding the contract event 0xf73542111561dc551cbbe9111c4dd3a040d53d7bc0339a53290f4d7f9a95c3cc.
//
// Solidity: event QuorumAdversaryThresholdPercentagesUpdated(bytes previousQuorumAdversaryThresholdPercentages, bytes newQuorumAdversaryThresholdPercentages)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) ParseQuorumAdversaryThresholdPercentagesUpdated(log types.Log) (*ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdated, error) {
	event := new(ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdated)
	if err := _ContractEigenDAThresholdRegistry.contract.UnpackLog(event, "QuorumAdversaryThresholdPercentagesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdatedIterator is returned from FilterQuorumConfirmationThresholdPercentagesUpdated and is used to iterate over the raw logs and unpacked data for QuorumConfirmationThresholdPercentagesUpdated events raised by the ContractEigenDAThresholdRegistry contract.
type ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdatedIterator struct {
	Event *ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdated // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdated)
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
		it.Event = new(ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdated)
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
func (it *ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdated represents a QuorumConfirmationThresholdPercentagesUpdated event raised by the ContractEigenDAThresholdRegistry contract.
type ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdated struct {
	PreviousQuorumConfirmationThresholdPercentages []byte
	NewQuorumConfirmationThresholdPercentages      []byte
	Raw                                            types.Log // Blockchain specific contextual infos
}

// FilterQuorumConfirmationThresholdPercentagesUpdated is a free log retrieval operation binding the contract event 0x9f1ea99a8363f2964c53c763811648354a8437441b30b39465f9d26118d6a5a0.
//
// Solidity: event QuorumConfirmationThresholdPercentagesUpdated(bytes previousQuorumConfirmationThresholdPercentages, bytes newQuorumConfirmationThresholdPercentages)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) FilterQuorumConfirmationThresholdPercentagesUpdated(opts *bind.FilterOpts) (*ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdatedIterator, error) {

	logs, sub, err := _ContractEigenDAThresholdRegistry.contract.FilterLogs(opts, "QuorumConfirmationThresholdPercentagesUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdatedIterator{contract: _ContractEigenDAThresholdRegistry.contract, event: "QuorumConfirmationThresholdPercentagesUpdated", logs: logs, sub: sub}, nil
}

// WatchQuorumConfirmationThresholdPercentagesUpdated is a free log subscription operation binding the contract event 0x9f1ea99a8363f2964c53c763811648354a8437441b30b39465f9d26118d6a5a0.
//
// Solidity: event QuorumConfirmationThresholdPercentagesUpdated(bytes previousQuorumConfirmationThresholdPercentages, bytes newQuorumConfirmationThresholdPercentages)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) WatchQuorumConfirmationThresholdPercentagesUpdated(opts *bind.WatchOpts, sink chan<- *ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDAThresholdRegistry.contract.WatchLogs(opts, "QuorumConfirmationThresholdPercentagesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdated)
				if err := _ContractEigenDAThresholdRegistry.contract.UnpackLog(event, "QuorumConfirmationThresholdPercentagesUpdated", log); err != nil {
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

// ParseQuorumConfirmationThresholdPercentagesUpdated is a log parse operation binding the contract event 0x9f1ea99a8363f2964c53c763811648354a8437441b30b39465f9d26118d6a5a0.
//
// Solidity: event QuorumConfirmationThresholdPercentagesUpdated(bytes previousQuorumConfirmationThresholdPercentages, bytes newQuorumConfirmationThresholdPercentages)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) ParseQuorumConfirmationThresholdPercentagesUpdated(log types.Log) (*ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdated, error) {
	event := new(ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdated)
	if err := _ContractEigenDAThresholdRegistry.contract.UnpackLog(event, "QuorumConfirmationThresholdPercentagesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdatedIterator is returned from FilterQuorumNumbersRequiredUpdated and is used to iterate over the raw logs and unpacked data for QuorumNumbersRequiredUpdated events raised by the ContractEigenDAThresholdRegistry contract.
type ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdatedIterator struct {
	Event *ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdated // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdated)
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
		it.Event = new(ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdated)
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
func (it *ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdated represents a QuorumNumbersRequiredUpdated event raised by the ContractEigenDAThresholdRegistry contract.
type ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdated struct {
	PreviousQuorumNumbersRequired []byte
	NewQuorumNumbersRequired      []byte
	Raw                           types.Log // Blockchain specific contextual infos
}

// FilterQuorumNumbersRequiredUpdated is a free log retrieval operation binding the contract event 0x60c0ba1da794fcbbf549d370512442cb8f3f3f774cb557205cc88c6f842cb36a.
//
// Solidity: event QuorumNumbersRequiredUpdated(bytes previousQuorumNumbersRequired, bytes newQuorumNumbersRequired)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) FilterQuorumNumbersRequiredUpdated(opts *bind.FilterOpts) (*ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdatedIterator, error) {

	logs, sub, err := _ContractEigenDAThresholdRegistry.contract.FilterLogs(opts, "QuorumNumbersRequiredUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdatedIterator{contract: _ContractEigenDAThresholdRegistry.contract, event: "QuorumNumbersRequiredUpdated", logs: logs, sub: sub}, nil
}

// WatchQuorumNumbersRequiredUpdated is a free log subscription operation binding the contract event 0x60c0ba1da794fcbbf549d370512442cb8f3f3f774cb557205cc88c6f842cb36a.
//
// Solidity: event QuorumNumbersRequiredUpdated(bytes previousQuorumNumbersRequired, bytes newQuorumNumbersRequired)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) WatchQuorumNumbersRequiredUpdated(opts *bind.WatchOpts, sink chan<- *ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDAThresholdRegistry.contract.WatchLogs(opts, "QuorumNumbersRequiredUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdated)
				if err := _ContractEigenDAThresholdRegistry.contract.UnpackLog(event, "QuorumNumbersRequiredUpdated", log); err != nil {
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

// ParseQuorumNumbersRequiredUpdated is a log parse operation binding the contract event 0x60c0ba1da794fcbbf549d370512442cb8f3f3f774cb557205cc88c6f842cb36a.
//
// Solidity: event QuorumNumbersRequiredUpdated(bytes previousQuorumNumbersRequired, bytes newQuorumNumbersRequired)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) ParseQuorumNumbersRequiredUpdated(log types.Log) (*ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdated, error) {
	event := new(ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdated)
	if err := _ContractEigenDAThresholdRegistry.contract.UnpackLog(event, "QuorumNumbersRequiredUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAThresholdRegistryVersionedBlobParamsAddedIterator is returned from FilterVersionedBlobParamsAdded and is used to iterate over the raw logs and unpacked data for VersionedBlobParamsAdded events raised by the ContractEigenDAThresholdRegistry contract.
type ContractEigenDAThresholdRegistryVersionedBlobParamsAddedIterator struct {
	Event *ContractEigenDAThresholdRegistryVersionedBlobParamsAdded // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAThresholdRegistryVersionedBlobParamsAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAThresholdRegistryVersionedBlobParamsAdded)
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
		it.Event = new(ContractEigenDAThresholdRegistryVersionedBlobParamsAdded)
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
func (it *ContractEigenDAThresholdRegistryVersionedBlobParamsAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAThresholdRegistryVersionedBlobParamsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAThresholdRegistryVersionedBlobParamsAdded represents a VersionedBlobParamsAdded event raised by the ContractEigenDAThresholdRegistry contract.
type ContractEigenDAThresholdRegistryVersionedBlobParamsAdded struct {
	Version             uint16
	VersionedBlobParams VersionedBlobParams
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterVersionedBlobParamsAdded is a free log retrieval operation binding the contract event 0xdbee9d337a6e5fde30966e157673aaeeb6a0134afaf774a4b6979b7c79d07da4.
//
// Solidity: event VersionedBlobParamsAdded(uint16 indexed version, (uint32,uint32,uint8) versionedBlobParams)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) FilterVersionedBlobParamsAdded(opts *bind.FilterOpts, version []uint16) (*ContractEigenDAThresholdRegistryVersionedBlobParamsAddedIterator, error) {

	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}

	logs, sub, err := _ContractEigenDAThresholdRegistry.contract.FilterLogs(opts, "VersionedBlobParamsAdded", versionRule)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAThresholdRegistryVersionedBlobParamsAddedIterator{contract: _ContractEigenDAThresholdRegistry.contract, event: "VersionedBlobParamsAdded", logs: logs, sub: sub}, nil
}

// WatchVersionedBlobParamsAdded is a free log subscription operation binding the contract event 0xdbee9d337a6e5fde30966e157673aaeeb6a0134afaf774a4b6979b7c79d07da4.
//
// Solidity: event VersionedBlobParamsAdded(uint16 indexed version, (uint32,uint32,uint8) versionedBlobParams)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) WatchVersionedBlobParamsAdded(opts *bind.WatchOpts, sink chan<- *ContractEigenDAThresholdRegistryVersionedBlobParamsAdded, version []uint16) (event.Subscription, error) {

	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}

	logs, sub, err := _ContractEigenDAThresholdRegistry.contract.WatchLogs(opts, "VersionedBlobParamsAdded", versionRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAThresholdRegistryVersionedBlobParamsAdded)
				if err := _ContractEigenDAThresholdRegistry.contract.UnpackLog(event, "VersionedBlobParamsAdded", log); err != nil {
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

// ParseVersionedBlobParamsAdded is a log parse operation binding the contract event 0xdbee9d337a6e5fde30966e157673aaeeb6a0134afaf774a4b6979b7c79d07da4.
//
// Solidity: event VersionedBlobParamsAdded(uint16 indexed version, (uint32,uint32,uint8) versionedBlobParams)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) ParseVersionedBlobParamsAdded(log types.Log) (*ContractEigenDAThresholdRegistryVersionedBlobParamsAdded, error) {
	event := new(ContractEigenDAThresholdRegistryVersionedBlobParamsAdded)
	if err := _ContractEigenDAThresholdRegistry.contract.UnpackLog(event, "VersionedBlobParamsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
