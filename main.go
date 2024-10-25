package main

import (
	"gimchat/router"
	"gimchat/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMySql()
	utils.InitRedis()
	r := router.Router()
	r.Run(":8081") // 默认端口号8080，端口号配置修改
}
