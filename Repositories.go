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

func writeRepos(repos Repos) error {
	home := os.Getenv("HOME")
	repoFile := filepath.Join(home, ".config", "redpanda", "repos.json")
	err := os.MkdirAll(filepath.Dir(repoFile), 0o755)
	if err != nil {
		return err
	}

	data, err := json.Marshal(repos)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(repoFile, data, 0o644); err != nil {
		return err
	}

	return nil
}

func ReposAdd(repos []string) {
	data, err := readRepos()
	if err != nil {
		panic(err)
	}

	for _, repo := range repos {
		c, _ := contains(data.List, repo)
		if !c {
			data.List = append(data.List, repos...)
		}
	}

	err = writeRepos(data)
	if err != nil {
		panic(err)
	}
}

func contains(s []string, str string) (bool, int) {
	for i, v := range s {
		if v == str {
			return true, i
		}
	}

	return false, -1
}

func ReposRemove(repos []string) {
	data, err := readRepos()
	if err != nil {
		panic(err)
	}

	for _, repo := range repos {
		c, i := contains(data.List, repo)
		if c {
			data.List = append(data.List[:i], data.List[i+1:]...)
		}
	}

	err = writeRepos(data)
	if err != nil {
		panic(err)
	}
}

func ReposList() []string {
	data, err := readRepos()
	if err != nil {
		panic(err)
	}

	return data.List
}
