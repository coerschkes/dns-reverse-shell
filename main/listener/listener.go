package main

import (
	"dns-reverse-shell/main/protocol/listener"
)

// todo: parameter port as command line arg
func main() {
	dnsListener := listener.NewDnsServer("8090")
	dnsListener.Initialize()
}
