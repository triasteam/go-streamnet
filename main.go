// Package main
package main

import (
	//"bytes"

	"fmt"
	"os"

	"github.com/triasteam/go-streamnet/dag"

	"github.com/triasteam/go-streamnet/tipselection"

	streamnet_conf "github.com/triasteam/go-streamnet/config"

	//"io"
	//cmd "github.com/triasteam/go-streamnet/commands"
	"github.com/triasteam/go-streamnet/server"
	"github.com/triasteam/go-streamnet/store"
)

// StreamNet is the biggest structure.
type StreamNet struct {
	dag   *dag.Dag
	Store *store.Storage
	ts    tipselection.TipSelector
}

// GlobalData is running through the daemon.
var GlobalData StreamNet

func main() {

	// start http server
	server.Start(GlobalData.Store)
}

func init() {
	// open DB
	store := store.Storage{}
	fmt.Println("Port: " + streamnet_conf.EnvConfig.Port + ", DBpath: " + streamnet_conf.EnvConfig.DBPath)
	err := store.Init(streamnet_conf.EnvConfig.DBPath)
	if err != nil {
		fmt.Printf("Open database failed!")
		os.Exit(-1)
	}
	GlobalData.Store = &store

	// init dag
	GlobalData.dag = dag.Init(&store)

	// init tipselector
	var ts tipselection.TipSelectorStreamWork
	ts.dag
	GlobalData.
}
