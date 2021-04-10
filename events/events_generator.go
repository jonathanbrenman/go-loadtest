package events

type Event interface {
	Validate(args ...string) error
	Execute()
}
