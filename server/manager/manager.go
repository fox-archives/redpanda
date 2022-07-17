package manager

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hyperupcall/redpanda/server/store"
)

func New(store *store.Store) Manager {
	var manager Manager
	manager.store = store

	return manager
}

type Manager struct {
	store *store.Store
}

func (m *Manager) Initialize() error {
	for i, transaction := range m.store.Transactions {
		fmt.Println("-- " + transaction.Name)
		for j, repo := range transaction.Repos {
			fmt.Println(repo.Status)
			if repo.Status == "uninitialized" {
				fmt.Printf("transaction %s: initializing %s\n", transaction.Name, repo.Name)

				url := fmt.Sprintf("git@github.com:%s", repo.Name)
				dir := filepath.Join(os.Getenv("HOME"), ".local", "share", "redpanda", "downloads", repo.Name)

				m.store.Transactions[i].Repos[j].URL = url
				m.store.Transactions[i].Repos[j].Dir = dir

				_, err := ioutil.ReadDir(dir)
				if errors.Is(err, os.ErrNotExist) {
					fmt.Printf("Cloning: %s to %s\n", url, dir)
					cmd := exec.Command("git", "clone", url, dir)
					cmd.Stdin = os.Stdin
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr

					if err := cmd.Run(); err != nil {
						return err
					}

				} else {

				}

				m.store.Transactions[i].Repos[j].Status = "initialized"

			}
		}
	}

	return nil
}

func (m *Manager) IdempotentApply(transactionName string) error {
	found := false

	for _, transaction := range m.store.Transactions {
		if transaction.Name == transactionName {
			found = true

			originalDir, err := os.Getwd()
			if err != nil {
				return err
			}

			fmt.Printf("Running transaction: %s\n", transaction.Name)
			fmt.Printf("  repos: %+v\n", transaction.Repos)
			for _, repo := range transaction.Repos {
				if err := os.Chdir(repo.Dir); err != nil {
					return err
				}

				cmd := exec.Command("git", "-C", repo.Dir, "reset", "--hard", "HEAD")
				if err := cmd.Run(); err != nil {
					return err
				}

				for _, former := range transaction.Transformers {
					if former.Type == "command" {
						l := strings.Split(former.Content, " ")

						cmd := exec.Command(l[0], l...)
						if err := cmd.Run(); err != nil {
							return err
						}

						cmd = exec.Command("git", "add", "-A")
						if err := cmd.Run(); err != nil {
							return err
						}
					} else if former.Type == "regex" {
						return fmt.Errorf("Type regex not supported (on server side)")
					}
				}
			}

			if err = os.Chdir(originalDir); err != nil {
				return err
			}

		}
	}

	if !found {
		return fmt.Errorf("Failed to find a transaction with that particular name")
	}

	return nil
}

func (m *Manager) Diff(transactionName string) error {
	found := false

	for _, transaction := range m.store.Transactions {
		if transaction.Name == transactionName {
			found = true

			originalDir, err := os.Getwd()
			if err != nil {
				return err
			}

			fmt.Printf("----------- DIFFFFF FORRRR: %s\n", transaction.Name)
			for _, repo := range transaction.Repos {
				if err := os.Chdir(repo.Dir); err != nil {
					return err
				}

				cmd := exec.Command("git", "diff", "--staged")
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Start()
			}

			if err = os.Chdir(originalDir); err != nil {
				return err
			}

		}
	}

	if !found {
		return fmt.Errorf("Failed to find a transaction with that particular name")
	}

	return nil
}
