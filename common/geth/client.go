package geth

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

const (
	RPC_SWITCH_TRIGGER = 3
)

var (
	FallbackGasTipCap       = big.NewInt(15000000000)
	ErrCannotGetECDSAPubKey = errors.New("ErrCannotGetECDSAPubKey")
	ErrTransactionFailed    = errors.New("ErrTransactionFailed")
)

type EthClient struct {
	*EthClientInstance // the first instance
	BackupInstances    []*EthClientInstance
	Controller         *FailoverController
}

func NewClient(config EthClientConfig, logger logging.Logger) (*EthClient, error) {
	controller := NewFailoverController(len(config.RPCURLBackup), RPC_SWITCH_TRIGGER, logger)

	primary, err := NewClientInstance(config, logger, controller)
	if err != nil {
		return nil, err
	}

	if len(config.PrivateKeyStringBackup) != 0 && len(config.PrivateKeyStringBackup) != len(config.RPCURLBackup) {
		return nil, fmt.Errorf("inconsistent number of %v private keys and %v url", len(config.PrivateKeyStringBackup), len(config.RPCURLBackup))
	}

	numBackup := len(config.RPCURLBackup)
	backups := make([]*EthClientInstance, numBackup)
	for i := 0; i < numBackup; i++ {
		// overwrite the default
		config.RPCURL = config.RPCURLBackup[i]
		if len(config.PrivateKeyStringBackup) != 0 {
			config.PrivateKeyString = config.PrivateKeyStringBackup[i]
		}

		backups[i], err = NewClientInstance(config, logger, controller)
		if err != nil {
			return nil, err
		}
	}

	logger.Info("Maintain", 1+numBackup, "Eth Client Instances")

	return &EthClient{
		EthClientInstance: primary,
		BackupInstances:   backups,
		Controller:        controller,
	}, nil
}

func (c *EthClient) GetEthClientInstance() *EthClientInstance {
	isPrimary, index := c.Controller.GetClientIndex()
	if isPrimary {
		return c.EthClientInstance
	} else {
		return c.BackupInstances[index]
	}
}
