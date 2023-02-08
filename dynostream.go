package dynostream

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodbstreams"
	"github.com/aws/aws-sdk-go/service/dynamodbstreams/dynamodbstreamsiface"
)

// DynoStream is a struct that contains the DynamoDB and DynamoDBStreams interfaces
// as well as the name of the table to stream from.
type DynoStream struct {
	db        dynamodbiface.DynamoDBAPI
	dbStreams dynamodbstreamsiface.DynamoDBStreamsAPI
	table     string
}

// New creates a new DynoStream struct from existing DynamoDB and DynamoDBStreams interfaces
// and the name of the table to stream from.
func New(db dynamodbiface.DynamoDBAPI, dbStreams dynamodbstreamsiface.DynamoDBStreamsAPI, table string) *DynoStream {
	return &DynoStream{
		db:        db,
		dbStreams: dbStreams,
		table:     table,
	}
}

// NewFromSession creates a new DynoStream struct from an existing AWS session and the name
// of the table to stream from.
func NewFromSession(sess *session.Session, table string) *DynoStream {
	return New(dynamodb.New(sess), dynamodbstreams.New(sess), table)
}
