package msgBuilder

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go-app/internal/domain/word"
)

func NewInline(rows []tgbotapi.InlineKeyboardButton) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(rows)
}

func AddInlineRow() []tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardRow()
}

func AddReplyRow() []tgbotapi.KeyboardButton {
	return tgbotapi.NewKeyboardButtonRow()
}

func AddInlineButton(row *[]tgbotapi.InlineKeyboardButton, callback *Callback) {
	*row = append(*row, tgbotapi.NewInlineKeyboardButtonData(callback.Key, CallbackDataToString(callback)))
}

func AddReplyButton(row *[]tgbotapi.KeyboardButton, key string) {
	*row = append(*row, tgbotapi.NewKeyboardButton(key))
}

// Old below

func BuildKeyboard(msg *tgbotapi.MessageConfig, keys []string, word string) {
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, k := range keys {
		btn := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(k, "translate/"+k+"/"+word))
		rows = append(rows, btn)
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("My Option", "translate/MyOption/"+word)))

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func BuildManageKeyboard(msg *tgbotapi.MessageConfig, keys []*word.Word) {
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, w := range keys {
		btn := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(w.Value, "manage/"+w.ID.Hex()))
		rows = append(rows, btn)
	}

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func BuildKeyboardByCallbacks(msg *tgbotapi.MessageConfig, callbacks []*Callback) {
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, c := range callbacks {
		btn := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(c.Key, CallbackDataToString(c)))
		rows = append(rows, btn)
	}

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
}
