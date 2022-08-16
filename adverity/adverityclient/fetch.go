package adverityclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
	"time"
	"strings"
)

func (client *Client) DoFetch(fetchConfig FetchConfig, id string) (*FetchResponse, error) {
	u := *client.restURL
	u.Path = u.Path + "datastreams/" + id + "/fetch_fixed/"
	body, _ := json.Marshal(fetchConfig)
	log.Println("[DEBUG] Sent body for scheduling fetch: " + string(body))
	response, err := client.sendRequestCreate(u, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return nil, errorString{"Failed doing fetch. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}
	}
	resMap := &FetchResponse{}
	err = getJSON(response, resMap)
	return resMap, nil
}

func (client *Client) FetchNumberOfDays(days_to_fetch int, id string) (*FetchResponse, error) {
	if days_to_fetch < 0 {
		return nil, errorString{"Days to fetch cannot be negative."}
	}
	currentTime := time.Now()
	endDate := currentTime.Format("2006-01-02")
	startDate := currentTime.AddDate(0, 0, -days_to_fetch).Format("2006-01-02")

	fetchConf := FetchConfig{
		StartDate: startDate,
		EndDate:   endDate,
	}

	return client.DoFetch(fetchConf, id)
}

func (client *Client) FetchOnDate(startDate string, endDate string, id string) (*FetchResponse, error) {
	if len(strings.TrimSpace(startDate)) == 0 ||  len(strings.TrimSpace(endDate)) == 0{
		return nil, errorString{"Given dates are empty."}
	}
	fetchConf := FetchConfig{
		StartDate: startDate,
		EndDate:   endDate,
	}
	return client.DoFetch(fetchConf, id)
}

func (client *Client) FetchPreviousMonths(days_to_fetch int, id string) (*FetchResponse, error) {
	currentTime := time.Now()

	// Take the last day of the previous month
	endDate := LastOfMonth(currentTime.AddDate(0, -1, 0))
	// Take the first day of the month after subtraction of the amount of days to fetch
	startDate := FirstOfMonth(currentTime.AddDate(0, 0, -days_to_fetch))
	// If the startdate is after the enddate (if after subtraction of days we're still in the current month), use the first day of the previous month
	if endDate.Before(startDate) {
		startDate = FirstOfMonth(currentTime.AddDate(0, -1, 0))
	}

	startDateText := startDate.Format("2006-01-02")
	endDateText := endDate.Format("2006-01-02")

	fetchConfig := FetchConfig{
		StartDate: startDateText,
		EndDate:   endDateText,
	}
	return client.DoFetch(fetchConfig, id)
}

func (client *Client) FetchCurrentMonth(id string) (*FetchResponse, error) {
	currentTime := time.Now()
	// End date is current day
	endDate := currentTime.Format("2006-01-02")
	// Start date is the first day of this month
	startDate := FirstOfMonth(currentTime).Format("2006-01-02")

	fetchConfig := FetchConfig{
		StartDate: startDate,
		EndDate:   endDate,
	}
	return client.DoFetch(fetchConfig, id)
}

func (client *Client) FetchPreviousWeeks(days_to_fetch int, id string) (*FetchResponse, error) {
	currentTime := time.Now()
	// End date is the sunday of the previous week
	endDate := GetSunday(currentTime).AddDate(0, 0, -7)
	// Start date is the monday of the week after subtraction of the amount of days to fetch
	startDate := GetMonday(currentTime.AddDate(0, 0, -days_to_fetch))
	// If the startdate is after the enddate (if after subtraction of days we're still in the current week), use the first day of the previous week
	if endDate.Before(startDate) {
		startDate = GetMonday(currentTime.AddDate(0, 0, -7))
	}

	startDateText := startDate.Format("2006-01-02")
	endDateText := endDate.Format("2006-01-02")

	fetchConfig := FetchConfig{
		StartDate: startDateText,
		EndDate:   endDateText,
	}
	return client.DoFetch(fetchConfig, id)
}

func (client *Client) FetchCurrentWeek(id string) (*FetchResponse, error) {
	currentTime := time.Now()
	// End date is current day
	endDate := currentTime.Format("2006-01-02")
	// Start date is the Monday of this week
	startDate := GetMonday(currentTime).Format("2006-01-02")

	fetchConfig := FetchConfig{
		StartDate: startDate,
		EndDate:   endDate,
	}
	return client.DoFetch(fetchConfig, id)
}

func FirstOfMonth(date time.Time) time.Time {
	year, month, _ := date.Date()
	location := date.Location()
	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, location)
	return firstOfMonth
}

func LastOfMonth(date time.Time) time.Time {
	firstOfMonth := FirstOfMonth(date)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	return lastOfMonth
}

func GetMonday(date time.Time) time.Time {
	difference := int(date.Weekday()) - 1
	if difference < 0 {
		difference = 6
	}
	monday := date.AddDate(0, 0, -difference)
	return monday
}

func GetSunday(date time.Time) time.Time {
	difference := (7 - int(date.Weekday())) % 7
	sunday := date.AddDate(0, 0, difference)
	return sunday
}

func (client *Client) ReadJob(ID int) (*Job, error, int) {
	u := *client.restURL
	u.Path = u.Path + "jobs/" + strconv.Itoa(ID) + "/"
	response, err := client.sendRequestRead(u)
	if err != nil {
		return nil, err, 0
	}
	resMap := &Job{}
	if !responseOK(response) {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return resMap, errorString{"Failed reading job. Got back statuscode: " + strconv.Itoa(response.StatusCode) + " with body: " + string(body)}, response.StatusCode
	}
	err = getJSON(response, resMap)
	if err != nil {
		return nil, err, 0
	}
	return resMap, nil, 0
}
