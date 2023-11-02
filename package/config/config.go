package config

import (
	"cloudMonitor/package/utils"
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigData struct {
	// COM config
	ComPort string `yaml:"comPort"`
	// cloud detector
	SampleSecond uint8   `yaml:"sampleSecond"`
	CenterWeight float32 `yaml:"centerWeight"`
	DetectRange  float32 `yaml:"detectRange"`
	// data management
	DataServer string `yaml:"dataServer"`
}

type Config interface {
	// config file management
	LoadFromYaml(filePath string) error
}

func (cfg *ConfigData) LoadFromYaml(filePath string) error {
	dataBytes, err := os.ReadFile(filePath)
	if nil != err {
		return errors.New(utils.ConcatStrings("Failed to read Yaml file ", filePath, ". Reason: ", err.Error()))
	}

	err = yaml.Unmarshal(dataBytes, cfg)
	if nil != err {
		return errors.New(utils.ConcatStrings("Failed to parse the Yaml content of ", filePath, ". Reason: ", err.Error()))
	}

	return nil
}
