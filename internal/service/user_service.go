package service

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"mime/multipart"
	"social-app/internal/dto"
	"social-app/internal/entity"
	customError "social-app/pkg/error"
	"social-app/pkg/util"
	"time"
)

type UserService interface {
	GetUserProfile(username string) (dto.UserProfileWithConnection, error)
	FollowUser(followingID uint, followerID uint) error
	UnfollowUser(followingID uint, followerID uint) error
	GetUserFollowers(username string) ([]dto.PreviewUser, error)
	GetUserFollowing(username string) ([]dto.PreviewUser, error)
	UpdateUserProfile(files []*multipart.FileHeader, userId uint, request dto.UpdateUserProfileDto) error
	GetPostsWatchedByUser(username string) ([]dto.PostWithUser, error)
	GetUserUploadedPosts(username string) ([]dto.PostWithUser, error)
	GetUserLikedPosts(username string) ([]dto.PostWithUser, error)
}

type userServiceStruct struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userServiceStruct{db}
}

func (s *userServiceStruct) IsFollowingUser(followingID uint, followerID uint) (bool, error) {
	var count int64

	if err := s.db.Model(&entity.UserFollower{}).Where("user_id = ? AND follower_id = ?", followingID, followerID).Count(&count).Error; err != nil {
		log.Errorf("Database error: %s", err.Error())
		return false, &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	return count > 0, nil
}

func (s *userServiceStruct) GetUserProfile(username string) (dto.UserProfileWithConnection, error) {
	var user entity.User

	if err := s.db.Select("id", "username", "full_name", "bio", "profile_picture").Model(entity.User{}).Preload("Followers").Preload("Following").Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.UserProfileWithConnection{}, &customError.ServiceError{Type: "value_doesnt_exist", Errors: fmt.Sprintf("account with username %s not found", username)}
		}

		log.Errorf("Database error: %s", err.Error())
		return dto.UserProfileWithConnection{}, &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	userFriendship := dto.UserProfileWithConnection{
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
	var user entity.User
	var follower entity.User

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
		follower := entity.UserFollower{
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
	var user entity.User
	var follower entity.User

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
		if err := s.db.Where("user_id = ? AND follower_id = ?", followingID, followerID).Delete(&entity.UserFollower{}).Error; err != nil {
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

func (s *userServiceStruct) GetUserFollowers(username string) ([]dto.PreviewUser, error) {
	var user entity.User
	var UserFollowers []dto.PreviewUser

	if err := s.db.Model(&entity.User{}).Select("id").Where("username = ?", username).First(&user).Error; err != nil {
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
		previewUser := dto.PreviewUser{
			ID:             follower.User.ID,
			Username:       follower.User.Username,
			ProfilePicture: follower.User.ProfilePicture,
		}

		UserFollowers = append(UserFollowers, previewUser)
	}

	return UserFollowers, nil
}

func (s *userServiceStruct) GetUserFollowing(username string) ([]dto.PreviewUser, error) {
	var user entity.User
	var userFollowing []dto.PreviewUser

	if err := s.db.Model(&entity.User{}).Select("id").Where("username = ?", username).First(&user).Error; err != nil {
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
		previewUser := dto.PreviewUser{
			ID:             following.User.ID,
			Username:       following.User.Username,
			ProfilePicture: following.User.ProfilePicture,
		}
		userFollowing = append(userFollowing, previewUser)
	}

	return userFollowing, nil
}

func (s *userServiceStruct) UpdateUserProfile(files []*multipart.FileHeader, userId uint, request dto.UpdateUserProfileDto) error {
	var user entity.User

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

	if len(files) != 0 {
		filePath, err := util.Upload(files[0], "avatar")
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

func (s *userServiceStruct) GetPostsWatchedByUser(username string) ([]dto.PostWithUser, error) {
	var user entity.User
	var posts []entity.Post
	var userPosts []dto.PostWithUser

	if err := s.db.Model(&entity.User{}).Select("id").Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &customError.ServiceError{Type: "value_doesnt_exist", Errors: fmt.Sprintf("account with username %s not found", username)}
		}
		log.Errorf("Database error: %s", err.Error())
		return nil, &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	if err := s.db.Model(&entity.Post{}).Select("posts.id, posts.image, posts.caption, posts.created_at, posts.user_id, posts.likes_count, posts.viewers_count").Preload("User", func(db *gorm.DB) *gorm.DB {
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

func (s *userServiceStruct) GetUserUploadedPosts(username string) ([]dto.PostWithUser, error) {
	var user entity.User
	var posts []entity.Post
	var userPost []dto.PostWithUser

	if err := s.db.Model(&entity.User{}).Select("id").Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &customError.ServiceError{Type: "value_doesnt_exist", Errors: fmt.Sprintf("account with username %s not found", username)}
		}

		log.Errorf("Database error: %s", err.Error())
		return nil, &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	if err := s.db.Model(&entity.Post{}).Select("DISTINCT posts.id, posts.image, posts.caption, posts.created_at, posts.user_id, posts.likes_count, posts.viewers_count").Preload("User", func(db *gorm.DB) *gorm.DB {
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

func (s *userServiceStruct) GetUserLikedPosts(username string) ([]dto.PostWithUser, error) {
	var user entity.User
	var posts []dto.PostWithUser

	if err := s.db.Model(&entity.User{}).Select("id").Where("username = ?", username).First(&user).Error; err != nil {
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
		previewUser := dto.PreviewUser{
			ID:             post.User.ID,
			Username:       post.User.Username,
			ProfilePicture: post.User.ProfilePicture,
		}

		post := dto.PostWithUser{
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
