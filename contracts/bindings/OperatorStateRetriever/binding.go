// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractOperatorStateRetriever

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

// OperatorStateRetrieverCheckSignaturesIndices is an auto generated low-level Go binding around an user-defined struct.
type OperatorStateRetrieverCheckSignaturesIndices struct {
	NonSignerQuorumBitmapIndices []uint32
	QuorumApkIndices             []uint32
	TotalStakeIndices            []uint32
	NonSignerStakeIndices        [][]uint32
}

// OperatorStateRetrieverOperator is an auto generated low-level Go binding around an user-defined struct.
type OperatorStateRetrieverOperator struct {
	Operator   common.Address
	OperatorId [32]byte
	Stake      *big.Int
}

// ContractOperatorStateRetrieverMetaData contains all meta data concerning the ContractOperatorStateRetriever contract.
var ContractOperatorStateRetrieverMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"getBatchOperatorFromId\",\"inputs\":[{\"name\":\"registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"},{\"name\":\"operatorIds\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[{\"name\":\"operators\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBatchOperatorId\",\"inputs\":[{\"name\":\"registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"},{\"name\":\"operators\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[{\"name\":\"operatorIds\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCheckSignaturesIndices\",\"inputs\":[{\"name\":\"registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"nonSignerOperatorIds\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structOperatorStateRetriever.CheckSignaturesIndices\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorState\",\"inputs\":[{\"name\":\"registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"blockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[][]\",\"internalType\":\"structOperatorStateRetriever.Operator[][]\",\"components\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"stake\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorState\",\"inputs\":[{\"name\":\"registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"},{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"blockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"tuple[][]\",\"internalType\":\"structOperatorStateRetriever.Operator[][]\",\"components\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"stake\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumBitmapsAtBlockNumber\",\"inputs\":[{\"name\":\"registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"},{\"name\":\"operatorIds\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"blockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"}]",
	Bin: "0x608060405234801561001057600080fd5b50611e03806100206000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c806331b36bd9146100675780633563b0d1146100905780634d2b57fe146100b05780634f739f74146100d05780635c155662146100f0578063cefdc1d414610110575b600080fd5b61007a6100753660046113fa565b610131565b60405161008791906114e8565b60405180910390f35b6100a361009e366004611524565b61024d565b604051610087919061167f565b6100c36100be3660046116f8565b6106e3565b6040516100879190611747565b6100e36100de3660046117df565b6107f8565b60405161008791906118d7565b6101036100fe366004611992565b610f22565b60405161008791906119f5565b61012361011e366004611a2d565b6110ea565b604051610087929190611a64565b606081516001600160401b0381111561014c5761014c611391565b604051908082528060200260200182016040528015610175578160200160208202803683370190505b50905060005b825181101561024657836001600160a01b03166313542a4e8483815181106101a5576101a5611a85565b60200260200101516040518263ffffffff1660e01b81526004016101d891906001600160a01b0391909116815260200190565b602060405180830381865afa1580156101f5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102199190611a9b565b82828151811061022b5761022b611a85565b602090810291909101015261023f81611aca565b905061017b565b5092915050565b60606000846001600160a01b031663683048356040518163ffffffff1660e01b8152600401602060405180830381865afa15801561028f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102b39190611ae5565b90506000856001600160a01b0316639e9923c26040518163ffffffff1660e01b8152600401602060405180830381865afa1580156102f5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103199190611ae5565b90506000866001600160a01b0316635df459466040518163ffffffff1660e01b8152600401602060405180830381865afa15801561035b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061037f9190611ae5565b9050600086516001600160401b0381111561039c5761039c611391565b6040519080825280602002602001820160405280156103cf57816020015b60608152602001906001900390816103ba5790505b50905060005b87518110156106d75760008882815181106103f2576103f2611a85565b0160200151604051638902624560e01b815260f89190911c6004820181905263ffffffff8a16602483015291506000906001600160a01b03871690638902624590604401600060405180830381865afa158015610453573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261047b9190810190611b02565b905080516001600160401b0381111561049657610496611391565b6040519080825280602002602001820160405280156104e157816020015b60408051606081018252600080825260208083018290529282015282526000199092019101816104b45790505b508484815181106104f4576104f4611a85565b602002602001018190525060005b81518110156106c1576040518060600160405280876001600160a01b03166347b314e885858151811061053757610537611a85565b60200260200101516040518263ffffffff1660e01b815260040161055d91815260200190565b602060405180830381865afa15801561057a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061059e9190611ae5565b6001600160a01b031681526020018383815181106105be576105be611a85565b60200260200101518152602001896001600160a01b031663fa28c6278585815181106105ec576105ec611a85565b60209081029190910101516040516001600160e01b031960e084901b168152600481019190915260ff8816602482015263ffffffff8f166044820152606401602060405180830381865afa158015610648573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061066c9190611b92565b6001600160601b031681525085858151811061068a5761068a611a85565b602002602001015182815181106106a3576106a3611a85565b602002602001018190525080806106b990611aca565b915050610502565b50505080806106cf90611aca565b9150506103d5565b50979650505050505050565b606081516001600160401b038111156106fe576106fe611391565b604051908082528060200260200182016040528015610727578160200160208202803683370190505b50905060005b825181101561024657836001600160a01b031663296bb06484838151811061075757610757611a85565b60200260200101516040518263ffffffff1660e01b815260040161077d91815260200190565b602060405180830381865afa15801561079a573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107be9190611ae5565b8282815181106107d0576107d0611a85565b6001600160a01b03909216602092830291909101909101526107f181611aca565b905061072d565b6108236040518060800160405280606081526020016060815260200160608152602001606081525090565b6000876001600160a01b031663683048356040518163ffffffff1660e01b8152600401602060405180830381865afa158015610863573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108879190611ae5565b90506108b46040518060800160405280606081526020016060815260200160608152602001606081525090565b6040516361c8a12f60e11b81526001600160a01b038a169063c391425e906108e4908b9089908990600401611bbb565b600060405180830381865afa158015610901573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526109299190810190611c05565b81526040516340e03a8160e11b81526001600160a01b038316906381c075029061095b908b908b908b90600401611cbc565b600060405180830381865afa158015610978573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526109a09190810190611c05565b6040820152856001600160401b038111156109bd576109bd611391565b6040519080825280602002602001820160405280156109f057816020015b60608152602001906001900390816109db5790505b50606082015260005b60ff8116871115610e33576000856001600160401b03811115610a1e57610a1e611391565b604051908082528060200260200182016040528015610a47578160200160208202803683370190505b5083606001518360ff1681518110610a6157610a61611a85565b602002602001018190525060005b86811015610d335760008c6001600160a01b03166304ec63518a8a85818110610a9a57610a9a611a85565b905060200201358e88600001518681518110610ab857610ab8611a85565b60200260200101516040518463ffffffff1660e01b8152600401610af59392919092835263ffffffff918216602084015216604082015260600190565b602060405180830381865afa158015610b12573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b369190611ce5565b90506001600160c01b038116610bde5760405162461bcd60e51b815260206004820152605c60248201527f4f70657261746f7253746174655265747269657665722e676574436865636b5360448201527f69676e617475726573496e64696365733a206f70657261746f72206d7573742060648201527f6265207265676973746572656420617420626c6f636b6e756d62657200000000608482015260a40160405180910390fd5b8a8a8560ff16818110610bf357610bf3611a85565b6001600160c01b03841692013560f81c9190911c600190811614159050610d2057856001600160a01b031663dd9846b98a8a85818110610c3557610c35611a85565b905060200201358d8d8860ff16818110610c5157610c51611a85565b6040516001600160e01b031960e087901b1681526004810194909452919091013560f81c60248301525063ffffffff8f166044820152606401602060405180830381865afa158015610ca7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ccb9190611d0e565b85606001518560ff1681518110610ce457610ce4611a85565b60200260200101518481518110610cfd57610cfd611a85565b63ffffffff9092166020928302919091019091015282610d1c81611aca565b9350505b5080610d2b81611aca565b915050610a6f565b506000816001600160401b03811115610d4e57610d4e611391565b604051908082528060200260200182016040528015610d77578160200160208202803683370190505b50905060005b82811015610df85784606001518460ff1681518110610d9e57610d9e611a85565b60200260200101518181518110610db757610db7611a85565b6020026020010151828281518110610dd157610dd1611a85565b63ffffffff9092166020928302919091019091015280610df081611aca565b915050610d7d565b508084606001518460ff1681518110610e1357610e13611a85565b602002602001018190525050508080610e2b90611d2b565b9150506109f9565b506000896001600160a01b0316635df459466040518163ffffffff1660e01b8152600401602060405180830381865afa158015610e74573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e989190611ae5565b60405163354952a360e21b81529091506001600160a01b0382169063d5254a8c90610ecb908b908b908e90600401611d4b565b600060405180830381865afa158015610ee8573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610f109190810190611c05565b60208301525098975050505050505050565b60606000846001600160a01b031663c391425e84866040518363ffffffff1660e01b8152600401610f54929190611d75565b600060405180830381865afa158015610f71573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610f999190810190611c05565b9050600084516001600160401b03811115610fb657610fb6611391565b604051908082528060200260200182016040528015610fdf578160200160208202803683370190505b50905060005b85518110156110e057866001600160a01b03166304ec635187838151811061100f5761100f611a85565b60200260200101518786858151811061102a5761102a611a85565b60200260200101516040518463ffffffff1660e01b81526004016110679392919092835263ffffffff918216602084015216604082015260600190565b602060405180830381865afa158015611084573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110a89190611ce5565b6001600160c01b03168282815181106110c3576110c3611a85565b6020908102919091010152806110d881611aca565b915050610fe5565b5095945050505050565b604080516001808252818301909252600091606091839160208083019080368337019050509050848160008151811061112557611125611a85565b60209081029190910101526040516361c8a12f60e11b81526000906001600160a01b0388169063c391425e906111619088908690600401611d75565b600060405180830381865afa15801561117e573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526111a69190810190611c05565b6000815181106111b8576111b8611a85565b60209081029190910101516040516304ec635160e01b81526004810188905263ffffffff87811660248301529091166044820181905291506000906001600160a01b038916906304ec635190606401602060405180830381865afa158015611224573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112489190611ce5565b6001600160c01b03169050600061125e8261127c565b90508161126c8a838a61024d565b9550955050505050935093915050565b606060008061128a84611348565b61ffff166001600160401b038111156112a5576112a5611391565b6040519080825280601f01601f1916602001820160405280156112cf576020820181803683370190505b5090506000805b8251821080156112e7575061010081105b1561133e576001811b93508584161561132e578060f81b83838151811061131057611310611a85565b60200101906001600160f81b031916908160001a9053508160010191505b61133781611aca565b90506112d6565b5090949350505050565b6000805b82156113735761135d600184611d94565b909216918061136b81611dab565b91505061134c565b92915050565b6001600160a01b038116811461138e57600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b03811182821017156113cf576113cf611391565b604052919050565b60006001600160401b038211156113f0576113f0611391565b5060051b60200190565b6000806040838503121561140d57600080fd5b823561141881611379565b91506020838101356001600160401b0381111561143457600080fd5b8401601f8101861361144557600080fd5b8035611458611453826113d7565b6113a7565b81815260059190911b8201830190838101908883111561147757600080fd5b928401925b8284101561149e57833561148f81611379565b8252928401929084019061147c565b80955050505050509250929050565b600081518084526020808501945080840160005b838110156114dd578151875295820195908201906001016114c1565b509495945050505050565b6020815260006114fb60208301846114ad565b9392505050565b63ffffffff8116811461138e57600080fd5b803561151f81611502565b919050565b60008060006060848603121561153957600080fd5b833561154481611379565b92506020848101356001600160401b038082111561156157600080fd5b818701915087601f83011261157557600080fd5b81358181111561158757611587611391565b611599601f8201601f191685016113a7565b915080825288848285010111156115af57600080fd5b80848401858401376000848284010152508094505050506115d260408501611514565b90509250925092565b600081518084526020808501808196508360051b810191508286016000805b86811015611671578385038a52825180518087529087019087870190845b8181101561165c57835180516001600160a01b031684528a8101518b8501526040908101516001600160601b03169084015292890192606090920191600101611618565b50509a87019a955050918501916001016115fa565b509298975050505050505050565b6020815260006114fb60208301846115db565b600082601f8301126116a357600080fd5b813560206116b3611453836113d7565b82815260059290921b840181019181810190868411156116d257600080fd5b8286015b848110156116ed57803583529183019183016116d6565b509695505050505050565b6000806040838503121561170b57600080fd5b823561171681611379565b915060208301356001600160401b0381111561173157600080fd5b61173d85828601611692565b9150509250929050565b6020808252825182820181905260009190848201906040850190845b818110156117885783516001600160a01b031683529284019291840191600101611763565b50909695505050505050565b60008083601f8401126117a657600080fd5b5081356001600160401b038111156117bd57600080fd5b6020830191508360208260051b85010111156117d857600080fd5b9250929050565b600080600080600080608087890312156117f857600080fd5b863561180381611379565b9550602087013561181381611502565b945060408701356001600160401b038082111561182f57600080fd5b818901915089601f83011261184357600080fd5b81358181111561185257600080fd5b8a602082850101111561186457600080fd5b60208301965080955050606089013591508082111561188257600080fd5b5061188f89828a01611794565b979a9699509497509295939492505050565b600081518084526020808501945080840160005b838110156114dd57815163ffffffff16875295820195908201906001016118b5565b6000602080835283516080828501526118f360a08501826118a1565b905081850151601f198086840301604087015261191083836118a1565b9250604087015191508086840301606087015261192d83836118a1565b60608801518782038301608089015280518083529194508501925084840190600581901b8501860160005b8281101561198457848783030184526119728287516118a1565b95880195938801939150600101611958565b509998505050505050505050565b6000806000606084860312156119a757600080fd5b83356119b281611379565b925060208401356001600160401b038111156119cd57600080fd5b6119d986828701611692565b92505060408401356119ea81611502565b809150509250925092565b6020808252825182820181905260009190848201906040850190845b8181101561178857835183529284019291840191600101611a11565b600080600060608486031215611a4257600080fd5b8335611a4d81611379565b92506020840135915060408401356119ea81611502565b828152604060208201526000611a7d60408301846115db565b949350505050565b634e487b7160e01b600052603260045260246000fd5b600060208284031215611aad57600080fd5b5051919050565b634e487b7160e01b600052601160045260246000fd5b6000600019821415611ade57611ade611ab4565b5060010190565b600060208284031215611af757600080fd5b81516114fb81611379565b60006020808385031215611b1557600080fd5b82516001600160401b03811115611b2b57600080fd5b8301601f81018513611b3c57600080fd5b8051611b4a611453826113d7565b81815260059190911b82018301908381019087831115611b6957600080fd5b928401925b82841015611b8757835182529284019290840190611b6e565b979650505050505050565b600060208284031215611ba457600080fd5b81516001600160601b03811681146114fb57600080fd5b63ffffffff84168152604060208201819052810182905260006001600160fb1b03831115611be857600080fd5b8260051b8085606085013760009201606001918252509392505050565b60006020808385031215611c1857600080fd5b82516001600160401b03811115611c2e57600080fd5b8301601f81018513611c3f57600080fd5b8051611c4d611453826113d7565b81815260059190911b82018301908381019087831115611c6c57600080fd5b928401925b82841015611b87578351611c8481611502565b82529284019290840190611c71565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b63ffffffff84168152604060208201526000611cdc604083018486611c93565b95945050505050565b600060208284031215611cf757600080fd5b81516001600160c01b03811681146114fb57600080fd5b600060208284031215611d2057600080fd5b81516114fb81611502565b600060ff821660ff811415611d4257611d42611ab4565b60010192915050565b604081526000611d5f604083018587611c93565b905063ffffffff83166020830152949350505050565b63ffffffff83168152604060208201526000611a7d60408301846114ad565b600082821015611da657611da6611ab4565b500390565b600061ffff80831681811415611dc357611dc3611ab4565b600101939250505056fea2646970667358221220b60f9c63dd6e254991b8775fd03325d59528eb1e3a3b531ca875ed05687e09a564736f6c634300080c0033",
}

// ContractOperatorStateRetrieverABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractOperatorStateRetrieverMetaData.ABI instead.
var ContractOperatorStateRetrieverABI = ContractOperatorStateRetrieverMetaData.ABI

// ContractOperatorStateRetrieverBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractOperatorStateRetrieverMetaData.Bin instead.
var ContractOperatorStateRetrieverBin = ContractOperatorStateRetrieverMetaData.Bin

// DeployContractOperatorStateRetriever deploys a new Ethereum contract, binding an instance of ContractOperatorStateRetriever to it.
func DeployContractOperatorStateRetriever(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ContractOperatorStateRetriever, error) {
	parsed, err := ContractOperatorStateRetrieverMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractOperatorStateRetrieverBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractOperatorStateRetriever{ContractOperatorStateRetrieverCaller: ContractOperatorStateRetrieverCaller{contract: contract}, ContractOperatorStateRetrieverTransactor: ContractOperatorStateRetrieverTransactor{contract: contract}, ContractOperatorStateRetrieverFilterer: ContractOperatorStateRetrieverFilterer{contract: contract}}, nil
}

// ContractOperatorStateRetriever is an auto generated Go binding around an Ethereum contract.
type ContractOperatorStateRetriever struct {
	ContractOperatorStateRetrieverCaller     // Read-only binding to the contract
	ContractOperatorStateRetrieverTransactor // Write-only binding to the contract
	ContractOperatorStateRetrieverFilterer   // Log filterer for contract events
}

// ContractOperatorStateRetrieverCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractOperatorStateRetrieverCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractOperatorStateRetrieverTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractOperatorStateRetrieverTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractOperatorStateRetrieverFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractOperatorStateRetrieverFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractOperatorStateRetrieverSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractOperatorStateRetrieverSession struct {
	Contract     *ContractOperatorStateRetriever // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                   // Call options to use throughout this session
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// ContractOperatorStateRetrieverCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractOperatorStateRetrieverCallerSession struct {
	Contract *ContractOperatorStateRetrieverCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                         // Call options to use throughout this session
}

// ContractOperatorStateRetrieverTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractOperatorStateRetrieverTransactorSession struct {
	Contract     *ContractOperatorStateRetrieverTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                         // Transaction auth options to use throughout this session
}

// ContractOperatorStateRetrieverRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractOperatorStateRetrieverRaw struct {
	Contract *ContractOperatorStateRetriever // Generic contract binding to access the raw methods on
}

// ContractOperatorStateRetrieverCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractOperatorStateRetrieverCallerRaw struct {
	Contract *ContractOperatorStateRetrieverCaller // Generic read-only contract binding to access the raw methods on
}

// ContractOperatorStateRetrieverTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractOperatorStateRetrieverTransactorRaw struct {
	Contract *ContractOperatorStateRetrieverTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractOperatorStateRetriever creates a new instance of ContractOperatorStateRetriever, bound to a specific deployed contract.
func NewContractOperatorStateRetriever(address common.Address, backend bind.ContractBackend) (*ContractOperatorStateRetriever, error) {
	contract, err := bindContractOperatorStateRetriever(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractOperatorStateRetriever{ContractOperatorStateRetrieverCaller: ContractOperatorStateRetrieverCaller{contract: contract}, ContractOperatorStateRetrieverTransactor: ContractOperatorStateRetrieverTransactor{contract: contract}, ContractOperatorStateRetrieverFilterer: ContractOperatorStateRetrieverFilterer{contract: contract}}, nil
}

// NewContractOperatorStateRetrieverCaller creates a new read-only instance of ContractOperatorStateRetriever, bound to a specific deployed contract.
func NewContractOperatorStateRetrieverCaller(address common.Address, caller bind.ContractCaller) (*ContractOperatorStateRetrieverCaller, error) {
	contract, err := bindContractOperatorStateRetriever(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractOperatorStateRetrieverCaller{contract: contract}, nil
}

// NewContractOperatorStateRetrieverTransactor creates a new write-only instance of ContractOperatorStateRetriever, bound to a specific deployed contract.
func NewContractOperatorStateRetrieverTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractOperatorStateRetrieverTransactor, error) {
	contract, err := bindContractOperatorStateRetriever(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractOperatorStateRetrieverTransactor{contract: contract}, nil
}

// NewContractOperatorStateRetrieverFilterer creates a new log filterer instance of ContractOperatorStateRetriever, bound to a specific deployed contract.
func NewContractOperatorStateRetrieverFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractOperatorStateRetrieverFilterer, error) {
	contract, err := bindContractOperatorStateRetriever(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractOperatorStateRetrieverFilterer{contract: contract}, nil
}

// bindContractOperatorStateRetriever binds a generic wrapper to an already deployed contract.
func bindContractOperatorStateRetriever(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractOperatorStateRetrieverMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractOperatorStateRetriever.Contract.ContractOperatorStateRetrieverCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractOperatorStateRetriever.Contract.ContractOperatorStateRetrieverTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractOperatorStateRetriever.Contract.ContractOperatorStateRetrieverTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractOperatorStateRetriever.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractOperatorStateRetriever.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractOperatorStateRetriever.Contract.contract.Transact(opts, method, params...)
}

// GetBatchOperatorFromId is a free data retrieval call binding the contract method 0x4d2b57fe.
//
// Solidity: function getBatchOperatorFromId(address registryCoordinator, bytes32[] operatorIds) view returns(address[] operators)
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverCaller) GetBatchOperatorFromId(opts *bind.CallOpts, registryCoordinator common.Address, operatorIds [][32]byte) ([]common.Address, error) {
	var out []interface{}
	err := _ContractOperatorStateRetriever.contract.Call(opts, &out, "getBatchOperatorFromId", registryCoordinator, operatorIds)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetBatchOperatorFromId is a free data retrieval call binding the contract method 0x4d2b57fe.
//
// Solidity: function getBatchOperatorFromId(address registryCoordinator, bytes32[] operatorIds) view returns(address[] operators)
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverSession) GetBatchOperatorFromId(registryCoordinator common.Address, operatorIds [][32]byte) ([]common.Address, error) {
	return _ContractOperatorStateRetriever.Contract.GetBatchOperatorFromId(&_ContractOperatorStateRetriever.CallOpts, registryCoordinator, operatorIds)
}

// GetBatchOperatorFromId is a free data retrieval call binding the contract method 0x4d2b57fe.
//
// Solidity: function getBatchOperatorFromId(address registryCoordinator, bytes32[] operatorIds) view returns(address[] operators)
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverCallerSession) GetBatchOperatorFromId(registryCoordinator common.Address, operatorIds [][32]byte) ([]common.Address, error) {
	return _ContractOperatorStateRetriever.Contract.GetBatchOperatorFromId(&_ContractOperatorStateRetriever.CallOpts, registryCoordinator, operatorIds)
}

// GetBatchOperatorId is a free data retrieval call binding the contract method 0x31b36bd9.
//
// Solidity: function getBatchOperatorId(address registryCoordinator, address[] operators) view returns(bytes32[] operatorIds)
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverCaller) GetBatchOperatorId(opts *bind.CallOpts, registryCoordinator common.Address, operators []common.Address) ([][32]byte, error) {
	var out []interface{}
	err := _ContractOperatorStateRetriever.contract.Call(opts, &out, "getBatchOperatorId", registryCoordinator, operators)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetBatchOperatorId is a free data retrieval call binding the contract method 0x31b36bd9.
//
// Solidity: function getBatchOperatorId(address registryCoordinator, address[] operators) view returns(bytes32[] operatorIds)
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverSession) GetBatchOperatorId(registryCoordinator common.Address, operators []common.Address) ([][32]byte, error) {
	return _ContractOperatorStateRetriever.Contract.GetBatchOperatorId(&_ContractOperatorStateRetriever.CallOpts, registryCoordinator, operators)
}

// GetBatchOperatorId is a free data retrieval call binding the contract method 0x31b36bd9.
//
// Solidity: function getBatchOperatorId(address registryCoordinator, address[] operators) view returns(bytes32[] operatorIds)
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverCallerSession) GetBatchOperatorId(registryCoordinator common.Address, operators []common.Address) ([][32]byte, error) {
	return _ContractOperatorStateRetriever.Contract.GetBatchOperatorId(&_ContractOperatorStateRetriever.CallOpts, registryCoordinator, operators)
}

// GetCheckSignaturesIndices is a free data retrieval call binding the contract method 0x4f739f74.
//
// Solidity: function getCheckSignaturesIndices(address registryCoordinator, uint32 referenceBlockNumber, bytes quorumNumbers, bytes32[] nonSignerOperatorIds) view returns((uint32[],uint32[],uint32[],uint32[][]))
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverCaller) GetCheckSignaturesIndices(opts *bind.CallOpts, registryCoordinator common.Address, referenceBlockNumber uint32, quorumNumbers []byte, nonSignerOperatorIds [][32]byte) (OperatorStateRetrieverCheckSignaturesIndices, error) {
	var out []interface{}
	err := _ContractOperatorStateRetriever.contract.Call(opts, &out, "getCheckSignaturesIndices", registryCoordinator, referenceBlockNumber, quorumNumbers, nonSignerOperatorIds)

	if err != nil {
		return *new(OperatorStateRetrieverCheckSignaturesIndices), err
	}

	out0 := *abi.ConvertType(out[0], new(OperatorStateRetrieverCheckSignaturesIndices)).(*OperatorStateRetrieverCheckSignaturesIndices)

	return out0, err

}

// GetCheckSignaturesIndices is a free data retrieval call binding the contract method 0x4f739f74.
//
// Solidity: function getCheckSignaturesIndices(address registryCoordinator, uint32 referenceBlockNumber, bytes quorumNumbers, bytes32[] nonSignerOperatorIds) view returns((uint32[],uint32[],uint32[],uint32[][]))
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverSession) GetCheckSignaturesIndices(registryCoordinator common.Address, referenceBlockNumber uint32, quorumNumbers []byte, nonSignerOperatorIds [][32]byte) (OperatorStateRetrieverCheckSignaturesIndices, error) {
	return _ContractOperatorStateRetriever.Contract.GetCheckSignaturesIndices(&_ContractOperatorStateRetriever.CallOpts, registryCoordinator, referenceBlockNumber, quorumNumbers, nonSignerOperatorIds)
}

// GetCheckSignaturesIndices is a free data retrieval call binding the contract method 0x4f739f74.
//
// Solidity: function getCheckSignaturesIndices(address registryCoordinator, uint32 referenceBlockNumber, bytes quorumNumbers, bytes32[] nonSignerOperatorIds) view returns((uint32[],uint32[],uint32[],uint32[][]))
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverCallerSession) GetCheckSignaturesIndices(registryCoordinator common.Address, referenceBlockNumber uint32, quorumNumbers []byte, nonSignerOperatorIds [][32]byte) (OperatorStateRetrieverCheckSignaturesIndices, error) {
	return _ContractOperatorStateRetriever.Contract.GetCheckSignaturesIndices(&_ContractOperatorStateRetriever.CallOpts, registryCoordinator, referenceBlockNumber, quorumNumbers, nonSignerOperatorIds)
}

// GetOperatorState is a free data retrieval call binding the contract method 0x3563b0d1.
//
// Solidity: function getOperatorState(address registryCoordinator, bytes quorumNumbers, uint32 blockNumber) view returns((address,bytes32,uint96)[][])
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverCaller) GetOperatorState(opts *bind.CallOpts, registryCoordinator common.Address, quorumNumbers []byte, blockNumber uint32) ([][]OperatorStateRetrieverOperator, error) {
	var out []interface{}
	err := _ContractOperatorStateRetriever.contract.Call(opts, &out, "getOperatorState", registryCoordinator, quorumNumbers, blockNumber)

	if err != nil {
		return *new([][]OperatorStateRetrieverOperator), err
	}

	out0 := *abi.ConvertType(out[0], new([][]OperatorStateRetrieverOperator)).(*[][]OperatorStateRetrieverOperator)

	return out0, err

}

// GetOperatorState is a free data retrieval call binding the contract method 0x3563b0d1.
//
// Solidity: function getOperatorState(address registryCoordinator, bytes quorumNumbers, uint32 blockNumber) view returns((address,bytes32,uint96)[][])
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverSession) GetOperatorState(registryCoordinator common.Address, quorumNumbers []byte, blockNumber uint32) ([][]OperatorStateRetrieverOperator, error) {
	return _ContractOperatorStateRetriever.Contract.GetOperatorState(&_ContractOperatorStateRetriever.CallOpts, registryCoordinator, quorumNumbers, blockNumber)
}

// GetOperatorState is a free data retrieval call binding the contract method 0x3563b0d1.
//
// Solidity: function getOperatorState(address registryCoordinator, bytes quorumNumbers, uint32 blockNumber) view returns((address,bytes32,uint96)[][])
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverCallerSession) GetOperatorState(registryCoordinator common.Address, quorumNumbers []byte, blockNumber uint32) ([][]OperatorStateRetrieverOperator, error) {
	return _ContractOperatorStateRetriever.Contract.GetOperatorState(&_ContractOperatorStateRetriever.CallOpts, registryCoordinator, quorumNumbers, blockNumber)
}

// GetOperatorState0 is a free data retrieval call binding the contract method 0xcefdc1d4.
//
// Solidity: function getOperatorState(address registryCoordinator, bytes32 operatorId, uint32 blockNumber) view returns(uint256, (address,bytes32,uint96)[][])
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverCaller) GetOperatorState0(opts *bind.CallOpts, registryCoordinator common.Address, operatorId [32]byte, blockNumber uint32) (*big.Int, [][]OperatorStateRetrieverOperator, error) {
	var out []interface{}
	err := _ContractOperatorStateRetriever.contract.Call(opts, &out, "getOperatorState0", registryCoordinator, operatorId, blockNumber)

	if err != nil {
		return *new(*big.Int), *new([][]OperatorStateRetrieverOperator), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new([][]OperatorStateRetrieverOperator)).(*[][]OperatorStateRetrieverOperator)

	return out0, out1, err

}

// GetOperatorState0 is a free data retrieval call binding the contract method 0xcefdc1d4.
//
// Solidity: function getOperatorState(address registryCoordinator, bytes32 operatorId, uint32 blockNumber) view returns(uint256, (address,bytes32,uint96)[][])
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverSession) GetOperatorState0(registryCoordinator common.Address, operatorId [32]byte, blockNumber uint32) (*big.Int, [][]OperatorStateRetrieverOperator, error) {
	return _ContractOperatorStateRetriever.Contract.GetOperatorState0(&_ContractOperatorStateRetriever.CallOpts, registryCoordinator, operatorId, blockNumber)
}

// GetOperatorState0 is a free data retrieval call binding the contract method 0xcefdc1d4.
//
// Solidity: function getOperatorState(address registryCoordinator, bytes32 operatorId, uint32 blockNumber) view returns(uint256, (address,bytes32,uint96)[][])
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverCallerSession) GetOperatorState0(registryCoordinator common.Address, operatorId [32]byte, blockNumber uint32) (*big.Int, [][]OperatorStateRetrieverOperator, error) {
	return _ContractOperatorStateRetriever.Contract.GetOperatorState0(&_ContractOperatorStateRetriever.CallOpts, registryCoordinator, operatorId, blockNumber)
}

// GetQuorumBitmapsAtBlockNumber is a free data retrieval call binding the contract method 0x5c155662.
//
// Solidity: function getQuorumBitmapsAtBlockNumber(address registryCoordinator, bytes32[] operatorIds, uint32 blockNumber) view returns(uint256[])
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverCaller) GetQuorumBitmapsAtBlockNumber(opts *bind.CallOpts, registryCoordinator common.Address, operatorIds [][32]byte, blockNumber uint32) ([]*big.Int, error) {
	var out []interface{}
	err := _ContractOperatorStateRetriever.contract.Call(opts, &out, "getQuorumBitmapsAtBlockNumber", registryCoordinator, operatorIds, blockNumber)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetQuorumBitmapsAtBlockNumber is a free data retrieval call binding the contract method 0x5c155662.
//
// Solidity: function getQuorumBitmapsAtBlockNumber(address registryCoordinator, bytes32[] operatorIds, uint32 blockNumber) view returns(uint256[])
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverSession) GetQuorumBitmapsAtBlockNumber(registryCoordinator common.Address, operatorIds [][32]byte, blockNumber uint32) ([]*big.Int, error) {
	return _ContractOperatorStateRetriever.Contract.GetQuorumBitmapsAtBlockNumber(&_ContractOperatorStateRetriever.CallOpts, registryCoordinator, operatorIds, blockNumber)
}

// GetQuorumBitmapsAtBlockNumber is a free data retrieval call binding the contract method 0x5c155662.
//
// Solidity: function getQuorumBitmapsAtBlockNumber(address registryCoordinator, bytes32[] operatorIds, uint32 blockNumber) view returns(uint256[])
func (_ContractOperatorStateRetriever *ContractOperatorStateRetrieverCallerSession) GetQuorumBitmapsAtBlockNumber(registryCoordinator common.Address, operatorIds [][32]byte, blockNumber uint32) ([]*big.Int, error) {
	return _ContractOperatorStateRetriever.Contract.GetQuorumBitmapsAtBlockNumber(&_ContractOperatorStateRetriever.CallOpts, registryCoordinator, operatorIds, blockNumber)
}
