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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_reservationBinInterval\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_reservationBinStartTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_priceUpdateCooldown\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositOnDemand\",\"inputs\":[{\"name\":\"_account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"getOnDemandAmount\",\"inputs\":[{\"name\":\"_account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOnDemandAmounts\",\"inputs\":[{\"name\":\"_accounts\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[{\"name\":\"_payments\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReservation\",\"inputs\":[{\"name\":\"_account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIPaymentVault.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReservations\",\"inputs\":[{\"name\":\"_accounts\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[{\"name\":\"_reservations\",\"type\":\"tuple[]\",\"internalType\":\"structIPaymentVault.Reservation[]\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"globalSymbolsPerSecond\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_minChargeableSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_globalSymbolsPerSecond\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_pricePerSymbol\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"lastPriceUpdateTime\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minChargeableSize\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"onDemandPayments\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pricePerSymbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"priceUpdateCooldown\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"reservationBinInterval\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"reservationBinStartTimestamp\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"reservations\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setGlobalSymbolsPerSecond\",\"inputs\":[{\"name\":\"_globalSymbolsPerSecond\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinChargeableSize\",\"inputs\":[{\"name\":\"_minChargeableSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPricePerSymbol\",\"inputs\":[{\"name\":\"_pricePerSymbol\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setReservation\",\"inputs\":[{\"name\":\"_account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_reservation\",\"type\":\"tuple\",\"internalType\":\"structIPaymentVault.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawERC20\",\"inputs\":[{\"name\":\"_token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"GlobalSymbolsPerSecondUpdated\",\"inputs\":[{\"name\":\"previousValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinChargeableSizeUpdated\",\"inputs\":[{\"name\":\"previousValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OnDemandPaymentUpdated\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"onDemandPayment\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"totalDeposit\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PricePerSymbolUpdated\",\"inputs\":[{\"name\":\"previousValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReservationUpdated\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"reservation\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIPaymentVault.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"anonymous\":false}]",
	Bin: "0x60e060405234801561001057600080fd5b5060405162001a3038038062001a3083398101604081905261003191610110565b608083905260a082905260c0819052610048610050565b50505061013e565b603254610100900460ff16156100bc5760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840160405180910390fd5b60325460ff908116101561010e576032805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b60008060006060848603121561012557600080fd5b8351925060208401519150604084015190509250925092565b60805160a05160c0516118bb62000175600039600081816101620152610e05015260006102a7015260006102db01526118bb6000f3fe60806040526004361061014b5760003560e01c80638bec7d02116100b6578063d996dc991161006f578063d996dc99146103da578063e60a8dd714610407578063efb435f814610427578063f2fde38b1461045d578063f323726a1461047d578063fd3dc53a1461049357600080fd5b80638bec7d02146103125780638da5cb5b146103255780639aec86401461034d578063a1db97821461036d578063b2066f801461038d578063bfafe8bf146103ba57600080fd5b80634486bfb7116101085780634486bfb71461023257806349b9a7af1461025f5780634ec81af114610275578063550571b4146102955780635a8a6869146102c9578063715018a6146102fd57600080fd5b8063039f091c14610150578063109f8fe51461019757806312a76a20146101c45780632e1a7d4d146101da578063316e0299146101fc5780633816885014610212575b600080fd5b34801561015c57600080fd5b506101847f000000000000000000000000000000000000000000000000000000000000000081565b6040519081526020015b60405180910390f35b3480156101a357600080fd5b506101b76101b2366004611354565b6104c4565b60405161018e91906114b1565b3480156101d057600080fd5b5061018460005481565b3480156101e657600080fd5b506101fa6101f5366004611513565b610717565b005b34801561020857600080fd5b5061018460015481565b34801561021e57600080fd5b506101fa61022d366004611513565b610794565b34801561023e57600080fd5b5061025261024d366004611354565b6107dd565b60405161018e919061152c565b34801561026b57600080fd5b5061018460035481565b34801561028157600080fd5b506101fa610290366004611570565b61089c565b3480156102a157600080fd5b506101847f000000000000000000000000000000000000000000000000000000000000000081565b3480156102d557600080fd5b506101847f000000000000000000000000000000000000000000000000000000000000000081565b34801561030957600080fd5b506101fa6109c6565b6101fa6103203660046115a9565b6109da565b34801561033157600080fd5b506065546040516001600160a01b03909116815260200161018e565b34801561035957600080fd5b506101fa610368366004611651565b610a5b565b34801561037957600080fd5b506101fa610388366004611727565b610b69565b34801561039957600080fd5b506103ad6103a83660046115a9565b610c08565b60405161018e9190611751565b3480156103c657600080fd5b506101fa6103d5366004611513565b610db2565b3480156103e657600080fd5b506101846103f53660046115a9565b60056020526000908152604090205481565b34801561041357600080fd5b506101fa610422366004611513565b610dfb565b34801561043357600080fd5b506101846104423660046115a9565b6001600160a01b031660009081526005602052604090205490565b34801561046957600080fd5b506101fa6104783660046115a9565b610ed1565b34801561048957600080fd5b5061018460025481565b34801561049f57600080fd5b506104b36104ae3660046115a9565b610f4a565b60405161018e959493929190611764565b606081516001600160401b038111156104df576104df6112ca565b60405190808252806020026020018201604052801561053757816020015b6040805160a08101825260008082526020808301829052928201526060808201819052608082015282526000199092019101816104fd5790505b50905060005b8251811015610711576004600084838151811061055c5761055c6117aa565b6020908102919091018101516001600160a01b03168252818101929092526040908101600020815160a08101835281546001600160401b038082168352600160401b8204811695830195909552600160801b9004909316918301919091526001810180546060840191906105cf906117c0565b80601f01602080910402602001604051908101604052809291908181526020018280546105fb906117c0565b80156106485780601f1061061d57610100808354040283529160200191610648565b820191906000526020600020905b81548152906001019060200180831161062b57829003601f168201915b50505050508152602001600282018054610661906117c0565b80601f016020809104026020016040519081016040528092919081815260200182805461068d906117c0565b80156106da5780601f106106af576101008083540402835291602001916106da565b820191906000526020600020905b8154815290600101906020018083116106bd57829003601f168201915b5050505050815250508282815181106106f5576106f56117aa565b60200260200101819052508061070a9061180b565b905061053d565b50919050565b61071f61109d565b60006107336065546001600160a01b031690565b6001600160a01b03168260405160006040518083038185875af1925050503d806000811461077d576040519150601f19603f3d011682016040523d82523d6000602084013e610782565b606091505b505090508061079057600080fd5b5050565b61079c61109d565b60025460408051918252602082018390527f590fdc6eef6046429b66dbdf71f0317122a79e082956279f70104cd907112c67910160405180910390a1600255565b606081516001600160401b038111156107f8576107f86112ca565b604051908082528060200260200182016040528015610821578160200160208202803683370190505b50905060005b82518110156107115760056000848381518110610846576108466117aa565b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002054828281518110610881576108816117aa565b60209081029190910101526108958161180b565b9050610827565b603254610100900460ff16158080156108bc5750603254600160ff909116105b806108d65750303b1580156108d6575060325460ff166001145b61093e5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b60648201526084015b60405180910390fd5b6032805460ff191660011790558015610961576032805461ff0019166101001790555b61096a85610ed1565b60008490556001839055600282905580156109bf576032805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b5050505050565b6109ce61109d565b6109d860006110f7565b565b6001600160a01b03811660009081526005602052604081208054349290610a02908490611826565b90915550506001600160a01b038116600081815260056020908152604091829020548251348152918201527f56b34df61acb18dada28b541448a4ff3faf4c0970eb58b9980468a2c75383322910160405180910390a250565b610a6361109d565b610a7581606001518260800151611149565b6001600160a01b0382166000908152600460209081526040918290208351815483860151948601516001600160401b03908116600160801b0267ffffffffffffffff60801b19968216600160401b026fffffffffffffffffffffffffffffffff199093169190931617179390931692909217825560608301518051849392610b04926001850192910190611231565b5060808201518051610b20916002840191602090910190611231565b50905050816001600160a01b03167fff3054d138559c39b4c0826c43e94b2b2c6bc9a33ea1d0b74f16c916c7b73ec182604051610b5d9190611751565b60405180910390a25050565b610b7161109d565b816001600160a01b031663a9059cbb610b926065546001600160a01b031690565b6040516001600160e01b031960e084901b1681526001600160a01b039091166004820152602481018490526044016020604051808303816000875af1158015610bdf573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c03919061183e565b505050565b6040805160a08082018352600080835260208084018290528385018290526060808501819052608085018190526001600160a01b038716835260048252918590208551938401865280546001600160401b038082168652600160401b8204811693860193909352600160801b9004909116948301949094526001840180549394929391840191610c97906117c0565b80601f0160208091040260200160405190810160405280929190818152602001828054610cc3906117c0565b8015610d105780601f10610ce557610100808354040283529160200191610d10565b820191906000526020600020905b815481529060010190602001808311610cf357829003601f168201915b50505050508152602001600282018054610d29906117c0565b80601f0160208091040260200160405190810160405280929190818152602001828054610d55906117c0565b8015610da25780601f10610d7757610100808354040283529160200191610da2565b820191906000526020600020905b815481529060010190602001808311610d8557829003601f168201915b5050505050815250509050919050565b610dba61109d565b60015460408051918252602082018390527f33ec972d20d3bef7c2f239466a44f6c7335a4faf21f7b135c2b30138d71f8807910160405180910390a1600155565b610e0361109d565b7f0000000000000000000000000000000000000000000000000000000000000000600354610e319190611826565b421015610e8c5760405162461bcd60e51b815260206004820152602360248201527f70726963652075706461746520636f6f6c646f776e206e6f74207375727061736044820152621cd95960ea1b6064820152608401610935565b60005460408051918252602082018390527f62caca228682fc3cff59bcae8bdf562027847c9295336694d56a0892bbeed0b9910160405180910390a142600355600055565b610ed961109d565b6001600160a01b038116610f3e5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610935565b610f47816110f7565b50565b600460205260009081526040902080546001820180546001600160401b0380841694600160401b8504821694600160801b9004909116929091610f8c906117c0565b80601f0160208091040260200160405190810160405280929190818152602001828054610fb8906117c0565b80156110055780601f10610fda57610100808354040283529160200191611005565b820191906000526020600020905b815481529060010190602001808311610fe857829003601f168201915b50505050509080600201805461101a906117c0565b80601f0160208091040260200160405190810160405280929190818152602001828054611046906117c0565b80156110935780601f1061106857610100808354040283529160200191611093565b820191906000526020600020905b81548152906001019060200180831161107657829003601f168201915b5050505050905085565b6065546001600160a01b031633146109d85760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610935565b606580546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b805182511461119a5760405162461bcd60e51b815260206004820181905260248201527f617272617973206d7573742068617665207468652073616d65206c656e6774686044820152606401610935565b6000805b82518110156111dd578281815181106111b9576111b96117aa565b01602001516111cb9060f81c83611860565b91506111d68161180b565b905061119e565b508060ff16606414610c035760405162461bcd60e51b815260206004820152601f60248201527f73756d206f662071756f72756d53706c697473206d75737420626520313030006044820152606401610935565b82805461123d906117c0565b90600052602060002090601f01602090048101928261125f57600085556112a5565b82601f1061127857805160ff19168380011785556112a5565b828001600101855582156112a5579182015b828111156112a557825182559160200191906001019061128a565b506112b19291506112b5565b5090565b5b808211156112b157600081556001016112b6565b634e487b7160e01b600052604160045260246000fd5b60405160a081016001600160401b0381118282101715611302576113026112ca565b60405290565b604051601f8201601f191681016001600160401b0381118282101715611330576113306112ca565b604052919050565b80356001600160a01b038116811461134f57600080fd5b919050565b6000602080838503121561136757600080fd5b82356001600160401b038082111561137e57600080fd5b818501915085601f83011261139257600080fd5b8135818111156113a4576113a46112ca565b8060051b91506113b5848301611308565b81815291830184019184810190888411156113cf57600080fd5b938501935b838510156113f4576113e585611338565b825293850193908501906113d4565b98975050505050505050565b6000815180845260005b818110156114265760208185018101518683018201520161140a565b81811115611438576000602083870101525b50601f01601f19169290920160200192915050565b60006001600160401b0380835116845280602084015116602085015280604084015116604085015250606082015160a0606085015261148f60a0850182611400565b9050608083015184820360808601526114a88282611400565b95945050505050565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b8281101561150657603f198886030184526114f485835161144d565b945092850192908501906001016114d8565b5092979650505050505050565b60006020828403121561152557600080fd5b5035919050565b6020808252825182820181905260009190848201906040850190845b8181101561156457835183529284019291840191600101611548565b50909695505050505050565b6000806000806080858703121561158657600080fd5b61158f85611338565b966020860135965060408601359560600135945092505050565b6000602082840312156115bb57600080fd5b6115c482611338565b9392505050565b80356001600160401b038116811461134f57600080fd5b600082601f8301126115f357600080fd5b81356001600160401b0381111561160c5761160c6112ca565b61161f601f8201601f1916602001611308565b81815284602083860101111561163457600080fd5b816020850160208301376000918101602001919091529392505050565b6000806040838503121561166457600080fd5b61166d83611338565b915060208301356001600160401b038082111561168957600080fd5b9084019060a0828703121561169d57600080fd5b6116a56112e0565b6116ae836115cb565b81526116bc602084016115cb565b60208201526116cd604084016115cb565b60408201526060830135828111156116e457600080fd5b6116f0888286016115e2565b60608301525060808301358281111561170857600080fd5b611714888286016115e2565b6080830152508093505050509250929050565b6000806040838503121561173a57600080fd5b61174383611338565b946020939093013593505050565b6020815260006115c4602083018461144d565b60006001600160401b038088168352808716602084015280861660408401525060a0606083015261179860a0830185611400565b82810360808401526113f48185611400565b634e487b7160e01b600052603260045260246000fd5b600181811c908216806117d457607f821691505b6020821081141561071157634e487b7160e01b600052602260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b600060001982141561181f5761181f6117f5565b5060010190565b60008219821115611839576118396117f5565b500190565b60006020828403121561185057600080fd5b815180151581146115c457600080fd5b600060ff821660ff84168060ff0382111561187d5761187d6117f5565b01939250505056fea2646970667358221220e500c71a831bcab2febcc1fc64db7e8b2d678b720b319304875542504ab04c0864736f6c634300080c0033",
}

// ContractPaymentVaultABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractPaymentVaultMetaData.ABI instead.
var ContractPaymentVaultABI = ContractPaymentVaultMetaData.ABI

// ContractPaymentVaultBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractPaymentVaultMetaData.Bin instead.
var ContractPaymentVaultBin = ContractPaymentVaultMetaData.Bin

// DeployContractPaymentVault deploys a new Ethereum contract, binding an instance of ContractPaymentVault to it.
func DeployContractPaymentVault(auth *bind.TransactOpts, backend bind.ContractBackend, _reservationBinInterval *big.Int, _reservationBinStartTimestamp *big.Int, _priceUpdateCooldown *big.Int) (common.Address, *types.Transaction, *ContractPaymentVault, error) {
	parsed, err := ContractPaymentVaultMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractPaymentVaultBin), backend, _reservationBinInterval, _reservationBinStartTimestamp, _priceUpdateCooldown)
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

// GetOnDemandAmount is a free data retrieval call binding the contract method 0xefb435f8.
//
// Solidity: function getOnDemandAmount(address _account) view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultCaller) GetOnDemandAmount(opts *bind.CallOpts, _account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "getOnDemandAmount", _account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetOnDemandAmount is a free data retrieval call binding the contract method 0xefb435f8.
//
// Solidity: function getOnDemandAmount(address _account) view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultSession) GetOnDemandAmount(_account common.Address) (*big.Int, error) {
	return _ContractPaymentVault.Contract.GetOnDemandAmount(&_ContractPaymentVault.CallOpts, _account)
}

// GetOnDemandAmount is a free data retrieval call binding the contract method 0xefb435f8.
//
// Solidity: function getOnDemandAmount(address _account) view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) GetOnDemandAmount(_account common.Address) (*big.Int, error) {
	return _ContractPaymentVault.Contract.GetOnDemandAmount(&_ContractPaymentVault.CallOpts, _account)
}

// GetOnDemandAmounts is a free data retrieval call binding the contract method 0x4486bfb7.
//
// Solidity: function getOnDemandAmounts(address[] _accounts) view returns(uint256[] _payments)
func (_ContractPaymentVault *ContractPaymentVaultCaller) GetOnDemandAmounts(opts *bind.CallOpts, _accounts []common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "getOnDemandAmounts", _accounts)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetOnDemandAmounts is a free data retrieval call binding the contract method 0x4486bfb7.
//
// Solidity: function getOnDemandAmounts(address[] _accounts) view returns(uint256[] _payments)
func (_ContractPaymentVault *ContractPaymentVaultSession) GetOnDemandAmounts(_accounts []common.Address) ([]*big.Int, error) {
	return _ContractPaymentVault.Contract.GetOnDemandAmounts(&_ContractPaymentVault.CallOpts, _accounts)
}

// GetOnDemandAmounts is a free data retrieval call binding the contract method 0x4486bfb7.
//
// Solidity: function getOnDemandAmounts(address[] _accounts) view returns(uint256[] _payments)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) GetOnDemandAmounts(_accounts []common.Address) ([]*big.Int, error) {
	return _ContractPaymentVault.Contract.GetOnDemandAmounts(&_ContractPaymentVault.CallOpts, _accounts)
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

// GlobalSymbolsPerSecond is a free data retrieval call binding the contract method 0x316e0299.
//
// Solidity: function globalSymbolsPerSecond() view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultCaller) GlobalSymbolsPerSecond(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "globalSymbolsPerSecond")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GlobalSymbolsPerSecond is a free data retrieval call binding the contract method 0x316e0299.
//
// Solidity: function globalSymbolsPerSecond() view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultSession) GlobalSymbolsPerSecond() (*big.Int, error) {
	return _ContractPaymentVault.Contract.GlobalSymbolsPerSecond(&_ContractPaymentVault.CallOpts)
}

// GlobalSymbolsPerSecond is a free data retrieval call binding the contract method 0x316e0299.
//
// Solidity: function globalSymbolsPerSecond() view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) GlobalSymbolsPerSecond() (*big.Int, error) {
	return _ContractPaymentVault.Contract.GlobalSymbolsPerSecond(&_ContractPaymentVault.CallOpts)
}

// LastPriceUpdateTime is a free data retrieval call binding the contract method 0x49b9a7af.
//
// Solidity: function lastPriceUpdateTime() view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultCaller) LastPriceUpdateTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "lastPriceUpdateTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastPriceUpdateTime is a free data retrieval call binding the contract method 0x49b9a7af.
//
// Solidity: function lastPriceUpdateTime() view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultSession) LastPriceUpdateTime() (*big.Int, error) {
	return _ContractPaymentVault.Contract.LastPriceUpdateTime(&_ContractPaymentVault.CallOpts)
}

// LastPriceUpdateTime is a free data retrieval call binding the contract method 0x49b9a7af.
//
// Solidity: function lastPriceUpdateTime() view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) LastPriceUpdateTime() (*big.Int, error) {
	return _ContractPaymentVault.Contract.LastPriceUpdateTime(&_ContractPaymentVault.CallOpts)
}

// MinChargeableSize is a free data retrieval call binding the contract method 0x12a76a20.
//
// Solidity: function minChargeableSize() view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultCaller) MinChargeableSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "minChargeableSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinChargeableSize is a free data retrieval call binding the contract method 0x12a76a20.
//
// Solidity: function minChargeableSize() view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultSession) MinChargeableSize() (*big.Int, error) {
	return _ContractPaymentVault.Contract.MinChargeableSize(&_ContractPaymentVault.CallOpts)
}

// MinChargeableSize is a free data retrieval call binding the contract method 0x12a76a20.
//
// Solidity: function minChargeableSize() view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) MinChargeableSize() (*big.Int, error) {
	return _ContractPaymentVault.Contract.MinChargeableSize(&_ContractPaymentVault.CallOpts)
}

// OnDemandPayments is a free data retrieval call binding the contract method 0xd996dc99.
//
// Solidity: function onDemandPayments(address ) view returns(uint256)
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
// Solidity: function onDemandPayments(address ) view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultSession) OnDemandPayments(arg0 common.Address) (*big.Int, error) {
	return _ContractPaymentVault.Contract.OnDemandPayments(&_ContractPaymentVault.CallOpts, arg0)
}

// OnDemandPayments is a free data retrieval call binding the contract method 0xd996dc99.
//
// Solidity: function onDemandPayments(address ) view returns(uint256)
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
// Solidity: function pricePerSymbol() view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultCaller) PricePerSymbol(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "pricePerSymbol")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PricePerSymbol is a free data retrieval call binding the contract method 0xf323726a.
//
// Solidity: function pricePerSymbol() view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultSession) PricePerSymbol() (*big.Int, error) {
	return _ContractPaymentVault.Contract.PricePerSymbol(&_ContractPaymentVault.CallOpts)
}

// PricePerSymbol is a free data retrieval call binding the contract method 0xf323726a.
//
// Solidity: function pricePerSymbol() view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) PricePerSymbol() (*big.Int, error) {
	return _ContractPaymentVault.Contract.PricePerSymbol(&_ContractPaymentVault.CallOpts)
}

// PriceUpdateCooldown is a free data retrieval call binding the contract method 0x039f091c.
//
// Solidity: function priceUpdateCooldown() view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultCaller) PriceUpdateCooldown(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "priceUpdateCooldown")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PriceUpdateCooldown is a free data retrieval call binding the contract method 0x039f091c.
//
// Solidity: function priceUpdateCooldown() view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultSession) PriceUpdateCooldown() (*big.Int, error) {
	return _ContractPaymentVault.Contract.PriceUpdateCooldown(&_ContractPaymentVault.CallOpts)
}

// PriceUpdateCooldown is a free data retrieval call binding the contract method 0x039f091c.
//
// Solidity: function priceUpdateCooldown() view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) PriceUpdateCooldown() (*big.Int, error) {
	return _ContractPaymentVault.Contract.PriceUpdateCooldown(&_ContractPaymentVault.CallOpts)
}

// ReservationBinInterval is a free data retrieval call binding the contract method 0x5a8a6869.
//
// Solidity: function reservationBinInterval() view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultCaller) ReservationBinInterval(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "reservationBinInterval")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ReservationBinInterval is a free data retrieval call binding the contract method 0x5a8a6869.
//
// Solidity: function reservationBinInterval() view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultSession) ReservationBinInterval() (*big.Int, error) {
	return _ContractPaymentVault.Contract.ReservationBinInterval(&_ContractPaymentVault.CallOpts)
}

// ReservationBinInterval is a free data retrieval call binding the contract method 0x5a8a6869.
//
// Solidity: function reservationBinInterval() view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) ReservationBinInterval() (*big.Int, error) {
	return _ContractPaymentVault.Contract.ReservationBinInterval(&_ContractPaymentVault.CallOpts)
}

// ReservationBinStartTimestamp is a free data retrieval call binding the contract method 0x550571b4.
//
// Solidity: function reservationBinStartTimestamp() view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultCaller) ReservationBinStartTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "reservationBinStartTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ReservationBinStartTimestamp is a free data retrieval call binding the contract method 0x550571b4.
//
// Solidity: function reservationBinStartTimestamp() view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultSession) ReservationBinStartTimestamp() (*big.Int, error) {
	return _ContractPaymentVault.Contract.ReservationBinStartTimestamp(&_ContractPaymentVault.CallOpts)
}

// ReservationBinStartTimestamp is a free data retrieval call binding the contract method 0x550571b4.
//
// Solidity: function reservationBinStartTimestamp() view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) ReservationBinStartTimestamp() (*big.Int, error) {
	return _ContractPaymentVault.Contract.ReservationBinStartTimestamp(&_ContractPaymentVault.CallOpts)
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

// Initialize is a paid mutator transaction binding the contract method 0x4ec81af1.
//
// Solidity: function initialize(address _initialOwner, uint256 _minChargeableSize, uint256 _globalSymbolsPerSecond, uint256 _pricePerSymbol) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) Initialize(opts *bind.TransactOpts, _initialOwner common.Address, _minChargeableSize *big.Int, _globalSymbolsPerSecond *big.Int, _pricePerSymbol *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "initialize", _initialOwner, _minChargeableSize, _globalSymbolsPerSecond, _pricePerSymbol)
}

// Initialize is a paid mutator transaction binding the contract method 0x4ec81af1.
//
// Solidity: function initialize(address _initialOwner, uint256 _minChargeableSize, uint256 _globalSymbolsPerSecond, uint256 _pricePerSymbol) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) Initialize(_initialOwner common.Address, _minChargeableSize *big.Int, _globalSymbolsPerSecond *big.Int, _pricePerSymbol *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.Initialize(&_ContractPaymentVault.TransactOpts, _initialOwner, _minChargeableSize, _globalSymbolsPerSecond, _pricePerSymbol)
}

// Initialize is a paid mutator transaction binding the contract method 0x4ec81af1.
//
// Solidity: function initialize(address _initialOwner, uint256 _minChargeableSize, uint256 _globalSymbolsPerSecond, uint256 _pricePerSymbol) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) Initialize(_initialOwner common.Address, _minChargeableSize *big.Int, _globalSymbolsPerSecond *big.Int, _pricePerSymbol *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.Initialize(&_ContractPaymentVault.TransactOpts, _initialOwner, _minChargeableSize, _globalSymbolsPerSecond, _pricePerSymbol)
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

// SetGlobalSymbolsPerSecond is a paid mutator transaction binding the contract method 0xbfafe8bf.
//
// Solidity: function setGlobalSymbolsPerSecond(uint256 _globalSymbolsPerSecond) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) SetGlobalSymbolsPerSecond(opts *bind.TransactOpts, _globalSymbolsPerSecond *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "setGlobalSymbolsPerSecond", _globalSymbolsPerSecond)
}

// SetGlobalSymbolsPerSecond is a paid mutator transaction binding the contract method 0xbfafe8bf.
//
// Solidity: function setGlobalSymbolsPerSecond(uint256 _globalSymbolsPerSecond) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) SetGlobalSymbolsPerSecond(_globalSymbolsPerSecond *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetGlobalSymbolsPerSecond(&_ContractPaymentVault.TransactOpts, _globalSymbolsPerSecond)
}

// SetGlobalSymbolsPerSecond is a paid mutator transaction binding the contract method 0xbfafe8bf.
//
// Solidity: function setGlobalSymbolsPerSecond(uint256 _globalSymbolsPerSecond) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) SetGlobalSymbolsPerSecond(_globalSymbolsPerSecond *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetGlobalSymbolsPerSecond(&_ContractPaymentVault.TransactOpts, _globalSymbolsPerSecond)
}

// SetMinChargeableSize is a paid mutator transaction binding the contract method 0xe60a8dd7.
//
// Solidity: function setMinChargeableSize(uint256 _minChargeableSize) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) SetMinChargeableSize(opts *bind.TransactOpts, _minChargeableSize *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "setMinChargeableSize", _minChargeableSize)
}

// SetMinChargeableSize is a paid mutator transaction binding the contract method 0xe60a8dd7.
//
// Solidity: function setMinChargeableSize(uint256 _minChargeableSize) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) SetMinChargeableSize(_minChargeableSize *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetMinChargeableSize(&_ContractPaymentVault.TransactOpts, _minChargeableSize)
}

// SetMinChargeableSize is a paid mutator transaction binding the contract method 0xe60a8dd7.
//
// Solidity: function setMinChargeableSize(uint256 _minChargeableSize) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) SetMinChargeableSize(_minChargeableSize *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetMinChargeableSize(&_ContractPaymentVault.TransactOpts, _minChargeableSize)
}

// SetPricePerSymbol is a paid mutator transaction binding the contract method 0x38168850.
//
// Solidity: function setPricePerSymbol(uint256 _pricePerSymbol) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) SetPricePerSymbol(opts *bind.TransactOpts, _pricePerSymbol *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "setPricePerSymbol", _pricePerSymbol)
}

// SetPricePerSymbol is a paid mutator transaction binding the contract method 0x38168850.
//
// Solidity: function setPricePerSymbol(uint256 _pricePerSymbol) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) SetPricePerSymbol(_pricePerSymbol *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetPricePerSymbol(&_ContractPaymentVault.TransactOpts, _pricePerSymbol)
}

// SetPricePerSymbol is a paid mutator transaction binding the contract method 0x38168850.
//
// Solidity: function setPricePerSymbol(uint256 _pricePerSymbol) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) SetPricePerSymbol(_pricePerSymbol *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetPricePerSymbol(&_ContractPaymentVault.TransactOpts, _pricePerSymbol)
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

// ContractPaymentVaultGlobalSymbolsPerSecondUpdatedIterator is returned from FilterGlobalSymbolsPerSecondUpdated and is used to iterate over the raw logs and unpacked data for GlobalSymbolsPerSecondUpdated events raised by the ContractPaymentVault contract.
type ContractPaymentVaultGlobalSymbolsPerSecondUpdatedIterator struct {
	Event *ContractPaymentVaultGlobalSymbolsPerSecondUpdated // Event containing the contract specifics and raw log

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
func (it *ContractPaymentVaultGlobalSymbolsPerSecondUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPaymentVaultGlobalSymbolsPerSecondUpdated)
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
		it.Event = new(ContractPaymentVaultGlobalSymbolsPerSecondUpdated)
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
func (it *ContractPaymentVaultGlobalSymbolsPerSecondUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPaymentVaultGlobalSymbolsPerSecondUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPaymentVaultGlobalSymbolsPerSecondUpdated represents a GlobalSymbolsPerSecondUpdated event raised by the ContractPaymentVault contract.
type ContractPaymentVaultGlobalSymbolsPerSecondUpdated struct {
	PreviousValue *big.Int
	NewValue      *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterGlobalSymbolsPerSecondUpdated is a free log retrieval operation binding the contract event 0x33ec972d20d3bef7c2f239466a44f6c7335a4faf21f7b135c2b30138d71f8807.
//
// Solidity: event GlobalSymbolsPerSecondUpdated(uint256 previousValue, uint256 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) FilterGlobalSymbolsPerSecondUpdated(opts *bind.FilterOpts) (*ContractPaymentVaultGlobalSymbolsPerSecondUpdatedIterator, error) {

	logs, sub, err := _ContractPaymentVault.contract.FilterLogs(opts, "GlobalSymbolsPerSecondUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultGlobalSymbolsPerSecondUpdatedIterator{contract: _ContractPaymentVault.contract, event: "GlobalSymbolsPerSecondUpdated", logs: logs, sub: sub}, nil
}

// WatchGlobalSymbolsPerSecondUpdated is a free log subscription operation binding the contract event 0x33ec972d20d3bef7c2f239466a44f6c7335a4faf21f7b135c2b30138d71f8807.
//
// Solidity: event GlobalSymbolsPerSecondUpdated(uint256 previousValue, uint256 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) WatchGlobalSymbolsPerSecondUpdated(opts *bind.WatchOpts, sink chan<- *ContractPaymentVaultGlobalSymbolsPerSecondUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractPaymentVault.contract.WatchLogs(opts, "GlobalSymbolsPerSecondUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPaymentVaultGlobalSymbolsPerSecondUpdated)
				if err := _ContractPaymentVault.contract.UnpackLog(event, "GlobalSymbolsPerSecondUpdated", log); err != nil {
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

// ParseGlobalSymbolsPerSecondUpdated is a log parse operation binding the contract event 0x33ec972d20d3bef7c2f239466a44f6c7335a4faf21f7b135c2b30138d71f8807.
//
// Solidity: event GlobalSymbolsPerSecondUpdated(uint256 previousValue, uint256 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) ParseGlobalSymbolsPerSecondUpdated(log types.Log) (*ContractPaymentVaultGlobalSymbolsPerSecondUpdated, error) {
	event := new(ContractPaymentVaultGlobalSymbolsPerSecondUpdated)
	if err := _ContractPaymentVault.contract.UnpackLog(event, "GlobalSymbolsPerSecondUpdated", log); err != nil {
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

// ContractPaymentVaultMinChargeableSizeUpdatedIterator is returned from FilterMinChargeableSizeUpdated and is used to iterate over the raw logs and unpacked data for MinChargeableSizeUpdated events raised by the ContractPaymentVault contract.
type ContractPaymentVaultMinChargeableSizeUpdatedIterator struct {
	Event *ContractPaymentVaultMinChargeableSizeUpdated // Event containing the contract specifics and raw log

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
func (it *ContractPaymentVaultMinChargeableSizeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPaymentVaultMinChargeableSizeUpdated)
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
		it.Event = new(ContractPaymentVaultMinChargeableSizeUpdated)
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
func (it *ContractPaymentVaultMinChargeableSizeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPaymentVaultMinChargeableSizeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPaymentVaultMinChargeableSizeUpdated represents a MinChargeableSizeUpdated event raised by the ContractPaymentVault contract.
type ContractPaymentVaultMinChargeableSizeUpdated struct {
	PreviousValue *big.Int
	NewValue      *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterMinChargeableSizeUpdated is a free log retrieval operation binding the contract event 0x62caca228682fc3cff59bcae8bdf562027847c9295336694d56a0892bbeed0b9.
//
// Solidity: event MinChargeableSizeUpdated(uint256 previousValue, uint256 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) FilterMinChargeableSizeUpdated(opts *bind.FilterOpts) (*ContractPaymentVaultMinChargeableSizeUpdatedIterator, error) {

	logs, sub, err := _ContractPaymentVault.contract.FilterLogs(opts, "MinChargeableSizeUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultMinChargeableSizeUpdatedIterator{contract: _ContractPaymentVault.contract, event: "MinChargeableSizeUpdated", logs: logs, sub: sub}, nil
}

// WatchMinChargeableSizeUpdated is a free log subscription operation binding the contract event 0x62caca228682fc3cff59bcae8bdf562027847c9295336694d56a0892bbeed0b9.
//
// Solidity: event MinChargeableSizeUpdated(uint256 previousValue, uint256 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) WatchMinChargeableSizeUpdated(opts *bind.WatchOpts, sink chan<- *ContractPaymentVaultMinChargeableSizeUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractPaymentVault.contract.WatchLogs(opts, "MinChargeableSizeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPaymentVaultMinChargeableSizeUpdated)
				if err := _ContractPaymentVault.contract.UnpackLog(event, "MinChargeableSizeUpdated", log); err != nil {
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

// ParseMinChargeableSizeUpdated is a log parse operation binding the contract event 0x62caca228682fc3cff59bcae8bdf562027847c9295336694d56a0892bbeed0b9.
//
// Solidity: event MinChargeableSizeUpdated(uint256 previousValue, uint256 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) ParseMinChargeableSizeUpdated(log types.Log) (*ContractPaymentVaultMinChargeableSizeUpdated, error) {
	event := new(ContractPaymentVaultMinChargeableSizeUpdated)
	if err := _ContractPaymentVault.contract.UnpackLog(event, "MinChargeableSizeUpdated", log); err != nil {
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

// FilterOnDemandPaymentUpdated is a free log retrieval operation binding the contract event 0x56b34df61acb18dada28b541448a4ff3faf4c0970eb58b9980468a2c75383322.
//
// Solidity: event OnDemandPaymentUpdated(address indexed account, uint256 onDemandPayment, uint256 totalDeposit)
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

// WatchOnDemandPaymentUpdated is a free log subscription operation binding the contract event 0x56b34df61acb18dada28b541448a4ff3faf4c0970eb58b9980468a2c75383322.
//
// Solidity: event OnDemandPaymentUpdated(address indexed account, uint256 onDemandPayment, uint256 totalDeposit)
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

// ParseOnDemandPaymentUpdated is a log parse operation binding the contract event 0x56b34df61acb18dada28b541448a4ff3faf4c0970eb58b9980468a2c75383322.
//
// Solidity: event OnDemandPaymentUpdated(address indexed account, uint256 onDemandPayment, uint256 totalDeposit)
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

// ContractPaymentVaultPricePerSymbolUpdatedIterator is returned from FilterPricePerSymbolUpdated and is used to iterate over the raw logs and unpacked data for PricePerSymbolUpdated events raised by the ContractPaymentVault contract.
type ContractPaymentVaultPricePerSymbolUpdatedIterator struct {
	Event *ContractPaymentVaultPricePerSymbolUpdated // Event containing the contract specifics and raw log

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
func (it *ContractPaymentVaultPricePerSymbolUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPaymentVaultPricePerSymbolUpdated)
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
		it.Event = new(ContractPaymentVaultPricePerSymbolUpdated)
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
func (it *ContractPaymentVaultPricePerSymbolUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPaymentVaultPricePerSymbolUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPaymentVaultPricePerSymbolUpdated represents a PricePerSymbolUpdated event raised by the ContractPaymentVault contract.
type ContractPaymentVaultPricePerSymbolUpdated struct {
	PreviousValue *big.Int
	NewValue      *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterPricePerSymbolUpdated is a free log retrieval operation binding the contract event 0x590fdc6eef6046429b66dbdf71f0317122a79e082956279f70104cd907112c67.
//
// Solidity: event PricePerSymbolUpdated(uint256 previousValue, uint256 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) FilterPricePerSymbolUpdated(opts *bind.FilterOpts) (*ContractPaymentVaultPricePerSymbolUpdatedIterator, error) {

	logs, sub, err := _ContractPaymentVault.contract.FilterLogs(opts, "PricePerSymbolUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultPricePerSymbolUpdatedIterator{contract: _ContractPaymentVault.contract, event: "PricePerSymbolUpdated", logs: logs, sub: sub}, nil
}

// WatchPricePerSymbolUpdated is a free log subscription operation binding the contract event 0x590fdc6eef6046429b66dbdf71f0317122a79e082956279f70104cd907112c67.
//
// Solidity: event PricePerSymbolUpdated(uint256 previousValue, uint256 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) WatchPricePerSymbolUpdated(opts *bind.WatchOpts, sink chan<- *ContractPaymentVaultPricePerSymbolUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractPaymentVault.contract.WatchLogs(opts, "PricePerSymbolUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPaymentVaultPricePerSymbolUpdated)
				if err := _ContractPaymentVault.contract.UnpackLog(event, "PricePerSymbolUpdated", log); err != nil {
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

// ParsePricePerSymbolUpdated is a log parse operation binding the contract event 0x590fdc6eef6046429b66dbdf71f0317122a79e082956279f70104cd907112c67.
//
// Solidity: event PricePerSymbolUpdated(uint256 previousValue, uint256 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) ParsePricePerSymbolUpdated(log types.Log) (*ContractPaymentVaultPricePerSymbolUpdated, error) {
	event := new(ContractPaymentVaultPricePerSymbolUpdated)
	if err := _ContractPaymentVault.contract.UnpackLog(event, "PricePerSymbolUpdated", log); err != nil {
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
