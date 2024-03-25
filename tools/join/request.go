package main

import "github.com/HardDie/DeckBuilder/internal/tts_entity"

type Root struct {
	ObjectStates []GameBag `json:"ObjectStates"`
}

type GameBag struct {
	Name             string               `json:"Name"`
	Transform        tts_entity.Transform `json:"Transform"`
	Nickname         string               `json:"Nickname"`
	Description      string               `json:"Description"`
	ContainedObjects []CollectionBag      `json:"ContainedObjects"`
}

type CollectionBag struct {
	Name             string                 `json:"Name"`
	Transform        tts_entity.Transform   `json:"Transform"`
	Nickname         string                 `json:"Nickname"`
	Description      string                 `json:"Description"`
	ContainedObjects []tts_entity.TTSObject `json:"ContainedObjects"`
}
