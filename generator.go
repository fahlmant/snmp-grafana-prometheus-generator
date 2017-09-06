package main

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
)

//Top-level struct holding all configs
type PrometheusJobs struct {
	Global struct {
		EvaluationInterval string `json:"evaluation_interval"`
		ScrapeInterval     string `json:"scrape_interval"`
		ScrapeTimeout      string `json:"scrape_timeout"`
	} `json:"global"`
	PromScrapeConfigs ScrapeConfigs `json:"scrape_configs"`
	//Catch all for other data
	X map[string]interface{} `yaml:",inline"`
}

type ScrapeConfigs []ScrapeConfig

//Holds all scrape_config data
type ScrapeConfig struct {
	JobName     string `json:"job_name"`
	MetricsPath string `json:"metrics_path"`
	Params      struct {
		Module []string `json:"module"`
	} `json:"params"`
	RelabelConfigs []struct {
		SourceLabels []string `json:"source_labels,omitempty"`
		TargetLabel  string   `json:"target_label"`
		Replacement  string   `json:"replacement,omitempty"`
	} `json:"relabel_configs"`
	/*	StaticConfigs []struct {
		Targets []string `json:"targets"`
	} `json:"static_configs"`*/
	PromStaticConfigs StaticConfigs `json:"static_configs"`
	ScrapeInterval    string        `json:"scrape_interval,omitempty"`
	ScrapeTimeout     string        `json:"scrape_timeout,omitempty"`
}

//Holds all the targets to scrape for a given job
type StaticConfig struct {
	Targets []string `json:"targets"`
}

type StaticConfigs []StaticConfig

//Since a scrape_config can be an object or an array, must be custom-unmarshalled
func (sc *ScrapeConfigs) UnmarshalJSON(b []byte) error {

	var oneConfig ScrapeConfig
	if err := json.Unmarshal(b, &oneConfig); err == nil {
		*sc = ScrapeConfigs{oneConfig}
		return nil
	}
	var multiConfigs []ScrapeConfig
	if err := json.Unmarshal(b, &multiConfigs); err == nil {
		*sc = multiConfigs
		return nil
	} else {
		fmt.Println(err)
		fmt.Println("Badly formatted YAML: Exiting")
		os.Exit(1)
	}
	return nil
}

//See explination for scrape_config above
func (stc *StaticConfigs) UnmarshalJSON(b []byte) error {

	var oneConfig StaticConfig
	if err := json.Unmarshal(b, &oneConfig); err == nil {
		*stc = StaticConfigs{oneConfig}
		return nil
	}

	var multiConfigs []StaticConfig
	if err := json.Unmarshal(b, &multiConfigs); err == nil {
		*stc = multiConfigs
	} else {
		fmt.Println(err)
		fmt.Println("Badly formatted YAML: Exiting")
		os.Exit(1)

	}
	return nil
}

//Unroll struct and get input from user on what to do for each target in a given job
func PromptUser(res *PrometheusJobs) {

	for _, val := range res.PromScrapeConfigs {
		fmt.Println(val.JobName)

		for _, configs := range val.PromStaticConfigs {
			for _, targets := range configs.Targets {
				fmt.Printf("For job '%s', target %s what would you like to do\n", val.JobName, targets)
			}
		}
	}
}

func main() {

	//Ensure some file is provided
	//TODO: Validate file's existance and that it's valid YAML
	if len(os.Args) < 2 {

		fmt.Println("Please provide a file")
		os.Exit(1)
	}

	//Convert YAML to JSON for easier parsing
	data, _ := ioutil.ReadFile(os.Args[1])
	j, _ := yaml.YAMLToJSON(data)

	//The PrometheusJobs struct can't only handle 1 job for some reason
	//Unmarshal JSON into the proper Prometheus job structs
	res := PrometheusJobs{}
	json.Unmarshal(j, &res)
	PromptUser(&res)
}
