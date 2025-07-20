package main

import (
	"fastdb/internals/command"
	"fastdb/internals/engine"
	"fastdb/internals/tcp"
	"runtime"
)

func main() {
	db := engine.New(uint32(runtime.GOMAXPROCS(runtime.NumCPU())))
	ex := command.NewExecutor(db)
	// command.StartREPL(ex)
	tcp.Start(":6380", ex)
}
