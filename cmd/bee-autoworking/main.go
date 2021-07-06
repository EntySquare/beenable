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
	//restartingCMD := restart()

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
				Name:        "start",
				Usage:       "bee manual start",
				Subcommands: []*cli.Command{manualStartingCMD},
			},
			//{
			//	Name:        "restart",
			//	Usage:       "bee restart",
			//	Subcommands: []*cli.Command{restartingCMD},
			//},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func manual() *cli.Command {
	cmd := &cli.Command{
		Name:  "bee",
		Usage: "starting bee based on manual options",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "swap-endpoint",
				Aliases: []string{"sp"},
				Usage:   "used for bee start for swap-endpoint",
			},
			&cli.StringFlag{
				Name:    "swap-enable",
				Aliases: []string{"se"},
				Value:   "true",
				Usage:   "used for bee start for swap-enable",
			},
			&cli.StringFlag{
				Name:    "swap-deployment-gas-price",
				Aliases: []string{"sg"},
				Value:   "1000000000000",
				Usage:   "used for bee start for swap-deployment-gas-price",
			},
			&cli.StringFlag{
				Name:    "swap-initial-deposit",
				Aliases: []string{"sid"},
				Value:   "0",
				Usage:   "used for bee start for swap-initial-deposit",
			},
			&cli.StringFlag{
				Name:    "debug-api-enable",
				Aliases: []string{"d"},
				Value:   "true",
				Usage:   "used for bee start for debug-api-enable",
			},
			&cli.StringFlag{
				Name:    "network-id",
				Aliases: []string{"n"},
				Value:   "1",
				Usage:   "used for bee start for network-id",
			},
			&cli.StringFlag{
				Name:    "mainnet",
				Aliases: []string{"m"},
				Value:   "true",
				Usage:   "used for bee start for main net",
			},
			&cli.StringFlag{
				Name:    "full-node",
				Aliases: []string{"f"},
				Value:   "true",
				Usage:   "used for bee start for full-node",
			},
			&cli.StringFlag{
				Name:    "verbosity",
				Aliases: []string{"v"},
				Value:   "info",
				Usage:   "used for bee start for log verbosity",
			},
			&cli.StringFlag{
				Name:    "clef-signer-enable",
				Aliases: []string{"c"},
				Value:   "false",
				Usage:   "used for bee start for clef-signer-enable",
			},
			&cli.StringFlag{
				Name:    "docker-image",
				Aliases: []string{"i"},
				Usage:   "docker image name",
			},
			&cli.StringFlag{
				Name:    "password",
				Aliases: []string{"p"},
				Usage:   "used for bee start for password",
			},
			&cli.StringFlag{
				Name:    "data-dir",
				Aliases: []string{"dd"},
				Value:   "/bee/file",
				Usage:   "used for bee start for data-dir",
			},
			&cli.StringFlag{
				Name:    "api-addr",
				Aliases: []string{"ap"},
				Value:   "1633",
				Usage:   "HTTP API listen address port",
			},
			&cli.StringFlag{
				Name:    "p2p-addr",
				Aliases: []string{"pp"},
				Value:   "1634",
				Usage:   "P2P listen address port",
			},
			&cli.StringFlag{
				Name:    "debug-api-addr",
				Aliases: []string{"dp"},
				Value:   "1635",
				Usage:   "debug HTTP API listen address port",
			},
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

//func restart() *cli.Command {
//
//}
