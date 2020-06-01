package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type MysqlConfig struct {
	U string            `json:"url"`
	Un string	    `json:"username"`
	Pw string           `json:"password"`
}

type Config struct {
	Po int   `json: "port"`
	MySQl MysqlConfig `json: "mysql"`
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
