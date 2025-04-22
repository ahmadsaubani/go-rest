package user

import (
	"fmt"
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
	if err := helpers.GetModelByID(&user, userID); err != nil {
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

// func GetAllUsers(ctx *gin.Context) {
// 	// Extract pagination parameters from query string
// 	page, limit, offset := helpers.GetPaginationParams(ctx)

// 	// var usersList []users.User
// 	var total int64

// 	// Get total count of users
// 	database.DB.Model(&users.User{}).Count(&total)

// 	// Retrieve users with pagination and ordering
// 	var usersList []users.User
// 	database.DB.
// 		Limit(limit).
// 		Offset(offset).
// 		Order("created_at desc").
// 		Find(&usersList)

// 	// Convert to response struct
// 	var profileResponses []users.ProfileResponse
// 	for _, user := range usersList {
// 		profileResponses = append(profileResponses, users.ProfileResponse{
// 			ID:       user.ID,
// 			Email:    user.Email,
// 			Username: user.Username,
// 		})
// 	}

// 	// Return consistent structure (empty array if nothing found)
// 	if len(profileResponses) == 0 {
// 		profileResponses = []users.ProfileResponse{}
// 	}

// 	// Send success response with pagination
// 	helpers.SuccessResponse(ctx, "Data Found!", profileResponses, helpers.PaginationMeta{
// 		Page:  page,
// 		Limit: limit,
// 		Total: total,
// 	})
// }

func GetAllUsers(ctx *gin.Context) {
	page, limit, offset := helpers.GetPaginationParams(ctx)
	var usersList []users.User

	total, err := helpers.CountModel[users.User]()
	if err != nil {
		helpers.ErrorResponse(err, ctx, http.StatusInternalServerError)
		return
	}

	err = helpers.GetAllModels(&usersList, limit, offset, "created_at DESC")
	if err != nil {
		helpers.ErrorResponse(err, ctx, http.StatusInternalServerError)
		return
	}

	var response []users.ProfileResponse
	for _, u := range usersList {
		response = append(response, users.ProfileResponse{
			ID:       u.ID,
			Email:    u.Email,
			Username: u.Username,
		})
	}

	helpers.SuccessResponse(ctx, "Data Found!", response, helpers.PaginationMeta{
		Page:  page,
		Limit: limit,
		Total: total,
	})
}
