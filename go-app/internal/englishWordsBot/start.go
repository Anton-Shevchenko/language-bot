package englishWordsBot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go-app/internal/domain/user"
	"strconv"
)

func (b EnglishWordsBot) commandStart(u tgbotapi.Update) {
	savedUser := b.userRepository.GetByChatId(u.Message.Chat.ID)

	if savedUser.ChatId == 0 {
		_, err := b.userRepository.Create(&user.User{
			ChatId:   u.Message.Chat.ID,
			MaxRate:  5,
			Interval: user.Interval60,
		})

		if err != nil {
			return
		}

		b.askNativeLang(u.Message.Chat.ID)

		return
	}

	if savedUser.LangFrom == "" {
		b.askNativeLang(u.Message.Chat.ID)

		return
	}

	if savedUser.LangTo == "" {
		b.askTargetLang(u.Message.Chat.ID)

		return
	}

	if savedUser.Level == "" {
		b.askLangLevel(u.Message.Chat.ID)

		return
	}

	b.SendMsg(tgbotapi.NewMessage(u.Message.Chat.ID, "Hi again"))
	b.sendMenu(u.Message.Chat.ID)
}

func (b EnglishWordsBot) askLangLevel(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Choose your level language you want to study")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Low", "level/low"),
			tgbotapi.NewInlineKeyboardButtonData("Medium", "level/medium"),
			tgbotapi.NewInlineKeyboardButtonData("High", "level/high"),
		),
	)

	b.SendMsg(msg)
}

func (b EnglishWordsBot) askMaxRate(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Choose your max rate for words")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("3", "max-rate/3"),
			tgbotapi.NewInlineKeyboardButtonData("5", "max-rate/5"),
			tgbotapi.NewInlineKeyboardButtonData("10", "max-rate/10"),
		),
	)

	b.SendMsg(msg)
}

func (b EnglishWordsBot) askNativeLang(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Choose your native language")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Ukrainian", "native/uk"),
			tgbotapi.NewInlineKeyboardButtonData("Russian", "native/ru"),
		),
	)

	_, _ = b.api.Send(msg)
}

func (b EnglishWordsBot) askTargetLang(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Choose language you want to study")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("English", "target/en"),
			tgbotapi.NewInlineKeyboardButtonData("French", "target/fr"),
			tgbotapi.NewInlineKeyboardButtonData("Germany", "target/de"),
		),
	)

	b.SendMsg(msg)
}

func (b EnglishWordsBot) updateNativeLang(u tgbotapi.Update, param string) {
	savedUser := b.userRepository.GetByChatId(u.CallbackQuery.Message.Chat.ID)

	savedUser.LangFrom = param
	_, err := b.userRepository.Update(savedUser)
	if err != nil {
		fmt.Println(err.Error())

		return
	}

	if savedUser.LangTo == "" {
		b.askTargetLang(u.CallbackQuery.Message.Chat.ID)

		return
	}

	b.SendMsg(tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, "✅"))
}

func (b EnglishWordsBot) updateLevelLang(u tgbotapi.Update, param string) {
	savedUser := b.userRepository.GetByChatId(u.CallbackQuery.Message.Chat.ID)
	savedUser.Level = param
	_, err := b.userRepository.Update(savedUser)

	if err != nil {
		panic(err)
	}

	b.SendMsg(tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, "✅"))
}

func (b EnglishWordsBot) updateMaxRate(u tgbotapi.Update, param string) {
	savedUser := b.userRepository.GetByChatId(u.CallbackQuery.Message.Chat.ID)
	atoi, err := strconv.Atoi(param)
	if err != nil {
		return
	}
	savedUser.MaxRate = int8(atoi)
	_, err = b.userRepository.Update(savedUser)

	if err != nil {
		panic(err)
	}

	b.SendMsg(tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, "✅"))
}

func (b EnglishWordsBot) updateTargetLang(u tgbotapi.Update, param string) {
	savedUser := b.userRepository.GetByChatId(u.CallbackQuery.Message.Chat.ID)
	savedUser.LangTo = param
	_, err := b.userRepository.Update(savedUser)

	if err != nil {
		panic(err)
	}

	if savedUser.Level == "" {
		b.askLangLevel(u.CallbackQuery.Message.Chat.ID)

		return
	}

	b.SendMsg(tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, "✅"))
}

func (b EnglishWordsBot) askNotDisturbInterval(u tgbotapi.Update) {
	msg := tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, "Choose an interval")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Unset", "not-disturb-time/0"),
			tgbotapi.NewInlineKeyboardButtonData("3h", fmt.Sprintf("not-disturb-time/%d", 3*60)),
			tgbotapi.NewInlineKeyboardButtonData("8h", fmt.Sprintf("not-disturb-time/%d", 8*60)),
			tgbotapi.NewInlineKeyboardButtonData("10h", fmt.Sprintf("not-disturb-time/%d", 10*60)),
		),
	)

	b.SendMsg(msg)
}
