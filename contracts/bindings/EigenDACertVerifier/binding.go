// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractEigenDACertVerifier

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

// EigenDATypesV1SecurityThresholds is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV1SecurityThresholds struct {
	ConfirmationThreshold uint8
	AdversaryThreshold    uint8
}

// ContractEigenDACertVerifierMetaData contains all meta data concerning the ContractEigenDACertVerifier contract.
var ContractEigenDACertVerifierMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"initEigenDAThresholdRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"},{\"name\":\"initEigenDASignatureVerifier\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"},{\"name\":\"initSecurityThresholds\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.SecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"initQuorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"certVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"checkDACert\",\"inputs\":[{\"name\":\"abiEncodedCert\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDASignatureVerifier\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"securityThresholds\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.SecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"error\",\"name\":\"InvalidSecurityThresholds\",\"inputs\":[]}]",
	Bin: "0x60c06040523480156200001157600080fd5b5060405162001f6338038062001f63833981016040819052620000349162000204565b816020015160ff16826000015160ff161162000063576040516308a6997560e01b815260040160405180910390fd5b6001600160a01b03808516608052831660a05281516000805460208086015160ff9081166101000261ffff199093169416939093171790558151620000af9160019190840190620000ba565b505050505062000374565b828054620000c89062000337565b90600052602060002090601f016020900481019282620000ec576000855562000137565b82601f106200010757805160ff191683800117855562000137565b8280016001018555821562000137579182015b82811115620001375782518255916020019190600101906200011a565b506200014592915062000149565b5090565b5b808211156200014557600081556001016200014a565b6001600160a01b03811681146200017657600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b0381118282101715620001b457620001b462000179565b60405290565b604051601f8201601f191681016001600160401b0381118282101715620001e557620001e562000179565b604052919050565b805160ff81168114620001ff57600080fd5b919050565b60008060008084860360a08112156200021c57600080fd5b8551620002298162000160565b809550506020808701516200023e8162000160565b94506040603f19830112156200025357600080fd5b6200025d6200018f565b91506200026d60408801620001ed565b82526200027d60608801620001ed565b8282015260808701519193506001600160401b03808311156200029f57600080fd5b828801925088601f840112620002b457600080fd5b825181811115620002c957620002c962000179565b620002dd601f8201601f19168401620001ba565b91508082528983828601011115620002f457600080fd5b60005b8181101562000314578481018401518382018501528301620002f7565b81811115620003265760008483850101525b505080935050505092959194509250565b600181811c908216806200034c57607f821691505b602082108114156200036e57634e487b7160e01b600052602260045260246000fd5b50919050565b60805160a051611bbb620003a86000396000818161010301526101b201526000818161013d01526101900152611bbb6000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c806321b9b2fb146100675780632ead0b96146100b85780639077193b146100c7578063e15234ff146100ec578063efd4532b14610101578063f8c668141461013b575b600080fd5b6040805180820182526000808252602091820181905282518084018452905460ff80821680845261010090920481169284019283528451918252915190911691810191909152015b60405180910390f35b604051600381526020016100af565b6100da6100d5366004610d8d565b610161565b60405160ff90911681526020016100af565b6100f4610280565b6040516100af9190610e4b565b7f00000000000000000000000000000000000000000000000000000000000000005b6040516001600160a01b0390911681526020016100af565b7f0000000000000000000000000000000000000000000000000000000000000000610123565b604080518082019091526000805460ff80821684526101009091041660208301526001805491928392610263927f0000000000000000000000000000000000000000000000000000000000000000927f00000000000000000000000000000000000000000000000000000000000000009289928992916101e090610e5e565b80601f016020809104026020016040519081016040528092919081815260200182805461020c90610e5e565b80156102595780601f1061022e57610100808354040283529160200191610259565b820191906000526020600020905b81548152906001019060200180831161023c57829003601f168201915b5050505050610312565b50905080600581111561027857610278610e99565b949350505050565b60606001805461028f90610e5e565b80601f01602080910402602001604051908101604052809291908181526020018280546102bb90610e5e565b80156103085780601f106102dd57610100808354040283529160200191610308565b820191906000526020600020905b8154815290600101906020018083116102eb57829003601f168201915b5050505050905090565b6000606060006103228787610354565b905061034489898360000151846020015185604001518a8a886060015161036f565b9250925050965096945050505050565b61035c610bee565b61036882840184611530565b9392505050565b6000606061037d88886104d9565b9092509050600182600581111561039657610396610e99565b146103a0576104cc565b86515151604051632413487f60e21b815261ffff909116600482015261041b906001600160a01b038c169063904d21fc9060240160a060405180830381865afa1580156103f1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104159190611620565b866105b2565b9092509050600182600581111561043457610434610e99565b1461043e576104cc565b600061045a8a61044d8b610651565b868c602001518b8b610699565b91945092509050600183600581111561047557610475610e99565b1461048057506104cc565b8751516020015160009061049490836107fb565b9195509350905060018460058111156104af576104af610e99565b146104bb5750506104cc565b6104c58682610862565b9350935050505b9850989650505050505050565b6000606060006104ec84600001516108c7565b905060008160405160200161050391815260200190565b6040516020818303038152906040528051906020012090506000866000015190506000610540876040015183858a6020015163ffffffff166108ef565b905080156105675760016040518060200160405280600081525095509550505050506105ab565b6020808801516040805163ffffffff909216928201929092529081018490526060810183905260029060800160405160208183030381529060405295509550505050505b9250929050565b60006060836020015163ffffffff16836020015184600001516105d591906116cd565b60ff1611156105f75750506040805160208101909152600081526001906105ab565b60038360000151846020015186602001516040516020016106399392919060ff938416815291909216602082015263ffffffff91909116604082015260600190565b604051602081830303815290604052915091506105ab565b60008160405160200161067c91908151815260209182015163ffffffff169181019190915260400190565b604051602081830303815290604052805190602001209050919050565b60006060600080896001600160a01b0316636efb46368a8a8a8a6040518563ffffffff1660e01b81526004016106d29493929190611815565b600060405180830381865afa1580156106ef573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526107179190810190611982565b5090506000915060005b88518110156107d857856000015160ff168260200151828151811061074857610748611a1e565b602002602001015161075a9190611a34565b6001600160601b031660648360000151838151811061077b5761077b611a1e565b60200260200101516001600160601b03166107969190611a63565b106107c6576107c3838a83815181106107b1576107b1611a1e565b0160200151600160f89190911c1b1790565b92505b806107d081611a82565b915050610721565b505060408051602081019091526000815260019350915096509650969350505050565b60006060600061080a85610907565b90508381168114156108305760408051602081019091526000815260019350915061085b565b6040805160208101929092528181018590528051808303820181526060909201905260049250905060005b9250925092565b60006060600061087185610907565b9050838116811415610897575050604080516020810190915260008152600191506105ab565b604080516020810183905290810185905260059060600160405160208183030381529060405292509250506105ab565b60006108d68260000151610a99565b602080840151604080860151905161067c949301611a9d565b6000836108fd868585610aeb565b1495945050505050565b6000610100825111156109955760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a4015b60405180910390fd5b81516109a357506000919050565b600080836000815181106109b9576109b9611a1e565b0160200151600160f89190911c81901b92505b8451811015610a90578481815181106109e7576109e7611a1e565b0160200151600160f89190911c1b9150828211610a7c5760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a40161098c565b91811791610a8981611a82565b90506109cc565b50909392505050565b6000816000015182602001518360400151604051602001610abc93929190611ad2565b60408051601f19818403018152828252805160209182012060608087015192850191909152918301520161067c565b600060208451610afb9190611b4b565b15610b825760405162461bcd60e51b815260206004820152604b60248201527f4d65726b6c652e70726f63657373496e636c7573696f6e50726f6f664b65636360448201527f616b3a2070726f6f66206c656e6774682073686f756c642062652061206d756c60648201526a3a34b836329037b310199960a91b608482015260a40161098c565b8260205b85518111610be557610b99600285611b4b565b610bba57816000528086015160205260406000209150600284049350610bd3565b8086015160005281602052604060002091506002840493505b610bde602082611b6d565b9050610b86565b50949350505050565b6040805160c0810190915260006080820181815260a0830191909152815260208101610c18610c32565b8152602001610c25610c59565b8152602001606081525090565b6040518060600160405280610c45610cbf565b815260006020820152606060409091015290565b604051806101000160405280606081526020016060815260200160608152602001610c82610ce6565b8152602001610ca4604051806040016040528060008152602001600081525090565b81526020016060815260200160608152602001606081525090565b6040518060600160405280610cd2610d0b565b815260200160608152602001606081525090565b6040518060400160405280610cf9610d38565b8152602001610d06610d38565b905290565b604080516080810182526000815260606020820152908101610d2b610d56565b8152600060209091015290565b60405180604001604052806002906020820280368337509192915050565b6040805160c0810190915260006080820181815260a0830191909152815260208101610d80610ce6565b8152602001610d2b610ce6565b60008060208385031215610da057600080fd5b82356001600160401b0380821115610db757600080fd5b818501915085601f830112610dcb57600080fd5b813581811115610dda57600080fd5b866020828501011115610dec57600080fd5b60209290920196919550909350505050565b6000815180845260005b81811015610e2457602081850181015186830182015201610e08565b81811115610e36576000602083870101525b50601f01601f19169290920160200192915050565b6020815260006103686020830184610dfe565b600181811c90821680610e7257607f821691505b60208210811415610e9357634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052602160045260246000fd5b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b0381118282101715610ee757610ee7610eaf565b60405290565b604051606081016001600160401b0381118282101715610ee757610ee7610eaf565b604051608081016001600160401b0381118282101715610ee757610ee7610eaf565b60405161010081016001600160401b0381118282101715610ee757610ee7610eaf565b604051601f8201601f191681016001600160401b0381118282101715610f7c57610f7c610eaf565b604052919050565b63ffffffff81168114610f9657600080fd5b50565b8035610fa481610f84565b919050565b600082601f830112610fba57600080fd5b81356001600160401b03811115610fd357610fd3610eaf565b610fe6601f8201601f1916602001610f54565b818152846020838601011115610ffb57600080fd5b816020850160208301376000918101602001919091529392505050565b60006040828403121561102a57600080fd5b611032610ec5565b9050813581526020820135602082015292915050565b600082601f83011261105957600080fd5b611061610ec5565b80604084018581111561107357600080fd5b845b8181101561108d578035845260209384019301611075565b509095945050505050565b6000608082840312156110aa57600080fd5b6110b2610ec5565b90506110be8383611048565b81526110cd8360408401611048565b602082015292915050565b60006001600160401b038211156110f1576110f1610eaf565b5060051b60200190565b600082601f83011261110c57600080fd5b8135602061112161111c836110d8565b610f54565b82815260059290921b8401810191818101908684111561114057600080fd5b8286015b8481101561116457803561115781610f84565b8352918301918301611144565b509695505050505050565b60006060828403121561118157600080fd5b611189610eed565b905081356001600160401b03808211156111a257600080fd5b90830190606082860312156111b657600080fd5b6111be610eed565b8235828111156111cd57600080fd5b83018087036101c08112156111e157600080fd5b6111e9610f0f565b823561ffff811681146111fb57600080fd5b81526020838101358681111561121057600080fd5b61121c8b828701610fa9565b83830152506040610160603f198501121561123657600080fd5b61123e610f0f565b935061124c8b828701611018565b845261125b8b60808701611098565b8285015261126d8b6101008701611098565b8185015261018085013561128081610f84565b8060608601525083818401526101a08501356060840152828652818801359450868511156112ad57600080fd5b6112b98b868a01610fa9565b82870152808801359450868511156112d057600080fd5b6112dc8b868a016110fb565b818701528589526112ee828b01610f99565b828a0152808a013597508688111561130557600080fd5b6113118b898c01610fa9565b818a0152505050505050505092915050565b600082601f83011261133457600080fd5b8135602061134461111c836110d8565b82815260069290921b8401810191818101908684111561136357600080fd5b8286015b84811015611164576113798882611018565b835291830191604001611367565b600082601f83011261139857600080fd5b813560206113a861111c836110d8565b82815260059290921b840181019181810190868411156113c757600080fd5b8286015b848110156111645780356001600160401b038111156113ea5760008081fd5b6113f88986838b01016110fb565b8452509183019183016113cb565b6000610180828403121561141957600080fd5b611421610f31565b905081356001600160401b038082111561143a57600080fd5b611446858386016110fb565b8352602084013591508082111561145c57600080fd5b61146885838601611323565b6020840152604084013591508082111561148157600080fd5b61148d85838601611323565b604084015261149f8560608601611098565b60608401526114b18560e08601611018565b60808401526101208401359150808211156114cb57600080fd5b6114d7858386016110fb565b60a08401526101408401359150808211156114f157600080fd5b6114fd858386016110fb565b60c084015261016084013591508082111561151757600080fd5b5061152484828501611387565b60e08301525092915050565b60006020828403121561154257600080fd5b81356001600160401b038082111561155957600080fd5b9083019081850360a081121561156e57600080fd5b611576610f0f565b604082121561158457600080fd5b61158c610ec5565b91508335825260208401356115a081610f84565b6020830152908152604083013590828211156115bb57600080fd5b6115c78783860161116f565b602082015260608401359150828211156115e057600080fd5b6115ec87838601611406565b6040820152608084013591508282111561160557600080fd5b61161187838601610fa9565b60608201529695505050505050565b600060a0828403121561163257600080fd5b60405160a081018181106001600160401b038211171561165457611654610eaf565b604052825161166281610f84565b8152602083015161167281610f84565b6020820152604083015161168581610f84565b6040820152606083015161169881610f84565b606082015260808301516116ab81610f84565b60808201529392505050565b634e487b7160e01b600052601160045260246000fd5b600060ff821660ff8416808210156116e7576116e76116b7565b90039392505050565b600081518084526020808501945080840160005b8381101561172657815163ffffffff1687529582019590820190600101611704565b509495945050505050565b600081518084526020808501945080840160005b838110156117265761176287835180518252602090810151910152565b6040969096019590820190600101611745565b8060005b6002811015611798578151845260209384019390910190600101611779565b50505050565b6117a9828251611775565b60208101516117bb6040840182611775565b505050565b600081518084526020808501808196508360051b8101915082860160005b858110156118085782840389526117f68483516116f0565b988501989350908401906001016117de565b5091979650505050505050565b84815260806020820152600061182e6080830186610dfe565b63ffffffff8516604084015282810360608401526101808451818352611856828401826116f0565b915050602085015182820360208401526118708282611731565b9150506040850151828203604084015261188a8282611731565b915050606085015161189f606084018261179e565b506080850151805160e08401526020015161010083015260a08501518282036101208401526118ce82826116f0565b91505060c08501518282036101408401526118e982826116f0565b91505060e085015182820361016084015261190482826117c0565b9998505050505050505050565b600082601f83011261192257600080fd5b8151602061193261111c836110d8565b82815260059290921b8401810191818101908684111561195157600080fd5b8286015b848110156111645780516001600160601b03811681146119755760008081fd5b8352918301918301611955565b6000806040838503121561199557600080fd5b82516001600160401b03808211156119ac57600080fd5b90840190604082870312156119c057600080fd5b6119c8610ec5565b8251828111156119d757600080fd5b6119e388828601611911565b8252506020830151828111156119f857600080fd5b611a0488828601611911565b602083015250809450505050602083015190509250929050565b634e487b7160e01b600052603260045260246000fd5b60006001600160601b0380831681851681830481118215151615611a5a57611a5a6116b7565b02949350505050565b6000816000190483118215151615611a7d57611a7d6116b7565b500290565b6000600019821415611a9657611a966116b7565b5060010190565b838152606060208201526000611ab66060830185610dfe565b8281036040840152611ac881856116f0565b9695505050505050565b60006101a061ffff86168352806020840152611af081840186610dfe565b8451805160408601526020015160608501529150611b0b9050565b6020830151611b1d608084018261179e565b506040830151611b3161010084018261179e565b5063ffffffff606084015116610180830152949350505050565b600082611b6857634e487b7160e01b600052601260045260246000fd5b500690565b60008219821115611b8057611b806116b7565b50019056fea2646970667358221220e7b84d43684228b838e1e4302d70a60e6a2ef3fa725e4c457b2d5f60554b37aa64736f6c634300080c0033",
}

// ContractEigenDACertVerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEigenDACertVerifierMetaData.ABI instead.
var ContractEigenDACertVerifierABI = ContractEigenDACertVerifierMetaData.ABI

// ContractEigenDACertVerifierBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractEigenDACertVerifierMetaData.Bin instead.
var ContractEigenDACertVerifierBin = ContractEigenDACertVerifierMetaData.Bin

// DeployContractEigenDACertVerifier deploys a new Ethereum contract, binding an instance of ContractEigenDACertVerifier to it.
func DeployContractEigenDACertVerifier(auth *bind.TransactOpts, backend bind.ContractBackend, initEigenDAThresholdRegistry common.Address, initEigenDASignatureVerifier common.Address, initSecurityThresholds EigenDATypesV1SecurityThresholds, initQuorumNumbersRequired []byte) (common.Address, *types.Transaction, *ContractEigenDACertVerifier, error) {
	parsed, err := ContractEigenDACertVerifierMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractEigenDACertVerifierBin), backend, initEigenDAThresholdRegistry, initEigenDASignatureVerifier, initSecurityThresholds, initQuorumNumbersRequired)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractEigenDACertVerifier{ContractEigenDACertVerifierCaller: ContractEigenDACertVerifierCaller{contract: contract}, ContractEigenDACertVerifierTransactor: ContractEigenDACertVerifierTransactor{contract: contract}, ContractEigenDACertVerifierFilterer: ContractEigenDACertVerifierFilterer{contract: contract}}, nil
}

// ContractEigenDACertVerifier is an auto generated Go binding around an Ethereum contract.
type ContractEigenDACertVerifier struct {
	ContractEigenDACertVerifierCaller     // Read-only binding to the contract
	ContractEigenDACertVerifierTransactor // Write-only binding to the contract
	ContractEigenDACertVerifierFilterer   // Log filterer for contract events
}

// ContractEigenDACertVerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractEigenDACertVerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractEigenDACertVerifierSession struct {
	Contract     *ContractEigenDACertVerifier // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                // Call options to use throughout this session
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ContractEigenDACertVerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractEigenDACertVerifierCallerSession struct {
	Contract *ContractEigenDACertVerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                      // Call options to use throughout this session
}

// ContractEigenDACertVerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractEigenDACertVerifierTransactorSession struct {
	Contract     *ContractEigenDACertVerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                      // Transaction auth options to use throughout this session
}

// ContractEigenDACertVerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractEigenDACertVerifierRaw struct {
	Contract *ContractEigenDACertVerifier // Generic contract binding to access the raw methods on
}

// ContractEigenDACertVerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierCallerRaw struct {
	Contract *ContractEigenDACertVerifierCaller // Generic read-only contract binding to access the raw methods on
}

// ContractEigenDACertVerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierTransactorRaw struct {
	Contract *ContractEigenDACertVerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractEigenDACertVerifier creates a new instance of ContractEigenDACertVerifier, bound to a specific deployed contract.
func NewContractEigenDACertVerifier(address common.Address, backend bind.ContractBackend) (*ContractEigenDACertVerifier, error) {
	contract, err := bindContractEigenDACertVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifier{ContractEigenDACertVerifierCaller: ContractEigenDACertVerifierCaller{contract: contract}, ContractEigenDACertVerifierTransactor: ContractEigenDACertVerifierTransactor{contract: contract}, ContractEigenDACertVerifierFilterer: ContractEigenDACertVerifierFilterer{contract: contract}}, nil
}

// NewContractEigenDACertVerifierCaller creates a new read-only instance of ContractEigenDACertVerifier, bound to a specific deployed contract.
func NewContractEigenDACertVerifierCaller(address common.Address, caller bind.ContractCaller) (*ContractEigenDACertVerifierCaller, error) {
	contract, err := bindContractEigenDACertVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierCaller{contract: contract}, nil
}

// NewContractEigenDACertVerifierTransactor creates a new write-only instance of ContractEigenDACertVerifier, bound to a specific deployed contract.
func NewContractEigenDACertVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractEigenDACertVerifierTransactor, error) {
	contract, err := bindContractEigenDACertVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierTransactor{contract: contract}, nil
}

// NewContractEigenDACertVerifierFilterer creates a new log filterer instance of ContractEigenDACertVerifier, bound to a specific deployed contract.
func NewContractEigenDACertVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractEigenDACertVerifierFilterer, error) {
	contract, err := bindContractEigenDACertVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierFilterer{contract: contract}, nil
}

// bindContractEigenDACertVerifier binds a generic wrapper to an already deployed contract.
func bindContractEigenDACertVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractEigenDACertVerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDACertVerifier.Contract.ContractEigenDACertVerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDACertVerifier.Contract.ContractEigenDACertVerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDACertVerifier.Contract.ContractEigenDACertVerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDACertVerifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDACertVerifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDACertVerifier.Contract.contract.Transact(opts, method, params...)
}

// CertVersion is a free data retrieval call binding the contract method 0x2ead0b96.
//
// Solidity: function certVersion() pure returns(uint64)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) CertVersion(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "certVersion")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// CertVersion is a free data retrieval call binding the contract method 0x2ead0b96.
//
// Solidity: function certVersion() pure returns(uint64)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) CertVersion() (uint64, error) {
	return _ContractEigenDACertVerifier.Contract.CertVersion(&_ContractEigenDACertVerifier.CallOpts)
}

// CertVersion is a free data retrieval call binding the contract method 0x2ead0b96.
//
// Solidity: function certVersion() pure returns(uint64)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) CertVersion() (uint64, error) {
	return _ContractEigenDACertVerifier.Contract.CertVersion(&_ContractEigenDACertVerifier.CallOpts)
}

// CheckDACert is a free data retrieval call binding the contract method 0x9077193b.
//
// Solidity: function checkDACert(bytes abiEncodedCert) view returns(uint8)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) CheckDACert(opts *bind.CallOpts, abiEncodedCert []byte) (uint8, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "checkDACert", abiEncodedCert)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// CheckDACert is a free data retrieval call binding the contract method 0x9077193b.
//
// Solidity: function checkDACert(bytes abiEncodedCert) view returns(uint8)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) CheckDACert(abiEncodedCert []byte) (uint8, error) {
	return _ContractEigenDACertVerifier.Contract.CheckDACert(&_ContractEigenDACertVerifier.CallOpts, abiEncodedCert)
}

// CheckDACert is a free data retrieval call binding the contract method 0x9077193b.
//
// Solidity: function checkDACert(bytes abiEncodedCert) view returns(uint8)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) CheckDACert(abiEncodedCert []byte) (uint8, error) {
	return _ContractEigenDACertVerifier.Contract.CheckDACert(&_ContractEigenDACertVerifier.CallOpts, abiEncodedCert)
}

// EigenDASignatureVerifier is a free data retrieval call binding the contract method 0xefd4532b.
//
// Solidity: function eigenDASignatureVerifier() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) EigenDASignatureVerifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "eigenDASignatureVerifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDASignatureVerifier is a free data retrieval call binding the contract method 0xefd4532b.
//
// Solidity: function eigenDASignatureVerifier() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) EigenDASignatureVerifier() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.EigenDASignatureVerifier(&_ContractEigenDACertVerifier.CallOpts)
}

// EigenDASignatureVerifier is a free data retrieval call binding the contract method 0xefd4532b.
//
// Solidity: function eigenDASignatureVerifier() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) EigenDASignatureVerifier() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.EigenDASignatureVerifier(&_ContractEigenDACertVerifier.CallOpts)
}

// EigenDAThresholdRegistry is a free data retrieval call binding the contract method 0xf8c66814.
//
// Solidity: function eigenDAThresholdRegistry() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) EigenDAThresholdRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "eigenDAThresholdRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDAThresholdRegistry is a free data retrieval call binding the contract method 0xf8c66814.
//
// Solidity: function eigenDAThresholdRegistry() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) EigenDAThresholdRegistry() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.EigenDAThresholdRegistry(&_ContractEigenDACertVerifier.CallOpts)
}

// EigenDAThresholdRegistry is a free data retrieval call binding the contract method 0xf8c66814.
//
// Solidity: function eigenDAThresholdRegistry() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) EigenDAThresholdRegistry() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.EigenDAThresholdRegistry(&_ContractEigenDACertVerifier.CallOpts)
}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) QuorumNumbersRequired(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "quorumNumbersRequired")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) QuorumNumbersRequired() ([]byte, error) {
	return _ContractEigenDACertVerifier.Contract.QuorumNumbersRequired(&_ContractEigenDACertVerifier.CallOpts)
}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) QuorumNumbersRequired() ([]byte, error) {
	return _ContractEigenDACertVerifier.Contract.QuorumNumbersRequired(&_ContractEigenDACertVerifier.CallOpts)
}

// SecurityThresholds is a free data retrieval call binding the contract method 0x21b9b2fb.
//
// Solidity: function securityThresholds() view returns((uint8,uint8))
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) SecurityThresholds(opts *bind.CallOpts) (EigenDATypesV1SecurityThresholds, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "securityThresholds")

	if err != nil {
		return *new(EigenDATypesV1SecurityThresholds), err
	}

	out0 := *abi.ConvertType(out[0], new(EigenDATypesV1SecurityThresholds)).(*EigenDATypesV1SecurityThresholds)

	return out0, err

}

// SecurityThresholds is a free data retrieval call binding the contract method 0x21b9b2fb.
//
// Solidity: function securityThresholds() view returns((uint8,uint8))
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) SecurityThresholds() (EigenDATypesV1SecurityThresholds, error) {
	return _ContractEigenDACertVerifier.Contract.SecurityThresholds(&_ContractEigenDACertVerifier.CallOpts)
}

// SecurityThresholds is a free data retrieval call binding the contract method 0x21b9b2fb.
//
// Solidity: function securityThresholds() view returns((uint8,uint8))
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) SecurityThresholds() (EigenDATypesV1SecurityThresholds, error) {
	return _ContractEigenDACertVerifier.Contract.SecurityThresholds(&_ContractEigenDACertVerifier.CallOpts)
}
