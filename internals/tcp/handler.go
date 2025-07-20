// internals/tcp/handler.go
package tcp

import (
	"bufio"
	"fastdb/internals/command"
	"net"
	"strings"
)

func handleConnection(conn net.Conn, executor *command.Executor) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		conn.Write([]byte("> "))
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.EqualFold(line, "exit") || strings.EqualFold(line, "quit") {
			conn.Write([]byte("👋 Bye!\n"))
			return // Close the connection
		}

		args := strings.Split(line, " ")
		cmd := strings.ToUpper(args[0])
		response, err := executor.ExecuteCommand(&command.Command {
			Name: cmd,
			Args: args[1:],
		})
		if err != nil {
			conn.Write([]byte("ERR " + err.Error() + "\n"))
			continue
		}
		conn.Write([]byte(response + "\n"))

	}
}
