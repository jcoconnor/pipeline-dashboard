package web

import (
	"os"
	"time"

	"github.com/puppetlabs/pipeline-dashboard/backend/lib/report/constants"
)

func lastUpdated() float64 {

    statinfo, err := os.Stat(constants.Results_filename)
    if os.IsNotExist(err) {
        return 0
    }
    intLastUpdated := statinfo.ModTime().Unix()

    secondsSinceLastUpdate := float64(time.Now().Unix() - intLastUpdated)
    return secondsSinceLastUpdate
}
