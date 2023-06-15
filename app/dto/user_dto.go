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
