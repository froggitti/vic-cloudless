package vtr

import "strings"

func ProcessTextAll(voiceText string, intents []JsonIntent) (string, map[string]string, bool) {
	var botSerial string
	var intentNum int = 0
	var successMatched bool = false
	voiceText = strings.ToLower(voiceText)
	// Look for a perfect match first
	for _, b := range intents {
		for _, c := range b.Keyphrases {
			if voiceText == strings.ToLower(c) {
				return ParamChecker(b.Name, voiceText, botSerial)
			}
		}
		intentNum = intentNum + 1
	}
	// Not found? Then let's be happy with a bare substring search
	if !successMatched {
		intentNum = 0
		for _, b := range intents {
			for _, c := range b.Keyphrases {
				if strings.Contains(voiceText, strings.ToLower(c)) && !b.RequireExactMatch {
					return ParamChecker(b.Name, voiceText, botSerial)
				}
			}
			intentNum = intentNum + 1
		}
	}
	return "intent_system_noaudio", map[string]string{}, false
}
