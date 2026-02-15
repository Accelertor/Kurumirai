package backend

import (
	"bufio"
	"bytes"
	"fmt"
)

func (c *Client) readPump(cs *Server) {
	defer func() {
		cs.unregister <- c
		c.conn.Close()
		fmt.Println("Disconnected:", c.conn.RemoteAddr())
	}()
	joinMsg := fmt.Appendf(nil, ">> %s joined.", c.conn.RemoteAddr().String())
	cs.broadcast <- Message{Sender: c, Payload: joinMsg}

	scanner := bufio.NewScanner(c.conn)

	for scanner.Scan() {
		rawMsg := scanner.Bytes()

		if bytes.Equal(rawMsg, []byte("/quit")) {
			break
		}
		formattedMsg := fmt.Appendf(nil, "[%s]: %s", c.conn.RemoteAddr(), string(rawMsg))
		cs.broadcast <- Message{Sender: c, Payload: formattedMsg}
	}

	leaveMsg := fmt.Appendf(nil, ">> %s left.", c.conn.RemoteAddr())
	cs.broadcast <- Message{Sender: c, Payload: leaveMsg}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for message := range c.send {
		_, err := c.conn.Write(append(message, '\n'))
		if err != nil {
			break
		}
	}
}
