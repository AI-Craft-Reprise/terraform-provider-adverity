package adverityclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (client *Client) ReadDatastream(id string, datastream_type_id int) (*Datastream, error, int) {
	u := *client.restURL
	u.Path = u.Path + "datastream-types/" + strconv.Itoa(datastream_type_id) + "/datastreams/" + id + "/"
	response, err := client.sendRequestRead(u)

	resMap := &Datastream{}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return resMap, errorString{"Failed reading datastream. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}, response.StatusCode
	}

	err = getJSON(response, resMap)
	if err != nil {
		return nil, err, response.StatusCode
	}

	return resMap, nil, response.StatusCode

}

func (c *DatastreamConfig) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"name": c.Name,
	}
	if c.Schedules != nil {
		m["schedules"] = c.Schedules
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
	if c.Description != nil {
		m["description"] = *c.Description
	}
	if c.RetentionType != nil {
		m["retention_type"] = *c.RetentionType
	}
	if c.RetentionNumber != nil {
		m["retention_number"] = *c.RetentionNumber
	}
	if c.OverwriteKeyColumns != nil {
		m["overwrite_key_columns"] = *c.OverwriteKeyColumns
	}
	if c.OverwriteDatastream != nil {
		m["overwrite_datastream"] = *c.OverwriteDatastream
	}
	if c.OverwriteFileName != nil {
		m["overwrite_filename"] = *c.OverwriteFileName
	}
	if c.IsInsightsMediaplan != nil {
		m["is_insights_mediaplan"] = *c.IsInsightsMediaplan
	}
	if c.ManageExtractNames != nil {
		m["manage_extract_names"] = *c.ManageExtractNames
	}
	if c.ExtractNameKeys != nil {
		m["extract_name_keys"] = *c.ExtractNameKeys
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

func (c *DatastreamSpecificConfig) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}
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
	log.Println("[DEBUG] Sending body for create: " + string(body))
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

// Deprecated
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

func (client *Client) UpdateDatastreamCommon(conf DatastreamCommonUpdateConfig, id string) (*http.Response, error) {
	u := *client.restURL
	u.Path = u.Path + "datastreams/" + id + "/"
	body, _ := json.Marshal(&conf)
	log.Println("[DEBUG] Sent body for common config: " + string(body))
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

func (client *Client) UpdateDatastreamSpecific(conf DatastreamSpecificConfig, id string, datastream_type_id int) (*http.Response, error) {
	u := *client.restURL
	u.Path = u.Path + "datastream-types/" + strconv.Itoa(datastream_type_id) + "/datastreams/" + id + "/"
	body, _ := json.Marshal(&conf)
	log.Println("[DEBUG] Sent body for specific config: " + string(body))
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

func (client *Client) DataStreamChangeDatatype(conf DatastreamDatatypeConfig, id string) (*http.Response, error) {
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

func (client *Client) ScheduleFetch(days_to_fetch int, id string) (*http.Response, error) {
	u := *client.restURL
	u.Path = u.Path + "datastreams/" + id + "/fetch_fixed/"

	currentTime := time.Now()
	endDate := currentTime.Format("2006-01-02")
	startDate := currentTime.AddDate(0, 0, -days_to_fetch).Format("2006-01-02")

	fetchConf := FetchConfig{
		StartDate: startDate,
		EndDate:   endDate,
	}

	body, _ := json.Marshal(fetchConf)
	log.Println("[DEBUG] Sent body for echeduling fetch: " + string(body))
	response, err := client.sendRequestCreate(u, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return response, errorString{"Failed doing fetch. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}
	return response, nil
}
