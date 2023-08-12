package word

import "go.mongodb.org/mongo-driver/bson/primitive"

type Word struct {
	ID              primitive.ObjectID
	Value           string `json:"value" bson:"value"`
	ValueLang       string `json:"value_lang" bson:"valueLang"`
	Translation     string `json:"translation" bson:"translation"`
	TranslationLang string `json:"translation_lang" bson:"translationLang"`
	Rate            int8   `json:"rate" bson:"rate"`
	ChatId          int64  `json:"chat_id" bson:"chatId"`
}
