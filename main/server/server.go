package main

import (
	"dns-shellcode/main/encoder"
	"dns-shellcode/main/protocol"
)

func main() {
	dnsServer := protocol.NewDnsServer("8090", encoder.NewBase64Encoder())
	dnsServer.Serve()
}

//todo create shell impl -> have command buffer/stack? prefix? -> exec command "cd .., ls" for instance
// => currently working: "cd .. && ls"
//todo: clear stack
//todo: navigation stack -> only for navigation!
//
//
//todo: implement msg splitter
// => current problem: header size too big, reicv fails with header overflow
//todo: tcp instead of udp?

//todo: revert master/slave?
