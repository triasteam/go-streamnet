package main

import (
	//"bytes"
	//"fmt"
	//"github.com/triasteam/go-streamnet/types"

	"fmt"
	"os"
	"flag"
	"io/ioutil"
	streamnet_conf "github.com/triasteam/go-streamnet/config"
	//"io"
	//cmd "github.com/triasteam/go-streamnet/commands"
	"github.com/triasteam/go-streamnet/server"
	"github.com/triasteam/go-streamnet/store"
	//"github.com/triasteam/go-streamnet/store"
)

type StreamNet struct {
	Store *store.Storage
}

var GlobalData StreamNet

func init(){
	var filePath = "config.yml";
	data, err := ioutil.ReadFile(filePath);
	if(err != nil) {
		yaml.Unmarshal(data, &streamnet_conf.EnvConfig)
	} else {
		&streamnet_conf.EnvConfig{
			Port: flag.String("port",":14700","start port")
			DBPath: flag.String("dbpath","./db")
		}
	}
}

func main() {
	// open DB
	flag.Parse()
	store := store.Storage{}
	err := store.Init(streamnet_conf.EnvConfig.DBPath)
	if err != nil {
		fmt.Printf("Open database failed!")
		os.Exit(-1)
	}
	GlobalData.Store = &store

	// start http server
	server.Start(&store)

	/*
		rootCmd := cmd.RootCmd
		rootCmd.AddCommand(
			cmd.InitFilesCmd,
		)

		// parse config.  examples: sng --mwm 1 -p 14700 &>  sng.log &
		// other parameters like '--enable-streaming-graph' '--entrypoint-selector-algorithm "KATZ"' '--tip-sel-algo "CONFLUX"' '--weight-calculation-algorithm "IN_MEM"'
		// will be considered later.

		// start server
		// Create & start node
		rootCmd.AddCommand(cmd.NewRunNodeCmd())

	*/
}
