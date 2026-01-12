package model

import (
	"time"
)

// User 用户实体
type User struct {
	ID           int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username     string     `gorm:"column:username;type:varchar(256)" json:"username"`
	UserAccount  string     `gorm:"column:userAccount;type:varchar(256);index" json:"userAccount"`
	AvatarUrl    string     `gorm:"column:avatarUrl;type:varchar(1024)" json:"avatarUrl"`
	Gender       int        `gorm:"column:gender;type:tinyint" json:"gender"`
	UserPassword string     `gorm:"column:userPassword;type:varchar(512);not null" json:"-"`
	Phone        string     `gorm:"column:phone;type:varchar(128)" json:"phone"`
	Email        string     `gorm:"column:email;type:varchar(512)" json:"email"`
	UserStatus   int        `gorm:"column:userStatus;default:0;not null" json:"userStatus"`
	CreateTime   time.Time  `gorm:"column:createTime;autoCreateTime" json:"createTime"`
	UpdateTime   time.Time  `gorm:"column:updateTime;autoUpdateTime" json:"updateTime"`
	IsDelete     int        `gorm:"column:isDelete;type:tinyint;default:0;not null" json:"-"`
	UserRole     int        `gorm:"column:userRole;default:0;not null" json:"userRole"`
	PlanetCode   string     `gorm:"column:planetCode;type:varchar(512);index" json:"planetCode"`
}

// TableName 指定表名
func (User) TableName() string {
	return "user"
}

// SafetyUser 用户脱敏后的信息
type SafetyUser struct {
	ID          int64     `json:"id"`
	Username    string    `json:"username"`
	UserAccount string    `json:"userAccount"`
	AvatarUrl   string    `json:"avatarUrl"`
	Gender      int       `json:"gender"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	UserStatus  int       `json:"userStatus"`
	CreateTime  time.Time `json:"createTime"`
	UserRole    int       `json:"userRole"`
	PlanetCode  string    `json:"planetCode"`
}
