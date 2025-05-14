package vtr

import (
	"encoding/json"
	"fmt"
	"log"

	vosk "github.com/os-vector/vosk-api/go"
)

var model *vosk.VoskModel
var rec *vosk.VoskRecognizer

func InitVosk() {
	loadIntents()
	var err error
	model, err = vosk.NewModel("/anki/data/assets/cozmo_resources/cloudless/en-US/model")
	if err != nil {
		log.Fatal("model not found", err)
	}
	rec, err = vosk.NewRecognizerGrm(model, 16000, GetGrammerList("en-US"))
	if err != nil {
		log.Fatal("error making rec:", err)
	}
	// does this actually do anything
	rec.SetMaxAlternatives(0)
	rec.SetPartialWords(0)
	rec.SetWords(0)
	rec.SetEndpointerDelays(3, 0, 0)
}

func Process(chunk []byte) string {
	if len(chunk) == 0 {
		fmt.Println("empty ahh chunk")
		return ""
	}
	// todo: experiment with giving acceptwaveform smaller or bigger chunks
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
