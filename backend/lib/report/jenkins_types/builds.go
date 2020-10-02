package jenkins_types

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/table"
	"github.com/puppetlabs/pipeline-dashboard/backend/lib/report/utils"
)

type Builds struct {
    List []Build
}

func ProcessTopLevelBuilds(jd BuildsAndJobs) Builds {
    /*
     * First grab all top level jobs.
     */
    var retVal Builds

    for _, build := range jd.Builds {
        jobForBuild := FindJob(build, jd)

        if ( build.URL == "" || jobForBuild.URL == "") {
            continue
        }

        jobForBuild.Fetch()
        build.Fetch()

        // var buildDownstreamJobs []Job
        var buildDownstreamBuilds []Build

        if len(jobForBuild.DownstreamProjects) > 0 {
            buildDownstreamBuilds, _ = JobsFromDownstreamProjects(build, jobForBuild.DownstreamProjects)
        }

        build.JobName = jobForBuild.Name
        retVal.List = append(retVal.List, build)
        retVal.List = append(retVal.List, buildDownstreamBuilds...)
    }
    return retVal
}

// AllBuildsTriggerMatchesParent iterates through builds checking struct of a job and check
// if any matches the parent
func AllBuildsTriggerMatchesParent(child Job, parent Build) (bool, Build) {

    for _, build := range child.Builds {
        utils.LogTree(fmt.Sprintf("Trying %s", build.FullDisplayName), "", 1)
        if BuildTriggerMatchesParent(build, parent) {
            build.Fetch()
            utils.LogTree(fmt.Sprintf("Returning Build with Number %d", build.Number), "", 3)
            return true, build
        }
    }

    var build Build

    return false, build
}

func BuildTriggerMatchesParent(child Build, parent Build) bool {
    triggered_url_stub, triggered_by := child.TriggeredBy()
    if triggered_by == 0 {
        child.Fetch()
        triggered_url_stub, triggered_by = child.TriggeredBy()
    }

    if parent.Number == 0 || parent.URL == "" {
        parent.Fetch()
    }

    if parent.URL == "" {
        panic("No Parent URL")
    }
    utils.LogTree(fmt.Sprintf("Child was triggered by %d, parent is %d (%s should contain %s)", triggered_by, parent.Number, parent.URL, triggered_url_stub), "", 2)

    return (triggered_by == parent.Number) && strings.Contains(parent.URL, triggered_url_stub)
}

func LogTree(trainData map[int][]Train) {
    var jobs []string
    for _, topTrain := range trainData {
        for _, train := range topTrain {
            if train.JobName != "" {
                jobs = append(jobs, train.JobName)
            }
        }
    }

    for _, job := range jobs {
        utils.LogTree(job, "", 1)
        for _, topTrain := range trainData {
            for _, train := range topTrain {
                if train.JobName == job {
                    utils.LogTree(fmt.Sprintf("%s %d %s %s", train.Name, train.BuildNumber, train.StartTime, train.EndTime), "", 2)
                }
            }
        }
    }
}

func parseBuildToTrain(build Build) (train Train) {
    train.BuildNumber = build.Number
    train.DurationMinutes = float32(build.Duration) / (60 * 1000)
    train.EndTime = train.GetEndTime()
    train.JobName = build.JobName
    train.Name = build.FullDisplayName
    train.QueueTimeMinutes = float32(build.TimeInQueue.QueueTime()) / (60 * 1000)
    train.StartTime = time.Unix(build.Timestamp/1000, 0)
    train.Timestamp = build.Timestamp
    train.URL = build.URL

    return train
}

func addTrain(trains map[int][]Train, newTrain Train, i int) map[int][]Train {
    if len(trains[i]) > 0 {
        trains[i] = append(trains[i], newTrain)
    } else {
        trains[i] = []Train{newTrain}
    }

    return trains
}

// GetJobData gets trains from matrix builds and builds.
func (b *Builds) GetJobData(pipelineName string, pipelineVersion string) (JobData, map[int][]Train) {
    trainData := make(map[int][]Train)

    for _, build := range b.List {
        build.Fetch()
        i := 0

        if build.Class == "hudson.matrix.MatrixBuild" {
            // Here is where we get trains from Matrix Builds
            if len(build.Runs) > 0 {
                for _, cellBuild := range BuildsFromMatrixRuns(build, build.Runs) {
                    train := parseBuildToTrain(cellBuild)
                    trainData = addTrain(trainData, train, i)
                }
            }

            i++
        }

        train := parseBuildToTrain(build)
        trainData = addTrain(trainData, train, i)
        i++

    }

    var jobData JobData
    LogTree(trainData)

    var totalMinutes float32
    var queueTimeMinutes float32

    var startTime int64 = 9999999999999
    var endTime int64
    for _, train := range trainData {

        for _, t := range train {
            totalMinutes = totalMinutes + t.DurationMinutes
            queueTimeMinutes = queueTimeMinutes + t.QueueTimeMinutes

            timeOfEvent := time.Unix(t.Timestamp/1000, 0)

            if time.Now().Sub(timeOfEvent).Hours()/24 <= 365 {

                if startTime > t.Timestamp {
                    startTime = t.Timestamp
                }

                if endTime < t.EndTimeSeconds() {
                    endTime = t.EndTimeSeconds()
                }
            } 

        }
        jobData.AssignJobValues(startTime, endTime, totalMinutes, queueTimeMinutes)
    }

    fmt.Printf ("Summary table for %s/%s\n", pipelineName, pipelineVersion)
    t := table.NewWriter()
    t.SetOutputMirror(os.Stdout)
    t.AppendHeader(table.Row{"Start Time", "End Time", "Wall Clock Time Hours", "Wall Clock Time Minutes", "Queue Time Hours", "Queue Time Minutes"})
    t.AppendRows([]table.Row{
        {jobData.StartTime, jobData.EndTime, jobData.WallClockTimeHours, jobData.WallClockTimeMinutes, jobData.QueueTimeHours, jobData.QueueTimeMinutes},
    })
    t.Render()

    return jobData, trainData

}
