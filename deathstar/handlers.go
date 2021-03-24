package main

import (
	"fmt"
	"net/http"
	"roava-test/pb"
	"strconv"

	"github.com/go-chi/render"
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

	_, e := destroyer.AcquireTargets(ctx, req)
	if e != nil {
		internalError(rw, r, fmt.Sprintf("Request failed with error: %v", e.Error()))
		return
	}

	res := map[string]interface{}{"message": "Target acquired and published."}
	render.JSON(rw, r, res)
}
