package main

import (
	"context"
	"roava-test/common"
	"roava-test/pb"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jackc/pgx/v4"
)

type deathStarService struct {
	DbConn       *pgx.Conn
	PubSubClient pulsar.Client
}

func (s *deathStarService) HealthCheck(ctx context.Context, _ *empty.Empty) (*pb.HealthCheckResponse, error) {
	res := pb.HealthCheckResponse{
		Id:        common.UUID(),
		Service:   "deathstar",
		Timestamp: common.Timestamp(),
	}
	return &res, nil
}

func (s *deathStarService) ServiceReadiness(ctx context.Context, _ *empty.Empty) (*pb.ServiceReadinessResponse, error) {
	res := pb.ServiceReadinessResponse{
		Id:        common.UUID(),
		Service:   "deathstar",
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
