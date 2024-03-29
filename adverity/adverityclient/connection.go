package adverityclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	// "log"
	"net/http"
	"strconv"
)

func (client *Client) ReadConnection(id string, connection_type_id int) (*Connection, error, int) {
	u := *client.restURL
	u.Path = u.Path + "connection-types/" + strconv.Itoa(connection_type_id) + "/connections/" + id + "/"
	response, err := client.sendRequestRead(u)
	if err != nil {
		return nil, err, 0
	}

	resMap := &Connection{}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return resMap, errorString{"Failed reading connection. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}, response.StatusCode
	}

	err = getJSON(response, resMap)
	if err != nil {
		return nil, err, response.StatusCode
	}

	return resMap, nil, response.StatusCode

}

func (c *ConnectionConfig) MarshalJSON() ([]byte, error) {
	m := map[string]string{
		"name":  c.Name,
		"stack": fmt.Sprintf("%d", c.Stack),
	}
	for _, param := range c.Parameters {
		m[param.Name] = param.Value
	}
	return json.Marshal(m)
}

func (client *Client) CreateConnection(conf ConnectionConfig, connection_type_id int) (*Connection, error) {
	u := *client.restURL
	u.Path = u.Path + "connection-types/" + strconv.Itoa(connection_type_id) + "/connections/"

	body, _ := json.Marshal(&conf)
	response, err := client.sendRequestCreate(u, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	resMap := &Connection{}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return resMap, errorString{"Failed creating connection. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}

	err = getJSON(response, resMap)
	if err != nil {
		return nil, err
	}

	return resMap, nil
}

func (client *Client) UpdateConnection(conf ConnectionConfig, id string, connection_type_id int) (*http.Response, error) {
	u := *client.restURL
	u.Path = u.Path + "connection-types/" + strconv.Itoa(connection_type_id) + "/connections/" + id + "/"

	body, _ := json.Marshal(&conf)
	response, err := client.sendRequestUpdate(u, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return response, errorString{"Failed updating connection. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}

	return response, nil

}

func (client *Client) DeleteConnection(id string, connection_type_id int) (*http.Response, error) {
	u := *client.restURL
	u.Path = u.Path + "connection-types/" + strconv.Itoa(connection_type_id) + "/connections/" + id + "/"
	response, err := client.sendRequestDelete(u)
	if err != nil {
		return nil, err
	}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return response, errorString{"Failed deleting connection. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}

	return response, nil

}
