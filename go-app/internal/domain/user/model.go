package user

type User struct {
	LangFrom           string `json:"lang_from" bson:"langFrom"`
	LangTo             string `json:"lang_to" bson:"langTo"`
	Level              string `json:"level" bson:"level"`
	MaxRate            int8   `json:"max_rate" bson:"maxRate"`
	Interval           uint16 `json:"interval" bson:"interval"`
	ChatId             int64  `json:"chat_id" bson:"chatId"`
	NotDisturbFrom     string `json:"not_disturb_from" bson:"notDisturbFrom"`
	NotDisturbInterval int16  `json:"not_disturb_interval" bson:"notDisturbInterval"`
	WaitingType        string `json:"waiting_type" bson:"waitingType"`
}

type interval uint16

const (
	Interval30  interval = 30
	Interval60  interval = 60
	Interval120 interval = 120
	Interval180 interval = 180
)

func GetIntervals() [4]interval {
	return [...]interval{Interval30, Interval60, Interval120, Interval180}
}

func (u *User) GetUserLangs() []string {
	return []string{u.LangTo, u.LangFrom}
}
