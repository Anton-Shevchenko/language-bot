package configs

import (
	"os"
)

type MongoDBConfig struct {
	URI string `env:"MONGO_URI"`
}

func GetMongoDBConfig() MongoDBConfig {
	return MongoDBConfig{
		URI: os.Getenv("MONGO_URI"),
	}
}
