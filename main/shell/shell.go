package shell

import (
	"bufio"
	"fmt"
	"os"
)

type Shell struct {
	scanner        *bufio.Scanner
	inputProcessor func(string)
}

func NewShell(inputProcessor func(string)) *Shell {
	return &Shell{scanner: bufio.NewScanner(os.Stdin), inputProcessor: inputProcessor}
}

func (s Shell) Start() {
	fmt.Println("Enter command. Empty string exits the program")
	s.loopScanner()
	if s.scanner.Err() != nil {
		s.handleScannerError()
	}
}

func (s Shell) loopScanner() {
	for {
		s.scanner.Scan()
		text := s.scanner.Text()
		if len(text) != 0 {
			s.inputProcessor(text)
		} else {
			break
		}
	}
}

func (s Shell) handleScannerError() {
	fmt.Println("Error: ", s.scanner.Err())
}
