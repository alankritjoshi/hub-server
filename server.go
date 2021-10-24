package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"sync"
)

func handleConnection(id uint64, connMap *sync.Map, c net.Conn) {
	defer func() {
		c.Close()
		connMap.Delete(id)
	}()
	fmt.Printf("Client %d connected!", id)
}

func generateClientId() uint64 {
	buf := make([]byte, 8)
	rand.Read(buf)
	return binary.LittleEndian.Uint64(buf)
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer l.Close()

	var connMap = &sync.Map{}

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		clientId := generateClientId()
		connMap.Store(clientId, c)
		go handleConnection(clientId, connMap, c)
	}
}
