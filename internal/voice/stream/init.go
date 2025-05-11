package stream

import (
	"strings"

	chippergrpc2 "github.com/digital-dream-labs/api/go/chipperpb"
)

func (strm *Streamer) init(streamSize int) {
	// set up error response if context times out/is canceled
	go strm.cancelResponse()

	// // start routine to buffer communication between main routine and upload routine
	go strm.bufferRoutine(streamSize)
	// if strm.opts.checkOpts != nil {
	// 	go strm.testRoutine(streamSize)
	// }

	// // connect to server
	// var err *CloudError
	// if strm.conn, err = strm.opts.connectFn(strm.ctx); err != nil {
	// 	strm.receiver.OnError(err.Kind, err.Err)
	// 	strm.cancel()
	// 	return
	// }

	// start routine to upload audio via GRPC until response or error
	// go func() {
	// 	responseInited := false
	// 	for data := range strm.audioStream {
	// 		if err := strm.sendAudio(data); err != nil {
	// 			return
	// 		}
	// 		if !responseInited {
	// 			go strm.responseRoutine()
	// 			responseInited = true
	// 		}
	// 	}
	// }()

	go func() {
		for data := range strm.audioStream {
			text := Process(data)
			if text != "" {
				intent, iParam, _ := ProcessTextAll(text, IntentList)
				sendIntentGraphResponse(&chippergrpc2.IntentGraphResponse{
					ResponseType: chippergrpc2.IntentGraphMode_INTENT,
					IsFinal:      true,
					IntentResult: &chippergrpc2.IntentResult{
						Action:     intent,
						Parameters: iParam,
					},
				}, strm.receiver)
				return
			}
		}
	}()
}

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
