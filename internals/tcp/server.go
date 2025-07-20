package tcp

import (
	"fastdb/internals/command"
	"log"
	"net"
)

func Start(address string, executor *command.Executor) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("TCP Server: Failed to listen on %s: %v", address, err)
	}
	log.Printf("🚀 FastDB TCP server started on %s\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting: %v", err)
			continue
		}
		go handleConnection(conn, executor)
	}
}