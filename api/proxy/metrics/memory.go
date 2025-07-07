package metrics

import (
	"fmt"
	"sort"
	"sync"

	"github.com/ethereum-optimism/optimism/op-service/metrics"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

// fingerprint ... Construct a deterministic hash key for a label set
func fingerprint(labels []string) (common.Hash, error) {
	sort.Strings(labels) // in-place sort strings so keys are order agnostic

	encodedBytes, err := rlp.EncodeToBytes(labels)
	if err != nil {
		return common.Hash{}, err
	}

	hash := crypto.Keccak256Hash(encodedBytes)

	return hash, nil
}

// CountMap ... In memory representation of a prometheus Count metric type
type CountMap struct {
	m *sync.Map
}

// NewCountMap ... Init
func NewCountMap() *CountMap {
	return &CountMap{
		m: new(sync.Map),
	}
}

// insert ... increments or sets value associated with fingerprint
func (cm *CountMap) insert(labels ...string) error {
	key, err := fingerprint(labels)

	if err != nil {
		return err
	}

	// update or add count entry
	value, exists := cm.m.Load(key.Hex())
	if !exists {
		cm.m.Store(key.Hex(), uint64(1))
		return nil
	}
	uint64Val, ok := value.(uint64)
	if !ok {
		return fmt.Errorf("could not read uint64 from sync map")
	}

	cm.m.Store(key.Hex(), uint64Val+uint64(1))
	return nil
}

// Get ... fetches the value count associated with a deterministic label key
func (cm *CountMap) Get(labels ...string) (uint64, error) {
	key, err := fingerprint(labels)

	if err != nil {
		return 0, err
	}

	val, exists := cm.m.Load(key.Hex())
	if !exists {
		return 0, fmt.Errorf("value doesn't exist for key %s", key.String())
	}
	uint64Val, ok := val.(uint64)
	if !ok {
		return 0, fmt.Errorf("could not read uint64 from sync map")
	}

	return uint64Val, nil
}

// EmulatedMetricer ... allows for tracking count metrics in memory
// and is only used for E2E testing. This is needed since prometheus/client_golang doesn't provide
// an interface for reading the count values from the codified metric.
type EmulatedMetricer struct {
	HTTPServerRequestsTotal *CountMap
	// secondary metrics
	SecondaryRequestsTotal *CountMap
}

// NewEmulatedMetricer ... constructor
func NewEmulatedMetricer() *EmulatedMetricer {
	return &EmulatedMetricer{
		HTTPServerRequestsTotal: NewCountMap(),
		SecondaryRequestsTotal:  NewCountMap(),
	}
}

var _ Metricer = NewEmulatedMetricer()

// Document ... noop
func (n *EmulatedMetricer) Document() []metrics.DocumentedMetric {
	return nil
}

// RecordInfo ... noop
func (n *EmulatedMetricer) RecordInfo(_ string) {
}

// RecordUp ... noop
func (n *EmulatedMetricer) RecordUp() {
}

// RecordRPCServerRequest ... updates server requests counter associated with label fingerprint
func (n *EmulatedMetricer) RecordRPCServerRequest(method string) func(status, mode, ver string) {
	return func(_ string, mode string, _ string) {
		err := n.HTTPServerRequestsTotal.insert(method, mode)
		if err != nil { // panicking here is ok since this is only ran per E2E testing and never in server logic.
			panic(err)
		}
	}
}

// RecordSecondaryRequest ... updates secondary insertion counter associated with label fingerprint
func (n *EmulatedMetricer) RecordSecondaryRequest(x string, y string) func(status string) {
	return func(z string) {
		err := n.SecondaryRequestsTotal.insert(x, y, z)
		if err != nil {
			panic(err)
		}
	}
}
