package test

import (
	"fmt"
	"gimchat/models"
	"gimchat/utils"
)

func TestCreat() {
	db := utils.DB
	// 迁移 schema
	db.AutoMigrate(models.UserBasic{})

	// 创建用户
	user := &models.UserBasic{Name: "haha", Password: "123456aaa"}
	db.Create(user)

	// Read
	var product models.UserBasic
	db.First(&product, 1) // 根据整型主键查找
	fmt.Println(product)
	// db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录

	// Update - 将 product 的 price 更新为 200
	db.Model(user).Update("Name", "测试用户")
	// Update - 更新多个字段
	// db.Model(user).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	// db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - 删除 product
	// db.Delete(&product, 1)
}
