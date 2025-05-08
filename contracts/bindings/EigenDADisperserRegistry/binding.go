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
	ABI: "[{\"type\":\"function\",\"name\":\"deregisterDisperser\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getDepositParams\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV3.LockedDisperserDeposit\",\"components\":[{\"name\":\"deposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"refund\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lockPeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDisperserInfo\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV3.DisperserInfo\",\"components\":[{\"name\":\"disperser\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"registered\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"disperserURL\",\"type\":\"string\",\"internalType\":\"string\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLockedDeposit\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV3.LockedDisperserDeposit\",\"components\":[{\"name\":\"deposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"refund\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lockPeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"unlockTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"depositParams\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV3.LockedDisperserDeposit\",\"components\":[{\"name\":\"deposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"refund\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lockPeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerDisperser\",\"inputs\":[{\"name\":\"disperserAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"disperserURL\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDepositParams\",\"inputs\":[{\"name\":\"depositParams\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV3.LockedDisperserDeposit\",\"components\":[{\"name\":\"deposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"refund\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lockPeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferDisperserOwnership\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"disperserAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateDisperserURL\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"disperserURL\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"error\",\"name\":\"AlreadyInitialized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MissingRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x608060405234801561001057600080fd5b50611b2e806100206000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c8063a450815811610071578063a450815814610138578063aa9224cd1461014b578063bb4c7dda1461015e578063d5f3e3f214610171578063dc4767c2146101d1578063f2fde38b146101e457600080fd5b80631b4b940a146100ae5780633fcf1fb4146100db5780634fe5de32146100fb5780636923ef3a146101105780638e1e482914610123575b600080fd5b6100c16100bc3660046116bd565b6101f7565b60405163ffffffff90911681526020015b60405180910390f35b6100ee6100e936600461171e565b61020c565b6040516100d29190611791565b61010e61010936600461171e565b610231565b005b61010e61011e366004611850565b610281565b61012b61028d565b6040516100d2919061186c565b61010e6101463660046118a9565b6102a2565b61010e61015936600461171e565b6102df565b61010e61016c3660046118dd565b610322565b61018461017f36600461171e565b61036b565b604080518351815260208085015190820152838201516001600160a01b0316918101919091526060928301516001600160401b03908116938201939093529116608082015260a0016100d2565b61010e6101df366004611907565b610387565b61010e6101f2366004611923565b6103cb565b60006102038383610420565b90505b92915050565b6040805160608082018352600080835260208301529181019190915261020682610627565b8061023b81610728565b6001600160a01b0316336001600160a01b0316146102745760405162461bcd60e51b815260040161026b9061193e565b60405180910390fd5b61027d82610757565b5050565b61028a81610819565b50565b61029561152c565b61029d61091f565b905090565b6102ac600161097b565b6102d67fe579393920b888b1e4a7e1afdd7d58fa4f3101113547ac874aefa75ff4a960f9836109f7565b61027d81610819565b806102e981610728565b6001600160a01b0316336001600160a01b0316146103195760405162461bcd60e51b815260040161026b9061193e565b61027d82610a53565b8161032c81610728565b6001600160a01b0316336001600160a01b03161461035c5760405162461bcd60e51b815260040161026b9061193e565b6103668383610b73565b505050565b61037361152c565b600061037e83610cb8565b91509150915091565b8161039181610728565b6001600160a01b0316336001600160a01b0316146103c15760405162461bcd60e51b815260040161026b9061193e565b6103668383610d3d565b6103f57fe579393920b888b1e4a7e1afdd7d58fa4f3101113547ac874aefa75ff4a960f933610e2a565b61028a7fe579393920b888b1e4a7e1afdd7d58fa4f3101113547ac874aefa75ff4a960f93383610e63565b600061042a610e77565b90506000610436610ed3565b63ffffffff831660009081526020919091526040902090506001600160a01b0384166104a05760405162461bcd60e51b8152602060048201526019602482015278496e76616c696420646973706572736572206164647265737360381b604482015260640161026b565b60006104aa610ed3565b600401546001600160a01b0316905060006104c3610ed3565b60020154905060006104d3610ed3565b6003015490508115610541576104f46001600160a01b038416333085610edd565b6104fe818361198b565b610506610ed3565b6001016000856001600160a01b03166001600160a01b03168152602001908152602001600020600082825461053b91906119a2565b90915550505b604080516060810182526001600160a01b038916808252600160208084018290529383018a905287546001600160a81b031916909117600160a01b178755885191928792610594928401918b0190611566565b509050506105a0610ed3565b60028181015490860155600380820154908601556004908101805491860180546001600160a01b039093166001600160a01b031984168117825591546001600160e01b0319909316909117600160a01b928390046001600160401b039081169093021790556005909401805467ffffffffffffffff19169094179093555091949350505050565b6040805160608082018352600080835260208301529181019190915261064b610ed3565b63ffffffff831660009081526020918252604090819020815160608101835281546001600160a01b0381168252600160a01b900460ff16151593810193909352600181018054919284019161069f906119ba565b80601f01602080910402602001604051908101604052809291908181526020018280546106cb906119ba565b80156107185780601f106106ed57610100808354040283529160200191610718565b820191906000526020600020905b8154815290600101906020018083116106fb57829003601f168201915b5050505050815250509050919050565b6000610732610ed3565b63ffffffff90921660009081526020929092525060409020546001600160a01b031690565b6000610761610ed3565b63ffffffff83166000908152602091909152604081209150610781610ed3565b63ffffffff8416600090815260209190915260409020825460029091019150600160a01b900460ff166107c65760405162461bcd60e51b815260040161026b906119f5565b815460ff60a01b1916825560028101546107f090600160a01b90046001600160401b031642611a2c565b600592909201805467ffffffffffffffff19166001600160401b03909316929092179091555050565b60208101518151101561086e5760405162461bcd60e51b815260206004820152601f60248201527f4465706f736974206d757374206265206174206c6561737420726566756e6400604482015260640161026b565b60408101516001600160a01b03166108c05760405162461bcd60e51b8152602060048201526015602482015274496e76616c696420746f6b656e206164647265737360581b604482015260640161026b565b806108c9610ed3565b81516002820155602082015160038201556040820151600490910180546060909301516001600160401b0316600160a01b026001600160e01b03199093166001600160a01b039092169190911791909117905550565b61092761152c565b61092f610ed3565b6040805160808101825260028301548152600383015460208201526004909201546001600160a01b03811691830191909152600160a01b90046001600160401b03166060820152919050565b8060ff16610987610f4e565b5460ff16106109a85760405162dc149f60e41b815260040160405180910390fd5b806109b1610f4e565b805460ff191660ff92831617905560405190821681527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a150565b610a1881610a03610f58565b60008581526020919091526040902090610f62565b506040516001600160a01b0382169083907f2ae6a113c0ed5b78a53413ffbb7679881f11145ccfba4fb92e863dfcd5a1d2f390600090a35050565b6000610a5d610ed3565b63ffffffff83166000908152602091909152604081209150610a7d610ed3565b63ffffffff8416600090815260209190915260409020600381015460029091019150610ae45760405162461bcd60e51b81526020600482015260166024820152754e6f206465706f73697420746f20776974686472617760501b604482015260640161026b565b6005820154426001600160401b039091161115610b435760405162461bcd60e51b815260206004820152601760248201527f4465706f736974206973207374696c6c206c6f636b6564000000000000000000604482015260640161026b565b815460018201546002830154610b67926001600160a01b0391821692911690610f77565b60006001909101555050565b6000610b7d610ed3565b63ffffffff84166000908152602091909152604081209150610b9d610ed3565b600401546001600160a01b031690506000610bb6610ed3565b6005015490508015610c1b57610bd76001600160a01b038316333084610edd565b80610be0610ed3565b6001016000846001600160a01b03166001600160a01b031681526020019081526020016000206000828254610c1591906119a2565b90915550505b8254600160a01b900460ff16610c435760405162461bcd60e51b815260040161026b906119f5565b6001600160a01b038416610c955760405162461bcd60e51b8152602060048201526019602482015278496e76616c696420646973706572736572206164647265737360381b604482015260640161026b565b505080546001600160a01b0319166001600160a01b039290921691909117905550565b610cc061152c565b600080610ccb610ed3565b63ffffffff9094166000908152602094855260409081902060058101548251608081018452600283015481526003830154978101979097526004909101546001600160a01b03811692870192909252600160a01b9091046001600160401b039081166060870152949594169392505050565b6000610d47610ed3565b63ffffffff84166000908152602091909152604081209150610d67610ed3565b600401546001600160a01b031690506000610d80610ed3565b6005015490508015610de557610da16001600160a01b038316333084610edd565b80610daa610ed3565b6001016000846001600160a01b03166001600160a01b031681526020019081526020016000206000828254610ddf91906119a2565b90915550505b8254600160a01b900460ff16610e0d5760405162461bcd60e51b815260040161026b906119f5565b8351610e229060018501906020870190611566565b505050505050565b610e348282610fa7565b61027d576040516301d4003760e61b8152600481018390526001600160a01b038216602482015260440161026b565b610e6d8383610fca565b61036683826109f7565b600080610e82610ed3565b6006015463ffffffff169050610e96610ed3565b600601805463ffffffff16906000610ead83611a57565b91906101000a81548163ffffffff021916908363ffffffff160217905550508091505090565b600061029d611026565b6040516001600160a01b0380851660248301528316604482015260648101829052610f489085906323b872dd60e01b906084015b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b0319909316929092179091526110cb565b50505050565b600061029d61119d565b600061029d6111e6565b6000610203836001600160a01b038416611230565b6040516001600160a01b03831660248201526044810182905261036690849063a9059cbb60e01b90606401610f11565b600061020382610fb5610f58565b6000868152602091909152604090209061127f565b610feb81610fd6610f58565b600085815260209190915260409020906112a1565b506040516001600160a01b0382169083907f155aaafb6329a2098580462df33ec4b7441b19729b9601c5fc17ae1cf99a8a5290600090a35050565b60008060ff60001b1960016040518060400160405280601b81526020017f656967656e2e64612e6469737065727365722e726567697374727900000000008152506040516020016110779190611a7b565b6040516020818303038152906040528051906020012060001c61109a919061198b565b6040516020016110ac91815260200190565b60408051601f1981840301815291905280516020909101201692915050565b6000611120826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564815250856001600160a01b03166112b69092919063ffffffff16565b805190915015610366578080602001905181019061113e9190611a97565b6103665760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6044820152691bdd081cdd58d8d9595960b21b606482015260840161026b565b60008060ff60001b19600160405180604001604052806015815260200174696e697469616c697a61626c652e73746f7261676560581b8152506040516020016110779190611a7b565b60008060ff60001b196001604051806040016040528060168152602001756163636573732e636f6e74726f6c2e73746f7261676560501b8152506040516020016110779190611a7b565b600081815260018301602052604081205461127757508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610206565b506000610206565b6001600160a01b03811660009081526001830160205260408120541515610203565b6000610203836001600160a01b0384166112cf565b60606112c584846000856113c2565b90505b9392505050565b600081815260018301602052604081205480156113b85760006112f360018361198b565b85549091506000906113079060019061198b565b905081811461136c57600086600001828154811061132757611327611ab9565b906000526020600020015490508087600001848154811061134a5761134a611ab9565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061137d5761137d611acf565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610206565b6000915050610206565b6060824710156114235760405162461bcd60e51b815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f6044820152651c8818d85b1b60d21b606482015260840161026b565b6001600160a01b0385163b61147a5760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640161026b565b600080866001600160a01b031685876040516114969190611a7b565b60006040518083038185875af1925050503d80600081146114d3576040519150601f19603f3d011682016040523d82523d6000602084013e6114d8565b606091505b50915091506114e88282866114f3565b979650505050505050565b606083156115025750816112c8565b8251156115125782518084602001fd5b8160405162461bcd60e51b815260040161026b9190611ae5565b6040518060800160405280600081526020016000815260200160006001600160a01b0316815260200160006001600160401b031681525090565b828054611572906119ba565b90600052602060002090601f01602090048101928261159457600085556115da565b82601f106115ad57805160ff19168380011785556115da565b828001600101855582156115da579182015b828111156115da5782518255916020019190600101906115bf565b506115e69291506115ea565b5090565b5b808211156115e657600081556001016115eb565b80356001600160a01b038116811461161657600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b600082601f83011261164257600080fd5b81356001600160401b038082111561165c5761165c61161b565b604051601f8301601f19908116603f011681019082821181831017156116845761168461161b565b8160405283815286602085880101111561169d57600080fd5b836020870160208301376000602085830101528094505050505092915050565b600080604083850312156116d057600080fd5b6116d9836115ff565b915060208301356001600160401b038111156116f457600080fd5b61170085828601611631565b9150509250929050565b803563ffffffff8116811461161657600080fd5b60006020828403121561173057600080fd5b6102038261170a565b60005b8381101561175457818101518382015260200161173c565b83811115610f485750506000910152565b6000815180845261177d816020860160208601611739565b601f01601f19169290920160200192915050565b6020815260018060a01b038251166020820152602082015115156040820152600060408301516060808401526117ca6080840182611765565b949350505050565b6000608082840312156117e457600080fd5b604051608081016001600160401b0382821081831117156118075761180761161b565b816040528293508435835260208501356020840152611828604086016115ff565b604084015260608501359150808216821461184257600080fd5b506060919091015292915050565b60006080828403121561186257600080fd5b61020383836117d2565b81518152602080830151908201526040808301516001600160a01b0316908201526060808301516001600160401b03169082015260808101610206565b60008060a083850312156118bc57600080fd5b6118c5836115ff565b91506118d484602085016117d2565b90509250929050565b600080604083850312156118f057600080fd5b6118f98361170a565b91506118d4602084016115ff565b6000806040838503121561191a57600080fd5b6116d98361170a565b60006020828403121561193557600080fd5b610203826115ff565b6020808252601b908201527f43616c6c6572206973206e6f7420746865206469737065727365720000000000604082015260600190565b634e487b7160e01b600052601160045260246000fd5b60008282101561199d5761199d611975565b500390565b600082198211156119b5576119b5611975565b500190565b600181811c908216806119ce57607f821691505b602082108114156119ef57634e487b7160e01b600052602260045260246000fd5b50919050565b60208082526018908201527f446973706572736572206e6f7420726567697374657265640000000000000000604082015260600190565b60006001600160401b03808316818516808303821115611a4e57611a4e611975565b01949350505050565b600063ffffffff80831681811415611a7157611a71611975565b6001019392505050565b60008251611a8d818460208701611739565b9190910192915050565b600060208284031215611aa957600080fd5b815180151581146112c857600080fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052603160045260246000fd5b602081526000610203602083018461176556fea26469706673582212204f54fd422c1a7cdc930a069969cdc973b13098b3917b06782f99c9729400624364736f6c634300080c0033",
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

// TransferDisperserOwnership is a paid mutator transaction binding the contract method 0xbb4c7dda.
//
// Solidity: function transferDisperserOwnership(uint32 disperserKey, address disperserAddress) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryTransactor) TransferDisperserOwnership(opts *bind.TransactOpts, disperserKey uint32, disperserAddress common.Address) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.contract.Transact(opts, "transferDisperserOwnership", disperserKey, disperserAddress)
}

// TransferDisperserOwnership is a paid mutator transaction binding the contract method 0xbb4c7dda.
//
// Solidity: function transferDisperserOwnership(uint32 disperserKey, address disperserAddress) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistrySession) TransferDisperserOwnership(disperserKey uint32, disperserAddress common.Address) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.TransferDisperserOwnership(&_ContractEigenDADisperserRegistry.TransactOpts, disperserKey, disperserAddress)
}

// TransferDisperserOwnership is a paid mutator transaction binding the contract method 0xbb4c7dda.
//
// Solidity: function transferDisperserOwnership(uint32 disperserKey, address disperserAddress) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryTransactorSession) TransferDisperserOwnership(disperserKey uint32, disperserAddress common.Address) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.TransferDisperserOwnership(&_ContractEigenDADisperserRegistry.TransactOpts, disperserKey, disperserAddress)
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

// UpdateDisperserURL is a paid mutator transaction binding the contract method 0xdc4767c2.
//
// Solidity: function updateDisperserURL(uint32 disperserKey, string disperserURL) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryTransactor) UpdateDisperserURL(opts *bind.TransactOpts, disperserKey uint32, disperserURL string) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.contract.Transact(opts, "updateDisperserURL", disperserKey, disperserURL)
}

// UpdateDisperserURL is a paid mutator transaction binding the contract method 0xdc4767c2.
//
// Solidity: function updateDisperserURL(uint32 disperserKey, string disperserURL) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistrySession) UpdateDisperserURL(disperserKey uint32, disperserURL string) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.UpdateDisperserURL(&_ContractEigenDADisperserRegistry.TransactOpts, disperserKey, disperserURL)
}

// UpdateDisperserURL is a paid mutator transaction binding the contract method 0xdc4767c2.
//
// Solidity: function updateDisperserURL(uint32 disperserKey, string disperserURL) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryTransactorSession) UpdateDisperserURL(disperserKey uint32, disperserURL string) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.UpdateDisperserURL(&_ContractEigenDADisperserRegistry.TransactOpts, disperserKey, disperserURL)
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
