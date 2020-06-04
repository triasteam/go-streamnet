package main

import (
	"fmt"
	cmd "github.com/triasteam/StreamNet-go/commands"
)
func main() {
	fmt.Println("hello, streamnet-go")

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
}
