package entity

import "time"

type Post struct {
	ID      uint   `gorm:"primaryKey"`
	Image   string `gorm:"type:varchar(250);not null"`
	Caption string `gorm:"type:varchar(250);not null"`

	ViewersCount int `json:"viewers_count"`
	LikesCount   int `json:"likes_count"`

	UserID uint
	User   User

	ViewedBy []UserViewedPost `gorm:"foreignKey:PostID" json:"viewed_by"`

	Likes []UserLikedPost `gorm:"many2many:user_liked_posts"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
