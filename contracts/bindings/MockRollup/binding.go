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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_eigenDAServiceManager\",\"type\":\"address\",\"internalType\":\"contractIEigenDAServiceManager\"},{\"name\":\"_tau\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"_illegalValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_stakeRequired\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"blacklist\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"challengeCommitment\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"point\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"commitments\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"polynomialCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAServiceManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAServiceManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"illegalValue\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"postCommitment\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structIEigenDAServiceManager.QuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"quorumThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structEigenDARollupUtils.BlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumThresholdPercentages\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"fee\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumThresholdIndexes\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerValidator\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"stakeRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tau\",\"inputs\":[],\"outputs\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validators\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"}]",
	Bin: "0x608060405234801561001057600080fd5b506040516118f23803806118f283398101604081905261002f91610069565b600080546001600160a01b0319166001600160a01b0395909516949094179093558151600155602090910151600255600355600455610107565b60008060008084860360a081121561008057600080fd5b85516001600160a01b038116811461009757600080fd5b94506040601f19820112156100ab57600080fd5b50604080519081016001600160401b03811182821017156100dc57634e487b7160e01b600052604160045260246000fd5b6040908152602087810151835290870151908201526060860151608090960151949790965092505050565b6117dc806101166000396000f3fe6080604052600436106100915760003560e01c8063bcc6587f11610059578063bcc6587f1461018e578063cfc4af5514610196578063f9f92be4146101c6578063fa52c7d814610206578063fc30cad01461023657600080fd5b8063050101051461009657806321553525146100bf5780634440bc5c146100e157806349ce8997146100f75780636281e63b1461016e575b600080fd5b3480156100a257600080fd5b506100ac60045481565b6040519081526020015b60405180910390f35b3480156100cb57600080fd5b506100df6100da366004611219565b61026e565b005b3480156100ed57600080fd5b506100ac60035481565b34801561010357600080fd5b5061015f6101123660046113a6565b6007602090815260009182526040918290208054835180850190945260018201548452600290910154918301919091526001600160a01b03811691600160a01b90910463ffffffff169083565b6040516100b6939291906113bf565b34801561017a57600080fd5b506100df610189366004611442565b610467565b6100df610704565b3480156101a257600080fd5b506001546002546101b1919082565b604080519283526020830191909152016100b6565b3480156101d257600080fd5b506101f66101e13660046114aa565b60066020526000908152604090205460ff1681565b60405190151581526020016100b6565b34801561021257600080fd5b506101f66102213660046114aa565b60056020526000908152604090205460ff1681565b34801561024257600080fd5b50600054610256906001600160a01b031681565b6040516001600160a01b0390911681526020016100b6565b3360009081526005602052604090205460ff166102ee5760405162461bcd60e51b815260206004820152603360248201527f4d6f636b526f6c6c75702e706f7374436f6d6d69746d656e743a2056616c6964604482015272185d1bdc881b9bdd081c9959da5cdd195c9959606a1b60648201526084015b60405180910390fd5b426000908152600760205260409020546001600160a01b0316156103715760405162461bcd60e51b815260206004820152603460248201527f4d6f636b526f6c6c75702e706f7374436f6d6d69746d656e743a20436f6d6d696044820152731d1b595b9d08185b1c9958591e481c1bdcdd195960621b60648201526084016102e5565b60005460405163219460e160e21b815273__$32f04d18c688c2c57b0347c20a77f3d6c9$__916386518384916103b99186916001600160a01b03909116908690600401611624565b60006040518083038186803b1580156103d157600080fd5b505af41580156103e5573d6000803e3d6000fd5b50506040805160608101825233815260208681015163ffffffff90811682840190815297518385019081524260009081526007845294909420925183549851909116600160a01b026001600160c01b03199098166001600160a01b03919091161796909617815590518051600183015590940151600290940193909355505050565b600083815260076020908152604091829020825160608101845281546001600160a01b038082168352600160a01b90910463ffffffff16828501528451808601865260018401548152600290930154938301939093529283015281511661052e5760405162461bcd60e51b815260206004820152603560248201527f4d6f636b526f6c6c75702e6368616c6c656e6765436f6d6d69746d656e743a2060448201527410dbdb5b5a5d1b595b9d081b9bdd081c1bdcdd1959605a1b60648201526084016102e5565b806020015163ffffffff1683106105b95760405162461bcd60e51b815260206004820152604360248201527f4d6f636b526f6c6c75702e6368616c6c656e6765436f6d6d69746d656e743a2060448201527f506f696e74206d757374206265206c657373207468616e2064617461206c656e6064820152620cee8d60eb1b608482015260a4016102e5565b60035460408051808201825260015481526002546020820152908301516105e49286929091866108a4565b6106615760405162461bcd60e51b815260206004820152604260248201527f4d6f636b526f6c6c75702e6368616c6c656e6765436f6d6d69746d656e743a2060448201527f446f6573206e6f74206576616c7561746520746f20696c6c6567616c2076616c606482015261756560f01b608482015260a4016102e5565b80516001600160a01b039081166000908152600560209081526040808320805460ff19908116909155855190941683526006909152808220805490931660011790925590513390670de0b6b3a7640000908381818185875af1925050503d80600081146106ea576040519150601f19603f3d011682016040523d82523d6000602084013e6106ef565b606091505b50509050806106fd57600080fd5b5050505050565b60045434146107865760405162461bcd60e51b815260206004820152604260248201527f4d6f636b526f6c6c75702e726567697374657256616c696461746f723a204d7560448201527f73742073656e64207374616b6520726571756972656420746f2072656769737460648201526132b960f11b608482015260a4016102e5565b3360009081526005602052604090205460ff161561080c5760405162461bcd60e51b815260206004820152603a60248201527f4d6f636b526f6c6c75702e726567697374657256616c696461746f723a20566160448201527f6c696461746f7220616c7265616479207265676973746572656400000000000060648201526084016102e5565b3360009081526006602052604090205460ff16156108885760405162461bcd60e51b815260206004820152603360248201527f4d6f636b526f6c6c75702e726567697374657256616c696461746f723a2056616044820152721b1a59185d1bdc88189b1858dadb1a5cdd1959606a1b60648201526084016102e5565b336000908152600560205260409020805460ff19166001179055565b6000806108db6108d6604080518082018252600080825260209182015281518083019092526001825260029082015290565b610921565b90506109166108f46108ed838a6109e0565b8790610a77565b84610909610902858b6109e0565b8890610a77565b610911610b0b565b610bcb565b979650505050505050565b6040805180820190915260008082526020820152815115801561094657506020820151155b15610964575050604080518082019091526000808252602082015290565b6040518060400160405280836000015181526020017f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4784602001516109a991906116ef565b6109d3907f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47611727565b905292915050565b919050565b60408051808201909152600080825260208201526109fc610e3a565b835181526020808501519082015260408082018490526000908360608460076107d05a03fa9050808015610a2f57610a31565bfe5b5080610a6f5760405162461bcd60e51b815260206004820152600d60248201526c1958cb5b5d5b0b59985a5b1959609a1b60448201526064016102e5565b505092915050565b6040805180820190915260008082526020820152610a93610e58565b835181526020808501518183015283516040808401919091529084015160608301526000908360808460066107d05a03fa9050808015610a2f575080610a6f5760405162461bcd60e51b815260206004820152600d60248201526c1958cb5859190b59985a5b1959609a1b60448201526064016102e5565b610b13610e76565b50604080516080810182527f198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c28183019081527f1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed6060830152815281518083019092527f275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec82527f1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d60208381019190915281019190915290565b604080518082018252858152602080820185905282518084019093528583528201839052600091610bfa610e9b565b60005b6002811015610dbf576000610c13826006611754565b9050848260028110610c2757610c2761173e565b60200201515183610c39836000611773565b600c8110610c4957610c4961173e565b6020020152848260028110610c6057610c6061173e565b60200201516020015183826001610c779190611773565b600c8110610c8757610c8761173e565b6020020152838260028110610c9e57610c9e61173e565b6020020151515183610cb1836002611773565b600c8110610cc157610cc161173e565b6020020152838260028110610cd857610cd861173e565b6020020151516001602002015183610cf1836003611773565b600c8110610d0157610d0161173e565b6020020152838260028110610d1857610d1861173e565b602002015160200151600060028110610d3357610d3361173e565b602002015183610d44836004611773565b600c8110610d5457610d5461173e565b6020020152838260028110610d6b57610d6b61173e565b602002015160200151600160028110610d8657610d8661173e565b602002015183610d97836005611773565b600c8110610da757610da761173e565b60200201525080610db78161178b565b915050610bfd565b50610dc8610eba565b60006020826101808560086107d05a03fa9050808015610a2f575080610e285760405162461bcd60e51b81526020600482015260156024820152741c185a5c9a5b99cb5bdc18dbd9194b59985a5b1959605a1b60448201526064016102e5565b5051151593505050505b949350505050565b60405180606001604052806003906020820280368337509192915050565b60405180608001604052806004906020820280368337509192915050565b6040518060400160405280610e89610ed8565b8152602001610e96610ed8565b905290565b604051806101800160405280600c906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b60405180604001604052806002906020820280368337509192915050565b634e487b7160e01b600052604160045260246000fd5b6040516080810167ffffffffffffffff81118282101715610f2f57610f2f610ef6565b60405290565b60405160a0810167ffffffffffffffff81118282101715610f2f57610f2f610ef6565b6040516060810167ffffffffffffffff81118282101715610f2f57610f2f610ef6565b6040805190810167ffffffffffffffff81118282101715610f2f57610f2f610ef6565b604051601f8201601f1916810167ffffffffffffffff81118282101715610fc757610fc7610ef6565b604052919050565b803563ffffffff811681146109db57600080fd5b803560ff811681146109db57600080fd5b600082601f83011261100557600080fd5b813567ffffffffffffffff81111561101f5761101f610ef6565b611032601f8201601f1916602001610f9e565b81815284602083860101111561104757600080fd5b816020850160208301376000918101602001919091529392505050565b80356bffffffffffffffffffffffff811681146109db57600080fd5b60006080828403121561109257600080fd5b61109a610f0c565b9050813567ffffffffffffffff808211156110b457600080fd5b90830190608082860312156110c857600080fd5b6110d0610f0c565b823581526020830135828111156110e657600080fd5b6110f287828601610ff4565b60208301525060408301358281111561110a57600080fd5b61111687828601610ff4565b60408301525061112860608401610fcf565b6060820152835250506020828101359082015261114760408301611064565b604082015261115860608301610fcf565b606082015292915050565b600060a0828403121561117557600080fd5b61117d610f35565b905061118882610fcf565b815261119660208301610fe3565b6020820152604082013567ffffffffffffffff808211156111b657600080fd5b6111c285838601611080565b604084015260608401359150808211156111db57600080fd5b6111e785838601610ff4565b6060840152608084013591508082111561120057600080fd5b5061120d84828501610ff4565b60808301525092915050565b600080604080848603121561122d57600080fd5b833567ffffffffffffffff8082111561124557600080fd5b9085019081870360808082121561125b57600080fd5b611263610f58565b8583121561127057600080fd5b611278610f7b565b92508435835260208086013581850152838252611296878701610fcf565b818301526060935083860135858111156112af57600080fd5b8087019650508a601f8701126112c457600080fd5b8535858111156112d6576112d6610ef6565b6112e4828260051b01610f9e565b81815260079190911b8701820190828101908d83111561130357600080fd5b978301975b8289101561136f5785898f0312156113205760008081fd5b611328610f0c565b6113318a610fe3565b815261133e858b01610fe3565b8582015261134d8b8b01610fe3565b8b82015261135c888b01610fcf565b8189015282529785019790830190611308565b9884019890985250909750880135945050508083111561138e57600080fd5b505061139c85828601611163565b9150509250929050565b6000602082840312156113b857600080fd5b5035919050565b6001600160a01b038416815263ffffffff8316602082015260808101610e32604083018480518252602090810151910152565b600082601f83011261140357600080fd5b61140b610f7b565b80604084018581111561141d57600080fd5b845b8181101561143757803584526020938401930161141f565b509095945050505050565b600080600083850360c081121561145857600080fd5b84359350602085013592506080603f198201121561147557600080fd5b5061147e610f7b565b61148b86604087016113f2565b815261149a86608087016113f2565b6020820152809150509250925092565b6000602082840312156114bc57600080fd5b81356001600160a01b03811681146114d357600080fd5b9392505050565b6000815180845260005b81811015611500576020818501810151868301820152016114e4565b81811115611512576000602083870101525b50601f01601f19169290920160200192915050565b600063ffffffff80835116845260ff6020840151166020850152604083015160a060408601528051608060a08701528051610120870152602081015160806101408801526115796101a08801826114da565b9050604082015161011f198883030161016089015261159882826114da565b60608401519095166101808901525050602082015160c08701526040820151926115d260e08801856bffffffffffffffffffffffff169052565b606083015163ffffffff811661010089015293506060860151935086810360608801526115ff81856114da565b93505050506080830151848203608086015261161b82826114da565b95945050505050565b6060808252845180518383015260200151608083015260009060e0830160208088015163ffffffff80821660a088015260409150818a015160808060c08a01528582518088526101008b0191508684019750600093505b808410156116be578751805160ff90811684528882015181168985015287820151168784015289015185168983015296860196600193909301929082019061167b565b506001600160a01b038c168a870152898103858b01526116de818c611527565b9d9c50505050505050505050505050565b60008261170c57634e487b7160e01b600052601260045260246000fd5b500690565b634e487b7160e01b600052601160045260246000fd5b60008282101561173957611739611711565b500390565b634e487b7160e01b600052603260045260246000fd5b600081600019048311821515161561176e5761176e611711565b500290565b6000821982111561178657611786611711565b500190565b600060001982141561179f5761179f611711565b506001019056fea2646970667358221220c8a548011aebeca51cd15c4692b010198baa95740777cc9e008c6b45bf03468864736f6c634300080c0033",
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
