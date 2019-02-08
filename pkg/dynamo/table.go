package dynamo

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/tylfin/dynatomic/pkg/types"
)

// CreateTable creates dynamoDB table, this will be used for testing purposes
func CreateTable(svc *dynamodb.DynamoDB, schema *types.Schema) error {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: schema.HashKey,
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: schema.RangeKey,
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: schema.HashKey,
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: schema.RangeKey,
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: schema.TableName,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	_, err := svc.CreateTableWithContext(ctx, input)
	return err
}

// DeleteTable deletes the DynamoDB table
func DeleteTable(svc *dynamodb.DynamoDB, schema *types.Schema) error {
	input := &dynamodb.DeleteTableInput{
		TableName: schema.TableName,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	_, err := svc.DeleteTableWithContext(ctx, input)

	return err
}
