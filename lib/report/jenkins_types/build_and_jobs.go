package jenkins_types

type BuildsAndJobs struct {
    Builds []Build `json:"builds"`
    Jobs   []Job   `json:"jobs"`
}

func (baj *BuildsAndJobs) GetBuilds() Builds {
    var builds Builds

    for _, job := range baj.Jobs {
        if len(job.LastBuild.URL) == 0 {
            continue
        }
        build := job.LastBuild

        build.Fetch()

        builds.List = append(builds.List, build)
    }

    return builds
}
