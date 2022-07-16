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
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "transaction",
						Aliases:  []string{"-t"},
						Usage:    "name of the transaction",
						Required: true,
					},
				},
				Subcommands: []*cli.Command{
					{
						Name:  "add",
						Usage: "Add repository to the transaction",
						Action: func(ctx *cli.Context) error {
							repo := ctx.Args().First()
							transaction := ctx.String("transaction")

							result, err := client.RepoAdd(transaction, repo)
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
							transaction := ctx.String("transaction")

							result, err := client.RepoRemove(transaction, value)
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
							transaction := ctx.String("transaction")

							result, err := client.RepoList(transaction)
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
				Name:    "transaction",
				Aliases: []string{"tr"},
				Usage:   "Manage transactions",
				Subcommands: []*cli.Command{
					{
						Name:  "get",
						Usage: "get a transaction",
						Action: func(ctx *cli.Context) error {
							name := ctx.Args().First()

							result, err := client.TransactionGet(name)
							if err != nil {
								return err
							}

							fmt.Println(result)

							return nil
						},
					},
					{
						Name:  "add",
						Usage: "add a transaction",
						Action: func(ctx *cli.Context) error {
							name := ctx.Args().First()

							if err := client.TransactionAdd(name); err != nil {
								return err
							}

							return nil
						},
					},
					{
						Name:  "remove",
						Usage: "remove a transaction",
						Action: func(ctx *cli.Context) error {
							name := ctx.Args().First()

							if err := client.TransactionRemove(name); err != nil {
								return err
							}

							return nil
						},
					},
					{
						Name:  "rename",
						Usage: "rename a transaction",
						Action: func(ctx *cli.Context) error {
							oldName := ctx.Args().First()
							newName := ctx.Args().Get(1)

							if err := client.TransactionRename(oldName, newName); err != nil {
								return err
							}

							return nil
						},
					},
					{
						Name:  "list",
						Usage: "list a transaction",
						Action: func(ctx *cli.Context) error {
							result, err := client.TransactionList()
							if err != nil {
								return err
							}

							fmt.Println(result)

							return nil
						},
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
