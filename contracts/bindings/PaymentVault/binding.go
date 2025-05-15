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

// PaymentVaultTypesQuorumConfig is an auto generated low-level Go binding around an user-defined struct.
type PaymentVaultTypesQuorumConfig struct {
	Token                       common.Address
	Recipient                   common.Address
	ReservationSymbolsPerSecond uint64
	OnDemandSymbolsPerPeriod    uint64
	OnDemandPricePerSymbol      uint64
}

// PaymentVaultTypesQuorumProtocolConfig is an auto generated low-level Go binding around an user-defined struct.
type PaymentVaultTypesQuorumProtocolConfig struct {
	MinNumSymbols              uint64
	ReservationAdvanceWindow   uint64
	ReservationRateLimitWindow uint64
	OnDemandRateLimitWindow    uint64
	OnDemandEnabled            bool
}

// PaymentVaultTypesReservation is an auto generated low-level Go binding around an user-defined struct.
type PaymentVaultTypesReservation struct {
	SymbolsPerSecond uint64
	StartTimestamp   uint64
	EndTimestamp     uint64
}

// ContractPaymentVaultMetaData contains all meta data concerning the ContractPaymentVault contract.
var ContractPaymentVaultMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"schedulePeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"SCHEDULE_PERIOD\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"createReservation\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reservation\",\"type\":\"tuple\",\"internalType\":\"structPaymentVaultTypes.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"decreaseReservation\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"reservation\",\"type\":\"tuple\",\"internalType\":\"structPaymentVaultTypes.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositOnDemand\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getOnDemandDeposit\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumPaymentConfig\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPaymentVaultTypes.QuorumConfig\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reservationSymbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandSymbolsPerPeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandPricePerSymbol\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumProtocolConfig\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPaymentVaultTypes.QuorumProtocolConfig\",\"components\":[{\"name\":\"minNumSymbols\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"reservationAdvanceWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"reservationRateLimitWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandRateLimitWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumReservedSymbols\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"period\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReservation\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPaymentVaultTypes.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"increaseReservation\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reservation\",\"type\":\"tuple\",\"internalType\":\"structPaymentVaultTypes.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initializeQuorum\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"protocolCfg\",\"type\":\"tuple\",\"internalType\":\"structPaymentVaultTypes.QuorumProtocolConfig\",\"components\":[{\"name\":\"minNumSymbols\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"reservationAdvanceWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"reservationRateLimitWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandRateLimitWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setOnDemandEnabled\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"protocolCfg\",\"type\":\"tuple\",\"internalType\":\"structPaymentVaultTypes.QuorumProtocolConfig\",\"components\":[{\"name\":\"minNumSymbols\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"reservationAdvanceWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"reservationRateLimitWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandRateLimitWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setQuorumPaymentConfig\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"paymentConfig\",\"type\":\"tuple\",\"internalType\":\"structPaymentVaultTypes.QuorumConfig\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reservationSymbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandSymbolsPerPeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandPricePerSymbol\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setReservationAdvanceWindow\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"protocolCfg\",\"type\":\"tuple\",\"internalType\":\"structPaymentVaultTypes.QuorumProtocolConfig\",\"components\":[{\"name\":\"minNumSymbols\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"reservationAdvanceWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"reservationRateLimitWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandRateLimitWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferQuorumOwnership\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162002197380380620021978339810160408190526200003491620000b3565b6000816001600160401b031611620000a15760405162461bcd60e51b815260206004820152602660248201527f5363686564756c6520706572696f64206d75737420626520677265617465722060448201526507468616e20360d41b606482015260840160405180910390fd5b6001600160401b0316608052620000e5565b600060208284031215620000c657600080fd5b81516001600160401b0381168114620000de57600080fd5b9392505050565b60805161208162000116600039600081816102e0015281816105900152818161080e01526108fb01526120816000f3fe608060405234801561001057600080fd5b506004361061010a5760003560e01c806389a06b35116100a2578063e453a05811610071578063e453a05814610315578063e631f82c14610328578063ecfc2d8e1461033b578063ede57a2d1461034e578063f2fde38b1461036157600080fd5b806389a06b3514610260578063ac990b8a146102c8578063c1eafe51146102db578063c4d66de81461030257600080fd5b8063400c2322116100de578063400c2322146101945780634023a200146101b5578063710f48be146101e05780637a9426ca146101f357600080fd5b8062e691aa1461010f578063063a5d56146101595780631269a9d71461016e5780633b50e8b814610181575b600080fd5b61012261011d366004611a6d565b610374565b6040805182516001600160401b03908116825260208085015182169083015292820151909216908201526060015b60405180910390f35b61016c610167366004611b57565b610403565b005b61016c61017c366004611b9b565b61052e565b61016c61018f366004611c42565b610588565b6101a76101a2366004611a6d565b6105b8565b604051908152602001610150565b6101c86101c3366004611c6d565b6105f7565b6040516001600160401b039091168152602001610150565b61016c6101ee366004611b9b565b610634565b610206610201366004611c97565b61067e565b6040805182516001600160a01b03908116825260208085015190911690820152828201516001600160401b03908116928201929092526060808401518316908201526080928301519091169181019190915260a001610150565b61027361026e366004611c97565b610725565b6040516101509190600060a0820190506001600160401b038084511683528060208501511660208401528060408501511660408401528060608501511660608401525060808301511515608083015292915050565b61016c6102d6366004611cb2565b6107c4565b6101c87f000000000000000000000000000000000000000000000000000000000000000081565b61016c610310366004611cdc565b6107cf565b61016c610323366004611cf7565b6107fc565b61016c610336366004611d32565b610838565b61016c610349366004611cf7565b6108e9565b61016c61035c366004611a6d565b61091f565b61016c61036f366004611cdc565b610997565b6040805160608101825260008082526020820181905291810191909152610399610a20565b6001600160401b038085166000908152602092835260408082206001600160a01b038716835260050184529081902081516060810183526001909101548084168252600160401b8104841694820194909452600160801b9093049091169082015290505b92915050565b61040b610a2f565b61041c61041784610a93565b610ac8565b1561046e5760405162461bcd60e51b815260206004820152601860248201527f51756f72756d206f776e657220616c726561647920736574000000000000000060448201526064015b60405180910390fd5b61048061047a84610a93565b83610ae9565b80610489610a20565b6001600160401b039485166000908152602091825260409081902083518154938501519285015160608601518916600160c01b026001600160c01b03918a16600160801b02919091166fffffffffffffffffffffffffffffffff948a16600160401b026001600160801b03199096169290991691909117939093179190911695909517178455608001516001909301805493151560ff19909416939093179092555050565b8161053881610b45565b8160200151610545610a20565b6001600160401b03948516600090815260209190915260409020805491909416600160401b026fffffffffffffffff000000000000000019909116179092555050565b6105b48233837f0000000000000000000000000000000000000000000000000000000000000000610b96565b5050565b60006105c2610a20565b6001600160401b0384166000908152602091825260408082206001600160a01b03861683526005019092522054905092915050565b6000610601610a20565b6001600160401b038085166000908152602092835260408082208684168352600601909352919091205416905092915050565b8161063e81610b45565b816080015161064b610a20565b6001600160401b039490941660009081526020949094526040909320600101805460ff1916931515939093179092555050565b6040805160a0810182526000808252602082018190529181018290526060810182905260808101919091526106b1610a20565b6001600160401b0392831660009081526020918252604090819020815160a08101835260028201546001600160a01b039081168252600383015490811694820194909452600160a01b909304851691830191909152600401548084166060830152600160401b900490921660808301525090565b6040805160a081018252600080825260208201819052918101829052606081018290526080810191909152610758610a20565b6001600160401b0392831660009081526020918252604090819020815160a08101835281548087168252600160401b8104871694820194909452600160801b8404861692810192909252600160c01b90920490931660608401526001015460ff16151560808301525090565b6105b4823383610d2e565b6107f97fe579393920b888b1e4a7e1afdd7d58fa4f3101113547ac874aefa75ff4a960f982610ae9565b50565b8261080681610b45565b6108328484847f0000000000000000000000000000000000000000000000000000000000000000610df6565b50505050565b8161084281610b45565b8161084b610a20565b6001600160401b039485166000908152602091825260409081902083516002820180546001600160a01b039283166001600160a01b031990911617905592840151600382018054938601518916600160a01b026001600160e01b03199094169190941617919091179091556060820151600490910180546080909301518616600160401b026001600160801b03199093169190951617179092555050565b826108f381610b45565b6108328484847f0000000000000000000000000000000000000000000000000000000000000000610edd565b8161092981610b45565b6001600160a01b03821661097f5760405162461bcd60e51b815260206004820152601d60248201527f4e6577206f776e657220697320746865207a65726f20616464726573730000006044820152606401610465565b61099261098b84610a93565b338461101f565b505050565b61099f610a2f565b6001600160a01b0381166109f55760405162461bcd60e51b815260206004820152601d60248201527f4e6577206f776e657220697320746865207a65726f20616464726573730000006044820152606401610465565b6107f97fe579393920b888b1e4a7e1afdd7d58fa4f3101113547ac874aefa75ff4a960f9338361101f565b6000610a2a611033565b905090565b610a597fe579393920b888b1e4a7e1afdd7d58fa4f3101113547ac874aefa75ff4a960f9336110d1565b610a915760405162461bcd60e51b81526020600482015260096024820152682737ba1037bbb732b960b91b6044820152606401610465565b565b60006103fd6001600160401b0383167f9cb79c1d0fdfada3ea04142fe992963bc303c019dac5ad1fb95c78752893db12611ddf565b60006103fd610ad56110fb565b600084815260209190915260409020611105565b610b0a81610af56110fb565b6000858152602091909152604090209061110f565b506040516001600160a01b0382169083907f2ae6a113c0ed5b78a53413ffbb7679881f11145ccfba4fb92e863dfcd5a1d2f390600090a35050565b610b57610b5182610a93565b336110d1565b6107f95760405162461bcd60e51b815260206004820152601060248201526f2737ba1038bab7b93ab69037bbb732b960811b6044820152606401610465565b6000610ba0610a20565b6001600160401b0386166000908152602091825260408082206001600160a01b0388168352600501909252206001019050610bdd85858585611124565b50805460208401516001600160401b03908116600160401b9092041614610c165760405162461bcd60e51b815260040161046590611df7565b805460408401516001600160401b03600160801b909204821691161115610c4f5760405162461bcd60e51b815260040161046590611e2e565b805483516001600160401b0391821691161115610c7e5760405162461bcd60e51b815260040161046590611e5d565b805460408401518451610ca5928892600160801b9091046001600160401b031691866113cc565b82610cae610a20565b6001600160401b039687166000908152602091825260408082206001600160a01b039098168252600590970182528690208251600190910180549284015193909701518816600160801b0267ffffffffffffffff60801b19938916600160401b026001600160801b03199093169190981617171694909417909255505050565b6000610d38610a20565b6001600160401b0385166000908152602091825260408082206001600160a01b0387168352600581019093528120805492935091600284019190610d7d908690611ddf565b905069ffffffffffffffffffff811115610dcc5760405162461bcd60e51b815260206004820152601060248201526f416d6f756e7420746f6f206c6172676560801b6044820152606401610465565b60018201548254610dec916001600160a01b039182169189911688611576565b9091555050505050565b6000610e00610a20565b6001600160401b0386166000908152602091825260408082206001600160a01b0388168352600501909252206001019050610e3d85858585611124565b50805460208401516001600160401b03908116600160401b9092041614610e765760405162461bcd60e51b815260040161046590611df7565b805460408401516001600160401b03600160801b9092048216911611610eae5760405162461bcd60e51b815260040161046590611e2e565b805483516001600160401b0391821691161015610c7e5760405162461bcd60e51b815260040161046590611e5d565b610ee984848484611124565b50610ef2610a20565b6001600160401b038086166000908152602092835260408082206001600160a01b0388168352600501845290206001015491840151600160801b909204811691161015610f515760405162461bcd60e51b815260040161046590611df7565b4282602001516001600160401b03161015610f7e5760405162461bcd60e51b815260040161046590611df7565b610f9784836020015184604001518560000151856113cc565b81610fa0610a20565b6001600160401b039586166000908152602091825260408082206001600160a01b039097168252600590960182528590208251600190910180549284015193909601518716600160801b0267ffffffffffffffff60801b19938816600160401b026001600160801b031990931691909716171716939093179091555050565b61102983836115d0565b6109928382610ae9565b60008060ff60001b19600160405180604001604052806016815260200175195a59d95b8b99184b9c185e5b595b9d0b9d985d5b1d60521b81525060405160200161107d9190611ec0565b6040516020818303038152906040528051906020012060001c6110a09190611edc565b6040516020016110b291815260200190565b60408051601f1981840301815291905280516020909101201692915050565b60006110f4826110df6110fb565b6000868152602091909152604090209061162c565b9392505050565b6000610a2a61164e565b60006103fd825490565b60006110f4836001600160a01b038416611698565b60008061112f610a20565b6001600160401b0387166000908152602091825260408082206001600160a01b038916835260058101845291209186015190925061116e908590611f09565b6001600160401b0316156111945760405162461bcd60e51b815260040161046590611df7565b8385604001516111a49190611f09565b6001600160401b0316156111ca5760405162461bcd60e51b815260040161046590611e2e565b84602001516001600160401b031685604001516001600160401b0316116112335760405162461bcd60e51b815260206004820152601a60248201527f496e76616c6964207265736572766174696f6e20706572696f640000000000006044820152606401610465565b815460208601516040870151600160401b9092046001600160401b03169161125b9190611f2f565b6001600160401b031611156112b25760405162461bcd60e51b815260206004820152601b60248201527f5265736572766174696f6e20706572696f6420746f6f206c6f6e6700000000006044820152606401610465565b84604001516001600160401b031642116113bb576040805160608101825260018301546001600160401b038082168352600160401b820481166020808501829052600160801b9093048216948401949094529088015191929116148015611332575080604001516001600160401b031686604001516001600160401b0316115b61137e5760405162461bcd60e51b815260206004820152601a60248201527f496e76616c6964207265736572766174696f6e207570646174650000000000006044820152606401610465565b805186516001600160401b03918216911610156113ad5760405162461bcd60e51b815260040161046590611e5d565b6040015192506113c4915050565b50505060208201515b949350505050565b6113d68185611f09565b6001600160401b0316156113fc5760405162461bcd60e51b815260040161046590611df7565b6114068184611f09565b6001600160401b03161561142c5760405162461bcd60e51b815260040161046590611e2e565b60006114388286611f57565b905060006114468386611f57565b90506000611452610a20565b6001600160401b03808a16600090815260209290925260409091206003810154909250600160a01b900416835b836001600160401b0316816001600160401b0316101561156a576001600160401b03808216600090815260068501602052604081205490916114c3918a9116611f7d565b9050826001600160401b0316816001600160401b031611156115275760405162461bcd60e51b815260206004820152601c60248201527f4e6f7420656e6f7567682073796d626f6c7320617661696c61626c65000000006044820152606401610465565b6001600160401b038083166000908152600686016020526040902080549190921667ffffffffffffffff199091161790558061156281611fa8565b91505061147f565b50505050505050505050565b604080516001600160a01b0385811660248301528416604482015260648082018490528251808303909101815260849091019091526020810180516001600160e01b03166323b872dd60e01b1790526108329085906116e7565b6115f1816115dc6110fb565b600085815260209190915260409020906117b9565b506040516001600160a01b0382169083907f155aaafb6329a2098580462df33ec4b7441b19729b9601c5fc17ae1cf99a8a5290600090a35050565b6001600160a01b038116600090815260018301602052604081205415156110f4565b60008060ff60001b196001604051806040016040528060168152602001756163636573732e636f6e74726f6c2e73746f7261676560501b81525060405160200161107d9190611ec0565b60008181526001830160205260408120546116df575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556103fd565b5060006103fd565b600061173c826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564815250856001600160a01b03166117ce9092919063ffffffff16565b805190915015610992578080602001905181019061175a9190611fcf565b6109925760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6044820152691bdd081cdd58d8d9595960b21b6064820152608401610465565b60006110f4836001600160a01b0384166117dd565b60606113c484846000856118d0565b600081815260018301602052604081205480156118c6576000611801600183611edc565b855490915060009061181590600190611edc565b905081811461187a57600086600001828154811061183557611835611fec565b906000526020600020015490508087600001848154811061185857611858611fec565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061188b5761188b612002565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506103fd565b60009150506103fd565b6060824710156119315760405162461bcd60e51b815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f6044820152651c8818d85b1b60d21b6064820152608401610465565b6001600160a01b0385163b6119885760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610465565b600080866001600160a01b031685876040516119a49190611ec0565b60006040518083038185875af1925050503d80600081146119e1576040519150601f19603f3d011682016040523d82523d6000602084013e6119e6565b606091505b50915091506119f6828286611a01565b979650505050505050565b60608315611a105750816110f4565b825115611a205782518084602001fd5b8160405162461bcd60e51b81526004016104659190612018565b80356001600160401b0381168114611a5157600080fd5b919050565b80356001600160a01b0381168114611a5157600080fd5b60008060408385031215611a8057600080fd5b611a8983611a3a565b9150611a9760208401611a56565b90509250929050565b60405160a081016001600160401b0381118282101715611ad057634e487b7160e01b600052604160045260246000fd5b60405290565b80151581146107f957600080fd5b600060a08284031215611af657600080fd5b611afe611aa0565b9050611b0982611a3a565b8152611b1760208301611a3a565b6020820152611b2860408301611a3a565b6040820152611b3960608301611a3a565b60608201526080820135611b4c81611ad6565b608082015292915050565b600080600060e08486031215611b6c57600080fd5b611b7584611a3a565b9250611b8360208501611a56565b9150611b928560408601611ae4565b90509250925092565b60008060c08385031215611bae57600080fd5b611bb783611a3a565b9150611a978460208501611ae4565b600060608284031215611bd857600080fd5b604051606081018181106001600160401b0382111715611c0857634e487b7160e01b600052604160045260246000fd5b604052905080611c1783611a3a565b8152611c2560208401611a3a565b6020820152611c3660408401611a3a565b60408201525092915050565b60008060808385031215611c5557600080fd5b611c5e83611a3a565b9150611a978460208501611bc6565b60008060408385031215611c8057600080fd5b611c8983611a3a565b9150611a9760208401611a3a565b600060208284031215611ca957600080fd5b6110f482611a3a565b60008060408385031215611cc557600080fd5b611cce83611a3a565b946020939093013593505050565b600060208284031215611cee57600080fd5b6110f482611a56565b600080600060a08486031215611d0c57600080fd5b611d1584611a3a565b9250611d2360208501611a56565b9150611b928560408601611bc6565b60008082840360c0811215611d4657600080fd5b611d4f84611a3a565b925060a0601f1982011215611d6357600080fd5b50611d6c611aa0565b611d7860208501611a56565b8152611d8660408501611a56565b6020820152611d9760608501611a3a565b6040820152611da860808501611a3a565b6060820152611db960a08501611a3a565b6080820152809150509250929050565b634e487b7160e01b600052601160045260246000fd5b60008219821115611df257611df2611dc9565b500190565b60208082526017908201527f496e76616c69642073746172742074696d657374616d70000000000000000000604082015260600190565b6020808252601590820152740496e76616c696420656e642074696d657374616d7605c1b604082015260600190565b6020808252601a908201527f496e76616c69642073796d626f6c7320706572207365636f6e64000000000000604082015260600190565b60005b83811015611eaf578181015183820152602001611e97565b838111156108325750506000910152565b60008251611ed2818460208701611e94565b9190910192915050565b600082821015611eee57611eee611dc9565b500390565b634e487b7160e01b600052601260045260246000fd5b60006001600160401b0380841680611f2357611f23611ef3565b92169190910692915050565b60006001600160401b0383811690831681811015611f4f57611f4f611dc9565b039392505050565b60006001600160401b0380841680611f7157611f71611ef3565b92169190910492915050565b60006001600160401b03808316818516808303821115611f9f57611f9f611dc9565b01949350505050565b60006001600160401b0380831681811415611fc557611fc5611dc9565b6001019392505050565b600060208284031215611fe157600080fd5b81516110f481611ad6565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052603160045260246000fd5b6020815260008251806020840152612037816040850160208701611e94565b601f01601f1916919091016040019291505056fea26469706673582212206630f6c411c697d44dec4016f82e31e7707df3372e7f7ac5322e5fe165c032dc64736f6c634300080c0033",
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
func (_ContractPaymentVault *ContractPaymentVaultCaller) GetQuorumPaymentConfig(opts *bind.CallOpts, quorumId uint64) (PaymentVaultTypesQuorumConfig, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "getQuorumPaymentConfig", quorumId)

	if err != nil {
		return *new(PaymentVaultTypesQuorumConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(PaymentVaultTypesQuorumConfig)).(*PaymentVaultTypesQuorumConfig)

	return out0, err

}

// GetQuorumPaymentConfig is a free data retrieval call binding the contract method 0x7a9426ca.
//
// Solidity: function getQuorumPaymentConfig(uint64 quorumId) view returns((address,address,uint64,uint64,uint64))
func (_ContractPaymentVault *ContractPaymentVaultSession) GetQuorumPaymentConfig(quorumId uint64) (PaymentVaultTypesQuorumConfig, error) {
	return _ContractPaymentVault.Contract.GetQuorumPaymentConfig(&_ContractPaymentVault.CallOpts, quorumId)
}

// GetQuorumPaymentConfig is a free data retrieval call binding the contract method 0x7a9426ca.
//
// Solidity: function getQuorumPaymentConfig(uint64 quorumId) view returns((address,address,uint64,uint64,uint64))
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) GetQuorumPaymentConfig(quorumId uint64) (PaymentVaultTypesQuorumConfig, error) {
	return _ContractPaymentVault.Contract.GetQuorumPaymentConfig(&_ContractPaymentVault.CallOpts, quorumId)
}

// GetQuorumProtocolConfig is a free data retrieval call binding the contract method 0x89a06b35.
//
// Solidity: function getQuorumProtocolConfig(uint64 quorumId) view returns((uint64,uint64,uint64,uint64,bool))
func (_ContractPaymentVault *ContractPaymentVaultCaller) GetQuorumProtocolConfig(opts *bind.CallOpts, quorumId uint64) (PaymentVaultTypesQuorumProtocolConfig, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "getQuorumProtocolConfig", quorumId)

	if err != nil {
		return *new(PaymentVaultTypesQuorumProtocolConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(PaymentVaultTypesQuorumProtocolConfig)).(*PaymentVaultTypesQuorumProtocolConfig)

	return out0, err

}

// GetQuorumProtocolConfig is a free data retrieval call binding the contract method 0x89a06b35.
//
// Solidity: function getQuorumProtocolConfig(uint64 quorumId) view returns((uint64,uint64,uint64,uint64,bool))
func (_ContractPaymentVault *ContractPaymentVaultSession) GetQuorumProtocolConfig(quorumId uint64) (PaymentVaultTypesQuorumProtocolConfig, error) {
	return _ContractPaymentVault.Contract.GetQuorumProtocolConfig(&_ContractPaymentVault.CallOpts, quorumId)
}

// GetQuorumProtocolConfig is a free data retrieval call binding the contract method 0x89a06b35.
//
// Solidity: function getQuorumProtocolConfig(uint64 quorumId) view returns((uint64,uint64,uint64,uint64,bool))
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) GetQuorumProtocolConfig(quorumId uint64) (PaymentVaultTypesQuorumProtocolConfig, error) {
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
func (_ContractPaymentVault *ContractPaymentVaultCaller) GetReservation(opts *bind.CallOpts, quorumId uint64, account common.Address) (PaymentVaultTypesReservation, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "getReservation", quorumId, account)

	if err != nil {
		return *new(PaymentVaultTypesReservation), err
	}

	out0 := *abi.ConvertType(out[0], new(PaymentVaultTypesReservation)).(*PaymentVaultTypesReservation)

	return out0, err

}

// GetReservation is a free data retrieval call binding the contract method 0x00e691aa.
//
// Solidity: function getReservation(uint64 quorumId, address account) view returns((uint64,uint64,uint64))
func (_ContractPaymentVault *ContractPaymentVaultSession) GetReservation(quorumId uint64, account common.Address) (PaymentVaultTypesReservation, error) {
	return _ContractPaymentVault.Contract.GetReservation(&_ContractPaymentVault.CallOpts, quorumId, account)
}

// GetReservation is a free data retrieval call binding the contract method 0x00e691aa.
//
// Solidity: function getReservation(uint64 quorumId, address account) view returns((uint64,uint64,uint64))
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) GetReservation(quorumId uint64, account common.Address) (PaymentVaultTypesReservation, error) {
	return _ContractPaymentVault.Contract.GetReservation(&_ContractPaymentVault.CallOpts, quorumId, account)
}

// CreateReservation is a paid mutator transaction binding the contract method 0xecfc2d8e.
//
// Solidity: function createReservation(uint64 quorumId, address account, (uint64,uint64,uint64) reservation) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) CreateReservation(opts *bind.TransactOpts, quorumId uint64, account common.Address, reservation PaymentVaultTypesReservation) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "createReservation", quorumId, account, reservation)
}

// CreateReservation is a paid mutator transaction binding the contract method 0xecfc2d8e.
//
// Solidity: function createReservation(uint64 quorumId, address account, (uint64,uint64,uint64) reservation) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) CreateReservation(quorumId uint64, account common.Address, reservation PaymentVaultTypesReservation) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.CreateReservation(&_ContractPaymentVault.TransactOpts, quorumId, account, reservation)
}

// CreateReservation is a paid mutator transaction binding the contract method 0xecfc2d8e.
//
// Solidity: function createReservation(uint64 quorumId, address account, (uint64,uint64,uint64) reservation) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) CreateReservation(quorumId uint64, account common.Address, reservation PaymentVaultTypesReservation) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.CreateReservation(&_ContractPaymentVault.TransactOpts, quorumId, account, reservation)
}

// DecreaseReservation is a paid mutator transaction binding the contract method 0x3b50e8b8.
//
// Solidity: function decreaseReservation(uint64 quorumId, (uint64,uint64,uint64) reservation) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) DecreaseReservation(opts *bind.TransactOpts, quorumId uint64, reservation PaymentVaultTypesReservation) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "decreaseReservation", quorumId, reservation)
}

// DecreaseReservation is a paid mutator transaction binding the contract method 0x3b50e8b8.
//
// Solidity: function decreaseReservation(uint64 quorumId, (uint64,uint64,uint64) reservation) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) DecreaseReservation(quorumId uint64, reservation PaymentVaultTypesReservation) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.DecreaseReservation(&_ContractPaymentVault.TransactOpts, quorumId, reservation)
}

// DecreaseReservation is a paid mutator transaction binding the contract method 0x3b50e8b8.
//
// Solidity: function decreaseReservation(uint64 quorumId, (uint64,uint64,uint64) reservation) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) DecreaseReservation(quorumId uint64, reservation PaymentVaultTypesReservation) (*types.Transaction, error) {
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
func (_ContractPaymentVault *ContractPaymentVaultTransactor) IncreaseReservation(opts *bind.TransactOpts, quorumId uint64, account common.Address, reservation PaymentVaultTypesReservation) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "increaseReservation", quorumId, account, reservation)
}

// IncreaseReservation is a paid mutator transaction binding the contract method 0xe453a058.
//
// Solidity: function increaseReservation(uint64 quorumId, address account, (uint64,uint64,uint64) reservation) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) IncreaseReservation(quorumId uint64, account common.Address, reservation PaymentVaultTypesReservation) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.IncreaseReservation(&_ContractPaymentVault.TransactOpts, quorumId, account, reservation)
}

// IncreaseReservation is a paid mutator transaction binding the contract method 0xe453a058.
//
// Solidity: function increaseReservation(uint64 quorumId, address account, (uint64,uint64,uint64) reservation) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) IncreaseReservation(quorumId uint64, account common.Address, reservation PaymentVaultTypesReservation) (*types.Transaction, error) {
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

// InitializeQuorum is a paid mutator transaction binding the contract method 0x063a5d56.
//
// Solidity: function initializeQuorum(uint64 quorumId, address newOwner, (uint64,uint64,uint64,uint64,bool) protocolCfg) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) InitializeQuorum(opts *bind.TransactOpts, quorumId uint64, newOwner common.Address, protocolCfg PaymentVaultTypesQuorumProtocolConfig) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "initializeQuorum", quorumId, newOwner, protocolCfg)
}

// InitializeQuorum is a paid mutator transaction binding the contract method 0x063a5d56.
//
// Solidity: function initializeQuorum(uint64 quorumId, address newOwner, (uint64,uint64,uint64,uint64,bool) protocolCfg) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) InitializeQuorum(quorumId uint64, newOwner common.Address, protocolCfg PaymentVaultTypesQuorumProtocolConfig) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.InitializeQuorum(&_ContractPaymentVault.TransactOpts, quorumId, newOwner, protocolCfg)
}

// InitializeQuorum is a paid mutator transaction binding the contract method 0x063a5d56.
//
// Solidity: function initializeQuorum(uint64 quorumId, address newOwner, (uint64,uint64,uint64,uint64,bool) protocolCfg) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) InitializeQuorum(quorumId uint64, newOwner common.Address, protocolCfg PaymentVaultTypesQuorumProtocolConfig) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.InitializeQuorum(&_ContractPaymentVault.TransactOpts, quorumId, newOwner, protocolCfg)
}

// SetOnDemandEnabled is a paid mutator transaction binding the contract method 0x710f48be.
//
// Solidity: function setOnDemandEnabled(uint64 quorumId, (uint64,uint64,uint64,uint64,bool) protocolCfg) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) SetOnDemandEnabled(opts *bind.TransactOpts, quorumId uint64, protocolCfg PaymentVaultTypesQuorumProtocolConfig) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "setOnDemandEnabled", quorumId, protocolCfg)
}

// SetOnDemandEnabled is a paid mutator transaction binding the contract method 0x710f48be.
//
// Solidity: function setOnDemandEnabled(uint64 quorumId, (uint64,uint64,uint64,uint64,bool) protocolCfg) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) SetOnDemandEnabled(quorumId uint64, protocolCfg PaymentVaultTypesQuorumProtocolConfig) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetOnDemandEnabled(&_ContractPaymentVault.TransactOpts, quorumId, protocolCfg)
}

// SetOnDemandEnabled is a paid mutator transaction binding the contract method 0x710f48be.
//
// Solidity: function setOnDemandEnabled(uint64 quorumId, (uint64,uint64,uint64,uint64,bool) protocolCfg) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) SetOnDemandEnabled(quorumId uint64, protocolCfg PaymentVaultTypesQuorumProtocolConfig) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetOnDemandEnabled(&_ContractPaymentVault.TransactOpts, quorumId, protocolCfg)
}

// SetQuorumPaymentConfig is a paid mutator transaction binding the contract method 0xe631f82c.
//
// Solidity: function setQuorumPaymentConfig(uint64 quorumId, (address,address,uint64,uint64,uint64) paymentConfig) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) SetQuorumPaymentConfig(opts *bind.TransactOpts, quorumId uint64, paymentConfig PaymentVaultTypesQuorumConfig) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "setQuorumPaymentConfig", quorumId, paymentConfig)
}

// SetQuorumPaymentConfig is a paid mutator transaction binding the contract method 0xe631f82c.
//
// Solidity: function setQuorumPaymentConfig(uint64 quorumId, (address,address,uint64,uint64,uint64) paymentConfig) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) SetQuorumPaymentConfig(quorumId uint64, paymentConfig PaymentVaultTypesQuorumConfig) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetQuorumPaymentConfig(&_ContractPaymentVault.TransactOpts, quorumId, paymentConfig)
}

// SetQuorumPaymentConfig is a paid mutator transaction binding the contract method 0xe631f82c.
//
// Solidity: function setQuorumPaymentConfig(uint64 quorumId, (address,address,uint64,uint64,uint64) paymentConfig) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) SetQuorumPaymentConfig(quorumId uint64, paymentConfig PaymentVaultTypesQuorumConfig) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetQuorumPaymentConfig(&_ContractPaymentVault.TransactOpts, quorumId, paymentConfig)
}

// SetReservationAdvanceWindow is a paid mutator transaction binding the contract method 0x1269a9d7.
//
// Solidity: function setReservationAdvanceWindow(uint64 quorumId, (uint64,uint64,uint64,uint64,bool) protocolCfg) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) SetReservationAdvanceWindow(opts *bind.TransactOpts, quorumId uint64, protocolCfg PaymentVaultTypesQuorumProtocolConfig) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "setReservationAdvanceWindow", quorumId, protocolCfg)
}

// SetReservationAdvanceWindow is a paid mutator transaction binding the contract method 0x1269a9d7.
//
// Solidity: function setReservationAdvanceWindow(uint64 quorumId, (uint64,uint64,uint64,uint64,bool) protocolCfg) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) SetReservationAdvanceWindow(quorumId uint64, protocolCfg PaymentVaultTypesQuorumProtocolConfig) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetReservationAdvanceWindow(&_ContractPaymentVault.TransactOpts, quorumId, protocolCfg)
}

// SetReservationAdvanceWindow is a paid mutator transaction binding the contract method 0x1269a9d7.
//
// Solidity: function setReservationAdvanceWindow(uint64 quorumId, (uint64,uint64,uint64,uint64,bool) protocolCfg) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) SetReservationAdvanceWindow(quorumId uint64, protocolCfg PaymentVaultTypesQuorumProtocolConfig) (*types.Transaction, error) {
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
