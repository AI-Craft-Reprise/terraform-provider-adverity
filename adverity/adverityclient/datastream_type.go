package adverityclient

import (
	"io/ioutil"
	"strconv"
)

func (client *Client) LookupDatastreamTypes(searchTerm string) ([]DatastreamType, error) {
	u := *client.restURL
	u.Path = u.Path + "datastream-types/"
	queries := []Query{
		{
			Key:   "search",
			Value: searchTerm,
		},
	}
	response, err := client.sendRequestQuery(u, queries)
	if err != nil {
		return nil, err
	}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return nil, errorString{"Failed querying datastream types. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}
	resMap := &DatastreamTypeResults{}
	err = getJSON(response, resMap)
	if err != nil {
		return nil, err
	}
	return resMap.Results, nil
}
