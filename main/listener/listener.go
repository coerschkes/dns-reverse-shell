package main

import (
	"dns-reverse-shell/main/encoder"
	"dns-reverse-shell/main/protocol"
)

func main() {
	listener := protocol.NewDnsServer("8090", encoder.NewBase64Encoder())
	listener.Serve()
}
