package englishWordsBot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go-app/configs"
	"go-app/internal/domain/user"
	"go-app/internal/domain/word"
	"log"
)

type Flow struct {
	Step uint8
	Word string
}

const (
	defaultBotOffset  = 0
	defaultBotTimeout = 60
)

type wordRepository interface {
	GetAllByChatId(chatId int64) ([]*word.Word, error)
	DeleteById(id string) error
	GetById(id string) *word.Word
	Update(w *word.Word) (*word.Word, error)
	GetRandomFive(chatId int64, langTo string) []*word.Word
	GetRandomTranslations(w *word.Word) []*word.Word
	GetByChatIdAndValue(chatId int64, value string) *word.Word
	GetRandom(chatId int64, maxRate int8) *word.Word
}

type userRepository interface {
	GetByChatId(chatId int64) *user.User
	Create(u *user.User) (*user.User, error)
	Update(u *user.User) (*user.User, error)
}

type wordService interface {
	AddWord(word *word.Word) (*word.Word, error)
	GetTranslations(w string, u *user.User) ([]string, error)
	GetRandomWords(count int) []string
	GetParagraph() string
}

type EnglishWordsBot struct {
	api            *tgbotapi.BotAPI
	wordRepository wordRepository
	wordService    wordService
	userRepository userRepository
}

func NewEnglishBot(
	config configs.BaseBotConfig,
	wordsRepository wordRepository,
	wordService wordService,
	userRepository userRepository,
) *EnglishWordsBot {
	api, err := tgbotapi.NewBotAPI(config.Token)

	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", api.Self.UserName)

	return &EnglishWordsBot{
		api:            api,
		wordRepository: wordsRepository,
		wordService:    wordService,
		userRepository: userRepository,
	}
}

func (b EnglishWordsBot) Build() {
	u := tgbotapi.NewUpdate(defaultBotOffset)
	u.Timeout = defaultBotTimeout
	updates := b.api.GetUpdatesChan(u)
	b.handleUpdates(updates)
}

func (b EnglishWordsBot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		switch {
		case update.CallbackQuery != nil:
			b.handleCallback(update)
		case update.Message == nil:
			continue
		case update.Message.IsCommand():
			b.handleCommand(update)
		case update.Message.Text != "":
			b.handleMenu(update)
		default:
			b.handleDefault(update)
		}
	}
}

func (b EnglishWordsBot) SendMsg(msg tgbotapi.MessageConfig) {
	_, err := b.api.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func (b EnglishWordsBot) handleDefault(u tgbotapi.Update) {
	b.SendMsg(tgbotapi.NewMessage(u.Message.Chat.ID, "I don`t know this action"))
	log.Print("Default action", u)
}

func (b EnglishWordsBot) GetChatIdFromUpdate(u tgbotapi.Update) {
	//TODO
}
