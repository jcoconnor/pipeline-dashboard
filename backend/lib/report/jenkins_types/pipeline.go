package jenkins_types

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Pipeline struct {
    JobData              JobData
    TrainData            map[int][]Train
    JobDataStrings       *JobDataStrings
    Server               string
    BuildNumber          int
    Pipeline             string
    PipelineJob          string
    Version              string
    URL                  string
    StartTime            time.Time
    EndTime              time.Time
    WallClockTimeMinutes int64
    WallClockTimeHours   int64
    TotalMinutes         int
    TotalHours           int
    QueueTimeMinutes     int
    QueueTimeHours       int
    Errors               int
    Transients           int
}

func (ajd *Pipeline) AssignVals() {
    var modifiedURL string
    var jobData []string

    modifiedURL = strings.TrimPrefix(ajd.URL, "https://")
    butURL := strings.Split(modifiedURL, "/api/json")

    modifiedURL = butURL[0]

    var server, pipeline string

    if strings.Contains(modifiedURL, "/job/") {
        jobData = strings.Split(modifiedURL, "/job/")
        server = jobData[0]
        pipeline = jobData[1]
    } else {
        jobData = strings.Split(modifiedURL, "/view/")
        server = jobData[0]
        pipeline = jobData[2]
        ajd.PipelineJob = jobData[1]
        ajd.Version = jobData[3]
    }

    ajd.Server = server
    ajd.Pipeline = pipeline
}

func (ajd *Pipeline) StringArray() []string {
    ajd.AssignVals()

    return []string{
        ajd.URL,
        ajd.Server,
        ajd.Pipeline,
        ajd.PipelineJob,
        ajd.Version,
        strconv.Itoa(ajd.BuildNumber),
        fmt.Sprintf("%s", ajd.JobData.StartTime),
        fmt.Sprintf("%s", ajd.JobData.EndTime),
        fmt.Sprintf("%d", ajd.JobData.WallClockTimeHours),
        fmt.Sprintf("%d", ajd.JobData.WallClockTimeMinutes),
        fmt.Sprintf("%d", ajd.JobData.TotalHours),
        fmt.Sprintf("%d", ajd.JobData.TotalMinutes),
        fmt.Sprintf("%d", ajd.JobData.QueueTimeHours),
        fmt.Sprintf("%d", ajd.JobData.QueueTimeMinutes),
        fmt.Sprintf("%d", ajd.Errors),
        fmt.Sprintf("%d", ajd.Transients),
    }
}

func FindJob(build Build, jd BuildsAndJobs) Job {
    var retVal Job
    for _, job := range jd.Jobs {
        if strings.Contains(build.URL, job.URL) {
            retVal = job
            break
        }
    }

    return retVal
}
