package jenkins_types

import (
	"github.com/puppetlabs/pipeline-dashboard/config"
)

type Link struct {
    URL   string
    Title string
}

func GetLinks() []Link {
    configData := config.GetConfig()

    var retVal []Link
    for _, link := range configData.Links {
        retVal = append(retVal, Link{
            URL:   link.URL,
            Title: link.Title,
        })
    }

    return retVal
}
