package config

import "encoding/xml"

func unmarshalXML(content []byte, v interface{}) error {
	return xml.Unmarshal(content, v)
}

func marshalXML(v interface{}) (out []byte, err error) {
	return xml.Marshal(v)
}

func marshalXMLString(v interface{}) (out string) {
	marshal, err := xml.Marshal(v)
	if err != nil {
		return ""
	}

	return string(marshal)
}
