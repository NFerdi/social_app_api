package entity

import "time"

type UserLikedPost struct {
	ID      uint `gorm:"primaryKey"`
	UserID  uint
	PostID  uint
	User    User `gorm:"foreignKey:UserID"`
	Post    Post `gorm:"foreignKey:PostID"`
	LikedAt time.Time
}
