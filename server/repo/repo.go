package repos

import (
	"github.com/hyperupcall/redpanda/server/store"
)

// Stub
type Changeset interface{}

type Repo struct {
	RemoteURL string
	LocalDir  string
	status    string
}

func (r *Repo) ResetFiles() {

}

func (r *Repo) ApplyChangeset(changeset Changeset) {

}

func (r *Repo) Commit() {

}

func (r *Repo) Push() {

}

func (r *Repo) Pull() {

}

type Repos struct {
	store store.Store
	repos []string
}
