# Incremental PR Plan for Unified Payments Architecture

## **Phase 1: Foundation Interfaces (3 PRs)**

### **PR #1: Core Payment Interfaces**
```go
// File: core/meterer/interfaces.go
type AccountLedger interface {
    CreateDebit(debitSlip DebitSlip) error
    CreateReservationPaymentHeader(accountID string, timestamp time.Time, 
        transactionID string, paymentUsageType PaymentUsageType) (*PaymentHeader, error)
    GetAccountBalance(accountID string) (*AccountBalance, error)
}

type PaymentUsageType int
const (
    ReservationPayment PaymentUsageType = iota
    OnDemandPayment
)

type DebitSlip struct {
    PaymentHeader    PaymentHeader
    PaymentUsageType PaymentUsageType
    SourceID         string
    Amount           *big.Int
}

type PaymentHeader struct {
    AccountID     string
    Timestamp     time.Time
    TransactionID string
    UsageType     PaymentUsageType
}
```
**Impact**: Zero - pure interface definitions
**Files**: `core/meterer/interfaces.go`

### **PR #2: ThroughputMeter Interface**
```go
// File: core/meterer/throughput.go
type ThroughputMeter interface {
    RecordThroughput(debitSlip DebitSlip) error
    GetCurrentUsage() (*ThroughputStats, error)
    IsThrottled(accountID string) bool
}

type ThroughputStats struct {
    GlobalUsage     uint64
    AccountUsage    map[string]uint64
    ThrottledAccounts []string
}
```
**Impact**: Zero - interface only
**Files**: `core/meterer/throughput.go`

### **PR #3: BatchPaymentProcessor Interface**
```go
// File: core/meterer/batch_processor.go
type BatchPaymentProcessor interface {
    DebitBatch(debitSlips []DebitSlip) error
    ProcessBatch(batch PaymentBatch) (*BatchResult, error)
    GetAccountLedgers() map[string]AccountLedger
}

type PaymentBatch struct {
    BatchID    string
    DebitSlips []DebitSlip
    Timestamp  time.Time
}

type BatchResult struct {
    SuccessCount int
    FailedSlips  []DebitSlip
    Errors       []error
}
```
**Impact**: Zero - interface only
**Files**: `core/meterer/batch_processor.go`

## **Phase 2: Adapter Implementation (4 PRs)**

### **PR #4: DynamoAccountLedger Adapter**
```go
// File: core/meterer/dynamo_account_ledger.go
type DynamoAccountLedger struct {
    store MeteringStore // Existing interface
    mutex sync.RWMutex
}

func (d *DynamoAccountLedger) CreateDebit(debitSlip DebitSlip) error {
    // Delegate to existing store.UpdateCumulativePayment() or store.MeterRequest()
    return d.delegateToExistingLogic(debitSlip)
}
```
**Impact**: Low - wraps existing `DynamoMeteringStore`
**Files**: `core/meterer/dynamo_account_ledger.go`
**Dependencies**: PR #1

### **PR #5: EphemeralAccountLedger Implementation**
```go
// File: core/meterer/ephemeral_account_ledger.go
type EphemeralAccountLedger struct {
    accounts map[string]*AccountState
    mutex    sync.RWMutex
}

type AccountState struct {
    Balance          *big.Int
    ReservationUsage map[uint32]*PeriodRecord
    CumulativePayment *big.Int
}
```
**Impact**: Low - new in-memory implementation for testing
**Files**: `core/meterer/ephemeral_account_ledger.go`
**Dependencies**: PR #1

### **PR #6: ThroughputMeter Implementation**
```go
// File: core/meterer/throughput_meter_impl.go
type ThroughputMeterImpl struct {
    meterer *Meterer // Existing meterer
}

func (t *ThroughputMeterImpl) RecordThroughput(debitSlip DebitSlip) error {
    // Delegate to existing meterer rate limiting logic
    return t.meterer.checkRateLimits(debitSlip)
}
```
**Impact**: Low - wraps existing rate limiting
**Files**: `core/meterer/throughput_meter_impl.go`
**Dependencies**: PR #2

### **PR #7: BatchPaymentProcessor Implementation**
```go
// File: core/meterer/batch_processor_impl.go
type BatchPaymentProcessorImpl struct {
    accountLedgers map[string]AccountLedger
    throughputMeter ThroughputMeter
    mutex          sync.RWMutex
}

func (b *BatchPaymentProcessorImpl) DebitBatch(debitSlips []DebitSlip) error {
    // Group by account, process in parallel
    return b.processDebitSlipsInBatch(debitSlips)
}
```
**Impact**: Medium - new batch coordination logic
**Files**: `core/meterer/batch_processor_impl.go`
**Dependencies**: PR #1, #2, #3

## **Phase 3: Integration Layer (3 PRs)**

### **PR #8: Meterer Integration**
```go
// File: core/meterer/meterer.go
type Meterer struct {
    // Existing fields...
    
    // New unified interface
    accountLedger AccountLedger
    throughputMeter ThroughputMeter
    batchProcessor BatchPaymentProcessor
}

// Add new method alongside existing ones
func (m *Meterer) MeterRequestUnified(debitSlip DebitSlip) error {
    return m.batchProcessor.DebitBatch([]DebitSlip{debitSlip})
}
```
**Impact**: Low - adds new methods without changing existing ones
**Files**: `core/meterer/meterer.go`
**Dependencies**: PR #4, #5, #6, #7

### **PR #9: Client Accountant Integration**
```go
// File: api/clients/v2/accountant.go
func (accountant *Accountant) CreateDebitSlip(
    paymentHeader PaymentHeader,
    usageType PaymentUsageType,
    sourceID string,
) (*DebitSlip, error) {
    // Convert existing PaymentMetadata to DebitSlip
    return &DebitSlip{
        PaymentHeader:    paymentHeader,
        PaymentUsageType: usageType,
        SourceID:         sourceID,
        Amount:           accountant.calculateAmount(paymentHeader),
    }, nil
}
```
**Impact**: Low - adds new method alongside existing ones
**Files**: `api/clients/v2/accountant.go`
**Dependencies**: PR #1

### **PR #10: Disperser API Integration**
```go
// File: disperser/apiserver/server_v2.go
func (s *DispersalServerV2) disperseBlobUnified(
    ctx context.Context,
    blob *core.Blob,
    debitSlip *meterer.DebitSlip,
) (*disperser_rpc.BlobStatus, error) {
    // Use new unified path alongside existing logic
    return s.disperseBlobWithUnifiedPayments(ctx, blob, debitSlip)
}
```
**Impact**: Medium - new API path without breaking existing
**Files**: `disperser/apiserver/server_v2.go`
**Dependencies**: PR #8

## **Phase 4: Migration & Testing (3 PRs)**

### **PR #11: Comprehensive Testing Suite**
```go
// File: core/meterer/unified_payments_test.go
func TestUnifiedPaymentFlow(t *testing.T) {
    // Integration tests for complete flow
    ephemeralLedger := NewEphemeralAccountLedger()
    throughputMeter := NewThroughputMeterImpl()
    batchProcessor := NewBatchPaymentProcessorImpl(ephemeralLedger, throughputMeter)
    
    // Test reservation and on-demand flows
}
```
**Impact**: Zero - test-only changes
**Files**: `core/meterer/*_test.go`
**Dependencies**: All previous PRs

### **PR #12: Migration Utilities**
```go
// File: core/meterer/migration.go
type PaymentDataMigrator struct {
    oldStore MeteringStore
    newLedger AccountLedger
}

func (m *PaymentDataMigrator) MigrateAccountData(accountID string) error {
    // Migrate existing DynamoDB data to new format
    return m.migrateAccount(accountID)
}
```
**Impact**: Low - optional migration tools
**Files**: `core/meterer/migration.go`
**Dependencies**: PR #4, #5

### **PR #13: Feature Flag & Rollout**
```go
// File: core/meterer/meterer.go
func (m *Meterer) MeterRequest(
    ctx context.Context,
    accountID string,
    paymentMetadata core.PaymentMetadata,
    symbolsCharged uint64,
    quorumNumbers []uint8,
) error {
    if m.config.UseUnifiedPayments {
        debitSlip := m.convertToDebitSlip(paymentMetadata, symbolsCharged)
        return m.MeterRequestUnified(debitSlip)
    }
    
    // Existing logic unchanged
    return m.meterRequestLegacy(ctx, accountID, paymentMetadata, symbolsCharged, quorumNumbers)
}
```
**Impact**: Low - feature flag controls rollout
**Files**: `core/meterer/meterer.go`, `core/meterer/config.go`
**Dependencies**: PR #8

## **Phase 5: Deprecation & Cleanup (2 PRs)**

### **PR #14: Legacy Method Deprecation**
```go
// Mark existing methods as deprecated
// @deprecated Use MeterRequestUnified instead
func (m *Meterer) MeterReservationRequest(...) error {
    return m.MeterRequestUnified(convertToDebitSlip(...))
}
```
**Impact**: Medium - requires coordination with client updates
**Dependencies**: Full rollout of unified system

### **PR #15: Legacy Code Removal**
**Impact**: High - removes old implementation
**Dependencies**: All clients migrated to unified APIs

## **Risk Mitigation Strategy**

### **Per-PR Safety Measures**
- **Interface PRs**: Zero risk - no implementation changes
- **Adapter PRs**: Low risk - delegate to existing proven logic  
- **Integration PRs**: Feature flags enable gradual rollout
- **Testing PRs**: Comprehensive coverage before migration

### **Rollback Plan**
- Each PR maintains backward compatibility
- Feature flags allow instant rollback
- Existing APIs remain unchanged until Phase 5

### **Testing Strategy**
- Unit tests for each new component
- Integration tests for complete flows
- Performance benchmarks vs. existing implementation
- Chaos testing for error scenarios

## **Abstraction Analysis Summary**

### **Strong Alignment Areas**
1. **Payment Processing Flow** - Current `Meterer.MeterRequest()` maps to `BatchPaymentProcessor.DebitBatch()`
2. **Account Tracking** - Current `DynamoMeteringStore` can implement `AccountLedger` interface
3. **Rate Limiting** - Current token bucket logic aligns with `ThroughputMeter`

### **Key Benefits**
- **80% component alignment** between current and proposed systems
- **Direct interface mapping** for core payment processing logic  
- **Incremental migration path** preserving existing functionality
- **Enhanced testability** through interface abstraction
- **Improved maintainability** via unified payment processing flow

### **Success Metrics**
- All PRs maintain backward compatibility until Phase 5
- Feature flags enable safe rollout and rollback
- Comprehensive testing validates each increment
- Performance remains equivalent or improves
- Code complexity reduces after migration completion