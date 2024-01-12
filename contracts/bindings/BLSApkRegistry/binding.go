// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractBLSApkRegistry

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

// BN254G1Point is an auto generated low-level Go binding around an user-defined struct.
type BN254G1Point struct {
	X *big.Int
	Y *big.Int
}

// BN254G2Point is an auto generated low-level Go binding around an user-defined struct.
type BN254G2Point struct {
	X [2]*big.Int
	Y [2]*big.Int
}

// IBLSApkRegistryApkUpdate is an auto generated low-level Go binding around an user-defined struct.
type IBLSApkRegistryApkUpdate struct {
	ApkHash               [24]byte
	UpdateBlockNumber     uint32
	NextUpdateBlockNumber uint32
}

// IBLSApkRegistryPubkeyRegistrationParams is an auto generated low-level Go binding around an user-defined struct.
type IBLSApkRegistryPubkeyRegistrationParams struct {
	PubkeyRegistrationSignature BN254G1Point
	PubkeyG1                    BN254G1Point
	PubkeyG2                    BN254G2Point
}

// ContractBLSApkRegistryMetaData contains all meta data concerning the ContractBLSApkRegistry contract.
var ContractBLSApkRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"apkHistory\",\"inputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"apkHash\",\"type\":\"bytes24\",\"internalType\":\"bytes24\"},{\"name\":\"updateBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"nextUpdateBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"currentApk\",\"inputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deregisterOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getApk\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApkHashAtBlockNumberAndIndex\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"blockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes24\",\"internalType\":\"bytes24\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApkHistoryLength\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApkIndicesAtBlockNumber\",\"inputs\":[{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApkUpdateAtIndex\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIBLSApkRegistry.ApkUpdate\",\"components\":[{\"name\":\"apkHash\",\"type\":\"bytes24\",\"internalType\":\"bytes24\"},{\"name\":\"updateBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"nextUpdateBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorFromPubkeyHash\",\"inputs\":[{\"name\":\"pubkeyHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorId\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRegisteredPubkey\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initializeQuorum\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"operatorToPubkey\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operatorToPubkeyHash\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pubkeyHashToOperator\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerBLSPublicKey\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIBLSApkRegistry.PubkeyRegistrationParams\",\"components\":[{\"name\":\"pubkeyRegistrationSignature\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"pubkeyG1\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"pubkeyG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]}]},{\"name\":\"pubkeyRegistrationMessageHash\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registryCoordinator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NewPubkeyRegistration\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"pubkeyG1\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"pubkeyG2\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorAddedToQuorums\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorRemovedFromQuorums\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false}]",
	Bin: "0x60a06040523480156200001157600080fd5b506040516200213738038062002137833981016040819052620000349162000116565b6001600160a01b038116608052806200004c62000054565b505062000148565b600054610100900460ff1615620000c15760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840160405180910390fd5b60005460ff908116101562000114576000805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b6000602082840312156200012957600080fd5b81516001600160a01b03811681146200014157600080fd5b9392505050565b608051611fb7620001806000396000818161030f01528181610466015281816105bf015281816109a501526110420152611fb76000f3fe608060405234801561001057600080fd5b50600436106101155760003560e01c80636d14a987116100a2578063bf79ce5811610071578063bf79ce58146103cc578063d5254a8c146103df578063de29fac0146103ff578063e8bb9ae61461041f578063f4e24fe51461044857600080fd5b80636d14a9871461030a5780637916cea6146103315780637ff81a8714610372578063a3db80e2146103a557600080fd5b80633fb27952116100e95780633fb27952146101df57806347b314e8146101f25780635f61a88414610233578063605747d51461028f57806368bccaac146102dd57600080fd5b8062a1f4cb1461011a57806313542a4e1461015b57806326d941f214610192578063377ed99d146101a7575b600080fd5b610141610128366004611905565b6003602052600090815260409020805460019091015482565b604080519283526020830191909152015b60405180910390f35b610184610169366004611905565b6001600160a01b031660009081526001602052604090205490565b604051908152602001610152565b6101a56101a0366004611938565b61045b565b005b6101ca6101b5366004611938565b60ff1660009081526004602052604090205490565b60405163ffffffff9091168152602001610152565b6101a56101ed3660046119c3565b6105b4565b61021b610200366004611a69565b6000908152600260205260409020546001600160a01b031690565b6040516001600160a01b039091168152602001610152565b610282610241366004611938565b60408051808201909152600080825260208201525060ff16600090815260056020908152604091829020825180840190935280548352600101549082015290565b6040516101529190611a82565b6102a261029d366004611a99565b610652565b60408051825167ffffffffffffffff1916815260208084015163ffffffff908116918301919091529282015190921690820152606001610152565b6102f06102eb366004611ac3565b6106e5565b60405167ffffffffffffffff199091168152602001610152565b61021b7f000000000000000000000000000000000000000000000000000000000000000081565b61034461033f366004611a99565b610880565b6040805167ffffffffffffffff19909416845263ffffffff9283166020850152911690820152606001610152565b610385610380366004611905565b6108cb565b604080518351815260209384015193810193909352820152606001610152565b6101416103b3366004611938565b6005602052600090815260409020805460019091015482565b6101846103da366004611b0b565b610998565b6103f26103ed366004611b68565b610dec565b6040516101529190611be0565b61018461040d366004611905565b60016020526000908152604090205481565b61021b61042d366004611a69565b6002602052600090815260409020546001600160a01b031681565b6101a56104563660046119c3565b611037565b336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146104ac5760405162461bcd60e51b81526004016104a390611c2a565b60405180910390fd5b60ff81166000908152600460205260409020541561052b5760405162461bcd60e51b815260206004820152603660248201527f424c5341706b52656769737472792e696e697469616c697a6551756f72756d3a6044820152752071756f72756d20616c72656164792065786973747360501b60648201526084016104a3565b60ff166000908152600460209081526040808320815160608101835284815263ffffffff4381168286019081528285018781528454600181018655948852959096209151919092018054955194518316600160e01b026001600160e01b0395909316600160c01b026001600160e01b03199096169190931c179390931791909116919091179055565b336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146105fc5760405162461bcd60e51b81526004016104a390611c2a565b6000610607836108cb565b50905061061482826110d0565b7f5358c5b42179178c8fc757734ac2a3198f9071c765ee0d8389211525f50052468383604051610645929190611c9e565b60405180910390a1505050565b604080516060810182526000808252602080830182905282840182905260ff86168252600490529190912080548390811061068f5761068f611d03565b600091825260209182902060408051606081018252919092015467ffffffffffffffff1981841b16825263ffffffff600160c01b8204811694830194909452600160e01b90049092169082015290505b92915050565b60ff8316600090815260046020526040812080548291908490811061070c5761070c611d03565b600091825260209182902060408051606081018252919092015467ffffffffffffffff1981841b16825263ffffffff600160c01b82048116948301859052600160e01b9091048116928201929092529250851610156107d35760405162461bcd60e51b815260206004820152603e60248201527f424c5341706b52656769737472792e5f76616c696461746541706b486173684160448201527f74426c6f636b4e756d6265723a20696e64657820746f6f20726563656e74000060648201526084016104a3565b604081015163ffffffff1615806107f95750806040015163ffffffff168463ffffffff16105b6108775760405162461bcd60e51b815260206004820152604360248201527f424c5341706b52656769737472792e5f76616c696461746541706b486173684160448201527f74426c6f636b4e756d6265723a206e6f74206c61746573742061706b2075706460648201526261746560e81b608482015260a4016104a3565b51949350505050565b6004602052816000526040600020818154811061089c57600080fd5b600091825260209091200154604081901b925063ffffffff600160c01b820481169250600160e01b9091041683565b60408051808201909152600080825260208201526001600160a01b03821660008181526003602090815260408083208151808301835281548152600191820154818501529484529091528120549091908061098e5760405162461bcd60e51b815260206004820152603e60248201527f424c5341706b52656769737472792e676574526567697374657265645075626b60448201527f65793a206f70657261746f72206973206e6f742072656769737465726564000060648201526084016104a3565b9094909350915050565b6000336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146109e25760405162461bcd60e51b81526004016104a390611c2a565b6000610a106109f936869003860160408701611d19565b805160009081526020918201519091526040902090565b90507fad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5811415610a98576040805162461bcd60e51b8152602060048201526024810191909152600080516020611f6283398151915260448201527f4b65793a2063616e6e6f74207265676973746572207a65726f207075626b657960648201526084016104a3565b6001600160a01b03851660009081526001602052604090205415610b225760405162461bcd60e51b81526020600482015260476024820152600080516020611f6283398151915260448201527f4b65793a206f70657261746f7220616c72656164792072656769737465726564606482015266207075626b657960c81b608482015260a4016104a3565b6000818152600260205260409020546001600160a01b031615610ba65760405162461bcd60e51b81526020600482015260426024820152600080516020611f6283398151915260448201527f4b65793a207075626c6963206b657920616c7265616479207265676973746572606482015261195960f21b608482015260a4016104a3565b604080516000917f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f000000191610bff918835916020808b0135928b01359160608c01359160808d019160c08e01918d35918e8201359101611d4b565b6040516020818303038152906040528051906020012060001c610c229190611d96565b9050610cbc610c5b610c4683610c40368a90038a0160408b01611d19565b9061131b565b610c5536899003890189611d19565b906113b2565b610c63611446565b610ca5610c9685610c40604080518082018252600080825260209182015281518083019092526001825260029082015290565b610c55368a90038a018a611d19565b610cb7368a90038a0160808b01611e08565b611506565b610d575760405162461bcd60e51b815260206004820152606c6024820152600080516020611f6283398151915260448201527f4b65793a2065697468657220746865204731207369676e61747572652069732060648201527f77726f6e672c206f7220473120616e642047322070726976617465206b65792060848201526b0c8de40dcdee840dac2e8c6d60a31b60a482015260c4016104a3565b6001600160a01b03861660008181526003602090815260408083208982018035825560608b013560019283015590835281842087905586845260029092529182902080546001600160a01b0319168417905590517fe3fb6613af2e8930cf85d47fcf6db10192224a64c6cbe8023e0eee1ba382804191610ddb9160808a0190611e65565b60405180910390a250949350505050565b606060008367ffffffffffffffff811115610e0957610e09611953565b604051908082528060200260200182016040528015610e32578160200160208202803683370190505b50905060005b8481101561102e576000868683818110610e5457610e54611d03565b919091013560f81c600081815260046020526040902054909250905063ffffffff81161580610ebd575060ff821660009081526004602052604081208054909190610ea157610ea1611d03565b600091825260209091200154600160c01b900463ffffffff1686105b15610f4a5760405162461bcd60e51b815260206004820152605160248201527f424c5341706b52656769737472792e67657441706b496e64696365734174426c60448201527f6f636b4e756d6265723a20626c6f636b4e756d626572206973206265666f7265606482015270207468652066697273742075706461746560781b608482015260a4016104a3565b60005b8163ffffffff168163ffffffff1610156110185760ff8316600090815260046020526040902087906001610f818486611eaf565b610f8b9190611eaf565b63ffffffff1681548110610fa157610fa1611d03565b600091825260209091200154600160c01b900463ffffffff1611611006576001610fcb8284611eaf565b610fd59190611eaf565b858581518110610fe757610fe7611d03565b602002602001019063ffffffff16908163ffffffff1681525050611018565b8061101081611ed4565b915050610f4d565b505050808061102690611ef8565b915050610e38565b50949350505050565b336001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161461107f5760405162461bcd60e51b81526004016104a390611c2a565b600061108a836108cb565b50905061109f8261109a83611773565b6110d0565b7f14a5172b312e9d2c22b8468f9c70ec2caa9de934fe380734fbc6f3beff2b14ba8383604051610645929190611c9e565b604080518082019091526000808252602082015260005b835181101561131557600084828151811061110457611104611d03565b0160209081015160f81c60008181526004909252604090912054909150806111945760405162461bcd60e51b815260206004820152603d60248201527f424c5341706b52656769737472792e5f70726f6365737351756f72756d41706b60448201527f5570646174653a2071756f72756d20646f6573206e6f7420657869737400000060648201526084016104a3565b60ff821660009081526005602090815260409182902082518084019093528054835260010154908201526111c890866113b2565b60ff831660008181526005602090815260408083208551808255868401805160019384015590855251835281842094845260049092528220939750919290916112119085611f13565b8154811061122157611221611d03565b600091825260209091200180549091504363ffffffff908116600160c01b9092041614156112625780546001600160c01b031916604083901c1781556112fe565b805463ffffffff438116600160e01b8181026001600160e01b0394851617855560ff88166000908152600460209081526040808320815160608101835267ffffffffffffffff198b16815280840196875280830185815282546001810184559286529390942093519301805495519251871690940291909516600160c01b026001600160e01b0319949094169190941c17919091179092161790555b50505050808061130d90611ef8565b9150506110e7565b50505050565b6040805180820190915260008082526020820152611337611832565b835181526020808501519082015260408082018490526000908360608460076107d05a03fa905080801561136a5761136c565bfe5b50806113aa5760405162461bcd60e51b815260206004820152600d60248201526c1958cb5b5d5b0b59985a5b1959609a1b60448201526064016104a3565b505092915050565b60408051808201909152600080825260208201526113ce611850565b835181526020808501518183015283516040808401919091529084015160608301526000908360808460066107d05a03fa905080801561136a5750806113aa5760405162461bcd60e51b815260206004820152600d60248201526c1958cb5859190b59985a5b1959609a1b60448201526064016104a3565b61144e61186e565b50604080516080810182527f198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c28183019081527f1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed6060830152815281518083019092527f275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec82527f1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d60208381019190915281019190915290565b604080518082018252858152602080820185905282518084019093528583528201839052600091611535611893565b60005b60028110156116fa57600061154e826006611f2a565b905084826002811061156257611562611d03565b60200201515183611574836000611f49565b600c811061158457611584611d03565b602002015284826002811061159b5761159b611d03565b602002015160200151838260016115b29190611f49565b600c81106115c2576115c2611d03565b60200201528382600281106115d9576115d9611d03565b60200201515151836115ec836002611f49565b600c81106115fc576115fc611d03565b602002015283826002811061161357611613611d03565b602002015151600160200201518361162c836003611f49565b600c811061163c5761163c611d03565b602002015283826002811061165357611653611d03565b60200201516020015160006002811061166e5761166e611d03565b60200201518361167f836004611f49565b600c811061168f5761168f611d03565b60200201528382600281106116a6576116a6611d03565b6020020151602001516001600281106116c1576116c1611d03565b6020020151836116d2836005611f49565b600c81106116e2576116e2611d03565b602002015250806116f281611ef8565b915050611538565b506117036118b2565b60006020826101808560086107d05a03fa905080801561136a5750806117635760405162461bcd60e51b81526020600482015260156024820152741c185a5c9a5b99cb5bdc18dbd9194b59985a5b1959605a1b60448201526064016104a3565b5051151598975050505050505050565b6040805180820190915260008082526020820152815115801561179857506020820151155b156117b6575050604080518082019091526000808252602082015290565b6040518060400160405280836000015181526020017f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4784602001516117fb9190611d96565b611825907f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47611f13565b905292915050565b919050565b60405180606001604052806003906020820280368337509192915050565b60405180608001604052806004906020820280368337509192915050565b60405180604001604052806118816118d0565b815260200161188e6118d0565b905290565b604051806101800160405280600c906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b60405180604001604052806002906020820280368337509192915050565b80356001600160a01b038116811461182d57600080fd5b60006020828403121561191757600080fd5b611920826118ee565b9392505050565b803560ff8116811461182d57600080fd5b60006020828403121561194a57600080fd5b61192082611927565b634e487b7160e01b600052604160045260246000fd5b6040805190810167ffffffffffffffff8111828210171561198c5761198c611953565b60405290565b604051601f8201601f1916810167ffffffffffffffff811182821017156119bb576119bb611953565b604052919050565b600080604083850312156119d657600080fd5b6119df836118ee565b915060208084013567ffffffffffffffff808211156119fd57600080fd5b818601915086601f830112611a1157600080fd5b813581811115611a2357611a23611953565b611a35601f8201601f19168501611992565b91508082528784828501011115611a4b57600080fd5b80848401858401376000848284010152508093505050509250929050565b600060208284031215611a7b57600080fd5b5035919050565b8151815260208083015190820152604081016106df565b60008060408385031215611aac57600080fd5b611ab583611927565b946020939093013593505050565b600080600060608486031215611ad857600080fd5b611ae184611927565b9250602084013563ffffffff81168114611afa57600080fd5b929592945050506040919091013590565b6000806000838503610160811215611b2257600080fd5b611b2b856118ee565b9350610100601f1982011215611b4057600080fd5b602085019250604061011f1982011215611b5957600080fd5b50610120840190509250925092565b600080600060408486031215611b7d57600080fd5b833567ffffffffffffffff80821115611b9557600080fd5b818601915086601f830112611ba957600080fd5b813581811115611bb857600080fd5b876020828501011115611bca57600080fd5b6020928301989097509590910135949350505050565b6020808252825182820181905260009190848201906040850190845b81811015611c1e57835163ffffffff1683529284019291840191600101611bfc565b50909695505050505050565b6020808252604e908201527f424c5341706b52656769737472792e6f6e6c795265676973747279436f6f726460408201527f696e61746f723a2063616c6c6572206973206e6f74207468652072656769737460608201526d393c9031b7b7b93234b730ba37b960911b608082015260a00190565b60018060a01b038316815260006020604081840152835180604085015260005b81811015611cda57858101830151858201606001528201611cbe565b81811115611cec576000606083870101525b50601f01601f191692909201606001949350505050565b634e487b7160e01b600052603260045260246000fd5b600060408284031215611d2b57600080fd5b611d33611969565b82358152602083013560208201528091505092915050565b8881528760208201528660408201528560608201526040856080830137600060c082016000815260408682375050610100810192909252610120820152610140019695505050505050565b600082611db357634e487b7160e01b600052601260045260246000fd5b500690565b600082601f830112611dc957600080fd5b611dd1611969565b806040840185811115611de357600080fd5b845b81811015611dfd578035845260209384019301611de5565b509095945050505050565b600060808284031215611e1a57600080fd5b6040516040810181811067ffffffffffffffff82111715611e3d57611e3d611953565b604052611e4a8484611db8565b8152611e598460408501611db8565b60208201529392505050565b823581526020808401359082015260c081016040838184013760808201600081526040808501823750600081529392505050565b634e487b7160e01b600052601160045260246000fd5b600063ffffffff83811690831681811015611ecc57611ecc611e99565b039392505050565b600063ffffffff80831681811415611eee57611eee611e99565b6001019392505050565b6000600019821415611f0c57611f0c611e99565b5060010190565b600082821015611f2557611f25611e99565b500390565b6000816000190483118215151615611f4457611f44611e99565b500290565b60008219821115611f5c57611f5c611e99565b50019056fe424c5341706b52656769737472792e7265676973746572424c535075626c6963a2646970667358221220d60912434f8da97fa12b719aaca1890878f8f5822bd4c255c1dee557b8e627ce64736f6c634300080c0033",
}

// ContractBLSApkRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractBLSApkRegistryMetaData.ABI instead.
var ContractBLSApkRegistryABI = ContractBLSApkRegistryMetaData.ABI

// ContractBLSApkRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractBLSApkRegistryMetaData.Bin instead.
var ContractBLSApkRegistryBin = ContractBLSApkRegistryMetaData.Bin

// DeployContractBLSApkRegistry deploys a new Ethereum contract, binding an instance of ContractBLSApkRegistry to it.
func DeployContractBLSApkRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, _registryCoordinator common.Address) (common.Address, *types.Transaction, *ContractBLSApkRegistry, error) {
	parsed, err := ContractBLSApkRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractBLSApkRegistryBin), backend, _registryCoordinator)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractBLSApkRegistry{ContractBLSApkRegistryCaller: ContractBLSApkRegistryCaller{contract: contract}, ContractBLSApkRegistryTransactor: ContractBLSApkRegistryTransactor{contract: contract}, ContractBLSApkRegistryFilterer: ContractBLSApkRegistryFilterer{contract: contract}}, nil
}

// ContractBLSApkRegistry is an auto generated Go binding around an Ethereum contract.
type ContractBLSApkRegistry struct {
	ContractBLSApkRegistryCaller     // Read-only binding to the contract
	ContractBLSApkRegistryTransactor // Write-only binding to the contract
	ContractBLSApkRegistryFilterer   // Log filterer for contract events
}

// ContractBLSApkRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractBLSApkRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractBLSApkRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractBLSApkRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractBLSApkRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractBLSApkRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractBLSApkRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractBLSApkRegistrySession struct {
	Contract     *ContractBLSApkRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// ContractBLSApkRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractBLSApkRegistryCallerSession struct {
	Contract *ContractBLSApkRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// ContractBLSApkRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractBLSApkRegistryTransactorSession struct {
	Contract     *ContractBLSApkRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// ContractBLSApkRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractBLSApkRegistryRaw struct {
	Contract *ContractBLSApkRegistry // Generic contract binding to access the raw methods on
}

// ContractBLSApkRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractBLSApkRegistryCallerRaw struct {
	Contract *ContractBLSApkRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// ContractBLSApkRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractBLSApkRegistryTransactorRaw struct {
	Contract *ContractBLSApkRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractBLSApkRegistry creates a new instance of ContractBLSApkRegistry, bound to a specific deployed contract.
func NewContractBLSApkRegistry(address common.Address, backend bind.ContractBackend) (*ContractBLSApkRegistry, error) {
	contract, err := bindContractBLSApkRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractBLSApkRegistry{ContractBLSApkRegistryCaller: ContractBLSApkRegistryCaller{contract: contract}, ContractBLSApkRegistryTransactor: ContractBLSApkRegistryTransactor{contract: contract}, ContractBLSApkRegistryFilterer: ContractBLSApkRegistryFilterer{contract: contract}}, nil
}

// NewContractBLSApkRegistryCaller creates a new read-only instance of ContractBLSApkRegistry, bound to a specific deployed contract.
func NewContractBLSApkRegistryCaller(address common.Address, caller bind.ContractCaller) (*ContractBLSApkRegistryCaller, error) {
	contract, err := bindContractBLSApkRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractBLSApkRegistryCaller{contract: contract}, nil
}

// NewContractBLSApkRegistryTransactor creates a new write-only instance of ContractBLSApkRegistry, bound to a specific deployed contract.
func NewContractBLSApkRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractBLSApkRegistryTransactor, error) {
	contract, err := bindContractBLSApkRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractBLSApkRegistryTransactor{contract: contract}, nil
}

// NewContractBLSApkRegistryFilterer creates a new log filterer instance of ContractBLSApkRegistry, bound to a specific deployed contract.
func NewContractBLSApkRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractBLSApkRegistryFilterer, error) {
	contract, err := bindContractBLSApkRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractBLSApkRegistryFilterer{contract: contract}, nil
}

// bindContractBLSApkRegistry binds a generic wrapper to an already deployed contract.
func bindContractBLSApkRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractBLSApkRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractBLSApkRegistry *ContractBLSApkRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractBLSApkRegistry.Contract.ContractBLSApkRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractBLSApkRegistry *ContractBLSApkRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractBLSApkRegistry.Contract.ContractBLSApkRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractBLSApkRegistry *ContractBLSApkRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractBLSApkRegistry.Contract.ContractBLSApkRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractBLSApkRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractBLSApkRegistry *ContractBLSApkRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractBLSApkRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractBLSApkRegistry *ContractBLSApkRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractBLSApkRegistry.Contract.contract.Transact(opts, method, params...)
}

// ApkHistory is a free data retrieval call binding the contract method 0x7916cea6.
//
// Solidity: function apkHistory(uint8 , uint256 ) view returns(bytes24 apkHash, uint32 updateBlockNumber, uint32 nextUpdateBlockNumber)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCaller) ApkHistory(opts *bind.CallOpts, arg0 uint8, arg1 *big.Int) (struct {
	ApkHash               [24]byte
	UpdateBlockNumber     uint32
	NextUpdateBlockNumber uint32
}, error) {
	var out []interface{}
	err := _ContractBLSApkRegistry.contract.Call(opts, &out, "apkHistory", arg0, arg1)

	outstruct := new(struct {
		ApkHash               [24]byte
		UpdateBlockNumber     uint32
		NextUpdateBlockNumber uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ApkHash = *abi.ConvertType(out[0], new([24]byte)).(*[24]byte)
	outstruct.UpdateBlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.NextUpdateBlockNumber = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

// ApkHistory is a free data retrieval call binding the contract method 0x7916cea6.
//
// Solidity: function apkHistory(uint8 , uint256 ) view returns(bytes24 apkHash, uint32 updateBlockNumber, uint32 nextUpdateBlockNumber)
func (_ContractBLSApkRegistry *ContractBLSApkRegistrySession) ApkHistory(arg0 uint8, arg1 *big.Int) (struct {
	ApkHash               [24]byte
	UpdateBlockNumber     uint32
	NextUpdateBlockNumber uint32
}, error) {
	return _ContractBLSApkRegistry.Contract.ApkHistory(&_ContractBLSApkRegistry.CallOpts, arg0, arg1)
}

// ApkHistory is a free data retrieval call binding the contract method 0x7916cea6.
//
// Solidity: function apkHistory(uint8 , uint256 ) view returns(bytes24 apkHash, uint32 updateBlockNumber, uint32 nextUpdateBlockNumber)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCallerSession) ApkHistory(arg0 uint8, arg1 *big.Int) (struct {
	ApkHash               [24]byte
	UpdateBlockNumber     uint32
	NextUpdateBlockNumber uint32
}, error) {
	return _ContractBLSApkRegistry.Contract.ApkHistory(&_ContractBLSApkRegistry.CallOpts, arg0, arg1)
}

// CurrentApk is a free data retrieval call binding the contract method 0xa3db80e2.
//
// Solidity: function currentApk(uint8 ) view returns(uint256 X, uint256 Y)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCaller) CurrentApk(opts *bind.CallOpts, arg0 uint8) (struct {
	X *big.Int
	Y *big.Int
}, error) {
	var out []interface{}
	err := _ContractBLSApkRegistry.contract.Call(opts, &out, "currentApk", arg0)

	outstruct := new(struct {
		X *big.Int
		Y *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.X = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Y = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// CurrentApk is a free data retrieval call binding the contract method 0xa3db80e2.
//
// Solidity: function currentApk(uint8 ) view returns(uint256 X, uint256 Y)
func (_ContractBLSApkRegistry *ContractBLSApkRegistrySession) CurrentApk(arg0 uint8) (struct {
	X *big.Int
	Y *big.Int
}, error) {
	return _ContractBLSApkRegistry.Contract.CurrentApk(&_ContractBLSApkRegistry.CallOpts, arg0)
}

// CurrentApk is a free data retrieval call binding the contract method 0xa3db80e2.
//
// Solidity: function currentApk(uint8 ) view returns(uint256 X, uint256 Y)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCallerSession) CurrentApk(arg0 uint8) (struct {
	X *big.Int
	Y *big.Int
}, error) {
	return _ContractBLSApkRegistry.Contract.CurrentApk(&_ContractBLSApkRegistry.CallOpts, arg0)
}

// GetApk is a free data retrieval call binding the contract method 0x5f61a884.
//
// Solidity: function getApk(uint8 quorumNumber) view returns((uint256,uint256))
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCaller) GetApk(opts *bind.CallOpts, quorumNumber uint8) (BN254G1Point, error) {
	var out []interface{}
	err := _ContractBLSApkRegistry.contract.Call(opts, &out, "getApk", quorumNumber)

	if err != nil {
		return *new(BN254G1Point), err
	}

	out0 := *abi.ConvertType(out[0], new(BN254G1Point)).(*BN254G1Point)

	return out0, err

}

// GetApk is a free data retrieval call binding the contract method 0x5f61a884.
//
// Solidity: function getApk(uint8 quorumNumber) view returns((uint256,uint256))
func (_ContractBLSApkRegistry *ContractBLSApkRegistrySession) GetApk(quorumNumber uint8) (BN254G1Point, error) {
	return _ContractBLSApkRegistry.Contract.GetApk(&_ContractBLSApkRegistry.CallOpts, quorumNumber)
}

// GetApk is a free data retrieval call binding the contract method 0x5f61a884.
//
// Solidity: function getApk(uint8 quorumNumber) view returns((uint256,uint256))
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCallerSession) GetApk(quorumNumber uint8) (BN254G1Point, error) {
	return _ContractBLSApkRegistry.Contract.GetApk(&_ContractBLSApkRegistry.CallOpts, quorumNumber)
}

// GetApkHashAtBlockNumberAndIndex is a free data retrieval call binding the contract method 0x68bccaac.
//
// Solidity: function getApkHashAtBlockNumberAndIndex(uint8 quorumNumber, uint32 blockNumber, uint256 index) view returns(bytes24)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCaller) GetApkHashAtBlockNumberAndIndex(opts *bind.CallOpts, quorumNumber uint8, blockNumber uint32, index *big.Int) ([24]byte, error) {
	var out []interface{}
	err := _ContractBLSApkRegistry.contract.Call(opts, &out, "getApkHashAtBlockNumberAndIndex", quorumNumber, blockNumber, index)

	if err != nil {
		return *new([24]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([24]byte)).(*[24]byte)

	return out0, err

}

// GetApkHashAtBlockNumberAndIndex is a free data retrieval call binding the contract method 0x68bccaac.
//
// Solidity: function getApkHashAtBlockNumberAndIndex(uint8 quorumNumber, uint32 blockNumber, uint256 index) view returns(bytes24)
func (_ContractBLSApkRegistry *ContractBLSApkRegistrySession) GetApkHashAtBlockNumberAndIndex(quorumNumber uint8, blockNumber uint32, index *big.Int) ([24]byte, error) {
	return _ContractBLSApkRegistry.Contract.GetApkHashAtBlockNumberAndIndex(&_ContractBLSApkRegistry.CallOpts, quorumNumber, blockNumber, index)
}

// GetApkHashAtBlockNumberAndIndex is a free data retrieval call binding the contract method 0x68bccaac.
//
// Solidity: function getApkHashAtBlockNumberAndIndex(uint8 quorumNumber, uint32 blockNumber, uint256 index) view returns(bytes24)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCallerSession) GetApkHashAtBlockNumberAndIndex(quorumNumber uint8, blockNumber uint32, index *big.Int) ([24]byte, error) {
	return _ContractBLSApkRegistry.Contract.GetApkHashAtBlockNumberAndIndex(&_ContractBLSApkRegistry.CallOpts, quorumNumber, blockNumber, index)
}

// GetApkHistoryLength is a free data retrieval call binding the contract method 0x377ed99d.
//
// Solidity: function getApkHistoryLength(uint8 quorumNumber) view returns(uint32)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCaller) GetApkHistoryLength(opts *bind.CallOpts, quorumNumber uint8) (uint32, error) {
	var out []interface{}
	err := _ContractBLSApkRegistry.contract.Call(opts, &out, "getApkHistoryLength", quorumNumber)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// GetApkHistoryLength is a free data retrieval call binding the contract method 0x377ed99d.
//
// Solidity: function getApkHistoryLength(uint8 quorumNumber) view returns(uint32)
func (_ContractBLSApkRegistry *ContractBLSApkRegistrySession) GetApkHistoryLength(quorumNumber uint8) (uint32, error) {
	return _ContractBLSApkRegistry.Contract.GetApkHistoryLength(&_ContractBLSApkRegistry.CallOpts, quorumNumber)
}

// GetApkHistoryLength is a free data retrieval call binding the contract method 0x377ed99d.
//
// Solidity: function getApkHistoryLength(uint8 quorumNumber) view returns(uint32)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCallerSession) GetApkHistoryLength(quorumNumber uint8) (uint32, error) {
	return _ContractBLSApkRegistry.Contract.GetApkHistoryLength(&_ContractBLSApkRegistry.CallOpts, quorumNumber)
}

// GetApkIndicesAtBlockNumber is a free data retrieval call binding the contract method 0xd5254a8c.
//
// Solidity: function getApkIndicesAtBlockNumber(bytes quorumNumbers, uint256 blockNumber) view returns(uint32[])
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCaller) GetApkIndicesAtBlockNumber(opts *bind.CallOpts, quorumNumbers []byte, blockNumber *big.Int) ([]uint32, error) {
	var out []interface{}
	err := _ContractBLSApkRegistry.contract.Call(opts, &out, "getApkIndicesAtBlockNumber", quorumNumbers, blockNumber)

	if err != nil {
		return *new([]uint32), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint32)).(*[]uint32)

	return out0, err

}

// GetApkIndicesAtBlockNumber is a free data retrieval call binding the contract method 0xd5254a8c.
//
// Solidity: function getApkIndicesAtBlockNumber(bytes quorumNumbers, uint256 blockNumber) view returns(uint32[])
func (_ContractBLSApkRegistry *ContractBLSApkRegistrySession) GetApkIndicesAtBlockNumber(quorumNumbers []byte, blockNumber *big.Int) ([]uint32, error) {
	return _ContractBLSApkRegistry.Contract.GetApkIndicesAtBlockNumber(&_ContractBLSApkRegistry.CallOpts, quorumNumbers, blockNumber)
}

// GetApkIndicesAtBlockNumber is a free data retrieval call binding the contract method 0xd5254a8c.
//
// Solidity: function getApkIndicesAtBlockNumber(bytes quorumNumbers, uint256 blockNumber) view returns(uint32[])
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCallerSession) GetApkIndicesAtBlockNumber(quorumNumbers []byte, blockNumber *big.Int) ([]uint32, error) {
	return _ContractBLSApkRegistry.Contract.GetApkIndicesAtBlockNumber(&_ContractBLSApkRegistry.CallOpts, quorumNumbers, blockNumber)
}

// GetApkUpdateAtIndex is a free data retrieval call binding the contract method 0x605747d5.
//
// Solidity: function getApkUpdateAtIndex(uint8 quorumNumber, uint256 index) view returns((bytes24,uint32,uint32))
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCaller) GetApkUpdateAtIndex(opts *bind.CallOpts, quorumNumber uint8, index *big.Int) (IBLSApkRegistryApkUpdate, error) {
	var out []interface{}
	err := _ContractBLSApkRegistry.contract.Call(opts, &out, "getApkUpdateAtIndex", quorumNumber, index)

	if err != nil {
		return *new(IBLSApkRegistryApkUpdate), err
	}

	out0 := *abi.ConvertType(out[0], new(IBLSApkRegistryApkUpdate)).(*IBLSApkRegistryApkUpdate)

	return out0, err

}

// GetApkUpdateAtIndex is a free data retrieval call binding the contract method 0x605747d5.
//
// Solidity: function getApkUpdateAtIndex(uint8 quorumNumber, uint256 index) view returns((bytes24,uint32,uint32))
func (_ContractBLSApkRegistry *ContractBLSApkRegistrySession) GetApkUpdateAtIndex(quorumNumber uint8, index *big.Int) (IBLSApkRegistryApkUpdate, error) {
	return _ContractBLSApkRegistry.Contract.GetApkUpdateAtIndex(&_ContractBLSApkRegistry.CallOpts, quorumNumber, index)
}

// GetApkUpdateAtIndex is a free data retrieval call binding the contract method 0x605747d5.
//
// Solidity: function getApkUpdateAtIndex(uint8 quorumNumber, uint256 index) view returns((bytes24,uint32,uint32))
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCallerSession) GetApkUpdateAtIndex(quorumNumber uint8, index *big.Int) (IBLSApkRegistryApkUpdate, error) {
	return _ContractBLSApkRegistry.Contract.GetApkUpdateAtIndex(&_ContractBLSApkRegistry.CallOpts, quorumNumber, index)
}

// GetOperatorFromPubkeyHash is a free data retrieval call binding the contract method 0x47b314e8.
//
// Solidity: function getOperatorFromPubkeyHash(bytes32 pubkeyHash) view returns(address)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCaller) GetOperatorFromPubkeyHash(opts *bind.CallOpts, pubkeyHash [32]byte) (common.Address, error) {
	var out []interface{}
	err := _ContractBLSApkRegistry.contract.Call(opts, &out, "getOperatorFromPubkeyHash", pubkeyHash)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOperatorFromPubkeyHash is a free data retrieval call binding the contract method 0x47b314e8.
//
// Solidity: function getOperatorFromPubkeyHash(bytes32 pubkeyHash) view returns(address)
func (_ContractBLSApkRegistry *ContractBLSApkRegistrySession) GetOperatorFromPubkeyHash(pubkeyHash [32]byte) (common.Address, error) {
	return _ContractBLSApkRegistry.Contract.GetOperatorFromPubkeyHash(&_ContractBLSApkRegistry.CallOpts, pubkeyHash)
}

// GetOperatorFromPubkeyHash is a free data retrieval call binding the contract method 0x47b314e8.
//
// Solidity: function getOperatorFromPubkeyHash(bytes32 pubkeyHash) view returns(address)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCallerSession) GetOperatorFromPubkeyHash(pubkeyHash [32]byte) (common.Address, error) {
	return _ContractBLSApkRegistry.Contract.GetOperatorFromPubkeyHash(&_ContractBLSApkRegistry.CallOpts, pubkeyHash)
}

// GetOperatorId is a free data retrieval call binding the contract method 0x13542a4e.
//
// Solidity: function getOperatorId(address operator) view returns(bytes32)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCaller) GetOperatorId(opts *bind.CallOpts, operator common.Address) ([32]byte, error) {
	var out []interface{}
	err := _ContractBLSApkRegistry.contract.Call(opts, &out, "getOperatorId", operator)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetOperatorId is a free data retrieval call binding the contract method 0x13542a4e.
//
// Solidity: function getOperatorId(address operator) view returns(bytes32)
func (_ContractBLSApkRegistry *ContractBLSApkRegistrySession) GetOperatorId(operator common.Address) ([32]byte, error) {
	return _ContractBLSApkRegistry.Contract.GetOperatorId(&_ContractBLSApkRegistry.CallOpts, operator)
}

// GetOperatorId is a free data retrieval call binding the contract method 0x13542a4e.
//
// Solidity: function getOperatorId(address operator) view returns(bytes32)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCallerSession) GetOperatorId(operator common.Address) ([32]byte, error) {
	return _ContractBLSApkRegistry.Contract.GetOperatorId(&_ContractBLSApkRegistry.CallOpts, operator)
}

// GetRegisteredPubkey is a free data retrieval call binding the contract method 0x7ff81a87.
//
// Solidity: function getRegisteredPubkey(address operator) view returns((uint256,uint256), bytes32)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCaller) GetRegisteredPubkey(opts *bind.CallOpts, operator common.Address) (BN254G1Point, [32]byte, error) {
	var out []interface{}
	err := _ContractBLSApkRegistry.contract.Call(opts, &out, "getRegisteredPubkey", operator)

	if err != nil {
		return *new(BN254G1Point), *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(BN254G1Point)).(*BN254G1Point)
	out1 := *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)

	return out0, out1, err

}

// GetRegisteredPubkey is a free data retrieval call binding the contract method 0x7ff81a87.
//
// Solidity: function getRegisteredPubkey(address operator) view returns((uint256,uint256), bytes32)
func (_ContractBLSApkRegistry *ContractBLSApkRegistrySession) GetRegisteredPubkey(operator common.Address) (BN254G1Point, [32]byte, error) {
	return _ContractBLSApkRegistry.Contract.GetRegisteredPubkey(&_ContractBLSApkRegistry.CallOpts, operator)
}

// GetRegisteredPubkey is a free data retrieval call binding the contract method 0x7ff81a87.
//
// Solidity: function getRegisteredPubkey(address operator) view returns((uint256,uint256), bytes32)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCallerSession) GetRegisteredPubkey(operator common.Address) (BN254G1Point, [32]byte, error) {
	return _ContractBLSApkRegistry.Contract.GetRegisteredPubkey(&_ContractBLSApkRegistry.CallOpts, operator)
}

// OperatorToPubkey is a free data retrieval call binding the contract method 0x00a1f4cb.
//
// Solidity: function operatorToPubkey(address ) view returns(uint256 X, uint256 Y)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCaller) OperatorToPubkey(opts *bind.CallOpts, arg0 common.Address) (struct {
	X *big.Int
	Y *big.Int
}, error) {
	var out []interface{}
	err := _ContractBLSApkRegistry.contract.Call(opts, &out, "operatorToPubkey", arg0)

	outstruct := new(struct {
		X *big.Int
		Y *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.X = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Y = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// OperatorToPubkey is a free data retrieval call binding the contract method 0x00a1f4cb.
//
// Solidity: function operatorToPubkey(address ) view returns(uint256 X, uint256 Y)
func (_ContractBLSApkRegistry *ContractBLSApkRegistrySession) OperatorToPubkey(arg0 common.Address) (struct {
	X *big.Int
	Y *big.Int
}, error) {
	return _ContractBLSApkRegistry.Contract.OperatorToPubkey(&_ContractBLSApkRegistry.CallOpts, arg0)
}

// OperatorToPubkey is a free data retrieval call binding the contract method 0x00a1f4cb.
//
// Solidity: function operatorToPubkey(address ) view returns(uint256 X, uint256 Y)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCallerSession) OperatorToPubkey(arg0 common.Address) (struct {
	X *big.Int
	Y *big.Int
}, error) {
	return _ContractBLSApkRegistry.Contract.OperatorToPubkey(&_ContractBLSApkRegistry.CallOpts, arg0)
}

// OperatorToPubkeyHash is a free data retrieval call binding the contract method 0xde29fac0.
//
// Solidity: function operatorToPubkeyHash(address ) view returns(bytes32)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCaller) OperatorToPubkeyHash(opts *bind.CallOpts, arg0 common.Address) ([32]byte, error) {
	var out []interface{}
	err := _ContractBLSApkRegistry.contract.Call(opts, &out, "operatorToPubkeyHash", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// OperatorToPubkeyHash is a free data retrieval call binding the contract method 0xde29fac0.
//
// Solidity: function operatorToPubkeyHash(address ) view returns(bytes32)
func (_ContractBLSApkRegistry *ContractBLSApkRegistrySession) OperatorToPubkeyHash(arg0 common.Address) ([32]byte, error) {
	return _ContractBLSApkRegistry.Contract.OperatorToPubkeyHash(&_ContractBLSApkRegistry.CallOpts, arg0)
}

// OperatorToPubkeyHash is a free data retrieval call binding the contract method 0xde29fac0.
//
// Solidity: function operatorToPubkeyHash(address ) view returns(bytes32)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCallerSession) OperatorToPubkeyHash(arg0 common.Address) ([32]byte, error) {
	return _ContractBLSApkRegistry.Contract.OperatorToPubkeyHash(&_ContractBLSApkRegistry.CallOpts, arg0)
}

// PubkeyHashToOperator is a free data retrieval call binding the contract method 0xe8bb9ae6.
//
// Solidity: function pubkeyHashToOperator(bytes32 ) view returns(address)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCaller) PubkeyHashToOperator(opts *bind.CallOpts, arg0 [32]byte) (common.Address, error) {
	var out []interface{}
	err := _ContractBLSApkRegistry.contract.Call(opts, &out, "pubkeyHashToOperator", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PubkeyHashToOperator is a free data retrieval call binding the contract method 0xe8bb9ae6.
//
// Solidity: function pubkeyHashToOperator(bytes32 ) view returns(address)
func (_ContractBLSApkRegistry *ContractBLSApkRegistrySession) PubkeyHashToOperator(arg0 [32]byte) (common.Address, error) {
	return _ContractBLSApkRegistry.Contract.PubkeyHashToOperator(&_ContractBLSApkRegistry.CallOpts, arg0)
}

// PubkeyHashToOperator is a free data retrieval call binding the contract method 0xe8bb9ae6.
//
// Solidity: function pubkeyHashToOperator(bytes32 ) view returns(address)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCallerSession) PubkeyHashToOperator(arg0 [32]byte) (common.Address, error) {
	return _ContractBLSApkRegistry.Contract.PubkeyHashToOperator(&_ContractBLSApkRegistry.CallOpts, arg0)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCaller) RegistryCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractBLSApkRegistry.contract.Call(opts, &out, "registryCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractBLSApkRegistry *ContractBLSApkRegistrySession) RegistryCoordinator() (common.Address, error) {
	return _ContractBLSApkRegistry.Contract.RegistryCoordinator(&_ContractBLSApkRegistry.CallOpts)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryCallerSession) RegistryCoordinator() (common.Address, error) {
	return _ContractBLSApkRegistry.Contract.RegistryCoordinator(&_ContractBLSApkRegistry.CallOpts)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xf4e24fe5.
//
// Solidity: function deregisterOperator(address operator, bytes quorumNumbers) returns()
func (_ContractBLSApkRegistry *ContractBLSApkRegistryTransactor) DeregisterOperator(opts *bind.TransactOpts, operator common.Address, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractBLSApkRegistry.contract.Transact(opts, "deregisterOperator", operator, quorumNumbers)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xf4e24fe5.
//
// Solidity: function deregisterOperator(address operator, bytes quorumNumbers) returns()
func (_ContractBLSApkRegistry *ContractBLSApkRegistrySession) DeregisterOperator(operator common.Address, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractBLSApkRegistry.Contract.DeregisterOperator(&_ContractBLSApkRegistry.TransactOpts, operator, quorumNumbers)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xf4e24fe5.
//
// Solidity: function deregisterOperator(address operator, bytes quorumNumbers) returns()
func (_ContractBLSApkRegistry *ContractBLSApkRegistryTransactorSession) DeregisterOperator(operator common.Address, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractBLSApkRegistry.Contract.DeregisterOperator(&_ContractBLSApkRegistry.TransactOpts, operator, quorumNumbers)
}

// InitializeQuorum is a paid mutator transaction binding the contract method 0x26d941f2.
//
// Solidity: function initializeQuorum(uint8 quorumNumber) returns()
func (_ContractBLSApkRegistry *ContractBLSApkRegistryTransactor) InitializeQuorum(opts *bind.TransactOpts, quorumNumber uint8) (*types.Transaction, error) {
	return _ContractBLSApkRegistry.contract.Transact(opts, "initializeQuorum", quorumNumber)
}

// InitializeQuorum is a paid mutator transaction binding the contract method 0x26d941f2.
//
// Solidity: function initializeQuorum(uint8 quorumNumber) returns()
func (_ContractBLSApkRegistry *ContractBLSApkRegistrySession) InitializeQuorum(quorumNumber uint8) (*types.Transaction, error) {
	return _ContractBLSApkRegistry.Contract.InitializeQuorum(&_ContractBLSApkRegistry.TransactOpts, quorumNumber)
}

// InitializeQuorum is a paid mutator transaction binding the contract method 0x26d941f2.
//
// Solidity: function initializeQuorum(uint8 quorumNumber) returns()
func (_ContractBLSApkRegistry *ContractBLSApkRegistryTransactorSession) InitializeQuorum(quorumNumber uint8) (*types.Transaction, error) {
	return _ContractBLSApkRegistry.Contract.InitializeQuorum(&_ContractBLSApkRegistry.TransactOpts, quorumNumber)
}

// RegisterBLSPublicKey is a paid mutator transaction binding the contract method 0xbf79ce58.
//
// Solidity: function registerBLSPublicKey(address operator, ((uint256,uint256),(uint256,uint256),(uint256[2],uint256[2])) params, (uint256,uint256) pubkeyRegistrationMessageHash) returns(bytes32 operatorId)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryTransactor) RegisterBLSPublicKey(opts *bind.TransactOpts, operator common.Address, params IBLSApkRegistryPubkeyRegistrationParams, pubkeyRegistrationMessageHash BN254G1Point) (*types.Transaction, error) {
	return _ContractBLSApkRegistry.contract.Transact(opts, "registerBLSPublicKey", operator, params, pubkeyRegistrationMessageHash)
}

// RegisterBLSPublicKey is a paid mutator transaction binding the contract method 0xbf79ce58.
//
// Solidity: function registerBLSPublicKey(address operator, ((uint256,uint256),(uint256,uint256),(uint256[2],uint256[2])) params, (uint256,uint256) pubkeyRegistrationMessageHash) returns(bytes32 operatorId)
func (_ContractBLSApkRegistry *ContractBLSApkRegistrySession) RegisterBLSPublicKey(operator common.Address, params IBLSApkRegistryPubkeyRegistrationParams, pubkeyRegistrationMessageHash BN254G1Point) (*types.Transaction, error) {
	return _ContractBLSApkRegistry.Contract.RegisterBLSPublicKey(&_ContractBLSApkRegistry.TransactOpts, operator, params, pubkeyRegistrationMessageHash)
}

// RegisterBLSPublicKey is a paid mutator transaction binding the contract method 0xbf79ce58.
//
// Solidity: function registerBLSPublicKey(address operator, ((uint256,uint256),(uint256,uint256),(uint256[2],uint256[2])) params, (uint256,uint256) pubkeyRegistrationMessageHash) returns(bytes32 operatorId)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryTransactorSession) RegisterBLSPublicKey(operator common.Address, params IBLSApkRegistryPubkeyRegistrationParams, pubkeyRegistrationMessageHash BN254G1Point) (*types.Transaction, error) {
	return _ContractBLSApkRegistry.Contract.RegisterBLSPublicKey(&_ContractBLSApkRegistry.TransactOpts, operator, params, pubkeyRegistrationMessageHash)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x3fb27952.
//
// Solidity: function registerOperator(address operator, bytes quorumNumbers) returns()
func (_ContractBLSApkRegistry *ContractBLSApkRegistryTransactor) RegisterOperator(opts *bind.TransactOpts, operator common.Address, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractBLSApkRegistry.contract.Transact(opts, "registerOperator", operator, quorumNumbers)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x3fb27952.
//
// Solidity: function registerOperator(address operator, bytes quorumNumbers) returns()
func (_ContractBLSApkRegistry *ContractBLSApkRegistrySession) RegisterOperator(operator common.Address, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractBLSApkRegistry.Contract.RegisterOperator(&_ContractBLSApkRegistry.TransactOpts, operator, quorumNumbers)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x3fb27952.
//
// Solidity: function registerOperator(address operator, bytes quorumNumbers) returns()
func (_ContractBLSApkRegistry *ContractBLSApkRegistryTransactorSession) RegisterOperator(operator common.Address, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractBLSApkRegistry.Contract.RegisterOperator(&_ContractBLSApkRegistry.TransactOpts, operator, quorumNumbers)
}

// ContractBLSApkRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ContractBLSApkRegistry contract.
type ContractBLSApkRegistryInitializedIterator struct {
	Event *ContractBLSApkRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *ContractBLSApkRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractBLSApkRegistryInitialized)
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
		it.Event = new(ContractBLSApkRegistryInitialized)
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
func (it *ContractBLSApkRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractBLSApkRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractBLSApkRegistryInitialized represents a Initialized event raised by the ContractBLSApkRegistry contract.
type ContractBLSApkRegistryInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*ContractBLSApkRegistryInitializedIterator, error) {

	logs, sub, err := _ContractBLSApkRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ContractBLSApkRegistryInitializedIterator{contract: _ContractBLSApkRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ContractBLSApkRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _ContractBLSApkRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractBLSApkRegistryInitialized)
				if err := _ContractBLSApkRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ContractBLSApkRegistry *ContractBLSApkRegistryFilterer) ParseInitialized(log types.Log) (*ContractBLSApkRegistryInitialized, error) {
	event := new(ContractBLSApkRegistryInitialized)
	if err := _ContractBLSApkRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractBLSApkRegistryNewPubkeyRegistrationIterator is returned from FilterNewPubkeyRegistration and is used to iterate over the raw logs and unpacked data for NewPubkeyRegistration events raised by the ContractBLSApkRegistry contract.
type ContractBLSApkRegistryNewPubkeyRegistrationIterator struct {
	Event *ContractBLSApkRegistryNewPubkeyRegistration // Event containing the contract specifics and raw log

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
func (it *ContractBLSApkRegistryNewPubkeyRegistrationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractBLSApkRegistryNewPubkeyRegistration)
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
		it.Event = new(ContractBLSApkRegistryNewPubkeyRegistration)
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
func (it *ContractBLSApkRegistryNewPubkeyRegistrationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractBLSApkRegistryNewPubkeyRegistrationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractBLSApkRegistryNewPubkeyRegistration represents a NewPubkeyRegistration event raised by the ContractBLSApkRegistry contract.
type ContractBLSApkRegistryNewPubkeyRegistration struct {
	Operator common.Address
	PubkeyG1 BN254G1Point
	PubkeyG2 BN254G2Point
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterNewPubkeyRegistration is a free log retrieval operation binding the contract event 0xe3fb6613af2e8930cf85d47fcf6db10192224a64c6cbe8023e0eee1ba3828041.
//
// Solidity: event NewPubkeyRegistration(address indexed operator, (uint256,uint256) pubkeyG1, (uint256[2],uint256[2]) pubkeyG2)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryFilterer) FilterNewPubkeyRegistration(opts *bind.FilterOpts, operator []common.Address) (*ContractBLSApkRegistryNewPubkeyRegistrationIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ContractBLSApkRegistry.contract.FilterLogs(opts, "NewPubkeyRegistration", operatorRule)
	if err != nil {
		return nil, err
	}
	return &ContractBLSApkRegistryNewPubkeyRegistrationIterator{contract: _ContractBLSApkRegistry.contract, event: "NewPubkeyRegistration", logs: logs, sub: sub}, nil
}

// WatchNewPubkeyRegistration is a free log subscription operation binding the contract event 0xe3fb6613af2e8930cf85d47fcf6db10192224a64c6cbe8023e0eee1ba3828041.
//
// Solidity: event NewPubkeyRegistration(address indexed operator, (uint256,uint256) pubkeyG1, (uint256[2],uint256[2]) pubkeyG2)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryFilterer) WatchNewPubkeyRegistration(opts *bind.WatchOpts, sink chan<- *ContractBLSApkRegistryNewPubkeyRegistration, operator []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ContractBLSApkRegistry.contract.WatchLogs(opts, "NewPubkeyRegistration", operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractBLSApkRegistryNewPubkeyRegistration)
				if err := _ContractBLSApkRegistry.contract.UnpackLog(event, "NewPubkeyRegistration", log); err != nil {
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

// ParseNewPubkeyRegistration is a log parse operation binding the contract event 0xe3fb6613af2e8930cf85d47fcf6db10192224a64c6cbe8023e0eee1ba3828041.
//
// Solidity: event NewPubkeyRegistration(address indexed operator, (uint256,uint256) pubkeyG1, (uint256[2],uint256[2]) pubkeyG2)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryFilterer) ParseNewPubkeyRegistration(log types.Log) (*ContractBLSApkRegistryNewPubkeyRegistration, error) {
	event := new(ContractBLSApkRegistryNewPubkeyRegistration)
	if err := _ContractBLSApkRegistry.contract.UnpackLog(event, "NewPubkeyRegistration", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractBLSApkRegistryOperatorAddedToQuorumsIterator is returned from FilterOperatorAddedToQuorums and is used to iterate over the raw logs and unpacked data for OperatorAddedToQuorums events raised by the ContractBLSApkRegistry contract.
type ContractBLSApkRegistryOperatorAddedToQuorumsIterator struct {
	Event *ContractBLSApkRegistryOperatorAddedToQuorums // Event containing the contract specifics and raw log

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
func (it *ContractBLSApkRegistryOperatorAddedToQuorumsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractBLSApkRegistryOperatorAddedToQuorums)
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
		it.Event = new(ContractBLSApkRegistryOperatorAddedToQuorums)
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
func (it *ContractBLSApkRegistryOperatorAddedToQuorumsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractBLSApkRegistryOperatorAddedToQuorumsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractBLSApkRegistryOperatorAddedToQuorums represents a OperatorAddedToQuorums event raised by the ContractBLSApkRegistry contract.
type ContractBLSApkRegistryOperatorAddedToQuorums struct {
	Operator      common.Address
	QuorumNumbers []byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOperatorAddedToQuorums is a free log retrieval operation binding the contract event 0x5358c5b42179178c8fc757734ac2a3198f9071c765ee0d8389211525f5005246.
//
// Solidity: event OperatorAddedToQuorums(address operator, bytes quorumNumbers)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryFilterer) FilterOperatorAddedToQuorums(opts *bind.FilterOpts) (*ContractBLSApkRegistryOperatorAddedToQuorumsIterator, error) {

	logs, sub, err := _ContractBLSApkRegistry.contract.FilterLogs(opts, "OperatorAddedToQuorums")
	if err != nil {
		return nil, err
	}
	return &ContractBLSApkRegistryOperatorAddedToQuorumsIterator{contract: _ContractBLSApkRegistry.contract, event: "OperatorAddedToQuorums", logs: logs, sub: sub}, nil
}

// WatchOperatorAddedToQuorums is a free log subscription operation binding the contract event 0x5358c5b42179178c8fc757734ac2a3198f9071c765ee0d8389211525f5005246.
//
// Solidity: event OperatorAddedToQuorums(address operator, bytes quorumNumbers)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryFilterer) WatchOperatorAddedToQuorums(opts *bind.WatchOpts, sink chan<- *ContractBLSApkRegistryOperatorAddedToQuorums) (event.Subscription, error) {

	logs, sub, err := _ContractBLSApkRegistry.contract.WatchLogs(opts, "OperatorAddedToQuorums")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractBLSApkRegistryOperatorAddedToQuorums)
				if err := _ContractBLSApkRegistry.contract.UnpackLog(event, "OperatorAddedToQuorums", log); err != nil {
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

// ParseOperatorAddedToQuorums is a log parse operation binding the contract event 0x5358c5b42179178c8fc757734ac2a3198f9071c765ee0d8389211525f5005246.
//
// Solidity: event OperatorAddedToQuorums(address operator, bytes quorumNumbers)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryFilterer) ParseOperatorAddedToQuorums(log types.Log) (*ContractBLSApkRegistryOperatorAddedToQuorums, error) {
	event := new(ContractBLSApkRegistryOperatorAddedToQuorums)
	if err := _ContractBLSApkRegistry.contract.UnpackLog(event, "OperatorAddedToQuorums", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractBLSApkRegistryOperatorRemovedFromQuorumsIterator is returned from FilterOperatorRemovedFromQuorums and is used to iterate over the raw logs and unpacked data for OperatorRemovedFromQuorums events raised by the ContractBLSApkRegistry contract.
type ContractBLSApkRegistryOperatorRemovedFromQuorumsIterator struct {
	Event *ContractBLSApkRegistryOperatorRemovedFromQuorums // Event containing the contract specifics and raw log

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
func (it *ContractBLSApkRegistryOperatorRemovedFromQuorumsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractBLSApkRegistryOperatorRemovedFromQuorums)
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
		it.Event = new(ContractBLSApkRegistryOperatorRemovedFromQuorums)
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
func (it *ContractBLSApkRegistryOperatorRemovedFromQuorumsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractBLSApkRegistryOperatorRemovedFromQuorumsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractBLSApkRegistryOperatorRemovedFromQuorums represents a OperatorRemovedFromQuorums event raised by the ContractBLSApkRegistry contract.
type ContractBLSApkRegistryOperatorRemovedFromQuorums struct {
	Operator      common.Address
	QuorumNumbers []byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOperatorRemovedFromQuorums is a free log retrieval operation binding the contract event 0x14a5172b312e9d2c22b8468f9c70ec2caa9de934fe380734fbc6f3beff2b14ba.
//
// Solidity: event OperatorRemovedFromQuorums(address operator, bytes quorumNumbers)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryFilterer) FilterOperatorRemovedFromQuorums(opts *bind.FilterOpts) (*ContractBLSApkRegistryOperatorRemovedFromQuorumsIterator, error) {

	logs, sub, err := _ContractBLSApkRegistry.contract.FilterLogs(opts, "OperatorRemovedFromQuorums")
	if err != nil {
		return nil, err
	}
	return &ContractBLSApkRegistryOperatorRemovedFromQuorumsIterator{contract: _ContractBLSApkRegistry.contract, event: "OperatorRemovedFromQuorums", logs: logs, sub: sub}, nil
}

// WatchOperatorRemovedFromQuorums is a free log subscription operation binding the contract event 0x14a5172b312e9d2c22b8468f9c70ec2caa9de934fe380734fbc6f3beff2b14ba.
//
// Solidity: event OperatorRemovedFromQuorums(address operator, bytes quorumNumbers)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryFilterer) WatchOperatorRemovedFromQuorums(opts *bind.WatchOpts, sink chan<- *ContractBLSApkRegistryOperatorRemovedFromQuorums) (event.Subscription, error) {

	logs, sub, err := _ContractBLSApkRegistry.contract.WatchLogs(opts, "OperatorRemovedFromQuorums")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractBLSApkRegistryOperatorRemovedFromQuorums)
				if err := _ContractBLSApkRegistry.contract.UnpackLog(event, "OperatorRemovedFromQuorums", log); err != nil {
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

// ParseOperatorRemovedFromQuorums is a log parse operation binding the contract event 0x14a5172b312e9d2c22b8468f9c70ec2caa9de934fe380734fbc6f3beff2b14ba.
//
// Solidity: event OperatorRemovedFromQuorums(address operator, bytes quorumNumbers)
func (_ContractBLSApkRegistry *ContractBLSApkRegistryFilterer) ParseOperatorRemovedFromQuorums(log types.Log) (*ContractBLSApkRegistryOperatorRemovedFromQuorums, error) {
	event := new(ContractBLSApkRegistryOperatorRemovedFromQuorums)
	if err := _ContractBLSApkRegistry.contract.UnpackLog(event, "OperatorRemovedFromQuorums", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
