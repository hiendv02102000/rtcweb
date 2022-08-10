package entity

type Room struct {
	ID         string `gorm:"column:id;primary_key;not null" json:"id"`
	Title      string `gorm:"column:id;not null" json:"title"`
	IsStream   bool   `gorm:"column:id;not null" json:"is_stream"`
	StreamerID int    `gorm:"column:streamer_id;not null"  json:"streamer_id"`
	Streamer   Users  `gorm:"foreignKey:streamer_id;references:id" json:"streamer"`
	BaseModel
}

func (u *Room) TableName() string {
	return "room"
}
