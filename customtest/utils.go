package customtest

import (
	"encoding/json"
	"encoding/xml"
	"testing"
)

func ToJson(t *testing.T, v interface{}) string {
	b, err := json.Marshal(v)
	Nil(t, err)
	return string(b)
}

func toXML(t *testing.T, v interface{}) string {
	b, err := xml.Marshal(v)
	Nil(t, err)

	return string(b)
}

func toDefault(t *testing.T, v interface{}) string {
	return ""
}
