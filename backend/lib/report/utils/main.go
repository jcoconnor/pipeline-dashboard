package utils

import (
	"log"
)

func CheckError(message string, err error) {
    if err != nil {
        log.Fatal(message, err)
    }
}

func StringSliceContains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}
