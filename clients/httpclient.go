package clients

import (
	"bytes"
	"log"
	"net/http"
	"sync"
)

type HttClient interface {
	Get(ch chan<- int, wg *sync.WaitGroup)
	Post(ch chan<- int, wg *sync.WaitGroup)
}

type httpClientImpl struct {
	url string
	method string
	payload string
}

func NewHttpClient(url, method, payload string) HttClient {
	return httpClientImpl{
		url: url,
		payload: payload,
	}
}

func (h httpClientImpl) Get(ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(h.url)
	if err != nil {
		log.Println("Error doing get request", err.Error())
	}
	if resp != nil {
		ch <- resp.StatusCode
	} else {
		ch <- 0
	}
}

func (h httpClientImpl) Post(ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	req, err := http.NewRequest("POST", h.url, bytes.NewBuffer([]byte(h.payload)))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error doing post request", err.Error())
	}
	defer resp.Body.Close()
	if resp != nil {
		ch <- resp.StatusCode
	} else {
		ch <- 0
	}
}