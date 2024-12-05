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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addVersionedBlobParams\",\"inputs\":[{\"name\":\"_versionedBlobParams\",\"type\":\"tuple\",\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"defaultSecurityThresholdsV2\",\"inputs\":[],\"outputs\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBlobParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDefaultSecurityThresholdsV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIsQuorumRequired\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumAdversaryThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumConfirmationThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_quorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_quorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_quorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_versionedBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structVersionedBlobParams[]\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"_defaultSecurityThresholdsV2\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"nextBlobVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumAdversaryThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumConfirmationThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateDefaultSecurityThresholdsV2\",\"inputs\":[{\"name\":\"_defaultSecurityThresholdsV2\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateQuorumAdversaryThresholdPercentages\",\"inputs\":[{\"name\":\"_quorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateQuorumConfirmationThresholdPercentages\",\"inputs\":[{\"name\":\"_quorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateQuorumNumbersRequired\",\"inputs\":[{\"name\":\"_quorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"versionedBlobParams\",\"inputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VersionedBlobParamsAdded\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"indexed\":true,\"internalType\":\"uint16\"},{\"name\":\"versionedBlobParams\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"anonymous\":false}]",
	Bin: "0x608060405234801561001057600080fd5b5061001961001e565b6100de565b603254610100900460ff161561008a5760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840160405180910390fd5b60325460ff90811610156100dc576032805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b6111b8806100ed6000396000f3fe608060405234801561001057600080fd5b506004361061012b5760003560e01c80638a476982116100ad578063ee6c3bcf11610071578063ee6c3bcf14610317578063ef6355291461032a578063f2fde38b14610375578063f74e363c14610388578063fb87355e146103ed57600080fd5b80638a476982146102c65780638da5cb5b146102d9578063a5e9b2eb146102f4578063bafa910714610307578063e15234ff1461030f57600080fd5b806332430f14116100f457806332430f14146102625780634a96aaa014610283578063715018a6146102965780637c6ee6ab1461029e5780638687feae146102b157600080fd5b806239859914610130578063048886d2146101455780631429c7c21461016d5780631c3970fa146101925780632ecfe72b146101c2575b600080fd5b61014361013e366004610d52565b610400565b005b610158610153366004610da0565b61041f565b60405190151581526020015b60405180910390f35b61018061017b366004610da0565b6104c9565b60405160ff9091168152602001610164565b6005546101a89060ff8082169161010090041682565b6040805160ff938416815292909116602083015201610164565b6102326101d0366004610dc2565b60408051606080820183526000808352602080840182905292840181905261ffff9490941684526004825292829020825193840183525463ffffffff808216855264010000000082041691840191909152600160401b900460ff169082015290565b60408051825163ffffffff9081168252602080850151909116908201529181015160ff1690820152606001610164565b6003546102709061ffff1681565b60405161ffff9091168152602001610164565b610143610291366004610d52565b610537565b610143610552565b6101436102ac366004610e44565b610566565b6102b9610598565b6040516101649190610e60565b6102706102d4366004610f38565b610626565b6065546040516001600160a01b039091168152602001610164565b610143610302366004610d52565b61063f565b6102b961065a565b6102b9610667565b610180610325366004610da0565b610674565b60408051808201825260008082526020918201528151808301835260055460ff8082168084526101009092048116928401928352845191825291519091169181019190915201610164565b610143610383366004610f6b565b6106a0565b6103c7610396366004610dc2565b60046020526000908152604090205463ffffffff80821691640100000000810490911690600160401b900460ff1683565b6040805163ffffffff948516815293909216602084015260ff1690820152606001610164565b6101436103fb366004611010565b61071e565b6104086108d9565b805161041b906001906020840190610c02565b5050565b600080600160ff84161b9050806104bf6002805461043c906110de565b80601f0160208091040260200160405190810160405280929190818152602001828054610468906110de565b80156104b55780601f1061048a576101008083540402835291602001916104b5565b820191906000526020600020905b81548152906001019060200180831161049857829003601f168201915b5050505050610933565b9091161492915050565b60008160ff16600180546104dc906110de565b905011156105325760018260ff1681546104f5906110de565b811061050357610503611119565b8154600116156105225790600052602060002090602091828204019190065b9054901a600160f81b0260f81c90505b919050565b61053f6108d9565b805161041b906000906020840190610c02565b61055a6108d9565b6105646000610ac0565b565b61056e6108d9565b80516005805460209093015160ff9081166101000261ffff19909416921691909117919091179055565b600080546105a5906110de565b80601f01602080910402602001604051908101604052809291908181526020018280546105d1906110de565b801561061e5780601f106105f35761010080835404028352916020019161061e565b820191906000526020600020905b81548152906001019060200180831161060157829003601f168201915b505050505081565b60006106306108d9565b61063982610b12565b92915050565b6106476108d9565b805161041b906002906020840190610c02565b600180546105a5906110de565b600280546105a5906110de565b60008160ff1660008054610687906110de565b905011156105325760008260ff1681546104f5906110de565b6106a86108d9565b6001600160a01b0381166107125760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b60648201526084015b60405180910390fd5b61071b81610ac0565b50565b603254610100900460ff161580801561073e5750603254600160ff909116105b806107585750303b158015610758575060325460ff166001145b6107bb5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b6064820152608401610709565b6032805460ff1916600117905580156107de576032805461ff0019166101001790555b6107e787610ac0565b85516107fa906000906020890190610c02565b50845161080e906001906020880190610c02565b508351610822906002906020870190610c02565b50815160058054602085015160ff9081166101000261ffff1990921693169290921791909117905560005b83518110156108895761087884828151811061086b5761086b611119565b6020026020010151610b12565b5061088281611145565b905061084d565b5080156108d0576032805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50505050505050565b6065546001600160a01b031633146105645760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610709565b6000610100825111156109bc5760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a401610709565b81516109ca57506000919050565b600080836000815181106109e0576109e0611119565b0160200151600160f89190911c81901b92505b8451811015610ab757848181518110610a0e57610a0e611119565b0160200151600160f89190911c1b9150828211610aa35760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a401610709565b91811791610ab081611145565b90506109f3565b50909392505050565b606580546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b6003805461ffff90811660009081526004602090815260408083208651815488850180518a8601805163ffffffff95861667ffffffffffffffff199095168517640100000000938716939093029290921768ff00000000000000001916600160401b60ff9384160217909555985485519283529051909216948101949094529051909516908201529092909116907fdbee9d337a6e5fde30966e157673aaeeb6a0134afaf774a4b6979b7c79d07da49060600160405180910390a26003805461ffff16906000610be183611160565b91906101000a81548161ffff021916908361ffff1602179055509050919050565b828054610c0e906110de565b90600052602060002090601f016020900481019282610c305760008555610c76565b82601f10610c4957805160ff1916838001178555610c76565b82800160010185558215610c76579182015b82811115610c76578251825591602001919060010190610c5b565b50610c82929150610c86565b5090565b5b80821115610c825760008155600101610c87565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff81118282101715610cda57610cda610c9b565b604052919050565b600082601f830112610cf357600080fd5b813567ffffffffffffffff811115610d0d57610d0d610c9b565b610d20601f8201601f1916602001610cb1565b818152846020838601011115610d3557600080fd5b816020850160208301376000918101602001919091529392505050565b600060208284031215610d6457600080fd5b813567ffffffffffffffff811115610d7b57600080fd5b610d8784828501610ce2565b949350505050565b803560ff8116811461053257600080fd5b600060208284031215610db257600080fd5b610dbb82610d8f565b9392505050565b600060208284031215610dd457600080fd5b813561ffff81168114610dbb57600080fd5b600060408284031215610df857600080fd5b6040516040810181811067ffffffffffffffff82111715610e1b57610e1b610c9b565b604052905080610e2a83610d8f565b8152610e3860208401610d8f565b60208201525092915050565b600060408284031215610e5657600080fd5b610dbb8383610de6565b600060208083528351808285015260005b81811015610e8d57858101830151858201604001528201610e71565b81811115610e9f576000604083870101525b50601f01601f1916929092016040019392505050565b803563ffffffff8116811461053257600080fd5b600060608284031215610edb57600080fd5b6040516060810181811067ffffffffffffffff82111715610efe57610efe610c9b565b604052905080610f0d83610eb5565b8152610f1b60208401610eb5565b6020820152610f2c60408401610d8f565b60408201525092915050565b600060608284031215610f4a57600080fd5b610dbb8383610ec9565b80356001600160a01b038116811461053257600080fd5b600060208284031215610f7d57600080fd5b610dbb82610f54565b600082601f830112610f9757600080fd5b8135602067ffffffffffffffff821115610fb357610fb3610c9b565b610fc1818360051b01610cb1565b82815260609283028501820192828201919087851115610fe057600080fd5b8387015b8581101561100357610ff68982610ec9565b8452928401928101610fe4565b5090979650505050505050565b60008060008060008060e0878903121561102957600080fd5b61103287610f54565b9550602087013567ffffffffffffffff8082111561104f57600080fd5b61105b8a838b01610ce2565b9650604089013591508082111561107157600080fd5b61107d8a838b01610ce2565b9550606089013591508082111561109357600080fd5b61109f8a838b01610ce2565b945060808901359150808211156110b557600080fd5b506110c289828a01610f86565b9250506110d28860a08901610de6565b90509295509295509295565b600181811c908216806110f257607f821691505b6020821081141561111357634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b60006000198214156111595761115961112f565b5060010190565b600061ffff808316818114156111785761117861112f565b600101939250505056fea264697066735822122089ddb4d1d016b99e8d0b49592d60555fec0f671de64b7bb17f2afc4c1e7aeee064736f6c634300080c0033",
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

// DefaultSecurityThresholdsV2 is a free data retrieval call binding the contract method 0x1c3970fa.
//
// Solidity: function defaultSecurityThresholdsV2() view returns(uint8 confirmationThreshold, uint8 adversaryThreshold)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCaller) DefaultSecurityThresholdsV2(opts *bind.CallOpts) (struct {
	ConfirmationThreshold uint8
	AdversaryThreshold    uint8
}, error) {
	var out []interface{}
	err := _ContractEigenDAThresholdRegistry.contract.Call(opts, &out, "defaultSecurityThresholdsV2")

	outstruct := new(struct {
		ConfirmationThreshold uint8
		AdversaryThreshold    uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfirmationThreshold = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.AdversaryThreshold = *abi.ConvertType(out[1], new(uint8)).(*uint8)

	return *outstruct, err

}

// DefaultSecurityThresholdsV2 is a free data retrieval call binding the contract method 0x1c3970fa.
//
// Solidity: function defaultSecurityThresholdsV2() view returns(uint8 confirmationThreshold, uint8 adversaryThreshold)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) DefaultSecurityThresholdsV2() (struct {
	ConfirmationThreshold uint8
	AdversaryThreshold    uint8
}, error) {
	return _ContractEigenDAThresholdRegistry.Contract.DefaultSecurityThresholdsV2(&_ContractEigenDAThresholdRegistry.CallOpts)
}

// DefaultSecurityThresholdsV2 is a free data retrieval call binding the contract method 0x1c3970fa.
//
// Solidity: function defaultSecurityThresholdsV2() view returns(uint8 confirmationThreshold, uint8 adversaryThreshold)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCallerSession) DefaultSecurityThresholdsV2() (struct {
	ConfirmationThreshold uint8
	AdversaryThreshold    uint8
}, error) {
	return _ContractEigenDAThresholdRegistry.Contract.DefaultSecurityThresholdsV2(&_ContractEigenDAThresholdRegistry.CallOpts)
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

// GetDefaultSecurityThresholdsV2 is a free data retrieval call binding the contract method 0xef635529.
//
// Solidity: function getDefaultSecurityThresholdsV2() view returns((uint8,uint8))
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCaller) GetDefaultSecurityThresholdsV2(opts *bind.CallOpts) (SecurityThresholds, error) {
	var out []interface{}
	err := _ContractEigenDAThresholdRegistry.contract.Call(opts, &out, "getDefaultSecurityThresholdsV2")

	if err != nil {
		return *new(SecurityThresholds), err
	}

	out0 := *abi.ConvertType(out[0], new(SecurityThresholds)).(*SecurityThresholds)

	return out0, err

}

// GetDefaultSecurityThresholdsV2 is a free data retrieval call binding the contract method 0xef635529.
//
// Solidity: function getDefaultSecurityThresholdsV2() view returns((uint8,uint8))
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) GetDefaultSecurityThresholdsV2() (SecurityThresholds, error) {
	return _ContractEigenDAThresholdRegistry.Contract.GetDefaultSecurityThresholdsV2(&_ContractEigenDAThresholdRegistry.CallOpts)
}

// GetDefaultSecurityThresholdsV2 is a free data retrieval call binding the contract method 0xef635529.
//
// Solidity: function getDefaultSecurityThresholdsV2() view returns((uint8,uint8))
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCallerSession) GetDefaultSecurityThresholdsV2() (SecurityThresholds, error) {
	return _ContractEigenDAThresholdRegistry.Contract.GetDefaultSecurityThresholdsV2(&_ContractEigenDAThresholdRegistry.CallOpts)
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

// Initialize is a paid mutator transaction binding the contract method 0xfb87355e.
//
// Solidity: function initialize(address _initialOwner, bytes _quorumAdversaryThresholdPercentages, bytes _quorumConfirmationThresholdPercentages, bytes _quorumNumbersRequired, (uint32,uint32,uint8)[] _versionedBlobParams, (uint8,uint8) _defaultSecurityThresholdsV2) returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactor) Initialize(opts *bind.TransactOpts, _initialOwner common.Address, _quorumAdversaryThresholdPercentages []byte, _quorumConfirmationThresholdPercentages []byte, _quorumNumbersRequired []byte, _versionedBlobParams []VersionedBlobParams, _defaultSecurityThresholdsV2 SecurityThresholds) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.contract.Transact(opts, "initialize", _initialOwner, _quorumAdversaryThresholdPercentages, _quorumConfirmationThresholdPercentages, _quorumNumbersRequired, _versionedBlobParams, _defaultSecurityThresholdsV2)
}

// Initialize is a paid mutator transaction binding the contract method 0xfb87355e.
//
// Solidity: function initialize(address _initialOwner, bytes _quorumAdversaryThresholdPercentages, bytes _quorumConfirmationThresholdPercentages, bytes _quorumNumbersRequired, (uint32,uint32,uint8)[] _versionedBlobParams, (uint8,uint8) _defaultSecurityThresholdsV2) returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) Initialize(_initialOwner common.Address, _quorumAdversaryThresholdPercentages []byte, _quorumConfirmationThresholdPercentages []byte, _quorumNumbersRequired []byte, _versionedBlobParams []VersionedBlobParams, _defaultSecurityThresholdsV2 SecurityThresholds) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.Initialize(&_ContractEigenDAThresholdRegistry.TransactOpts, _initialOwner, _quorumAdversaryThresholdPercentages, _quorumConfirmationThresholdPercentages, _quorumNumbersRequired, _versionedBlobParams, _defaultSecurityThresholdsV2)
}

// Initialize is a paid mutator transaction binding the contract method 0xfb87355e.
//
// Solidity: function initialize(address _initialOwner, bytes _quorumAdversaryThresholdPercentages, bytes _quorumConfirmationThresholdPercentages, bytes _quorumNumbersRequired, (uint32,uint32,uint8)[] _versionedBlobParams, (uint8,uint8) _defaultSecurityThresholdsV2) returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactorSession) Initialize(_initialOwner common.Address, _quorumAdversaryThresholdPercentages []byte, _quorumConfirmationThresholdPercentages []byte, _quorumNumbersRequired []byte, _versionedBlobParams []VersionedBlobParams, _defaultSecurityThresholdsV2 SecurityThresholds) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.Initialize(&_ContractEigenDAThresholdRegistry.TransactOpts, _initialOwner, _quorumAdversaryThresholdPercentages, _quorumConfirmationThresholdPercentages, _quorumNumbersRequired, _versionedBlobParams, _defaultSecurityThresholdsV2)
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

// UpdateDefaultSecurityThresholdsV2 is a paid mutator transaction binding the contract method 0x7c6ee6ab.
//
// Solidity: function updateDefaultSecurityThresholdsV2((uint8,uint8) _defaultSecurityThresholdsV2) returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactor) UpdateDefaultSecurityThresholdsV2(opts *bind.TransactOpts, _defaultSecurityThresholdsV2 SecurityThresholds) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.contract.Transact(opts, "updateDefaultSecurityThresholdsV2", _defaultSecurityThresholdsV2)
}

// UpdateDefaultSecurityThresholdsV2 is a paid mutator transaction binding the contract method 0x7c6ee6ab.
//
// Solidity: function updateDefaultSecurityThresholdsV2((uint8,uint8) _defaultSecurityThresholdsV2) returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) UpdateDefaultSecurityThresholdsV2(_defaultSecurityThresholdsV2 SecurityThresholds) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.UpdateDefaultSecurityThresholdsV2(&_ContractEigenDAThresholdRegistry.TransactOpts, _defaultSecurityThresholdsV2)
}

// UpdateDefaultSecurityThresholdsV2 is a paid mutator transaction binding the contract method 0x7c6ee6ab.
//
// Solidity: function updateDefaultSecurityThresholdsV2((uint8,uint8) _defaultSecurityThresholdsV2) returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactorSession) UpdateDefaultSecurityThresholdsV2(_defaultSecurityThresholdsV2 SecurityThresholds) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.UpdateDefaultSecurityThresholdsV2(&_ContractEigenDAThresholdRegistry.TransactOpts, _defaultSecurityThresholdsV2)
}

// UpdateQuorumAdversaryThresholdPercentages is a paid mutator transaction binding the contract method 0x4a96aaa0.
//
// Solidity: function updateQuorumAdversaryThresholdPercentages(bytes _quorumAdversaryThresholdPercentages) returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactor) UpdateQuorumAdversaryThresholdPercentages(opts *bind.TransactOpts, _quorumAdversaryThresholdPercentages []byte) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.contract.Transact(opts, "updateQuorumAdversaryThresholdPercentages", _quorumAdversaryThresholdPercentages)
}

// UpdateQuorumAdversaryThresholdPercentages is a paid mutator transaction binding the contract method 0x4a96aaa0.
//
// Solidity: function updateQuorumAdversaryThresholdPercentages(bytes _quorumAdversaryThresholdPercentages) returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) UpdateQuorumAdversaryThresholdPercentages(_quorumAdversaryThresholdPercentages []byte) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.UpdateQuorumAdversaryThresholdPercentages(&_ContractEigenDAThresholdRegistry.TransactOpts, _quorumAdversaryThresholdPercentages)
}

// UpdateQuorumAdversaryThresholdPercentages is a paid mutator transaction binding the contract method 0x4a96aaa0.
//
// Solidity: function updateQuorumAdversaryThresholdPercentages(bytes _quorumAdversaryThresholdPercentages) returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactorSession) UpdateQuorumAdversaryThresholdPercentages(_quorumAdversaryThresholdPercentages []byte) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.UpdateQuorumAdversaryThresholdPercentages(&_ContractEigenDAThresholdRegistry.TransactOpts, _quorumAdversaryThresholdPercentages)
}

// UpdateQuorumConfirmationThresholdPercentages is a paid mutator transaction binding the contract method 0x00398599.
//
// Solidity: function updateQuorumConfirmationThresholdPercentages(bytes _quorumConfirmationThresholdPercentages) returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactor) UpdateQuorumConfirmationThresholdPercentages(opts *bind.TransactOpts, _quorumConfirmationThresholdPercentages []byte) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.contract.Transact(opts, "updateQuorumConfirmationThresholdPercentages", _quorumConfirmationThresholdPercentages)
}

// UpdateQuorumConfirmationThresholdPercentages is a paid mutator transaction binding the contract method 0x00398599.
//
// Solidity: function updateQuorumConfirmationThresholdPercentages(bytes _quorumConfirmationThresholdPercentages) returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) UpdateQuorumConfirmationThresholdPercentages(_quorumConfirmationThresholdPercentages []byte) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.UpdateQuorumConfirmationThresholdPercentages(&_ContractEigenDAThresholdRegistry.TransactOpts, _quorumConfirmationThresholdPercentages)
}

// UpdateQuorumConfirmationThresholdPercentages is a paid mutator transaction binding the contract method 0x00398599.
//
// Solidity: function updateQuorumConfirmationThresholdPercentages(bytes _quorumConfirmationThresholdPercentages) returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactorSession) UpdateQuorumConfirmationThresholdPercentages(_quorumConfirmationThresholdPercentages []byte) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.UpdateQuorumConfirmationThresholdPercentages(&_ContractEigenDAThresholdRegistry.TransactOpts, _quorumConfirmationThresholdPercentages)
}

// UpdateQuorumNumbersRequired is a paid mutator transaction binding the contract method 0xa5e9b2eb.
//
// Solidity: function updateQuorumNumbersRequired(bytes _quorumNumbersRequired) returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactor) UpdateQuorumNumbersRequired(opts *bind.TransactOpts, _quorumNumbersRequired []byte) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.contract.Transact(opts, "updateQuorumNumbersRequired", _quorumNumbersRequired)
}

// UpdateQuorumNumbersRequired is a paid mutator transaction binding the contract method 0xa5e9b2eb.
//
// Solidity: function updateQuorumNumbersRequired(bytes _quorumNumbersRequired) returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) UpdateQuorumNumbersRequired(_quorumNumbersRequired []byte) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.UpdateQuorumNumbersRequired(&_ContractEigenDAThresholdRegistry.TransactOpts, _quorumNumbersRequired)
}

// UpdateQuorumNumbersRequired is a paid mutator transaction binding the contract method 0xa5e9b2eb.
//
// Solidity: function updateQuorumNumbersRequired(bytes _quorumNumbersRequired) returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactorSession) UpdateQuorumNumbersRequired(_quorumNumbersRequired []byte) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.UpdateQuorumNumbersRequired(&_ContractEigenDAThresholdRegistry.TransactOpts, _quorumNumbersRequired)
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
