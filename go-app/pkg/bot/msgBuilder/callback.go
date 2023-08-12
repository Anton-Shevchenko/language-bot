package msgBuilder

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Callback struct {
	Key     string
	Type    string
	Action  string
	Options map[string]string
}

func CallbackDataToString(call *Callback) string {
	options, _ := json.Marshal(call.Options)

	return fmt.Sprintf("%s/%s/%s", call.Type, call.Action, string(options))
}

func CallbackStringToData(call string) *Callback {
	var options map[string]string

	parts := strings.Split(call, "/")

	if len(parts) > 2 {
		_ = json.Unmarshal([]byte(parts[2]), &options)
	}

	return &Callback{
		Type:    parts[0],
		Action:  parts[1],
		Options: options,
	}
}
