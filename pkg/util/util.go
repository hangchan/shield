package util

import "log"

func LogError(err interface{}) {
	if err != nil {
		log.Fatal(err)
	}
}
