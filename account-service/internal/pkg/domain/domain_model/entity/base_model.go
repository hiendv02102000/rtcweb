package entity

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel struct
type BaseModel struct {
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}
type BaseModelWithDeleteAt struct {
	BaseModel
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (m *BaseModel) BeforeCreate(db *gorm.DB) error {
	m.CreatedAt = time.Now().UTC()
	return nil
}

func (m *BaseModel) BeforeUpdate(db *gorm.DB) error {
	m.UpdatedAt = time.Now().UTC()
	return nil
}
func (m *BaseModelWithDeleteAt) BeforeDelete(db *gorm.DB) error {
	m.CreatedAt = time.Now().UTC()
	return nil
}
