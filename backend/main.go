package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"

	"github.com/puppetlabs/pipeline-dashboard/backend/config"
	"github.com/puppetlabs/pipeline-dashboard/backend/lib/report"
	"github.com/puppetlabs/pipeline-dashboard/backend/lib/report/cith"
	"github.com/puppetlabs/pipeline-dashboard/backend/lib/report/jenkins_types"
	"github.com/puppetlabs/pipeline-dashboard/backend/lib/report/utils"
	"github.com/puppetlabs/pipeline-dashboard/backend/lib/web"
)

// CithFailures Exported function
func CithFailures(config config.Config) []cith.CithFailure {
    var cithFailures []cith.CithFailure

    if len(config.CithURL) == 0 {
        return cithFailures
    }

    g := utils.Getable{
        URL:    "CITHFAILURES",
        Config: config,
    }

    client := g.GetRedisClient()

    var cached bool
    if client != nil {
        defer client.Close()

        var body []byte
        cached, body = g.Cached(client, "CITHFAILURES")
        json.Unmarshal(body, &cithFailures)
    } else {
        cached = false
    }

    if cached {
        return cithFailures
    }

    cithFailures = report.CompileCith(config.CithURL)

    if len(cithFailures) == 0 {
        return cithFailures
    }

    body, err := json.Marshal(cithFailures)

    if err != nil {
        panic(err)
    }


    g.Cache(client, "CITHFAILURES", body)

    return cithFailures
}

// JenkinsData ...
func JenkinsData(config config.Config) []jenkins_types.Pipeline {
    var allJenkinsData []jenkins_types.Pipeline

    g := utils.Getable{
        URL:    "ALLJENKINSDATA",
        Config: config,
    }

    client := g.GetRedisClient()

    if client != nil {
        defer client.Close()

        cached, body := g.Cached(client, "ALLJENKINSDATA")
        json.Unmarshal(body, &allJenkinsData)

        if cached {
            return allJenkinsData
        }
    }

    allJenkinsData = report.Compile(config)

    body, err := json.Marshal(allJenkinsData)

    if err != nil {
        panic(err)
    }

    g.Cache(client, "ALLJENKINSDATA", body)

    return allJenkinsData
}

// UpdateData ...
func UpdateData(runConfig config.Config) {

    fmt.Printf("[%s] Starting Jenkins/CITH Scrape\n", time.Now().Format("2006-01-02 15:04:05"))

    cithFailures := CithFailures(runConfig)
    jenkinsData := JenkinsData(runConfig)

    var compiledData []jenkins_types.Pipeline
    if len(runConfig.CithURL) > 0 {
        compiledData = report.ApplyCith(jenkinsData, cithFailures)
    } else {
        compiledData = jenkinsData
    }

    report.WriteToCSV(compiledData)
    fmt.Printf("[%s] Ending Jenkins/CITH Scrape\n", time.Now().Format("2006-01-02 15:04:05") )
}


func main() {

    // See important notes in GetConfig w.r.t. configuration parameters and making
    // them consistent.
    runConfig := config.GetConfig()

    fmt.Printf ("Scrape Interval is every %d seconds\n", runConfig.ScrapeInterval)
    
    // Schedule UPdateData  to run every interval, with an initial first run.
    // These are run async so that we can run web server concurrently.
    s1 := gocron.NewScheduler(time.UTC)
    s1.Every(runConfig.ScrapeInterval).Seconds().StartImmediately().Do(UpdateData, runConfig)
    s1.StartAsync()

    // Run Web server in this thread until end.
    web.Serve()
}
