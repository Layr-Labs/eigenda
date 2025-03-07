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
	Bin: "0x60806040523480156200001157600080fd5b50620000226200002860201b60201c565b620001d6565b603260019054906101000a900460ff16156200007b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401620000729062000179565b60405180910390fd5b60ff8016603260009054906101000a900460ff1660ff161015620000f05760ff603260006101000a81548160ff021916908360ff1602179055507f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb384740249860ff604051620000e79190620001b9565b60405180910390a15b565b600082825260208201905092915050565b7f496e697469616c697a61626c653a20636f6e747261637420697320696e69746960008201527f616c697a696e6700000000000000000000000000000000000000000000000000602082015250565b600062000161602783620000f2565b91506200016e8262000103565b604082019050919050565b60006020820190508181036000830152620001948162000152565b9050919050565b600060ff82169050919050565b620001b3816200019b565b82525050565b6000602082019050620001d06000830184620001a8565b92915050565b611af380620001e66000396000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c80638a4769821161008c578063e15234ff11610066578063e15234ff1461024d578063ee6c3bcf1461026b578063f2fde38b1461029b578063f74e363c146102b7576100ea565b80638a476982146101e15780638da5cb5b14610211578063bafa91071461022f576100ea565b806332430f14116100c857806332430f141461017f578063715018a61461019d5780638491bad6146101a75780638687feae146101c3576100ea565b8063048886d2146100ef5780631429c7c21461011f5780632ecfe72b1461014f575b600080fd5b61010960048036038101906101049190610f06565b6102e9565b6040516101169190610f4e565b60405180910390f35b61013960048036038101906101349190610f06565b610398565b6040516101469190610f78565b60405180910390f35b61016960048036038101906101649190610fcd565b610423565b604051610176919061106a565b60405180910390f35b6101876104c1565b6040516101949190611094565b60405180910390f35b6101a56104d5565b005b6101c160048036038101906101bc91906113b0565b6104e9565b005b6101cb6106bb565b6040516101d89190611523565b60405180910390f35b6101fb60048036038101906101f69190611545565b610749565b6040516102089190611094565b60405180910390f35b610219610763565b6040516102269190611581565b60405180910390f35b61023761078d565b6040516102449190611523565b60405180910390f35b61025561081b565b6040516102629190611523565b60405180910390f35b61028560048036038101906102809190610f06565b6108a9565b6040516102929190610f78565b60405180910390f35b6102b560048036038101906102b0919061159c565b610934565b005b6102d160048036038101906102cc9190610fcd565b6109b8565b6040516102e0939291906115d8565b60405180910390f35b6000806102f7600084610a0f565b90508061038d6002805461030a9061163e565b80601f01602080910402602001604051908101604052809291908181526020018280546103369061163e565b80156103835780601f1061035857610100808354040283529160200191610383565b820191906000526020600020905b81548152906001019060200180831161036657829003601f168201915b5050505050610a23565b821614915050919050565b60008160ff16600180546103ab9061163e565b9050111561041e5760018260ff1681546103c49061163e565b81106103d3576103d2611670565b5b8154600116156103f25790600052602060002090602091828204019190065b9054901a7f01000000000000000000000000000000000000000000000000000000000000000260f81c90505b919050565b61042b610de6565b600460008361ffff1661ffff1681526020019081526020016000206040518060600160405290816000820160009054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160049054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160089054906101000a900460ff1660ff1660ff16815250509050919050565b600360009054906101000a900461ffff1681565b6104dd610b4a565b6104e76000610bc8565b565b6000603260019054906101000a900460ff1615905080801561051d57506001603260009054906101000a900460ff1660ff16105b8061054c575061052c30610c8e565b15801561054b57506001603260009054906101000a900460ff1660ff16145b5b61058b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161058290611722565b60405180910390fd5b6001603260006101000a81548160ff021916908360ff16021790555080156105c9576001603260016101000a81548160ff0219169083151502179055505b6105d286610bc8565b84600090805190602001906105e8929190610e16565b5083600190805190602001906105ff929190610e16565b508260029080519060200190610616929190610e16565b5060005b82518110156106585761064683828151811061063957610638611670565b5b6020026020010151610cb1565b50806106519061177b565b905061061a565b5080156106b3576000603260016101000a81548160ff0219169083151502179055507f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb384740249860016040516106aa9190611809565b60405180910390a15b505050505050565b600080546106c89061163e565b80601f01602080910402602001604051908101604052809291908181526020018280546106f49061163e565b80156107415780601f1061071657610100808354040283529160200191610741565b820191906000526020600020905b81548152906001019060200180831161072457829003601f168201915b505050505081565b6000610753610b4a565b61075c82610cb1565b9050919050565b6000606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6001805461079a9061163e565b80601f01602080910402602001604051908101604052809291908181526020018280546107c69061163e565b80156108135780601f106107e857610100808354040283529160200191610813565b820191906000526020600020905b8154815290600101906020018083116107f657829003601f168201915b505050505081565b600280546108289061163e565b80601f01602080910402602001604051908101604052809291908181526020018280546108549061163e565b80156108a15780601f10610876576101008083540402835291602001916108a1565b820191906000526020600020905b81548152906001019060200180831161088457829003601f168201915b505050505081565b60008160ff16600080546108bc9061163e565b9050111561092f5760008260ff1681546108d59061163e565b81106108e4576108e3611670565b5b8154600116156109035790600052602060002090602091828204019190065b9054901a7f01000000000000000000000000000000000000000000000000000000000000000260f81c90505b919050565b61093c610b4a565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614156109ac576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109a390611896565b60405180910390fd5b6109b581610bc8565b50565b60046020528060005260406000206000915090508060000160009054906101000a900463ffffffff16908060000160049054906101000a900463ffffffff16908060000160089054906101000a900460ff16905083565b60008160ff166001901b8317905092915050565b600061010082511115610a6b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a629061194e565b60405180910390fd5b600082511415610a7e5760009050610b45565b60008083600081518110610a9557610a94611670565b5b602001015160f81c60f81b60f81c60ff166001901b91506000600190505b8451811015610b3e57848181518110610acf57610ace611670565b5b602001015160f81c60f81b60f81c60ff166001901b9150828211610b28576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b1f90611a06565b60405180910390fd5b818317925080610b379061177b565b9050610ab3565b5081925050505b919050565b610b52610dde565b73ffffffffffffffffffffffffffffffffffffffff16610b70610763565b73ffffffffffffffffffffffffffffffffffffffff1614610bc6576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610bbd90611a72565b60405180910390fd5b565b6000606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905081606560006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b6000808273ffffffffffffffffffffffffffffffffffffffff163b119050919050565b60008160046000600360009054906101000a900461ffff1661ffff1661ffff16815260200190815260200160002060008201518160000160006101000a81548163ffffffff021916908363ffffffff16021790555060208201518160000160046101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160086101000a81548160ff021916908360ff160217905550905050600360009054906101000a900461ffff1661ffff167fdbee9d337a6e5fde30966e157673aaeeb6a0134afaf774a4b6979b7c79d07da483604051610d95919061106a565b60405180910390a26003600081819054906101000a900461ffff1680929190610dbd90611a92565b91906101000a81548161ffff021916908361ffff1602179055509050919050565b600033905090565b6040518060600160405280600063ffffffff168152602001600063ffffffff168152602001600060ff1681525090565b828054610e229061163e565b90600052602060002090601f016020900481019282610e445760008555610e8b565b82601f10610e5d57805160ff1916838001178555610e8b565b82800160010185558215610e8b579182015b82811115610e8a578251825591602001919060010190610e6f565b5b509050610e989190610e9c565b5090565b5b80821115610eb5576000816000905550600101610e9d565b5090565b6000604051905090565b600080fd5b600080fd5b600060ff82169050919050565b610ee381610ecd565b8114610eee57600080fd5b50565b600081359050610f0081610eda565b92915050565b600060208284031215610f1c57610f1b610ec3565b5b6000610f2a84828501610ef1565b91505092915050565b60008115159050919050565b610f4881610f33565b82525050565b6000602082019050610f636000830184610f3f565b92915050565b610f7281610ecd565b82525050565b6000602082019050610f8d6000830184610f69565b92915050565b600061ffff82169050919050565b610faa81610f93565b8114610fb557600080fd5b50565b600081359050610fc781610fa1565b92915050565b600060208284031215610fe357610fe2610ec3565b5b6000610ff184828501610fb8565b91505092915050565b600063ffffffff82169050919050565b61101381610ffa565b82525050565b61102281610ecd565b82525050565b60608201600082015161103e600085018261100a565b506020820151611051602085018261100a565b5060408201516110646040850182611019565b50505050565b600060608201905061107f6000830184611028565b92915050565b61108e81610f93565b82525050565b60006020820190506110a96000830184611085565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006110da826110af565b9050919050565b6110ea816110cf565b81146110f557600080fd5b50565b600081359050611107816110e1565b92915050565b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b61116082611117565b810181811067ffffffffffffffff8211171561117f5761117e611128565b5b80604052505050565b6000611192610eb9565b905061119e8282611157565b919050565b600067ffffffffffffffff8211156111be576111bd611128565b5b6111c782611117565b9050602081019050919050565b82818337600083830152505050565b60006111f66111f1846111a3565b611188565b90508281526020810184848401111561121257611211611112565b5b61121d8482856111d4565b509392505050565b600082601f83011261123a5761123961110d565b5b813561124a8482602086016111e3565b91505092915050565b600067ffffffffffffffff82111561126e5761126d611128565b5b602082029050602081019050919050565b600080fd5b600080fd5b61129281610ffa565b811461129d57600080fd5b50565b6000813590506112af81611289565b92915050565b6000606082840312156112cb576112ca611284565b5b6112d56060611188565b905060006112e5848285016112a0565b60008301525060206112f9848285016112a0565b602083015250604061130d84828501610ef1565b60408301525092915050565b600061132c61132784611253565b611188565b9050808382526020820190506060840283018581111561134f5761134e61127f565b5b835b81811015611378578061136488826112b5565b845260208401935050606081019050611351565b5050509392505050565b600082601f8301126113975761139661110d565b5b81356113a7848260208601611319565b91505092915050565b600080600080600060a086880312156113cc576113cb610ec3565b5b60006113da888289016110f8565b955050602086013567ffffffffffffffff8111156113fb576113fa610ec8565b5b61140788828901611225565b945050604086013567ffffffffffffffff81111561142857611427610ec8565b5b61143488828901611225565b935050606086013567ffffffffffffffff81111561145557611454610ec8565b5b61146188828901611225565b925050608086013567ffffffffffffffff81111561148257611481610ec8565b5b61148e88828901611382565b9150509295509295909350565b600081519050919050565b600082825260208201905092915050565b60005b838110156114d55780820151818401526020810190506114ba565b838111156114e4576000848401525b50505050565b60006114f58261149b565b6114ff81856114a6565b935061150f8185602086016114b7565b61151881611117565b840191505092915050565b6000602082019050818103600083015261153d81846114ea565b905092915050565b60006060828403121561155b5761155a610ec3565b5b6000611569848285016112b5565b91505092915050565b61157b816110cf565b82525050565b60006020820190506115966000830184611572565b92915050565b6000602082840312156115b2576115b1610ec3565b5b60006115c0848285016110f8565b91505092915050565b6115d281610ffa565b82525050565b60006060820190506115ed60008301866115c9565b6115fa60208301856115c9565b6116076040830184610f69565b949350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061165657607f821691505b6020821081141561166a5761166961160f565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600082825260208201905092915050565b7f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160008201527f647920696e697469616c697a6564000000000000000000000000000000000000602082015250565b600061170c602e8361169f565b9150611717826116b0565b604082019050919050565b6000602082019050818103600083015261173b816116ff565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000819050919050565b600061178682611771565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156117b9576117b8611742565b5b600182019050919050565b6000819050919050565b6000819050919050565b60006117f36117ee6117e9846117c4565b6117ce565b610ecd565b9050919050565b611803816117d8565b82525050565b600060208201905061181e60008301846117fa565b92915050565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b600061188060268361169f565b915061188b82611824565b604082019050919050565b600060208201905081810360008301526118af81611873565b9050919050565b7f4269746d61705574696c732e6f72646572656442797465734172726179546f4260008201527f69746d61703a206f7264657265644279746573417272617920697320746f6f2060208201527f6c6f6e6700000000000000000000000000000000000000000000000000000000604082015250565b600061193860448361169f565b9150611943826118b6565b606082019050919050565b600060208201905081810360008301526119678161192b565b9050919050565b7f4269746d61705574696c732e6f72646572656442797465734172726179546f4260008201527f69746d61703a206f72646572656442797465734172726179206973206e6f742060208201527f6f72646572656400000000000000000000000000000000000000000000000000604082015250565b60006119f060478361169f565b91506119fb8261196e565b606082019050919050565b60006020820190508181036000830152611a1f816119e3565b9050919050565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b6000611a5c60208361169f565b9150611a6782611a26565b602082019050919050565b60006020820190508181036000830152611a8b81611a4f565b9050919050565b6000611a9d82610f93565b915061ffff821415611ab257611ab1611742565b5b60018201905091905056fea264697066735822122066cb7ef3bd896b75eff059336d8184f0601754ba19d0b33b96a12c1cd510de0c64736f6c634300080c0033",
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
