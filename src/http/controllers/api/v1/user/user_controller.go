package user

import (
	"fmt"
	"gin/src/configs/database"
	"gin/src/entities/users"
	"gin/src/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProfile(ctx *gin.Context) {
	// Extract user_id from context set by JWT middleware
	userID, exists := ctx.Get("user_id")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Find user by ID
	var user users.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		fmt.Println(err)
		helpers.ErrorResponse(err, ctx, http.StatusNotFound)
		return
	}

	response := users.ProfileResponse{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}

	helpers.SuccessResponse(ctx, "Data Found!", response)
}

// GetAllUsers handles the request to retrieve a paginated list of users.
// It extracts pagination parameters from the query string, retrieves the total
// count of users, and fetches a paginated list of users from the database.
// The function returns a paginated response containing the users list along
// with pagination metadata.

func GetAllUsers(ctx *gin.Context) {
	page, limit, offset := helpers.GetPaginationParams(ctx)

	// Use the generic helper to get paginated users
	usersList, meta, _ := helpers.GetPaginatedData[users.User](ctx, database.DB, "created_at desc", page, limit, offset)

	// Convert to response struct
	var profileResponses []users.ProfileResponse
	for _, user := range usersList {
		profileResponses = append(profileResponses, users.ProfileResponse{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		})
	}

	if len(profileResponses) == 0 {
		profileResponses = []users.ProfileResponse{}
	}

	helpers.SuccessResponse(ctx, "Data Found!", profileResponses, meta)
}
