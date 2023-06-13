package service

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"social-app/internal/dto"
	"social-app/internal/entity"
	customError "social-app/pkg/error"
	"social-app/pkg/util"
	"time"
)

type PostService interface {
	CreatePost(ctx *fiber.Ctx, userID uint, request dto.CreatePostDto) error
	ViewPost(userID uint, postID uint) error
	GetUserFeed(userID uint) ([]dto.PostWithUser, error)
	GetViewersOnPost(postID uint) ([]dto.PreviewUser, error)
	LikePost(UserID uint, PostID uint) error
	UnLikePost(UserID uint, PostID uint) error
	GetUserWhoLikedPost(postID uint) ([]dto.PreviewUser, error)
}

type PostServiceStruct struct {
	db *gorm.DB
}

func NewPostService(db *gorm.DB) PostService {
	return &PostServiceStruct{db: db}
}

func (s *PostServiceStruct) CreatePost(ctx *fiber.Ctx, userID uint, request dto.CreatePostDto) error {
	_, filePath, err := util.UploadImage(ctx, "post", "image")

	if err != nil {
		log.Errorf("Upload file error: %s", err.Error())
		return &customError.ServiceError{Type: "upload_file_error", Errors: err.Error()}
	}

	post := entity.Post{
		Image:   filePath,
		Caption: request.Caption,
		UserID:  userID,
	}

	if err := s.db.Create(&post).Error; err != nil {
		log.Errorf("Database error: %s", err.Error())
		return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	return nil
}

func (s *PostServiceStruct) ViewPost(userID uint, postID uint) error {
	var user entity.User

	tx := s.db.Begin()
	if tx.Error != nil {
		log.Errorf("Database error: %s", tx.Error.Error())
		return &customError.ServiceError{Type: "database_error", Errors: tx.Error.Error()}
	}

	if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return &customError.ServiceError{Type: "value_doesnt_exist", Errors: "user not found"}
		}

		tx.Rollback()
		log.Errorf("Database error: %s", err.Error())
		return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	var post entity.Post
	if err := tx.Where("id = ?", postID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return &customError.ServiceError{Type: "value_doesnt_exist", Errors: "post not found"}
		}

		tx.Rollback()
		log.Errorf("Database error: %s", err.Error())
		return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	var countUserViewed int64
	if err := tx.Model(&entity.UserViewedPost{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&countUserViewed).Error; err != nil {
		tx.Rollback()
		log.Errorf("Database error: %s", err.Error())
		return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	if countUserViewed == 0 {
		addUserViewedPost := entity.UserViewedPost{
			UserID:   userID,
			PostID:   postID,
			ViewedAt: time.Now(),
		}

		if err := tx.Create(&addUserViewedPost).Error; err != nil {
			tx.Rollback()
			log.Errorf("Database error: %s", err.Error())
			return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
		}

		post.ViewersCount += 1

		if err := tx.Save(&post).Error; err != nil {
			tx.Rollback()
			log.Errorf("Database error: %s", err.Error())
			return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
		}

		return nil
	}

	if err := tx.Commit().Error; err != nil {
		log.Errorf("Database error: %s", err.Error())
		return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	return nil
}

func (s *PostServiceStruct) GetUserFeed(userID uint) ([]dto.PostWithUser, error) {
	var userViewedPost []entity.UserViewedPost
	var userPost []entity.Post
	var userFeedPost []dto.PostWithUser

	if err := s.db.Where("user_id = ?", userID).
		Find(&userViewedPost).
		Error; err != nil {
		return nil, err
	}

	if len(userViewedPost) > 0 {
		var postIDs []uint

		for _, userViewed := range userViewedPost {
			postIDs = append(postIDs, userViewed.PostID)
		}

		if err := s.db.Select("id, image, caption, created_at, user_id, likes_count, viewers_count").
			Preload("User", func(db *gorm.DB) *gorm.DB {
				return db.Select("id, username, profile_picture")
			}).
			Where("id IN (?)", postIDs).
			Find(&userPost).Error; err != nil {
			log.Errorf("Database error: %s", err.Error())
			return nil, &customError.ServiceError{Type: "database_error", Errors: err.Error()}
		}

		userFeedPost = util.ConvertToUserPosts(userPost)

		return userFeedPost, nil
	}

	if err := s.db.Select("id, image, caption, created_at, user_id, likes_count, viewers_count").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username, profile_picture")
		}).
		Order("viewers_count DESC").
		Limit(10).
		Find(&userPost).Error; err != nil {
		log.Errorf("Database error: %s", err.Error())
		return nil, &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	userFeedPost = util.ConvertToUserPosts(userPost)

	return userFeedPost, nil
}

func (s *PostServiceStruct) GetViewersOnPost(postID uint) ([]dto.PreviewUser, error) {
	var post entity.Post
	var usersPreview []dto.PreviewUser

	if err := s.db.Model(&entity.Post{}).
		Preload("ViewedBy.User").
		Find(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		log.Errorf("Database error: %s", err.Error())
		return nil, &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	if post.ID != 0 {
		for _, viewedBy := range post.ViewedBy {
			userPreview := dto.PreviewUser{
				ID:             viewedBy.User.ID,
				Username:       viewedBy.User.Username,
				ProfilePicture: viewedBy.User.ProfilePicture,
			}

			usersPreview = append(usersPreview, userPreview)
		}
	}

	return usersPreview, nil
}

func (s *PostServiceStruct) LikePost(UserID uint, PostID uint) error {
	var user entity.User
	var post entity.Post

	tx := s.db.Begin()
	if tx.Error != nil {
		log.Errorf("Database error: %s", tx.Error.Error())
		return &customError.ServiceError{Type: "database_error", Errors: tx.Error.Error()}
	}

	if err := tx.First(&user, UserID).Error; err != nil {
		tx.Rollback()
		log.Errorf("Database error: %s", err.Error())
		return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	if err := tx.First(&post, PostID).Error; err != nil {
		tx.Rollback()
		log.Errorf("Database error: %s", err.Error())
		return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	createLikePost := entity.UserLikedPost{
		UserID:  UserID,
		PostID:  PostID,
		LikedAt: time.Now(),
	}

	var countLikedPost int64
	if err := tx.Model(&entity.UserLikedPost{}).
		Where("user_id = ? AND post_id = ?", UserID, PostID).
		Count(&countLikedPost).
		Error; err != nil {
		tx.Rollback()
		log.Errorf("Database error: %s", err.Error())
		return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	if countLikedPost == 0 {
		if err := tx.Create(&createLikePost).Error; err != nil {
			tx.Rollback()
			log.Errorf("Database error: %s", err.Error())
			return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
		}

		post.LikesCount += 1

		if err := tx.Save(&post).Error; err != nil {
			tx.Rollback()
			log.Errorf("Database error: %s", err.Error())
			return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			log.Errorf("Database error: %s", err.Error())
			return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
		}

		return nil
	}

	return nil
}

func (s *PostServiceStruct) UnLikePost(UserID uint, PostID uint) error {
	var user entity.User
	var post entity.Post

	tx := s.db.Begin()
	if tx.Error != nil {
		log.Errorf("Database error: %s", tx.Error.Error())
		return &customError.ServiceError{Type: "database_error", Errors: tx.Error.Error()}
	}

	if err := tx.First(&user, UserID).Error; err != nil {
		tx.Rollback()
		log.Errorf("Database error: %s", err.Error())
		return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	if err := tx.First(&post, PostID).Error; err != nil {
		tx.Rollback()
		log.Errorf("Database error: %s", err.Error())
		return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	var countLikedPost int64
	if err := tx.Model(&entity.UserLikedPost{}).
		Where("user_id = ? AND post_id = ?", UserID, PostID).
		Count(&countLikedPost).Error; err != nil {
		tx.Rollback()
		log.Errorf("Database error: %s", err.Error())
		return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	if countLikedPost != 0 {
		if err := tx.Where("user_id = ? AND post_id = ?", UserID, PostID).Delete(&entity.UserLikedPost{}).Error; err != nil {
			tx.Rollback()
			log.Errorf("Database error: %s", err.Error())
			return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
		}

		post.LikesCount -= 1

		if err := tx.Save(&post).Error; err != nil {
			tx.Rollback()
			log.Errorf("Database error: %s", err.Error())
			return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			log.Errorf("Database error: %s", err.Error())
			return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
		}

		return nil
	}

	return nil
}

func (s *PostServiceStruct) GetUserWhoLikedPost(postID uint) ([]dto.PreviewUser, error) {
	var post entity.Post
	var previewUsers []dto.PreviewUser

	if err := s.db.Model(&entity.Post{}).
		Preload("Likes.User").
		First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		log.Errorf("Database error: %s", err.Error())
		return nil, &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	for _, user := range post.Likes {
		previewUser := dto.PreviewUser{
			ID:             user.User.ID,
			Username:       user.User.Username,
			ProfilePicture: user.User.ProfilePicture,
		}

		previewUsers = append(previewUsers, previewUser)
	}

	return previewUsers, nil
}
