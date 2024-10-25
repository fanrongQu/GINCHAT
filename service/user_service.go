package service

import (
	"fmt"
	"gimchat/models"
	"gimchat/utils"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// FindUserList
// @Summary 获取用户列表
// @Schemes
// @Description 获取用户列表信息
// @Tags 用户
// @param userIds query string false "需要查询的用户id字符串(例如xxx,xxx)"
// @Accept json
// @Produce json
// @Success 200 {object} integer 成功
// @Router /user/findUserList [get]
func FindUserList(c *gin.Context) {
	userIds := c.Query("userIds")
	idList := strings.Split(userIds, ",")
	intIdList := make([]int64, 10)
	for i := 0; i < len(idList); i++ {
		intId, err := strconv.ParseInt(idList[i], 10, 64)
		if err == nil {
			intIdList = append(intIdList, int64(intId))
		}
	}
	data, err := models.FindUserList(intIdList)
	if err != nil {
		c.JSON(http.StatusAccepted, gin.H{
			"state":   false,
			"message": "获取用户列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"state":   true,
		"message": data,
	})
}

// CreateUser
// @Summary 创建用户
// @Schemes
// @Description 创建用户
// @Tags 用户
// @param name query string false "用户名"
// @param email query string false "邮箱"
// @param phone query string true "手机号"
// @param passWord query string true "密码"
// @param rePassWord query string true "确认密码"
// @Produce json
// @Success 200 {object} integer 成功
// @Router /user/registed [post]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Email = c.Query("email")

	name := c.Query("name")
	phone := c.Query("phone")
	password := c.Query("passWord")
	rePassword := c.Query("rePassWord")

	if len(phone) == 0 {
		c.JSON(http.StatusAccepted, gin.H{
			"state":   false,
			"message": "输入手机号码不能为空",
		})
		return
	}
	if len(password) == 0 {
		c.JSON(http.StatusAccepted, gin.H{
			"state":   false,
			"message": "输入密码不能为空",
		})
		return
	}
	if password != rePassword {
		c.JSON(http.StatusAccepted, gin.H{
			"state":   false,
			"message": "密码不一致",
		})
		return
	}
	_, err := models.FindUserByPhone(phone)
	if err == nil {
		c.JSON(http.StatusAccepted, gin.H{
			"state":   false,
			"message": "手机号已注册",
		})
		return
	}
	if name == "" {
		user.Name = phone
	} else {
		user.Name = name
	}
	user.Phone = phone
	salt := fmt.Sprintf("%06d", rand.Int31())
	user.Salt = salt
	user.Password = utils.Md5ScalEncode(password, salt)
	user.LoginTime = time.Now()
	user.IsLogin = true
	indetity := models.GetUserUniqueIDs()[0]
	user.Identity = indetity
	// 用户鉴权，token处理
	str := fmt.Sprintf("%d", time.Now().Unix())
	token := utils.Md5ScalEncode(phone, str)
	user.Token = token
	_, err = govalidator.ValidateStruct(user)
	if err != nil {
		c.JSON(http.StatusAccepted, gin.H{
			"state":   false,
			"message": "参数校验失败，数据格式错误",
		})
		return
	}

	user, err = models.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusAccepted, gin.H{
			"state":   false,
			"message": "创建用户失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"state": true,
		"data":  models.OwnData(user),
	})
}

// DeleteUser
// @Summary 删除用户
// @Schemes
// @Description 删除用户
// @Tags 用户
// @param token header string true "用户token"
// @Produce json
// @Success 200 {object} integer 成功
// @Router /user/deleteUser [post]
func DeleteUser(c *gin.Context) {
	token := c.GetHeader("token")
	user, err := models.FindUserByToken(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"state":   false,
			"message": "用户不存在",
		})
		return
	}
	_, err = models.DeleteUser(user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"state":   false,
			"message": "删除用户失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"state":   true,
		"message": "删除用户成功",
	})
}

// UpdateUser
// @Summary 修改用户
// @Schemes
// @Description 修改用户
// @Tags 用户
// @param token header string true "用户token"
// @param name query string false "用户名"
// @param email query string false "邮箱"
// @param phone query string true "手机号"
// @param passWord query string true "密码"
// @Produce json
// @Success 200 {object} integer 成功
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	token := c.GetHeader("token")
	user, err := models.FindUserByToken(token)
	if err != nil {
		c.JSON(http.StatusAccepted, gin.H{
			"state":   false,
			"message": "用户不存在",
		})
		return
	}
	dbUser := user
	user.Name = c.Query("name")
	user.Email = c.Query("email")
	user.Phone = c.Query("phone")
	password := c.Query("passWord")
	salt := fmt.Sprintf("%06d", rand.Int31())
	user.Salt = salt
	user.Password = utils.Md5ScalEncode(password, salt)
	_, err = govalidator.ValidateStruct(user)
	if err != nil {
		c.JSON(http.StatusAccepted, gin.H{
			"state":   false,
			"message": "格式错误",
			"error":   err,
		})
		return
	}
	newUser, err := models.UpdateUser(dbUser, user)
	if err != nil {
		c.JSON(http.StatusAccepted, gin.H{
			"state":   false,
			"message": "更新用户信息失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"state": true,
		"data":  models.OwnData(newUser),
	})
}

// FindUser
// @Summary 查找用户
// @Schemes
// @Description 查找用户
// @Tags 用户
// @param identify query string true "用户ID"
// @Produce json
// @Success 200 {object} integer 成功
// @Router /user/findUser [get]
func FindUser(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("identify"), 10, 64)
	user, error := models.FindUserByIdentity(id)
	if error != nil {
		c.JSON(http.StatusAccepted, gin.H{
			"state":   false,
			"message": "用户不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"state": true,
		"data":  user,
	})
}

// Login
// @Summary 登录
// @Schemes
// @Description 根据手机号和密码登录
// @Tags 用户
// @param phone query string true "手机号"
// @param passWord query string true "密码"
// @Produce json
// @Success 200 {object} integer 成功
// @Router /account/login [post]
func Login(c *gin.Context) {
	phone := c.Query("phone")
	passWord := c.Query("passWord")
	user, error := models.FindUserByPhone(phone)
	if error != nil {
		c.JSON(http.StatusAccepted, gin.H{
			"state":   false,
			"message": "用户不存在，请确认用户信息",
		})
		return
	}
	dbUser := user
	salt := user.Salt
	dbPwd := user.Password
	if !utils.ValidMd5Scal(passWord, salt, dbPwd) {
		c.JSON(http.StatusAccepted, gin.H{
			"state":   false,
			"message": "密码错误，请确认",
		})
		return
	}
	user.LoginTime = time.Now()
	user.IsLogin = true
	// 用户鉴权，token处理
	str := fmt.Sprintf("%d", time.Now().Unix())
	token := utils.Md5ScalEncode(phone, str)
	user.Token = token
	user.LoginTime = time.Now()
	user.IsLogin = true
	user, error = models.LoginUser(dbUser, user)
	if error != nil {
		c.JSON(http.StatusOK, gin.H{
			"state":   false,
			"message": "登录失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"state": true,
		"data":  models.OwnData(user),
	})
}

// Logout
// @Summary 登出
// @Schemes
// @Description 退出登录
// @Tags 用户
// @param token header string true "用户token"
// @Produce json
// @Success 200 {object} integer 成功
// @Router /account/logout [post]
func Logout(c *gin.Context) {
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
	dbUser := user
	user.LogoutTime = time.Now()
	user.IsLogin = false
	user.Token = ""
	user, error = models.LogoutUser(dbUser, user)
	if error != nil {
		c.JSON(http.StatusOK, gin.H{
			"state":   false,
			"message": "退出失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"state":   true,
		"message": "退出成功",
	})
}

// 防止跨域站点伪请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
