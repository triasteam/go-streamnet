package main

import (
	"fmt"
	"os"
	"github.com/urfave/cli"
)

func main() {
	/* usage: cli -url 192.168.111.11 -dp 6666
	*/
	var mysqlUrl string
	var debugPort uint

	app := cli.NewApp()
	app.Name = "sysconfig"
	app.Usage = "Setting basic configuration"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag {
			Name: "mysqlurl, url",
			Usage: "specify mysql address",
			Value: "182.122.11.11",
			Destination: &mysqlUrl,
		},
		cli.UintFlag {
			Name: "debugport, dp",
			Usage: "specify debug port",
			Value: 9666,
			Destination: &debugPort,
		},
	}

	app.Action = func(c *cli.Context) error {
		fmt.Println("Prepare applying basic configuration")

		fmt.Println("Mysql Url:", mysqlUrl)
		fmt.Println("Debug port:", debugPort)

		return nil
	}

	app.Run(os.Args)
}
