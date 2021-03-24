package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "roava-test/pb"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
)

// destroyer - implements all destroyer methods
type destroyer struct{}

// AcquireTargets - implements AcquireTargets
func (s *destroyer) AcquireTargets(ctx context.Context, req *pb.AcquireTargetRequest) (*empty.Empty, error) {
	targets := []pb.Target{}

	for i := int64(0); i < req.GetNumber(); i++ {
		id := uuid.New()
		t := time.Now().UTC().Format(time.RFC3339)

		targets = append(targets, pb.Target{
			Id:        id.String(),
			Message:   fmt.Sprintf("Unique target %v", id),
			CreatedOn: t,
			UpdatedOn: t,
		})
	}

	log.Printf("%d target(s) acquired\n", req.GetNumber())
	return &empty.Empty{}, nil
}

func (s *destroyer) ListTargets(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}
