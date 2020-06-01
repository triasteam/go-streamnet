package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type MysqlConfig struct {
	URL string
	UserName string
	Password string
}

type Config struct {
	Port int
	MySql MysqlConfig
}

func main() {
	var config Config

	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	fmt.Println(config)
}
