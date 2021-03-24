package common

import "log"

func ExitOnError(err error, message string) {
	if err != nil {
		log.Fatalf(message, err)
	}
}
