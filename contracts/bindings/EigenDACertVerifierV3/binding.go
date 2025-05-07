// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractEigenDACertVerifierV3

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

// ContractEigenDACertVerifierV3MetaData contains all meta data concerning the ContractEigenDACertVerifierV3 contract.
var ContractEigenDACertVerifierV3MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_eigenDAThresholdRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"},{\"name\":\"_eigenDASignatureVerifier\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"},{\"name\":\"_securityThresholds\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.SecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"_quorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"certVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"checkDACert\",\"inputs\":[{\"name\":\"certBytes\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDASignatureVerifier\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"securityThresholds\",\"inputs\":[],\"outputs\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"error\",\"name\":\"InvalidSecurityThresholds\",\"inputs\":[]}]",
	Bin: "0x60c06040523480156200001157600080fd5b5060405162001fb038038062001fb0833981016040819052620000349162000204565b816020015160ff16826000015160ff161162000063576040516308a6997560e01b815260040160405180910390fd5b6001600160a01b03808516608052831660a05281516000805460208086015160ff9081166101000261ffff199093169416939093171790558151620000af9160019190840190620000ba565b505050505062000374565b828054620000c89062000337565b90600052602060002090601f016020900481019282620000ec576000855562000137565b82601f106200010757805160ff191683800117855562000137565b8280016001018555821562000137579182015b82811115620001375782518255916020019190600101906200011a565b506200014592915062000149565b5090565b5b808211156200014557600081556001016200014a565b6001600160a01b03811681146200017657600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b0381118282101715620001b457620001b462000179565b60405290565b604051601f8201601f191681016001600160401b0381118282101715620001e557620001e562000179565b604052919050565b805160ff81168114620001ff57600080fd5b919050565b60008060008084860360a08112156200021c57600080fd5b8551620002298162000160565b809550506020808701516200023e8162000160565b94506040603f19830112156200025357600080fd5b6200025d6200018f565b91506200026d60408801620001ed565b82526200027d60608801620001ed565b8282015260808701519193506001600160401b03808311156200029f57600080fd5b828801925088601f840112620002b457600080fd5b825181811115620002c957620002c962000179565b620002dd601f8201601f19168401620001ba565b91508082528983828601011115620002f457600080fd5b60005b8181101562000314578481018401518382018501528301620002f7565b81811115620003265760008483850101525b505080935050505092959194509250565b600181811c908216806200034c57607f821691505b602082108114156200036e57634e487b7160e01b600052602260045260246000fd5b50919050565b60805160a051611c09620003a76000396000818160de015261019001526000818161011d015261016e0152611c096000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c806321b9b2fb146100675780632ead0b961461009c5780639077193b146100b1578063e15234ff146100c4578063efd4532b146100d9578063f8c6681414610118575b600080fd5b60005461007d9060ff8082169161010090041682565b6040805160ff9384168152929091166020830152015b60405180910390f35b60035b60405160ff9091168152602001610093565b61009f6100bf366004610da5565b61013f565b6100cc61025e565b6040516100939190610e63565b6101007f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b039091168152602001610093565b6101007f000000000000000000000000000000000000000000000000000000000000000081565b604080518082019091526000805460ff80821684526101009091041660208301526001805491928392610241927f0000000000000000000000000000000000000000000000000000000000000000927f00000000000000000000000000000000000000000000000000000000000000009289928992916101be90610e76565b80601f01602080910402602001604051908101604052809291908181526020018280546101ea90610e76565b80156102375780601f1061020c57610100808354040283529160200191610237565b820191906000526020600020905b81548152906001019060200180831161021a57829003601f168201915b50505050506102ec565b50905080600581111561025657610256610eb1565b949350505050565b6001805461026b90610e76565b80601f016020809104026020016040519081016040528092919081815260200182805461029790610e76565b80156102e45780601f106102b9576101008083540402835291602001916102e4565b820191906000526020600020905b8154815290600101906020018083116102c757829003601f168201915b505050505081565b6000606060006102fc878761032e565b905061031e89898360000151846020015185604001518a8a8860600151610349565b9250925050965096945050505050565b610336610c06565b61034282840184611548565b9392505050565b6000606061035788886104b3565b9092509050600182600581111561037057610370610eb1565b1461037a576104a6565b86515151604051632ecfe72b60e01b815261ffff90911660048201526103f5906001600160a01b038c1690632ecfe72b90602401606060405180830381865afa1580156103cb573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103ef9190611638565b8661058c565b9092509050600182600581111561040e5761040e610eb1565b14610418576104a6565b60006104348a6104278b610669565b868c602001518b8b6106b1565b91945092509050600183600581111561044f5761044f610eb1565b1461045a57506104a6565b8751516020015160009061046e9083610813565b91955093509050600184600581111561048957610489610eb1565b146104955750506104a6565b61049f868261087a565b9350935050505b9850989650505050505050565b6000606060006104c684600001516108df565b90506000816040516020016104dd91815260200190565b604051602081830303815290604052805190602001209050600086600001519050600061051a876040015183858a6020015163ffffffff16610907565b90508015610541576001604051806020016040528060008152509550955050505050610585565b6020808801516040805163ffffffff909216928201929092529081018490526060810183905260029060800160405160208183030381529060405295509550505050505b9250929050565b600060606000836020015184600001516105a691906116c5565b60ff1690506000856020015163ffffffff16866040015160ff1683620f42406105cf91906116fe565b6105d991906116fe565b6105e590612710611712565b6105ef9190611729565b865190915060009061060390612710611748565b63ffffffff1690508082106106305760016040518060200160405280600081525094509450505050610585565b60408051602081018590529081018390526060810182905260039060800160405160208183030381529060405294509450505050610585565b60008160405160200161069491908151815260209182015163ffffffff169181019190915260400190565b604051602081830303815290604052805190602001209050919050565b60006060600080896001600160a01b0316636efb46368a8a8a8a6040518563ffffffff1660e01b81526004016106ea9493929190611899565b600060405180830381865afa158015610707573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261072f9190810190611a06565b5090506000915060005b88518110156107f057856000015160ff168260200151828151811061076057610760611aa2565b60200260200101516107729190611ab8565b6001600160601b031660648360000151838151811061079357610793611aa2565b60200260200101516001600160601b03166107ae9190611729565b106107de576107db838a83815181106107c9576107c9611aa2565b0160200151600160f89190911c1b1790565b92505b806107e881611ade565b915050610739565b505060408051602081019091526000815260019350915096509650969350505050565b6000606060006108228561091f565b905083811681141561084857604080516020810190915260008152600193509150610873565b6040805160208101929092528181018590528051808303820181526060909201905260049250905060005b9250925092565b6000606060006108898561091f565b90508381168114156108af57505060408051602081019091526000815260019150610585565b60408051602081018390529081018590526005906060016040516020818303038152906040529250925050610585565b60006108ee8260000151610ab1565b6020808401516040808601519051610694949301611af9565b600083610915868585610b03565b1495945050505050565b6000610100825111156109ad5760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a4015b60405180910390fd5b81516109bb57506000919050565b600080836000815181106109d1576109d1611aa2565b0160200151600160f89190911c81901b92505b8451811015610aa8578481815181106109ff576109ff611aa2565b0160200151600160f89190911c1b9150828211610a945760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a4016109a4565b91811791610aa181611ade565b90506109e4565b50909392505050565b6000816000015182602001518360400151604051602001610ad493929190611b2e565b60408051601f198184030181528282528051602091820120606080870151928501919091529183015201610694565b600060208451610b139190611ba7565b15610b9a5760405162461bcd60e51b815260206004820152604b60248201527f4d65726b6c652e70726f63657373496e636c7573696f6e50726f6f664b65636360448201527f616b3a2070726f6f66206c656e6774682073686f756c642062652061206d756c60648201526a3a34b836329037b310199960a91b608482015260a4016109a4565b8260205b85518111610bfd57610bb1600285611ba7565b610bd257816000528086015160205260406000209150600284049350610beb565b8086015160005281602052604060002091506002840493505b610bf6602082611bbb565b9050610b9e565b50949350505050565b6040805160c0810190915260006080820181815260a0830191909152815260208101610c30610c4a565b8152602001610c3d610c71565b8152602001606081525090565b6040518060600160405280610c5d610cd7565b815260006020820152606060409091015290565b604051806101000160405280606081526020016060815260200160608152602001610c9a610cfe565b8152602001610cbc604051806040016040528060008152602001600081525090565b81526020016060815260200160608152602001606081525090565b6040518060600160405280610cea610d23565b815260200160608152602001606081525090565b6040518060400160405280610d11610d50565b8152602001610d1e610d50565b905290565b604080516080810182526000815260606020820152908101610d43610d6e565b8152600060209091015290565b60405180604001604052806002906020820280368337509192915050565b6040805160c0810190915260006080820181815260a0830191909152815260208101610d98610cfe565b8152602001610d43610cfe565b60008060208385031215610db857600080fd5b82356001600160401b0380821115610dcf57600080fd5b818501915085601f830112610de357600080fd5b813581811115610df257600080fd5b866020828501011115610e0457600080fd5b60209290920196919550909350505050565b6000815180845260005b81811015610e3c57602081850181015186830182015201610e20565b81811115610e4e576000602083870101525b50601f01601f19169290920160200192915050565b6020815260006103426020830184610e16565b600181811c90821680610e8a57607f821691505b60208210811415610eab57634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052602160045260246000fd5b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b0381118282101715610eff57610eff610ec7565b60405290565b604051606081016001600160401b0381118282101715610eff57610eff610ec7565b604051608081016001600160401b0381118282101715610eff57610eff610ec7565b60405161010081016001600160401b0381118282101715610eff57610eff610ec7565b604051601f8201601f191681016001600160401b0381118282101715610f9457610f94610ec7565b604052919050565b63ffffffff81168114610fae57600080fd5b50565b8035610fbc81610f9c565b919050565b600082601f830112610fd257600080fd5b81356001600160401b03811115610feb57610feb610ec7565b610ffe601f8201601f1916602001610f6c565b81815284602083860101111561101357600080fd5b816020850160208301376000918101602001919091529392505050565b60006040828403121561104257600080fd5b61104a610edd565b9050813581526020820135602082015292915050565b600082601f83011261107157600080fd5b611079610edd565b80604084018581111561108b57600080fd5b845b818110156110a557803584526020938401930161108d565b509095945050505050565b6000608082840312156110c257600080fd5b6110ca610edd565b90506110d68383611060565b81526110e58360408401611060565b602082015292915050565b60006001600160401b0382111561110957611109610ec7565b5060051b60200190565b600082601f83011261112457600080fd5b81356020611139611134836110f0565b610f6c565b82815260059290921b8401810191818101908684111561115857600080fd5b8286015b8481101561117c57803561116f81610f9c565b835291830191830161115c565b509695505050505050565b60006060828403121561119957600080fd5b6111a1610f05565b905081356001600160401b03808211156111ba57600080fd5b90830190606082860312156111ce57600080fd5b6111d6610f05565b8235828111156111e557600080fd5b83018087036101c08112156111f957600080fd5b611201610f27565b823561ffff8116811461121357600080fd5b81526020838101358681111561122857600080fd5b6112348b828701610fc1565b83830152506040610160603f198501121561124e57600080fd5b611256610f27565b93506112648b828701611030565b84526112738b608087016110b0565b828501526112858b61010087016110b0565b8185015261018085013561129881610f9c565b8060608601525083818401526101a08501356060840152828652818801359450868511156112c557600080fd5b6112d18b868a01610fc1565b82870152808801359450868511156112e857600080fd5b6112f48b868a01611113565b81870152858952611306828b01610fb1565b828a0152808a013597508688111561131d57600080fd5b6113298b898c01610fc1565b818a0152505050505050505092915050565b600082601f83011261134c57600080fd5b8135602061135c611134836110f0565b82815260069290921b8401810191818101908684111561137b57600080fd5b8286015b8481101561117c576113918882611030565b83529183019160400161137f565b600082601f8301126113b057600080fd5b813560206113c0611134836110f0565b82815260059290921b840181019181810190868411156113df57600080fd5b8286015b8481101561117c5780356001600160401b038111156114025760008081fd5b6114108986838b0101611113565b8452509183019183016113e3565b6000610180828403121561143157600080fd5b611439610f49565b905081356001600160401b038082111561145257600080fd5b61145e85838601611113565b8352602084013591508082111561147457600080fd5b6114808583860161133b565b6020840152604084013591508082111561149957600080fd5b6114a58583860161133b565b60408401526114b785606086016110b0565b60608401526114c98560e08601611030565b60808401526101208401359150808211156114e357600080fd5b6114ef85838601611113565b60a084015261014084013591508082111561150957600080fd5b61151585838601611113565b60c084015261016084013591508082111561152f57600080fd5b5061153c8482850161139f565b60e08301525092915050565b60006020828403121561155a57600080fd5b81356001600160401b038082111561157157600080fd5b9083019081850360a081121561158657600080fd5b61158e610f27565b604082121561159c57600080fd5b6115a4610edd565b91508335825260208401356115b881610f9c565b6020830152908152604083013590828211156115d357600080fd5b6115df87838601611187565b602082015260608401359150828211156115f857600080fd5b6116048783860161141e565b6040820152608084013591508282111561161d57600080fd5b61162987838601610fc1565b60608201529695505050505050565b60006060828403121561164a57600080fd5b604051606081018181106001600160401b038211171561166c5761166c610ec7565b604052825161167a81610f9c565b8152602083015161168a81610f9c565b6020820152604083015160ff811681146116a357600080fd5b60408201529392505050565b634e487b7160e01b600052601160045260246000fd5b600060ff821660ff8416808210156116df576116df6116af565b90039392505050565b634e487b7160e01b600052601260045260246000fd5b60008261170d5761170d6116e8565b500490565b600082821015611724576117246116af565b500390565b6000816000190483118215151615611743576117436116af565b500290565b600063ffffffff8083168185168183048111821515161561176b5761176b6116af565b02949350505050565b600081518084526020808501945080840160005b838110156117aa57815163ffffffff1687529582019590820190600101611788565b509495945050505050565b600081518084526020808501945080840160005b838110156117aa576117e687835180518252602090810151910152565b60409690960195908201906001016117c9565b8060005b600281101561181c5781518452602093840193909101906001016117fd565b50505050565b61182d8282516117f9565b602081015161183f60408401826117f9565b505050565b600081518084526020808501808196508360051b8101915082860160005b8581101561188c57828403895261187a848351611774565b98850198935090840190600101611862565b5091979650505050505050565b8481526080602082015260006118b26080830186610e16565b63ffffffff85166040840152828103606084015261018084518183526118da82840182611774565b915050602085015182820360208401526118f482826117b5565b9150506040850151828203604084015261190e82826117b5565b91505060608501516119236060840182611822565b506080850151805160e08401526020015161010083015260a08501518282036101208401526119528282611774565b91505060c085015182820361014084015261196d8282611774565b91505060e08501518282036101608401526119888282611844565b9998505050505050505050565b600082601f8301126119a657600080fd5b815160206119b6611134836110f0565b82815260059290921b840181019181810190868411156119d557600080fd5b8286015b8481101561117c5780516001600160601b03811681146119f95760008081fd5b83529183019183016119d9565b60008060408385031215611a1957600080fd5b82516001600160401b0380821115611a3057600080fd5b9084019060408287031215611a4457600080fd5b611a4c610edd565b825182811115611a5b57600080fd5b611a6788828601611995565b825250602083015182811115611a7c57600080fd5b611a8888828601611995565b602083015250809450505050602083015190509250929050565b634e487b7160e01b600052603260045260246000fd5b60006001600160601b038083168185168183048111821515161561176b5761176b6116af565b6000600019821415611af257611af26116af565b5060010190565b838152606060208201526000611b126060830185610e16565b8281036040840152611b248185611774565b9695505050505050565b60006101a061ffff86168352806020840152611b4c81840186610e16565b8451805160408601526020015160608501529150611b679050565b6020830151611b796080840182611822565b506040830151611b8d610100840182611822565b5063ffffffff606084015116610180830152949350505050565b600082611bb657611bb66116e8565b500690565b60008219821115611bce57611bce6116af565b50019056fea2646970667358221220d88d9f3efcd921adfd00ffcf4dd4c5dff91afb03973fb9f706966e55a5bda08c64736f6c634300080c0033",
}

// ContractEigenDACertVerifierV3ABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEigenDACertVerifierV3MetaData.ABI instead.
var ContractEigenDACertVerifierV3ABI = ContractEigenDACertVerifierV3MetaData.ABI

// ContractEigenDACertVerifierV3Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractEigenDACertVerifierV3MetaData.Bin instead.
var ContractEigenDACertVerifierV3Bin = ContractEigenDACertVerifierV3MetaData.Bin

// DeployContractEigenDACertVerifierV3 deploys a new Ethereum contract, binding an instance of ContractEigenDACertVerifierV3 to it.
func DeployContractEigenDACertVerifierV3(auth *bind.TransactOpts, backend bind.ContractBackend, _eigenDAThresholdRegistry common.Address, _eigenDASignatureVerifier common.Address, _securityThresholds EigenDATypesV1SecurityThresholds, _quorumNumbersRequired []byte) (common.Address, *types.Transaction, *ContractEigenDACertVerifierV3, error) {
	parsed, err := ContractEigenDACertVerifierV3MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractEigenDACertVerifierV3Bin), backend, _eigenDAThresholdRegistry, _eigenDASignatureVerifier, _securityThresholds, _quorumNumbersRequired)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractEigenDACertVerifierV3{ContractEigenDACertVerifierV3Caller: ContractEigenDACertVerifierV3Caller{contract: contract}, ContractEigenDACertVerifierV3Transactor: ContractEigenDACertVerifierV3Transactor{contract: contract}, ContractEigenDACertVerifierV3Filterer: ContractEigenDACertVerifierV3Filterer{contract: contract}}, nil
}

// ContractEigenDACertVerifierV3 is an auto generated Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV3 struct {
	ContractEigenDACertVerifierV3Caller     // Read-only binding to the contract
	ContractEigenDACertVerifierV3Transactor // Write-only binding to the contract
	ContractEigenDACertVerifierV3Filterer   // Log filterer for contract events
}

// ContractEigenDACertVerifierV3Caller is an auto generated read-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV3Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierV3Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV3Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierV3Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractEigenDACertVerifierV3Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierV3Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractEigenDACertVerifierV3Session struct {
	Contract     *ContractEigenDACertVerifierV3 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                  // Call options to use throughout this session
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// ContractEigenDACertVerifierV3CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractEigenDACertVerifierV3CallerSession struct {
	Contract *ContractEigenDACertVerifierV3Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                        // Call options to use throughout this session
}

// ContractEigenDACertVerifierV3TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractEigenDACertVerifierV3TransactorSession struct {
	Contract     *ContractEigenDACertVerifierV3Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                        // Transaction auth options to use throughout this session
}

// ContractEigenDACertVerifierV3Raw is an auto generated low-level Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV3Raw struct {
	Contract *ContractEigenDACertVerifierV3 // Generic contract binding to access the raw methods on
}

// ContractEigenDACertVerifierV3CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV3CallerRaw struct {
	Contract *ContractEigenDACertVerifierV3Caller // Generic read-only contract binding to access the raw methods on
}

// ContractEigenDACertVerifierV3TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV3TransactorRaw struct {
	Contract *ContractEigenDACertVerifierV3Transactor // Generic write-only contract binding to access the raw methods on
}

// NewContractEigenDACertVerifierV3 creates a new instance of ContractEigenDACertVerifierV3, bound to a specific deployed contract.
func NewContractEigenDACertVerifierV3(address common.Address, backend bind.ContractBackend) (*ContractEigenDACertVerifierV3, error) {
	contract, err := bindContractEigenDACertVerifierV3(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierV3{ContractEigenDACertVerifierV3Caller: ContractEigenDACertVerifierV3Caller{contract: contract}, ContractEigenDACertVerifierV3Transactor: ContractEigenDACertVerifierV3Transactor{contract: contract}, ContractEigenDACertVerifierV3Filterer: ContractEigenDACertVerifierV3Filterer{contract: contract}}, nil
}

// NewContractEigenDACertVerifierV3Caller creates a new read-only instance of ContractEigenDACertVerifierV3, bound to a specific deployed contract.
func NewContractEigenDACertVerifierV3Caller(address common.Address, caller bind.ContractCaller) (*ContractEigenDACertVerifierV3Caller, error) {
	contract, err := bindContractEigenDACertVerifierV3(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierV3Caller{contract: contract}, nil
}

// NewContractEigenDACertVerifierV3Transactor creates a new write-only instance of ContractEigenDACertVerifierV3, bound to a specific deployed contract.
func NewContractEigenDACertVerifierV3Transactor(address common.Address, transactor bind.ContractTransactor) (*ContractEigenDACertVerifierV3Transactor, error) {
	contract, err := bindContractEigenDACertVerifierV3(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierV3Transactor{contract: contract}, nil
}

// NewContractEigenDACertVerifierV3Filterer creates a new log filterer instance of ContractEigenDACertVerifierV3, bound to a specific deployed contract.
func NewContractEigenDACertVerifierV3Filterer(address common.Address, filterer bind.ContractFilterer) (*ContractEigenDACertVerifierV3Filterer, error) {
	contract, err := bindContractEigenDACertVerifierV3(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierV3Filterer{contract: contract}, nil
}

// bindContractEigenDACertVerifierV3 binds a generic wrapper to an already deployed contract.
func bindContractEigenDACertVerifierV3(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractEigenDACertVerifierV3MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDACertVerifierV3 *ContractEigenDACertVerifierV3Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDACertVerifierV3.Contract.ContractEigenDACertVerifierV3Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDACertVerifierV3 *ContractEigenDACertVerifierV3Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierV3.Contract.ContractEigenDACertVerifierV3Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDACertVerifierV3 *ContractEigenDACertVerifierV3Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierV3.Contract.ContractEigenDACertVerifierV3Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDACertVerifierV3 *ContractEigenDACertVerifierV3CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDACertVerifierV3.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDACertVerifierV3 *ContractEigenDACertVerifierV3TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierV3.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDACertVerifierV3 *ContractEigenDACertVerifierV3TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierV3.Contract.contract.Transact(opts, method, params...)
}

// CertVersion is a free data retrieval call binding the contract method 0x2ead0b96.
//
// Solidity: function certVersion() pure returns(uint8)
func (_ContractEigenDACertVerifierV3 *ContractEigenDACertVerifierV3Caller) CertVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV3.contract.Call(opts, &out, "certVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// CertVersion is a free data retrieval call binding the contract method 0x2ead0b96.
//
// Solidity: function certVersion() pure returns(uint8)
func (_ContractEigenDACertVerifierV3 *ContractEigenDACertVerifierV3Session) CertVersion() (uint8, error) {
	return _ContractEigenDACertVerifierV3.Contract.CertVersion(&_ContractEigenDACertVerifierV3.CallOpts)
}

// CertVersion is a free data retrieval call binding the contract method 0x2ead0b96.
//
// Solidity: function certVersion() pure returns(uint8)
func (_ContractEigenDACertVerifierV3 *ContractEigenDACertVerifierV3CallerSession) CertVersion() (uint8, error) {
	return _ContractEigenDACertVerifierV3.Contract.CertVersion(&_ContractEigenDACertVerifierV3.CallOpts)
}

// CheckDACert is a free data retrieval call binding the contract method 0x9077193b.
//
// Solidity: function checkDACert(bytes certBytes) view returns(uint8)
func (_ContractEigenDACertVerifierV3 *ContractEigenDACertVerifierV3Caller) CheckDACert(opts *bind.CallOpts, certBytes []byte) (uint8, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV3.contract.Call(opts, &out, "checkDACert", certBytes)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// CheckDACert is a free data retrieval call binding the contract method 0x9077193b.
//
// Solidity: function checkDACert(bytes certBytes) view returns(uint8)
func (_ContractEigenDACertVerifierV3 *ContractEigenDACertVerifierV3Session) CheckDACert(certBytes []byte) (uint8, error) {
	return _ContractEigenDACertVerifierV3.Contract.CheckDACert(&_ContractEigenDACertVerifierV3.CallOpts, certBytes)
}

// CheckDACert is a free data retrieval call binding the contract method 0x9077193b.
//
// Solidity: function checkDACert(bytes certBytes) view returns(uint8)
func (_ContractEigenDACertVerifierV3 *ContractEigenDACertVerifierV3CallerSession) CheckDACert(certBytes []byte) (uint8, error) {
	return _ContractEigenDACertVerifierV3.Contract.CheckDACert(&_ContractEigenDACertVerifierV3.CallOpts, certBytes)
}

// EigenDASignatureVerifier is a free data retrieval call binding the contract method 0xefd4532b.
//
// Solidity: function eigenDASignatureVerifier() view returns(address)
func (_ContractEigenDACertVerifierV3 *ContractEigenDACertVerifierV3Caller) EigenDASignatureVerifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV3.contract.Call(opts, &out, "eigenDASignatureVerifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDASignatureVerifier is a free data retrieval call binding the contract method 0xefd4532b.
//
// Solidity: function eigenDASignatureVerifier() view returns(address)
func (_ContractEigenDACertVerifierV3 *ContractEigenDACertVerifierV3Session) EigenDASignatureVerifier() (common.Address, error) {
	return _ContractEigenDACertVerifierV3.Contract.EigenDASignatureVerifier(&_ContractEigenDACertVerifierV3.CallOpts)
}

// EigenDASignatureVerifier is a free data retrieval call binding the contract method 0xefd4532b.
//
// Solidity: function eigenDASignatureVerifier() view returns(address)
func (_ContractEigenDACertVerifierV3 *ContractEigenDACertVerifierV3CallerSession) EigenDASignatureVerifier() (common.Address, error) {
	return _ContractEigenDACertVerifierV3.Contract.EigenDASignatureVerifier(&_ContractEigenDACertVerifierV3.CallOpts)
}

// EigenDAThresholdRegistry is a free data retrieval call binding the contract method 0xf8c66814.
//
// Solidity: function eigenDAThresholdRegistry() view returns(address)
func (_ContractEigenDACertVerifierV3 *ContractEigenDACertVerifierV3Caller) EigenDAThresholdRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV3.contract.Call(opts, &out, "eigenDAThresholdRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDAThresholdRegistry is a free data retrieval call binding the contract method 0xf8c66814.
//
// Solidity: function eigenDAThresholdRegistry() view returns(address)
func (_ContractEigenDACertVerifierV3 *ContractEigenDACertVerifierV3Session) EigenDAThresholdRegistry() (common.Address, error) {
	return _ContractEigenDACertVerifierV3.Contract.EigenDAThresholdRegistry(&_ContractEigenDACertVerifierV3.CallOpts)
}

// EigenDAThresholdRegistry is a free data retrieval call binding the contract method 0xf8c66814.
//
// Solidity: function eigenDAThresholdRegistry() view returns(address)
func (_ContractEigenDACertVerifierV3 *ContractEigenDACertVerifierV3CallerSession) EigenDAThresholdRegistry() (common.Address, error) {
	return _ContractEigenDACertVerifierV3.Contract.EigenDAThresholdRegistry(&_ContractEigenDACertVerifierV3.CallOpts)
}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDACertVerifierV3 *ContractEigenDACertVerifierV3Caller) QuorumNumbersRequired(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV3.contract.Call(opts, &out, "quorumNumbersRequired")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDACertVerifierV3 *ContractEigenDACertVerifierV3Session) QuorumNumbersRequired() ([]byte, error) {
	return _ContractEigenDACertVerifierV3.Contract.QuorumNumbersRequired(&_ContractEigenDACertVerifierV3.CallOpts)
}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDACertVerifierV3 *ContractEigenDACertVerifierV3CallerSession) QuorumNumbersRequired() ([]byte, error) {
	return _ContractEigenDACertVerifierV3.Contract.QuorumNumbersRequired(&_ContractEigenDACertVerifierV3.CallOpts)
}

// SecurityThresholds is a free data retrieval call binding the contract method 0x21b9b2fb.
//
// Solidity: function securityThresholds() view returns(uint8 confirmationThreshold, uint8 adversaryThreshold)
func (_ContractEigenDACertVerifierV3 *ContractEigenDACertVerifierV3Caller) SecurityThresholds(opts *bind.CallOpts) (struct {
	ConfirmationThreshold uint8
	AdversaryThreshold    uint8
}, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV3.contract.Call(opts, &out, "securityThresholds")

	outstruct := new(struct {
		ConfirmationThreshold uint8
		AdversaryThreshold    uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfirmationThreshold = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.AdversaryThreshold = *abi.ConvertType(out[1], new(uint8)).(*uint8)

	return *outstruct, err

}

// SecurityThresholds is a free data retrieval call binding the contract method 0x21b9b2fb.
//
// Solidity: function securityThresholds() view returns(uint8 confirmationThreshold, uint8 adversaryThreshold)
func (_ContractEigenDACertVerifierV3 *ContractEigenDACertVerifierV3Session) SecurityThresholds() (struct {
	ConfirmationThreshold uint8
	AdversaryThreshold    uint8
}, error) {
	return _ContractEigenDACertVerifierV3.Contract.SecurityThresholds(&_ContractEigenDACertVerifierV3.CallOpts)
}

// SecurityThresholds is a free data retrieval call binding the contract method 0x21b9b2fb.
//
// Solidity: function securityThresholds() view returns(uint8 confirmationThreshold, uint8 adversaryThreshold)
func (_ContractEigenDACertVerifierV3 *ContractEigenDACertVerifierV3CallerSession) SecurityThresholds() (struct {
	ConfirmationThreshold uint8
	AdversaryThreshold    uint8
}, error) {
	return _ContractEigenDACertVerifierV3.Contract.SecurityThresholds(&_ContractEigenDACertVerifierV3.CallOpts)
}
