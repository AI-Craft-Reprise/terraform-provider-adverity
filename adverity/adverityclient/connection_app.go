package adverityclient

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func (client *Client) LookupConnectionApp(connectionTypeID int, selector string) (int, error) {
	u := *client.restURL
	u.Path = u.Path + "connection-types/" + strconv.Itoa(connectionTypeID) + "/connections/"
	response, err := client.sendRequestOptions(u)
	if err != nil {
		return -1, err
	}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return -1, errorString{"Failed querying connection apps. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}
	resMap := &ConnectionOptions{}
	err = getJSON(response, resMap)
	if err != nil {
		return -1, err
	}
	if len(resMap.Actions["POST"].App.Choices) > 1 {
		stringList := []string{}
		for _, choice := range resMap.Actions["POST"].App.Choices {
			if choice.DisplayName == selector {
				return choice.Value, nil
			}
			stringList = append(stringList, choice.DisplayName)
		}
		return -1, errorString{fmt.Sprintf("Multiple app options found for connection type, none matched the selector (or no selector was given): %s", strings.Join(stringList, ", "))}
	} else if len(resMap.Actions["POST"].App.Choices) < 1 {
		return -1, errorString{"No app options found for connection type."}
	}
	return resMap.Actions["POST"].App.Choices[0].Value, nil
}
