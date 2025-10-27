package api

import (
	"testing"

	"github.com/Layr-Labs/eigenda/api/grpc/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestNewErrorResourceExhaustedWithRateLimitDetails(t *testing.T) {
	// Test creating a rate limit error with details
	msg := "request ratelimited: System throughput rate limit for quorum 0"
	rateLimitType := common.RateLimitType_SYSTEM_THROUGHPUT_RATE_LIMIT
	quorumId := uint32(0)
	retryAfterSeconds := uint32(30)
	context := "Rate limit exceeded for System throughput rate limit"

	err := NewErrorResourceExhaustedWithRateLimitDetails(msg, rateLimitType, quorumId, retryAfterSeconds, context)

	// Check that it's a gRPC error with RESOURCE_EXHAUSTED code
	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("Expected gRPC status error")
	}

	if st.Code() != codes.ResourceExhausted {
		t.Errorf("Expected RESOURCE_EXHAUSTED code, got %v", st.Code())
	}

	if st.Message() != msg {
		t.Errorf("Expected message %q, got %q", msg, st.Message())
	}

	// Check that details are present
	details := st.Details()
	if len(details) != 1 {
		t.Errorf("Expected 1 detail, got %d", len(details))
	}

	// Try to extract the rate limit details
	anyDetail := details[0]
	var rateLimitDetails common.RateLimitDetails
	if err := anypb.UnmarshalTo(anyDetail, &rateLimitDetails, nil); err != nil {
		t.Fatalf("Failed to unmarshal rate limit details: %v", err)
	}

	if rateLimitDetails.RateLimitType != rateLimitType {
		t.Errorf("Expected rate limit type %v, got %v", rateLimitType, rateLimitDetails.RateLimitType)
	}

	if rateLimitDetails.QuorumId != quorumId {
		t.Errorf("Expected quorum ID %d, got %d", quorumId, rateLimitDetails.QuorumId)
	}

	if rateLimitDetails.RetryAfterSeconds != retryAfterSeconds {
		t.Errorf("Expected retry after seconds %d, got %d", retryAfterSeconds, rateLimitDetails.RetryAfterSeconds)
	}

	if rateLimitDetails.Context != context {
		t.Errorf("Expected context %q, got %q", context, rateLimitDetails.Context)
	}
}

func TestNewErrorResourceExhaustedWithRateLimitDetailsFallback(t *testing.T) {
	// Test fallback behavior when details creation fails
	// This is hard to test directly, but we can test that the function doesn't panic
	msg := "test error"
	rateLimitType := common.RateLimitType_RATE_LIMIT_TYPE_UNSPECIFIED
	quorumId := uint32(0)
	retryAfterSeconds := uint32(0)
	context := ""

	err := NewErrorResourceExhaustedWithRateLimitDetails(msg, rateLimitType, quorumId, retryAfterSeconds, context)

	// Should still be a valid gRPC error
	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("Expected gRPC status error")
	}

	if st.Code() != codes.ResourceExhausted {
		t.Errorf("Expected RESOURCE_EXHAUSTED code, got %v", st.Code())
	}
}
