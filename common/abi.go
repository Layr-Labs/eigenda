package common

import (
	_ "embed"

	"github.com/ethereum/go-ethereum/crypto"
)

//go:embed abis/EigenDAServiceManager.json
var ServiceManagerAbi []byte

//go:embed abis/RegistryCoordinator.json
var RegistryCoordinatorAbi []byte

var BatchConfirmedEventSigHash = crypto.Keccak256Hash([]byte("BatchConfirmed(bytes32,uint32)"))
var OperatorSocketUpdateEventSigHash = crypto.Keccak256Hash([]byte("OperatorSocketUpdate(bytes32,string)"))

// TODO: consider adding deregistration for limiting size of socket map
// var OperatorDeregisteredEventSigHash = crypto.Keccak256Hash([]byte("OperatorDeregistered(address,bytes32)"))
