package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type UserInfo struct {
	UserName string
	Address string
	Sex byte
	Company Company
}

type Company struct {
	Name string
	EmployeeId int
	Department []interface{}
}


func main() {
	//读取yaml文件
	v := viper.New()
	//设置读取的配置文件名
	v.SetConfigName("userInfo")
	//windows环境下为%GOPATH，linux环境下为$GOPATH
	v.AddConfigPath("/Users/trust/go/src/github.com/triasteam/StreamNet-go/examples/config/viper")
	//设置配置文件类型
	v.SetConfigType("yaml")

	if err := v.ReadInConfig();err != nil {
		fmt.Printf("err:%s\n",err)
	}

	fmt.Printf("userName:%s sex:%s company.name:%s \n", v.Get("userName"), v.Get("sex"), v.Get("company.name"))


  //也可以直接反序列化为Struct

	var userInfo UserInfo
	if err := v.Unmarshal(&userInfo) ; err != nil{
		fmt.Printf("err:%s",err)
	}
	fmt.Println(userInfo)
}
