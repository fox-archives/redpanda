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
		Name:         name,
		Repos:        []Repo{},
		Transformers: []Transformer{},
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
	Name         string        `json:"name"`
	Repos        []Repo        `json:"repos"`
	Transformers []Transformer `json:"transformers"`
}

type Transformer struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

func (s *Store) TransformerAdd(transactionName string, typ string, name string, content string) error {
	found := false

	for i, transaction := range s.Transactions {
		if transaction.Name == transactionName {
			s.Transactions[i].Transformers = append(s.Transactions[i].Transformers, Transformer{
				Type:    typ,
				Name:    name,
				Content: content,
			})
			found = true
		}
	}

	if !found {
		return fmt.Errorf("Failed to find a transaction with that particular name")
	}

	return s.Save()
}

func (s *Store) TransformerRemove(transactionName string, transformerName string) error {
	foundTransformer := false
	foundRepo := false

	for i, t := range s.Transactions {
		if t.Name == transactionName {
			foundTransformer = true

			newTransformers := []Transformer{}
			for _, transformer := range t.Transformers {
				if transformer.Name == transformerName {
					foundRepo = true
					continue
				}

				newTransformers = append(newTransformers, transformer)
			}

			s.Transactions[i].Transformers = newTransformers
		}
	}

	if !foundRepo {
		return fmt.Errorf("Failed to find a repo with that particular name")
	}

	if !foundTransformer {
		return fmt.Errorf("Failed to find a transaction with that particular name")
	}

	return s.Save()
}

func (s *Store) TransformerEdit(transactionName string, transformerName string, newContent string) error {
	foundTransaction := false
	foundTransformer := false

	for i, t := range s.Transactions {
		if t.Name == transactionName {
			foundTransaction = true

			for j, transformer := range t.Transformers {
				if transformer.Name == transformerName {
					foundTransformer = true
					s.Transactions[i].Transformers[j] = Transformer{
						Type:    transformer.Type,
						Name:    transformer.Name,
						Content: newContent,
					}
					break
				}
			}
		}
	}

	if !foundTransformer {
		return fmt.Errorf("Failed to find a repo with that particular name")
	}

	if !foundTransaction {
		return fmt.Errorf("Failed to find a transaction with that particular name")
	}

	return s.Save()
}

func (s *Store) TransformerOrder(transformerName string, transformer string) error {
	return s.Save()
}

func (s *Store) RepoAdd(transactionName string, repoName string) error {
	found := false

	for i, t := range s.Transactions {
		if t.Name == transactionName {
			found = true
			s.Transactions[i].Repos = append(s.Transactions[i].Repos, Repo{
				Name:   repoName,
				Status: "uninitialized",
			})
		}
	}

	if !found {
		return fmt.Errorf("Failed to find a transaction with that particular name")
	}

	return s.Save()
}

func (s *Store) RepoRemove(transactionName string, repoName string) error {
	foundTransaction := false
	foundRepo := false

	for i, t := range s.Transactions {
		if t.Name == transactionName {
			foundTransaction = true

			newRepos := []Repo{}
			for _, repo := range t.Repos {
				if repo.Name == repoName {
					foundRepo = true
					continue
				}

				newRepos = append(newRepos, repo)
			}

			s.Transactions[i].Repos = newRepos
		}
	}

	if !foundRepo {
		return fmt.Errorf("Failed to find a repo with that particular name")
	}

	if !foundTransaction {
		return fmt.Errorf("Failed to find a transaction with that particular name")
	}

	return s.Save()
}

type Repo struct {
	Name   string
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
