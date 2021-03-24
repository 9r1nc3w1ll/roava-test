package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"roava-test/common"

	"github.com/jackc/pgx/v4"
)

func main() {
	psUrl := "pulsar://pulsar:6650"
	dbUrl := "host=db port=5432 user=postgres dbname=roava_test"

	// Connect to database
	dbConn, err := pgx.Connect(context.Background(), dbUrl)
	common.ExitOnError(err, "Database initialization error %v")
	defer dbConn.Close(context.Background())

	consumer := pubSubConnect(psUrl)
	defer consumer.Close()
	go pubSubListener(consumer, dbConn)

	// Prevents main routine exit
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	<-sigint
	log.Println("Destroyer shutting down")
}
