// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractIndexRegistry

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

// IIndexRegistryOperatorUpdate is an auto generated low-level Go binding around an user-defined struct.
type IIndexRegistryOperatorUpdate struct {
	FromBlockNumber uint32
	OperatorId      [32]byte
}

// IIndexRegistryQuorumUpdate is an auto generated low-level Go binding around an user-defined struct.
type IIndexRegistryQuorumUpdate struct {
	FromBlockNumber uint32
	NumOperators    uint32
}

// ContractIndexRegistryMetaData contains all meta data concerning the ContractIndexRegistry contract.
var ContractIndexRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"OPERATOR_DOES_NOT_EXIST_ID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deregisterOperator\",\"inputs\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getOperatorIndexUpdateOfIndexForQuorumAtIndex\",\"inputs\":[{\"name\":\"operatorIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIIndexRegistry.OperatorUpdate\",\"components\":[{\"name\":\"fromBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorListForQuorumAtBlockNumber\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"blockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumUpdateAtIndex\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIIndexRegistry.QuorumUpdate\",\"components\":[{\"name\":\"fromBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTotalOperatorsForQuorumAtBlockNumberByIndex\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"blockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"globalOperatorList\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operatorIdToIndex\",\"inputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerOperator\",\"inputs\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registryCoordinator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalOperatorsForQuorum\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumIndexUpdate\",\"inputs\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"newIndex\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false}]",
	Bin: "0x60a060405234801561001057600080fd5b5060405161142538038061142583398101604081905261002f9161010c565b6001600160a01b0381166080528061004561004c565b505061013c565b600054610100900460ff16156100b85760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840160405180910390fd5b60005460ff908116101561010a576000805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b60006020828403121561011e57600080fd5b81516001600160a01b038116811461013557600080fd5b9392505050565b6080516112c06101656000396000818161013001528181610257015261097801526112c06000f3fe608060405234801561001057600080fd5b50600436106100a85760003560e01c8063a454b3be11610071578063a454b3be1461018a578063a48bb0ac146101b2578063b81b2d3e146101e9578063bd29b8cd1461021a578063caa3cd761461022f578063f34109221461023757600080fd5b8062bff04d146100ad5780633a5c3c41146100d65780636653b53b1461010a5780636d14a9871461012b578063889ae3e51461016a575b600080fd5b6100c06100bb366004610f2e565b61024a565b6040516100cd9190610faa565b60405180910390f35b6100e96100e436600461101e565b610453565b60408051825163ffffffff16815260209283015192810192909252016100cd565b61011d610118366004611061565b6104d9565b6040519081526020016100cd565b6101527f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b0390911681526020016100cd565b61017d61017836600461107a565b6104fa565b6040516100cd91906110ad565b61019d6101983660046110e5565b61066a565b60405163ffffffff90911681526020016100cd565b6101c56101c036600461107a565b6108f6565b60408051825163ffffffff90811682526020938401511692810192909252016100cd565b61019d6101f7366004611111565b600260209081526000928352604080842090915290825290205463ffffffff1681565b61022d610228366004610f2e565b61096d565b005b61011d600081565b61019d61024536600461113b565b610a83565b6060336001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161461029d5760405162461bcd60e51b815260040161029490611156565b60405180910390fd5b60008267ffffffffffffffff8111156102b8576102b86111c9565b6040519080825280602002602001820160405280156102e1578160200160208202803683370190505b50905060005b83811015610448576000858583818110610303576103036111df565b919091013560f81c6000818152600460205260408120549193509091508161032c576000610373565b60ff8316600090815260046020526040902061034960018461120b565b81548110610359576103596111df565b600091825260209091200154600160201b900463ffffffff165b9050610380898483610a8e565b6103fb8361038f836001611222565b60408051808201825263ffffffff9283166020808301918252438516835260ff9590951660009081526004865292832080546001810182559084529490922090519301805491518316600160201b0267ffffffffffffffff199092169390921692909217919091179055565b610406816001611222565b858581518110610418576104186111df565b602002602001019063ffffffff16908163ffffffff168152505050505080806104409061124a565b9150506102e7565b5090505b9392505050565b604080518082019091526000808252602082015260ff8316600090815260036020908152604080832063ffffffff808916855292529091208054909184169081106104a0576104a06111df565b600091825260209182902060408051808201909152600290920201805463ffffffff168252600101549181019190915290509392505050565b600181815481106104e957600080fd5b600091825260209091200154905081565b606060006105088484610b56565b63ffffffff1667ffffffffffffffff811115610526576105266111c9565b60405190808252806020026020018201604052801561054f578160200160208202803683370190505b50905060005b815181101561066057610569818686610cf9565b82828151811061057b5761057b6111df565b6020026020010181815250506000801b82828151811061059d5761059d6111df565b6020026020010151141561064e5760405162461bcd60e51b815260206004820152606660248201527f496e64657852656769737472792e6765744f70657261746f724c697374466f7260448201527f51756f72756d4174426c6f636b4e756d6265723a206f70657261746f7220646f60648201527f6573206e6f742065786973742061742074686520676976656e20626c6f636b20608482015265373ab6b132b960d11b60a482015260c401610294565b806106588161124a565b915050610555565b5090505b92915050565b60ff83166000908152600460205260408120805482919063ffffffff8516908110610697576106976111df565b60009182526020918290206040805180820190915291015463ffffffff808216808452600160201b90920481169383019390935290925090851610156107915760405162461bcd60e51b815260206004820152607d60248201527f496e64657852656769737472792e676574546f74616c4f70657261746f72734660448201527f6f7251756f72756d4174426c6f636b4e756d6265724279496e6465783a20707260648201527f6f766964656420696e64657820697320746f6f2066617220696e20746865207060848201527f61737420666f722070726f766964656420626c6f636b206e756d62657200000060a482015260c401610294565b60ff85166000908152600460205260409020546107b09060019061120b565b8363ffffffff16146108ea5760ff851660009081526004602052604081206107d9856001611222565b63ffffffff16815481106107ef576107ef6111df565b60009182526020918290206040805180820190915291015463ffffffff808216808452600160201b909204811693830193909352909250908616106108e85760405162461bcd60e51b815260206004820152607f60248201527f496e64657852656769737472792e676574546f74616c4f70657261746f72734660448201527f6f7251756f72756d4174426c6f636b4e756d6265724279496e6465783a20707260648201527f6f766964656420696e64657820697320746f6f2066617220696e20746865206660848201527f757475726520666f722070726f766964656420626c6f636b206e756d6265720060a482015260c401610294565b505b60200151949350505050565b604080518082019091526000808252602082015260ff83166000908152600460205260409020805463ffffffff8416908110610934576109346111df565b60009182526020918290206040805180820190915291015463ffffffff8082168352600160201b90910416918101919091529392505050565b336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146109b55760405162461bcd60e51b815260040161029490611156565b60005b81811015610a7d5760008383838181106109d4576109d46111df565b919091013560f81c60008181526002602090815260408083208a845290915290205490925063ffffffff169050610a0c868383610de9565b60ff821660009081526004602052604090208054610a6891849160019190610a3590839061120b565b81548110610a4557610a456111df565b60009182526020909120015461038f9190600160201b900463ffffffff16611265565b50508080610a759061124a565b9150506109b8565b50505050565b600061066482610ec0565b604080518082018252602080820186815263ffffffff438116845260ff87166000818152600385528681208884168083529086528782208054600180820183559184528784208951600292830290910180549190971663ffffffff19918216178755965195909101949094558282529285528681208a8252855286902080549093168217909255845191825291810191909152909185917f6ee1e4f4075f3d067176140d34e87874244dd273294c05b2218133e49a2ba6f6910160405180910390a250505050565b60ff821660009081526004602052604081205480610b78576000915050610664565b60ff841660009081526004602052604081208054909190610b9b57610b9b6111df565b60009182526020909120015463ffffffff9081169084161015610bc2576000915050610664565b60005b610bd060018361120b565b8111610cb457600081610be460018561120b565b610bee919061120b565b60ff871660009081526004602052604081208054929350909183908110610c1757610c176111df565b60009182526020918290206040805180820190915291015463ffffffff808216808452600160201b90920481169383019390935290925090871610610c9f5760ff87166000908152600460205260409020805483908110610c7a57610c7a6111df565b600091825260209091200154600160201b900463ffffffff1694506106649350505050565b50508080610cac9061124a565b915050610bc5565b5060ff841660009081526004602052604081208054909190610cd857610cd86111df565b600091825260209091200154600160201b900463ffffffff16949350505050565b60ff8216600090815260036020908152604080832063ffffffff87168452909152812054815b81811015610ddd57600081610d3560018561120b565b610d3f919061120b565b60ff8716600090815260036020908152604080832063ffffffff8c16845290915281208054929350909183908110610d7957610d796111df565b600091825260209182902060408051808201909152600290920201805463ffffffff9081168084526001909201549383019390935290925090871610610dc85760200151935061044c92505050565b50508080610dd59061124a565b915050610d1f565b50600095945050505050565b6000610df483610ec0565b60ff841660009081526003602052604081209192509081610e16600185611265565b63ffffffff1681526020808201929092526040908101600090812060ff881682526003909352908120600191610e4c8387611265565b63ffffffff168152602081019190915260400160002054610e6d919061120b565b81548110610e7d57610e7d6111df565b9060005260206000209060020201600101549050808514610ea357610ea3818585610a8e565b610eb9600085610eb4600186611265565b610a8e565b5050505050565b60ff811660009081526004602052604081205480610ee15750600092915050565b60ff83166000908152600460205260409020610efe60018361120b565b81548110610f0e57610f0e6111df565b600091825260209091200154600160201b900463ffffffff169392505050565b600080600060408486031215610f4357600080fd5b83359250602084013567ffffffffffffffff80821115610f6257600080fd5b818601915086601f830112610f7657600080fd5b813581811115610f8557600080fd5b876020828501011115610f9757600080fd5b6020830194508093505050509250925092565b6020808252825182820181905260009190848201906040850190845b81811015610fe857835163ffffffff1683529284019291840191600101610fc6565b50909695505050505050565b803563ffffffff8116811461100857600080fd5b919050565b803560ff8116811461100857600080fd5b60008060006060848603121561103357600080fd5b61103c84610ff4565b925061104a6020850161100d565b915061105860408501610ff4565b90509250925092565b60006020828403121561107357600080fd5b5035919050565b6000806040838503121561108d57600080fd5b6110968361100d565b91506110a460208401610ff4565b90509250929050565b6020808252825182820181905260009190848201906040850190845b81811015610fe8578351835292840192918401916001016110c9565b6000806000606084860312156110fa57600080fd5b6111038461100d565b925061104a60208501610ff4565b6000806040838503121561112457600080fd5b61112d8361100d565b946020939093013593505050565b60006020828403121561114d57600080fd5b61044c8261100d565b6020808252604d908201527f496e64657852656769737472792e6f6e6c795265676973747279436f6f72646960408201527f6e61746f723a2063616c6c6572206973206e6f7420746865207265676973747260608201526c3c9031b7b7b93234b730ba37b960991b608082015260a00190565b634e487b7160e01b600052604160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b60008282101561121d5761121d6111f5565b500390565b600063ffffffff808316818516808303821115611241576112416111f5565b01949350505050565b600060001982141561125e5761125e6111f5565b5060010190565b600063ffffffff83811690831681811015611282576112826111f5565b03939250505056fea2646970667358221220c2b30bd928dde59d05244bf7421f8952673593242c53aa336a8ec1a811bf9d8064736f6c634300080c0033",
}

// ContractIndexRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractIndexRegistryMetaData.ABI instead.
var ContractIndexRegistryABI = ContractIndexRegistryMetaData.ABI

// ContractIndexRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractIndexRegistryMetaData.Bin instead.
var ContractIndexRegistryBin = ContractIndexRegistryMetaData.Bin

// DeployContractIndexRegistry deploys a new Ethereum contract, binding an instance of ContractIndexRegistry to it.
func DeployContractIndexRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, _registryCoordinator common.Address) (common.Address, *types.Transaction, *ContractIndexRegistry, error) {
	parsed, err := ContractIndexRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractIndexRegistryBin), backend, _registryCoordinator)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractIndexRegistry{ContractIndexRegistryCaller: ContractIndexRegistryCaller{contract: contract}, ContractIndexRegistryTransactor: ContractIndexRegistryTransactor{contract: contract}, ContractIndexRegistryFilterer: ContractIndexRegistryFilterer{contract: contract}}, nil
}

// ContractIndexRegistry is an auto generated Go binding around an Ethereum contract.
type ContractIndexRegistry struct {
	ContractIndexRegistryCaller     // Read-only binding to the contract
	ContractIndexRegistryTransactor // Write-only binding to the contract
	ContractIndexRegistryFilterer   // Log filterer for contract events
}

// ContractIndexRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractIndexRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIndexRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractIndexRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIndexRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractIndexRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIndexRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractIndexRegistrySession struct {
	Contract     *ContractIndexRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// ContractIndexRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractIndexRegistryCallerSession struct {
	Contract *ContractIndexRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// ContractIndexRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractIndexRegistryTransactorSession struct {
	Contract     *ContractIndexRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// ContractIndexRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractIndexRegistryRaw struct {
	Contract *ContractIndexRegistry // Generic contract binding to access the raw methods on
}

// ContractIndexRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractIndexRegistryCallerRaw struct {
	Contract *ContractIndexRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// ContractIndexRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractIndexRegistryTransactorRaw struct {
	Contract *ContractIndexRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractIndexRegistry creates a new instance of ContractIndexRegistry, bound to a specific deployed contract.
func NewContractIndexRegistry(address common.Address, backend bind.ContractBackend) (*ContractIndexRegistry, error) {
	contract, err := bindContractIndexRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractIndexRegistry{ContractIndexRegistryCaller: ContractIndexRegistryCaller{contract: contract}, ContractIndexRegistryTransactor: ContractIndexRegistryTransactor{contract: contract}, ContractIndexRegistryFilterer: ContractIndexRegistryFilterer{contract: contract}}, nil
}

// NewContractIndexRegistryCaller creates a new read-only instance of ContractIndexRegistry, bound to a specific deployed contract.
func NewContractIndexRegistryCaller(address common.Address, caller bind.ContractCaller) (*ContractIndexRegistryCaller, error) {
	contract, err := bindContractIndexRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractIndexRegistryCaller{contract: contract}, nil
}

// NewContractIndexRegistryTransactor creates a new write-only instance of ContractIndexRegistry, bound to a specific deployed contract.
func NewContractIndexRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractIndexRegistryTransactor, error) {
	contract, err := bindContractIndexRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractIndexRegistryTransactor{contract: contract}, nil
}

// NewContractIndexRegistryFilterer creates a new log filterer instance of ContractIndexRegistry, bound to a specific deployed contract.
func NewContractIndexRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractIndexRegistryFilterer, error) {
	contract, err := bindContractIndexRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractIndexRegistryFilterer{contract: contract}, nil
}

// bindContractIndexRegistry binds a generic wrapper to an already deployed contract.
func bindContractIndexRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractIndexRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractIndexRegistry *ContractIndexRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractIndexRegistry.Contract.ContractIndexRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractIndexRegistry *ContractIndexRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIndexRegistry.Contract.ContractIndexRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractIndexRegistry *ContractIndexRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractIndexRegistry.Contract.ContractIndexRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractIndexRegistry *ContractIndexRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractIndexRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractIndexRegistry *ContractIndexRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIndexRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractIndexRegistry *ContractIndexRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractIndexRegistry.Contract.contract.Transact(opts, method, params...)
}

// OPERATORDOESNOTEXISTID is a free data retrieval call binding the contract method 0xcaa3cd76.
//
// Solidity: function OPERATOR_DOES_NOT_EXIST_ID() view returns(bytes32)
func (_ContractIndexRegistry *ContractIndexRegistryCaller) OPERATORDOESNOTEXISTID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ContractIndexRegistry.contract.Call(opts, &out, "OPERATOR_DOES_NOT_EXIST_ID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// OPERATORDOESNOTEXISTID is a free data retrieval call binding the contract method 0xcaa3cd76.
//
// Solidity: function OPERATOR_DOES_NOT_EXIST_ID() view returns(bytes32)
func (_ContractIndexRegistry *ContractIndexRegistrySession) OPERATORDOESNOTEXISTID() ([32]byte, error) {
	return _ContractIndexRegistry.Contract.OPERATORDOESNOTEXISTID(&_ContractIndexRegistry.CallOpts)
}

// OPERATORDOESNOTEXISTID is a free data retrieval call binding the contract method 0xcaa3cd76.
//
// Solidity: function OPERATOR_DOES_NOT_EXIST_ID() view returns(bytes32)
func (_ContractIndexRegistry *ContractIndexRegistryCallerSession) OPERATORDOESNOTEXISTID() ([32]byte, error) {
	return _ContractIndexRegistry.Contract.OPERATORDOESNOTEXISTID(&_ContractIndexRegistry.CallOpts)
}

// GetOperatorIndexUpdateOfIndexForQuorumAtIndex is a free data retrieval call binding the contract method 0x3a5c3c41.
//
// Solidity: function getOperatorIndexUpdateOfIndexForQuorumAtIndex(uint32 operatorIndex, uint8 quorumNumber, uint32 index) view returns((uint32,bytes32))
func (_ContractIndexRegistry *ContractIndexRegistryCaller) GetOperatorIndexUpdateOfIndexForQuorumAtIndex(opts *bind.CallOpts, operatorIndex uint32, quorumNumber uint8, index uint32) (IIndexRegistryOperatorUpdate, error) {
	var out []interface{}
	err := _ContractIndexRegistry.contract.Call(opts, &out, "getOperatorIndexUpdateOfIndexForQuorumAtIndex", operatorIndex, quorumNumber, index)

	if err != nil {
		return *new(IIndexRegistryOperatorUpdate), err
	}

	out0 := *abi.ConvertType(out[0], new(IIndexRegistryOperatorUpdate)).(*IIndexRegistryOperatorUpdate)

	return out0, err

}

// GetOperatorIndexUpdateOfIndexForQuorumAtIndex is a free data retrieval call binding the contract method 0x3a5c3c41.
//
// Solidity: function getOperatorIndexUpdateOfIndexForQuorumAtIndex(uint32 operatorIndex, uint8 quorumNumber, uint32 index) view returns((uint32,bytes32))
func (_ContractIndexRegistry *ContractIndexRegistrySession) GetOperatorIndexUpdateOfIndexForQuorumAtIndex(operatorIndex uint32, quorumNumber uint8, index uint32) (IIndexRegistryOperatorUpdate, error) {
	return _ContractIndexRegistry.Contract.GetOperatorIndexUpdateOfIndexForQuorumAtIndex(&_ContractIndexRegistry.CallOpts, operatorIndex, quorumNumber, index)
}

// GetOperatorIndexUpdateOfIndexForQuorumAtIndex is a free data retrieval call binding the contract method 0x3a5c3c41.
//
// Solidity: function getOperatorIndexUpdateOfIndexForQuorumAtIndex(uint32 operatorIndex, uint8 quorumNumber, uint32 index) view returns((uint32,bytes32))
func (_ContractIndexRegistry *ContractIndexRegistryCallerSession) GetOperatorIndexUpdateOfIndexForQuorumAtIndex(operatorIndex uint32, quorumNumber uint8, index uint32) (IIndexRegistryOperatorUpdate, error) {
	return _ContractIndexRegistry.Contract.GetOperatorIndexUpdateOfIndexForQuorumAtIndex(&_ContractIndexRegistry.CallOpts, operatorIndex, quorumNumber, index)
}

// GetOperatorListForQuorumAtBlockNumber is a free data retrieval call binding the contract method 0x889ae3e5.
//
// Solidity: function getOperatorListForQuorumAtBlockNumber(uint8 quorumNumber, uint32 blockNumber) view returns(bytes32[])
func (_ContractIndexRegistry *ContractIndexRegistryCaller) GetOperatorListForQuorumAtBlockNumber(opts *bind.CallOpts, quorumNumber uint8, blockNumber uint32) ([][32]byte, error) {
	var out []interface{}
	err := _ContractIndexRegistry.contract.Call(opts, &out, "getOperatorListForQuorumAtBlockNumber", quorumNumber, blockNumber)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetOperatorListForQuorumAtBlockNumber is a free data retrieval call binding the contract method 0x889ae3e5.
//
// Solidity: function getOperatorListForQuorumAtBlockNumber(uint8 quorumNumber, uint32 blockNumber) view returns(bytes32[])
func (_ContractIndexRegistry *ContractIndexRegistrySession) GetOperatorListForQuorumAtBlockNumber(quorumNumber uint8, blockNumber uint32) ([][32]byte, error) {
	return _ContractIndexRegistry.Contract.GetOperatorListForQuorumAtBlockNumber(&_ContractIndexRegistry.CallOpts, quorumNumber, blockNumber)
}

// GetOperatorListForQuorumAtBlockNumber is a free data retrieval call binding the contract method 0x889ae3e5.
//
// Solidity: function getOperatorListForQuorumAtBlockNumber(uint8 quorumNumber, uint32 blockNumber) view returns(bytes32[])
func (_ContractIndexRegistry *ContractIndexRegistryCallerSession) GetOperatorListForQuorumAtBlockNumber(quorumNumber uint8, blockNumber uint32) ([][32]byte, error) {
	return _ContractIndexRegistry.Contract.GetOperatorListForQuorumAtBlockNumber(&_ContractIndexRegistry.CallOpts, quorumNumber, blockNumber)
}

// GetQuorumUpdateAtIndex is a free data retrieval call binding the contract method 0xa48bb0ac.
//
// Solidity: function getQuorumUpdateAtIndex(uint8 quorumNumber, uint32 index) view returns((uint32,uint32))
func (_ContractIndexRegistry *ContractIndexRegistryCaller) GetQuorumUpdateAtIndex(opts *bind.CallOpts, quorumNumber uint8, index uint32) (IIndexRegistryQuorumUpdate, error) {
	var out []interface{}
	err := _ContractIndexRegistry.contract.Call(opts, &out, "getQuorumUpdateAtIndex", quorumNumber, index)

	if err != nil {
		return *new(IIndexRegistryQuorumUpdate), err
	}

	out0 := *abi.ConvertType(out[0], new(IIndexRegistryQuorumUpdate)).(*IIndexRegistryQuorumUpdate)

	return out0, err

}

// GetQuorumUpdateAtIndex is a free data retrieval call binding the contract method 0xa48bb0ac.
//
// Solidity: function getQuorumUpdateAtIndex(uint8 quorumNumber, uint32 index) view returns((uint32,uint32))
func (_ContractIndexRegistry *ContractIndexRegistrySession) GetQuorumUpdateAtIndex(quorumNumber uint8, index uint32) (IIndexRegistryQuorumUpdate, error) {
	return _ContractIndexRegistry.Contract.GetQuorumUpdateAtIndex(&_ContractIndexRegistry.CallOpts, quorumNumber, index)
}

// GetQuorumUpdateAtIndex is a free data retrieval call binding the contract method 0xa48bb0ac.
//
// Solidity: function getQuorumUpdateAtIndex(uint8 quorumNumber, uint32 index) view returns((uint32,uint32))
func (_ContractIndexRegistry *ContractIndexRegistryCallerSession) GetQuorumUpdateAtIndex(quorumNumber uint8, index uint32) (IIndexRegistryQuorumUpdate, error) {
	return _ContractIndexRegistry.Contract.GetQuorumUpdateAtIndex(&_ContractIndexRegistry.CallOpts, quorumNumber, index)
}

// GetTotalOperatorsForQuorumAtBlockNumberByIndex is a free data retrieval call binding the contract method 0xa454b3be.
//
// Solidity: function getTotalOperatorsForQuorumAtBlockNumberByIndex(uint8 quorumNumber, uint32 blockNumber, uint32 index) view returns(uint32)
func (_ContractIndexRegistry *ContractIndexRegistryCaller) GetTotalOperatorsForQuorumAtBlockNumberByIndex(opts *bind.CallOpts, quorumNumber uint8, blockNumber uint32, index uint32) (uint32, error) {
	var out []interface{}
	err := _ContractIndexRegistry.contract.Call(opts, &out, "getTotalOperatorsForQuorumAtBlockNumberByIndex", quorumNumber, blockNumber, index)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// GetTotalOperatorsForQuorumAtBlockNumberByIndex is a free data retrieval call binding the contract method 0xa454b3be.
//
// Solidity: function getTotalOperatorsForQuorumAtBlockNumberByIndex(uint8 quorumNumber, uint32 blockNumber, uint32 index) view returns(uint32)
func (_ContractIndexRegistry *ContractIndexRegistrySession) GetTotalOperatorsForQuorumAtBlockNumberByIndex(quorumNumber uint8, blockNumber uint32, index uint32) (uint32, error) {
	return _ContractIndexRegistry.Contract.GetTotalOperatorsForQuorumAtBlockNumberByIndex(&_ContractIndexRegistry.CallOpts, quorumNumber, blockNumber, index)
}

// GetTotalOperatorsForQuorumAtBlockNumberByIndex is a free data retrieval call binding the contract method 0xa454b3be.
//
// Solidity: function getTotalOperatorsForQuorumAtBlockNumberByIndex(uint8 quorumNumber, uint32 blockNumber, uint32 index) view returns(uint32)
func (_ContractIndexRegistry *ContractIndexRegistryCallerSession) GetTotalOperatorsForQuorumAtBlockNumberByIndex(quorumNumber uint8, blockNumber uint32, index uint32) (uint32, error) {
	return _ContractIndexRegistry.Contract.GetTotalOperatorsForQuorumAtBlockNumberByIndex(&_ContractIndexRegistry.CallOpts, quorumNumber, blockNumber, index)
}

// GlobalOperatorList is a free data retrieval call binding the contract method 0x6653b53b.
//
// Solidity: function globalOperatorList(uint256 ) view returns(bytes32)
func (_ContractIndexRegistry *ContractIndexRegistryCaller) GlobalOperatorList(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _ContractIndexRegistry.contract.Call(opts, &out, "globalOperatorList", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GlobalOperatorList is a free data retrieval call binding the contract method 0x6653b53b.
//
// Solidity: function globalOperatorList(uint256 ) view returns(bytes32)
func (_ContractIndexRegistry *ContractIndexRegistrySession) GlobalOperatorList(arg0 *big.Int) ([32]byte, error) {
	return _ContractIndexRegistry.Contract.GlobalOperatorList(&_ContractIndexRegistry.CallOpts, arg0)
}

// GlobalOperatorList is a free data retrieval call binding the contract method 0x6653b53b.
//
// Solidity: function globalOperatorList(uint256 ) view returns(bytes32)
func (_ContractIndexRegistry *ContractIndexRegistryCallerSession) GlobalOperatorList(arg0 *big.Int) ([32]byte, error) {
	return _ContractIndexRegistry.Contract.GlobalOperatorList(&_ContractIndexRegistry.CallOpts, arg0)
}

// OperatorIdToIndex is a free data retrieval call binding the contract method 0xb81b2d3e.
//
// Solidity: function operatorIdToIndex(uint8 , bytes32 ) view returns(uint32)
func (_ContractIndexRegistry *ContractIndexRegistryCaller) OperatorIdToIndex(opts *bind.CallOpts, arg0 uint8, arg1 [32]byte) (uint32, error) {
	var out []interface{}
	err := _ContractIndexRegistry.contract.Call(opts, &out, "operatorIdToIndex", arg0, arg1)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// OperatorIdToIndex is a free data retrieval call binding the contract method 0xb81b2d3e.
//
// Solidity: function operatorIdToIndex(uint8 , bytes32 ) view returns(uint32)
func (_ContractIndexRegistry *ContractIndexRegistrySession) OperatorIdToIndex(arg0 uint8, arg1 [32]byte) (uint32, error) {
	return _ContractIndexRegistry.Contract.OperatorIdToIndex(&_ContractIndexRegistry.CallOpts, arg0, arg1)
}

// OperatorIdToIndex is a free data retrieval call binding the contract method 0xb81b2d3e.
//
// Solidity: function operatorIdToIndex(uint8 , bytes32 ) view returns(uint32)
func (_ContractIndexRegistry *ContractIndexRegistryCallerSession) OperatorIdToIndex(arg0 uint8, arg1 [32]byte) (uint32, error) {
	return _ContractIndexRegistry.Contract.OperatorIdToIndex(&_ContractIndexRegistry.CallOpts, arg0, arg1)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractIndexRegistry *ContractIndexRegistryCaller) RegistryCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractIndexRegistry.contract.Call(opts, &out, "registryCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractIndexRegistry *ContractIndexRegistrySession) RegistryCoordinator() (common.Address, error) {
	return _ContractIndexRegistry.Contract.RegistryCoordinator(&_ContractIndexRegistry.CallOpts)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractIndexRegistry *ContractIndexRegistryCallerSession) RegistryCoordinator() (common.Address, error) {
	return _ContractIndexRegistry.Contract.RegistryCoordinator(&_ContractIndexRegistry.CallOpts)
}

// TotalOperatorsForQuorum is a free data retrieval call binding the contract method 0xf3410922.
//
// Solidity: function totalOperatorsForQuorum(uint8 quorumNumber) view returns(uint32)
func (_ContractIndexRegistry *ContractIndexRegistryCaller) TotalOperatorsForQuorum(opts *bind.CallOpts, quorumNumber uint8) (uint32, error) {
	var out []interface{}
	err := _ContractIndexRegistry.contract.Call(opts, &out, "totalOperatorsForQuorum", quorumNumber)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// TotalOperatorsForQuorum is a free data retrieval call binding the contract method 0xf3410922.
//
// Solidity: function totalOperatorsForQuorum(uint8 quorumNumber) view returns(uint32)
func (_ContractIndexRegistry *ContractIndexRegistrySession) TotalOperatorsForQuorum(quorumNumber uint8) (uint32, error) {
	return _ContractIndexRegistry.Contract.TotalOperatorsForQuorum(&_ContractIndexRegistry.CallOpts, quorumNumber)
}

// TotalOperatorsForQuorum is a free data retrieval call binding the contract method 0xf3410922.
//
// Solidity: function totalOperatorsForQuorum(uint8 quorumNumber) view returns(uint32)
func (_ContractIndexRegistry *ContractIndexRegistryCallerSession) TotalOperatorsForQuorum(quorumNumber uint8) (uint32, error) {
	return _ContractIndexRegistry.Contract.TotalOperatorsForQuorum(&_ContractIndexRegistry.CallOpts, quorumNumber)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xbd29b8cd.
//
// Solidity: function deregisterOperator(bytes32 operatorId, bytes quorumNumbers) returns()
func (_ContractIndexRegistry *ContractIndexRegistryTransactor) DeregisterOperator(opts *bind.TransactOpts, operatorId [32]byte, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractIndexRegistry.contract.Transact(opts, "deregisterOperator", operatorId, quorumNumbers)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xbd29b8cd.
//
// Solidity: function deregisterOperator(bytes32 operatorId, bytes quorumNumbers) returns()
func (_ContractIndexRegistry *ContractIndexRegistrySession) DeregisterOperator(operatorId [32]byte, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractIndexRegistry.Contract.DeregisterOperator(&_ContractIndexRegistry.TransactOpts, operatorId, quorumNumbers)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xbd29b8cd.
//
// Solidity: function deregisterOperator(bytes32 operatorId, bytes quorumNumbers) returns()
func (_ContractIndexRegistry *ContractIndexRegistryTransactorSession) DeregisterOperator(operatorId [32]byte, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractIndexRegistry.Contract.DeregisterOperator(&_ContractIndexRegistry.TransactOpts, operatorId, quorumNumbers)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x00bff04d.
//
// Solidity: function registerOperator(bytes32 operatorId, bytes quorumNumbers) returns(uint32[])
func (_ContractIndexRegistry *ContractIndexRegistryTransactor) RegisterOperator(opts *bind.TransactOpts, operatorId [32]byte, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractIndexRegistry.contract.Transact(opts, "registerOperator", operatorId, quorumNumbers)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x00bff04d.
//
// Solidity: function registerOperator(bytes32 operatorId, bytes quorumNumbers) returns(uint32[])
func (_ContractIndexRegistry *ContractIndexRegistrySession) RegisterOperator(operatorId [32]byte, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractIndexRegistry.Contract.RegisterOperator(&_ContractIndexRegistry.TransactOpts, operatorId, quorumNumbers)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x00bff04d.
//
// Solidity: function registerOperator(bytes32 operatorId, bytes quorumNumbers) returns(uint32[])
func (_ContractIndexRegistry *ContractIndexRegistryTransactorSession) RegisterOperator(operatorId [32]byte, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractIndexRegistry.Contract.RegisterOperator(&_ContractIndexRegistry.TransactOpts, operatorId, quorumNumbers)
}

// ContractIndexRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ContractIndexRegistry contract.
type ContractIndexRegistryInitializedIterator struct {
	Event *ContractIndexRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *ContractIndexRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIndexRegistryInitialized)
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
		it.Event = new(ContractIndexRegistryInitialized)
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
func (it *ContractIndexRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIndexRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIndexRegistryInitialized represents a Initialized event raised by the ContractIndexRegistry contract.
type ContractIndexRegistryInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractIndexRegistry *ContractIndexRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*ContractIndexRegistryInitializedIterator, error) {

	logs, sub, err := _ContractIndexRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ContractIndexRegistryInitializedIterator{contract: _ContractIndexRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractIndexRegistry *ContractIndexRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ContractIndexRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _ContractIndexRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIndexRegistryInitialized)
				if err := _ContractIndexRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ContractIndexRegistry *ContractIndexRegistryFilterer) ParseInitialized(log types.Log) (*ContractIndexRegistryInitialized, error) {
	event := new(ContractIndexRegistryInitialized)
	if err := _ContractIndexRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIndexRegistryQuorumIndexUpdateIterator is returned from FilterQuorumIndexUpdate and is used to iterate over the raw logs and unpacked data for QuorumIndexUpdate events raised by the ContractIndexRegistry contract.
type ContractIndexRegistryQuorumIndexUpdateIterator struct {
	Event *ContractIndexRegistryQuorumIndexUpdate // Event containing the contract specifics and raw log

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
func (it *ContractIndexRegistryQuorumIndexUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIndexRegistryQuorumIndexUpdate)
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
		it.Event = new(ContractIndexRegistryQuorumIndexUpdate)
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
func (it *ContractIndexRegistryQuorumIndexUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIndexRegistryQuorumIndexUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIndexRegistryQuorumIndexUpdate represents a QuorumIndexUpdate event raised by the ContractIndexRegistry contract.
type ContractIndexRegistryQuorumIndexUpdate struct {
	OperatorId   [32]byte
	QuorumNumber uint8
	NewIndex     uint32
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterQuorumIndexUpdate is a free log retrieval operation binding the contract event 0x6ee1e4f4075f3d067176140d34e87874244dd273294c05b2218133e49a2ba6f6.
//
// Solidity: event QuorumIndexUpdate(bytes32 indexed operatorId, uint8 quorumNumber, uint32 newIndex)
func (_ContractIndexRegistry *ContractIndexRegistryFilterer) FilterQuorumIndexUpdate(opts *bind.FilterOpts, operatorId [][32]byte) (*ContractIndexRegistryQuorumIndexUpdateIterator, error) {

	var operatorIdRule []interface{}
	for _, operatorIdItem := range operatorId {
		operatorIdRule = append(operatorIdRule, operatorIdItem)
	}

	logs, sub, err := _ContractIndexRegistry.contract.FilterLogs(opts, "QuorumIndexUpdate", operatorIdRule)
	if err != nil {
		return nil, err
	}
	return &ContractIndexRegistryQuorumIndexUpdateIterator{contract: _ContractIndexRegistry.contract, event: "QuorumIndexUpdate", logs: logs, sub: sub}, nil
}

// WatchQuorumIndexUpdate is a free log subscription operation binding the contract event 0x6ee1e4f4075f3d067176140d34e87874244dd273294c05b2218133e49a2ba6f6.
//
// Solidity: event QuorumIndexUpdate(bytes32 indexed operatorId, uint8 quorumNumber, uint32 newIndex)
func (_ContractIndexRegistry *ContractIndexRegistryFilterer) WatchQuorumIndexUpdate(opts *bind.WatchOpts, sink chan<- *ContractIndexRegistryQuorumIndexUpdate, operatorId [][32]byte) (event.Subscription, error) {

	var operatorIdRule []interface{}
	for _, operatorIdItem := range operatorId {
		operatorIdRule = append(operatorIdRule, operatorIdItem)
	}

	logs, sub, err := _ContractIndexRegistry.contract.WatchLogs(opts, "QuorumIndexUpdate", operatorIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIndexRegistryQuorumIndexUpdate)
				if err := _ContractIndexRegistry.contract.UnpackLog(event, "QuorumIndexUpdate", log); err != nil {
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

// ParseQuorumIndexUpdate is a log parse operation binding the contract event 0x6ee1e4f4075f3d067176140d34e87874244dd273294c05b2218133e49a2ba6f6.
//
// Solidity: event QuorumIndexUpdate(bytes32 indexed operatorId, uint8 quorumNumber, uint32 newIndex)
func (_ContractIndexRegistry *ContractIndexRegistryFilterer) ParseQuorumIndexUpdate(log types.Log) (*ContractIndexRegistryQuorumIndexUpdate, error) {
	event := new(ContractIndexRegistryQuorumIndexUpdate)
	if err := _ContractIndexRegistry.contract.UnpackLog(event, "QuorumIndexUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
