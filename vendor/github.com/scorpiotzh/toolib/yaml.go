package toolib

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func UnmarshalYamlFile(filePath string, r interface{}) error {
	if file, err := ioutil.ReadFile(filePath); err != nil {
		return err
	} else if err = yaml.Unmarshal(file, r); err != nil {
		return err
	}
	return nil
}
