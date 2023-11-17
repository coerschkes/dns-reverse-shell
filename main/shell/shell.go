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
		if len(text) == 0 {
			break
		}
		s.handleInput(text)
	}
}

func (s Shell) handleInput(text string) {
	if strings.ContainsAny(text, "cd") {
		s.handleNavigationCommand(text)
	} else {
		s.processInput(text)
	}
}

func (s Shell) handleNavigationCommand(text string) {
	err := s.navigator.AddNavigation(text)
	if err != nil {
		fmt.Println(err)
	}
}

func (s Shell) processInput(text string) {
	navigation := s.navigator.BuildCommand()
	if len(navigation) != 0 {
		s.inputProcessor(navigation + " && " + text)
	} else {
		s.inputProcessor(text)
	}
}

func (s Shell) handleScannerError() {
	fmt.Println("Error: ", s.scanner.Err())
}
