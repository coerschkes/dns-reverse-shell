package protocol

import (
	"github.com/miekg/dns"
	"strings"
)

const maxMessageSliceLen = 59

type MessageSplitter interface {
	Split(string) []dns.RR
	Collect([]dns.RR) string
}

type SimpleMessageSplitter struct{}

func NewSimpleMessageSplitter() *SimpleMessageSplitter {
	return &SimpleMessageSplitter{}
}

func (s SimpleMessageSplitter) Split(message string) []dns.RR {
	splitMessage := strings.Split(s.markSliceBoundaries(message), ".")
	return s.createRRs(splitMessage)
}

func (s SimpleMessageSplitter) Collect(rrs []dns.RR) string {
	var message string
	for _, rr := range rrs {
		message += strings.TrimSuffix(rr.Header().Name, ".")
	}
	return message
}

func (s SimpleMessageSplitter) markSliceBoundaries(message string) string {
	if len(message) < maxMessageSliceLen {
		return message
	} else {
		return message[:maxMessageSliceLen] + "." + s.markSliceBoundaries(message[maxMessageSliceLen:])
	}
}

func (s SimpleMessageSplitter) createRRs(message []string) []dns.RR {
	var rrs []dns.RR
	for _, m := range message {
		rr, err := dns.NewRR(m + " 3600 IN MX 10 example.com")
		if err != nil {
			println(err.Error())
		}
		rrs = append(rrs, rr)
	}
	return rrs
}
