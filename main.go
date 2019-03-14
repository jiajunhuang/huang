package main

import (
	"log"
	"os"

	"github.com/jiajunhuang/huang/cmd"
	"github.com/jiajunhuang/huang/master"
	"github.com/jiajunhuang/huang/worker"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:  "master",
			Usage: "run as master",
			Action: func(c *cli.Context) error {
				return master.Main(c)
			},
		},
		{
			Name:  "worker",
			Usage: "run as worker",
			Action: func(c *cli.Context) error {
				return worker.Main(c)
			},
		},
	}

	app.Name = "huang"
	app.Usage = "$ huang"
	app.Action = func(c *cli.Context) error {
		return cmd.Main(c)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
