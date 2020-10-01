package web

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

func lastUpdated() float64 {
    content, err := ioutil.ReadFile("updated")

    // Get Last update from stat of file.
    // See https://gist.github.com/alexisrobert/982674

    if err != nil {
        log.Fatal(err)
    }

    intLastUpdated, err := strconv.ParseFloat(strings.TrimSpace(string(content)), 32)

    if err != nil {
        fmt.Println(err)
    }

    secondsSinceLastUpdate := float64(time.Now().Unix()) - intLastUpdated
    return secondsSinceLastUpdate

}
