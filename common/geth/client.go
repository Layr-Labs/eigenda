package geth

import (
	"errors"
	"fmt"
	"math/big"

	dacommon "github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

var (
	FallbackGasTipCap       = big.NewInt(15000000000)
	ErrCannotGetECDSAPubKey = errors.New("ErrCannotGetECDSAPubKey")
	ErrTransactionFailed    = errors.New("ErrTransactionFailed")
)

type EthClient struct {
	*EthClientBase  // the first instance
	BackupInstances []*EthClientBase
	Controller      *FailoverController
}

var _ dacommon.EthClient = (*EthClient)(nil)

func NewClient(config EthClientConfig, logger logging.Logger) (*EthClient, error) {
	controller := NewFailoverController(len(config.RPCURLBackup), config.RpcSwitchTrigger, logger)

	primary, err := NewClientInstance(config, logger, controller)
	if err != nil {
		return nil, err
	}

	if len(config.PrivateKeyStringBackup) != 0 && len(config.PrivateKeyStringBackup) != len(config.RPCURLBackup) {
		return nil, fmt.Errorf("inconsistent number of %v private keys and %v url", len(config.PrivateKeyStringBackup), len(config.RPCURLBackup))
	}

	numBackup := len(config.RPCURLBackup)
	backups := make([]*EthClientBase, numBackup)
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
		EthClientBase:   primary,
		BackupInstances: backups,
		Controller:      controller,
	}, nil
}

func (c *EthClient) GetEthClientInstance() dacommon.EthClient {
	isPrimary, index := c.Controller.GetClientIndex()
	if isPrimary {
		return c.EthClientBase
	} else {
		return c.BackupInstances[index]
	}
}
