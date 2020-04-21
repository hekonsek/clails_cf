package clails

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

type YmlProject struct {
}

func NewYmlProject() *YmlProject {
	return &YmlProject{}
}

func (ymlProject *YmlProject) LoadFromFile(path string) (*Project, error) {
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

func (ymlProject *YmlProject) FriendlyMessage(err error) error {
	if err != nil {
		if strings.HasPrefix(err.Error(), "no such file") {
			return errors.New(fmt.Sprintf(
				"There is no " + color.GreenString("clails.yml") + " file in current directory. " +
					"Please create Clails project file and save it in current directory as " +
					color.GreenString("clails.yml") + " file."))
		} else {
			return errors.New("Something went wrong: " + err.Error())
		}
	}

	return nil
}