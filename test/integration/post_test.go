package integration

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"social-app/app/model"
	"strconv"
	"testing"
)

var postID string

func TestCreatePost(t *testing.T) {
	token, err := getBearerToken(app)

	assert.Nil(t, err)

	t.Run("SuccessCreatePost", func(t *testing.T) {
		body, writer, err := getBodyRequestMultipart("default_post.jpg")
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/post", body)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Content-Type", writer.FormDataContentType())
		assert.Nil(t, err)

		res, err := app.Test(req)
		assert.Nil(t, err)

		assert.Equal(t, res.StatusCode, http.StatusOK)

		assert.Nil(t, err)

		response, err := GetResponseSuccess(res)
		assert.Nil(t, err)

		assert.Equal(t, "successfully added post", response.Message)
	})

	t.Run("ErrorCreatePostWhenEmptyRequest", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/api/v1/post", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Content-Type", "multipart/form-data")
		assert.Nil(t, err)

		res, err := app.Test(req)
		assert.Nil(t, err)

		assert.Equal(t, res.StatusCode, http.StatusBadRequest)

		assert.Nil(t, err)

		response, err := GetResponseError(res)
		assert.Nil(t, err)

		assert.Equal(t, "validation_error", response.Type)
	})

	var post model.Post

	err = getDatabaseConnection().Model(&model.Post{}).Where("caption = ?", "test").First(&post).Error
	assert.Nil(t, err)

	postID = strconv.FormatUint(uint64(post.ID), 10)
}
func TestGetUserFeed(t *testing.T) {
	token, err := getBearerToken(app)
	assert.Nil(t, err)

	t.Run("SuccessGetUserFeed", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/v1/post", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		assert.Nil(t, err)

		res, err := app.Test(req)
		assert.Nil(t, err)

		response, err := GetResponseSuccess(res)
		assert.Nil(t, err)

		assert.Equal(t, "successfully to get feed post", response.Message)
		assert.NotNil(t, response.Data)
	})
}

func TestAddUserViewedPost(t *testing.T) {
	token, err := getBearerToken(app)

	assert.Nil(t, err)

	t.Run("SuccessAddUserViewedPost", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/post/%s/view", postID), nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		assert.Nil(t, err)

		res, err := app.Test(req)
		assert.Nil(t, err)

		response, err := GetResponseSuccess(res)
		assert.Nil(t, err)

		assert.Equal(t, fmt.Sprintf("successfully to see post with id %s", postID), response.Message)
	})
}
func TestGetViewersOnPost(t *testing.T) {
	token, err := getBearerToken(app)
	assert.Nil(t, err)

	t.Run("SuccessGetUsersOnPost", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/post/%s/viewers", postID), nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		assert.Nil(t, err)

		res, err := app.Test(req)
		assert.Nil(t, err)

		response, err := GetResponseSuccess(res)
		assert.Nil(t, err)

		assert.Equal(t, fmt.Sprintf("successfully to get viewers from post id %s", postID), response.Message)
	})
}

func TestLikePost(t *testing.T) {
	token, err := getBearerToken(app)
	assert.Nil(t, err)

	t.Run("SuccessLikedPost", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/post/%s/like", postID), nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		assert.Nil(t, err)

		res, err := app.Test(req)
		assert.Nil(t, err)

		response, err := GetResponseSuccess(res)
		assert.Nil(t, err)

		assert.Equal(t, fmt.Sprintf("successfully to like a post with id %s", postID), response.Message)
	})
}
func TestGetUserLikedPost(t *testing.T) {
	token, err := getBearerToken(app)
	assert.Nil(t, err)

	t.Run("SuccessGetUserLikedPost", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/post/%s/likes", postID), nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		assert.Nil(t, err)

		res, err := app.Test(req)
		assert.Nil(t, err)

		response, err := GetResponseSuccess(res)
		assert.Nil(t, err)

		assert.Equal(t, fmt.Sprintf("managed to fetch user data who liked posts with id %s", postID), response.Message)
	})
}
func TestUnlikePost(t *testing.T) {
	token, err := getBearerToken(app)
	assert.Nil(t, err)

	t.Run("SuccessUnlikedPost", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/post/%s/unlike", postID), nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		assert.Nil(t, err)

		res, err := app.Test(req)
		assert.Nil(t, err)

		response, err := GetResponseSuccess(res)
		assert.Nil(t, err)

		assert.Equal(t, fmt.Sprintf("successfully to unlike a post with id %s", postID), response.Message)
	})
}

func TestCreatePostComment(t *testing.T) {

}
func TestGetPostComments(t *testing.T) {

}
