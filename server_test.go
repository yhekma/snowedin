package snowedin

import (
	"bytes"
	"io/ioutil"

	//"github.com/prometheus/alertmanager/template"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETMap(t *testing.T) {
	server := &SnowServer{}
	server.FieldConfig = map[string]string{"bla": "{{.Status}}"}
	//testjson, _ := json.Marshal(testmap)
	testjson, _ := ioutil.ReadFile("test.json")

	t.Run("test if we get a map of a json back", func(t *testing.T) {
		request := NewJsonPostRequest(testjson, "/webhook")
		response := httptest.NewRecorder()

		server.ServerHTTP(response, request)
	})
}

func NewJsonPostRequest(json []byte, url string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(json))
	return req
}