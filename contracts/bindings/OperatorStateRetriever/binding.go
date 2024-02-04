// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractOperatorStateRetriever

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

// OperatorStateRetrieverCheckSignaturesIndices is an auto generated low-level Go binding around an user-defined struct.
type OperatorStateRetrieverCheckSignaturesIndices struct {
	NonSignerQuorumBitmapIndices []uint32
	QuorumApkIndices             []uint32
	TotalStakeIndices            []uint32
	NonSignerStakeIndices        [][]uint32
}

// OperatorStateRetrieverOperator is an auto generated low-level Go binding around an user-defined struct.
type OperatorStateRetrieverOperator struct {
	Operator   common.Address
	OperatorId [32]byte
	Stake      *big.Int
}

// ContractOperatorStateRetrieverMetaData contains all meta data concerning the ContractOperatorStateRetriever contract.
var ContractOperatorStateRetrieverMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"getCheckSignaturesIndices\",\"inputs\":[{\"name\":\"registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"nonSignerOperatorIds\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structOperatorStateRetriever.CheckSignaturesIndices\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorState\",\"inputs\":[{\"name\":\"registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"blockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[][]\",\"internalType\":\"structOperatorStateRetriever.Operator[][]\",\"components\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"stake\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorState\",\"inputs\":[{\"name\":\"registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"},{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"blockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"tuple[][]\",\"internalType\":\"structOperatorStateRetriever.Operator[][]\",\"components\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"stake\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"stateMutability\":\"view\"}]",
	Bin: "0x608060405234801561001057600080fd5b50611722806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80633563b0d1146100465780634f739f741461006f578063cefdc1d41461008f575b600080fd5b610059610054366004610f7f565b6100b0565b60405161006691906110da565b60405180910390f35b61008261007d36600461113f565b610546565b6040516100669190611242565b6100a261009d3660046112fd565b610c70565b60405161006692919061133f565b60606000846001600160a01b031663683048356040518163ffffffff1660e01b8152600401602060405180830381865afa1580156100f2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101169190611360565b90506000856001600160a01b0316639e9923c26040518163ffffffff1660e01b8152600401602060405180830381865afa158015610158573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061017c9190611360565b90506000866001600160a01b0316635df459466040518163ffffffff1660e01b8152600401602060405180830381865afa1580156101be573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101e29190611360565b9050600086516001600160401b038111156101ff576101ff610f17565b60405190808252806020026020018201604052801561023257816020015b606081526020019060019003908161021d5790505b50905060005b875181101561053a5760008882815181106102555761025561137d565b0160200151604051638902624560e01b815260f89190911c6004820181905263ffffffff8a16602483015291506000906001600160a01b03871690638902624590604401600060405180830381865afa1580156102b6573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526102de91908101906113b6565b905080516001600160401b038111156102f9576102f9610f17565b60405190808252806020026020018201604052801561034457816020015b60408051606081018252600080825260208083018290529282015282526000199092019101816103175790505b508484815181106103575761035761137d565b602002602001018190525060005b8151811015610524576040518060600160405280876001600160a01b03166347b314e885858151811061039a5761039a61137d565b60200260200101516040518263ffffffff1660e01b81526004016103c091815260200190565b602060405180830381865afa1580156103dd573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104019190611360565b6001600160a01b031681526020018383815181106104215761042161137d565b60200260200101518152602001896001600160a01b031663fa28c62785858151811061044f5761044f61137d565b60209081029190910101516040516001600160e01b031960e084901b168152600481019190915260ff8816602482015263ffffffff8f166044820152606401602060405180830381865afa1580156104ab573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104cf919061144b565b6001600160601b03168152508585815181106104ed576104ed61137d565b602002602001015182815181106105065761050661137d565b6020026020010181905250808061051c9061148a565b915050610365565b50505080806105329061148a565b915050610238565b50979650505050505050565b6105716040518060800160405280606081526020016060815260200160608152602001606081525090565b6000876001600160a01b031663683048356040518163ffffffff1660e01b8152600401602060405180830381865afa1580156105b1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105d59190611360565b90506106026040518060800160405280606081526020016060815260200160608152602001606081525090565b6040516361c8a12f60e11b81526001600160a01b038a169063c391425e90610632908b90899089906004016114a5565b600060405180830381865afa15801561064f573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261067791908101906114ef565b81526040516340e03a8160e11b81526001600160a01b038316906381c07502906106a9908b908b908b906004016115a6565b600060405180830381865afa1580156106c6573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526106ee91908101906114ef565b6040820152856001600160401b0381111561070b5761070b610f17565b60405190808252806020026020018201604052801561073e57816020015b60608152602001906001900390816107295790505b50606082015260005b60ff8116871115610b81576000856001600160401b0381111561076c5761076c610f17565b604051908082528060200260200182016040528015610795578160200160208202803683370190505b5083606001518360ff16815181106107af576107af61137d565b602002602001018190525060005b86811015610a815760008c6001600160a01b03166304ec63518a8a858181106107e8576107e861137d565b905060200201358e886000015186815181106108065761080661137d565b60200260200101516040518463ffffffff1660e01b81526004016108439392919092835263ffffffff918216602084015216604082015260600190565b602060405180830381865afa158015610860573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061088491906115cf565b90506001600160c01b03811661092c5760405162461bcd60e51b815260206004820152605c60248201527f4f70657261746f7253746174655265747269657665722e676574436865636b5360448201527f69676e617475726573496e64696365733a206f70657261746f72206d7573742060648201527f6265207265676973746572656420617420626c6f636b6e756d62657200000000608482015260a40160405180910390fd5b8a8a8560ff168181106109415761094161137d565b6001600160c01b03841692013560f81c9190911c600190811614159050610a6e57856001600160a01b031663dd9846b98a8a858181106109835761098361137d565b905060200201358d8d8860ff1681811061099f5761099f61137d565b6040516001600160e01b031960e087901b1681526004810194909452919091013560f81c60248301525063ffffffff8f166044820152606401602060405180830381865afa1580156109f5573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a1991906115f8565b85606001518560ff1681518110610a3257610a3261137d565b60200260200101518481518110610a4b57610a4b61137d565b63ffffffff9092166020928302919091019091015282610a6a8161148a565b9350505b5080610a798161148a565b9150506107bd565b506000816001600160401b03811115610a9c57610a9c610f17565b604051908082528060200260200182016040528015610ac5578160200160208202803683370190505b50905060005b82811015610b465784606001518460ff1681518110610aec57610aec61137d565b60200260200101518181518110610b0557610b0561137d565b6020026020010151828281518110610b1f57610b1f61137d565b63ffffffff9092166020928302919091019091015280610b3e8161148a565b915050610acb565b508084606001518460ff1681518110610b6157610b6161137d565b602002602001018190525050508080610b7990611615565b915050610747565b506000896001600160a01b0316635df459466040518163ffffffff1660e01b8152600401602060405180830381865afa158015610bc2573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610be69190611360565b60405163354952a360e21b81529091506001600160a01b0382169063d5254a8c90610c19908b908b908e90600401611635565b600060405180830381865afa158015610c36573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610c5e91908101906114ef565b60208301525098975050505050505050565b6040805160018082528183019092526000916060918391602080830190803683370190505090508481600081518110610cab57610cab61137d565b60209081029190910101526040516361c8a12f60e11b81526000906001600160a01b0388169063c391425e90610ce7908890869060040161165f565b600060405180830381865afa158015610d04573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610d2c91908101906114ef565b600081518110610d3e57610d3e61137d565b60209081029190910101516040516304ec635160e01b81526004810188905263ffffffff87811660248301529091166044820181905291506000906001600160a01b038916906304ec635190606401602060405180830381865afa158015610daa573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610dce91906115cf565b6001600160c01b031690506000610de482610e02565b905081610df28a838a6100b0565b9550955050505050935093915050565b6060600080610e1084610ece565b61ffff166001600160401b03811115610e2b57610e2b610f17565b6040519080825280601f01601f191660200182016040528015610e55576020820181803683370190505b5090506000805b825182108015610e6d575061010081105b15610ec4576001811b935085841615610eb4578060f81b838381518110610e9657610e9661137d565b60200101906001600160f81b031916908160001a9053508160010191505b610ebd8161148a565b9050610e5c565b5090949350505050565b6000805b8215610ef957610ee36001846116b3565b9092169180610ef1816116ca565b915050610ed2565b92915050565b6001600160a01b0381168114610f1457600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b0381118282101715610f5557610f55610f17565b604052919050565b63ffffffff81168114610f1457600080fd5b8035610f7a81610f5d565b919050565b600080600060608486031215610f9457600080fd5b8335610f9f81610eff565b92506020848101356001600160401b0380821115610fbc57600080fd5b818701915087601f830112610fd057600080fd5b813581811115610fe257610fe2610f17565b610ff4601f8201601f19168501610f2d565b9150808252888482850101111561100a57600080fd5b808484018584013760008482840101525080945050505061102d60408501610f6f565b90509250925092565b600081518084526020808501808196508360051b810191508286016000805b868110156110cc578385038a52825180518087529087019087870190845b818110156110b757835180516001600160a01b031684528a8101518b8501526040908101516001600160601b03169084015292890192606090920191600101611073565b50509a87019a95505091850191600101611055565b509298975050505050505050565b6020815260006110ed6020830184611036565b9392505050565b60008083601f84011261110657600080fd5b5081356001600160401b0381111561111d57600080fd5b6020830191508360208260051b850101111561113857600080fd5b9250929050565b6000806000806000806080878903121561115857600080fd5b863561116381610eff565b9550602087013561117381610f5d565b945060408701356001600160401b038082111561118f57600080fd5b818901915089601f8301126111a357600080fd5b8135818111156111b257600080fd5b8a60208285010111156111c457600080fd5b6020830196508095505060608901359150808211156111e257600080fd5b506111ef89828a016110f4565b979a9699509497509295939492505050565b600081518084526020808501945080840160005b8381101561123757815163ffffffff1687529582019590820190600101611215565b509495945050505050565b60006020808352835160808285015261125e60a0850182611201565b905081850151601f198086840301604087015261127b8383611201565b925060408701519150808684030160608701526112988383611201565b60608801518782038301608089015280518083529194508501925084840190600581901b8501860160005b828110156112ef57848783030184526112dd828751611201565b958801959388019391506001016112c3565b509998505050505050505050565b60008060006060848603121561131257600080fd5b833561131d81610eff565b925060208401359150604084013561133481610f5d565b809150509250925092565b8281526040602082015260006113586040830184611036565b949350505050565b60006020828403121561137257600080fd5b81516110ed81610eff565b634e487b7160e01b600052603260045260246000fd5b60006001600160401b038211156113ac576113ac610f17565b5060051b60200190565b600060208083850312156113c957600080fd5b82516001600160401b038111156113df57600080fd5b8301601f810185136113f057600080fd5b80516114036113fe82611393565b610f2d565b81815260059190911b8201830190838101908783111561142257600080fd5b928401925b8284101561144057835182529284019290840190611427565b979650505050505050565b60006020828403121561145d57600080fd5b81516001600160601b03811681146110ed57600080fd5b634e487b7160e01b600052601160045260246000fd5b600060001982141561149e5761149e611474565b5060010190565b63ffffffff84168152604060208201819052810182905260006001600160fb1b038311156114d257600080fd5b8260051b8085606085013760009201606001918252509392505050565b6000602080838503121561150257600080fd5b82516001600160401b0381111561151857600080fd5b8301601f8101851361152957600080fd5b80516115376113fe82611393565b81815260059190911b8201830190838101908783111561155657600080fd5b928401925b8284101561144057835161156e81610f5d565b8252928401929084019061155b565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b63ffffffff841681526040602082015260006115c660408301848661157d565b95945050505050565b6000602082840312156115e157600080fd5b81516001600160c01b03811681146110ed57600080fd5b60006020828403121561160a57600080fd5b81516110ed81610f5d565b600060ff821660ff81141561162c5761162c611474565b60010192915050565b60408152600061164960408301858761157d565b905063ffffffff83166020830152949350505050565b60006040820163ffffffff851683526020604081850152818551808452606086019150828701935060005b818110156116a65784518352938301939183019160010161168a565b5090979650505050505050565b6000828210156116c5576116c5611474565b500390565b600061ffff808316818114156116e2576116e2611474565b600101939250505056fea26469706673582212205bfd930a1b7c8f8471802a2cf81cf4219840fbe6918482e6ae4f8270b52d86e864736f6c634300080c0033",
}

// ContractOperatorStateRetrieverABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractOperatorStateRetrieverMetaData.ABI instead.
var ContractOperatorStateRetrieverABI = ContractOperatorStateRetrieverMetaData.ABI

// ContractOperatorStateRetrieverBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractOperatorStateRetrieverMetaData.Bin instead.
var ContractOperatorStateRetrieverBin = ContractOperatorStateRetrieverMetaData.Bin

// DeployContractOperatorStateRetriever deploys a new Ethereum contract, binding an instance of ContractOperatorStateRetriever to it.
func DeployContractOperatorStateRetriever(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ContractOperatorStateRetriever, error) {
	parsed, err := ContractOperatorStateRetrieverMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractOperatorStateRetrieverBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractOperatorStateRetriever{ContractOperatorStateRetrieverCaller: ContractOperatorStateRetrieverCaller{contract: contract}, ContractOperatorStateRetrieverTransactor: ContractOperatorStateRetrieverTransactor{contract: contract}, ContractOperatorStateRetrieverFilterer: ContractOperatorStateRetrieverFilterer{contract: contract}}, nil
}

// ContractOperatorStateRetriever is an auto generated Go binding around an Ethereum contract.
type ContractOperatorStateRetriever struct {
	ContractOperatorStateRetrieverCaller     // Read-only binding to the contract
	ContractOperatorStateRetrieverTransactor // Write-only binding to the contract
	ContractOperatorStateRetrieverFilterer   // Log filterer for contract events
}

// ContractOperatorStateRetrieverCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractOperatorStateRetrieverCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractOperatorStateRetrieverTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractOperatorStateRetrieverTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractOperatorStateRetrieverFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractOperatorStateRetrieverFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractOperatorStateRetrieverSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractOperatorStateRetrieverSession struct {
	Contract     *ContractOperatorStateRetriever // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                   // Call options to use throughout this session
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// ContractOperatorStateRetrieverCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractOperatorStateRetrieverCallerSession struct {
	Contract *ContractOperatorStateRetrieverCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                         // Call options to use throughout this session
}

// ContractOperatorStateRetrieverTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractOperatorStateRetrieverTransactorSession struct {
	Contract     *ContractOperatorStateRetrieverTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                         // Transaction auth options to use throughout this session
}

// ContractOperatorStateRetrieverRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractOperatorStateRetrieverRaw struct {
	Contract *ContractOperatorStateRetriever // Generic contract binding to access the raw methods on
}

// ContractOperatorStateRetrieverCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractOperatorStateRetrieverCallerRaw struct {
	Contract *ContractOperatorStateRetrieverCaller // Generic read-only contract binding to access the raw methods on
}

// ContractOperatorStateRetrieverTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractOperatorStateRetrieverTransactorRaw struct {
	Contract *ContractOperatorStateRetrieverTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractOperatorStateRetriever creates a new instance of ContractOperatorStateRetriever, bound to a specific deployed contract.
func NewContractOperatorStateRetriever(address common.Address, backend bind.ContractBackend) (*ContractOperatorStateRetriever, error) {
	contract, err := bindContractOperatorStateRetriever(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractOperatorStateRetriever{ContractOperatorStateRetrieverCaller: ContractOperatorStateRetrieverCaller{contract: contract}, ContractOperatorStateRetrieverTransactor: ContractOperatorStateRetrieverTransactor{contract: contract}, ContractOperatorStateRetrieverFilterer: ContractOperatorStateRetrieverFilterer{contract: contract}}, nil
}

// NewContractOperatorStateRetrieverCaller creates a new read-only instance of ContractOperatorStateRetriever, bound to a specific deployed contract.
func NewContractOperatorStateRetrieverCaller(address common.Address, caller bind.ContractCaller) (*ContractOperatorStateRetrieverCaller, error) {
	contract, err := bindContractOperatorStateRetriever(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractOperatorStateRetrieverCaller{contract: contract}, nil
}

// NewContractOperatorStateRetrieverTransactor creates a new write-only instance of ContractOperatorStateRetriever, bound to a specific deployed contract.
func NewContractOperatorStateRetrieverTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractOperatorStateRetrieverTransactor, error) {
	contract, err := bindContractOperatorStateRetriever(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractOperatorStateRetrieverTransactor{contract: contract}, nil
}

// NewContractOperatorStateRetrieverFilterer creates a new log filterer instance of ContractOperatorStateRetriever, bound to a specific deployed contract.
func NewContractOperatorStateRetrieverFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractOperatorStateRetrieverFilterer, error) {
	contract, err := bindContractOperatorStateRetriever(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractOperatorStateRetrieverFilterer{contract: contract}, nil
}

// bindContractOperatorStateRetriever binds a generic wrapper to an already deployed contract.
func bindContractOperatorStateRetriever(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractOperatorStateRetrieverMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractOperatorStateRetriever.Contract.ContractOperatorStateRetrieverCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractOperatorStateRetriever.Contract.ContractOperatorStateRetrieverTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractOperatorStateRetriever.Contract.ContractOperatorStateRetrieverTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractOperatorStateRetriever.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractOperatorStateRetriever.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractOperatorStateRetriever.Contract.contract.Transact(opts, method, params...)
}

// GetCheckSignaturesIndices is a free data retrieval call binding the contract method 0x4f739f74.
//
// Solidity: function getCheckSignaturesIndices(address registryCoordinator, uint32 referenceBlockNumber, bytes quorumNumbers, bytes32[] nonSignerOperatorIds) view returns((uint32[],uint32[],uint32[],uint32[][]))
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverCaller) GetCheckSignaturesIndices(opts *bind.CallOpts, registryCoordinator common.Address, referenceBlockNumber uint32, quorumNumbers []byte, nonSignerOperatorIds [][32]byte) (OperatorStateRetrieverCheckSignaturesIndices, error) {
	var out []interface{}
	err := _ContractOperatorStateRetriever.contract.Call(opts, &out, "getCheckSignaturesIndices", registryCoordinator, referenceBlockNumber, quorumNumbers, nonSignerOperatorIds)

	if err != nil {
		return *new(OperatorStateRetrieverCheckSignaturesIndices), err
	}

	out0 := *abi.ConvertType(out[0], new(OperatorStateRetrieverCheckSignaturesIndices)).(*OperatorStateRetrieverCheckSignaturesIndices)

	return out0, err

}

// GetCheckSignaturesIndices is a free data retrieval call binding the contract method 0x4f739f74.
//
// Solidity: function getCheckSignaturesIndices(address registryCoordinator, uint32 referenceBlockNumber, bytes quorumNumbers, bytes32[] nonSignerOperatorIds) view returns((uint32[],uint32[],uint32[],uint32[][]))
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverSession) GetCheckSignaturesIndices(registryCoordinator common.Address, referenceBlockNumber uint32, quorumNumbers []byte, nonSignerOperatorIds [][32]byte) (OperatorStateRetrieverCheckSignaturesIndices, error) {
	return _ContractOperatorStateRetriever.Contract.GetCheckSignaturesIndices(&_ContractOperatorStateRetriever.CallOpts, registryCoordinator, referenceBlockNumber, quorumNumbers, nonSignerOperatorIds)
}

// GetCheckSignaturesIndices is a free data retrieval call binding the contract method 0x4f739f74.
//
// Solidity: function getCheckSignaturesIndices(address registryCoordinator, uint32 referenceBlockNumber, bytes quorumNumbers, bytes32[] nonSignerOperatorIds) view returns((uint32[],uint32[],uint32[],uint32[][]))
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverCallerSession) GetCheckSignaturesIndices(registryCoordinator common.Address, referenceBlockNumber uint32, quorumNumbers []byte, nonSignerOperatorIds [][32]byte) (OperatorStateRetrieverCheckSignaturesIndices, error) {
	return _ContractOperatorStateRetriever.Contract.GetCheckSignaturesIndices(&_ContractOperatorStateRetriever.CallOpts, registryCoordinator, referenceBlockNumber, quorumNumbers, nonSignerOperatorIds)
}

// GetOperatorState is a free data retrieval call binding the contract method 0x3563b0d1.
//
// Solidity: function getOperatorState(address registryCoordinator, bytes quorumNumbers, uint32 blockNumber) view returns((address,bytes32,uint96)[][])
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverCaller) GetOperatorState(opts *bind.CallOpts, registryCoordinator common.Address, quorumNumbers []byte, blockNumber uint32) ([][]OperatorStateRetrieverOperator, error) {
	var out []interface{}
	err := _ContractOperatorStateRetriever.contract.Call(opts, &out, "getOperatorState", registryCoordinator, quorumNumbers, blockNumber)

	if err != nil {
		return *new([][]OperatorStateRetrieverOperator), err
	}

	out0 := *abi.ConvertType(out[0], new([][]OperatorStateRetrieverOperator)).(*[][]OperatorStateRetrieverOperator)

	return out0, err

}

// GetOperatorState is a free data retrieval call binding the contract method 0x3563b0d1.
//
// Solidity: function getOperatorState(address registryCoordinator, bytes quorumNumbers, uint32 blockNumber) view returns((address,bytes32,uint96)[][])
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverSession) GetOperatorState(registryCoordinator common.Address, quorumNumbers []byte, blockNumber uint32) ([][]OperatorStateRetrieverOperator, error) {
	return _ContractOperatorStateRetriever.Contract.GetOperatorState(&_ContractOperatorStateRetriever.CallOpts, registryCoordinator, quorumNumbers, blockNumber)
}

// GetOperatorState is a free data retrieval call binding the contract method 0x3563b0d1.
//
// Solidity: function getOperatorState(address registryCoordinator, bytes quorumNumbers, uint32 blockNumber) view returns((address,bytes32,uint96)[][])
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverCallerSession) GetOperatorState(registryCoordinator common.Address, quorumNumbers []byte, blockNumber uint32) ([][]OperatorStateRetrieverOperator, error) {
	return _ContractOperatorStateRetriever.Contract.GetOperatorState(&_ContractOperatorStateRetriever.CallOpts, registryCoordinator, quorumNumbers, blockNumber)
}

// GetOperatorState0 is a free data retrieval call binding the contract method 0xcefdc1d4.
//
// Solidity: function getOperatorState(address registryCoordinator, bytes32 operatorId, uint32 blockNumber) view returns(uint256, (address,bytes32,uint96)[][])
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverCaller) GetOperatorState0(opts *bind.CallOpts, registryCoordinator common.Address, operatorId [32]byte, blockNumber uint32) (*big.Int, [][]OperatorStateRetrieverOperator, error) {
	var out []interface{}
	err := _ContractOperatorStateRetriever.contract.Call(opts, &out, "getOperatorState0", registryCoordinator, operatorId, blockNumber)

	if err != nil {
		return *new(*big.Int), *new([][]OperatorStateRetrieverOperator), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new([][]OperatorStateRetrieverOperator)).(*[][]OperatorStateRetrieverOperator)

	return out0, out1, err

}

// GetOperatorState0 is a free data retrieval call binding the contract method 0xcefdc1d4.
//
// Solidity: function getOperatorState(address registryCoordinator, bytes32 operatorId, uint32 blockNumber) view returns(uint256, (address,bytes32,uint96)[][])
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverSession) GetOperatorState0(registryCoordinator common.Address, operatorId [32]byte, blockNumber uint32) (*big.Int, [][]OperatorStateRetrieverOperator, error) {
	return _ContractOperatorStateRetriever.Contract.GetOperatorState0(&_ContractOperatorStateRetriever.CallOpts, registryCoordinator, operatorId, blockNumber)
}

// GetOperatorState0 is a free data retrieval call binding the contract method 0xcefdc1d4.
//
// Solidity: function getOperatorState(address registryCoordinator, bytes32 operatorId, uint32 blockNumber) view returns(uint256, (address,bytes32,uint96)[][])
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverCallerSession) GetOperatorState0(registryCoordinator common.Address, operatorId [32]byte, blockNumber uint32) (*big.Int, [][]OperatorStateRetrieverOperator, error) {
	return _ContractOperatorStateRetriever.Contract.GetOperatorState0(&_ContractOperatorStateRetriever.CallOpts, registryCoordinator, operatorId, blockNumber)
}
