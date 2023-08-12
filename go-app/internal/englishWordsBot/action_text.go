package englishWordsBot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/pemistahl/lingua-go"
	"go-app/internal/domain/word"
	"go-app/pkg/bot/msgBuilder"
	"go-app/pkg/languageDetector"
)

func (b EnglishWordsBot) getWordAndSendTranslateOptions(u tgbotapi.Update) {
	var chatId int64
	var msgText string

	if u.CallbackQuery != nil {
		call := msgBuilder.CallbackStringToData(u.CallbackData())
		chatId = u.CallbackQuery.Message.Chat.ID
		msgText = call.Action
	} else {
		chatId = u.Message.Chat.ID
		msgText = u.Message.Text
	}
	chatUser := b.userRepository.GetByChatId(chatId)
	w := chatUser.WaitingType

	switch {
	case w != "":
		valueLang, err := languageDetector.Detect(w, chatUser.GetUserLangs())
		translationLang, err := languageDetector.Detect(msgText, chatUser.GetUserLangs())

		_, err = b.wordService.AddWord(&word.Word{
			Value:           w,
			Translation:     msgText,
			ChatId:          chatId,
			ValueLang:       valueLang,
			TranslationLang: translationLang,
		})

		if err != nil {
			b.SendError(chatId, err)
		}

		msg := tgbotapi.NewMessage(chatId, w+" - "+msgText)
		chatUser.WaitingType = ""
		_, err = b.userRepository.Update(chatUser)
		if err != nil {
			b.SendError(chatId, err)
		}
		b.SendMsg(msg)

		return
	}

	translations, err := b.wordService.GetTranslations(msgText, chatUser)

	if err != nil {
		b.SendError(chatId, err)
	}

	msg := tgbotapi.NewMessage(chatId, msgText)
	msgBuilder.BuildKeyboard(&msg, translations, msgText)

	b.SendMsg(msg)
}
