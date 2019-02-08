package dynamo_test

import (
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/tylfin/dynatomic/pkg/dynamo"
	"github.com/tylfin/dynatomic/pkg/types"
)

var (
	schema = &types.Schema{
		HashKey:   aws.String("Key"),
		RangeKey:  aws.String("MonthDay"),
		TableName: aws.String("test"),
		AtomicKey: aws.String("Incr"),
	}
)

func setup() (func(), *dynamodb.DynamoDB) {
	conf := &aws.Config{
		Region:      aws.String("us-east"),
		Endpoint:    aws.String("http://dynamodb:8000"),
		Credentials: credentials.NewStaticCredentials("fake", "fake", ""),
	}

	svc, err := dynamo.New(conf)
	if err != nil {
		log.Fatal(err)
	}

	err = dynamo.CreateTable(svc, schema)
	if err != nil {
		log.Fatal(err)
	}

	return func() {
		err = dynamo.DeleteTable(svc, schema)
		if err != nil {
			log.Fatal(err)
		}
	}, svc
}

func TestIncrement(t *testing.T) {
	teardown, svc := setup()
	defer teardown()

	row := &types.Row{
		HashValue:  aws.String("test"),
		RangeValue: aws.String("12-31"),
		Schema:     schema,
		Incr:       aws.String("1"),
	}

	// Try a single increment
	c, err := dynamo.Insert(svc, row)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), *c)

	// Try incrementing by a higher amount
	row.Incr = aws.String("10000")
	c, err = dynamo.Insert(svc, row)
	assert.Nil(t, err)
	assert.Equal(t, int64(10001), *c)

	// Foobar should not be used as a valid increment
	row.Incr = aws.String("foobar")
	c, err = dynamo.Insert(svc, row)
	assert.NotNil(t, err)
	if awsErr, ok := err.(awserr.Error); ok {
		assert.Equal(t, "ValidationException", awsErr.Code())
	}
}
