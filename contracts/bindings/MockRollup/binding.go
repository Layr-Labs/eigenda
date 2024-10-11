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

// EigenDABlobVerificationUtilsBlobVerificationProof is an auto generated low-level Go binding around an user-defined struct.
type EigenDABlobVerificationUtilsBlobVerificationProof struct {
	BatchId        uint32
	BlobIndex      uint32
	BatchMetadata  IEigenDAServiceManagerBatchMetadata
	InclusionProof []byte
	QuorumIndices  []byte
}

// IEigenDAServiceManagerBatchHeader is an auto generated low-level Go binding around an user-defined struct.
type IEigenDAServiceManagerBatchHeader struct {
	BlobHeadersRoot       [32]byte
	QuorumNumbers         []byte
	SignedStakeForQuorums []byte
	ReferenceBlockNumber  uint32
}

// IEigenDAServiceManagerBatchMetadata is an auto generated low-level Go binding around an user-defined struct.
type IEigenDAServiceManagerBatchMetadata struct {
	BatchHeader             IEigenDAServiceManagerBatchHeader
	SignatoryRecordHash     [32]byte
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
	QuorumNumber                    uint8
	AdversaryThresholdPercentage    uint8
	ConfirmationThresholdPercentage uint8
	ChunkLength                     uint32
}

// ContractMockRollupMetaData contains all meta data concerning the ContractMockRollup contract.
var ContractMockRollupMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_eigenDABlobVerifier\",\"type\":\"address\",\"internalType\":\"contractIEigenDABlobVerifier\"},{\"name\":\"_tau\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"challengeCommitment\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"point\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"challengeValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"commitments\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"confirmer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"polynomialCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDABlobVerifier\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDABlobVerifier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"postCommitment\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structIEigenDAServiceManager.QuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structEigenDABlobVerificationUtils.BlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"tau\",\"inputs\":[],\"outputs\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"}]",
	Bin: "0x608060405234801561001057600080fd5b506040516113cc3803806113cc83398101604081905261002f91610061565b600080546001600160a01b0319166001600160a01b0393909316929092179091558051600155602001516002556100ef565b600080828403606081121561007557600080fd5b83516001600160a01b038116811461008c57600080fd5b92506040601f19820112156100a057600080fd5b50604080519081016001600160401b03811182821017156100d157634e487b7160e01b600052604160045260246000fd5b60409081526020858101518352940151938101939093525092909150565b6112ce806100fe6000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c806349ce89971461005c578063680c784e146100cf578063b5144c73146100fa578063cfc4af551461010f578063d2d16eb214610132575b600080fd5b6100b761006a366004610a7b565b6003602090815260009182526040918290208054835180850190945260018201548452600290910154918301919091526001600160a01b03811691600160a01b90910463ffffffff169083565b6040516100c693929190610a94565b60405180910390f35b6000546100e2906001600160a01b031681565b6040516001600160a01b0390911681526020016100c6565b61010d610108366004610dbd565b610155565b005b60015460025461011d919082565b604080519283526020830191909152016100c6565b610145610140366004610f9a565b6102bd565b60405190151581526020016100c6565b426000908152600360205260409020546001600160a01b0316156101dd5760405162461bcd60e51b815260206004820152603460248201527f4d6f636b526f6c6c75702e706f7374436f6d6d69746d656e743a20436f6d6d696044820152731d1b595b9d08185b1c9958591e481c1bdcdd195960621b60648201526084015b60405180910390fd5b600054604051638d67b90960e01b81526001600160a01b0390911690638d67b9099061020f9085908590600401611128565b60006040518083038186803b15801561022757600080fd5b505afa15801561023b573d6000803e3d6000fd5b50506040805160608101825233815260208681015163ffffffff90811682840190815297518385019081524260009081526003845294909420925183549851909116600160a01b026001600160c01b03199098166001600160a01b03919091161796909617815590518051600183015590940151600290940193909355505050565b6000848152600360209081526040808320815160608101835281546001600160a01b038082168352600160a01b90910463ffffffff16828601528351808501855260018401548152600290930154948301949094529182015280519091166103855760405162461bcd60e51b815260206004820152603560248201527f4d6f636b526f6c6c75702e6368616c6c656e6765436f6d6d69746d656e743a2060448201527410dbdb5b5a5d1b595b9d081b9bdd081c1bdcdd1959605a1b60648201526084016101d4565b806020015163ffffffff1685106104105760405162461bcd60e51b815260206004820152604360248201527f4d6f636b526f6c6c75702e6368616c6c656e6765436f6d6d69746d656e743a2060448201527f506f696e74206d757374206265206c657373207468616e2064617461206c656e6064820152620cee8d60eb1b608482015260a4016101d4565b6000610446610441604080518082018252600080825260209182015281518083019092526001825260029082015290565b6104a8565b905061049b6104756104588389610567565b6040805180820190915260015481526002546020820152906105fe565b8661048e6104838589610567565b6040870151906105fe565b610496610692565b610752565b925050505b949350505050565b604080518082019091526000808252602082015281511580156104cd57506020820151155b156104eb575050604080518082019091526000808252602082015290565b6040518060400160405280836000015181526020017f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47846020015161053091906111e1565b61055a907f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47611219565b905292915050565b919050565b60408051808201909152600080825260208201526105836109bf565b835181526020808501519082015260408082018490526000908360608460076107d05a03fa90508080156105b6576105b8565bfe5b50806105f65760405162461bcd60e51b815260206004820152600d60248201526c1958cb5b5d5b0b59985a5b1959609a1b60448201526064016101d4565b505092915050565b604080518082019091526000808252602082015261061a6109dd565b835181526020808501518183015283516040808401919091529084015160608301526000908360808460066107d05a03fa90508080156105b65750806105f65760405162461bcd60e51b815260206004820152600d60248201526c1958cb5859190b59985a5b1959609a1b60448201526064016101d4565b61069a6109fb565b50604080516080810182527f198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c28183019081527f1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed6060830152815281518083019092527f275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec82527f1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d60208381019190915281019190915290565b604080518082018252858152602080820185905282518084019093528583528201839052600091610781610a20565b60005b600281101561094657600061079a826006611246565b90508482600281106107ae576107ae611230565b602002015151836107c0836000611265565b600c81106107d0576107d0611230565b60200201528482600281106107e7576107e7611230565b602002015160200151838260016107fe9190611265565b600c811061080e5761080e611230565b602002015283826002811061082557610825611230565b6020020151515183610838836002611265565b600c811061084857610848611230565b602002015283826002811061085f5761085f611230565b6020020151516001602002015183610878836003611265565b600c811061088857610888611230565b602002015283826002811061089f5761089f611230565b6020020151602001516000600281106108ba576108ba611230565b6020020151836108cb836004611265565b600c81106108db576108db611230565b60200201528382600281106108f2576108f2611230565b60200201516020015160016002811061090d5761090d611230565b60200201518361091e836005611265565b600c811061092e5761092e611230565b6020020152508061093e8161127d565b915050610784565b5061094f610a3f565b60006020826101808560086107d05a03fa90508080156105b65750806109af5760405162461bcd60e51b81526020600482015260156024820152741c185a5c9a5b99cb5bdc18dbd9194b59985a5b1959605a1b60448201526064016101d4565b5051151598975050505050505050565b60405180606001604052806003906020820280368337509192915050565b60405180608001604052806004906020820280368337509192915050565b6040518060400160405280610a0e610a5d565b8152602001610a1b610a5d565b905290565b604051806101800160405280600c906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b60405180604001604052806002906020820280368337509192915050565b600060208284031215610a8d57600080fd5b5035919050565b6001600160a01b038416815263ffffffff83166020820152608081016104a0604083018480518252602090810151910152565b634e487b7160e01b600052604160045260246000fd5b6040516060810167ffffffffffffffff81118282101715610b0057610b00610ac7565b60405290565b6040516080810167ffffffffffffffff81118282101715610b0057610b00610ac7565b60405160a0810167ffffffffffffffff81118282101715610b0057610b00610ac7565b6040805190810167ffffffffffffffff81118282101715610b0057610b00610ac7565b604051601f8201601f1916810167ffffffffffffffff81118282101715610b9857610b98610ac7565b604052919050565b803563ffffffff8116811461056257600080fd5b803560ff8116811461056257600080fd5b600082601f830112610bd657600080fd5b813567ffffffffffffffff811115610bf057610bf0610ac7565b610c03601f8201601f1916602001610b6f565b818152846020838601011115610c1857600080fd5b816020850160208301376000918101602001919091529392505050565b600060608284031215610c4757600080fd5b610c4f610add565b9050813567ffffffffffffffff80821115610c6957600080fd5b9083019060808286031215610c7d57600080fd5b610c85610b06565b82358152602083013582811115610c9b57600080fd5b610ca787828601610bc5565b602083015250604083013582811115610cbf57600080fd5b610ccb87828601610bc5565b604083015250610cdd60608401610ba0565b60608201528352505060208281013590820152610cfc60408301610ba0565b604082015292915050565b600060a08284031215610d1957600080fd5b610d21610b29565b9050610d2c82610ba0565b8152610d3a60208301610ba0565b6020820152604082013567ffffffffffffffff80821115610d5a57600080fd5b610d6685838601610c35565b60408401526060840135915080821115610d7f57600080fd5b610d8b85838601610bc5565b60608401526080840135915080821115610da457600080fd5b50610db184828501610bc5565b60808301525092915050565b6000806040808486031215610dd157600080fd5b833567ffffffffffffffff80821115610de957600080fd5b90850190818703608080821215610dff57600080fd5b610e07610add565b85831215610e1457600080fd5b610e1c610b4c565b92508435835260208086013581850152838252610e3a878701610ba0565b81830152606093508386013585811115610e5357600080fd5b8087019650508a601f870112610e6857600080fd5b853585811115610e7a57610e7a610ac7565b610e88828260051b01610b6f565b81815260079190911b8701820190828101908d831115610ea757600080fd5b978301975b82891015610f135785898f031215610ec45760008081fd5b610ecc610b06565b610ed58a610bb4565b8152610ee2858b01610bb4565b85820152610ef18b8b01610bb4565b8b820152610f00888b01610ba0565b8189015282529785019790830190610eac565b98840198909852509097508801359450505080831115610f3257600080fd5b5050610f4085828601610d07565b9150509250929050565b600082601f830112610f5b57600080fd5b610f63610b4c565b806040840185811115610f7557600080fd5b845b81811015610f8f578035845260209384019301610f77565b509095945050505050565b60008060008084860360e0811215610fb157600080fd5b85359450602086013593506080603f1982011215610fce57600080fd5b50610fd7610b4c565b610fe48760408801610f4a565b8152610ff38760808801610f4a565b60208201529396929550929360c00135925050565b6000815180845260005b8181101561102e57602081850181015186830182015201611012565b81811115611040576000602083870101525b50601f01601f19169290920160200192915050565b600063ffffffff808351168452806020840151166020850152604083015160a060408601528051606060a08701528051610100870152602081015160806101208801526110a6610180880182611008565b9050604082015160ff19888303016101408901526110c48282611008565b91505083606083015116610160880152602083015160c08801528360408401511660e08801526060860151935086810360608801526111038185611008565b93505050506080830151848203608086015261111f8282611008565b95945050505050565b6040808252835180518383015260200151606083015260009060c0830160208681015163ffffffff9081166080878101919091528489015160a088018290528051948590528301936000929060e08901905b808510156111c0578651805160ff90811684528782015181168885015289820151168984015260609081015184169083015295850195600194909401939083019061117a565b50888103858a01526111d2818b611055565b9b9a5050505050505050505050565b6000826111fe57634e487b7160e01b600052601260045260246000fd5b500690565b634e487b7160e01b600052601160045260246000fd5b60008282101561122b5761122b611203565b500390565b634e487b7160e01b600052603260045260246000fd5b600081600019048311821515161561126057611260611203565b500290565b6000821982111561127857611278611203565b500190565b600060001982141561129157611291611203565b506001019056fea264697066735822122060856746ad35203930d7ed214e1f5bf22fd420d7c66a18f8eb472ebf3e21193864736f6c634300080c0033",
}

// ContractMockRollupABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractMockRollupMetaData.ABI instead.
var ContractMockRollupABI = ContractMockRollupMetaData.ABI

// ContractMockRollupBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractMockRollupMetaData.Bin instead.
var ContractMockRollupBin = ContractMockRollupMetaData.Bin

// DeployContractMockRollup deploys a new Ethereum contract, binding an instance of ContractMockRollup to it.
func DeployContractMockRollup(auth *bind.TransactOpts, backend bind.ContractBackend, _eigenDABlobVerifier common.Address, _tau BN254G1Point) (common.Address, *types.Transaction, *ContractMockRollup, error) {
	parsed, err := ContractMockRollupMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractMockRollupBin), backend, _eigenDABlobVerifier, _tau)
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

// EigenDABlobVerifier is a free data retrieval call binding the contract method 0x680c784e.
//
// Solidity: function eigenDABlobVerifier() view returns(address)
func (_ContractMockRollup *ContractMockRollupCaller) EigenDABlobVerifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractMockRollup.contract.Call(opts, &out, "eigenDABlobVerifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDABlobVerifier is a free data retrieval call binding the contract method 0x680c784e.
//
// Solidity: function eigenDABlobVerifier() view returns(address)
func (_ContractMockRollup *ContractMockRollupSession) EigenDABlobVerifier() (common.Address, error) {
	return _ContractMockRollup.Contract.EigenDABlobVerifier(&_ContractMockRollup.CallOpts)
}

// EigenDABlobVerifier is a free data retrieval call binding the contract method 0x680c784e.
//
// Solidity: function eigenDABlobVerifier() view returns(address)
func (_ContractMockRollup *ContractMockRollupCallerSession) EigenDABlobVerifier() (common.Address, error) {
	return _ContractMockRollup.Contract.EigenDABlobVerifier(&_ContractMockRollup.CallOpts)
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
func (_ContractMockRollup *ContractMockRollupTransactor) PostCommitment(opts *bind.TransactOpts, blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof EigenDABlobVerificationUtilsBlobVerificationProof) (*types.Transaction, error) {
	return _ContractMockRollup.contract.Transact(opts, "postCommitment", blobHeader, blobVerificationProof)
}

// PostCommitment is a paid mutator transaction binding the contract method 0xb5144c73.
//
// Solidity: function postCommitment(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) returns()
func (_ContractMockRollup *ContractMockRollupSession) PostCommitment(blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof EigenDABlobVerificationUtilsBlobVerificationProof) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.PostCommitment(&_ContractMockRollup.TransactOpts, blobHeader, blobVerificationProof)
}

// PostCommitment is a paid mutator transaction binding the contract method 0xb5144c73.
//
// Solidity: function postCommitment(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) returns()
func (_ContractMockRollup *ContractMockRollupTransactorSession) PostCommitment(blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof EigenDABlobVerificationUtilsBlobVerificationProof) (*types.Transaction, error) {
	return _ContractMockRollup.Contract.PostCommitment(&_ContractMockRollup.TransactOpts, blobHeader, blobVerificationProof)
}
