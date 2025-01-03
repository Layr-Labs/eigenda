// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractPaymentVault

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

// IPaymentVaultReservation is an auto generated low-level Go binding around an user-defined struct.
type IPaymentVaultReservation struct {
	SymbolsPerSecond uint64
	StartTimestamp   uint64
	EndTimestamp     uint64
	QuorumNumbers    []byte
	QuorumSplits     []byte
}

// ContractPaymentVaultMetaData contains all meta data concerning the ContractPaymentVault contract.
var ContractPaymentVaultMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"depositOnDemand\",\"inputs\":[{\"name\":\"_account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"getOnDemandTotalDeposit\",\"inputs\":[{\"name\":\"_account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOnDemandTotalDeposits\",\"inputs\":[{\"name\":\"_accounts\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[{\"name\":\"_payments\",\"type\":\"uint80[]\",\"internalType\":\"uint80[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReservation\",\"inputs\":[{\"name\":\"_account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIPaymentVault.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReservations\",\"inputs\":[{\"name\":\"_accounts\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[{\"name\":\"_reservations\",\"type\":\"tuple[]\",\"internalType\":\"structIPaymentVault.Reservation[]\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"globalRatePeriodInterval\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"globalSymbolsPerPeriod\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_minNumSymbols\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_pricePerSymbol\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_priceUpdateCooldown\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_globalSymbolsPerPeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_reservationPeriodInterval\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_globalRatePeriodInterval\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"lastPriceUpdateTime\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minNumSymbols\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"onDemandPayments\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"totalDeposit\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pricePerSymbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"priceUpdateCooldown\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"reservationPeriodInterval\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"reservations\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setGlobalRatePeriodInterval\",\"inputs\":[{\"name\":\"_globalRatePeriodInterval\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setGlobalSymbolsPerPeriod\",\"inputs\":[{\"name\":\"_globalSymbolsPerPeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPriceParams\",\"inputs\":[{\"name\":\"_minNumSymbols\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_pricePerSymbol\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_priceUpdateCooldown\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setReservation\",\"inputs\":[{\"name\":\"_account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_reservation\",\"type\":\"tuple\",\"internalType\":\"structIPaymentVault.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setReservationPeriodInterval\",\"inputs\":[{\"name\":\"_reservationPeriodInterval\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawERC20\",\"inputs\":[{\"name\":\"_token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"GlobalRatePeriodIntervalUpdated\",\"inputs\":[{\"name\":\"previousValue\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"newValue\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GlobalSymbolsPerPeriodUpdated\",\"inputs\":[{\"name\":\"previousValue\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"newValue\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OnDemandPaymentUpdated\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"onDemandPayment\",\"type\":\"uint80\",\"indexed\":false,\"internalType\":\"uint80\"},{\"name\":\"totalDeposit\",\"type\":\"uint80\",\"indexed\":false,\"internalType\":\"uint80\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PriceParamsUpdated\",\"inputs\":[{\"name\":\"previousMinNumSymbols\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"newMinNumSymbols\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"previousPricePerSymbol\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"newPricePerSymbol\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"previousPriceUpdateCooldown\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"newPriceUpdateCooldown\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReservationPeriodIntervalUpdated\",\"inputs\":[{\"name\":\"previousValue\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"newValue\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReservationUpdated\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"reservation\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIPaymentVault.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"anonymous\":false}]",
	Bin: "0x608060405234801561001057600080fd5b5061001961001e565b6100de565b600054610100900460ff161561008a5760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840160405180910390fd5b60005460ff90811610156100dc576000805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b611e08806100ed6000396000f3fe60806040526004361061016a5760003560e01c80639aec8640116100d1578063c98d97dd1161008a578063f2fde38b11610064578063f2fde38b146104c2578063f323726a146104e2578063fba2b1d114610509578063fd3dc53a146105295761017b565b8063c98d97dd14610415578063d1c1fdcd14610435578063d996dc991461048c5761017b565b80639aec864014610341578063a16cf88414610361578063a1db978214610381578063aa788bd7146103a1578063b2066f80146103c1578063bff8a3d4146103ee5761017b565b806372228ab21161012357806372228ab21461027f578063761dab89146102a6578063897218fc146102c65780638bec7d02146102e65780638da5cb5b146102f95780639a1bbf37146103215761017b565b8063039f091c14610185578063109f8fe5146101c95780632e1a7d4d146101f65780634184a6741461021657806349b9a7af14610243578063715018a61461026a5761017b565b3661017b57610179333461055a565b005b610179333461055a565b34801561019157600080fd5b506065546101ac90600160801b90046001600160401b031681565b6040516001600160401b0390911681526020015b60405180910390f35b3480156101d557600080fd5b506101e96101e43660046117b8565b610679565b6040516101c09190611915565b34801561020257600080fd5b50610179610211366004611977565b6108cc565b34801561022257600080fd5b506102366102313660046117b8565b610949565b6040516101c09190611990565b34801561024f57600080fd5b506065546101ac90600160c01b90046001600160401b031681565b34801561027657600080fd5b50610179610a2b565b34801561028b57600080fd5b506066546101ac90600160401b90046001600160401b031681565b3480156102b257600080fd5b506065546101ac906001600160401b031681565b3480156102d257600080fd5b506101796102e13660046119f4565b610a3f565b6101796102f4366004611a16565b610ac7565b34801561030557600080fd5b506033546040516001600160a01b0390911681526020016101c0565b34801561032d57600080fd5b5061017961033c366004611a31565b610ad4565b34801561034d57600080fd5b5061017961035c366004611b26565b610cfe565b34801561036d57600080fd5b5061017961037c3660046119f4565b610e84565b34801561038d57600080fd5b5061017961039c366004611bfc565b610ef6565b3480156103ad57600080fd5b506101796103bc3660046119f4565b610f95565b3480156103cd57600080fd5b506103e16103dc366004611a16565b611018565b6040516101c09190611c26565b3480156103fa57600080fd5b506066546101ac90600160801b90046001600160401b031681565b34801561042157600080fd5b506066546101ac906001600160401b031681565b34801561044157600080fd5b50610474610450366004611a16565b6001600160a01b03166000908152606860205260409020546001600160501b031690565b6040516001600160501b0390911681526020016101c0565b34801561049857600080fd5b506104746104a7366004611a16565b6068602052600090815260409020546001600160501b031681565b3480156104ce57600080fd5b506101796104dd366004611a16565b6111c2565b3480156104ee57600080fd5b506065546101ac90600160401b90046001600160401b031681565b34801561051557600080fd5b50610179610524366004611c39565b611238565b34801561053557600080fd5b50610549610544366004611a16565b6113ae565b6040516101c0959493929190611c7c565b6001600160501b038111156105cb5760405162461bcd60e51b815260206004820152602c60248201527f616d6f756e74206d757374206265206c657373207468616e206f72206571756160448201526b6c20746f203830206269747360a01b60648201526084015b60405180910390fd5b6001600160a01b038216600090815260686020526040812080548392906105fc9084906001600160501b0316611cd8565b82546101009290920a6001600160501b038181021990931691831602179091556001600160a01b03841660008181526068602090815260409182902054825187861681529416908401529092507f6fbb447a2c09b8901d70b0d5b9fbce159ee8fda4460e5af2570cab3fe0adf26891015b60405180910390a25050565b606081516001600160401b038111156106945761069461172e565b6040519080825280602002602001820160405280156106ec57816020015b6040805160a08101825260008082526020808301829052928201526060808201819052608082015282526000199092019101816106b25790505b50905060005b82518110156108c6576067600084838151811061071157610711611d03565b6020908102919091018101516001600160a01b03168252818101929092526040908101600020815160a08101835281546001600160401b038082168352600160401b8204811695830195909552600160801b90049093169183019190915260018101805460608401919061078490611d19565b80601f01602080910402602001604051908101604052809291908181526020018280546107b090611d19565b80156107fd5780601f106107d2576101008083540402835291602001916107fd565b820191906000526020600020905b8154815290600101906020018083116107e057829003601f168201915b5050505050815260200160028201805461081690611d19565b80601f016020809104026020016040519081016040528092919081815260200182805461084290611d19565b801561088f5780601f106108645761010080835404028352916020019161088f565b820191906000526020600020905b81548152906001019060200180831161087257829003601f168201915b5050505050815250508282815181106108aa576108aa611d03565b6020026020010181905250806108bf90611d4e565b90506106f2565b50919050565b6108d4611501565b60006108e86033546001600160a01b031690565b6001600160a01b03168260405160006040518083038185875af1925050503d8060008114610932576040519150601f19603f3d011682016040523d82523d6000602084013e610937565b606091505b505090508061094557600080fd5b5050565b606081516001600160401b038111156109645761096461172e565b60405190808252806020026020018201604052801561098d578160200160208202803683370190505b50905060005b82518110156108c657606860008483815181106109b2576109b2611d03565b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060000160009054906101000a90046001600160501b0316828281518110610a0357610a03611d03565b6001600160501b0390921660209283029190910190910152610a2481611d4e565b9050610993565b610a33611501565b610a3d600061155b565b565b610a47611501565b606654604080516001600160401b03600160401b9093048316815291831660208301527f1ef4a1ce7d8e50959d15578b346bb20a5b049e5ee1978014a4ba66476265c957910160405180910390a1606680546001600160401b03909216600160401b026fffffffffffffffff000000000000000019909216919091179055565b610ad1813461055a565b50565b600054610100900460ff1615808015610af45750600054600160ff909116105b80610b0e5750303b158015610b0e575060005460ff166001145b610b715760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b60648201526084016105c2565b6000805460ff191660011790558015610b94576000805461ff0019166101001790555b610b9d8861155b565b86606560006101000a8154816001600160401b0302191690836001600160401b0316021790555085606560086101000a8154816001600160401b0302191690836001600160401b0316021790555084606560106101000a8154816001600160401b0302191690836001600160401b0316021790555042606560186101000a8154816001600160401b0302191690836001600160401b0316021790555083606660006101000a8154816001600160401b0302191690836001600160401b0316021790555082606660086101000a8154816001600160401b0302191690836001600160401b0316021790555081606660106101000a8154816001600160401b0302191690836001600160401b031602179055508015610cf4576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b5050505050505050565b610d06611501565b610d18816060015182608001516115ad565b80602001516001600160401b031681604001516001600160401b031611610d9c5760405162461bcd60e51b815260206004820152603260248201527f656e642074696d657374616d70206d75737420626520677265617465722074686044820152710616e2073746172742074696d657374616d760741b60648201526084016105c2565b6001600160a01b0382166000908152606760209081526040918290208351815483860151948601516001600160401b03908116600160801b0267ffffffffffffffff60801b19968216600160401b026fffffffffffffffffffffffffffffffff199093169190931617179390931692909217825560608301518051849392610e2b926001850192910190611695565b5060808201518051610e47916002840191602090910190611695565b50905050816001600160a01b03167fff3054d138559c39b4c0826c43e94b2b2c6bc9a33ea1d0b74f16c916c7b73ec18260405161066d9190611c26565b610e8c611501565b606654604080516001600160401b03928316815291831660208301527f3edf3b79e74d9e583ff51df95fbabefe15f504d33475b2cc77cffba292268aae910160405180910390a16066805467ffffffffffffffff19166001600160401b0392909216919091179055565b610efe611501565b816001600160a01b031663a9059cbb610f1f6033546001600160a01b031690565b6040516001600160e01b031960e084901b1681526001600160a01b039091166004820152602481018490526044016020604051808303816000875af1158015610f6c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f909190611d69565b505050565b610f9d611501565b606654604080516001600160401b03600160801b9093048316815291831660208301527f833819c38214ef9f462f88b5c27a21bf201f394572a14da3e63c77ee15f0e93a910160405180910390a1606680546001600160401b03909216600160801b0267ffffffffffffffff60801b19909216919091179055565b6040805160a08082018352600080835260208084018290528385018290526060808501819052608085018190526001600160a01b038716835260678252918590208551938401865280546001600160401b038082168652600160401b8204811693860193909352600160801b90049091169483019490945260018401805493949293918401916110a790611d19565b80601f01602080910402602001604051908101604052809291908181526020018280546110d390611d19565b80156111205780601f106110f557610100808354040283529160200191611120565b820191906000526020600020905b81548152906001019060200180831161110357829003601f168201915b5050505050815260200160028201805461113990611d19565b80601f016020809104026020016040519081016040528092919081815260200182805461116590611d19565b80156111b25780601f10611187576101008083540402835291602001916111b2565b820191906000526020600020905b81548152906001019060200180831161119557829003601f168201915b5050505050815250509050919050565b6111ca611501565b6001600160a01b03811661122f5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b60648201526084016105c2565b610ad18161155b565b611240611501565b606554611266906001600160401b03600160801b8204811691600160c01b900416611d8b565b6001600160401b03164210156112ca5760405162461bcd60e51b815260206004820152602360248201527f70726963652075706461746520636f6f6c646f776e206e6f74207375727061736044820152621cd95960ea1b60648201526084016105c2565b606554604080516001600160401b0380841682528681166020830152600160401b84048116828401528581166060830152600160801b9093048316608082015291831660a0830152517f9b97ed982ea5820e21bfc9578505e78068a5333487583460ad56ff72defef77a9181900360c00190a160658054426001600160401b03908116600160c01b026001600160c01b03948216600160801b0277ffffffffffffffff0000000000000000ffffffffffffffff19968316600160401b02969096166001600160c01b0319909316929092179516949094179290921716919091179055565b606760205260009081526040902080546001820180546001600160401b0380841694600160401b8504821694600160801b90049091169290916113f090611d19565b80601f016020809104026020016040519081016040528092919081815260200182805461141c90611d19565b80156114695780601f1061143e57610100808354040283529160200191611469565b820191906000526020600020905b81548152906001019060200180831161144c57829003601f168201915b50505050509080600201805461147e90611d19565b80601f01602080910402602001604051908101604052809291908181526020018280546114aa90611d19565b80156114f75780601f106114cc576101008083540402835291602001916114f7565b820191906000526020600020905b8154815290600101906020018083116114da57829003601f168201915b5050505050905085565b6033546001600160a01b03163314610a3d5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016105c2565b603380546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b80518251146115fe5760405162461bcd60e51b815260206004820181905260248201527f617272617973206d7573742068617665207468652073616d65206c656e67746860448201526064016105c2565b6000805b82518110156116415782818151811061161d5761161d611d03565b016020015161162f9060f81c83611dad565b915061163a81611d4e565b9050611602565b508060ff16606414610f905760405162461bcd60e51b815260206004820152601f60248201527f73756d206f662071756f72756d53706c697473206d757374206265203130300060448201526064016105c2565b8280546116a190611d19565b90600052602060002090601f0160209004810192826116c35760008555611709565b82601f106116dc57805160ff1916838001178555611709565b82800160010185558215611709579182015b828111156117095782518255916020019190600101906116ee565b50611715929150611719565b5090565b5b80821115611715576000815560010161171a565b634e487b7160e01b600052604160045260246000fd5b60405160a081016001600160401b03811182821017156117665761176661172e565b60405290565b604051601f8201601f191681016001600160401b03811182821017156117945761179461172e565b604052919050565b80356001600160a01b03811681146117b357600080fd5b919050565b600060208083850312156117cb57600080fd5b82356001600160401b03808211156117e257600080fd5b818501915085601f8301126117f657600080fd5b8135818111156118085761180861172e565b8060051b915061181984830161176c565b818152918301840191848101908884111561183357600080fd5b938501935b83851015611858576118498561179c565b82529385019390850190611838565b98975050505050505050565b6000815180845260005b8181101561188a5760208185018101518683018201520161186e565b8181111561189c576000602083870101525b50601f01601f19169290920160200192915050565b60006001600160401b0380835116845280602084015116602085015280604084015116604085015250606082015160a060608501526118f360a0850182611864565b90506080830151848203608086015261190c8282611864565b95945050505050565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b8281101561196a57603f198886030184526119588583516118b1565b9450928501929085019060010161193c565b5092979650505050505050565b60006020828403121561198957600080fd5b5035919050565b6020808252825182820181905260009190848201906040850190845b818110156119d15783516001600160501b0316835292840192918401916001016119ac565b50909695505050505050565b80356001600160401b03811681146117b357600080fd5b600060208284031215611a0657600080fd5b611a0f826119dd565b9392505050565b600060208284031215611a2857600080fd5b611a0f8261179c565b600080600080600080600060e0888a031215611a4c57600080fd5b611a558861179c565b9650611a63602089016119dd565b9550611a71604089016119dd565b9450611a7f606089016119dd565b9350611a8d608089016119dd565b9250611a9b60a089016119dd565b9150611aa960c089016119dd565b905092959891949750929550565b600082601f830112611ac857600080fd5b81356001600160401b03811115611ae157611ae161172e565b611af4601f8201601f191660200161176c565b818152846020838601011115611b0957600080fd5b816020850160208301376000918101602001919091529392505050565b60008060408385031215611b3957600080fd5b611b428361179c565b915060208301356001600160401b0380821115611b5e57600080fd5b9084019060a08287031215611b7257600080fd5b611b7a611744565b611b83836119dd565b8152611b91602084016119dd565b6020820152611ba2604084016119dd565b6040820152606083013582811115611bb957600080fd5b611bc588828601611ab7565b606083015250608083013582811115611bdd57600080fd5b611be988828601611ab7565b6080830152508093505050509250929050565b60008060408385031215611c0f57600080fd5b611c188361179c565b946020939093013593505050565b602081526000611a0f60208301846118b1565b600080600060608486031215611c4e57600080fd5b611c57846119dd565b9250611c65602085016119dd565b9150611c73604085016119dd565b90509250925092565b60006001600160401b038088168352808716602084015280861660408401525060a06060830152611cb060a0830185611864565b82810360808401526118588185611864565b634e487b7160e01b600052601160045260246000fd5b60006001600160501b03808316818516808303821115611cfa57611cfa611cc2565b01949350505050565b634e487b7160e01b600052603260045260246000fd5b600181811c90821680611d2d57607f821691505b602082108114156108c657634e487b7160e01b600052602260045260246000fd5b6000600019821415611d6257611d62611cc2565b5060010190565b600060208284031215611d7b57600080fd5b81518015158114611a0f57600080fd5b60006001600160401b03808316818516808303821115611cfa57611cfa611cc2565b600060ff821660ff84168060ff03821115611dca57611dca611cc2565b01939250505056fea2646970667358221220bf54aa3bf7c7e22eba135628944495414f1d13862e145a046854f3f77b6736dc64736f6c634300080c0033",
}

// ContractPaymentVaultABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractPaymentVaultMetaData.ABI instead.
var ContractPaymentVaultABI = ContractPaymentVaultMetaData.ABI

// ContractPaymentVaultBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractPaymentVaultMetaData.Bin instead.
var ContractPaymentVaultBin = ContractPaymentVaultMetaData.Bin

// DeployContractPaymentVault deploys a new Ethereum contract, binding an instance of ContractPaymentVault to it.
func DeployContractPaymentVault(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ContractPaymentVault, error) {
	parsed, err := ContractPaymentVaultMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractPaymentVaultBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractPaymentVault{ContractPaymentVaultCaller: ContractPaymentVaultCaller{contract: contract}, ContractPaymentVaultTransactor: ContractPaymentVaultTransactor{contract: contract}, ContractPaymentVaultFilterer: ContractPaymentVaultFilterer{contract: contract}}, nil
}

// ContractPaymentVault is an auto generated Go binding around an Ethereum contract.
type ContractPaymentVault struct {
	ContractPaymentVaultCaller     // Read-only binding to the contract
	ContractPaymentVaultTransactor // Write-only binding to the contract
	ContractPaymentVaultFilterer   // Log filterer for contract events
}

// ContractPaymentVaultCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractPaymentVaultCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractPaymentVaultTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractPaymentVaultTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractPaymentVaultFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractPaymentVaultFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractPaymentVaultSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractPaymentVaultSession struct {
	Contract     *ContractPaymentVault // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ContractPaymentVaultCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractPaymentVaultCallerSession struct {
	Contract *ContractPaymentVaultCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// ContractPaymentVaultTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractPaymentVaultTransactorSession struct {
	Contract     *ContractPaymentVaultTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// ContractPaymentVaultRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractPaymentVaultRaw struct {
	Contract *ContractPaymentVault // Generic contract binding to access the raw methods on
}

// ContractPaymentVaultCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractPaymentVaultCallerRaw struct {
	Contract *ContractPaymentVaultCaller // Generic read-only contract binding to access the raw methods on
}

// ContractPaymentVaultTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractPaymentVaultTransactorRaw struct {
	Contract *ContractPaymentVaultTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractPaymentVault creates a new instance of ContractPaymentVault, bound to a specific deployed contract.
func NewContractPaymentVault(address common.Address, backend bind.ContractBackend) (*ContractPaymentVault, error) {
	contract, err := bindContractPaymentVault(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVault{ContractPaymentVaultCaller: ContractPaymentVaultCaller{contract: contract}, ContractPaymentVaultTransactor: ContractPaymentVaultTransactor{contract: contract}, ContractPaymentVaultFilterer: ContractPaymentVaultFilterer{contract: contract}}, nil
}

// NewContractPaymentVaultCaller creates a new read-only instance of ContractPaymentVault, bound to a specific deployed contract.
func NewContractPaymentVaultCaller(address common.Address, caller bind.ContractCaller) (*ContractPaymentVaultCaller, error) {
	contract, err := bindContractPaymentVault(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultCaller{contract: contract}, nil
}

// NewContractPaymentVaultTransactor creates a new write-only instance of ContractPaymentVault, bound to a specific deployed contract.
func NewContractPaymentVaultTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractPaymentVaultTransactor, error) {
	contract, err := bindContractPaymentVault(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultTransactor{contract: contract}, nil
}

// NewContractPaymentVaultFilterer creates a new log filterer instance of ContractPaymentVault, bound to a specific deployed contract.
func NewContractPaymentVaultFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractPaymentVaultFilterer, error) {
	contract, err := bindContractPaymentVault(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultFilterer{contract: contract}, nil
}

// bindContractPaymentVault binds a generic wrapper to an already deployed contract.
func bindContractPaymentVault(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractPaymentVaultMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractPaymentVault *ContractPaymentVaultRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractPaymentVault.Contract.ContractPaymentVaultCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractPaymentVault *ContractPaymentVaultRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.ContractPaymentVaultTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractPaymentVault *ContractPaymentVaultRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.ContractPaymentVaultTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractPaymentVault *ContractPaymentVaultCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractPaymentVault.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractPaymentVault *ContractPaymentVaultTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractPaymentVault *ContractPaymentVaultTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.contract.Transact(opts, method, params...)
}

// GetOnDemandTotalDeposit is a free data retrieval call binding the contract method 0xd1c1fdcd.
//
// Solidity: function getOnDemandTotalDeposit(address _account) view returns(uint80)
func (_ContractPaymentVault *ContractPaymentVaultCaller) GetOnDemandTotalDeposit(opts *bind.CallOpts, _account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "getOnDemandTotalDeposit", _account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetOnDemandTotalDeposit is a free data retrieval call binding the contract method 0xd1c1fdcd.
//
// Solidity: function getOnDemandTotalDeposit(address _account) view returns(uint80)
func (_ContractPaymentVault *ContractPaymentVaultSession) GetOnDemandTotalDeposit(_account common.Address) (*big.Int, error) {
	return _ContractPaymentVault.Contract.GetOnDemandTotalDeposit(&_ContractPaymentVault.CallOpts, _account)
}

// GetOnDemandTotalDeposit is a free data retrieval call binding the contract method 0xd1c1fdcd.
//
// Solidity: function getOnDemandTotalDeposit(address _account) view returns(uint80)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) GetOnDemandTotalDeposit(_account common.Address) (*big.Int, error) {
	return _ContractPaymentVault.Contract.GetOnDemandTotalDeposit(&_ContractPaymentVault.CallOpts, _account)
}

// GetOnDemandTotalDeposits is a free data retrieval call binding the contract method 0x4184a674.
//
// Solidity: function getOnDemandTotalDeposits(address[] _accounts) view returns(uint80[] _payments)
func (_ContractPaymentVault *ContractPaymentVaultCaller) GetOnDemandTotalDeposits(opts *bind.CallOpts, _accounts []common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "getOnDemandTotalDeposits", _accounts)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetOnDemandTotalDeposits is a free data retrieval call binding the contract method 0x4184a674.
//
// Solidity: function getOnDemandTotalDeposits(address[] _accounts) view returns(uint80[] _payments)
func (_ContractPaymentVault *ContractPaymentVaultSession) GetOnDemandTotalDeposits(_accounts []common.Address) ([]*big.Int, error) {
	return _ContractPaymentVault.Contract.GetOnDemandTotalDeposits(&_ContractPaymentVault.CallOpts, _accounts)
}

// GetOnDemandTotalDeposits is a free data retrieval call binding the contract method 0x4184a674.
//
// Solidity: function getOnDemandTotalDeposits(address[] _accounts) view returns(uint80[] _payments)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) GetOnDemandTotalDeposits(_accounts []common.Address) ([]*big.Int, error) {
	return _ContractPaymentVault.Contract.GetOnDemandTotalDeposits(&_ContractPaymentVault.CallOpts, _accounts)
}

// GetReservation is a free data retrieval call binding the contract method 0xb2066f80.
//
// Solidity: function getReservation(address _account) view returns((uint64,uint64,uint64,bytes,bytes))
func (_ContractPaymentVault *ContractPaymentVaultCaller) GetReservation(opts *bind.CallOpts, _account common.Address) (IPaymentVaultReservation, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "getReservation", _account)

	if err != nil {
		return *new(IPaymentVaultReservation), err
	}

	out0 := *abi.ConvertType(out[0], new(IPaymentVaultReservation)).(*IPaymentVaultReservation)

	return out0, err

}

// GetReservation is a free data retrieval call binding the contract method 0xb2066f80.
//
// Solidity: function getReservation(address _account) view returns((uint64,uint64,uint64,bytes,bytes))
func (_ContractPaymentVault *ContractPaymentVaultSession) GetReservation(_account common.Address) (IPaymentVaultReservation, error) {
	return _ContractPaymentVault.Contract.GetReservation(&_ContractPaymentVault.CallOpts, _account)
}

// GetReservation is a free data retrieval call binding the contract method 0xb2066f80.
//
// Solidity: function getReservation(address _account) view returns((uint64,uint64,uint64,bytes,bytes))
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) GetReservation(_account common.Address) (IPaymentVaultReservation, error) {
	return _ContractPaymentVault.Contract.GetReservation(&_ContractPaymentVault.CallOpts, _account)
}

// GetReservations is a free data retrieval call binding the contract method 0x109f8fe5.
//
// Solidity: function getReservations(address[] _accounts) view returns((uint64,uint64,uint64,bytes,bytes)[] _reservations)
func (_ContractPaymentVault *ContractPaymentVaultCaller) GetReservations(opts *bind.CallOpts, _accounts []common.Address) ([]IPaymentVaultReservation, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "getReservations", _accounts)

	if err != nil {
		return *new([]IPaymentVaultReservation), err
	}

	out0 := *abi.ConvertType(out[0], new([]IPaymentVaultReservation)).(*[]IPaymentVaultReservation)

	return out0, err

}

// GetReservations is a free data retrieval call binding the contract method 0x109f8fe5.
//
// Solidity: function getReservations(address[] _accounts) view returns((uint64,uint64,uint64,bytes,bytes)[] _reservations)
func (_ContractPaymentVault *ContractPaymentVaultSession) GetReservations(_accounts []common.Address) ([]IPaymentVaultReservation, error) {
	return _ContractPaymentVault.Contract.GetReservations(&_ContractPaymentVault.CallOpts, _accounts)
}

// GetReservations is a free data retrieval call binding the contract method 0x109f8fe5.
//
// Solidity: function getReservations(address[] _accounts) view returns((uint64,uint64,uint64,bytes,bytes)[] _reservations)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) GetReservations(_accounts []common.Address) ([]IPaymentVaultReservation, error) {
	return _ContractPaymentVault.Contract.GetReservations(&_ContractPaymentVault.CallOpts, _accounts)
}

// GlobalRatePeriodInterval is a free data retrieval call binding the contract method 0xbff8a3d4.
//
// Solidity: function globalRatePeriodInterval() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCaller) GlobalRatePeriodInterval(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "globalRatePeriodInterval")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GlobalRatePeriodInterval is a free data retrieval call binding the contract method 0xbff8a3d4.
//
// Solidity: function globalRatePeriodInterval() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultSession) GlobalRatePeriodInterval() (uint64, error) {
	return _ContractPaymentVault.Contract.GlobalRatePeriodInterval(&_ContractPaymentVault.CallOpts)
}

// GlobalRatePeriodInterval is a free data retrieval call binding the contract method 0xbff8a3d4.
//
// Solidity: function globalRatePeriodInterval() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) GlobalRatePeriodInterval() (uint64, error) {
	return _ContractPaymentVault.Contract.GlobalRatePeriodInterval(&_ContractPaymentVault.CallOpts)
}

// GlobalSymbolsPerPeriod is a free data retrieval call binding the contract method 0xc98d97dd.
//
// Solidity: function globalSymbolsPerPeriod() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCaller) GlobalSymbolsPerPeriod(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "globalSymbolsPerPeriod")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GlobalSymbolsPerPeriod is a free data retrieval call binding the contract method 0xc98d97dd.
//
// Solidity: function globalSymbolsPerPeriod() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultSession) GlobalSymbolsPerPeriod() (uint64, error) {
	return _ContractPaymentVault.Contract.GlobalSymbolsPerPeriod(&_ContractPaymentVault.CallOpts)
}

// GlobalSymbolsPerPeriod is a free data retrieval call binding the contract method 0xc98d97dd.
//
// Solidity: function globalSymbolsPerPeriod() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) GlobalSymbolsPerPeriod() (uint64, error) {
	return _ContractPaymentVault.Contract.GlobalSymbolsPerPeriod(&_ContractPaymentVault.CallOpts)
}

// LastPriceUpdateTime is a free data retrieval call binding the contract method 0x49b9a7af.
//
// Solidity: function lastPriceUpdateTime() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCaller) LastPriceUpdateTime(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "lastPriceUpdateTime")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// LastPriceUpdateTime is a free data retrieval call binding the contract method 0x49b9a7af.
//
// Solidity: function lastPriceUpdateTime() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultSession) LastPriceUpdateTime() (uint64, error) {
	return _ContractPaymentVault.Contract.LastPriceUpdateTime(&_ContractPaymentVault.CallOpts)
}

// LastPriceUpdateTime is a free data retrieval call binding the contract method 0x49b9a7af.
//
// Solidity: function lastPriceUpdateTime() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) LastPriceUpdateTime() (uint64, error) {
	return _ContractPaymentVault.Contract.LastPriceUpdateTime(&_ContractPaymentVault.CallOpts)
}

// MinNumSymbols is a free data retrieval call binding the contract method 0x761dab89.
//
// Solidity: function minNumSymbols() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCaller) MinNumSymbols(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "minNumSymbols")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// MinNumSymbols is a free data retrieval call binding the contract method 0x761dab89.
//
// Solidity: function minNumSymbols() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultSession) MinNumSymbols() (uint64, error) {
	return _ContractPaymentVault.Contract.MinNumSymbols(&_ContractPaymentVault.CallOpts)
}

// MinNumSymbols is a free data retrieval call binding the contract method 0x761dab89.
//
// Solidity: function minNumSymbols() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) MinNumSymbols() (uint64, error) {
	return _ContractPaymentVault.Contract.MinNumSymbols(&_ContractPaymentVault.CallOpts)
}

// OnDemandPayments is a free data retrieval call binding the contract method 0xd996dc99.
//
// Solidity: function onDemandPayments(address ) view returns(uint80 totalDeposit)
func (_ContractPaymentVault *ContractPaymentVaultCaller) OnDemandPayments(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "onDemandPayments", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OnDemandPayments is a free data retrieval call binding the contract method 0xd996dc99.
//
// Solidity: function onDemandPayments(address ) view returns(uint80 totalDeposit)
func (_ContractPaymentVault *ContractPaymentVaultSession) OnDemandPayments(arg0 common.Address) (*big.Int, error) {
	return _ContractPaymentVault.Contract.OnDemandPayments(&_ContractPaymentVault.CallOpts, arg0)
}

// OnDemandPayments is a free data retrieval call binding the contract method 0xd996dc99.
//
// Solidity: function onDemandPayments(address ) view returns(uint80 totalDeposit)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) OnDemandPayments(arg0 common.Address) (*big.Int, error) {
	return _ContractPaymentVault.Contract.OnDemandPayments(&_ContractPaymentVault.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractPaymentVault *ContractPaymentVaultCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractPaymentVault *ContractPaymentVaultSession) Owner() (common.Address, error) {
	return _ContractPaymentVault.Contract.Owner(&_ContractPaymentVault.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) Owner() (common.Address, error) {
	return _ContractPaymentVault.Contract.Owner(&_ContractPaymentVault.CallOpts)
}

// PricePerSymbol is a free data retrieval call binding the contract method 0xf323726a.
//
// Solidity: function pricePerSymbol() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCaller) PricePerSymbol(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "pricePerSymbol")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// PricePerSymbol is a free data retrieval call binding the contract method 0xf323726a.
//
// Solidity: function pricePerSymbol() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultSession) PricePerSymbol() (uint64, error) {
	return _ContractPaymentVault.Contract.PricePerSymbol(&_ContractPaymentVault.CallOpts)
}

// PricePerSymbol is a free data retrieval call binding the contract method 0xf323726a.
//
// Solidity: function pricePerSymbol() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) PricePerSymbol() (uint64, error) {
	return _ContractPaymentVault.Contract.PricePerSymbol(&_ContractPaymentVault.CallOpts)
}

// PriceUpdateCooldown is a free data retrieval call binding the contract method 0x039f091c.
//
// Solidity: function priceUpdateCooldown() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCaller) PriceUpdateCooldown(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "priceUpdateCooldown")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// PriceUpdateCooldown is a free data retrieval call binding the contract method 0x039f091c.
//
// Solidity: function priceUpdateCooldown() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultSession) PriceUpdateCooldown() (uint64, error) {
	return _ContractPaymentVault.Contract.PriceUpdateCooldown(&_ContractPaymentVault.CallOpts)
}

// PriceUpdateCooldown is a free data retrieval call binding the contract method 0x039f091c.
//
// Solidity: function priceUpdateCooldown() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) PriceUpdateCooldown() (uint64, error) {
	return _ContractPaymentVault.Contract.PriceUpdateCooldown(&_ContractPaymentVault.CallOpts)
}

// ReservationPeriodInterval is a free data retrieval call binding the contract method 0x72228ab2.
//
// Solidity: function reservationPeriodInterval() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCaller) ReservationPeriodInterval(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "reservationPeriodInterval")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// ReservationPeriodInterval is a free data retrieval call binding the contract method 0x72228ab2.
//
// Solidity: function reservationPeriodInterval() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultSession) ReservationPeriodInterval() (uint64, error) {
	return _ContractPaymentVault.Contract.ReservationPeriodInterval(&_ContractPaymentVault.CallOpts)
}

// ReservationPeriodInterval is a free data retrieval call binding the contract method 0x72228ab2.
//
// Solidity: function reservationPeriodInterval() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) ReservationPeriodInterval() (uint64, error) {
	return _ContractPaymentVault.Contract.ReservationPeriodInterval(&_ContractPaymentVault.CallOpts)
}

// Reservations is a free data retrieval call binding the contract method 0xfd3dc53a.
//
// Solidity: function reservations(address ) view returns(uint64 symbolsPerSecond, uint64 startTimestamp, uint64 endTimestamp, bytes quorumNumbers, bytes quorumSplits)
func (_ContractPaymentVault *ContractPaymentVaultCaller) Reservations(opts *bind.CallOpts, arg0 common.Address) (struct {
	SymbolsPerSecond uint64
	StartTimestamp   uint64
	EndTimestamp     uint64
	QuorumNumbers    []byte
	QuorumSplits     []byte
}, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "reservations", arg0)

	outstruct := new(struct {
		SymbolsPerSecond uint64
		StartTimestamp   uint64
		EndTimestamp     uint64
		QuorumNumbers    []byte
		QuorumSplits     []byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.SymbolsPerSecond = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.StartTimestamp = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.EndTimestamp = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.QuorumNumbers = *abi.ConvertType(out[3], new([]byte)).(*[]byte)
	outstruct.QuorumSplits = *abi.ConvertType(out[4], new([]byte)).(*[]byte)

	return *outstruct, err

}

// Reservations is a free data retrieval call binding the contract method 0xfd3dc53a.
//
// Solidity: function reservations(address ) view returns(uint64 symbolsPerSecond, uint64 startTimestamp, uint64 endTimestamp, bytes quorumNumbers, bytes quorumSplits)
func (_ContractPaymentVault *ContractPaymentVaultSession) Reservations(arg0 common.Address) (struct {
	SymbolsPerSecond uint64
	StartTimestamp   uint64
	EndTimestamp     uint64
	QuorumNumbers    []byte
	QuorumSplits     []byte
}, error) {
	return _ContractPaymentVault.Contract.Reservations(&_ContractPaymentVault.CallOpts, arg0)
}

// Reservations is a free data retrieval call binding the contract method 0xfd3dc53a.
//
// Solidity: function reservations(address ) view returns(uint64 symbolsPerSecond, uint64 startTimestamp, uint64 endTimestamp, bytes quorumNumbers, bytes quorumSplits)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) Reservations(arg0 common.Address) (struct {
	SymbolsPerSecond uint64
	StartTimestamp   uint64
	EndTimestamp     uint64
	QuorumNumbers    []byte
	QuorumSplits     []byte
}, error) {
	return _ContractPaymentVault.Contract.Reservations(&_ContractPaymentVault.CallOpts, arg0)
}

// DepositOnDemand is a paid mutator transaction binding the contract method 0x8bec7d02.
//
// Solidity: function depositOnDemand(address _account) payable returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) DepositOnDemand(opts *bind.TransactOpts, _account common.Address) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "depositOnDemand", _account)
}

// DepositOnDemand is a paid mutator transaction binding the contract method 0x8bec7d02.
//
// Solidity: function depositOnDemand(address _account) payable returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) DepositOnDemand(_account common.Address) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.DepositOnDemand(&_ContractPaymentVault.TransactOpts, _account)
}

// DepositOnDemand is a paid mutator transaction binding the contract method 0x8bec7d02.
//
// Solidity: function depositOnDemand(address _account) payable returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) DepositOnDemand(_account common.Address) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.DepositOnDemand(&_ContractPaymentVault.TransactOpts, _account)
}

// Initialize is a paid mutator transaction binding the contract method 0x9a1bbf37.
//
// Solidity: function initialize(address _initialOwner, uint64 _minNumSymbols, uint64 _pricePerSymbol, uint64 _priceUpdateCooldown, uint64 _globalSymbolsPerPeriod, uint64 _reservationPeriodInterval, uint64 _globalRatePeriodInterval) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) Initialize(opts *bind.TransactOpts, _initialOwner common.Address, _minNumSymbols uint64, _pricePerSymbol uint64, _priceUpdateCooldown uint64, _globalSymbolsPerPeriod uint64, _reservationPeriodInterval uint64, _globalRatePeriodInterval uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "initialize", _initialOwner, _minNumSymbols, _pricePerSymbol, _priceUpdateCooldown, _globalSymbolsPerPeriod, _reservationPeriodInterval, _globalRatePeriodInterval)
}

// Initialize is a paid mutator transaction binding the contract method 0x9a1bbf37.
//
// Solidity: function initialize(address _initialOwner, uint64 _minNumSymbols, uint64 _pricePerSymbol, uint64 _priceUpdateCooldown, uint64 _globalSymbolsPerPeriod, uint64 _reservationPeriodInterval, uint64 _globalRatePeriodInterval) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) Initialize(_initialOwner common.Address, _minNumSymbols uint64, _pricePerSymbol uint64, _priceUpdateCooldown uint64, _globalSymbolsPerPeriod uint64, _reservationPeriodInterval uint64, _globalRatePeriodInterval uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.Initialize(&_ContractPaymentVault.TransactOpts, _initialOwner, _minNumSymbols, _pricePerSymbol, _priceUpdateCooldown, _globalSymbolsPerPeriod, _reservationPeriodInterval, _globalRatePeriodInterval)
}

// Initialize is a paid mutator transaction binding the contract method 0x9a1bbf37.
//
// Solidity: function initialize(address _initialOwner, uint64 _minNumSymbols, uint64 _pricePerSymbol, uint64 _priceUpdateCooldown, uint64 _globalSymbolsPerPeriod, uint64 _reservationPeriodInterval, uint64 _globalRatePeriodInterval) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) Initialize(_initialOwner common.Address, _minNumSymbols uint64, _pricePerSymbol uint64, _priceUpdateCooldown uint64, _globalSymbolsPerPeriod uint64, _reservationPeriodInterval uint64, _globalRatePeriodInterval uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.Initialize(&_ContractPaymentVault.TransactOpts, _initialOwner, _minNumSymbols, _pricePerSymbol, _priceUpdateCooldown, _globalSymbolsPerPeriod, _reservationPeriodInterval, _globalRatePeriodInterval)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.RenounceOwnership(&_ContractPaymentVault.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.RenounceOwnership(&_ContractPaymentVault.TransactOpts)
}

// SetGlobalRatePeriodInterval is a paid mutator transaction binding the contract method 0xaa788bd7.
//
// Solidity: function setGlobalRatePeriodInterval(uint64 _globalRatePeriodInterval) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) SetGlobalRatePeriodInterval(opts *bind.TransactOpts, _globalRatePeriodInterval uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "setGlobalRatePeriodInterval", _globalRatePeriodInterval)
}

// SetGlobalRatePeriodInterval is a paid mutator transaction binding the contract method 0xaa788bd7.
//
// Solidity: function setGlobalRatePeriodInterval(uint64 _globalRatePeriodInterval) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) SetGlobalRatePeriodInterval(_globalRatePeriodInterval uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetGlobalRatePeriodInterval(&_ContractPaymentVault.TransactOpts, _globalRatePeriodInterval)
}

// SetGlobalRatePeriodInterval is a paid mutator transaction binding the contract method 0xaa788bd7.
//
// Solidity: function setGlobalRatePeriodInterval(uint64 _globalRatePeriodInterval) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) SetGlobalRatePeriodInterval(_globalRatePeriodInterval uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetGlobalRatePeriodInterval(&_ContractPaymentVault.TransactOpts, _globalRatePeriodInterval)
}

// SetGlobalSymbolsPerPeriod is a paid mutator transaction binding the contract method 0xa16cf884.
//
// Solidity: function setGlobalSymbolsPerPeriod(uint64 _globalSymbolsPerPeriod) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) SetGlobalSymbolsPerPeriod(opts *bind.TransactOpts, _globalSymbolsPerPeriod uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "setGlobalSymbolsPerPeriod", _globalSymbolsPerPeriod)
}

// SetGlobalSymbolsPerPeriod is a paid mutator transaction binding the contract method 0xa16cf884.
//
// Solidity: function setGlobalSymbolsPerPeriod(uint64 _globalSymbolsPerPeriod) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) SetGlobalSymbolsPerPeriod(_globalSymbolsPerPeriod uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetGlobalSymbolsPerPeriod(&_ContractPaymentVault.TransactOpts, _globalSymbolsPerPeriod)
}

// SetGlobalSymbolsPerPeriod is a paid mutator transaction binding the contract method 0xa16cf884.
//
// Solidity: function setGlobalSymbolsPerPeriod(uint64 _globalSymbolsPerPeriod) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) SetGlobalSymbolsPerPeriod(_globalSymbolsPerPeriod uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetGlobalSymbolsPerPeriod(&_ContractPaymentVault.TransactOpts, _globalSymbolsPerPeriod)
}

// SetPriceParams is a paid mutator transaction binding the contract method 0xfba2b1d1.
//
// Solidity: function setPriceParams(uint64 _minNumSymbols, uint64 _pricePerSymbol, uint64 _priceUpdateCooldown) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) SetPriceParams(opts *bind.TransactOpts, _minNumSymbols uint64, _pricePerSymbol uint64, _priceUpdateCooldown uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "setPriceParams", _minNumSymbols, _pricePerSymbol, _priceUpdateCooldown)
}

// SetPriceParams is a paid mutator transaction binding the contract method 0xfba2b1d1.
//
// Solidity: function setPriceParams(uint64 _minNumSymbols, uint64 _pricePerSymbol, uint64 _priceUpdateCooldown) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) SetPriceParams(_minNumSymbols uint64, _pricePerSymbol uint64, _priceUpdateCooldown uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetPriceParams(&_ContractPaymentVault.TransactOpts, _minNumSymbols, _pricePerSymbol, _priceUpdateCooldown)
}

// SetPriceParams is a paid mutator transaction binding the contract method 0xfba2b1d1.
//
// Solidity: function setPriceParams(uint64 _minNumSymbols, uint64 _pricePerSymbol, uint64 _priceUpdateCooldown) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) SetPriceParams(_minNumSymbols uint64, _pricePerSymbol uint64, _priceUpdateCooldown uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetPriceParams(&_ContractPaymentVault.TransactOpts, _minNumSymbols, _pricePerSymbol, _priceUpdateCooldown)
}

// SetReservation is a paid mutator transaction binding the contract method 0x9aec8640.
//
// Solidity: function setReservation(address _account, (uint64,uint64,uint64,bytes,bytes) _reservation) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) SetReservation(opts *bind.TransactOpts, _account common.Address, _reservation IPaymentVaultReservation) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "setReservation", _account, _reservation)
}

// SetReservation is a paid mutator transaction binding the contract method 0x9aec8640.
//
// Solidity: function setReservation(address _account, (uint64,uint64,uint64,bytes,bytes) _reservation) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) SetReservation(_account common.Address, _reservation IPaymentVaultReservation) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetReservation(&_ContractPaymentVault.TransactOpts, _account, _reservation)
}

// SetReservation is a paid mutator transaction binding the contract method 0x9aec8640.
//
// Solidity: function setReservation(address _account, (uint64,uint64,uint64,bytes,bytes) _reservation) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) SetReservation(_account common.Address, _reservation IPaymentVaultReservation) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetReservation(&_ContractPaymentVault.TransactOpts, _account, _reservation)
}

// SetReservationPeriodInterval is a paid mutator transaction binding the contract method 0x897218fc.
//
// Solidity: function setReservationPeriodInterval(uint64 _reservationPeriodInterval) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) SetReservationPeriodInterval(opts *bind.TransactOpts, _reservationPeriodInterval uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "setReservationPeriodInterval", _reservationPeriodInterval)
}

// SetReservationPeriodInterval is a paid mutator transaction binding the contract method 0x897218fc.
//
// Solidity: function setReservationPeriodInterval(uint64 _reservationPeriodInterval) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) SetReservationPeriodInterval(_reservationPeriodInterval uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetReservationPeriodInterval(&_ContractPaymentVault.TransactOpts, _reservationPeriodInterval)
}

// SetReservationPeriodInterval is a paid mutator transaction binding the contract method 0x897218fc.
//
// Solidity: function setReservationPeriodInterval(uint64 _reservationPeriodInterval) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) SetReservationPeriodInterval(_reservationPeriodInterval uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetReservationPeriodInterval(&_ContractPaymentVault.TransactOpts, _reservationPeriodInterval)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.TransferOwnership(&_ContractPaymentVault.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.TransferOwnership(&_ContractPaymentVault.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _amount) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) Withdraw(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "withdraw", _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _amount) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) Withdraw(_amount *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.Withdraw(&_ContractPaymentVault.TransactOpts, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _amount) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) Withdraw(_amount *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.Withdraw(&_ContractPaymentVault.TransactOpts, _amount)
}

// WithdrawERC20 is a paid mutator transaction binding the contract method 0xa1db9782.
//
// Solidity: function withdrawERC20(address _token, uint256 _amount) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) WithdrawERC20(opts *bind.TransactOpts, _token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "withdrawERC20", _token, _amount)
}

// WithdrawERC20 is a paid mutator transaction binding the contract method 0xa1db9782.
//
// Solidity: function withdrawERC20(address _token, uint256 _amount) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) WithdrawERC20(_token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.WithdrawERC20(&_ContractPaymentVault.TransactOpts, _token, _amount)
}

// WithdrawERC20 is a paid mutator transaction binding the contract method 0xa1db9782.
//
// Solidity: function withdrawERC20(address _token, uint256 _amount) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) WithdrawERC20(_token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.WithdrawERC20(&_ContractPaymentVault.TransactOpts, _token, _amount)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.Fallback(&_ContractPaymentVault.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.Fallback(&_ContractPaymentVault.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) Receive() (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.Receive(&_ContractPaymentVault.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) Receive() (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.Receive(&_ContractPaymentVault.TransactOpts)
}

// ContractPaymentVaultGlobalRatePeriodIntervalUpdatedIterator is returned from FilterGlobalRatePeriodIntervalUpdated and is used to iterate over the raw logs and unpacked data for GlobalRatePeriodIntervalUpdated events raised by the ContractPaymentVault contract.
type ContractPaymentVaultGlobalRatePeriodIntervalUpdatedIterator struct {
	Event *ContractPaymentVaultGlobalRatePeriodIntervalUpdated // Event containing the contract specifics and raw log

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
func (it *ContractPaymentVaultGlobalRatePeriodIntervalUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPaymentVaultGlobalRatePeriodIntervalUpdated)
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
		it.Event = new(ContractPaymentVaultGlobalRatePeriodIntervalUpdated)
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
func (it *ContractPaymentVaultGlobalRatePeriodIntervalUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPaymentVaultGlobalRatePeriodIntervalUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPaymentVaultGlobalRatePeriodIntervalUpdated represents a GlobalRatePeriodIntervalUpdated event raised by the ContractPaymentVault contract.
type ContractPaymentVaultGlobalRatePeriodIntervalUpdated struct {
	PreviousValue uint64
	NewValue      uint64
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterGlobalRatePeriodIntervalUpdated is a free log retrieval operation binding the contract event 0x833819c38214ef9f462f88b5c27a21bf201f394572a14da3e63c77ee15f0e93a.
//
// Solidity: event GlobalRatePeriodIntervalUpdated(uint64 previousValue, uint64 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) FilterGlobalRatePeriodIntervalUpdated(opts *bind.FilterOpts) (*ContractPaymentVaultGlobalRatePeriodIntervalUpdatedIterator, error) {

	logs, sub, err := _ContractPaymentVault.contract.FilterLogs(opts, "GlobalRatePeriodIntervalUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultGlobalRatePeriodIntervalUpdatedIterator{contract: _ContractPaymentVault.contract, event: "GlobalRatePeriodIntervalUpdated", logs: logs, sub: sub}, nil
}

// WatchGlobalRatePeriodIntervalUpdated is a free log subscription operation binding the contract event 0x833819c38214ef9f462f88b5c27a21bf201f394572a14da3e63c77ee15f0e93a.
//
// Solidity: event GlobalRatePeriodIntervalUpdated(uint64 previousValue, uint64 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) WatchGlobalRatePeriodIntervalUpdated(opts *bind.WatchOpts, sink chan<- *ContractPaymentVaultGlobalRatePeriodIntervalUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractPaymentVault.contract.WatchLogs(opts, "GlobalRatePeriodIntervalUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPaymentVaultGlobalRatePeriodIntervalUpdated)
				if err := _ContractPaymentVault.contract.UnpackLog(event, "GlobalRatePeriodIntervalUpdated", log); err != nil {
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

// ParseGlobalRatePeriodIntervalUpdated is a log parse operation binding the contract event 0x833819c38214ef9f462f88b5c27a21bf201f394572a14da3e63c77ee15f0e93a.
//
// Solidity: event GlobalRatePeriodIntervalUpdated(uint64 previousValue, uint64 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) ParseGlobalRatePeriodIntervalUpdated(log types.Log) (*ContractPaymentVaultGlobalRatePeriodIntervalUpdated, error) {
	event := new(ContractPaymentVaultGlobalRatePeriodIntervalUpdated)
	if err := _ContractPaymentVault.contract.UnpackLog(event, "GlobalRatePeriodIntervalUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractPaymentVaultGlobalSymbolsPerPeriodUpdatedIterator is returned from FilterGlobalSymbolsPerPeriodUpdated and is used to iterate over the raw logs and unpacked data for GlobalSymbolsPerPeriodUpdated events raised by the ContractPaymentVault contract.
type ContractPaymentVaultGlobalSymbolsPerPeriodUpdatedIterator struct {
	Event *ContractPaymentVaultGlobalSymbolsPerPeriodUpdated // Event containing the contract specifics and raw log

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
func (it *ContractPaymentVaultGlobalSymbolsPerPeriodUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPaymentVaultGlobalSymbolsPerPeriodUpdated)
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
		it.Event = new(ContractPaymentVaultGlobalSymbolsPerPeriodUpdated)
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
func (it *ContractPaymentVaultGlobalSymbolsPerPeriodUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPaymentVaultGlobalSymbolsPerPeriodUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPaymentVaultGlobalSymbolsPerPeriodUpdated represents a GlobalSymbolsPerPeriodUpdated event raised by the ContractPaymentVault contract.
type ContractPaymentVaultGlobalSymbolsPerPeriodUpdated struct {
	PreviousValue uint64
	NewValue      uint64
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterGlobalSymbolsPerPeriodUpdated is a free log retrieval operation binding the contract event 0x3edf3b79e74d9e583ff51df95fbabefe15f504d33475b2cc77cffba292268aae.
//
// Solidity: event GlobalSymbolsPerPeriodUpdated(uint64 previousValue, uint64 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) FilterGlobalSymbolsPerPeriodUpdated(opts *bind.FilterOpts) (*ContractPaymentVaultGlobalSymbolsPerPeriodUpdatedIterator, error) {

	logs, sub, err := _ContractPaymentVault.contract.FilterLogs(opts, "GlobalSymbolsPerPeriodUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultGlobalSymbolsPerPeriodUpdatedIterator{contract: _ContractPaymentVault.contract, event: "GlobalSymbolsPerPeriodUpdated", logs: logs, sub: sub}, nil
}

// WatchGlobalSymbolsPerPeriodUpdated is a free log subscription operation binding the contract event 0x3edf3b79e74d9e583ff51df95fbabefe15f504d33475b2cc77cffba292268aae.
//
// Solidity: event GlobalSymbolsPerPeriodUpdated(uint64 previousValue, uint64 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) WatchGlobalSymbolsPerPeriodUpdated(opts *bind.WatchOpts, sink chan<- *ContractPaymentVaultGlobalSymbolsPerPeriodUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractPaymentVault.contract.WatchLogs(opts, "GlobalSymbolsPerPeriodUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPaymentVaultGlobalSymbolsPerPeriodUpdated)
				if err := _ContractPaymentVault.contract.UnpackLog(event, "GlobalSymbolsPerPeriodUpdated", log); err != nil {
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

// ParseGlobalSymbolsPerPeriodUpdated is a log parse operation binding the contract event 0x3edf3b79e74d9e583ff51df95fbabefe15f504d33475b2cc77cffba292268aae.
//
// Solidity: event GlobalSymbolsPerPeriodUpdated(uint64 previousValue, uint64 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) ParseGlobalSymbolsPerPeriodUpdated(log types.Log) (*ContractPaymentVaultGlobalSymbolsPerPeriodUpdated, error) {
	event := new(ContractPaymentVaultGlobalSymbolsPerPeriodUpdated)
	if err := _ContractPaymentVault.contract.UnpackLog(event, "GlobalSymbolsPerPeriodUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractPaymentVaultInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ContractPaymentVault contract.
type ContractPaymentVaultInitializedIterator struct {
	Event *ContractPaymentVaultInitialized // Event containing the contract specifics and raw log

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
func (it *ContractPaymentVaultInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPaymentVaultInitialized)
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
		it.Event = new(ContractPaymentVaultInitialized)
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
func (it *ContractPaymentVaultInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPaymentVaultInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPaymentVaultInitialized represents a Initialized event raised by the ContractPaymentVault contract.
type ContractPaymentVaultInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) FilterInitialized(opts *bind.FilterOpts) (*ContractPaymentVaultInitializedIterator, error) {

	logs, sub, err := _ContractPaymentVault.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultInitializedIterator{contract: _ContractPaymentVault.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ContractPaymentVaultInitialized) (event.Subscription, error) {

	logs, sub, err := _ContractPaymentVault.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPaymentVaultInitialized)
				if err := _ContractPaymentVault.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ContractPaymentVault *ContractPaymentVaultFilterer) ParseInitialized(log types.Log) (*ContractPaymentVaultInitialized, error) {
	event := new(ContractPaymentVaultInitialized)
	if err := _ContractPaymentVault.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractPaymentVaultOnDemandPaymentUpdatedIterator is returned from FilterOnDemandPaymentUpdated and is used to iterate over the raw logs and unpacked data for OnDemandPaymentUpdated events raised by the ContractPaymentVault contract.
type ContractPaymentVaultOnDemandPaymentUpdatedIterator struct {
	Event *ContractPaymentVaultOnDemandPaymentUpdated // Event containing the contract specifics and raw log

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
func (it *ContractPaymentVaultOnDemandPaymentUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPaymentVaultOnDemandPaymentUpdated)
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
		it.Event = new(ContractPaymentVaultOnDemandPaymentUpdated)
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
func (it *ContractPaymentVaultOnDemandPaymentUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPaymentVaultOnDemandPaymentUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPaymentVaultOnDemandPaymentUpdated represents a OnDemandPaymentUpdated event raised by the ContractPaymentVault contract.
type ContractPaymentVaultOnDemandPaymentUpdated struct {
	Account         common.Address
	OnDemandPayment *big.Int
	TotalDeposit    *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterOnDemandPaymentUpdated is a free log retrieval operation binding the contract event 0x6fbb447a2c09b8901d70b0d5b9fbce159ee8fda4460e5af2570cab3fe0adf268.
//
// Solidity: event OnDemandPaymentUpdated(address indexed account, uint80 onDemandPayment, uint80 totalDeposit)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) FilterOnDemandPaymentUpdated(opts *bind.FilterOpts, account []common.Address) (*ContractPaymentVaultOnDemandPaymentUpdatedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractPaymentVault.contract.FilterLogs(opts, "OnDemandPaymentUpdated", accountRule)
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultOnDemandPaymentUpdatedIterator{contract: _ContractPaymentVault.contract, event: "OnDemandPaymentUpdated", logs: logs, sub: sub}, nil
}

// WatchOnDemandPaymentUpdated is a free log subscription operation binding the contract event 0x6fbb447a2c09b8901d70b0d5b9fbce159ee8fda4460e5af2570cab3fe0adf268.
//
// Solidity: event OnDemandPaymentUpdated(address indexed account, uint80 onDemandPayment, uint80 totalDeposit)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) WatchOnDemandPaymentUpdated(opts *bind.WatchOpts, sink chan<- *ContractPaymentVaultOnDemandPaymentUpdated, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractPaymentVault.contract.WatchLogs(opts, "OnDemandPaymentUpdated", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPaymentVaultOnDemandPaymentUpdated)
				if err := _ContractPaymentVault.contract.UnpackLog(event, "OnDemandPaymentUpdated", log); err != nil {
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

// ParseOnDemandPaymentUpdated is a log parse operation binding the contract event 0x6fbb447a2c09b8901d70b0d5b9fbce159ee8fda4460e5af2570cab3fe0adf268.
//
// Solidity: event OnDemandPaymentUpdated(address indexed account, uint80 onDemandPayment, uint80 totalDeposit)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) ParseOnDemandPaymentUpdated(log types.Log) (*ContractPaymentVaultOnDemandPaymentUpdated, error) {
	event := new(ContractPaymentVaultOnDemandPaymentUpdated)
	if err := _ContractPaymentVault.contract.UnpackLog(event, "OnDemandPaymentUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractPaymentVaultOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ContractPaymentVault contract.
type ContractPaymentVaultOwnershipTransferredIterator struct {
	Event *ContractPaymentVaultOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ContractPaymentVaultOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPaymentVaultOwnershipTransferred)
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
		it.Event = new(ContractPaymentVaultOwnershipTransferred)
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
func (it *ContractPaymentVaultOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPaymentVaultOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPaymentVaultOwnershipTransferred represents a OwnershipTransferred event raised by the ContractPaymentVault contract.
type ContractPaymentVaultOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ContractPaymentVaultOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractPaymentVault.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultOwnershipTransferredIterator{contract: _ContractPaymentVault.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ContractPaymentVaultOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractPaymentVault.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPaymentVaultOwnershipTransferred)
				if err := _ContractPaymentVault.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ContractPaymentVault *ContractPaymentVaultFilterer) ParseOwnershipTransferred(log types.Log) (*ContractPaymentVaultOwnershipTransferred, error) {
	event := new(ContractPaymentVaultOwnershipTransferred)
	if err := _ContractPaymentVault.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractPaymentVaultPriceParamsUpdatedIterator is returned from FilterPriceParamsUpdated and is used to iterate over the raw logs and unpacked data for PriceParamsUpdated events raised by the ContractPaymentVault contract.
type ContractPaymentVaultPriceParamsUpdatedIterator struct {
	Event *ContractPaymentVaultPriceParamsUpdated // Event containing the contract specifics and raw log

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
func (it *ContractPaymentVaultPriceParamsUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPaymentVaultPriceParamsUpdated)
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
		it.Event = new(ContractPaymentVaultPriceParamsUpdated)
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
func (it *ContractPaymentVaultPriceParamsUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPaymentVaultPriceParamsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPaymentVaultPriceParamsUpdated represents a PriceParamsUpdated event raised by the ContractPaymentVault contract.
type ContractPaymentVaultPriceParamsUpdated struct {
	PreviousMinNumSymbols       uint64
	NewMinNumSymbols            uint64
	PreviousPricePerSymbol      uint64
	NewPricePerSymbol           uint64
	PreviousPriceUpdateCooldown uint64
	NewPriceUpdateCooldown      uint64
	Raw                         types.Log // Blockchain specific contextual infos
}

// FilterPriceParamsUpdated is a free log retrieval operation binding the contract event 0x9b97ed982ea5820e21bfc9578505e78068a5333487583460ad56ff72defef77a.
//
// Solidity: event PriceParamsUpdated(uint64 previousMinNumSymbols, uint64 newMinNumSymbols, uint64 previousPricePerSymbol, uint64 newPricePerSymbol, uint64 previousPriceUpdateCooldown, uint64 newPriceUpdateCooldown)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) FilterPriceParamsUpdated(opts *bind.FilterOpts) (*ContractPaymentVaultPriceParamsUpdatedIterator, error) {

	logs, sub, err := _ContractPaymentVault.contract.FilterLogs(opts, "PriceParamsUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultPriceParamsUpdatedIterator{contract: _ContractPaymentVault.contract, event: "PriceParamsUpdated", logs: logs, sub: sub}, nil
}

// WatchPriceParamsUpdated is a free log subscription operation binding the contract event 0x9b97ed982ea5820e21bfc9578505e78068a5333487583460ad56ff72defef77a.
//
// Solidity: event PriceParamsUpdated(uint64 previousMinNumSymbols, uint64 newMinNumSymbols, uint64 previousPricePerSymbol, uint64 newPricePerSymbol, uint64 previousPriceUpdateCooldown, uint64 newPriceUpdateCooldown)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) WatchPriceParamsUpdated(opts *bind.WatchOpts, sink chan<- *ContractPaymentVaultPriceParamsUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractPaymentVault.contract.WatchLogs(opts, "PriceParamsUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPaymentVaultPriceParamsUpdated)
				if err := _ContractPaymentVault.contract.UnpackLog(event, "PriceParamsUpdated", log); err != nil {
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

// ParsePriceParamsUpdated is a log parse operation binding the contract event 0x9b97ed982ea5820e21bfc9578505e78068a5333487583460ad56ff72defef77a.
//
// Solidity: event PriceParamsUpdated(uint64 previousMinNumSymbols, uint64 newMinNumSymbols, uint64 previousPricePerSymbol, uint64 newPricePerSymbol, uint64 previousPriceUpdateCooldown, uint64 newPriceUpdateCooldown)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) ParsePriceParamsUpdated(log types.Log) (*ContractPaymentVaultPriceParamsUpdated, error) {
	event := new(ContractPaymentVaultPriceParamsUpdated)
	if err := _ContractPaymentVault.contract.UnpackLog(event, "PriceParamsUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractPaymentVaultReservationPeriodIntervalUpdatedIterator is returned from FilterReservationPeriodIntervalUpdated and is used to iterate over the raw logs and unpacked data for ReservationPeriodIntervalUpdated events raised by the ContractPaymentVault contract.
type ContractPaymentVaultReservationPeriodIntervalUpdatedIterator struct {
	Event *ContractPaymentVaultReservationPeriodIntervalUpdated // Event containing the contract specifics and raw log

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
func (it *ContractPaymentVaultReservationPeriodIntervalUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPaymentVaultReservationPeriodIntervalUpdated)
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
		it.Event = new(ContractPaymentVaultReservationPeriodIntervalUpdated)
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
func (it *ContractPaymentVaultReservationPeriodIntervalUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPaymentVaultReservationPeriodIntervalUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPaymentVaultReservationPeriodIntervalUpdated represents a ReservationPeriodIntervalUpdated event raised by the ContractPaymentVault contract.
type ContractPaymentVaultReservationPeriodIntervalUpdated struct {
	PreviousValue uint64
	NewValue      uint64
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterReservationPeriodIntervalUpdated is a free log retrieval operation binding the contract event 0x1ef4a1ce7d8e50959d15578b346bb20a5b049e5ee1978014a4ba66476265c957.
//
// Solidity: event ReservationPeriodIntervalUpdated(uint64 previousValue, uint64 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) FilterReservationPeriodIntervalUpdated(opts *bind.FilterOpts) (*ContractPaymentVaultReservationPeriodIntervalUpdatedIterator, error) {

	logs, sub, err := _ContractPaymentVault.contract.FilterLogs(opts, "ReservationPeriodIntervalUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultReservationPeriodIntervalUpdatedIterator{contract: _ContractPaymentVault.contract, event: "ReservationPeriodIntervalUpdated", logs: logs, sub: sub}, nil
}

// WatchReservationPeriodIntervalUpdated is a free log subscription operation binding the contract event 0x1ef4a1ce7d8e50959d15578b346bb20a5b049e5ee1978014a4ba66476265c957.
//
// Solidity: event ReservationPeriodIntervalUpdated(uint64 previousValue, uint64 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) WatchReservationPeriodIntervalUpdated(opts *bind.WatchOpts, sink chan<- *ContractPaymentVaultReservationPeriodIntervalUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractPaymentVault.contract.WatchLogs(opts, "ReservationPeriodIntervalUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPaymentVaultReservationPeriodIntervalUpdated)
				if err := _ContractPaymentVault.contract.UnpackLog(event, "ReservationPeriodIntervalUpdated", log); err != nil {
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

// ParseReservationPeriodIntervalUpdated is a log parse operation binding the contract event 0x1ef4a1ce7d8e50959d15578b346bb20a5b049e5ee1978014a4ba66476265c957.
//
// Solidity: event ReservationPeriodIntervalUpdated(uint64 previousValue, uint64 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) ParseReservationPeriodIntervalUpdated(log types.Log) (*ContractPaymentVaultReservationPeriodIntervalUpdated, error) {
	event := new(ContractPaymentVaultReservationPeriodIntervalUpdated)
	if err := _ContractPaymentVault.contract.UnpackLog(event, "ReservationPeriodIntervalUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractPaymentVaultReservationUpdatedIterator is returned from FilterReservationUpdated and is used to iterate over the raw logs and unpacked data for ReservationUpdated events raised by the ContractPaymentVault contract.
type ContractPaymentVaultReservationUpdatedIterator struct {
	Event *ContractPaymentVaultReservationUpdated // Event containing the contract specifics and raw log

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
func (it *ContractPaymentVaultReservationUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPaymentVaultReservationUpdated)
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
		it.Event = new(ContractPaymentVaultReservationUpdated)
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
func (it *ContractPaymentVaultReservationUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPaymentVaultReservationUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPaymentVaultReservationUpdated represents a ReservationUpdated event raised by the ContractPaymentVault contract.
type ContractPaymentVaultReservationUpdated struct {
	Account     common.Address
	Reservation IPaymentVaultReservation
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterReservationUpdated is a free log retrieval operation binding the contract event 0xff3054d138559c39b4c0826c43e94b2b2c6bc9a33ea1d0b74f16c916c7b73ec1.
//
// Solidity: event ReservationUpdated(address indexed account, (uint64,uint64,uint64,bytes,bytes) reservation)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) FilterReservationUpdated(opts *bind.FilterOpts, account []common.Address) (*ContractPaymentVaultReservationUpdatedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractPaymentVault.contract.FilterLogs(opts, "ReservationUpdated", accountRule)
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultReservationUpdatedIterator{contract: _ContractPaymentVault.contract, event: "ReservationUpdated", logs: logs, sub: sub}, nil
}

// WatchReservationUpdated is a free log subscription operation binding the contract event 0xff3054d138559c39b4c0826c43e94b2b2c6bc9a33ea1d0b74f16c916c7b73ec1.
//
// Solidity: event ReservationUpdated(address indexed account, (uint64,uint64,uint64,bytes,bytes) reservation)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) WatchReservationUpdated(opts *bind.WatchOpts, sink chan<- *ContractPaymentVaultReservationUpdated, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractPaymentVault.contract.WatchLogs(opts, "ReservationUpdated", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPaymentVaultReservationUpdated)
				if err := _ContractPaymentVault.contract.UnpackLog(event, "ReservationUpdated", log); err != nil {
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

// ParseReservationUpdated is a log parse operation binding the contract event 0xff3054d138559c39b4c0826c43e94b2b2c6bc9a33ea1d0b74f16c916c7b73ec1.
//
// Solidity: event ReservationUpdated(address indexed account, (uint64,uint64,uint64,bytes,bytes) reservation)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) ParseReservationUpdated(log types.Log) (*ContractPaymentVaultReservationUpdated, error) {
	event := new(ContractPaymentVaultReservationUpdated)
	if err := _ContractPaymentVault.contract.UnpackLog(event, "ReservationUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
