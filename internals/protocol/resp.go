package protocol

import (
	"bufio"
	"bytes"
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

func ParseRESPCommandPartial(data []byte) (args []string, consumed int, err error) {
	if len(data) == 0 || data[0] != '*' {
		return nil, 0, nil 
	}

	lineEnd := bytes.Index(data, []byte("\r\n"))
	if lineEnd == -1 {
		return nil, 0, nil
	}

	numArgs, err := strconv.Atoi(string(data[1:lineEnd]))
	if err != nil {
		return nil, 0, fmt.Errorf("invalid array length: %s", string(data[1:lineEnd]))
	}

	pos := lineEnd + 2
	args = make([]string, 0, numArgs)

	for i := 0; i < numArgs; i++ {
		if len(data) < pos || data[pos] != '$' {
			return nil, 0, fmt.Errorf("expected bulk string")
		}

		lineEnd = bytes.Index(data[pos:], []byte("\r\n"))
		if lineEnd == -1 {
			return nil, 0, nil
		}

		argLen, err := strconv.Atoi(string(data[pos+1 : pos+lineEnd]))
		if err != nil {
			return nil, 0, fmt.Errorf("invalid bulk string length: %s", string(data[pos+1:pos+lineEnd]))
		}

		if len(data) < pos+lineEnd+2+argLen+2 {
			return nil, 0, nil
		}

		pos += lineEnd + 2
		args = append(args, string(data[pos:pos+argLen]))
		pos += argLen + 2
	}

	return args, pos, nil
}

func EncodeString(s string) []byte {
    return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(s), s))
}

func EncodeSimpleString(s string) []byte {
    return []byte(fmt.Sprintf("+%s\r\n", s))
}

func EncodeError(s string) []byte {
    return []byte(fmt.Sprintf("-%s\r\n", s))
}

var EncodeEmpty = []byte("+OK\r\n")

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