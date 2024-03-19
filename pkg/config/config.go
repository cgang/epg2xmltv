package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type CrawlerConfig struct {
	Type string `yaml:"type"`
	Id   string `yaml:"id"`
	Arg  string `yaml:"arg"`
}

type ChannelConfig struct {
	Id   string `yaml:"id"`
	Name string `yaml:"name"`
}

type XmlConfig struct {
	Name     string          `yaml:"name"`
	Channels []ChannelConfig `yaml:"channels"`
}

type AppConfig struct {
	CrawlersConfig []CrawlerConfig `yaml:"crawlers"`
	OutputsConfig  []XmlConfig     `yaml:"outputs"`
}

func LoadConfig(filename string) (*AppConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg AppConfig
	if err = yaml.NewDecoder(file).Decode(&cfg); err == nil {
		return &cfg, nil
	} else {
		return nil, err
	}
}
