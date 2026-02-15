package backend

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func NewServer(address string) *Server {
	return &Server{
		peers:         make(map[*Client]bool),
		boardCastList: make([]*Client, 0),
		broadcast:     make(chan Message, 256),
		address:       address,
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		quit:          make(chan struct{}),
	}
}

func (cs *Server) Run() {
	for {
		select {
		case client := <-cs.register:
			cs.peers[client] = true
			cs.TotalConnections.Add(1)
			cs.rebuildBroadcastList()

		case client := <-cs.unregister:
			if _, ok := cs.peers[client]; ok {
				delete(cs.peers, client)
				close(client.send)
				cs.TotalConnections.Add(-1)
				cs.rebuildBroadcastList()
			}

		case msg := <-cs.broadcast:
			laggingClients := false

			for _, client := range cs.boardCastList {
				if client == msg.Sender {
					continue
				}

				select {
				case client.send <- msg.Payload:
				default:
					close(client.send)
					delete(cs.peers, client)
					cs.TotalConnections.Add(-1)
					laggingClients = true
				}
			}
			if laggingClients {
				cs.rebuildBroadcastList()
			}

		case <-cs.quit:
			for client := range cs.peers {
				client.conn.Close()
			}
			return
		}
	}
}

func (cs *Server) rebuildBroadcastList() {
	newList := make([]*Client, 0, len(cs.peers))
	for client := range cs.peers {
		newList = append(newList, client)
	}
	cs.boardCastList = newList
}

func (cs *Server) Start() {
	listener, err := net.Listen("tcp", cs.address)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer listener.Close()

	go handleShutdown(cs, listener)
	go cs.Run()

	fmt.Println("Server listening on", cs.address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-cs.quit:
				return
			default:
				fmt.Println("Connection error:", err)
			}
			continue
		}

		fmt.Println("Connected:", conn.RemoteAddr())
		client := &Client{
			conn: conn,
			send: make(chan []byte, 256),
		}
		cs.register <- client

		go client.readPump(cs)
		go client.writePump()
	}
}

func handleShutdown(cs *Server, listener net.Listener) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nShutting down server...")
	listener.Close()
	close(cs.quit)
}
