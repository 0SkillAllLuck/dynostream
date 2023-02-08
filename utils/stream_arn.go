package utils

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	// ErrNoStreamArn is returned when no stream ARN is found
	ErrNoStreamArn = errors.New("no stream ARN found")
)

func GetLastStreamArn(ctx context.Context, db dynamodbiface.DynamoDBAPI, table string) (*string, error) {
	// Get the table information
	tableInfo, err := db.DescribeTableWithContext(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(table),
	})
	if err != nil {
		return nil, err
	}

	// Get the stream ARN
	if tableInfo.Table.LatestStreamArn == nil {
		return nil, ErrNoStreamArn
	}

	return tableInfo.Table.LatestStreamArn, nil
}
