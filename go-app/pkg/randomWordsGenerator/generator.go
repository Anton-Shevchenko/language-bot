package randomWordsGenerator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type Config struct {
	Url    string
	ApiKey string
}

type Response struct {
	Word []string `json:"word"`
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

func (rw *Service) doRequest() (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", rw.url, nil)

	if err != nil {
		return "", fmt.Errorf("request creation failed: %w", err)
	}

	req.Header.Set("X-Api-Key", rw.apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request execution failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)

	var result Response

	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("HHHHHH", result)
		return "", fmt.Errorf("json decoding failed: %w", err)
	}

	if len(result.Word) == 0 {
		return "", errors.New("received empty word")
	}

	return result.Word[0], nil
}

func (rw *Service) GetRandomWords(count int) []string {
	words := make([]string, 0, count)
	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(count)

	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()

			word, err := rw.doRequest()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			mu.Lock()
			words = append(words, word)
			mu.Unlock()
		}()
	}

	wg.Wait()
	return words
}
