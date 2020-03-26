package snowedin

import (
	"bytes"
	"io/ioutil"

	//"github.com/prometheus/alertmanager/template"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetParsedString(t *testing.T) {
	server := &SnowServer{}
	server.FieldConfig = map[string]string{"bla": "{{.Status}}"}
	//testjson, _ := json.Marshal(testmap)
	testjson, _ := ioutil.ReadFile("test.json")

	t.Run("test if we get a string back", func(t *testing.T) {
		request := NewJsonPostRequest(testjson, "/webhook")
		response := httptest.NewRecorder()

		server.ServerHTTP(response, request)

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
