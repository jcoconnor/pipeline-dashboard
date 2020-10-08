/*
   Sample JSON being processed here:

   {
     "_class": "hudson.model.FreeStyleProject",
     "downstreamProjects": [
       {
         "_class": "hudson.model.FreeStyleProject",
         "name": "project_2019.1.x",
         "url": "https://jenkins.example.com/job/project_2019.1.x/"
       },
       {
         "_class": "hudson.model.FreeStyleProject",
         "name": "project_2019.1.x",
         "url": "https://jenkins.example.com/job/project_2019.1.x/"
       }
     ]
   }
*/

package runners

import (
	"encoding/json"

	"fmt"
	"strings"

	"github.com/puppetlabs/pipeline-dashboard/backend/lib/report/jenkins_types"
	"github.com/puppetlabs/pipeline-dashboard/backend/lib/report/utils"
)

type Ordered struct {
    utils.Getable
    LastBuild          jenkins_types.Build     `json:"lastBuild"`
    DownstreamProjects []jenkins_types.Project `json:"downstreamProjects"`
}

func FromOrdered(url string, job string, version string) jenkins_types.Pipeline {
    ko := Ordered{Getable: utils.Getable{URL: url}}
    JobData := ko.Pull(job, version)

    JobData.Version = version
    JobData.PipelineJob = job

    return JobData
}

func (k *Ordered) filter() string {
    return "/api/json?tree=lastBuild[url],downstreamProjects[name,url]"
}

func (o *Ordered) ProcessTopLevelBuilds(jd jenkins_types.BuildsAndJobs) jenkins_types.Builds {
    var retVal jenkins_types.Builds

    for _, build := range jd.Builds {
        jobForBuild := jenkins_types.FindJob(build, jd)

        build.Fetch()

        var buildDownstreamJobs []jenkins_types.Job
        if len(jobForBuild.DownstreamProjects) > 0 {
            _, buildDownstreamJobs = jenkins_types.JobsFromDownstreamProjects(jobForBuild.LastBuild, jobForBuild.DownstreamProjects)
        }

        for _, job := range buildDownstreamJobs {
            lastBuild := job.LastBuild
            if strings.Contains(lastBuild.URL, "http") {
                lastBuild.Fetch()

                if jenkins_types.BuildTriggerMatchesParent(lastBuild, build) {
                    retVal.List = append(retVal.List, lastBuild)
                }
            } else {
                continue
            }
        }
    }

    return retVal
}

// Pull will pull all jobs starting at the top., those jobs then need to be aggregated.
func (o *Ordered) Pull(pipeline_name string, pipeline_version string) jenkins_types.Pipeline {
    urlWithFilter := fmt.Sprintf(o.URL + o.filter())

    body := o.Get(urlWithFilter)

    json.Unmarshal(body, &o)

    // This needs to be processed by a job aggregator.

    var jd jenkins_types.BuildsAndJobs

    lastBuild := o.LastBuild
    lastBuild.Fetch()

    /*
     * This gets all the builds and jobs from the Ordered Job, for example:
     * https://jenkins.example.com/job/project_2016.4.x
     */
    jd.Jobs, jd.Builds = jenkins_types.OrderedJobsAndBuildsFromDownstreamProjects(lastBuild, o.DownstreamProjects)
    // jd.Builds = BuildsFromAllJobs(o.LastBuild, jd.Jobs)

    lastBuild.TopLevelBuild = true
    jd.Builds = append(jd.Builds, lastBuild)
    // subJobsBuilds := o.ProcessTopLevelBuilds(jd)

    subJobsBuilds := jenkins_types.Builds{List: jd.Builds}

    jobData, trainData := subJobsBuilds.GetJobData(pipeline_name, pipeline_version)

    /*
     *  Here is where I left off, I need to process all the builds for the sub jobs and turn it into job data.
     */

    return jenkins_types.Pipeline{JobData: jobData, URL: o.URL, BuildNumber: lastBuild.Number, TrainData: trainData}
}
