package helpers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaginationMeta struct {
	Page  int    `json:"page"`
	Limit int    `json:"per_page"`
	Total int64  `json:"total"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
	First string `json:"first,omitempty"`
	Last  string `json:"last,omitempty"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Links   interface{} `json:"links,omitempty"`
}

func SuccessResponse(ctx *gin.Context, message string, data interface{}, pagination ...PaginationMeta) {
	// Default message if not provided
	if message == "" {
		message = "Data Found"
	}

	// Create response object
	webResponse := Response{
		Success: true,
		Message: message,
		Data:    data,
	}

	// If pagination meta is provided, add meta and links
	if len(pagination) > 0 {
		webResponse.Meta = map[string]interface{}{
			"pagination": pagination[0],
		}
		webResponse.Links = buildPaginationLinks(ctx, pagination[0])
	}

	// Return the response
	JSONResponse(ctx, webResponse)
}

// PaginatedResponse handles paginated data and metadata
// PaginatedResponse handles paginated data and metadata
func PaginatedResponse(ctx *gin.Context, message string, data interface{}, page, limit int, total int64) {
	meta := PaginationMeta{
		Page:  page,
		Limit: limit,
		Total: total,
	}

	SuccessResponse(ctx, message, data, meta)
}

// buildPaginationLinks generates pagination links (next, prev, first, last)
func buildPaginationLinks(ctx *gin.Context, meta PaginationMeta) map[string]string {
	links := make(map[string]string)

	// Calculate total pages
	totalPages := int(math.Ceil(float64(meta.Total) / float64(meta.Limit)))

	if meta.Page < totalPages {
		links["next"] = buildPaginationLink(ctx, meta.Page+1, meta.Limit)
	}
	if meta.Page > 1 {
		links["prev"] = buildPaginationLink(ctx, meta.Page-1, meta.Limit)
	}
	if meta.Page > 1 {
		links["first"] = buildPaginationLink(ctx, 1, meta.Limit)
	}
	if meta.Page < totalPages {
		links["last"] = buildPaginationLink(ctx, totalPages, meta.Limit)
	}

	return links
}

// buildPaginationLink constructs a pagination URL for the given page and limit
func buildPaginationLink(ctx *gin.Context, page, limit int) string {
	return fmt.Sprintf("%s?page=%d&per_page=%d", ctx.Request.URL.Path, page, limit)
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

// JSONResponse handles sending the JSON response with the correct status code
func JSONResponse(ctx *gin.Context, data interface{}) {
	isCreate := ctx.Request.Method == http.MethodPost
	statusCode := http.StatusOK
	if isCreate {
		statusCode = http.StatusCreated
	}

	ctx.JSON(statusCode, data)
}

func GetPaginationParams(ctx *gin.Context) (page, limit, offset int) {
	page, _ = strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ = strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset = (page - 1) * limit
	return
}
