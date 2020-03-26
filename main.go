package main

import (
	"github.com/prometheus/common/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
)

// TODO: make this an argument
const configFile = "config.yaml"

type Config struct {
	DefaultIncident map[string]string `yaml:"default_incident"`
	ServiceNow      SnowConfig        `yaml:"servicenow_config"`
}

type SnowConfig struct {
	InstanceName string `yaml:"instance_name"`
	UserName     string `yaml:"user_name"`
	Password     string `yaml:"password"`
}

func main() {
	//username := os.Getenv("SEVICENOW_USERNAME")
	//password := os.Getenv("SERVICENOW_PASSWORD")
	config := Config{}
	configYaml, _ := ioutil.ReadFile(configFile)
	_ = yaml.Unmarshal(configYaml, &config)
	server := CreateSnowServer(config)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen to port 5000 %v", err)
	}

}
