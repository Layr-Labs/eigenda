package ondemand_test

import (
	"testing"

	"github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/require"
)

func TestNewOnDemandPaymentValidator(t *testing.T) {
	mockOnChainState := &mock.MockOnchainPaymentState{}
	dynamoClient := &dynamodb.Client{}
	tableName := "test-table"
	maxLedgers := 100

	t.Run("nil onChainState", func(t *testing.T) {
		validator, err := ondemand.NewOnDemandPaymentValidator(
			testutils.GetLogger(),
			maxLedgers,
			nil,
			dynamoClient,
			tableName,
		)
		require.Error(t, err)
		require.Nil(t, validator)
	})

	t.Run("nil dynamoClient", func(t *testing.T) {
		validator, err := ondemand.NewOnDemandPaymentValidator(
			testutils.GetLogger(),
			maxLedgers,
			mockOnChainState,
			nil,
			tableName,
		)
		require.Error(t, err)
		require.Nil(t, validator)
	})

	t.Run("empty table name", func(t *testing.T) {
		validator, err := ondemand.NewOnDemandPaymentValidator(
			testutils.GetLogger(),
			maxLedgers,
			mockOnChainState,
			dynamoClient,
			"",
		)
		require.Error(t, err)
		require.Nil(t, validator)
	})

	t.Run("zero max ledgers", func(t *testing.T) {
		validator, err := ondemand.NewOnDemandPaymentValidator(
			testutils.GetLogger(),
			0,
			mockOnChainState,
			dynamoClient,
			tableName,
		)
		require.Error(t, err)
		require.Nil(t, validator)
	})

	t.Run("negative max ledgers", func(t *testing.T) {
		validator, err := ondemand.NewOnDemandPaymentValidator(
			testutils.GetLogger(),
			-1,
			mockOnChainState,
			dynamoClient,
			tableName,
		)
		require.Error(t, err)
		require.Nil(t, validator)
	})
}
