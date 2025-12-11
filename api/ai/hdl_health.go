package ai

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// hdlHealth 健康检查
func (a *App) hdlHealth(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
