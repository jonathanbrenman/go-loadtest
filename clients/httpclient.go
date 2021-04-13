package clients

import (
	"bytes"
	"go-loadtest/models"
	"log"
	"net/http"
	"sync"
	"time"
)

type HttClient interface {
	Get(ch chan<-models.ResponseRoutine, wg *sync.WaitGroup)
	Post(ch chan<-models.ResponseRoutine, wg *sync.WaitGroup)
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

func (h httpClientImpl) Get(ch chan<-models.ResponseRoutine, wg *sync.WaitGroup) {
	start := time.Now()
	defer wg.Done()
	resp, err := http.Get(h.url)
	if err != nil {
		log.Println("Error doing get request", err.Error())
	}
	since := time.Since(start)
	response := models.ResponseRoutine{
		Code: resp.StatusCode,
		Time: since,
	}
	if resp != nil {
		ch <- response
	} else {
		ch <- models.ResponseRoutine{}
	}
}

func (h httpClientImpl) Post(ch chan<-models.ResponseRoutine, wg *sync.WaitGroup) {
	start := time.Now()
	defer wg.Done()
	req, err := http.NewRequest("POST", h.url, bytes.NewBuffer([]byte(h.payload)))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error doing post request", err.Error())
	}
	defer resp.Body.Close()
	since := time.Since(start)
	response := models.ResponseRoutine{
		Code: resp.StatusCode,
		Time: since,
	}
	if resp != nil {
		ch <- response
	} else {
		ch <- models.ResponseRoutine{}
	}
}