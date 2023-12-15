package protocol

import (
	"dns-reverse-shell/main/encoder"
	"github.com/miekg/dns"
)

type MessageHandler struct {
	encoder         encoder.StringEncoder
	messageSplitter MessageSplitter
}

func NewMessageHandler(encoder encoder.StringEncoder, messageSplitter MessageSplitter) *MessageHandler {
	return &MessageHandler{encoder: encoder, messageSplitter: messageSplitter}
}

func (m MessageHandler) CreateAnswerMessage(r *dns.Msg) *dns.Msg {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Authoritative = true
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
