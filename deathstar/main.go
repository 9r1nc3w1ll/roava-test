package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"roava-test/common"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/jackc/pgx/v4"
)

func main() {
	psUrl := "pulsar://pulsar:6650"
	dbUrl := "host=db port=5432 user=postgres dbname=roava_test"

	// Connect to database
	dbConn, err := pgx.Connect(context.Background(), dbUrl)
	common.ExitOnError(err, "Database initialization error %v")
	defer dbConn.Close(context.Background())

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
	defer consumer.Close()
	log.Println("deathstar is online.")

	// Pub-Sub listener routine
	go pubSubListener(consumer, dbConn)

	// Prevents main routine exit
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	<-sigint
	log.Println("Destroyer shutting down")
}
