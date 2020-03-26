package main

import (
	"bytes"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubClient struct{}

func (s StubClient) doRequest(req *http.Request) ([]byte, error) {
	return []byte("200"), nil // TODO: Set this to something we expect
}

func (s StubClient) create(body []byte) ([]byte, error) {
	req, _ := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(body))
	return s.doRequest(req)
}

func TestGetParsedString(t *testing.T) {
	config := Config{}
	configYaml, _ := ioutil.ReadFile("config.yaml")
	_ = yaml.Unmarshal(configYaml, &config)

	server := CreateSnowServer(config, StubClient{})

	testjson, _ := ioutil.ReadFile("test.json")

	t.Run("test if we get a string back", func(t *testing.T) {
		request := NewJsonPostRequest(testjson, "/webhook")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "200"

		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})
}

func NewJsonPostRequest(json []byte, url string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(json))
	return req
}
