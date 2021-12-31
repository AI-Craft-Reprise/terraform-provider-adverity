package adverityclient

import (
	"io/ioutil"
	"strconv"
	"strings"
)

func (client *Client) DoLookup(url string, queries []Query) (*Lookup, error) {
	u := *client.restURL
	u.Path = strings.ReplaceAll(u.Path, "api/", url)

	response, err := client.sendRequestQuery(u, queries)

	resMap := &Lookup{}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return resMap, errorString{"Failed doing lookup. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}

	err = getJSON(response, resMap)
	if err != nil {
		return nil, err
	}

	return resMap, nil
}
