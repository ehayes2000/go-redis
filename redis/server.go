package redis

import (
	"fmt"
	"io"
	"net"
	"time"
)

func echoHandler(c net.Conn) error {
	defer c.Close()
	c.SetReadDeadline(time.Now().Add(5 * time.Second))
	buff := make([]byte, 1024)
	var response []byte
	for {
		n, err := c.Read(buff)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		response = append(response, buff[:n]...)
		if buff[n-1] == '\n' {
			c.Write(response)
			break
		}
	}
	return nil
}

func StartEchoServer() error {
	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		fmt.Println("MAKE CONNECTION")
		if err != nil {
			return err
		}
		go echoHandler(conn)
	}
	return nil
}
