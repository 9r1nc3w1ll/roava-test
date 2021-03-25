package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "roava-test/pb"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
)

// destroyer - implements all destroyer methods
type destroyer struct {
	DbConn       *pgx.Conn
	PubSubClient pulsar.Client
}

// AcquireTargets - implements AcquireTargets
func (s *destroyer) AcquireTargets(ctx context.Context, req *pb.AcquireTargetsRequest) (*empty.Empty, error) {
	targets := []*pb.Target{}

	for i := int64(0); i < req.Number; i++ {
		// Create reusable Target properties
		t := time.Now().UTC().Format(time.RFC3339)
		id := uuid.New().String()

		// Create and append Target to slice
		targets = append(targets, &pb.Target{
			Id:        id,
			Message:   fmt.Sprintf("Unique target %v", id),
			CreatedOn: t,
			UpdatedOn: t,
		})

		// Create Pulsar producer
		producer, err := s.PubSubClient.CreateProducer(pulsar.ProducerOptions{
			Topic: "targets-acquired-event",
		})
		if err != nil {
			return &emptypb.Empty{}, err
		}

		// Build message payload
		payload := pb.TargetsAcquiredPayload{
			Id:        uuid.New().String(),
			Name:      "targets.acquired",
			Data:      targets,
			CreatedOn: time.Now().UTC().Format(time.RFC3339),
		}

		jsonBytes, err := protojson.Marshal(&payload)
		if err != nil {
			return &emptypb.Empty{}, err
		}

		// Send pub-sub message
		producer.Send(context.Background(), &pulsar.ProducerMessage{
			Payload: jsonBytes,
		})

		// Close producer
		defer producer.Close()
	}

	log.Printf("%d target(s) acquired\n", req.GetNumber())
	return &empty.Empty{}, nil
}

func (s *destroyer) ListTargets(ctx context.Context, _ *empty.Empty) (*pb.ListTargetsResponse, error) {
	res := pb.ListTargetsResponse{}
	rows, err := s.DbConn.Query(context.Background(), "SELECT id, message, created_on, updated_on FROM targets")
	if err != nil {
		return &res, err
	}

	for rows.Next() {
		t := pb.Target{}

		if err := rows.Scan(&t.Id, &t.Message, &t.CreatedOn, &t.UpdatedOn); err != nil {
			return &res, err
		}
		res.Data = append(res.Data, &t)
	}

	return &res, nil
}
