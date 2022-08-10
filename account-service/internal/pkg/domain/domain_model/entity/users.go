package entity

import (
	"time"
)

type userRole string

const (
	AdminRole    userRole = "admin"
	CustomerRole userRole = "customer"
)

// Users struct
type Users struct {
	ID             int        `gorm:"column:id;primary_key;auto_increment;not null" json:"id"`
	Email          string     `gorm:"column:email;not null;unique;type:varchar(255)" json:"email"`
	Password       string     `gorm:"column:password;not null;type:varchar(255)" json:"password"`
	FirstName      string     `gorm:"column:first_name;type:varchar(255)" json:"first_name"`
	LastName       string     `gorm:"column:last_name;type:varchar(255)" json:"last_name"`
	Role           userRole   `gorm:"column:role;type:varchar(255)" json:"role"`
	AvatarUrl      *string    `gorm:"column:avatar_url;type:varchar(255)" json:"avatar_url"`
	Token          *string    `gorm:"column:token;type:varchar(255)" json:"token"`
	TokenExpiredAt *time.Time `gorm:"column:token_expired_at" json:"token_expired_at"`
	BaseModel
}

func (u *Users) TableName() string {
	return "users"
}
