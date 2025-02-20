package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

const socketPath = "/tmp/unix_socket_example.sock"

type Request struct {
	Message string `json:"message"`
}

type Response struct {
	Reply string `json:"reply"`
}

func main() {
	if err := os.RemoveAll(socketPath); err != nil {
		fmt.Println("Error removing existing socket file:", err)
		return
	}

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		fmt.Println("Error creating Unix socket:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Unix socket server listening on", socketPath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	var req Request
	if err := decoder.Decode(&req); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	fmt.Println("Received:", req.Message)

	resp := Response{Reply: "Hello, " + req.Message}
	if err := encoder.Encode(resp); err != nil {
		fmt.Println("Error encoding JSON:", err)
	}
}
