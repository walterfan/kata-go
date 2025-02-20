package main

import (
	"encoding/json"
	"fmt"
	"net"
)

const socketPath = "/tmp/unix_socket_example.sock"

type Request struct {
	Message string `json:"message"`
}

type Response struct {
	Reply string `json:"reply"`
}

func main() {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		fmt.Println("Error connecting to Unix socket:", err)
		return
	}
	defer conn.Close()

	req := Request{Message: "Walter"}
	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)

	if err := encoder.Encode(req); err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	var resp Response
	if err := decoder.Decode(&resp); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	fmt.Println("Received response:", resp.Reply)
}
