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
	ApiPath      string `yaml:"api_path"`
}

type Client interface {
	doRequest(req *http.Request) ([]byte, error)
	create(body []byte) ([]byte, error)
}

func main() {
	//username := os.Getenv("SEVICENOW_USERNAME")
	//password := os.Getenv("SERVICENOW_PASSWORD")
	config := Config{}
	configYaml, _ := ioutil.ReadFile(configFile)
	_ = yaml.Unmarshal(configYaml, &config)

	snowConfig := config.ServiceNow
	snowClient, err := NewServiceNowClient(
		snowConfig.InstanceName,
		snowConfig.ApiPath,
		snowConfig.UserName,
		snowConfig.Password,
	)
	if err != nil {
		log.Errorf("could not create servicnowclient. %v", err)
	}

	server := CreateSnowServer(config, snowClient)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen to port 5000 %v", err)
	}

}
