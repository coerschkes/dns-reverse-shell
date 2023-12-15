package main

import (
	"dns-reverse-shell/main/protocol/listener"
	"os"
)

func main() {
	params := os.Args[1:]
	if len(params) != 1 {
		panic("usage: <port>")
	}
	dnsListener := listener.NewDnsServer(params[0])
	dnsListener.Initialize()
}
