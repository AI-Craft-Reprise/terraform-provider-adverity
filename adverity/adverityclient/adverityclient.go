package adverityclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	// 	"log"
	"io/ioutil"
)

// CreateClientFromLogin can be used to create an alooma client from an API key
func CreateClientFromLogin(instanceURL, token string) (*Client, error) {
	var client Client
	client.token = token
	baseURL := instanceURL
	restPath := "/api/"
	restURL, err := url.Parse(baseURL + restPath)

	if err != nil {
		return nil, err
	}
	client.restURL = restURL
	params := make(map[string]string)
	params["timeout"] = "60"
	client.requestsParams = params
	opt := &cookiejar.Options{}
	jar, _ := cookiejar.New(opt)
	client.httpClient = &http.Client{Jar: jar, Timeout: time.Second * 60}
	return &client, nil
}

func responseOK(response *http.Response) bool {
	return response.StatusCode >= 200 && response.StatusCode < 300
}

func getJSON(r *http.Response, target interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func (client *Client) sendRequestUpdate(u url.URL, body *bytes.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPatch, u.String(), ioutil.NopCloser(body))
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", client.token))
	req.Header.Add("Content-Type", "application/json")

	response, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return response, err
}

func (client *Client) sendRequestCreate(u url.URL, body *bytes.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, u.String(), ioutil.NopCloser(body))
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", client.token))
	req.Header.Add("Content-Type", "application/json")

	response, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return response, err
}

func (client *Client) sendRequestDelete(u url.URL) (*http.Response, error) {
	// log.Println(u.String())
	req, err := http.NewRequest(http.MethodDelete, u.String(), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", client.token))
	req.Header.Add("Content-Type", "application/json")

	response, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return response, err
}

func (client *Client) sendRequestRead(u url.URL) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", client.token))
	req.Header.Add("Content-Type", "application/json")

	response, err := client.httpClient.Do(req)

	if err != nil {
		return nil, err
	}
	return response, err
}

func (client *Client) sendRequestQuery(u url.URL, queries []Query) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", client.token))

	q := req.URL.Query()

	for _, query := range queries {
		q.Add(query.Key, query.Value)
	}

	req.URL.RawQuery = q.Encode()

	response, err := client.httpClient.Do(req)

	if err != nil {
		return nil, err
	}
	return response, err
}

func (client *Client) sendRequestOptions(u url.URL) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodOptions, u.String(), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", client.token))
	req.Header.Add("Content-Type", "application/json")

	response, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return response, err
}
