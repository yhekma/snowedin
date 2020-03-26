package main

import (
	"bytes"
	"github.com/cloudflare/cfssl/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetParsedString(t *testing.T) {
	config := Config{}
	configYaml, _ := ioutil.ReadFile("config.yaml")
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

	testjson, _ := ioutil.ReadFile("test.json")

	t.Run("test if we get a string back", func(t *testing.T) {
		request := NewJsonPostRequest(testjson, "/webhook")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "bla: firing"

		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})
}

func NewJsonPostRequest(json []byte, url string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(json))
	return req
}
