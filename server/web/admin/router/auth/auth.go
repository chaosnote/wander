package auth

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	// 使用者身份驗證
	// 加入至 Session

	_ = sessions.Default(c)
}

func Logout(c *gin.Context) {
	// 使用者身份驗證
	// 加入至 Session
}
