package shell

import (
	"errors"
	"testing"
)

func TestNavigator_AddNavigation(t *testing.T) {
	type args struct {
		navCommand string
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{name: "simple path", args: args{navCommand: "cd slt"}, err: nil},
		{name: "test relative nested path", args: args{navCommand: "cd slt/sjkd"}, err: nil},
		{name: "test absolute nested path", args: args{navCommand: "cd /slt/sjkd"}, err: nil},
		{name: "navigate home", args: args{navCommand: "cd ~"}, err: nil},
		{name: "navigate relative dot notation", args: args{navCommand: "cd ../.."}, err: nil},
		{name: "invalid command", args: args{navCommand: "d ~"}, err: errors.New("invalid command 'd ~'")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := UnixNavigator{
				navStack: NewNavigationStack(),
			}
			if err := n.AddNavigationStep(tt.args.navCommand); err == nil && tt.err != nil {
				t.Errorf("BuildCommand() = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestNavigator_BuildCommand(t *testing.T) {
	type args struct {
		navCommand string
	}
	tests := []struct {
		name string
		want string
		args args
	}{
		{name: "simple path", want: "cd slt/", args: args{navCommand: "cd slt"}},
		{name: "test relative nested path", want: "cd slt/sjkd/", args: args{navCommand: "cd slt/sjkd"}},
		{name: "test absolute nested path", want: "cd /slt/sjkd/", args: args{navCommand: "cd /slt/sjkd"}},
		{name: "navigate home", want: "cd ~/", args: args{navCommand: "cd ~"}},
		{name: "navigate relative dot notation", want: "cd ../../", args: args{navCommand: "cd ../.."}},
		{name: "navigate empty stack", want: "", args: args{navCommand: ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := UnixNavigator{
				navStack: NewNavigationStack(),
			}
			if tt.args.navCommand != "" {
				err := n.AddNavigationStep(tt.args.navCommand)
				if err != nil {
					t.Errorf("BuildCommand() = %v, want %v", err, tt.want)
				}
			}
			if got := n.BuildCommand(); got != tt.want {
				t.Errorf("BuildCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNavigator_getNavigationPath(t *testing.T) {
	type args struct {
		navCommand string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "nested navigation", args: args{navCommand: "cd slt/sjkd"}, want: "slt/sjkd"},
		{name: "simple navigation", args: args{navCommand: "cd slt"}, want: "slt"},
		{name: "nested navigation with whitespace", args: args{navCommand: "cd slt/dsa as/djksl"}, want: "slt/dsa as/djksl"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := UnixNavigator{
				navStack: NewNavigationStack(),
			}
			if got := n.getNavigationPath(tt.args.navCommand); got != tt.want {
				t.Errorf("getNavigationPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNavigator_isValidCommand(t *testing.T) {
	type args struct {
		navCommand string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "command without spaces", args: args{navCommand: "cd/"}, want: false},
		{name: "command with space", args: args{navCommand: "cd /"}, want: true},
		{name: "nested command without spaces", args: args{navCommand: "cd/asdjkflöasjf/jsakldfj/"}, want: false},
		{name: "nested command with spaces", args: args{navCommand: "cd /asdjkflöasjf/jfkasl"}, want: true},
		{name: "invalid command without content", args: args{navCommand: "cd"}, want: false},
		{name: "invalid command", args: args{navCommand: "c /"}, want: false},
		{name: "nested invalid command", args: args{navCommand: "d/asdjkflöasjf/jsakldfj"}, want: false},
		{name: "nested chained command", args: args{navCommand: "cd sui && ls"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := UnixNavigator{
				navStack: NewNavigationStack(),
			}
			if got := n.isValidCommand(tt.args.navCommand); got != tt.want {
				t.Errorf("isNavigationHome() = %v, want %v", got, tt.want)
			}
		})
	}
}
