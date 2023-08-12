package wordService

import (
	"go-app/internal/domain/user"
	"go-app/internal/domain/word"
	"go-app/internal/repositories/wordRepository"
	"go-app/pkg/languageDetector"
	"go-app/pkg/wordTranslator"
)

type WordService interface {
	AddWord(word *word.Word) (*word.Word, error)
	GetTranslations(w string, u *user.User) ([]string, error)
	GetRandomWords(count int) []string
}

type RandomWordsGenerator interface {
	GetRandomWords(count int) []string
}

type wordService struct {
	Repository           wordRepository.WordRepository
	RandomWordsGenerator RandomWordsGenerator
}

func NewWordService(
	repo wordRepository.WordRepository,
	randWordsGenerator RandomWordsGenerator,
) WordService {
	return &wordService{
		Repository:           repo,
		RandomWordsGenerator: randWordsGenerator,
	}
}

func (bs *wordService) AddWord(w *word.Word) (*word.Word, error) {
	return bs.Repository.AddWord(w)
}

func (bs *wordService) GetTranslations(w string, u *user.User) ([]string, error) {
	var langTo string

	inputLang, err := languageDetector.Detect(w, u.GetUserLangs())

	if err != nil {
		return nil, err
	}

	if u.LangFrom == inputLang {
		langTo = u.LangTo
	} else {
		langTo = u.LangFrom
	}

	translations, err := wordTranslator.Translate(w, inputLang, langTo)

	if err != nil {
		return nil, err
	}

	return translations, nil
}

func (bs *wordService) GetRandomWords(count int) []string {
	return bs.RandomWordsGenerator.GetRandomWords(count)
}
