package main

import (
	"github.com/prometheus/common/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
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
	ApiPath      string `yaml:"api_path"`
}

type Client interface {
	doRequest(req *http.Request) ([]byte, error)
	create(body []byte) ([]byte, error)
}

func main() {
	username := os.Getenv("SERVICENOW_USERNAME")
	password := os.Getenv("SERVICENOW_PASSWORD")
	config := Config{}
	configYaml, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Errorf("Could not read configfile %s. %v", configYaml, err)
	}

	err = yaml.Unmarshal(configYaml, &config)
	if err != nil {
		log.Errorf("Could not parse configfile %s. %v", configYaml, err)
	}

	snowConfig := config.ServiceNow
	if username == "" {
		username = snowConfig.UserName
	}
	if password == "" {
		password = snowConfig.Password
	}

	snowClient, err := NewServiceNowClient(
		snowConfig.InstanceName,
		snowConfig.ApiPath,
		username,
		password,
	)
	if err != nil {
		log.Errorf("could not create servicnowclient. %v", err)
	}

	server := CreateSnowServer(config, snowClient)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen to port 5000 %v", err)
	}

}
