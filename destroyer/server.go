package main

import (
	"context"
	"log"

	pb "roava-test/pb"

	"github.com/golang/protobuf/ptypes/empty"
)

// destroyer - implements all destroyer methods
type destroyer struct{}

// AcquireTargets - implements AcquireTargets
func (s *destroyer) AcquireTargets(ctx context.Context, req *pb.AcquireTargetRequest) (*empty.Empty, error) {
	log.Printf("%d target(s) acquired\n", req.GetNumber())
	return nil, nil
}

func (s *destroyer) ListTargets(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	return nil, nil
}
