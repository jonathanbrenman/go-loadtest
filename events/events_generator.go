package events

type Event interface {
	Execute(args ...string) error
}
