package core

type Core interface {
	Init() error
	Drop() error
}
