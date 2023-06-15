package service

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"social-app/app/dto"
	"social-app/app/model"
	customError "social-app/pkg/error"
	"social-app/pkg/util"
	"time"
)

type PostService interface {
	CreatePost(ctx *fiber.Ctx, userID uint, request dto.CreatePostDto) error
	ViewPost(userID uint, postID uint) error
	GetUserFeed(userID uint) ([]model.PostWithUser, error)
	GetViewersOnPost(postID uint) ([]model.PreviewUser, error)
	LikePost(UserID uint, PostID uint) error
	UnLikePost(UserID uint, PostID uint) error
	GetUserWhoLikedPost(postID uint) ([]model.PreviewUser, error)
	CreateCommentPost(postID uint, userID uint, request dto.CreateCommentPostDto) error
	GetPostComments(postID uint) ([]model.PreviewComment, error)
	DeletePostComment(commentID uint) error
}

type PostServiceImpl struct {
	db *gorm.DB
}

func NewPostService(db *gorm.DB) PostService {
	return &PostServiceImpl{db: db}
}

func (s *PostServiceImpl) CreatePost(ctx *fiber.Ctx, userID uint, request dto.CreatePostDto) error {
	_, filePath, err := util.UploadImage(ctx, "post", "image")

	if err != nil {
		log.Errorf("Upload file error: %s", err.Error())
		return &customError.ServiceError{Type: "upload_file_error", Errors: err.Error()}
	}

	post := model.Post{
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

func (s *PostServiceImpl) ViewPost(userID uint, postID uint) error {
	var user model.User

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

	var post model.Post
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
	if err := tx.Model(&model.UserViewedPost{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&countUserViewed).Error; err != nil {
		tx.Rollback()
		log.Errorf("Database error: %s", err.Error())
		return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	if countUserViewed == 0 {
		addUserViewedPost := model.UserViewedPost{
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

func (s *PostServiceImpl) GetUserFeed(userID uint) ([]model.PostWithUser, error) {
	var userViewedPost []model.UserViewedPost
	var userPost []model.Post
	var userFeedPost []model.PostWithUser

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

func (s *PostServiceImpl) GetViewersOnPost(postID uint) ([]model.PreviewUser, error) {
	var post model.Post
	var usersPreview []model.PreviewUser

	if err := s.db.Model(&model.Post{}).
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
			userPreview := model.PreviewUser{
				ID:             viewedBy.User.ID,
				Username:       viewedBy.User.Username,
				ProfilePicture: viewedBy.User.ProfilePicture,
			}

			usersPreview = append(usersPreview, userPreview)
		}
	}

	return usersPreview, nil
}

func (s *PostServiceImpl) LikePost(UserID uint, PostID uint) error {
	var user model.User
	var post model.Post

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

	createLikePost := model.UserLikedPost{
		UserID:  UserID,
		PostID:  PostID,
		LikedAt: time.Now(),
	}

	var countLikedPost int64
	if err := tx.Model(&model.UserLikedPost{}).
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

func (s *PostServiceImpl) UnLikePost(UserID uint, PostID uint) error {
	var user model.User
	var post model.Post

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
	if err := tx.Model(&model.UserLikedPost{}).
		Where("user_id = ? AND post_id = ?", UserID, PostID).
		Count(&countLikedPost).Error; err != nil {
		tx.Rollback()
		log.Errorf("Database error: %s", err.Error())
		return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	if countLikedPost != 0 {
		if err := tx.Where("user_id = ? AND post_id = ?", UserID, PostID).Delete(&model.UserLikedPost{}).Error; err != nil {
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

func (s *PostServiceImpl) GetUserWhoLikedPost(postID uint) ([]model.PreviewUser, error) {
	var post model.Post
	var previewUsers []model.PreviewUser

	if err := s.db.Model(&model.Post{}).
		Preload("Likes.User").
		First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		log.Errorf("Database error: %s", err.Error())
		return nil, &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	for _, user := range post.Likes {
		previewUser := model.PreviewUser{
			ID:             user.User.ID,
			Username:       user.User.Username,
			ProfilePicture: user.User.ProfilePicture,
		}

		previewUsers = append(previewUsers, previewUser)
	}

	return previewUsers, nil
}

func (s *PostServiceImpl) CreateCommentPost(postID uint, userID uint, request dto.CreateCommentPostDto) error {
	var post model.Post
	var user model.User

	tx := s.db.Begin()
	if tx.Error != nil {
		log.Errorf("Database error: %s", tx.Error.Error())
		return &customError.ServiceError{Type: "database_error", Errors: tx.Error.Error()}
	}

	if err := tx.First(&post, postID).Error; err != nil {
		log.Errorf("Database error: %s", err.Error())
		return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	if err := tx.First(&user, userID).Error; err != nil {
		log.Errorf("Database error: %s", err.Error())
		return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	commentPost := &model.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: request.Comment,
	}

	if err := tx.Create(&commentPost).Error; err != nil {
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

func (s *PostServiceImpl) GetPostComments(postID uint) ([]model.PreviewComment, error) {
	var post model.Post
	var comments []model.PreviewComment

	if err := s.db.Preload("Comments.User").First(&post, postID).Error; err != nil {
		log.Errorf("Database error: %s", err.Error())
		return nil, &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	for _, comment := range post.Comments {
		commentItem := model.PreviewComment{
			ID:        comment.ID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
			User: model.PreviewUser{
				ID:             comment.User.ID,
				Username:       comment.User.Username,
				ProfilePicture: comment.User.ProfilePicture,
			},
		}

		comments = append(comments, commentItem)
	}

	return comments, nil
}

func (s *PostServiceImpl) DeletePostComment(commentID uint) error {
	if err := s.db.Delete(&model.Comment{}, commentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &customError.ServiceError{Type: "value_doesnt_exist", Errors: "comment not found"}
		}

		log.Errorf("Database error: %s", err.Error())
		return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	return nil
}
