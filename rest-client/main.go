package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"roava-test/common"
	"roava-test/pb"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"google.golang.org/grpc"
)

var destroyerClient pb.DestroyerClient

func main() {
	port := ":3000"
	destroyerServerAddress := "0.0.0.0:5000"

	// Connect to Destroyer
	conn, err := grpc.Dial(destroyerServerAddress, grpc.WithInsecure())
	common.ExitOnError(err, "Failed to connect with destroyer. %v")
	defer conn.Close()

	// Destroyer gRCP Client
	destroyerClient = pb.NewDestroyerClient(conn)

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

	// Serve REST endpoints
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

	// Prevents main routine exit
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
