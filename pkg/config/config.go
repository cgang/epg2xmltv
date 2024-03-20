package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ChannelConfig struct {
	Id     string `yaml:"id"`
	Name   string `yaml:"name"`
	Source string `yaml:"source"`
}

type XmlConfig struct {
	Name     string          `yaml:"name"`
	Channels []ChannelConfig `yaml:"channels"`
}

func LoadConfigs(filename string) ([]XmlConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfgs []XmlConfig
	if err = yaml.NewDecoder(file).Decode(&cfgs); err == nil {
		return cfgs, nil
	} else {
		return nil, err
	}
}
