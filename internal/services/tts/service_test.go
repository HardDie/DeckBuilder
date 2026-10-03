package tts

import (
	"encoding/json"
	"io"
	"net"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/HardDie/DeckBuilder/internal/apperr"
)

// fakeEditor stands in for the TTS External Editor port.
// Every connection body is sent to messages.
func fakeEditor(t *testing.T) (string, <-chan []byte) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = ln.Close() })

	messages := make(chan []byte, 128)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			body, _ := io.ReadAll(conn)
			_ = conn.Close()
			select {
			case messages <- body:
			default:
			}
		}
	}()
	return ln.Addr().String(), messages
}

func newTestTTS(addr string) *tts {
	s := New().(*tts)
	s.editorAddr = addr
	return s
}

func TestSendToTTS(t *testing.T) {
	addr, messages := fakeEditor(t)
	s := newTestTTS(addr)
	s.SetHTTPPort(5003)

	s.SendToTTS(map[string]string{"Name": "Bag"})

	var msg Message
	require.NoError(t, json.Unmarshal(<-messages, &msg))
	assert.Equal(t, 3, msg.MessageID)
	assert.Equal(t, "-1", msg.GUID)
	assert.True(t, strings.Contains(msg.Script, "http://127.0.0.1:5003/api/tts/data"), msg.Script)

	data, err := s.DataForTTS()
	require.NoError(t, err)
	assert.JSONEq(t, `{"Name":"Bag"}`, string(data))

	_, err = s.DataForTTS()
	assert.ErrorIs(t, err, apperr.ErrNothingForTTS, "data is served once")
}

func TestSendToTTSNoEditor(t *testing.T) {
	// Closed port: nothing listens there any more.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	s := newTestTTS(ln.Addr().String())
	require.NoError(t, ln.Close())

	s.SendToTTS(map[string]string{"Name": "Bag"})

	_, err = s.DataForTTS()
	assert.Error(t, err, "nothing is stored when TTS is not reachable")
}

func TestSendToTTSConcurrent(t *testing.T) {
	addr, _ := fakeEditor(t)
	s := newTestTTS(addr)

	const workers = 20
	var wg sync.WaitGroup
	for range workers {
		wg.Add(2)
		go func() {
			defer wg.Done()
			s.SendToTTS(map[string]string{"Name": "Bag"})
		}()
		go func() {
			defer wg.Done()
			_, _ = s.DataForTTS()
		}()
	}
	wg.Wait()
}
