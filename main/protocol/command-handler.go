package protocol

type CommandHandler interface {
	HandleCommand(value string, pollCallback func(), answerCallback func(string), exitCallback func())
	Poll(pollCallback func())
	Answer(value string, answerCallback func(string))
	Exit(exitCallback func())
}
