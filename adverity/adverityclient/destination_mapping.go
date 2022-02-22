package adverityclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (client *Client) ReadDestinationMapping(id int, destination_type int, destination_id int) (*DestinationMapping, error, int) {
	u := *client.restURL
	u.Path = u.Path + "target-types/" + strconv.Itoa(destination_type) + "/targets/" + strconv.Itoa(destination_id) + "/mappings/" + strconv.Itoa(id) + "/"

	response, err := client.sendRequestRead(u)

	resMap := &DestinationMapping{}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return resMap, errorString{"Failed reading Destination Mapping. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}, response.StatusCode
	}

	err = getJSON(response, resMap)
	if err != nil {
		return nil, err, response.StatusCode
	}

	return resMap, nil, response.StatusCode
}

func (client *Client) CreateDestinationMapping(conf DestinationMappingConfig, destination_type int, destination_id int) (*DestinationMapping, error) {
	u := *client.restURL
	u.Path = u.Path + "target-types/" + strconv.Itoa(destination_type) + "/targets/" + strconv.Itoa(destination_id) + "/mappings/"

	body, _ := json.Marshal(conf)
	response, err := client.sendRequestCreate(u, bytes.NewReader(body))
	resMap := &DestinationMapping{}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return resMap, errorString{"Failed creating Destination Mapping. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}

	err = getJSON(response, resMap)
	if err != nil {
		return nil, err
	}

	return resMap, nil
}

func (client *Client) UpdateDestinationMapping(conf DestinationMappingConfig, destination_type int, destination_id int, id int) (*DestinationMapping, error) {
	u := *client.restURL
	u.Path = u.Path + "target-types/" + strconv.Itoa(destination_type) + "/targets/" + strconv.Itoa(destination_id) + "/mappings/" + strconv.Itoa(id) + "/"

	body, _ := json.Marshal(conf)
	response, err := client.sendRequestUpdate(u, bytes.NewReader(body))
	resMap := &DestinationMapping{}
	if err != nil {
		return nil, err
	}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return resMap, errorString{"Failed updating Destination Mapping. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}

	err = getJSON(response, resMap)
	if err != nil {
		return nil, err
	}

	return resMap, nil
}

func (client *Client) DeleteDestinationMapping(id int, destination_type int, destination_id int) (*http.Response, error) {
	u := *client.restURL
	u.Path = u.Path + "target-types/" + strconv.Itoa(destination_type) + "/targets/" + strconv.Itoa(destination_id) + "/mappings/" + strconv.Itoa(id) + "/"

	response, err := client.sendRequestDelete(u)
	if err != nil {
		return nil, err
	}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return response, errorString{"Failed deleting Destination Mapping. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}

	return response, nil
}
