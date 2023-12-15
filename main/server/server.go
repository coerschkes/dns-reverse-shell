package main

import (
	"dns-reverse-shell/main/protocol/server"
	"os"
)

func main() {
	params := os.Args[1:]
	if len(params) != 1 {
		panic("usage: <port>")
	}
	server.NewDnsServer(params[0]).Initialize()
}
