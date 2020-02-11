package clails

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func LoadProjectFromYmlFile(path string) (*Project, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("no such file: %s", path))
	}

	projectFileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	model := &Project{}
	err = yaml.Unmarshal(projectFileBytes, model)
	if err != nil {
		return nil, err
	}

	return model, nil
}
