package adverityclient

import (
	"io/ioutil"
	"strconv"
)

func (client *Client) ReadAuthUrl(connectionTypeId int, connectionId int) (*AuthUrl, error) {
	u := *client.restURL
	u.Path = u.Path + "connection-types" + strconv.Itoa(connectionTypeId) + "/connections/" + strconv.Itoa(connectionId) + "/authorize"

	response, err := client.sendRequestRead(u)

	resMap := &AuthUrl{}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return resMap, errorString{"Failed reading authentication URL. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}

	err = getJSON(response, resMap)
	if err != nil {
		return nil, err
	}

	return resMap, nil
}
