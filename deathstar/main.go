package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"roava-test/common"
	"roava-test/pb"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := 5001
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

	// Initialize gRPC TCP Listener
	tcpListener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	common.ExitOnError(err, "Failed to start listener . %v")

	// Initialize gRPC server
	server := grpc.NewServer()
	pb.RegisterDeathstarServer(server, &deathStarService{
		PubSubClient: pubSubClient,
		DbConn:       dbConn,
	})
	reflection.Register(server)

	// gRPC blocks main routine so I move it to it's on routine
	go func() {
		if e := server.Serve(tcpListener); e != nil {
			log.Fatalf("Failed to serve gRPC %v", e.Error())
		}
	}()

	// Prints info to CLI
	log.Printf("Deathstar is running on port %d\n", port)

	// Prevents main routine exit
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	<-sigint
	log.Println("Destroyer shutting down")
}
