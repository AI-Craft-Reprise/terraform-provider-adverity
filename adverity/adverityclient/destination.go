package adverityclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (client *Client) ReadDestination(id string, destination_type_id int) (*Destination, error, int) {
	u := *client.restURL
	u.Path = u.Path + "target-types/" + strconv.Itoa(destination_type_id) + "/targets/" + id + "/"

	response, err := client.sendRequestRead(u)

	resMap := &Destination{}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return resMap, errorString{"Failed reading Destination. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}, response.StatusCode
	}

	err = getJSON(response, resMap)
	if err != nil {
		return nil, err, response.StatusCode
	}

	return resMap, nil, response.StatusCode

}

func (client *Client) CreateDestination(conf DestinationConfig, destination_type_id int) (*Destination, error) {
	u := *client.restURL
	u.Path = u.Path + "target-types/" + strconv.Itoa(destination_type_id) + "/targets/"

	body, _ := json.Marshal(conf)
	response, err := client.sendRequestCreate(u, bytes.NewReader(body))
	resMap := &Destination{}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return resMap, errorString{"Failed creating Destination. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}

	err = getJSON(response, resMap)
	if err != nil {
		return nil, err
	}

	return resMap, nil
}

func (client *Client) UpdateDestination(conf DestinationConfig, destination_type_id int, id string) (*http.Response, error) {
	u := *client.restURL

	u.Path = u.Path + "target-types/" + strconv.Itoa(destination_type_id) + "/targets/" + id + "/"

	body, _ := json.Marshal(conf)
	response, err := client.sendRequestUpdate(u, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return response, errorString{"Failed deleting Destination. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}

	return response, nil

}

func (client *Client) DeleteDestination(id string, destination_type_id int) (*http.Response, error) {
	u := *client.restURL

	u.Path = u.Path + "target-types/" + strconv.Itoa(destination_type_id) + "/targets/" + id + "/"

	response, err := client.sendRequestDelete(u)
	if err != nil {
		return nil, err
	}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return response, errorString{"Failed deleting Destination. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}

	return response, nil

}
