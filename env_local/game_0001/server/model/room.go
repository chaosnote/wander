package model

import "github.com/looplab/fsm"

type GameRoom struct {
	ID   string
	Name string

	// 上下限
	// 自開房 - 自定義密碼

	state *fsm.FSM
}
