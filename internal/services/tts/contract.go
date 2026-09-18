package tts

type TTS interface {
	SetHTTPPort(port int)
	SendToTTS(data any)
	DataForTTS() ([]byte, error)
}
