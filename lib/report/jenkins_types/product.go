package jenkins_types

import (
	"fmt"
	"strconv"
	"time"

	"github.com/puppetlabs/pipeline-dashboard/config"
)

type Product struct {
	Name                 string
	Pipeline             string
	WallClockTime        time.Duration
	TotalTime            string
	StartTime            time.Time
	EndTime              time.Time
	WallClockTimeMinutes int
	TotalTimeMinutes     int
	TotalTimeDuration    string
	Errors               int
	Transients           int
}

func GetProducts() []Product {
	configData := config.GetConfig()

	var retVal []Product
	for _, product := range configData.Products {
		retVal = append(retVal, Product{
			Name:     product.Name,
			Pipeline: product.Pipeline,
		})
	}

	return retVal
}

func (p *Product) SetVals(jobs []Pipeline) {
	timeFormat := "2006-01-02 15:04:05 -0700 MST"
	// 2019-09-06 10:45:32 -0700 PDT

	p.StartTime = time.Now().AddDate(0, 0, 365)
	p.EndTime = time.Now().AddDate(0, 0, -365)

	p.TotalTimeMinutes = 0
	p.Errors = 0
	p.Transients = 0

	for _, job := range jobs {
		if job.Pipeline == p.Pipeline {
			jobStartTime, err := time.Parse(timeFormat, job.JobDataStrings.StartTime)
			if err != nil {
				fmt.Println(err)
			}

			jobEndTime, err := time.Parse(timeFormat, job.JobDataStrings.EndTime)
			if err != nil {
				fmt.Println(err)
			}

			if jobStartTime.Before(p.StartTime) && jobStartTime.After(p.StartTime.AddDate(0, 0, -1825)) {
				fmt.Println("Changing Start")
				p.StartTime = jobStartTime
			}
			if jobEndTime.After(p.EndTime) && jobEndTime.After(p.EndTime.AddDate(0, 0, -1825)) {
				fmt.Println("Changing End")
				p.EndTime = jobEndTime
			}

			p.WallClockTime = p.EndTime.Sub(p.StartTime)

			totalJobMinutes, _ := strconv.Atoi(job.JobDataStrings.TotalMinutes)
			totalJobHours, _ := strconv.Atoi(job.JobDataStrings.TotalHours)

			p.Errors = p.Errors + job.Errors
			p.Transients = p.Transients + job.Transients

			p.TotalTimeMinutes = p.TotalTimeMinutes + totalJobMinutes + totalJobHours*60
		}
	}

	duration, _ := time.ParseDuration(fmt.Sprintf("%dm", p.TotalTimeMinutes))
	p.TotalTimeDuration = fmt.Sprintf("%s", duration)

}