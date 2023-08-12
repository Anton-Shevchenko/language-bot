package configs

import (
	"go-app/pkg/randomWordsGenerator"
	"os"
)

const (
	prod = "production"
)

type Config struct {
	Env                       string
	MongoDB                   MongoDBConfig
	Host                      string
	Port                      string
	RandomWordGeneratorConfig randomWordsGenerator.RandomWordGeneratorConfig
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
		RandomWordGeneratorConfig: randomWordsGenerator.RandomWordGeneratorConfig{
			Url:    os.Getenv("RW_URL"),
			ApiKey: os.Getenv("RW_API_KEY"),
		},
	}
}
