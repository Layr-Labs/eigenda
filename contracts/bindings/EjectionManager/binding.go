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
	Bin: "0x60c06040523480156200001157600080fd5b506040516200289b3803806200289b833981810160405281019062000037919062000241565b8173ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff16815250508073ffffffffffffffffffffffffffffffffffffffff1660a08173ffffffffffffffffffffffffffffffffffffffff1681525050620000af620000b760201b60201c565b50506200036c565b600060019054906101000a900460ff16156200010a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040162000101906200030f565b60405180910390fd5b60ff801660008054906101000a900460ff1660ff1610156200017c5760ff6000806101000a81548160ff021916908360ff1602179055507f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb384740249860ff6040516200017391906200034f565b60405180910390a15b565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000620001b08262000183565b9050919050565b6000620001c482620001a3565b9050919050565b620001d681620001b7565b8114620001e257600080fd5b50565b600081519050620001f681620001cb565b92915050565b60006200020982620001a3565b9050919050565b6200021b81620001fc565b81146200022757600080fd5b50565b6000815190506200023b8162000210565b92915050565b600080604083850312156200025b576200025a6200017e565b5b60006200026b85828601620001e5565b92505060206200027e858286016200022a565b9150509250929050565b600082825260208201905092915050565b7f496e697469616c697a61626c653a20636f6e747261637420697320696e69746960008201527f616c697a696e6700000000000000000000000000000000000000000000000000602082015250565b6000620002f760278362000288565b9150620003048262000299565b604082019050919050565b600060208201905081810360008301526200032a81620002e8565b9050919050565b600060ff82169050919050565b620003498162000331565b82525050565b60006020820190506200036660008301846200033e565b92915050565b60805160a0516124df620003bc6000396000818161042501528181610b4e0152610e2d0152600081816105cc01528181610608015281816107d8015281816108140152610b9201526124df6000f3fe608060405234801561001057600080fd5b50600436106100ce5760003560e01c80636d14a9871161008c5780638b88a024116100665780638b88a024146101ff5780638da5cb5b1461021b578063b13f450414610239578063f2fde38b14610269576100ce565b80636d14a987146101bb578063715018a6146101d957806377d17586146101e3576100ce565b8062482569146100d35780630a0593d11461010457806310ea4f8a146101205780633a0b0ddd1461013c578063683048351461016d5780636c08a8791461018b575b600080fd5b6100ed60048036038101906100e891906113d3565b610285565b6040516100fb92919061143c565b60405180910390f35b61011e600480360381019061011991906116d5565b6102c7565b005b61013a600480360381019061013591906117b4565b610af5565b005b6101566004803603810190610151919061182a565b610b0b565b604051610164929190611879565b60405180910390f35b610175610b4c565b6040516101829190611901565b60405180910390f35b6101a560048036038101906101a0919061191c565b610b70565b6040516101b29190611958565b60405180910390f35b6101c3610b90565b6040516101d09190611994565b60405180910390f35b6101e1610bb4565b005b6101fd60048036038101906101f89190611a5c565b610bc8565b005b61021960048036038101906102149190611c22565b610bde565b005b610223610db5565b6040516102309190611cbc565b60405180910390f35b610253600480360381019061024e91906113d3565b610ddf565b6040516102609190611cd7565b60405180910390f35b610283600480360381019061027e919061191c565b611051565b005b60676020528060005260406000206000915090508060000160009054906101000a900463ffffffff16908060000160049054906101000a900461ffff16905082565b606560003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16806103515750610322610db5565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b610390576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161038790611d75565b60405180910390fd5b60005b8151811015610af157600081905060006103ac82610ddf565b90506000806000808411806103f357506103c4610db5565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b156109d65760005b87878151811061040e5761040d611d95565b5b6020026020010151518160ff1610156109d45760007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16635401ed278a8a8151811061047257610471611d95565b5b60200260200101518460ff168151811061048f5761048e611d95565b5b6020026020010151896040518363ffffffff1660e01b81526004016104b5929190611de2565b602060405180830381865afa1580156104d2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104f69190611e4f565b6bffffffffffffffffffffffff169050606560003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16801561059157506000606760008960ff1660ff16815260200190815260200160002060000160009054906101000a900463ffffffff1663ffffffff16115b80156105a757508581866105a59190611eab565b115b156107bc576001925080856105bc9190611eab565b9450836105c890611f01565b93507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16636e3b17db7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663296bb0648c8c8151811061065557610654611d95565b5b60200260200101518660ff168151811061067257610671611d95565b5b60200260200101516040518263ffffffff1660e01b81526004016106969190611f2e565b602060405180830381865afa1580156106b3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106d79190611f5e565b896040516020016106e89190611fc1565b6040516020818303038152906040526040518363ffffffff1660e01b8152600401610714929190612064565b600060405180830381600087803b15801561072e57600080fd5b505af1158015610742573d6000803e3d6000fd5b505050507f97ddb711c61a9d2d7effcba3e042a33862297f898d555655cca39ec4451f53b489898151811061077a57610779611d95565b5b60200260200101518360ff168151811061079757610796611d95565b5b6020026020010151886040516107ae929190611de2565b60405180910390a1506109d4565b80856107c89190611eab565b9450836107d490611f01565b93507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16636e3b17db7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663296bb0648c8c8151811061086157610860611d95565b5b60200260200101518660ff168151811061087e5761087d611d95565b5b60200260200101516040518263ffffffff1660e01b81526004016108a29190611f2e565b602060405180830381865afa1580156108bf573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108e39190611f5e565b896040516020016108f49190611fc1565b6040516020818303038152906040526040518363ffffffff1660e01b8152600401610920929190612064565b600060405180830381600087803b15801561093a57600080fd5b505af115801561094e573d6000803e3d6000fd5b505050507f97ddb711c61a9d2d7effcba3e042a33862297f898d555655cca39ec4451f53b489898151811061098657610985611d95565b5b60200260200101518360ff16815181106109a3576109a2611d95565b5b6020026020010151886040516109ba929190611de2565b60405180910390a150806109cd90612094565b90506103fb565b505b606560003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff168015610a2f5750600083115b15610aa257606660008660ff1660ff1681526020019081526020016000206040518060400160405280428152602001858152509080600181540180825580915050600190039060005260206000209060020201600090919091909150600082015181600001556020820151816001015550505b7f19dd87ae49ed14a795f8c2d5e8055bf2a4a9d01641a00a2f8f0a5a7bf7f702498282604051610ad39291906120be565b60405180910390a1505050505080610aea906120e7565b9050610393565b5050565b610afd6110d5565b610b078282611153565b5050565b60666020528160005260406000208181548110610b2757600080fd5b9060005260206000209060020201600091509150508060000154908060010154905082565b7f000000000000000000000000000000000000000000000000000000000000000081565b60656020528060005260406000206000915054906101000a900460ff1681565b7f000000000000000000000000000000000000000000000000000000000000000081565b610bbc6110d5565b610bc660006111e7565b565b610bd06110d5565b610bda82826112ad565b5050565b60008060019054906101000a900460ff16159050808015610c0f5750600160008054906101000a900460ff1660ff16105b80610c3c5750610c1e3061135b565b158015610c3b5750600160008054906101000a900460ff1660ff16145b5b610c7b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c72906121a2565b60405180910390fd5b60016000806101000a81548160ff021916908360ff1602179055508015610cb8576001600060016101000a81548160ff0219169083151502179055505b610cc1846111e7565b60005b83518160ff161015610d0b57610cf8848260ff1681518110610ce957610ce8611d95565b5b60200260200101516001611153565b8080610d0390612094565b915050610cc4565b5060005b82518160ff161015610d5557610d4281848360ff1681518110610d3557610d34611d95565b5b60200260200101516112ad565b8080610d4d90612094565b915050610d0f565b508015610daf5760008060016101000a81548160ff0219169083151502179055507f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024986001604051610da691906121fd565b60405180910390a15b50505050565b6000603360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b600080606760008460ff1660ff16815260200190815260200160002060000160009054906101000a900463ffffffff1663ffffffff1642610e209190612218565b9050600061271061ffff167f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663d5eccc05866040518263ffffffff1660e01b8152600401610e84919061224c565b602060405180830381865afa158015610ea1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ec59190611e4f565b6bffffffffffffffffffffffff16606760008760ff1660ff16815260200190815260200160002060000160049054906101000a900461ffff1661ffff16610f0c9190612267565b610f1691906122f0565b90506000806000606660008860ff1660ff168152602001908152602001600020805490501415610f4c578294505050505061104c565b6001606660008860ff1660ff16815260200190815260200160002080549050610f759190612218565b90505b83606660008860ff1660ff1681526020019081526020016000208281548110610fa457610fa3611d95565b5b906000526020600020906002020160000154111561102557606660008760ff1660ff1681526020019081526020016000208181548110610fe757610fe6611d95565b5b906000526020600020906002020160010154826110049190611eab565b9150600081141561101457611025565b8061101e90612321565b9050610f78565b82821061103957600094505050505061104c565b81836110459190612218565b9450505050505b919050565b6110596110d5565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614156110c9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016110c0906123bd565b60405180910390fd5b6110d2816111e7565b50565b6110dd61137e565b73ffffffffffffffffffffffffffffffffffffffff166110fb610db5565b73ffffffffffffffffffffffffffffffffffffffff1614611151576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161114890612429565b60405180910390fd5b565b80606560008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055507f7676686b6d22e112412bd874d70177e011ab06602c26063f19f0386c9a3cee4282826040516111db929190612449565b60405180910390a15050565b6000603360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905081603360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b80606760008460ff1660ff16815260200190815260200160002060008201518160000160006101000a81548163ffffffff021916908363ffffffff16021790555060208201518160000160046101000a81548161ffff021916908361ffff1602179055509050507fe69c2827a1e2fdd32265ebb4eeea5ee564f0551cf5dfed4150f8e116a67209eb828260000151836020015160405161134f93929190612472565b60405180910390a15050565b6000808273ffffffffffffffffffffffffffffffffffffffff163b119050919050565b600033905090565b6000604051905090565b600080fd5b600080fd5b600060ff82169050919050565b6113b08161139a565b81146113bb57600080fd5b50565b6000813590506113cd816113a7565b92915050565b6000602082840312156113e9576113e8611390565b5b60006113f7848285016113be565b91505092915050565b600063ffffffff82169050919050565b61141981611400565b82525050565b600061ffff82169050919050565b6114368161141f565b82525050565b60006040820190506114516000830185611410565b61145e602083018461142d565b9392505050565b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6114b38261146a565b810181811067ffffffffffffffff821117156114d2576114d161147b565b5b80604052505050565b60006114e5611386565b90506114f182826114aa565b919050565b600067ffffffffffffffff8211156115115761151061147b565b5b602082029050602081019050919050565b600080fd5b600067ffffffffffffffff8211156115425761154161147b565b5b602082029050602081019050919050565b6000819050919050565b61156681611553565b811461157157600080fd5b50565b6000813590506115838161155d565b92915050565b600061159c61159784611527565b6114db565b905080838252602082019050602084028301858111156115bf576115be611522565b5b835b818110156115e857806115d48882611574565b8452602084019350506020810190506115c1565b5050509392505050565b600082601f83011261160757611606611465565b5b8135611617848260208601611589565b91505092915050565b600061163361162e846114f6565b6114db565b9050808382526020820190506020840283018581111561165657611655611522565b5b835b8181101561169d57803567ffffffffffffffff81111561167b5761167a611465565b5b80860161168889826115f2565b85526020850194505050602081019050611658565b5050509392505050565b600082601f8301126116bc576116bb611465565b5b81356116cc848260208601611620565b91505092915050565b6000602082840312156116eb576116ea611390565b5b600082013567ffffffffffffffff81111561170957611708611395565b5b611715848285016116a7565b91505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006117498261171e565b9050919050565b6117598161173e565b811461176457600080fd5b50565b60008135905061177681611750565b92915050565b60008115159050919050565b6117918161177c565b811461179c57600080fd5b50565b6000813590506117ae81611788565b92915050565b600080604083850312156117cb576117ca611390565b5b60006117d985828601611767565b92505060206117ea8582860161179f565b9150509250929050565b6000819050919050565b611807816117f4565b811461181257600080fd5b50565b600081359050611824816117fe565b92915050565b6000806040838503121561184157611840611390565b5b600061184f858286016113be565b925050602061186085828601611815565b9150509250929050565b611873816117f4565b82525050565b600060408201905061188e600083018561186a565b61189b602083018461186a565b9392505050565b6000819050919050565b60006118c76118c26118bd8461171e565b6118a2565b61171e565b9050919050565b60006118d9826118ac565b9050919050565b60006118eb826118ce565b9050919050565b6118fb816118e0565b82525050565b600060208201905061191660008301846118f2565b92915050565b60006020828403121561193257611931611390565b5b600061194084828501611767565b91505092915050565b6119528161177c565b82525050565b600060208201905061196d6000830184611949565b92915050565b600061197e826118ce565b9050919050565b61198e81611973565b82525050565b60006020820190506119a96000830184611985565b92915050565b600080fd5b6119bd81611400565b81146119c857600080fd5b50565b6000813590506119da816119b4565b92915050565b6119e98161141f565b81146119f457600080fd5b50565b600081359050611a06816119e0565b92915050565b600060408284031215611a2257611a216119af565b5b611a2c60406114db565b90506000611a3c848285016119cb565b6000830152506020611a50848285016119f7565b60208301525092915050565b60008060608385031215611a7357611a72611390565b5b6000611a81858286016113be565b9250506020611a9285828601611a0c565b9150509250929050565b600067ffffffffffffffff821115611ab757611ab661147b565b5b602082029050602081019050919050565b6000611adb611ad684611a9c565b6114db565b90508083825260208201905060208402830185811115611afe57611afd611522565b5b835b81811015611b275780611b138882611767565b845260208401935050602081019050611b00565b5050509392505050565b600082601f830112611b4657611b45611465565b5b8135611b56848260208601611ac8565b91505092915050565b600067ffffffffffffffff821115611b7a57611b7961147b565b5b602082029050602081019050919050565b6000611b9e611b9984611b5f565b6114db565b90508083825260208201905060408402830185811115611bc157611bc0611522565b5b835b81811015611bea5780611bd68882611a0c565b845260208401935050604081019050611bc3565b5050509392505050565b600082601f830112611c0957611c08611465565b5b8135611c19848260208601611b8b565b91505092915050565b600080600060608486031215611c3b57611c3a611390565b5b6000611c4986828701611767565b935050602084013567ffffffffffffffff811115611c6a57611c69611395565b5b611c7686828701611b31565b925050604084013567ffffffffffffffff811115611c9757611c96611395565b5b611ca386828701611bf4565b9150509250925092565b611cb68161173e565b82525050565b6000602082019050611cd16000830184611cad565b92915050565b6000602082019050611cec600083018461186a565b92915050565b600082825260208201905092915050565b7f456a6563746f723a204f6e6c79206f776e6572206f7220656a6563746f72206360008201527f616e20656a656374000000000000000000000000000000000000000000000000602082015250565b6000611d5f602883611cf2565b9150611d6a82611d03565b604082019050919050565b60006020820190508181036000830152611d8e81611d52565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b611dcd81611553565b82525050565b611ddc8161139a565b82525050565b6000604082019050611df76000830185611dc4565b611e046020830184611dd3565b9392505050565b60006bffffffffffffffffffffffff82169050919050565b611e2c81611e0b565b8114611e3757600080fd5b50565b600081519050611e4981611e23565b92915050565b600060208284031215611e6557611e64611390565b5b6000611e7384828501611e3a565b91505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000611eb6826117f4565b9150611ec1836117f4565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03821115611ef657611ef5611e7c565b5b828201905092915050565b6000611f0c82611400565b915063ffffffff821415611f2357611f22611e7c565b5b600182019050919050565b6000602082019050611f436000830184611dc4565b92915050565b600081519050611f5881611750565b92915050565b600060208284031215611f7457611f73611390565b5b6000611f8284828501611f49565b91505092915050565b60008160f81b9050919050565b6000611fa382611f8b565b9050919050565b611fbb611fb68261139a565b611f98565b82525050565b6000611fcd8284611faa565b60018201915081905092915050565b600081519050919050565b600082825260208201905092915050565b60005b83811015612016578082015181840152602081019050611ffb565b83811115612025576000848401525b50505050565b600061203682611fdc565b6120408185611fe7565b9350612050818560208601611ff8565b6120598161146a565b840191505092915050565b60006040820190506120796000830185611cad565b818103602083015261208b818461202b565b90509392505050565b600061209f8261139a565b915060ff8214156120b3576120b2611e7c565b5b600182019050919050565b60006040820190506120d36000830185611410565b6120e06020830184611949565b9392505050565b60006120f2826117f4565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82141561212557612124611e7c565b5b600182019050919050565b7f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160008201527f647920696e697469616c697a6564000000000000000000000000000000000000602082015250565b600061218c602e83611cf2565b915061219782612130565b604082019050919050565b600060208201905081810360008301526121bb8161217f565b9050919050565b6000819050919050565b60006121e76121e26121dd846121c2565b6118a2565b61139a565b9050919050565b6121f7816121cc565b82525050565b600060208201905061221260008301846121ee565b92915050565b6000612223826117f4565b915061222e836117f4565b92508282101561224157612240611e7c565b5b828203905092915050565b60006020820190506122616000830184611dd3565b92915050565b6000612272826117f4565b915061227d836117f4565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156122b6576122b5611e7c565b5b828202905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60006122fb826117f4565b9150612306836117f4565b925082612316576123156122c1565b5b828204905092915050565b600061232c826117f4565b915060008214156123405761233f611e7c565b5b600182039050919050565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b60006123a7602683611cf2565b91506123b28261234b565b604082019050919050565b600060208201905081810360008301526123d68161239a565b9050919050565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b6000612413602083611cf2565b915061241e826123dd565b602082019050919050565b6000602082019050818103600083015261244281612406565b9050919050565b600060408201905061245e6000830185611cad565b61246b6020830184611949565b9392505050565b60006060820190506124876000830186611dd3565b6124946020830185611410565b6124a1604083018461142d565b94935050505056fea26469706673582212201b4273a8bbf7c8d2db8c80d90b198eebc24916360160eaa1c57165191a61d32e64736f6c634300080c0033",
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
