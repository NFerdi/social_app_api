package model

import "time"

type UserLikedPost struct {
	ID      uint `gorm:"primaryKey"`
	UserID  uint
	PostID  uint
	LikedAt time.Time

	User User `gorm:"foreignKey:UserID"`
	Post Post `gorm:"foreignKey:PostID"`
}
