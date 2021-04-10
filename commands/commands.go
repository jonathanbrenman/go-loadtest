package commands

import "go-loadtest/events"

func NewCmd(cmd string) events.Event {
	eventMap := make(map[string]events.Event)
	eventMap["start"] = events.NewLoadTest()
	return eventMap[cmd]
}