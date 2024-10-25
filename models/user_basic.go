package models

import (
	"fmt"
	"gimchat/utils"
	"time"

	"gorm.io/gorm"
)

type LoginUserBasic struct {
	Avatar   string `json:"avatar"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Identity int64  `json:"identity"`
	Token    string `json:"token"`
}

// 用户信息
type UserBasic struct {
	gorm.Model    `json:"-"`
	Avatar        string    `json:"avatar"`
	Name          string    `json:"name"`
	Email         string    `valid:"email" json:"email"`
	Phone         string    `valid:"matches(^1[3-9]{1}\\d{9}$)" json:"phone"`
	Password      string    `json:"-"`
	Identity      int64     `json:"identity"`
	Token         string    `json:"-"`
	ClientIP      string    `json:"-"`
	ClientPort    string    `json:"-"`
	Salt          string    `json:"-"`
	LoginTime     time.Time `gorm:"default:null;comment:登录时间" json:"-"`
	LogoutTime    time.Time `gorm:"column:logout_time;default:null;comment:退出登录时间" json:"-"` // 设置别名
	HeartbeatTime time.Time `gorm:"default:null;comment:socket心跳时间" json:"-"`
	IsLogin       bool      `json:"-"`
	DeviceInfo    string    `json:"-"`
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

// 自己的用户信息
func OwnData(data UserBasic) LoginUserBasic {
	return LoginUserBasic{
		Avatar:   data.Avatar,
		Name:     data.Name,
		Email:    data.Email,
		Phone:    data.Phone,
		Identity: data.Identity,
		Token:    data.Token,
	}
}

// 初始化数据表
func InitDB(dst interface{}) {
	if !utils.DB.Migrator().HasTable(dst) {
		utils.DB.AutoMigrate(dst)
	}
}

// 创建用户
func CreateUser(user UserBasic) (UserBasic, error) {
	InitDB(&UserBasic{})
	db := utils.DB.Create(&user)
	return user, db.Error
}

// 删除用户
func DeleteUser(user UserBasic) (UserBasic, error) {
	db := utils.DB.Delete(&user)
	return user, db.Error
}

// 更新数据
func UpdateUser(dbData UserBasic, data UserBasic) (UserBasic, error) {
	// utils.DB.Model(&dbData).Updates(data) // 仅更新非零值字段，bool=false或string=”“时不会更新
	db := utils.DB.Model(&dbData).Updates(
		map[string]interface{}{
			"name":           data.Name,
			"phone":          data.Phone,
			"password":       data.Password,
			"salt":           data.Salt,
			"identity":       data.Identity,
			"token":          data.Token,
			"client_ip":      data.ClientIP,
			"client_port":    data.ClientPort,
			"is_login":       data.IsLogin,
			"device_info":    data.DeviceInfo,
			"login_time":     data.LoginTime,
			"logout_time":    data.LogoutTime,
			"heartbeat_time": data.HeartbeatTime,
		},
	)
	return dbData, db.Error
}

// 登录用户
func LoginUser(dbData UserBasic, data UserBasic) (UserBasic, error) {
	// utils.DB.Model(&dbData).Updates(data) // 仅更新非零值字段，bool=false或string=”“时不会更新
	db := utils.DB.Model(&dbData).Updates(
		map[string]interface{}{
			"token":       data.Token,
			"client_ip":   data.ClientIP,
			"client_port": data.ClientPort,
			"is_login":    data.IsLogin,
			"device_info": data.DeviceInfo,
			"login_time":  data.LoginTime,
		},
	)
	return dbData, db.Error
}

// 登出用户
func LogoutUser(dbData UserBasic, data UserBasic) (UserBasic, error) {
	// utils.DB.Model(&dbData).Updates(data) // 仅更新非零值字段，bool=false或string=”“时不会更新
	db := utils.DB.Model(&dbData).Updates(
		map[string]interface{}{
			"token":       data.Token,
			"client_ip":   data.ClientIP,
			"client_port": data.ClientPort,
			"is_login":    data.IsLogin,
			"device_info": data.DeviceInfo,
			"logout_time": data.LogoutTime,
		},
	)
	return dbData, db.Error
}

// 获取用户列表
func FindUserList(idList []int64) ([]*UserBasic, error) {
	data := make([]*UserBasic, 20)
	db := utils.DB.Select("name,email,phone,identity").Where("identity in (?)", idList).Find(&data)
	for _, v := range data {
		fmt.Println("\n用户数据:", v)
	}
	return data, db.Error
}

// 查找
func FindUserByID(id uint) (UserBasic, error) {
	user := UserBasic{}
	db := utils.DB.First(&user, "id = ?", id)
	return user, db.Error
}

// 通过手机号查找用户
func FindUserByPhone(phone string) (UserBasic, error) {
	user := UserBasic{}
	db := utils.DB.First(&user, "phone = ?", phone)
	return user, db.Error
}

// 通过identity查找用户
func FindUserByIdentity(identity int64) (UserBasic, error) {
	user := UserBasic{}
	db := utils.DB.First(&user, "identity = ?", identity)
	return user, db.Error
}

// 通过token查找用户
func FindUserByToken(token string) (UserBasic, error) {
	user := UserBasic{}
	db := utils.DB.First(&user, "token = ?", token)
	return user, db.Error
}

// 通过手机号密码查找用户
func FindUserByPhoneAndPwd(phone string, pwd string) (UserBasic, error) {
	user := UserBasic{}
	db := utils.DB.First(&user, "phone = ? and password = ?", phone, pwd)
	return user, db.Error
}

// 获取用户id列表
func GetUserUniqueIDs() []int64 {
	// 最小6位id，初始用户id为111111(startId + 1)
	var startId int64 = 111110
	// 步长为1
	step := 1

	dbData := UniqueID{}
	// 读取旧数据
	db := utils.DB.Where(&UniqueID{Key: "userUniqueID"}).First(&dbData)
	if db.Error == nil {
		startId = dbData.MaxId
		step = dbData.Step
	}
	// 创建需要的数据
	uniqueIDs := createUniqueIDs(startId, step)
	// 更新数据
	newData := UniqueID{}
	newData.Key = "userUniqueID"
	newData.MaxId = uniqueIDs[len(uniqueIDs)-1]
	newData.Step = step
	newData.Description = "用户id生成，根据步长生成一定量的用户id，在服务器缓存里用于分配给新创建的用户"
	newData.UpdateTime = time.Now()
	// 写入数据库
	if db.Error != nil {
		CreateUniqueID(newData)
	} else {
		UpdateUniqueID(dbData, newData)
	}
	return uniqueIDs
}
