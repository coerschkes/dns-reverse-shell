package listener

type connectionHandler struct {
	hasConnection bool
}

func newConnectionHandler() *connectionHandler {
	handler := connectionHandler{}
	handler.hasConnection = false
	return &handler
}

func (c *connectionHandler) setConnectionStatus(status bool) {
	c.hasConnection = status
}
