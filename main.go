package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/puppetlabs/pipeline-dashboard/config"
	"github.com/puppetlabs/pipeline-dashboard/lib/report"
	"github.com/puppetlabs/pipeline-dashboard/lib/report/cith"
	"github.com/puppetlabs/pipeline-dashboard/lib/report/constants"
	"github.com/puppetlabs/pipeline-dashboard/lib/report/jenkins_types"
	"github.com/puppetlabs/pipeline-dashboard/lib/report/utils"
	"github.com/puppetlabs/pipeline-dashboard/lib/web"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

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
        fmt.Println("Checking to see if CITHFAILURES are cached")

        var body []byte
        cached, body = g.Cached(client, "CITHFAILURES")
        json.Unmarshal(body, &cithFailures)
    } else {
        cached = false
    }

    if cached {
        fmt.Println("CITHFAILURES are cached")
        return cithFailures
    }
    fmt.Println("CITHFAILURES are not cached")

    cithFailures = report.CompileCith(config.CithURL)

    fmt.Printf("Found %d failures for today from Cith.\n", len(cithFailures))

    if len(cithFailures) == 0 {
        fmt.Println("Found no Cith Failures")
        return cithFailures
    }

    body, err := json.Marshal(cithFailures)

    if err != nil {
        panic(err)
    }

    fmt.Println("Caching Cith Failures")

    g.Cache(client, "CITHFAILURES", body)

    return cithFailures
}

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


func UpdateData(runConfig config.Config) {

    cithFailures := CithFailures(runConfig)
    jenkinsData := JenkinsData(runConfig)

    var compiledData []jenkins_types.Pipeline
    if len(runConfig.CithURL) > 0 {
        compiledData = report.ApplyCith(jenkinsData, cithFailures)
    } else {
        compiledData = jenkinsData
    }

    report.WriteToCSV(compiledData)
}


func main() {

    pflag.Bool("hidetreelog", false, "Whether or not to hide the long tree log")
    pflag.Bool("no-cache", true, "Whether or not to use a cache")
    pflag.Uint64("scrape-interval", constants.SCRAPE_INTERVAL, "Interval in Seconds to scrape Jenkins/CITH")
    pflag.Parse()
    viper.BindPFlags(pflag.CommandLine)
    runConfig := config.GetConfig()

    runConfig.SetUseCache(viper.GetBool("no-cache"))
    
    // Schedule UPdateData  to run every interval, with an initial first run.
    // These are run async so that we can run web server concurrently.
    s1 := gocron.NewScheduler(time.UTC)
    s1.Every(viper.GetUint64("scrape-interval")).Seconds().StartImmediately().Do(UpdateData, runConfig)
    s1.StartAsync()

    // Run Web server in this thread until end.
    web.Serve()
}
