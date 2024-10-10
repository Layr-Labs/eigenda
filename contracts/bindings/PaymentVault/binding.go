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
	DataRate       uint64
	StartTimestamp uint64
	EndTimestamp   uint64
	QuorumNumbers  []byte
	QuorumSplits   []byte
}

// ContractPaymentVaultMetaData contains all meta data concerning the ContractPaymentVault contract.
var ContractPaymentVaultMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_reservationBinInterval\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_reservationBinStartTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_priceUpdateCooldown\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositOnDemand\",\"inputs\":[{\"name\":\"_account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"getOnDemandAmount\",\"inputs\":[{\"name\":\"_account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOnDemandAmounts\",\"inputs\":[{\"name\":\"_accounts\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[{\"name\":\"_payments\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReservation\",\"inputs\":[{\"name\":\"_account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIPaymentVault.Reservation\",\"components\":[{\"name\":\"dataRate\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReservations\",\"inputs\":[{\"name\":\"_accounts\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[{\"name\":\"_reservations\",\"type\":\"tuple[]\",\"internalType\":\"structIPaymentVault.Reservation[]\",\"components\":[{\"name\":\"dataRate\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"globalBytesPerSecond\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_minChargeableSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_globalBytesPerSecond\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"lastPriceUpdateTime\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minChargeableSize\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"onDemandPayments\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"priceUpdateCooldown\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"reservationBinInterval\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"reservationBinStartTimestamp\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"reservations\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"dataRate\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setGlobalBytesPerSec\",\"inputs\":[{\"name\":\"_globalBytesPerSecond\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinChargeableSize\",\"inputs\":[{\"name\":\"_minChargeableSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setReservation\",\"inputs\":[{\"name\":\"_account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_reservation\",\"type\":\"tuple\",\"internalType\":\"structIPaymentVault.Reservation\",\"components\":[{\"name\":\"dataRate\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"GlobalBytesPerSecondUpdated\",\"inputs\":[{\"name\":\"previousValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinChargeableSizeUpdated\",\"inputs\":[{\"name\":\"previousValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OnDemandPaymentUpdated\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"onDemandPayment\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"totalDeposit\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReservationUpdated\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"reservation\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIPaymentVault.Reservation\",\"components\":[{\"name\":\"dataRate\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"anonymous\":false}]",
	Bin: "0x60e060405234801561001057600080fd5b5060405161187b38038061187b83398101604081905261002f9161010e565b608083905260a082905260c081905261004661004e565b50505061013c565b600554610100900460ff16156100ba5760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840160405180910390fd5b60055460ff908116101561010c576005805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b60008060006060848603121561012357600080fd5b8351925060208401519150604084015190509250925092565b60805160a05160c051611709610172600039600081816101410152610ca0015260006102460152600061027a01526117096000f3fe60806040526004361061012a5760003560e01c80637a1ac61e116100ab578063d996dc991161006f578063d996dc9914610359578063e30c880614610386578063e60a8dd7146103a6578063efb435f8146103c6578063f2fde38b146103fc578063fd3dc53a1461041c57600080fd5b80637a1ac61e146102b15780638bec7d02146102d15780638da5cb5b146102e45780639aec86401461030c578063b2066f801461032c57600080fd5b80634486bfb7116100f25780634486bfb7146101f157806349b9a7af1461021e578063550571b4146102345780635a8a686914610268578063715018a61461029c57600080fd5b8063039f091c1461012f578063109f8fe51461017657806312a76a20146101a357806312cc61b8146101b95780632e1a7d4d146101cf575b600080fd5b34801561013b57600080fd5b506101637f000000000000000000000000000000000000000000000000000000000000000081565b6040519081526020015b60405180910390f35b34801561018257600080fd5b506101966101913660046111f4565b61044d565b60405161016d9190611351565b3480156101af57600080fd5b5061016360005481565b3480156101c557600080fd5b5061016360015481565b3480156101db57600080fd5b506101ef6101ea3660046113b3565b6106a0565b005b3480156101fd57600080fd5b5061021161020c3660046111f4565b61071d565b60405161016d91906113cc565b34801561022a57600080fd5b5061016360025481565b34801561024057600080fd5b506101637f000000000000000000000000000000000000000000000000000000000000000081565b34801561027457600080fd5b506101637f000000000000000000000000000000000000000000000000000000000000000081565b3480156102a857600080fd5b506101ef6107dc565b3480156102bd57600080fd5b506101ef6102cc366004611410565b6107f0565b6101ef6102df366004611443565b610914565b3480156102f057600080fd5b506038546040516001600160a01b03909116815260200161016d565b34801561031857600080fd5b506101ef6103273660046114eb565b610995565b34801561033857600080fd5b5061034c610347366004611443565b610aa3565b60405161016d91906115c1565b34801561036557600080fd5b50610163610374366004611443565b60046020526000908152604090205481565b34801561039257600080fd5b506101ef6103a13660046113b3565b610c4d565b3480156103b257600080fd5b506101ef6103c13660046113b3565b610c96565b3480156103d257600080fd5b506101636103e1366004611443565b6001600160a01b031660009081526004602052604090205490565b34801561040857600080fd5b506101ef610417366004611443565b610d6c565b34801561042857600080fd5b5061043c610437366004611443565b610de5565b60405161016d9594939291906115d4565b606081516001600160401b038111156104685761046861116a565b6040519080825280602002602001820160405280156104c057816020015b6040805160a08101825260008082526020808301829052928201526060808201819052608082015282526000199092019101816104865790505b50905060005b825181101561069a57600360008483815181106104e5576104e561161a565b6020908102919091018101516001600160a01b03168252818101929092526040908101600020815160a08101835281546001600160401b038082168352600160401b8204811695830195909552600160801b90049093169183019190915260018101805460608401919061055890611630565b80601f016020809104026020016040519081016040528092919081815260200182805461058490611630565b80156105d15780601f106105a6576101008083540402835291602001916105d1565b820191906000526020600020905b8154815290600101906020018083116105b457829003601f168201915b505050505081526020016002820180546105ea90611630565b80601f016020809104026020016040519081016040528092919081815260200182805461061690611630565b80156106635780601f1061063857610100808354040283529160200191610663565b820191906000526020600020905b81548152906001019060200180831161064657829003601f168201915b50505050508152505082828151811061067e5761067e61161a565b6020026020010181905250806106939061167b565b90506104c6565b50919050565b6106a8610f38565b60006106bc6038546001600160a01b031690565b6001600160a01b03168260405160006040518083038185875af1925050503d8060008114610706576040519150601f19603f3d011682016040523d82523d6000602084013e61070b565b606091505b505090508061071957600080fd5b5050565b606081516001600160401b038111156107385761073861116a565b604051908082528060200260200182016040528015610761578160200160208202803683370190505b50905060005b825181101561069a57600460008483815181106107865761078661161a565b60200260200101516001600160a01b03166001600160a01b03168152602001908152602001600020548282815181106107c1576107c161161a565b60209081029190910101526107d58161167b565b9050610767565b6107e4610f38565b6107ee6000610f92565b565b600554610100900460ff16158080156108105750600554600160ff909116105b8061082a5750303b15801561082a575060055460ff166001145b6108925760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b60648201526084015b60405180910390fd5b6005805460ff1916600117905580156108b5576005805461ff0019166101001790555b6108be84610d6c565b60008390556001829055801561090e576005805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50505050565b6001600160a01b0381166000908152600460205260408120805434929061093c908490611696565b90915550506001600160a01b038116600081815260046020908152604091829020548251348152918201527f56b34df61acb18dada28b541448a4ff3faf4c0970eb58b9980468a2c75383322910160405180910390a250565b61099d610f38565b6109af81606001518260800151610fe4565b6001600160a01b0382166000908152600360209081526040918290208351815483860151948601516001600160401b03908116600160801b0267ffffffffffffffff60801b19968216600160401b026fffffffffffffffffffffffffffffffff199093169190931617179390931692909217825560608301518051849392610a3e9260018501929101906110d1565b5060808201518051610a5a9160028401916020909101906110d1565b50905050816001600160a01b03167fff3054d138559c39b4c0826c43e94b2b2c6bc9a33ea1d0b74f16c916c7b73ec182604051610a9791906115c1565b60405180910390a25050565b6040805160a08082018352600080835260208084018290528385018290526060808501819052608085018190526001600160a01b038716835260038252918590208551938401865280546001600160401b038082168652600160401b8204811693860193909352600160801b9004909116948301949094526001840180549394929391840191610b3290611630565b80601f0160208091040260200160405190810160405280929190818152602001828054610b5e90611630565b8015610bab5780601f10610b8057610100808354040283529160200191610bab565b820191906000526020600020905b815481529060010190602001808311610b8e57829003601f168201915b50505050508152602001600282018054610bc490611630565b80601f0160208091040260200160405190810160405280929190818152602001828054610bf090611630565b8015610c3d5780601f10610c1257610100808354040283529160200191610c3d565b820191906000526020600020905b815481529060010190602001808311610c2057829003601f168201915b5050505050815250509050919050565b610c55610f38565b60015460408051918252602082018390527f8a2a96caad50bb77e6c390a954ed96a23dc59b9c8a7e2c6fdf243b34559c2346910160405180910390a1600155565b610c9e610f38565b7f0000000000000000000000000000000000000000000000000000000000000000600254610ccc9190611696565b421015610d275760405162461bcd60e51b815260206004820152602360248201527f70726963652075706461746520636f6f6c646f776e206e6f74207375727061736044820152621cd95960ea1b6064820152608401610889565b60005460408051918252602082018390527f62caca228682fc3cff59bcae8bdf562027847c9295336694d56a0892bbeed0b9910160405180910390a142600255600055565b610d74610f38565b6001600160a01b038116610dd95760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610889565b610de281610f92565b50565b600360205260009081526040902080546001820180546001600160401b0380841694600160401b8504821694600160801b9004909116929091610e2790611630565b80601f0160208091040260200160405190810160405280929190818152602001828054610e5390611630565b8015610ea05780601f10610e7557610100808354040283529160200191610ea0565b820191906000526020600020905b815481529060010190602001808311610e8357829003601f168201915b505050505090806002018054610eb590611630565b80601f0160208091040260200160405190810160405280929190818152602001828054610ee190611630565b8015610f2e5780601f10610f0357610100808354040283529160200191610f2e565b820191906000526020600020905b815481529060010190602001808311610f1157829003601f168201915b5050505050905085565b6038546001600160a01b031633146107ee5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610889565b603880546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b80518251146110355760405162461bcd60e51b815260206004820181905260248201527f617272617973206d7573742068617665207468652073616d65206c656e6774686044820152606401610889565b6000805b8251811015611078578281815181106110545761105461161a565b01602001516110669060f81c836116ae565b91506110718161167b565b9050611039565b508060ff166064146110cc5760405162461bcd60e51b815260206004820152601f60248201527f73756d206f662071756f72756d53706c697473206d75737420626520313030006044820152606401610889565b505050565b8280546110dd90611630565b90600052602060002090601f0160209004810192826110ff5760008555611145565b82601f1061111857805160ff1916838001178555611145565b82800160010185558215611145579182015b8281111561114557825182559160200191906001019061112a565b50611151929150611155565b5090565b5b808211156111515760008155600101611156565b634e487b7160e01b600052604160045260246000fd5b60405160a081016001600160401b03811182821017156111a2576111a261116a565b60405290565b604051601f8201601f191681016001600160401b03811182821017156111d0576111d061116a565b604052919050565b80356001600160a01b03811681146111ef57600080fd5b919050565b6000602080838503121561120757600080fd5b82356001600160401b038082111561121e57600080fd5b818501915085601f83011261123257600080fd5b8135818111156112445761124461116a565b8060051b91506112558483016111a8565b818152918301840191848101908884111561126f57600080fd5b938501935b8385101561129457611285856111d8565b82529385019390850190611274565b98975050505050505050565b6000815180845260005b818110156112c6576020818501810151868301820152016112aa565b818111156112d8576000602083870101525b50601f01601f19169290920160200192915050565b60006001600160401b0380835116845280602084015116602085015280604084015116604085015250606082015160a0606085015261132f60a08501826112a0565b90506080830151848203608086015261134882826112a0565b95945050505050565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b828110156113a657603f198886030184526113948583516112ed565b94509285019290850190600101611378565b5092979650505050505050565b6000602082840312156113c557600080fd5b5035919050565b6020808252825182820181905260009190848201906040850190845b81811015611404578351835292840192918401916001016113e8565b50909695505050505050565b60008060006060848603121561142557600080fd5b61142e846111d8565b95602085013595506040909401359392505050565b60006020828403121561145557600080fd5b61145e826111d8565b9392505050565b80356001600160401b03811681146111ef57600080fd5b600082601f83011261148d57600080fd5b81356001600160401b038111156114a6576114a661116a565b6114b9601f8201601f19166020016111a8565b8181528460208386010111156114ce57600080fd5b816020850160208301376000918101602001919091529392505050565b600080604083850312156114fe57600080fd5b611507836111d8565b915060208301356001600160401b038082111561152357600080fd5b9084019060a0828703121561153757600080fd5b61153f611180565b61154883611465565b815261155660208401611465565b602082015261156760408401611465565b604082015260608301358281111561157e57600080fd5b61158a8882860161147c565b6060830152506080830135828111156115a257600080fd5b6115ae8882860161147c565b6080830152508093505050509250929050565b60208152600061145e60208301846112ed565b60006001600160401b038088168352808716602084015280861660408401525060a0606083015261160860a08301856112a0565b828103608084015261129481856112a0565b634e487b7160e01b600052603260045260246000fd5b600181811c9082168061164457607f821691505b6020821081141561069a57634e487b7160e01b600052602260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b600060001982141561168f5761168f611665565b5060010190565b600082198211156116a9576116a9611665565b500190565b600060ff821660ff84168060ff038211156116cb576116cb611665565b01939250505056fea2646970667358221220507c3b0360e07502c3577854b76e10f776b4af1bb404b3b0bbb432c98c23d7fa64736f6c634300080c0033",
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

// GlobalBytesPerSecond is a free data retrieval call binding the contract method 0x12cc61b8.
//
// Solidity: function globalBytesPerSecond() view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultCaller) GlobalBytesPerSecond(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "globalBytesPerSecond")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GlobalBytesPerSecond is a free data retrieval call binding the contract method 0x12cc61b8.
//
// Solidity: function globalBytesPerSecond() view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultSession) GlobalBytesPerSecond() (*big.Int, error) {
	return _ContractPaymentVault.Contract.GlobalBytesPerSecond(&_ContractPaymentVault.CallOpts)
}

// GlobalBytesPerSecond is a free data retrieval call binding the contract method 0x12cc61b8.
//
// Solidity: function globalBytesPerSecond() view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) GlobalBytesPerSecond() (*big.Int, error) {
	return _ContractPaymentVault.Contract.GlobalBytesPerSecond(&_ContractPaymentVault.CallOpts)
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
// Solidity: function reservations(address ) view returns(uint64 dataRate, uint64 startTimestamp, uint64 endTimestamp, bytes quorumNumbers, bytes quorumSplits)
func (_ContractPaymentVault *ContractPaymentVaultCaller) Reservations(opts *bind.CallOpts, arg0 common.Address) (struct {
	DataRate       uint64
	StartTimestamp uint64
	EndTimestamp   uint64
	QuorumNumbers  []byte
	QuorumSplits   []byte
}, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "reservations", arg0)

	outstruct := new(struct {
		DataRate       uint64
		StartTimestamp uint64
		EndTimestamp   uint64
		QuorumNumbers  []byte
		QuorumSplits   []byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.DataRate = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.StartTimestamp = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.EndTimestamp = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.QuorumNumbers = *abi.ConvertType(out[3], new([]byte)).(*[]byte)
	outstruct.QuorumSplits = *abi.ConvertType(out[4], new([]byte)).(*[]byte)

	return *outstruct, err

}

// Reservations is a free data retrieval call binding the contract method 0xfd3dc53a.
//
// Solidity: function reservations(address ) view returns(uint64 dataRate, uint64 startTimestamp, uint64 endTimestamp, bytes quorumNumbers, bytes quorumSplits)
func (_ContractPaymentVault *ContractPaymentVaultSession) Reservations(arg0 common.Address) (struct {
	DataRate       uint64
	StartTimestamp uint64
	EndTimestamp   uint64
	QuorumNumbers  []byte
	QuorumSplits   []byte
}, error) {
	return _ContractPaymentVault.Contract.Reservations(&_ContractPaymentVault.CallOpts, arg0)
}

// Reservations is a free data retrieval call binding the contract method 0xfd3dc53a.
//
// Solidity: function reservations(address ) view returns(uint64 dataRate, uint64 startTimestamp, uint64 endTimestamp, bytes quorumNumbers, bytes quorumSplits)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) Reservations(arg0 common.Address) (struct {
	DataRate       uint64
	StartTimestamp uint64
	EndTimestamp   uint64
	QuorumNumbers  []byte
	QuorumSplits   []byte
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

// Initialize is a paid mutator transaction binding the contract method 0x7a1ac61e.
//
// Solidity: function initialize(address _initialOwner, uint256 _minChargeableSize, uint256 _globalBytesPerSecond) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) Initialize(opts *bind.TransactOpts, _initialOwner common.Address, _minChargeableSize *big.Int, _globalBytesPerSecond *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "initialize", _initialOwner, _minChargeableSize, _globalBytesPerSecond)
}

// Initialize is a paid mutator transaction binding the contract method 0x7a1ac61e.
//
// Solidity: function initialize(address _initialOwner, uint256 _minChargeableSize, uint256 _globalBytesPerSecond) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) Initialize(_initialOwner common.Address, _minChargeableSize *big.Int, _globalBytesPerSecond *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.Initialize(&_ContractPaymentVault.TransactOpts, _initialOwner, _minChargeableSize, _globalBytesPerSecond)
}

// Initialize is a paid mutator transaction binding the contract method 0x7a1ac61e.
//
// Solidity: function initialize(address _initialOwner, uint256 _minChargeableSize, uint256 _globalBytesPerSecond) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) Initialize(_initialOwner common.Address, _minChargeableSize *big.Int, _globalBytesPerSecond *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.Initialize(&_ContractPaymentVault.TransactOpts, _initialOwner, _minChargeableSize, _globalBytesPerSecond)
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

// SetGlobalBytesPerSec is a paid mutator transaction binding the contract method 0xe30c8806.
//
// Solidity: function setGlobalBytesPerSec(uint256 _globalBytesPerSecond) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) SetGlobalBytesPerSec(opts *bind.TransactOpts, _globalBytesPerSecond *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "setGlobalBytesPerSec", _globalBytesPerSecond)
}

// SetGlobalBytesPerSec is a paid mutator transaction binding the contract method 0xe30c8806.
//
// Solidity: function setGlobalBytesPerSec(uint256 _globalBytesPerSecond) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) SetGlobalBytesPerSec(_globalBytesPerSecond *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetGlobalBytesPerSec(&_ContractPaymentVault.TransactOpts, _globalBytesPerSecond)
}

// SetGlobalBytesPerSec is a paid mutator transaction binding the contract method 0xe30c8806.
//
// Solidity: function setGlobalBytesPerSec(uint256 _globalBytesPerSecond) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) SetGlobalBytesPerSec(_globalBytesPerSecond *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetGlobalBytesPerSec(&_ContractPaymentVault.TransactOpts, _globalBytesPerSecond)
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

// ContractPaymentVaultGlobalBytesPerSecondUpdatedIterator is returned from FilterGlobalBytesPerSecondUpdated and is used to iterate over the raw logs and unpacked data for GlobalBytesPerSecondUpdated events raised by the ContractPaymentVault contract.
type ContractPaymentVaultGlobalBytesPerSecondUpdatedIterator struct {
	Event *ContractPaymentVaultGlobalBytesPerSecondUpdated // Event containing the contract specifics and raw log

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
func (it *ContractPaymentVaultGlobalBytesPerSecondUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPaymentVaultGlobalBytesPerSecondUpdated)
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
		it.Event = new(ContractPaymentVaultGlobalBytesPerSecondUpdated)
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
func (it *ContractPaymentVaultGlobalBytesPerSecondUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPaymentVaultGlobalBytesPerSecondUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPaymentVaultGlobalBytesPerSecondUpdated represents a GlobalBytesPerSecondUpdated event raised by the ContractPaymentVault contract.
type ContractPaymentVaultGlobalBytesPerSecondUpdated struct {
	PreviousValue *big.Int
	NewValue      *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterGlobalBytesPerSecondUpdated is a free log retrieval operation binding the contract event 0x8a2a96caad50bb77e6c390a954ed96a23dc59b9c8a7e2c6fdf243b34559c2346.
//
// Solidity: event GlobalBytesPerSecondUpdated(uint256 previousValue, uint256 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) FilterGlobalBytesPerSecondUpdated(opts *bind.FilterOpts) (*ContractPaymentVaultGlobalBytesPerSecondUpdatedIterator, error) {

	logs, sub, err := _ContractPaymentVault.contract.FilterLogs(opts, "GlobalBytesPerSecondUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultGlobalBytesPerSecondUpdatedIterator{contract: _ContractPaymentVault.contract, event: "GlobalBytesPerSecondUpdated", logs: logs, sub: sub}, nil
}

// WatchGlobalBytesPerSecondUpdated is a free log subscription operation binding the contract event 0x8a2a96caad50bb77e6c390a954ed96a23dc59b9c8a7e2c6fdf243b34559c2346.
//
// Solidity: event GlobalBytesPerSecondUpdated(uint256 previousValue, uint256 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) WatchGlobalBytesPerSecondUpdated(opts *bind.WatchOpts, sink chan<- *ContractPaymentVaultGlobalBytesPerSecondUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractPaymentVault.contract.WatchLogs(opts, "GlobalBytesPerSecondUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPaymentVaultGlobalBytesPerSecondUpdated)
				if err := _ContractPaymentVault.contract.UnpackLog(event, "GlobalBytesPerSecondUpdated", log); err != nil {
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

// ParseGlobalBytesPerSecondUpdated is a log parse operation binding the contract event 0x8a2a96caad50bb77e6c390a954ed96a23dc59b9c8a7e2c6fdf243b34559c2346.
//
// Solidity: event GlobalBytesPerSecondUpdated(uint256 previousValue, uint256 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) ParseGlobalBytesPerSecondUpdated(log types.Log) (*ContractPaymentVaultGlobalBytesPerSecondUpdated, error) {
	event := new(ContractPaymentVaultGlobalBytesPerSecondUpdated)
	if err := _ContractPaymentVault.contract.UnpackLog(event, "GlobalBytesPerSecondUpdated", log); err != nil {
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
