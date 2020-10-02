package config

import (
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type JobDefinition struct {
    URL     string
    Product string
    Version string
}

type ProductConfig struct {
    Name     string
    Pipeline string
}

// Config structure that we should probably viperise in the future
type Config struct {
    UseCache       bool            `toml:"use_cache"`
    CithURL        string          `toml:"cith_url"`
    ScrapeInterval uint64          `toml:"scrape_interval"`
    Products       []ProductConfig `toml:"products"`
    KickoffJobs    []JobDefinition `toml:"kickoff_jobs"`
    OrderedJobs    []JobDefinition `toml:"ordered_jobs"`
}

func GetConfig() Config {

    // This is an interim solution for configuration parameters.
    // Can clearly do better but want to get this out for use asap.
    tomlData, err := ioutil.ReadFile("conf/config.toml")

    if err != nil {
        panic(err)
    }

    var conf Config

    if _, err := toml.Decode(string(tomlData), &conf); err != nil {
        fmt.Println(err)
    }

    // Hardwire this one off for the moment until a later revision.
    conf.UseCache = false

    return conf
}
