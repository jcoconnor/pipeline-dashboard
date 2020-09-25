package jenkins_types

import (
	"testing"
)

func TestUnmarshallingJson(t *testing.T) {
    jsonBody := []byte("{\"_class\":\"hudson.model.FreeStyleBuild\",\"actions\":[{\"_class\":\"hudson.model.ParametersAction\"},{},{\"_class\":\"hudson.model.CauseAction\",\"causes\":[{\"_class\":\"hudson.model.Cause$UpstreamCause\",\"shortDescription\":\"Started by upstream project \\\"upstream_project_job_2016.4.x\\\" build number 54\",\"upstreamBuild\":54,\"upstreamProject\":\"upstream_project_job_2016.4.x\",\"upstreamUrl\":\"job/upstream_project_job_2016.4.x/\"}]},{\"_class\":\"jenkins.metrics.impl.TimeInQueueAction\",\"blockedDurationMillis\":0,\"blockedTimeMillis\":0,\"buildableDurationMillis\":3789,\"buildableTimeMillis\":3789,\"buildingDurationMillis\":3470,\"executingTimeMillis\":3470,\"waitingDurationMillis\":0,\"waitingTimeMillis\":0},{},{},{},{},{},{},{}],\"duration\":3470,\"fullDisplayName\":\"pe-acceptance-tests (2016.4.x) Init 31: Analytics Integration Init pre-calculates resources Kickoff #54\",\"number\":54,\"timestamp\":1586224930523,\"url\":\"https://jenkins.example.com/job/upstream_project_job_2016.4.x/54/\"}")

    var jenkinsBuild Build

    jenkinsBuild.UnmarshalFetchedBuild(jsonBody)

    if jenkinsBuild.TimeInQueue.BuildableDurationMillis != 3789 {
        t.Errorf("BuildableDurationMillis is wrong")
    }
    if jenkinsBuild.TimeInQueue.BuildableTimeMillis != 3789 {
        t.Errorf("BuildableTimeMillis is wrong")
    }

    if jenkinsBuild.TimeInQueue.BuildingDurationMillis != 3470 {
        t.Errorf("BuildingDurationMillis is wrong")
    }

    if jenkinsBuild.TimeInQueue.ExecutingTimeMillis != 3470 {
        t.Errorf("ExecutingTimeMillis is wrong")
    }

}
