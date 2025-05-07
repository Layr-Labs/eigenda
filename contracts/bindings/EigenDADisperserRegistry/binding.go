// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractEigenDADisperserRegistry

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

// EigenDATypesV3DisperserInfo is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV3DisperserInfo struct {
	Disperser    common.Address
	Registered   bool
	DisperserURL string
}

// EigenDATypesV3LockedDisperserDeposit is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV3LockedDisperserDeposit struct {
	Deposit    *big.Int
	Refund     *big.Int
	Token      common.Address
	LockPeriod uint64
}

// ContractEigenDADisperserRegistryMetaData contains all meta data concerning the ContractEigenDADisperserRegistry contract.
var ContractEigenDADisperserRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"deregisterDisperser\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getDepositParams\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV3.LockedDisperserDeposit\",\"components\":[{\"name\":\"deposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"refund\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lockPeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDisperserInfo\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV3.DisperserInfo\",\"components\":[{\"name\":\"disperser\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"registered\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"disperserURL\",\"type\":\"string\",\"internalType\":\"string\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLockedDeposit\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV3.LockedDisperserDeposit\",\"components\":[{\"name\":\"deposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"refund\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lockPeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"unlockTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"depositParams\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV3.LockedDisperserDeposit\",\"components\":[{\"name\":\"deposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"refund\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lockPeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerDisperser\",\"inputs\":[{\"name\":\"disperserAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"disperserURL\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDepositParams\",\"inputs\":[{\"name\":\"depositParams\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV3.LockedDisperserDeposit\",\"components\":[{\"name\":\"deposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"refund\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lockPeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"error\",\"name\":\"AlreadyInitialized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MissingRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x608060405234801561001057600080fd5b50611761806100206000396000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c80638e1e4829116100665780638e1e48291461010d578063a450815814610122578063aa9224cd14610135578063d5f3e3f214610148578063f2fde38b146101a857600080fd5b80631b4b940a146100985780633fcf1fb4146100c55780634fe5de32146100e55780636923ef3a146100fa575b600080fd5b6100ab6100a636600461136d565b6101bb565b60405163ffffffff90911681526020015b60405180910390f35b6100d86100d3366004611412565b6101d0565b6040516100bc9190611490565b6100f86100f3366004611412565b6101f5565b005b6100f861010836600461154f565b610271565b61011561027d565b6040516100bc919061156b565b6100f86101303660046115a8565b610292565b6100f8610143366004611412565b6102cf565b61015b610156366004611412565b610342565b604080518351815260208085015190820152838201516001600160a01b0316918101919091526060928301516001600160401b03908116938201939093529116608082015260a0016100bc565b6100f86101b63660046115dc565b61035e565b60006101c783836103b3565b90505b92915050565b604080516060808201835260008083526020830152918101919091526101ca82610530565b806101ff81610631565b6001600160a01b0316336001600160a01b0316146102645760405162461bcd60e51b815260206004820152601b60248201527f43616c6c6572206973206e6f742074686520646973706572736572000000000060448201526064015b60405180910390fd5b61026d82610660565b5050565b61027a81610752565b50565b610285611238565b61028d610858565b905090565b61029c60016108b4565b6102c67fe579393920b888b1e4a7e1afdd7d58fa4f3101113547ac874aefa75ff4a960f983610930565b61026d81610752565b806102d981610631565b6001600160a01b0316336001600160a01b0316146103395760405162461bcd60e51b815260206004820152601b60248201527f43616c6c6572206973206e6f7420746865206469737065727365720000000000604482015260640161025b565b61026d8261098c565b61034a611238565b600061035583610aac565b91509150915091565b6103887fe579393920b888b1e4a7e1afdd7d58fa4f3101113547ac874aefa75ff4a960f933610b31565b61027a7fe579393920b888b1e4a7e1afdd7d58fa4f3101113547ac874aefa75ff4a960f93383610b6a565b60006103bd610b83565b905060006103c9610bdf565b63ffffffff831660009081526020919091526040902090506001600160a01b0384166104375760405162461bcd60e51b815260206004820152601960248201527f496e76616c696420646973706572736572206164647265737300000000000000604482015260640161025b565b6104653330610444610bdf565b60010154610450610bdf565b600301546001600160a01b0316929190610be9565b604080516060810182526001600160a01b0386168082526001602080840182905293830187905284546001600160a81b031916909117600160a01b1784558551919284926104b892840191880190611272565b509050506104c4610bdf565b6001810154600280840191909155810154600380840191909155018054600490920180546001600160a01b031981166001600160a01b03909416938417825591546001600160401b03600160a01b9182900416026001600160e01b031990921690921717905592915050565b60408051606080820183526000808352602083015291810191909152610554610bdf565b63ffffffff831660009081526020918252604090819020815160608101835281546001600160a01b0381168252600160a01b900460ff1615159381019390935260018101805491928401916105a8906115f7565b80601f01602080910402602001604051908101604052809291908181526020018280546105d4906115f7565b80156106215780601f106105f657610100808354040283529160200191610621565b820191906000526020600020905b81548152906001019060200180831161060457829003601f168201915b5050505050815250509050919050565b600061063b610bdf565b63ffffffff90921660009081526020929092525060409020546001600160a01b031690565b600061066a610bdf565b63ffffffff8316600090815260209190915260408120915061068a610bdf565b63ffffffff8416600090815260209190915260409020825460029091019150600160a01b900460ff166106ff5760405162461bcd60e51b815260206004820152601860248201527f446973706572736572206e6f7420726567697374657265640000000000000000604482015260640161025b565b815460ff60a01b19168255600281015461072990600160a01b90046001600160401b031642611648565b600592909201805467ffffffffffffffff19166001600160401b03909316929092179091555050565b6020810151815110156107a75760405162461bcd60e51b815260206004820152601f60248201527f4465706f736974206d757374206265206174206c6561737420726566756e6400604482015260640161025b565b60408101516001600160a01b03166107f95760405162461bcd60e51b8152602060048201526015602482015274496e76616c696420746f6b656e206164647265737360581b604482015260640161025b565b80610802610bdf565b81516001820155602082015160028201556040820151600390910180546060909301516001600160401b0316600160a01b026001600160e01b03199093166001600160a01b039092169190911791909117905550565b610860611238565b610868610bdf565b6040805160808101825260018301548152600283015460208201526003909201546001600160a01b03811691830191909152600160a01b90046001600160401b03166060820152919050565b8060ff166108c0610c5a565b5460ff16106108e15760405162dc149f60e41b815260040160405180910390fd5b806108ea610c5a565b805460ff191660ff92831617905560405190821681527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a150565b6109518161093c610c64565b60008581526020919091526040902090610c6e565b506040516001600160a01b0382169083907f2ae6a113c0ed5b78a53413ffbb7679881f11145ccfba4fb92e863dfcd5a1d2f390600090a35050565b6000610996610bdf565b63ffffffff831660009081526020919091526040812091506109b6610bdf565b63ffffffff8416600090815260209190915260409020600381015460029091019150610a1d5760405162461bcd60e51b81526020600482015260166024820152754e6f206465706f73697420746f20776974686472617760501b604482015260640161025b565b6005820154426001600160401b039091161115610a7c5760405162461bcd60e51b815260206004820152601760248201527f4465706f736974206973207374696c6c206c6f636b6564000000000000000000604482015260640161025b565b815460018201546002830154610aa0926001600160a01b0391821692911690610c83565b60006001909101555050565b610ab4611238565b600080610abf610bdf565b63ffffffff9094166000908152602094855260409081902060058101548251608081018452600283015481526003830154978101979097526004909101546001600160a01b03811692870192909252600160a01b9091046001600160401b039081166060870152949594169392505050565b610b3b8282610cb3565b61026d576040516301d4003760e61b8152600481018390526001600160a01b038216602482015260440161025b565b610b748383610cd6565b610b7e8382610930565b505050565b600080610b8e610bdf565b6004015463ffffffff169050610ba2610bdf565b600401805463ffffffff16906000610bb983611673565b91906101000a81548163ffffffff021916908363ffffffff160217905550508091505090565b600061028d610d32565b6040516001600160a01b0380851660248301528316604482015260648101829052610c549085906323b872dd60e01b906084015b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b031990931692909217909152610dd7565b50505050565b600061028d610ea9565b600061028d610ef2565b60006101c7836001600160a01b038416610f3c565b6040516001600160a01b038316602482015260448101829052610b7e90849063a9059cbb60e01b90606401610c1d565b60006101c782610cc1610c64565b60008681526020919091526040902090610f8b565b610cf781610ce2610c64565b60008581526020919091526040902090610fad565b506040516001600160a01b0382169083907f155aaafb6329a2098580462df33ec4b7441b19729b9601c5fc17ae1cf99a8a5290600090a35050565b60008060ff60001b1960016040518060400160405280601b81526020017f656967656e2e64612e6469737065727365722e72656769737472790000000000815250604051602001610d839190611697565b6040516020818303038152906040528051906020012060001c610da691906116b3565b604051602001610db891815260200190565b60408051601f1981840301815291905280516020909101201692915050565b6000610e2c826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564815250856001600160a01b0316610fc29092919063ffffffff16565b805190915015610b7e5780806020019051810190610e4a91906116ca565b610b7e5760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6044820152691bdd081cdd58d8d9595960b21b606482015260840161025b565b60008060ff60001b19600160405180604001604052806015815260200174696e697469616c697a61626c652e73746f7261676560581b815250604051602001610d839190611697565b60008060ff60001b196001604051806040016040528060168152602001756163636573732e636f6e74726f6c2e73746f7261676560501b815250604051602001610d839190611697565b6000818152600183016020526040812054610f83575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556101ca565b5060006101ca565b6001600160a01b038116600090815260018301602052604081205415156101c7565b60006101c7836001600160a01b038416610fdb565b6060610fd184846000856110ce565b90505b9392505050565b600081815260018301602052604081205480156110c4576000610fff6001836116b3565b8554909150600090611013906001906116b3565b9050818114611078576000866000018281548110611033576110336116ec565b9060005260206000200154905080876000018481548110611056576110566116ec565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061108957611089611702565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506101ca565b60009150506101ca565b60608247101561112f5760405162461bcd60e51b815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f6044820152651c8818d85b1b60d21b606482015260840161025b565b6001600160a01b0385163b6111865760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640161025b565b600080866001600160a01b031685876040516111a29190611697565b60006040518083038185875af1925050503d80600081146111df576040519150601f19603f3d011682016040523d82523d6000602084013e6111e4565b606091505b50915091506111f48282866111ff565b979650505050505050565b6060831561120e575081610fd4565b82511561121e5782518084602001fd5b8160405162461bcd60e51b815260040161025b9190611718565b6040518060800160405280600081526020016000815260200160006001600160a01b0316815260200160006001600160401b031681525090565b82805461127e906115f7565b90600052602060002090601f0160209004810192826112a057600085556112e6565b82601f106112b957805160ff19168380011785556112e6565b828001600101855582156112e6579182015b828111156112e65782518255916020019190600101906112cb565b506112f29291506112f6565b5090565b5b808211156112f257600081556001016112f7565b80356001600160a01b038116811461132257600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b038111828210171561136557611365611327565b604052919050565b6000806040838503121561138057600080fd5b6113898361130b565b91506020808401356001600160401b03808211156113a657600080fd5b818601915086601f8301126113ba57600080fd5b8135818111156113cc576113cc611327565b6113de601f8201601f1916850161133d565b915080825287848285010111156113f457600080fd5b80848401858401376000848284010152508093505050509250929050565b60006020828403121561142457600080fd5b813563ffffffff81168114610fd457600080fd5b60005b8381101561145357818101518382015260200161143b565b83811115610c545750506000910152565b6000815180845261147c816020860160208601611438565b601f01601f19169290920160200192915050565b6020815260018060a01b038251166020820152602082015115156040820152600060408301516060808401526114c96080840182611464565b949350505050565b6000608082840312156114e357600080fd5b604051608081016001600160401b03828210818311171561150657611506611327565b8160405282935084358352602085013560208401526115276040860161130b565b604084015260608501359150808216821461154157600080fd5b506060919091015292915050565b60006080828403121561156157600080fd5b6101c783836114d1565b81518152602080830151908201526040808301516001600160a01b0316908201526060808301516001600160401b031690820152608081016101ca565b60008060a083850312156115bb57600080fd5b6115c48361130b565b91506115d384602085016114d1565b90509250929050565b6000602082840312156115ee57600080fd5b6101c78261130b565b600181811c9082168061160b57607f821691505b6020821081141561162c57634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b60006001600160401b0380831681851680830382111561166a5761166a611632565b01949350505050565b600063ffffffff8083168181141561168d5761168d611632565b6001019392505050565b600082516116a9818460208701611438565b9190910192915050565b6000828210156116c5576116c5611632565b500390565b6000602082840312156116dc57600080fd5b81518015158114610fd457600080fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052603160045260246000fd5b6020815260006101c7602083018461146456fea2646970667358221220abffb17de5ecdac66da92e6257f4d3009bcb963f09a4717a05f69e9f9678d54064736f6c634300080c0033",
}

// ContractEigenDADisperserRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEigenDADisperserRegistryMetaData.ABI instead.
var ContractEigenDADisperserRegistryABI = ContractEigenDADisperserRegistryMetaData.ABI

// ContractEigenDADisperserRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractEigenDADisperserRegistryMetaData.Bin instead.
var ContractEigenDADisperserRegistryBin = ContractEigenDADisperserRegistryMetaData.Bin

// DeployContractEigenDADisperserRegistry deploys a new Ethereum contract, binding an instance of ContractEigenDADisperserRegistry to it.
func DeployContractEigenDADisperserRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ContractEigenDADisperserRegistry, error) {
	parsed, err := ContractEigenDADisperserRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractEigenDADisperserRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractEigenDADisperserRegistry{ContractEigenDADisperserRegistryCaller: ContractEigenDADisperserRegistryCaller{contract: contract}, ContractEigenDADisperserRegistryTransactor: ContractEigenDADisperserRegistryTransactor{contract: contract}, ContractEigenDADisperserRegistryFilterer: ContractEigenDADisperserRegistryFilterer{contract: contract}}, nil
}

// ContractEigenDADisperserRegistry is an auto generated Go binding around an Ethereum contract.
type ContractEigenDADisperserRegistry struct {
	ContractEigenDADisperserRegistryCaller     // Read-only binding to the contract
	ContractEigenDADisperserRegistryTransactor // Write-only binding to the contract
	ContractEigenDADisperserRegistryFilterer   // Log filterer for contract events
}

// ContractEigenDADisperserRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractEigenDADisperserRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDADisperserRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractEigenDADisperserRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDADisperserRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractEigenDADisperserRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDADisperserRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractEigenDADisperserRegistrySession struct {
	Contract     *ContractEigenDADisperserRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                     // Call options to use throughout this session
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// ContractEigenDADisperserRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractEigenDADisperserRegistryCallerSession struct {
	Contract *ContractEigenDADisperserRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                           // Call options to use throughout this session
}

// ContractEigenDADisperserRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractEigenDADisperserRegistryTransactorSession struct {
	Contract     *ContractEigenDADisperserRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                           // Transaction auth options to use throughout this session
}

// ContractEigenDADisperserRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractEigenDADisperserRegistryRaw struct {
	Contract *ContractEigenDADisperserRegistry // Generic contract binding to access the raw methods on
}

// ContractEigenDADisperserRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractEigenDADisperserRegistryCallerRaw struct {
	Contract *ContractEigenDADisperserRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// ContractEigenDADisperserRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractEigenDADisperserRegistryTransactorRaw struct {
	Contract *ContractEigenDADisperserRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractEigenDADisperserRegistry creates a new instance of ContractEigenDADisperserRegistry, bound to a specific deployed contract.
func NewContractEigenDADisperserRegistry(address common.Address, backend bind.ContractBackend) (*ContractEigenDADisperserRegistry, error) {
	contract, err := bindContractEigenDADisperserRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDADisperserRegistry{ContractEigenDADisperserRegistryCaller: ContractEigenDADisperserRegistryCaller{contract: contract}, ContractEigenDADisperserRegistryTransactor: ContractEigenDADisperserRegistryTransactor{contract: contract}, ContractEigenDADisperserRegistryFilterer: ContractEigenDADisperserRegistryFilterer{contract: contract}}, nil
}

// NewContractEigenDADisperserRegistryCaller creates a new read-only instance of ContractEigenDADisperserRegistry, bound to a specific deployed contract.
func NewContractEigenDADisperserRegistryCaller(address common.Address, caller bind.ContractCaller) (*ContractEigenDADisperserRegistryCaller, error) {
	contract, err := bindContractEigenDADisperserRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDADisperserRegistryCaller{contract: contract}, nil
}

// NewContractEigenDADisperserRegistryTransactor creates a new write-only instance of ContractEigenDADisperserRegistry, bound to a specific deployed contract.
func NewContractEigenDADisperserRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractEigenDADisperserRegistryTransactor, error) {
	contract, err := bindContractEigenDADisperserRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDADisperserRegistryTransactor{contract: contract}, nil
}

// NewContractEigenDADisperserRegistryFilterer creates a new log filterer instance of ContractEigenDADisperserRegistry, bound to a specific deployed contract.
func NewContractEigenDADisperserRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractEigenDADisperserRegistryFilterer, error) {
	contract, err := bindContractEigenDADisperserRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDADisperserRegistryFilterer{contract: contract}, nil
}

// bindContractEigenDADisperserRegistry binds a generic wrapper to an already deployed contract.
func bindContractEigenDADisperserRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractEigenDADisperserRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDADisperserRegistry.Contract.ContractEigenDADisperserRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.ContractEigenDADisperserRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.ContractEigenDADisperserRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDADisperserRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.contract.Transact(opts, method, params...)
}

// GetDepositParams is a free data retrieval call binding the contract method 0x8e1e4829.
//
// Solidity: function getDepositParams() view returns((uint256,uint256,address,uint64))
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryCaller) GetDepositParams(opts *bind.CallOpts) (EigenDATypesV3LockedDisperserDeposit, error) {
	var out []interface{}
	err := _ContractEigenDADisperserRegistry.contract.Call(opts, &out, "getDepositParams")

	if err != nil {
		return *new(EigenDATypesV3LockedDisperserDeposit), err
	}

	out0 := *abi.ConvertType(out[0], new(EigenDATypesV3LockedDisperserDeposit)).(*EigenDATypesV3LockedDisperserDeposit)

	return out0, err

}

// GetDepositParams is a free data retrieval call binding the contract method 0x8e1e4829.
//
// Solidity: function getDepositParams() view returns((uint256,uint256,address,uint64))
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistrySession) GetDepositParams() (EigenDATypesV3LockedDisperserDeposit, error) {
	return _ContractEigenDADisperserRegistry.Contract.GetDepositParams(&_ContractEigenDADisperserRegistry.CallOpts)
}

// GetDepositParams is a free data retrieval call binding the contract method 0x8e1e4829.
//
// Solidity: function getDepositParams() view returns((uint256,uint256,address,uint64))
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryCallerSession) GetDepositParams() (EigenDATypesV3LockedDisperserDeposit, error) {
	return _ContractEigenDADisperserRegistry.Contract.GetDepositParams(&_ContractEigenDADisperserRegistry.CallOpts)
}

// GetDisperserInfo is a free data retrieval call binding the contract method 0x3fcf1fb4.
//
// Solidity: function getDisperserInfo(uint32 disperserKey) view returns((address,bool,string))
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryCaller) GetDisperserInfo(opts *bind.CallOpts, disperserKey uint32) (EigenDATypesV3DisperserInfo, error) {
	var out []interface{}
	err := _ContractEigenDADisperserRegistry.contract.Call(opts, &out, "getDisperserInfo", disperserKey)

	if err != nil {
		return *new(EigenDATypesV3DisperserInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(EigenDATypesV3DisperserInfo)).(*EigenDATypesV3DisperserInfo)

	return out0, err

}

// GetDisperserInfo is a free data retrieval call binding the contract method 0x3fcf1fb4.
//
// Solidity: function getDisperserInfo(uint32 disperserKey) view returns((address,bool,string))
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistrySession) GetDisperserInfo(disperserKey uint32) (EigenDATypesV3DisperserInfo, error) {
	return _ContractEigenDADisperserRegistry.Contract.GetDisperserInfo(&_ContractEigenDADisperserRegistry.CallOpts, disperserKey)
}

// GetDisperserInfo is a free data retrieval call binding the contract method 0x3fcf1fb4.
//
// Solidity: function getDisperserInfo(uint32 disperserKey) view returns((address,bool,string))
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryCallerSession) GetDisperserInfo(disperserKey uint32) (EigenDATypesV3DisperserInfo, error) {
	return _ContractEigenDADisperserRegistry.Contract.GetDisperserInfo(&_ContractEigenDADisperserRegistry.CallOpts, disperserKey)
}

// GetLockedDeposit is a free data retrieval call binding the contract method 0xd5f3e3f2.
//
// Solidity: function getLockedDeposit(uint32 disperserKey) view returns((uint256,uint256,address,uint64), uint64 unlockTimestamp)
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryCaller) GetLockedDeposit(opts *bind.CallOpts, disperserKey uint32) (EigenDATypesV3LockedDisperserDeposit, uint64, error) {
	var out []interface{}
	err := _ContractEigenDADisperserRegistry.contract.Call(opts, &out, "getLockedDeposit", disperserKey)

	if err != nil {
		return *new(EigenDATypesV3LockedDisperserDeposit), *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(EigenDATypesV3LockedDisperserDeposit)).(*EigenDATypesV3LockedDisperserDeposit)
	out1 := *abi.ConvertType(out[1], new(uint64)).(*uint64)

	return out0, out1, err

}

// GetLockedDeposit is a free data retrieval call binding the contract method 0xd5f3e3f2.
//
// Solidity: function getLockedDeposit(uint32 disperserKey) view returns((uint256,uint256,address,uint64), uint64 unlockTimestamp)
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistrySession) GetLockedDeposit(disperserKey uint32) (EigenDATypesV3LockedDisperserDeposit, uint64, error) {
	return _ContractEigenDADisperserRegistry.Contract.GetLockedDeposit(&_ContractEigenDADisperserRegistry.CallOpts, disperserKey)
}

// GetLockedDeposit is a free data retrieval call binding the contract method 0xd5f3e3f2.
//
// Solidity: function getLockedDeposit(uint32 disperserKey) view returns((uint256,uint256,address,uint64), uint64 unlockTimestamp)
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryCallerSession) GetLockedDeposit(disperserKey uint32) (EigenDATypesV3LockedDisperserDeposit, uint64, error) {
	return _ContractEigenDADisperserRegistry.Contract.GetLockedDeposit(&_ContractEigenDADisperserRegistry.CallOpts, disperserKey)
}

// DeregisterDisperser is a paid mutator transaction binding the contract method 0x4fe5de32.
//
// Solidity: function deregisterDisperser(uint32 disperserKey) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryTransactor) DeregisterDisperser(opts *bind.TransactOpts, disperserKey uint32) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.contract.Transact(opts, "deregisterDisperser", disperserKey)
}

// DeregisterDisperser is a paid mutator transaction binding the contract method 0x4fe5de32.
//
// Solidity: function deregisterDisperser(uint32 disperserKey) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistrySession) DeregisterDisperser(disperserKey uint32) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.DeregisterDisperser(&_ContractEigenDADisperserRegistry.TransactOpts, disperserKey)
}

// DeregisterDisperser is a paid mutator transaction binding the contract method 0x4fe5de32.
//
// Solidity: function deregisterDisperser(uint32 disperserKey) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryTransactorSession) DeregisterDisperser(disperserKey uint32) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.DeregisterDisperser(&_ContractEigenDADisperserRegistry.TransactOpts, disperserKey)
}

// Initialize is a paid mutator transaction binding the contract method 0xa4508158.
//
// Solidity: function initialize(address owner, (uint256,uint256,address,uint64) depositParams) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryTransactor) Initialize(opts *bind.TransactOpts, owner common.Address, depositParams EigenDATypesV3LockedDisperserDeposit) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.contract.Transact(opts, "initialize", owner, depositParams)
}

// Initialize is a paid mutator transaction binding the contract method 0xa4508158.
//
// Solidity: function initialize(address owner, (uint256,uint256,address,uint64) depositParams) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistrySession) Initialize(owner common.Address, depositParams EigenDATypesV3LockedDisperserDeposit) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.Initialize(&_ContractEigenDADisperserRegistry.TransactOpts, owner, depositParams)
}

// Initialize is a paid mutator transaction binding the contract method 0xa4508158.
//
// Solidity: function initialize(address owner, (uint256,uint256,address,uint64) depositParams) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryTransactorSession) Initialize(owner common.Address, depositParams EigenDATypesV3LockedDisperserDeposit) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.Initialize(&_ContractEigenDADisperserRegistry.TransactOpts, owner, depositParams)
}

// RegisterDisperser is a paid mutator transaction binding the contract method 0x1b4b940a.
//
// Solidity: function registerDisperser(address disperserAddress, string disperserURL) returns(uint32 disperserKey)
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryTransactor) RegisterDisperser(opts *bind.TransactOpts, disperserAddress common.Address, disperserURL string) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.contract.Transact(opts, "registerDisperser", disperserAddress, disperserURL)
}

// RegisterDisperser is a paid mutator transaction binding the contract method 0x1b4b940a.
//
// Solidity: function registerDisperser(address disperserAddress, string disperserURL) returns(uint32 disperserKey)
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistrySession) RegisterDisperser(disperserAddress common.Address, disperserURL string) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.RegisterDisperser(&_ContractEigenDADisperserRegistry.TransactOpts, disperserAddress, disperserURL)
}

// RegisterDisperser is a paid mutator transaction binding the contract method 0x1b4b940a.
//
// Solidity: function registerDisperser(address disperserAddress, string disperserURL) returns(uint32 disperserKey)
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryTransactorSession) RegisterDisperser(disperserAddress common.Address, disperserURL string) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.RegisterDisperser(&_ContractEigenDADisperserRegistry.TransactOpts, disperserAddress, disperserURL)
}

// SetDepositParams is a paid mutator transaction binding the contract method 0x6923ef3a.
//
// Solidity: function setDepositParams((uint256,uint256,address,uint64) depositParams) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryTransactor) SetDepositParams(opts *bind.TransactOpts, depositParams EigenDATypesV3LockedDisperserDeposit) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.contract.Transact(opts, "setDepositParams", depositParams)
}

// SetDepositParams is a paid mutator transaction binding the contract method 0x6923ef3a.
//
// Solidity: function setDepositParams((uint256,uint256,address,uint64) depositParams) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistrySession) SetDepositParams(depositParams EigenDATypesV3LockedDisperserDeposit) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.SetDepositParams(&_ContractEigenDADisperserRegistry.TransactOpts, depositParams)
}

// SetDepositParams is a paid mutator transaction binding the contract method 0x6923ef3a.
//
// Solidity: function setDepositParams((uint256,uint256,address,uint64) depositParams) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryTransactorSession) SetDepositParams(depositParams EigenDATypesV3LockedDisperserDeposit) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.SetDepositParams(&_ContractEigenDADisperserRegistry.TransactOpts, depositParams)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.TransferOwnership(&_ContractEigenDADisperserRegistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.TransferOwnership(&_ContractEigenDADisperserRegistry.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0xaa9224cd.
//
// Solidity: function withdraw(uint32 disperserKey) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryTransactor) Withdraw(opts *bind.TransactOpts, disperserKey uint32) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.contract.Transact(opts, "withdraw", disperserKey)
}

// Withdraw is a paid mutator transaction binding the contract method 0xaa9224cd.
//
// Solidity: function withdraw(uint32 disperserKey) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistrySession) Withdraw(disperserKey uint32) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.Withdraw(&_ContractEigenDADisperserRegistry.TransactOpts, disperserKey)
}

// Withdraw is a paid mutator transaction binding the contract method 0xaa9224cd.
//
// Solidity: function withdraw(uint32 disperserKey) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryTransactorSession) Withdraw(disperserKey uint32) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.Withdraw(&_ContractEigenDADisperserRegistry.TransactOpts, disperserKey)
}
