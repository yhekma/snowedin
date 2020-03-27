package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/common/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
)

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
	flagUsername := flag.String("username", "", "username for servicenow")
	flagPassword := flag.String("password", "", "password for servicenow")
	flagInstanceName := flag.String("instance", "", "name if the instance to use. <name>.service-now.com")
	listenPort := flag.String("port", "5000", "port to listen on")
	configFile := flag.String("config", "config.yaml", "configfile")
	flag.Parse()

	envUsername := os.Getenv("SERVICENOW_USERNAME")
	envPassword := os.Getenv("SERVICENOW_PASSWORD")

	config := Config{}
	configYaml, err := ioutil.ReadFile(*configFile)
	err = yaml.Unmarshal(configYaml, &config)
	if err != nil {
		log.Errorf("Could not read configfile %s. %v", configYaml, err)
	}

	snowConfig := config.ServiceNow
	var (
		username     string
		password     string
		instanceName string
	)

	if err != nil {
		log.Errorf("Could not parse configfile %s. %v", configYaml, err)
	}

	switch {
	case *flagUsername == "":
		username = envUsername
	default:
		username = *flagUsername
	}
	switch {
	case *flagPassword == "":
		password = envPassword
	default:
		password = envPassword
	}
	switch {
	case *flagInstanceName == "":
		instanceName = snowConfig.InstanceName
	default:
		instanceName = *flagInstanceName
	}

	snowClient, err := NewServiceNowClient(
		instanceName,
		snowConfig.ApiPath,
		username,
		password,
	)
	if err != nil {
		log.Errorf("could not create servicnowclient. %v", err)
	}

	server := CreateSnowServer(config, snowClient)

	if err := http.ListenAndServe(fmt.Sprintf(":"+*listenPort), server); err != nil {
		log.Fatalf("could not listen to port %s %v", *listenPort, err)
	}

}
