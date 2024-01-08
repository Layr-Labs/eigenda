// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractMockRollup

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

// EigenDABlobUtilsBlobVerificationProof is an auto generated low-level Go binding around an user-defined struct.
type EigenDABlobUtilsBlobVerificationProof struct {
	BatchId                uint32
	BlobIndex              uint8
	BatchMetadata          IEigenDAServiceManagerBatchMetadata
	InclusionProof         []byte
	QuorumThresholdIndexes []byte
}

// IEigenDAServiceManagerBatchHeader is an auto generated low-level Go binding around an user-defined struct.
type IEigenDAServiceManagerBatchHeader struct {
	BlobHeadersRoot            [32]byte
	QuorumNumbers              []byte
	QuorumThresholdPercentages []byte
	ReferenceBlockNumber       uint32
}

// IEigenDAServiceManagerBatchMetadata is an auto generated low-level Go binding around an user-defined struct.
type IEigenDAServiceManagerBatchMetadata struct {
	BatchHeader             IEigenDAServiceManagerBatchHeader
	SignatoryRecordHash     [32]byte
	Fee                     *big.Int
	ConfirmationBlockNumber uint32
}

// IEigenDAServiceManagerBlobHeader is an auto generated low-level Go binding around an user-defined struct.
type IEigenDAServiceManagerBlobHeader struct {
	Commitment       BN254G1Point
	DataLength       uint32
	QuorumBlobParams []IEigenDAServiceManagerQuorumBlobParam
}

// IEigenDAServiceManagerQuorumBlobParam is an auto generated low-level Go binding around an user-defined struct.
type IEigenDAServiceManagerQuorumBlobParam struct {
	QuorumNumber                 uint8
	AdversaryThresholdPercentage uint8
	QuorumThresholdPercentage    uint8
	QuantizationParameter        uint8
}

// ContractMockRollupMetaData contains all meta data concerning the ContractMockRollup contract.
var ContractMockRollupMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_eigenDAServiceManager\",\"type\":\"address\",\"internalType\":\"contractIEigenDAServiceManager\"},{\"name\":\"_tau\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"_illegalValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_stakeRequired\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"blacklist\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"challengeCommitment\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"point\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"commitments\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"polynomialCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"computeQuorumBlobParamsHash\",\"inputs\":[{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structIEigenDAServiceManager.QuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"quorumThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"quantizationParameter\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"deRegisterValidator\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"eigenDAServiceManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAServiceManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"illegalValue\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"postCommitment\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structIEigenDAServiceManager.QuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"quorumThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"quantizationParameter\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structEigenDABlobUtils.BlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumThresholdPercentages\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"fee\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumThresholdIndexes\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"blobParamsHashInput\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerValidator\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"stakeRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tau\",\"inputs\":[],\"outputs\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validators\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405162001c0f38038062001c0f8339810160408190526100319161006b565b600080546001600160a01b0319166001600160a01b0395909516949094179093558151600155602090910151600255600355600455610109565b60008060008084860360a081121561008257600080fd5b85516001600160a01b038116811461009957600080fd5b94506040601f19820112156100ad57600080fd5b50604080519081016001600160401b03811182821017156100de57634e487b7160e01b600052604160045260246000fd5b6040908152602087810151835290870151908201526060860151608090960151949790965092505050565b611af680620001196000396000f3fe6080604052600436106100a65760003560e01c80636281e63b116100645780636281e63b146101b8578063bcc6587f146101d8578063cfc4af55146101e0578063f9f92be414610210578063fa52c7d814610250578063fc30cad01461028057600080fd5b8062363779146100ab57806305010105146100de5780632ff46cf3146100f45780634440bc5c1461010b57806349ce899714610121578063586f01bb14610198575b600080fd5b3480156100b757600080fd5b506100cb6100c636600461126a565b6102b8565b6040519081526020015b60405180910390f35b3480156100ea57600080fd5b506100cb60045481565b34801561010057600080fd5b506101096103ac565b005b34801561011757600080fd5b506100cb60035481565b34801561012d57600080fd5b5061018961013c36600461129f565b6007602090815260009182526040918290208054835180850190945260018201548452600290910154918301919091526001600160a01b03811691600160a01b90910463ffffffff169083565b6040516100d5939291906112b8565b3480156101a457600080fd5b506101096101b3366004611524565b61042d565b3480156101c457600080fd5b506101096101d3366004611652565b6106bf565b61010961095c565b3480156101ec57600080fd5b506001546002546101fb919082565b604080519283526020830191909152016100d5565b34801561021c57600080fd5b5061024061022b3660046116ba565b60066020526000908152604090205460ff1681565b60405190151581526020016100d5565b34801561025c57600080fd5b5061024061026b3660046116ba565b60056020526000908152604090205460ff1681565b34801561028c57600080fd5b506000546102a0906001600160a01b031681565b6040516001600160a01b0390911681526020016100d5565b60008082516040516020016102cf91815260200190565b604051602081830303815290604052905060005b835181101561039d57818482815181106102ff576102ff6116ea565b60200260200101516000015185838151811061031d5761031d6116ea565b60200260200101516020015186848151811061033b5761033b6116ea565b602002602001015160400151878581518110610359576103596116ea565b602002602001015160600151604051602001610379959493929190611730565b6040516020818303038152906040529150808061039590611796565b9150506102e3565b50805160209091012092915050565b3360009081526005602052604090205460ff166103e45760405162461bcd60e51b81526004016103db906117b1565b60405180910390fd5b3360009081526006602052604090205460ff16156104145760405162461bcd60e51b81526004016103db9061180e565b336000908152600560205260409020805460ff19169055565b3360009081526005602052604090205460ff166104a85760405162461bcd60e51b815260206004820152603360248201527f4d6f636b526f6c6c75702e706f7374436f6d6d69746d656e743a2056616c6964604482015272185d1bdc881b9bdd081c9959da5cdd195c9959606a1b60648201526084016103db565b426000908152600760205260409020546001600160a01b03161561052b5760405162461bcd60e51b815260206004820152603460248201527f4d6f636b526f6c6c75702e706f7374436f6d6d69746d656e743a20436f6d6d696044820152731d1b595b9d08185b1c9958591e481c1bdcdd195960621b60648201526084016103db565b600061053a84604001516102b8565b90508181146105c75760405162461bcd60e51b815260206004820152604d60248201527f4d6f636b526f6c6c75702e706f7374436f6d6d69746d656e743a2051756f727560448201527f6d426c6f62506172616d7320646f206e6f74206d617463682071756f72756d4260648201526c0d8dec4a0c2e4c2dae690c2e6d609b1b608482015260a4016103db565b600054604051633ddc42f360e11b815273__$3652595bba153b674dc03d93006671d207$__91637bb885e69161060f9188916001600160a01b0390911690889060040161198a565b60006040518083038186803b15801561062757600080fd5b505af415801561063b573d6000803e3d6000fd5b50506040805160608101825233815260208881015163ffffffff90811682840190815299518385019081524260009081526007845294909420925183549a51909116600160a01b026001600160c01b0319909a166001600160a01b039190911617989098178155905180516001830155909601516002909601959095555050505050565b600083815260076020908152604091829020825160608101845281546001600160a01b038082168352600160a01b90910463ffffffff1682850152845180860186526001840154815260029093015493830193909352928301528151166107865760405162461bcd60e51b815260206004820152603560248201527f4d6f636b526f6c6c75702e6368616c6c656e6765436f6d6d69746d656e743a2060448201527410dbdb5b5a5d1b595b9d081b9bdd081c1bdcdd1959605a1b60648201526084016103db565b806020015163ffffffff1683106108115760405162461bcd60e51b815260206004820152604360248201527f4d6f636b526f6c6c75702e6368616c6c656e6765436f6d6d69746d656e743a2060448201527f506f696e74206d757374206265206c657373207468616e2064617461206c656e6064820152620cee8d60eb1b608482015260a4016103db565b600354604080518082018252600154815260025460208201529083015161083c928692909186610a5a565b6108b95760405162461bcd60e51b815260206004820152604260248201527f4d6f636b526f6c6c75702e6368616c6c656e6765436f6d6d69746d656e743a2060448201527f446f6573206e6f74206576616c7561746520746f20696c6c6567616c2076616c606482015261756560f01b608482015260a4016103db565b80516001600160a01b039081166000908152600560209081526040808320805460ff19908116909155855190941683526006909152808220805490931660011790925590513390670de0b6b3a7640000908381818185875af1925050503d8060008114610942576040519150601f19603f3d011682016040523d82523d6000602084013e610947565b606091505b505090508061095557600080fd5b5050505050565b60045434146109de5760405162461bcd60e51b815260206004820152604260248201527f4d6f636b526f6c6c75702e726567697374657256616c696461746f723a204d7560448201527f73742073656e64207374616b6520726571756972656420746f2072656769737460648201526132b960f11b608482015260a4016103db565b3360009081526005602052604090205460ff1615610a0e5760405162461bcd60e51b81526004016103db906117b1565b3360009081526006602052604090205460ff1615610a3e5760405162461bcd60e51b81526004016103db9061180e565b336000908152600560205260409020805460ff19166001179055565b600080610a91610a8c604080518082018252600080825260209182015281518083019092526001825260029082015290565b610ad7565b9050610acc610aaa610aa3838a610b96565b8790610c2d565b84610abf610ab8858b610b96565b8890610c2d565b610ac7610cc1565b610d81565b979650505050505050565b60408051808201909152600080825260208201528151158015610afc57506020820151155b15610b1a575050604080518082019091526000808252602082015290565b6040518060400160405280836000015181526020017f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478460200151610b5f9190611a50565b610b89907f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47611a72565b905292915050565b919050565b6040805180820190915260008082526020820152610bb2610ff0565b835181526020808501519082015260408082018490526000908360608460076107d05a03fa9050808015610be557610be7565bfe5b5080610c255760405162461bcd60e51b815260206004820152600d60248201526c1958cb5b5d5b0b59985a5b1959609a1b60448201526064016103db565b505092915050565b6040805180820190915260008082526020820152610c4961100e565b835181526020808501518183015283516040808401919091529084015160608301526000908360808460066107d05a03fa9050808015610be5575080610c255760405162461bcd60e51b815260206004820152600d60248201526c1958cb5859190b59985a5b1959609a1b60448201526064016103db565b610cc961102c565b50604080516080810182527f198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c28183019081527f1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed6060830152815281518083019092527f275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec82527f1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d60208381019190915281019190915290565b604080518082018252858152602080820185905282518084019093528583528201839052600091610db0611051565b60005b6002811015610f75576000610dc9826006611a89565b9050848260028110610ddd57610ddd6116ea565b60200201515183610def836000611aa8565b600c8110610dff57610dff6116ea565b6020020152848260028110610e1657610e166116ea565b60200201516020015183826001610e2d9190611aa8565b600c8110610e3d57610e3d6116ea565b6020020152838260028110610e5457610e546116ea565b6020020151515183610e67836002611aa8565b600c8110610e7757610e776116ea565b6020020152838260028110610e8e57610e8e6116ea565b6020020151516001602002015183610ea7836003611aa8565b600c8110610eb757610eb76116ea565b6020020152838260028110610ece57610ece6116ea565b602002015160200151600060028110610ee957610ee96116ea565b602002015183610efa836004611aa8565b600c8110610f0a57610f0a6116ea565b6020020152838260028110610f2157610f216116ea565b602002015160200151600160028110610f3c57610f3c6116ea565b602002015183610f4d836005611aa8565b600c8110610f5d57610f5d6116ea565b60200201525080610f6d81611796565b915050610db3565b50610f7e611070565b60006020826101808560086107d05a03fa9050808015610be5575080610fde5760405162461bcd60e51b81526020600482015260156024820152741c185a5c9a5b99cb5bdc18dbd9194b59985a5b1959605a1b60448201526064016103db565b5051151593505050505b949350505050565b60405180606001604052806003906020820280368337509192915050565b60405180608001604052806004906020820280368337509192915050565b604051806040016040528061103f61108e565b815260200161104c61108e565b905290565b604051806101800160405280600c906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b60405180604001604052806002906020820280368337509192915050565b634e487b7160e01b600052604160045260246000fd5b6040516080810167ffffffffffffffff811182821017156110e5576110e56110ac565b60405290565b60405160a0810167ffffffffffffffff811182821017156110e5576110e56110ac565b6040516060810167ffffffffffffffff811182821017156110e5576110e56110ac565b6040805190810167ffffffffffffffff811182821017156110e5576110e56110ac565b604051601f8201601f1916810167ffffffffffffffff8111828210171561117d5761117d6110ac565b604052919050565b803560ff81168114610b9157600080fd5b600082601f8301126111a757600080fd5b8135602067ffffffffffffffff8211156111c3576111c36110ac565b6111d1818360051b01611154565b82815260079290921b840181019181810190868411156111f057600080fd5b8286015b8481101561125f576080818903121561120d5760008081fd5b6112156110c2565b61121e82611185565b815261122b858301611185565b85820152604061123c818401611185565b90820152606061124d838201611185565b908201528352918301916080016111f4565b509695505050505050565b60006020828403121561127c57600080fd5b813567ffffffffffffffff81111561129357600080fd5b610fe884828501611196565b6000602082840312156112b157600080fd5b5035919050565b6001600160a01b038416815263ffffffff8316602082015260808101610fe8604083018480518252602090810151910152565b803563ffffffff81168114610b9157600080fd5b600082601f83011261131057600080fd5b813567ffffffffffffffff81111561132a5761132a6110ac565b61133d601f8201601f1916602001611154565b81815284602083860101111561135257600080fd5b816020850160208301376000918101602001919091529392505050565b80356bffffffffffffffffffffffff81168114610b9157600080fd5b60006080828403121561139d57600080fd5b6113a56110c2565b9050813567ffffffffffffffff808211156113bf57600080fd5b90830190608082860312156113d357600080fd5b6113db6110c2565b823581526020830135828111156113f157600080fd5b6113fd878286016112ff565b60208301525060408301358281111561141557600080fd5b611421878286016112ff565b604083015250611433606084016112eb565b606082015283525050602082810135908201526114526040830161136f565b6040820152611463606083016112eb565b606082015292915050565b600060a0828403121561148057600080fd5b6114886110eb565b9050611493826112eb565b81526114a160208301611185565b6020820152604082013567ffffffffffffffff808211156114c157600080fd5b6114cd8583860161138b565b604084015260608401359150808211156114e657600080fd5b6114f2858386016112ff565b6060840152608084013591508082111561150b57600080fd5b50611518848285016112ff565b60808301525092915050565b60008060006060848603121561153957600080fd5b833567ffffffffffffffff8082111561155157600080fd5b90850190818703608081121561156657600080fd5b61156e61110e565b604082121561157c57600080fd5b611584611131565b915083358252602084013560208301528181526115a3604085016112eb565b602082015260608401359150828211156115bc57600080fd5b6115c889838601611196565b604082015295505060208601359150808211156115e457600080fd5b506115f18682870161146e565b925050604084013590509250925092565b600082601f83011261161357600080fd5b61161b611131565b80604084018581111561162d57600080fd5b845b8181101561164757803584526020938401930161162f565b509095945050505050565b600080600083850360c081121561166857600080fd5b84359350602085013592506080603f198201121561168557600080fd5b5061168e611131565b61169b8660408701611602565b81526116aa8660808701611602565b6020820152809150509250925092565b6000602082840312156116cc57600080fd5b81356001600160a01b03811681146116e357600080fd5b9392505050565b634e487b7160e01b600052603260045260246000fd5b60005b8381101561171b578181015183820152602001611703565b8381111561172a576000848401525b50505050565b60008651611742818460208b01611700565b60f896871b6001600160f81b03199081169390910192835294861b851660018301525091841b8316600283015290921b166003820152600401919050565b634e487b7160e01b600052601160045260246000fd5b60006000198214156117aa576117aa611780565b5060010190565b6020808252603a908201527f4d6f636b526f6c6c75702e726567697374657256616c696461746f723a20566160408201527f6c696461746f7220616c72656164792072656769737465726564000000000000606082015260800190565b60208082526033908201527f4d6f636b526f6c6c75702e726567697374657256616c696461746f723a2056616040820152721b1a59185d1bdc88189b1858dadb1a5cdd1959606a1b606082015260800190565b60008151808452611879816020860160208601611700565b601f01601f19169290920160200192915050565b600063ffffffff80835116845260ff6020840151166020850152604083015160a060408601528051608060a08701528051610120870152602081015160806101408801526118df6101a0880182611861565b9050604082015161011f19888303016101608901526118fe8282611861565b60608401519095166101808901525050602082015160c087015260408201519261193860e08801856bffffffffffffffffffffffff169052565b606083015163ffffffff811661010089015293506060860151935086810360608801526119658185611861565b9350505050608083015184820360808601526119818282611861565b95945050505050565b6060808252845180518383015260200151608083015260009060e0830160208781015163ffffffff1660a0860152604080890151608060c0880181905281519485905290830193600091906101008901905b80841015611a20578651805160ff908116845287820151811688850152868201518116878501529089015116888301529585019560019390930192908201906119dc565b506001600160a01b038b1689860152888103848a0152611a40818b61188d565b9c9b505050505050505050505050565b600082611a6d57634e487b7160e01b600052601260045260246000fd5b500690565b600082821015611a8457611a84611780565b500390565b6000816000190483118215151615611aa357611aa3611780565b500290565b60008219821115611abb57611abb611780565b50019056fea2646970667358221220eaf8dd586640d4c481a7eef54a3009b336226201b645c2f3c33def8322a641be64736f6c634300080c0033",
}

// ContractMockRollupABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractMockRollupMetaData.ABI instead.
var ContractMockRollupABI = ContractMockRollupMetaData.ABI

// ContractMockRollupBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractMockRollupMetaData.Bin instead.
var ContractMockRollupBin = ContractMockRollupMetaData.Bin

// DeployContractMockRollup deploys a new Ethereum contract, binding an instance of ContractMockRollup to it.
func DeployContractMockRollup(auth *bind.TransactOpts, backend bind.ContractBackend, _eigenDAServiceManager common.Address, _tau BN254G1Point, _illegalValue *big.Int, _stakeRequired *big.Int) (common.Address, *types.Transaction, *ContractMockRollup, error) {
	parsed, err := ContractMockRollupMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractMockRollupBin), backend, _eigenDAServiceManager, _tau, _illegalValue, _stakeRequired)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractMockRollup{ContractMockRollupCaller: ContractMockRollupCaller{contract: contract}, ContractMockRollupTransactor: ContractMockRollupTransactor{contract: contract}, ContractMockRollupFilterer: ContractMockRollupFilterer{contract: contract}}, nil
}

// ContractMockRollup is an auto generated Go binding around an Ethereum contract.
type ContractMockRollup struct {
	ContractMockRollupCaller     // Read-only binding to the contract
	ContractMockRollupTransactor // Write-only binding to the contract
	ContractMockRollupFilterer   // Log filterer for contract events
}

// ContractMockRollupCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractMockRollupCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractMockRollupTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractMockRollupTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractMockRollupFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractMockRollupFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractMockRollupSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractMockRollupSession struct {
	Contract     *ContractMockRollup // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ContractMockRollupCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractMockRollupCallerSession struct {
	Contract *ContractMockRollupCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// ContractMockRollupTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractMockRollupTransactorSession struct {
	Contract     *ContractMockRollupTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// ContractMockRollupRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractMockRollupRaw struct {
	Contract *ContractMockRollup // Generic contract binding to access the raw methods on
}

// ContractMockRollupCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractMockRollupCallerRaw struct {
	Contract *ContractMockRollupCaller // Generic read-only contract binding to access the raw methods on
}

// ContractMockRollupTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractMockRollupTransactorRaw struct {
	Contract *ContractMockRollupTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractMockRollup creates a new instance of ContractMockRollup, bound to a specific deployed contract.
func NewContractMockRollup(address common.Address, backend bind.ContractBackend) (*ContractMockRollup, error) {
	contract, err := bindContractMockRollup(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractMockRollup{ContractMockRollupCaller: ContractMockRollupCaller{contract: contract}, ContractMockRollupTransactor: ContractMockRollupTransactor{contract: contract}, ContractMockRollupFilterer: ContractMockRollupFilterer{contract: contract}}, nil
}

// NewContractMockRollupCaller creates a new read-only instance of ContractMockRollup, bound to a specific deployed contract.
func NewContractMockRollupCaller(address common.Address, caller bind.ContractCaller) (*ContractMockRollupCaller, error) {
	contract, err := bindContractMockRollup(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractMockRollupCaller{contract: contract}, nil
}

// NewContractMockRollupTransactor creates a new write-only instance of ContractMockRollup, bound to a specific deployed contract.
func NewContractMockRollupTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractMockRollupTransactor, error) {
	contract, err := bindContractMockRollup(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractMockRollupTransactor{contract: contract}, nil
}

// NewContractMockRollupFilterer creates a new log filterer instance of ContractMockRollup, bound to a specific deployed contract.
func NewContractMockRollupFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractMockRollupFilterer, error) {
	contract, err := bindContractMockRollup(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractMockRollupFilterer{contract: contract}, nil
}

// bindContractMockRollup binds a generic wrapper to an already deployed contract.
func bindContractMockRollup(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractMockRollupMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractMockRollup *ContractMockRollupRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractMockRollup.Contract.ContractMockRollupCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractMockRollup *ContractMockRollupRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.ContractMockRollupTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractMockRollup *ContractMockRollupRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.ContractMockRollupTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractMockRollup *ContractMockRollupCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractMockRollup.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractMockRollup *ContractMockRollupTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractMockRollup *ContractMockRollupTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.contract.Transact(opts, method, params...)
}

// Blacklist is a free data retrieval call binding the contract method 0xf9f92be4.
//
// Solidity: function blacklist(address ) view returns(bool)
func (_ContractMockRollup *ContractMockRollupCaller) Blacklist(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _ContractMockRollup.contract.Call(opts, &out, "blacklist", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Blacklist is a free data retrieval call binding the contract method 0xf9f92be4.
//
// Solidity: function blacklist(address ) view returns(bool)
func (_ContractMockRollup *ContractMockRollupSession) Blacklist(arg0 common.Address) (bool, error) {
	return _ContractMockRollup.Contract.Blacklist(&_ContractMockRollup.CallOpts, arg0)
}

// Blacklist is a free data retrieval call binding the contract method 0xf9f92be4.
//
// Solidity: function blacklist(address ) view returns(bool)
func (_ContractMockRollup *ContractMockRollupCallerSession) Blacklist(arg0 common.Address) (bool, error) {
	return _ContractMockRollup.Contract.Blacklist(&_ContractMockRollup.CallOpts, arg0)
}

// Commitments is a free data retrieval call binding the contract method 0x49ce8997.
//
// Solidity: function commitments(uint256 ) view returns(address validator, uint32 dataLength, (uint256,uint256) polynomialCommitment)
func (_ContractMockRollup *ContractMockRollupCaller) Commitments(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Validator            common.Address
	DataLength           uint32
	PolynomialCommitment BN254G1Point
}, error) {
	var out []interface{}
	err := _ContractMockRollup.contract.Call(opts, &out, "commitments", arg0)

	outstruct := new(struct {
		Validator            common.Address
		DataLength           uint32
		PolynomialCommitment BN254G1Point
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Validator = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.DataLength = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.PolynomialCommitment = *abi.ConvertType(out[2], new(BN254G1Point)).(*BN254G1Point)

	return *outstruct, err

}

// Commitments is a free data retrieval call binding the contract method 0x49ce8997.
//
// Solidity: function commitments(uint256 ) view returns(address validator, uint32 dataLength, (uint256,uint256) polynomialCommitment)
func (_ContractMockRollup *ContractMockRollupSession) Commitments(arg0 *big.Int) (struct {
	Validator            common.Address
	DataLength           uint32
	PolynomialCommitment BN254G1Point
}, error) {
	return _ContractMockRollup.Contract.Commitments(&_ContractMockRollup.CallOpts, arg0)
}

// Commitments is a free data retrieval call binding the contract method 0x49ce8997.
//
// Solidity: function commitments(uint256 ) view returns(address validator, uint32 dataLength, (uint256,uint256) polynomialCommitment)
func (_ContractMockRollup *ContractMockRollupCallerSession) Commitments(arg0 *big.Int) (struct {
	Validator            common.Address
	DataLength           uint32
	PolynomialCommitment BN254G1Point
}, error) {
	return _ContractMockRollup.Contract.Commitments(&_ContractMockRollup.CallOpts, arg0)
}

// ComputeQuorumBlobParamsHash is a free data retrieval call binding the contract method 0x00363779.
//
// Solidity: function computeQuorumBlobParamsHash((uint8,uint8,uint8,uint8)[] quorumBlobParams) pure returns(bytes32)
func (_ContractMockRollup *ContractMockRollupCaller) ComputeQuorumBlobParamsHash(opts *bind.CallOpts, quorumBlobParams []IEigenDAServiceManagerQuorumBlobParam) ([32]byte, error) {
	var out []interface{}
	err := _ContractMockRollup.contract.Call(opts, &out, "computeQuorumBlobParamsHash", quorumBlobParams)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ComputeQuorumBlobParamsHash is a free data retrieval call binding the contract method 0x00363779.
//
// Solidity: function computeQuorumBlobParamsHash((uint8,uint8,uint8,uint8)[] quorumBlobParams) pure returns(bytes32)
func (_ContractMockRollup *ContractMockRollupSession) ComputeQuorumBlobParamsHash(quorumBlobParams []IEigenDAServiceManagerQuorumBlobParam) ([32]byte, error) {
	return _ContractMockRollup.Contract.ComputeQuorumBlobParamsHash(&_ContractMockRollup.CallOpts, quorumBlobParams)
}

// ComputeQuorumBlobParamsHash is a free data retrieval call binding the contract method 0x00363779.
//
// Solidity: function computeQuorumBlobParamsHash((uint8,uint8,uint8,uint8)[] quorumBlobParams) pure returns(bytes32)
func (_ContractMockRollup *ContractMockRollupCallerSession) ComputeQuorumBlobParamsHash(quorumBlobParams []IEigenDAServiceManagerQuorumBlobParam) ([32]byte, error) {
	return _ContractMockRollup.Contract.ComputeQuorumBlobParamsHash(&_ContractMockRollup.CallOpts, quorumBlobParams)
}

// EigenDAServiceManager is a free data retrieval call binding the contract method 0xfc30cad0.
//
// Solidity: function eigenDAServiceManager() view returns(address)
func (_ContractMockRollup *ContractMockRollupCaller) EigenDAServiceManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractMockRollup.contract.Call(opts, &out, "eigenDAServiceManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDAServiceManager is a free data retrieval call binding the contract method 0xfc30cad0.
//
// Solidity: function eigenDAServiceManager() view returns(address)
func (_ContractMockRollup *ContractMockRollupSession) EigenDAServiceManager() (common.Address, error) {
	return _ContractMockRollup.Contract.EigenDAServiceManager(&_ContractMockRollup.CallOpts)
}

// EigenDAServiceManager is a free data retrieval call binding the contract method 0xfc30cad0.
//
// Solidity: function eigenDAServiceManager() view returns(address)
func (_ContractMockRollup *ContractMockRollupCallerSession) EigenDAServiceManager() (common.Address, error) {
	return _ContractMockRollup.Contract.EigenDAServiceManager(&_ContractMockRollup.CallOpts)
}

// IllegalValue is a free data retrieval call binding the contract method 0x4440bc5c.
//
// Solidity: function illegalValue() view returns(uint256)
func (_ContractMockRollup *ContractMockRollupCaller) IllegalValue(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractMockRollup.contract.Call(opts, &out, "illegalValue")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// IllegalValue is a free data retrieval call binding the contract method 0x4440bc5c.
//
// Solidity: function illegalValue() view returns(uint256)
func (_ContractMockRollup *ContractMockRollupSession) IllegalValue() (*big.Int, error) {
	return _ContractMockRollup.Contract.IllegalValue(&_ContractMockRollup.CallOpts)
}

// IllegalValue is a free data retrieval call binding the contract method 0x4440bc5c.
//
// Solidity: function illegalValue() view returns(uint256)
func (_ContractMockRollup *ContractMockRollupCallerSession) IllegalValue() (*big.Int, error) {
	return _ContractMockRollup.Contract.IllegalValue(&_ContractMockRollup.CallOpts)
}

// StakeRequired is a free data retrieval call binding the contract method 0x05010105.
//
// Solidity: function stakeRequired() view returns(uint256)
func (_ContractMockRollup *ContractMockRollupCaller) StakeRequired(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractMockRollup.contract.Call(opts, &out, "stakeRequired")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StakeRequired is a free data retrieval call binding the contract method 0x05010105.
//
// Solidity: function stakeRequired() view returns(uint256)
func (_ContractMockRollup *ContractMockRollupSession) StakeRequired() (*big.Int, error) {
	return _ContractMockRollup.Contract.StakeRequired(&_ContractMockRollup.CallOpts)
}

// StakeRequired is a free data retrieval call binding the contract method 0x05010105.
//
// Solidity: function stakeRequired() view returns(uint256)
func (_ContractMockRollup *ContractMockRollupCallerSession) StakeRequired() (*big.Int, error) {
	return _ContractMockRollup.Contract.StakeRequired(&_ContractMockRollup.CallOpts)
}

// Tau is a free data retrieval call binding the contract method 0xcfc4af55.
//
// Solidity: function tau() view returns(uint256 X, uint256 Y)
func (_ContractMockRollup *ContractMockRollupCaller) Tau(opts *bind.CallOpts) (struct {
	X *big.Int
	Y *big.Int
}, error) {
	var out []interface{}
	err := _ContractMockRollup.contract.Call(opts, &out, "tau")

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

// Tau is a free data retrieval call binding the contract method 0xcfc4af55.
//
// Solidity: function tau() view returns(uint256 X, uint256 Y)
func (_ContractMockRollup *ContractMockRollupSession) Tau() (struct {
	X *big.Int
	Y *big.Int
}, error) {
	return _ContractMockRollup.Contract.Tau(&_ContractMockRollup.CallOpts)
}

// Tau is a free data retrieval call binding the contract method 0xcfc4af55.
//
// Solidity: function tau() view returns(uint256 X, uint256 Y)
func (_ContractMockRollup *ContractMockRollupCallerSession) Tau() (struct {
	X *big.Int
	Y *big.Int
}, error) {
	return _ContractMockRollup.Contract.Tau(&_ContractMockRollup.CallOpts)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) view returns(bool)
func (_ContractMockRollup *ContractMockRollupCaller) Validators(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _ContractMockRollup.contract.Call(opts, &out, "validators", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) view returns(bool)
func (_ContractMockRollup *ContractMockRollupSession) Validators(arg0 common.Address) (bool, error) {
	return _ContractMockRollup.Contract.Validators(&_ContractMockRollup.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) view returns(bool)
func (_ContractMockRollup *ContractMockRollupCallerSession) Validators(arg0 common.Address) (bool, error) {
	return _ContractMockRollup.Contract.Validators(&_ContractMockRollup.CallOpts, arg0)
}

// ChallengeCommitment is a paid mutator transaction binding the contract method 0x6281e63b.
//
// Solidity: function challengeCommitment(uint256 timestamp, uint256 point, (uint256[2],uint256[2]) proof) returns()
func (_ContractMockRollup *ContractMockRollupTransactor) ChallengeCommitment(opts *bind.TransactOpts, timestamp *big.Int, point *big.Int, proof BN254G2Point) (*types.Transaction, error) {
	return _ContractMockRollup.contract.Transact(opts, "challengeCommitment", timestamp, point, proof)
}

// ChallengeCommitment is a paid mutator transaction binding the contract method 0x6281e63b.
//
// Solidity: function challengeCommitment(uint256 timestamp, uint256 point, (uint256[2],uint256[2]) proof) returns()
func (_ContractMockRollup *ContractMockRollupSession) ChallengeCommitment(timestamp *big.Int, point *big.Int, proof BN254G2Point) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.ChallengeCommitment(&_ContractMockRollup.TransactOpts, timestamp, point, proof)
}

// ChallengeCommitment is a paid mutator transaction binding the contract method 0x6281e63b.
//
// Solidity: function challengeCommitment(uint256 timestamp, uint256 point, (uint256[2],uint256[2]) proof) returns()
func (_ContractMockRollup *ContractMockRollupTransactorSession) ChallengeCommitment(timestamp *big.Int, point *big.Int, proof BN254G2Point) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.ChallengeCommitment(&_ContractMockRollup.TransactOpts, timestamp, point, proof)
}

// DeRegisterValidator is a paid mutator transaction binding the contract method 0x2ff46cf3.
//
// Solidity: function deRegisterValidator() returns()
func (_ContractMockRollup *ContractMockRollupTransactor) DeRegisterValidator(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractMockRollup.contract.Transact(opts, "deRegisterValidator")
}

// DeRegisterValidator is a paid mutator transaction binding the contract method 0x2ff46cf3.
//
// Solidity: function deRegisterValidator() returns()
func (_ContractMockRollup *ContractMockRollupSession) DeRegisterValidator() (*types.Transaction, error) {
	return _ContractMockRollup.Contract.DeRegisterValidator(&_ContractMockRollup.TransactOpts)
}

// DeRegisterValidator is a paid mutator transaction binding the contract method 0x2ff46cf3.
//
// Solidity: function deRegisterValidator() returns()
func (_ContractMockRollup *ContractMockRollupTransactorSession) DeRegisterValidator() (*types.Transaction, error) {
	return _ContractMockRollup.Contract.DeRegisterValidator(&_ContractMockRollup.TransactOpts)
}

// PostCommitment is a paid mutator transaction binding the contract method 0x586f01bb.
//
// Solidity: function postCommitment(((uint256,uint256),uint32,(uint8,uint8,uint8,uint8)[]) blobHeader, (uint32,uint8,((bytes32,bytes,bytes,uint32),bytes32,uint96,uint32),bytes,bytes) blobVerificationProof, bytes32 blobParamsHashInput) returns()
func (_ContractMockRollup *ContractMockRollupTransactor) PostCommitment(opts *bind.TransactOpts, blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof EigenDABlobUtilsBlobVerificationProof, blobParamsHashInput [32]byte) (*types.Transaction, error) {
	return _ContractMockRollup.contract.Transact(opts, "postCommitment", blobHeader, blobVerificationProof, blobParamsHashInput)
}

// PostCommitment is a paid mutator transaction binding the contract method 0x586f01bb.
//
// Solidity: function postCommitment(((uint256,uint256),uint32,(uint8,uint8,uint8,uint8)[]) blobHeader, (uint32,uint8,((bytes32,bytes,bytes,uint32),bytes32,uint96,uint32),bytes,bytes) blobVerificationProof, bytes32 blobParamsHashInput) returns()
func (_ContractMockRollup *ContractMockRollupSession) PostCommitment(blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof EigenDABlobUtilsBlobVerificationProof, blobParamsHashInput [32]byte) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.PostCommitment(&_ContractMockRollup.TransactOpts, blobHeader, blobVerificationProof, blobParamsHashInput)
}

// PostCommitment is a paid mutator transaction binding the contract method 0x586f01bb.
//
// Solidity: function postCommitment(((uint256,uint256),uint32,(uint8,uint8,uint8,uint8)[]) blobHeader, (uint32,uint8,((bytes32,bytes,bytes,uint32),bytes32,uint96,uint32),bytes,bytes) blobVerificationProof, bytes32 blobParamsHashInput) returns()
func (_ContractMockRollup *ContractMockRollupTransactorSession) PostCommitment(blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof EigenDABlobUtilsBlobVerificationProof, blobParamsHashInput [32]byte) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.PostCommitment(&_ContractMockRollup.TransactOpts, blobHeader, blobVerificationProof, blobParamsHashInput)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0xbcc6587f.
//
// Solidity: function registerValidator() payable returns()
func (_ContractMockRollup *ContractMockRollupTransactor) RegisterValidator(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractMockRollup.contract.Transact(opts, "registerValidator")
}

// RegisterValidator is a paid mutator transaction binding the contract method 0xbcc6587f.
//
// Solidity: function registerValidator() payable returns()
func (_ContractMockRollup *ContractMockRollupSession) RegisterValidator() (*types.Transaction, error) {
	return _ContractMockRollup.Contract.RegisterValidator(&_ContractMockRollup.TransactOpts)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0xbcc6587f.
//
// Solidity: function registerValidator() payable returns()
func (_ContractMockRollup *ContractMockRollupTransactorSession) RegisterValidator() (*types.Transaction, error) {
	return _ContractMockRollup.Contract.RegisterValidator(&_ContractMockRollup.TransactOpts)
}
