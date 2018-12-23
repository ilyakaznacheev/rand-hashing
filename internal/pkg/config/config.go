package config

import (
	"errors"
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

// Config is an application config structure
type Config struct {
	Redis struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`
}

// ReadConfig reads configuration from yaml config file
func ReadConfig(path string) (*Config, error) {
	confFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("config error: " + err.Error())
	}

	conf := &Config{}
	err = yaml.Unmarshal(confFile, conf)
	if err != nil {
		return nil, errors.New("config error: " + err.Error())
	}
	return conf, nil
}
