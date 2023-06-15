package model

import "time"

type User struct {
	ID             uint      `gorm:"primaryKey"`
	Username       string    `gorm:"type:varchar(250);not null"`
	FullName       string    `gorm:"type:varchar(250);not null"`
	Email          string    `gorm:"type:varchar(250);not null"`
	Password       string    `gorm:"type:varchar(250);not null"`
	Bio            string    `gorm:"type:varchar(150)"`
	Gender         string    `gorm:"type:varchar(100)"`
	IsPrivate      bool      `json:"is_private" gorm:"default:false"`
	ProfilePicture string    `json:"profile_picture" gorm:"type:varchar(250)"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	ViewedPosts []UserViewedPost `gorm:"foreignKey:UserID"`
	Posts       []Post           `gorm:"foreignKey:UserID"`
	LikedPosts  []Post           `gorm:"many2many:user_liked_posts"`
	Followers   []UserFollower   `gorm:"foreignKey:UserID"`
	Following   []UserFollower   `gorm:"foreignKey:FollowerID"`
	Comments    []Comment        `gorm:"foreignKey:UserID"`
}

type UserProfile struct {
	ID             uint   `json:"id"`
	Username       string `json:"username"`
	FullName       string `json:"full_name"`
	Bio            string `json:"bio"`
	ProfilePicture string `json:"profile_Picture"`
}

type UserConnection struct {
	Followers int `json:"followers"`
	Following int `json:"following"`
}

type PreviewUser struct {
	ID             uint   `json:"id"`
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_Picture"`
}

type UserProfileWithConnection struct {
	ID             uint   `json:"id"`
	Username       string `json:"username"`
	FullName       string `json:"full_name"`
	Bio            string `json:"bio"`
	ProfilePicture string `json:"profile_Picture"`
	Followers      int    `json:"followers"`
	Following      int    `json:"following"`
}
