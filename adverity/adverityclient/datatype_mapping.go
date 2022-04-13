package adverityclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strconv"
)

func (client *Client) ReadColumns(datastreamID string) ([]Column, error) {
	u := *client.restURL
	u.Path = u.Path + "columns/"
	page := 1
	queries := []Query{
		{
			Key:   "datastream_id",
			Value: datastreamID,
		},
		{
			Key:   "page",
			Value: strconv.Itoa(page),
		},
	}

	columns := []Column{}
	for {
		response, err := client.sendRequestQuery(u, queries)
		resultsMap := &ColumnResults{}
		if !responseOK(response) {
			defer response.Body.Close()
			body, _ := ioutil.ReadAll(response.Body)
			return nil, errorString{"Failed reading column mapping. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
		}
		err = getJSON(response, resultsMap)
		if err != nil {
			return nil, err
		}
		for _, column := range resultsMap.Results {
			columns = append(columns, column)
		}
		if resultsMap.Next == "" {
			break
		} else {
			page = page + 1
			queries[1] = Query{
				Key:   "page",
				Value: strconv.Itoa(page),
			}
		}
	}
	return columns, nil
}

func (client *Client) PatchColumn(columnID string, dataType string) error {
	u := *client.restURL
	u.Path = u.Path + "columns/" + columnID + "/"
	datatypeMap := map[string]string{"datatype": dataType}
	body, _ := json.Marshal(datatypeMap)
	response, err := client.sendRequestUpdate(u, bytes.NewReader(body))
	if err != nil {
		return err
	}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return errorString{"Failed patching column. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}
	return nil
}
