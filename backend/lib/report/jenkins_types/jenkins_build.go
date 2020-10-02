package jenkins_types

import (
	"encoding/json"
	"fmt"

	"github.com/puppetlabs/pipeline-dashboard/backend/lib/report/utils"
)

type Build struct {
    utils.Getable
    JobName         string
    Actions         []Action  `json:"actions"`
    Class           string    `json:"_class"`
    Number          int       `json:"number"`
    Duration        int       `json:"duration"`
    Name            string    `json:"name"`
    FullDisplayName string    `json:"fullDisplayName"`
    Runs            []Project `json:"runs"`
    Timestamp       int64     `json:"timestamp"`
    URL             string    `json:"url"`
    TopLevelBuild   bool
    TimeInQueue     TimeInQueue
}

func TriggeredBuildsFromActions(actions []Action) []Build {
    var retVal []Build
    for _, a := range actions {
        if a.Class == "hudson.plugins.parameterizedtrigger.BuildInfoExporterAction" {
            retVal = a.TriggeredBuilds
        }
    }

    return retVal
}

func BuildsFromMatrixRuns(parent Build, projects []Project) []Build {
    var builds []Build
    for _, dp := range projects {

        build := Build{URL: dp.URL}
        build.Fetch()

        if BuildTriggerMatchesParent(build, parent) {
            builds = append(builds, build)
        }

    }

    return builds
}

func (b *Build) filter() string {
    waitingStatNames := "blockedDurationMillis,blockedTimeMillis,buildableDurationMillis,buildableTimeMillis,buildingDurationMillis,executingTimeMillis,waitingDurationMillis,waitingTimeMillis"
    return fmt.Sprintf("api/json?tree=actions[triggeredBuilds[*],causes[*],%s],duration,fullDisplayName,name,number,timestamp,url,runs[number,url]&depth=2", waitingStatNames)
}

func (b *Build) SetQueueTimes() {
    // var project string
    for _, action := range b.Actions {
        if action.Class == "jenkins.metrics.impl.TimeInQueueAction" {

            b.TimeInQueue.BlockedDurationMillis = action.BlockedDurationMillis
            b.TimeInQueue.BlockedTimeMillis = action.BlockedTimeMillis
            b.TimeInQueue.BuildableDurationMillis = action.BuildableDurationMillis
            b.TimeInQueue.BuildableTimeMillis = action.BuildableTimeMillis
            b.TimeInQueue.BuildingDurationMillis = action.BuildingDurationMillis
            b.TimeInQueue.ExecutingTimeMillis = action.ExecutingTimeMillis
            b.TimeInQueue.WaitingDurationMillis = action.WaitingDurationMillis
            b.TimeInQueue.WaitingTimeMillis = action.WaitingTimeMillis
        }
    }
}

func (b *Build) TriggeredBy() (string, int) {
    // var project string
    var buildNumber int
    var url string
    for _, action := range b.Actions {
        if action.Class == "hudson.model.CauseAction" {
            // project = action.Causes[0].UpstreamProject
            buildNumber = action.Causes[0].UpstreamBuild
            url = action.Causes[0].UpstreamURL

            return url, buildNumber
        }
    }

    return url, buildNumber
}

func (b *Build) UnmarshalFetchedBuild(jsonBody []byte) {
    json.Unmarshal(jsonBody, &b)
    b.SetQueueTimes()
}

func (b *Build) Fetch() Build {

    urlWithFilter := fmt.Sprint(b.URL + b.filter())

    if (b.URL != "") {
        body := b.Get(urlWithFilter)

        b.UnmarshalFetchedBuild(body)
    }

    return *b
}
