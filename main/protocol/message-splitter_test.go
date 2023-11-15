package protocol

import (
	"dns-shellcode/main/encoder"
	"testing"
)

func TestSimpleMessageSplitter_Split(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name         string
		args         args
		wantedLength int
	}{
		{
			name:         "split simple message (less than 60 bytes)",
			args:         args{message: "simple message"},
			wantedLength: 1,
		},
		{
			name:         "split message (more than 60 bytes)",
			args:         args{message: "TEST123456789_123456789_123456789_123456789_1_1234567890"},
			wantedLength: 2,
		},
		{
			name:         "complicated message (more than 60 bytes)",
			args:         args{message: "TEST123456789_123456789_123456789_123456789_1_1234567890!(%)/&§)%&&?&?`*'*`$&(      \n\n kjalöfjahsklghalsdög"},
			wantedLength: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SimpleMessageSplitter{}
			encodedMessage := encoder.NewBase64Encoder().Encode(tt.args.message)
			if len(s.Split(encodedMessage)) != tt.wantedLength {
				t.Errorf("Split() = %v, want %v", len(s.Split(encodedMessage)), tt.wantedLength)
			}
			for i := range s.Split(encodedMessage) {
				if len(s.Split(encodedMessage)[i].Header().Name) > 60 {
					t.Errorf("len(Split()[i=%v] = %v), want %v", i, s.Split(encodedMessage)[i].Header().Name, encodedMessage)
				}
			}
		})
	}
}
