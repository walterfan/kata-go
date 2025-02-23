package main

type Request struct {
	Message string `json:"message"`
}

type Response struct {
	Reply string `json:"reply"`
}
