package adverityclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	//     "log"
	"fmt"
)

func (client *Client) ReadDatastream(id string, datastream_type_id int) (*Datastream, error) {
	u := *client.restURL
	u.Path = u.Path + "datastream-types/" + strconv.Itoa(datastream_type_id) + "/datastreams/" + id + "/"
	response, err := client.sendRequestRead(u)

	resMap := &Datastream{}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return resMap, errorString{"Failed reading datastream. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}

	err = getJSON(response, resMap)
	if err != nil {
		return nil, err
	}

	return resMap, nil

}

func (c *DatastreamConfig) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"name":      c.Name,
		"stack":     fmt.Sprintf("%d", c.Stack),
		"schedules": c.Schedules,
	}

	for _, param := range c.Parameters {
		m[param.Name] = param.Value
	}
	for _, p := range c.ParametersListInt {
		arr_int := []int{}
		for _, v := range p.Value {
			arr_int = append(arr_int, v)
		}
		m[p.Name] = arr_int
	}

	for _, p := range c.ParametersListStr {
		arr_str := []string{}
		for _, v := range p.Value {
			arr_str = append(arr_str, v)
		}
		m[p.Name] = arr_str
	}

	return json.Marshal(m)
}

func (client *Client) CreateDatastream(conf DatastreamConfig, datastream_type_id int) (*Datastream, error) {
	u := *client.restURL
	u.Path = u.Path + "datastream-types/" + strconv.Itoa(datastream_type_id) + "/datastreams/"

	body, _ := json.Marshal(&conf)
	response, err := client.sendRequestCreate(u, bytes.NewReader(body))
	resMap := &Datastream{}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return resMap, errorString{"Failed creating datastream. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}

	err = getJSON(response, resMap)
	if err != nil {
		return nil, err
	}

	return resMap, nil
}

func (client *Client) UpdateDatastream(conf DatastreamConfig, id string, datastream_type_id int) (*http.Response, error) {
	u := *client.restURL
	u.Path = u.Path + "datastream-types/" + strconv.Itoa(datastream_type_id) + "/datastreams/" + id + "/"

	body, _ := json.Marshal(&conf)
	response, err := client.sendRequestUpdate(u, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return response, errorString{"Failed updating workspace. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}

	return response, nil

}

func (client *Client) DeleteDatastream(id string, datastream_type_id int) (*http.Response, error) {
	u := *client.restURL
	u.Path = u.Path + "datastream-types/" + strconv.Itoa(datastream_type_id) + "/datastreams/" + id + "/"
	response, err := client.sendRequestDelete(u)
	if err != nil {
		return nil, err
	}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return response, errorString{"Failed deleting datastream. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}

	return response, nil

}

func (client *Client) EnableDatastream(conf DataStreamEnablingConfig, id string, datastream_type_id int) (*http.Response, error) {
	u := *client.restURL
	u.Path = u.Path + "datastream-types/" + strconv.Itoa(datastream_type_id) + "/datastreams/" + id + "/"

	body, _ := json.Marshal(&conf)
	response, err := client.sendRequestUpdate(u, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return response, errorString{"Failed enabling or disabling workspace. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}

	return response, nil
}
