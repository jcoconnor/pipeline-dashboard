package jenkins_types

type ViewsJenkinsData struct {
	Views []BuildsAndJobs `json: "views"`
}
