package main

import (
	"encoding/json"
)

type jsonFormatter struct{}

func (j *jsonFormatter) Format(iput string, fields []string) string {
	var in, out map[string]interface{}
	out = make(map[string]interface{})
	err := json.Unmarshal([]byte(iput), &in)
	if err != nil {
		return "invalid line"
	}

	for k, v := range out {
		if contains(k, fields) {
			out[k] = v
		}
	}

	res, err := json.Marshal(out)
	if err != nil {
		return "invalid line"
	}

	return string(res)
}

func contains(needle string, haystack []string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}

	return false
}
