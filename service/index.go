package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetIndex
// @Summary 获取配置index
// @Schemes
// @Description 获取配置index信息
// @Tags 首页
// @Accept json
// @Produce json
// @Success 200 {object} integer 成功
// @Router /index [get]
func GetIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "get index success",
	})
}
