package version

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Output(c *gin.Context) {
	// 使用者身份驗證
	// 加入至 Session

	c.String(http.StatusOK, "Server v0.0.0")
}
