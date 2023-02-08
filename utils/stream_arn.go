package utils

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

func GetLastStreamArn(db dynamodbiface.DynamoDBAPI, table string) (*string, error) {
	// Get the table information
	tableInfo, err := db.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(table),
	})
	if err != nil {
		return nil, err
	}

	// Get the stream ARN
	if tableInfo.Table.LatestStreamArn == nil {
		return nil, errors.New("no stream ARN found")
	}

	return tableInfo.Table.LatestStreamArn, nil
}
