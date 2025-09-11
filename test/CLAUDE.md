# CLAUDE.md - Testing Guidelines for EigenDA

> **Purpose** – Idiomatic Go testing patterns and best practices for the EigenDA codebase.
> Apply these consistently across all test files to improve maintainability, reliability, and speed

---

## 1. Core Testing Principles

| Principle | Description |
|-----------|-------------|
| **Test Independence** | Each test must be self-contained and not rely on execution order. |
| **Clear Failures** | Failure output should make it obvious what went wrong and where |
| **Resource Cleanup** | Containers, files, and connections must be cleaned up deterministically |
| **Subtests** | Use t.Run() for logical grouping and granular execution (go test -run) |
| **Descriptive Names** | Name tests for the behavior under test, not the implementation |

---

## 2. Anti-Patterns & Fixes

| Anti-Pattern | Problem | Prefer |
|--------------|--------------|-----------------|
| Global test state | Interdependent tests; order sensitivity | Test-scoped setup/fixtures |
| Silent failures | Hard to debug | Include meaningful failure messages |
| Missing cleanup | Resource leaks; flakiness | Use t.Cleanup() with timeouts |
| Testing internals | Brittle/overfitted tests | Test behavior via public interfaces |
| Fixed sleeps | Flaky under load | Polling with timeout/backoff |
| Unmarked helpers | Noisy stack traces | t.Helper() in all helpers |
| Error string matching | Brittle; breaks on wording changes | Check error types/sentinel errors |

---


### 3 Test Categorization

- Unit tests: Fast, isolated, no external dependencies
- Integration tests: Test component interactions, may use containers
- E2E tests: Full system tests, require complete environment


## 3. Setup Patterns

### Test-Scoped Setup (Preferred)

Use when setup is quick and isolation is required.

```go
func setupTest(t *testing.T) (*SomeContainer, *Config) {
    t.Helper()

    ctx := t.Context()
    container, err := NewContainer(ctx, options)
    require.NoError(t, err, "failed to start container")

    cfg := NewConfig()

    t.Cleanup(func() {
        logger.Info("Cleaning up test resources")
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        _ = container.Terminate(ctx)
    })

    return container, cfg
}

func TestSomething(t *testing.T) {
    container, cfg := setupTest(t)
    // ... test logic
}
```

### 3.2 TestMain for Expensive, Shareable Setup

Use when setup is slow and safe to share across the package:

EigenDA examples: inabox deployments; contract deployments on Anvil; multi-operator setups; subgraph indexing.
```go
var (
    anvilContainer      *testbed.AnvilContainer
    localstackContainer *testbed.LocalStackContainer
    testConfig          *deploy.Config
    contractAddresses   *ContractAddresses

    logger = testutils.GetLogger()
)

func TestMain(m *testing.M) {
var (
    anvilContainer      *testbed.AnvilContainer
    localstackContainer *testbed.LocalStackContainer
    testConfig          *deploy.Config
    contractAddresses   *ContractAddresses
)

func TestMain(m *testing.M) {
    if testing.Short() {
        fmt.Println("Skipping inabox deployment in short mode")
        os.Exit(0)
    }

    // Use Background() here — there is no test context in TestMain.
    ctx := context.Background()

    var err error
    anvilContainer, err = testbed.NewAnvilContainerWithOptions(ctx, testbed.AnvilOptions{
        ExposeHostPort: true, Logger: logger,
    })
    if err != nil { log.Fatal("Failed to start Anvil:", err) }

    localstackContainer, err = testbed.NewLocalStackContainerWithOptions(ctx, testbed.LocalStackOptions{
        ExposeHostPort: true, Services: []string{"s3", "dynamodb", "kms"}, Logger: logger,
    })
    if err != nil { log.Fatal("Failed to start LocalStack:", err) }

    testConfig = deploy.NewTestConfig("testname", "../../../")
    testConfig.Deployers[0].DeploySubgraphs = false

    fmt.Println("Deploying inabox experiment...")
    testConfig.DeployExperiment()

    contractAddresses = &ContractAddresses{
        ServiceManager:         testConfig.EigenDA.ServiceManager,
        OperatorStateRetriever: testConfig.EigenDA.OperatorStateRetriever,
        Churner:                testConfig.EigenDA.Churner,
    }

    code := m.Run()

    // Cleanup (best-effort)
    cctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    _ = anvilContainer.Terminate(cctx)
    _ = localstackContainer.Terminate(cctx)
    os.Exit(code)
}
```
Example tests reusing the shared environment:

```go
func TestOperatorRegistration(t *testing.T) {
    tx, err := createTransactorFromScratch(
        operatorKey,
        contractAddresses.OperatorStateRetriever,
        contractAddresses.ServiceManager,
        logger,
    )
    require.NoError(t, err, "create transactor")

    err = tx.RegisterOperator(t.Context(), signer, socket, quorums, key, salt, expiry)
    require.NoError(t, err, "operator registration should succeed")
}

func TestChurning(t *testing.T) {
    srv := churner.NewServer(testConfig, contractAddresses.Churner, logger, metrics)
    reply, err := srv.Churn(t.Context(), request)
    require.NoError(t, err, "churn should not error")
    require.NotNil(t, reply, "reply should not be nil")
}
```

## 4. Assertions

### 4.1 require vs assert

- **`require`**: Stop immediately on failure (use in setup or preconditions).
- **`assert`**: Continue to gather multiple failures in one run.

```go
// Use require when setup must succeed
container, err := StartContainer()
require.NoError(t, err, "failed to start container")

// Use assert for multiple checks that should all run
assert.Equal(t, expected1, actual1, "first check failed")
assert.Equal(t, expected2, actual2, "second check failed")
assert.Equal(t, expected3, actual3, "third check failed")
```

### 4.2 Meaningful Assertion Messages

Include messages that add context about what's being tested or why it matters. Avoid redundant messages that merely restate the assertion.

```go
// Bad - redundant messages
require.Equal(t, a, b, "a should equal b")
require.NoError(t, err, "err should be nil")

// Good - adds context
require.NoError(t, err, "failed to create signer")
require.Equal(t, expected, actual, "hash mismatch for valid signature")
require.True(t, isValid, "signature verification failed")
```

Messages are optional when the assertion itself is self-explanatory:
```go
// Acceptable - self-explanatory
require.NotNil(t, client)
require.Len(t, results, 3)
```

---

## 5. Organizing Tests with Subtests

### 5.1 Table-Driven Tests

```go
func TestCalculation(t *testing.T) {
    tests := []struct {
        name     string
        input    int
        expected int
    }{
        {"zero", 0, 0},
        {"positive", 5, 25},
        {"negative", -3, 9},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Calculate(tt.input)
            require.Equal(t, tt.expected, result, 
                "Calculate(%d) = %d, want %d", tt.input, result, tt.expected)
        })
    }
}
```

### 5.2 Scenario-Based Subtests

```go
func TestSignature(t *testing.T) {
    ctx := t.Context()
    signer := setupSigner(t)
    
    t.Run("valid_signature", func(t *testing.T) {
        // Test valid case
    })
    
    t.Run("corrupted_signature", func(t *testing.T) {
        // Test corrupted signature
    })
    
    t.Run("modified_request", func(t *testing.T) {
        // Test modified request
    })
}
```

---

## 6. Context and Cleanup Patterns

### 6.1 Prefer t.Context()

```go
// Bad
ctx := context.Background()

// Good
ctx := t.Context()
```

Exception: In TestMain, use context.Background() (no test context exists).

```go
func TestMain(m *testing.M) {
    // Must use context.Background() here - no test context available
    ctx := context.Background()
    container, err := startContainer(ctx)
    // ...
}

func TestSomething(t *testing.T) {
    // Use t.Context() in actual test functions
    ctx := t.Context()
    // ...
}
```

### 6.2 Cleanup with Timeouts

```go
t.Cleanup(func() {
    logger.Info("Stopping container")
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := container.Stop(ctx); err != nil {
        logger.Warn("Failed to stop container", "error", err)
    }
})
```

---

## 7. Test Helpers & Utilities

### 7.1 Mark Helpers

```go
func setupTestEnvironment(t *testing.T) *Environment {
    t.Helper() // Improves stack traces on failure
    
    env, err := NewEnvironment()
    require.NoError(t, err, "failed to create environment")
    return env
}
```

Using t.Helper() will make it so that the stack trace will point to the actual test function that called it — not the helper itself.

### 7.2 Common Test Utilities

Create reusable test utilities for common operations. Store them in the `test` package.

```go
func requireEventuallyTrue(t *testing.T, condition func() bool, timeout time.Duration, msg string) {
    t.Helper()
    
    deadline := time.Now().Add(timeout)
    for time.Now().Before(deadline) {
        if condition() {
            return
        }
        time.Sleep(100 * time.Millisecond)
    }
    t.Fatalf("Condition not met within %v: %s", timeout, msg)
}
```

---

## 8. Containers & Integration

### 8.1 Container Management

We use the testbed package for container management and integration tests. It utilizes https://testcontainers.com/ for the 
container lifecycle and Docker for container runtime.

```go
func setupLocalStack(t *testing.T) *testbed.LocalStackContainer {
    t.Helper()
    
    ctx := t.Context()
    container, err := testbed.NewLocalStackContainerWithOptions(ctx, testbed.LocalStackOptions{
        ExposeHostPort: true,
        HostPort:       "4579",
        Services:       []string{"kms", "s3"},
        Logger:         logger,
    })
    require.NoError(t, err, "failed to start LocalStack")
    
    t.Cleanup(func() {
        logger.Info("Stopping LocalStack")
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        _ = container.Terminate(ctx)
    })
    
    return container
}
```

## 9. Test Documentation

### 9.1 Document Complex Test Scenarios

```go
func TestComplexScenario(t *testing.T) {
    // This test verifies that when a node receives a signature request
    // while in a degraded state, it properly queues the request and
    // processes it once the node recovers.
    
    // Setup: Create a node in degraded state
    node := setupDegradedNode(t)
    
    // When: Send signature request
    response := node.SignRequest(request)
    
    // Then: Verify queued response
    require.Equal(t, StatusQueued, response.Status)
    
    // And: Verify processing after recovery
    node.Recover()
    requireEventually(t, func() bool {
        return node.IsRequestProcessed(request.ID)
    }, 5*time.Second, "request should be processed after recovery")
}
```

---

## 12. Logging in Tests

### 12.1 Test Logger Usage

Use EigenDA’s test logger for consistent formatting and behavior.

```go
import (
    "github.com/Layr-Labs/eigenda/common/testutils"
)

var logger = testutils.GetLogger()
```

Notes:
- Fatal() logs and terminates (use in TestMain setup failures).
- Logger defaults to text with colors; includes source info and debug level.

```go
var logger = testutils.GetLogger()

func TestMain(m *testing.M) {
    // Use logger.Fatal for setup failures in TestMain
    container, err := startContainer()
    if err != nil {
        logger.Fatal("Failed to start container:", err)
    }
    
    // Use logger for informational messages
    logger.Info("Container started successfully")
    
    code := m.Run()
    
    // Cleanup...
    if container != nil {
        if err := container.Terminate(context.Background()); err != nil {
            logger.Error("Failed to terminate container:", err)
        } else {
            logger.Info("Container terminated successfully")
        }
    }

    os.Exit(code)
}
```

---

## 13. CI/CD Considerations

If something should only run in CI or be skipped in CI, use the `$CI` environment variable:

```go
func TestIntegration(t *testing.T) {
    if os.Getenv("CI") == "" {
        t.Skip("skipping integration test outside CI")
    }
    // ... integration test logic ...
}
```

Note: Github actions automatically set the `$CI` environment variable.
