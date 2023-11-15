package shell

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Shell struct {
	scanner        *bufio.Scanner
	inputProcessor func(string)
	navigator      *UnixNavigator
}

func NewShell(inputProcessor func(string)) *Shell {
	return &Shell{scanner: bufio.NewScanner(os.Stdin), inputProcessor: inputProcessor, navigator: NewUnixNavigator()}
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
			if strings.ContainsAny(text, "cd") {
				err := s.navigator.AddNavigation(text)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				s.inputProcessor(s.navigator.BuildCommand() + " && " + text)
			}
		} else {
			break
		}
	}
}

func (s Shell) handleScannerError() {
	fmt.Println("Error: ", s.scanner.Err())
}
