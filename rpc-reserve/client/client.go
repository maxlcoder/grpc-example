package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

func main() {
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("Listen TCP error:", err)
	}
	clientChan := make(chan *rpc.Client)
	go func() {
		for  {
			conn, err := l.Accept()
			if err != nil {
				log.Fatal("Accept error:", err)
			}
			clientChan <- rpc.NewClient(conn)

		}
	}()

	doClientWork(clientChan)
}

func doClientWork(clientChan <-chan *rpc.Client) {
	client := <-clientChan
	defer client.Close()

	var reply string
	err := client.Call("HelloService.Hello", "hello", &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)
}
