package main

import (
	"encoding/json"
	. "github.com/tj/go-debug"
	"strings"
)

type jsonFormatter struct {
	debug DebugFunction
}

func newJsonFormatter() *jsonFormatter {
	return &jsonFormatter{
		debug: Debug("jsonFormatter"),
	}
}

func (j *jsonFormatter) Format(iput string, fields []string) string {
	var in, out map[string]interface{}
	out = make(map[string]interface{})
	err := json.Unmarshal([]byte(iput), &in)
	if err != nil {
		return "invalid line"
	}

	for _, field := range fields {
		j.addJsonField(field, in, out)
	}

	res, err := json.Marshal(out)
	if err != nil {
		return "invalid line"
	}

	return string(res)
}

// adds one of the specified fields to the output
func (j *jsonFormatter) addJsonField(field string, source, dest map[string]interface{}) {
	fieldParts := strings.Split(field, ".")

	j.debug("have some fields", len(fieldParts))

	if len(fieldParts) == 1 && source[fieldParts[0]] != nil {
		//only one part? just shove it in the destination map
		dest[fieldParts[0]] = source[fieldParts[0]]
	} else if len(fieldParts) > 1 {
		//anything else gets more complicated...
		var prevSource, prevDest map[string]interface{}

		if source[fieldParts[0]] == nil {
			j.debug("first part of field not found in source: " + fieldParts[0])
			return
		}

		dest[fieldParts[0]] = make(map[string]interface{})
		prevDest = dest[fieldParts[0]].(map[string]interface{})
		prevSource = source[fieldParts[0]].(map[string]interface{})

		for i := 1; i < len(fieldParts)-1; i += 1 {
			curName := fieldParts[i]

			if prevSource[curName] == nil {
				j.debug(field + " failed on " + curName)
				return
			}

			prevDest[curName] = make(map[string]interface{})

			prevSource = prevSource[curName].(map[string]interface{})
			prevDest = prevDest[curName].(map[string]interface{})
		}

		j.debug("hopefully setting " + fieldParts[len(fieldParts)-1])
		prevDest[fieldParts[len(fieldParts)-1]] = prevSource[fieldParts[len(fieldParts)-1]]
	}
}
