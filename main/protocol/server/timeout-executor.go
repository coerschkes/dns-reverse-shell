package server

import "time"

type timeoutExecutor struct {
	counter  int
	callback func()
	ex       bool
}

func newTimeoutExecutor() *timeoutExecutor {
	handler := timeoutExecutor{}
	handler.counter = 0
	handler.ex = false
	return &handler
}

func (t *timeoutExecutor) start() {
	for {
		if t.ex {
			return
		}
		time.Sleep(1 * time.Second)
		t.counter++
		if t.counter >= 10 {
			t.exit()
			t.callback()
		}
	}
}

func (t *timeoutExecutor) reset() {
	t.counter = 0
}

func (t *timeoutExecutor) exit() {
	t.ex = true
}
