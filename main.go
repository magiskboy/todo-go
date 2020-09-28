package main

import (
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name: "worker",
			Action: func(ctx *cli.Context) {
				err := StartWorker()
				if err != nil {
					panic(err)
				}
			},
		},
		{
			Name: "web",
			Action: func(ctx *cli.Context) {
				InitDB()
				router := SetupRouter()
				router.Run()
			},
		},
	}
	app.Run(os.Args)
}
