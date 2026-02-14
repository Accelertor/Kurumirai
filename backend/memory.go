package backend

import (
	"net"
	"sync"
)

var arena = sync.Pool{
	New: func() interface{} {
		buf := make([]byte, 4096)
		return &buf
	},
}

func (cs *ChatServer) broadcast(message string, sender net.Conn) {
	bufPtr := arena.Get().(*[]byte)
	defer arena.Put(bufPtr)

	// Copy message into buffer without append()
	msgLen := copy(*bufPtr, message+"\n")
	msgBytes := (*bufPtr)[:msgLen] // Slice only the relevant part

	cs.peers.Range(func(_, value any) bool {
		peer := value.(net.Conn)
		if peer != sender {
			_, _ = peer.Write(msgBytes) // Direct zero-copy write
		}
		return true
	})
}
