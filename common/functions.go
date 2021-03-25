package common

import (
	"log"
	"time"

	"github.com/google/uuid"
)

func ExitOnError(err error, message string) {
	if err != nil {
		log.Fatalf(message, err)
	}
}

func UUID() string {
	return uuid.New().String()
}

func Timestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}
