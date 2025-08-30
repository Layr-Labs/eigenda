package gas_exhaustion_cert_meter

import (
	"context"
	"fmt"

	proxycommon "github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/common"
	blsapkregistry "github.com/Layr-Labs/eigenda/contracts/bindings/BLSApkRegistry"
	contractIEigenDADirectory "github.com/Layr-Labs/eigenda/contracts/bindings/IEigenDADirectory"
	opstateretriever "github.com/Layr-Labs/eigenda/contracts/bindings/OperatorStateRetriever"
	"github.com/Layr-Labs/eigenda/tools/gas_exhaustion_cert_meter/flags"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli"

	certVerifierBinding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifier"
)

type Config struct {
	Logger    logging.Logger
	EthClient *ethclient.Client

	OpStateRetrCaller    *opstateretriever.ContractOperatorStateRetrieverCaller
	BLSApkRegistryCaller *blsapkregistry.ContractBLSApkRegistryCaller
	CertVerifierCaller   *certVerifierBinding.ContractEigenDACertVerifierCaller

	CertVerifierAddr        gethcommon.Address
	RegistryCoordinatorAddr gethcommon.Address

	Ctx context.Context

	CertPath string
}

func GetAddressByName(
	ctx context.Context,
	client *ethclient.Client,
	directoryAddress gethcommon.Address,
	name string,
) (gethcommon.Address, error) {
	caller, err := contractIEigenDADirectory.NewContractIEigenDADirectoryCaller(directoryAddress, client)
	if err != nil {
		return gethcommon.Address{}, fmt.Errorf("failed to create EigenDA directory contract caller: %w", err)
	}

	operatorStateRetrieverAddr, err := caller.GetAddress0(&bind.CallOpts{Context: ctx}, name)
	if err != nil {
		return gethcommon.Address{}, fmt.Errorf("failed to get address for name %v: %w", name, err)
	}

	return operatorStateRetrieverAddr, nil
}

func ReadConfig(ctx *cli.Context, logger logging.Logger) (*Config, error) {

	rpcURL := ctx.GlobalString(flags.EthRpcUrlFlag.Name)
	ethContext := context.Background()
	client, err := ethclient.DialContext(ethContext, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("dial Ethereum node at %s: %w", rpcURL, err)
	}

	networkString := ctx.String(flags.NetworkFlag.Name)
	eigenDANetwork, err := proxycommon.EigenDANetworkFromString(networkString)
	if err != nil {
		return nil, fmt.Errorf("parse eigenDANetwork: %w", err)
	}

	directoryAddress := gethcommon.HexToAddress(eigenDANetwork.GetEigenDADirectory())

	operatorStateRetrieverAddr, err := GetAddressByName(
		ethContext, client, directoryAddress, "OPERATOR_STATE_RETRIEVER")
	if err != nil {
		return nil, err
	}

	blsApkRegistryAddr, err := GetAddressByName(ethContext, client, directoryAddress, "BLS_APK_REGISTRY")
	if err != nil {
		return nil, err
	}

	registryCoordinatorAddr, err := GetAddressByName(ethContext, client, directoryAddress, "REGISTRY_COORDINATOR")
	if err != nil {
		return nil, err
	}

	v3CertVerifierAddr, err := GetAddressByName(ethContext, client, directoryAddress, "CERT_VERIFIER")
	if err != nil {
		return nil, err
	}

	opStateRetrCaller, err := opstateretriever.NewContractOperatorStateRetrieverCaller(
		operatorStateRetrieverAddr, client)
	if err != nil {
		logger.Error("Failed to fetch OperatorStateRetriever contract", "err", err)
		return nil, fmt.Errorf("failed to create operator state retriever caller: %w", err)
	}

	blsApkRegistryCaller, err := blsapkregistry.NewContractBLSApkRegistryCaller(blsApkRegistryAddr, client)
	if err != nil {
		logger.Error("Failed to fetch NewContractBLSApkRegistry contract", "err", err)
		return nil, fmt.Errorf("failed to create BLS APK registry caller: %w", err)
	}

	certVerifierCaller, err := certVerifierBinding.NewContractEigenDACertVerifierCaller(v3CertVerifierAddr, client)
	if err != nil {
		return nil, fmt.Errorf("bind to verifier contract at %s: %w", v3CertVerifierAddr.Hex(), err)
	}

	return &Config{
		EthClient:               client,
		OpStateRetrCaller:       opStateRetrCaller,
		BLSApkRegistryCaller:    blsApkRegistryCaller,
		CertVerifierCaller:      certVerifierCaller,
		RegistryCoordinatorAddr: registryCoordinatorAddr,
		CertVerifierAddr:        v3CertVerifierAddr,
		CertPath:                ctx.GlobalString(flags.CertRlpFileFlag.Name),
		Logger:                  logger,
		Ctx:                     ethContext,
	}, nil
}

func NewConfig(ctx *cli.Context) (*Config, error) {
	loggerConfig, err := common.ReadLoggerCLIConfig(ctx, flags.FlagPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to read logger config: %w", err)
	}

	logger, err := common.NewLogger(loggerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	config, err := ReadConfig(ctx, logger)
	if err != nil {
		return nil, fmt.Errorf("cannot read config %w", err)
	}

	return config, nil
}
