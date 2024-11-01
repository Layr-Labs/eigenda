package api

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorFailoverErrorsIs(t *testing.T) {
	baseErr := fmt.Errorf("some error")
	failoverErr := NewErrorFailover(baseErr)
	wrappedFailoverErr := fmt.Errorf("some extra context: %w", failoverErr)
	if !errors.Is(wrappedFailoverErr, &ErrorFailover{}) {
		// do something...
	}

	// error comparison only checks the type of the error
	if !errors.Is(failoverErr, &ErrorFailover{}) {
		t.Error("should match ErrorFailover type")
	}

	// can also compare if failover error is wrapped
	wrapped := fmt.Errorf("wrapped: %w", failoverErr)
	if !errors.Is(wrapped, &ErrorFailover{}) {
		t.Error("should match ErrorFailover type even when wrapped")
	}

}

func TestErrorFailoverZeroValue(t *testing.T) {
	var failoverErr ErrorFailover
	if failoverErr.Error() != "Failover" {
		t.Error("should return 'Failover' for zero value")
	}
}
