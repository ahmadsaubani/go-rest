package auth

import (
	"gin/src/entities/users"
	"gin/src/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx *gin.Context) {
	var user users.User

	// Bind JSON directly to struct
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Hash the user's password before saving it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		helpers.ErrorResponse(err, ctx, http.StatusBadRequest)
		return
	}

	// Update the password with the hashed password
	user.Password = string(hashedPassword)

	if err := helpers.InsertModel(&user); err != nil {
		helpers.ErrorResponse(err, ctx, http.StatusBadRequest)
		return
	}

	response := users.ResponseRegister{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}

	helpers.SuccessResponse(ctx, "Data Found", response)
}
