package tts

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/HardDie/DeckBuilder/internal/apperr"
	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/logger"
)

type tts struct {
	editorAddr string

	// mu guards dataForTTS and httpPort.
	// SendToTTS runs on the generate or binding goroutine, DataForTTS on the HTTP one.
	mu         sync.Mutex
	dataForTTS []byte
	httpPort   int
}

type Message struct {
	MessageID int    `json:"messageID"`
	GUID      string `json:"guid"`
	Script    string `json:"script"`
}

func New() TTS {
	return &tts{
		httpPort:   config.HTTPPort,
		editorAddr: "127.0.0.1:39999",
	}
}

func (s *tts) SetHTTPPort(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.httpPort = port
}

func (s *tts) SendToTTS(data any) {
	// Try to open TCP socket
	conn, err := net.Dial("tcp", s.editorAddr)
	if err != nil {
		logger.Info.Printf("Can't connect to TTS via tcp connection %q: %s", s.editorAddr, err.Error())
		return
	}
	defer func() { conn.Close() }()

	dataForTTS, err := json.Marshal(data)
	if err != nil {
		logger.Warn.Println("error marshal data for TTS:", err.Error())
		return
	}
	s.mu.Lock()
	s.dataForTTS = dataForTTS
	httpPort := s.httpPort
	s.mu.Unlock()

	msg := Message{
		MessageID: 3,
		GUID:      "-1",
		Script: fmt.Sprintf(`
WebRequest.get("http://%s:%d/api/tts/data", function(request)
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
end)`, config.HTTPHost, httpPort),
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

func (s *tts) DataForTTS() ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.dataForTTS == nil {
		return nil, apperr.ErrNothingForTTS
	}
	res := s.dataForTTS
	s.dataForTTS = nil
	return res, nil
}
