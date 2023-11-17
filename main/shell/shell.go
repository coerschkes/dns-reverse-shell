package shell

import (
	"bufio"
	"dns-shellcode/main/shell/navigation"
	"fmt"
	"os"
	"strings"
)

type Shell struct {
	scanner        *bufio.Scanner
	inputProcessor func(string)
	navigator      *navigation.UnixNavigator
}

func NewShell(inputProcessor func(string)) *Shell {
	return &Shell{scanner: bufio.NewScanner(os.Stdin), inputProcessor: inputProcessor, navigator: navigation.NewUnixNavigator()}
}

func (s Shell) Start() {
	fmt.Println("Enter command. Empty string exits the program")
	s.printPrompt()
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
	err := s.navigator.AddNavigationStep(text)
	if err != nil {
		fmt.Println(err)
	}
	s.printPrompt()
}

func (s Shell) processInput(text string) {
	navCommand := s.navigator.BuildCommand()
	if len(navCommand) != 0 {
		s.inputProcessor(navCommand + " && " + text)
	} else {
		s.inputProcessor(text + "\n" + navCommand)
	}
	s.printPrompt()
}

func (s Shell) printPrompt() {
	path := s.navigator.BuildPath()
	if len(path) != 0 {
		fmt.Print(path + " > ")
		return
	}
	fmt.Print("> ")
}

func (s Shell) handleScannerError() {
	fmt.Println("Error: ", s.scanner.Err())
}
