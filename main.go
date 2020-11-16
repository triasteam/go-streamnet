// Package main
package main

import (
	//"bytes"

	"fmt"
	"log"
	"os"

	network "github.com/triasteam/go-streamnet/network/v2"
	"github.com/triasteam/go-streamnet/streamnet"

	"github.com/triasteam/go-streamnet/dag"
	"github.com/triasteam/go-streamnet/tipselection"

	streamnet_conf "github.com/triasteam/go-streamnet/config"

	//"io"
	//cmd "github.com/triasteam/go-streamnet/commands"

	"github.com/triasteam/go-streamnet/server"
	"github.com/triasteam/go-streamnet/store"
)

// GlobalData is running through the daemon.
var GlobalData streamnet.StreamNet

func main() {
	// set log config
	// todo: in debug mode, set the log module as following; if not in debug mode, don't set it.
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	initStreamWork()

	// start http server
	server.Start(&GlobalData)
}

func initStreamWork() {
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
	dag := dag.Dag{}
	dag.Init(&store)
	GlobalData.Dag = &dag

	// init tipselector
	tips := tipselection.TipSelectorStreamWork{}
	tips.Init(&dag)
	GlobalData.Tips = &tips

	// Set genesis trunk and branch

	// init libp2p
	node, err := network.NewNode(server.OnReceived)
	if err != nil {
		fmt.Printf("New Node error! err: %+v \n", err)
		os.Exit(-1)
	}
	// node.Init(server.OnReceived)
	GlobalData.Network = node

}
