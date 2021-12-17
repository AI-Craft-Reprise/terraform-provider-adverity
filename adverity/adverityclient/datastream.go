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
		"schedules": c.Schedules,
	}

	if c.Stack != 0 {
		m["stack"] = fmt.Sprintf("%d", c.Stack)
	}
	if c.Auth != 0 {
		m["auth"] = fmt.Sprintf("%d", c.Auth)
	}
	if c.Datatype != "" {
		m["datatype"] = c.Datatype
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

func (client *Client) UpdateDatastream(conf DatastreamConfig, id string) (*http.Response, error) {
	u := *client.restURL
	u.Path = u.Path + "datastreams/" + id + "/"

	body, _ := json.Marshal(&conf)

	response, err := client.sendRequestUpdate(u, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return response, errorString{"Failed updating datastream. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
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

func (client *Client) EnableDatastream(conf DataStreamEnablingConfig, id string) (*http.Response, error) {
	u := *client.restURL
	u.Path = u.Path + "datastreams/" + id + "/"

	body, _ := json.Marshal(&conf)
	response, err := client.sendRequestUpdate(u, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return response, errorString{"Failed enabling or disabling datastream. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}

	return response, nil
}

func (client *Client) DataStreamChanegDatatype(conf DatastreamDatatypeConfig, id string) (*http.Response, error) {
	u := *client.restURL
	u.Path = u.Path + "datastreams/" + id + "/"

	body, _ := json.Marshal(&conf)
	response, err := client.sendRequestUpdate(u, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return response, errorString{"Failed setting datatype of Datastream. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}

	return response, nil
}
