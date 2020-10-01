package jenkins_types

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/puppetlabs/pipeline-dashboard/lib/report/constants"
	"github.com/puppetlabs/pipeline-dashboard/lib/report/utils"
)

type Train struct {
    JobName          string
    BuildNumber      int
    URL              string
    Name             string
    DurationMinutes  float32
    QueueTimeMinutes float32
    StartTime        time.Time
    EndTime          time.Time
    Timestamp        int64
    Errors           int
    Transients       int
    Platform         string
    PlatformVersion  string
}

type TrainStrings struct {
    URL                  string
    Pipeline             string
    Version              string
    Name                 string
    DurationHours        string
    DurationMinutes      string
    DurationSortMinutes  int
    QueueTimeSortMinutes int
    QueueTimeMinutes     string
    QueueTimeHours       string
    StartTime            string
    EndTime              string
    Timestamp            string
    Errors               int
    Transients           int
    Platform             string
    PlatformVersion      string
}

func (t *Train) EndTimeSeconds() int64 {
    return t.Timestamp + (int64(t.DurationMinutes) * 60 * 1000)
}

func (t *Train) StringArray() []string {
    return []string{
        t.URL,
        t.Name,
        fmt.Sprintf("%g", t.DurationMinutes),
        fmt.Sprintf("%g", t.QueueTimeMinutes),
        fmt.Sprintf("%s", t.StartTime),
        fmt.Sprintf("%s", t.EndTime),
        fmt.Sprintf("%d", t.Timestamp),
        fmt.Sprintf("%d", t.Errors),
        fmt.Sprintf("%d", t.Transients),
        t.GetPlatform(),
        t.GetPlatformVersion(),
    }
}

// GetPlatform pulls a platform the from the train URL.
func (t *Train) GetPlatform() string {
    osList := "aix|osx|fedora|redhat|redhatfips|centos|sles|ubuntu|debian|windows|windowsfips|solaris|amazon|arista"

    var envVar string
    if strings.Contains(t.URL, "TEST_TARGET") {
        envVar = "TEST_TARGET"
    }

    if strings.Contains(t.URL, "BUILD_TARGET") {
        envVar = "BUILD_TARGET"
    }

    if strings.Contains(t.URL, "LAYOUT") {
        envVar = "LAYOUT"
    }

    if envVar != "" {
        var rgx = regexp.MustCompile(fmt.Sprintf("%s=(%s)", envVar, osList))
        rs := rgx.FindStringSubmatch(t.URL)
        if len(rs) > 1 && len(rs[1]) > 0 {
            return rs[1]
        }
        log.Printf("&&&&&&&&&%s Could not pull platform data from URL %s", envVar, t.URL)
    }

    if strings.Contains(t.URL, "redhat") {
        log.Printf("&&&&&&&&&ad Could not pull platform data from URL %s", t.URL)
    }
    return ""
}

// GetPlatformVersion pulls a platform version string from the from the train URL.
func (t *Train) GetPlatformVersion() string {
    osList := "aix|osx|el|fedora|redhat|redhatfips|centos|sles|ubuntu|debian|windows|windowsfips|solaris|amazon|arista"
    versionsList := "[0-9]*|-[0-9*]|-6.1|10ent|2008r2|2012r2|2019_ja|-2012r2"
    platformsList := "-i386|-x86|-x64|-ppc64le|-amd64|-x86_64|-32bolt|-64bolt|-64mda|-64mcd|-64m|-64a|-POWERfa|-32a|-POWERa|-AARCH64a|-aarch64|-ppc"

    var envVar string
    if strings.Contains(t.URL, "TEST_TARGET") {
        envVar = "TEST_TARGET"
    }

    if strings.Contains(t.URL, "BUILD_TARGET") {
        envVar = "BUILD_TARGET"
    }

    if strings.Contains(t.URL, "LAYOUT") {
        envVar = "LAYOUT"
    }

    if envVar != "" {
        var rgx = regexp.MustCompile(fmt.Sprintf("%s=(%s)(%s)(%s)", envVar, osList, versionsList, platformsList))
        rs := rgx.FindStringSubmatch(t.URL)
        if len(rs) > 1 && len(rs[2]) > 0 && len(rs[3]) > 0 {
            return fmt.Sprintf("%s%s", rs[2], rs[3])
        }
        log.Printf("&&&&&&&&&%s Could not pull platform data from URL %s", envVar, t.URL)
    }
    if strings.Contains(t.URL, "redhat") {
        log.Printf("&&&&&&&&&bd Could not pull platform data from URL %s", t.URL)
    }
    return ""
}

func (t *Train) GetEndTime() time.Time {
    retVal := t.StartTime.Add(time.Minute * time.Duration(t.DurationMinutes))
    fmt.Println(t.DurationMinutes)
    fmt.Printf("%s %s\n", t.StartTime, retVal)

    return retVal
}

func OpenTrainCSV() *csv.Writer {
    file, err := os.OpenFile(constants.Trains_filename, os.O_WRONLY|os.O_CREATE, 0644)
    utils.CheckError("Cannot create file", err)

    file.WriteString("\n")

    writer := csv.NewWriter(file)
    writer.Write([]string{
        "Job Name",
        "Job Version",
        "URL",
        "Name",
        "Duration Minutes",
        "Queue Time Minutes",
        "Start Time",
        "End Time",
        "Timestamp",
        "Errors",
        "Transients",
        "Platform",
        "PlatformVersion",
    })
    return writer
}

func WriteTrainCSV(writer *csv.Writer, train Train, jobName string, jobVersion string) {
    writer.Write(append([]string{jobName, jobVersion}, train.StringArray()...))
    writer.Flush()
}
