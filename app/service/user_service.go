package service

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"social-app/app/dto"
	"social-app/app/model"
	customError "social-app/pkg/error"
	"social-app/pkg/util"
	"time"
)

type UserService interface {
	GetUserProfile(username string) (model.UserProfileWithConnection, error)
	FollowUser(followingID uint, followerID uint) error
	UnfollowUser(followingID uint, followerID uint) error
	GetUserFollowers(username string) ([]model.PreviewUser, error)
	GetUserFollowing(username string) ([]model.PreviewUser, error)
	UpdateUserProfile(userId uint, request dto.UpdateUserProfileDto) error
	GetPostsWatchedByUser(username string) ([]model.PostWithUser, error)
	GetUserUploadedPosts(username string) ([]model.PostWithUser, error)
	GetUserLikedPosts(username string) ([]model.PostWithUser, error)
}

type userServiceStruct struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userServiceStruct{db}
}

func (s *userServiceStruct) IsFollowingUser(followingID uint, followerID uint) (bool, error) {
	var count int64

	if err := s.db.Model(&model.UserFollower{}).Where("user_id = ? AND follower_id = ?", followingID, followerID).Count(&count).Error; err != nil {
		log.Errorf("Database error: %s", err.Error())
		return false, &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	return count > 0, nil
}

func (s *userServiceStruct) GetUserProfile(username string) (model.UserProfileWithConnection, error) {
	var user model.User

	if err := s.db.Select("id", "username", "full_name", "bio", "profile_picture").Model(model.User{}).Preload("Followers").Preload("Following").Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.UserProfileWithConnection{}, &customError.ServiceError{Type: "value_doesnt_exist", Errors: fmt.Sprintf("account with username %s not found", username)}
		}

		log.Errorf("Database error: %s", err.Error())
		return model.UserProfileWithConnection{}, &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	userFriendship := model.UserProfileWithConnection{
		ID:             user.ID,
		Username:       user.Username,
		FullName:       user.FullName,
		Bio:            user.Bio,
		ProfilePicture: user.ProfilePicture,
		Followers:      len(user.Followers),
		Following:      len(user.Following),
	}

	return userFriendship, nil
}

func (s *userServiceStruct) FollowUser(followingID uint, followerID uint) error {
	var user model.User
	var follower model.User

	tx := s.db.Begin()

	if err := tx.First(&user, followingID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &customError.ServiceError{Type: "value_doesnt_exist", Errors: fmt.Sprintf("account with id %d not found", followingID)}
		}

		tx.Rollback()
		log.Errorf("Database error: %s", err.Error())
		return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	if err := tx.First(&follower, followerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &customError.ServiceError{Type: "value_doesnt_exist", Errors: fmt.Sprintf("account with id %d not found", followerID)}
		}

		tx.Rollback()
		log.Errorf("Database error: %s", err.Error())
		return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	isFollowing, err := s.IsFollowingUser(followingID, followerID)
	if err != nil {
		log.Errorf("Database error: %s", err.Error())
		return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	if !isFollowing {
		follower := model.UserFollower{
			UserID:     followingID,
			FollowerID: followerID,
			FollowedAt: time.Now(),
		}

		if err := tx.Create(&follower).Error; err != nil {
			tx.Rollback()
			log.Errorf("Database error: %s", err.Error())
			return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
		}

		if err := tx.Commit().Error; err != nil {
			log.Errorf("Database error: %s", err.Error())
			return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
		}

		return nil
	}

	tx.Rollback()
	return &customError.ServiceError{Type: "user_follower_error", Errors: "You have done following"}
}

func (s *userServiceStruct) UnfollowUser(followingID uint, followerID uint) error {
	var user model.User
	var follower model.User

	tx := s.db.Begin()

	if err := s.db.Preload("Following").First(&user, followingID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &customError.ServiceError{Type: "value_doesnt_exist", Errors: fmt.Sprintf("account with id %d not found", followingID)}
		}

		tx.Rollback()
		log.Errorf("Database error: %s", err.Error())
		return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	if err := s.db.First(&follower, followerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &customError.ServiceError{Type: "value_doesnt_exist", Errors: fmt.Sprintf("account with id %d not found", followingID)}
		}

		tx.Rollback()
		log.Errorf("Database error: %s", err.Error())
		return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	isFollowing, err := s.IsFollowingUser(followingID, followerID)
	if err != nil {
		log.Errorf("Database error: %s", err.Error())
		return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	if isFollowing {
		if err := s.db.Where("user_id = ? AND follower_id = ?", followingID, followerID).Delete(&model.UserFollower{}).Error; err != nil {
			tx.Rollback()
			log.Errorf("Database error: %s", err.Error())
			return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
		}

		if err := tx.Commit().Error; err != nil {
			log.Errorf("Database error: %s", err.Error())
			return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
		}

		return nil
	}

	tx.Rollback()
	return &customError.ServiceError{Type: "user_follower_error", Errors: "You haven't followed yet"}
}

func (s *userServiceStruct) GetUserFollowers(username string) ([]model.PreviewUser, error) {
	var user model.User
	var UserFollowers []model.PreviewUser

	if err := s.db.Model(&model.User{}).Select("id").Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &customError.ServiceError{Type: "value_doesnt_exist", Errors: fmt.Sprintf("account with username %s not found", username)}
		}

		log.Errorf("Database error: %s", err.Error())
		return nil, &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	if err := s.db.Select("id, username, profile_picture").Preload("Followers.User").First(&user, user.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &customError.ServiceError{Type: "value_doesnt_exist", Errors: fmt.Sprintf("account with username %s not found", username)}
		}

		log.Errorf("Database error: %s", err.Error())
		return nil, &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	for _, follower := range user.Followers {
		previewUser := model.PreviewUser{
			ID:             follower.User.ID,
			Username:       follower.User.Username,
			ProfilePicture: follower.User.ProfilePicture,
		}

		UserFollowers = append(UserFollowers, previewUser)
	}

	return UserFollowers, nil
}

func (s *userServiceStruct) GetUserFollowing(username string) ([]model.PreviewUser, error) {
	var user model.User
	var userFollowing []model.PreviewUser

	if err := s.db.Model(&model.User{}).Select("id").Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &customError.ServiceError{Type: "value_doesnt_exist", Errors: fmt.Sprintf("account with username %s not found", username)}
		}

		log.Errorf("Database error: %s", err.Error())
		return nil, &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	if err := s.db.Select("id, username, profile_picture").Preload("Following.User").First(&user, user.ID).Error; err != nil {
		log.Errorf("Database error: %s", err.Error())
		return nil, &customError.ServiceError{Type: "user_follower_error", Errors: err.Error()}
	}

	for _, following := range user.Following {
		previewUser := model.PreviewUser{
			ID:             following.User.ID,
			Username:       following.User.Username,
			ProfilePicture: following.User.ProfilePicture,
		}
		userFollowing = append(userFollowing, previewUser)
	}

	return userFollowing, nil
}

func (s *userServiceStruct) UpdateUserProfile(userId uint, request dto.UpdateUserProfileDto) error {
	var user model.User

	tx := s.db.Begin()

	if err := tx.First(&user, userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return &customError.ServiceError{Type: "value_doesnt_exist", Errors: "account could not be found"}
		}

		log.Errorf("Database error: %s", err.Error())
		tx.Rollback()
		return &customError.ServiceError{Type: "database_error", Errors: err}
	}

	if request.ProfilePicture != nil {
		filePath, err := util.Upload(request.ProfilePicture, "avatar")
		if err != nil {
			log.Errorf("Upload file error: %s", err.Error())
			return &customError.ServiceError{Type: "upload_file_error", Errors: err.Error()}
		}

		user.ProfilePicture = filePath
	}
	if request.Bio != "" || request.Bio != user.Bio {
		user.Bio = request.Bio
	}
	if request.Gender != "" && request.Gender != user.Gender {
		if request.Gender != "Male" || request.Gender != "Female" || request.Gender != "Unknown" {
			user.Gender = request.Gender
		} else {
			return &customError.ServiceError{Type: "invalid_value_error", Errors: "invalid gender"}
		}
	}

	if err := tx.Updates(&user).Error; err != nil {
		log.Errorf("Database error: %s", err.Error())
		tx.Rollback()
		return &customError.ServiceError{Type: "database_error", Errors: err}
	}

	if err := tx.Commit().Error; err != nil {
		log.Errorf("Database error: %s", err.Error())
		return &customError.ServiceError{Type: "database_error", Errors: err}
	}

	return nil
}

func (s *userServiceStruct) GetPostsWatchedByUser(username string) ([]model.PostWithUser, error) {
	var user model.User
	var posts []model.Post
	var userPosts []model.PostWithUser

	if err := s.db.Model(&model.User{}).Select("id").Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &customError.ServiceError{Type: "value_doesnt_exist", Errors: fmt.Sprintf("account with username %s not found", username)}
		}
		log.Errorf("Database error: %s", err.Error())
		return nil, &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	if err := s.db.Model(&model.Post{}).Select("posts.id, posts.image, posts.caption, posts.created_at, posts.user_id, posts.likes_count, posts.viewers_count").Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username, profile_picture")
	}).Joins("JOIN user_viewed_posts ON posts.id = user_viewed_posts.post_id").Find(&posts).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		log.Errorf("Database error: %s", err.Error())
		return nil, err
	}

	userPosts = util.ConvertToUserPosts(posts)

	return userPosts, nil
}

func (s *userServiceStruct) GetUserUploadedPosts(username string) ([]model.PostWithUser, error) {
	var user model.User
	var posts []model.Post
	var userPost []model.PostWithUser

	if err := s.db.Model(&model.User{}).Select("id").Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &customError.ServiceError{Type: "value_doesnt_exist", Errors: fmt.Sprintf("account with username %s not found", username)}
		}

		log.Errorf("Database error: %s", err.Error())
		return nil, &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	if err := s.db.Model(&model.Post{}).Select("DISTINCT posts.id, posts.image, posts.caption, posts.created_at, posts.user_id, posts.likes_count, posts.viewers_count").Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username, profile_picture")
	}).Where("posts.user_id = ?", user.ID).Find(&posts).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		log.Errorf("Database error: %s", err.Error())
		return nil, &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	userPost = util.ConvertToUserPosts(posts)

	return userPost, nil
}

func (s *userServiceStruct) GetUserLikedPosts(username string) ([]model.PostWithUser, error) {
	var user model.User
	var posts []model.PostWithUser

	if err := s.db.Model(&model.User{}).Select("id").Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &customError.ServiceError{Type: "value_doesnt_exist", Errors: fmt.Sprintf("account with username %s not found", username)}
		}

		log.Errorf("Database error: %s", err.Error())
		return nil, &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	if err := s.db.Preload("LikedPosts.User").First(&user, user.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	for _, post := range user.LikedPosts {
		previewUser := model.PreviewUser{
			ID:             post.User.ID,
			Username:       post.User.Username,
			ProfilePicture: post.User.ProfilePicture,
		}

		post := model.PostWithUser{
			ID:           post.ID,
			Image:        post.Image,
			Caption:      post.Caption,
			User:         previewUser,
			ViewersCount: post.ViewersCount,
			LikesCount:   post.LikesCount,
			CreatedAt:    post.CreatedAt,
		}

		posts = append(posts, post)
	}

	return posts, nil
}
