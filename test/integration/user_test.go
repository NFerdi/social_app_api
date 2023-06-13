package integration

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"net/http"
	dto2 "social-app/internal/dto"
	"social-app/internal/entity"
	"testing"
)

func getFollowerID() (uint, error) {
	signupDto := dto2.SignupDto{
		Username: "test1",
		FullName: "testing",
		Email:    "test1@test.com",
		Password: "test",
	}
	requestBodySignup, err := json.Marshal(signupDto)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewBuffer(requestBodySignup))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	_, err = app.Test(req)
	if err != nil {
		return 0, err
	}

	var user2 entity.User

	if err := getDatabaseConnection().Model(&entity.User{}).Where("username = ?", signupDto.Username).First(&user2).Error; err != nil {
		return 0, err
	}

	return user2.ID, nil
}

func TestFindUserByParams(t *testing.T) {
	token, err := getBearerToken(app)

	assert.Nil(t, err)

	t.Run("SuccessFindUserByParams", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/v1/user/test", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		assert.Nil(t, err)

		res, err := app.Test(req)
		assert.Nil(t, err)

		assert.Equal(t, res.StatusCode, http.StatusOK)

		response, err := GetResponseSuccess(res)
		assert.Nil(t, err)

		assert.Equal(t, "managed to fetch user based on username", response.Message)
		assert.NotNil(t, response.Data)
	})

	t.Run("ErrorFindUserByParamsWhenUsernameDoesn'tExists", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/v1/user/tests", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		assert.Nil(t, err)

		res, err := app.Test(req)
		assert.Nil(t, err)

		assert.Equal(t, res.StatusCode, http.StatusBadRequest)

		response, err := GetResponseError(res)
		assert.Nil(t, err)

		assert.Equal(t, "account with username tests not found", response.Errors)
	})
}

func TestUpdateUserProfile(t *testing.T) {
	token, err := getBearerToken(app)
	assert.Nil(t, err)

	t.Run("SuccessUpdateUserProfile", func(t *testing.T) {
		body, writer, err := getBodyRequestMultipart("default_avatar.png")
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPatch, "/api/v1/user", body)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Content-Type", writer.FormDataContentType())
		assert.Nil(t, err)

		res, err := app.Test(req)
		assert.Nil(t, err)

		assert.Equal(t, res.StatusCode, http.StatusOK)

		assert.Nil(t, err)

		response, err := GetResponseSuccess(res)
		assert.Nil(t, err)

		assert.Equal(t, "managed to change the user profile", response.Message)
	})
}

func TestFollowUser2(t *testing.T) {
	token, err := getBearerToken(app)

	assert.Nil(t, err)

	userFollowerId, err := getFollowerID()

	t.Run("SuccessFollowUser2", func(t *testing.T) {
		followUserDto := dto2.FollowUserDto{
			FollowerId: userFollowerId,
		}

		requestBody, err := json.Marshal(followUserDto)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/user/follow", bytes.NewBuffer(requestBody))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Content-Type", "application/json")
		assert.Nil(t, err)

		res, err := app.Test(req)
		assert.Nil(t, err)

		assert.Equal(t, res.StatusCode, http.StatusOK)

		response, err := GetResponseSuccess(res)
		assert.Nil(t, err)

		assert.Equal(t, "managed to follow user", response.Message)
	})

	t.Run("ErrorFollowUser2WhenAlreadyFollowed", func(t *testing.T) {
		followUserDto := dto2.FollowUserDto{
			FollowerId: userFollowerId,
		}

		requestBody, err := json.Marshal(followUserDto)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/user/follow", bytes.NewBuffer(requestBody))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Content-Type", "application/json")
		assert.Nil(t, err)

		res, err := app.Test(req)
		assert.Nil(t, err)

		assert.Equal(t, res.StatusCode, http.StatusBadRequest)

		response, err := GetResponseError(res)
		assert.Nil(t, err)

		assert.Equal(t, "You have done following", response.Errors)
	})
}

func TestUnfollowUser2(t *testing.T) {
	token, err := getBearerToken(app)

	assert.Nil(t, err)

	userFollowerId, err := getFollowerID()

	t.Run("SuccessUnfollowUser2", func(t *testing.T) {
		followUserDto := dto2.FollowUserDto{
			FollowerId: userFollowerId,
		}

		requestBody, err := json.Marshal(followUserDto)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/user/unfollow", bytes.NewBuffer(requestBody))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Content-Type", "application/json")
		assert.Nil(t, err)

		res, err := app.Test(req)
		assert.Nil(t, err)

		assert.Equal(t, res.StatusCode, http.StatusOK)

		response, err := GetResponseSuccess(res)
		assert.Nil(t, err)

		assert.Equal(t, "successfully unfollowed user", response.Message)
	})

	t.Run("ErrorUnfollowUser2WhenHaven'tFollow", func(t *testing.T) {
		followUserDto := dto2.FollowUserDto{
			FollowerId: userFollowerId,
		}

		requestBody, err := json.Marshal(followUserDto)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/user/unfollow", bytes.NewBuffer(requestBody))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Content-Type", "application/json")
		assert.Nil(t, err)

		res, err := app.Test(req)
		assert.Nil(t, err)

		assert.Equal(t, res.StatusCode, http.StatusBadRequest)

		response, err := GetResponseError(res)
		assert.Nil(t, err)

		assert.Equal(t, "You haven't followed yet", response.Errors)
	})
}

func TestGetUserFollower(t *testing.T) {
	token, err := getBearerToken(app)

	assert.Nil(t, err)

	t.Run("SuccessGetUserFollowers", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/v1/user/test/followers", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		assert.Nil(t, err)

		res, err := app.Test(req)
		assert.Nil(t, err)

		assert.Equal(t, res.StatusCode, http.StatusOK)

		response, err := GetResponseSuccess(res)
		assert.Nil(t, err)

		assert.Equal(t, "managed to fetch user followers", response.Message)
	})
}

func TestGetUserFollowing(t *testing.T) {
	token, err := getBearerToken(app)

	assert.Nil(t, err)

	t.Run("SuccessGetUserFollowing", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/v1/user/test/following", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		assert.Nil(t, err)

		res, err := app.Test(req)
		assert.Nil(t, err)

		assert.Equal(t, res.StatusCode, http.StatusOK)

		response, err := GetResponseSuccess(res)
		assert.Nil(t, err)

		assert.Equal(t, "managed to fetch user following", response.Message)
	})
}

func TestGetPostsWatchedByUser(t *testing.T) {
	token, err := getBearerToken(app)

	assert.Nil(t, err)

	t.Run("SuccessGetPostsWatchedByUser", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/v1/user/test/post/viewed", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		assert.Nil(t, err)

		res, err := app.Test(req)
		assert.Nil(t, err)

		assert.Equal(t, res.StatusCode, http.StatusOK)

		response, err := GetResponseSuccess(res)
		assert.Nil(t, err)

		assert.Equal(t, "managed to get posts that have been watched by users", response.Message)
	})
}

func TestGetUserLikedPosts(t *testing.T) {
	token, err := getBearerToken(app)

	assert.Nil(t, err)

	t.Run("SuccessGetUserLikedPosts", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/v1/user/test/post/liked", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		assert.Nil(t, err)

		res, err := app.Test(req)
		assert.Nil(t, err)

		assert.Equal(t, res.StatusCode, http.StatusOK)

		response, err := GetResponseSuccess(res)
		assert.Nil(t, err)

		assert.Equal(t, "managed to get posts that users like", response.Message)
	})
}

func TestGetUserUploadedPosts(t *testing.T) {
	token, err := getBearerToken(app)

	assert.Nil(t, err)

	t.Run("SuccessTestGetUserUploadedPosts", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/v1/user/test/post/uploaded", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		assert.Nil(t, err)

		res, err := app.Test(req)
		assert.Nil(t, err)

		assert.Equal(t, res.StatusCode, http.StatusOK)

		response, err := GetResponseSuccess(res)
		assert.Nil(t, err)

		assert.Equal(t, "managed to fetch all posts uploaded by users", response.Message)
	})
}
