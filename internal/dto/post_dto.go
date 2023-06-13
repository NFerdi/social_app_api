package dto

import (
	"mime/multipart"
	"time"
)

type CreatePostDto struct {
	Image   *multipart.FileHeader `form:"image"`
	Caption string                `json:"caption" validate:"required"`
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
