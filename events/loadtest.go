package events

import (
	"errors"
	"fmt"
	"go-loadtest/clients"
	"strconv"
	"sync"
	"time"
)

type LoadTestImpl struct {
	concurrency int64
	host string
	method string
	payload string
	time time.Duration
}

func NewLoadTest() Event {
	return &LoadTestImpl{}
}

func (l LoadTestImpl) Execute(args ...string) error {
	if len(args) < 4 {
		return errors.New("Arguments not valid for load test.\nplease try with -h (host) -c (concurrency workers)")
	}
	for i, arg := range args {
		switch arg {
			case "-c":
				num, err := strconv.ParseInt(args[i+1], 10, 64)
				if err != nil {
					return errors.New(" -c arg is not a number.")
				}
				l.concurrency = num
				break
			case "-h":
				l.host = args[i+1]
				break
			case "-m":
				l.method = args[i+1]
				break
			case "-p":
				l.payload = args[i+1]
				break
			case "-t":
				var parseErr error
				l.time, parseErr = time.ParseDuration(args[i+1])
				if parseErr != nil {
					return errors.New(" -t arg is not a duration format")
				}
				break
		}
	}

	if l.concurrency == 0 || l.host == "" || (l.method == "post" && l.payload == ""){
		return errors.New(" Something wrong with the parameters")
	}

	fmt.Println("Starting load test...")
	fmt.Println(l.Start())
	return nil
}

func (l LoadTestImpl) Start() string {
	ch := make(chan int)
	var wg sync.WaitGroup

	httpClient := clients.NewHttpClient(l.host, l.method, l.payload)

	if l.time == 0 {
		l.time = 1 * time.Second
	}

	start := time.Now()
	end := start.Add(l.time)

	for end.After(time.Now()) {
		fmt.Printf("performing %d requests to %s \n", l.concurrency, l.host)
		for i := 0; i < int(l.concurrency); i++ {
			wg.Add(1)
			switch l.method {
			case "get":
				go httpClient.Get(ch, &wg)
				break
			case "post":
				go httpClient.Post(ch, &wg)
				break
			default:
				go httpClient.Get(ch, &wg)
				break
			}
		}
		time.Sleep(1 * time.Second)
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
	return fmt.Sprintf("Successed: %d, Failed: %d\n", cnSuccess, cnError)
}