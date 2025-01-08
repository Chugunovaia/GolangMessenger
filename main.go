package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello\n If you want to start server, please, write 0. If you want to start client, write 1\n Client won't be able to work until you start the server\n")
	s := ""
	fmt.Scan(&s)
	for s != "1" && s != "0" {
		fmt.Println("Please, write 0 for server or 1 for client:\n ")
		fmt.Scan(&s)
	}
	if s == "1" {
		client()
	} else if s == "0" {
		server()
	}
	return
}
