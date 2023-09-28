package englishWordsBot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go-app/pkg/bot/msgBuilder"
)

const (
	Settings  = "Settings ⚙️"
	Cabinet   = "Cabinet 👤"
	Back      = "Back 🔼"
	NewWords  = "New Words 🆕"
	Paragraph = "Paragraph 📝️"
	TestMe    = "Test me 💻️"
)

func (b EnglishWordsBot) handleMenu(u tgbotapi.Update) {
	switch u.Message.Text {
	case Cabinet:
		b.sendCabinetMenu(u)
	case NewWords:
		b.commandRandomWords(u)
	case Paragraph:
		b.commandParagraph(u)
	case Settings:
		b.sendSettingsMenu(u)
		return
	case Back:
		b.sendMainMenu(u.Message.Chat.ID)
		return
	case TestMe:
		b.commandExam(u)
		return
	case "My words ✏️":
		b.commandList(u)
		return
	default:
		b.getWordAndSendTranslateOptions(u)
	}
}

func buildMainMenu() tgbotapi.ReplyKeyboardMarkup {
	firstRow := msgBuilder.AddReplyRow()
	msgBuilder.AddReplyButton(&firstRow, Settings)
	msgBuilder.AddReplyButton(&firstRow, Cabinet)

	return tgbotapi.NewReplyKeyboard(firstRow)
}

func buildCabinetMenu() tgbotapi.ReplyKeyboardMarkup {
	firstRow := msgBuilder.AddReplyRow()
	msgBuilder.AddReplyButton(&firstRow, NewWords)
	msgBuilder.AddReplyButton(&firstRow, Paragraph)
	msgBuilder.AddReplyButton(&firstRow, TestMe)
	secondRow := msgBuilder.AddReplyRow()
	msgBuilder.AddReplyButton(&secondRow, Back)

	return tgbotapi.NewReplyKeyboard(firstRow, secondRow)
}

func (b EnglishWordsBot) sendMainMenu(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Menu")
	msg.ReplyMarkup = buildMainMenu()

	b.SendMsg(msg)
}

func (b EnglishWordsBot) sendCabinetMenu(u tgbotapi.Update) {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, Cabinet)
	msg.ReplyMarkup = buildCabinetMenu()

	b.SendMsg(msg)
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
