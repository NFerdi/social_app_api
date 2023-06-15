package model

import "time"

type UserFollower struct {
	ID         uint `gorm:"primaryKey"`
	UserID     uint
	FollowerID uint
	FollowedAt time.Time

	User     User `gorm:"foreignKey:UserID"`
	Follower User `gorm:"foreignKey:FollowerID"`
}
