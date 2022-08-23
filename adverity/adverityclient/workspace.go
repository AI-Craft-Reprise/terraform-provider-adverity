package adverityclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (client *Client) ReadWorkspace(id string) (*Workspace, error, int) {
	u := *client.restURL
	u.Path = u.Path + "stacks/" + id + "/"

	response, err := client.sendRequestRead(u)
	if err != nil {
		return nil, err, 0
	}

	resMap := &Workspace{}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return resMap, errorString{"Failed reading workspace. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}, response.StatusCode
	}

	err = getJSON(response, resMap)
	if err != nil {
		return nil, err, response.StatusCode
	}

	return resMap, nil, response.StatusCode

}

func (client *Client) CreateWorkspace(conf CreateWorkspaceConfig) (*Workspace, error) {
	u := *client.restURL
	u.Path = u.Path + "stacks/"

	body, _ := json.Marshal(conf)
	response, err := client.sendRequestCreate(u, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	resMap := &Workspace{}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return resMap, errorString{"Failed creating workspace. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}

	err = getJSON(response, resMap)
	if err != nil {
		return nil, err
	}

	return resMap, nil
}

func (client *Client) UpdateWorkspace(conf UpdateWorkspaceConfig) (*http.Response, error) {
	u := *client.restURL
	u.Path = u.Path + "stacks/" + conf.StackSlug + "/"

	body, _ := json.Marshal(conf)
	response, err := client.sendRequestUpdate(u, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return response, errorString{"Failed deleting workspace. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}

	return response, nil

}

func (client *Client) DeleteWorkspace(conf DeleteWorkspaceConfig) (*http.Response, error) {
	u := *client.restURL
	u.Path = u.Path + "stacks/" + conf.StackSlug + "/"
	response, err := client.sendRequestDelete(u)
	if err != nil {
		return nil, err
	}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return response, errorString{"Failed deleting workspace. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}

	return response, nil

}
