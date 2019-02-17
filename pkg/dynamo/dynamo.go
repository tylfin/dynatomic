// Package dynamo wraps the AWS DynamoDB SDK
//
// To insert a row:
//  svc, _ := dynamo.New(conf)
//  row := &types.Row{...}
//  count, err := dynamo.Insert(svc, row)
package dynamo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/tylfin/dynatomic/pkg/types"
)

var (
	db *dynamodb.DynamoDB
)

// New creates a dynamodb client on first call,
// and returns the existing client on subsequent calls.
func New(config *aws.Config) (*dynamodb.DynamoDB, error) {
	if db == nil {
		sess, err := session.NewSession(config)

		if err != nil {
			return nil, err
		}

		// Create DynamoDB client
		db = dynamodb.New(sess)
	}

	return db, nil
}

// Insert takes in a row, and increments the atomic key value by row.Incr using the query:
// 	SET [ATOMIC KEY] = IF_NOT_EXISTS(VALUE, 0) + VALUE
// This will row lock and can create hotspots if used on the same key too often,
// for high volume writes on a single key use the dynatomic async wrapper for an
// eventually consistent result
func Insert(svc *dynamodb.DynamoDB, row *types.Row) (*int64, error) {
	// Avoid an off-by-one error by using the if_not_exists to create a new row if it doesn't exist
	updateExpr := fmt.Sprintf("SET %s = if_not_exists(%s, :zero) + :incr", *row.Schema.AtomicKey, *row.Schema.AtomicKey)

	input := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			*row.Schema.HashKey:  &dynamodb.AttributeValue{S: row.HashValue},
			*row.Schema.RangeKey: &dynamodb.AttributeValue{S: row.RangeValue},
		},
		UpdateExpression: aws.String(updateExpr),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":incr": &dynamodb.AttributeValue{N: row.Incr},
			":zero": &dynamodb.AttributeValue{N: aws.String("0")},
		},
		TableName:    row.Schema.TableName,
		ReturnValues: aws.String("UPDATED_NEW"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	res, err := svc.UpdateItemWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	countStr, ok := res.Attributes[*row.Schema.AtomicKey]
	if !ok {
		return nil, types.ErrAtomicAttribute
	}

	count, err := strconv.ParseInt(*countStr.N, 10, 64)
	if err != nil {
		return nil, err
	}

	return &count, nil
}
