package main

import (
	"bufio"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

func handlerWhoAmI(id uint64, c net.Conn) {
	c.Write([]byte(strconv.FormatUint(id, 10) + "\n"))
}

func handlerWhoElse(id uint64, connMap *sync.Map, c net.Conn) {
	connMap.Range(func(key, value interface{}) bool {
		if key == id {
			return true
		}
		_, ok := value.(net.Conn)
		if !ok {
			fmt.Printf("Error sending message to %s!", key)
			connMap.Delete(key)
			return true
		}
		k, ok := key.(uint64)
		if !ok {
			fmt.Printf("Error retrieving Client ID for %s!", key)
			connMap.Delete(key)
			return true
		}
		c.Write([]byte(strconv.FormatUint(k, 10) + " "))
		return true
	})
	c.Write([]byte("\n"))
}

func parseSendCommand(input string) (string, map[uint64]bool, error) {
	inputList := strings.Fields(input)
	command := inputList[0]
	if command != "send" || len(inputList) < 3 {
		return "", nil, errors.New("Invalid command: " + command + "\n")
	}
	message, clientIds := inputList[1], inputList[2:]
	if len(clientIds) > 255 {
		return "", nil, errors.New("Clients exceeded. Expected less than 256\n")
	}
	clientMap := make(map[uint64]bool)
	errorGettingClientIds := false
	for _, clientId := range clientIds {
		clientId, err := strconv.ParseUint(clientId, 10, 64)
		if err != nil {
			errorGettingClientIds = true
			break
		}
		clientMap[clientId] = true
	}
	if errorGettingClientIds {
		return "", nil, errors.New("Invalid clientIds\n")
	}
	return message, clientMap, nil
}

func handlerSendCommand(input string, id uint64, connMap *sync.Map, c net.Conn) {
	message, clientMap, err := parseSendCommand(input)
	if err != nil {
		c.Write([]byte("Invalid command. Use `send [message] [clientID 1] [clientID 2] ... [clientID N]" + "\n"))
		return
	}

	connMap.Range(func(key, value interface{}) bool {
		if key == id {
			return true
		}
		k, ok := key.(uint64)
		if !ok {
			fmt.Printf("Error retrieving Client ID for %s!", key)
			connMap.Delete(key)
			return true
		}
		if _, ok := clientMap[k]; !ok {
			return true
		}
		conn, ok := value.(net.Conn)
		if !ok {
			fmt.Printf("Error sending message to %s!", key)
			connMap.Delete(key)
		}
		conn.Write([]byte(fmt.Sprintf("[FROM %d] %s\n", id, message)))
		return true
	})
}

func handleConnection(id uint64, connMap *sync.Map, c net.Conn) {
	defer func() {
		c.Close()
		connMap.Delete(id)
	}()
	fmt.Printf("Client %d connected!", id)
	for {
		input, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		input = strings.TrimSpace(string(input))
		switch input {
		case "whoami":
			handlerWhoAmI(id, c)
			continue
		case "whoelse":
			handlerWhoElse(id, connMap, c)
			continue
		default:
			handlerSendCommand(input, id, connMap, c)
			continue
		}
	}
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
