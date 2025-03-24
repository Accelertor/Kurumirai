package main

import (
	"fmt"
	"net"
	"os"
	"sync"
)

type ChatServer struct {
	peers   sync.Map
	address string
	quit    chan struct{}
}

func NewChatServer(address string) *ChatServer {
	return &ChatServer{
		address: address,
		quit:    make(chan struct{}),
	}
}

func (cs *ChatServer) Start() {
	listener, err := net.Listen("tcp", cs.address)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer listener.Close()
	go handleShutdown(cs, listener)

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
		cs.peers.Store(conn.RemoteAddr().String(), conn)
		go handleConnection(cs, conn)
	}
}
