package main

import (
	"os"

	yaml "gopkg.in/yaml.v3"
)

type (
	ConfigFile struct {
		Tapper *TapperConfig `yaml:"tapper"`
	}

	TapperConfig struct {
		Options *Options
		Device  *DeviceConfig
		Actions []ActionConfig
	}

	Options struct {
		Interval int64
	}

	DeviceConfig struct {
		Name string
		Path string
	}

	ActionConfig struct {
		Pattern []string
		Command []string
	}
)

func ReadConfiguration(path string) (*TapperConfig, error) {
	var config ConfigFile

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(content, &config); err != nil {
		return nil, err
	}

	return config.Tapper, nil
}
