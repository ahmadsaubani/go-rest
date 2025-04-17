package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func SuccessResponse(ctx *gin.Context, data interface{}) {
	var webResponse = Response{}

	webResponse = Response{
		Success: true,
		Message: "Success",
		Data:    data,
	}
	JSONResponse(ctx, webResponse)
	return
}

func ErrorResponse(err error, ctx *gin.Context, httpCode ...int) {
	if len(httpCode) == 0 {
		httpCode = append(httpCode, http.StatusBadRequest)
	}

	webResponse := Response{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	}
	ctx.JSON(httpCode[0], webResponse)
}

func JSONResponse(ctx *gin.Context, data interface{}) {
	isCreate := ctx.Request.Method == http.MethodPost
	statusCode := http.StatusOK
	if isCreate {
		statusCode = http.StatusCreated
	}

	ctx.JSON(statusCode, data)
}
