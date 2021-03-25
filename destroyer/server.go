package main

import (
	"context"
	"fmt"
	"log"

	"roava-test/common"
	pb "roava-test/pb"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jackc/pgx/v4"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
)

// destroyerService - implements all destroyerService methods
type destroyerService struct {
	DbConn       *pgx.Conn
	PubSubClient pulsar.Client
}

// AcquireTargets - implements AcquireTargets RPC
func (s *destroyerService) AcquireTargets(ctx context.Context, req *pb.AcquireTargetsRequest) (*empty.Empty, error) {
	targets := []*pb.Target{}

	for i := int64(0); i < req.Number; i++ {
		// Create reusable Target properties
		t := common.Timestamp()
		id := common.UUID()

		// Create and append Target to slice
		targets = append(targets, &pb.Target{
			Id:        id,
			Message:   fmt.Sprintf("Unique target %v", id),
			CreatedOn: t,
			UpdatedOn: t,
		})
	}

	// Create Pulsar producer
	producer, err := s.PubSubClient.CreateProducer(pulsar.ProducerOptions{
		Topic: "targets-acquired-event",
	})
	if err != nil {
		return &emptypb.Empty{}, err
	}
	defer producer.Close()

	// Build message payload
	payload := pb.TargetsAcquiredPayload{
		Id:        common.UUID(),
		Name:      "targets.acquired",
		Data:      targets,
		CreatedOn: common.Timestamp(),
	}

	jsonBytes, err := protojson.Marshal(&payload)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	// Send pub-sub message
	producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: jsonBytes,
	})

	log.Printf("%d target(s) acquired\n", req.GetNumber())
	return &empty.Empty{}, nil
}

// ListTargets - Implements ListTargets RPC
func (s *destroyerService) ListTargets(ctx context.Context, _ *empty.Empty) (*pb.ListTargetsResponse, error) {
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

func (s *destroyerService) HealthCheck(ctx context.Context, _ *empty.Empty) (*pb.HealthCheckResponse, error) {
	res := pb.HealthCheckResponse{
		Id:        common.UUID(),
		Service:   "destroyer",
		Timestamp: common.Timestamp(),
	}
	return &res, nil
}

func (s *destroyerService) ServiceReadiness(ctx context.Context, _ *empty.Empty) (*pb.ServiceReadinessResponse, error) {
	res := pb.ServiceReadinessResponse{
		Id:        common.UUID(),
		Service:   "destroyer",
		Timestamp: common.Timestamp(),
	}

	res.Dependencies = append(res.Dependencies, &pb.DependencyStatus{
		Name:   "postgres",
		Url:    s.DbConn.Config().ConnString(),
		Active: !s.DbConn.IsClosed(),
	})

	res.Dependencies = append(res.Dependencies, &pb.DependencyStatus{
		Name:   "apache pulsar",
		Url:    "pulsar://pulsar:6650",
		Active: s.PubSubClient != nil,
	})

	return &res, nil
}
