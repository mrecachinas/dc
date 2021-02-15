package api

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

// QueryExternal submits a GET request to the external
// API to get all available task options.
func QueryExternal(taskURL string) ([]Task, error) {
	resp, err := http.Get(taskURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tasks []Task
	err = xml.Unmarshal(body, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
