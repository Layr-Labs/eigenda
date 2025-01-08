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

// BatchHeader is an auto generated low-level Go binding around an user-defined struct.
type BatchHeader struct {
	BlobHeadersRoot       [32]byte
	QuorumNumbers         []byte
	SignedStakeForQuorums []byte
	ReferenceBlockNumber  uint32
}

// BatchMetadata is an auto generated low-level Go binding around an user-defined struct.
type BatchMetadata struct {
	BatchHeader             BatchHeader
	SignatoryRecordHash     [32]byte
	ConfirmationBlockNumber uint32
}

// BlobHeader is an auto generated low-level Go binding around an user-defined struct.
type BlobHeader struct {
	Commitment       BN254G1Point
	DataLength       uint32
	QuorumBlobParams []QuorumBlobParam
}

// BlobVerificationProof is an auto generated low-level Go binding around an user-defined struct.
type BlobVerificationProof struct {
	BatchId        uint32
	BlobIndex      uint32
	BatchMetadata  BatchMetadata
	InclusionProof []byte
	QuorumIndices  []byte
}

// QuorumBlobParam is an auto generated low-level Go binding around an user-defined struct.
type QuorumBlobParam struct {
	QuorumNumber                    uint8
	AdversaryThresholdPercentage    uint8
	ConfirmationThresholdPercentage uint8
	ChunkLength                     uint32
}

// ContractMockRollupMetaData contains all meta data concerning the ContractMockRollup contract.
var ContractMockRollupMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_eigenDAServiceManager\",\"type\":\"address\",\"internalType\":\"contractIEigenDAServiceManager\"},{\"name\":\"_tau\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"challengeCommitment\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"point\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"challengeValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"commitments\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"confirmer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"polynomialCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAServiceManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAServiceManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"postCommitment\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"tau\",\"inputs\":[],\"outputs\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162001e4938038062001e49833981016040819052620000349162000067565b600080546001600160a01b0319166001600160a01b039390931692909217909155805160015560200151600255620000f9565b60008082840360608112156200007c57600080fd5b83516001600160a01b03811681146200009457600080fd5b92506040601f1982011215620000a957600080fd5b50604080519081016001600160401b0381118282101715620000db57634e487b7160e01b600052604160045260246000fd5b60409081526020858101518352940151938101939093525092909150565b611d4080620001096000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c806349ce89971461005c578063b5144c73146100cf578063cfc4af55146100e4578063d2d16eb214610107578063fc30cad01461012a575b600080fd5b6100b761006a3660046114cb565b6003602090815260009182526040918290208054835180850190945260018201548452600290910154918301919091526001600160a01b03811691600160a01b90910463ffffffff169083565b6040516100c6939291906114e4565b60405180910390f35b6100e26100dd36600461181b565b610155565b005b6001546002546100f2919082565b604080519283526020830191909152016100c6565b61011a6101153660046119f8565b6101e8565b60405190151581526020016100c6565b60005461013d906001600160a01b031681565b6040516001600160a01b0390911681526020016100c6565b60005461016d9083906001600160a01b031683610375565b506040805160608101825233815260208381015163ffffffff90811682840190815294518385019081524260009081526003845294909420925183549551909116600160a01b026001600160c01b03199095166001600160a01b03919091161793909317815590518051600183015590910151600290910155565b6000848152600360209081526040808320815160608101835281546001600160a01b038082168352600160a01b90910463ffffffff16828601528351808501855260018401548152600290930154948301949094529182015280519091166102b55760405162461bcd60e51b815260206004820152603560248201527f4d6f636b526f6c6c75702e6368616c6c656e6765436f6d6d69746d656e743a2060448201527410dbdb5b5a5d1b595b9d081b9bdd081c1bdcdd1959605a1b60648201526084015b60405180910390fd5b806020015163ffffffff1685106103405760405162461bcd60e51b815260206004820152604360248201527f4d6f636b526f6c6c75702e6368616c6c656e6765436f6d6d69746d656e743a2060448201527f506f696e74206d757374206265206c657373207468616e2064617461206c656e6064820152620cee8d60eb1b608482015260a4016102ac565b604080518082018252600154815260025460208201529082015161036991879186919088610a35565b9150505b949350505050565b805160405163eccbbfc960e01b815263ffffffff90911660048201526001600160a01b0383169063eccbbfc990602401602060405180830381865afa1580156103c2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103e69190611a66565b6103f38260400151610ab2565b1461047a5760405162461bcd60e51b815260206004820152604b60248201527f456967656e4441526f6c6c75705574696c732e766572696679426c6f623a206260448201527f617463684d6574616461746120646f6573206e6f74206d617463682073746f7260648201526a6564206d6574616461746160a81b608482015260a4016102ac565b6060810151604082015151516104cc919061049486610b29565b6040516020016104a691815260200190565b60405160208183030381529060405280519060200120846020015163ffffffff16610b59565b61053e5760405162461bcd60e51b815260206004820152603960248201527f456967656e4441526f6c6c75705574696c732e766572696679426c6f623a206960448201527f6e636c7573696f6e2070726f6f6620697320696e76616c69640000000000000060648201526084016102ac565b6000805b84604001515181101561091e578460400151818151811061056557610565611a7f565b60200260200101516000015160ff16836040015160000151602001518460800151838151811061059757610597611a7f565b0160200151815160f89190911c9081106105b3576105b3611a7f565b016020015160f81c1461062e5760405162461bcd60e51b815260206004820152603a60248201527f456967656e4441526f6c6c75705574696c732e766572696679426c6f623a207160448201527f756f72756d4e756d62657220646f6573206e6f74206d6174636800000000000060648201526084016102ac565b8460400151818151811061064457610644611a7f565b60200260200101516040015160ff168560400151828151811061066957610669611a7f565b60200260200101516020015160ff16106106fc5760405162461bcd60e51b815260206004820152604860248201527f456967656e4441526f6c6c75705574696c732e766572696679426c6f623a206160448201527f64766572736172795468726573686f6c6450657263656e74616765206973206e6064820152671bdd081d985b1a5960c21b608482015260a4016102ac565b600061072c858760400151848151811061071857610718611a7f565b60200260200101516000015160ff16610b71565b905060ff8116156107e3578060ff168660400151838151811061075157610751611a7f565b60200260200101516020015160ff1610156107e35760405162461bcd60e51b815260206004820152604660248201527f456967656e4441526f6c6c75705574696c732e766572696679426c6f623a206160448201527f64766572736172795468726573686f6c6450657263656e74616765206973206e6064820152651bdd081b595d60d21b608482015260a4016102ac565b856040015182815181106107f9576107f9611a7f565b60200260200101516040015160ff16846040015160000151604001518560800151848151811061082b5761082b611a7f565b0160200151815160f89190911c90811061084757610847611a7f565b016020015160f81c10156108d55760405162461bcd60e51b815260206004820152604960248201527f456967656e4441526f6c6c75705574696c732e766572696679426c6f623a206360448201527f6f6e6669726d6174696f6e5468726573686f6c6450657263656e7461676520696064820152681cc81b9bdd081b595d60ba1b608482015260a4016102ac565b61090883876040015184815181106108ef576108ef611a7f565b602002602001015160000151600160ff919091161b1790565b925050808061091690611aab565b915050610542565b50610997610990846001600160a01b031663e15234ff6040518163ffffffff1660e01b8152600401600060405180830381865afa158015610963573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261098b9190810190611af2565b610c68565b8281161490565b610a2f5760405162461bcd60e51b815260206004820152605960248201527f456967656e4441526f6c6c75705574696c732e766572696679426c6f623a207260448201527f657175697265642071756f72756d7320617265206e6f7420612073756273657460648201527f206f662074686520636f6e6669726d65642071756f72756d7300000000000000608482015260a4016102ac565b50505050565b600080610a6c610a67604080518082018252600080825260209182015281518083019092526001825260029082015290565b610df5565b9050610aa7610a85610a7e838a610eb4565b8790610f4b565b84610a9a610a93858b610eb4565b8890610f4b565b610aa2610fdf565b61109f565b979650505050505050565b6000610b238260000151604051602001610acc9190611b95565b60408051808303601f1901815282825280516020918201208682015187840151838601929092528484015260e01b6001600160e01b0319166060840152815160448185030181526064909301909152815191012090565b92915050565b600081604051602001610b3c9190611bf5565b604051602081830303815290604052805190602001209050919050565b600083610b6786858561130c565b1495945050505050565b600081836001600160a01b0316638687feae6040518163ffffffff1660e01b8152600401600060405180830381865afa158015610bb2573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610bda9190810190611af2565b511115610b2357826001600160a01b0316638687feae6040518163ffffffff1660e01b8152600401600060405180830381865afa158015610c1f573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610c479190810190611af2565b8281518110610c5857610c58611a7f565b016020015160f81c905092915050565b600061010082511115610cf15760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a4016102ac565b8151610cff57506000919050565b60008083600081518110610d1557610d15611a7f565b0160200151600160f89190911c81901b92505b8451811015610dec57848181518110610d4357610d43611a7f565b0160200151600160f89190911c1b9150828211610dd85760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a4016102ac565b91811791610de581611aab565b9050610d28565b50909392505050565b60408051808201909152600080825260208201528151158015610e1a57506020820151155b15610e38575050604080518082019091526000808252602082015290565b6040518060400160405280836000015181526020017f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478460200151610e7d9190611c9a565b610ea7907f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47611cbc565b905292915050565b919050565b6040805180820190915260008082526020820152610ed061140f565b835181526020808501519082015260408082018490526000908360608460076107d05a03fa9050808015610f0357610f05565bfe5b5080610f435760405162461bcd60e51b815260206004820152600d60248201526c1958cb5b5d5b0b59985a5b1959609a1b60448201526064016102ac565b505092915050565b6040805180820190915260008082526020820152610f6761142d565b835181526020808501518183015283516040808401919091529084015160608301526000908360808460066107d05a03fa9050808015610f03575080610f435760405162461bcd60e51b815260206004820152600d60248201526c1958cb5859190b59985a5b1959609a1b60448201526064016102ac565b610fe761144b565b50604080516080810182527f198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c28183019081527f1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed6060830152815281518083019092527f275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec82527f1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d60208381019190915281019190915290565b6040805180820182528581526020808201859052825180840190935285835282018390526000916110ce611470565b60005b60028110156112935760006110e7826006611cd3565b90508482600281106110fb576110fb611a7f565b6020020151518361110d836000611cf2565b600c811061111d5761111d611a7f565b602002015284826002811061113457611134611a7f565b6020020151602001518382600161114b9190611cf2565b600c811061115b5761115b611a7f565b602002015283826002811061117257611172611a7f565b6020020151515183611185836002611cf2565b600c811061119557611195611a7f565b60200201528382600281106111ac576111ac611a7f565b60200201515160016020020151836111c5836003611cf2565b600c81106111d5576111d5611a7f565b60200201528382600281106111ec576111ec611a7f565b60200201516020015160006002811061120757611207611a7f565b602002015183611218836004611cf2565b600c811061122857611228611a7f565b602002015283826002811061123f5761123f611a7f565b60200201516020015160016002811061125a5761125a611a7f565b60200201518361126b836005611cf2565b600c811061127b5761127b611a7f565b6020020152508061128b81611aab565b9150506110d1565b5061129c61148f565b60006020826101808560086107d05a03fa9050808015610f035750806112fc5760405162461bcd60e51b81526020600482015260156024820152741c185a5c9a5b99cb5bdc18dbd9194b59985a5b1959605a1b60448201526064016102ac565b5051151598975050505050505050565b60006020845161131c9190611c9a565b156113a35760405162461bcd60e51b815260206004820152604b60248201527f4d65726b6c652e70726f63657373496e636c7573696f6e50726f6f664b65636360448201527f616b3a2070726f6f66206c656e6774682073686f756c642062652061206d756c60648201526a3a34b836329037b310199960a91b608482015260a4016102ac565b8260205b85518111611406576113ba600285611c9a565b6113db578160005280860151602052604060002091506002840493506113f4565b8086015160005281602052604060002091506002840493505b6113ff602082611cf2565b90506113a7565b50949350505050565b60405180606001604052806003906020820280368337509192915050565b60405180608001604052806004906020820280368337509192915050565b604051806040016040528061145e6114ad565b815260200161146b6114ad565b905290565b604051806101800160405280600c906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b60405180604001604052806002906020820280368337509192915050565b6000602082840312156114dd57600080fd5b5035919050565b6001600160a01b038416815263ffffffff831660208201526080810161036d604083018480518252602090810151910152565b634e487b7160e01b600052604160045260246000fd5b6040516060810167ffffffffffffffff8111828210171561155057611550611517565b60405290565b6040516080810167ffffffffffffffff8111828210171561155057611550611517565b60405160a0810167ffffffffffffffff8111828210171561155057611550611517565b6040805190810167ffffffffffffffff8111828210171561155057611550611517565b604051601f8201601f1916810167ffffffffffffffff811182821017156115e8576115e8611517565b604052919050565b803563ffffffff81168114610eaf57600080fd5b803560ff81168114610eaf57600080fd5b600067ffffffffffffffff82111561162f5761162f611517565b50601f01601f191660200190565b600082601f83011261164e57600080fd5b813561166161165c82611615565b6115bf565b81815284602083860101111561167657600080fd5b816020850160208301376000918101602001919091529392505050565b6000606082840312156116a557600080fd5b6116ad61152d565b9050813567ffffffffffffffff808211156116c757600080fd5b90830190608082860312156116db57600080fd5b6116e3611556565b823581526020830135828111156116f957600080fd5b6117058782860161163d565b60208301525060408301358281111561171d57600080fd5b6117298782860161163d565b60408301525061173b606084016115f0565b6060820152835250506020828101359082015261175a604083016115f0565b604082015292915050565b600060a0828403121561177757600080fd5b61177f611579565b905061178a826115f0565b8152611798602083016115f0565b6020820152604082013567ffffffffffffffff808211156117b857600080fd5b6117c485838601611693565b604084015260608401359150808211156117dd57600080fd5b6117e98583860161163d565b6060840152608084013591508082111561180257600080fd5b5061180f8482850161163d565b60808301525092915050565b600080604080848603121561182f57600080fd5b833567ffffffffffffffff8082111561184757600080fd5b9085019081870360808082121561185d57600080fd5b61186561152d565b8583121561187257600080fd5b61187a61159c565b925084358352602080860135818501528382526118988787016115f0565b818301526060935083860135858111156118b157600080fd5b8087019650508a601f8701126118c657600080fd5b8535858111156118d8576118d8611517565b6118e6828260051b016115bf565b81815260079190911b8701820190828101908d83111561190557600080fd5b978301975b828910156119715785898f0312156119225760008081fd5b61192a611556565b6119338a611604565b8152611940858b01611604565b8582015261194f8b8b01611604565b8b82015261195e888b016115f0565b818901528252978501979083019061190a565b9884019890985250909750880135945050508083111561199057600080fd5b505061199e85828601611765565b9150509250929050565b600082601f8301126119b957600080fd5b6119c161159c565b8060408401858111156119d357600080fd5b845b818110156119ed5780358452602093840193016119d5565b509095945050505050565b60008060008084860360e0811215611a0f57600080fd5b85359450602086013593506080603f1982011215611a2c57600080fd5b50611a3561159c565b611a4287604088016119a8565b8152611a5187608088016119a8565b60208201529396929550929360c00135925050565b600060208284031215611a7857600080fd5b5051919050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b6000600019821415611abf57611abf611a95565b5060010190565b60005b83811015611ae1578181015183820152602001611ac9565b83811115610a2f5750506000910152565b600060208284031215611b0457600080fd5b815167ffffffffffffffff811115611b1b57600080fd5b8201601f81018413611b2c57600080fd5b8051611b3a61165c82611615565b818152856020838501011115611b4f57600080fd5b611b60826020830160208601611ac6565b95945050505050565b60008151808452611b81816020860160208601611ac6565b601f01601f19169290920160200192915050565b60208152815160208201526000602083015160806040840152611bbb60a0840182611b69565b90506040840151601f19848303016060850152611bd88282611b69565b91505063ffffffff60608501511660808401528091505092915050565b60208082528251805183830152810151604083015260009060a0830181850151606063ffffffff808316828801526040925082880151608080818a015285825180885260c08b0191508884019750600093505b80841015611c8b578751805160ff90811684528a82015181168b850152888201511688840152860151851686830152968801966001939093019290820190611c48565b509a9950505050505050505050565b600082611cb757634e487b7160e01b600052601260045260246000fd5b500690565b600082821015611cce57611cce611a95565b500390565b6000816000190483118215151615611ced57611ced611a95565b500290565b60008219821115611d0557611d05611a95565b50019056fea26469706673582212200aeeed4f9b021817e07eb202ba618f7ffb4b9e2c5bc60727a154f432c8fc54d564736f6c634300080c0033",
}

// ContractMockRollupABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractMockRollupMetaData.ABI instead.
var ContractMockRollupABI = ContractMockRollupMetaData.ABI

// ContractMockRollupBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractMockRollupMetaData.Bin instead.
var ContractMockRollupBin = ContractMockRollupMetaData.Bin

// DeployContractMockRollup deploys a new Ethereum contract, binding an instance of ContractMockRollup to it.
func DeployContractMockRollup(auth *bind.TransactOpts, backend bind.ContractBackend, _eigenDAServiceManager common.Address, _tau BN254G1Point) (common.Address, *types.Transaction, *ContractMockRollup, error) {
	parsed, err := ContractMockRollupMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractMockRollupBin), backend, _eigenDAServiceManager, _tau)
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

// Commitments is a free data retrieval call binding the contract method 0x49ce8997.
//
// Solidity: function commitments(uint256 ) view returns(address confirmer, uint32 dataLength, (uint256,uint256) polynomialCommitment)
func (_ContractMockRollup *ContractMockRollupCaller) Commitments(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Confirmer            common.Address
	DataLength           uint32
	PolynomialCommitment BN254G1Point
}, error) {
	var out []interface{}
	err := _ContractMockRollup.contract.Call(opts, &out, "commitments", arg0)

	outstruct := new(struct {
		Confirmer            common.Address
		DataLength           uint32
		PolynomialCommitment BN254G1Point
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Confirmer = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.DataLength = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.PolynomialCommitment = *abi.ConvertType(out[2], new(BN254G1Point)).(*BN254G1Point)

	return *outstruct, err

}

// Commitments is a free data retrieval call binding the contract method 0x49ce8997.
//
// Solidity: function commitments(uint256 ) view returns(address confirmer, uint32 dataLength, (uint256,uint256) polynomialCommitment)
func (_ContractMockRollup *ContractMockRollupSession) Commitments(arg0 *big.Int) (struct {
	Confirmer            common.Address
	DataLength           uint32
	PolynomialCommitment BN254G1Point
}, error) {
	return _ContractMockRollup.Contract.Commitments(&_ContractMockRollup.CallOpts, arg0)
}

// Commitments is a free data retrieval call binding the contract method 0x49ce8997.
//
// Solidity: function commitments(uint256 ) view returns(address confirmer, uint32 dataLength, (uint256,uint256) polynomialCommitment)
func (_ContractMockRollup *ContractMockRollupCallerSession) Commitments(arg0 *big.Int) (struct {
	Confirmer            common.Address
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

// ChallengeCommitment is a paid mutator transaction binding the contract method 0xd2d16eb2.
//
// Solidity: function challengeCommitment(uint256 timestamp, uint256 point, (uint256[2],uint256[2]) proof, uint256 challengeValue) returns(bool)
func (_ContractMockRollup *ContractMockRollupTransactor) ChallengeCommitment(opts *bind.TransactOpts, timestamp *big.Int, point *big.Int, proof BN254G2Point, challengeValue *big.Int) (*types.Transaction, error) {
	return _ContractMockRollup.contract.Transact(opts, "challengeCommitment", timestamp, point, proof, challengeValue)
}

// ChallengeCommitment is a paid mutator transaction binding the contract method 0xd2d16eb2.
//
// Solidity: function challengeCommitment(uint256 timestamp, uint256 point, (uint256[2],uint256[2]) proof, uint256 challengeValue) returns(bool)
func (_ContractMockRollup *ContractMockRollupSession) ChallengeCommitment(timestamp *big.Int, point *big.Int, proof BN254G2Point, challengeValue *big.Int) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.ChallengeCommitment(&_ContractMockRollup.TransactOpts, timestamp, point, proof, challengeValue)
}

// ChallengeCommitment is a paid mutator transaction binding the contract method 0xd2d16eb2.
//
// Solidity: function challengeCommitment(uint256 timestamp, uint256 point, (uint256[2],uint256[2]) proof, uint256 challengeValue) returns(bool)
func (_ContractMockRollup *ContractMockRollupTransactorSession) ChallengeCommitment(timestamp *big.Int, point *big.Int, proof BN254G2Point, challengeValue *big.Int) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.ChallengeCommitment(&_ContractMockRollup.TransactOpts, timestamp, point, proof, challengeValue)
}

// PostCommitment is a paid mutator transaction binding the contract method 0xb5144c73.
//
// Solidity: function postCommitment(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) returns()
func (_ContractMockRollup *ContractMockRollupTransactor) PostCommitment(opts *bind.TransactOpts, blobHeader BlobHeader, blobVerificationProof BlobVerificationProof) (*types.Transaction, error) {
	return _ContractMockRollup.contract.Transact(opts, "postCommitment", blobHeader, blobVerificationProof)
}

// PostCommitment is a paid mutator transaction binding the contract method 0xb5144c73.
//
// Solidity: function postCommitment(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) returns()
func (_ContractMockRollup *ContractMockRollupSession) PostCommitment(blobHeader BlobHeader, blobVerificationProof BlobVerificationProof) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.PostCommitment(&_ContractMockRollup.TransactOpts, blobHeader, blobVerificationProof)
}

// PostCommitment is a paid mutator transaction binding the contract method 0xb5144c73.
//
// Solidity: function postCommitment(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) returns()
func (_ContractMockRollup *ContractMockRollupTransactorSession) PostCommitment(blobHeader BlobHeader, blobVerificationProof BlobVerificationProof) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.PostCommitment(&_ContractMockRollup.TransactOpts, blobHeader, blobVerificationProof)
}
