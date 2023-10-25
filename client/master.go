package main

import (
	"dns-shellcode/encoder"
	"fmt"
	"github.com/miekg/dns"
)

func main() {
	command := encoder.Encode("ls")
	resolve(command, dns.TypeA)
}

func resolve(domain string, qtype uint16) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), qtype)
	m.RecursionDesired = true

	c := new(dns.Client)
	in, _, err := c.Exchange(m, "127.0.0.1:8090")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, ans := range in.Answer {
		fmt.Println(encoder.Decode(ans.Header().Name))
	}
}
