package config

import(
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

var EnvConfig = StreamNetconf{ //default value
	DBPath: "./db",
	Port: "14700",
}

type StreamNetconf struct {
	DBPath     string  `yaml:"DBPath"`
	Port       string  `yaml:"Port"`
}

func init(){
	var filePath = "config.yml";
	data, err := ioutil.ReadFile(filePath);
	if(err == nil) {
		yaml.Unmarshal(data, &EnvConfig)
	}
}
