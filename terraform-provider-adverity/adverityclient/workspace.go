package adverityclient

import (
	"io/ioutil"
	"net/http"
	"strconv"
    "encoding/json"
    "bytes"
    "log"
// 	"github.com/google/go-querystring/query"
)

func (client *Client) ReadWorkspace(id string) (*GetWorkspace, error) {
	u := *client.restURL
	u.Path = u.Path + "stacks/" + id + "/"

	resMap, err := client.sendRequestRead(u)
	if err != nil {
		return nil, err
	}
	return resMap, nil

}

func (client *Client) CreateWorkspace(conf CreateWorkspaceConfig) (*GetWorkspace, error) {
	u := *client.restURL
	u.Path = u.Path + "stacks/"

	body, _ := json.Marshal(conf)
	log.Println(string(body))
	resMap, err := client.sendRequestCreate(u, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	return resMap, nil
}

// func (client *Client) CreateWorkspace(conf CreateWorkspaceConfig) (*http.Response, error) {
// 	u := *client.restURL
// 	u.Path = u.Path + "stacks/"
//
// 	body, _ := json.Marshal(conf)
// 	log.Println(string(body))
// 	response, err := client.sendRequestCreate(u, bytes.NewReader(body))
// 	if err != nil {
// 		return nil, err
// 	}
// 	if !responseOK(response) {
// 		defer response.Body.Close()
// 		body, _ := ioutil.ReadAll(response.Body)
// 		return response, errorString{"Failed creating workspace. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
// 	}
//
// 	return response, nil
//
// }


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

