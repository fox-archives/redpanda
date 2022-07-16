package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hyperupcall/redpanda/client-cli/client"
	cli "github.com/urfave/cli/v2"
)

func main() {
	client := client.New()

	app := &cli.App{
		Name:    "redpanda",
		Usage:   "Multi-repository refactoring tool",
		Version: "0.1.0",
		Authors: []*cli.Author{
			{
				Name:  "Edwin Kofler",
				Email: "edwin@kofler.dev",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "repo",
				Usage: "Manager repos in current transaction",
				Subcommands: []*cli.Command{
					{
						Name:  "add",
						Usage: "Add repository to the transaction",
						Action: func(ctx *cli.Context) error {
							value := ctx.Args().First()
							result, err := client.RepoAdd(value)
							if err != nil {
								return err
							}
							fmt.Println(result)

							return nil
						},
					},
					{
						Name:  "remove",
						Usage: "Remove repository to the current transaction",
						Action: func(ctx *cli.Context) error {
							value := ctx.Args().First()
							result, err := client.RepoRemove(value)
							if err != nil {
								return err
							}
							fmt.Println(result)

							return nil
						},
					},
					{
						Name:  "list",
						Usage: "List all repositories in the current transaction",
						Action: func(ctx *cli.Context) error {
							result, err := client.RepoList()
							if err != nil {
								return err
							}
							fmt.Println(result)

							return nil
						},
					},
				},
			},
			{
				Name:    "transactions",
				Aliases: []string{"tr"},
				Usage:   "Manage transactions",
				Subcommands: []*cli.Command{
					{
						Name: "repo",
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
