package shell

import (
	"errors"
	"regexp"
	"strings"
)

type UnixNavigator struct {
	navStack *NavigationStack
}

//todo: implement up/downwards navigation with "cd .." instead of putting the ".." on the navigation stack

func NewUnixNavigator() *UnixNavigator {
	return &UnixNavigator{navStack: NewNavigationStack()}
}

func (n UnixNavigator) AddNavigationStep(navCommand string) error {
	if !n.isValidCommand(navCommand) {
		return errors.New("invalid command '" + navCommand + "'")
	}
	navPath := n.getNavigationPath(navCommand)
	n.navStack.Push(navPath)
	return nil
}

func (n UnixNavigator) BuildCommand() string {
	path := n.navStack.Build()
	if path == "" {
		return path
	}
	return "cd " + path
}

func (n UnixNavigator) getNavigationPath(command string) string {
	splitN := strings.SplitN(command, " ", 2)
	return splitN[1]
}

func (n UnixNavigator) isValidCommand(command string) bool {
	startsWithCd, _ := regexp.MatchString("cd .+", command)
	notContainingChains := !strings.ContainsAny(command, "&&")
	return startsWithCd && notContainingChains
}
