package randomParagraphGenerator

import (
	"fmt"
	"io"
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
	resp, err := http.Get("http://metaphorpsum.com/paragraphs/2/4")

	if err != nil {
		fmt.Println(err)
		return "Paragraph function error."
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	
	if err != nil {
		fmt.Println(err)
		return "Paragraph function error."
	}

	return string(body)
}
