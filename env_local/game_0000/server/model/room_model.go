package model

import (
	"github.com/chaosnote/wander/model/member"
)

type ViewModel struct {
	// 贏分
	// ∟ 線獎
	// ∟ 其它獎項
}

type RoomModel struct {
	Player member.Player // 玩家進入時更新

	// 共用資訊
	// ∟ 押注資訊

	// ViewModel
	// ∟ N(ormal)G(ame)
	// ∟ F(ree)G(ame)
}
