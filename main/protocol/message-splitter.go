package protocol

import (
	"github.com/miekg/dns"
	"strings"
)

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

func (s SimpleMessageSplitter) markSliceBoundaries(message string) string {
	if len(message) < 59 {
		return message
	} else {
		return message[:59] + "." + s.markSliceBoundaries(message[59:])
	}
}

func (s SimpleMessageSplitter) createRRs(message []string) []dns.RR {
	var rrs []dns.RR
	for _, m := range message {
		//todo: newRR seems to add a dot at the end of the name -> is message reconstruction affected?
		rr, err := dns.NewRR(m + " 3600 IN MX 10 example.com")
		if err != nil {
			println(err.Error())
		}
		rrs = append(rrs, rr)
	}
	return rrs
}

/*
todo:
	* collect method
	* constant for max message slice length
*/
