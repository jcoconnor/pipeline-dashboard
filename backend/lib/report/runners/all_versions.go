package runners

import (
	"encoding/json"
	"fmt"

	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/puppetlabs/pipeline-dashboard/backend/lib/report/jenkins_types"
	"github.com/puppetlabs/pipeline-dashboard/backend/lib/report/utils"
)

type VersionDataView struct {
    Class string `json: "_class"`
    Name  string `json: "name"`
    URL   string `json: "url"`
}

type VersionData struct {
    Views []VersionDataView `json: "views"`
}

// AllVersions iterates through a list of versions in a Jenkins View and runs the "OneVersion" script against it.
// Oneversion can probably go away in place of "ordered"
func AllVersions(url string, fromView bool, matchString string) []jenkins_types.Pipeline {
    var vd VersionData

    urlWithOptions := fmt.Sprintf("%s?tree=views[name,url]", url)
    resp, err := http.Get(urlWithOptions)

    utils.CheckError("Error Getting All Versions", err)

    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)

    err = json.Unmarshal(body, &vd)

    if err != nil {
        fmt.Println("error 1:", err)
    }

    var allLines []jenkins_types.Pipeline

    for _, v := range vd.Views {
        re := regexp.MustCompile(`(pe-[0-9]{4}|[0-9]+)\.[0-9]+\.x`)
        if re.MatchString(v.Name) {
            allLines = append(allLines, OneVersion(fmt.Sprintf("%sapi/json?depth=2", v.URL), fromView, matchString))
        }
    }

    return allLines
}

// OneVersion  Starts getting jenkins_types.Pipeline Data from a single version.
func OneVersion(url string, fromView bool, matchString string) jenkins_types.Pipeline {
    var urlWithOptions string
    if fromView {
        urlWithOptions = fmt.Sprintf("%s&tree=views[jobs[buildable,builds[number,duration,fullDisplayName,name,timestamp,url],lastBuild[duration,number,fullDisplayName,name,timestamp,url]]]", url)
    } else {
        urlWithOptions = fmt.Sprintf("%s&tree=jobs[buildable,builds[duration,number,fullDisplayName,name,timestamp,url],lastBuild]", url)
    }

    resp, err := http.Get(urlWithOptions)

    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)

    var jobTimingData jenkins_types.JobData

    preJobData := jenkins_types.Pipeline{URL: urlWithOptions}
    preJobData.AssignVals()

    var trainData map[int][]jenkins_types.Train

    if fromView {
        var vjd jenkins_types.ViewsJenkinsData
        err = json.Unmarshal(body, &vjd)
        for _, view := range vjd.Views {
            builds := view.GetBuilds()

            jobTimingData, trainData = builds.GetJobData(preJobData.PipelineJob, preJobData.Version)
        }

    } else {
        var jd jenkins_types.BuildsAndJobs
        err = json.Unmarshal(body, &jd)
        // jobTimingData = processView(jd, matchString)

        builds := jd.GetBuilds()

        jobTimingData, trainData = builds.GetJobData(preJobData.PipelineJob, preJobData.Version)
    }

    pipelineData := jenkins_types.Pipeline{JobData: jobTimingData, URL: urlWithOptions, TrainData: trainData}

    if err != nil {
        fmt.Println("error 2:", err)
    }

    return pipelineData
}
