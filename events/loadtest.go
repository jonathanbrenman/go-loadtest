package events

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type LoadTestImpl struct {
	concurrency int64
	host string
}

func NewLoadTest() Event {
	return &LoadTestImpl{}
}

func (l LoadTestImpl) Execute(args ...string) error {
	if len(args) < 4 {
		return errors.New("Arguments not valid for load test.\nplease try with -h (host) -c (concurrency workers)")
	}
	for i, arg := range args {
		if arg == "-c" {
			num, err := strconv.ParseInt(args[i+1], 10, 64)
			if err != nil {
				return errors.New(" -c arg is not a number.")
			}
			l.concurrency = num
		}
		if arg == "-h" {
			l.host = args[i+1]
		}
	}

	if l.concurrency == 0 || l.host == "" {
		return errors.New(" Something wrong with the parameters")
	}

	fmt.Println("Starting load test...")
	ch := make(chan int)
	var wg sync.WaitGroup
	for i := 0; i < int(l.concurrency); i++ {
		wg.Add(1)
		go l.DoRequest(ch, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	cnSuccess := 0
	cnError := 0
	for res := range ch {
		if res == 200 {
			cnSuccess++
		} else {
			cnError++
		}
	}
	fmt.Printf("Successed: %d, Failed: %d\n", cnSuccess, cnError)
	return nil
}

func (l LoadTestImpl) DoRequest(ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(l.host)
	if err != nil {
		log.Println("Error doing request", err.Error())
	}
	if resp != nil {
		ch <- resp.StatusCode
	} else {
		ch <- 0
	}
}