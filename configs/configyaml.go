package config

import "gopkg.in/yaml.v2"

func unmarshalYaml(content []byte, v interface{}) error {
	return yaml.Unmarshal(content, v)
}

func marshalYaml(v interface{}) (out []byte, err error) {
	return yaml.Marshal(v)
}

func marshalYamlString(v interface{}) (out string) {
	marshal, err := yaml.Marshal(v)
	if err != nil {
		return ""
	}
	return string(marshal)
}
