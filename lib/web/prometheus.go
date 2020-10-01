package web

import (
	"fmt"

	"net/http"
	"regexp"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func (handlers *Handlers) GenerateMetrics(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Our middleware logic goes here...
        handlers.GeneratePageData()
        if (handlers.Page == nil) {
            // Not totally sure if this is the best treatment here, but
            // in testing it seems to avoid the stack dump errors on the
            // console and web requests to http://localhost:8080/metrics seem
            // to behave beningly (i.e. return no data)
            return 
        }
        // TBD - Check that Page is Not Null - or do we wait and loop ?
        // Probably a Null check and return - safer.
        for _, job := range handlers.Page.Jobs {
            tH, _ := strconv.Atoi(job.JobDataStrings.TotalHours)
            tM, _ := strconv.Atoi(job.JobDataStrings.TotalMinutes)
            wCTH, _ := strconv.Atoi(job.JobDataStrings.WallClockTimeHours)
            wCTM, _ := strconv.Atoi(job.JobDataStrings.WallClockTimeMinutes)
            wQTH, _ := strconv.Atoi(job.JobDataStrings.QueueTimeHours)
            wQTM, _ := strconv.Atoi(job.JobDataStrings.QueueTimeMinutes)

            queueTimeSeconds.
                With(prometheus.Labels{
                    "pipeline":         job.Pipeline,
                    "pipeline_job":     job.PipelineJob,
                    "version":          job.Version,
                    "build_number":     strconv.Itoa(job.BuildNumber),
                    "scope":            "pipeline",
                    "platform":         "n/a",
                    "platform_version": "n/a",
                }).
                Set(float64((60 * wQTM) + (3600 * wQTH)))
            totalSeconds.
                With(prometheus.Labels{
                    "pipeline":         job.Pipeline,
                    "pipeline_job":     job.PipelineJob,
                    "version":          job.Version,
                    "build_number":     strconv.Itoa(job.BuildNumber),
                    "scope":            "pipeline",
                    "platform":         "n/a",
                    "platform_version": "n/a",
                }).
                Set(float64((60 * tM) + (3600 * tH)))
            wallClockSeconds.
                With(prometheus.Labels{
                    "pipeline":         job.Pipeline,
                    "pipeline_job":     job.PipelineJob,
                    "version":          job.Version,
                    "build_number":     strconv.Itoa(job.BuildNumber),
                    "scope":            "pipeline",
                    "platform":         "n/a",
                    "platform_version": "n/a",
                }).
                Set(float64((60 * wCTM) + (3600 * wCTH)))
        }

        for _, train := range handlers.Page.Trains {
            var matchBuildNumber = regexp.MustCompile(` #[0-9]+$`)

            var trainNameWithoutBuildNumber = string(matchBuildNumber.ReplaceAll([]byte(train.Name), []byte("")))
            fmt.Println(trainNameWithoutBuildNumber)

            queueTimeSeconds.
                With(prometheus.Labels{
                    "pipeline":         train.Pipeline,
                    "pipeline_job":     trainNameWithoutBuildNumber,
                    "version":          train.Version,
                    "build_number":     "0",
                    "scope":            "job",
                    "platform":         train.Platform,
                    "platform_version": train.PlatformVersion,
                }).
                Set(float64(train.QueueTimeSortMinutes * 60))
            totalSeconds.
                With(prometheus.Labels{
                    "pipeline":         train.Pipeline,
                    "pipeline_job":     trainNameWithoutBuildNumber,
                    "version":          train.Version,
                    "build_number":     "0",
                    "scope":            "job",
                    "platform":         train.Platform,
                    "platform_version": train.PlatformVersion,
                }).
                Set(float64(train.DurationSortMinutes * 60))
            wallClockSeconds.
                With(prometheus.Labels{
                    "pipeline":         train.Pipeline,
                    "pipeline_job":     trainNameWithoutBuildNumber,
                    "version":          train.Version,
                    "build_number":     "0",
                    "scope":            "job",
                    "platform":         train.Platform,
                    "platform_version": train.PlatformVersion,
                }).
                Set(float64(train.DurationSortMinutes * 60))
        }

        timeSinceLastUpdate.
            With(prometheus.Labels{}).Set(lastUpdated())

        next.ServeHTTP(w, r)
    })
}

var (
    queueTimeSeconds = promauto.NewGaugeVec(prometheus.GaugeOpts{
        Name: "cidashboard_queue_time_seconds",
        Help: "The amount of time a job is queued.",
    }, []string{"pipeline", "pipeline_job", "version", "build_number", "scope", "platform", "platform_version"},
    )

    wallClockSeconds = promauto.NewGaugeVec(prometheus.GaugeOpts{
        Name: "cidashboard_wall_clock_seconds",
        Help: "The time it has taken for a job or pipeline from start to finish.",
    }, []string{"pipeline", "pipeline_job", "version", "build_number", "scope", "platform", "platform_version"},
    )
    totalSeconds = promauto.NewGaugeVec(prometheus.GaugeOpts{
        Name: "cidashboard_total_seconds",
        Help: "The consecutive time taken for a job or pipeline",
    }, []string{"pipeline", "pipeline_job", "version", "build_number", "scope", "platform", "platform_version"},
    )
    timeSinceLastUpdate = promauto.NewGaugeVec(prometheus.GaugeOpts{
        Name: "cidashboard_seconds_since_last_update",
        Help: "The amount of time that has passed since the last update.",
    }, []string{},
    )
)
