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
				Name:  "swap-endpoint, sp",
				Usage: "used for bee start for swap-endpoint",
			},
			&cli.StringFlag{
				Name:  "swap-enable, se",
				Value: "true",
				Usage: "used for bee start for swap-enable",
			},
			&cli.StringFlag{
				Name:  "swap-deployment-gas-price, sg",
				Value: "1000000000000",
				Usage: "used for bee start for swap-deployment-gas-price",
			},
			&cli.StringFlag{
				Name:  "swap-initial-deposit, sid",
				Value: "0",
				Usage: "used for bee start for swap-initial-deposit",
			},
			&cli.StringFlag{
				Name:  "debug-api-enable, d",
				Value: "true",
				Usage: "used for bee start for debug-api-enable",
			},
			&cli.StringFlag{
				Name:  "network-id, n",
				Value: "1",
				Usage: "used for bee start for network-id",
			},
			&cli.StringFlag{
				Name:  "main-net, m",
				Value: "true",
				Usage: "used for bee start for main net",
			},
			&cli.StringFlag{
				Name:  "full-node, f",
				Value: "true",
				Usage: "used for bee start for full-node",
			},
			&cli.StringFlag{
				Name:  "verbosity, v",
				Value: "info",
				Usage: "used for bee start for log verbosity",
			},
			&cli.StringFlag{
				Name:  "clef-signer-enable, c",
				Value: "false",
				Usage: "used for bee start for clef-signer-enable",
			},
			&cli.StringFlag{
				Name:  "i",
				Usage: "docker image name",
			},
			&cli.StringFlag{
				Name:  "password, p",
				Usage: "used for bee start for password",
			},
			&cli.StringFlag{
				Name:  "data-dir, dd",
				Value: "/bee/file",
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
			return core.NewStaticStrategy(context.String("sp"), context.String("se"), context.String("sg"),
				context.String("sid"), context.String("d"), context.String("n"), context.String("m"),
				context.String("f"), context.String("v"), context.String("c"), context.String("i"),
				context.String("p"), context.String("dd")).Run()
		},
	}
	return cmd
}
