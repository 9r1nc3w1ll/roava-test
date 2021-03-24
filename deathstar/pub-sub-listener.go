package main

import (
	"context"
	"log"
	"roava-test/common"
	"roava-test/pb"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/jackc/pgx/v4"
	"google.golang.org/protobuf/encoding/protojson"
)

// Watches topic for changes
func pubSubListener(consumer pulsar.Consumer, dbConn *pgx.Conn) {
	for {
		// Blocks until a message is received
		msg, err := consumer.Receive(context.Background())
		common.ExitOnError(err, "PubSub message error: %v")

		payload := pb.TargetsAcquiredPayload{}

		// Unmarshalls payload
		if err := protojson.Unmarshal(msg.Payload(), &payload); err == nil {
			// Store Payload in the database
			for _, t := range payload.Data {
				if _, err := dbConn.Exec(context.Background(), "INSERT INTO targets (id, message, created_on, updated_on) VALUES ($1, $2, $3, $4)", t.Id, t.Message, t.CreatedOn, t.UpdatedOn); err != nil {
					log.Printf("Failed to save target with error: %v\n", err.Error())
					continue //Skip the rest of the logic, do not acknowledge message
				}
			}

			log.Println("Targets saved to database.")
			consumer.Ack(msg)
		}
	}
}
