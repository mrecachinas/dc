package api

import (
	"encoding/xml"
	"net/http"
)

// QueryExternal submits a GET request to the external
// XML API to get all available task options.
func QueryExternal(taskURL string, client *http.Client) ([]Task, error) {
	// Submit the HTTP request
	response, err := client.Get(taskURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Deserialize the XML response to a list of `Task`s
	// Note: Using xml.NewDecoder(...).Decode(...) allows
	// pulling from a stream rather than reading it all
	// into memory at once (as in xml.Unmarshal(...))
	var tasks []Task
	err = xml.NewDecoder(response.Body).Decode(&tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
