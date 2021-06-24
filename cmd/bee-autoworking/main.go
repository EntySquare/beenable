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
			//&cli.StringFlag{
			//	Name:  "p",
			//	Usage: "used for plotting for cert pool",
			//},
			//&cli.StringFlag{
			//	Name:  "f",
			//	Usage: "used for plotting for cert farmer",
			//},
			//&cli.StringFlag{
			//	Name:  "d",
			//	Usage: "user dir use to match move",
			//},
			&cli.StringFlag{
				Name:  "i",
				Usage: "docker image name",
			},
			//&cli.StringFlag{
			//	Name:  "k",
			//	Usage: "k of plot",
			//},
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
