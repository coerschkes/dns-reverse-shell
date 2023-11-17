package shell

import (
	"github.com/golang-collections/collections/stack"
	"strings"
)

type NavigationStack struct {
	stack *stack.Stack
}

func NewNavigationStack() *NavigationStack {
	return &NavigationStack{stack: stack.New()}
}

func (n *NavigationStack) Push(path string) {
	n.handleExtraNavigation(path)
	splitNavPath := strings.Split(path, "/")
	for i := range splitNavPath {
		if splitNavPath[i] == "" {
			continue
		}
		n.stack.Push(splitNavPath[i])
		n.stack.Push("/")
	}
}

func (n *NavigationStack) Clear() {
	for n.stack.Len() > 0 {
		n.stack.Pop()
	}
}

func (n *NavigationStack) Len() int {
	return n.stack.Len()
}

func (n *NavigationStack) Build() string {
	if n.Len() == 0 {
		return ""
	} else {
		return n.calculatePath()
	}
}

func (n *NavigationStack) handleExtraNavigation(path string) {
	if n.isAbsolute(path) {
		n.Clear()
		n.stack.Push("/")
	} else if n.isHome(path) {
		n.Clear()
	}
}

func (n *NavigationStack) isAbsolute(path string) bool {
	trimmedPath := strings.ReplaceAll(path, " ", "")
	return strings.HasPrefix(trimmedPath, "/")
}

func (n *NavigationStack) isHome(path string) bool {
	trimmedPath := strings.ReplaceAll(path, " ", "")
	return trimmedPath == "~" || trimmedPath == "~/"
}

func (n *NavigationStack) calculatePath() string {
	navigationPath := ""
	tmpStack := n.copyInternalStack()
	for tmpStack.Len() > 0 {
		currentElement := tmpStack.Pop()
		navigationPath = currentElement.(string) + navigationPath
	}
	return navigationPath
}

func (n *NavigationStack) copyInternalStack() *stack.Stack {
	stackBuffer := stack.New()
	tmpStack := stack.New()
	n.revertInternalStack()
	for n.stack.Len() > 0 {
		currentValue := n.stack.Pop()
		tmpStack.Push(currentValue)
		stackBuffer.Push(currentValue)
	}
	n.stack = stackBuffer
	return tmpStack
}

func (n *NavigationStack) revertInternalStack() {
	stackBuffer := stack.New()
	for n.stack.Len() > 0 {
		stackBuffer.Push(n.stack.Pop())
	}
	n.stack = stackBuffer
}
