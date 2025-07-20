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
	id := conn.RemoteAddr().String()
	subscribedKeys := make(map[string]struct{}) // Track keys for cleanup

	defer func() {
		// Unsubscribe from all keys on disconnect
		for key := range subscribedKeys {
			executor.Engine.PubSub().Unsubscribe(key, id)
		}
	}()


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

		// Handle SUBSCRIBE mode
		if strings.HasPrefix(response, "__SUB__:") {
			key := strings.TrimPrefix(response, "__SUB__:")

			sub := executor.Engine.PubSub().Subscribe(key, conn.RemoteAddr().String())
			conn.Write([]byte("SUBSCRIBED to " + key + "\n"))

			// Listen to channel and write updates to TCP client
			go func() {
				for msg := range sub.Chan {
					conn.Write([]byte("[" + key + "] " + string(msg) + "\n"))
				}
			}()

			continue // wait for more input or pub messages
		}

		// Handle UNSUBSCRIBE mode
		if strings.HasPrefix(response, "__UNSUB__:") {
			key := strings.TrimPrefix(response, "__UNSUB__:")

			id := conn.RemoteAddr().String()
			executor.Engine.PubSub().Unsubscribe(key, id)
			delete(subscribedKeys, key)
			conn.Write([]byte("UNSUBSCRIBED from " + key + "\n"))
			continue
		}


		conn.Write([]byte(response + "\n"))

	}
}
