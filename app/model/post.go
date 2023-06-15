package model

import "time"

type Post struct {
	ID           uint      `gorm:"primaryKey"`
	Image        string    `gorm:"type:varchar(250);not null"`
	Caption      string    `gorm:"type:varchar(250);not null"`
	ViewersCount int       `json:"viewers_count"`
	LikesCount   int       `json:"likes_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	UserID   uint
	User     User
	ViewedBy []UserViewedPost `gorm:"foreignKey:PostID" json:"viewed_by"`
	Likes    []UserLikedPost  `gorm:"many2many:user_liked_posts"`
	Comments []Comment        `gorm:"foreignKey:PostID"`
}

type CountViewerPost struct {
	PostId uint `json:"post_id"`
	Viewer int64
}

type PostWithUser struct {
	ID           uint        `json:"id"`
	Image        string      `json:"image"`
	Caption      string      `json:"caption"`
	User         PreviewUser `json:"user"`
	ViewersCount int         `json:"viewers_count"`
	LikesCount   int         `json:"likes_count"`
	CreatedAt    time.Time   `json:"created_at"`
}
