package model

type bot struct {
	Mode mode
}

func NewBot() *bot {
	return &bot{
		Mode: MODE_1,
	}
}
