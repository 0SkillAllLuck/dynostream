package dynostream

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodbstreams"
)

func (ds *DynoStream) processShard(input *dynamodbstreams.GetShardIteratorInput, channel chan<- *dynamodbstreams.Record) error {
	// Get the shard iterator
	shardIterator, err := ds.dbStreams.GetShardIterator(input)
	if err != nil {
		return err
	}
	// Check if results are emptry
	if shardIterator.ShardIterator == nil {
		return nil
	}
	// Process the shard records
	iterator := shardIterator.ShardIterator
	for iterator != nil {
		records, err := ds.dbStreams.GetRecords(&dynamodbstreams.GetRecordsInput{
			ShardIterator: iterator,
		})

		// In DynamoDB Streams, there is a 24 hour limit on data retention. Stream records whose age exceeds this limit are subject to removal (trimming) from the stream.
		// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_streams_GetRecords.html
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "TrimmedDataAccessException" {
			return nil
		}

		if err != nil {
			return err
		}

		// Publish records to the channel
		for _, record := range records.Records {
			channel <- record
		}
		iterator = records.NextShardIterator

		// Calculate sleep time
		sleepTime := time.Second
		if iterator == nil {
			// Shard closed, get new shard
			sleepTime = time.Millisecond * 10
		} else if len(records.Records) == 0 {
			// No records, sleep for 10 seconds
			sleepTime = time.Second * 10
		}

		time.Sleep(sleepTime)
	}
	return nil
}
