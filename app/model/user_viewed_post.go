package model

import "time"

type UserViewedPost struct {
	ID       uint `gorm:"primaryKey"`
	UserID   uint
	PostID   uint ``
	ViewedAt time.Time

	User User `gorm:"foreignKey:UserID"`
	Post Post `gorm:"foreignKey:PostID"`
}
