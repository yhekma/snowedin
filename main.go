package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
)

const listenPort = "5000"

var log = logrus.New()

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
	log.SetFormatter(&logrus.JSONFormatter{})

	configFile := flag.String("config", "config.yaml", "configfile")
	debug := flag.Bool("debug", false, "run in debug mode")
	flag.Parse()

	log.Infof("Using '%s' as config yaml", *configFile)

	if *debug {
		log.SetLevel(logrus.DebugLevel)
		log.Debug("Running in debug mode")
	}

	envUsername := os.Getenv("SERVICENOW_USERNAME")
	envPassword := os.Getenv("SERVICENOW_PASSWORD")
	envInstanceName := os.Getenv("SERVICENOW_INSTANCE_NAME")

	config := Config{}
	configYaml, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("Could not read configfile %s. %v", configYaml, err)
	}

	err = yaml.Unmarshal(configYaml, &config)

	if err != nil {
		log.Fatalf("Could not parse configfile %s. %v", configYaml, err)
	}

	snowConfig := config.ServiceNow
	var (
		username     string
		password     string
		instanceName string
	)

	switch {
	case envUsername == "":
		username = snowConfig.UserName
		log.Debug("Using username from config")
	default:
		username = envUsername
		log.Debug("Using username from env")
	}
	switch {
	case envPassword == "":
		password = snowConfig.Password
		log.Debug("Using password from config")
	default:
		password = envPassword
		log.Debug("Using password from env")
	}
	switch {
	case envInstanceName == "":
		instanceName = snowConfig.InstanceName
		log.Debug("Using instancename from config")
	default:
		instanceName = envInstanceName
		log.Debug("Using instancename from env")
	}

	snowClient, err := NewServiceNowClient(
		instanceName,
		snowConfig.ApiPath,
		username,
		password,
		log,
	)
	if err != nil {
		log.Fatalf("could not create servicnowclient. %v", err)
	}

	server := CreateSnowServer(config, snowClient, log)

	if err := http.ListenAndServe(fmt.Sprintf(":"+listenPort), server); err != nil {
		log.Fatalf("could not listen to port %s %v", listenPort, err)
	}

}
