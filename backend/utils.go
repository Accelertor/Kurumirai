package backend

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func handleShutdown(cs *ChatServer, listener net.Listener) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	fmt.Println("\nShutting down server...")
	listener.Close()
	cs.peers.Range(func(_, value any) bool {
		conn := value.(net.Conn)
		conn.Close()
		return true
	})
	close(cs.quit)
	os.Exit(0)
}
