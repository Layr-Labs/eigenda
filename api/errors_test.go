package api

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorFailoverErrorsIs(t *testing.T) {
	baseErr := fmt.Errorf("base error")
	failoverErr := NewErrorFailover(baseErr)
	otherFailoverErr := NewErrorFailover(fmt.Errorf("some other error"))
	wrappedFailoverErr := fmt.Errorf("wrapped: %w", failoverErr)

	if !errors.Is(failoverErr, failoverErr) {
		t.Error("should match itself")
	}

	if !errors.Is(failoverErr, baseErr) {
		t.Error("should match base error")
	}

	if errors.Is(failoverErr, fmt.Errorf("some other error")) {
		t.Error("should not match other errors")
	}

	if !errors.Is(failoverErr, otherFailoverErr) {
		t.Error("should match other failover error")
	}

	if !errors.Is(failoverErr, &ErrorFailover{}) {
		t.Error("should match ErrorFailover type")
	}

	if !errors.Is(wrappedFailoverErr, &ErrorFailover{}) {
		t.Error("should match ErrorFailover type even when wrapped")
	}

}

func TestErrorFailoverZeroValue(t *testing.T) {
	var failoverErr ErrorFailover
	if failoverErr.Error() != "Failover" {
		t.Error("should return 'Failover' for zero value")
	}
}
