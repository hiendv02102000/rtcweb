package dto

import "time"

type MessageChat struct {
	Type    string     `json:"type" form:"type"`
	Content string     `json:"content" form:"content"`
	Sender  string     `json:"sender" form:"sender"`
	SendAt  *time.Time `json:"send_at" `
}
