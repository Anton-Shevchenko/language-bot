package englishWordsBot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	Settings = "Settings ⚙️"
)

func (b EnglishWordsBot) handleMenu(u tgbotapi.Update) {
	switch u.Message.Text {
	case Settings:
		b.sendSettingsMenu(u)
		return
	case "New Words 📖":
		b.commandRandomWords(u)
		return
	case "Test me 💻️":
		b.commandExam(u)
		return
	case "My words ✏️":
		b.commandList(u)
		return
	default:
		b.getWordAndSendTranslateOptions(u)
	}
}

func (b EnglishWordsBot) sendMenu(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Menu")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(fmt.Sprintf(Settings)),

			tgbotapi.NewKeyboardButton("New Words 📖"),
			tgbotapi.NewKeyboardButton("Test me 💻️"),
			tgbotapi.NewKeyboardButton("My words ✏️"),
		),
	)

	_, err := b.api.Send(msg)

	if err != nil {
		fmt.Println("ERRR", err.Error())
	}
}

func (b EnglishWordsBot) sendSettingsMenu(u tgbotapi.Update) {
	var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Change interval", "change/interval"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Change target language", "change/target-language"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Change target language level", "change/language-level"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Change native language", "change/native-language"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Change max rate", "change/max-rate"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Usually don`t disturb me from now", "change/not-disturb-time"),
		),
	)
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
	msg.ReplyMarkup = numericKeyboard

	if _, err := b.api.Send(msg); err != nil {
		panic(err)
	}
}
