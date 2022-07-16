# shellcheck shell=bash

task.run() { 
	bake.cfg 'big-print' 'no'

	go run . "$@"
}