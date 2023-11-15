package shell

import (
	"errors"
	"github.com/golang-collections/collections/stack"
	"strings"
)

type UnixNavigator struct {
	navigationStack *stack.Stack
}

func NewUnixNavigator() *UnixNavigator {
	return &UnixNavigator{navigationStack: stack.New()}
}

func (n UnixNavigator) AddNavigation(navCommand string) error {
	if !n.isValidCommand(navCommand) {
		return errors.New("invalid command '" + navCommand + "'")
	}
	if n.isNavigationHome(navCommand) || n.isNavigationAbsolute(navCommand) {
		n.clearNavigationStack()
	}
	navPath := n.getNavigationPath(navCommand)
	n.pushPathToStack(navPath)
	return nil
}

func (n UnixNavigator) BuildCommand() string {
	navigationPath := ""
	if n.navigationStack.Len() == 0 {
		return navigationPath
	}
	for n.navigationStack.Len() > 0 {
		navigationPath = n.navigationStack.Pop().(string) + navigationPath
	}
	return "cd " + navigationPath
}

func (n UnixNavigator) isNavigationAbsolute(command string) bool {
	trimmedCommand := strings.ReplaceAll(command, " ", "")
	return strings.HasPrefix(trimmedCommand, "cd/")
}

func (n UnixNavigator) isNavigationHome(command string) bool {
	trimmedCommand := strings.ReplaceAll(command, " ", "")
	return strings.HasPrefix(trimmedCommand, "cd~")
}

func (n UnixNavigator) getNavigationPath(command string) string {
	splitN := strings.SplitN(command, " ", 2)
	return splitN[1]
}

func (n UnixNavigator) pushPathToStack(navPath string) {
	if n.isNavigationAbsolute(navPath) {
		n.navigationStack.Push("/")
	}
	splitNavPath := strings.Split(navPath, "/")
	for i := range splitNavPath {
		n.navigationStack.Push(splitNavPath[i])
		n.navigationStack.Push("/")
	}
}

func (n UnixNavigator) clearNavigationStack() {
	for n.navigationStack.Len() > 0 {
		n.navigationStack.Pop()
	}
}

func (n UnixNavigator) isValidCommand(command string) bool {
	trimmedCommand := strings.ReplaceAll(command, " ", "")
	noEmptyCommand := strings.Compare("cd", trimmedCommand) != 0
	beginsWithCd := strings.HasPrefix(trimmedCommand, "cd")
	notContainingChains := !strings.ContainsAny(trimmedCommand, "&&")
	return beginsWithCd && noEmptyCommand && notContainingChains
}
