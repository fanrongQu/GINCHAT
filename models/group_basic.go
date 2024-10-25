package models

import (
	"time"

	"gorm.io/gorm"
)

type GroupBasic struct {
	gorm.Model
	Name      string
	Icon      int
	Type      int
	Desc      string
	Identity  string
	OwnerID   uint
	Owner     UserBasic   `gorm:"foreignKey:Identity;references:OwnerUserID"`
	Admins    []UserBasic `gorm:"many2many:group_admins;foreignKey:Identity;references:AdminUserID"`
	Members   []UserBasic `gorm:"many2many:group_members;foreignKey:Identity;references:MemberUserID"`
	CreatTime time.Time
}

func (table *GroupBasic) TableName() string {
	return "group_basic"
}
