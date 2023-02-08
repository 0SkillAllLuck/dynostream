package utils

import (
	"context"

	"github.com/aws/aws-sdk-go/service/dynamodbstreams"
	"github.com/aws/aws-sdk-go/service/dynamodbstreams/dynamodbstreamsiface"
)

func GetShardId(ctx context.Context, dbStreams dynamodbstreamsiface.DynamoDBStreamsAPI, previousShardId *string, streamArn *string) (*string, error) {
	streamInfo, err := dbStreams.DescribeStreamWithContext(ctx, &dynamodbstreams.DescribeStreamInput{
		StreamArn: streamArn,
	})
	if err != nil {
		return nil, err
	}
	if len(streamInfo.StreamDescription.Shards) == 0 {
		return nil, nil
	}

	// If no previous shard id, return the first shard
	if previousShardId == nil {
		return streamInfo.StreamDescription.Shards[0].ShardId, nil
	}

	var shardId *string
	for _, shard := range streamInfo.StreamDescription.Shards {
		shardId = shard.ShardId
		if shard.ParentShardId != nil && shard.ParentShardId == previousShardId {
			break
		}
	}
	return shardId, nil
}
