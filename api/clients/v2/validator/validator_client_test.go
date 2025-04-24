package validator

import (
	"testing"
)

// TODO
//  - basic workflow
//  - test with pessimism factors of 1
//  - test with pessimism factors of 8
//  - test with pessimism factors between 1 and 8
//  - test what happens when a pessimistic download timeout triggers
//  - test what happens when a long download timeout triggers
//  - test what happens when chunks are invalid
//  - pessimistic timeout -> slow download eventually finishes -> validation on a chunk fails
//  - respecting thread pool limits

func TestBasicWorkflow(t *testing.T) {
	//rand := testrandom.NewTestRandom()
	//start := rand.Time()
	//
	//logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	//require.NoError(t, err)
	//
	//fakeClock := atomic.Pointer[time.Time]{}
	//fakeClock.Store(&start)
	//
	//config := DefaultClientConfig()
	//config.ControlLoopPeriod = 50 * time.Microsecond
	//config.timeSource = func() time.Time {
	//	return *fakeClock.Load()
	//}

	//client := NewValidatorClient(logger, nil, nil, nil, config)

	// TODO

}
