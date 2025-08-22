// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractEigenDACertVerifierRouter

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

// ContractEigenDACertVerifierRouterMetaData contains all meta data concerning the ContractEigenDACertVerifierRouter contract.
var ContractEigenDACertVerifierRouterMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"addCertVerifier\",\"inputs\":[{\"name\":\"activationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"certVerifier\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"certVerifierABNs\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"certVerifiers\",\"inputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"checkDACert\",\"inputs\":[{\"name\":\"abiEncodedCert\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCertVerifierAt\",\"inputs\":[{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"initABNs\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"initCertVerifiers\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"CertVerifierAdded\",\"inputs\":[{\"name\":\"activationBlockNumber\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"},{\"name\":\"certVerifier\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ABNNotGreaterThanLast\",\"inputs\":[{\"name\":\"activationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"ABNNotInFuture\",\"inputs\":[{\"name\":\"activationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"InvalidCertLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RBNInFuture\",\"inputs\":[{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]",
	ID:  "ContractEigenDACertVerifierRouter",
}

// ContractEigenDACertVerifierRouter is an auto generated Go binding around an Ethereum contract.
type ContractEigenDACertVerifierRouter struct {
	abi abi.ABI
}

// NewContractEigenDACertVerifierRouter creates a new instance of ContractEigenDACertVerifierRouter.
func NewContractEigenDACertVerifierRouter() *ContractEigenDACertVerifierRouter {
	parsed, err := ContractEigenDACertVerifierRouterMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ContractEigenDACertVerifierRouter{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ContractEigenDACertVerifierRouter) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackAddCertVerifier is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbfda00de.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function addCertVerifier(uint32 activationBlockNumber, address certVerifier) returns()
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) PackAddCertVerifier(activationBlockNumber uint32, certVerifier common.Address) []byte {
	enc, err := contractEigenDACertVerifierRouter.abi.Pack("addCertVerifier", activationBlockNumber, certVerifier)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAddCertVerifier is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbfda00de.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function addCertVerifier(uint32 activationBlockNumber, address certVerifier) returns()
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) TryPackAddCertVerifier(activationBlockNumber uint32, certVerifier common.Address) ([]byte, error) {
	return contractEigenDACertVerifierRouter.abi.Pack("addCertVerifier", activationBlockNumber, certVerifier)
}

// PackCertVerifierABNs is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf0df66df.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function certVerifierABNs(uint256 ) view returns(uint32)
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) PackCertVerifierABNs(arg0 *big.Int) []byte {
	enc, err := contractEigenDACertVerifierRouter.abi.Pack("certVerifierABNs", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCertVerifierABNs is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf0df66df.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function certVerifierABNs(uint256 ) view returns(uint32)
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) TryPackCertVerifierABNs(arg0 *big.Int) ([]byte, error) {
	return contractEigenDACertVerifierRouter.abi.Pack("certVerifierABNs", arg0)
}

// UnpackCertVerifierABNs is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf0df66df.
//
// Solidity: function certVerifierABNs(uint256 ) view returns(uint32)
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) UnpackCertVerifierABNs(data []byte) (uint32, error) {
	out, err := contractEigenDACertVerifierRouter.abi.Unpack("certVerifierABNs", data)
	if err != nil {
		return *new(uint32), err
	}
	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)
	return out0, nil
}

// PackCertVerifiers is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4c046566.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function certVerifiers(uint32 ) view returns(address)
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) PackCertVerifiers(arg0 uint32) []byte {
	enc, err := contractEigenDACertVerifierRouter.abi.Pack("certVerifiers", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCertVerifiers is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4c046566.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function certVerifiers(uint32 ) view returns(address)
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) TryPackCertVerifiers(arg0 uint32) ([]byte, error) {
	return contractEigenDACertVerifierRouter.abi.Pack("certVerifiers", arg0)
}

// UnpackCertVerifiers is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4c046566.
//
// Solidity: function certVerifiers(uint32 ) view returns(address)
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) UnpackCertVerifiers(data []byte) (common.Address, error) {
	out, err := contractEigenDACertVerifierRouter.abi.Unpack("certVerifiers", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackCheckDACert is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9077193b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkDACert(bytes abiEncodedCert) view returns(uint8)
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) PackCheckDACert(abiEncodedCert []byte) []byte {
	enc, err := contractEigenDACertVerifierRouter.abi.Pack("checkDACert", abiEncodedCert)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckDACert is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9077193b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkDACert(bytes abiEncodedCert) view returns(uint8)
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) TryPackCheckDACert(abiEncodedCert []byte) ([]byte, error) {
	return contractEigenDACertVerifierRouter.abi.Pack("checkDACert", abiEncodedCert)
}

// UnpackCheckDACert is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9077193b.
//
// Solidity: function checkDACert(bytes abiEncodedCert) view returns(uint8)
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) UnpackCheckDACert(data []byte) (uint8, error) {
	out, err := contractEigenDACertVerifierRouter.abi.Unpack("checkDACert", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, nil
}

// PackGetCertVerifierAt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4a4ae0e2.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getCertVerifierAt(uint32 referenceBlockNumber) view returns(address)
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) PackGetCertVerifierAt(referenceBlockNumber uint32) []byte {
	enc, err := contractEigenDACertVerifierRouter.abi.Pack("getCertVerifierAt", referenceBlockNumber)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetCertVerifierAt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4a4ae0e2.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getCertVerifierAt(uint32 referenceBlockNumber) view returns(address)
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) TryPackGetCertVerifierAt(referenceBlockNumber uint32) ([]byte, error) {
	return contractEigenDACertVerifierRouter.abi.Pack("getCertVerifierAt", referenceBlockNumber)
}

// UnpackGetCertVerifierAt is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4a4ae0e2.
//
// Solidity: function getCertVerifierAt(uint32 referenceBlockNumber) view returns(address)
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) UnpackGetCertVerifierAt(data []byte) (common.Address, error) {
	out, err := contractEigenDACertVerifierRouter.abi.Unpack("getCertVerifierAt", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9d8ecd85.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function initialize(address initialOwner, uint32[] initABNs, address[] initCertVerifiers) returns()
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) PackInitialize(initialOwner common.Address, initABNs []uint32, initCertVerifiers []common.Address) []byte {
	enc, err := contractEigenDACertVerifierRouter.abi.Pack("initialize", initialOwner, initABNs, initCertVerifiers)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9d8ecd85.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function initialize(address initialOwner, uint32[] initABNs, address[] initCertVerifiers) returns()
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) TryPackInitialize(initialOwner common.Address, initABNs []uint32, initCertVerifiers []common.Address) ([]byte, error) {
	return contractEigenDACertVerifierRouter.abi.Pack("initialize", initialOwner, initABNs, initCertVerifiers)
}

// PackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function owner() view returns(address)
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) PackOwner() []byte {
	enc, err := contractEigenDACertVerifierRouter.abi.Pack("owner")
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
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) TryPackOwner() ([]byte, error) {
	return contractEigenDACertVerifierRouter.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) UnpackOwner(data []byte) (common.Address, error) {
	out, err := contractEigenDACertVerifierRouter.abi.Unpack("owner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackRenounceOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715018a6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceOwnership() returns()
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) PackRenounceOwnership() []byte {
	enc, err := contractEigenDACertVerifierRouter.abi.Pack("renounceOwnership")
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
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) TryPackRenounceOwnership() ([]byte, error) {
	return contractEigenDACertVerifierRouter.abi.Pack("renounceOwnership")
}

// PackTransferOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2fde38b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) PackTransferOwnership(newOwner common.Address) []byte {
	enc, err := contractEigenDACertVerifierRouter.abi.Pack("transferOwnership", newOwner)
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
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) TryPackTransferOwnership(newOwner common.Address) ([]byte, error) {
	return contractEigenDACertVerifierRouter.abi.Pack("transferOwnership", newOwner)
}

// ContractEigenDACertVerifierRouterCertVerifierAdded represents a CertVerifierAdded event raised by the ContractEigenDACertVerifierRouter contract.
type ContractEigenDACertVerifierRouterCertVerifierAdded struct {
	ActivationBlockNumber uint32
	CertVerifier          common.Address
	Raw                   *types.Log // Blockchain specific contextual infos
}

const ContractEigenDACertVerifierRouterCertVerifierAddedEventName = "CertVerifierAdded"

// ContractEventName returns the user-defined event name.
func (ContractEigenDACertVerifierRouterCertVerifierAdded) ContractEventName() string {
	return ContractEigenDACertVerifierRouterCertVerifierAddedEventName
}

// UnpackCertVerifierAddedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event CertVerifierAdded(uint32 indexed activationBlockNumber, address indexed certVerifier)
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) UnpackCertVerifierAddedEvent(log *types.Log) (*ContractEigenDACertVerifierRouterCertVerifierAdded, error) {
	event := "CertVerifierAdded"
	if log.Topics[0] != contractEigenDACertVerifierRouter.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ContractEigenDACertVerifierRouterCertVerifierAdded)
	if len(log.Data) > 0 {
		if err := contractEigenDACertVerifierRouter.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range contractEigenDACertVerifierRouter.abi.Events[event].Inputs {
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

// ContractEigenDACertVerifierRouterInitialized represents a Initialized event raised by the ContractEigenDACertVerifierRouter contract.
type ContractEigenDACertVerifierRouterInitialized struct {
	Version uint8
	Raw     *types.Log // Blockchain specific contextual infos
}

const ContractEigenDACertVerifierRouterInitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (ContractEigenDACertVerifierRouterInitialized) ContractEventName() string {
	return ContractEigenDACertVerifierRouterInitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint8 version)
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) UnpackInitializedEvent(log *types.Log) (*ContractEigenDACertVerifierRouterInitialized, error) {
	event := "Initialized"
	if log.Topics[0] != contractEigenDACertVerifierRouter.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ContractEigenDACertVerifierRouterInitialized)
	if len(log.Data) > 0 {
		if err := contractEigenDACertVerifierRouter.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range contractEigenDACertVerifierRouter.abi.Events[event].Inputs {
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

// ContractEigenDACertVerifierRouterOwnershipTransferred represents a OwnershipTransferred event raised by the ContractEigenDACertVerifierRouter contract.
type ContractEigenDACertVerifierRouterOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const ContractEigenDACertVerifierRouterOwnershipTransferredEventName = "OwnershipTransferred"

// ContractEventName returns the user-defined event name.
func (ContractEigenDACertVerifierRouterOwnershipTransferred) ContractEventName() string {
	return ContractEigenDACertVerifierRouterOwnershipTransferredEventName
}

// UnpackOwnershipTransferredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) UnpackOwnershipTransferredEvent(log *types.Log) (*ContractEigenDACertVerifierRouterOwnershipTransferred, error) {
	event := "OwnershipTransferred"
	if log.Topics[0] != contractEigenDACertVerifierRouter.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ContractEigenDACertVerifierRouterOwnershipTransferred)
	if len(log.Data) > 0 {
		if err := contractEigenDACertVerifierRouter.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range contractEigenDACertVerifierRouter.abi.Events[event].Inputs {
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

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], contractEigenDACertVerifierRouter.abi.Errors["ABNNotGreaterThanLast"].ID.Bytes()[:4]) {
		return contractEigenDACertVerifierRouter.UnpackABNNotGreaterThanLastError(raw[4:])
	}
	if bytes.Equal(raw[:4], contractEigenDACertVerifierRouter.abi.Errors["ABNNotInFuture"].ID.Bytes()[:4]) {
		return contractEigenDACertVerifierRouter.UnpackABNNotInFutureError(raw[4:])
	}
	if bytes.Equal(raw[:4], contractEigenDACertVerifierRouter.abi.Errors["InvalidCertLength"].ID.Bytes()[:4]) {
		return contractEigenDACertVerifierRouter.UnpackInvalidCertLengthError(raw[4:])
	}
	if bytes.Equal(raw[:4], contractEigenDACertVerifierRouter.abi.Errors["LengthMismatch"].ID.Bytes()[:4]) {
		return contractEigenDACertVerifierRouter.UnpackLengthMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], contractEigenDACertVerifierRouter.abi.Errors["RBNInFuture"].ID.Bytes()[:4]) {
		return contractEigenDACertVerifierRouter.UnpackRBNInFutureError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// ContractEigenDACertVerifierRouterABNNotGreaterThanLast represents a ABNNotGreaterThanLast error raised by the ContractEigenDACertVerifierRouter contract.
type ContractEigenDACertVerifierRouterABNNotGreaterThanLast struct {
	ActivationBlockNumber uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ABNNotGreaterThanLast(uint32 activationBlockNumber)
func ContractEigenDACertVerifierRouterABNNotGreaterThanLastErrorID() common.Hash {
	return common.HexToHash("0xfaf9cb693ddee49e99a1d7129f71d5d21abc4f92bccc0e02edfaca70adb8351b")
}

// UnpackABNNotGreaterThanLastError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ABNNotGreaterThanLast(uint32 activationBlockNumber)
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) UnpackABNNotGreaterThanLastError(raw []byte) (*ContractEigenDACertVerifierRouterABNNotGreaterThanLast, error) {
	out := new(ContractEigenDACertVerifierRouterABNNotGreaterThanLast)
	if err := contractEigenDACertVerifierRouter.abi.UnpackIntoInterface(out, "ABNNotGreaterThanLast", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ContractEigenDACertVerifierRouterABNNotInFuture represents a ABNNotInFuture error raised by the ContractEigenDACertVerifierRouter contract.
type ContractEigenDACertVerifierRouterABNNotInFuture struct {
	ActivationBlockNumber uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ABNNotInFuture(uint32 activationBlockNumber)
func ContractEigenDACertVerifierRouterABNNotInFutureErrorID() common.Hash {
	return common.HexToHash("0x5526abd48e9c714865dd2e2cb86e9955904725ef8b1e57b388bcdadfeb05c690")
}

// UnpackABNNotInFutureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ABNNotInFuture(uint32 activationBlockNumber)
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) UnpackABNNotInFutureError(raw []byte) (*ContractEigenDACertVerifierRouterABNNotInFuture, error) {
	out := new(ContractEigenDACertVerifierRouterABNNotInFuture)
	if err := contractEigenDACertVerifierRouter.abi.UnpackIntoInterface(out, "ABNNotInFuture", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ContractEigenDACertVerifierRouterInvalidCertLength represents a InvalidCertLength error raised by the ContractEigenDACertVerifierRouter contract.
type ContractEigenDACertVerifierRouterInvalidCertLength struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidCertLength()
func ContractEigenDACertVerifierRouterInvalidCertLengthErrorID() common.Hash {
	return common.HexToHash("0x75e774a059dc5d2eea2bbef60bb7b5348cdc66fef7e76ba7a4349c17420374ed")
}

// UnpackInvalidCertLengthError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidCertLength()
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) UnpackInvalidCertLengthError(raw []byte) (*ContractEigenDACertVerifierRouterInvalidCertLength, error) {
	out := new(ContractEigenDACertVerifierRouterInvalidCertLength)
	if err := contractEigenDACertVerifierRouter.abi.UnpackIntoInterface(out, "InvalidCertLength", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ContractEigenDACertVerifierRouterLengthMismatch represents a LengthMismatch error raised by the ContractEigenDACertVerifierRouter contract.
type ContractEigenDACertVerifierRouterLengthMismatch struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LengthMismatch()
func ContractEigenDACertVerifierRouterLengthMismatchErrorID() common.Hash {
	return common.HexToHash("0xff633a3803c58b9bc21e58efecee59f27e033cc0b1883fccb4969c76146fe60f")
}

// UnpackLengthMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LengthMismatch()
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) UnpackLengthMismatchError(raw []byte) (*ContractEigenDACertVerifierRouterLengthMismatch, error) {
	out := new(ContractEigenDACertVerifierRouterLengthMismatch)
	if err := contractEigenDACertVerifierRouter.abi.UnpackIntoInterface(out, "LengthMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ContractEigenDACertVerifierRouterRBNInFuture represents a RBNInFuture error raised by the ContractEigenDACertVerifierRouter contract.
type ContractEigenDACertVerifierRouterRBNInFuture struct {
	ReferenceBlockNumber uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RBNInFuture(uint32 referenceBlockNumber)
func ContractEigenDACertVerifierRouterRBNInFutureErrorID() common.Hash {
	return common.HexToHash("0x47fa9454ab3064f0cfe43e8ac7575ad819689ba71a0b4335ab64868763f2bdbb")
}

// UnpackRBNInFutureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RBNInFuture(uint32 referenceBlockNumber)
func (contractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouter) UnpackRBNInFutureError(raw []byte) (*ContractEigenDACertVerifierRouterRBNInFuture, error) {
	out := new(ContractEigenDACertVerifierRouterRBNInFuture)
	if err := contractEigenDACertVerifierRouter.abi.UnpackIntoInterface(out, "RBNInFuture", raw); err != nil {
		return nil, err
	}
	return out, nil
}
