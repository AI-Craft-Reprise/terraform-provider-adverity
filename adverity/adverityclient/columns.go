package adverityclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
)

func (client *Client) CreateColumns(datastreamID string, columns []ColumnConfig) ([]Column, error) {
	u := *client.restURL
	u.Path = u.Path + "datastreams/" + datastreamID + "/columns/"
	body, err := json.Marshal(columns)
	if err != nil {
		return nil, err
	}
	response, err := client.sendRequestCreate(u, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Trying to create columns by sending request to " + u.Path)
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return nil, errorString{"Failed creating columns. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}
	resultColumns := &[]Column{}
	err = getJSON(response, resultColumns)
	if err != nil {
		return nil, err
	}
	return *resultColumns, nil
}
