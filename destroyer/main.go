package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	pb "roava-test/pb"

	// "github.com/jackc/pgx/v4"
	"github.com/apache/pulsar-client-go/pulsar"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := 5000
	psUrl := "pulsar://pulsar:6650"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// _, e := pgx.Connect(ctx, "") //TODO: Setup PGSQL connection

	// if e != nil {
	// 	log.Fatalf("Database initialization error %v", e.Error())
	// }

	// Create pub-sub client, using pulsar connection
	pubSub, e := pulsar.NewClient(pulsar.ClientOptions{
		URL:              psUrl,
		OperationTimeout: 5 * time.Second,
	})

	if e != nil {
		log.Fatalf("Could not instantiate Pulsar client: %v", e)
	}

	defer pubSub.Close()

	listen, e := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if e != nil {
		log.Fatalf("Failed to listen . %v", e.Error())
	}

	s := grpc.NewServer()

	pb.RegisterDestroyerServer(s, &destroyer{
		Ctx:          ctx,
		PubSubClient: pubSub,
	})
	reflection.Register(s)

	// gRCP blocks main routine so I move it to it's on routine
	go func() {
		if e := s.Serve(listen); e != nil {
			log.Fatalf("Failed to serve gRPC %v", e.Error())
		}
	}()

	// Prints info to CLI
	log.Printf("Destroyer is running on port %d\n", port)

	/*
		Prevent destroyer from shutting down when main routine is done
		because gRPC is listening on another routine
	*/
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	<-sigint
	log.Println("Destroyer shutting down")
}
