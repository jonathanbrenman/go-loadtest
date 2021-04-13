package events

import (
	"errors"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"go-loadtest/clients"
	"go-loadtest/models"
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
	locker sync.Mutex
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
	ch := make(chan models.ResponseRoutine)
	var wg sync.WaitGroup

	httpClient := clients.NewHttpClient(l.host, l.method, l.payload)

	if l.time == 0 {
		l.time = 1 * time.Second
	}

	start := time.Now()
	end := start.Add(l.time)

	bar := progressbar.Default(100)
	for end.After(time.Now()) {
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
		bar.Add(int(l.concurrency) * 100 / (int(l.concurrency) * int(l.time.Seconds())))
		time.Sleep(1 * time.Second)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	responseCodes := make(map[int]struct{
		Code int
		Counter int
	})

	var avgTime time.Duration
	total := 0

	for res := range ch {
		l.locker.Lock()
		if responseCodes[res.Code].Counter != 0 {
			responseCodes[res.Code] = struct {
				Code    int
				Counter int
			}{Code: res.Code, Counter: responseCodes[res.Code].Counter + 1}
		} else {
			responseCodes[res.Code] = struct {
				Code    int
				Counter int
			}{Code: res.Code, Counter: 1}
		}
		avgTime += res.Time
		total++
		l.locker.Unlock()
	}
	response := "\n-----------------------------------\nResult:\n-----------------------------------\n"
	for k, _ := range responseCodes {
		response += fmt.Sprintf("http code: %d count: %d \n", responseCodes[k].Code, responseCodes[k].Counter)
	}
	response += fmt.Sprintf("response time avg: %vms\n", avgTime.Milliseconds() / int64(total))
	response += "-----------------------------------\n"
	return response
}