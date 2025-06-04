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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"initEigenDAThresholdRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"},{\"name\":\"initEigenDASignatureVerifier\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"},{\"name\":\"initSecurityThresholds\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.SecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"initQuorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"certVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"checkDACert\",\"inputs\":[{\"name\":\"abiEncodedCert\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDASignatureVerifier\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"securityThresholds\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.SecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"error\",\"name\":\"InvalidSecurityThresholds\",\"inputs\":[]}]",
	Bin: "0x60c06040523480156200001157600080fd5b5060405162001fca38038062001fca833981016040819052620000349162000204565b816020015160ff16826000015160ff161162000063576040516308a6997560e01b815260040160405180910390fd5b6001600160a01b03808516608052831660a05281516000805460208086015160ff9081166101000261ffff199093169416939093171790558151620000af9160019190840190620000ba565b505050505062000374565b828054620000c89062000337565b90600052602060002090601f016020900481019282620000ec576000855562000137565b82601f106200010757805160ff191683800117855562000137565b8280016001018555821562000137579182015b82811115620001375782518255916020019190600101906200011a565b506200014592915062000149565b5090565b5b808211156200014557600081556001016200014a565b6001600160a01b03811681146200017657600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b0381118282101715620001b457620001b462000179565b60405290565b604051601f8201601f191681016001600160401b0381118282101715620001e557620001e562000179565b604052919050565b805160ff81168114620001ff57600080fd5b919050565b60008060008084860360a08112156200021c57600080fd5b8551620002298162000160565b809550506020808701516200023e8162000160565b94506040603f19830112156200025357600080fd5b6200025d6200018f565b91506200026d60408801620001ed565b82526200027d60608801620001ed565b8282015260808701519193506001600160401b03808311156200029f57600080fd5b828801925088601f840112620002b457600080fd5b825181811115620002c957620002c962000179565b620002dd601f8201601f19168401620001ba565b91508082528983828601011115620002f457600080fd5b60005b8181101562000314578481018401518382018501528301620002f7565b81811115620003265760008483850101525b505080935050505092959194509250565b600181811c908216806200034c57607f821691505b602082108114156200036e57634e487b7160e01b600052602260045260246000fd5b50919050565b60805160a051611c23620003a76000396000818160f701526101a601526000818161013101526101840152611c236000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c806321b9b2fb146100675780632ead0b96146100b85780639077193b146100cd578063e15234ff146100e0578063efd4532b146100f5578063f8c668141461012f575b600080fd5b6040805180820182526000808252602091820181905282518084018452905460ff80821680845261010090920481169284019283528451918252915190911691810191909152015b60405180910390f35b60035b60405160ff90911681526020016100af565b6100bb6100db366004610dbf565b610155565b6100e8610274565b6040516100af9190610e7d565b7f00000000000000000000000000000000000000000000000000000000000000005b6040516001600160a01b0390911681526020016100af565b7f0000000000000000000000000000000000000000000000000000000000000000610117565b604080518082019091526000805460ff80821684526101009091041660208301526001805491928392610257927f0000000000000000000000000000000000000000000000000000000000000000927f00000000000000000000000000000000000000000000000000000000000000009289928992916101d490610e90565b80601f016020809104026020016040519081016040528092919081815260200182805461020090610e90565b801561024d5780601f106102225761010080835404028352916020019161024d565b820191906000526020600020905b81548152906001019060200180831161023057829003601f168201915b5050505050610306565b50905080600581111561026c5761026c610ecb565b949350505050565b60606001805461028390610e90565b80601f01602080910402602001604051908101604052809291908181526020018280546102af90610e90565b80156102fc5780601f106102d1576101008083540402835291602001916102fc565b820191906000526020600020905b8154815290600101906020018083116102df57829003601f168201915b5050505050905090565b6000606060006103168787610348565b905061033889898360000151846020015185604001518a8a8860600151610363565b9250925050965096945050505050565b610350610c20565b61035c82840184611562565b9392505050565b6000606061037188886104cd565b9092509050600182600581111561038a5761038a610ecb565b14610394576104c0565b86515151604051632ecfe72b60e01b815261ffff909116600482015261040f906001600160a01b038c1690632ecfe72b90602401606060405180830381865afa1580156103e5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104099190611652565b866105a6565b9092509050600182600581111561042857610428610ecb565b14610432576104c0565b600061044e8a6104418b610683565b868c602001518b8b6106cb565b91945092509050600183600581111561046957610469610ecb565b1461047457506104c0565b87515160200151600090610488908361082d565b9195509350905060018460058111156104a3576104a3610ecb565b146104af5750506104c0565b6104b98682610894565b9350935050505b9850989650505050505050565b6000606060006104e084600001516108f9565b90506000816040516020016104f791815260200190565b6040516020818303038152906040528051906020012090506000866000015190506000610534876040015183858a6020015163ffffffff16610921565b9050801561055b57600160405180602001604052806000815250955095505050505061059f565b6020808801516040805163ffffffff909216928201929092529081018490526060810183905260029060800160405160208183030381529060405295509550505050505b9250929050565b600060606000836020015184600001516105c091906116df565b60ff1690506000856020015163ffffffff16866040015160ff1683620f42406105e99190611718565b6105f39190611718565b6105ff9061271061172c565b6106099190611743565b865190915060009061061d90612710611762565b63ffffffff16905080821061064a576001604051806020016040528060008152509450945050505061059f565b6040805160208101859052908101839052606081018290526003906080016040516020818303038152906040529450945050505061059f565b6000816040516020016106ae91908151815260209182015163ffffffff169181019190915260400190565b604051602081830303815290604052805190602001209050919050565b60006060600080896001600160a01b0316636efb46368a8a8a8a6040518563ffffffff1660e01b815260040161070494939291906118b3565b600060405180830381865afa158015610721573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526107499190810190611a20565b5090506000915060005b885181101561080a57856000015160ff168260200151828151811061077a5761077a611abc565b602002602001015161078c9190611ad2565b6001600160601b03166064836000015183815181106107ad576107ad611abc565b60200260200101516001600160601b03166107c89190611743565b106107f8576107f5838a83815181106107e3576107e3611abc565b0160200151600160f89190911c1b1790565b92505b8061080281611af8565b915050610753565b505060408051602081019091526000815260019350915096509650969350505050565b60006060600061083c85610939565b90508381168114156108625760408051602081019091526000815260019350915061088d565b6040805160208101929092528181018590528051808303820181526060909201905260049250905060005b9250925092565b6000606060006108a385610939565b90508381168114156108c95750506040805160208101909152600081526001915061059f565b6040805160208101839052908101859052600590606001604051602081830303815290604052925092505061059f565b60006109088260000151610acb565b60208084015160408086015190516106ae949301611b13565b60008361092f868585610b1d565b1495945050505050565b6000610100825111156109c75760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a4015b60405180910390fd5b81516109d557506000919050565b600080836000815181106109eb576109eb611abc565b0160200151600160f89190911c81901b92505b8451811015610ac257848181518110610a1957610a19611abc565b0160200151600160f89190911c1b9150828211610aae5760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a4016109be565b91811791610abb81611af8565b90506109fe565b50909392505050565b6000816000015182602001518360400151604051602001610aee93929190611b48565b60408051601f1981840301815282825280516020918201206060808701519285019190915291830152016106ae565b600060208451610b2d9190611bc1565b15610bb45760405162461bcd60e51b815260206004820152604b60248201527f4d65726b6c652e70726f63657373496e636c7573696f6e50726f6f664b65636360448201527f616b3a2070726f6f66206c656e6774682073686f756c642062652061206d756c60648201526a3a34b836329037b310199960a91b608482015260a4016109be565b8260205b85518111610c1757610bcb600285611bc1565b610bec57816000528086015160205260406000209150600284049350610c05565b8086015160005281602052604060002091506002840493505b610c10602082611bd5565b9050610bb8565b50949350505050565b6040805160c0810190915260006080820181815260a0830191909152815260208101610c4a610c64565b8152602001610c57610c8b565b8152602001606081525090565b6040518060600160405280610c77610cf1565b815260006020820152606060409091015290565b604051806101000160405280606081526020016060815260200160608152602001610cb4610d18565b8152602001610cd6604051806040016040528060008152602001600081525090565b81526020016060815260200160608152602001606081525090565b6040518060600160405280610d04610d3d565b815260200160608152602001606081525090565b6040518060400160405280610d2b610d6a565b8152602001610d38610d6a565b905290565b604080516080810182526000815260606020820152908101610d5d610d88565b8152600060209091015290565b60405180604001604052806002906020820280368337509192915050565b6040805160c0810190915260006080820181815260a0830191909152815260208101610db2610d18565b8152602001610d5d610d18565b60008060208385031215610dd257600080fd5b82356001600160401b0380821115610de957600080fd5b818501915085601f830112610dfd57600080fd5b813581811115610e0c57600080fd5b866020828501011115610e1e57600080fd5b60209290920196919550909350505050565b6000815180845260005b81811015610e5657602081850181015186830182015201610e3a565b81811115610e68576000602083870101525b50601f01601f19169290920160200192915050565b60208152600061035c6020830184610e30565b600181811c90821680610ea457607f821691505b60208210811415610ec557634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052602160045260246000fd5b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b0381118282101715610f1957610f19610ee1565b60405290565b604051606081016001600160401b0381118282101715610f1957610f19610ee1565b604051608081016001600160401b0381118282101715610f1957610f19610ee1565b60405161010081016001600160401b0381118282101715610f1957610f19610ee1565b604051601f8201601f191681016001600160401b0381118282101715610fae57610fae610ee1565b604052919050565b63ffffffff81168114610fc857600080fd5b50565b8035610fd681610fb6565b919050565b600082601f830112610fec57600080fd5b81356001600160401b0381111561100557611005610ee1565b611018601f8201601f1916602001610f86565b81815284602083860101111561102d57600080fd5b816020850160208301376000918101602001919091529392505050565b60006040828403121561105c57600080fd5b611064610ef7565b9050813581526020820135602082015292915050565b600082601f83011261108b57600080fd5b611093610ef7565b8060408401858111156110a557600080fd5b845b818110156110bf5780358452602093840193016110a7565b509095945050505050565b6000608082840312156110dc57600080fd5b6110e4610ef7565b90506110f0838361107a565b81526110ff836040840161107a565b602082015292915050565b60006001600160401b0382111561112357611123610ee1565b5060051b60200190565b600082601f83011261113e57600080fd5b8135602061115361114e8361110a565b610f86565b82815260059290921b8401810191818101908684111561117257600080fd5b8286015b8481101561119657803561118981610fb6565b8352918301918301611176565b509695505050505050565b6000606082840312156111b357600080fd5b6111bb610f1f565b905081356001600160401b03808211156111d457600080fd5b90830190606082860312156111e857600080fd5b6111f0610f1f565b8235828111156111ff57600080fd5b83018087036101c081121561121357600080fd5b61121b610f41565b823561ffff8116811461122d57600080fd5b81526020838101358681111561124257600080fd5b61124e8b828701610fdb565b83830152506040610160603f198501121561126857600080fd5b611270610f41565b935061127e8b82870161104a565b845261128d8b608087016110ca565b8285015261129f8b61010087016110ca565b818501526101808501356112b281610fb6565b8060608601525083818401526101a08501356060840152828652818801359450868511156112df57600080fd5b6112eb8b868a01610fdb565b828701528088013594508685111561130257600080fd5b61130e8b868a0161112d565b81870152858952611320828b01610fcb565b828a0152808a013597508688111561133757600080fd5b6113438b898c01610fdb565b818a0152505050505050505092915050565b600082601f83011261136657600080fd5b8135602061137661114e8361110a565b82815260069290921b8401810191818101908684111561139557600080fd5b8286015b84811015611196576113ab888261104a565b835291830191604001611399565b600082601f8301126113ca57600080fd5b813560206113da61114e8361110a565b82815260059290921b840181019181810190868411156113f957600080fd5b8286015b848110156111965780356001600160401b0381111561141c5760008081fd5b61142a8986838b010161112d565b8452509183019183016113fd565b6000610180828403121561144b57600080fd5b611453610f63565b905081356001600160401b038082111561146c57600080fd5b6114788583860161112d565b8352602084013591508082111561148e57600080fd5b61149a85838601611355565b602084015260408401359150808211156114b357600080fd5b6114bf85838601611355565b60408401526114d185606086016110ca565b60608401526114e38560e0860161104a565b60808401526101208401359150808211156114fd57600080fd5b6115098583860161112d565b60a084015261014084013591508082111561152357600080fd5b61152f8583860161112d565b60c084015261016084013591508082111561154957600080fd5b50611556848285016113b9565b60e08301525092915050565b60006020828403121561157457600080fd5b81356001600160401b038082111561158b57600080fd5b9083019081850360a08112156115a057600080fd5b6115a8610f41565b60408212156115b657600080fd5b6115be610ef7565b91508335825260208401356115d281610fb6565b6020830152908152604083013590828211156115ed57600080fd5b6115f9878386016111a1565b6020820152606084013591508282111561161257600080fd5b61161e87838601611438565b6040820152608084013591508282111561163757600080fd5b61164387838601610fdb565b60608201529695505050505050565b60006060828403121561166457600080fd5b604051606081018181106001600160401b038211171561168657611686610ee1565b604052825161169481610fb6565b815260208301516116a481610fb6565b6020820152604083015160ff811681146116bd57600080fd5b60408201529392505050565b634e487b7160e01b600052601160045260246000fd5b600060ff821660ff8416808210156116f9576116f96116c9565b90039392505050565b634e487b7160e01b600052601260045260246000fd5b60008261172757611727611702565b500490565b60008282101561173e5761173e6116c9565b500390565b600081600019048311821515161561175d5761175d6116c9565b500290565b600063ffffffff80831681851681830481118215151615611785576117856116c9565b02949350505050565b600081518084526020808501945080840160005b838110156117c457815163ffffffff16875295820195908201906001016117a2565b509495945050505050565b600081518084526020808501945080840160005b838110156117c45761180087835180518252602090810151910152565b60409690960195908201906001016117e3565b8060005b6002811015611836578151845260209384019390910190600101611817565b50505050565b611847828251611813565b60208101516118596040840182611813565b505050565b600081518084526020808501808196508360051b8101915082860160005b858110156118a657828403895261189484835161178e565b9885019893509084019060010161187c565b5091979650505050505050565b8481526080602082015260006118cc6080830186610e30565b63ffffffff85166040840152828103606084015261018084518183526118f48284018261178e565b9150506020850151828203602084015261190e82826117cf565b9150506040850151828203604084015261192882826117cf565b915050606085015161193d606084018261183c565b506080850151805160e08401526020015161010083015260a085015182820361012084015261196c828261178e565b91505060c0850151828203610140840152611987828261178e565b91505060e08501518282036101608401526119a2828261185e565b9998505050505050505050565b600082601f8301126119c057600080fd5b815160206119d061114e8361110a565b82815260059290921b840181019181810190868411156119ef57600080fd5b8286015b848110156111965780516001600160601b0381168114611a135760008081fd5b83529183019183016119f3565b60008060408385031215611a3357600080fd5b82516001600160401b0380821115611a4a57600080fd5b9084019060408287031215611a5e57600080fd5b611a66610ef7565b825182811115611a7557600080fd5b611a81888286016119af565b825250602083015182811115611a9657600080fd5b611aa2888286016119af565b602083015250809450505050602083015190509250929050565b634e487b7160e01b600052603260045260246000fd5b60006001600160601b0380831681851681830481118215151615611785576117856116c9565b6000600019821415611b0c57611b0c6116c9565b5060010190565b838152606060208201526000611b2c6060830185610e30565b8281036040840152611b3e818561178e565b9695505050505050565b60006101a061ffff86168352806020840152611b6681840186610e30565b8451805160408601526020015160608501529150611b819050565b6020830151611b93608084018261183c565b506040830151611ba761010084018261183c565b5063ffffffff606084015116610180830152949350505050565b600082611bd057611bd0611702565b500690565b60008219821115611be857611be86116c9565b50019056fea264697066735822122051491a9fe956df3c29f199f8ce05bb1914b0eccadfa09e2e34e450b87499da5264736f6c634300080c0033",
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
// Solidity: function certVersion() pure returns(uint8)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) CertVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "certVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// CertVersion is a free data retrieval call binding the contract method 0x2ead0b96.
//
// Solidity: function certVersion() pure returns(uint8)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) CertVersion() (uint8, error) {
	return _ContractEigenDACertVerifier.Contract.CertVersion(&_ContractEigenDACertVerifier.CallOpts)
}

// CertVersion is a free data retrieval call binding the contract method 0x2ead0b96.
//
// Solidity: function certVersion() pure returns(uint8)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) CertVersion() (uint8, error) {
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
