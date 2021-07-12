package adverityclient

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"encoding/json"
	"time"
	"bytes"
	"fmt"
	"log"
	"io/ioutil"
	"strconv"
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
	return response, err
}


func (client *Client) sendRequestCreate(u url.URL, body *bytes.Reader) (*GetWorkspace, error) {
    req, err := http.NewRequest("POST", u.String(), ioutil.NopCloser(body))
    req.Header.Add("Authorization", fmt.Sprintf("Token %s", client.token))
    req.Header.Add("Content-Type", "application/json")

	response, err := client.httpClient.Do(req)
    resMap:= &GetWorkspace{}
    if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return resMap, errorString{"Failed reading workspace. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}


	err = getJSON(response, resMap)
	if err != nil {
	    return nil, err
    }

	return resMap, err
}


func (client *Client) sendRequestDelete(u url.URL) (*http.Response, error) {
    log.Println(u.String())
    req, err := http.NewRequest("DELETE", u.String(), nil)
    req.Header.Add("Authorization", fmt.Sprintf("Token %s", client.token))
    req.Header.Add("Content-Type", "application/json")

	response, err := client.httpClient.Do(req)
	return response, err
}

func (client *Client) sendRequestRead(u url.URL) (*GetWorkspace, error) {
    req, err := http.NewRequest("GET", u.String(), nil)
    req.Header.Add("Authorization", fmt.Sprintf("Token %s", client.token))
    req.Header.Add("Content-Type", "application/json")

	response, err := client.httpClient.Do(req)
    resMap:= &GetWorkspace{}
    if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return resMap, errorString{"Failed reading workspace. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}


	err = getJSON(response, resMap)
	if err != nil {
	    return nil, err
    }

	return resMap, err
}
