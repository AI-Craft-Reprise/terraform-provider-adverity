package adverityclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func (client *Client) ReadStorage(id string) (*Storage, error) {
	u := *client.restURL
	u.Path = u.Path + "storage/" + id + "/"
	response, err := client.sendRequestRead(u)

	resMap := &Storage{}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return resMap, errorString{"Failed reading Storage. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}

	err = getJSON(response, resMap)
	if err != nil {
		return nil, err
	}

	return resMap, nil

}

func (client *Client) CreateStorage(conf StorageConfig) (*Storage, error) {
	u := *client.restURL
	u.Path = u.Path + "storage/"
	body, _ := json.Marshal(conf)

	log.Println(string(body))
	response, err := client.sendRequestCreate(u, bytes.NewReader(body))
	resMap := &Storage{}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return resMap, errorString{"Failed creating Storage. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}

	err = getJSON(response, resMap)
	if err != nil {
		return nil, err
	}

	return resMap, nil
}

func (client *Client) UpdateStorage(conf StorageConfig, id string) (*http.Response, error) {
	u := *client.restURL
	u.Path = u.Path + "storage/" + id + "/"

	body, _ := json.Marshal(&conf)
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

func (client *Client) DeleteStorage(id string) (*http.Response, error) {
	u := *client.restURL
	u.Path = u.Path + "storage/" + id + "/"
	response, err := client.sendRequestDelete(u)
	if err != nil {
		return nil, err
	}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return response, errorString{"Failed deleting storage. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}

	return response, nil

}
