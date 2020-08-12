package config

import(
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

var EnvConfig = StreamNetconf{ //default value
	DBPath: "./db",
	Port: ":14700",
	Redis: redisConfig,
}

var redisConfig = RedisConfig{ //default value
	Url: "localhost:6379",
	Password: "",
	DB: 0,
}

type RedisConfig struct {
	Url 		string	`yaml:"Url"`
	Password 	string	`yaml:"Password"`
	DB 			int		`yaml:"DB"`
}
/*

Redis: 
  Port: 6379
  Password:
  DB: 0
*/
type StreamNetconf struct {
	DBPath		string		`yaml:"DBPath"`
	Port		string		`yaml:"Port"`
	Redis		RedisConfig	`yaml:"Redis"`
}

func init(){
	var filePath = "config.yml";
	data, err := ioutil.ReadFile(filePath);
	if(err == nil) {
		yaml.Unmarshal(data, &EnvConfig)
	}
}
