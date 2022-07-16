package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func New() Store {
	store := Store{
		Transactions: []Transaction{},
	}
	if err := initializeStore(&store); err != nil {
		log.Fatalln(err)
	}

	return store
}

type Store struct {
	Transactions []Transaction `json:"transactions"`
}

func (s *Store) TransactionGet(name string) (Transaction, error) {
	for _, t := range s.Transactions {
		if t.Name == name {
			return t, nil
		}
	}

	return Transaction{}, fmt.Errorf("A transaction with the specified name does not exist")
}

func (s *Store) TransactionAdd(name string) error {
	for _, t := range s.Transactions {
		if t.Name == name {
			return fmt.Errorf("A transaction with the specified name already exists")
		}
	}

	s.Transactions = append(s.Transactions, Transaction{
		Name:  name,
		Repos: []Repo{},
	})

	return s.Save()
}

func (s *Store) TransactionRemove(name string) error {
	idx := -1

	for i, t := range s.Transactions {
		if t.Name == name {
			idx = i
			break
		}
	}

	if idx == -1 {
		return fmt.Errorf("A transaction with the specified name does not exist")
	}

	s.Transactions = append(s.Transactions[:idx], s.Transactions[idx+1:]...)

	return s.Save()
}

func (s *Store) TransactionRename(oldName string, newName string) error {
	success := false

	for i, t := range s.Transactions {
		if t.Name == oldName {
			s.Transactions[i].Name = newName
			success = true
			break
		}
	}

	if !success {
		return fmt.Errorf("A transaction with the specified name does not exist")
	}

	return s.Save()
}

func (s *Store) TransactionList() []Transaction {
	return s.Transactions
}

type Transaction struct {
	Name  string `json:"name"`
	Repos []Repo `json:"repos"`
}

func (s *Store) ReposAdd(name string)    {}
func (s *Store) ReposRemove(name string) {}
func (s *Store) ReposList()              {}

type Repo struct {
	URL    string
	Dir    string
	Status string
}

func (s *Store) Save() error {
	home := os.Getenv("HOME")
	repoFile := filepath.Join(home, ".config", "redpanda", "data.json")
	err := os.MkdirAll(filepath.Dir(repoFile), 0o755)
	if err != nil {
		return err
	}

	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(repoFile, data, 0o644); err != nil {
		return err
	}

	return nil
}
func initializeStore(store *Store) error {
	home := os.Getenv("HOME") // TODO
	dataFile := filepath.Join(home, ".config", "redpanda", "data.json")
	content, err := ioutil.ReadFile(dataFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		} else {
			return err
		}
	}

	err = json.Unmarshal(content, store)
	if err != nil {
		return err
	}
	return nil
}
