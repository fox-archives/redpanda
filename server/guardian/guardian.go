package manager

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/hyperupcall/redpanda/server/logger"
	"github.com/hyperupcall/redpanda/server/store"
)

func New(store *store.Store) Guardian {
	l := logger.New("redpanda.log")

	return Guardian{
		store:  store,
		logger: &l,
	}
}

func forEachRepo(m Guardian, transactionName string, fn func(repo store.Repo) error) error {
	found := false

	for _, transaction := range m.store.Transactions {
		if transaction.Name == transactionName {
			found = true

			originalDir, err := os.Getwd()
			if err != nil {
				return err
			}

			for _, repo := range transaction.Repos {
				if err := os.Chdir(repo.Dir); err != nil {
					return err
				}

				if err = fn(repo); err != nil {
					return err
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

func RepoIsCloned(dir string) (error, bool) {
	infos, err := ioutil.ReadDir(dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, false
	} else if err != nil {
		return err, false
	}

	if len(infos) == 0 {
		return nil, false
	} else {
		return nil, true
	}
}

func gitDiff(g *Guardian, transactionName string) (string, error) {
	contents := ""
	if err := g.forEachRepoInTransactionCd(transactionName, func(transaction *store.Transaction, repo *store.Repo) error {
		cmd := exec.Command("git", "diff", "--staged")
		content, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}

		contents = contents + string(content)

		return nil
	}); err != nil {
		return "", err
	}

	return contents, nil
}

func gitReset(g *Guardian, transactionName string) (string, error) {
	contents := ""
	if err := g.forEachRepoInTransactionCd(transactionName, func(transaction *store.Transaction, repo *store.Repo) error {
		cmd := exec.Command("git", "reset", "--hard", "HEAD")
		content, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}

		contents = contents + string(content)

		return nil
	}); err != nil {
		return "", err
	}

	return contents, nil
}

func executeModifiers(g *Guardian, transactionName string) error {
	if err := g.forEachRepoInTransactionCd(transactionName, func(transaction *store.Transaction, repo *store.Repo) error {
		cmd := exec.Command("git", "merge-base", "origin", "HEAD")
		_, err := cmd.Output()
		if err != nil {
			return err
		}

		for _, former := range transaction.Transformers {
			if former.Type == "command" {
				fileName := "/tmp/redpanda-script.sh"

				if err := ioutil.WriteFile(fileName, []byte(former.Content), 0o755); err != nil {
					return err
				}

				cmd := exec.Command("bash", fileName)
				if err := cmd.Run(); err != nil {
					return err
				}

				// l := strings.Split(former.Content, " ")

				// cmd := exec.Command(l[0], l...)
				// if err := cmd.Run(); err != nil {
				// 	return err
				// }

				cmd = exec.Command("git", "add", "-A")
				if err := cmd.Run(); err != nil {
					return err
				}
			} else {
				panic("Unknown type")
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

type Guardian struct {
	store  *store.Store
	logger logger.Logger
}

func (g *Guardian) forEachRepoInTransaction(transactionName string, fn func(transaction *store.Transaction, repo *store.Repo) error) error {
	foundTransaction := false

	for i := range g.store.Transactions {
		transaction := g.store.Transactions[i]

		if transaction.Name == transactionName {
			foundTransaction = true
			for j := range transaction.Repos {
				repo := transaction.Repos[j]

				if err := fn(&transaction, &repo); err != nil {
					return err
				}
			}
		}
	}

	if !foundTransaction {
		return fmt.Errorf("Failed to find transaction with a name of %s", transactionName)
	}

	return nil
}

func (g *Guardian) forEachRepoInTransactionCd(transactionName string, fn func(transaction *store.Transaction, repo *store.Repo) error) error {
	g.forEachRepoInTransaction(transactionName, func(transaction *store.Transaction, repo *store.Repo) error {
		originalDir, err := os.Getwd()
		if err != nil {
			return err
		}

		if err = os.Chdir(repo.Dir); err != nil {
			return err
		}

		if err := fn(transaction, repo); err != nil {
			return err
		}

		if err = os.Chdir(originalDir); err != nil {
			return err
		}

		return nil
	})
	return nil
}

func (g *Guardian) ActionApply(transactionName string) (string, error) {
	if err := g.forEachRepoInTransaction(transactionName, func(transaction *store.Transaction, repo *store.Repo) error {
		if repo.Status == "uninitialized" {
			g.logger.Info("Initializing " + repo.Name)

			url := fmt.Sprintf("git@github.com:%s", repo.Name)
			dir := filepath.Join(os.Getenv("HOME"), ".local", "share", "redpanda", "downloads", repo.Name)

			repo.URL = url
			repo.Dir = dir

			err, isCloned := RepoIsCloned(dir)
			if err != nil {
				return err
			}

			if !isCloned {
				cmd := exec.Command("git", "clone", url, dir)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				if err := cmd.Run(); err != nil {
					return err
				}
			}

			repo.Status = "initialized"
		}

		return nil
	}); err != nil {
		return "", err
	}

	if _, err := gitReset(g, transactionName); err != nil {
		return "", err
	}

	if err := executeModifiers(g, transactionName); err != nil {
		return "", err
	}

	return gitDiff(g, transactionName)
}

func (g *Guardian) ActionRefresh(transactionName string) (string, error) {
	if err := g.forEachRepoInTransactionCd(transactionName, func(transaction *store.Transaction, repo *store.Repo) error {
		if err := os.Chdir(repo.Dir); err != nil {
			fmt.Println(repo.Dir)
			return err
		}

		g.logger.Trace("git fetch: " + repo.Name)
		cmd := exec.Command("git", "fetch", "origin")
		if err := cmd.Run(); err != nil {
			return err
		}

		g.logger.Trace("git merge-base: " + repo.Name)
		cmd = exec.Command("git", "merge-base", "origin", "HEAD")
		mergeBase, err := cmd.Output()
		if err != nil {
			return err
		}

		g.logger.Trace("git reset: " + repo.Name)
		cmd = exec.Command("git", "-C", repo.Dir, "reset", "--hard", strings.TrimSpace(string(mergeBase)))
		if err := cmd.Run(); err != nil {
			return err
		}

		g.logger.Trace("git pull: " + repo.Name)
		cmd = exec.Command("git", "pull", "origin")
		if err := cmd.Run(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return "", err
	}

	if _, err := gitReset(g, transactionName); err != nil {
		return "", err
	}

	if err := executeModifiers(g, transactionName); err != nil {
		return "", err
	}

	return gitDiff(g, transactionName)
}

func (m *Guardian) ActionCommit(transactionName string) error {
	return nil
}

func (m *Guardian) ActionPush(transactionName string) error {
	return nil
}

func (m *Guardian) Setup(transactionName string) error {

	return nil
}

func (m *Guardian) Commit(transactionName string, commitMessage string) (string, error) {
	transactionRepoDir := filepath.Join(os.Getenv("HOME"), ".local", "share", "redpanda", "transaction-repo")

	err, isCloned := RepoIsCloned(transactionRepoDir)
	if err != nil {
		return "", err
	}
	if !isCloned {
		cmd := exec.Command("git", "clone", "https://github.com/hyperupcall/transactions", transactionRepoDir)
		if err := cmd.Run(); err != nil {
			return "", err
		}
	}

	config := struct {
		gitAuthor string
		gpgId     string
	}{
		gitAuthor: "Captain Woofers <99463792+captain-woofers@users.noreply.github.com>",
		gpgId:     "0xF1BBE0168CC63A97",
	}
	id := uuid.NewString()

	os.Setenv("GIT_COMMITTER_NAME", "Captain Woofers")
	os.Setenv("GIT_COMMITTER_EMAIL", "99463792+captain-woofers@users.noreply.github.com")

	if err := forEachRepo(*m, transactionName, func(repo store.Repo) error {
		message := commitMessage + `

Transaction-Id: ` + id + ``
		fmt.Println(message)
		cmd := exec.Command("git", "commit", "--allow-empty", "-m", message, "--author", config.gitAuthor, "--gpg-sign="+config.gpgId)
		content, err := cmd.CombinedOutput()
		fmt.Println(string(content))
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return "{}", err
	}

	// Now, make a record in transactions
	transactionMessage := ""
	cmd := exec.Command("git", "-C", transactionRepoDir, "add", "-m", transactionMessage)
	if err := cmd.Run(); err != nil {
		return "", err
	}

	transactionMessage := ""
	cmd := exec.Command("git", "-C", transactionRepoDir, "commit", "-m", transactionMessage)
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return "{}", nil
}

func (m *Guardian) Push(transactionName string) (string, error) {
	err := forEachRepo(*m, transactionName, func(repo store.Repo) error {
		cmd := exec.Command("git", "push", "origin")
		content, err := cmd.CombinedOutput()
		fmt.Println(string(content))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "{}", err
	}

	transactionRepoDir := filepath.Join(os.Getenv("HOME"), ".local", "share", "redpanda", "transaction-repo")
	cmd := exec.Command("git", "-C", transactionRepoDir, "push", "origin")
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return "{}", nil
}
