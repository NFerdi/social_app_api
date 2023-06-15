package dto

import (
	"mime/multipart"
)

type CreatePostDto struct {
	Image   *multipart.FileHeader `form:"image"`
	Caption string                `json:"caption" validate:"required"`
}

type CreateCommentPostDto struct {
	Comment string `json:"comment"`
}
