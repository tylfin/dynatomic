package dynatomic

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/tylfin/dynatomic/pkg/types"
)

func errHandler(loc string, err error) {
	fmt.Println(loc, err)
}

func setup() func() {
	origNew := new
	origInsert := insert

	new = func(*aws.Config) (*dynamodb.DynamoDB, error) {
		return nil, nil
	}

	insert = func(svc *dynamodb.DynamoDB, row *types.Row) (*int64, error) {
		return nil, nil
	}

	return func() {
		new = origNew
		insert = origInsert
	}
}

func TestNew(t *testing.T) {
	teardown := setup()
	defer teardown()

	dynatomic := New(300, time.Second, nil, errHandler)
	assert.NotNil(t, dynatomic.RowChan)
	assert.NotNil(t, dynatomic.done)
	time.Sleep(time.Second * 2)
	dynatomic.Done()

	// Ensure the done channel is closed
	v, ok := <-dynatomic.done
	assert.Zero(t, v)
	assert.False(t, ok)
}

func TestBrokenConfig(t *testing.T) {
	teardown := setup()
	defer teardown()

	new = func(*aws.Config) (*dynamodb.DynamoDB, error) {
		return nil, errors.New("fake error")
	}

	dynatomic := New(300, time.Second, nil, errHandler)
	v, ok := <-dynatomic.done
	assert.Zero(t, v)
	assert.False(t, ok)
}

func TestBulkInsert(t *testing.T) {
	teardown := setup()
	defer teardown()

	var lastRow *types.Row
	insert = func(svc *dynamodb.DynamoDB, row *types.Row) (*int64, error) {
		lastRow = row
		return nil, nil
	}

	dynatomic := New(10, 10*time.Millisecond, nil, errHandler)
	dynatomic.RowChan <- &types.Row{Schema: &types.Schema{TableName: aws.String("fakeTable")}, Incr: aws.String("5")}

	time.Sleep(300 * time.Millisecond)
	assert.Equal(t, *lastRow.Incr, "5")

	for i := 0; i < 2; i++ {
		dynatomic.RowChan <- &types.Row{Schema: &types.Schema{TableName: aws.String("fakeTable")}, Incr: aws.String("5")}
	}

	time.Sleep(300 * time.Millisecond)
	assert.Equal(t, *lastRow.Incr, "10")
}
