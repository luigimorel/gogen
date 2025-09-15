package fsutils

import (
	"encoding/json"
	"errors"
)

func NewJobFromJSON(b []byte) (*CreateJob, error) {
	if len(b) == 0 {
		return nil, errors.New("empty byte array")
	}

	var cj = CreateJob{}
	err := json.Unmarshal(b, &cj)
	if err != nil {
		return nil, err
	}

	return &cj, nil
}
