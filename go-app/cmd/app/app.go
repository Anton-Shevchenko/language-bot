package app

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go-app/configs"
	"go-app/internal/englishWordsBot"
	"go-app/internal/repositories/userRepository"
	"go-app/internal/repositories/wordRepository"
	"go-app/internal/services/wordService"
	"go-app/pkg/jobManager"
	"go-app/pkg/randomParagraphGenerator"
	"go-app/pkg/randomWordsGenerator"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var (
	r = gin.Default()
)

func Run() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	config := configs.GetConfig()

	clientOptions := options.Client().ApplyURI(config.MongoDB.URI)
	mongoDB, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		panic(err)
	}

	randWordsGenerator := randomWordsGenerator.NewRandomWordsGenerator(config.RandomWordGeneratorConfig)
	paragraphGenerator := randomParagraphGenerator.NewRandomParagraphGenerator(config.RandomParagraphGeneratorConfig)

	wordRepo := wordRepository.NewWordRepository(mongoDB)
	userRepo := userRepository.NewUserRepository(mongoDB)

	newWordService := wordService.NewWordService(wordRepo, randWordsGenerator, paragraphGenerator)

	botTokens := configs.GetBotsTokens()
	englishBot := englishWordsBot.NewEnglishBot(
		configs.BaseBotConfig{
			Debug: true,
			Token: botTokens.EnglishBotToken,
		},
		wordRepo,
		newWordService,
		userRepo,
	)
	go englishBot.Build()
	manager := jobManager.NewJobManager()
	//Jobs
	wordJob := englishWordsBot.NewWordJob(wordRepo, userRepo, englishBot)

	manager.Add(wordJob.WordJob, "word-job", 1*time.Minute)

	manager.Scheduler()
	r.Run()
}
