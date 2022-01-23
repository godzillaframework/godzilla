/**
 * reference links:
 *	 - https://gobyexample.com/json
 *	 - https://golang.cafe/blog/golang-json-marshal-example.html
 *	 - https://pkg.go.dev/encoding/json
 */

package config

import "encoding/json"

func unmarshalJSON(content []byte, v interface{}) error {
	return json.Unmarshal(content, v)
}

func marshalJSON(v interface{}) (out []byte, err error) {
	return json.Marshal(v)
}

func marshalJSONString(v interface{}) (out string) {
	marshal, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(marshal)
}
