package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/prometheus/common/log"
	"io/ioutil"
	"net/http"
)

const (
	servicenowBaseUrl = "https://%s.service-now.com"
	servicenowAPIPath = "%s"
)

type Incident map[string]string

type ServiceNowClient struct {
	baseURL    string
	apiPath    string
	authHeader string
	client     *http.Client
}

func NewServiceNowClient(instanceName, apiPath, userName, password string) (*ServiceNowClient, error) {
	if instanceName == "" {
		return nil, errors.New("no instancename specified")
	}

	if userName == "" {
		return nil, errors.New("no username specified")
	}

	if password == "" {
		return nil, errors.New("no password specified")
	}

	return &ServiceNowClient{
		baseURL:    fmt.Sprintf(servicenowBaseUrl, instanceName),
		apiPath:    apiPath,
		authHeader: fmt.Sprintf("Basic %s", base64.URLEncoding.EncodeToString([]byte(userName+":"+password))),
		client:     http.DefaultClient,
	}, nil
}

func (snClient *ServiceNowClient) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", snClient.authHeader)
	resp, err := snClient.client.Do(req)

	if err != nil {
		log.Errorf("Error sending the request. %v", err)
		return nil, err
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading response body. %s", err)
		return nil, err
	}

	return responseBody, nil

}

func (snClient *ServiceNowClient) create(body []byte) ([]byte, error) {
	url := fmt.Sprint(snClient.baseURL, snClient.apiPath)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Errorf("Error creating request. %s", err)
		return nil, err
	}

	return snClient.doRequest(req)
}
