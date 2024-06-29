package wordRepository

import (
	"context"
	"go-app/internal/domain/word"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type WordRepository interface {
	AddWord(word *word.Word) (*word.Word, error)
	GetAllByChatId(chatId int64) ([]*word.Word, error)
	DeleteById(id string) error
	GetById(id string) *word.Word
	GetRandom(chatId int64, maxRate int8) *word.Word
	GetRandomFive(chatId int64, langTo string) []*word.Word
	GetRandomTranslations(w *word.Word) []*word.Word
	Update(w *word.Word) (*word.Word, error)
	GetByChatIdAndValue(chatId int64, value string) *word.Word
	GetByValue(value string) *word.Word
}

type wordRepository struct {
	db *mongo.Client
}

func NewWordRepository(db *mongo.Client) WordRepository {
	return &wordRepository{
		db: db,
	}
}

func (r *wordRepository) Update(w *word.Word) (*word.Word, error) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel() // releases resources if AddWord completes before timeout elapses
	collection := r.db.Database("words-db").Collection("words")
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"id": w.ID}, bson.M{"$set": bson.M{"rate": w.Rate, "translation": w.Translation}},
	)

	if err != nil {
		return w, err
	}

	return w, nil
}

func (r *wordRepository) AddWord(word *word.Word) (*word.Word, error) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel() // releases resources if AddWord completes before timeout elapses
	collection := r.db.Database("words-db").Collection("words")
	word.ID = primitive.NewObjectID()
	_, err := collection.InsertOne(ctx, *word)

	if err != nil {
		panic(err)
	}

	return word, nil
}

func (r *wordRepository) GetAllByChatId(chatId int64) ([]*word.Word, error) {
	var words []*word.Word
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	option := options.Find()
	collection := r.db.Database("words-db").Collection("words")
	option.SetLimit(10)
	option.SetSort(bson.D{{"_id", -1}})
	cur, err := collection.Find(ctx, bson.M{"chatId": chatId}, option)

	for cur.Next(context.TODO()) {
		var w word.Word
		err := cur.Decode(&w)
		if err != nil {
			log.Fatal(err)
		}

		words = append(words, &w)

	}

	if err != nil {
		panic(err)
	}
	return words, nil
}

func (r *wordRepository) GetById(id string) *word.Word {
	var entity word.Word
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	objId, _ := primitive.ObjectIDFromHex(id)
	collection := r.db.Database("words-db").Collection("words")
	res := collection.FindOne(ctx, bson.M{"id": objId})

	err := res.Decode(&entity)

	if err != nil {
		return nil
	}

	return &entity
}

func (r *wordRepository) GetByValue(value string) *word.Word {
	var entity word.Word
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	collection := r.db.Database("words-db").Collection("words")
	res := collection.FindOne(ctx, bson.M{"value": value})

	err := res.Decode(&entity)

	if err != nil {
		return nil
	}

	return &entity
}

func (r *wordRepository) GetByChatIdAndValue(chatId int64, value string) *word.Word {
	var entity word.Word
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	collection := r.db.Database("words-db").Collection("words")
	res := collection.FindOne(ctx, bson.M{"chatId": chatId, "value": value})

	err := res.Decode(&entity)

	if err != nil {
		return nil
	}

	return &entity
}

func (r *wordRepository) DeleteById(id string) error {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	objId, _ := primitive.ObjectIDFromHex(id)
	collection := r.db.Database("words-db").Collection("words")
	_, err := collection.DeleteOne(ctx, bson.M{"id": objId})

	if err != nil {
		panic(err)
	}

	return nil
}

func (r *wordRepository) GetRandom(chatId int64, maxRate int8) *word.Word {
	var entity word.Word
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	collection := r.db.Database("words-db").Collection("words")
	aggregate, err := collection.Aggregate(
		ctx,
		[]bson.M{
			bson.M{
				"$match": bson.M{"chatId": chatId, "rate": bson.M{"$lt": maxRate}},
			},
			bson.M{"$sample": bson.M{"size": 1}},
		},
	)
	if err != nil {
		return nil
	}

	for aggregate.Next(context.TODO()) {
		err := aggregate.Decode(&entity)
		if err != nil {
			log.Fatal(err)
		}

	}

	return &entity
}

func (r *wordRepository) GetRandomFive(chatId int64, langTo string) []*word.Word {
	var entities []*word.Word
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	collection := r.db.Database("words-db").Collection("words")
	aggregate, err := collection.Aggregate(
		ctx,
		[]bson.M{
			bson.M{"$match": bson.M{
				"chatId":    chatId,
				"valueLang": langTo,
			}},
			bson.M{"$sample": bson.M{"size": 5}},
		},
	)
	if err != nil {
		return nil
	}

	for aggregate.Next(context.TODO()) {
		var wrd word.Word
		err := aggregate.Decode(&wrd)
		if err != nil {
			log.Fatal(err)
		}
		entities = append(entities, &wrd)
	}

	return entities
}

func (r *wordRepository) GetRandomTranslations(w *word.Word) []*word.Word {
	var entities []*word.Word

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	collection := r.db.Database("words-db").Collection("words")
	aggregate, err := collection.Aggregate(
		ctx,
		[]bson.M{
			bson.M{"$match": bson.M{
				"chatId":          w.ChatId,
				"translationLang": w.TranslationLang,
				"value":           bson.M{"$ne": w.Value}},
			},
			bson.M{"$sample": bson.M{"size": 3}},
		},
	)
	if err != nil {
		return nil
	}

	for aggregate.Next(context.TODO()) {
		var wrd word.Word
		err := aggregate.Decode(&wrd)
		if err != nil {
			log.Fatal(err)
		}

		entities = append(entities, &wrd)
	}

	return entities
}
