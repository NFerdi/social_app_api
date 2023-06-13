package dto

import (
	"mime/multipart"
)

//request

type FollowUserDto struct {
	FollowerId uint `json:"follower_id"`
}

type UpdateUserProfileDto struct {
	ProfilePicture *multipart.FileHeader `json:"profile_picture" validator:"required_without:Bio,Gender"`
	Bio            string                `validator:"required_without:ProfilePicture,Gender"`
	Gender         string                `validator:"required_without:ProfilePicture,Bio"`
}

//response

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
