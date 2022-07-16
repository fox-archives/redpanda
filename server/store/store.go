package store

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hyperupcall/redpanda/server/util"
)

func initializeStore(store *Store) error {
	home := os.Getenv("HOME") // TODO
	dataFile := filepath.Join(home, ".config", "redpanda", "data.json")
	data, err := ioutil.ReadFile(dataFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		} else {
			return err
		}
	}

	err = json.Unmarshal(data, store)
	if err != nil {
		return err
	}
	return nil
}

func New() Store {
	var store Store
	if err := initializeStore(&store); err != nil {
		log.Fatalln(err)
	}

	return store
}

type Store struct {
	Repos []string `json:"repos"`
}

func (s *Store) RepoList() []string {
	return s.Repos
}

func (s *Store) RepoAdd(repos []string) error {
	for _, repo := range repos {
		c, _ := util.Contains(s.Repos, repo)
		if !c {
			s.Repos = append(s.Repos, repos...)
		}
	}

	return s.Save()
}

func (s *Store) RepoRemove(repos []string) error {
	for _, repo := range repos {
		c, i := util.Contains(s.Repos, repo)
		if c {
			s.Repos = append(s.Repos[:i], s.Repos[i+1:]...)
		}
	}

	return s.Save()
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
