package entity

import "time"

type User struct {
	ID             uint   `gorm:"primaryKey"`
	Username       string `gorm:"type:varchar(250);not null"`
	FullName       string `gorm:"type:varchar(250);not null"`
	Email          string `gorm:"type:varchar(250);not null"`
	Password       string `gorm:"type:varchar(250);not null"`
	Bio            string `gorm:"type:varchar(150)"`
	Gender         string `gorm:"type:varchar(100)"`
	IsPrivate      bool   `json:"is_private" gorm:"default:false"`
	ProfilePicture string `json:"profile_picture" gorm:"type:varchar(250)"`

	ViewedPosts []UserViewedPost `gorm:"foreignKey:UserID"`
	Posts       []Post           `gorm:"foreignKey:UserID"`

	LikedPosts []Post `gorm:"many2many:user_liked_posts"`

	Followers []UserFollower `gorm:"foreignKey:UserID"`
	Following []UserFollower `gorm:"foreignKey:FollowerID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
