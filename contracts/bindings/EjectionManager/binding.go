// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractEjectionManager

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

// IEjectionManagerQuorumEjectionParams is an auto generated low-level Go binding around an user-defined struct.
type IEjectionManagerQuorumEjectionParams struct {
	RateLimitWindow       uint32
	EjectableStakePercent uint16
}

// ContractEjectionManagerMetaData contains all meta data concerning the ContractEjectionManager contract.
var ContractEjectionManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"},{\"name\":\"_stakeRegistry\",\"type\":\"address\",\"internalType\":\"contractIStakeRegistry\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"amountEjectableForQuorum\",\"inputs\":[{\"name\":\"_quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ejectOperators\",\"inputs\":[{\"name\":\"_operatorIds\",\"type\":\"bytes32[][]\",\"internalType\":\"bytes32[][]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_ejectors\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"_quorumEjectionParams\",\"type\":\"tuple[]\",\"internalType\":\"structIEjectionManager.QuorumEjectionParams[]\",\"components\":[{\"name\":\"rateLimitWindow\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"ejectableStakePercent\",\"type\":\"uint16\",\"internalType\":\"uint16\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isEjector\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumEjectionParams\",\"inputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"rateLimitWindow\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"ejectableStakePercent\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registryCoordinator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setEjector\",\"inputs\":[{\"name\":\"_ejector\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_status\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setQuorumEjectionParams\",\"inputs\":[{\"name\":\"_quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"_quorumEjectionParams\",\"type\":\"tuple\",\"internalType\":\"structIEjectionManager.QuorumEjectionParams\",\"components\":[{\"name\":\"rateLimitWindow\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"ejectableStakePercent\",\"type\":\"uint16\",\"internalType\":\"uint16\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stakeEjectedForQuorum\",\"inputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"stakeEjected\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"stakeRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIStakeRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"EjectorUpdated\",\"inputs\":[{\"name\":\"ejector\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"status\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorEjected\",\"inputs\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumEjection\",\"inputs\":[{\"name\":\"ejectedOperators\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"ratelimitHit\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumEjectionParamsSet\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"rateLimitWindow\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"ejectableStakePercent\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"}],\"anonymous\":false}]",
	Bin: "0x60c06040523480156200001157600080fd5b506040516200178d3803806200178d833981016040819052620000349162000134565b6001600160a01b03808316608052811660a0526200005162000059565b505062000173565b600054610100900460ff1615620000c65760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840160405180910390fd5b60005460ff908116101562000119576000805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b6001600160a01b03811681146200013157600080fd5b50565b600080604083850312156200014857600080fd5b825162000155816200011b565b602084015190925062000168816200011b565b809150509250929050565b60805160a0516115ca620001c360003960008181610180015281816103790152610be10152600081816101f2015281816104c1015281816104f0015281816106c301526106f201526115ca6000f3fe608060405234801561001057600080fd5b50600436106100ce5760003560e01c80636d14a9871161008c5780638b88a024116100665780638b88a0241461022f5780638da5cb5b14610242578063b13f450414610253578063f2fde38b1461027457600080fd5b80636d14a987146101ed578063715018a61461021457806377d175861461021c57600080fd5b8062482569146100d35780630a0593d11461012b57806310ea4f8a146101405780633a0b0ddd14610153578063683048351461017b5780636c08a879146101ba575b600080fd5b6101076100e1366004610fb8565b60676020526000908152604090205463ffffffff811690640100000000900461ffff1682565b6040805163ffffffff909316835261ffff9091166020830152015b60405180910390f35b61013e610139366004611045565b610287565b005b61013e61014e366004611166565b610980565b6101666101613660046111a4565b610992565b60408051928352602083019190915201610122565b6101a27f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b039091168152602001610122565b6101dd6101c83660046111ce565b60656020526000908152604090205460ff1681565b6040519015158152602001610122565b6101a27f000000000000000000000000000000000000000000000000000000000000000081565b61013e6109ce565b61013e61022a36600461125e565b6109e2565b61013e61023d366004611301565b6109f4565b6033546001600160a01b03166101a2565b610266610261366004610fb8565b610b9a565b604051908152602001610122565b61013e6102823660046111ce565b610d9c565b3360009081526065602052604090205460ff16806102af57506033546001600160a01b031633145b6103115760405162461bcd60e51b815260206004820152602860248201527f456a6563746f723a204f6e6c79206f776e6572206f7220656a6563746f722063604482015267185b88195a9958dd60c21b60648201526084015b60405180910390fd5b60005b815181101561097c5780600061032982610b9a565b905060008060008084118061034857506033546001600160a01b031633145b156108b95760005b878781518110610362576103626113d7565b6020026020010151518160ff1610156108b75760007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316635401ed278a8a815181106103b8576103b86113d7565b60200260200101518460ff16815181106103d4576103d46113d7565b6020026020010151896040518363ffffffff1660e01b815260040161040692919091825260ff16602082015260400190565b602060405180830381865afa158015610423573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061044791906113ed565b336000908152606560205260409020546001600160601b0391909116915060ff16801561048b575060ff871660009081526067602052604090205463ffffffff1615155b801561049f57508561049d828761142c565b115b156106aa57600192506104b2818661142c565b94506104bd84611444565b93507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316636e3b17db7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663296bb0648c8c8151811061052f5761052f6113d7565b60200260200101518660ff168151811061054b5761054b6113d7565b60200260200101516040518263ffffffff1660e01b815260040161057191815260200190565b602060405180830381865afa15801561058e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105b29190611468565b6040516001600160f81b031960f88c901b1660208201526021016040516020818303038152906040526040518363ffffffff1660e01b81526004016105f8929190611485565b600060405180830381600087803b15801561061257600080fd5b505af1158015610626573d6000803e3d6000fd5b505050507f97ddb711c61a9d2d7effcba3e042a33862297f898d555655cca39ec4451f53b489898151811061065d5761065d6113d7565b60200260200101518360ff1681518110610679576106796113d7565b60200260200101518860405161069c92919091825260ff16602082015260400190565b60405180910390a1506108b7565b6106b4818661142c565b94506106bf84611444565b93507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316636e3b17db7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663296bb0648c8c81518110610731576107316113d7565b60200260200101518660ff168151811061074d5761074d6113d7565b60200260200101516040518263ffffffff1660e01b815260040161077391815260200190565b602060405180830381865afa158015610790573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107b49190611468565b6040516001600160f81b031960f88c901b1660208201526021016040516020818303038152906040526040518363ffffffff1660e01b81526004016107fa929190611485565b600060405180830381600087803b15801561081457600080fd5b505af1158015610828573d6000803e3d6000fd5b505050507f97ddb711c61a9d2d7effcba3e042a33862297f898d555655cca39ec4451f53b489898151811061085f5761085f6113d7565b60200260200101518360ff168151811061087b5761087b6113d7565b60200260200101518860405161089e92919091825260ff16602082015260400190565b60405180910390a1506108b0816114ea565b9050610350565b505b3360009081526065602052604090205460ff1680156108d85750600083115b156109265760ff851660009081526066602090815260408083208151808301909252428252818301878152815460018181018455928652939094209151600290930290910191825591519101555b6040805163ffffffff8416815282151560208201527f19dd87ae49ed14a795f8c2d5e8055bf2a4a9d01641a00a2f8f0a5a7bf7f70249910160405180910390a15050505050806109759061150a565b9050610314565b5050565b610988610e15565b61097c8282610e6f565b606660205281600052604060002081815481106109ae57600080fd5b600091825260209091206002909102018054600190910154909250905082565b6109d6610e15565b6109e06000610ed3565b565b6109ea610e15565b61097c8282610f25565b600054610100900460ff1615808015610a145750600054600160ff909116105b80610a2e5750303b158015610a2e575060005460ff166001145b610a915760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b6064820152608401610308565b6000805460ff191660011790558015610ab4576000805461ff0019166101001790555b610abd84610ed3565b60005b83518160ff161015610b0557610af3848260ff1681518110610ae457610ae46113d7565b60200260200101516001610e6f565b80610afd816114ea565b915050610ac0565b5060005b82518160ff161015610b4d57610b3b81848360ff1681518110610b2e57610b2e6113d7565b6020026020010151610f25565b80610b45816114ea565b915050610b09565b508015610b94576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50505050565b60ff81166000908152606760205260408120548190610bbf9063ffffffff1642611525565b60405163d5eccc0560e01b815260ff85166004820152909150600090612710907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063d5eccc0590602401602060405180830381865afa158015610c30573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c5491906113ed565b60ff8616600090815260676020526040902054610c86916001600160601b031690640100000000900461ffff1661153c565b610c90919061155b565b60ff8516600090815260666020526040812054919250908190610cb7575090949350505050565b60ff8616600090815260666020526040902054610cd690600190611525565b90505b60ff86166000908152606660205260409020805485919083908110610d0057610d006113d7565b9060005260206000209060020201600001541115610d755760ff86166000908152606660205260409020805482908110610d3c57610d3c6113d7565b90600052602060002090600202016001015482610d59919061142c565b915080610d6557610d75565b610d6e8161157d565b9050610cd9565b828210610d885750600095945050505050565b610d928284611525565b9695505050505050565b610da4610e15565b6001600160a01b038116610e095760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610308565b610e1281610ed3565b50565b6033546001600160a01b031633146109e05760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610308565b6001600160a01b038216600081815260656020908152604091829020805460ff19168515159081179091558251938452908301527f7676686b6d22e112412bd874d70177e011ab06602c26063f19f0386c9a3cee4291015b60405180910390a15050565b603380546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b60ff8216600081815260676020908152604091829020845181548684015163ffffffff90921665ffffffffffff19909116811764010000000061ffff90931692830217909255835194855291840152908201527fe69c2827a1e2fdd32265ebb4eeea5ee564f0551cf5dfed4150f8e116a67209eb90606001610ec7565b803560ff81168114610fb357600080fd5b919050565b600060208284031215610fca57600080fd5b610fd382610fa2565b9392505050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff8111828210171561101957611019610fda565b604052919050565b600067ffffffffffffffff82111561103b5761103b610fda565b5060051b60200190565b6000602080838503121561105857600080fd5b823567ffffffffffffffff8082111561107057600080fd5b818501915085601f83011261108457600080fd5b813561109761109282611021565b610ff0565b818152600591821b84018501918582019190898411156110b657600080fd5b8686015b84811015611142578035868111156110d25760008081fd5b8701603f81018c136110e45760008081fd5b8881013560406110f661109283611021565b82815291851b83018101918b8101908f8411156111135760008081fd5b938201935b838510156111315784358252938c0193908c0190611118565b8852505050938801935087016110ba565b50909998505050505050505050565b6001600160a01b0381168114610e1257600080fd5b6000806040838503121561117957600080fd5b823561118481611151565b91506020830135801515811461119957600080fd5b809150509250929050565b600080604083850312156111b757600080fd5b6111c083610fa2565b946020939093013593505050565b6000602082840312156111e057600080fd5b8135610fd381611151565b6000604082840312156111fd57600080fd5b6040516040810181811067ffffffffffffffff8211171561122057611220610fda565b604052905080823563ffffffff8116811461123a57600080fd5b8152602083013561ffff8116811461125157600080fd5b6020919091015292915050565b6000806060838503121561127157600080fd5b61127a83610fa2565b915061128984602085016111eb565b90509250929050565b600082601f8301126112a357600080fd5b813560206112b361109283611021565b82815260069290921b840181019181810190868411156112d257600080fd5b8286015b848110156112f6576112e888826111eb565b8352918301916040016112d6565b509695505050505050565b60008060006060848603121561131657600080fd5b833561132181611151565b925060208481013567ffffffffffffffff8082111561133f57600080fd5b818701915087601f83011261135357600080fd5b813561136161109282611021565b81815260059190911b8301840190848101908a83111561138057600080fd5b938501935b828510156113a757843561139881611151565b82529385019390850190611385565b9650505060408701359250808311156113bf57600080fd5b50506113cd86828701611292565b9150509250925092565b634e487b7160e01b600052603260045260246000fd5b6000602082840312156113ff57600080fd5b81516001600160601b0381168114610fd357600080fd5b634e487b7160e01b600052601160045260246000fd5b6000821982111561143f5761143f611416565b500190565b600063ffffffff8083168181141561145e5761145e611416565b6001019392505050565b60006020828403121561147a57600080fd5b8151610fd381611151565b60018060a01b038316815260006020604081840152835180604085015260005b818110156114c1578581018301518582016060015282016114a5565b818111156114d3576000606083870101525b50601f01601f191692909201606001949350505050565b600060ff821660ff81141561150157611501611416565b60010192915050565b600060001982141561151e5761151e611416565b5060010190565b60008282101561153757611537611416565b500390565b600081600019048311821515161561155657611556611416565b500290565b60008261157857634e487b7160e01b600052601260045260246000fd5b500490565b60008161158c5761158c611416565b50600019019056fea2646970667358221220c47d5f2ce7854d29f53646568c1ec7b91981cfd571775c7986c6c41bd8402d8364736f6c634300080c0033",
}

// ContractEjectionManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEjectionManagerMetaData.ABI instead.
var ContractEjectionManagerABI = ContractEjectionManagerMetaData.ABI

// ContractEjectionManagerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractEjectionManagerMetaData.Bin instead.
var ContractEjectionManagerBin = ContractEjectionManagerMetaData.Bin

// DeployContractEjectionManager deploys a new Ethereum contract, binding an instance of ContractEjectionManager to it.
func DeployContractEjectionManager(auth *bind.TransactOpts, backend bind.ContractBackend, _registryCoordinator common.Address, _stakeRegistry common.Address) (common.Address, *types.Transaction, *ContractEjectionManager, error) {
	parsed, err := ContractEjectionManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractEjectionManagerBin), backend, _registryCoordinator, _stakeRegistry)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractEjectionManager{ContractEjectionManagerCaller: ContractEjectionManagerCaller{contract: contract}, ContractEjectionManagerTransactor: ContractEjectionManagerTransactor{contract: contract}, ContractEjectionManagerFilterer: ContractEjectionManagerFilterer{contract: contract}}, nil
}

// ContractEjectionManager is an auto generated Go binding around an Ethereum contract.
type ContractEjectionManager struct {
	ContractEjectionManagerCaller     // Read-only binding to the contract
	ContractEjectionManagerTransactor // Write-only binding to the contract
	ContractEjectionManagerFilterer   // Log filterer for contract events
}

// ContractEjectionManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractEjectionManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEjectionManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractEjectionManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEjectionManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractEjectionManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEjectionManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractEjectionManagerSession struct {
	Contract     *ContractEjectionManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// ContractEjectionManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractEjectionManagerCallerSession struct {
	Contract *ContractEjectionManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// ContractEjectionManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractEjectionManagerTransactorSession struct {
	Contract     *ContractEjectionManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// ContractEjectionManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractEjectionManagerRaw struct {
	Contract *ContractEjectionManager // Generic contract binding to access the raw methods on
}

// ContractEjectionManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractEjectionManagerCallerRaw struct {
	Contract *ContractEjectionManagerCaller // Generic read-only contract binding to access the raw methods on
}

// ContractEjectionManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractEjectionManagerTransactorRaw struct {
	Contract *ContractEjectionManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractEjectionManager creates a new instance of ContractEjectionManager, bound to a specific deployed contract.
func NewContractEjectionManager(address common.Address, backend bind.ContractBackend) (*ContractEjectionManager, error) {
	contract, err := bindContractEjectionManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractEjectionManager{ContractEjectionManagerCaller: ContractEjectionManagerCaller{contract: contract}, ContractEjectionManagerTransactor: ContractEjectionManagerTransactor{contract: contract}, ContractEjectionManagerFilterer: ContractEjectionManagerFilterer{contract: contract}}, nil
}

// NewContractEjectionManagerCaller creates a new read-only instance of ContractEjectionManager, bound to a specific deployed contract.
func NewContractEjectionManagerCaller(address common.Address, caller bind.ContractCaller) (*ContractEjectionManagerCaller, error) {
	contract, err := bindContractEjectionManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEjectionManagerCaller{contract: contract}, nil
}

// NewContractEjectionManagerTransactor creates a new write-only instance of ContractEjectionManager, bound to a specific deployed contract.
func NewContractEjectionManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractEjectionManagerTransactor, error) {
	contract, err := bindContractEjectionManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEjectionManagerTransactor{contract: contract}, nil
}

// NewContractEjectionManagerFilterer creates a new log filterer instance of ContractEjectionManager, bound to a specific deployed contract.
func NewContractEjectionManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractEjectionManagerFilterer, error) {
	contract, err := bindContractEjectionManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractEjectionManagerFilterer{contract: contract}, nil
}

// bindContractEjectionManager binds a generic wrapper to an already deployed contract.
func bindContractEjectionManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractEjectionManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEjectionManager *ContractEjectionManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEjectionManager.Contract.ContractEjectionManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEjectionManager *ContractEjectionManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.ContractEjectionManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEjectionManager *ContractEjectionManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.ContractEjectionManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEjectionManager *ContractEjectionManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEjectionManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEjectionManager *ContractEjectionManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEjectionManager *ContractEjectionManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.contract.Transact(opts, method, params...)
}

// AmountEjectableForQuorum is a free data retrieval call binding the contract method 0xb13f4504.
//
// Solidity: function amountEjectableForQuorum(uint8 _quorumNumber) view returns(uint256)
func (_ContractEjectionManager *ContractEjectionManagerCaller) AmountEjectableForQuorum(opts *bind.CallOpts, _quorumNumber uint8) (*big.Int, error) {
	var out []interface{}
	err := _ContractEjectionManager.contract.Call(opts, &out, "amountEjectableForQuorum", _quorumNumber)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AmountEjectableForQuorum is a free data retrieval call binding the contract method 0xb13f4504.
//
// Solidity: function amountEjectableForQuorum(uint8 _quorumNumber) view returns(uint256)
func (_ContractEjectionManager *ContractEjectionManagerSession) AmountEjectableForQuorum(_quorumNumber uint8) (*big.Int, error) {
	return _ContractEjectionManager.Contract.AmountEjectableForQuorum(&_ContractEjectionManager.CallOpts, _quorumNumber)
}

// AmountEjectableForQuorum is a free data retrieval call binding the contract method 0xb13f4504.
//
// Solidity: function amountEjectableForQuorum(uint8 _quorumNumber) view returns(uint256)
func (_ContractEjectionManager *ContractEjectionManagerCallerSession) AmountEjectableForQuorum(_quorumNumber uint8) (*big.Int, error) {
	return _ContractEjectionManager.Contract.AmountEjectableForQuorum(&_ContractEjectionManager.CallOpts, _quorumNumber)
}

// IsEjector is a free data retrieval call binding the contract method 0x6c08a879.
//
// Solidity: function isEjector(address ) view returns(bool)
func (_ContractEjectionManager *ContractEjectionManagerCaller) IsEjector(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _ContractEjectionManager.contract.Call(opts, &out, "isEjector", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsEjector is a free data retrieval call binding the contract method 0x6c08a879.
//
// Solidity: function isEjector(address ) view returns(bool)
func (_ContractEjectionManager *ContractEjectionManagerSession) IsEjector(arg0 common.Address) (bool, error) {
	return _ContractEjectionManager.Contract.IsEjector(&_ContractEjectionManager.CallOpts, arg0)
}

// IsEjector is a free data retrieval call binding the contract method 0x6c08a879.
//
// Solidity: function isEjector(address ) view returns(bool)
func (_ContractEjectionManager *ContractEjectionManagerCallerSession) IsEjector(arg0 common.Address) (bool, error) {
	return _ContractEjectionManager.Contract.IsEjector(&_ContractEjectionManager.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEjectionManager *ContractEjectionManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEjectionManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEjectionManager *ContractEjectionManagerSession) Owner() (common.Address, error) {
	return _ContractEjectionManager.Contract.Owner(&_ContractEjectionManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEjectionManager *ContractEjectionManagerCallerSession) Owner() (common.Address, error) {
	return _ContractEjectionManager.Contract.Owner(&_ContractEjectionManager.CallOpts)
}

// QuorumEjectionParams is a free data retrieval call binding the contract method 0x00482569.
//
// Solidity: function quorumEjectionParams(uint8 ) view returns(uint32 rateLimitWindow, uint16 ejectableStakePercent)
func (_ContractEjectionManager *ContractEjectionManagerCaller) QuorumEjectionParams(opts *bind.CallOpts, arg0 uint8) (struct {
	RateLimitWindow       uint32
	EjectableStakePercent uint16
}, error) {
	var out []interface{}
	err := _ContractEjectionManager.contract.Call(opts, &out, "quorumEjectionParams", arg0)

	outstruct := new(struct {
		RateLimitWindow       uint32
		EjectableStakePercent uint16
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RateLimitWindow = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.EjectableStakePercent = *abi.ConvertType(out[1], new(uint16)).(*uint16)

	return *outstruct, err

}

// QuorumEjectionParams is a free data retrieval call binding the contract method 0x00482569.
//
// Solidity: function quorumEjectionParams(uint8 ) view returns(uint32 rateLimitWindow, uint16 ejectableStakePercent)
func (_ContractEjectionManager *ContractEjectionManagerSession) QuorumEjectionParams(arg0 uint8) (struct {
	RateLimitWindow       uint32
	EjectableStakePercent uint16
}, error) {
	return _ContractEjectionManager.Contract.QuorumEjectionParams(&_ContractEjectionManager.CallOpts, arg0)
}

// QuorumEjectionParams is a free data retrieval call binding the contract method 0x00482569.
//
// Solidity: function quorumEjectionParams(uint8 ) view returns(uint32 rateLimitWindow, uint16 ejectableStakePercent)
func (_ContractEjectionManager *ContractEjectionManagerCallerSession) QuorumEjectionParams(arg0 uint8) (struct {
	RateLimitWindow       uint32
	EjectableStakePercent uint16
}, error) {
	return _ContractEjectionManager.Contract.QuorumEjectionParams(&_ContractEjectionManager.CallOpts, arg0)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractEjectionManager *ContractEjectionManagerCaller) RegistryCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEjectionManager.contract.Call(opts, &out, "registryCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractEjectionManager *ContractEjectionManagerSession) RegistryCoordinator() (common.Address, error) {
	return _ContractEjectionManager.Contract.RegistryCoordinator(&_ContractEjectionManager.CallOpts)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractEjectionManager *ContractEjectionManagerCallerSession) RegistryCoordinator() (common.Address, error) {
	return _ContractEjectionManager.Contract.RegistryCoordinator(&_ContractEjectionManager.CallOpts)
}

// StakeEjectedForQuorum is a free data retrieval call binding the contract method 0x3a0b0ddd.
//
// Solidity: function stakeEjectedForQuorum(uint8 , uint256 ) view returns(uint256 timestamp, uint256 stakeEjected)
func (_ContractEjectionManager *ContractEjectionManagerCaller) StakeEjectedForQuorum(opts *bind.CallOpts, arg0 uint8, arg1 *big.Int) (struct {
	Timestamp    *big.Int
	StakeEjected *big.Int
}, error) {
	var out []interface{}
	err := _ContractEjectionManager.contract.Call(opts, &out, "stakeEjectedForQuorum", arg0, arg1)

	outstruct := new(struct {
		Timestamp    *big.Int
		StakeEjected *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Timestamp = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.StakeEjected = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// StakeEjectedForQuorum is a free data retrieval call binding the contract method 0x3a0b0ddd.
//
// Solidity: function stakeEjectedForQuorum(uint8 , uint256 ) view returns(uint256 timestamp, uint256 stakeEjected)
func (_ContractEjectionManager *ContractEjectionManagerSession) StakeEjectedForQuorum(arg0 uint8, arg1 *big.Int) (struct {
	Timestamp    *big.Int
	StakeEjected *big.Int
}, error) {
	return _ContractEjectionManager.Contract.StakeEjectedForQuorum(&_ContractEjectionManager.CallOpts, arg0, arg1)
}

// StakeEjectedForQuorum is a free data retrieval call binding the contract method 0x3a0b0ddd.
//
// Solidity: function stakeEjectedForQuorum(uint8 , uint256 ) view returns(uint256 timestamp, uint256 stakeEjected)
func (_ContractEjectionManager *ContractEjectionManagerCallerSession) StakeEjectedForQuorum(arg0 uint8, arg1 *big.Int) (struct {
	Timestamp    *big.Int
	StakeEjected *big.Int
}, error) {
	return _ContractEjectionManager.Contract.StakeEjectedForQuorum(&_ContractEjectionManager.CallOpts, arg0, arg1)
}

// StakeRegistry is a free data retrieval call binding the contract method 0x68304835.
//
// Solidity: function stakeRegistry() view returns(address)
func (_ContractEjectionManager *ContractEjectionManagerCaller) StakeRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEjectionManager.contract.Call(opts, &out, "stakeRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StakeRegistry is a free data retrieval call binding the contract method 0x68304835.
//
// Solidity: function stakeRegistry() view returns(address)
func (_ContractEjectionManager *ContractEjectionManagerSession) StakeRegistry() (common.Address, error) {
	return _ContractEjectionManager.Contract.StakeRegistry(&_ContractEjectionManager.CallOpts)
}

// StakeRegistry is a free data retrieval call binding the contract method 0x68304835.
//
// Solidity: function stakeRegistry() view returns(address)
func (_ContractEjectionManager *ContractEjectionManagerCallerSession) StakeRegistry() (common.Address, error) {
	return _ContractEjectionManager.Contract.StakeRegistry(&_ContractEjectionManager.CallOpts)
}

// EjectOperators is a paid mutator transaction binding the contract method 0x0a0593d1.
//
// Solidity: function ejectOperators(bytes32[][] _operatorIds) returns()
func (_ContractEjectionManager *ContractEjectionManagerTransactor) EjectOperators(opts *bind.TransactOpts, _operatorIds [][][32]byte) (*types.Transaction, error) {
	return _ContractEjectionManager.contract.Transact(opts, "ejectOperators", _operatorIds)
}

// EjectOperators is a paid mutator transaction binding the contract method 0x0a0593d1.
//
// Solidity: function ejectOperators(bytes32[][] _operatorIds) returns()
func (_ContractEjectionManager *ContractEjectionManagerSession) EjectOperators(_operatorIds [][][32]byte) (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.EjectOperators(&_ContractEjectionManager.TransactOpts, _operatorIds)
}

// EjectOperators is a paid mutator transaction binding the contract method 0x0a0593d1.
//
// Solidity: function ejectOperators(bytes32[][] _operatorIds) returns()
func (_ContractEjectionManager *ContractEjectionManagerTransactorSession) EjectOperators(_operatorIds [][][32]byte) (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.EjectOperators(&_ContractEjectionManager.TransactOpts, _operatorIds)
}

// Initialize is a paid mutator transaction binding the contract method 0x8b88a024.
//
// Solidity: function initialize(address _owner, address[] _ejectors, (uint32,uint16)[] _quorumEjectionParams) returns()
func (_ContractEjectionManager *ContractEjectionManagerTransactor) Initialize(opts *bind.TransactOpts, _owner common.Address, _ejectors []common.Address, _quorumEjectionParams []IEjectionManagerQuorumEjectionParams) (*types.Transaction, error) {
	return _ContractEjectionManager.contract.Transact(opts, "initialize", _owner, _ejectors, _quorumEjectionParams)
}

// Initialize is a paid mutator transaction binding the contract method 0x8b88a024.
//
// Solidity: function initialize(address _owner, address[] _ejectors, (uint32,uint16)[] _quorumEjectionParams) returns()
func (_ContractEjectionManager *ContractEjectionManagerSession) Initialize(_owner common.Address, _ejectors []common.Address, _quorumEjectionParams []IEjectionManagerQuorumEjectionParams) (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.Initialize(&_ContractEjectionManager.TransactOpts, _owner, _ejectors, _quorumEjectionParams)
}

// Initialize is a paid mutator transaction binding the contract method 0x8b88a024.
//
// Solidity: function initialize(address _owner, address[] _ejectors, (uint32,uint16)[] _quorumEjectionParams) returns()
func (_ContractEjectionManager *ContractEjectionManagerTransactorSession) Initialize(_owner common.Address, _ejectors []common.Address, _quorumEjectionParams []IEjectionManagerQuorumEjectionParams) (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.Initialize(&_ContractEjectionManager.TransactOpts, _owner, _ejectors, _quorumEjectionParams)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEjectionManager *ContractEjectionManagerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEjectionManager.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEjectionManager *ContractEjectionManagerSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.RenounceOwnership(&_ContractEjectionManager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEjectionManager *ContractEjectionManagerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.RenounceOwnership(&_ContractEjectionManager.TransactOpts)
}

// SetEjector is a paid mutator transaction binding the contract method 0x10ea4f8a.
//
// Solidity: function setEjector(address _ejector, bool _status) returns()
func (_ContractEjectionManager *ContractEjectionManagerTransactor) SetEjector(opts *bind.TransactOpts, _ejector common.Address, _status bool) (*types.Transaction, error) {
	return _ContractEjectionManager.contract.Transact(opts, "setEjector", _ejector, _status)
}

// SetEjector is a paid mutator transaction binding the contract method 0x10ea4f8a.
//
// Solidity: function setEjector(address _ejector, bool _status) returns()
func (_ContractEjectionManager *ContractEjectionManagerSession) SetEjector(_ejector common.Address, _status bool) (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.SetEjector(&_ContractEjectionManager.TransactOpts, _ejector, _status)
}

// SetEjector is a paid mutator transaction binding the contract method 0x10ea4f8a.
//
// Solidity: function setEjector(address _ejector, bool _status) returns()
func (_ContractEjectionManager *ContractEjectionManagerTransactorSession) SetEjector(_ejector common.Address, _status bool) (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.SetEjector(&_ContractEjectionManager.TransactOpts, _ejector, _status)
}

// SetQuorumEjectionParams is a paid mutator transaction binding the contract method 0x77d17586.
//
// Solidity: function setQuorumEjectionParams(uint8 _quorumNumber, (uint32,uint16) _quorumEjectionParams) returns()
func (_ContractEjectionManager *ContractEjectionManagerTransactor) SetQuorumEjectionParams(opts *bind.TransactOpts, _quorumNumber uint8, _quorumEjectionParams IEjectionManagerQuorumEjectionParams) (*types.Transaction, error) {
	return _ContractEjectionManager.contract.Transact(opts, "setQuorumEjectionParams", _quorumNumber, _quorumEjectionParams)
}

// SetQuorumEjectionParams is a paid mutator transaction binding the contract method 0x77d17586.
//
// Solidity: function setQuorumEjectionParams(uint8 _quorumNumber, (uint32,uint16) _quorumEjectionParams) returns()
func (_ContractEjectionManager *ContractEjectionManagerSession) SetQuorumEjectionParams(_quorumNumber uint8, _quorumEjectionParams IEjectionManagerQuorumEjectionParams) (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.SetQuorumEjectionParams(&_ContractEjectionManager.TransactOpts, _quorumNumber, _quorumEjectionParams)
}

// SetQuorumEjectionParams is a paid mutator transaction binding the contract method 0x77d17586.
//
// Solidity: function setQuorumEjectionParams(uint8 _quorumNumber, (uint32,uint16) _quorumEjectionParams) returns()
func (_ContractEjectionManager *ContractEjectionManagerTransactorSession) SetQuorumEjectionParams(_quorumNumber uint8, _quorumEjectionParams IEjectionManagerQuorumEjectionParams) (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.SetQuorumEjectionParams(&_ContractEjectionManager.TransactOpts, _quorumNumber, _quorumEjectionParams)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEjectionManager *ContractEjectionManagerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ContractEjectionManager.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEjectionManager *ContractEjectionManagerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.TransferOwnership(&_ContractEjectionManager.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEjectionManager *ContractEjectionManagerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.TransferOwnership(&_ContractEjectionManager.TransactOpts, newOwner)
}

// ContractEjectionManagerEjectorUpdatedIterator is returned from FilterEjectorUpdated and is used to iterate over the raw logs and unpacked data for EjectorUpdated events raised by the ContractEjectionManager contract.
type ContractEjectionManagerEjectorUpdatedIterator struct {
	Event *ContractEjectionManagerEjectorUpdated // Event containing the contract specifics and raw log

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
func (it *ContractEjectionManagerEjectorUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEjectionManagerEjectorUpdated)
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
		it.Event = new(ContractEjectionManagerEjectorUpdated)
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
func (it *ContractEjectionManagerEjectorUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEjectionManagerEjectorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEjectionManagerEjectorUpdated represents a EjectorUpdated event raised by the ContractEjectionManager contract.
type ContractEjectionManagerEjectorUpdated struct {
	Ejector common.Address
	Status  bool
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterEjectorUpdated is a free log retrieval operation binding the contract event 0x7676686b6d22e112412bd874d70177e011ab06602c26063f19f0386c9a3cee42.
//
// Solidity: event EjectorUpdated(address ejector, bool status)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) FilterEjectorUpdated(opts *bind.FilterOpts) (*ContractEjectionManagerEjectorUpdatedIterator, error) {

	logs, sub, err := _ContractEjectionManager.contract.FilterLogs(opts, "EjectorUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractEjectionManagerEjectorUpdatedIterator{contract: _ContractEjectionManager.contract, event: "EjectorUpdated", logs: logs, sub: sub}, nil
}

// WatchEjectorUpdated is a free log subscription operation binding the contract event 0x7676686b6d22e112412bd874d70177e011ab06602c26063f19f0386c9a3cee42.
//
// Solidity: event EjectorUpdated(address ejector, bool status)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) WatchEjectorUpdated(opts *bind.WatchOpts, sink chan<- *ContractEjectionManagerEjectorUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractEjectionManager.contract.WatchLogs(opts, "EjectorUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEjectionManagerEjectorUpdated)
				if err := _ContractEjectionManager.contract.UnpackLog(event, "EjectorUpdated", log); err != nil {
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

// ParseEjectorUpdated is a log parse operation binding the contract event 0x7676686b6d22e112412bd874d70177e011ab06602c26063f19f0386c9a3cee42.
//
// Solidity: event EjectorUpdated(address ejector, bool status)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) ParseEjectorUpdated(log types.Log) (*ContractEjectionManagerEjectorUpdated, error) {
	event := new(ContractEjectionManagerEjectorUpdated)
	if err := _ContractEjectionManager.contract.UnpackLog(event, "EjectorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEjectionManagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ContractEjectionManager contract.
type ContractEjectionManagerInitializedIterator struct {
	Event *ContractEjectionManagerInitialized // Event containing the contract specifics and raw log

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
func (it *ContractEjectionManagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEjectionManagerInitialized)
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
		it.Event = new(ContractEjectionManagerInitialized)
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
func (it *ContractEjectionManagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEjectionManagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEjectionManagerInitialized represents a Initialized event raised by the ContractEjectionManager contract.
type ContractEjectionManagerInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*ContractEjectionManagerInitializedIterator, error) {

	logs, sub, err := _ContractEjectionManager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ContractEjectionManagerInitializedIterator{contract: _ContractEjectionManager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ContractEjectionManagerInitialized) (event.Subscription, error) {

	logs, sub, err := _ContractEjectionManager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEjectionManagerInitialized)
				if err := _ContractEjectionManager.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ContractEjectionManager *ContractEjectionManagerFilterer) ParseInitialized(log types.Log) (*ContractEjectionManagerInitialized, error) {
	event := new(ContractEjectionManagerInitialized)
	if err := _ContractEjectionManager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEjectionManagerOperatorEjectedIterator is returned from FilterOperatorEjected and is used to iterate over the raw logs and unpacked data for OperatorEjected events raised by the ContractEjectionManager contract.
type ContractEjectionManagerOperatorEjectedIterator struct {
	Event *ContractEjectionManagerOperatorEjected // Event containing the contract specifics and raw log

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
func (it *ContractEjectionManagerOperatorEjectedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEjectionManagerOperatorEjected)
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
		it.Event = new(ContractEjectionManagerOperatorEjected)
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
func (it *ContractEjectionManagerOperatorEjectedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEjectionManagerOperatorEjectedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEjectionManagerOperatorEjected represents a OperatorEjected event raised by the ContractEjectionManager contract.
type ContractEjectionManagerOperatorEjected struct {
	OperatorId   [32]byte
	QuorumNumber uint8
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterOperatorEjected is a free log retrieval operation binding the contract event 0x97ddb711c61a9d2d7effcba3e042a33862297f898d555655cca39ec4451f53b4.
//
// Solidity: event OperatorEjected(bytes32 operatorId, uint8 quorumNumber)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) FilterOperatorEjected(opts *bind.FilterOpts) (*ContractEjectionManagerOperatorEjectedIterator, error) {

	logs, sub, err := _ContractEjectionManager.contract.FilterLogs(opts, "OperatorEjected")
	if err != nil {
		return nil, err
	}
	return &ContractEjectionManagerOperatorEjectedIterator{contract: _ContractEjectionManager.contract, event: "OperatorEjected", logs: logs, sub: sub}, nil
}

// WatchOperatorEjected is a free log subscription operation binding the contract event 0x97ddb711c61a9d2d7effcba3e042a33862297f898d555655cca39ec4451f53b4.
//
// Solidity: event OperatorEjected(bytes32 operatorId, uint8 quorumNumber)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) WatchOperatorEjected(opts *bind.WatchOpts, sink chan<- *ContractEjectionManagerOperatorEjected) (event.Subscription, error) {

	logs, sub, err := _ContractEjectionManager.contract.WatchLogs(opts, "OperatorEjected")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEjectionManagerOperatorEjected)
				if err := _ContractEjectionManager.contract.UnpackLog(event, "OperatorEjected", log); err != nil {
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

// ParseOperatorEjected is a log parse operation binding the contract event 0x97ddb711c61a9d2d7effcba3e042a33862297f898d555655cca39ec4451f53b4.
//
// Solidity: event OperatorEjected(bytes32 operatorId, uint8 quorumNumber)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) ParseOperatorEjected(log types.Log) (*ContractEjectionManagerOperatorEjected, error) {
	event := new(ContractEjectionManagerOperatorEjected)
	if err := _ContractEjectionManager.contract.UnpackLog(event, "OperatorEjected", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEjectionManagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ContractEjectionManager contract.
type ContractEjectionManagerOwnershipTransferredIterator struct {
	Event *ContractEjectionManagerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ContractEjectionManagerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEjectionManagerOwnershipTransferred)
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
		it.Event = new(ContractEjectionManagerOwnershipTransferred)
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
func (it *ContractEjectionManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEjectionManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEjectionManagerOwnershipTransferred represents a OwnershipTransferred event raised by the ContractEjectionManager contract.
type ContractEjectionManagerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ContractEjectionManagerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractEjectionManager.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ContractEjectionManagerOwnershipTransferredIterator{contract: _ContractEjectionManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ContractEjectionManagerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractEjectionManager.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEjectionManagerOwnershipTransferred)
				if err := _ContractEjectionManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ContractEjectionManager *ContractEjectionManagerFilterer) ParseOwnershipTransferred(log types.Log) (*ContractEjectionManagerOwnershipTransferred, error) {
	event := new(ContractEjectionManagerOwnershipTransferred)
	if err := _ContractEjectionManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEjectionManagerQuorumEjectionIterator is returned from FilterQuorumEjection and is used to iterate over the raw logs and unpacked data for QuorumEjection events raised by the ContractEjectionManager contract.
type ContractEjectionManagerQuorumEjectionIterator struct {
	Event *ContractEjectionManagerQuorumEjection // Event containing the contract specifics and raw log

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
func (it *ContractEjectionManagerQuorumEjectionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEjectionManagerQuorumEjection)
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
		it.Event = new(ContractEjectionManagerQuorumEjection)
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
func (it *ContractEjectionManagerQuorumEjectionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEjectionManagerQuorumEjectionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEjectionManagerQuorumEjection represents a QuorumEjection event raised by the ContractEjectionManager contract.
type ContractEjectionManagerQuorumEjection struct {
	EjectedOperators uint32
	RatelimitHit     bool
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterQuorumEjection is a free log retrieval operation binding the contract event 0x19dd87ae49ed14a795f8c2d5e8055bf2a4a9d01641a00a2f8f0a5a7bf7f70249.
//
// Solidity: event QuorumEjection(uint32 ejectedOperators, bool ratelimitHit)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) FilterQuorumEjection(opts *bind.FilterOpts) (*ContractEjectionManagerQuorumEjectionIterator, error) {

	logs, sub, err := _ContractEjectionManager.contract.FilterLogs(opts, "QuorumEjection")
	if err != nil {
		return nil, err
	}
	return &ContractEjectionManagerQuorumEjectionIterator{contract: _ContractEjectionManager.contract, event: "QuorumEjection", logs: logs, sub: sub}, nil
}

// WatchQuorumEjection is a free log subscription operation binding the contract event 0x19dd87ae49ed14a795f8c2d5e8055bf2a4a9d01641a00a2f8f0a5a7bf7f70249.
//
// Solidity: event QuorumEjection(uint32 ejectedOperators, bool ratelimitHit)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) WatchQuorumEjection(opts *bind.WatchOpts, sink chan<- *ContractEjectionManagerQuorumEjection) (event.Subscription, error) {

	logs, sub, err := _ContractEjectionManager.contract.WatchLogs(opts, "QuorumEjection")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEjectionManagerQuorumEjection)
				if err := _ContractEjectionManager.contract.UnpackLog(event, "QuorumEjection", log); err != nil {
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

// ParseQuorumEjection is a log parse operation binding the contract event 0x19dd87ae49ed14a795f8c2d5e8055bf2a4a9d01641a00a2f8f0a5a7bf7f70249.
//
// Solidity: event QuorumEjection(uint32 ejectedOperators, bool ratelimitHit)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) ParseQuorumEjection(log types.Log) (*ContractEjectionManagerQuorumEjection, error) {
	event := new(ContractEjectionManagerQuorumEjection)
	if err := _ContractEjectionManager.contract.UnpackLog(event, "QuorumEjection", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEjectionManagerQuorumEjectionParamsSetIterator is returned from FilterQuorumEjectionParamsSet and is used to iterate over the raw logs and unpacked data for QuorumEjectionParamsSet events raised by the ContractEjectionManager contract.
type ContractEjectionManagerQuorumEjectionParamsSetIterator struct {
	Event *ContractEjectionManagerQuorumEjectionParamsSet // Event containing the contract specifics and raw log

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
func (it *ContractEjectionManagerQuorumEjectionParamsSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEjectionManagerQuorumEjectionParamsSet)
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
		it.Event = new(ContractEjectionManagerQuorumEjectionParamsSet)
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
func (it *ContractEjectionManagerQuorumEjectionParamsSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEjectionManagerQuorumEjectionParamsSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEjectionManagerQuorumEjectionParamsSet represents a QuorumEjectionParamsSet event raised by the ContractEjectionManager contract.
type ContractEjectionManagerQuorumEjectionParamsSet struct {
	QuorumNumber          uint8
	RateLimitWindow       uint32
	EjectableStakePercent uint16
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterQuorumEjectionParamsSet is a free log retrieval operation binding the contract event 0xe69c2827a1e2fdd32265ebb4eeea5ee564f0551cf5dfed4150f8e116a67209eb.
//
// Solidity: event QuorumEjectionParamsSet(uint8 quorumNumber, uint32 rateLimitWindow, uint16 ejectableStakePercent)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) FilterQuorumEjectionParamsSet(opts *bind.FilterOpts) (*ContractEjectionManagerQuorumEjectionParamsSetIterator, error) {

	logs, sub, err := _ContractEjectionManager.contract.FilterLogs(opts, "QuorumEjectionParamsSet")
	if err != nil {
		return nil, err
	}
	return &ContractEjectionManagerQuorumEjectionParamsSetIterator{contract: _ContractEjectionManager.contract, event: "QuorumEjectionParamsSet", logs: logs, sub: sub}, nil
}

// WatchQuorumEjectionParamsSet is a free log subscription operation binding the contract event 0xe69c2827a1e2fdd32265ebb4eeea5ee564f0551cf5dfed4150f8e116a67209eb.
//
// Solidity: event QuorumEjectionParamsSet(uint8 quorumNumber, uint32 rateLimitWindow, uint16 ejectableStakePercent)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) WatchQuorumEjectionParamsSet(opts *bind.WatchOpts, sink chan<- *ContractEjectionManagerQuorumEjectionParamsSet) (event.Subscription, error) {

	logs, sub, err := _ContractEjectionManager.contract.WatchLogs(opts, "QuorumEjectionParamsSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEjectionManagerQuorumEjectionParamsSet)
				if err := _ContractEjectionManager.contract.UnpackLog(event, "QuorumEjectionParamsSet", log); err != nil {
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

// ParseQuorumEjectionParamsSet is a log parse operation binding the contract event 0xe69c2827a1e2fdd32265ebb4eeea5ee564f0551cf5dfed4150f8e116a67209eb.
//
// Solidity: event QuorumEjectionParamsSet(uint8 quorumNumber, uint32 rateLimitWindow, uint16 ejectableStakePercent)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) ParseQuorumEjectionParamsSet(log types.Log) (*ContractEjectionManagerQuorumEjectionParamsSet, error) {
	event := new(ContractEjectionManagerQuorumEjectionParamsSet)
	if err := _ContractEjectionManager.contract.UnpackLog(event, "QuorumEjectionParamsSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
