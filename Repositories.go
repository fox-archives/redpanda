package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	git "github.com/libgit2/git2go/v33"
)

func utilForEachRepo(reposDir string, repos []string) []*git.Repository {
	var arr []*git.Repository

	for _, repoName := range repos {
		repoDir := filepath.Join(reposDir, repoName)
		gitRepo, err := git.OpenRepository(repoDir)
		if err != nil {
			panic(err)
		}

		arr = append(arr, gitRepo)
	}

	return arr
}

type Repos struct {
	List []string `json:"repos"`
}

func readRepos() (Repos, error) {
	home := os.Getenv("HOME")
	repoFile := filepath.Join(home, ".config", "redpanda", "repos.json")
	data, err := ioutil.ReadFile(repoFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Repos{}, nil
		} else {
			return Repos{}, err
		}
	}

	var repos Repos
	err = json.Unmarshal(data, &repos)
	if err != nil {
		return Repos{}, err
	}

	return repos, nil
}

func ReposList() []string {
	data, err := readRepos()
	if err != nil {
		panic(err)
	}

	return data.List
}
