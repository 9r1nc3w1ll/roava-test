package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"roava-test/pb"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestAcquireTargets(t *testing.T) {
	mockMux := Mux{
		DestroyerClient: &mockDestoyerClient{},
		DeathstarClient: &mockDeathstarClient{},
	}

	req, err := http.NewRequest("POST", "/acquire-targets", bytes.NewReader(nil))
	if err != nil {
		t.Fatalf("Request creation failed. %s", err)
	}

	rr := httptest.NewRecorder()
	h := http.HandlerFunc(mockMux.acquireTargets)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 status code but got %d", req.Response.StatusCode)
	}
}

func TestListTargets(t *testing.T) {
	mockMux := Mux{
		DestroyerClient: &mockDestoyerClient{},
		DeathstarClient: &mockDeathstarClient{},
	}

	req, err := http.NewRequest("POST", "/list-targets", bytes.NewReader(nil))
	if err != nil {
		t.Fatalf("Request creation failed. %s", err)
	}

	rr := httptest.NewRecorder()
	h := http.HandlerFunc(mockMux.listTargets)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 status code but got %d", req.Response.StatusCode)
	}
}

/*
	Find Mocks below this line
*/

type mockDestoyerClient struct{}

func (d *mockDestoyerClient) AcquireTargets(ctx context.Context, in *pb.AcquireTargetsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return &empty.Empty{}, nil
}

func (d *mockDestoyerClient) ListTargets(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pb.ListTargetsResponse, error) {
	return &pb.ListTargetsResponse{}, nil
}

func (d *mockDestoyerClient) HealthCheck(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{}, nil
}

func (d *mockDestoyerClient) ServiceReadiness(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pb.ServiceReadinessResponse, error) {
	return &pb.ServiceReadinessResponse{}, nil
}

type mockDeathstarClient struct{}

func (mockDeathstarClient) HealthCheck(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{}, nil
}

func (mockDeathstarClient) ServiceReadiness(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pb.ServiceReadinessResponse, error) {
	return &pb.ServiceReadinessResponse{}, nil
}
