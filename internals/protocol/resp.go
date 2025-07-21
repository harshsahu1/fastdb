package protocol

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func ParseRESPCommand(r *bufio.Reader) ([]string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = strings.TrimSpace(line)

	if !strings.HasPrefix(line, "*") {
		return nil, errors.New("expected array (*)")
	}

	numArgs, err := strconv.Atoi(line[1:])
	if err != nil {
		return nil, errors.New("invalid array length")
	}

	args := make([]string, 0, numArgs)
	for i := 0; i < numArgs; i++ {
		lenLine, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		if !strings.HasPrefix(lenLine, "$") {
			return nil, errors.New("expected bulk string ($)")
		}

		argLen, err := strconv.Atoi(strings.TrimSpace(lenLine[1:]))
		if err != nil {
			return nil, fmt.Errorf("invalid length in bulk string: %v", err)
		}

		arg := make([]byte, argLen+2) // +2 for \r\n
		if _, err := io.ReadFull(r, arg); err != nil {
			return nil, err
		}
		args = append(args, string(arg[:argLen])) // exclude \r\n
	}

	return args, nil
}

// Example:
// *3\r\n
// $3\r\n
// SET\r\n
// $4\r\n
// name\r\n
// $5\r\n
// harsh\r\n

// Step-by-step:
// *3 → Array of 3 elements
// $3 → Next arg is 3 bytes → SET
// $4 → Next arg is 4 bytes → name
// $5 → Next arg is 5 bytes → harsh
// Final Result: []string{"SET", "name", "harsh"}