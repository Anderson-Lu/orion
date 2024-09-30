package parser

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

func ParseConfigFile(fsPath string, target interface{}) error {
	bs, err := os.ReadFile(fsPath)
	if err != nil {
		return err
	}
	if strings.HasSuffix(fsPath, ".yaml") {
		if err := yaml.Unmarshal(bs, target); err != nil {
			return err
		}
		return nil
	}
	if strings.HasSuffix(fsPath, ".json") {
		return json.Unmarshal(bs, target)
	}
	if strings.HasSuffix(fsPath, ".toml") {
		_, err := toml.Decode(string(bs), target)
		return err
	}
	return errors.New("unsupport config file type")
}
