package navigation

import (
	"testing"
)

func TestNavigationStack_Push(t *testing.T) {
	tests := []struct {
		name                 string
		path                 string
		wantedLen            int
		wantedStringsInOrder []string
	}{
		{name: "simple path", path: "slt", wantedLen: 2, wantedStringsInOrder: []string{"/", "slt"}},
		{name: "simple path with space", path: "sl t", wantedLen: 2, wantedStringsInOrder: []string{"/", "sl t"}},
		{name: "nested path", path: "slt/tes/tui", wantedLen: 6, wantedStringsInOrder: []string{"/", "tui", "/", "tes", "/", "slt"}},
		{name: "nested path with space", path: "slt/te s/tui", wantedLen: 6, wantedStringsInOrder: []string{"/", "tui", "/", "te s", "/", "slt"}},
		{name: "nested absolute path with space", path: "/slt/te s/tui", wantedLen: 7, wantedStringsInOrder: []string{"/", "tui", "/", "te s", "/", "slt", "/"}},
		{name: "home", path: "~", wantedLen: 2, wantedStringsInOrder: []string{"/", "~"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := newNavigationStack()
			n.push(tt.path)

			if n.stack.Len() != tt.wantedLen {
				t.Errorf("Push() = %v, want %v", n.len(), tt.wantedLen)
			}
			for i := 0; i < n.stack.Len()+1; i++ {
				stackValue := n.stack.Pop().(string)
				if stackValue != tt.wantedStringsInOrder[i] {
					t.Errorf("pushPathToStack() = %v, want %v", stackValue, tt.wantedStringsInOrder[i])
				}
			}
		})
	}
}

func TestNavigationStack_Build(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "simple path", path: "slt", want: "slt/"},
		{name: "test relative nested path", path: "slt/sjkd", want: "slt/sjkd/"},
		{name: "test absolute nested path", path: "/slt/sjkd", want: "/slt/sjkd/"},
		{name: "navigate home", path: "~", want: "~/"},
		{name: "navigate relative dot notation", path: "../..", want: "../../"},
		{name: "navigate empty stack", path: "", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := newNavigationStack()
			n.push(tt.path)
			if got := n.build(); got != tt.want {
				t.Errorf("build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNavigationStack_BuildWithPersistence(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "simple path", path: "slt", want: "slt/slt/"},
		{name: "test relative nested path", path: "slt/sjkd", want: "slt/sjkd/slt/sjkd/"},
		{name: "test absolute nested path", path: "/slt/sjkd", want: "/slt/sjkd/"},
		{name: "navigate home", path: "~", want: "~/"},
		{name: "navigate relative dot notation", path: "../..", want: "../../../../"},
		{name: "navigate empty stack", path: "", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := newNavigationStack()
			n.push(tt.path)
			n.build()
			n.push(tt.path)
			if got := n.build(); got != tt.want {
				t.Errorf("build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNavigationStack_Len(t *testing.T) {
	tests := []struct {
		name string
		path string
		want int
	}{
		{name: "simple path", path: "slt", want: 2},
		{name: "test relative nested path", path: "slt/sjkd", want: 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := newNavigationStack()
			n.push(tt.path)
			if got := n.len(); got != tt.want {
				t.Errorf("len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNavigationStack_Clear(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{name: "simple path", path: "slt"},
		{name: "test relative nested path", path: "slt/sjkd"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := newNavigationStack()
			n.push(tt.path)
			n.clear()

			if n.len() != 0 {
				t.Errorf("clear() = %v, want 0", n.len())
			}

		})
	}
}

func TestNavigationStack_isAbsolute(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{name: "simple path", path: "slt", want: false},
		{name: "test relative nested path", path: "slt/sjkd", want: false},
		{name: "test absolute nested path", path: "/slt/sjkd", want: true},
		{name: "navigate home", path: "~", want: false},
		{name: "navigate relative dot notation", path: "../..", want: false},
		{name: "navigate empty stack", path: "", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := newNavigationStack()
			if got := n.isAbsolute(tt.path); got != tt.want {
				t.Errorf("isAbsolute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNavigationStack_isHome(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{name: "simple path", path: "slt", want: false},
		{name: "test relative nested path", path: "slt/sjkd", want: false},
		{name: "test absolute nested path", path: "/slt/sjkd", want: false},
		{name: "navigate home", path: "~", want: true},
		{name: "navigate home with slash", path: "~/", want: true},
		{name: "navigate relative dot notation", path: "../..", want: false},
		{name: "navigate empty stack", path: "", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := newNavigationStack()
			if got := n.isHome(tt.path); got != tt.want {
				t.Errorf("isHome() = %v, want %v", got, tt.want)
			}
		})
	}
}
