package protocol

import (
	"dns-shellcode/main/encoder"
	"fmt"
	"github.com/miekg/dns"
)

type DNSClient struct {
	address string
	encoder encoder.StringEncoder
}

func NewDNSClient(address string, encoder encoder.StringEncoder) *DNSClient {
	return &DNSClient{address: address, encoder: encoder}
}

func (d DNSClient) SendMessage(message string) {
	m := d.createMessage(d.encoder.Encode(message))

	c := new(dns.Client)
	in, _, err := c.Exchange(m, d.address)
	if err != nil {
		fmt.Println(err)
		return
	}
	d.handleAnswer(in)
}

func (d DNSClient) createMessage(message string) *dns.Msg {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(message), dns.TypeA)
	m.RecursionDesired = true
	return m
}

func (d DNSClient) handleAnswer(answerMsg *dns.Msg) {
	for _, ans := range answerMsg.Answer {
		fmt.Println(d.encoder.Decode(ans.Header().Name))
	}
}