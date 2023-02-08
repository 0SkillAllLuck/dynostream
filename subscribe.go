package dynostream

import (
	"context"
	"time"

	"github.com/0skillallluck/dynostream/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodbstreams"
)

// Subscribe creates a new subscription to the DynamoDB stream for the table specified in the
// DynoStream struct. It returns a channel of dynamodbstreams.Record structs and a channel of
// errors. The dynamodbstreams.Record channel will be closed when the context is closed.
func (ds *DynoStream) Subscribe(ctx context.Context) (<-chan *dynamodbstreams.Record, <-chan error) {
	// Create the channels
	channel := make(chan *dynamodbstreams.Record, 1)
	errorChannel := make(chan error, 1)

	// Create a goroutine to handle the subscription
	go func(channel chan<- *dynamodbstreams.Record, errorChannel chan<- error) {
		// Close the channels when the context is closed
		defer close(channel)
		defer close(errorChannel)

		var streamArn *string
		var shardId *string
		var previousShardId *string
		var err error

		for {
			select {
			case <-ctx.Done():
				return
			default:
				previousShardId = shardId
				// Get Stream ARN
				streamArn, err = utils.GetLastStreamArn(ds.db, ds.table)
				if err != nil {
					errorChannel <- err
					continue
				}

				// Get Shard ID
				shardId, err = utils.GetShardId(ds.dbStreams, previousShardId, streamArn)
				if err != nil {
					errorChannel <- err
					continue
				}

				// Get Records
				if shardId == nil {
					time.Sleep(time.Second * 10)
					continue
				}

				// TODO: Handle shard splitting
				err = ds.processShard(&dynamodbstreams.GetShardIteratorInput{
					StreamArn:         streamArn,
					ShardId:           shardId,
					ShardIteratorType: aws.String(dynamodbstreams.ShardIteratorTypeLatest),
				}, channel)
				if err != nil {
					errorChannel <- err
					// Reset shardID to previousShardId to try and process it again
					shardId = previousShardId
				}
			}
		}
	}(channel, errorChannel)

	return channel, errorChannel
}
