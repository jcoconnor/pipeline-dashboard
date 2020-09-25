package jenkins_types

/*
 * _class: "jenkins.metrics.impl.TimeInQueueAction"
 */

type TimeInQueue struct {
    BlockedDurationMillis   int
    BlockedTimeMillis       int
    BuildableDurationMillis int
    BuildableTimeMillis     int
    BuildingDurationMillis  int
    ExecutingTimeMillis     int
    ExecutorUtilization     int
    SubTaskCount            int
    WaitingDurationMillis   int
    WaitingTimeMillis       int
}

// Milliseconds
func (t *TimeInQueue) QueueTime() int {
    // var project string
    return t.BlockedDurationMillis + t.BuildableDurationMillis
}
