package containers

// AnvilConfig configures the Anvil blockchain container
type AnvilConfig struct {
	Enabled   bool   `json:"enabled"`
	ChainID   int    `json:"chain_id"`
	BlockTime int    `json:"block_time"` // seconds between blocks, 0 for instant mining
	GasLimit  uint64 `json:"gas_limit"`
	GasPrice  uint64 `json:"gas_price"`
	Accounts  int    `json:"accounts"`   // number of pre-funded accounts
	Mnemonic  string `json:"mnemonic"`   // custom mnemonic for deterministic accounts
	Fork      string `json:"fork"`       // fork from this RPC URL
	ForkBlock uint64 `json:"fork_block"` // fork from specific block
}

// LocalStackConfig configures the LocalStack AWS simulation container
type LocalStackConfig struct {
	Enabled  bool     `json:"enabled"`
	Services []string `json:"services"` // AWS services to enable: s3, dynamodb, kms, secretsmanager
	Region   string   `json:"region"`
	Debug    bool     `json:"debug"`
}

// GraphNodeConfig configures The Graph node container
type GraphNodeConfig struct {
	Enabled      bool   `json:"enabled"`
	PostgresDB   string `json:"postgres_db"`
	PostgresUser string `json:"postgres_user"`
	PostgresPass string `json:"postgres_pass"`
	EthereumRPC  string `json:"ethereum_rpc"` // will be set to Anvil RPC if Anvil is enabled
	IPFSEndpoint string `json:"ipfs_endpoint"`
}
