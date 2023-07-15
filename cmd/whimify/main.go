package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/Hofled/whimify/internal/function/sort"
)

func main() {
	app := &cli.App{
		HelpName: "whimify",
		Usage:    "satisfies every Go code whim",
		Commands: []*cli.Command{
			{
				Name:    "sort",
				Aliases: []string{"s"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "path",
						Usage:    "the path to the file which will get sorted",
						Required: true,
					},
				},
				Action: func(ctx *cli.Context) error {
					path := ctx.String("path")
					if path == "" {
						log.Println("No path provided")
						return nil
					}
					return sort.FileByDependencyOrder(path)
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
