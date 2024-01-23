package ratelimit

import (
	"context"
	"fmt"
	"strings"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/cenkalti/backoff/v4"
)

func isConcurrentUpdateError(err error) bool {
	if awsErr, ok := err.(awserr.Error); ok {
		fmt.Printf("Error code: %s\n", awsErr.Code())
		// Check if the error code is ValidationException, which is used for a variety of input errors
		if awsErr.Code() == "ValidationException" {
			// Check if the message contains the specific error detail
			return strings.Contains(awsErr.Message(), "Two document paths overlap with each other")
		}
	}
	return false
}

func updateItemWithRetry(ctx context.Context, requesterID string, initialUpdateBuilder *expression.UpdateBuilder, bs common.KVStoreVersioned[common.RateBucketParamsConcurrencySafe], retryPolicy *backoff.ExponentialBackOff) error {
	// Perform the initial update
	err := bs.UpdateItemWithExpression(ctx, requesterID, initialUpdateBuilder)
	if err != nil && isConcurrentUpdateError(err) {
		// Retry logic for concurrent update error
		return retryUpdateItemAfterConcurrentFailure(ctx, requesterID, initialUpdateBuilder, retryPolicy, bs)
	}
	// Return any non-concurrent update errors or nil if successful
	return err
}

func retryUpdateItemAfterConcurrentFailure(ctx context.Context, requesterID string, updateBuilder *expression.UpdateBuilder, retryPolicy *backoff.ExponentialBackOff, bs common.KVStoreVersioned[common.RateBucketParamsConcurrencySafe]) error {

	retryOperation := func() error {
		// Fetch the latest version when a concurrent update error occurs
		_, version, err := bs.GetItemWithVersion(ctx, requesterID)
		if err != nil {
			return backoff.Permanent(err) // Stop retrying on version fetch error
		}

		// Update the builder with the new version
		versionName := expression.Name("Version")
		updatedBuilder := updateBuilder.Set(versionName, expression.Value(version+1))

		// Retry the update with the new version
		return bs.UpdateItemWithExpression(ctx, requesterID, &updatedBuilder)
	}

	return backoff.Retry(retryOperation, retryPolicy)
}
