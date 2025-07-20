package main

import (
	"fastdb/internals/command"
	"fastdb/internals/engine"
	"fastdb/internals/tcp"
)

func main() {
	db := engine.New(256)
	ex := command.NewExecutor(db)
	// command.StartREPL(ex)
	tcp.Start(":6380", ex)
}
