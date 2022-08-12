package dto

import "time"

type MessageChat struct {
	Type    string     `json:"type" form:"type" binding:"required"`
	Content string     `json:"content" form:"content" binding:"required"`
	Sender  string     `json:"sender" form:"sender" binding:"omitempty"`
	SendAt  *time.Time `json:"send_at" binding:"omitempty"`
}
