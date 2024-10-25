package models

import (
	"gimchat/utils"

	"gorm.io/gorm"
)

// 用户关系表
type ContactBasic struct {
	gorm.Model   `json:"-"`
	OwnerUserId  int64     `json:"ownerUserId"`                                                    // 谁的关系用户id
	TargetUserId int64     `json:"targetUserId"`                                                   // 对应的用户数据id
	TargetUser   UserBasic `json:"targetUser" gorm:"foreignKey:Identity;references:TargetUserId;"` // 对应的用户数据
	Type         int       `json:"type"`                                                           // 关系类型（0:陌生人, 1:好友, 2:好友申请已发送, 3:拉黑用户）
	Remarks      string    `json:"remarks"`                                                        // 备注信息
}

func (table *ContactBasic) TableName() string {
	return "contact_basic"
}

// 创建会话关系
func CreateContact(contact ContactBasic) (ContactBasic, error) {
	InitDB(&ContactBasic{})
	db := utils.DB.Create(&contact)
	return contact, db.Error
}

// 删除会话关系
func DeleteContact(contact ContactBasic) (ContactBasic, error) {
	db := utils.DB.Delete(&contact)
	return contact, db.Error
}

// 查找会话关系
func FindContact(ownerUserId int64, targetUserId int64) (ContactBasic, error) {
	contact := ContactBasic{}
	db := utils.DB.Preload("TargetUser").First(&contact, "owner_user_id = ? and target_user_id = ?", ownerUserId, targetUserId)
	return contact, db.Error
}

// 查找所有会话关系
func FindAllContact(ownerUserId int64) ([]*ContactBasic, error) {
	var contacts []*ContactBasic
	db := utils.DB.Where("owner_user_id = ?", ownerUserId).Preload("TargetUser").Find(&contacts)
	return contacts, db.Error
}

// 更新数据
func UpdateContact(dbData ContactBasic, data ContactBasic) (ContactBasic, error) {
	db := utils.DB.Model(&dbData).Updates(
		map[string]interface{}{
			"type":    data.Type,
			"remarks": data.Remarks,
		},
	)
	return data, db.Error
}
