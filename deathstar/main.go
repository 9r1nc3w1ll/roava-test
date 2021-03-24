package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"roava-test/pb"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"google.golang.org/grpc"
)

// Declare package wide globals
var destroyer pb.DestroyerClient
var ctx context.Context
var cancelCtx context.CancelFunc

func main() {
	// Variables
	port := ":3000"
	destroyerAddress := "0.0.0.0:5000"

	// Create context
	ctx, cancelCtx = context.WithCancel(context.Background())
	defer cancelCtx()

	// Connect to Destroyer
	conn, e := grpc.Dial(destroyerAddress, grpc.WithInsecure())
	if e != nil {
		log.Fatalf("Failed to connect with destroyer. %v", e.Error())
	}

	destroyer = pb.NewDestroyerClient(conn)

	// Setup HTTP endpoints
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Acquire targets endpoint
	r.Get("/acquire-targets", acquireTargets)

	// Start REST server
	srv := http.Server{
		Addr:    port,
		Handler: r,
	}

	// This blocks the main routine so, let it run separately.
	go func() {
		if e := srv.ListenAndServe(); e != http.ErrServerClosed {
			log.Fatalf("HTTP server listen and serve error: %v", e)
		}
	}()

	// Print out routes
	if err := chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("➪ %s %s\n", method, route)
		return nil
	}); err != nil {
		log.Panicf("Chi walk err: %s\n", err.Error())
	}

	/*
		Prevent destroyer from shutting down when main routine is done
		because gRPC is listening on another routine
	*/
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	<-sigint
	log.Println("Destroyer shutting down")
}

// Handle internal server errors to avoid repeatition
func internalError(rw http.ResponseWriter, r *http.Request, message string) {
	e := map[string]interface{}{"errors": message}
	render.Status(r, http.StatusInternalServerError)
	render.JSON(rw, r, e)
}
