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
	GetParagraph() string
}

type RandomWordsGenerator interface {
	GetRandomWords(count int) []string
}

type ParagraphGenerator interface {
	GetRandomParagraph() string
}

type wordService struct {
	Repository           wordRepository.WordRepository
	RandomWordsGenerator RandomWordsGenerator
	ParagraphGenerator   ParagraphGenerator
}

func NewWordService(
	repo wordRepository.WordRepository,
	randWordsGenerator RandomWordsGenerator,
	paragraphGenerator ParagraphGenerator,
) WordService {
	return &wordService{
		Repository:           repo,
		RandomWordsGenerator: randWordsGenerator,
		ParagraphGenerator:   paragraphGenerator,
	}
}

func (bs *wordService) AddWord(w *word.Word) (*word.Word, error) {
	oldWord := bs.Repository.GetByValueAndTranslationLang(w.Value, w.TranslationLang)

	if oldWord == nil {
		return bs.Repository.AddWord(w)
	}

	oldWord.Translation = w.Translation

	return bs.Repository.Update(oldWord)
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

func (bs *wordService) GetParagraph() string {
	return bs.ParagraphGenerator.GetRandomParagraph()
}
