package main

import (
	"fastdb/internals/command"
	"fastdb/internals/engine"
)

func main() {
	db := engine.New(256, nil)
	ex := command.NewExecutor(db)
	command.StartREPL(ex)
}
