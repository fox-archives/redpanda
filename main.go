package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/go-github/github"
	git "github.com/libgit2/git2go/v33"
	cli "github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func startServer() {
	key := "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" // TODO

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: key},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	repositories, _, err := client.Repositories.List(context.Background(), "", &github.RepositoryListOptions{
		Type: "all",
	})
	handle(err)

	for _, repo := range repositories {
		// defaultBranch = repo.GetDefaultBranch()
		// description := repo.GetDescription()
		// url := repo.GetGitCommitsURL()

		fmt.Printf("%s/%s\n", repo.GetOwner().GetLogin(), repo.GetName())
		fmt.Println(repo)
		fmt.Println()
	}

}

func main() {
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
				Name:  "server",
				Usage: "launch server",
				Action: func(*cli.Context) error {
					startServer()

					return nil
				},
			},
			{
				Name:  "repo",
				Usage: "Manager repos in current transaction",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "add",
						Usage: "Add repository to the transaction",
					},
					&cli.StringFlag{
						Name:  "remove",
						Usage: "Remove repository to the current transaction",
					},
					&cli.BoolFlag{
						Name:  "list",
						Usage: "List all repositories in the current transaction",
					},
				},
				Action: func(ctx *cli.Context) error {
					addFlag := ctx.String("add")
					removeFlag := ctx.String("remove")
					listFlag := ctx.Bool("list")

					if addFlag != "" {
						repos := strings.Split(addFlag, ",")
						ReposAdd(repos)
					} else if removeFlag != "" {
						repos := strings.Split(removeFlag, ",")
						ReposRemove(repos)
					} else if listFlag {
						repos := ReposList()
						fmt.Println(repos)
					} else {
						return fmt.Errorf("Must pass in a flag")
					}
					return nil
				},
			},
			{
				Name:  "action",
				Usage: "run an action on the repositories",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "sync",
						Usage: "synchronize repositories",
					},
					&cli.StringFlag{
						Name:  "cmd",
						Usage: "Run a raw command in all the repositories",
					},
					&cli.BoolFlag{
						Name:  "reset",
						Usage: "Reset to HEAD",
					},
					&cli.BoolFlag{
						Name:  "transaction",
						Usage: "Commit transaction",
					},
					&cli.BoolFlag{
						Name:  "transaction-post",
						Usage: "git push (be careful!)",
					},
				},
				Action: func(ctx *cli.Context) error {
					syncFlag := ctx.Bool("sync")
					cmdFlag := ctx.String("cmd")
					// resetFlag := ctx.Bool("reset")
					// transactionFlag := ctx.Bool("transaction")
					// transactionPostFlag := ctx.Bool("transaction-post")

					repos := ReposList()
					reposDir := filepath.Join(os.Getenv("HOME"), ".local", "share", "redpanda", "repositories")
					err := os.MkdirAll(reposDir, 0o755)
					if err != nil {
						return err
					}

					if syncFlag {
						for _, repoName := range repos {
							repoDir := filepath.Join(reposDir, repoName)

							_, err := os.Stat(filepath.Join(repoDir, ".git"))
							if errors.Is(err, os.ErrNotExist) {
								gitRepo, err := git.Clone("https://github.com/"+repoName, repoDir, &git.CloneOptions{
									CheckoutOptions: git.CheckoutOptions{
										Strategy: git.CheckoutSafe,
									},
								})
								if err != nil {
									return err
								}
								fmt.Println(repoDir)
								ref, err := gitRepo.Head()
								if err != nil {
									return err
								}
								fmt.Println(ref.Name())
							} else {
								gitRepo, err := git.OpenRepository(repoDir)
								if err != nil {
									return err
								}

								ref, err := gitRepo.Head()
								if err != nil {
									return err
								}
								fmt.Printf("repo: %s (%s)\n", repoDir, ref.Name())
							}
						}
					} else if cmdFlag != "" {
						for _, gitRepo := range utilForEachRepo(reposDir, repos) {

							ref, err := gitRepo.References.Lookup("HEAD")
							if err != nil {
								return err
							}

							headRef, err := ref.Peel(git.ObjectCommit)
							if err != nil {
								return err
							}

							c, err := headRef.AsCommit()
							if err != nil {
								return err
							}

							_, err = gitRepo.CreateBranch("redpanda-workbranch", c, false)
							if err != nil {
								if !git.IsErrorCode(err, git.ErrorCodeExists) {
									return err
								}
							}

							fmt.Println("---")
							fmt.Println(gitRepo.Path())
							fmt.Println("---")
							cmd := exec.Command("git", append([]string{"-C", filepath.Dir(strings.TrimRight(gitRepo.Path(), "/"))}, strings.Split(cmdFlag, " ")...)...)
							cmd.Stdin = os.Stdin
							cmd.Stdout = os.Stdout
							cmd.Stderr = os.Stderr
							err = cmd.Run()
							if err != nil {
								return err
							}
							// br, err := gitRepo.LookupBranch("redpanda-workbranch", git.BranchLocal)
							// if err != nil {
							// 	return err
							// }

							// p, err := br.Peel(git.ObjectTree)
							// if err != nil {
							// 	return err
							// }

							// t, err := p.AsTree()
							// if err != nil {
							// 	return err
							// }

							// err = gitRepo.CheckoutTree(t, &git.CheckoutOptions{
							// 	Strategy: git.CheckoutForce,
							// })
							// if err != nil {
							// 	return err
							// }
						}
					} else {
						return fmt.Errorf("Must pass flag\n")
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
