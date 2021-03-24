package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"roava-test/pb"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"google.golang.org/grpc"
)

func main() {
	// Variables
	port := ":3000"
	destroyerAddress := "0.0.0.0:5000"

	// Create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect to Destroyer
	conn, e := grpc.Dial(destroyerAddress, grpc.WithInsecure())
	if e != nil {
		log.Fatalf("Failed to connect with destroyer. %v", e.Error())
	}

	destroyer := pb.NewDestroyerClient(conn)

	// Setup HTTP endpoints
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Acquire targets endpoint
	r.Get("/acquire-targets", func(rw http.ResponseWriter, r *http.Request) {
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

		req := &pb.AcquireTargetRequest{
			Number: number,
		}

		_, e := destroyer.AcquireTargets(ctx, req)
		if e != nil {
			internalError(rw, r, fmt.Sprintf("Request failed with error: %v", e.Error()))
			return
		}

		res := map[string]interface{}{"message": "Target acquired and published."}
		render.JSON(rw, r, res)
	})

	// Start REST server
	srv := http.Server{
		Addr:    port,
		Handler: r,
	}

	if e := srv.ListenAndServe(); e != http.ErrServerClosed {
		log.Fatalf("HTTP server listen and serve error: %v", e)
	}
}

// Handle internal server errors to avoid repeatition
func internalError(rw http.ResponseWriter, r *http.Request, message string) {
	e := map[string]interface{}{"errors": message}
	render.Status(r, http.StatusInternalServerError)
	render.JSON(rw, r, e)
}
