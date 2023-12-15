package navigation

import (
	"errors"
	"regexp"
	"strings"
)

type UnixNavigator struct {
	navStack *navigationStack
}

func NewUnixNavigator() *UnixNavigator {
	return &UnixNavigator{navStack: newNavigationStack()}
}

func (n UnixNavigator) AddNavigationStep(navCommand string) error {
	if !n.isValidCommand(navCommand) {
		return errors.New("invalid command '" + navCommand + "'")
	}
	navPath := n.getNavigationPath(navCommand)
	n.navStack.push(navPath)
	return nil
}

func (n UnixNavigator) BuildPath() string {
	return n.navStack.build()
}

func (n UnixNavigator) BuildCommand() string {
	path := n.BuildPath()
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
