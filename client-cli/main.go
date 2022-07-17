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
				Name:  "step",
				Usage: "Performs a step",
				Subcommands: []*cli.Command{
					{
						Name:  "initialize",
						Usage: "Initializes uninitialized repositories",
						Action: func(ctx *cli.Context) error {
							result, err := client.StepInitialize()
							if err != nil {
								return err
							}
							fmt.Println(result)

							return nil
						},
					},
					{
						Name:  "idempotent-apply",
						Usage: "Idempotently apply transformations",
						Action: func(ctx *cli.Context) error {
							transaction := ctx.Args().First()

							result, err := client.StepIdempotentApply(transaction)
							if err != nil {
								return err
							}
							fmt.Println(result)

							return nil
						},
					},
					{
						Name:  "diff",
						Usage: "View resulting diff",
						Action: func(ctx *cli.Context) error {
							transaction := ctx.Args().First()

							result, err := client.StepDiff(transaction)
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
				Name:    "transformers",
				Aliases: []string{"tf"},
				Usage:   "Manage, of particular transaction, the transformers",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "transaction",
						Aliases:  []string{"t"},
						Usage:    "Name of the transaction",
						Required: true,
					},
				},
				Subcommands: []*cli.Command{
					{
						Name:  "add",
						Usage: "add a transformer",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "type",
								Usage:    "Type of transformer (command,regex)",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "content",
								Usage:    "Content of transformer",
								Required: true,
							},
						},
						Action: func(ctx *cli.Context) error {
							transformer := ctx.Args().First()
							transaction := ctx.String("transaction")
							typ := ctx.String("type")
							content := ctx.String("content")

							if typ != "command" && typ != "regex" {
								return fmt.Errorf("Type must be either command or regex")
							}

							result, err := client.TransformerAdd(transaction, typ, transformer, content)
							if err != nil {
								return err
							}

							fmt.Println(result)

							return nil
						},
					},
					{
						Name:  "remove",
						Usage: "remove a transformer",
						Action: func(ctx *cli.Context) error {
							transformer := ctx.Args().First()
							transaction := ctx.String("transaction")

							result, err := client.TransformerRemove(transaction, transformer)
							if err != nil {
								return err
							}

							fmt.Println(result)

							return nil
						},
					},
					{
						Name:  "edit",
						Usage: "edit a transformer",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "content",
								Usage:    "Content of transformer",
								Required: true,
							},
						},
						Action: func(ctx *cli.Context) error {
							transformer := ctx.Args().First()
							transaction := ctx.String("transaction")
							newContent := ctx.String("content")

							result, err := client.TransformerEdit(transaction, transformer, newContent)
							if err != nil {
								return err
							}

							fmt.Println(result)

							return nil
						},
					},
					{
						Name:  "order",
						Usage: "order transformers (TODO)",
						Action: func(ctx *cli.Context) error {
							order := ctx.Args().First()
							transaction := ctx.String("transaction")

							result, err := client.TransformerOrder(transaction, order)
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
				Name:  "repo",
				Usage: "Manage, of particular transaction, the repos",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "transaction",
						Aliases:  []string{"t"},
						Usage:    "Name of the transaction",
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
				},
			},
			{
				Name:    "transaction",
				Aliases: []string{"ta"},
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
