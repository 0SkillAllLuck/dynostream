package main

import (
	"context"
	"log"

	"github.com/0skillallluck/dynostream"
	"github.com/aws/aws-sdk-go/aws/session"
)

func main() {
	// Create a new AWS session
	sess := session.New()

	// Create a new DynoStream
	stream := dynostream.NewFromSession(sess, "my-table")

	// Create a new context for the subscription
	ctx := context.Background()

	// Subscribe to the stream
	channel, errorChannel := stream.Subscribe(ctx)

	// Handle errors
	go func(errCh <-chan error) {
		for err := range errCh {
			log.Println("Error: ", err)
		}
	}(errorChannel)

	// Handle records
	for record := range channel {
		log.Println("Record:", record)
	}
}
