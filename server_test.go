package main

import (
	"bytes"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type StubClient struct{}

func (s StubClient) doRequest(req *http.Request) ([]byte, error) {
	body, _ := ioutil.ReadAll(req.Body)
	return body, nil
}

func (s StubClient) create(body []byte) ([]byte, error) {
	req, _ := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(body))
	return s.doRequest(req)
}

func TestParsing(t *testing.T) {
	config := Config{}
	configYaml, _ := ioutil.ReadFile("tests/config.yaml")
	_ = yaml.Unmarshal(configYaml, &config)

	server := CreateSnowServer(config, StubClient{})

	testjson, _ := ioutil.ReadFile("tests/test.json")
	testjsonwant, _ := ioutil.ReadFile("tests/test_want.json")

	t.Run("test if json parsing/templating works", func(t *testing.T) {
		request := NewJsonPostRequest(testjson, "/webhook")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var (
			gotJson  map[string]string
			wantJson map[string]string
		)

		err := json.Unmarshal(response.Body.Bytes(), &gotJson)
		if err != nil {
			t.Errorf("could not parse json from response. %v", err)
		}
		_ = json.Unmarshal(testjsonwant, &wantJson)

		if _, ok := gotJson["u_correlation_id"]; !ok {
			t.Errorf("didn't get u_correlation_id back from json %v", wantJson)
		}
		// Remove u_correlation_id from returned map, since it's epoch time
		delete(gotJson, "u_correlation_id")

		if !reflect.DeepEqual(gotJson, wantJson) {
			t.Errorf("got '%s' want '%s'", gotJson, wantJson)
		}
	})

	t.Run("see if we get 200 back on get", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		if response.Code != 200 {
			t.Errorf("invalid response. Want 200, got %d", response.Code)
		}
	})
}

func NewJsonPostRequest(json []byte, url string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(json))
	return req
}
