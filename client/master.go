package main

import (
	"bufio"
	"dns-shellcode/encoder"
	"fmt"
	"github.com/miekg/dns"
	"os"
)

func main() {
	readInput()
}

func readInput() {
	fmt.Println("Enter command. Empty string exits the program")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		text := scanner.Text()
		if len(text) != 0 {
			command := encoder.Encode(text)
			resolve(command, dns.TypeA)
		} else {
			break
		}
	}
	if scanner.Err() != nil {
		fmt.Println("Error: ", scanner.Err())
	}
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
