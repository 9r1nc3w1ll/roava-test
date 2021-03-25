package main

import (
	"context"
	"fmt"
	"net/http"
	"roava-test/pb"
	"strconv"

	"github.com/go-chi/render"
	"github.com/golang/protobuf/ptypes/empty"
)

// REST request handlers
func acquireTargets(rw http.ResponseWriter, r *http.Request) {
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

	_, err := destroyerClient.AcquireTargets(context.Background(), req)
	if err != nil {
		internalError(rw, r, fmt.Sprintf("Request failed with error: %v", err.Error()))
		return
	}

	res := map[string]interface{}{"message": "Target acquired and published."}
	render.JSON(rw, r, res)
}

func listTargets(rw http.ResponseWriter, r *http.Request) {
	targets, err := destroyerClient.ListTargets(context.Background(), &empty.Empty{})
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
