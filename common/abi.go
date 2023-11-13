package common

import (
	_ "embed"

	"github.com/ethereum/go-ethereum/crypto"
)

//go:embed abis/EigenDAServiceManager.json
var ServiceManagerAbi []byte

var BatchConfirmedEventSigHash = crypto.Keccak256Hash([]byte("BatchConfirmed(bytes32,uint32,uint96)"))
