package englishWordsBot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go-app/internal/domain/user"
	"go-app/internal/domain/word"
	"go-app/pkg/bot/msgBuilder"
	"math/rand"
	"time"
)

type jobWordRepository interface {
	GetRandom(chatId int64) *word.Word
	GetRandomTranslations(w *word.Word) []*word.Word
}

type jobUserRepository interface {
	GetByIntervals(intervals []uint16) []*user.User
}

type WordJob struct {
	userRepository       jobUserRepository
	wordRepository       jobWordRepository
	intervals            []uint16
	currentIntervalIndex uint8
	bot                  *EnglishWordsBot
}

func NewWordJob(wordRepository jobWordRepository, userRepository jobUserRepository, bot *EnglishWordsBot) *WordJob {
	return &WordJob{
		userRepository:       userRepository,
		wordRepository:       wordRepository,
		intervals:            []uint16{30, 60, 120, 180},
		currentIntervalIndex: 0,
		bot:                  bot,
	}
}

func (j *WordJob) WordJob() {
	intervals := j.getCurrentIntervals()

	if len(intervals) == 0 {
		return
	}

	us := j.userRepository.GetByIntervals(intervals)

	for _, u := range us {
		if j.checkIsNotDisturbTime(u) {
			return
		}
		go j.SendWord(u)
	}
}

func (j *WordJob) checkIsNotDisturbTime(u *user.User) bool {
	if u.NotDisturbFrom == "" {
		return false
	}
	parse, err := time.Parse("15:04", u.NotDisturbFrom)
	if err != nil {
		return true
	}

	to := parse.Add(time.Duration(u.NotDisturbInterval/60) * time.Minute)

	if to.Sub(parse) >= time.Duration(u.NotDisturbInterval/60)*time.Minute {
		return false
	}

	return true
}

func (j *WordJob) getCurrentIntervals() []uint16 {
	var currentInterval []uint16
	nowMinute := uint16(time.Now().Minute())

	for _, interval := range user.GetIntervals() {
		if nowMinute%uint16(interval) == 0 {
			currentInterval = append(currentInterval, uint16(interval))
		}
	}

	return currentInterval
}

func (j *WordJob) SendWord(u *user.User) {
	w := j.wordRepository.GetRandom(u.ChatId)
	var calls []*msgBuilder.Callback
	trans := j.wordRepository.GetRandomTranslations(w)
	trans = append(trans, w)

	for _, t := range trans {
		calls = append(calls, &msgBuilder.Callback{
			Key:    t.Translation,
			Type:   "answer",
			Action: t.Translation,
		})
	}

	msg := tgbotapi.NewMessage(u.ChatId, w.Value)
	msgBuilder.BuildKeyboardByCallbacks(&msg, calls)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(calls), func(i, j int) { calls[i], calls[j] = calls[j], calls[i] })

	j.bot.SendMsg(msg)
}
