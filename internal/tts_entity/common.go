package tts_entity

var (
	transform = Transform{
		ScaleX: 1,
		ScaleY: 1,
		ScaleZ: 1,
	}
)

type TTSObject interface {
	GetName() string
	GetNickname() string
}
