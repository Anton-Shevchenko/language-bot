package userRepository

import (
	"context"
	"go-app/internal/domain/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Repository interface {
	GetByChatId(chatId int64) *user.User
	Create(u *user.User) (*user.User, error)
	Update(u *user.User) (*user.User, error)
	GetByIntervals(intervals []uint16) []*user.User
}

type userRepository struct {
	db *mongo.Client
}

func NewUserRepository(db *mongo.Client) Repository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(u *user.User) (*user.User, error) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel() // releases resources if AddWord completes before timeout elapses
	collection := r.db.Database("words-db").Collection("users")
	_, err := collection.InsertOne(ctx, *u)

	if err != nil {
		panic(err)
	}

	return u, nil
}

func (r *userRepository) Update(u *user.User) (*user.User, error) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel() // releases resources if AddWord completes before timeout elapses
	collection := r.db.Database("words-db").Collection("users")

	_, err := collection.UpdateOne(ctx, bson.M{"chatId": u.ChatId}, bson.M{"$set": u})

	if err != nil {
		return u, err
	}

	return u, nil
}

func (r *userRepository) GetByChatId(chatId int64) *user.User {
	var entity user.User
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	collection := r.db.Database("words-db").Collection("users")
	res := collection.FindOne(ctx, bson.M{"chatId": chatId})
	_ = res.Decode(&entity)

	return &entity
}

func (r *userRepository) GetByIntervals(intervals []uint16) []*user.User {
	var entities []*user.User
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	collection := r.db.Database("words-db").Collection("users")

	cursor, err := collection.Find(ctx, bson.M{"interval": bson.M{"$in": intervals}})
	if err != nil {
		return nil
	}

	err = cursor.All(ctx, &entities)

	return entities
}
