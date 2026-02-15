package backend

import (
	"net"
	"sync/atomic"
)

/* i couldn't decide if i should put this mf in a file or in connection.go so i flipped a coin 🥲*/
//probably in future this file will fill up woth data structure... unlike my life 🥲

// remember to make struct 64 bytes(512 bit) long. Cuz cpu cAche line L1 etc. Google it
type Client struct {
	conn net.Conn
	send chan []byte
}

// you're not alone any more
type Message struct {
	Sender  *Client
	Payload []byte
}
type Server struct {
	peers            map[*Client]bool
	boardCastList    []*Client
	broadcast        chan Message
	address          string
	register         chan *Client
	unregister       chan *Client
	quit             chan struct{}
	TotalConnections atomic.Int64
}
