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
	Bin: "0x60806040523480156200001157600080fd5b5060405162001ecc38038062001ecc833981016040819052620000349162000067565b600080546001600160a01b0319166001600160a01b039390931692909217909155805160015560200151600255620000f9565b60008082840360608112156200007c57600080fd5b83516001600160a01b03811681146200009457600080fd5b92506040601f1982011215620000a957600080fd5b50604080519081016001600160401b0381118282101715620000db57634e487b7160e01b600052604160045260246000fd5b60409081526020858101518352940151938101939093525092909150565b611dc380620001096000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c806349ce89971461005c578063b5144c73146100cf578063cfc4af55146100e4578063d2d16eb214610107578063fc30cad01461012a575b600080fd5b6100b761006a36600461154e565b6003602090815260009182526040918290208054835180850190945260018201548452600290910154918301919091526001600160a01b03811691600160a01b90910463ffffffff169083565b6040516100c693929190611567565b60405180910390f35b6100e26100dd36600461189e565b610155565b005b6001546002546100f2919082565b604080519283526020830191909152016100c6565b61011a610115366004611a7b565b610270565b60405190151581526020016100c6565b60005461013d906001600160a01b031681565b6040516001600160a01b0390911681526020016100c6565b426000908152600360205260409020546001600160a01b0316156101dd5760405162461bcd60e51b815260206004820152603460248201527f4d6f636b526f6c6c75702e706f7374436f6d6d69746d656e743a20436f6d6d696044820152731d1b595b9d08185b1c9958591e481c1bdcdd195960621b60648201526084015b60405180910390fd5b6000546101f59083906001600160a01b0316836103f8565b506040805160608101825233815260208381015163ffffffff90811682840190815294518385019081524260009081526003845294909420925183549551909116600160a01b026001600160c01b03199095166001600160a01b03919091161793909317815590518051600183015590910151600290910155565b6000848152600360209081526040808320815160608101835281546001600160a01b038082168352600160a01b90910463ffffffff16828601528351808501855260018401548152600290930154948301949094529182015280519091166103385760405162461bcd60e51b815260206004820152603560248201527f4d6f636b526f6c6c75702e6368616c6c656e6765436f6d6d69746d656e743a2060448201527410dbdb5b5a5d1b595b9d081b9bdd081c1bdcdd1959605a1b60648201526084016101d4565b806020015163ffffffff1685106103c35760405162461bcd60e51b815260206004820152604360248201527f4d6f636b526f6c6c75702e6368616c6c656e6765436f6d6d69746d656e743a2060448201527f506f696e74206d757374206265206c657373207468616e2064617461206c656e6064820152620cee8d60eb1b608482015260a4016101d4565b60408051808201825260015481526002546020820152908201516103ec91879186919088610ab8565b9150505b949350505050565b805160405163eccbbfc960e01b815263ffffffff90911660048201526001600160a01b0383169063eccbbfc990602401602060405180830381865afa158015610445573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104699190611ae9565b6104768260400151610b35565b146104fd5760405162461bcd60e51b815260206004820152604b60248201527f456967656e4441526f6c6c75705574696c732e766572696679426c6f623a206260448201527f617463684d6574616461746120646f6573206e6f74206d617463682073746f7260648201526a6564206d6574616461746160a81b608482015260a4016101d4565b60608101516040820151515161054f919061051786610bac565b60405160200161052991815260200190565b60405160208183030381529060405280519060200120846020015163ffffffff16610bdc565b6105c15760405162461bcd60e51b815260206004820152603960248201527f456967656e4441526f6c6c75705574696c732e766572696679426c6f623a206960448201527f6e636c7573696f6e2070726f6f6620697320696e76616c69640000000000000060648201526084016101d4565b6000805b8460400151518110156109a157846040015181815181106105e8576105e8611b02565b60200260200101516000015160ff16836040015160000151602001518460800151838151811061061a5761061a611b02565b0160200151815160f89190911c90811061063657610636611b02565b016020015160f81c146106b15760405162461bcd60e51b815260206004820152603a60248201527f456967656e4441526f6c6c75705574696c732e766572696679426c6f623a207160448201527f756f72756d4e756d62657220646f6573206e6f74206d6174636800000000000060648201526084016101d4565b846040015181815181106106c7576106c7611b02565b60200260200101516040015160ff16856040015182815181106106ec576106ec611b02565b60200260200101516020015160ff161061077f5760405162461bcd60e51b815260206004820152604860248201527f456967656e4441526f6c6c75705574696c732e766572696679426c6f623a206160448201527f64766572736172795468726573686f6c6450657263656e74616765206973206e6064820152671bdd081d985b1a5960c21b608482015260a4016101d4565b60006107af858760400151848151811061079b5761079b611b02565b60200260200101516000015160ff16610bf4565b905060ff811615610866578060ff16866040015183815181106107d4576107d4611b02565b60200260200101516020015160ff1610156108665760405162461bcd60e51b815260206004820152604660248201527f456967656e4441526f6c6c75705574696c732e766572696679426c6f623a206160448201527f64766572736172795468726573686f6c6450657263656e74616765206973206e6064820152651bdd081b595d60d21b608482015260a4016101d4565b8560400151828151811061087c5761087c611b02565b60200260200101516040015160ff1684604001516000015160400151856080015184815181106108ae576108ae611b02565b0160200151815160f89190911c9081106108ca576108ca611b02565b016020015160f81c10156109585760405162461bcd60e51b815260206004820152604960248201527f456967656e4441526f6c6c75705574696c732e766572696679426c6f623a206360448201527f6f6e6669726d6174696f6e5468726573686f6c6450657263656e7461676520696064820152681cc81b9bdd081b595d60ba1b608482015260a4016101d4565b61098b838760400151848151811061097257610972611b02565b602002602001015160000151600160ff919091161b1790565b925050808061099990611b2e565b9150506105c5565b50610a1a610a13846001600160a01b031663e15234ff6040518163ffffffff1660e01b8152600401600060405180830381865afa1580156109e6573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610a0e9190810190611b75565b610ceb565b8281161490565b610ab25760405162461bcd60e51b815260206004820152605960248201527f456967656e4441526f6c6c75705574696c732e766572696679426c6f623a207260448201527f657175697265642071756f72756d7320617265206e6f7420612073756273657460648201527f206f662074686520636f6e6669726d65642071756f72756d7300000000000000608482015260a4016101d4565b50505050565b600080610aef610aea604080518082018252600080825260209182015281518083019092526001825260029082015290565b610e78565b9050610b2a610b08610b01838a610f37565b8790610fce565b84610b1d610b16858b610f37565b8890610fce565b610b25611062565b611122565b979650505050505050565b6000610ba68260000151604051602001610b4f9190611c18565b60408051808303601f1901815282825280516020918201208682015187840151838601929092528484015260e01b6001600160e01b0319166060840152815160448185030181526064909301909152815191012090565b92915050565b600081604051602001610bbf9190611c78565b604051602081830303815290604052805190602001209050919050565b600083610bea86858561138f565b1495945050505050565b600081836001600160a01b0316638687feae6040518163ffffffff1660e01b8152600401600060405180830381865afa158015610c35573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610c5d9190810190611b75565b511115610ba657826001600160a01b0316638687feae6040518163ffffffff1660e01b8152600401600060405180830381865afa158015610ca2573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610cca9190810190611b75565b8281518110610cdb57610cdb611b02565b016020015160f81c905092915050565b600061010082511115610d745760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a4016101d4565b8151610d8257506000919050565b60008083600081518110610d9857610d98611b02565b0160200151600160f89190911c81901b92505b8451811015610e6f57848181518110610dc657610dc6611b02565b0160200151600160f89190911c1b9150828211610e5b5760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a4016101d4565b91811791610e6881611b2e565b9050610dab565b50909392505050565b60408051808201909152600080825260208201528151158015610e9d57506020820151155b15610ebb575050604080518082019091526000808252602082015290565b6040518060400160405280836000015181526020017f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478460200151610f009190611d1d565b610f2a907f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47611d3f565b905292915050565b919050565b6040805180820190915260008082526020820152610f53611492565b835181526020808501519082015260408082018490526000908360608460076107d05a03fa9050808015610f8657610f88565bfe5b5080610fc65760405162461bcd60e51b815260206004820152600d60248201526c1958cb5b5d5b0b59985a5b1959609a1b60448201526064016101d4565b505092915050565b6040805180820190915260008082526020820152610fea6114b0565b835181526020808501518183015283516040808401919091529084015160608301526000908360808460066107d05a03fa9050808015610f86575080610fc65760405162461bcd60e51b815260206004820152600d60248201526c1958cb5859190b59985a5b1959609a1b60448201526064016101d4565b61106a6114ce565b50604080516080810182527f198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c28183019081527f1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed6060830152815281518083019092527f275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec82527f1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d60208381019190915281019190915290565b6040805180820182528581526020808201859052825180840190935285835282018390526000916111516114f3565b60005b600281101561131657600061116a826006611d56565b905084826002811061117e5761117e611b02565b60200201515183611190836000611d75565b600c81106111a0576111a0611b02565b60200201528482600281106111b7576111b7611b02565b602002015160200151838260016111ce9190611d75565b600c81106111de576111de611b02565b60200201528382600281106111f5576111f5611b02565b6020020151515183611208836002611d75565b600c811061121857611218611b02565b602002015283826002811061122f5761122f611b02565b6020020151516001602002015183611248836003611d75565b600c811061125857611258611b02565b602002015283826002811061126f5761126f611b02565b60200201516020015160006002811061128a5761128a611b02565b60200201518361129b836004611d75565b600c81106112ab576112ab611b02565b60200201528382600281106112c2576112c2611b02565b6020020151602001516001600281106112dd576112dd611b02565b6020020151836112ee836005611d75565b600c81106112fe576112fe611b02565b6020020152508061130e81611b2e565b915050611154565b5061131f611512565b60006020826101808560086107d05a03fa9050808015610f8657508061137f5760405162461bcd60e51b81526020600482015260156024820152741c185a5c9a5b99cb5bdc18dbd9194b59985a5b1959605a1b60448201526064016101d4565b5051151598975050505050505050565b60006020845161139f9190611d1d565b156114265760405162461bcd60e51b815260206004820152604b60248201527f4d65726b6c652e70726f63657373496e636c7573696f6e50726f6f664b65636360448201527f616b3a2070726f6f66206c656e6774682073686f756c642062652061206d756c60648201526a3a34b836329037b310199960a91b608482015260a4016101d4565b8260205b855181116114895761143d600285611d1d565b61145e57816000528086015160205260406000209150600284049350611477565b8086015160005281602052604060002091506002840493505b611482602082611d75565b905061142a565b50949350505050565b60405180606001604052806003906020820280368337509192915050565b60405180608001604052806004906020820280368337509192915050565b60405180604001604052806114e1611530565b81526020016114ee611530565b905290565b604051806101800160405280600c906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b60405180604001604052806002906020820280368337509192915050565b60006020828403121561156057600080fd5b5035919050565b6001600160a01b038416815263ffffffff83166020820152608081016103f0604083018480518252602090810151910152565b634e487b7160e01b600052604160045260246000fd5b6040516060810167ffffffffffffffff811182821017156115d3576115d361159a565b60405290565b6040516080810167ffffffffffffffff811182821017156115d3576115d361159a565b60405160a0810167ffffffffffffffff811182821017156115d3576115d361159a565b6040805190810167ffffffffffffffff811182821017156115d3576115d361159a565b604051601f8201601f1916810167ffffffffffffffff8111828210171561166b5761166b61159a565b604052919050565b803563ffffffff81168114610f3257600080fd5b803560ff81168114610f3257600080fd5b600067ffffffffffffffff8211156116b2576116b261159a565b50601f01601f191660200190565b600082601f8301126116d157600080fd5b81356116e46116df82611698565b611642565b8181528460208386010111156116f957600080fd5b816020850160208301376000918101602001919091529392505050565b60006060828403121561172857600080fd5b6117306115b0565b9050813567ffffffffffffffff8082111561174a57600080fd5b908301906080828603121561175e57600080fd5b6117666115d9565b8235815260208301358281111561177c57600080fd5b611788878286016116c0565b6020830152506040830135828111156117a057600080fd5b6117ac878286016116c0565b6040830152506117be60608401611673565b606082015283525050602082810135908201526117dd60408301611673565b604082015292915050565b600060a082840312156117fa57600080fd5b6118026115fc565b905061180d82611673565b815261181b60208301611673565b6020820152604082013567ffffffffffffffff8082111561183b57600080fd5b61184785838601611716565b6040840152606084013591508082111561186057600080fd5b61186c858386016116c0565b6060840152608084013591508082111561188557600080fd5b50611892848285016116c0565b60808301525092915050565b60008060408084860312156118b257600080fd5b833567ffffffffffffffff808211156118ca57600080fd5b908501908187036080808212156118e057600080fd5b6118e86115b0565b858312156118f557600080fd5b6118fd61161f565b9250843583526020808601358185015283825261191b878701611673565b8183015260609350838601358581111561193457600080fd5b8087019650508a601f87011261194957600080fd5b85358581111561195b5761195b61159a565b611969828260051b01611642565b81815260079190911b8701820190828101908d83111561198857600080fd5b978301975b828910156119f45785898f0312156119a55760008081fd5b6119ad6115d9565b6119b68a611687565b81526119c3858b01611687565b858201526119d28b8b01611687565b8b8201526119e1888b01611673565b818901528252978501979083019061198d565b98840198909852509097508801359450505080831115611a1357600080fd5b5050611a21858286016117e8565b9150509250929050565b600082601f830112611a3c57600080fd5b611a4461161f565b806040840185811115611a5657600080fd5b845b81811015611a70578035845260209384019301611a58565b509095945050505050565b60008060008084860360e0811215611a9257600080fd5b85359450602086013593506080603f1982011215611aaf57600080fd5b50611ab861161f565b611ac58760408801611a2b565b8152611ad48760808801611a2b565b60208201529396929550929360c00135925050565b600060208284031215611afb57600080fd5b5051919050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b6000600019821415611b4257611b42611b18565b5060010190565b60005b83811015611b64578181015183820152602001611b4c565b83811115610ab25750506000910152565b600060208284031215611b8757600080fd5b815167ffffffffffffffff811115611b9e57600080fd5b8201601f81018413611baf57600080fd5b8051611bbd6116df82611698565b818152856020838501011115611bd257600080fd5b611be3826020830160208601611b49565b95945050505050565b60008151808452611c04816020860160208601611b49565b601f01601f19169290920160200192915050565b60208152815160208201526000602083015160806040840152611c3e60a0840182611bec565b90506040840151601f19848303016060850152611c5b8282611bec565b91505063ffffffff60608501511660808401528091505092915050565b60208082528251805183830152810151604083015260009060a0830181850151606063ffffffff808316828801526040925082880151608080818a015285825180885260c08b0191508884019750600093505b80841015611d0e578751805160ff90811684528a82015181168b850152888201511688840152860151851686830152968801966001939093019290820190611ccb565b509a9950505050505050505050565b600082611d3a57634e487b7160e01b600052601260045260246000fd5b500690565b600082821015611d5157611d51611b18565b500390565b6000816000190483118215151615611d7057611d70611b18565b500290565b60008219821115611d8857611d88611b18565b50019056fea26469706673582212203eb8469450e15c2f6053e675e9ca80aba4ccd8fdcd79c45cc463c45b7345f09f64736f6c634300080c0033",
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
