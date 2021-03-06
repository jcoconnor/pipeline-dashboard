package csv_writers

import (
	"encoding/csv"
	"os"

	"github.com/puppetlabs/pipeline-dashboard/backend/lib/report/constants"
	"github.com/puppetlabs/pipeline-dashboard/backend/lib/report/jenkins_types"
	"github.com/puppetlabs/pipeline-dashboard/backend/lib/report/utils"
)

func WritePipelines(allPipelines []jenkins_types.Pipeline) {
    file, err := os.OpenFile(constants.Results_filename, os.O_WRONLY|os.O_CREATE, 0644)
    utils.CheckError("Cannot create file", err)
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    writer.Write([]string{
        "URL",
        "Server",
        "Pipeline",
        "Pipeline Job",
        "Version",
        "Build Number",
        "Start Time",
        "End Time",
        "Wall Clock Time Hours",
        "Wall Clock Time Minutes",
        "Total Hours",
        "Total Minutes",
        "Queue Time Hours",
        "Queue Time Minutes",
        "Errors",
        "Transients",
    })

    for _, value := range allPipelines {
        err := writer.Write(value.StringArray())
        utils.CheckError("Cannot write to file", err)
    }
}
