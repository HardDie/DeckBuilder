package service

import (
	"encoding/json"
	"errors"
	"net"

	"github.com/HardDie/DeckBuilder/internal/logger"
)

type ITTSService interface {
	SendToTTS(data any)
	DataForTTS() ([]byte, error)
}

type TTSService struct {
	dataForTTS []byte
}

type TTSMessage struct {
	MessageID int    `json:"messageID"`
	GUID      string `json:"guid"`
	Script    string `json:"script"`
}

func NewTTSService() *TTSService {
	return &TTSService{}
}

func (s *TTSService) SendToTTS(data any) {
	// Try to open TCP socket
	conn, err := net.Dial("tcp", "127.0.0.1:39999")
	if err != nil {
		logger.Info.Println("Can't connect to TTS via tcp connection '127.0.0.1:39999':", err.Error())
		return
	}
	defer func() { conn.Close() }()

	dataForTTS, err := json.Marshal(data)
	if err != nil {
		logger.Warn.Println("error marshal data for TTS:", err.Error())
		return
	}
	s.dataForTTS = dataForTTS

	msg := TTSMessage{
		MessageID: 3,
		GUID:      "-1",
		Script: `
WebRequest.get("http://127.0.0.1:5000/api/tts/data", function(request)
	if request.is_error then
		print('Downloading json error: ', request.error)
		return
	end
	print('JSON were downloaded!')
	spawnObjectJSON({
		json = request.text,
		callback_function = function(spawned_object)
			print('Object were spawned! Done!')
		end
	})
end)`,
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		logger.Warn.Println("error marshal msg for TTS:", err.Error())
		return
	}

	_, err = conn.Write(jsonData)
	if err != nil {
		logger.Warn.Println("error write message into TTS socket:", err.Error())
		return
	}
}

func (s *TTSService) DataForTTS() ([]byte, error) {
	if s.dataForTTS == nil {
		return nil, errors.New("there is nothing to serve")
	}
	res := s.dataForTTS
	s.dataForTTS = nil
	return res, nil
}
