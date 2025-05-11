package stream

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	vosk "github.com/os-vector/vosk-api/go"
)

var model *vosk.VoskModel
var rec *vosk.VoskRecognizer

func InitVosk() {
	loadIntents()
	var err error
	model, err = vosk.NewModel("/data/en-US")
	if err != nil {
		log.Fatal("model not found", err)
	}
	rec, err = vosk.NewRecognizerGrm(model, 16000, GetGrammerList("en-US"))
	if err != nil {
		log.Fatal("error making rec:", err)
	}
	rec.SetMaxAlternatives(0)
	rec.SetEndpointerDelays(3, 0, 0)
}

func Process(chunk []byte) string {
	if len(chunk) == 0 {
		fmt.Println("empty ahh chunk")
		return ""
	}
	stop, isActive := DetectEndOfSpeech(chunk)
	if isActive {
		rec.AcceptWaveform(chunk)
	}
	if stop {
		var jres map[string]interface{}
		json.Unmarshal([]byte(rec.FinalResult()), &jres)
		transcribedText := jres["text"].(string)
		fmt.Println("transcribed text: " + transcribedText)
		return transcribedText
	}
	return ""
}

var NumbersEN_US []string = []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety", "hundred", "seconds", "minutes", "hours", "minute", "second", "hour"}

func removeDuplicates(strings []string) []string {
	occurred := map[string]bool{}
	var result []string
	for _, str := range strings {
		if !occurred[str] {
			result = append(result, str)
			occurred[str] = true
		}
	}
	return result
}
func GetGrammerList(lang string) string {
	var wordsList []string
	var grammer string
	// add words in intent json
	for _, words := range IntentList {
		for _, word := range words.Keyphrases {
			wors := strings.Split(word, " ")
			for _, wor := range wors {
				found := model.FindWord(wor)
				if found != -1 {
					wordsList = append(wordsList, wor)
				}
			}
		}
	}
	// add words in localization
	for _, str := range ALL_STR {
		text := GetText(str)
		wors := strings.Split(text, " ")
		for _, wor := range wors {
			found := model.FindWord(wor)
			if found != -1 {
				wordsList = append(wordsList, wor)
			}
		}
	}
	// add numbers
	for _, wor := range NumbersEN_US {
		found := model.FindWord(wor)
		if found != -1 {
			wordsList = append(wordsList, wor)
		}
	}

	wordsList = removeDuplicates(wordsList)
	for i, word := range wordsList {
		if i == len(wordsList)-1 {
			grammer = grammer + `"` + word + `"`
		} else {
			grammer = grammer + `"` + word + `"` + ", "
		}
	}
	grammer = "[" + grammer + "]"
	return grammer
}

type JsonIntent struct {
	Name              string   `json:"name"`
	Keyphrases        []string `json:"keyphrases"`
	RequireExactMatch bool     `json:"requiresexact"`
}

var IntentList []JsonIntent

func loadIntents() {
	file, err := os.ReadFile("/data/intent-data.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(file, &IntentList)
	if err != nil {
		fmt.Println(err)
	}
}
