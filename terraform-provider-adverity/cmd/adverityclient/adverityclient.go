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
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return true
	}
	return false
}

func getJSON(r *http.Response, target interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}


func (client *Client) sendRequestUpdate(u url.URL, body *bytes.Reader) (*http.Response, error) {
    req, err := http.NewRequest("PATCH", u.String(), ioutil.NopCloser(body))
    req.Header.Add("Authorization", fmt.Sprintf("Token %s", client.token))
    req.Header.Add("Content-Type", "application/json")

	response, err := client.httpClient.Do(req)
	if err != nil {
	    return nil, err
    }
	return response, err
}


func (client *Client) sendRequestCreate(u url.URL, body *bytes.Reader) (*http.Response, error) {
    req, err := http.NewRequest("POST", u.String(), ioutil.NopCloser(body))
    req.Header.Add("Authorization", fmt.Sprintf("Token %s", client.token))
    req.Header.Add("Content-Type", "application/json")

	response, err := client.httpClient.Do(req)
    if err != nil {
	    return nil, err
    }
	return response, err
}


func (client *Client) sendRequestDelete(u url.URL) (*http.Response, error) {
//     log.Println(u.String())
    req, err := http.NewRequest("DELETE", u.String(), nil)
    req.Header.Add("Authorization", fmt.Sprintf("Token %s", client.token))
    req.Header.Add("Content-Type", "application/json")

	response, err := client.httpClient.Do(req)
	if err != nil {
	    return nil, err
    }
	return response, err
}

func (client *Client) sendRequestRead(u url.URL) (*http.Response, error) {
    req, err := http.NewRequest("GET", u.String(), nil)
    req.Header.Add("Authorization", fmt.Sprintf("Token %s", client.token))
    req.Header.Add("Content-Type", "application/json")

	response, err := client.httpClient.Do(req)

	if err != nil{
	    return nil, err
	}
	return response, err
}
