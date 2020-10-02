package jenkins_types

import (
	"testing"
)

func TestApplyingValues(t *testing.T) {
    var jobData JobData

    startTime := int64(1573466400 * 1000) // 2019-11-11 10:00 AM
    endTime := int64(1573472700 * 1000)   // 2019-11-11 11:45 AM
    totalMinutes := float32(400)
    queueTime := float32(100)

    jobData.AssignJobValues(startTime, endTime, totalMinutes, queueTime)

    if jobData.StartTime.Format("2006-01-02 15:04:05") != "2019-11-11 02:00:00" && jobData.StartTime.Format("2006-01-02 15:04:05") != "2019-11-11 10:00:00" {
        t.Errorf("Start Time is wrong, got %s", jobData.StartTime.Format("2006-01-02 15:04:05"))
    }
    if jobData.EndTime.Format("2006-01-02 15:04:05") != "2019-11-11 03:45:00" && jobData.EndTime.Format("2006-01-02 15:04:05") != "2019-11-11 11:45:00" {
        t.Errorf("End Time is wrong, got %s", jobData.EndTime.Format("2006-01-02 15:04:05"))
    }

    if jobData.WallClockTimeMinutes != 45 {
        t.Errorf("WallClockTimeMinutes is wrong, got %d", jobData.WallClockTimeMinutes)
    }

    if jobData.WallClockTimeHours != 1 {
        t.Errorf("WallClockTimeHours is wrong")
    }

    if jobData.TotalMinutes != 40 {
        t.Errorf("TotalMinutes is wrong")
    }

    if jobData.TotalHours != 6 {
        t.Errorf("TotalHours is wrong")
    }

    if jobData.QueueTimeMinutes != 40 {
        t.Errorf("QueueTimeHours is wrong, got %d", jobData.QueueTimeMinutes)
    }

    if jobData.QueueTimeHours != 1 {
        t.Errorf("QueueTimeHours is wrong, got %d", jobData.QueueTimeHours)
    }

}
