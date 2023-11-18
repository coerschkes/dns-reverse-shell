package navigation

import (
	"github.com/golang-collections/collections/stack"
	"strings"
)

type navigationStack struct {
	stack *stack.Stack
}

func newNavigationStack() *navigationStack {
	return &navigationStack{stack: stack.New()}
}

func (n *navigationStack) push(path string) {
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

func (n *navigationStack) clear() {
	for n.stack.Len() > 0 {
		n.stack.Pop()
	}
}

func (n *navigationStack) len() int {
	return n.stack.Len()
}

func (n *navigationStack) build() string {
	if n.len() == 0 {
		return ""
	} else {
		return n.calculatePath()
	}
}

func (n *navigationStack) handleExtraNavigation(path string) {
	if n.isAbsolute(path) {
		n.clear()
		n.stack.Push("/")
	} else if n.isHome(path) {
		n.clear()
	}
}

func (n *navigationStack) isAbsolute(path string) bool {
	trimmedPath := strings.ReplaceAll(path, " ", "")
	return strings.HasPrefix(trimmedPath, "/")
}

func (n *navigationStack) isHome(path string) bool {
	trimmedPath := strings.ReplaceAll(path, " ", "")
	return trimmedPath == "~" || trimmedPath == "~/"
}

func (n *navigationStack) calculatePath() string {
	navigationPath := ""
	tmpStack := n.copyInternalStack()
	for tmpStack.Len() > 0 {
		currentElement := tmpStack.Pop()
		navigationPath = currentElement.(string) + navigationPath
	}
	return navigationPath
}

func (n *navigationStack) copyInternalStack() *stack.Stack {
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

func (n *navigationStack) revertInternalStack() {
	stackBuffer := stack.New()
	for n.stack.Len() > 0 {
		stackBuffer.Push(n.stack.Pop())
	}
	n.stack = stackBuffer
}
