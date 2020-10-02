package cith

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type CithFailure struct {
    BuildNumber  int
    CategoryName string
    CategorySlug string
    Master       string
    ProjectName  string
    YearMonthDay string
}

type CithFailureCauseData struct {
    Categories []struct {
        Name  string `json:"name"`
        Value string `json:"value"`
    } `json:"categories"`
    RecentOccurences []struct {
        ID           map[string]string `json:"_id"`
        ProjectName  string            `json:"projectName"`
        BuildNumber  int               `json:"buildNumber"`
        Master       string            `json:"master"`
        YearMonthDay string            `json:"yearMonthDay"`
        HostURL      string            `json:"hostURL"`
    } `json:"recent_occurrences"`
}

type DailyPrevalentFailureCause struct {
    CithURL string
    Count   int    `json:"count"`
    ID      string `json:"id"`
}

func (d *DailyPrevalentFailureCause) URL() string {
    return fmt.Sprintf("%s/api/v1/failure-cause/%s", d.CithURL, d.ID)
}

func (d *DailyPrevalentFailureCause) FailureCauseData() []CithFailure {

    resp, err := http.Get(d.URL())

    if err != nil {
        panic(err)
    }

    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)

    var cithData CithFailureCauseData

    json.Unmarshal(body, &cithData)

    var cithFailures []CithFailure

    for _, occurence := range cithData.RecentOccurences {
        cithFailures = append(cithFailures, CithFailure{
            CategoryName: cithData.Categories[0].Name,
            CategorySlug: cithData.Categories[0].Value,
            ProjectName:  occurence.ProjectName,
            BuildNumber:  occurence.BuildNumber,
            YearMonthDay: occurence.YearMonthDay,
            Master:       occurence.Master,
        })
    }

    return cithFailures
}

type DataByJenkinsServer struct {
    DailyFailureCauseOccurences int `json:"daily_failure_cause_occurences"`
    // DailyFailureByType map[string]interface{} `json:"daily_failures_by_type"`
    DailyPrevalentFailureCauses map[string]DailyPrevalentFailureCause `json:"daily_prevalent_failure_causes"`
}

type CithDashboardData struct {
    Data        map[string]DataByJenkinsServer
    RefreshRate int
}

func GetPrevalentFailureCauses(CithURL string) []CithFailure {
    currentTime := time.Now()
    today := currentTime.Format("2006-01-02")

    URL := fmt.Sprintf("%s/api/v1/jenkins/statistics-for-dashboard?from=%s&to=%s", CithURL, today, today)

    resp, err := http.Get(URL)

    if err != nil {
        panic(err)
    }

    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)

    var cithData CithDashboardData

    json.Unmarshal(body, &cithData)

    var allFailures []CithFailure

    for _, v := range cithData.Data["All Jenkins instances"].DailyPrevalentFailureCauses {
        v.CithURL = CithURL
        allFailures = append(allFailures, v.FailureCauseData()...)
    }

    return allFailures
}

func Remove(slice []CithFailure, s int) []CithFailure {
    if s+1 <= len(slice) {
        return append(slice[:s], slice[s+1:]...)
    }
    return slice[:s]
}
