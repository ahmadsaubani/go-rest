package user

import (
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
