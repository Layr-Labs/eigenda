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
	Bin: "0x60c06040523480156200001157600080fd5b50604051620015a8380380620015a8833981016040819052620000349162000134565b6001600160a01b03808316608052811660a0526200005162000059565b505062000173565b600054610100900460ff1615620000c65760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840160405180910390fd5b60005460ff908116101562000119576000805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b6001600160a01b03811681146200013157600080fd5b50565b600080604083850312156200014857600080fd5b825162000155816200011b565b602084015190925062000168816200011b565b809150509250929050565b60805160a0516113f3620001b5600039600081816101800152818161035a0152610a0a0152600081816101f2015281816104ef015261051e01526113f36000f3fe608060405234801561001057600080fd5b50600436106100ce5760003560e01c80636d14a9871161008c5780638b88a024116100665780638b88a0241461022f5780638da5cb5b14610242578063b13f450414610253578063f2fde38b1461027457600080fd5b80636d14a987146101ed578063715018a61461021457806377d175861461021c57600080fd5b8062482569146100d35780630a0593d11461012b57806310ea4f8a146101405780633a0b0ddd14610153578063683048351461017b5780636c08a879146101ba575b600080fd5b6101076100e1366004610de1565b60676020526000908152604090205463ffffffff811690640100000000900461ffff1682565b6040805163ffffffff909316835261ffff9091166020830152015b60405180910390f35b61013e610139366004610e6e565b610287565b005b61013e61014e366004610f8f565b6107a9565b610166610161366004610fcd565b6107bb565b60408051928352602083019190915201610122565b6101a27f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b039091168152602001610122565b6101dd6101c8366004610ff7565b60656020526000908152604090205460ff1681565b6040519015158152602001610122565b6101a27f000000000000000000000000000000000000000000000000000000000000000081565b61013e6107f7565b61013e61022a366004611087565b61080b565b61013e61023d36600461112a565b61081d565b6033546001600160a01b03166101a2565b610266610261366004610de1565b6109c3565b604051908152602001610122565b61013e610282366004610ff7565b610bc5565b3360009081526065602052604090205460ff16806102af57506033546001600160a01b031633145b6103115760405162461bcd60e51b815260206004820152602860248201527f456a6563746f723a204f6e6c79206f776e6572206f7220656a6563746f722063604482015267185b88195a9958dd60c21b60648201526084015b60405180910390fd5b60005b81518110156107a557806000610329826109c3565b905060008080805b87878151811061034357610343611200565b6020026020010151518160ff1610156106e35760007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316635401ed278a8a8151811061039957610399611200565b60200260200101518460ff16815181106103b5576103b5611200565b6020026020010151896040518363ffffffff1660e01b81526004016103e792919091825260ff16602082015260400190565b602060405180830381865afa158015610404573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104289190611216565b336000908152606560205260409020546001600160601b0391909116915060ff16801561046c575060ff871660009081526067602052604090205463ffffffff1615155b801561048057508561047e8287611255565b115b156104d6575060ff861660009081526066602090815260408083208151808301909252428252818301888152815460018181018455928652939094209151600290930290910191825591519082015591506106e3565b6104e08186611255565b94506104eb8461126d565b93507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316636e3b17db7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663296bb0648c8c8151811061055d5761055d611200565b60200260200101518660ff168151811061057957610579611200565b60200260200101516040518263ffffffff1660e01b815260040161059f91815260200190565b602060405180830381865afa1580156105bc573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105e09190611291565b6040516001600160f81b031960f88c901b1660208201526021016040516020818303038152906040526040518363ffffffff1660e01b81526004016106269291906112ae565b600060405180830381600087803b15801561064057600080fd5b505af1158015610654573d6000803e3d6000fd5b505050507f97ddb711c61a9d2d7effcba3e042a33862297f898d555655cca39ec4451f53b489898151811061068b5761068b611200565b60200260200101518360ff16815181106106a7576106a7611200565b6020026020010151886040516106ca92919091825260ff16602082015260400190565b60405180910390a1506106dc81611313565b9050610331565b508015801561070157503360009081526065602052604090205460ff165b1561074f5760ff851660009081526066602090815260408083208151808301909252428252818301878152815460018181018455928652939094209151600290930290910191825591519101555b6040805163ffffffff8416815282151560208201527f19dd87ae49ed14a795f8c2d5e8055bf2a4a9d01641a00a2f8f0a5a7bf7f70249910160405180910390a150505050508061079e90611333565b9050610314565b5050565b6107b1610c3e565b6107a58282610c98565b606660205281600052604060002081815481106107d757600080fd5b600091825260209091206002909102018054600190910154909250905082565b6107ff610c3e565b6108096000610cfc565b565b610813610c3e565b6107a58282610d4e565b600054610100900460ff161580801561083d5750600054600160ff909116105b806108575750303b158015610857575060005460ff166001145b6108ba5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b6064820152608401610308565b6000805460ff1916600117905580156108dd576000805461ff0019166101001790555b6108e684610cfc565b60005b83518160ff16101561092e5761091c848260ff168151811061090d5761090d611200565b60200260200101516001610c98565b8061092681611313565b9150506108e9565b5060005b82518160ff1610156109765761096481848360ff168151811061095757610957611200565b6020026020010151610d4e565b8061096e81611313565b915050610932565b5080156109bd576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50505050565b60ff811660009081526067602052604081205481906109e89063ffffffff164261134e565b60405163d5eccc0560e01b815260ff85166004820152909150600090612710907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063d5eccc0590602401602060405180830381865afa158015610a59573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a7d9190611216565b60ff8616600090815260676020526040902054610aaf916001600160601b031690640100000000900461ffff16611365565b610ab99190611384565b60ff8516600090815260666020526040812054919250908190610ae0575090949350505050565b60ff8616600090815260666020526040902054610aff9060019061134e565b90505b60ff86166000908152606660205260409020805485919083908110610b2957610b29611200565b9060005260206000209060020201600001541115610b9e5760ff86166000908152606660205260409020805482908110610b6557610b65611200565b90600052602060002090600202016001015482610b829190611255565b915080610b8e57610b9e565b610b97816113a6565b9050610b02565b828210610bb15750600095945050505050565b610bbb828461134e565b9695505050505050565b610bcd610c3e565b6001600160a01b038116610c325760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610308565b610c3b81610cfc565b50565b6033546001600160a01b031633146108095760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610308565b6001600160a01b038216600081815260656020908152604091829020805460ff19168515159081179091558251938452908301527f7676686b6d22e112412bd874d70177e011ab06602c26063f19f0386c9a3cee4291015b60405180910390a15050565b603380546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b60ff8216600081815260676020908152604091829020845181548684015163ffffffff90921665ffffffffffff19909116811764010000000061ffff90931692830217909255835194855291840152908201527fe69c2827a1e2fdd32265ebb4eeea5ee564f0551cf5dfed4150f8e116a67209eb90606001610cf0565b803560ff81168114610ddc57600080fd5b919050565b600060208284031215610df357600080fd5b610dfc82610dcb565b9392505050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff81118282101715610e4257610e42610e03565b604052919050565b600067ffffffffffffffff821115610e6457610e64610e03565b5060051b60200190565b60006020808385031215610e8157600080fd5b823567ffffffffffffffff80821115610e9957600080fd5b818501915085601f830112610ead57600080fd5b8135610ec0610ebb82610e4a565b610e19565b818152600591821b8401850191858201919089841115610edf57600080fd5b8686015b84811015610f6b57803586811115610efb5760008081fd5b8701603f81018c13610f0d5760008081fd5b888101356040610f1f610ebb83610e4a565b82815291851b83018101918b8101908f841115610f3c5760008081fd5b938201935b83851015610f5a5784358252938c0193908c0190610f41565b885250505093880193508701610ee3565b50909998505050505050505050565b6001600160a01b0381168114610c3b57600080fd5b60008060408385031215610fa257600080fd5b8235610fad81610f7a565b915060208301358015158114610fc257600080fd5b809150509250929050565b60008060408385031215610fe057600080fd5b610fe983610dcb565b946020939093013593505050565b60006020828403121561100957600080fd5b8135610dfc81610f7a565b60006040828403121561102657600080fd5b6040516040810181811067ffffffffffffffff8211171561104957611049610e03565b604052905080823563ffffffff8116811461106357600080fd5b8152602083013561ffff8116811461107a57600080fd5b6020919091015292915050565b6000806060838503121561109a57600080fd5b6110a383610dcb565b91506110b28460208501611014565b90509250929050565b600082601f8301126110cc57600080fd5b813560206110dc610ebb83610e4a565b82815260069290921b840181019181810190868411156110fb57600080fd5b8286015b8481101561111f576111118882611014565b8352918301916040016110ff565b509695505050505050565b60008060006060848603121561113f57600080fd5b833561114a81610f7a565b925060208481013567ffffffffffffffff8082111561116857600080fd5b818701915087601f83011261117c57600080fd5b813561118a610ebb82610e4a565b81815260059190911b8301840190848101908a8311156111a957600080fd5b938501935b828510156111d05784356111c181610f7a565b825293850193908501906111ae565b9650505060408701359250808311156111e857600080fd5b50506111f6868287016110bb565b9150509250925092565b634e487b7160e01b600052603260045260246000fd5b60006020828403121561122857600080fd5b81516001600160601b0381168114610dfc57600080fd5b634e487b7160e01b600052601160045260246000fd5b600082198211156112685761126861123f565b500190565b600063ffffffff808316818114156112875761128761123f565b6001019392505050565b6000602082840312156112a357600080fd5b8151610dfc81610f7a565b60018060a01b038316815260006020604081840152835180604085015260005b818110156112ea578581018301518582016060015282016112ce565b818111156112fc576000606083870101525b50601f01601f191692909201606001949350505050565b600060ff821660ff81141561132a5761132a61123f565b60010192915050565b60006000198214156113475761134761123f565b5060010190565b6000828210156113605761136061123f565b500390565b600081600019048311821515161561137f5761137f61123f565b500290565b6000826113a157634e487b7160e01b600052601260045260246000fd5b500490565b6000816113b5576113b561123f565b50600019019056fea2646970667358221220acbd8cd72aeec572935a1fe431c1a7916fb5a54ba6b3e6273dd51b7abb5334fb64736f6c634300080c0033",
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
