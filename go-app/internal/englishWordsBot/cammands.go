package englishWordsBot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go-app/internal/domain/user"
	"go-app/internal/domain/word"
	"go-app/pkg/bot/msgBuilder"
	"go-app/pkg/wordTranslator"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

type commandKey string

const (
	StartCmdKey = commandKey("start")
)

const DefaultRandWordsCount = 10

func (b EnglishWordsBot) handleCommand(u tgbotapi.Update) {
	switch commandKey(u.Message.Command()) {
	case StartCmdKey:
		b.commandStart(u)
	}
}

func (b EnglishWordsBot) commandExam(u tgbotapi.Update) {
	chatId := u.Message.Chat.ID
	randomFive := b.wordRepository.GetRandomFive(chatId)

	for _, w := range randomFive {
		var calls []*msgBuilder.Callback
		trans := b.wordRepository.GetRandomTranslations(w)
		trans = append(trans, w)

		for _, t := range trans {
			calls = append(calls, &msgBuilder.Callback{
				Key:    t.Translation,
				Type:   "answer",
				Action: t.Translation,
			})
		}

		msg := tgbotapi.NewMessage(u.Message.Chat.ID, w.Value)
		rand.New(rand.NewSource(time.Now().UnixNano()))
		rand.Shuffle(len(calls), func(i, j int) { calls[i], calls[j] = calls[j], calls[i] })
		msgBuilder.BuildKeyboardByCallbacks(&msg, calls)

		b.SendMsg(msg)
	}
}

func (b EnglishWordsBot) commandParagraph(u tgbotapi.Update) {
	chatId := u.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatId, b.wordService.GetParagraph())
	b.SendMsg(msg)
}

func (b EnglishWordsBot) commandList(u tgbotapi.Update) {
	words, err := b.wordRepository.GetAllByChatId(u.Message.Chat.ID)
	chatId := u.Message.Chat.ID

	if err != nil {
		b.SendError(chatId, err)
	}

	msg := tgbotapi.NewMessage(chatId, "Last words")
	msgBuilder.BuildManageKeyboard(&msg, words)

	b.SendMsg(msg)
}

func (b EnglishWordsBot) commandRandomWords(u tgbotapi.Update) {
	var words []word.Word
	var calls []*msgBuilder.Callback

	randWords := b.wordService.GetRandomWords(DefaultRandWordsCount)
	chatUser := b.userRepository.GetByChatId(u.Message.Chat.ID)

	var wg sync.WaitGroup
	wg.Add(len(randWords))

	for _, rw := range randWords {
		rw := rw
		go func() {
			translate, _ := wordTranslator.Translate(rw, "en", chatUser.LangFrom)

			words = append(words, word.Word{
				Translation: strings.Join(translate, ", "),
				Value:       rw,
			})
			wg.Done()
		}()
	}
	wg.Wait()
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, "List")

	for _, w := range words {
		calls = append(calls, &msgBuilder.Callback{
			Key:    w.Value + " - " + w.Translation,
			Type:   "manage-generated",
			Action: w.Value,
		})
	}

	msgBuilder.BuildKeyboardByCallbacks(&msg, calls)

	b.SendMsg(msg)
}

func (b EnglishWordsBot) commandUpdateInterval(u tgbotapi.Update, param string) {
	chatUser := b.userRepository.GetByChatId(u.CallbackQuery.Message.Chat.ID)
	interval, err := strconv.Atoi(param)
	if err != nil {
		return
	}
	chatUser.Interval = user.Interval(interval)
	_, err = b.userRepository.Update(chatUser)

	if err != nil {
		panic(err)
	}

	b.SendMsg(tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, "✅"))
}

func (b EnglishWordsBot) commandAskInterval(chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "Choose an interval to send words")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("30 minutes", "interval/30")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("1 hour", "interval/60")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("2 hours", "interval/120")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("3 hours", "interval/180")),
	)

	b.SendMsg(msg)
}
