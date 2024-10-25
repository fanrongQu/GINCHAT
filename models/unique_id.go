package models

import (
	"gimchat/utils"
	"time"

	"gorm.io/gorm"
)

func (table *UniqueID) TableName() string {
	return "unique_id"
}

// 唯一id标识库
type UniqueID struct {
	gorm.Model
	Key         string // 对应业务key值
	MaxId       int64  // 当前最大值
	Step        int    // 步长
	Description string // 业务描述
	UpdateTime  time.Time
}

// 创建标识数据
func CreateUniqueID(data UniqueID) (UniqueID, error) {
	InitDB(&UniqueID{})
	db := utils.DB.Create(&data)
	return data, db.Error
}

// 更新数据
func GetUniqueID(dbData UniqueID, data UniqueID) (UniqueID, error) {
	db := utils.DB.Model(&dbData).Updates(data)
	return dbData, db.Error
}

// 更新数据
func UpdateUniqueID(dbData UniqueID, data UniqueID) (UniqueID, error) {
	db := utils.DB.Model(&dbData).Updates(data)
	return dbData, db.Error
}

// 根据起始值和步长创建id列表
func createUniqueIDs(startId int64, step int) []int64 {
	if step < 1 {
		return []int64{}
	}
	idList := []int64{}
	for i := 0; i < step; i++ {
		newID := startId + int64(i) + 1
		idList = append(idList, newID)
	}
	return idList
}
