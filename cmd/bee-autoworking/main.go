package main

import (
	"beenable/core"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sort"
)

func main() {

	manualStartingCMD := manual()
	autoStartingCMD := auto()

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "debug",
				Value: "debug",
				Usage: "Debug commands",
			},
		},
		Commands: []*cli.Command{
			{
				Name:        "manual",
				Usage:       "manual start",
				Subcommands: []*cli.Command{manualStartingCMD},
			},
			{
				Name:        "auto",
				Usage:       "automatic start",
				Subcommands: []*cli.Command{autoStartingCMD},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func auto() *cli.Command {
	cmd := &cli.Command{
		Name:  "bee",
		Usage: "starting bee based on manual options",
		Flags: []cli.Flag{},
		Action: func(context *cli.Context) error {
			return core.NewDynamicStrategy().Run()
		},
	}
	return cmd

}

func manual() *cli.Command {
	cmd := &cli.Command{
		Name:  "bee",
		Usage: "starting bee based on manual options",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "sp",
				Usage: "used for bee start for swap-endpoint",
			},
			&cli.StringFlag{
				Name:  "se",
				Usage: "used for bee start for swap-enable",
			},
			&cli.StringFlag{
				Name:  "sg",
				Usage: "used for bee start for swap-deployment-gas-price",
			},
			&cli.StringFlag{
				Name:  "sid",
				Usage: "used for bee start for swap-initial-deposit",
			},
			&cli.StringFlag{
				Name:  "d",
				Usage: "used for bee start for debug-api-enable",
			},
			&cli.StringFlag{
				Name:  "n",
				Usage: "used for bee start for network-id",
			},
			&cli.StringFlag{
				Name:  "m",
				Usage: "used for bee start for main net",
			},
			&cli.StringFlag{
				Name:  "f",
				Usage: "used for bee start for full-node",
			},
			&cli.StringFlag{
				Name:  "v",
				Usage: "used for bee start for verbosity",
			},
			&cli.StringFlag{
				Name:  "c",
				Usage: "used for bee start for clef-signer-enable",
			},
			&cli.StringFlag{
				Name:  "i",
				Usage: "docker image name",
			},
			&cli.StringFlag{
				Name:  "p",
				Usage: "used for bee start for password",
			},
			&cli.StringFlag{
				Name:  "d",
				Usage: "used for bee start for data-dir",
			},
			//&cli.StringFlag{
			//	Name:  "rp",
			//	Usage: "ip to report",
			//},
			//&cli.StringFlag{
			//	Name:  "po",
			//	Usage: "port to report",
			//},
		},
		Action: func(context *cli.Context) error {
			return core.NewStaticStrategy(context.String("i")).Run()
		},
	}
	return cmd
}
