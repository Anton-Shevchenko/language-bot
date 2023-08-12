package configs

import "os"

type BaseBotConfig struct {
	Token string
	Debug bool
}

type BotTokens struct {
	EnglishBotToken string `env:"BOT_TOKEN"`
}

func GetBotsTokens() *BotTokens {
	return &BotTokens{
		EnglishBotToken: os.Getenv("BOT_TOKEN"),
	}
}
