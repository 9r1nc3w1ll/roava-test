package main

import (
	"context"
	"log"
	"net"

	pb "roava-test/pb"

	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := ":3000"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, e := pgx.Connect(ctx, "") //TODO: Setup PGSQL connection

	if e != nil {
		log.Fatalf("Database initialization error %v", e.Error())
	}

	listen, e := net.Listen("tcp", port)
	if e != nil {
		log.Fatalf("Failed to listen . %v", e.Error())
	}

	s := grpc.NewServer()

	pb.RegisterDestroyerServer(s, &destroyer{})
	reflection.Register(s)

	if e := s.Serve(listen); e != nil {
		log.Fatalf("Failed to serve gRPC %v", e.Error())
	}
}
