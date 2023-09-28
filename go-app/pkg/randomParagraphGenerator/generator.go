package randomParagraphGenerator

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Config struct {
	Url string
}

type Response struct {
	Paragraph string `json:"definition"`
}

type RandomParagraphGenerator interface {
	GetRandomParagraph() string
}

type Service struct {
	url string
}

func NewRandomParagraphGenerator(cf Config) RandomParagraphGenerator {
	return &Service{url: cf.Url}
}

func (rw *Service) GetRandomParagraph() string {
	var r []Response
	client := &http.Client{}
	req, err := http.NewRequest("GET", rw.url, nil)

	if err != nil {
		fmt.Println("Error create request:", err.Error())
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error do request:", err.Error())
	}

	err = json.NewDecoder(resp.Body).Decode(&r)

	if err != nil {
		fmt.Println("Error decode:", err.Error())
	}

	if len(r) > 0 {
		return r[0].Paragraph
	}

	return "Paragraph function error."
}
