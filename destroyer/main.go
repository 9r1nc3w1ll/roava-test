package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	pb "roava-test/pb"

	"github.com/apache/pulsar-client-go/pulsar"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := 5000
	pulsarUrl := "pulsar://pulsar:6650"

	// Create pub-sub client, using pulsar connection
	pubSubClient, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:              pulsarUrl,
		OperationTimeout: 5 * time.Second,
	})

	if err != nil {
		log.Fatalf("Could not instantiate Pulsar client: %v", err)
	}

	defer pubSubClient.Close()

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen . %v", err.Error())
	}

	server := grpc.NewServer()

	pb.RegisterDestroyerServer(server, &destroyer{
		PubSubClient: pubSubClient,
	})
	reflection.Register(server)

	// gRCP blocks main routine so I move it to it's on routine
	go func() {
		if e := server.Serve(listen); e != nil {
			log.Fatalf("Failed to serve gRPC %v", e.Error())
		}
	}()

	// Prints info to CLI
	log.Printf("Destroyer is running on port %d\n", port)

	// Prevents main routine exit
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	<-sigint
	log.Println("Destroyer shutting down")
}
