package main

import (
	"fmt"
	"net/http"
	//cmd "github.com/triasteam/StreamNet-go/commands"
)

func SaveHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello, save")
}

func GetHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello, get")
}

func main() {
	fmt.Println("hello, streamnet-go")

	http.HandleFunc("/save", SaveHandle)
	http.HandleFunc("/get", GetHandle)
	http.ListenAndServe(":14700", nil)

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
