/*
 * START HERE The Compile below com
 */
package report

import (
	"fmt"

	"strings"

	"github.com/puppetlabs/pipeline-dashboard/config"
	"github.com/puppetlabs/pipeline-dashboard/lib/report/cith"
	"github.com/puppetlabs/pipeline-dashboard/lib/report/csv_writers"
	"github.com/puppetlabs/pipeline-dashboard/lib/report/jenkins_types"
	"github.com/puppetlabs/pipeline-dashboard/lib/report/runners"
	"github.com/puppetlabs/pipeline-dashboard/lib/report/utils"
)

func CompileCith(CithURL string) []cith.CithFailure {
    return cith.GetPrevalentFailureCauses(CithURL)
}

// Compile is the entrypoint to the Jenkins Job Compilation.
func Compile(configData config.Config) []jenkins_types.Pipeline {
    var allPipelines []jenkins_types.Pipeline
    // Pass in kickoff job and follow downstreams

    for _, koJob := range configData.KickoffJobs {
        allPipelines = append(allPipelines, runners.FromKickoff(koJob.URL, koJob.Product, koJob.Version))
    }

    for _, ordJob := range configData.OrderedJobs {
        allPipelines = append(allPipelines, runners.FromOrdered(ordJob.URL, ordJob.Product, ordJob.Version))
    }

    return allPipelines
}

func ClosestMatch(failure cith.CithFailure, pipelines []jenkins_types.Pipeline) string {
    charsMatched := 0
    var matchedJob jenkins_types.Train

    for _, pipeline := range pipelines {
        for _, jobs := range pipeline.TrainData {
            for _, job := range jobs {

                tmpChars := utils.CharactersTheSame(job.URL, failure.ProjectName)
                if tmpChars > charsMatched {
                    charsMatched = tmpChars
                    matchedJob = job
                }
            }
        }
    }

    return fmt.Sprintf("Because %d character matched: %+v\n\n", charsMatched, matchedJob)
}

func URLEncodedContains(url string, testVal string) bool {
    if strings.Contains(url, "%252C") {
        url = strings.ReplaceAll(url, "%252C", "%2C")
    }

    return strings.Contains(url, testVal)
}

func ApplyCith(pipelines []jenkins_types.Pipeline, cithFailures []cith.CithFailure) []jenkins_types.Pipeline {
    var retVal []jenkins_types.Pipeline

    transients := []string{
        "Infrastructure (Suspected)",
        "Infrastructure (Confirmed)",
    }

    errors := []string{
        "Test (Suspected)",
        "Other",
        "Product (Suspected)",
        "Product (Confirmed)",
    }

    fmt.Printf("# of Failures: %d", len(cithFailures))

    for _, pipeline := range pipelines {
        fmt.Println("=============== A Pipeline =============")

        for ji, jobs := range pipeline.TrainData {
            fmt.Printf("# of Jobs: %d\n", len(jobs))
            for jj, job := range jobs {
                for ci, failure := range cithFailures {

                    if utils.CharactersTheSame(job.URL, failure.ProjectName) > 25 {
                        if URLEncodedContains(job.URL, failure.Master) {
                            fmt.Printf("%s and %s have %d characters the same ", job.URL, failure.ProjectName, utils.CharactersTheSame(job.URL, failure.ProjectName))
                            fmt.Printf("and are the same.\n\n")
                        } else {
                            fmt.Println("")
                            fmt.Println("")
                        }
                    }

                    if (job.BuildNumber == failure.BuildNumber || job.BuildNumber-1 == failure.BuildNumber || job.BuildNumber-2 == failure.BuildNumber) && URLEncodedContains(job.URL, failure.Master) && URLEncodedContains(job.URL, failure.ProjectName) {
                        if utils.StringSliceContains(transients, failure.CategoryName) {
                            pipeline.Transients++
                            pipeline.TrainData[ji][jj].Transients++
                            cithFailures = cith.Remove(cithFailures, ci)
                        } else if utils.StringSliceContains(errors, failure.CategoryName) {
                            pipeline.Errors++
                            pipeline.TrainData[ji][jj].Errors++
                            cithFailures = cith.Remove(cithFailures, ci)
                        } else {
                            panic(fmt.Sprintf("%s is not a transient or normal error", failure.CategoryName))
                        }
                    }
                }
            }
        }
        retVal = append(retVal, pipeline)
    }

    fmt.Printf("The following %d failures did not have any matches: ", len(cithFailures))
    for _, failure := range cithFailures {
        fmt.Printf("FAILURE: %+v\n", failure)
        fmt.Printf("Closest Match: %s", ClosestMatch(failure, pipelines))

    }

    // Locks added around file writes to prevent issues with updater.
    utils.Lock()
    defer utils.Unlock()
    csv := jenkins_types.OpenTrainCSV()

    for _, pipeline := range pipelines {
        for _, jobs := range pipeline.TrainData {
            for _, job := range jobs {
                jenkins_types.WriteTrainCSV(csv, job, pipeline.PipelineJob, pipeline.Version)
            }
        }
    }

    return retVal

}

func WriteToCSV(pipelines []jenkins_types.Pipeline) {
    // Locks added around file writes to prevent issues with updater.
    utils.Lock()
    defer utils.Unlock()
    csv_writers.WritePipelines(pipelines)
}
