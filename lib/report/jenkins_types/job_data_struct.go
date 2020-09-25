package jenkins_types

import (
	"time"
)

type JobData struct {
    StartTime            time.Time
    EndTime              time.Time
    WallClockTimeMinutes int64
    WallClockTimeHours   int64
    TotalMinutes         int
    TotalHours           int
    QueueTimeHours       int
    QueueTimeMinutes     int
}

type JobDataStrings struct {
    StartTime            string
    EndTime              string
    WallClockTimeHours   string
    WallClockTimeMinutes string
    TotalHours           string
    TotalMinutes         string
    QueueTimeHours       string
    QueueTimeMinutes     string
}

func (jD *JobData) AssignJobValues(startTime int64, endTime int64, totalMinutes float32, queueTimeMinutes float32) {
    jD.StartTime = time.Unix(startTime/1000, 0)
    jD.EndTime = time.Unix(endTime/1000, 0)
    jD.WallClockTimeMinutes = (((endTime/1000 - startTime/1000) / 60) % 60)
    jD.WallClockTimeHours = ((endTime/1000 - startTime/1000) / 60 / 60)
    jD.TotalMinutes = int(totalMinutes) % 60
    jD.TotalHours = int(totalMinutes / 60)
    jD.QueueTimeMinutes = int(queueTimeMinutes) % 60
    jD.QueueTimeHours = int(queueTimeMinutes / 60)
}
