# EigenDA Test Infrastructure

This package provides a modular, testcontainer-based infrastructure for EigenDA testing.

## Features

- **Modular Design**: Enable only the infrastructure components you need
- **Testcontainers Integration**: Automatic container lifecycle management  
- **Parallel Test Support**: Isolated containers for each test
- **Type-Safe Configuration**: Go structs instead of YAML/env files
- **Automatic Cleanup**: Containers destroyed after tests complete

## Infrastructure Components

### 1. Anvil Blockchain
- Foundry's Anvil for local Ethereum testing
- Configurable chain ID, block time, gas settings
- Pre-funded accounts with deterministic private keys
- Optional mainnet forking support

### 2. LocalStack 
- AWS service simulation (S3, DynamoDB, KMS, Secrets Manager)
- Configurable service selection
- Standard AWS SDK integration

### 3. Graph Node (Optional)
- Complete Graph Protocol stack with PostgreSQL and IPFS
- Subgraph deployment capabilities
- Ethereum chain integration

## Quick Start

### Minimal Setup (Anvil + LocalStack)
```go
func TestWithInfrastructure(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
    defer cancel()

    // Start minimal infrastructure
    manager, result, err := testinfra.StartMinimal(ctx)
    require.NoError(t, err)
    defer manager.Stop(ctx)

    // Use Anvil blockchain
    ethClient := ethclient.Dial(result.AnvilRPC)
    
    // Use LocalStack for AWS services
    awsConfig := manager.GetLocalStack().GetAWSConfig()
}
```

### Custom Configuration
```go
func TestWithCustomConfig(t *testing.T) {
    config := testinfra.DefaultConfig()
    config.Anvil.ChainID = 1337
    config.Anvil.Accounts = 5
    config.LocalStack.Services = []string{"s3", "dynamodb"}
    
    manager, result, err := testinfra.StartCustom(ctx, config)
    require.NoError(t, err)
    defer manager.Stop(ctx)
    
    // Your test code here...
}
```

## Configuration Options

### AnvilConfig
```go
type AnvilConfig struct {
    Enabled     bool     // Enable Anvil container
    ChainID     int      // Ethereum chain ID (default: 31337)
    BlockTime   int      // Seconds between blocks (0 = instant)
    GasLimit    uint64   // Block gas limit
    GasPrice    uint64   // Gas price (0 = free)
    Accounts    int      // Pre-funded accounts (default: 10)
    Mnemonic    string   // Custom mnemonic for deterministic accounts
    Fork        string   // Fork from this RPC URL
    ForkBlock   uint64   // Fork from specific block number
}
```

### LocalStackConfig
```go
type LocalStackConfig struct {
    Enabled  bool     // Enable LocalStack container
    Services []string // AWS services: s3, dynamodb, kms, secretsmanager
    Region   string   // AWS region (default: us-east-1)
    Debug    bool     // Enable debug logging
}
```

### GraphNodeConfig
```go
type GraphNodeConfig struct {
    Enabled       bool   // Enable Graph Node (disabled by default)
    PostgresDB    string // PostgreSQL database name
    PostgresUser  string // PostgreSQL username  
    PostgresPass  string // PostgreSQL password
    EthereumRPC   string // Ethereum RPC (auto-set to Anvil if available)
    IPFSEndpoint  string // IPFS endpoint (empty = embedded IPFS)
}
```

## Migration from inabox

### Before (inabox shell scripts)
```bash
# Terminal 1
anvil --host 0.0.0.0

# Terminal 2  
cd inabox/thegraph && docker compose up

# Terminal 3
LOCALSTACK_HOST=localhost.localstack.cloud:4570 localstack start

# Terminal 4
cd inabox && make exp && ./bin.sh start
```

### After (testinfra)
```go
func TestEigenDAFlow(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()

    // One line replaces all the above
    manager, result, err := testinfra.StartMinimal(ctx)
    require.NoError(t, err)
    defer manager.Stop(ctx)
    
    // All services ready to use
    runYourTestLogic(result.AnvilRPC, result.LocalStackURL)
}
```

## Integration Examples

### With EigenDA Services
```go
// Start infrastructure
manager, result, err := testinfra.StartMinimal(ctx)
require.NoError(t, err)
defer manager.Stop(ctx)

// Configure EigenDA services to use testcontainers
disperserConfig := disperser.Config{
    EthRPC: result.AnvilRPC,
    // ... other config
}

// Deploy contracts using Anvil
anvil := manager.GetAnvil()
deployerKey, _ := anvil.GetPrivateKey(0)
deployContracts(result.AnvilRPC, deployerKey)
```

### With AWS Services via LocalStack
```go
manager, result, err := testinfra.StartMinimal(ctx)
require.NoError(t, err)
defer manager.Stop(ctx)

localstack := manager.GetLocalStack()
awsConfig := localstack.GetAWSConfig()

// Create AWS clients pointing to LocalStack
s3Client := s3.NewFromConfig(aws.Config{
    Region:      awsConfig["AWS_DEFAULT_REGION"],
    Credentials: credentials.NewStaticCredentialsProvider(
        awsConfig["AWS_ACCESS_KEY_ID"],
        awsConfig["AWS_SECRET_ACCESS_KEY"],
        "",
    ),
    EndpointResolver: aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
        return aws.Endpoint{URL: awsConfig["AWS_ENDPOINT_URL"]}, nil
    }),
})
```

## Best Practices

1. **Use Context Timeouts**: Always set reasonable timeouts for container startup
2. **Defer Cleanup**: Always defer `manager.Stop()` to ensure containers are cleaned up
3. **Start Minimal**: Use `StartMinimal()` unless you specifically need Graph Node
4. **Parallel Tests**: Each test gets isolated containers - no shared state issues
5. **Error Handling**: Check errors from `Start*()` functions - container startup can fail

## Performance Considerations

- **Container Startup**: First run downloads images (~30s), subsequent runs are fast (~5s)
- **Parallel Tests**: Each test spawns new containers - consider `t.Parallel()` usage
- **Resource Usage**: Full infrastructure uses ~500MB RAM, minimal uses ~200MB
- **CI/CD**: Works great in GitHub Actions, no external dependencies

## Troubleshooting

### Container Startup Failures
```go
// Get container logs for debugging
manager, result, err := testinfra.StartMinimal(ctx)
if err != nil {
    if anvil := manager.GetAnvil(); anvil != nil {
        logs, _ := anvil.GetLogs(ctx)
        t.Logf("Anvil logs: %s", logs)
    }
}
```

### Network Issues
- Ensure Docker daemon is running
- Check that required ports aren't already in use
- On macOS, verify Docker Desktop resource limits

### Slow Tests
- Consider using `t.Parallel()` for independent tests
- Use `StartMinimal()` instead of `StartFull()` when possible  
- Cache container images in CI/CD environments

## Future Enhancements

- [ ] Support for multiple Anvil chains
- [ ] EigenDA service containers (disperser, node, etc.)
- [ ] Kubernetes testcontainers support
- [ ] Configuration validation and better error messages
- [ ] Performance optimizations for CI/CD

## Contributing

When adding new infrastructure components:
1. Create container wrapper in `containers/` package
2. Add configuration to main config structs  
3. Update `InfraManager` to orchestrate lifecycle
4. Add example tests and documentation
5. Consider backward compatibility with existing tests