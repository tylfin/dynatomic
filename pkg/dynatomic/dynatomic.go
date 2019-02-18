// Package dynatomic provides a convenient wrapper API around
// using DynamoDB as highly available, concurrent, and performant
// asynchronous atomic counter
//
// Basic usage:
//  // Initialize the dynatomic backround goroutine with a batch size of 100,
//  // a wait time of a second, an AWS config and a function that will
//  // notify the user of internal errors
//  d := New(100, time.Second, config, errHandler)
//  d.RowChan <- &types.Row{...}
//  d.RowChan <- &types.Row{...}
//  d.RowChan <- &types.Row{...}
//  ...
//  d.Done()
// Dynamo will update accordingly
// For example if you write the rows:
// 	Table: MyTable, Key: A, Range: A, Incr: 5
// 	Table: MyTable, Key: A, Range: A, Incr: 5
// 	Table: MyTable, Key: A, Range: A, Incr: 5
// 	Table: MyTable, Key: A, Range: A, Incr: 5
// Then MyTable Key A, Range A will now show a value of 20
package dynatomic

import (
	"strconv"
	"time"

	"github.com/tylfin/dynatomic/pkg/dynamo"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/tylfin/dynatomic/pkg/types"
)

var (
	new    = dynamo.New
	insert = dynamo.Insert
)

// Dynatomic struct contains all information
// necessary to perform bulk, asynchronous writes
type Dynatomic struct {
	done chan bool
	svc  *dynamodb.DynamoDB

	RowChan    chan *types.Row
	Config     *aws.Config
	ErrHandler func(loc string, err error)

	BatchSize int
	WaitTime  time.Duration
}

// New creates a dynatomic struct listening to the RowChan for writes
func New(batchSize int, waitTime time.Duration, config *aws.Config, errHandler func(loc string, err error)) *Dynatomic {
	dynatomic := &Dynatomic{}

	dynatomic.done = make(chan bool)
	dynatomic.RowChan = make(chan *types.Row)
	dynatomic.Config = config
	dynatomic.BatchSize = batchSize
	dynatomic.WaitTime = waitTime
	dynatomic.ErrHandler = errHandler

	go dynatomic.run()
	return dynatomic
}

// Done takes the last of the messages ready to be sent
// and destroys the running goroutine
func (d *Dynatomic) Done() {
	select {
	case <-d.done:
		return
	default:
	}
	close(d.done)
}

func (d *Dynatomic) run() {
	var (
		err      error
		finished bool
	)

	d.svc, err = new(d.Config)
	if err != nil {
		d.ErrHandler("run.dynamo.New", err)
		d.Done()
		return
	}

	for !finished {
		finished = d.batch()
	}
}

func (d *Dynatomic) batch() bool {
	group := map[string][]*types.Row{}

	for i := 0; i < d.BatchSize; i++ {
		select {
		case <-d.done:
			d.write(group)
			return true
		case row := <-d.RowChan:
			group[*row.Schema.TableName] = append(group[*row.Schema.TableName], row)
		case <-time.Tick(d.WaitTime):
			continue
		}
	}

	d.write(group)
	return false
}

func (d *Dynatomic) write(group map[string][]*types.Row) {
	for _, rows := range group {
		bulkRow := &types.Row{Schema: rows[0].Schema, HashValue: rows[0].HashValue, RangeValue: rows[0].RangeValue}
		incr := 0
		for _, row := range rows {
			v, err := strconv.Atoi(*row.Incr)
			if err != nil {
				d.ErrHandler("write.atoi", err)
				continue
			}
			incr += v
		}
		incrStr := strconv.Itoa(incr)
		bulkRow.Incr = &incrStr
		_, err := insert(d.svc, bulkRow)
		if err != nil {
			d.ErrHandler("write.insert", err)
		}
	}
}
