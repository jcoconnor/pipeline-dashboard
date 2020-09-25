package jenkins_types

// Action should match up exactly with the Jenkins "action" schema
type Action struct {
    Class                   string   `json:"_class"`
    TriggeredBuilds         []Build  `json:"triggeredBuilds"`
    Causes                  []Causes `json:"causes"`
    BlockedDurationMillis   int      `json:"blockedDurationMillis"`
    BlockedTimeMillis       int      `json:"blockedTimeMillis"`
    BuildableDurationMillis int      `json:"buildableDurationMillis"`
    BuildableTimeMillis     int      `json:"buildableTimeMillis"`
    BuildingDurationMillis  int      `json:"buildingDurationMillis"`
    ExecutingTimeMillis     int      `json:"executingTimeMillis"`
    WaitingDurationMillis   int      `json:"waitingDurationMillis"`
    WaitingTimeMillis       int      `json:"waitingTimeMillis"`
}

// Causes should match up exactly with the Jenkins "causes" schema
type Causes struct {
    UpstreamBuild   int    `json:"upstreamBuild"`
    UpstreamProject string `json:"upstreamProject"`
    UpstreamURL     string `json:"upstreamURL"`
}
