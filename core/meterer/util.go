package meterer

import (
	"context"

	commonaws "github.com/Layr-Labs/eigenda/common/aws"
	test_utils "github.com/Layr-Labs/eigenda/common/aws/dynamodb/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func CreateReservationTable(clientConfig commonaws.ClientConfig, tableName string) error {
	ctx := context.Background()
	_, err := test_utils.CreateTable(ctx, clientConfig, tableName, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("AccountID"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("ReservationPeriod"),
				AttributeType: types.ScalarAttributeTypeN,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("AccountID"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("ReservationPeriod"),
				KeyType:       types.KeyTypeRange,
			},
		},
		TableName: aws.String(tableName),
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	})
	return err
}

func CreateGlobalReservationTable(clientConfig commonaws.ClientConfig, tableName string) error {
	ctx := context.Background()
	_, err := test_utils.CreateTable(ctx, clientConfig, tableName, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("ReservationPeriod"),
				AttributeType: types.ScalarAttributeTypeN,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("ReservationPeriod"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName: aws.String(tableName),
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	})

	if err != nil {
		if err.Error() == "ResourceInUseException: Table already exists" {
			return nil
		}
		return err
	}
	return nil
}

func CreateOnDemandTable(clientConfig commonaws.ClientConfig, tableName string) error {
	ctx := context.Background()
	_, err := test_utils.CreateTable(ctx, clientConfig, tableName, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("AccountID"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("AccountID"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName: aws.String(tableName),
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	})
	if err != nil {
		if err.Error() == "ResourceInUseException: Table already exists" {
			return nil
		}
		return err
	}
	return nil
}
