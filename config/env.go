package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var EnvConfig = StreamNetconf{ //default value
	DBPath: "./db",
	Port:   ":14700",
	Redis:  redisConfig,
	GRPC:   grpcCOnfig,
}

var redisConfig = RedisConfig{ //default value
	Url:      "localhost:6379",
	Password: "",
	DB:       0,
}

var grpcCOnfig = GrpcConfig{ //default value
	Port: ":50051",
}

type RedisConfig struct {
	Url      string `yaml:"Url"`
	Password string `yaml:"Password"`
	DB       int    `yaml:"DB"`
}

type GrpcConfig struct {
	Port string `yaml:"Port"`
}

type StreamNetconf struct {
	DBPath string      `yaml:"DBPath"`
	Port   string      `yaml:"Port"`
	Redis  RedisConfig `yaml:"Redis"`
	GRPC   GrpcConfig  `yaml:"GRPC"`
}

func init() {
	var filePath = "config.yml"
	data, err := ioutil.ReadFile(filePath)
	if err == nil {
		yaml.Unmarshal(data, &EnvConfig)
	}
}
