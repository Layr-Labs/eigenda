// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractPaymentVault

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

// EigenDATypesV3QuorumPaymentConfig is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV3QuorumPaymentConfig struct {
	Token                       common.Address
	Recipient                   common.Address
	ReservationSymbolsPerSecond uint64
	OnDemandSymbolsPerPeriod    uint64
	OnDemandPricePerSymbol      uint64
}

// EigenDATypesV3QuorumPaymentProtocolConfig is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV3QuorumPaymentProtocolConfig struct {
	ReservationAdvanceWindow uint64
}

// EigenDATypesV3Reservation is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV3Reservation struct {
	SymbolsPerSecond uint64
	StartTimestamp   uint64
	EndTimestamp     uint64
}

// ContractPaymentVaultMetaData contains all meta data concerning the ContractPaymentVault contract.
var ContractPaymentVaultMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"schedulePeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"SCHEDULE_PERIOD\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"createReservation\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reservation\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV3.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"decreaseReservation\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"reservation\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV3.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositOnDemand\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getOnDemandDeposit\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumPaymentConfig\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV3.QuorumPaymentConfig\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reservationSymbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandSymbolsPerPeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandPricePerSymbol\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumProtocolConfig\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV3.QuorumPaymentProtocolConfig\",\"components\":[{\"name\":\"reservationAdvanceWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumReservedSymbols\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"period\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReservation\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV3.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"increaseReservation\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reservation\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV3.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initializeQuorum\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"protocolCfg\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV3.QuorumPaymentProtocolConfig\",\"components\":[{\"name\":\"reservationAdvanceWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setQuorumPaymentConfig\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"paymentConfig\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV3.QuorumPaymentConfig\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reservationSymbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandSymbolsPerPeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandPricePerSymbol\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setReservationAdvanceWindow\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"protocolCfg\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV3.QuorumPaymentProtocolConfig\",\"components\":[{\"name\":\"reservationAdvanceWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferQuorumOwnership\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"}]",
	Bin: "0x60a060405234801561001057600080fd5b5060405162001fda38038062001fda833981016040819052610031916100ae565b6000816001600160401b03161161009d5760405162461bcd60e51b815260206004820152602660248201527f5363686564756c6520706572696f64206d75737420626520677265617465722060448201526507468616e20360d41b606482015260840160405180910390fd5b6001600160401b03166080526100de565b6000602082840312156100c057600080fd5b81516001600160401b03811681146100d757600080fd5b9392505050565b608051611ecb6200010f6000396000818161025f015281816103b0015281816105d001526106bd0152611ecb6000f3fe608060405234801561001057600080fd5b50600436106100ff5760003560e01c8063c4d66de811610097578063ecfc2d8e11610066578063ecfc2d8e146102cd578063ede57a2d146102e0578063ef7531bd146102f3578063f2fde38b1461030657600080fd5b8063c4d66de814610281578063dc5fee2514610294578063e453a058146102a7578063e631f82c146102ba57600080fd5b80637a9426ca116100d35780637a9426ca146101af57806389a06b351461021c578063ac990b8a14610247578063c1eafe511461025a57600080fd5b8062e691aa146101045780633b50e8b81461014e578063400c2322146101635780634023a20014610184575b600080fd5b6101176101123660046118e5565b610319565b6040805182516001600160401b03908116825260208085015182169083015292820151909216908201526060015b60405180910390f35b61016161015c366004611994565b6103a8565b005b6101766101713660046118e5565b6103d8565b604051908152602001610145565b6101976101923660046119bf565b610417565b6040516001600160401b039091168152602001610145565b6101c26101bd3660046119e9565b610454565b6040805182516001600160a01b03908116825260208085015190911690820152828201516001600160401b03908116928201929092526060808401518316908201526080928301519091169181019190915260a001610145565b61022f61022a3660046119e9565b6104fb565b60405190516001600160401b03168152602001610145565b610161610255366004611a04565b61053f565b6101977f000000000000000000000000000000000000000000000000000000000000000081565b61016161028f366004611a2e565b61054a565b6101616102a2366004611aa2565b610577565b6101616102b5366004611acd565b6105be565b6101616102c8366004611b11565b6105fa565b6101616102db366004611acd565b6106ab565b6101616102ee3660046118e5565b6106e1565b610161610301366004611bd3565b61075e565b610161610314366004611a2e565b610815565b604080516060810182526000808252602082018190529181019190915261033e61089e565b6001600160401b038085166000908152602092835260408082206001600160a01b038716835260040184529081902081516060810183526001909101548084168252600160401b8104841694820194909452600160801b9093049091169082015290505b92915050565b6103d48233837f00000000000000000000000000000000000000000000000000000000000000006108ad565b5050565b60006103e261089e565b6001600160401b0384166000908152602091825260408082206001600160a01b03861683526004019092522054905092915050565b600061042161089e565b6001600160401b038085166000908152602092835260408082208684168352600501909352919091205416905092915050565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915261048761089e565b6001600160401b0392831660009081526020918252604090819020815160a08101835260018201546001600160a01b039081168252600283015490811694820194909452600160a01b909304851691830191909152600301548084166060830152600160401b900490921660808301525090565b60408051602081019091526000815261051261089e565b6001600160401b039283166000908152602091825260409081902081519283019091525490921682525090565b6103d4823383610a45565b6105747fe579393920b888b1e4a7e1afdd7d58fa4f3101113547ac874aefa75ff4a960f982610b0d565b50565b8161058181610b69565b815161058b61089e565b6001600160401b03948516600090815260209190915260409020805467ffffffffffffffff191691909416179092555050565b826105c881610b69565b6105f48484847f0000000000000000000000000000000000000000000000000000000000000000610bba565b50505050565b8161060481610b69565b8161060d61089e565b6001600160401b039485166000908152602091825260409081902083516001820180546001600160a01b039283166001600160a01b031990911617905592840151600282018054938601518916600160a01b026001600160e01b03199094169190941617919091179091556060820151600390910180546080909301518616600160401b026001600160801b03199093169190951617179092555050565b826106b581610b69565b6105f48484847f0000000000000000000000000000000000000000000000000000000000000000610ca1565b816106eb81610b69565b6001600160a01b0382166107465760405162461bcd60e51b815260206004820152601d60248201527f4e6577206f776e657220697320746865207a65726f206164647265737300000060448201526064015b60405180910390fd5b61075961075284610de3565b3384610e18565b505050565b610766610e2c565b61077761077284610de3565b610e90565b156107c45760405162461bcd60e51b815260206004820152601860248201527f51756f72756d206f776e657220616c7265616479207365740000000000000000604482015260640161073d565b6107d66107d084610de3565b83610b0d565b806107df61089e565b6001600160401b039485166000908152602091909152604090209051815467ffffffffffffffff19169416939093179092555050565b61081d610e2c565b6001600160a01b0381166108735760405162461bcd60e51b815260206004820152601d60248201527f4e6577206f776e657220697320746865207a65726f2061646472657373000000604482015260640161073d565b6105747fe579393920b888b1e4a7e1afdd7d58fa4f3101113547ac874aefa75ff4a960f93383610e18565b60006108a8610eb1565b905090565b60006108b761089e565b6001600160401b0386166000908152602091825260408082206001600160a01b03881683526004019092522060010190506108f485858585610f4f565b50805460208401516001600160401b03908116600160401b909204161461092d5760405162461bcd60e51b815260040161073d90611c0e565b805460408401516001600160401b03600160801b9092048216911611156109665760405162461bcd60e51b815260040161073d90611c45565b805483516001600160401b03918216911611156109955760405162461bcd60e51b815260040161073d90611c74565b8054604084015184516109bc928892600160801b9091046001600160401b031691866111f1565b826109c561089e565b6001600160401b039687166000908152602091825260408082206001600160a01b039098168252600490970182528690208251600190910180549284015193909701518816600160801b0267ffffffffffffffff60801b19938916600160401b026001600160801b03199093169190981617171694909417909255505050565b6000610a4f61089e565b6001600160401b0385166000908152602091825260408082206001600160a01b0387168352600481019093528120805492935091600184019190610a94908690611cc1565b905069ffffffffffffffffffff811115610ae35760405162461bcd60e51b815260206004820152601060248201526f416d6f756e7420746f6f206c6172676560801b604482015260640161073d565b60018201548254610b03916001600160a01b03918216918991168861139b565b9091555050505050565b610b2e81610b196113f5565b600085815260209190915260409020906113ff565b506040516001600160a01b0382169083907f2ae6a113c0ed5b78a53413ffbb7679881f11145ccfba4fb92e863dfcd5a1d2f390600090a35050565b610b7b610b7582610de3565b3361141b565b6105745760405162461bcd60e51b815260206004820152601060248201526f2737ba1038bab7b93ab69037bbb732b960811b604482015260640161073d565b6000610bc461089e565b6001600160401b0386166000908152602091825260408082206001600160a01b0388168352600401909252206001019050610c0185858585610f4f565b50805460208401516001600160401b03908116600160401b9092041614610c3a5760405162461bcd60e51b815260040161073d90611c0e565b805460408401516001600160401b03600160801b9092048216911611610c725760405162461bcd60e51b815260040161073d90611c45565b805483516001600160401b03918216911610156109955760405162461bcd60e51b815260040161073d90611c74565b610cad84848484610f4f565b50610cb661089e565b6001600160401b038086166000908152602092835260408082206001600160a01b0388168352600401845290206001015491840151600160801b909204811691161015610d155760405162461bcd60e51b815260040161073d90611c0e565b4282602001516001600160401b03161015610d425760405162461bcd60e51b815260040161073d90611c0e565b610d5b84836020015184604001518560000151856111f1565b81610d6461089e565b6001600160401b039586166000908152602091825260408082206001600160a01b039097168252600490960182528590208251600190910180549284015193909601518716600160801b0267ffffffffffffffff60801b19938816600160401b026001600160801b031990931691909716171716939093179091555050565b60006103a26001600160401b0383167f9cb79c1d0fdfada3ea04142fe992963bc303c019dac5ad1fb95c78752893db12611cc1565b610e22838361143e565b6107598382610b0d565b610e567fe579393920b888b1e4a7e1afdd7d58fa4f3101113547ac874aefa75ff4a960f93361141b565b610e8e5760405162461bcd60e51b81526020600482015260096024820152682737ba1037bbb732b960b91b604482015260640161073d565b565b60006103a2610e9d6113f5565b60008481526020919091526040902061149a565b60008060ff60001b19600160405180604001604052806016815260200175195a59d95b8b99184b9c185e5b595b9d0b9d985d5b1d60521b815250604051602001610efb9190611d05565b6040516020818303038152906040528051906020012060001c610f1e9190611d21565b604051602001610f3091815260200190565b60408051601f1981840301815291905280516020909101201692915050565b600080610f5a61089e565b6001600160401b0387166000908152602091825260408082206001600160a01b0389168352600481018452912091860151909250610f99908590611d4e565b6001600160401b031615610fbf5760405162461bcd60e51b815260040161073d90611c0e565b838560400151610fcf9190611d4e565b6001600160401b031615610ff55760405162461bcd60e51b815260040161073d90611c45565b84602001516001600160401b031685604001516001600160401b03161161105e5760405162461bcd60e51b815260206004820152601a60248201527f496e76616c6964207265736572766174696f6e20706572696f64000000000000604482015260640161073d565b8154602086015160408701516001600160401b03909216916110809190611d74565b6001600160401b031611156110d75760405162461bcd60e51b815260206004820152601b60248201527f5265736572766174696f6e20706572696f6420746f6f206c6f6e670000000000604482015260640161073d565b84604001516001600160401b031642116111e0576040805160608101825260018301546001600160401b038082168352600160401b820481166020808501829052600160801b9093048216948401949094529088015191929116148015611157575080604001516001600160401b031686604001516001600160401b0316115b6111a35760405162461bcd60e51b815260206004820152601a60248201527f496e76616c6964207265736572766174696f6e20757064617465000000000000604482015260640161073d565b805186516001600160401b03918216911610156111d25760405162461bcd60e51b815260040161073d90611c74565b6040015192506111e9915050565b50505060208201515b949350505050565b6111fb8185611d4e565b6001600160401b0316156112215760405162461bcd60e51b815260040161073d90611c0e565b61122b8184611d4e565b6001600160401b0316156112515760405162461bcd60e51b815260040161073d90611c45565b600061125d8286611d9c565b9050600061126b8386611d9c565b9050600061127761089e565b6001600160401b03808a16600090815260209290925260409091206002810154909250600160a01b900416835b836001600160401b0316816001600160401b0316101561138f576001600160401b03808216600090815260058501602052604081205490916112e8918a9116611dc2565b9050826001600160401b0316816001600160401b0316111561134c5760405162461bcd60e51b815260206004820152601c60248201527f4e6f7420656e6f7567682073796d626f6c7320617661696c61626c6500000000604482015260640161073d565b6001600160401b038083166000908152600586016020526040902080549190921667ffffffffffffffff199091161790558061138781611ded565b9150506112a4565b50505050505050505050565b604080516001600160a01b0385811660248301528416604482015260648082018490528251808303909101815260849091019091526020810180516001600160e01b03166323b872dd60e01b1790526105f49085906114a4565b60006108a8611576565b6000611414836001600160a01b0384166115c0565b9392505050565b6000611414826114296113f5565b6000868152602091909152604090209061160f565b61145f8161144a6113f5565b60008581526020919091526040902090611631565b506040516001600160a01b0382169083907f155aaafb6329a2098580462df33ec4b7441b19729b9601c5fc17ae1cf99a8a5290600090a35050565b60006103a2825490565b60006114f9826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564815250856001600160a01b03166116469092919063ffffffff16565b80519091501561075957808060200190518101906115179190611e14565b6107595760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6044820152691bdd081cdd58d8d9595960b21b606482015260840161073d565b60008060ff60001b196001604051806040016040528060168152602001756163636573732e636f6e74726f6c2e73746f7261676560501b815250604051602001610efb9190611d05565b6000818152600183016020526040812054611607575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556103a2565b5060006103a2565b6001600160a01b03811660009081526001830160205260408120541515611414565b6000611414836001600160a01b038416611655565b60606111e98484600085611748565b6000818152600183016020526040812054801561173e576000611679600183611d21565b855490915060009061168d90600190611d21565b90508181146116f25760008660000182815481106116ad576116ad611e36565b90600052602060002001549050808760000184815481106116d0576116d0611e36565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061170357611703611e4c565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506103a2565b60009150506103a2565b6060824710156117a95760405162461bcd60e51b815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f6044820152651c8818d85b1b60d21b606482015260840161073d565b6001600160a01b0385163b6118005760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640161073d565b600080866001600160a01b0316858760405161181c9190611d05565b60006040518083038185875af1925050503d8060008114611859576040519150601f19603f3d011682016040523d82523d6000602084013e61185e565b606091505b509150915061186e828286611879565b979650505050505050565b60608315611888575081611414565b8251156118985782518084602001fd5b8160405162461bcd60e51b815260040161073d9190611e62565b80356001600160401b03811681146118c957600080fd5b919050565b80356001600160a01b03811681146118c957600080fd5b600080604083850312156118f857600080fd5b611901836118b2565b915061190f602084016118ce565b90509250929050565b60006060828403121561192a57600080fd5b604051606081018181106001600160401b038211171561195a57634e487b7160e01b600052604160045260246000fd5b604052905080611969836118b2565b8152611977602084016118b2565b6020820152611988604084016118b2565b60408201525092915050565b600080608083850312156119a757600080fd5b6119b0836118b2565b915061190f8460208501611918565b600080604083850312156119d257600080fd5b6119db836118b2565b915061190f602084016118b2565b6000602082840312156119fb57600080fd5b611414826118b2565b60008060408385031215611a1757600080fd5b611a20836118b2565b946020939093013593505050565b600060208284031215611a4057600080fd5b611414826118ce565b600060208284031215611a5b57600080fd5b604051602081018181106001600160401b0382111715611a8b57634e487b7160e01b600052604160045260246000fd5b604052905080611a9a836118b2565b905292915050565b60008060408385031215611ab557600080fd5b611abe836118b2565b915061190f8460208501611a49565b600080600060a08486031215611ae257600080fd5b611aeb846118b2565b9250611af9602085016118ce565b9150611b088560408601611918565b90509250925092565b60008082840360c0811215611b2557600080fd5b611b2e846118b2565b925060a0601f1982011215611b4257600080fd5b5060405160a081018181106001600160401b0382111715611b7357634e487b7160e01b600052604160045260246000fd5b604052611b82602085016118ce565b8152611b90604085016118ce565b6020820152611ba1606085016118b2565b6040820152611bb2608085016118b2565b6060820152611bc360a085016118b2565b6080820152809150509250929050565b600080600060608486031215611be857600080fd5b611bf1846118b2565b9250611bff602085016118ce565b9150611b088560408601611a49565b60208082526017908201527f496e76616c69642073746172742074696d657374616d70000000000000000000604082015260600190565b6020808252601590820152740496e76616c696420656e642074696d657374616d7605c1b604082015260600190565b6020808252601a908201527f496e76616c69642073796d626f6c7320706572207365636f6e64000000000000604082015260600190565b634e487b7160e01b600052601160045260246000fd5b60008219821115611cd457611cd4611cab565b500190565b60005b83811015611cf4578181015183820152602001611cdc565b838111156105f45750506000910152565b60008251611d17818460208701611cd9565b9190910192915050565b600082821015611d3357611d33611cab565b500390565b634e487b7160e01b600052601260045260246000fd5b60006001600160401b0380841680611d6857611d68611d38565b92169190910692915050565b60006001600160401b0383811690831681811015611d9457611d94611cab565b039392505050565b60006001600160401b0380841680611db657611db6611d38565b92169190910492915050565b60006001600160401b03808316818516808303821115611de457611de4611cab565b01949350505050565b60006001600160401b0380831681811415611e0a57611e0a611cab565b6001019392505050565b600060208284031215611e2657600080fd5b8151801515811461141457600080fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052603160045260246000fd5b6020815260008251806020840152611e81816040850160208701611cd9565b601f01601f1916919091016040019291505056fea26469706673582212203ed84d521ecc95ecff2f31c55999099a2bef23716302bc6f9e366ba4286eb63064736f6c634300080c0033",
}

// ContractPaymentVaultABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractPaymentVaultMetaData.ABI instead.
var ContractPaymentVaultABI = ContractPaymentVaultMetaData.ABI

// ContractPaymentVaultBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractPaymentVaultMetaData.Bin instead.
var ContractPaymentVaultBin = ContractPaymentVaultMetaData.Bin

// DeployContractPaymentVault deploys a new Ethereum contract, binding an instance of ContractPaymentVault to it.
func DeployContractPaymentVault(auth *bind.TransactOpts, backend bind.ContractBackend, schedulePeriod uint64) (common.Address, *types.Transaction, *ContractPaymentVault, error) {
	parsed, err := ContractPaymentVaultMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractPaymentVaultBin), backend, schedulePeriod)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractPaymentVault{ContractPaymentVaultCaller: ContractPaymentVaultCaller{contract: contract}, ContractPaymentVaultTransactor: ContractPaymentVaultTransactor{contract: contract}, ContractPaymentVaultFilterer: ContractPaymentVaultFilterer{contract: contract}}, nil
}

// ContractPaymentVault is an auto generated Go binding around an Ethereum contract.
type ContractPaymentVault struct {
	ContractPaymentVaultCaller     // Read-only binding to the contract
	ContractPaymentVaultTransactor // Write-only binding to the contract
	ContractPaymentVaultFilterer   // Log filterer for contract events
}

// ContractPaymentVaultCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractPaymentVaultCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractPaymentVaultTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractPaymentVaultTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractPaymentVaultFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractPaymentVaultFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractPaymentVaultSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractPaymentVaultSession struct {
	Contract     *ContractPaymentVault // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ContractPaymentVaultCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractPaymentVaultCallerSession struct {
	Contract *ContractPaymentVaultCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// ContractPaymentVaultTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractPaymentVaultTransactorSession struct {
	Contract     *ContractPaymentVaultTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// ContractPaymentVaultRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractPaymentVaultRaw struct {
	Contract *ContractPaymentVault // Generic contract binding to access the raw methods on
}

// ContractPaymentVaultCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractPaymentVaultCallerRaw struct {
	Contract *ContractPaymentVaultCaller // Generic read-only contract binding to access the raw methods on
}

// ContractPaymentVaultTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractPaymentVaultTransactorRaw struct {
	Contract *ContractPaymentVaultTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractPaymentVault creates a new instance of ContractPaymentVault, bound to a specific deployed contract.
func NewContractPaymentVault(address common.Address, backend bind.ContractBackend) (*ContractPaymentVault, error) {
	contract, err := bindContractPaymentVault(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVault{ContractPaymentVaultCaller: ContractPaymentVaultCaller{contract: contract}, ContractPaymentVaultTransactor: ContractPaymentVaultTransactor{contract: contract}, ContractPaymentVaultFilterer: ContractPaymentVaultFilterer{contract: contract}}, nil
}

// NewContractPaymentVaultCaller creates a new read-only instance of ContractPaymentVault, bound to a specific deployed contract.
func NewContractPaymentVaultCaller(address common.Address, caller bind.ContractCaller) (*ContractPaymentVaultCaller, error) {
	contract, err := bindContractPaymentVault(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultCaller{contract: contract}, nil
}

// NewContractPaymentVaultTransactor creates a new write-only instance of ContractPaymentVault, bound to a specific deployed contract.
func NewContractPaymentVaultTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractPaymentVaultTransactor, error) {
	contract, err := bindContractPaymentVault(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultTransactor{contract: contract}, nil
}

// NewContractPaymentVaultFilterer creates a new log filterer instance of ContractPaymentVault, bound to a specific deployed contract.
func NewContractPaymentVaultFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractPaymentVaultFilterer, error) {
	contract, err := bindContractPaymentVault(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultFilterer{contract: contract}, nil
}

// bindContractPaymentVault binds a generic wrapper to an already deployed contract.
func bindContractPaymentVault(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractPaymentVaultMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractPaymentVault *ContractPaymentVaultRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractPaymentVault.Contract.ContractPaymentVaultCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractPaymentVault *ContractPaymentVaultRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.ContractPaymentVaultTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractPaymentVault *ContractPaymentVaultRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.ContractPaymentVaultTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractPaymentVault *ContractPaymentVaultCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractPaymentVault.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractPaymentVault *ContractPaymentVaultTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractPaymentVault *ContractPaymentVaultTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.contract.Transact(opts, method, params...)
}

// SCHEDULEPERIOD is a free data retrieval call binding the contract method 0xc1eafe51.
//
// Solidity: function SCHEDULE_PERIOD() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCaller) SCHEDULEPERIOD(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "SCHEDULE_PERIOD")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// SCHEDULEPERIOD is a free data retrieval call binding the contract method 0xc1eafe51.
//
// Solidity: function SCHEDULE_PERIOD() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultSession) SCHEDULEPERIOD() (uint64, error) {
	return _ContractPaymentVault.Contract.SCHEDULEPERIOD(&_ContractPaymentVault.CallOpts)
}

// SCHEDULEPERIOD is a free data retrieval call binding the contract method 0xc1eafe51.
//
// Solidity: function SCHEDULE_PERIOD() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) SCHEDULEPERIOD() (uint64, error) {
	return _ContractPaymentVault.Contract.SCHEDULEPERIOD(&_ContractPaymentVault.CallOpts)
}

// GetOnDemandDeposit is a free data retrieval call binding the contract method 0x400c2322.
//
// Solidity: function getOnDemandDeposit(uint64 quorumId, address account) view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultCaller) GetOnDemandDeposit(opts *bind.CallOpts, quorumId uint64, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "getOnDemandDeposit", quorumId, account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetOnDemandDeposit is a free data retrieval call binding the contract method 0x400c2322.
//
// Solidity: function getOnDemandDeposit(uint64 quorumId, address account) view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultSession) GetOnDemandDeposit(quorumId uint64, account common.Address) (*big.Int, error) {
	return _ContractPaymentVault.Contract.GetOnDemandDeposit(&_ContractPaymentVault.CallOpts, quorumId, account)
}

// GetOnDemandDeposit is a free data retrieval call binding the contract method 0x400c2322.
//
// Solidity: function getOnDemandDeposit(uint64 quorumId, address account) view returns(uint256)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) GetOnDemandDeposit(quorumId uint64, account common.Address) (*big.Int, error) {
	return _ContractPaymentVault.Contract.GetOnDemandDeposit(&_ContractPaymentVault.CallOpts, quorumId, account)
}

// GetQuorumPaymentConfig is a free data retrieval call binding the contract method 0x7a9426ca.
//
// Solidity: function getQuorumPaymentConfig(uint64 quorumId) view returns((address,address,uint64,uint64,uint64))
func (_ContractPaymentVault *ContractPaymentVaultCaller) GetQuorumPaymentConfig(opts *bind.CallOpts, quorumId uint64) (EigenDATypesV3QuorumPaymentConfig, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "getQuorumPaymentConfig", quorumId)

	if err != nil {
		return *new(EigenDATypesV3QuorumPaymentConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(EigenDATypesV3QuorumPaymentConfig)).(*EigenDATypesV3QuorumPaymentConfig)

	return out0, err

}

// GetQuorumPaymentConfig is a free data retrieval call binding the contract method 0x7a9426ca.
//
// Solidity: function getQuorumPaymentConfig(uint64 quorumId) view returns((address,address,uint64,uint64,uint64))
func (_ContractPaymentVault *ContractPaymentVaultSession) GetQuorumPaymentConfig(quorumId uint64) (EigenDATypesV3QuorumPaymentConfig, error) {
	return _ContractPaymentVault.Contract.GetQuorumPaymentConfig(&_ContractPaymentVault.CallOpts, quorumId)
}

// GetQuorumPaymentConfig is a free data retrieval call binding the contract method 0x7a9426ca.
//
// Solidity: function getQuorumPaymentConfig(uint64 quorumId) view returns((address,address,uint64,uint64,uint64))
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) GetQuorumPaymentConfig(quorumId uint64) (EigenDATypesV3QuorumPaymentConfig, error) {
	return _ContractPaymentVault.Contract.GetQuorumPaymentConfig(&_ContractPaymentVault.CallOpts, quorumId)
}

// GetQuorumProtocolConfig is a free data retrieval call binding the contract method 0x89a06b35.
//
// Solidity: function getQuorumProtocolConfig(uint64 quorumId) view returns((uint64))
func (_ContractPaymentVault *ContractPaymentVaultCaller) GetQuorumProtocolConfig(opts *bind.CallOpts, quorumId uint64) (EigenDATypesV3QuorumPaymentProtocolConfig, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "getQuorumProtocolConfig", quorumId)

	if err != nil {
		return *new(EigenDATypesV3QuorumPaymentProtocolConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(EigenDATypesV3QuorumPaymentProtocolConfig)).(*EigenDATypesV3QuorumPaymentProtocolConfig)

	return out0, err

}

// GetQuorumProtocolConfig is a free data retrieval call binding the contract method 0x89a06b35.
//
// Solidity: function getQuorumProtocolConfig(uint64 quorumId) view returns((uint64))
func (_ContractPaymentVault *ContractPaymentVaultSession) GetQuorumProtocolConfig(quorumId uint64) (EigenDATypesV3QuorumPaymentProtocolConfig, error) {
	return _ContractPaymentVault.Contract.GetQuorumProtocolConfig(&_ContractPaymentVault.CallOpts, quorumId)
}

// GetQuorumProtocolConfig is a free data retrieval call binding the contract method 0x89a06b35.
//
// Solidity: function getQuorumProtocolConfig(uint64 quorumId) view returns((uint64))
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) GetQuorumProtocolConfig(quorumId uint64) (EigenDATypesV3QuorumPaymentProtocolConfig, error) {
	return _ContractPaymentVault.Contract.GetQuorumProtocolConfig(&_ContractPaymentVault.CallOpts, quorumId)
}

// GetQuorumReservedSymbols is a free data retrieval call binding the contract method 0x4023a200.
//
// Solidity: function getQuorumReservedSymbols(uint64 quorumId, uint64 period) view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCaller) GetQuorumReservedSymbols(opts *bind.CallOpts, quorumId uint64, period uint64) (uint64, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "getQuorumReservedSymbols", quorumId, period)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetQuorumReservedSymbols is a free data retrieval call binding the contract method 0x4023a200.
//
// Solidity: function getQuorumReservedSymbols(uint64 quorumId, uint64 period) view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultSession) GetQuorumReservedSymbols(quorumId uint64, period uint64) (uint64, error) {
	return _ContractPaymentVault.Contract.GetQuorumReservedSymbols(&_ContractPaymentVault.CallOpts, quorumId, period)
}

// GetQuorumReservedSymbols is a free data retrieval call binding the contract method 0x4023a200.
//
// Solidity: function getQuorumReservedSymbols(uint64 quorumId, uint64 period) view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) GetQuorumReservedSymbols(quorumId uint64, period uint64) (uint64, error) {
	return _ContractPaymentVault.Contract.GetQuorumReservedSymbols(&_ContractPaymentVault.CallOpts, quorumId, period)
}

// GetReservation is a free data retrieval call binding the contract method 0x00e691aa.
//
// Solidity: function getReservation(uint64 quorumId, address account) view returns((uint64,uint64,uint64))
func (_ContractPaymentVault *ContractPaymentVaultCaller) GetReservation(opts *bind.CallOpts, quorumId uint64, account common.Address) (EigenDATypesV3Reservation, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "getReservation", quorumId, account)

	if err != nil {
		return *new(EigenDATypesV3Reservation), err
	}

	out0 := *abi.ConvertType(out[0], new(EigenDATypesV3Reservation)).(*EigenDATypesV3Reservation)

	return out0, err

}

// GetReservation is a free data retrieval call binding the contract method 0x00e691aa.
//
// Solidity: function getReservation(uint64 quorumId, address account) view returns((uint64,uint64,uint64))
func (_ContractPaymentVault *ContractPaymentVaultSession) GetReservation(quorumId uint64, account common.Address) (EigenDATypesV3Reservation, error) {
	return _ContractPaymentVault.Contract.GetReservation(&_ContractPaymentVault.CallOpts, quorumId, account)
}

// GetReservation is a free data retrieval call binding the contract method 0x00e691aa.
//
// Solidity: function getReservation(uint64 quorumId, address account) view returns((uint64,uint64,uint64))
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) GetReservation(quorumId uint64, account common.Address) (EigenDATypesV3Reservation, error) {
	return _ContractPaymentVault.Contract.GetReservation(&_ContractPaymentVault.CallOpts, quorumId, account)
}

// CreateReservation is a paid mutator transaction binding the contract method 0xecfc2d8e.
//
// Solidity: function createReservation(uint64 quorumId, address account, (uint64,uint64,uint64) reservation) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) CreateReservation(opts *bind.TransactOpts, quorumId uint64, account common.Address, reservation EigenDATypesV3Reservation) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "createReservation", quorumId, account, reservation)
}

// CreateReservation is a paid mutator transaction binding the contract method 0xecfc2d8e.
//
// Solidity: function createReservation(uint64 quorumId, address account, (uint64,uint64,uint64) reservation) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) CreateReservation(quorumId uint64, account common.Address, reservation EigenDATypesV3Reservation) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.CreateReservation(&_ContractPaymentVault.TransactOpts, quorumId, account, reservation)
}

// CreateReservation is a paid mutator transaction binding the contract method 0xecfc2d8e.
//
// Solidity: function createReservation(uint64 quorumId, address account, (uint64,uint64,uint64) reservation) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) CreateReservation(quorumId uint64, account common.Address, reservation EigenDATypesV3Reservation) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.CreateReservation(&_ContractPaymentVault.TransactOpts, quorumId, account, reservation)
}

// DecreaseReservation is a paid mutator transaction binding the contract method 0x3b50e8b8.
//
// Solidity: function decreaseReservation(uint64 quorumId, (uint64,uint64,uint64) reservation) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) DecreaseReservation(opts *bind.TransactOpts, quorumId uint64, reservation EigenDATypesV3Reservation) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "decreaseReservation", quorumId, reservation)
}

// DecreaseReservation is a paid mutator transaction binding the contract method 0x3b50e8b8.
//
// Solidity: function decreaseReservation(uint64 quorumId, (uint64,uint64,uint64) reservation) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) DecreaseReservation(quorumId uint64, reservation EigenDATypesV3Reservation) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.DecreaseReservation(&_ContractPaymentVault.TransactOpts, quorumId, reservation)
}

// DecreaseReservation is a paid mutator transaction binding the contract method 0x3b50e8b8.
//
// Solidity: function decreaseReservation(uint64 quorumId, (uint64,uint64,uint64) reservation) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) DecreaseReservation(quorumId uint64, reservation EigenDATypesV3Reservation) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.DecreaseReservation(&_ContractPaymentVault.TransactOpts, quorumId, reservation)
}

// DepositOnDemand is a paid mutator transaction binding the contract method 0xac990b8a.
//
// Solidity: function depositOnDemand(uint64 quorumId, uint256 amount) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) DepositOnDemand(opts *bind.TransactOpts, quorumId uint64, amount *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "depositOnDemand", quorumId, amount)
}

// DepositOnDemand is a paid mutator transaction binding the contract method 0xac990b8a.
//
// Solidity: function depositOnDemand(uint64 quorumId, uint256 amount) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) DepositOnDemand(quorumId uint64, amount *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.DepositOnDemand(&_ContractPaymentVault.TransactOpts, quorumId, amount)
}

// DepositOnDemand is a paid mutator transaction binding the contract method 0xac990b8a.
//
// Solidity: function depositOnDemand(uint64 quorumId, uint256 amount) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) DepositOnDemand(quorumId uint64, amount *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.DepositOnDemand(&_ContractPaymentVault.TransactOpts, quorumId, amount)
}

// IncreaseReservation is a paid mutator transaction binding the contract method 0xe453a058.
//
// Solidity: function increaseReservation(uint64 quorumId, address account, (uint64,uint64,uint64) reservation) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) IncreaseReservation(opts *bind.TransactOpts, quorumId uint64, account common.Address, reservation EigenDATypesV3Reservation) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "increaseReservation", quorumId, account, reservation)
}

// IncreaseReservation is a paid mutator transaction binding the contract method 0xe453a058.
//
// Solidity: function increaseReservation(uint64 quorumId, address account, (uint64,uint64,uint64) reservation) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) IncreaseReservation(quorumId uint64, account common.Address, reservation EigenDATypesV3Reservation) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.IncreaseReservation(&_ContractPaymentVault.TransactOpts, quorumId, account, reservation)
}

// IncreaseReservation is a paid mutator transaction binding the contract method 0xe453a058.
//
// Solidity: function increaseReservation(uint64 quorumId, address account, (uint64,uint64,uint64) reservation) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) IncreaseReservation(quorumId uint64, account common.Address, reservation EigenDATypesV3Reservation) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.IncreaseReservation(&_ContractPaymentVault.TransactOpts, quorumId, account, reservation)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address owner) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) Initialize(opts *bind.TransactOpts, owner common.Address) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "initialize", owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address owner) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) Initialize(owner common.Address) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.Initialize(&_ContractPaymentVault.TransactOpts, owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address owner) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) Initialize(owner common.Address) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.Initialize(&_ContractPaymentVault.TransactOpts, owner)
}

// InitializeQuorum is a paid mutator transaction binding the contract method 0xef7531bd.
//
// Solidity: function initializeQuorum(uint64 quorumId, address newOwner, (uint64) protocolCfg) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) InitializeQuorum(opts *bind.TransactOpts, quorumId uint64, newOwner common.Address, protocolCfg EigenDATypesV3QuorumPaymentProtocolConfig) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "initializeQuorum", quorumId, newOwner, protocolCfg)
}

// InitializeQuorum is a paid mutator transaction binding the contract method 0xef7531bd.
//
// Solidity: function initializeQuorum(uint64 quorumId, address newOwner, (uint64) protocolCfg) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) InitializeQuorum(quorumId uint64, newOwner common.Address, protocolCfg EigenDATypesV3QuorumPaymentProtocolConfig) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.InitializeQuorum(&_ContractPaymentVault.TransactOpts, quorumId, newOwner, protocolCfg)
}

// InitializeQuorum is a paid mutator transaction binding the contract method 0xef7531bd.
//
// Solidity: function initializeQuorum(uint64 quorumId, address newOwner, (uint64) protocolCfg) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) InitializeQuorum(quorumId uint64, newOwner common.Address, protocolCfg EigenDATypesV3QuorumPaymentProtocolConfig) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.InitializeQuorum(&_ContractPaymentVault.TransactOpts, quorumId, newOwner, protocolCfg)
}

// SetQuorumPaymentConfig is a paid mutator transaction binding the contract method 0xe631f82c.
//
// Solidity: function setQuorumPaymentConfig(uint64 quorumId, (address,address,uint64,uint64,uint64) paymentConfig) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) SetQuorumPaymentConfig(opts *bind.TransactOpts, quorumId uint64, paymentConfig EigenDATypesV3QuorumPaymentConfig) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "setQuorumPaymentConfig", quorumId, paymentConfig)
}

// SetQuorumPaymentConfig is a paid mutator transaction binding the contract method 0xe631f82c.
//
// Solidity: function setQuorumPaymentConfig(uint64 quorumId, (address,address,uint64,uint64,uint64) paymentConfig) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) SetQuorumPaymentConfig(quorumId uint64, paymentConfig EigenDATypesV3QuorumPaymentConfig) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetQuorumPaymentConfig(&_ContractPaymentVault.TransactOpts, quorumId, paymentConfig)
}

// SetQuorumPaymentConfig is a paid mutator transaction binding the contract method 0xe631f82c.
//
// Solidity: function setQuorumPaymentConfig(uint64 quorumId, (address,address,uint64,uint64,uint64) paymentConfig) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) SetQuorumPaymentConfig(quorumId uint64, paymentConfig EigenDATypesV3QuorumPaymentConfig) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetQuorumPaymentConfig(&_ContractPaymentVault.TransactOpts, quorumId, paymentConfig)
}

// SetReservationAdvanceWindow is a paid mutator transaction binding the contract method 0xdc5fee25.
//
// Solidity: function setReservationAdvanceWindow(uint64 quorumId, (uint64) protocolCfg) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) SetReservationAdvanceWindow(opts *bind.TransactOpts, quorumId uint64, protocolCfg EigenDATypesV3QuorumPaymentProtocolConfig) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "setReservationAdvanceWindow", quorumId, protocolCfg)
}

// SetReservationAdvanceWindow is a paid mutator transaction binding the contract method 0xdc5fee25.
//
// Solidity: function setReservationAdvanceWindow(uint64 quorumId, (uint64) protocolCfg) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) SetReservationAdvanceWindow(quorumId uint64, protocolCfg EigenDATypesV3QuorumPaymentProtocolConfig) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetReservationAdvanceWindow(&_ContractPaymentVault.TransactOpts, quorumId, protocolCfg)
}

// SetReservationAdvanceWindow is a paid mutator transaction binding the contract method 0xdc5fee25.
//
// Solidity: function setReservationAdvanceWindow(uint64 quorumId, (uint64) protocolCfg) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) SetReservationAdvanceWindow(quorumId uint64, protocolCfg EigenDATypesV3QuorumPaymentProtocolConfig) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetReservationAdvanceWindow(&_ContractPaymentVault.TransactOpts, quorumId, protocolCfg)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.TransferOwnership(&_ContractPaymentVault.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.TransferOwnership(&_ContractPaymentVault.TransactOpts, newOwner)
}

// TransferQuorumOwnership is a paid mutator transaction binding the contract method 0xede57a2d.
//
// Solidity: function transferQuorumOwnership(uint64 quorumId, address newOwner) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) TransferQuorumOwnership(opts *bind.TransactOpts, quorumId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "transferQuorumOwnership", quorumId, newOwner)
}

// TransferQuorumOwnership is a paid mutator transaction binding the contract method 0xede57a2d.
//
// Solidity: function transferQuorumOwnership(uint64 quorumId, address newOwner) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) TransferQuorumOwnership(quorumId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.TransferQuorumOwnership(&_ContractPaymentVault.TransactOpts, quorumId, newOwner)
}

// TransferQuorumOwnership is a paid mutator transaction binding the contract method 0xede57a2d.
//
// Solidity: function transferQuorumOwnership(uint64 quorumId, address newOwner) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) TransferQuorumOwnership(quorumId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.TransferQuorumOwnership(&_ContractPaymentVault.TransactOpts, quorumId, newOwner)
}
