package router

import (
	"gimchat/docs"
	"gimchat/service"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.GET("/index", service.GetIndex)
	r.GET("/user/findUserList", service.FindUserList)
	r.POST("/user/registed", service.CreateUser)
	r.POST("/user/deleteUser", service.DeleteUser)
	r.POST("/user/updateUser", service.UpdateUser)
	r.GET("/user/findUser", service.FindUser)
	r.POST("/account/login", service.Login)
	r.POST("/account/logout", service.Logout)

	r.POST("/contact/setContact", service.SetContact)
	r.GET("/contact/findContact", service.FindContact)
	r.GET("/contact/findAllContact", service.FindAllContact)

	// socket通信，发送消息
	r.GET("/chat", service.Chat)

	return r
}
