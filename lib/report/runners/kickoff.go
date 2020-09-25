/*
    file: kickoff_job.go
    author: Mikker Gimenez-Peterson

    {
      "_class": "hudson.model.FreeStyleProject",
      "downstreamProjects": [
        {
          "_class": "hudson.model.FreeStyleProject",
          "name": "job_2019.1.x",
          "url": "https://jenkins.example.com/job/job_2019.1.x/"
        },
        {
          "_class": "hudson.model.FreeStyleProject",
          "name": "project_2019.1.x",
          "url": "https://jenkins.example.com/job/project_2019.1.x/"
        }
      ]
    }

 When passed in a kickoff job like https://jenkins.example.com/job/job_2019.1.x/

 this will grab a list of all sub-jobs from that kickoff job.
*/
package runners

import (
	"encoding/json"
	"fmt"

	"github.com/puppetlabs/pipeline-dashboard/lib/report/jenkins_types"
	"github.com/puppetlabs/pipeline-dashboard/lib/report/utils"
)

func FromKickoff(url string, job string, version string) jenkins_types.Pipeline {
    utils.LogHeading(fmt.Sprintf("Starting Kickoff Job with URL: %s", url), "")
    ko := Kickoff{Getable: utils.Getable{URL: url}}

    body := ko.Pull()
    ko.Unmarshal(body)
    jd := ko.PullDownstreams()
    jobData := ko.ProcessBuilds(jd, job, version)

    jobData.Version = version
    jobData.PipelineJob = job

    return jobData
}

type Kickoff struct {
    utils.Getable
    LastBuild          jenkins_types.Build     `json:"lastBuild"`
    DownstreamProjects []jenkins_types.Project `json:"downstreamProjects"`
}

func (k *Kickoff) filter() string {
    return "/api/json?tree=lastBuild[url],downstreamProjects[name,url]"
}

// Pull will pull all jobs in a kickoff, those jobs then need to be aggregated.
func (k *Kickoff) Pull() []byte {
    urlWithFilter := fmt.Sprintf(k.URL + k.filter())
    body := k.Get(urlWithFilter)

    return body
}

func (k *Kickoff) Unmarshal(body []byte) {
    json.Unmarshal(body, &k)
}

func (k *Kickoff) PullDownstreams() jenkins_types.BuildsAndJobs {

    // This needs to be processed by a job aggregator.

    var jd jenkins_types.BuildsAndJobs

    k.LastBuild.Fetch()

    utils.Log(fmt.Sprintf("Last build for kickoff job is %d", k.LastBuild.Number), "")

    /*
     * This gets all the builds and jobs from the Kickoff Job, for example:
     * https://jenkins.example.com/job/job_2019.1.x
     */
    jd.Builds = jenkins_types.TriggeredBuildsFromActions(k.LastBuild.Actions)
    _, jd.Jobs = jenkins_types.JobsFromDownstreamProjects(k.LastBuild, k.DownstreamProjects)

    return jd
}

func (k *Kickoff) ProcessBuilds(jd jenkins_types.BuildsAndJobs, pipeline_name string, pipeline_version string) jenkins_types.Pipeline {
    utils.LogTree(fmt.Sprintf("Found %d Jobs\n", len(jd.Jobs)), "", 1)
    subJobsBuilds := jenkins_types.ProcessTopLevelBuilds(jd)

    jobData, trainData := subJobsBuilds.GetJobData(pipeline_name, pipeline_version)

    /*
     *  Here is where I left off, I need to process all the builds for the sub jobs and turn it into job data.
     */

    return jenkins_types.Pipeline{JobData: jobData, URL: k.URL, BuildNumber: k.LastBuild.Number, TrainData: trainData}
}
