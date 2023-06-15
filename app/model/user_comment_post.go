package model

import "time"

type Comment struct {
	ID        uint `gorm:"primary_key"`
	PostID    uint
	UserID    uint
	Content   string
	CreatedAt time.Time

	User User `gorm:"foreignKey:UserID"`
	Post Post `gorm:"foreignKey:PostID"`
}

type PreviewComment struct {
	ID        uint        `json:"id"`
	Content   string      `json:"content"`
	User      PreviewUser `json:"user"`
	CreatedAt time.Time   `json:"createdAt"`
}
