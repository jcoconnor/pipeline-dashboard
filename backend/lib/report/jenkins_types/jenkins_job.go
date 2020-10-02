package jenkins_types

import (
	"encoding/json"

	"github.com/puppetlabs/pipeline-dashboard/backend/lib/report/utils"
)

type Job struct {
    utils.Getable
    Name               string  `json:"name"`
    Class              string  `json:"_class"`
    Buildable          bool    `json:"buildable"`
    Builds             []Build `json:"builds"`
    LastBuild          Build   `json:"lastBuild"`
    URL                string
    DownstreamProjects []Project `json:"downstreamProjects"`
    Runs               []Project `json:"runs"`
}

// func BuildsFromAllJobs(topBuild jenkins_types.Build, jobs []jenkins_types.Job) []jenkins_types.Build {
// 	var builds []jenkins_types.Build
//
// 	for _, dp := range jobs {
//
// 	}
//
// 	return builds
// }

/*
 * The purpose of this function is that when working through ordered jobs, we have to check that the sub job matches the parent job.
 * Since each job extends from the previous, we have to check it against the job in sequence, not completely seperately.
 */
func OrderedJobsAndBuildsFromDownstreamProjects(parentBuild Build, projects []Project) ([]Job, []Build) {
    var jobs []Job
    var builds []Build
    for _, dp := range projects {
        job := Job{URL: dp.URL}
        job.Fetch()

        // Add the first job in the list
        jobs = append(jobs, job)

        // Check to see which builds matches the aformentioned subjob
        matches, matchingBuild := AllBuildsTriggerMatchesParent(job, parentBuild)

        if matches {
            allSubJobs, allSubBuilds := job.AllJobs(matchingBuild)
            jobs = append(jobs, allSubJobs...)
            builds = append(builds, allSubBuilds...)
            builds = append(builds, matchingBuild)
        }
    }

    return jobs, builds
}

func JobsFromDownstreamProjects(parentBuild Build, projects []Project) ([]Build, []Job) {
    var jobs []Job
    var builds []Build

    for _, dp := range projects {
        job := Job{URL: dp.URL}
        job.Fetch()

        matches, matchedBuild := AllBuildsTriggerMatchesParent(job, parentBuild)
        if matches {
            builds = append(builds, matchedBuild)
            foundJobs, foundBuilds := job.AllJobs(matchedBuild)
            jobs = append(jobs, foundJobs...)
            builds = append(builds, foundBuilds...)
        }

    }

    return builds, jobs
}

func (j *Job) filter() string {
    return "api/json?tree=name,buildable,builds[duration,fullDisplayName,name,runs[number,url],timestamp,url],lastBuild[duration,fullDisplayName,name,runs[number,url],timestamp,url],downstreamProjects[url],runs[number,url]&depth=2"
}

func (j *Job) AllJobs(buildWhichMatchesParent Build) ([]Job, []Build) {
    var jobs []Job
    var builds []Build

    if len(j.DownstreamProjects) > 0 {
        foundJobs, foundBuilds := OrderedJobsAndBuildsFromDownstreamProjects(buildWhichMatchesParent, j.DownstreamProjects)
        jobs = append(jobs, foundJobs...)
        builds = append(builds, foundBuilds...)
    }

    jobs = append(jobs, *j)

    return jobs, builds
}

func (j *Job) Fetch() Job {
    urlWithFilter := j.URL + j.filter()

    body := j.Get(urlWithFilter)

    json.Unmarshal(body, &j)

    return *j
}
