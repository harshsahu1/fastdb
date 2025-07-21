package tcp

import (
	"bufio"
	"fastdb/internals/command"
	"fastdb/internals/engine"
	"fastdb/internals/protocol"
	"fmt"
	"net"
)

func handleConnection(conn net.Conn, executor *command.Executor) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	id := conn.RemoteAddr().String()
	subscribedKeys := make(map[string]struct{})

	writeChan := make(chan []byte, 128) // Buffered channel for all writes

	// Only one goroutine writes to the connection
	go func() {
		for msg := range writeChan {
			conn.Write(msg)
		}
	}()

	defer func() {
		for key := range subscribedKeys {
			executor.Engine.PubSub().Unsubscribe(key, id)
		}
		close(writeChan)
	}()

	for {
		args, err := protocol.ParseRESPCommand(reader)
		if err != nil {
			writeChan <- []byte("-ERR " + err.Error() + "\r\n")
			return
		}

		if len(args) == 0 {
			writeChan <- []byte("-ERR Empty command\r\n")
			continue
		}

		cmd := args[0]
		response, err := executor.ExecuteCommand(&command.Command{
			Name: cmd,
			Args: args[1:],
		})
		if err != nil {
			writeChan <- []byte("-ERR " + err.Error() + "\r\n")
			continue
		}

		// Handle SUBSCRIBE
		if len(response) > 14 && response[:14] == "__SUBSCRIBE__:" {
			key := response[14:]
			sub := executor.Engine.PubSub().Subscribe(key, id)
			subscribedKeys[key] = struct{}{}

			subResp := fmt.Sprintf("*3\r\n$9\r\nsubscribe\r\n$%d\r\n%s\r\n:%d\r\n", len(key), key, 1)
			writeChan <- []byte(subResp)

			// Push messages from pubsub through writeChan (NOT conn directly)
			go func(key string, sub *engine.Subscriber) {
				for msg := range sub.Chan {
					pubMsg := fmt.Sprintf("*3\r\n$7\r\nmessage\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n",
						len(key), key, len(msg), msg)
					writeChan <- []byte(pubMsg)
				}
			}(key, sub)

			continue
		}

		// Handle UNSUBSCRIBE
		if len(response) > 16 && response[:16] == "__UNSUBSCRIBE__:" {
			key := response[16:]
			executor.Engine.PubSub().Unsubscribe(key, id)
			delete(subscribedKeys, key)

			unsubResp := fmt.Sprintf("*3\r\n$11\r\nunsubscribe\r\n$%d\r\n%s\r\n:%d\r\n",
				len(key), key, 0)
			writeChan <- []byte(unsubResp)
			continue
		}

		// Default reply
		writeChan <- []byte("+OK " + response + "\r\n")
	}
}
