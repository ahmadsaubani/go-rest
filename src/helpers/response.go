package helpers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"gin/src/configs/database"
	"gin/src/utils/loggers"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
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
	Message interface{} `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Links   interface{} `json:"links,omitempty"`
}

type ErrorData struct {
	Error       interface{} `json:"error"`
	Path        string      `json:"path"`
	Method      string      `json:"method"`
	ClientIP    string      `json:"clientIP"`
	Status      int         `json:"status"`
	RequestBody interface{} `json:"requestBody"`
}

// SuccessResponse sends a JSON response with a success status, a message, and the provided data.
// If pagination is provided, it will include pagination metadata and links.
// The message will default to "Data Found" if not provided.
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

// GetPaginatedData retrieves data with pagination or returns an empty array if the page is too high
// GetPaginatedData fetches paginated data and returns it
func GetPaginatedData[T any](ctx *gin.Context, db *gorm.DB, order string, page, limit, offset int) ([]T, PaginationMeta, int64) {
	var data []T
	var total int64

	db.Model(new(T)).Count(&total)

	// If page is out of range, return empty data but still return meta
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if page > totalPages && totalPages != 0 {
		page = totalPages
	}

	// Use provided order or default to "created_at desc"
	if order == "" {
		order = "created_at desc"
	}

	db.Limit(limit).
		Offset(offset).
		Order(order).
		Find(&data)

	// db.Limit(limit).Offset(offset).Order("created_at desc").Find(&data)

	meta := PaginationMeta{
		Page:  page,
		Limit: limit,
		Total: total,
	}

	return data, meta, total
}

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

// ErrorResponse sends a JSON response with the given error and HTTP status code.
// If the HTTP status code is not provided, it defaults to 400 Bad Request.
func ErrorResponse(ctx *gin.Context, err error, httpCode ...int) {
	if len(httpCode) == 0 {
		httpCode = append(httpCode, http.StatusBadRequest)
	}

	// Baca ulang body jika ingin log data request
	var requestData interface{}
	if rawBody, exists := ctx.Get("RequestBody"); exists {
		if bodyBytes, ok := rawBody.([]byte); ok {
			var jsonBody map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &jsonBody); err == nil {
				requestData = jsonBody
			} else {
				requestData = string(bodyBytes) // fallback as string
			}
		}
	} else if formData, exists := ctx.Get("RequestForm"); exists {
		requestData = formData
	}

	// check apakah error validation atau error hardcode
	message := ParseValidationError(err)

	// Buat LogData dan panggil logError
	logData := ErrorData{
		Error:       message,
		Path:        ctx.FullPath(),
		Method:      ctx.Request.Method,
		ClientIP:    ctx.ClientIP(),
		Status:      httpCode[0],
		RequestBody: requestData,
	}

	// Logging dengan logError
	logError(logData)

	webResponse := Response{
		Success: false,
		Message: message,
		Data:    nil,
	}
	ctx.JSON(httpCode[0], webResponse)
}

func logError(data ErrorData) {
	// Jika message adalah string JSON, coba unmarshal
	if messageStr, ok := data.Error.(string); ok {
		var unmarshalledMessage map[string]string
		if err := json.Unmarshal([]byte(messageStr), &unmarshalledMessage); err == nil {
			// Validasi berhasil di-unmarshal
			loggers.Log.Error("Validation Error: ", map[string]interface{}{
				"error":       unmarshalledMessage,
				"path":        data.Path,
				"method":      data.Method,
				"clientIP":    data.ClientIP,
				"status":      data.Status,
				"requestBody": data.RequestBody,
			})
			return
		}
		// Gagal unmarshal JSON string
		loggers.Log.Error("Validation Error Unmarshal Failed: ", map[string]interface{}{
			"error":       messageStr,
			"path":        data.Path,
			"method":      data.Method,
			"clientIP":    data.ClientIP,
			"status":      data.Status,
			"requestBody": data.RequestBody,
		})
		return
	}

	// Bukan string, langsung log biasa
	loggers.Log.Error("ErrorResponse: ", map[string]interface{}{
		"error":       data.Error,
		"path":        data.Path,
		"method":      data.Method,
		"clientIP":    data.ClientIP,
		"status":      data.Status,
		"requestBody": data.RequestBody,
	})
}

func FormatLogError(message interface{}, ctx *gin.Context, httpCode int, requestData interface{}) {
	// Jika message adalah string JSON, coba unmarshal
	if messageStr, ok := message.(string); ok {
		var unmarshalledMessage map[string]string
		if err := json.Unmarshal([]byte(messageStr), &unmarshalledMessage); err == nil {
			// Log objek yang sudah di-unmarshal
			loggers.Log.Error("Validation Error: ", map[string]interface{}{
				"error":       unmarshalledMessage,
				"path":        ctx.FullPath(),
				"method":      ctx.Request.Method,
				"clientIP":    ctx.ClientIP(),
				"status":      httpCode,
				"requestBody": requestData,
			})
		} else {
			// Jika gagal unmarshal, log sebagai string biasa
			loggers.Log.Error("Validation Error Unmarshal Failed: ", map[string]interface{}{
				"error":       messageStr,
				"path":        ctx.FullPath(),
				"method":      ctx.Request.Method,
				"clientIP":    ctx.ClientIP(),
				"status":      httpCode,
				"requestBody": requestData,
			})
		}
	} else {
		// Jika bukan error validasi, log error biasa
		loggers.Log.Error("ErrorResponse: ", map[string]interface{}{
			"error":       message,
			"path":        ctx.FullPath(),
			"method":      ctx.Request.Method,
			"clientIP":    ctx.ClientIP(),
			"status":      httpCode,
			"requestBody": requestData,
		})
	}
}

func ParseValidationError(err error) interface{} {
	var ve validator.ValidationErrors
	var message interface{}

	if errors.As(err, &ve) {
		errorMap := map[string]string{}
		for _, fe := range ve {
			field := strings.ToLower(fe.Field())
			switch fe.Tag() {
			case "required":
				errorMap[field] = fmt.Sprintf("%s is required", field)
			case "email":
				errorMap[field] = fmt.Sprintf("%s must be a valid email", field)
			case "min":
				errorMap[field] = fmt.Sprintf("%s must be at least %s characters", field, fe.Param())
			default:
				errorMap[field] = fmt.Sprintf("%s is invalid", field)
			}
		}
		// Convert errorMap menjadi JSON string
		jsonMessage, _ := json.Marshal(errorMap)
		message = string(jsonMessage) // convert map ke string
		return message
	}
	return err.Error()
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

// GetPaginationParams parses the page and per_page query parameters from the given context.
// The page query parameter is required to be a positive integer, and the per_page query parameter
// is required to be a positive integer between 1 and 100. If the query parameters are not valid,
// the function will set the page to 1 and the per_page to 10. The function returns the parsed page,
// limit, and offset as integers.
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

func CountModel[T any]() (int64, error) {
	if database.GormDB != nil {
		var total int64
		err := database.GormDB.Model(new(T)).Count(&total).Error
		return total, err
	}

	if database.SQLDB == nil {
		return 0, sql.ErrConnDone
	}

	var model T
	table := GetTableName(&model)

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	row := database.SQLDB.QueryRow(query)

	var total int64
	err := row.Scan(&total)
	return total, err
}
