// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractEigenDACertVerifier

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

// EigenDATypesV1SecurityThresholds is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV1SecurityThresholds struct {
	ConfirmationThreshold uint8
	AdversaryThreshold    uint8
}

// ContractEigenDACertVerifierMetaData contains all meta data concerning the ContractEigenDACertVerifier contract.
var ContractEigenDACertVerifierMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"initEigenDAThresholdRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"},{\"name\":\"initEigenDASignatureVerifier\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"},{\"name\":\"initSecurityThresholds\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.SecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"initQuorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"certVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"checkDACert\",\"inputs\":[{\"name\":\"abiEncodedCert\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDASignatureVerifier\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"securityThresholds\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.SecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"semver\",\"inputs\":[],\"outputs\":[{\"name\":\"major\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"minor\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"patch\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"pure\"},{\"type\":\"error\",\"name\":\"InvalidSecurityThresholds\",\"inputs\":[]}]",
	ID:  "ContractEigenDACertVerifier",
}

// ContractEigenDACertVerifier is an auto generated Go binding around an Ethereum contract.
type ContractEigenDACertVerifier struct {
	abi abi.ABI
}

// NewContractEigenDACertVerifier creates a new instance of ContractEigenDACertVerifier.
func NewContractEigenDACertVerifier() *ContractEigenDACertVerifier {
	parsed, err := ContractEigenDACertVerifierMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ContractEigenDACertVerifier{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ContractEigenDACertVerifier) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address initEigenDAThresholdRegistry, address initEigenDASignatureVerifier, (uint8,uint8) initSecurityThresholds, bytes initQuorumNumbersRequired) returns()
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) PackConstructor(initEigenDAThresholdRegistry common.Address, initEigenDASignatureVerifier common.Address, initSecurityThresholds EigenDATypesV1SecurityThresholds, initQuorumNumbersRequired []byte) []byte {
	enc, err := contractEigenDACertVerifier.abi.Pack("", initEigenDAThresholdRegistry, initEigenDASignatureVerifier, initSecurityThresholds, initQuorumNumbersRequired)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackCertVersion is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2ead0b96.
//
// Solidity: function certVersion() pure returns(uint8)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) PackCertVersion() []byte {
	enc, err := contractEigenDACertVerifier.abi.Pack("certVersion")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackCertVersion is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2ead0b96.
//
// Solidity: function certVersion() pure returns(uint8)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) UnpackCertVersion(data []byte) (uint8, error) {
	out, err := contractEigenDACertVerifier.abi.Unpack("certVersion", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, err
}

// PackCheckDACert is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9077193b.
//
// Solidity: function checkDACert(bytes abiEncodedCert) view returns(uint8)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) PackCheckDACert(abiEncodedCert []byte) []byte {
	enc, err := contractEigenDACertVerifier.abi.Pack("checkDACert", abiEncodedCert)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackCheckDACert is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9077193b.
//
// Solidity: function checkDACert(bytes abiEncodedCert) view returns(uint8)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) UnpackCheckDACert(data []byte) (uint8, error) {
	out, err := contractEigenDACertVerifier.abi.Unpack("checkDACert", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, err
}

// PackEigenDASignatureVerifier is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xefd4532b.
//
// Solidity: function eigenDASignatureVerifier() view returns(address)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) PackEigenDASignatureVerifier() []byte {
	enc, err := contractEigenDACertVerifier.abi.Pack("eigenDASignatureVerifier")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackEigenDASignatureVerifier is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xefd4532b.
//
// Solidity: function eigenDASignatureVerifier() view returns(address)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) UnpackEigenDASignatureVerifier(data []byte) (common.Address, error) {
	out, err := contractEigenDACertVerifier.abi.Unpack("eigenDASignatureVerifier", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackEigenDAThresholdRegistry is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf8c66814.
//
// Solidity: function eigenDAThresholdRegistry() view returns(address)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) PackEigenDAThresholdRegistry() []byte {
	enc, err := contractEigenDACertVerifier.abi.Pack("eigenDAThresholdRegistry")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackEigenDAThresholdRegistry is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf8c66814.
//
// Solidity: function eigenDAThresholdRegistry() view returns(address)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) UnpackEigenDAThresholdRegistry(data []byte) (common.Address, error) {
	out, err := contractEigenDACertVerifier.abi.Unpack("eigenDAThresholdRegistry", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackQuorumNumbersRequired is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) PackQuorumNumbersRequired() []byte {
	enc, err := contractEigenDACertVerifier.abi.Pack("quorumNumbersRequired")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackQuorumNumbersRequired is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) UnpackQuorumNumbersRequired(data []byte) ([]byte, error) {
	out, err := contractEigenDACertVerifier.abi.Unpack("quorumNumbersRequired", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, err
}

// PackSecurityThresholds is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x21b9b2fb.
//
// Solidity: function securityThresholds() view returns((uint8,uint8))
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) PackSecurityThresholds() []byte {
	enc, err := contractEigenDACertVerifier.abi.Pack("securityThresholds")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackSecurityThresholds is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x21b9b2fb.
//
// Solidity: function securityThresholds() view returns((uint8,uint8))
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) UnpackSecurityThresholds(data []byte) (EigenDATypesV1SecurityThresholds, error) {
	out, err := contractEigenDACertVerifier.abi.Unpack("securityThresholds", data)
	if err != nil {
		return *new(EigenDATypesV1SecurityThresholds), err
	}
	out0 := *abi.ConvertType(out[0], new(EigenDATypesV1SecurityThresholds)).(*EigenDATypesV1SecurityThresholds)
	return out0, err
}

// PackSemver is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcda493c8.
//
// Solidity: function semver() pure returns(uint8 major, uint8 minor, uint8 patch)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) PackSemver() []byte {
	enc, err := contractEigenDACertVerifier.abi.Pack("semver")
	if err != nil {
		panic(err)
	}
	return enc
}

// SemverOutput serves as a container for the return parameters of contract
// method Semver.
type SemverOutput struct {
	Major uint8
	Minor uint8
	Patch uint8
}

// UnpackSemver is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xcda493c8.
//
// Solidity: function semver() pure returns(uint8 major, uint8 minor, uint8 patch)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) UnpackSemver(data []byte) (SemverOutput, error) {
	out, err := contractEigenDACertVerifier.abi.Unpack("semver", data)
	outstruct := new(SemverOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Major = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.Minor = *abi.ConvertType(out[1], new(uint8)).(*uint8)
	outstruct.Patch = *abi.ConvertType(out[2], new(uint8)).(*uint8)
	return *outstruct, err

}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], contractEigenDACertVerifier.abi.Errors["InvalidSecurityThresholds"].ID.Bytes()[:4]) {
		return contractEigenDACertVerifier.UnpackInvalidSecurityThresholdsError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// ContractEigenDACertVerifierInvalidSecurityThresholds represents a InvalidSecurityThresholds error raised by the ContractEigenDACertVerifier contract.
type ContractEigenDACertVerifierInvalidSecurityThresholds struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidSecurityThresholds()
func ContractEigenDACertVerifierInvalidSecurityThresholdsErrorID() common.Hash {
	return common.HexToHash("0x08a69975c4c065dd20db258fd793a9eb4231cd659928ecfc755e5cc8047fe11b")
}

// UnpackInvalidSecurityThresholdsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidSecurityThresholds()
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) UnpackInvalidSecurityThresholdsError(raw []byte) (*ContractEigenDACertVerifierInvalidSecurityThresholds, error) {
	out := new(ContractEigenDACertVerifierInvalidSecurityThresholds)
	if err := contractEigenDACertVerifier.abi.UnpackIntoInterface(out, "InvalidSecurityThresholds", raw); err != nil {
		return nil, err
	}
	return out, nil
}
