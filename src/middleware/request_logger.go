package middleware

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func SaveRequestBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				c.Set("RequestBody", bodyBytes)
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // Set lagi agar bisa dibaca di ShouldBindJSON
			}
		}
		c.Next()
	}
}
