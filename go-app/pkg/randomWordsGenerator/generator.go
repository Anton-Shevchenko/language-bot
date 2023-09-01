package randomWordsGenerator

import (
	"encoding/json"
	"net/http"
	"sync"
)

type Config struct {
	Url    string
	ApiKey string
}

type Response struct {
	Word string `json:"word"`
}

type RandomWordsGenerator interface {
	GetRandomWords(count int) []string
}

type Service struct {
	url    string
	apiKey string
}

func NewRandomWordsGenerator(cf Config) RandomWordsGenerator {
	return &Service{url: cf.Url, apiKey: cf.ApiKey}
}

func (rw *Service) doRequest() string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", rw.url, nil)

	if err != nil {

	}

	req.Header.Set("X-Api-Key", rw.apiKey)
	resp, err := client.Do(req)

	if err != nil {

	}

	r := &Response{}

	err = json.NewDecoder(resp.Body).Decode(r)

	if err != nil {
		//TODO
	}

	return r.Word
}

func (rw *Service) GetRandomWords(count int) []string {
	var words []string

	var wg sync.WaitGroup
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()

			words = append(words, rw.doRequest())
		}()
	}
	wg.Wait()

	return words
}
