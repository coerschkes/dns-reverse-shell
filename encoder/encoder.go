package encoder

import "encoding/base64"

func Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func Decode(str string) string {
	decodeString, _ := base64.StdEncoding.DecodeString(str)
	return string(decodeString)
}
