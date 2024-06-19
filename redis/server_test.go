package redis

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestEchoServer(t *testing.T) {
	fmt.Println("Testing server")
	pingData := []byte("hello world" + "\n")
	go StartEchoServer()

	time.Sleep(2 * time.Second)
	conn, err := net.Dial("tcp", ":6379")
	defer conn.Close()

	if err != nil {
		t.Fatalf("Could not connect to server %v", err)
	}
	if _, err := conn.Write(pingData); err != nil {
		t.Fatalf("Could not send data %v", err)
	}
	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	fmt.Printf("PONG %v\n[%s]\n", buff[:n], string(buff[:n]))
	if err != nil {
		t.Fatalf("Error reading response %v", err)
	}
	if n != len(pingData) {
		t.Errorf("Sent and recieved length does not match %d != %d", n, len(pingData))
	}
	for i := range pingData {
		if buff[i] != pingData[i] {
			t.Fatalf("Response does not match")
		}
	}
}
