package command

import (
	"strings"
)

type Command struct {
	Name string
	Args []string
}

// ParseCommand splits input into command + args
func ParseCommand(input string) (*Command, error) {
	tokens := strings.Fields(input)
	if len(tokens) == 0 {
		return nil, nil
	}
	cmd := &Command{
		Name: strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}
	return cmd, nil
}
