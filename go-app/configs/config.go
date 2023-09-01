package configs

import (
	"go-app/pkg/randomParagraphGenerator"
	"go-app/pkg/randomWordsGenerator"
	"os"
)

const (
	prod = "production"
)

type Config struct {
	Env                            string
	MongoDB                        MongoDBConfig
	Host                           string
	Port                           string
	RandomWordGeneratorConfig      randomWordsGenerator.Config
	RandomParagraphGeneratorConfig randomParagraphGenerator.Config
}

func (c Config) IsProd() bool {
	return c.Env == prod
}

func GetConfig() Config {
	return Config{
		Env:     os.Getenv("ENV"),
		MongoDB: GetMongoDBConfig(),
		Host:    os.Getenv("APP_HOST"),
		Port:    os.Getenv("APP_PORT"),
		RandomWordGeneratorConfig: randomWordsGenerator.Config{
			Url:    os.Getenv("RW_URL"),
			ApiKey: os.Getenv("RW_API_KEY"),
		},
		RandomParagraphGeneratorConfig: randomParagraphGenerator.Config{
			Url: os.Getenv("RP_API_URL"),
		},
	}
}
