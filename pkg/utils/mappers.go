package utils

import (
	"encoding/json"
	"fmt"
)


func TaskMapper(values map[string]interface{}) (any, error) {
	// Check if the "data" key exists in the map
	if _, ok := values["data"]; !ok {
		return nil, fmt.Errorf("key 'data' not found in the map")
	}
	// Check if the value associated with "data" is a string
	if _, ok := values["data"].(string); !ok {
		return nil, fmt.Errorf("value of 'data' is not a string")
	}

	var job any
	json.Unmarshal([]byte(values["data"].(string)), &job)

    return job, nil
}