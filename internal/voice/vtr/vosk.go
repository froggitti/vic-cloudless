package vtr

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

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
	//rec, err = vosk.NewRecognizer(model, 16000)
	if err != nil {
		log.Fatal("error making rec:", err)
	}
	// does this actually do anything
	rec.SetMaxAlternatives(0)
	rec.SetEndpointerDelays(3, 0, 0)
}

func Process(chunk []byte) string {
	if len(chunk) == 0 {
		fmt.Println("empty chunk")
		return ""
	}
	// todo: experiment with giving acceptwaveform smaller or bigger chunks
	stop, _ := DetectEndOfSpeech(chunk)
	rec.AcceptWaveform(chunk)
	if stop {
		var jres map[string]interface{}
		json.Unmarshal([]byte(rec.FinalResult()), &jres)
		transcribedText := jres["text"].(string)
		fmt.Println("transcribed text: " + transcribedText)
		go rec.Reset()
		return transcribedText
	}
	return ""
}

func GetFreq() string {
	file, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cpufreq/scaling_max_freq")
	if err == nil {
		return strings.TrimSpace(string(file))
	}
	return "533333"
}

func SetFreq(cpu, ram string) {
	go exec.Command("/usr/bin/sudo", "/usr/sbin/setfreq", cpu, ram).Run()
}
