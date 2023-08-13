package user

type User struct {
	LangFrom           string   `json:"lang_from" bson:"langFrom"`
	LangTo             string   `json:"lang_to" bson:"langTo"`
	Level              string   `json:"level" bson:"level"`
	MaxRate            int8     `json:"max_rate" bson:"maxRate"`
	Interval           Interval `json:"interval" bson:"interval"`
	ChatId             int64    `json:"chat_id" bson:"chatId"`
	NotDisturbFrom     string   `json:"not_disturb_from" bson:"notDisturbFrom"`
	NotDisturbInterval int16    `json:"not_disturb_interval" bson:"notDisturbInterval"`
	WaitingType        string   `json:"waiting_type" bson:"waitingType"`
}

type Interval uint16

const (
	Interval2   Interval = 2
	Interval30  Interval = 30
	Interval60  Interval = 60 //default
	Interval120 Interval = 120
	Interval180 Interval = 180
)

func GetIntervals() [5]Interval {
	return [...]Interval{Interval2, Interval30, Interval60, Interval120, Interval180}
}

func (u *User) GetUserLangs() []string {
	return []string{u.LangTo, u.LangFrom}
}
