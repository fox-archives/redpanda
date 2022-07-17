package main

import (
	"github.com/hyperupcall/redpanda/server/serve"
	"github.com/hyperupcall/redpanda/server/store"
)

func main() {
	store := store.New()
	serve.Serve(&store)
}
