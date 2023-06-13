package entity

import "time"

type UserViewedPost struct {
	ID       uint `gorm:"primaryKey"`
	UserID   uint
	PostID   uint
	User     User `gorm:"foreignKey:UserID"`
	Post     Post `gorm:"foreignKey:PostID"`
	ViewedAt time.Time
}
