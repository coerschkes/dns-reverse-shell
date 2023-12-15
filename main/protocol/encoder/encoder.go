package encoder

import "encoding/base64"

type StringEncoder interface {
	Encode(string) string
	Decode(string) string
}

type Base64Encoder struct{}

func NewBase64Encoder() *Base64Encoder {
	return &Base64Encoder{}
}

func (e Base64Encoder) Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func (e Base64Encoder) Decode(str string) string {
	decodeString, _ := base64.StdEncoding.DecodeString(str)
	return string(decodeString)
}
