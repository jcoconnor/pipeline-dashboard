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
