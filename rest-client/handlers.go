package main

import (
	"context"
	"fmt"
	"net/http"
	"roava-test/pb"
	"strconv"

	"github.com/go-chi/render"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Mux struct {
	DestroyerClient pb.DestroyerClient
	DeathstarClient pb.DeathstarClient
}

// REST request handlers
func (m *Mux) acquireTargets(rw http.ResponseWriter, r *http.Request) {
	number := int64(1) // Default number of targets
	targets := r.URL.Query().Get("targets")

	// Extract number of targets from request
	if targets != "" {
		n, err := strconv.ParseInt(targets, 10, 64)
		if err != nil {
			res := map[string]interface{}{"message": "Failed to parse targets parameter. An integer is required."}
			render.Status(r, http.StatusBadRequest)
			render.JSON(rw, r, res)
			return
		}

		if n > 1 {
			number = n
		}
	}

	req := &pb.AcquireTargetsRequest{
		Number: number,
	}

	_, err := m.DestroyerClient.AcquireTargets(context.Background(), req)
	if err != nil {
		internalError(rw, r, fmt.Sprintf("Request failed with error: %v", err.Error()))
		return
	}

	res := map[string]interface{}{"message": "Target acquired and published."}
	render.JSON(rw, r, res)
}

func (m *Mux) listTargets(rw http.ResponseWriter, r *http.Request) {
	targets, err := m.DestroyerClient.ListTargets(context.Background(), &empty.Empty{})
	if err != nil {
		internalError(rw, r, fmt.Sprintf("Request failed with error: %v", err.Error()))
		return
	}

	res := map[string]interface{}{
		"message": "Target acquired and published.",
		"targets": targets,
	}
	render.JSON(rw, r, res)
}

func (m *Mux) healthChecks(rw http.ResponseWriter, r *http.Request) {
	destoyerHealth, err := m.DestroyerClient.HealthCheck(context.Background(), &emptypb.Empty{})
	if err != nil {
		internalError(rw, r, fmt.Sprintf("Request failed with error: %v", err.Error()))
		return
	}

	deathstarHealth, err := m.DeathstarClient.HealthCheck(context.Background(), &emptypb.Empty{})
	if err != nil {
		internalError(rw, r, fmt.Sprintf("Request failed with error: %v", err.Error()))
		return
	}

	res := map[string]interface{}{
		"message":   "Health check results.",
		"destroyer": destoyerHealth,
		"deathstar": deathstarHealth,
	}
	render.JSON(rw, r, res)
}

func (m *Mux) serviceReadiness(rw http.ResponseWriter, r *http.Request) {
	destoyerReadiness, err := m.DestroyerClient.ServiceReadiness(context.Background(), &emptypb.Empty{})
	if err != nil {
		internalError(rw, r, fmt.Sprintf("Request failed with error: %v", err.Error()))
		return
	}

	deathstarReadiness, err := m.DeathstarClient.ServiceReadiness(context.Background(), &emptypb.Empty{})
	if err != nil {
		internalError(rw, r, fmt.Sprintf("Request failed with error: %v", err.Error()))
		return
	}

	res := map[string]interface{}{
		"message":   "Health check results.",
		"destroyer": destoyerReadiness,
		"deathstar": deathstarReadiness,
	}
	render.JSON(rw, r, res)
}
