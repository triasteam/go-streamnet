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

	app.Commands = []cli.Command {
		{
			Name: "debug",
			Aliases: []string{"d"},
			Usage: "debug settings",
			Action: func(c *cli.Context) error {
				fmt.Println("Set debug grade to ", c.Args().First())
				return nil
			},

		},
		{
			Name: "mysql",
			Aliases: []string{"q"},
			Usage: "mysql operations",
			Subcommands: []cli.Command {
				{
					Name: "query",
					Usage: "query userid",
					Action: func(c *cli.Context) error {
						fmt.Println("query userid=", c.Args().First())
						return nil
					},
				},
			},

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
