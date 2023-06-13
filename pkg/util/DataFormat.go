package util

import (
	dto2 "social-app/internal/dto"
	"social-app/internal/entity"
)

func ConvertToUserPosts(posts []entity.Post) []dto2.PostWithUser {
	var displayPosts []dto2.PostWithUser

	for _, post := range posts {
		previewUser := dto2.PreviewUser{
			ID:             post.User.ID,
			Username:       post.User.Username,
			ProfilePicture: post.User.ProfilePicture,
		}

		var newDisplayPost = dto2.PostWithUser{
			ID:           post.ID,
			Image:        post.Image,
			Caption:      post.Caption,
			User:         previewUser,
			ViewersCount: post.ViewersCount,
			LikesCount:   post.LikesCount,
			CreatedAt:    post.CreatedAt,
		}
		displayPosts = append(displayPosts, newDisplayPost)
	}

	return displayPosts
}
