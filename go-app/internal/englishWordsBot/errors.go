package englishWordsBot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func (b EnglishWordsBot) SendError(chatId int64, err error) {
	b.SendMsg(tgbotapi.NewMessage(chatId, "Something wrong. Try to write me later"))
	log.Panicln(err)
}
