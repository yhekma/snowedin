package snowedin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/prometheus/alertmanager/template"
	"net/http"
	tmpltext "text/template"
)

type FieldConfig map[string]string

type SnowServer struct {
	FieldConfig FieldConfig
}

func (s *SnowServer) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	data, _ := readRequestBody(r)
	for _, alert := range data.Alerts {
		for k, v := range s.FieldConfig {
			parsedText, _ := applyTemplate(v, alert)
			_, _ = fmt.Fprintf(w, "%s: %s", k, parsedText)
		}
	}
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
