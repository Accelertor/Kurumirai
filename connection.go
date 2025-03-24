package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func handleConnection(cs *ChatServer, conn net.Conn) {
	defer func() {
		cs.peers.Delete(conn.RemoteAddr().String())
		informPeers(cs, fmt.Sprintf("%s has left the chat.", conn.RemoteAddr()))
		conn.Close()
		fmt.Println("Disconnected:", conn.RemoteAddr())
	}()

	scanner := bufio.NewScanner(conn)
	writer := bufio.NewWriter(conn)

	informPeers(cs, fmt.Sprintf("%s has joined the chat.", conn.RemoteAddr()))

	for scanner.Scan() {
		msg := strings.TrimSpace(scanner.Text())
		if msg == "/quit" {
			fmt.Fprintln(writer, "Goodbye!")
			writer.Flush()
			break
		}
		cs.broadcast(msg, conn)
	}
}

// informPeers sends a message to all users except the sender
func informPeers(cs *ChatServer, message string) {
	cs.peers.Range(func(_, value any) bool {
		peer := value.(net.Conn)
		fmt.Fprintln(peer, message)
		return true
	})
}
