package entity

import "time"

type UserFollower struct {
	ID         uint `gorm:"primaryKey"`
	UserID     uint
	FollowerID uint
	User       User `gorm:"foreignKey:UserID"`
	Follower   User `gorm:"foreignKey:FollowerID"`
	FollowedAt time.Time
}
