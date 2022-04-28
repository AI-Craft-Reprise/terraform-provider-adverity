package adverityclient

import (
	"io/ioutil"
	"strconv"
)

func (client *Client) LookupConnectionTypes(searchTerm string) ([]ConnectionType, error) {
	u := *client.restURL
	u.Path = u.Path + "connection-types/"
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
		return nil, errorString{"Failed querying connection types. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}
	resMap := &ConnectionTypeResults{}
	err = getJSON(response, resMap)
	if err != nil {
		return nil, err
	}
	return resMap.Results, nil
}
