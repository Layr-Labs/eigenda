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
	ChunkLength                  uint32
}

// ContractMockRollupMetaData contains all meta data concerning the ContractMockRollup contract.
var ContractMockRollupMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_eigenDAServiceManager\",\"type\":\"address\",\"internalType\":\"contractIEigenDAServiceManager\"},{\"name\":\"_tau\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"_illegalValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_quorumBlobParamsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_stakeRequired\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"blacklist\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"challengeCommitment\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"point\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"commitments\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"polynomialCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAServiceManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAServiceManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"illegalValue\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"postCommitment\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structIEigenDAServiceManager.QuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"quorumThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structEigenDABlobUtils.BlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumThresholdPercentages\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"fee\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumThresholdIndexes\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"quorumBlobParamsHash\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerValidator\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"stakeRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tau\",\"inputs\":[],\"outputs\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validators\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"}]",
	Bin: "0x608060405234801561001057600080fd5b50604051611a67380380611a6783398101604081905261002f9161006c565b600080546001600160a01b0319166001600160a01b0396909616959095179094558251600155602090920151600255600355600455600555610115565b600080600080600085870360c081121561008557600080fd5b86516001600160a01b038116811461009c57600080fd5b95506040601f19820112156100b057600080fd5b50604080519081016001600160401b03811182821017156100e157634e487b7160e01b600052604160045260246000fd5b6040908152602088810151835290880151908201526060870151608088015160a09098015196999198509695945092505050565b611943806101246000396000f3fe60806040526004361061009c5760003560e01c80636281e63b116100645780636281e63b1461018f578063bcc6587f146101af578063cfc4af55146101b7578063f9f92be4146101e7578063fa52c7d814610227578063fc30cad01461025757600080fd5b806305010105146100a157806321553525146100ca57806337507265146100ec5780634440bc5c1461010257806349ce899714610118575b600080fd5b3480156100ad57600080fd5b506100b760055481565b6040519081526020015b60405180910390f35b3480156100d657600080fd5b506100ea6100e53660046112f1565b61028f565b005b3480156100f857600080fd5b506100b760045481565b34801561010e57600080fd5b506100b760035481565b34801561012457600080fd5b5061018061013336600461147e565b6008602090815260009182526040918290208054835180850190945260018201548452600290910154918301919091526001600160a01b03811691600160a01b90910463ffffffff169083565b6040516100c193929190611497565b34801561019b57600080fd5b506100ea6101aa36600461151a565b61053f565b6100ea6107dc565b3480156101c357600080fd5b506001546002546101d2919082565b604080519283526020830191909152016100c1565b3480156101f357600080fd5b50610217610202366004611582565b60076020526000908152604090205460ff1681565b60405190151581526020016100c1565b34801561023357600080fd5b50610217610242366004611582565b60066020526000908152604090205460ff1681565b34801561026357600080fd5b50600054610277906001600160a01b031681565b6040516001600160a01b0390911681526020016100c1565b3360009081526006602052604090205460ff1661030f5760405162461bcd60e51b815260206004820152603360248201527f4d6f636b526f6c6c75702e706f7374436f6d6d69746d656e743a2056616c6964604482015272185d1bdc881b9bdd081c9959da5cdd195c9959606a1b60648201526084015b60405180910390fd5b426000908152600860205260409020546001600160a01b0316156103925760405162461bcd60e51b815260206004820152603460248201527f4d6f636b526f6c6c75702e706f7374436f6d6d69746d656e743a20436f6d6d696044820152731d1b595b9d08185b1c9958591e481c1bdcdd195960621b6064820152608401610306565b60045482604001516040516020016103aa91906115b2565b60405160208183030381529060405280519060200120146104495760405162461bcd60e51b815260206004820152604d60248201527f4d6f636b526f6c6c75702e706f7374436f6d6d69746d656e743a2051756f727560448201527f6d426c6f62506172616d7320646f206e6f74206d617463682071756f72756d4260648201526c0d8dec4a0c2e4c2dae690c2e6d609b1b608482015260a401610306565b600054604051633ddc42f360e11b815273__$3652595bba153b674dc03d93006671d207$__91637bb885e6916104919186916001600160a01b0390911690869060040161177a565b60006040518083038186803b1580156104a957600080fd5b505af41580156104bd573d6000803e3d6000fd5b50506040805160608101825233815260208681015163ffffffff90811682840190815297518385019081524260009081526008845294909420925183549851909116600160a01b026001600160c01b03199098166001600160a01b03919091161796909617815590518051600183015590940151600290940193909355505050565b600083815260086020908152604091829020825160608101845281546001600160a01b038082168352600160a01b90910463ffffffff1682850152845180860186526001840154815260029093015493830193909352928301528151166106065760405162461bcd60e51b815260206004820152603560248201527f4d6f636b526f6c6c75702e6368616c6c656e6765436f6d6d69746d656e743a2060448201527410dbdb5b5a5d1b595b9d081b9bdd081c1bdcdd1959605a1b6064820152608401610306565b806020015163ffffffff1683106106915760405162461bcd60e51b815260206004820152604360248201527f4d6f636b526f6c6c75702e6368616c6c656e6765436f6d6d69746d656e743a2060448201527f506f696e74206d757374206265206c657373207468616e2064617461206c656e6064820152620cee8d60eb1b608482015260a401610306565b60035460408051808201825260015481526002546020820152908301516106bc92869290918661097c565b6107395760405162461bcd60e51b815260206004820152604260248201527f4d6f636b526f6c6c75702e6368616c6c656e6765436f6d6d69746d656e743a2060448201527f446f6573206e6f74206576616c7561746520746f20696c6c6567616c2076616c606482015261756560f01b608482015260a401610306565b80516001600160a01b039081166000908152600660209081526040808320805460ff19908116909155855190941683526007909152808220805490931660011790925590513390670de0b6b3a7640000908381818185875af1925050503d80600081146107c2576040519150601f19603f3d011682016040523d82523d6000602084013e6107c7565b606091505b50509050806107d557600080fd5b5050505050565b600554341461085e5760405162461bcd60e51b815260206004820152604260248201527f4d6f636b526f6c6c75702e726567697374657256616c696461746f723a204d7560448201527f73742073656e64207374616b6520726571756972656420746f2072656769737460648201526132b960f11b608482015260a401610306565b3360009081526006602052604090205460ff16156108e45760405162461bcd60e51b815260206004820152603a60248201527f4d6f636b526f6c6c75702e726567697374657256616c696461746f723a20566160448201527f6c696461746f7220616c726561647920726567697374657265640000000000006064820152608401610306565b3360009081526007602052604090205460ff16156109605760405162461bcd60e51b815260206004820152603360248201527f4d6f636b526f6c6c75702e726567697374657256616c696461746f723a2056616044820152721b1a59185d1bdc88189b1858dadb1a5cdd1959606a1b6064820152608401610306565b336000908152600660205260409020805460ff19166001179055565b6000806109b36109ae604080518082018252600080825260209182015281518083019092526001825260029082015290565b6109f9565b90506109ee6109cc6109c5838a610ab8565b8790610b4f565b846109e16109da858b610ab8565b8890610b4f565b6109e9610be3565b610ca3565b979650505050505050565b60408051808201909152600080825260208201528151158015610a1e57506020820151155b15610a3c575050604080518082019091526000808252602082015290565b6040518060400160405280836000015181526020017f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478460200151610a819190611856565b610aab907f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4761188e565b905292915050565b919050565b6040805180820190915260008082526020820152610ad4610f12565b835181526020808501519082015260408082018490526000908360608460076107d05a03fa9050808015610b0757610b09565bfe5b5080610b475760405162461bcd60e51b815260206004820152600d60248201526c1958cb5b5d5b0b59985a5b1959609a1b6044820152606401610306565b505092915050565b6040805180820190915260008082526020820152610b6b610f30565b835181526020808501518183015283516040808401919091529084015160608301526000908360808460066107d05a03fa9050808015610b07575080610b475760405162461bcd60e51b815260206004820152600d60248201526c1958cb5859190b59985a5b1959609a1b6044820152606401610306565b610beb610f4e565b50604080516080810182527f198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c28183019081527f1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed6060830152815281518083019092527f275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec82527f1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d60208381019190915281019190915290565b604080518082018252858152602080820185905282518084019093528583528201839052600091610cd2610f73565b60005b6002811015610e97576000610ceb8260066118bb565b9050848260028110610cff57610cff6118a5565b60200201515183610d118360006118da565b600c8110610d2157610d216118a5565b6020020152848260028110610d3857610d386118a5565b60200201516020015183826001610d4f91906118da565b600c8110610d5f57610d5f6118a5565b6020020152838260028110610d7657610d766118a5565b6020020151515183610d898360026118da565b600c8110610d9957610d996118a5565b6020020152838260028110610db057610db06118a5565b6020020151516001602002015183610dc98360036118da565b600c8110610dd957610dd96118a5565b6020020152838260028110610df057610df06118a5565b602002015160200151600060028110610e0b57610e0b6118a5565b602002015183610e1c8360046118da565b600c8110610e2c57610e2c6118a5565b6020020152838260028110610e4357610e436118a5565b602002015160200151600160028110610e5e57610e5e6118a5565b602002015183610e6f8360056118da565b600c8110610e7f57610e7f6118a5565b60200201525080610e8f816118f2565b915050610cd5565b50610ea0610f92565b60006020826101808560086107d05a03fa9050808015610b07575080610f005760405162461bcd60e51b81526020600482015260156024820152741c185a5c9a5b99cb5bdc18dbd9194b59985a5b1959605a1b6044820152606401610306565b5051151593505050505b949350505050565b60405180606001604052806003906020820280368337509192915050565b60405180608001604052806004906020820280368337509192915050565b6040518060400160405280610f61610fb0565b8152602001610f6e610fb0565b905290565b604051806101800160405280600c906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b60405180604001604052806002906020820280368337509192915050565b634e487b7160e01b600052604160045260246000fd5b6040516080810167ffffffffffffffff8111828210171561100757611007610fce565b60405290565b60405160a0810167ffffffffffffffff8111828210171561100757611007610fce565b6040516060810167ffffffffffffffff8111828210171561100757611007610fce565b6040805190810167ffffffffffffffff8111828210171561100757611007610fce565b604051601f8201601f1916810167ffffffffffffffff8111828210171561109f5761109f610fce565b604052919050565b803563ffffffff81168114610ab357600080fd5b803560ff81168114610ab357600080fd5b600082601f8301126110dd57600080fd5b813567ffffffffffffffff8111156110f7576110f7610fce565b61110a601f8201601f1916602001611076565b81815284602083860101111561111f57600080fd5b816020850160208301376000918101602001919091529392505050565b80356bffffffffffffffffffffffff81168114610ab357600080fd5b60006080828403121561116a57600080fd5b611172610fe4565b9050813567ffffffffffffffff8082111561118c57600080fd5b90830190608082860312156111a057600080fd5b6111a8610fe4565b823581526020830135828111156111be57600080fd5b6111ca878286016110cc565b6020830152506040830135828111156111e257600080fd5b6111ee878286016110cc565b604083015250611200606084016110a7565b6060820152835250506020828101359082015261121f6040830161113c565b6040820152611230606083016110a7565b606082015292915050565b600060a0828403121561124d57600080fd5b61125561100d565b9050611260826110a7565b815261126e602083016110bb565b6020820152604082013567ffffffffffffffff8082111561128e57600080fd5b61129a85838601611158565b604084015260608401359150808211156112b357600080fd5b6112bf858386016110cc565b606084015260808401359150808211156112d857600080fd5b506112e5848285016110cc565b60808301525092915050565b600080604080848603121561130557600080fd5b833567ffffffffffffffff8082111561131d57600080fd5b9085019081870360808082121561133357600080fd5b61133b611030565b8583121561134857600080fd5b611350611053565b9250843583526020808601358185015283825261136e8787016110a7565b8183015260609350838601358581111561138757600080fd5b8087019650508a601f87011261139c57600080fd5b8535858111156113ae576113ae610fce565b6113bc828260051b01611076565b81815260079190911b8701820190828101908d8311156113db57600080fd5b978301975b828910156114475785898f0312156113f85760008081fd5b611400610fe4565b6114098a6110bb565b8152611416858b016110bb565b858201526114258b8b016110bb565b8b820152611434888b016110a7565b81890152825297850197908301906113e0565b9884019890985250909750880135945050508083111561146657600080fd5b50506114748582860161123b565b9150509250929050565b60006020828403121561149057600080fd5b5035919050565b6001600160a01b038416815263ffffffff8316602082015260808101610f0a604083018480518252602090810151910152565b600082601f8301126114db57600080fd5b6114e3611053565b8060408401858111156114f557600080fd5b845b8181101561150f5780358452602093840193016114f7565b509095945050505050565b600080600083850360c081121561153057600080fd5b84359350602085013592506080603f198201121561154d57600080fd5b50611556611053565b61156386604087016114ca565b815261157286608087016114ca565b6020820152809150509250925092565b60006020828403121561159457600080fd5b81356001600160a01b03811681146115ab57600080fd5b9392505050565b6020808252825182820181905260009190848201906040850190845b818110156116245761161183855160ff815116825260ff602082015116602083015260ff604082015116604083015263ffffffff60608201511660608301525050565b92840192608092909201916001016115ce565b50909695505050505050565b6000815180845260005b818110156116565760208185018101518683018201520161163a565b81811115611668576000602083870101525b50601f01601f19169290920160200192915050565b600063ffffffff80835116845260ff6020840151166020850152604083015160a060408601528051608060a08701528051610120870152602081015160806101408801526116cf6101a0880182611630565b9050604082015161011f19888303016101608901526116ee8282611630565b60608401519095166101808901525050602082015160c087015260408201519261172860e08801856bffffffffffffffffffffffff169052565b606083015163ffffffff811661010089015293506060860151935086810360608801526117558185611630565b9350505050608083015184820360808601526117718282611630565b95945050505050565b60608152600060e0820161179c60608401875180518252602090810151910152565b60208681015163ffffffff1660a08501526040870151608060c0860181905281519384905290820192600091906101008701905b808410156118275761181382875160ff815116825260ff602082015116602083015260ff604082015116604083015263ffffffff60608201511660608301525050565b9484019460019390930192908201906117d0565b506001600160a01b038916878501528681036040880152611848818961167d565b9a9950505050505050505050565b60008261187357634e487b7160e01b600052601260045260246000fd5b500690565b634e487b7160e01b600052601160045260246000fd5b6000828210156118a0576118a0611878565b500390565b634e487b7160e01b600052603260045260246000fd5b60008160001904831182151516156118d5576118d5611878565b500290565b600082198211156118ed576118ed611878565b500190565b600060001982141561190657611906611878565b506001019056fea264697066735822122035f24f3d0500638f1c05f39b15cb1c1e39c17f0600a5da5cd2053663477694b264736f6c634300080c0033",
}

// ContractMockRollupABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractMockRollupMetaData.ABI instead.
var ContractMockRollupABI = ContractMockRollupMetaData.ABI

// ContractMockRollupBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractMockRollupMetaData.Bin instead.
var ContractMockRollupBin = ContractMockRollupMetaData.Bin

// DeployContractMockRollup deploys a new Ethereum contract, binding an instance of ContractMockRollup to it.
func DeployContractMockRollup(auth *bind.TransactOpts, backend bind.ContractBackend, _eigenDAServiceManager common.Address, _tau BN254G1Point, _illegalValue *big.Int, _quorumBlobParamsHash [32]byte, _stakeRequired *big.Int) (common.Address, *types.Transaction, *ContractMockRollup, error) {
	parsed, err := ContractMockRollupMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractMockRollupBin), backend, _eigenDAServiceManager, _tau, _illegalValue, _quorumBlobParamsHash, _stakeRequired)
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

// QuorumBlobParamsHash is a free data retrieval call binding the contract method 0x37507265.
//
// Solidity: function quorumBlobParamsHash() view returns(bytes32)
func (_ContractMockRollup *ContractMockRollupCaller) QuorumBlobParamsHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ContractMockRollup.contract.Call(opts, &out, "quorumBlobParamsHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// QuorumBlobParamsHash is a free data retrieval call binding the contract method 0x37507265.
//
// Solidity: function quorumBlobParamsHash() view returns(bytes32)
func (_ContractMockRollup *ContractMockRollupSession) QuorumBlobParamsHash() ([32]byte, error) {
	return _ContractMockRollup.Contract.QuorumBlobParamsHash(&_ContractMockRollup.CallOpts)
}

// QuorumBlobParamsHash is a free data retrieval call binding the contract method 0x37507265.
//
// Solidity: function quorumBlobParamsHash() view returns(bytes32)
func (_ContractMockRollup *ContractMockRollupCallerSession) QuorumBlobParamsHash() ([32]byte, error) {
	return _ContractMockRollup.Contract.QuorumBlobParamsHash(&_ContractMockRollup.CallOpts)
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

// PostCommitment is a paid mutator transaction binding the contract method 0x21553525.
//
// Solidity: function postCommitment(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint8,((bytes32,bytes,bytes,uint32),bytes32,uint96,uint32),bytes,bytes) blobVerificationProof) returns()
func (_ContractMockRollup *ContractMockRollupTransactor) PostCommitment(opts *bind.TransactOpts, blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof EigenDABlobUtilsBlobVerificationProof) (*types.Transaction, error) {
	return _ContractMockRollup.contract.Transact(opts, "postCommitment", blobHeader, blobVerificationProof)
}

// PostCommitment is a paid mutator transaction binding the contract method 0x21553525.
//
// Solidity: function postCommitment(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint8,((bytes32,bytes,bytes,uint32),bytes32,uint96,uint32),bytes,bytes) blobVerificationProof) returns()
func (_ContractMockRollup *ContractMockRollupSession) PostCommitment(blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof EigenDABlobUtilsBlobVerificationProof) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.PostCommitment(&_ContractMockRollup.TransactOpts, blobHeader, blobVerificationProof)
}

// PostCommitment is a paid mutator transaction binding the contract method 0x21553525.
//
// Solidity: function postCommitment(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint8,((bytes32,bytes,bytes,uint32),bytes32,uint96,uint32),bytes,bytes) blobVerificationProof) returns()
func (_ContractMockRollup *ContractMockRollupTransactorSession) PostCommitment(blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof EigenDABlobUtilsBlobVerificationProof) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.PostCommitment(&_ContractMockRollup.TransactOpts, blobHeader, blobVerificationProof)
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
