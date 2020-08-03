// Package main
package main

import (
	//"bytes"
	//"fmt"
	//"github.com/triasteam/go-streamnet/types"

	"fmt"
	"os"

	//"io"
	//cmd "github.com/triasteam/go-streamnet/commands"
	"github.com/triasteam/go-streamnet/server"
	"github.com/triasteam/go-streamnet/store"
	//"github.com/triasteam/go-streamnet/store"
)

// StreamNet is the biggest structure.
type StreamNet struct {
	Store *store.Storage
}

// GlobalData is running through the daemon.
var GlobalData StreamNet

// the daemon main function.
func main() {
	// open DB
	store := store.Storage{}
	err := store.Init("./db")
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
