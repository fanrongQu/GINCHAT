package service

import (
	"gimchat/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SetContact
// @Summary 设置用户关系
// @Schemes
// @Description 设置用户关系
// @Tags 用户关系
// @param token header string true "当前用户token"
// @param targetID query string true "用户ID"
// @param type query string true "用户关系"
// @param remarks query string true "备注信息"
// @Produce json
// @Success 200 {object} integer 成功
// @Router /contact/setContact [post]
func SetContact(c *gin.Context) {
	token := c.GetHeader("token")
	user, error := models.FindUserByToken(token)
	if error != nil {
		c.JSON(http.StatusOK, gin.H{
			"state":       true,
			"needRelogin": true,
			"message":     "token已过期",
		})
		return
	}
	id, error := strconv.ParseInt(c.Query("targetID"), 10, 64)
	if error != nil {
		c.JSON(http.StatusAccepted, gin.H{
			"state":   false,
			"message": "缺少targetID",
		})
		return
	}
	targetUser, error := models.FindUserByIdentity(id)
	if error != nil {
		c.JSON(http.StatusAccepted, gin.H{
			"state":   false,
			"message": "用户不存在",
		})
		return
	}
	typeId, error := strconv.Atoi(c.Query("type"))
	if error != nil {
		c.JSON(http.StatusAccepted, gin.H{
			"state":   false,
			"message": "关系类型错误",
		})
		return
	}

	remarks := c.Query("remarks")
	content := models.ContactBasic{}
	content.OwnerUserId = user.Identity
	content.TargetUserId = targetUser.Identity
	content.TargetUser = targetUser
	content.Type = typeId
	content.Remarks = remarks
	dbContact, error := models.FindContact(user.Identity, targetUser.Identity)
	if error == nil { // 存在关系，更新关系类型
		dbContact, error = models.UpdateContact(dbContact, content)
		if error != nil {
			c.JSON(http.StatusAccepted, gin.H{
				"state":   false,
				"message": "更新用户关系失败",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"state": true,
			"data":  dbContact,
		})
	} else { // 不存在关系，创建关系类型
		dbContact, error = models.CreateContact(content)
		if error != nil {
			c.JSON(http.StatusAccepted, gin.H{
				"state":   false,
				"message": "设置用户关系失败",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"state": true,
			"data":  dbContact,
		})
	}
}

// FindContact
// @Summary 查找用户关系
// @Schemes
// @Description 查找用户关系
// @Tags 用户关系
// @param token header string true "当前用户token"
// @param targetID query string true "用户ID"
// @Produce json
// @Success 200 {object} integer 成功
// @Router /contact/findContact [get]
func FindContact(c *gin.Context) {
	token := c.GetHeader("token")
	user, error := models.FindUserByToken(token)
	if error != nil {
		c.JSON(http.StatusOK, gin.H{
			"state":       true,
			"needRelogin": true,
			"message":     "token已过期",
		})
		return
	}
	id, error := strconv.ParseInt(c.Query("targetID"), 10, 64)
	if error != nil {
		c.JSON(http.StatusAccepted, gin.H{
			"state":   false,
			"message": "缺少targetID",
		})
		return
	}

	dbContact, error := models.FindContact(user.Identity, id)
	if error != nil {
		// 陌生人关系
		c.JSON(http.StatusOK, gin.H{
			"state": true,
			"data":  map[string]any{"ownerUserId": user.Identity, "targetUserId": id, "targetUser": user, "type": 0, "remarks": ""},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"state": true,
		"data":  dbContact,
	})
}

// FindAllContact
// @Summary 查找全部好友
// @Schemes
// @Description 查找全部好友
// @Tags 用户关系
// @param token header string true "当前用户token"
// @Produce json
// @Success 200 {object} integer 成功
// @Router /contact/findAllContact [get]
func FindAllContact(c *gin.Context) {
	token := c.GetHeader("token")
	user, error := models.FindUserByToken(token)
	if error != nil {
		c.JSON(http.StatusOK, gin.H{
			"state":       true,
			"needRelogin": true,
			"message":     "token已过期",
		})
		return
	}

	contacts, error := models.FindAllContact(user.Identity)
	if error != nil {
		// 没有设置过用户关系
		c.JSON(http.StatusOK, gin.H{
			"state": true,
			"data":  []string{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"state": true,
		"data":  contacts,
	})
}
