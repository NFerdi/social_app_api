package util

import (
	"social-app/app/model"
)

func ConvertToUserPosts(posts []model.Post) []model.PostWithUser {
	var displayPosts []model.PostWithUser

	for _, post := range posts {
		previewUser := model.PreviewUser{
			ID:             post.User.ID,
			Username:       post.User.Username,
			ProfilePicture: post.User.ProfilePicture,
		}

		var newDisplayPost = model.PostWithUser{
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
