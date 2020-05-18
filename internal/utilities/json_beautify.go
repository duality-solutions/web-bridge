package util

import (
	"encoding/json"
)

// BeautifyJSON takes raw json string and formats for readability
func BeautifyJSON(rawJSON string) (string, error) {
	mp := make(map[string]interface{})
	errUnmarshal := json.Unmarshal([]byte(rawJSON), &mp)
	if errUnmarshal != nil {
		return rawJSON, errUnmarshal
	}
	b, errMarshalIndent := json.MarshalIndent(mp, "", "   ")
	if errMarshalIndent != nil {
		return rawJSON, errMarshalIndent
	}
	return string(b), nil
}
