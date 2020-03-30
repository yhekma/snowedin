package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/prometheus/alertmanager/template"
	"net/http"
	"strconv"
	tmpltext "text/template"
	"time"
)

type snowServer struct {
	defaultIncident  map[string]string
	serviceNowClient Client
}

func CreateSnowServer(config Config, snowClient Client) *snowServer {
	return &snowServer{
		defaultIncident:  config.DefaultIncident,
		serviceNowClient: snowClient,
	}
}

func (s *snowServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		return
	}

	data, _ := readRequestBody(r)
	fmt.Printf("---------\nProcessessing:\n%v\n----------", data)
	incident := Incident{}
	for _, alert := range data.Alerts {
		for k, v := range s.defaultIncident {
			parsedText, _ := applyTemplate(v, alert)
			incident[k] = parsedText
		}
	}
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	incident["u_correlation_id"] = timestamp
	b, _ := json.Marshal(incident)
	resp, _ := s.serviceNowClient.create(b)
	_, _ = fmt.Fprintf(w, string(resp))
}

func readRequestBody(r *http.Request) (template.Data, error) {
	defer r.Body.Close()

	data := template.Data{}
	err := json.NewDecoder(r.Body).Decode(&data)

	return data, err
}

func applyTemplate(text string, data template.Alert) (string, error) {
	tmpl, err := tmpltext.New("n").Parse(text)
	if err != nil {
		return "", err
	}

	var result bytes.Buffer
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}
