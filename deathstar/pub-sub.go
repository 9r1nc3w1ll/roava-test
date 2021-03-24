package main

import (
	"context"
	"log"
	"roava-test/common"
	"roava-test/pb"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/jackc/pgx/v4"
	"google.golang.org/protobuf/encoding/protojson"
)

func pubSubConnect(psUrl string) pulsar.Consumer {
	// Create pub-sub client, using pulsar connection
	pubSubClient, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:              psUrl,
		OperationTimeout: 5 * time.Second,
	})
	common.ExitOnError(err, "Could not instantiate Pulsar client: %v")
	defer pubSubClient.Close()

	// Create pub-sub consumer
	consumer, err := pubSubClient.Subscribe(pulsar.ConsumerOptions{
		Topic:            "targets-acquired-event",
		SubscriptionName: "deathstar",
		Type:             pulsar.Shared,
	})
	common.ExitOnError(err, "Failed to connect to pulsar")

	return consumer
}

// Watches topic for changes
func pubSubListener(consumer pulsar.Consumer, dbConn *pgx.Conn) {
	for {
		// Blocks until a message is received
		msg, e := consumer.Receive(context.Background())
		if e != nil {
			log.Fatal(e)
		}

		payload := pb.TargetsAcquiredPayload{}

		// Unmarshalls payload
		if e := protojson.Unmarshal(msg.Payload(), &payload); e == nil {
			// Store Payload in the database
			for _, t := range payload.Data {
				if _, e := dbConn.Exec(context.Background(), "INSERT INTO targets (id, message, created_on, updated_on) VALUES ($1, $2, $3, $4)", t.Id, t.Message, t.CreatedOn, t.UpdatedOn); e != nil {
					log.Printf("Failed to save target with error: %v\n", e.Error())
					continue //Skip the rest of the logic, do not acknowledge message
				}
			}

			// Acknowledges message if there is no error
			consumer.Ack(msg) // Note if it is not acknowledged there will be an infinit loop
		}
	}
}
