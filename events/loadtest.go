package events

import (
	"errors"
	"fmt"
	"strconv"
)

type loadTestImpl struct {
	Concurrency int64
	Host string
}

func NewLoadTest() Event {
	return &loadTestImpl{}
}

func (l loadTestImpl) Execute() {
	fmt.Println("Executing")
}

func (l loadTestImpl) Validate(args ...string) error {
	if len(args) < 4 {
		return errors.New("Arguments not valid for load test.\nplease try with -h (host) -c (concurrency workers)")
	}
	for i, arg := range args {
		if arg == "-c" {
			num, err := strconv.ParseInt(args[i+1], 10, 64)
			if err != nil {
				return errors.New(" -c arg is not a number.")
			}
			l.Concurrency = num
		}
		if arg == "-h" {
			l.Host = args[i+1]
		}
	}
	if l.Concurrency == 0 || l.Host == "" {
		return errors.New(" Something wrong with the parameters")
	}
	return nil
}