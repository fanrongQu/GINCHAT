package models

import (
	"gimchat/utils"
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Identity   string
	FromUserID int64
	// FromUser      UserBasic `gorm:"foreignKey:ID;references:FromUserID"`
	TargetUserID int64
	// TargetUser    UserBasic `gorm:"foreignKey:ID;references:TargetUserID"`
	TargetGroupID int64
	// TargetGroup   Group `gorm:"foreignKey:ID;references:TargetGroupID"`
	Type     int
	Media    string
	Content  string
	Picture  string
	Url      string
	Desc     string
	SendTime time.Time
	Amount   int // 其他数字统计
}

func (table *Message) TableName() string {
	return "message"
}

// 创建消息
func CreateMessage(msg Message) {
	InitDB(&Message{})
	utils.DB.Create(&msg)
}
