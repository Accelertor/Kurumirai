package backend

import "net"

/* i couldn't decide if i should put this mf in a file or in connection.go so i flipped a coin 🥲*/
//probably in future this file will fill up woth data structure... unlike my life 🥲
type Client struct {
	sender net.Conn
	msg    chan string
}
