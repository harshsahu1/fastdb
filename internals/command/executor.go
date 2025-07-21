package command

import (
	"errors"
	"fastdb/internals/engine"
)

type Executor struct {
	Engine *engine.Engine
}

// NewExecutor creates a new command executor
func NewExecutor(e *engine.Engine) *Executor {
	return &Executor{Engine: e}
}

// ExecuteCommand dispatches command to engine
func (ex *Executor) ExecuteCommand(cmd *Command) (string, error) {
	switch cmd.Name {
	case "GET":
		if len(cmd.Args) != 1 {
			return "", errors.New("GET requires 1 argument")
		}
		val, ok := ex.Engine.Get(cmd.Args[0])
		if !ok {
			return "(nil)", nil
		}
		return string(val), nil

	case "SET":
		if len(cmd.Args) != 2 {
			return "", errors.New("SET requires 2 arguments")
		}
		ex.Engine.Set(cmd.Args[0], []byte(cmd.Args[1]))
		return "OK", nil

	case "DEL":
		if len(cmd.Args) != 1 {
			return "", errors.New("DEL requires 1 argument")
		}
		ex.Engine.Delete(cmd.Args[0])
		return "OK", nil
	case "SUBSCRIBE":
		if len(cmd.Args) != 1 {
			return "", errors.New("SUB requires 1 argument: key")
		}
		// Signal TCP handler that this is a live subscription
		return "__SUBSCRIBE__:" + cmd.Args[0], nil
	case "UNSUBSCRIBE":
		if len(cmd.Args) != 1 {
			return "", errors.New("UNSUB requires at least 1 key")
		}

		return "__UNSUBSCRIBE__:" + cmd.Args[0], nil
	default:
		return "", errors.New("unknown command: " + cmd.Name)
	}
}
