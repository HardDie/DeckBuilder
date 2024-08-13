package tts

type TTS interface {
	SendToTTS(data any)
	DataForTTS() ([]byte, error)
}
