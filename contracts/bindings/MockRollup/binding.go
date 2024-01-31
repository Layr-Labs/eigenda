// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractmockrollup

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

// EigenDARollupUtilsBlobVerificationProof is an auto generated low-level Go binding around an user-defined struct.
type EigenDARollupUtilsBlobVerificationProof struct {
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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_eigenDAServiceManager\",\"type\":\"address\",\"internalType\":\"contractIEigenDAServiceManager\"},{\"name\":\"_tau\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"_illegalValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_quorumBlobParamsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_stakeRequired\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"blacklist\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"challengeCommitment\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"point\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"commitments\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"polynomialCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAServiceManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAServiceManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"illegalValue\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"postCommitment\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structIEigenDAServiceManager.QuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"quorumThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structEigenDARollupUtils.BlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumThresholdPercentages\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"fee\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumThresholdIndexes\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"quorumBlobParamsHash\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerValidator\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"stakeRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tau\",\"inputs\":[],\"outputs\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validators\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"}]",
	Bin: "0x608060405234801561001057600080fd5b50604051611ab8380380611ab883398101604081905261002f9161006c565b600080546001600160a01b0319166001600160a01b0396909616959095179094558251600155602090920151600255600355600455600555610115565b600080600080600085870360c081121561008557600080fd5b86516001600160a01b038116811461009c57600080fd5b95506040601f19820112156100b057600080fd5b50604080519081016001600160401b03811182821017156100e157634e487b7160e01b600052604160045260246000fd5b6040908152602088810151835290880151908201526060870151608088015160a09098015196999198509695945092505050565b611994806101246000396000f3fe60806040526004361061009c5760003560e01c80636281e63b116100645780636281e63b1461018f578063bcc6587f146101af578063cfc4af55146101b7578063f9f92be4146101e7578063fa52c7d814610227578063fc30cad01461025757600080fd5b806305010105146100a157806321553525146100ca57806337507265146100ec5780634440bc5c1461010257806349ce899714610118575b600080fd5b3480156100ad57600080fd5b506100b760055481565b6040519081526020015b60405180910390f35b3480156100d657600080fd5b506100ea6100e5366004611342565b61028f565b005b3480156100f857600080fd5b506100b760045481565b34801561010e57600080fd5b506100b760035481565b34801561012457600080fd5b506101806101333660046114cf565b6008602090815260009182526040918290208054835180850190945260018201548452600290910154918301919091526001600160a01b03811691600160a01b90910463ffffffff169083565b6040516100c1939291906114e8565b34801561019b57600080fd5b506100ea6101aa36600461156b565b610590565b6100ea61082d565b3480156101c357600080fd5b506001546002546101d2919082565b604080519283526020830191909152016100c1565b3480156101f357600080fd5b506102176102023660046115d3565b60076020526000908152604090205460ff1681565b60405190151581526020016100c1565b34801561023357600080fd5b506102176102423660046115d3565b60066020526000908152604090205460ff1681565b34801561026357600080fd5b50600054610277906001600160a01b031681565b6040516001600160a01b0390911681526020016100c1565b3360009081526006602052604090205460ff1661030f5760405162461bcd60e51b815260206004820152603360248201527f4d6f636b526f6c6c75702e706f7374436f6d6d69746d656e743a2056616c6964604482015272185d1bdc881b9bdd081c9959da5cdd195c9959606a1b60648201526084015b60405180910390fd5b426000908152600860205260409020546001600160a01b0316156103925760405162461bcd60e51b815260206004820152603460248201527f4d6f636b526f6c6c75702e706f7374436f6d6d69746d656e743a20436f6d6d696044820152731d1b595b9d08185b1c9958591e481c1bdcdd195960621b6064820152608401610306565b60005460405163219460e160e21b815273__$32f04d18c688c2c57b0347c20a77f3d6c9$__916386518384916103da9186916001600160a01b0390911690869060040161174d565b60006040518083038186803b1580156103f257600080fd5b505af4158015610406573d6000803e3d6000fd5b5050505060005b82604001515181101561045d5760008360400151828151811061043257610432611829565b602090810291909101015163ffffffff9091166060909101528061045581611855565b91505061040d565b5060045482604001516040516020016104769190611870565b60405160208183030381529060405280519060200120146105155760405162461bcd60e51b815260206004820152604d60248201527f4d6f636b526f6c6c75702e706f7374436f6d6d69746d656e743a2051756f727560448201527f6d426c6f62506172616d7320646f206e6f74206d617463682071756f72756d4260648201526c0d8dec4a0c2e4c2dae690c2e6d609b1b608482015260a401610306565b506040805160608101825233815260208381015163ffffffff90811682840190815294518385019081524260009081526008845294909420925183549551909116600160a01b026001600160c01b03199095166001600160a01b03919091161793909317815590518051600183015590910151600290910155565b600083815260086020908152604091829020825160608101845281546001600160a01b038082168352600160a01b90910463ffffffff1682850152845180860186526001840154815260029093015493830193909352928301528151166106575760405162461bcd60e51b815260206004820152603560248201527f4d6f636b526f6c6c75702e6368616c6c656e6765436f6d6d69746d656e743a2060448201527410dbdb5b5a5d1b595b9d081b9bdd081c1bdcdd1959605a1b6064820152608401610306565b806020015163ffffffff1683106106e25760405162461bcd60e51b815260206004820152604360248201527f4d6f636b526f6c6c75702e6368616c6c656e6765436f6d6d69746d656e743a2060448201527f506f696e74206d757374206265206c657373207468616e2064617461206c656e6064820152620cee8d60eb1b608482015260a401610306565b600354604080518082018252600154815260025460208201529083015161070d9286929091866109cd565b61078a5760405162461bcd60e51b815260206004820152604260248201527f4d6f636b526f6c6c75702e6368616c6c656e6765436f6d6d69746d656e743a2060448201527f446f6573206e6f74206576616c7561746520746f20696c6c6567616c2076616c606482015261756560f01b608482015260a401610306565b80516001600160a01b039081166000908152600660209081526040808320805460ff19908116909155855190941683526007909152808220805490931660011790925590513390670de0b6b3a7640000908381818185875af1925050503d8060008114610813576040519150601f19603f3d011682016040523d82523d6000602084013e610818565b606091505b505090508061082657600080fd5b5050505050565b60055434146108af5760405162461bcd60e51b815260206004820152604260248201527f4d6f636b526f6c6c75702e726567697374657256616c696461746f723a204d7560448201527f73742073656e64207374616b6520726571756972656420746f2072656769737460648201526132b960f11b608482015260a401610306565b3360009081526006602052604090205460ff16156109355760405162461bcd60e51b815260206004820152603a60248201527f4d6f636b526f6c6c75702e726567697374657256616c696461746f723a20566160448201527f6c696461746f7220616c726561647920726567697374657265640000000000006064820152608401610306565b3360009081526007602052604090205460ff16156109b15760405162461bcd60e51b815260206004820152603360248201527f4d6f636b526f6c6c75702e726567697374657256616c696461746f723a2056616044820152721b1a59185d1bdc88189b1858dadb1a5cdd1959606a1b6064820152608401610306565b336000908152600660205260409020805460ff19166001179055565b600080610a046109ff604080518082018252600080825260209182015281518083019092526001825260029082015290565b610a4a565b9050610a3f610a1d610a16838a610b09565b8790610ba0565b84610a32610a2b858b610b09565b8890610ba0565b610a3a610c34565b610cf4565b979650505050505050565b60408051808201909152600080825260208201528151158015610a6f57506020820151155b15610a8d575050604080518082019091526000808252602082015290565b6040518060400160405280836000015181526020017f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478460200151610ad291906118ee565b610afc907f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47611910565b905292915050565b919050565b6040805180820190915260008082526020820152610b25610f63565b835181526020808501519082015260408082018490526000908360608460076107d05a03fa9050808015610b5857610b5a565bfe5b5080610b985760405162461bcd60e51b815260206004820152600d60248201526c1958cb5b5d5b0b59985a5b1959609a1b6044820152606401610306565b505092915050565b6040805180820190915260008082526020820152610bbc610f81565b835181526020808501518183015283516040808401919091529084015160608301526000908360808460066107d05a03fa9050808015610b58575080610b985760405162461bcd60e51b815260206004820152600d60248201526c1958cb5859190b59985a5b1959609a1b6044820152606401610306565b610c3c610f9f565b50604080516080810182527f198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c28183019081527f1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed6060830152815281518083019092527f275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec82527f1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d60208381019190915281019190915290565b604080518082018252858152602080820185905282518084019093528583528201839052600091610d23610fc4565b60005b6002811015610ee8576000610d3c826006611927565b9050848260028110610d5057610d50611829565b60200201515183610d62836000611946565b600c8110610d7257610d72611829565b6020020152848260028110610d8957610d89611829565b60200201516020015183826001610da09190611946565b600c8110610db057610db0611829565b6020020152838260028110610dc757610dc7611829565b6020020151515183610dda836002611946565b600c8110610dea57610dea611829565b6020020152838260028110610e0157610e01611829565b6020020151516001602002015183610e1a836003611946565b600c8110610e2a57610e2a611829565b6020020152838260028110610e4157610e41611829565b602002015160200151600060028110610e5c57610e5c611829565b602002015183610e6d836004611946565b600c8110610e7d57610e7d611829565b6020020152838260028110610e9457610e94611829565b602002015160200151600160028110610eaf57610eaf611829565b602002015183610ec0836005611946565b600c8110610ed057610ed0611829565b60200201525080610ee081611855565b915050610d26565b50610ef1610fe3565b60006020826101808560086107d05a03fa9050808015610b58575080610f515760405162461bcd60e51b81526020600482015260156024820152741c185a5c9a5b99cb5bdc18dbd9194b59985a5b1959605a1b6044820152606401610306565b5051151593505050505b949350505050565b60405180606001604052806003906020820280368337509192915050565b60405180608001604052806004906020820280368337509192915050565b6040518060400160405280610fb2611001565b8152602001610fbf611001565b905290565b604051806101800160405280600c906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b60405180604001604052806002906020820280368337509192915050565b634e487b7160e01b600052604160045260246000fd5b6040516080810167ffffffffffffffff811182821017156110585761105861101f565b60405290565b60405160a0810167ffffffffffffffff811182821017156110585761105861101f565b6040516060810167ffffffffffffffff811182821017156110585761105861101f565b6040805190810167ffffffffffffffff811182821017156110585761105861101f565b604051601f8201601f1916810167ffffffffffffffff811182821017156110f0576110f061101f565b604052919050565b803563ffffffff81168114610b0457600080fd5b803560ff81168114610b0457600080fd5b600082601f83011261112e57600080fd5b813567ffffffffffffffff8111156111485761114861101f565b61115b601f8201601f19166020016110c7565b81815284602083860101111561117057600080fd5b816020850160208301376000918101602001919091529392505050565b80356bffffffffffffffffffffffff81168114610b0457600080fd5b6000608082840312156111bb57600080fd5b6111c3611035565b9050813567ffffffffffffffff808211156111dd57600080fd5b90830190608082860312156111f157600080fd5b6111f9611035565b8235815260208301358281111561120f57600080fd5b61121b8782860161111d565b60208301525060408301358281111561123357600080fd5b61123f8782860161111d565b604083015250611251606084016110f8565b606082015283525050602082810135908201526112706040830161118d565b6040820152611281606083016110f8565b606082015292915050565b600060a0828403121561129e57600080fd5b6112a661105e565b90506112b1826110f8565b81526112bf6020830161110c565b6020820152604082013567ffffffffffffffff808211156112df57600080fd5b6112eb858386016111a9565b6040840152606084013591508082111561130457600080fd5b6113108583860161111d565b6060840152608084013591508082111561132957600080fd5b506113368482850161111d565b60808301525092915050565b600080604080848603121561135657600080fd5b833567ffffffffffffffff8082111561136e57600080fd5b9085019081870360808082121561138457600080fd5b61138c611081565b8583121561139957600080fd5b6113a16110a4565b925084358352602080860135818501528382526113bf8787016110f8565b818301526060935083860135858111156113d857600080fd5b8087019650508a601f8701126113ed57600080fd5b8535858111156113ff576113ff61101f565b61140d828260051b016110c7565b81815260079190911b8701820190828101908d83111561142c57600080fd5b978301975b828910156114985785898f0312156114495760008081fd5b611451611035565b61145a8a61110c565b8152611467858b0161110c565b858201526114768b8b0161110c565b8b820152611485888b016110f8565b8189015282529785019790830190611431565b988401989098525090975088013594505050808311156114b757600080fd5b50506114c58582860161128c565b9150509250929050565b6000602082840312156114e157600080fd5b5035919050565b6001600160a01b038416815263ffffffff8316602082015260808101610f5b604083018480518252602090810151910152565b600082601f83011261152c57600080fd5b6115346110a4565b80604084018581111561154657600080fd5b845b81811015611560578035845260209384019301611548565b509095945050505050565b600080600083850360c081121561158157600080fd5b84359350602085013592506080603f198201121561159e57600080fd5b506115a76110a4565b6115b4866040870161151b565b81526115c3866080870161151b565b6020820152809150509250925092565b6000602082840312156115e557600080fd5b81356001600160a01b03811681146115fc57600080fd5b9392505050565b6000815180845260005b818110156116295760208185018101518683018201520161160d565b8181111561163b576000602083870101525b50601f01601f19169290920160200192915050565b600063ffffffff80835116845260ff6020840151166020850152604083015160a060408601528051608060a08701528051610120870152602081015160806101408801526116a26101a0880182611603565b9050604082015161011f19888303016101608901526116c18282611603565b60608401519095166101808901525050602082015160c08701526040820151926116fb60e08801856bffffffffffffffffffffffff169052565b606083015163ffffffff811661010089015293506060860151935086810360608801526117288185611603565b9350505050608083015184820360808601526117448282611603565b95945050505050565b60608152600060e0820161176f60608401875180518252602090810151910152565b60208681015163ffffffff1660a08501526040870151608060c0860181905281519384905290820192600091906101008701905b808410156117fa576117e682875160ff815116825260ff602082015116602083015260ff604082015116604083015263ffffffff60608201511660608301525050565b9484019460019390930192908201906117a3565b506001600160a01b03891687850152868103604088015261181b8189611650565b9a9950505050505050505050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b60006000198214156118695761186961183f565b5060010190565b6020808252825182820181905260009190848201906040850190845b818110156118e2576118cf83855160ff815116825260ff602082015116602083015260ff604082015116604083015263ffffffff60608201511660608301525050565b928401926080929092019160010161188c565b50909695505050505050565b60008261190b57634e487b7160e01b600052601260045260246000fd5b500690565b6000828210156119225761192261183f565b500390565b60008160001904831182151516156119415761194161183f565b500290565b600082198211156119595761195961183f565b50019056fea26469706673582212208bd6d2816a5065ec167489e5490238a9e4fb18d147b052e68558bcc97f7bef5764736f6c634300080c0033",
}

// ContractmockrollupABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractmockrollupMetaData.ABI instead.
var ContractmockrollupABI = ContractmockrollupMetaData.ABI

// Contractmockrollup is an auto generated Go binding around an Ethereum contract.
type Contractmockrollup struct {
	ContractmockrollupCaller     // Read-only binding to the contract
	ContractmockrollupTransactor // Write-only binding to the contract
	ContractmockrollupFilterer   // Log filterer for contract events
}

// ContractmockrollupCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractmockrollupCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractmockrollupTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractmockrollupTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractmockrollupFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractmockrollupFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractmockrollupSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractmockrollupSession struct {
	Contract     *Contractmockrollup // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ContractmockrollupCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractmockrollupCallerSession struct {
	Contract *ContractmockrollupCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// ContractmockrollupTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractmockrollupTransactorSession struct {
	Contract     *ContractmockrollupTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// ContractmockrollupRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractmockrollupRaw struct {
	Contract *Contractmockrollup // Generic contract binding to access the raw methods on
}

// ContractmockrollupCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractmockrollupCallerRaw struct {
	Contract *ContractmockrollupCaller // Generic read-only contract binding to access the raw methods on
}

// ContractmockrollupTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractmockrollupTransactorRaw struct {
	Contract *ContractmockrollupTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractmockrollup creates a new instance of Contractmockrollup, bound to a specific deployed contract.
func NewContractmockrollup(address common.Address, backend bind.ContractBackend) (*Contractmockrollup, error) {
	contract, err := bindContractmockrollup(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contractmockrollup{ContractmockrollupCaller: ContractmockrollupCaller{contract: contract}, ContractmockrollupTransactor: ContractmockrollupTransactor{contract: contract}, ContractmockrollupFilterer: ContractmockrollupFilterer{contract: contract}}, nil
}

// NewContractmockrollupCaller creates a new read-only instance of Contractmockrollup, bound to a specific deployed contract.
func NewContractmockrollupCaller(address common.Address, caller bind.ContractCaller) (*ContractmockrollupCaller, error) {
	contract, err := bindContractmockrollup(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractmockrollupCaller{contract: contract}, nil
}

// NewContractmockrollupTransactor creates a new write-only instance of Contractmockrollup, bound to a specific deployed contract.
func NewContractmockrollupTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractmockrollupTransactor, error) {
	contract, err := bindContractmockrollup(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractmockrollupTransactor{contract: contract}, nil
}

// NewContractmockrollupFilterer creates a new log filterer instance of Contractmockrollup, bound to a specific deployed contract.
func NewContractmockrollupFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractmockrollupFilterer, error) {
	contract, err := bindContractmockrollup(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractmockrollupFilterer{contract: contract}, nil
}

// bindContractmockrollup binds a generic wrapper to an already deployed contract.
func bindContractmockrollup(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractmockrollupMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contractmockrollup *ContractmockrollupRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contractmockrollup.Contract.ContractmockrollupCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contractmockrollup *ContractmockrollupRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contractmockrollup.Contract.ContractmockrollupTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contractmockrollup *ContractmockrollupRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contractmockrollup.Contract.ContractmockrollupTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contractmockrollup *ContractmockrollupCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contractmockrollup.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contractmockrollup *ContractmockrollupTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contractmockrollup.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contractmockrollup *ContractmockrollupTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contractmockrollup.Contract.contract.Transact(opts, method, params...)
}

// Blacklist is a free data retrieval call binding the contract method 0xf9f92be4.
//
// Solidity: function blacklist(address ) view returns(bool)
func (_Contractmockrollup *ContractmockrollupCaller) Blacklist(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _Contractmockrollup.contract.Call(opts, &out, "blacklist", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Blacklist is a free data retrieval call binding the contract method 0xf9f92be4.
//
// Solidity: function blacklist(address ) view returns(bool)
func (_Contractmockrollup *ContractmockrollupSession) Blacklist(arg0 common.Address) (bool, error) {
	return _Contractmockrollup.Contract.Blacklist(&_Contractmockrollup.CallOpts, arg0)
}

// Blacklist is a free data retrieval call binding the contract method 0xf9f92be4.
//
// Solidity: function blacklist(address ) view returns(bool)
func (_Contractmockrollup *ContractmockrollupCallerSession) Blacklist(arg0 common.Address) (bool, error) {
	return _Contractmockrollup.Contract.Blacklist(&_Contractmockrollup.CallOpts, arg0)
}

// Commitments is a free data retrieval call binding the contract method 0x49ce8997.
//
// Solidity: function commitments(uint256 ) view returns(address validator, uint32 dataLength, (uint256,uint256) polynomialCommitment)
func (_Contractmockrollup *ContractmockrollupCaller) Commitments(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Validator            common.Address
	DataLength           uint32
	PolynomialCommitment BN254G1Point
}, error) {
	var out []interface{}
	err := _Contractmockrollup.contract.Call(opts, &out, "commitments", arg0)

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
func (_Contractmockrollup *ContractmockrollupSession) Commitments(arg0 *big.Int) (struct {
	Validator            common.Address
	DataLength           uint32
	PolynomialCommitment BN254G1Point
}, error) {
	return _Contractmockrollup.Contract.Commitments(&_Contractmockrollup.CallOpts, arg0)
}

// Commitments is a free data retrieval call binding the contract method 0x49ce8997.
//
// Solidity: function commitments(uint256 ) view returns(address validator, uint32 dataLength, (uint256,uint256) polynomialCommitment)
func (_Contractmockrollup *ContractmockrollupCallerSession) Commitments(arg0 *big.Int) (struct {
	Validator            common.Address
	DataLength           uint32
	PolynomialCommitment BN254G1Point
}, error) {
	return _Contractmockrollup.Contract.Commitments(&_Contractmockrollup.CallOpts, arg0)
}

// EigenDAServiceManager is a free data retrieval call binding the contract method 0xfc30cad0.
//
// Solidity: function eigenDAServiceManager() view returns(address)
func (_Contractmockrollup *ContractmockrollupCaller) EigenDAServiceManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Contractmockrollup.contract.Call(opts, &out, "eigenDAServiceManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDAServiceManager is a free data retrieval call binding the contract method 0xfc30cad0.
//
// Solidity: function eigenDAServiceManager() view returns(address)
func (_Contractmockrollup *ContractmockrollupSession) EigenDAServiceManager() (common.Address, error) {
	return _Contractmockrollup.Contract.EigenDAServiceManager(&_Contractmockrollup.CallOpts)
}

// EigenDAServiceManager is a free data retrieval call binding the contract method 0xfc30cad0.
//
// Solidity: function eigenDAServiceManager() view returns(address)
func (_Contractmockrollup *ContractmockrollupCallerSession) EigenDAServiceManager() (common.Address, error) {
	return _Contractmockrollup.Contract.EigenDAServiceManager(&_Contractmockrollup.CallOpts)
}

// IllegalValue is a free data retrieval call binding the contract method 0x4440bc5c.
//
// Solidity: function illegalValue() view returns(uint256)
func (_Contractmockrollup *ContractmockrollupCaller) IllegalValue(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contractmockrollup.contract.Call(opts, &out, "illegalValue")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// IllegalValue is a free data retrieval call binding the contract method 0x4440bc5c.
//
// Solidity: function illegalValue() view returns(uint256)
func (_Contractmockrollup *ContractmockrollupSession) IllegalValue() (*big.Int, error) {
	return _Contractmockrollup.Contract.IllegalValue(&_Contractmockrollup.CallOpts)
}

// IllegalValue is a free data retrieval call binding the contract method 0x4440bc5c.
//
// Solidity: function illegalValue() view returns(uint256)
func (_Contractmockrollup *ContractmockrollupCallerSession) IllegalValue() (*big.Int, error) {
	return _Contractmockrollup.Contract.IllegalValue(&_Contractmockrollup.CallOpts)
}

// StakeRequired is a free data retrieval call binding the contract method 0x05010105.
//
// Solidity: function stakeRequired() view returns(uint256)
func (_Contractmockrollup *ContractmockrollupCaller) StakeRequired(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contractmockrollup.contract.Call(opts, &out, "stakeRequired")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StakeRequired is a free data retrieval call binding the contract method 0x05010105.
//
// Solidity: function stakeRequired() view returns(uint256)
func (_Contractmockrollup *ContractmockrollupSession) StakeRequired() (*big.Int, error) {
	return _Contractmockrollup.Contract.StakeRequired(&_Contractmockrollup.CallOpts)
}

// StakeRequired is a free data retrieval call binding the contract method 0x05010105.
//
// Solidity: function stakeRequired() view returns(uint256)
func (_Contractmockrollup *ContractmockrollupCallerSession) StakeRequired() (*big.Int, error) {
	return _Contractmockrollup.Contract.StakeRequired(&_Contractmockrollup.CallOpts)
}

// Tau is a free data retrieval call binding the contract method 0xcfc4af55.
//
// Solidity: function tau() view returns(uint256 X, uint256 Y)
func (_Contractmockrollup *ContractmockrollupCaller) Tau(opts *bind.CallOpts) (struct {
	X *big.Int
	Y *big.Int
}, error) {
	var out []interface{}
	err := _Contractmockrollup.contract.Call(opts, &out, "tau")

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
func (_Contractmockrollup *ContractmockrollupSession) Tau() (struct {
	X *big.Int
	Y *big.Int
}, error) {
	return _Contractmockrollup.Contract.Tau(&_Contractmockrollup.CallOpts)
}

// Tau is a free data retrieval call binding the contract method 0xcfc4af55.
//
// Solidity: function tau() view returns(uint256 X, uint256 Y)
func (_Contractmockrollup *ContractmockrollupCallerSession) Tau() (struct {
	X *big.Int
	Y *big.Int
}, error) {
	return _Contractmockrollup.Contract.Tau(&_Contractmockrollup.CallOpts)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) view returns(bool)
func (_Contractmockrollup *ContractmockrollupCaller) Validators(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _Contractmockrollup.contract.Call(opts, &out, "validators", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) view returns(bool)
func (_Contractmockrollup *ContractmockrollupSession) Validators(arg0 common.Address) (bool, error) {
	return _Contractmockrollup.Contract.Validators(&_Contractmockrollup.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) view returns(bool)
func (_Contractmockrollup *ContractmockrollupCallerSession) Validators(arg0 common.Address) (bool, error) {
	return _Contractmockrollup.Contract.Validators(&_Contractmockrollup.CallOpts, arg0)
}

// ChallengeCommitment is a paid mutator transaction binding the contract method 0x6281e63b.
//
// Solidity: function challengeCommitment(uint256 timestamp, uint256 point, (uint256[2],uint256[2]) proof) returns()
func (_Contractmockrollup *ContractmockrollupTransactor) ChallengeCommitment(opts *bind.TransactOpts, timestamp *big.Int, point *big.Int, proof BN254G2Point) (*types.Transaction, error) {
	return _Contractmockrollup.contract.Transact(opts, "challengeCommitment", timestamp, point, proof)
}

// ChallengeCommitment is a paid mutator transaction binding the contract method 0x6281e63b.
//
// Solidity: function challengeCommitment(uint256 timestamp, uint256 point, (uint256[2],uint256[2]) proof) returns()
func (_Contractmockrollup *ContractmockrollupSession) ChallengeCommitment(timestamp *big.Int, point *big.Int, proof BN254G2Point) (*types.Transaction, error) {
	return _Contractmockrollup.Contract.ChallengeCommitment(&_Contractmockrollup.TransactOpts, timestamp, point, proof)
}

// ChallengeCommitment is a paid mutator transaction binding the contract method 0x6281e63b.
//
// Solidity: function challengeCommitment(uint256 timestamp, uint256 point, (uint256[2],uint256[2]) proof) returns()
func (_Contractmockrollup *ContractmockrollupTransactorSession) ChallengeCommitment(timestamp *big.Int, point *big.Int, proof BN254G2Point) (*types.Transaction, error) {
	return _Contractmockrollup.Contract.ChallengeCommitment(&_Contractmockrollup.TransactOpts, timestamp, point, proof)
}

// PostCommitment is a paid mutator transaction binding the contract method 0x21553525.
//
// Solidity: function postCommitment(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint8,((bytes32,bytes,bytes,uint32),bytes32,uint96,uint32),bytes,bytes) blobVerificationProof) returns()
func (_ContractMockRollup *ContractMockRollupTransactor) PostCommitment(opts *bind.TransactOpts, blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof EigenDARollupUtilsBlobVerificationProof) (*types.Transaction, error) {
	return _ContractMockRollup.contract.Transact(opts, "postCommitment", blobHeader, blobVerificationProof)
}

// PostCommitment is a paid mutator transaction binding the contract method 0x21553525.
//
// Solidity: function postCommitment(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint8,((bytes32,bytes,bytes,uint32),bytes32,uint96,uint32),bytes,bytes) blobVerificationProof) returns()
func (_ContractMockRollup *ContractMockRollupSession) PostCommitment(blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof EigenDARollupUtilsBlobVerificationProof) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.PostCommitment(&_ContractMockRollup.TransactOpts, blobHeader, blobVerificationProof)
}

// PostCommitment is a paid mutator transaction binding the contract method 0x21553525.
//
// Solidity: function postCommitment(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint8,((bytes32,bytes,bytes,uint32),bytes32,uint96,uint32),bytes,bytes) blobVerificationProof) returns()
func (_ContractMockRollup *ContractMockRollupTransactorSession) PostCommitment(blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof EigenDARollupUtilsBlobVerificationProof) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.PostCommitment(&_ContractMockRollup.TransactOpts, blobHeader, blobVerificationProof)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0xbcc6587f.
//
// Solidity: function registerValidator() payable returns()
func (_Contractmockrollup *ContractmockrollupTransactor) RegisterValidator(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contractmockrollup.contract.Transact(opts, "registerValidator")
}

// RegisterValidator is a paid mutator transaction binding the contract method 0xbcc6587f.
//
// Solidity: function registerValidator() payable returns()
func (_Contractmockrollup *ContractmockrollupSession) RegisterValidator() (*types.Transaction, error) {
	return _Contractmockrollup.Contract.RegisterValidator(&_Contractmockrollup.TransactOpts)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0xbcc6587f.
//
// Solidity: function registerValidator() payable returns()
func (_Contractmockrollup *ContractmockrollupTransactorSession) RegisterValidator() (*types.Transaction, error) {
	return _Contractmockrollup.Contract.RegisterValidator(&_Contractmockrollup.TransactOpts)
}
