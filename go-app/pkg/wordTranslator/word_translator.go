package wordTranslator

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/robertkrimen/otto"
	"io"
	"net/http"
	"strings"
)

func encodeURI(s string) (string, error) {
	eUri := `eUri = encodeURI(sourceText);`
	vm := otto.New()
	err := vm.Set("sourceText", s)
	if err != nil {
		return "err", errors.New("Error setting js variable")
	}
	_, err = vm.Run(eUri)
	if err != nil {
		return "err", errors.New("Error executing jscript")
	}
	val, err := vm.Get("eUri")
	if err != nil {
		return "err", errors.New("Error getting variable value from js")
	}
	v, err := val.ToString()
	if err != nil {
		return "err", errors.New("Error converting js var to string")
	}
	return v, nil
}

func Translate(source, sourceLang, targetLang string) ([]string, error) {
	var text []string
	var result []interface{}

	encodedSource, err := encodeURI(source)
	fmt.Println("JJJJJJ", source, encodedSource)
	if err != nil {
		return []string{"err"}, err
	}
	url := "https://translate.googleapis.com/translate_a/single?client=gtx&dt=t&dt=at&sl=" +
		sourceLang + "&tl=" + targetLang + "&dt=t&q=" + encodedSource

	r, err := http.Get(url)
	if err != nil {
		return []string{}, errors.New("Error getting translate.googleapis.com")
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return []string{}, errors.New("Error reading response body")
	}

	bReq := strings.Contains(string(body), `<title>Error 400 (Bad Request)`)
	if bReq {
		return []string{}, errors.New("Error 400 (Bad Request)")
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return []string{}, errors.New("Error unmarshaling data")
	}

	if len(result) > 0 {
		inner := result[5]

		fmt.Println("INNER", inner, result, url)

		for _, sliceOne := range inner.([]interface{}) {
			for i, sliceTwo := range sliceOne.([]interface{}) {
				if i == 2 {
					for _, sliceThree := range sliceTwo.([]interface{}) {
						for _, translate := range sliceThree.([]interface{}) {
							t := fmt.Sprintf("%v", translate)
							text = append(text, strings.ToLower(t))
							break
						}
					}
				}

			}
			break
		}

		return text, nil
	} else {
		return []string{}, errors.New("No translated data in responce")
	}
}
