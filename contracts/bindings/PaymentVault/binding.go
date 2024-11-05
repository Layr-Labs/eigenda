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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"depositOnDemand\",\"inputs\":[{\"name\":\"_account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"getOnDemandAmount\",\"inputs\":[{\"name\":\"_account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOnDemandAmounts\",\"inputs\":[{\"name\":\"_accounts\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[{\"name\":\"_payments\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReservation\",\"inputs\":[{\"name\":\"_account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIPaymentVault.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReservations\",\"inputs\":[{\"name\":\"_accounts\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[{\"name\":\"_reservations\",\"type\":\"tuple[]\",\"internalType\":\"structIPaymentVault.Reservation[]\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"globalRateBinInterval\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"globalSymbolsPerSecond\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_minChargeableSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_globalSymbolsPerSecond\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_pricePerSymbol\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_reservationBinInterval\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_priceUpdateCooldown\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_globalRateBinInterval\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"lastPriceUpdateTime\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minChargeableSize\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"onDemandPayments\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pricePerSymbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"priceUpdateCooldown\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"reservationBinInterval\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"reservations\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setGlobalRateBinInterval\",\"inputs\":[{\"name\":\"_globalRateBinInterval\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setGlobalSymbolsPerSecond\",\"inputs\":[{\"name\":\"_globalSymbolsPerSecond\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPriceParams\",\"inputs\":[{\"name\":\"_minChargeableSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_pricePerSymbol\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_priceUpdateCooldown\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setReservation\",\"inputs\":[{\"name\":\"_account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_reservation\",\"type\":\"tuple\",\"internalType\":\"structIPaymentVault.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setReservationBinInterval\",\"inputs\":[{\"name\":\"_reservationBinInterval\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawERC20\",\"inputs\":[{\"name\":\"_token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"GlobalRateBinIntervalUpdated\",\"inputs\":[{\"name\":\"previousValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GlobalSymbolsPerSecondUpdated\",\"inputs\":[{\"name\":\"previousValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OnDemandPaymentUpdated\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"onDemandPayment\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"totalDeposit\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PriceParamsUpdated\",\"inputs\":[{\"name\":\"previousMinChargeableSize\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newMinChargeableSize\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"previousPricePerSymbol\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newPricePerSymbol\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"previousPriceUpdateCooldown\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newPriceUpdateCooldown\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReservationBinIntervalUpdated\",\"inputs\":[{\"name\":\"previousValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReservationUpdated\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"reservation\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIPaymentVault.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"anonymous\":false}]",
	Bin: "0x608060405234801561001057600080fd5b5061001961001e565b6100de565b603354610100900460ff161561008a5760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840160405180910390fd5b60335460ff90811610156100dc576033805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b6119ec806100ed6000396000f3fe60806040526004361061016a5760003560e01c80638142951a116100d1578063bfafe8bf1161008a578063efb435f811610064578063efb435f814610420578063f2fde38b14610456578063f323726a14610476578063fd3dc53a1461048c5761017b565b8063bfafe8bf146103bd578063d882a5fd146103dd578063d996dc99146103f35761017b565b80638142951a146102f55780638bec7d02146103155780638da5cb5b146103285780639aec864014610350578063a1db978214610370578063b2066f80146103905761017b565b8063327fad0811610123578063327fad08146102475780634486bfb71461026757806349b9a7af146102945780635a8a6869146102aa578063715018a6146102c05780637201dd02146102d55761017b565b8063039f091c14610185578063109f8fe5146101ae57806312a76a20146101db5780632c1a33bc146101f15780632e1a7d4d14610211578063316e0299146102315761017b565b3661017b5761017933346104bd565b005b61017933346104bd565b34801561019157600080fd5b5061019b60025481565b6040519081526020015b60405180910390f35b3480156101ba57600080fd5b506101ce6101c936600461143f565b610540565b6040516101a5919061159c565b3480156101e757600080fd5b5061019b60005481565b3480156101fd57600080fd5b5061017961020c3660046115fe565b610793565b34801561021d57600080fd5b5061017961022c3660046115fe565b6107dc565b34801561023d57600080fd5b5061019b60035481565b34801561025357600080fd5b50610179610262366004611617565b610859565b34801561027357600080fd5b5061028761028236600461143f565b610943565b6040516101a59190611643565b3480156102a057600080fd5b5061019b60065481565b3480156102b657600080fd5b5061019b60045481565b3480156102cc57600080fd5b50610179610a02565b3480156102e157600080fd5b506101796102f03660046115fe565b610a16565b34801561030157600080fd5b50610179610310366004611687565b610a5f565b6101796103233660046116da565b610b9a565b34801561033457600080fd5b506066546040516001600160a01b0390911681526020016101a5565b34801561035c57600080fd5b5061017961036b366004611782565b610ba7565b34801561037c57600080fd5b5061017961038b366004611858565b610d2d565b34801561039c57600080fd5b506103b06103ab3660046116da565b610dcc565b6040516101a59190611882565b3480156103c957600080fd5b506101796103d83660046115fe565b610f76565b3480156103e957600080fd5b5061019b60055481565b3480156103ff57600080fd5b5061019b61040e3660046116da565b60086020526000908152604090205481565b34801561042c57600080fd5b5061019b61043b3660046116da565b6001600160a01b031660009081526008602052604090205490565b34801561046257600080fd5b506101796104713660046116da565b610fbf565b34801561048257600080fd5b5061019b60015481565b34801561049857600080fd5b506104ac6104a73660046116da565b611035565b6040516101a5959493929190611895565b6001600160a01b038216600090815260086020526040812080548392906104e59084906118f1565b90915550506001600160a01b038216600081815260086020908152604091829020548251858152918201527f56b34df61acb18dada28b541448a4ff3faf4c0970eb58b9980468a2c7538332291015b60405180910390a25050565b606081516001600160401b0381111561055b5761055b6113b5565b6040519080825280602002602001820160405280156105b357816020015b6040805160a08101825260008082526020808301829052928201526060808201819052608082015282526000199092019101816105795790505b50905060005b825181101561078d57600760008483815181106105d8576105d8611909565b6020908102919091018101516001600160a01b03168252818101929092526040908101600020815160a08101835281546001600160401b038082168352600160401b8204811695830195909552600160801b90049093169183019190915260018101805460608401919061064b9061191f565b80601f01602080910402602001604051908101604052809291908181526020018280546106779061191f565b80156106c45780601f10610699576101008083540402835291602001916106c4565b820191906000526020600020905b8154815290600101906020018083116106a757829003601f168201915b505050505081526020016002820180546106dd9061191f565b80601f01602080910402602001604051908101604052809291908181526020018280546107099061191f565b80156107565780601f1061072b57610100808354040283529160200191610756565b820191906000526020600020905b81548152906001019060200180831161073957829003601f168201915b50505050508152505082828151811061077157610771611909565b60200260200101819052508061078690611954565b90506105b9565b50919050565b61079b611188565b60045460408051918252602082018390527ff6b7e27129fefc5c7012885c85042d19038ba1abe27be1402518bafb73f55d1a910160405180910390a1600455565b6107e4611188565b60006107f86066546001600160a01b031690565b6001600160a01b03168260405160006040518083038185875af1925050503d8060008114610842576040519150601f19603f3d011682016040523d82523d6000602084013e610847565b606091505b505090508061085557600080fd5b5050565b610861611188565b60025460065461087191906118f1565b4210156108d15760405162461bcd60e51b815260206004820152602360248201527f70726963652075706461746520636f6f6c646f776e206e6f74207375727061736044820152621cd95960ea1b60648201526084015b60405180910390fd5b600054600154600254604080519384526020840187905283019190915260608201849052608082015260a081018290527f4e8b1ca1fe3c8cccee8ccb90aaf40042352429fa05c61c75b8e77836ed3af5519060c00160405180910390a160019190915560009190915560025542600655565b606081516001600160401b0381111561095e5761095e6113b5565b604051908082528060200260200182016040528015610987578160200160208202803683370190505b50905060005b825181101561078d57600860008483815181106109ac576109ac611909565b60200260200101516001600160a01b03166001600160a01b03168152602001908152602001600020548282815181106109e7576109e7611909565b60209081029190910101526109fb81611954565b905061098d565b610a0a611188565b610a1460006111e2565b565b610a1e611188565b60055460408051918252602082018390527fd8e1e8a19df284fc03b3fb160e79cd05bd9d732512fcfd248ffa3195433b932f910160405180910390a1600555565b603354610100900460ff1615808015610a7f5750603354600160ff909116105b80610a995750303b158015610a99575060335460ff166001145b610afc5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b60648201526084016108c8565b6033805460ff191660011790558015610b1f576033805461ff0019166101001790555b610b28886111e2565b600087905560038690556001859055600484905560028390556005829055426006558015610b90576033805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b5050505050505050565b610ba481346104bd565b50565b610baf611188565b610bc181606001518260800151611234565b80602001516001600160401b031681604001516001600160401b031611610c455760405162461bcd60e51b815260206004820152603260248201527f656e642074696d657374616d70206d75737420626520677265617465722074686044820152710616e2073746172742074696d657374616d760741b60648201526084016108c8565b6001600160a01b0382166000908152600760209081526040918290208351815483860151948601516001600160401b03908116600160801b0267ffffffffffffffff60801b19968216600160401b026fffffffffffffffffffffffffffffffff199093169190931617179390931692909217825560608301518051849392610cd492600185019291019061131c565b5060808201518051610cf091600284019160209091019061131c565b50905050816001600160a01b03167fff3054d138559c39b4c0826c43e94b2b2c6bc9a33ea1d0b74f16c916c7b73ec1826040516105349190611882565b610d35611188565b816001600160a01b031663a9059cbb610d566066546001600160a01b031690565b6040516001600160e01b031960e084901b1681526001600160a01b039091166004820152602481018490526044016020604051808303816000875af1158015610da3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610dc7919061196f565b505050565b6040805160a08082018352600080835260208084018290528385018290526060808501819052608085018190526001600160a01b038716835260078252918590208551938401865280546001600160401b038082168652600160401b8204811693860193909352600160801b9004909116948301949094526001840180549394929391840191610e5b9061191f565b80601f0160208091040260200160405190810160405280929190818152602001828054610e879061191f565b8015610ed45780601f10610ea957610100808354040283529160200191610ed4565b820191906000526020600020905b815481529060010190602001808311610eb757829003601f168201915b50505050508152602001600282018054610eed9061191f565b80601f0160208091040260200160405190810160405280929190818152602001828054610f199061191f565b8015610f665780601f10610f3b57610100808354040283529160200191610f66565b820191906000526020600020905b815481529060010190602001808311610f4957829003601f168201915b5050505050815250509050919050565b610f7e611188565b60035460408051918252602082018390527f33ec972d20d3bef7c2f239466a44f6c7335a4faf21f7b135c2b30138d71f8807910160405180910390a1600355565b610fc7611188565b6001600160a01b03811661102c5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b60648201526084016108c8565b610ba4816111e2565b600760205260009081526040902080546001820180546001600160401b0380841694600160401b8504821694600160801b90049091169290916110779061191f565b80601f01602080910402602001604051908101604052809291908181526020018280546110a39061191f565b80156110f05780601f106110c5576101008083540402835291602001916110f0565b820191906000526020600020905b8154815290600101906020018083116110d357829003601f168201915b5050505050908060020180546111059061191f565b80601f01602080910402602001604051908101604052809291908181526020018280546111319061191f565b801561117e5780601f106111535761010080835404028352916020019161117e565b820191906000526020600020905b81548152906001019060200180831161116157829003601f168201915b5050505050905085565b6066546001600160a01b03163314610a145760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016108c8565b606680546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b80518251146112855760405162461bcd60e51b815260206004820181905260248201527f617272617973206d7573742068617665207468652073616d65206c656e67746860448201526064016108c8565b6000805b82518110156112c8578281815181106112a4576112a4611909565b01602001516112b69060f81c83611991565b91506112c181611954565b9050611289565b508060ff16606414610dc75760405162461bcd60e51b815260206004820152601f60248201527f73756d206f662071756f72756d53706c697473206d757374206265203130300060448201526064016108c8565b8280546113289061191f565b90600052602060002090601f01602090048101928261134a5760008555611390565b82601f1061136357805160ff1916838001178555611390565b82800160010185558215611390579182015b82811115611390578251825591602001919060010190611375565b5061139c9291506113a0565b5090565b5b8082111561139c57600081556001016113a1565b634e487b7160e01b600052604160045260246000fd5b60405160a081016001600160401b03811182821017156113ed576113ed6113b5565b60405290565b604051601f8201601f191681016001600160401b038111828210171561141b5761141b6113b5565b604052919050565b80356001600160a01b038116811461143a57600080fd5b919050565b6000602080838503121561145257600080fd5b82356001600160401b038082111561146957600080fd5b818501915085601f83011261147d57600080fd5b81358181111561148f5761148f6113b5565b8060051b91506114a08483016113f3565b81815291830184019184810190888411156114ba57600080fd5b938501935b838510156114df576114d085611423565b825293850193908501906114bf565b98975050505050505050565b6000815180845260005b81811015611511576020818501810151868301820152016114f5565b81811115611523576000602083870101525b50601f01601f19169290920160200192915050565b60006001600160401b0380835116845280602084015116602085015280604084015116604085015250606082015160a0606085015261157a60a08501826114eb565b90506080830151848203608086015261159382826114eb565b95945050505050565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b828110156115f157603f198886030184526115df858351611538565b945092850192908501906001016115c3565b5092979650505050505050565b60006020828403121561161057600080fd5b5035919050565b60008060006060848603121561162c57600080fd5b505081359360208301359350604090920135919050565b6020808252825182820181905260009190848201906040850190845b8181101561167b5783518352928401929184019160010161165f565b50909695505050505050565b600080600080600080600060e0888a0312156116a257600080fd5b6116ab88611423565b9960208901359950604089013598606081013598506080810135975060a0810135965060c00135945092505050565b6000602082840312156116ec57600080fd5b6116f582611423565b9392505050565b80356001600160401b038116811461143a57600080fd5b600082601f83011261172457600080fd5b81356001600160401b0381111561173d5761173d6113b5565b611750601f8201601f19166020016113f3565b81815284602083860101111561176557600080fd5b816020850160208301376000918101602001919091529392505050565b6000806040838503121561179557600080fd5b61179e83611423565b915060208301356001600160401b03808211156117ba57600080fd5b9084019060a082870312156117ce57600080fd5b6117d66113cb565b6117df836116fc565b81526117ed602084016116fc565b60208201526117fe604084016116fc565b604082015260608301358281111561181557600080fd5b61182188828601611713565b60608301525060808301358281111561183957600080fd5b61184588828601611713565b6080830152508093505050509250929050565b6000806040838503121561186b57600080fd5b61187483611423565b946020939093013593505050565b6020815260006116f56020830184611538565b60006001600160401b038088168352808716602084015280861660408401525060a060608301526118c960a08301856114eb565b82810360808401526114df81856114eb565b634e487b7160e01b600052601160045260246000fd5b60008219821115611904576119046118db565b500190565b634e487b7160e01b600052603260045260246000fd5b600181811c9082168061193357607f821691505b6020821081141561078d57634e487b7160e01b600052602260045260246000fd5b6000600019821415611968576119686118db565b5060010190565b60006020828403121561198157600080fd5b815180151581146116f557600080fd5b600060ff821660ff84168060ff038211156119ae576119ae6118db565b01939250505056fea2646970667358221220b74e3b759f5411e05ca7f498a6de2d9e887ff3ff8670d064795dcdfaa9bbf3a964736f6c634300080c0033",
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

// GlobalRateBinInterval is a free data retrieval call binding the contract method 0xd882a5fd.
//
// Solidity: function globalRateBinInterval() view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultCaller) GlobalRateBinInterval(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "globalRateBinInterval")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GlobalRateBinInterval is a free data retrieval call binding the contract method 0xd882a5fd.
//
// Solidity: function globalRateBinInterval() view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultSession) GlobalRateBinInterval() (*big.Int, error) {
	return _ContractPaymentVault.Contract.GlobalRateBinInterval(&_ContractPaymentVault.CallOpts)
}

// GlobalRateBinInterval is a free data retrieval call binding the contract method 0xd882a5fd.
//
// Solidity: function globalRateBinInterval() view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) GlobalRateBinInterval() (*big.Int, error) {
	return _ContractPaymentVault.Contract.GlobalRateBinInterval(&_ContractPaymentVault.CallOpts)
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

// Initialize is a paid mutator transaction binding the contract method 0x8142951a.
//
// Solidity: function initialize(address _initialOwner, uint256 _minChargeableSize, uint256 _globalSymbolsPerSecond, uint256 _pricePerSymbol, uint256 _reservationBinInterval, uint256 _priceUpdateCooldown, uint256 _globalRateBinInterval) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) Initialize(opts *bind.TransactOpts, _initialOwner common.Address, _minChargeableSize *big.Int, _globalSymbolsPerSecond *big.Int, _pricePerSymbol *big.Int, _reservationBinInterval *big.Int, _priceUpdateCooldown *big.Int, _globalRateBinInterval *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "initialize", _initialOwner, _minChargeableSize, _globalSymbolsPerSecond, _pricePerSymbol, _reservationBinInterval, _priceUpdateCooldown, _globalRateBinInterval)
}

// Initialize is a paid mutator transaction binding the contract method 0x8142951a.
//
// Solidity: function initialize(address _initialOwner, uint256 _minChargeableSize, uint256 _globalSymbolsPerSecond, uint256 _pricePerSymbol, uint256 _reservationBinInterval, uint256 _priceUpdateCooldown, uint256 _globalRateBinInterval) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) Initialize(_initialOwner common.Address, _minChargeableSize *big.Int, _globalSymbolsPerSecond *big.Int, _pricePerSymbol *big.Int, _reservationBinInterval *big.Int, _priceUpdateCooldown *big.Int, _globalRateBinInterval *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.Initialize(&_ContractPaymentVault.TransactOpts, _initialOwner, _minChargeableSize, _globalSymbolsPerSecond, _pricePerSymbol, _reservationBinInterval, _priceUpdateCooldown, _globalRateBinInterval)
}

// Initialize is a paid mutator transaction binding the contract method 0x8142951a.
//
// Solidity: function initialize(address _initialOwner, uint256 _minChargeableSize, uint256 _globalSymbolsPerSecond, uint256 _pricePerSymbol, uint256 _reservationBinInterval, uint256 _priceUpdateCooldown, uint256 _globalRateBinInterval) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) Initialize(_initialOwner common.Address, _minChargeableSize *big.Int, _globalSymbolsPerSecond *big.Int, _pricePerSymbol *big.Int, _reservationBinInterval *big.Int, _priceUpdateCooldown *big.Int, _globalRateBinInterval *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.Initialize(&_ContractPaymentVault.TransactOpts, _initialOwner, _minChargeableSize, _globalSymbolsPerSecond, _pricePerSymbol, _reservationBinInterval, _priceUpdateCooldown, _globalRateBinInterval)
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

// SetGlobalRateBinInterval is a paid mutator transaction binding the contract method 0x7201dd02.
//
// Solidity: function setGlobalRateBinInterval(uint256 _globalRateBinInterval) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) SetGlobalRateBinInterval(opts *bind.TransactOpts, _globalRateBinInterval *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "setGlobalRateBinInterval", _globalRateBinInterval)
}

// SetGlobalRateBinInterval is a paid mutator transaction binding the contract method 0x7201dd02.
//
// Solidity: function setGlobalRateBinInterval(uint256 _globalRateBinInterval) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) SetGlobalRateBinInterval(_globalRateBinInterval *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetGlobalRateBinInterval(&_ContractPaymentVault.TransactOpts, _globalRateBinInterval)
}

// SetGlobalRateBinInterval is a paid mutator transaction binding the contract method 0x7201dd02.
//
// Solidity: function setGlobalRateBinInterval(uint256 _globalRateBinInterval) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) SetGlobalRateBinInterval(_globalRateBinInterval *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetGlobalRateBinInterval(&_ContractPaymentVault.TransactOpts, _globalRateBinInterval)
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

// SetPriceParams is a paid mutator transaction binding the contract method 0x327fad08.
//
// Solidity: function setPriceParams(uint256 _minChargeableSize, uint256 _pricePerSymbol, uint256 _priceUpdateCooldown) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) SetPriceParams(opts *bind.TransactOpts, _minChargeableSize *big.Int, _pricePerSymbol *big.Int, _priceUpdateCooldown *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "setPriceParams", _minChargeableSize, _pricePerSymbol, _priceUpdateCooldown)
}

// SetPriceParams is a paid mutator transaction binding the contract method 0x327fad08.
//
// Solidity: function setPriceParams(uint256 _minChargeableSize, uint256 _pricePerSymbol, uint256 _priceUpdateCooldown) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) SetPriceParams(_minChargeableSize *big.Int, _pricePerSymbol *big.Int, _priceUpdateCooldown *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetPriceParams(&_ContractPaymentVault.TransactOpts, _minChargeableSize, _pricePerSymbol, _priceUpdateCooldown)
}

// SetPriceParams is a paid mutator transaction binding the contract method 0x327fad08.
//
// Solidity: function setPriceParams(uint256 _minChargeableSize, uint256 _pricePerSymbol, uint256 _priceUpdateCooldown) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) SetPriceParams(_minChargeableSize *big.Int, _pricePerSymbol *big.Int, _priceUpdateCooldown *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetPriceParams(&_ContractPaymentVault.TransactOpts, _minChargeableSize, _pricePerSymbol, _priceUpdateCooldown)
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

// SetReservationBinInterval is a paid mutator transaction binding the contract method 0x2c1a33bc.
//
// Solidity: function setReservationBinInterval(uint256 _reservationBinInterval) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) SetReservationBinInterval(opts *bind.TransactOpts, _reservationBinInterval *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "setReservationBinInterval", _reservationBinInterval)
}

// SetReservationBinInterval is a paid mutator transaction binding the contract method 0x2c1a33bc.
//
// Solidity: function setReservationBinInterval(uint256 _reservationBinInterval) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) SetReservationBinInterval(_reservationBinInterval *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetReservationBinInterval(&_ContractPaymentVault.TransactOpts, _reservationBinInterval)
}

// SetReservationBinInterval is a paid mutator transaction binding the contract method 0x2c1a33bc.
//
// Solidity: function setReservationBinInterval(uint256 _reservationBinInterval) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) SetReservationBinInterval(_reservationBinInterval *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetReservationBinInterval(&_ContractPaymentVault.TransactOpts, _reservationBinInterval)
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

// ContractPaymentVaultGlobalRateBinIntervalUpdatedIterator is returned from FilterGlobalRateBinIntervalUpdated and is used to iterate over the raw logs and unpacked data for GlobalRateBinIntervalUpdated events raised by the ContractPaymentVault contract.
type ContractPaymentVaultGlobalRateBinIntervalUpdatedIterator struct {
	Event *ContractPaymentVaultGlobalRateBinIntervalUpdated // Event containing the contract specifics and raw log

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
func (it *ContractPaymentVaultGlobalRateBinIntervalUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPaymentVaultGlobalRateBinIntervalUpdated)
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
		it.Event = new(ContractPaymentVaultGlobalRateBinIntervalUpdated)
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
func (it *ContractPaymentVaultGlobalRateBinIntervalUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPaymentVaultGlobalRateBinIntervalUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPaymentVaultGlobalRateBinIntervalUpdated represents a GlobalRateBinIntervalUpdated event raised by the ContractPaymentVault contract.
type ContractPaymentVaultGlobalRateBinIntervalUpdated struct {
	PreviousValue *big.Int
	NewValue      *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterGlobalRateBinIntervalUpdated is a free log retrieval operation binding the contract event 0xd8e1e8a19df284fc03b3fb160e79cd05bd9d732512fcfd248ffa3195433b932f.
//
// Solidity: event GlobalRateBinIntervalUpdated(uint256 previousValue, uint256 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) FilterGlobalRateBinIntervalUpdated(opts *bind.FilterOpts) (*ContractPaymentVaultGlobalRateBinIntervalUpdatedIterator, error) {

	logs, sub, err := _ContractPaymentVault.contract.FilterLogs(opts, "GlobalRateBinIntervalUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultGlobalRateBinIntervalUpdatedIterator{contract: _ContractPaymentVault.contract, event: "GlobalRateBinIntervalUpdated", logs: logs, sub: sub}, nil
}

// WatchGlobalRateBinIntervalUpdated is a free log subscription operation binding the contract event 0xd8e1e8a19df284fc03b3fb160e79cd05bd9d732512fcfd248ffa3195433b932f.
//
// Solidity: event GlobalRateBinIntervalUpdated(uint256 previousValue, uint256 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) WatchGlobalRateBinIntervalUpdated(opts *bind.WatchOpts, sink chan<- *ContractPaymentVaultGlobalRateBinIntervalUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractPaymentVault.contract.WatchLogs(opts, "GlobalRateBinIntervalUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPaymentVaultGlobalRateBinIntervalUpdated)
				if err := _ContractPaymentVault.contract.UnpackLog(event, "GlobalRateBinIntervalUpdated", log); err != nil {
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

// ParseGlobalRateBinIntervalUpdated is a log parse operation binding the contract event 0xd8e1e8a19df284fc03b3fb160e79cd05bd9d732512fcfd248ffa3195433b932f.
//
// Solidity: event GlobalRateBinIntervalUpdated(uint256 previousValue, uint256 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) ParseGlobalRateBinIntervalUpdated(log types.Log) (*ContractPaymentVaultGlobalRateBinIntervalUpdated, error) {
	event := new(ContractPaymentVaultGlobalRateBinIntervalUpdated)
	if err := _ContractPaymentVault.contract.UnpackLog(event, "GlobalRateBinIntervalUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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
	PreviousMinChargeableSize   *big.Int
	NewMinChargeableSize        *big.Int
	PreviousPricePerSymbol      *big.Int
	NewPricePerSymbol           *big.Int
	PreviousPriceUpdateCooldown *big.Int
	NewPriceUpdateCooldown      *big.Int
	Raw                         types.Log // Blockchain specific contextual infos
}

// FilterPriceParamsUpdated is a free log retrieval operation binding the contract event 0x4e8b1ca1fe3c8cccee8ccb90aaf40042352429fa05c61c75b8e77836ed3af551.
//
// Solidity: event PriceParamsUpdated(uint256 previousMinChargeableSize, uint256 newMinChargeableSize, uint256 previousPricePerSymbol, uint256 newPricePerSymbol, uint256 previousPriceUpdateCooldown, uint256 newPriceUpdateCooldown)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) FilterPriceParamsUpdated(opts *bind.FilterOpts) (*ContractPaymentVaultPriceParamsUpdatedIterator, error) {

	logs, sub, err := _ContractPaymentVault.contract.FilterLogs(opts, "PriceParamsUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultPriceParamsUpdatedIterator{contract: _ContractPaymentVault.contract, event: "PriceParamsUpdated", logs: logs, sub: sub}, nil
}

// WatchPriceParamsUpdated is a free log subscription operation binding the contract event 0x4e8b1ca1fe3c8cccee8ccb90aaf40042352429fa05c61c75b8e77836ed3af551.
//
// Solidity: event PriceParamsUpdated(uint256 previousMinChargeableSize, uint256 newMinChargeableSize, uint256 previousPricePerSymbol, uint256 newPricePerSymbol, uint256 previousPriceUpdateCooldown, uint256 newPriceUpdateCooldown)
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

// ParsePriceParamsUpdated is a log parse operation binding the contract event 0x4e8b1ca1fe3c8cccee8ccb90aaf40042352429fa05c61c75b8e77836ed3af551.
//
// Solidity: event PriceParamsUpdated(uint256 previousMinChargeableSize, uint256 newMinChargeableSize, uint256 previousPricePerSymbol, uint256 newPricePerSymbol, uint256 previousPriceUpdateCooldown, uint256 newPriceUpdateCooldown)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) ParsePriceParamsUpdated(log types.Log) (*ContractPaymentVaultPriceParamsUpdated, error) {
	event := new(ContractPaymentVaultPriceParamsUpdated)
	if err := _ContractPaymentVault.contract.UnpackLog(event, "PriceParamsUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractPaymentVaultReservationBinIntervalUpdatedIterator is returned from FilterReservationBinIntervalUpdated and is used to iterate over the raw logs and unpacked data for ReservationBinIntervalUpdated events raised by the ContractPaymentVault contract.
type ContractPaymentVaultReservationBinIntervalUpdatedIterator struct {
	Event *ContractPaymentVaultReservationBinIntervalUpdated // Event containing the contract specifics and raw log

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
func (it *ContractPaymentVaultReservationBinIntervalUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPaymentVaultReservationBinIntervalUpdated)
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
		it.Event = new(ContractPaymentVaultReservationBinIntervalUpdated)
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
func (it *ContractPaymentVaultReservationBinIntervalUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPaymentVaultReservationBinIntervalUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPaymentVaultReservationBinIntervalUpdated represents a ReservationBinIntervalUpdated event raised by the ContractPaymentVault contract.
type ContractPaymentVaultReservationBinIntervalUpdated struct {
	PreviousValue *big.Int
	NewValue      *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterReservationBinIntervalUpdated is a free log retrieval operation binding the contract event 0xf6b7e27129fefc5c7012885c85042d19038ba1abe27be1402518bafb73f55d1a.
//
// Solidity: event ReservationBinIntervalUpdated(uint256 previousValue, uint256 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) FilterReservationBinIntervalUpdated(opts *bind.FilterOpts) (*ContractPaymentVaultReservationBinIntervalUpdatedIterator, error) {

	logs, sub, err := _ContractPaymentVault.contract.FilterLogs(opts, "ReservationBinIntervalUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultReservationBinIntervalUpdatedIterator{contract: _ContractPaymentVault.contract, event: "ReservationBinIntervalUpdated", logs: logs, sub: sub}, nil
}

// WatchReservationBinIntervalUpdated is a free log subscription operation binding the contract event 0xf6b7e27129fefc5c7012885c85042d19038ba1abe27be1402518bafb73f55d1a.
//
// Solidity: event ReservationBinIntervalUpdated(uint256 previousValue, uint256 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) WatchReservationBinIntervalUpdated(opts *bind.WatchOpts, sink chan<- *ContractPaymentVaultReservationBinIntervalUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractPaymentVault.contract.WatchLogs(opts, "ReservationBinIntervalUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPaymentVaultReservationBinIntervalUpdated)
				if err := _ContractPaymentVault.contract.UnpackLog(event, "ReservationBinIntervalUpdated", log); err != nil {
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

// ParseReservationBinIntervalUpdated is a log parse operation binding the contract event 0xf6b7e27129fefc5c7012885c85042d19038ba1abe27be1402518bafb73f55d1a.
//
// Solidity: event ReservationBinIntervalUpdated(uint256 previousValue, uint256 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) ParseReservationBinIntervalUpdated(log types.Log) (*ContractPaymentVaultReservationBinIntervalUpdated, error) {
	event := new(ContractPaymentVaultReservationBinIntervalUpdated)
	if err := _ContractPaymentVault.contract.UnpackLog(event, "ReservationBinIntervalUpdated", log); err != nil {
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
