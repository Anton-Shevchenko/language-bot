package englishWordsBot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go-app/internal/domain/word"
	"go-app/pkg/bot/msgBuilder"
	"go-app/pkg/languageDetector"
	"strconv"
	"strings"
	"time"
)

func (b EnglishWordsBot) handleCallback(u tgbotapi.Update) {
	call := msgBuilder.CallbackStringToData(u.CallbackData())
	callType := call.Type
	callAction := call.Action
	fmt.Println("TYPE", callType)
	switch callType {
	case "manage":
		b.manageWord(u, callAction)
	case "delete":
		b.deleteWord(u, callAction)
	case "edit":
		b.editWord(u)
	case "translate":
		b.callbackTranslation(u)
	case "native":
		b.updateNativeLang(u, callAction)
	case "target":
		b.updateTargetLang(u, callAction)
	case "level":
		b.updateLevelLang(u, callAction)
	case "max-rate":
		b.updateMaxRate(u, callAction)
	case "interval":
		b.commandUpdateInterval(u, callAction)
	case "manage-generated":
		b.getWordAndSendTranslateOptions(u)
	case "change":
		switch callAction {
		case "interval":
			b.commandAskInterval(u.CallbackQuery.Message.Chat.ID)
		case "target-language":
			b.askTargetLang(u.CallbackQuery.Message.Chat.ID)
		case "native-language":
			b.askNativeLang(u.CallbackQuery.Message.Chat.ID)
		case "language-level":
			b.askLangLevel(u.CallbackQuery.Message.Chat.ID)
		case "max-rate":
			b.askMaxRate(u.CallbackQuery.Message.Chat.ID)
		case "not-disturb-time":
			b.askNotDisturbInterval(u)
		}
	case "not-disturb-time":
		b.updateNotDisturbInterval(u, call)
	case "answer":
		b.handleAnswer(u)
	}
}

func (b EnglishWordsBot) updateNotDisturbInterval(u tgbotapi.Update, call *msgBuilder.Callback) {
	chatUser := b.userRepository.GetByChatId(u.CallbackQuery.Message.Chat.ID)
	ct := time.Now()
	chatUser.NotDisturbFrom = ct.Format("15:04")
	interval, err := strconv.Atoi(call.Action)

	fmt.Println("TIME", ct.Format("15:04"), interval)
	if err != nil {
		return
	}
	chatUser.NotDisturbInterval = int16(interval)
	_, err = b.userRepository.Update(chatUser)

	if err != nil {
		b.SendError(chatUser.ChatId, err)
	}

	b.SendMsg(tgbotapi.NewMessage(chatUser.ChatId, "✅"))
}

func (b EnglishWordsBot) manageWord(u tgbotapi.Update, param string) {
	entity := b.wordRepository.GetById(param)
	msg := tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, entity.Value+" - "+entity.Translation)
	editCallback := &msgBuilder.Callback{Key: "Edit", Type: "edit", Action: param}
	deleteCallback := &msgBuilder.Callback{Key: "Delete", Type: "delete", Action: param}
	msgBuilder.BuildKeyboardByCallbacks(&msg, []*msgBuilder.Callback{editCallback, deleteCallback})

	b.SendMsg(msg)
}

func (b EnglishWordsBot) deleteWord(u tgbotapi.Update, param string) {
	err := b.wordRepository.DeleteById(param)

	if err != nil {
		b.SendError(u.CallbackQuery.Message.Chat.ID, err)
	}

	msg := tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, "Deleted")

	b.SendMsg(msg)
}

func (b EnglishWordsBot) editWord(u tgbotapi.Update) {
	msg := tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, "🖕🏼")

	b.SendMsg(msg)
}

func (b EnglishWordsBot) callbackTranslation(u tgbotapi.Update) {
	req := strings.Split(u.CallbackQuery.Data, "/")
	param := req[1]
	chatUser := b.userRepository.GetByChatId(u.CallbackQuery.Message.Chat.ID)

	if param == "MyOption" {
		chatUser.WaitingType = u.CallbackQuery.Message.Text
		_, _ = b.userRepository.Update(chatUser)
		b.SendMsg(tgbotapi.NewMessage(chatUser.ChatId, "Write your translation"))

		return
	}

	//TODO
	chatUser.WaitingType = ""
	_, err := b.userRepository.Update(chatUser)

	if err != nil {
		b.SendMsg(tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, "Some error. Sorry..."))
	}

	valueLang, err := languageDetector.Detect(u.CallbackQuery.Message.Text, chatUser.GetUserLangs())
	transLang, err := languageDetector.Detect(param, chatUser.GetUserLangs())

	wordEntity, err := b.wordService.AddWord(&word.Word{
		Value:           req[2],
		ValueLang:       valueLang,
		Translation:     param,
		TranslationLang: transLang,
		ChatId:          chatUser.ChatId,
	})

	if err != nil {
		b.SendError(chatUser.ChatId, err)
	}

	b.SendMsg(tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, wordEntity.Value+" - "+wordEntity.Translation))
}

func (b EnglishWordsBot) handleAnswer(u tgbotapi.Update) {
	call := msgBuilder.CallbackStringToData(u.CallbackData())
	chatId := u.CallbackQuery.Message.Chat.ID
	chatUser := b.userRepository.GetByChatId(chatId)
	w := b.wordRepository.GetByChatIdAndValue(chatId, u.CallbackQuery.Message.Text)

	if w.Translation == call.Action {

		w.Rate++
		rateStr := fmt.Sprintf("rate: %d/%d", w.Rate, chatUser.MaxRate)
		_, err := b.wordRepository.Update(w)

		if err != nil {
			return
		}
		msg := tgbotapi.NewEditMessageText(
			chatId,
			u.CallbackQuery.Message.MessageID,
			w.Value+" - "+w.Translation+" 😀 "+rateStr,
		)
		if _, err := b.api.Send(msg); err != nil {
			panic(err)
		}

		return
	}

	w.Rate--
	rateStr := fmt.Sprintf("rate: %d/%d", w.Rate, chatUser.MaxRate)
	msg := tgbotapi.NewEditMessageText(
		chatId,
		u.CallbackQuery.Message.MessageID,
		w.Value+" - "+w.Translation+" 👺 "+rateStr,
	)

	_, _ = b.wordRepository.Update(w)

	if _, err := b.api.Send(msg); err != nil {
		panic(err)
	}
}
