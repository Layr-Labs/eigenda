// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractPaymentVault

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = bytes.Equal
	_ = errors.New
	_ = big.NewInt
	_ = common.Big1
	_ = types.BloomLookup
	_ = abi.ConvertType
)

// IPaymentVaultReservation is an auto generated low-level Go binding around an user-defined struct.
type IPaymentVaultReservation struct {
	SymbolsPerSecond uint64
	StartTimestamp   uint64
	EndTimestamp     uint64
	QuorumNumbers    []byte
	QuorumSplits     []byte
}

// ContractPaymentVaultMetaData contains all meta data concerning the ContractPaymentVault contract.
var ContractPaymentVaultMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"depositOnDemand\",\"inputs\":[{\"name\":\"_account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"getOnDemandTotalDeposit\",\"inputs\":[{\"name\":\"_account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOnDemandTotalDeposits\",\"inputs\":[{\"name\":\"_accounts\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[{\"name\":\"_payments\",\"type\":\"uint80[]\",\"internalType\":\"uint80[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReservation\",\"inputs\":[{\"name\":\"_account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIPaymentVault.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReservations\",\"inputs\":[{\"name\":\"_accounts\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[{\"name\":\"_reservations\",\"type\":\"tuple[]\",\"internalType\":\"structIPaymentVault.Reservation[]\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"globalRatePeriodInterval\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"globalSymbolsPerPeriod\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_minNumSymbols\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_pricePerSymbol\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_priceUpdateCooldown\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_globalSymbolsPerPeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_reservationPeriodInterval\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_globalRatePeriodInterval\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"lastPriceUpdateTime\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minNumSymbols\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"onDemandPayments\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"totalDeposit\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pricePerSymbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"priceUpdateCooldown\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"reservationPeriodInterval\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"reservations\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setGlobalRatePeriodInterval\",\"inputs\":[{\"name\":\"_globalRatePeriodInterval\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setGlobalSymbolsPerPeriod\",\"inputs\":[{\"name\":\"_globalSymbolsPerPeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPriceParams\",\"inputs\":[{\"name\":\"_minNumSymbols\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_pricePerSymbol\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_priceUpdateCooldown\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setReservation\",\"inputs\":[{\"name\":\"_account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_reservation\",\"type\":\"tuple\",\"internalType\":\"structIPaymentVault.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setReservationPeriodInterval\",\"inputs\":[{\"name\":\"_reservationPeriodInterval\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawERC20\",\"inputs\":[{\"name\":\"_token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"GlobalRatePeriodIntervalUpdated\",\"inputs\":[{\"name\":\"previousValue\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"newValue\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GlobalSymbolsPerPeriodUpdated\",\"inputs\":[{\"name\":\"previousValue\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"newValue\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OnDemandPaymentUpdated\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"onDemandPayment\",\"type\":\"uint80\",\"indexed\":false,\"internalType\":\"uint80\"},{\"name\":\"totalDeposit\",\"type\":\"uint80\",\"indexed\":false,\"internalType\":\"uint80\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PriceParamsUpdated\",\"inputs\":[{\"name\":\"previousMinNumSymbols\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"newMinNumSymbols\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"previousPricePerSymbol\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"newPricePerSymbol\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"previousPriceUpdateCooldown\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"newPriceUpdateCooldown\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReservationPeriodIntervalUpdated\",\"inputs\":[{\"name\":\"previousValue\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"newValue\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReservationUpdated\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"reservation\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIPaymentVault.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"anonymous\":false}]",
	ID:  "ContractPaymentVault",
}

// ContractPaymentVault is an auto generated Go binding around an Ethereum contract.
type ContractPaymentVault struct {
	abi abi.ABI
}

// NewContractPaymentVault creates a new instance of ContractPaymentVault.
func NewContractPaymentVault() *ContractPaymentVault {
	parsed, err := ContractPaymentVaultMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ContractPaymentVault{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ContractPaymentVault) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackDepositOnDemand is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8bec7d02.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function depositOnDemand(address _account) payable returns()
func (contractPaymentVault *ContractPaymentVault) PackDepositOnDemand(account common.Address) []byte {
	enc, err := contractPaymentVault.abi.Pack("depositOnDemand", account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDepositOnDemand is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8bec7d02.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function depositOnDemand(address _account) payable returns()
func (contractPaymentVault *ContractPaymentVault) TryPackDepositOnDemand(account common.Address) ([]byte, error) {
	return contractPaymentVault.abi.Pack("depositOnDemand", account)
}

// PackGetOnDemandTotalDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd1c1fdcd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getOnDemandTotalDeposit(address _account) view returns(uint80)
func (contractPaymentVault *ContractPaymentVault) PackGetOnDemandTotalDeposit(account common.Address) []byte {
	enc, err := contractPaymentVault.abi.Pack("getOnDemandTotalDeposit", account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetOnDemandTotalDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd1c1fdcd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getOnDemandTotalDeposit(address _account) view returns(uint80)
func (contractPaymentVault *ContractPaymentVault) TryPackGetOnDemandTotalDeposit(account common.Address) ([]byte, error) {
	return contractPaymentVault.abi.Pack("getOnDemandTotalDeposit", account)
}

// UnpackGetOnDemandTotalDeposit is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd1c1fdcd.
//
// Solidity: function getOnDemandTotalDeposit(address _account) view returns(uint80)
func (contractPaymentVault *ContractPaymentVault) UnpackGetOnDemandTotalDeposit(data []byte) (*big.Int, error) {
	out, err := contractPaymentVault.abi.Unpack("getOnDemandTotalDeposit", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetOnDemandTotalDeposits is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4184a674.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getOnDemandTotalDeposits(address[] _accounts) view returns(uint80[] _payments)
func (contractPaymentVault *ContractPaymentVault) PackGetOnDemandTotalDeposits(accounts []common.Address) []byte {
	enc, err := contractPaymentVault.abi.Pack("getOnDemandTotalDeposits", accounts)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetOnDemandTotalDeposits is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4184a674.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getOnDemandTotalDeposits(address[] _accounts) view returns(uint80[] _payments)
func (contractPaymentVault *ContractPaymentVault) TryPackGetOnDemandTotalDeposits(accounts []common.Address) ([]byte, error) {
	return contractPaymentVault.abi.Pack("getOnDemandTotalDeposits", accounts)
}

// UnpackGetOnDemandTotalDeposits is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4184a674.
//
// Solidity: function getOnDemandTotalDeposits(address[] _accounts) view returns(uint80[] _payments)
func (contractPaymentVault *ContractPaymentVault) UnpackGetOnDemandTotalDeposits(data []byte) ([]*big.Int, error) {
	out, err := contractPaymentVault.abi.Unpack("getOnDemandTotalDeposits", data)
	if err != nil {
		return *new([]*big.Int), err
	}
	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	return out0, nil
}

// PackGetReservation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb2066f80.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getReservation(address _account) view returns((uint64,uint64,uint64,bytes,bytes))
func (contractPaymentVault *ContractPaymentVault) PackGetReservation(account common.Address) []byte {
	enc, err := contractPaymentVault.abi.Pack("getReservation", account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetReservation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb2066f80.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getReservation(address _account) view returns((uint64,uint64,uint64,bytes,bytes))
func (contractPaymentVault *ContractPaymentVault) TryPackGetReservation(account common.Address) ([]byte, error) {
	return contractPaymentVault.abi.Pack("getReservation", account)
}

// UnpackGetReservation is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb2066f80.
//
// Solidity: function getReservation(address _account) view returns((uint64,uint64,uint64,bytes,bytes))
func (contractPaymentVault *ContractPaymentVault) UnpackGetReservation(data []byte) (IPaymentVaultReservation, error) {
	out, err := contractPaymentVault.abi.Unpack("getReservation", data)
	if err != nil {
		return *new(IPaymentVaultReservation), err
	}
	out0 := *abi.ConvertType(out[0], new(IPaymentVaultReservation)).(*IPaymentVaultReservation)
	return out0, nil
}

// PackGetReservations is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x109f8fe5.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getReservations(address[] _accounts) view returns((uint64,uint64,uint64,bytes,bytes)[] _reservations)
func (contractPaymentVault *ContractPaymentVault) PackGetReservations(accounts []common.Address) []byte {
	enc, err := contractPaymentVault.abi.Pack("getReservations", accounts)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetReservations is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x109f8fe5.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getReservations(address[] _accounts) view returns((uint64,uint64,uint64,bytes,bytes)[] _reservations)
func (contractPaymentVault *ContractPaymentVault) TryPackGetReservations(accounts []common.Address) ([]byte, error) {
	return contractPaymentVault.abi.Pack("getReservations", accounts)
}

// UnpackGetReservations is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x109f8fe5.
//
// Solidity: function getReservations(address[] _accounts) view returns((uint64,uint64,uint64,bytes,bytes)[] _reservations)
func (contractPaymentVault *ContractPaymentVault) UnpackGetReservations(data []byte) ([]IPaymentVaultReservation, error) {
	out, err := contractPaymentVault.abi.Unpack("getReservations", data)
	if err != nil {
		return *new([]IPaymentVaultReservation), err
	}
	out0 := *abi.ConvertType(out[0], new([]IPaymentVaultReservation)).(*[]IPaymentVaultReservation)
	return out0, nil
}

// PackGlobalRatePeriodInterval is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbff8a3d4.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function globalRatePeriodInterval() view returns(uint64)
func (contractPaymentVault *ContractPaymentVault) PackGlobalRatePeriodInterval() []byte {
	enc, err := contractPaymentVault.abi.Pack("globalRatePeriodInterval")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGlobalRatePeriodInterval is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbff8a3d4.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function globalRatePeriodInterval() view returns(uint64)
func (contractPaymentVault *ContractPaymentVault) TryPackGlobalRatePeriodInterval() ([]byte, error) {
	return contractPaymentVault.abi.Pack("globalRatePeriodInterval")
}

// UnpackGlobalRatePeriodInterval is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbff8a3d4.
//
// Solidity: function globalRatePeriodInterval() view returns(uint64)
func (contractPaymentVault *ContractPaymentVault) UnpackGlobalRatePeriodInterval(data []byte) (uint64, error) {
	out, err := contractPaymentVault.abi.Unpack("globalRatePeriodInterval", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackGlobalSymbolsPerPeriod is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc98d97dd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function globalSymbolsPerPeriod() view returns(uint64)
func (contractPaymentVault *ContractPaymentVault) PackGlobalSymbolsPerPeriod() []byte {
	enc, err := contractPaymentVault.abi.Pack("globalSymbolsPerPeriod")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGlobalSymbolsPerPeriod is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc98d97dd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function globalSymbolsPerPeriod() view returns(uint64)
func (contractPaymentVault *ContractPaymentVault) TryPackGlobalSymbolsPerPeriod() ([]byte, error) {
	return contractPaymentVault.abi.Pack("globalSymbolsPerPeriod")
}

// UnpackGlobalSymbolsPerPeriod is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc98d97dd.
//
// Solidity: function globalSymbolsPerPeriod() view returns(uint64)
func (contractPaymentVault *ContractPaymentVault) UnpackGlobalSymbolsPerPeriod(data []byte) (uint64, error) {
	out, err := contractPaymentVault.abi.Unpack("globalSymbolsPerPeriod", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9a1bbf37.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function initialize(address _initialOwner, uint64 _minNumSymbols, uint64 _pricePerSymbol, uint64 _priceUpdateCooldown, uint64 _globalSymbolsPerPeriod, uint64 _reservationPeriodInterval, uint64 _globalRatePeriodInterval) returns()
func (contractPaymentVault *ContractPaymentVault) PackInitialize(initialOwner common.Address, minNumSymbols uint64, pricePerSymbol uint64, priceUpdateCooldown uint64, globalSymbolsPerPeriod uint64, reservationPeriodInterval uint64, globalRatePeriodInterval uint64) []byte {
	enc, err := contractPaymentVault.abi.Pack("initialize", initialOwner, minNumSymbols, pricePerSymbol, priceUpdateCooldown, globalSymbolsPerPeriod, reservationPeriodInterval, globalRatePeriodInterval)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9a1bbf37.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function initialize(address _initialOwner, uint64 _minNumSymbols, uint64 _pricePerSymbol, uint64 _priceUpdateCooldown, uint64 _globalSymbolsPerPeriod, uint64 _reservationPeriodInterval, uint64 _globalRatePeriodInterval) returns()
func (contractPaymentVault *ContractPaymentVault) TryPackInitialize(initialOwner common.Address, minNumSymbols uint64, pricePerSymbol uint64, priceUpdateCooldown uint64, globalSymbolsPerPeriod uint64, reservationPeriodInterval uint64, globalRatePeriodInterval uint64) ([]byte, error) {
	return contractPaymentVault.abi.Pack("initialize", initialOwner, minNumSymbols, pricePerSymbol, priceUpdateCooldown, globalSymbolsPerPeriod, reservationPeriodInterval, globalRatePeriodInterval)
}

// PackLastPriceUpdateTime is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x49b9a7af.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function lastPriceUpdateTime() view returns(uint64)
func (contractPaymentVault *ContractPaymentVault) PackLastPriceUpdateTime() []byte {
	enc, err := contractPaymentVault.abi.Pack("lastPriceUpdateTime")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackLastPriceUpdateTime is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x49b9a7af.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function lastPriceUpdateTime() view returns(uint64)
func (contractPaymentVault *ContractPaymentVault) TryPackLastPriceUpdateTime() ([]byte, error) {
	return contractPaymentVault.abi.Pack("lastPriceUpdateTime")
}

// UnpackLastPriceUpdateTime is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x49b9a7af.
//
// Solidity: function lastPriceUpdateTime() view returns(uint64)
func (contractPaymentVault *ContractPaymentVault) UnpackLastPriceUpdateTime(data []byte) (uint64, error) {
	out, err := contractPaymentVault.abi.Unpack("lastPriceUpdateTime", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackMinNumSymbols is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x761dab89.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function minNumSymbols() view returns(uint64)
func (contractPaymentVault *ContractPaymentVault) PackMinNumSymbols() []byte {
	enc, err := contractPaymentVault.abi.Pack("minNumSymbols")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMinNumSymbols is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x761dab89.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function minNumSymbols() view returns(uint64)
func (contractPaymentVault *ContractPaymentVault) TryPackMinNumSymbols() ([]byte, error) {
	return contractPaymentVault.abi.Pack("minNumSymbols")
}

// UnpackMinNumSymbols is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x761dab89.
//
// Solidity: function minNumSymbols() view returns(uint64)
func (contractPaymentVault *ContractPaymentVault) UnpackMinNumSymbols(data []byte) (uint64, error) {
	out, err := contractPaymentVault.abi.Unpack("minNumSymbols", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackOnDemandPayments is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd996dc99.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function onDemandPayments(address ) view returns(uint80 totalDeposit)
func (contractPaymentVault *ContractPaymentVault) PackOnDemandPayments(arg0 common.Address) []byte {
	enc, err := contractPaymentVault.abi.Pack("onDemandPayments", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOnDemandPayments is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd996dc99.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function onDemandPayments(address ) view returns(uint80 totalDeposit)
func (contractPaymentVault *ContractPaymentVault) TryPackOnDemandPayments(arg0 common.Address) ([]byte, error) {
	return contractPaymentVault.abi.Pack("onDemandPayments", arg0)
}

// UnpackOnDemandPayments is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd996dc99.
//
// Solidity: function onDemandPayments(address ) view returns(uint80 totalDeposit)
func (contractPaymentVault *ContractPaymentVault) UnpackOnDemandPayments(data []byte) (*big.Int, error) {
	out, err := contractPaymentVault.abi.Unpack("onDemandPayments", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function owner() view returns(address)
func (contractPaymentVault *ContractPaymentVault) PackOwner() []byte {
	enc, err := contractPaymentVault.abi.Pack("owner")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function owner() view returns(address)
func (contractPaymentVault *ContractPaymentVault) TryPackOwner() ([]byte, error) {
	return contractPaymentVault.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (contractPaymentVault *ContractPaymentVault) UnpackOwner(data []byte) (common.Address, error) {
	out, err := contractPaymentVault.abi.Unpack("owner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackPricePerSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf323726a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pricePerSymbol() view returns(uint64)
func (contractPaymentVault *ContractPaymentVault) PackPricePerSymbol() []byte {
	enc, err := contractPaymentVault.abi.Pack("pricePerSymbol")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPricePerSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf323726a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pricePerSymbol() view returns(uint64)
func (contractPaymentVault *ContractPaymentVault) TryPackPricePerSymbol() ([]byte, error) {
	return contractPaymentVault.abi.Pack("pricePerSymbol")
}

// UnpackPricePerSymbol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf323726a.
//
// Solidity: function pricePerSymbol() view returns(uint64)
func (contractPaymentVault *ContractPaymentVault) UnpackPricePerSymbol(data []byte) (uint64, error) {
	out, err := contractPaymentVault.abi.Unpack("pricePerSymbol", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackPriceUpdateCooldown is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x039f091c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function priceUpdateCooldown() view returns(uint64)
func (contractPaymentVault *ContractPaymentVault) PackPriceUpdateCooldown() []byte {
	enc, err := contractPaymentVault.abi.Pack("priceUpdateCooldown")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPriceUpdateCooldown is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x039f091c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function priceUpdateCooldown() view returns(uint64)
func (contractPaymentVault *ContractPaymentVault) TryPackPriceUpdateCooldown() ([]byte, error) {
	return contractPaymentVault.abi.Pack("priceUpdateCooldown")
}

// UnpackPriceUpdateCooldown is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x039f091c.
//
// Solidity: function priceUpdateCooldown() view returns(uint64)
func (contractPaymentVault *ContractPaymentVault) UnpackPriceUpdateCooldown(data []byte) (uint64, error) {
	out, err := contractPaymentVault.abi.Unpack("priceUpdateCooldown", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackRenounceOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715018a6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceOwnership() returns()
func (contractPaymentVault *ContractPaymentVault) PackRenounceOwnership() []byte {
	enc, err := contractPaymentVault.abi.Pack("renounceOwnership")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRenounceOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715018a6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function renounceOwnership() returns()
func (contractPaymentVault *ContractPaymentVault) TryPackRenounceOwnership() ([]byte, error) {
	return contractPaymentVault.abi.Pack("renounceOwnership")
}

// PackReservationPeriodInterval is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x72228ab2.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function reservationPeriodInterval() view returns(uint64)
func (contractPaymentVault *ContractPaymentVault) PackReservationPeriodInterval() []byte {
	enc, err := contractPaymentVault.abi.Pack("reservationPeriodInterval")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackReservationPeriodInterval is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x72228ab2.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function reservationPeriodInterval() view returns(uint64)
func (contractPaymentVault *ContractPaymentVault) TryPackReservationPeriodInterval() ([]byte, error) {
	return contractPaymentVault.abi.Pack("reservationPeriodInterval")
}

// UnpackReservationPeriodInterval is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x72228ab2.
//
// Solidity: function reservationPeriodInterval() view returns(uint64)
func (contractPaymentVault *ContractPaymentVault) UnpackReservationPeriodInterval(data []byte) (uint64, error) {
	out, err := contractPaymentVault.abi.Unpack("reservationPeriodInterval", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackReservations is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfd3dc53a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function reservations(address ) view returns(uint64 symbolsPerSecond, uint64 startTimestamp, uint64 endTimestamp, bytes quorumNumbers, bytes quorumSplits)
func (contractPaymentVault *ContractPaymentVault) PackReservations(arg0 common.Address) []byte {
	enc, err := contractPaymentVault.abi.Pack("reservations", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackReservations is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfd3dc53a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function reservations(address ) view returns(uint64 symbolsPerSecond, uint64 startTimestamp, uint64 endTimestamp, bytes quorumNumbers, bytes quorumSplits)
func (contractPaymentVault *ContractPaymentVault) TryPackReservations(arg0 common.Address) ([]byte, error) {
	return contractPaymentVault.abi.Pack("reservations", arg0)
}

// ReservationsOutput serves as a container for the return parameters of contract
// method Reservations.
type ReservationsOutput struct {
	SymbolsPerSecond uint64
	StartTimestamp   uint64
	EndTimestamp     uint64
	QuorumNumbers    []byte
	QuorumSplits     []byte
}

// UnpackReservations is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfd3dc53a.
//
// Solidity: function reservations(address ) view returns(uint64 symbolsPerSecond, uint64 startTimestamp, uint64 endTimestamp, bytes quorumNumbers, bytes quorumSplits)
func (contractPaymentVault *ContractPaymentVault) UnpackReservations(data []byte) (ReservationsOutput, error) {
	out, err := contractPaymentVault.abi.Unpack("reservations", data)
	outstruct := new(ReservationsOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.SymbolsPerSecond = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.StartTimestamp = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.EndTimestamp = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.QuorumNumbers = *abi.ConvertType(out[3], new([]byte)).(*[]byte)
	outstruct.QuorumSplits = *abi.ConvertType(out[4], new([]byte)).(*[]byte)
	return *outstruct, nil
}

// PackSetGlobalRatePeriodInterval is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaa788bd7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setGlobalRatePeriodInterval(uint64 _globalRatePeriodInterval) returns()
func (contractPaymentVault *ContractPaymentVault) PackSetGlobalRatePeriodInterval(globalRatePeriodInterval uint64) []byte {
	enc, err := contractPaymentVault.abi.Pack("setGlobalRatePeriodInterval", globalRatePeriodInterval)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetGlobalRatePeriodInterval is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaa788bd7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setGlobalRatePeriodInterval(uint64 _globalRatePeriodInterval) returns()
func (contractPaymentVault *ContractPaymentVault) TryPackSetGlobalRatePeriodInterval(globalRatePeriodInterval uint64) ([]byte, error) {
	return contractPaymentVault.abi.Pack("setGlobalRatePeriodInterval", globalRatePeriodInterval)
}

// PackSetGlobalSymbolsPerPeriod is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa16cf884.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setGlobalSymbolsPerPeriod(uint64 _globalSymbolsPerPeriod) returns()
func (contractPaymentVault *ContractPaymentVault) PackSetGlobalSymbolsPerPeriod(globalSymbolsPerPeriod uint64) []byte {
	enc, err := contractPaymentVault.abi.Pack("setGlobalSymbolsPerPeriod", globalSymbolsPerPeriod)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetGlobalSymbolsPerPeriod is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa16cf884.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setGlobalSymbolsPerPeriod(uint64 _globalSymbolsPerPeriod) returns()
func (contractPaymentVault *ContractPaymentVault) TryPackSetGlobalSymbolsPerPeriod(globalSymbolsPerPeriod uint64) ([]byte, error) {
	return contractPaymentVault.abi.Pack("setGlobalSymbolsPerPeriod", globalSymbolsPerPeriod)
}

// PackSetPriceParams is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfba2b1d1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setPriceParams(uint64 _minNumSymbols, uint64 _pricePerSymbol, uint64 _priceUpdateCooldown) returns()
func (contractPaymentVault *ContractPaymentVault) PackSetPriceParams(minNumSymbols uint64, pricePerSymbol uint64, priceUpdateCooldown uint64) []byte {
	enc, err := contractPaymentVault.abi.Pack("setPriceParams", minNumSymbols, pricePerSymbol, priceUpdateCooldown)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetPriceParams is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfba2b1d1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setPriceParams(uint64 _minNumSymbols, uint64 _pricePerSymbol, uint64 _priceUpdateCooldown) returns()
func (contractPaymentVault *ContractPaymentVault) TryPackSetPriceParams(minNumSymbols uint64, pricePerSymbol uint64, priceUpdateCooldown uint64) ([]byte, error) {
	return contractPaymentVault.abi.Pack("setPriceParams", minNumSymbols, pricePerSymbol, priceUpdateCooldown)
}

// PackSetReservation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9aec8640.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setReservation(address _account, (uint64,uint64,uint64,bytes,bytes) _reservation) returns()
func (contractPaymentVault *ContractPaymentVault) PackSetReservation(account common.Address, reservation IPaymentVaultReservation) []byte {
	enc, err := contractPaymentVault.abi.Pack("setReservation", account, reservation)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetReservation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9aec8640.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setReservation(address _account, (uint64,uint64,uint64,bytes,bytes) _reservation) returns()
func (contractPaymentVault *ContractPaymentVault) TryPackSetReservation(account common.Address, reservation IPaymentVaultReservation) ([]byte, error) {
	return contractPaymentVault.abi.Pack("setReservation", account, reservation)
}

// PackSetReservationPeriodInterval is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x897218fc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setReservationPeriodInterval(uint64 _reservationPeriodInterval) returns()
func (contractPaymentVault *ContractPaymentVault) PackSetReservationPeriodInterval(reservationPeriodInterval uint64) []byte {
	enc, err := contractPaymentVault.abi.Pack("setReservationPeriodInterval", reservationPeriodInterval)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetReservationPeriodInterval is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x897218fc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setReservationPeriodInterval(uint64 _reservationPeriodInterval) returns()
func (contractPaymentVault *ContractPaymentVault) TryPackSetReservationPeriodInterval(reservationPeriodInterval uint64) ([]byte, error) {
	return contractPaymentVault.abi.Pack("setReservationPeriodInterval", reservationPeriodInterval)
}

// PackTransferOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2fde38b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (contractPaymentVault *ContractPaymentVault) PackTransferOwnership(newOwner common.Address) []byte {
	enc, err := contractPaymentVault.abi.Pack("transferOwnership", newOwner)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTransferOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2fde38b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (contractPaymentVault *ContractPaymentVault) TryPackTransferOwnership(newOwner common.Address) ([]byte, error) {
	return contractPaymentVault.abi.Pack("transferOwnership", newOwner)
}

// PackWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2e1a7d4d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdraw(uint256 _amount) returns()
func (contractPaymentVault *ContractPaymentVault) PackWithdraw(amount *big.Int) []byte {
	enc, err := contractPaymentVault.abi.Pack("withdraw", amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2e1a7d4d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withdraw(uint256 _amount) returns()
func (contractPaymentVault *ContractPaymentVault) TryPackWithdraw(amount *big.Int) ([]byte, error) {
	return contractPaymentVault.abi.Pack("withdraw", amount)
}

// PackWithdrawERC20 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa1db9782.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdrawERC20(address _token, uint256 _amount) returns()
func (contractPaymentVault *ContractPaymentVault) PackWithdrawERC20(token common.Address, amount *big.Int) []byte {
	enc, err := contractPaymentVault.abi.Pack("withdrawERC20", token, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithdrawERC20 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa1db9782.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withdrawERC20(address _token, uint256 _amount) returns()
func (contractPaymentVault *ContractPaymentVault) TryPackWithdrawERC20(token common.Address, amount *big.Int) ([]byte, error) {
	return contractPaymentVault.abi.Pack("withdrawERC20", token, amount)
}

// ContractPaymentVaultGlobalRatePeriodIntervalUpdated represents a GlobalRatePeriodIntervalUpdated event raised by the ContractPaymentVault contract.
type ContractPaymentVaultGlobalRatePeriodIntervalUpdated struct {
	PreviousValue uint64
	NewValue      uint64
	Raw           *types.Log // Blockchain specific contextual infos
}

const ContractPaymentVaultGlobalRatePeriodIntervalUpdatedEventName = "GlobalRatePeriodIntervalUpdated"

// ContractEventName returns the user-defined event name.
func (ContractPaymentVaultGlobalRatePeriodIntervalUpdated) ContractEventName() string {
	return ContractPaymentVaultGlobalRatePeriodIntervalUpdatedEventName
}

// UnpackGlobalRatePeriodIntervalUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event GlobalRatePeriodIntervalUpdated(uint64 previousValue, uint64 newValue)
func (contractPaymentVault *ContractPaymentVault) UnpackGlobalRatePeriodIntervalUpdatedEvent(log *types.Log) (*ContractPaymentVaultGlobalRatePeriodIntervalUpdated, error) {
	event := "GlobalRatePeriodIntervalUpdated"
	if log.Topics[0] != contractPaymentVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ContractPaymentVaultGlobalRatePeriodIntervalUpdated)
	if len(log.Data) > 0 {
		if err := contractPaymentVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range contractPaymentVault.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// ContractPaymentVaultGlobalSymbolsPerPeriodUpdated represents a GlobalSymbolsPerPeriodUpdated event raised by the ContractPaymentVault contract.
type ContractPaymentVaultGlobalSymbolsPerPeriodUpdated struct {
	PreviousValue uint64
	NewValue      uint64
	Raw           *types.Log // Blockchain specific contextual infos
}

const ContractPaymentVaultGlobalSymbolsPerPeriodUpdatedEventName = "GlobalSymbolsPerPeriodUpdated"

// ContractEventName returns the user-defined event name.
func (ContractPaymentVaultGlobalSymbolsPerPeriodUpdated) ContractEventName() string {
	return ContractPaymentVaultGlobalSymbolsPerPeriodUpdatedEventName
}

// UnpackGlobalSymbolsPerPeriodUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event GlobalSymbolsPerPeriodUpdated(uint64 previousValue, uint64 newValue)
func (contractPaymentVault *ContractPaymentVault) UnpackGlobalSymbolsPerPeriodUpdatedEvent(log *types.Log) (*ContractPaymentVaultGlobalSymbolsPerPeriodUpdated, error) {
	event := "GlobalSymbolsPerPeriodUpdated"
	if log.Topics[0] != contractPaymentVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ContractPaymentVaultGlobalSymbolsPerPeriodUpdated)
	if len(log.Data) > 0 {
		if err := contractPaymentVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range contractPaymentVault.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// ContractPaymentVaultInitialized represents a Initialized event raised by the ContractPaymentVault contract.
type ContractPaymentVaultInitialized struct {
	Version uint8
	Raw     *types.Log // Blockchain specific contextual infos
}

const ContractPaymentVaultInitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (ContractPaymentVaultInitialized) ContractEventName() string {
	return ContractPaymentVaultInitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint8 version)
func (contractPaymentVault *ContractPaymentVault) UnpackInitializedEvent(log *types.Log) (*ContractPaymentVaultInitialized, error) {
	event := "Initialized"
	if log.Topics[0] != contractPaymentVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ContractPaymentVaultInitialized)
	if len(log.Data) > 0 {
		if err := contractPaymentVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range contractPaymentVault.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// ContractPaymentVaultOnDemandPaymentUpdated represents a OnDemandPaymentUpdated event raised by the ContractPaymentVault contract.
type ContractPaymentVaultOnDemandPaymentUpdated struct {
	Account         common.Address
	OnDemandPayment *big.Int
	TotalDeposit    *big.Int
	Raw             *types.Log // Blockchain specific contextual infos
}

const ContractPaymentVaultOnDemandPaymentUpdatedEventName = "OnDemandPaymentUpdated"

// ContractEventName returns the user-defined event name.
func (ContractPaymentVaultOnDemandPaymentUpdated) ContractEventName() string {
	return ContractPaymentVaultOnDemandPaymentUpdatedEventName
}

// UnpackOnDemandPaymentUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OnDemandPaymentUpdated(address indexed account, uint80 onDemandPayment, uint80 totalDeposit)
func (contractPaymentVault *ContractPaymentVault) UnpackOnDemandPaymentUpdatedEvent(log *types.Log) (*ContractPaymentVaultOnDemandPaymentUpdated, error) {
	event := "OnDemandPaymentUpdated"
	if log.Topics[0] != contractPaymentVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ContractPaymentVaultOnDemandPaymentUpdated)
	if len(log.Data) > 0 {
		if err := contractPaymentVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range contractPaymentVault.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// ContractPaymentVaultOwnershipTransferred represents a OwnershipTransferred event raised by the ContractPaymentVault contract.
type ContractPaymentVaultOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const ContractPaymentVaultOwnershipTransferredEventName = "OwnershipTransferred"

// ContractEventName returns the user-defined event name.
func (ContractPaymentVaultOwnershipTransferred) ContractEventName() string {
	return ContractPaymentVaultOwnershipTransferredEventName
}

// UnpackOwnershipTransferredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (contractPaymentVault *ContractPaymentVault) UnpackOwnershipTransferredEvent(log *types.Log) (*ContractPaymentVaultOwnershipTransferred, error) {
	event := "OwnershipTransferred"
	if log.Topics[0] != contractPaymentVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ContractPaymentVaultOwnershipTransferred)
	if len(log.Data) > 0 {
		if err := contractPaymentVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range contractPaymentVault.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// ContractPaymentVaultPriceParamsUpdated represents a PriceParamsUpdated event raised by the ContractPaymentVault contract.
type ContractPaymentVaultPriceParamsUpdated struct {
	PreviousMinNumSymbols       uint64
	NewMinNumSymbols            uint64
	PreviousPricePerSymbol      uint64
	NewPricePerSymbol           uint64
	PreviousPriceUpdateCooldown uint64
	NewPriceUpdateCooldown      uint64
	Raw                         *types.Log // Blockchain specific contextual infos
}

const ContractPaymentVaultPriceParamsUpdatedEventName = "PriceParamsUpdated"

// ContractEventName returns the user-defined event name.
func (ContractPaymentVaultPriceParamsUpdated) ContractEventName() string {
	return ContractPaymentVaultPriceParamsUpdatedEventName
}

// UnpackPriceParamsUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PriceParamsUpdated(uint64 previousMinNumSymbols, uint64 newMinNumSymbols, uint64 previousPricePerSymbol, uint64 newPricePerSymbol, uint64 previousPriceUpdateCooldown, uint64 newPriceUpdateCooldown)
func (contractPaymentVault *ContractPaymentVault) UnpackPriceParamsUpdatedEvent(log *types.Log) (*ContractPaymentVaultPriceParamsUpdated, error) {
	event := "PriceParamsUpdated"
	if log.Topics[0] != contractPaymentVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ContractPaymentVaultPriceParamsUpdated)
	if len(log.Data) > 0 {
		if err := contractPaymentVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range contractPaymentVault.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// ContractPaymentVaultReservationPeriodIntervalUpdated represents a ReservationPeriodIntervalUpdated event raised by the ContractPaymentVault contract.
type ContractPaymentVaultReservationPeriodIntervalUpdated struct {
	PreviousValue uint64
	NewValue      uint64
	Raw           *types.Log // Blockchain specific contextual infos
}

const ContractPaymentVaultReservationPeriodIntervalUpdatedEventName = "ReservationPeriodIntervalUpdated"

// ContractEventName returns the user-defined event name.
func (ContractPaymentVaultReservationPeriodIntervalUpdated) ContractEventName() string {
	return ContractPaymentVaultReservationPeriodIntervalUpdatedEventName
}

// UnpackReservationPeriodIntervalUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ReservationPeriodIntervalUpdated(uint64 previousValue, uint64 newValue)
func (contractPaymentVault *ContractPaymentVault) UnpackReservationPeriodIntervalUpdatedEvent(log *types.Log) (*ContractPaymentVaultReservationPeriodIntervalUpdated, error) {
	event := "ReservationPeriodIntervalUpdated"
	if log.Topics[0] != contractPaymentVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ContractPaymentVaultReservationPeriodIntervalUpdated)
	if len(log.Data) > 0 {
		if err := contractPaymentVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range contractPaymentVault.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// ContractPaymentVaultReservationUpdated represents a ReservationUpdated event raised by the ContractPaymentVault contract.
type ContractPaymentVaultReservationUpdated struct {
	Account     common.Address
	Reservation IPaymentVaultReservation
	Raw         *types.Log // Blockchain specific contextual infos
}

const ContractPaymentVaultReservationUpdatedEventName = "ReservationUpdated"

// ContractEventName returns the user-defined event name.
func (ContractPaymentVaultReservationUpdated) ContractEventName() string {
	return ContractPaymentVaultReservationUpdatedEventName
}

// UnpackReservationUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ReservationUpdated(address indexed account, (uint64,uint64,uint64,bytes,bytes) reservation)
func (contractPaymentVault *ContractPaymentVault) UnpackReservationUpdatedEvent(log *types.Log) (*ContractPaymentVaultReservationUpdated, error) {
	event := "ReservationUpdated"
	if log.Topics[0] != contractPaymentVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ContractPaymentVaultReservationUpdated)
	if len(log.Data) > 0 {
		if err := contractPaymentVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range contractPaymentVault.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}
