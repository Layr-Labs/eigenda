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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addVersionedBlobParams\",\"inputs\":[{\"name\":\"_versionedBlobParams\",\"type\":\"tuple\",\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"defaultSecurityThresholdsV2\",\"inputs\":[],\"outputs\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBlobParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDefaultSecurityThresholdsV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIsQuorumRequired\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumAdversaryThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumConfirmationThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_quorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_quorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_quorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_versionedBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structVersionedBlobParams[]\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"_defaultSecurityThresholdsV2\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"nextBlobVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumAdversaryThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumConfirmationThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateDefaultSecurityThresholdsV2\",\"inputs\":[{\"name\":\"_defaultSecurityThresholdsV2\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateQuorumAdversaryThresholdPercentages\",\"inputs\":[{\"name\":\"_quorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateQuorumConfirmationThresholdPercentages\",\"inputs\":[{\"name\":\"_quorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateQuorumNumbersRequired\",\"inputs\":[{\"name\":\"_quorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"versionedBlobParams\",\"inputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"DefaultSecurityThresholdsV2Updated\",\"inputs\":[{\"name\":\"previousDefaultSecurityThresholdsV2\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"newDefaultSecurityThresholdsV2\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumAdversaryThresholdPercentagesUpdated\",\"inputs\":[{\"name\":\"previousQuorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumConfirmationThresholdPercentagesUpdated\",\"inputs\":[{\"name\":\"previousQuorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumNumbersRequiredUpdated\",\"inputs\":[{\"name\":\"previousQuorumNumbersRequired\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumNumbersRequired\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VersionedBlobParamsAdded\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"indexed\":true,\"internalType\":\"uint16\"},{\"name\":\"versionedBlobParams\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"anonymous\":false}]",
	Bin: "0x608060405234801561001057600080fd5b5061001961001e565b6100de565b603254610100900460ff161561008a5760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840160405180910390fd5b60325460ff90811610156100dc576032805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b6113ae806100ed6000396000f3fe608060405234801561001057600080fd5b506004361061012b5760003560e01c80638a476982116100ad578063ee6c3bcf11610071578063ee6c3bcf14610317578063ef6355291461032a578063f2fde38b14610368578063f74e363c1461037b578063fb87355e146103e057600080fd5b80638a476982146102c65780638da5cb5b146102d9578063a5e9b2eb146102f4578063bafa910714610307578063e15234ff1461030f57600080fd5b806332430f14116100f457806332430f14146102625780634a96aaa014610283578063715018a6146102965780637c6ee6ab1461029e5780638687feae146102b157600080fd5b806239859914610130578063048886d2146101455780631429c7c21461016d5780631c3970fa146101925780632ecfe72b146101c2575b600080fd5b61014361013e366004610e2d565b6103f3565b005b610158610153366004610e7b565b61044c565b60405190151581526020015b60405180910390f35b61018061017b366004610e7b565b6104f6565b60405160ff9091168152602001610164565b6005546101a89060ff8082169161010090041682565b6040805160ff938416815292909116602083015201610164565b6102326101d0366004610e9d565b60408051606080820183526000808352602080840182905292840181905261ffff9490941684526004825292829020825193840183525463ffffffff808216855264010000000082041691840191909152600160401b900460ff169082015290565b60408051825163ffffffff9081168252602080850151909116908201529181015160ff1690820152606001610164565b6003546102709061ffff1681565b60405161ffff9091168152602001610164565b610143610291366004610e2d565b610564565b6101436105b9565b6101436102ac366004610f1f565b6105cd565b6102b9610639565b6040516101649190610f88565b6102706102d436600461101e565b6106c7565b6065546040516001600160a01b039091168152602001610164565b610143610302366004610e2d565b6106e0565b6102b9610735565b6102b9610742565b610180610325366004610e7b565b61074f565b6040805180820182526000808252602091820152815180830190925260055460ff80821684526101009091041690820152604051610164919061103a565b61014361037636600461106f565b61077b565b6103ba610389366004610e9d565b60046020526000908152604090205463ffffffff80821691640100000000810490911690600160401b900460ff1683565b6040805163ffffffff948516815293909216602084015260ff1690820152606001610164565b6101436103ee366004611114565b6107f9565b6103fb6109b4565b7f9f1ea99a8363f2964c53c763811648354a8437441b30b39465f9d26118d6a5a060018260405161042d92919061121d565b60405180910390a18051610448906001906020840190610cdd565b5050565b600080600160ff84161b9050806104ec60028054610469906111e2565b80601f0160208091040260200160405190810160405280929190818152602001828054610495906111e2565b80156104e25780601f106104b7576101008083540402835291602001916104e2565b820191906000526020600020905b8154815290600101906020018083116104c557829003601f168201915b5050505050610a0e565b9091161492915050565b60008160ff1660018054610509906111e2565b9050111561055f5760018260ff168154610522906111e2565b8110610530576105306112d9565b81546001161561054f5790600052602060002090602091828204019190065b9054901a600160f81b0260f81c90505b919050565b61056c6109b4565b7ff73542111561dc551cbbe9111c4dd3a040d53d7bc0339a53290f4d7f9a95c3cc60008260405161059e92919061121d565b60405180910390a18051610448906000906020840190610cdd565b6105c16109b4565b6105cb6000610b9b565b565b6105d56109b4565b7ffe03afd62c76a6aed7376ae995cc55d073ba9d83d83ac8efc5446f8da4d509976005826040516106079291906112ef565b60405180910390a180516005805460209093015160ff9081166101000261ffff19909416921691909117919091179055565b60008054610646906111e2565b80601f0160208091040260200160405190810160405280929190818152602001828054610672906111e2565b80156106bf5780601f10610694576101008083540402835291602001916106bf565b820191906000526020600020905b8154815290600101906020018083116106a257829003601f168201915b505050505081565b60006106d16109b4565b6106da82610bed565b92915050565b6106e86109b4565b7f60c0ba1da794fcbbf549d370512442cb8f3f3f774cb557205cc88c6f842cb36a60028260405161071a92919061121d565b60405180910390a18051610448906002906020840190610cdd565b60018054610646906111e2565b60028054610646906111e2565b60008160ff1660008054610762906111e2565b9050111561055f5760008260ff168154610522906111e2565b6107836109b4565b6001600160a01b0381166107ed5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b60648201526084015b60405180910390fd5b6107f681610b9b565b50565b603254610100900460ff16158080156108195750603254600160ff909116105b806108335750303b158015610833575060325460ff166001145b6108965760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b60648201526084016107e4565b6032805460ff1916600117905580156108b9576032805461ff0019166101001790555b6108c287610b9b565b85516108d5906000906020890190610cdd565b5084516108e9906001906020880190610cdd565b5083516108fd906002906020870190610cdd565b50815160058054602085015160ff9081166101000261ffff1990921693169290921791909117905560005b835181101561096457610953848281518110610946576109466112d9565b6020026020010151610bed565b5061095d8161133b565b9050610928565b5080156109ab576032805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50505050505050565b6065546001600160a01b031633146105cb5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016107e4565b600061010082511115610a975760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a4016107e4565b8151610aa557506000919050565b60008083600081518110610abb57610abb6112d9565b0160200151600160f89190911c81901b92505b8451811015610b9257848181518110610ae957610ae96112d9565b0160200151600160f89190911c1b9150828211610b7e5760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a4016107e4565b91811791610b8b8161133b565b9050610ace565b50909392505050565b606580546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b6003805461ffff90811660009081526004602090815260408083208651815488850180518a8601805163ffffffff95861667ffffffffffffffff199095168517640100000000938716939093029290921768ff00000000000000001916600160401b60ff9384160217909555985485519283529051909216948101949094529051909516908201529092909116907fdbee9d337a6e5fde30966e157673aaeeb6a0134afaf774a4b6979b7c79d07da49060600160405180910390a26003805461ffff16906000610cbc83611356565b91906101000a81548161ffff021916908361ffff1602179055509050919050565b828054610ce9906111e2565b90600052602060002090601f016020900481019282610d0b5760008555610d51565b82601f10610d2457805160ff1916838001178555610d51565b82800160010185558215610d51579182015b82811115610d51578251825591602001919060010190610d36565b50610d5d929150610d61565b5090565b5b80821115610d5d5760008155600101610d62565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff81118282101715610db557610db5610d76565b604052919050565b600082601f830112610dce57600080fd5b813567ffffffffffffffff811115610de857610de8610d76565b610dfb601f8201601f1916602001610d8c565b818152846020838601011115610e1057600080fd5b816020850160208301376000918101602001919091529392505050565b600060208284031215610e3f57600080fd5b813567ffffffffffffffff811115610e5657600080fd5b610e6284828501610dbd565b949350505050565b803560ff8116811461055f57600080fd5b600060208284031215610e8d57600080fd5b610e9682610e6a565b9392505050565b600060208284031215610eaf57600080fd5b813561ffff81168114610e9657600080fd5b600060408284031215610ed357600080fd5b6040516040810181811067ffffffffffffffff82111715610ef657610ef6610d76565b604052905080610f0583610e6a565b8152610f1360208401610e6a565b60208201525092915050565b600060408284031215610f3157600080fd5b610e968383610ec1565b6000815180845260005b81811015610f6157602081850181015186830182015201610f45565b81811115610f73576000602083870101525b50601f01601f19169290920160200192915050565b602081526000610e966020830184610f3b565b803563ffffffff8116811461055f57600080fd5b600060608284031215610fc157600080fd5b6040516060810181811067ffffffffffffffff82111715610fe457610fe4610d76565b604052905080610ff383610f9b565b815261100160208401610f9b565b602082015261101260408401610e6a565b60408201525092915050565b60006060828403121561103057600080fd5b610e968383610faf565b604081016106da8284805160ff908116835260209182015116910152565b80356001600160a01b038116811461055f57600080fd5b60006020828403121561108157600080fd5b610e9682611058565b600082601f83011261109b57600080fd5b8135602067ffffffffffffffff8211156110b7576110b7610d76565b6110c5818360051b01610d8c565b828152606092830285018201928282019190878511156110e457600080fd5b8387015b85811015611107576110fa8982610faf565b84529284019281016110e8565b5090979650505050505050565b60008060008060008060e0878903121561112d57600080fd5b61113687611058565b9550602087013567ffffffffffffffff8082111561115357600080fd5b61115f8a838b01610dbd565b9650604089013591508082111561117557600080fd5b6111818a838b01610dbd565b9550606089013591508082111561119757600080fd5b6111a38a838b01610dbd565b945060808901359150808211156111b957600080fd5b506111c689828a0161108a565b9250506111d68860a08901610ec1565b90509295509295509295565b600181811c908216806111f657607f821691505b6020821081141561121757634e487b7160e01b600052602260045260246000fd5b50919050565b60408152600080845481600182811c91508083168061123d57607f831692505b602080841082141561125d57634e487b7160e01b86526022600452602486fd5b604088018490526060880182801561127c576001811461128d576112b8565b60ff198716825282820197506112b8565b60008c81526020902060005b878110156112b257815484820152908601908401611299565b83019850505b50508786038189015250505050506112d08185610f3b565b95945050505050565b634e487b7160e01b600052603260045260246000fd5b825460ff808216835260089190911c16602082015260808101610e966040830184805160ff908116835260209182015116910152565b634e487b7160e01b600052601160045260246000fd5b600060001982141561134f5761134f611325565b5060010190565b600061ffff8083168181141561136e5761136e611325565b600101939250505056fea26469706673582212208e4615301ccfefbb6997f01789117492a27a9f8b4f3611a98c8d3ecbafa3d14864736f6c634300080c0033",
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
