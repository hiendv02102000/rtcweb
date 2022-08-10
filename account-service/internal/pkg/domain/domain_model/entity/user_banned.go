package entity

type UserBanned struct {
	ID         int   `gorm:"column:id;primary_key;auto_increment;not null" json:"id"`
	UserID     int   `gorm:"column:user_id;not null"`
	User       Users `gorm:"foreignKey:user_id;references:id" json:"user"`
	StreamerID int   `gorm:"column:streamer_id;not null"`
	Streamer   Users `gorm:"foreignKey:streamer_id;references:id" json:"streamer"`
	BaseModelWithDeleteAt
}

func (u *UserBanned) TableName() string {
	return "user_banned"
}
