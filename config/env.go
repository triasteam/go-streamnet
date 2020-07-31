package config

import (
	"gopkg.in/yaml.v2"
)

var EnvConfig StreamNetconf{}

type StreamNetconf struct {
	DBPath     string  `yaml:"DBPath"`
	Port       string  `yaml:"Port"`
}
