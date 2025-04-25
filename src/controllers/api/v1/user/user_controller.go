package user

import (
	"fmt"
	"gin/src/entities/users"
	"gin/src/helpers"
	services "gin/src/services/user_services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Find user by ID
	var user users.User
	if err := helpers.GetModelByID(&user, userID); err != nil {
		helpers.ErrorResponse(ctx, err, http.StatusNotFound)
		return
	}

	response := users.ProfileResponse{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}

	helpers.SuccessResponse(ctx, "Data Found!", response)
}

// GetAllUsers returns all users with pagination.
// It will return a JSON response with pagination metadata and links.
// The response will contain an array of users with their UUID, ID, Email, and Username.
func GetAllUsers(service services.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		page, limit, offset := helpers.GetPaginationParams(ctx)

		userList, total, err := service.GetPaginatedUsers(ctx, limit, offset)
		if err != nil {
			helpers.ErrorResponse(ctx, err, http.StatusInternalServerError)
			return
		}

		var response []users.ProfileResponse
		for _, u := range userList {
			response = append(response, users.ProfileResponse{
				UUID:     u.UUID,
				ID:       u.ID,
				Email:    u.Email,
				Username: u.Username,
			})
		}

		helpers.SuccessResponse(ctx, "Data found!", response, helpers.PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
		})
	}
}

func UploadAvatar(userService services.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("user_id")

		if !exists {
			helpers.ErrorResponse(ctx, fmt.Errorf("user not exists"), http.StatusBadRequest)
			return
		}

		var userIDInt64 int64
		switch v := userID.(type) {
		case int64:
			userIDInt64 = v
		case uint:
			userIDInt64 = int64(v)
		case uint64:
			userIDInt64 = int64(v)
		default:
			helpers.ErrorResponse(ctx, fmt.Errorf("unexpected type for user_id: %T", v), http.StatusBadRequest)
			return
		}

		// Bind file from the request
		file, _, err := ctx.Request.FormFile("file")

		if err != nil {
			helpers.ErrorResponse(ctx, fmt.Errorf("failed to get file from form-data: %w", err), http.StatusBadRequest)
			return
		}

		// Call the service to upload the file and get the URL
		avatarURL, err := userService.UploadAvatar(ctx, userIDInt64, file, "avatars")
		if err != nil {
			helpers.ErrorResponse(ctx, err, http.StatusInternalServerError)
			return
		}

		// Return success response
		helpers.SuccessResponse(ctx, "Avatar uploaded successfully", gin.H{"avatar_url": avatarURL})
	}
}
