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
	Bin: "0x60c06040523480156200001157600080fd5b5060405162001fd738038062001fd7833981016040819052620000349162000204565b816020015160ff16826000015160ff161162000063576040516308a6997560e01b815260040160405180910390fd5b6001600160a01b03808516608052831660a05281516000805460208086015160ff9081166101000261ffff199093169416939093171790558151620000af9160019190840190620000ba565b505050505062000374565b828054620000c89062000337565b90600052602060002090601f016020900481019282620000ec576000855562000137565b82601f106200010757805160ff191683800117855562000137565b8280016001018555821562000137579182015b82811115620001375782518255916020019190600101906200011a565b506200014592915062000149565b5090565b5b808211156200014557600081556001016200014a565b6001600160a01b03811681146200017657600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b0381118282101715620001b457620001b462000179565b60405290565b604051601f8201601f191681016001600160401b0381118282101715620001e557620001e562000179565b604052919050565b805160ff81168114620001ff57600080fd5b919050565b60008060008084860360a08112156200021c57600080fd5b8551620002298162000160565b809550506020808701516200023e8162000160565b94506040603f19830112156200025357600080fd5b6200025d6200018f565b91506200026d60408801620001ed565b82526200027d60608801620001ed565b8282015260808701519193506001600160401b03808311156200029f57600080fd5b828801925088601f840112620002b457600080fd5b825181811115620002c957620002c962000179565b620002dd601f8201601f19168401620001ba565b91508082528983828601011115620002f457600080fd5b60005b8181101562000314578481018401518382018501528301620002f7565b81811115620003265760008483850101525b505080935050505092959194509250565b600181811c908216806200034c57607f821691505b602082108114156200036e57634e487b7160e01b600052602260045260246000fd5b50919050565b60805160a051611c2f620003a86000396000818161010301526101b201526000818161013d01526101900152611c2f6000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c806321b9b2fb146100675780632ead0b96146100b85780639077193b146100c7578063e15234ff146100ec578063efd4532b14610101578063f8c668141461013b575b600080fd5b6040805180820182526000808252602091820181905282518084018452905460ff80821680845261010090920481169284019283528451918252915190911691810191909152015b60405180910390f35b604051600381526020016100af565b6100da6100d5366004610dcb565b610161565b60405160ff90911681526020016100af565b6100f4610280565b6040516100af9190610e89565b7f00000000000000000000000000000000000000000000000000000000000000005b6040516001600160a01b0390911681526020016100af565b7f0000000000000000000000000000000000000000000000000000000000000000610123565b604080518082019091526000805460ff80821684526101009091041660208301526001805491928392610263927f0000000000000000000000000000000000000000000000000000000000000000927f00000000000000000000000000000000000000000000000000000000000000009289928992916101e090610e9c565b80601f016020809104026020016040519081016040528092919081815260200182805461020c90610e9c565b80156102595780601f1061022e57610100808354040283529160200191610259565b820191906000526020600020905b81548152906001019060200180831161023c57829003601f168201915b5050505050610312565b50905080600581111561027857610278610ed7565b949350505050565b60606001805461028f90610e9c565b80601f01602080910402602001604051908101604052809291908181526020018280546102bb90610e9c565b80156103085780601f106102dd57610100808354040283529160200191610308565b820191906000526020600020905b8154815290600101906020018083116102eb57829003601f168201915b5050505050905090565b6000606060006103228787610354565b905061034489898360000151846020015185604001518a8a886060015161036f565b9250925050965096945050505050565b61035c610c2c565b6103688284018461156e565b9392505050565b6000606061037d88886104d9565b9092509050600182600581111561039657610396610ed7565b146103a0576104cc565b86515151604051632ecfe72b60e01b815261ffff909116600482015261041b906001600160a01b038c1690632ecfe72b90602401606060405180830381865afa1580156103f1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610415919061165e565b866105b2565b9092509050600182600581111561043457610434610ed7565b1461043e576104cc565b600061045a8a61044d8b61068f565b868c602001518b8b6106d7565b91945092509050600183600581111561047557610475610ed7565b1461048057506104cc565b875151602001516000906104949083610839565b9195509350905060018460058111156104af576104af610ed7565b146104bb5750506104cc565b6104c586826108a0565b9350935050505b9850989650505050505050565b6000606060006104ec8460000151610905565b905060008160405160200161050391815260200190565b6040516020818303038152906040528051906020012090506000866000015190506000610540876040015183858a6020015163ffffffff1661092d565b905080156105675760016040518060200160405280600081525095509550505050506105ab565b6020808801516040805163ffffffff909216928201929092529081018490526060810183905260029060800160405160208183030381529060405295509550505050505b9250929050565b600060606000836020015184600001516105cc91906116eb565b60ff1690506000856020015163ffffffff16866040015160ff1683620f42406105f59190611724565b6105ff9190611724565b61060b90612710611738565b610615919061174f565b86519091506000906106299061271061176e565b63ffffffff16905080821061065657600160405180602001604052806000815250945094505050506105ab565b604080516020810185905290810183905260608101829052600390608001604051602081830303815290604052945094505050506105ab565b6000816040516020016106ba91908151815260209182015163ffffffff169181019190915260400190565b604051602081830303815290604052805190602001209050919050565b60006060600080896001600160a01b0316636efb46368a8a8a8a6040518563ffffffff1660e01b815260040161071094939291906118bf565b600060405180830381865afa15801561072d573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526107559190810190611a2c565b5090506000915060005b885181101561081657856000015160ff168260200151828151811061078657610786611ac8565b60200260200101516107989190611ade565b6001600160601b03166064836000015183815181106107b9576107b9611ac8565b60200260200101516001600160601b03166107d4919061174f565b1061080457610801838a83815181106107ef576107ef611ac8565b0160200151600160f89190911c1b1790565b92505b8061080e81611b04565b91505061075f565b505060408051602081019091526000815260019350915096509650969350505050565b60006060600061084885610945565b905083811681141561086e57604080516020810190915260008152600193509150610899565b6040805160208101929092528181018590528051808303820181526060909201905260049250905060005b9250925092565b6000606060006108af85610945565b90508381168114156108d5575050604080516020810190915260008152600191506105ab565b604080516020810183905290810185905260059060600160405160208183030381529060405292509250506105ab565b60006109148260000151610ad7565b60208084015160408086015190516106ba949301611b1f565b60008361093b868585610b29565b1495945050505050565b6000610100825111156109d35760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a4015b60405180910390fd5b81516109e157506000919050565b600080836000815181106109f7576109f7611ac8565b0160200151600160f89190911c81901b92505b8451811015610ace57848181518110610a2557610a25611ac8565b0160200151600160f89190911c1b9150828211610aba5760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a4016109ca565b91811791610ac781611b04565b9050610a0a565b50909392505050565b6000816000015182602001518360400151604051602001610afa93929190611b54565b60408051601f1981840301815282825280516020918201206060808701519285019190915291830152016106ba565b600060208451610b399190611bcd565b15610bc05760405162461bcd60e51b815260206004820152604b60248201527f4d65726b6c652e70726f63657373496e636c7573696f6e50726f6f664b65636360448201527f616b3a2070726f6f66206c656e6774682073686f756c642062652061206d756c60648201526a3a34b836329037b310199960a91b608482015260a4016109ca565b8260205b85518111610c2357610bd7600285611bcd565b610bf857816000528086015160205260406000209150600284049350610c11565b8086015160005281602052604060002091506002840493505b610c1c602082611be1565b9050610bc4565b50949350505050565b6040805160c0810190915260006080820181815260a0830191909152815260208101610c56610c70565b8152602001610c63610c97565b8152602001606081525090565b6040518060600160405280610c83610cfd565b815260006020820152606060409091015290565b604051806101000160405280606081526020016060815260200160608152602001610cc0610d24565b8152602001610ce2604051806040016040528060008152602001600081525090565b81526020016060815260200160608152602001606081525090565b6040518060600160405280610d10610d49565b815260200160608152602001606081525090565b6040518060400160405280610d37610d76565b8152602001610d44610d76565b905290565b604080516080810182526000815260606020820152908101610d69610d94565b8152600060209091015290565b60405180604001604052806002906020820280368337509192915050565b6040805160c0810190915260006080820181815260a0830191909152815260208101610dbe610d24565b8152602001610d69610d24565b60008060208385031215610dde57600080fd5b82356001600160401b0380821115610df557600080fd5b818501915085601f830112610e0957600080fd5b813581811115610e1857600080fd5b866020828501011115610e2a57600080fd5b60209290920196919550909350505050565b6000815180845260005b81811015610e6257602081850181015186830182015201610e46565b81811115610e74576000602083870101525b50601f01601f19169290920160200192915050565b6020815260006103686020830184610e3c565b600181811c90821680610eb057607f821691505b60208210811415610ed157634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052602160045260246000fd5b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b0381118282101715610f2557610f25610eed565b60405290565b604051606081016001600160401b0381118282101715610f2557610f25610eed565b604051608081016001600160401b0381118282101715610f2557610f25610eed565b60405161010081016001600160401b0381118282101715610f2557610f25610eed565b604051601f8201601f191681016001600160401b0381118282101715610fba57610fba610eed565b604052919050565b63ffffffff81168114610fd457600080fd5b50565b8035610fe281610fc2565b919050565b600082601f830112610ff857600080fd5b81356001600160401b0381111561101157611011610eed565b611024601f8201601f1916602001610f92565b81815284602083860101111561103957600080fd5b816020850160208301376000918101602001919091529392505050565b60006040828403121561106857600080fd5b611070610f03565b9050813581526020820135602082015292915050565b600082601f83011261109757600080fd5b61109f610f03565b8060408401858111156110b157600080fd5b845b818110156110cb5780358452602093840193016110b3565b509095945050505050565b6000608082840312156110e857600080fd5b6110f0610f03565b90506110fc8383611086565b815261110b8360408401611086565b602082015292915050565b60006001600160401b0382111561112f5761112f610eed565b5060051b60200190565b600082601f83011261114a57600080fd5b8135602061115f61115a83611116565b610f92565b82815260059290921b8401810191818101908684111561117e57600080fd5b8286015b848110156111a257803561119581610fc2565b8352918301918301611182565b509695505050505050565b6000606082840312156111bf57600080fd5b6111c7610f2b565b905081356001600160401b03808211156111e057600080fd5b90830190606082860312156111f457600080fd5b6111fc610f2b565b82358281111561120b57600080fd5b83018087036101c081121561121f57600080fd5b611227610f4d565b823561ffff8116811461123957600080fd5b81526020838101358681111561124e57600080fd5b61125a8b828701610fe7565b83830152506040610160603f198501121561127457600080fd5b61127c610f4d565b935061128a8b828701611056565b84526112998b608087016110d6565b828501526112ab8b61010087016110d6565b818501526101808501356112be81610fc2565b8060608601525083818401526101a08501356060840152828652818801359450868511156112eb57600080fd5b6112f78b868a01610fe7565b828701528088013594508685111561130e57600080fd5b61131a8b868a01611139565b8187015285895261132c828b01610fd7565b828a0152808a013597508688111561134357600080fd5b61134f8b898c01610fe7565b818a0152505050505050505092915050565b600082601f83011261137257600080fd5b8135602061138261115a83611116565b82815260069290921b840181019181810190868411156113a157600080fd5b8286015b848110156111a2576113b78882611056565b8352918301916040016113a5565b600082601f8301126113d657600080fd5b813560206113e661115a83611116565b82815260059290921b8401810191818101908684111561140557600080fd5b8286015b848110156111a25780356001600160401b038111156114285760008081fd5b6114368986838b0101611139565b845250918301918301611409565b6000610180828403121561145757600080fd5b61145f610f6f565b905081356001600160401b038082111561147857600080fd5b61148485838601611139565b8352602084013591508082111561149a57600080fd5b6114a685838601611361565b602084015260408401359150808211156114bf57600080fd5b6114cb85838601611361565b60408401526114dd85606086016110d6565b60608401526114ef8560e08601611056565b608084015261012084013591508082111561150957600080fd5b61151585838601611139565b60a084015261014084013591508082111561152f57600080fd5b61153b85838601611139565b60c084015261016084013591508082111561155557600080fd5b50611562848285016113c5565b60e08301525092915050565b60006020828403121561158057600080fd5b81356001600160401b038082111561159757600080fd5b9083019081850360a08112156115ac57600080fd5b6115b4610f4d565b60408212156115c257600080fd5b6115ca610f03565b91508335825260208401356115de81610fc2565b6020830152908152604083013590828211156115f957600080fd5b611605878386016111ad565b6020820152606084013591508282111561161e57600080fd5b61162a87838601611444565b6040820152608084013591508282111561164357600080fd5b61164f87838601610fe7565b60608201529695505050505050565b60006060828403121561167057600080fd5b604051606081018181106001600160401b038211171561169257611692610eed565b60405282516116a081610fc2565b815260208301516116b081610fc2565b6020820152604083015160ff811681146116c957600080fd5b60408201529392505050565b634e487b7160e01b600052601160045260246000fd5b600060ff821660ff841680821015611705576117056116d5565b90039392505050565b634e487b7160e01b600052601260045260246000fd5b6000826117335761173361170e565b500490565b60008282101561174a5761174a6116d5565b500390565b6000816000190483118215151615611769576117696116d5565b500290565b600063ffffffff80831681851681830481118215151615611791576117916116d5565b02949350505050565b600081518084526020808501945080840160005b838110156117d057815163ffffffff16875295820195908201906001016117ae565b509495945050505050565b600081518084526020808501945080840160005b838110156117d05761180c87835180518252602090810151910152565b60409690960195908201906001016117ef565b8060005b6002811015611842578151845260209384019390910190600101611823565b50505050565b61185382825161181f565b6020810151611865604084018261181f565b505050565b600081518084526020808501808196508360051b8101915082860160005b858110156118b25782840389526118a084835161179a565b98850198935090840190600101611888565b5091979650505050505050565b8481526080602082015260006118d86080830186610e3c565b63ffffffff85166040840152828103606084015261018084518183526119008284018261179a565b9150506020850151828203602084015261191a82826117db565b9150506040850151828203604084015261193482826117db565b91505060608501516119496060840182611848565b506080850151805160e08401526020015161010083015260a0850151828203610120840152611978828261179a565b91505060c0850151828203610140840152611993828261179a565b91505060e08501518282036101608401526119ae828261186a565b9998505050505050505050565b600082601f8301126119cc57600080fd5b815160206119dc61115a83611116565b82815260059290921b840181019181810190868411156119fb57600080fd5b8286015b848110156111a25780516001600160601b0381168114611a1f5760008081fd5b83529183019183016119ff565b60008060408385031215611a3f57600080fd5b82516001600160401b0380821115611a5657600080fd5b9084019060408287031215611a6a57600080fd5b611a72610f03565b825182811115611a8157600080fd5b611a8d888286016119bb565b825250602083015182811115611aa257600080fd5b611aae888286016119bb565b602083015250809450505050602083015190509250929050565b634e487b7160e01b600052603260045260246000fd5b60006001600160601b0380831681851681830481118215151615611791576117916116d5565b6000600019821415611b1857611b186116d5565b5060010190565b838152606060208201526000611b386060830185610e3c565b8281036040840152611b4a818561179a565b9695505050505050565b60006101a061ffff86168352806020840152611b7281840186610e3c565b8451805160408601526020015160608501529150611b8d9050565b6020830151611b9f6080840182611848565b506040830151611bb3610100840182611848565b5063ffffffff606084015116610180830152949350505050565b600082611bdc57611bdc61170e565b500690565b60008219821115611bf457611bf46116d5565b50019056fea2646970667358221220b9fbc740ecdcffa93ea1d2f09729391a298b9f7d9dc713d5c573756699f6766164736f6c634300080c0033",
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
