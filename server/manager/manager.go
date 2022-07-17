package manager

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

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

				fmt.Printf("Cloning: %s\n", dir)
				cmd := exec.Command("git", "clone", url, dir)

				if err := cmd.Run(); err != nil {
					return err
				}

				m.store.Transactions[i].Repos[j].Status = "initialized"
			}
		}
	}

	return nil
}

func (m *Manager) IdempotentApply(transaction string) error {
	for _, transaction m.store.Transactions {
		
	}
}

func (m *Manager) Diff(transaction string) error {
	return nil
}
